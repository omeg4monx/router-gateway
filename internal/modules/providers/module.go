// Package providers implements the pluggable ProviderAdapter registry,
// vault-encrypted upstream provider connections, custom compatible provider
// nodes, connectivity quick-testing, and headless-safe PKCE OAuth flows.
package providers

import (
	"context"
	"crypto/rand"
	"embed"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"io/fs"
	"net/http"
	"sort"
	"strings"
	"sync"
	"time"

	"github.com/omeg4monx/router-gateway/internal/contracts"
	"github.com/omeg4monx/router-gateway/internal/platform/crypto"
	"github.com/omeg4monx/router-gateway/internal/platform/db"
	"github.com/omeg4monx/router-gateway/internal/platform/eventbus"
	"github.com/omeg4monx/router-gateway/internal/platform/httpserver"
	"github.com/omeg4monx/router-gateway/internal/platform/lifecycle"
)

const (
	moduleName = "providers"

	// TopicProvidersChanged is published on the event bus whenever a provider
	// connection or custom provider node is created, updated, refreshed, or deleted.
	TopicProvidersChanged = "providers.changed"

	defaultAdminSubject  = "admin"
	defaultHTTPTimeout   = 30 * time.Second
	idRandomBytes        = 12
	defaultPriorityValue = 100
)

//go:embed migrations/*.sql
var migrationsFS embed.FS

// ChangedEvent is published on TopicProvidersChanged after committing a
// provider connection or node mutation. It never contains secret material.
type ChangedEvent struct {
	Action       string `json:"action"`
	ResourceType string `json:"resource_type"`
	ResourceID   string `json:"resource_id"`
	ProviderID   string `json:"provider_id,omitempty"`
	ChangedBy    string `json:"changed_by"`
	UpdatedAt    string `json:"updated_at"`
}

// Options configures a new providers Module instance.
type Options struct {
	DB             *db.DB
	Vault          *crypto.Vault
	HTTPClient     *http.Client
	Bus            *eventbus.Bus
	OAuthOverrides map[string]OAuthConfig
}

// Module implements lifecycle.Module and contracts.ProviderStore.
type Module struct {
	db         *db.DB
	vault      *crypto.Vault
	httpClient *http.Client
	bus        *eventbus.Bus
	repo       *repository

	mu              sync.RWMutex
	adapters        map[string]ProviderAdapter
	oauthOverrides  map[string]OAuthConfig
	oauthStrategies map[string]*oauthProviderStrategy

	refreshMu        sync.Mutex
	refreshCalls     map[string]*refreshCall
	refreshCooldowns map[string]refreshCooldownEntry

	awsQuick    *awsQuickRuntime
	deviceOAuth *deviceOAuthRuntime
	oauthProxy  *oauthLoopbackProxy
}

var (
	_ lifecycle.Module        = (*Module)(nil)
	_ contracts.ProviderStore = (*Module)(nil)
)

// New constructs a providers Module pre-populated with built-in API-key and
// OAuth adapters.
func New(opts Options) *Module {
	client := opts.HTTPClient
	if client == nil {
		client = &http.Client{Timeout: defaultHTTPTimeout}
	}

	overrides := make(map[string]OAuthConfig, len(opts.OAuthOverrides))
	for k, v := range opts.OAuthOverrides {
		cloned := cloneOAuthConfig(&v)
		overrides[k] = *cloned
	}

	m := &Module{
		db:             opts.DB,
		vault:          opts.Vault,
		httpClient:     client,
		bus:            opts.Bus,
		repo:           newRepository(opts.DB),
		adapters:         make(map[string]ProviderAdapter),
		oauthOverrides:   overrides,
		oauthStrategies:  defaultOAuthStrategies(),
		refreshCalls:     make(map[string]*refreshCall),
		refreshCooldowns: make(map[string]refreshCooldownEntry),
		awsQuick:         newAwsQuickRuntime(),
		deviceOAuth:      newDeviceOAuthRuntime(),
		oauthProxy:       newOAuthLoopbackProxy(),
	}

	for _, adapter := range allBuiltInAdapters() {
		m.RegisterAdapter(adapter)
	}
	return m
}

// Name returns the canonical module name.
func (m *Module) Name() string {
	return moduleName
}

// Dependencies declares that providers depends on the system module.
func (m *Module) Dependencies() []string {
	return []string{"system"}
}

// Migrations returns the embedded SQL migration filesystem rooted at migrations/.
func (m *Module) Migrations() fs.FS {
	sub, err := fs.Sub(migrationsFS, "migrations")
	if err != nil {
		panic(fmt.Sprintf("providers: load embedded migrations: %v", err))
	}
	return sub
}

// RegisterRoutes mounts administrative provider routes onto AdminAPI and the
// browser OAuth callback handoff onto PublicAPI.
func (m *Module) RegisterRoutes(mounts lifecycle.RouteMounts) {
	if mounts.AdminAPI != nil {
		mounts.AdminAPI.Get("/providers/catalog", m.handleGetCatalog)

		mounts.AdminAPI.Get("/providers/nodes", m.handleListNodes)
		mounts.AdminAPI.Post("/providers/nodes", m.handleCreateNode)
		mounts.AdminAPI.Post("/providers/node-drafts/test", m.handleTestNodeDraft)
		mounts.AdminAPI.Delete("/providers/nodes/{id}", m.handleDeleteNode)

		mounts.AdminAPI.Get("/providers/connections", m.handleListConnections)
		mounts.AdminAPI.Post("/providers/connections", m.handleCreateConnection)
		mounts.AdminAPI.Get("/providers/connections/{id}", m.handleGetConnection)
		mounts.AdminAPI.Patch("/providers/connections/{id}", m.handlePatchConnection)
		mounts.AdminAPI.Delete("/providers/connections/{id}", m.handleDeleteConnection)
		mounts.AdminAPI.Post("/providers/connections/{id}/test", m.handleTestConnection)
		mounts.AdminAPI.Get("/providers/connections/{id}/models", m.handleFetchConnectionModels)

		mounts.AdminAPI.Post("/providers/oauth/{provider}/start", m.handleOAuthStart)
		mounts.AdminAPI.Post("/providers/oauth/{provider}/exchange", m.handleOAuthExchange)
		mounts.AdminAPI.Get("/providers/oauth/{provider}/status", m.handleOAuthStatus)
		mounts.AdminAPI.Post("/providers/oauth/{provider}/paste", m.handleOAuthPaste)
		mounts.AdminAPI.Post("/providers/oauth/{provider}/stop", m.handleOAuthStop)
	}

	if mounts.PublicAPI != nil {
		mounts.PublicAPI.Get("/providers/oauth/callback", m.handleOAuthCallback)
	}
}

// Start validates that required dependencies are wired.
func (m *Module) Start(context.Context) error {
	if m.db == nil {
		return errors.New("providers: database is required")
	}
	if m.vault == nil {
		return errors.New("providers: vault is required")
	}
	return nil
}

// Stop performs graceful module teardown.
func (m *Module) Stop(context.Context) error {
	m.shutdownAwsQuickProxy()
	m.oauthProxy.stopAll()
	return nil
}

// RegisterAdapter registers or replaces a pluggable ProviderAdapter by its ID.
func (m *Module) RegisterAdapter(adapter ProviderAdapter) {
	if adapter == nil {
		return
	}
	id := strings.TrimSpace(adapter.ID())
	if id == "" {
		return
	}
	m.mu.Lock()
	m.adapters[id] = adapter
	m.mu.Unlock()
}

// GetAdapter retrieves a registered ProviderAdapter by ID.
func (m *Module) GetAdapter(id string) (ProviderAdapter, bool) {
	m.mu.RLock()
	adapter, ok := m.adapters[id]
	m.mu.RUnlock()
	return adapter, ok
}

// ListAdapters returns all registered ProviderAdapters sorted deterministically by ID.
func (m *Module) ListAdapters() []ProviderAdapter {
	m.mu.RLock()
	out := make([]ProviderAdapter, 0, len(m.adapters))
	for _, adapter := range m.adapters {
		out = append(out, adapter)
	}
	m.mu.RUnlock()

	sort.Slice(out, func(i, j int) bool {
		return out[i].ID() < out[j].ID()
	})
	return out
}

// ListEnabledConnections returns all enabled upstream provider connections
// ordered by ascending priority without exposing secret credentials.
func (m *Module) ListEnabledConnections(ctx context.Context) ([]contracts.ProviderConnection, error) {
	return m.repo.ListEnabledConnections(ctx)
}

// ResolveCredentials decrypts and returns upstream credentials for connectionID.
func (m *Module) ResolveCredentials(ctx context.Context, connectionID string) (contracts.UpstreamCredentials, error) {
	rec, err := m.repo.GetConnectionRecord(ctx, connectionID)
	if err != nil {
		return contracts.UpstreamCredentials{}, err
	}

	stored, err := m.unsealCredentials(rec.CredentialsEnc)
	if err != nil {
		return contracts.UpstreamCredentials{}, err
	}

	return contracts.UpstreamCredentials{
		ConnectionID: rec.Summary.ID,
		ProviderID:   rec.Summary.ProviderID,
		Protocol:     rec.Summary.Protocol,
		BaseURL:      rec.Summary.BaseURL,
		AuthType:     rec.Summary.AuthType,
		APIKey:       stored.APIKey,
		AccessToken:  stored.AccessToken,
		ExtraHeaders: cloneStringMap(rec.Summary.Headers),
	}, nil
}

func (m *Module) sealCredentials(creds storedCredentials) (string, error) {
	if m.vault == nil {
		return "", errors.New("providers: vault is not initialized")
	}
	raw, err := json.Marshal(creds) //nolint:gosec // G117: serialized solely to immediately encrypt with AES-256-GCM Vault
	if err != nil {
		return "", fmt.Errorf("providers: marshal credentials: %w", err)
	}
	enc, err := m.vault.Encrypt(raw)
	if err != nil {
		return "", fmt.Errorf("providers: encrypt credentials: %w", err)
	}
	return enc, nil
}

func (m *Module) unsealCredentials(enc string) (storedCredentials, error) {
	if m.vault == nil {
		return storedCredentials{}, errors.New("providers: vault is not initialized")
	}
	raw, err := m.vault.Decrypt(enc)
	if err != nil {
		return storedCredentials{}, fmt.Errorf("providers: decrypt credentials: %w", err)
	}
	var creds storedCredentials
	if err := json.Unmarshal(raw, &creds); err != nil {
		return storedCredentials{}, fmt.Errorf("providers: unmarshal credentials: %w", err)
	}
	return creds, nil
}

func (m *Module) publishChanged(action, resourceType, resourceID, providerID, changedBy, updatedAt string) {
	if m.bus == nil {
		return
	}
	m.bus.Publish(TopicProvidersChanged, ChangedEvent{
		Action:       action,
		ResourceType: resourceType,
		ResourceID:   resourceID,
		ProviderID:   providerID,
		ChangedBy:    changedBy,
		UpdatedAt:    updatedAt,
	})
}

func actorFromContext(ctx context.Context) string {
	if subj := strings.TrimSpace(httpserver.AdminSubject(ctx)); subj != "" {
		return subj
	}
	return defaultAdminSubject
}

func generateID(prefix string) (string, error) {
	buf := make([]byte, idRandomBytes)
	if _, err := rand.Read(buf); err != nil {
		return "", fmt.Errorf("providers: generate id: %w", err)
	}
	return prefix + hex.EncodeToString(buf), nil
}
