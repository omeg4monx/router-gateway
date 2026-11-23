// Package contracts defines narrow, domain-agnostic interfaces and lightweight
// value types for cross-module communication so feature modules never import
// each other directly.
package contracts

import (
	"context"
	"errors"
	"net/http"
	"net/url"
	"time"
)

// ErrNotInitialized indicates that the system administrator password has not
// been configured yet.
var ErrNotInitialized = errors.New("contracts: system is not initialized")

// ErrModelNotFound indicates that a requested model, alias, or combo could not
// be resolved to any configured upstream target.
var ErrModelNotFound = errors.New("contracts: model not found")

// ErrModelForbidden indicates that the caller's API key is not permitted to
// access the requested model.
var ErrModelForbidden = errors.New("contracts: model forbidden")

// ErrAllAccountsCoolingDown indicates that all matching upstream provider
// connections are currently in cooldown.
var ErrAllAccountsCoolingDown = errors.New("contracts: all accounts cooling down")

// ErrConnectionNotFound indicates that the requested provider connection ID
// does not exist.
var ErrConnectionNotFound = errors.New("contracts: connection not found")

// ErrOAuthRefreshUnsupported indicates that the requested provider connection
// does not support OAuth token refresh.
var ErrOAuthRefreshUnsupported = errors.New("contracts: oauth refresh unsupported")

// ErrJobNotFound indicates that no registered background job matches the name.
var ErrJobNotFound = errors.New("contracts: job not found")

// ErrJobAlreadyRunning indicates that the requested background job is currently executing.
var ErrJobAlreadyRunning = errors.New("contracts: job already running")

// AdminCredentialStore exposes system-owned administrator credential and login
// enforcement state to the authentication module.
type AdminCredentialStore interface {
	AdminPasswordHash(ctx context.Context) (string, error)
	UpdateAdminPasswordHash(ctx context.Context, hash string) error
	RequireLogin(ctx context.Context) (bool, error)
	SetRequireLogin(ctx context.Context, required bool) error
}

// AuthVerifier validates administrative dashboard requests and reports whether
// dashboard authentication is currently enforced.
type AuthVerifier interface {
	AuthenticateAdmin(r *http.Request) (string, error)
	RequireLogin(ctx context.Context) bool
}

// KeyInfo describes an authenticated gateway API key and its access limits.
type KeyInfo struct {
	ID            string
	Name          string
	AllowedModels []string
	Scopes        []string
	Tags          []string
	RateLimitRPM  int
	QuotaUSD      float64
	Disabled      bool
}

// KeyValidator authenticates raw gateway API keys and returns their metadata.
type KeyValidator interface {
	ValidateKey(ctx context.Context, rawKey string) (KeyInfo, error)
}

// ResolvedTarget identifies a concrete upstream provider connection candidate
// returned by ProviderResolver in fallback priority order.
type ResolvedTarget struct {
	ProviderID    string
	ProviderType  string
	ConnectionID  string
	Model         string
	UpstreamModel string
	Endpoint      string
	Priority      int
	Metadata      map[string]string
}

// ProviderResolver resolves a requested model identifier (combo, alias,
// provider/model, or bare model) into an ordered list of upstream targets
// permitted by the caller's allowed model rules.
type ProviderResolver interface {
	Resolve(ctx context.Context, model string, allowed []string) ([]ResolvedTarget, error)
}

// CanonicalRequest carries normalized inbound gateway request state through
// registered GatewayInterceptor hooks.
type CanonicalRequest struct {
	RequestID string
	Protocol  string
	Model     string
	Stream    bool
	KeyID     string
	Headers   http.Header
	Body      []byte
	Metadata  map[string]any
}

// ResponseMeta carries upstream response header metadata through registered
// GatewayInterceptor hooks before the response body is streamed to the client.
type ResponseMeta struct {
	RequestID    string
	StatusCode   int
	Headers      http.Header
	ProviderID   string
	ConnectionID string
	Model        string
	Metadata     map[string]any
}

// StreamEvent represents a single streaming chunk passed through registered
// GatewayInterceptor hooks during SSE forwarding.
type StreamEvent struct {
	RequestID string
	Event     string
	Data      []byte
	Done      bool
	Metadata  map[string]any
}

// GatewayInterceptor allows extension modules to inspect or mutate gateway
// requests, response headers, and stream chunks in deterministic priority order.
type GatewayInterceptor interface {
	Name() string
	Priority() int
	OnRequest(ctx context.Context, req *CanonicalRequest) error
	OnResponseHeaders(ctx context.Context, meta *ResponseMeta) error
	OnStreamChunk(ctx context.Context, ev *StreamEvent) error
}

// UsageRecord captures per-request telemetry emitted asynchronously by the
// gateway to the usage module.
type UsageRecord struct {
	RequestID    string
	KeyID        string
	Model        string
	ProviderID   string
	ConnectionID string
	Protocol     string
	InputTokens  int64
	OutputTokens int64
	CachedTokens int64
	CostUSD      float64
	Latency      time.Duration
	TTFT         time.Duration
	StatusCode   int
	ErrorSummary string
	Stream       bool
	Timestamp    time.Time
}

// UsageRecorder records gateway request usage asynchronously without blocking
// response streaming.
type UsageRecorder interface {
	Record(ctx context.Context, rec UsageRecord)
}

// ProxyDialerResolver resolves an optional outbound proxy URL for a given
// upstream provider. A nil *url.URL with a nil error indicates direct dial.
type ProxyDialerResolver interface {
	ProxyURLForProvider(ctx context.Context, providerID string) (*url.URL, error)
}

// ProviderConnection describes an enabled upstream provider connection and its
// routing metadata without exposing secret credentials.
type ProviderConnection struct {
	ID           string
	ProviderID   string
	Name         string
	Protocol     string
	BaseURL      string
	AuthType     string
	Priority     int
	Enabled      bool
	Models       []string
	ExtraHeaders map[string]string
}

// UpstreamCredentials holds resolved upstream connection parameters and secret
// authentication material for executing an upstream request.
type UpstreamCredentials struct {
	ConnectionID string
	ProviderID   string
	Protocol     string
	BaseURL      string
	AuthType     string
	APIKey       string
	AccessToken  string
	ExtraHeaders map[string]string
}

// ProviderStore exposes enabled upstream provider connections and credential
// resolution/refresh operations to routing and gateway modules.
type ProviderStore interface {
	ListEnabledConnections(ctx context.Context) ([]ProviderConnection, error)
	ResolveCredentials(ctx context.Context, connectionID string) (UpstreamCredentials, error)
	RefreshOAuthToken(ctx context.Context, connectionID string) (UpstreamCredentials, error)
}

// AccountHealthTracker tracks transient failure cooldown state per upstream
// provider connection.
type AccountHealthTracker interface {
	IsInCooldown(connectionID string, now time.Time) bool
	RecordFailure(connectionID string, statusCode int, retryAfter time.Duration)
	RecordSuccess(connectionID string)
}

// ModelDescriptor describes an available model, combo, or alias exposed by the
// gateway model listing endpoints. Kind is one of "combo", "alias", or "model".
type ModelDescriptor struct {
	ID      string
	OwnedBy string
	Kind    string
}

// ModelLister lists available models, combos, and aliases permitted by the
// caller's allowed model rules.
type ModelLister interface {
	ListModels(ctx context.Context, allowed []string) ([]ModelDescriptor, error)
}

// ProtocolTranslator exposes stateless wire-protocol request and response
// translation across supported gateway protocols without cross-module coupling.
type ProtocolTranslator interface {
	SupportedProtocols() []string
	TranslateRequest(fromProto, toProto string, body []byte) ([]byte, error)
	TranslateResponse(fromProto, toProto string, body []byte) ([]byte, error)
}

// JobStatus describes a point-in-time snapshot of a registered background job.
type JobStatus struct {
	Name           string
	IntervalMs     int64
	Running        bool
	RunCount       int64
	ErrorCount     int64
	LastStatus     string
	LastError      string
	LastDurationMs int64
	LastStartedAt  *time.Time
	LastFinishedAt *time.Time
	NextRunAt      *time.Time
}

// JobInspector provides read-only job telemetry snapshots and manual on-demand
// job triggering without coupling feature modules to internal/platform/jobs.
type JobInspector interface {
	Snapshot() []JobStatus
	Trigger(ctx context.Context, name string) error
}

// ProvisionKeyRequest describes parameters for minting a new gateway API key
// from another module (such as vouchers) without importing internal/modules/keys.
type ProvisionKeyRequest struct {
	Name          string
	AllowedModels []string
	Scopes        []string
	Tags          []string
	RateLimitRPM  int
	QuotaUSD      float64
	ExpiresAt     *string
	Actor         string
}

// ProvisionedKey returns the identifier, prefix, and one-time raw secret of a
// newly minted gateway API key.
type ProvisionedKey struct {
	ID        string
	Name      string
	KeyPrefix string
	Secret    string
}

// KeyProvisioner mints and revokes gateway API keys on behalf of authorized modules.
type KeyProvisioner interface {
	ProvisionKey(ctx context.Context, req ProvisionKeyRequest) (ProvisionedKey, error)
	RevokeKey(ctx context.Context, keyID string) error
}

// TeamAssigner binds a gateway API key to a team for shared budget/RPM grouping.
type TeamAssigner interface {
	AssignKeyToTeam(ctx context.Context, teamID, keyID string) error
}
