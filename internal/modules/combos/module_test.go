package combos_test

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"io"
	"io/fs"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"path/filepath"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/omeg4monx/router-gateway/internal/contracts"
	"github.com/omeg4monx/router-gateway/internal/modules/combos"
	"github.com/omeg4monx/router-gateway/internal/platform/config"
	"github.com/omeg4monx/router-gateway/internal/platform/db"
	"github.com/omeg4monx/router-gateway/internal/platform/eventbus"
	"github.com/omeg4monx/router-gateway/internal/platform/httpserver"
	"github.com/omeg4monx/router-gateway/internal/platform/lifecycle"
)

var (
	_ lifecycle.Module               = (*combos.Module)(nil)
	_ contracts.ProviderResolver     = (*combos.Module)(nil)
	_ contracts.AccountHealthTracker = (*combos.Module)(nil)
	_ contracts.ModelLister          = (*combos.Module)(nil)
)

type stubModule struct {
	name string
	deps []string
}

func (s stubModule) Name() string                         { return s.name }
func (s stubModule) Dependencies() []string               { return s.deps }
func (s stubModule) Migrations() fs.FS                    { return nil }
func (s stubModule) RegisterRoutes(lifecycle.RouteMounts) {}
func (s stubModule) Start(context.Context) error          { return nil }
func (s stubModule) Stop(context.Context) error           { return nil }

type stubAdminAuth struct {
	subject string
}

func (a *stubAdminAuth) AuthenticateAdmin(r *http.Request) (string, error) {
	if r.Header.Get("Authorization") == "Bearer test-admin-token" {
		return a.subject, nil
	}
	return "", errors.New("unauthorized")
}

func (a *stubAdminAuth) RequireLogin(context.Context) bool {
	return true
}

type fakeClock struct {
	mu  sync.Mutex
	cur time.Time
}

func newFakeClock(start time.Time) *fakeClock {
	return &fakeClock{cur: start.UTC()}
}

func (c *fakeClock) Now() time.Time {
	c.mu.Lock()
	defer c.mu.Unlock()
	return c.cur
}

func (c *fakeClock) Advance(d time.Duration) time.Time {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.cur = c.cur.Add(d)
	return c.cur
}

type fakeProviderStore struct {
	mu    sync.RWMutex
	conns []contracts.ProviderConnection
	err   error
}

func newFakeProviderStore(conns ...contracts.ProviderConnection) *fakeProviderStore {
	return &fakeProviderStore{conns: conns}
}

func (s *fakeProviderStore) SetConnections(conns ...contracts.ProviderConnection) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.conns = append([]contracts.ProviderConnection(nil), conns...)
}

func (s *fakeProviderStore) ListEnabledConnections(context.Context) ([]contracts.ProviderConnection, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	if s.err != nil {
		return nil, s.err
	}
	out := make([]contracts.ProviderConnection, 0, len(s.conns))
	for _, c := range s.conns {
		if c.Enabled {
			cloned := c
			cloned.Models = append([]string(nil), c.Models...)
			out = append(out, cloned)
		}
	}
	return out, nil
}

func (s *fakeProviderStore) ResolveCredentials(context.Context, string) (contracts.UpstreamCredentials, error) {
	return contracts.UpstreamCredentials{}, contracts.ErrConnectionNotFound
}

func (s *fakeProviderStore) RefreshOAuthToken(context.Context, string) (contracts.UpstreamCredentials, error) {
	return contracts.UpstreamCredentials{}, contracts.ErrOAuthRefreshUnsupported
}

type testHarness struct {
	db        *db.DB
	bus       *eventbus.Bus
	clock     *fakeClock
	providers *fakeProviderStore
	reg       *lifecycle.Registry
	server    *httpserver.Server
	mod       *combos.Module
	auth      *stubAdminAuth
}

func newTestHarness(t *testing.T, auth *stubAdminAuth, conns ...contracts.ProviderConnection) *testHarness {
	t.Helper()

	cfg := config.Config{
		Host:       "127.0.0.1",
		Port:       20128,
		APITimeout: 30 * time.Second,
	}

	dbPath := filepath.Join(t.TempDir(), "router.sqlite")
	database, err := db.Open(dbPath)
	if err != nil {
		t.Fatalf("db.Open: %v", err)
	}
	t.Cleanup(func() {
		_ = database.Close()
	})

	clock := newFakeClock(time.Date(2026, 11, 23, 12, 0, 0, 0, time.UTC))
	bus := eventbus.New()
	store := newFakeProviderStore(conns...)
	reg := lifecycle.NewRegistry(nil)
	mod := combos.New(combos.Options{
		DB:        database,
		Providers: store,
		Bus:       bus,
		Now:       clock.Now,
	})

	if err := reg.Register(stubModule{name: "system"}); err != nil {
		t.Fatalf("reg.Register(system): %v", err)
	}
	if err := reg.Register(stubModule{name: "providers", deps: []string{"system"}}); err != nil {
		t.Fatalf("reg.Register(providers): %v", err)
	}
	if err := reg.Register(mod); err != nil {
		t.Fatalf("reg.Register(combos): %v", err)
	}

	ctx := context.Background()
	if err := reg.MigrateAll(ctx, database); err != nil {
		t.Fatalf("reg.MigrateAll: %v", err)
	}

	srvOpts := httpserver.Options{
		Config: cfg,
		Logger: slog.New(slog.NewTextHandler(io.Discard, nil)),
	}
	if auth != nil {
		srvOpts.Auth = auth
	}
	srv := httpserver.New(srvOpts)
	if err := reg.MountAll(srv.Mounts()); err != nil {
		t.Fatalf("reg.MountAll: %v", err)
	}
	if err := reg.StartAll(ctx); err != nil {
		t.Fatalf("reg.StartAll: %v", err)
	}
	t.Cleanup(func() {
		_ = reg.StopAll(context.Background())
	})

	return &testHarness{
		db:        database,
		bus:       bus,
		clock:     clock,
		providers: store,
		reg:       reg,
		server:    srv,
		mod:       mod,
		auth:      auth,
	}
}

func (h *testHarness) doJSON(t *testing.T, method, path string, body any, authenticated bool) *httptest.ResponseRecorder {
	t.Helper()

	var reader io.Reader = http.NoBody
	if body != nil {
		switch v := body.(type) {
		case string:
			reader = strings.NewReader(v)
		case []byte:
			reader = bytes.NewReader(v)
		default:
			raw, err := json.Marshal(body)
			if err != nil {
				t.Fatalf("json.Marshal: %v", err)
			}
			reader = bytes.NewReader(raw)
		}
	}

	req := httptest.NewRequest(method, path, reader)
	if body != nil {
		req.Header.Set("Content-Type", "application/json")
	}
	if authenticated {
		req.Header.Set("Authorization", "Bearer test-admin-token")
	}
	rec := httptest.NewRecorder()
	h.server.Handler().ServeHTTP(rec, req)
	return rec
}

func decodeJSON[T any](t *testing.T, rec *httptest.ResponseRecorder) T {
	t.Helper()
	var out T
	if err := json.Unmarshal(rec.Body.Bytes(), &out); err != nil {
		t.Fatalf("unmarshal response (%d %s): %v", rec.Code, rec.Body.String(), err)
	}
	return out
}

func TestModuleMetadataAndStartValidation(t *testing.T) {
	t.Parallel()

	mod := combos.New(combos.Options{})
	if mod.Name() != "combos" {
		t.Fatalf("Name() = %q, want %q", mod.Name(), "combos")
	}
	deps := mod.Dependencies()
	if len(deps) != 2 || deps[0] != "system" || deps[1] != "providers" {
		t.Fatalf("Dependencies() = %v, want [system providers]", deps)
	}
	if err := mod.Start(context.Background()); err == nil {
		t.Fatal("Start with nil DB should fail")
	}

	dbPath := filepath.Join(t.TempDir(), "router.sqlite")
	database, err := db.Open(dbPath)
	if err != nil {
		t.Fatalf("db.Open: %v", err)
	}
	t.Cleanup(func() { _ = database.Close() })

	modNoProviders := combos.New(combos.Options{DB: database})
	if err := modNoProviders.Start(context.Background()); err == nil {
		t.Fatal("Start with nil Providers should fail")
	}
}

func TestCombosAliasesDisabledCRUDAndAuth(t *testing.T) {
	t.Parallel()

	h := newTestHarness(t, &stubAdminAuth{subject: "operator"})

	evCh, unsub := h.bus.Subscribe(combos.TopicCombosChanged, 32)
	t.Cleanup(unsub)

	// Unauthenticated request rejected with 401.
	if rec := h.doJSON(t, http.MethodGet, "/api/combos", nil, false); rec.Code != http.StatusUnauthorized {
		t.Fatalf("GET /api/combos unauthenticated = %d, want 401", rec.Code)
	}

	// Empty list initially.
	rec := h.doJSON(t, http.MethodGet, "/api/combos", nil, true)
	if rec.Code != http.StatusOK {
		t.Fatalf("GET /api/combos = %d: %s", rec.Code, rec.Body.String())
	}
	listResp := decodeJSON[struct {
		Combos []combos.Combo `json:"combos"`
	}](t, rec)
	if len(listResp.Combos) != 0 {
		t.Fatalf("initial combos = %d, want 0", len(listResp.Combos))
	}

	// Validation errors on create.
	badPayloads := []string{
		`{"name":"","targets":["openai/gpt-4o"]}`,
		`{"name":"*","targets":["openai/gpt-4o"]}`,
		`{"name":"bad-empty","targets":[]}`,
		`{"name":"bad-dup","targets":["openai/gpt-4o","openai/gpt-4o"]}`,
		`{"name":"self-ref","targets":["self-ref"]}`,
		`{"name":"bad-strategy","targets":["openai/gpt-4o"],"strategy":"random"}`,
		`{"name":"unknown-field","targets":["openai/gpt-4o"],"extra":true}`,
	}
	for _, body := range badPayloads {
		bad := h.doJSON(t, http.MethodPost, "/api/combos", body, true)
		if bad.Code != http.StatusBadRequest {
			t.Fatalf("POST /api/combos %s = %d, want 400 (body: %s)", body, bad.Code, bad.Body.String())
		}
	}

	// Oversized request body returns 413.
	huge := `{"name":"big","description":"` + strings.Repeat("a", 70*1024) + `","targets":["openai/gpt-4o"]}`
	if rec := h.doJSON(t, http.MethodPost, "/api/combos", huge, true); rec.Code != http.StatusRequestEntityTooLarge {
		t.Fatalf("oversized POST /api/combos = %d, want 413", rec.Code)
	}

	// Create valid combo.
	rec = h.doJSON(t, http.MethodPost, "/api/combos", map[string]any{
		"name":        "prod-fallback",
		"description": "Primary OpenAI with Anthropic fallback",
		"targets":     []string{"openai/gpt-4o", "anthropic/claude-sonnet-4"},
	}, true)
	if rec.Code != http.StatusCreated {
		t.Fatalf("POST /api/combos = %d: %s", rec.Code, rec.Body.String())
	}
	created := decodeJSON[combos.Combo](t, rec)
	if created.ID == "" || created.Name != "prod-fallback" || created.Strategy != "ordered" || !created.Enabled {
		t.Fatalf("unexpected created combo: %+v", created)
	}
	if len(created.Targets) != 2 || created.Targets[0] != "openai/gpt-4o" || created.Targets[1] != "anthropic/claude-sonnet-4" {
		t.Fatalf("unexpected targets: %+v", created.Targets)
	}

	// Duplicate combo name returns 409.
	dupRec := h.doJSON(t, http.MethodPost, "/api/combos", map[string]any{
		"name":    "prod-fallback",
		"targets": []string{"openai/gpt-4o"},
	}, true)
	if dupRec.Code != http.StatusConflict {
		t.Fatalf("duplicate POST /api/combos = %d, want 409", dupRec.Code)
	}

	// GET by ID.
	getRec := h.doJSON(t, http.MethodGet, "/api/combos/"+created.ID, nil, true)
	if getRec.Code != http.StatusOK {
		t.Fatalf("GET /api/combos/{id} = %d: %s", getRec.Code, getRec.Body.String())
	}
	fetched := decodeJSON[combos.Combo](t, getRec)
	if fetched.ID != created.ID || fetched.Name != created.Name {
		t.Fatalf("fetched = %+v, want %+v", fetched, created)
	}

	// GET unknown -> 404.
	if rec := h.doJSON(t, http.MethodGet, "/api/combos/combo_missing", nil, true); rec.Code != http.StatusNotFound {
		t.Fatalf("GET missing combo = %d, want 404", rec.Code)
	}

	// PATCH empty body -> 400.
	if rec := h.doJSON(t, http.MethodPatch, "/api/combos/"+created.ID, map[string]any{}, true); rec.Code != http.StatusBadRequest {
		t.Fatalf("empty PATCH /api/combos/{id} = %d, want 400", rec.Code)
	}

	// PATCH updates targets and enabled.
	h.clock.Advance(5 * time.Second)
	patchRec := h.doJSON(t, http.MethodPatch, "/api/combos/"+created.ID, map[string]any{
		"targets": []string{"anthropic/claude-sonnet-4", "openai/gpt-4o"},
		"enabled": false,
	}, true)
	if patchRec.Code != http.StatusOK {
		t.Fatalf("PATCH /api/combos/{id} = %d: %s", patchRec.Code, patchRec.Body.String())
	}
	updated := decodeJSON[combos.Combo](t, patchRec)
	if updated.Enabled || updated.Targets[0] != "anthropic/claude-sonnet-4" || updated.UpdatedAt == created.UpdatedAt {
		t.Fatalf("unexpected patched combo: %+v", updated)
	}

	// DELETE combo.
	if rec := h.doJSON(t, http.MethodDelete, "/api/combos/"+created.ID, nil, true); rec.Code != http.StatusNoContent {
		t.Fatalf("DELETE /api/combos/{id} = %d, want 204", rec.Code)
	}
	if rec := h.doJSON(t, http.MethodDelete, "/api/combos/"+created.ID, nil, true); rec.Code != http.StatusNotFound {
		t.Fatalf("second DELETE /api/combos/{id} = %d, want 404", rec.Code)
	}

	// Aliases CRUD.
	if rec := h.doJSON(t, http.MethodPut, "/api/aliases/fast", map[string]any{"target_model": "fast"}, true); rec.Code != http.StatusBadRequest {
		t.Fatalf("self-referential PUT /api/aliases/fast = %d, want 400", rec.Code)
	}
	putAlias := h.doJSON(t, http.MethodPut, "/api/aliases/fast", map[string]any{"target_model": "openai/gpt-4o-mini"}, true)
	if putAlias.Code != http.StatusOK {
		t.Fatalf("PUT /api/aliases/fast = %d: %s", putAlias.Code, putAlias.Body.String())
	}
	aliasObj := decodeJSON[combos.ModelAlias](t, putAlias)
	if aliasObj.Alias != "fast" || aliasObj.TargetModel != "openai/gpt-4o-mini" {
		t.Fatalf("unexpected alias: %+v", aliasObj)
	}

	// Upsert existing alias updates target_model.
	h.clock.Advance(2 * time.Second)
	putAlias2 := h.doJSON(t, http.MethodPut, "/api/aliases/fast", map[string]any{"target_model": "gemini/gemini-2.5-flash"}, true)
	if putAlias2.Code != http.StatusOK {
		t.Fatalf("second PUT /api/aliases/fast = %d: %s", putAlias2.Code, putAlias2.Body.String())
	}
	aliasObj2 := decodeJSON[combos.ModelAlias](t, putAlias2)
	if aliasObj2.TargetModel != "gemini/gemini-2.5-flash" || aliasObj2.CreatedAt != aliasObj.CreatedAt {
		t.Fatalf("unexpected upserted alias: %+v", aliasObj2)
	}

	listAliasesRec := h.doJSON(t, http.MethodGet, "/api/aliases", nil, true)
	if listAliasesRec.Code != http.StatusOK {
		t.Fatalf("GET /api/aliases = %d", listAliasesRec.Code)
	}
	aliasesResp := decodeJSON[struct {
		Aliases []combos.ModelAlias `json:"aliases"`
	}](t, listAliasesRec)
	if len(aliasesResp.Aliases) != 1 || aliasesResp.Aliases[0].Alias != "fast" {
		t.Fatalf("unexpected aliases list: %+v", aliasesResp.Aliases)
	}

	if rec := h.doJSON(t, http.MethodDelete, "/api/aliases/fast", nil, true); rec.Code != http.StatusNoContent {
		t.Fatalf("DELETE /api/aliases/fast = %d, want 204", rec.Code)
	}
	if rec := h.doJSON(t, http.MethodDelete, "/api/aliases/fast", nil, true); rec.Code != http.StatusNotFound {
		t.Fatalf("second DELETE /api/aliases/fast = %d, want 404", rec.Code)
	}

	// Disabled models CRUD (query param and body deletion for models containing '/').
	dis1 := h.doJSON(t, http.MethodPut, "/api/models/disabled", map[string]any{
		"model":  "openai/gpt-4o",
		"reason": "maintenance window",
	}, true)
	if dis1.Code != http.StatusOK {
		t.Fatalf("PUT /api/models/disabled (qualified) = %d: %s", dis1.Code, dis1.Body.String())
	}
	dis2 := h.doJSON(t, http.MethodPut, "/api/models/disabled", map[string]any{
		"model":  "o3-mini",
		"reason": "cost control",
	}, true)
	if dis2.Code != http.StatusOK {
		t.Fatalf("PUT /api/models/disabled (bare) = %d: %s", dis2.Code, dis2.Body.String())
	}

	listDisRec := h.doJSON(t, http.MethodGet, "/api/models/disabled", nil, true)
	if listDisRec.Code != http.StatusOK {
		t.Fatalf("GET /api/models/disabled = %d", listDisRec.Code)
	}
	disResp := decodeJSON[struct {
		DisabledModels []combos.DisabledModel `json:"disabled_models"`
	}](t, listDisRec)
	if len(disResp.DisabledModels) != 2 {
		t.Fatalf("disabled_models count = %d, want 2", len(disResp.DisabledModels))
	}

	// Delete via query param (slash preserved).
	if rec := h.doJSON(t, http.MethodDelete, "/api/models/disabled?model=openai/gpt-4o", nil, true); rec.Code != http.StatusNoContent {
		t.Fatalf("DELETE /api/models/disabled?model=openai/gpt-4o = %d: %s", rec.Code, rec.Body.String())
	}
	// Delete via JSON body.
	if rec := h.doJSON(t, http.MethodDelete, "/api/models/disabled", map[string]any{"model": "o3-mini"}, true); rec.Code != http.StatusNoContent {
		t.Fatalf("DELETE /api/models/disabled (body) = %d: %s", rec.Code, rec.Body.String())
	}
	// Missing model on DELETE returns 400.
	if rec := h.doJSON(t, http.MethodDelete, "/api/models/disabled", nil, true); rec.Code != http.StatusBadRequest {
		t.Fatalf("DELETE /api/models/disabled without model = %d, want 400", rec.Code)
	}
	// Deleting non-existent disabled model returns 404.
	if rec := h.doJSON(t, http.MethodDelete, "/api/models/disabled?model=o3-mini", nil, true); rec.Code != http.StatusNotFound {
		t.Fatalf("DELETE missing disabled model = %d, want 404", rec.Code)
	}

	gotEvents := 0
	for {
		select {
		case ev := <-evCh:
			ce, ok := ev.Payload.(combos.ChangedEvent)
			if !ok {
				t.Fatalf("unexpected event payload type %T", ev.Payload)
			}
			if ce.ChangedBy != "operator" {
				t.Fatalf("ChangedBy = %q, want operator", ce.ChangedBy)
			}
			gotEvents++
		default:
			goto doneDraining
		}
	}
doneDraining:
	if gotEvents < 8 {
		t.Fatalf("expected >= 8 combos.changed events, got %d", gotEvents)
	}
}

func TestResolverFourTierPrecedenceAndTargetFields(t *testing.T) {
	t.Parallel()

	conns := []contracts.ProviderConnection{
		{
			ID:         "conn_openai_low",
			ProviderID: "openai",
			Name:       "OpenAI Secondary",
			Protocol:   "openai-chat",
			BaseURL:    "https://api.openai.com/v1",
			AuthType:   "api_key",
			Priority:   20,
			Enabled:    true,
			Models:     []string{"shared-model", "gpt-4o"},
		},
		{
			ID:         "conn_openai_b",
			ProviderID: "openai",
			Name:       "OpenAI Primary B",
			Protocol:   "openai-chat",
			BaseURL:    "https://api.openai.com/v1",
			AuthType:   "api_key",
			Priority:   10,
			Enabled:    true,
			Models:     []string{"shared-model", "gpt-4o"},
		},
		{
			ID:         "conn_openai_a",
			ProviderID: "openai",
			Name:       "OpenAI Primary A",
			Protocol:   "openai-chat",
			BaseURL:    "https://api.openai.com/v1",
			AuthType:   "api_key",
			Priority:   10,
			Enabled:    true,
			Models:     []string{"shared-model", "gpt-4o"},
		},
		{
			ID:         "conn_anthropic",
			ProviderID: "anthropic",
			Name:       "Anthropic OAuth",
			Protocol:   "anthropic",
			BaseURL:    "https://api.anthropic.com",
			AuthType:   "oauth",
			Priority:   5,
			Enabled:    true,
			Models:     []string{"claude-sonnet-4"},
		},
		{
			ID:         "conn_gemini",
			ProviderID: "gemini",
			Name:       "Gemini Primary",
			Protocol:   "gemini",
			BaseURL:    "https://generativelanguage.googleapis.com",
			AuthType:   "api_key",
			Priority:   15,
			Enabled:    true,
			Models:     []string{"gemini-2.5-pro"},
		},
		{
			ID:         "conn_compat_empty",
			ProviderID: "openai-compatible",
			Name:       "Local vLLM",
			Protocol:   "openai-chat",
			BaseURL:    "http://127.0.0.1:8000/v1",
			AuthType:   "api_key",
			Priority:   1,
			Enabled:    true,
			Models:     []string{},
		},
		{
			ID:         "conn_openrouter",
			ProviderID: "openrouter",
			Name:       "OpenRouter",
			Protocol:   "openai-chat",
			BaseURL:    "https://openrouter.ai/api/v1",
			AuthType:   "api_key",
			Priority:   30,
			Enabled:    true,
			Models:     []string{"openai/gpt-4o"},
		},
	}

	h := newTestHarness(t, nil, conns...)
	ctx := context.Background()

	// Configure Tier 1 (Active Combo) AND Tier 2 (Alias) both named "shared-model".
	comboRec := h.doJSON(t, http.MethodPost, "/api/combos", map[string]any{
		"name":    "shared-model",
		"targets": []string{"anthropic/claude-sonnet-4"},
	}, false)
	if comboRec.Code != http.StatusCreated {
		t.Fatalf("create combo = %d: %s", comboRec.Code, comboRec.Body.String())
	}
	createdCombo := decodeJSON[combos.Combo](t, comboRec)

	if rec := h.doJSON(t, http.MethodPut, "/api/aliases/shared-model", map[string]any{
		"target_model": "gemini/gemini-2.5-pro",
	}, false); rec.Code != http.StatusOK {
		t.Fatalf("create alias = %d: %s", rec.Code, rec.Body.String())
	}

	// Tier 1 beats Tier 2 and Tier 4!
	targets, err := h.mod.Resolve(ctx, "shared-model", nil)
	if err != nil {
		t.Fatalf("Resolve(Tier 1) err = %v", err)
	}
	if len(targets) != 1 || targets[0].ConnectionID != "conn_anthropic" {
		t.Fatalf("Tier 1 targets = %+v, want [conn_anthropic]", targets)
	}
	// Verify Pre-Flight Ruling exact ResolvedTarget fields.
	got := targets[0]
	if got.ProviderID != "anthropic" ||
		got.ProviderType != "anthropic" ||
		got.ConnectionID != "conn_anthropic" ||
		got.Model != "shared-model" ||
		got.UpstreamModel != "claude-sonnet-4" ||
		got.Endpoint != "https://api.anthropic.com" ||
		got.Priority != 5 ||
		got.Metadata["auth_type"] != "oauth" ||
		got.Metadata["connection_name"] != "Anthropic OAuth" {
		t.Fatalf("ResolvedTarget fields mismatch: %+v", got)
	}

	// Disable the combo -> Tier 2 (Alias) now beats Tier 4!
	if rec := h.doJSON(t, http.MethodPatch, "/api/combos/"+createdCombo.ID, map[string]any{
		"enabled": false,
	}, false); rec.Code != http.StatusOK {
		t.Fatalf("disable combo = %d", rec.Code)
	}
	targets, err = h.mod.Resolve(ctx, "shared-model", nil)
	if err != nil {
		t.Fatalf("Resolve(Tier 2) err = %v", err)
	}
	if len(targets) != 1 || targets[0].ConnectionID != "conn_gemini" || targets[0].UpstreamModel != "gemini-2.5-pro" || targets[0].ProviderType != "gemini" {
		t.Fatalf("Tier 2 targets = %+v, want [conn_gemini]", targets)
	}

	// Delete the alias -> Tier 4 (Bare model) ordered strictly by (Priority ASC, ID ASC), skipping empty-models compat connection!
	if rec := h.doJSON(t, http.MethodDelete, "/api/aliases/shared-model", nil, false); rec.Code != http.StatusNoContent {
		t.Fatalf("delete alias = %d", rec.Code)
	}
	targets, err = h.mod.Resolve(ctx, "shared-model", nil)
	if err != nil {
		t.Fatalf("Resolve(Tier 4) err = %v", err)
	}
	if len(targets) != 3 {
		t.Fatalf("Tier 4 len = %d (%+v), want 3", len(targets), targets)
	}
	if targets[0].ConnectionID != "conn_openai_a" || targets[1].ConnectionID != "conn_openai_b" || targets[2].ConnectionID != "conn_openai_low" {
		t.Fatalf("Tier 4 ordering = [%s, %s, %s], want [conn_openai_a, conn_openai_b, conn_openai_low]",
			targets[0].ConnectionID, targets[1].ConnectionID, targets[2].ConnectionID)
	}

	// Tier 3 (Explicit <provider>/<model>) supports empty Models list on matched provider and preserves nested model IDs.
	targets, err = h.mod.Resolve(ctx, "openai-compatible/llama-3.3-70b", nil)
	if err != nil {
		t.Fatalf("Resolve(Tier 3 empty models) err = %v", err)
	}
	if len(targets) != 1 || targets[0].ConnectionID != "conn_compat_empty" || targets[0].UpstreamModel != "llama-3.3-70b" || targets[0].Model != "openai-compatible/llama-3.3-70b" {
		t.Fatalf("Tier 3 empty-models target = %+v", targets)
	}

	// Tier 3 with nested slash model ID (openrouter/openai/gpt-4o).
	targets, err = h.mod.Resolve(ctx, "openrouter/openai/gpt-4o", nil)
	if err != nil {
		t.Fatalf("Resolve(Tier 3 nested slash) err = %v", err)
	}
	if len(targets) != 1 || targets[0].ConnectionID != "conn_openrouter" || targets[0].UpstreamModel != "openai/gpt-4o" {
		t.Fatalf("Tier 3 nested slash target = %+v", targets)
	}

	// Tier 2 (Alias) beats Tier 3 (Explicit <provider>/<model>) when alias is named "<provider>/<model>".
	if rec := h.doJSON(t, http.MethodPut, "/api/aliases/openai%2Fgpt-4o", map[string]any{
		"target_model": "anthropic/claude-sonnet-4",
	}, false); rec.Code != http.StatusOK {
		t.Fatalf("create explicit-named alias = %d: %s", rec.Code, rec.Body.String())
	}
	targets, err = h.mod.Resolve(ctx, "openai/gpt-4o", nil)
	if err != nil {
		t.Fatalf("Resolve(Tier 2 beating Tier 3) err = %v", err)
	}
	if len(targets) != 1 || targets[0].ConnectionID != "conn_anthropic" {
		t.Fatalf("Tier 2 beating Tier 3 = %+v, want [conn_anthropic]", targets)
	}
}

func TestResolverAllowedModelsAuthorization(t *testing.T) {
	t.Parallel()

	conns := []contracts.ProviderConnection{
		{
			ID:         "conn_openai",
			ProviderID: "openai",
			Name:       "OpenAI",
			Protocol:   "openai-chat",
			BaseURL:    "https://api.openai.com/v1",
			AuthType:   "api_key",
			Priority:   10,
			Enabled:    true,
			Models:     []string{"gpt-4o"},
		},
		{
			ID:         "conn_openai_compat",
			ProviderID: "openai-compat",
			Name:       "OpenAI Compat",
			Protocol:   "openai-chat",
			BaseURL:    "https://compat.example.test/v1",
			AuthType:   "api_key",
			Priority:   10,
			Enabled:    true,
			Models:     []string{"gpt-4o"},
		},
		{
			ID:         "conn_anthropic",
			ProviderID: "anthropic",
			Name:       "Anthropic",
			Protocol:   "anthropic",
			BaseURL:    "https://api.anthropic.com",
			AuthType:   "api_key",
			Priority:   10,
			Enabled:    true,
			Models:     []string{"claude-sonnet-4"},
		},
	}

	h := newTestHarness(t, nil, conns...)
	ctx := context.Background()

	if rec := h.doJSON(t, http.MethodPost, "/api/combos", map[string]any{
		"name":    "team-combo",
		"targets": []string{"openai/gpt-4o", "anthropic/claude-sonnet-4"},
	}, false); rec.Code != http.StatusCreated {
		t.Fatalf("create combo = %d", rec.Code)
	}

	// Empty allowed and wildcard "*" allow all.
	if targets, err := h.mod.Resolve(ctx, "openai/gpt-4o", []string{}); err != nil || len(targets) != 1 {
		t.Fatalf("empty allowed: targets=%v err=%v", targets, err)
	}
	if targets, err := h.mod.Resolve(ctx, "anthropic/claude-sonnet-4", []string{"*"}); err != nil || len(targets) != 1 {
		t.Fatalf("wildcard allowed: targets=%v err=%v", targets, err)
	}

	// Authorizing combo name authorizes all of its member targets.
	targets, err := h.mod.Resolve(ctx, "team-combo", []string{"team-combo"})
	if err != nil {
		t.Fatalf("Resolve(team-combo) with allowed=[team-combo] err = %v", err)
	}
	if len(targets) != 2 || targets[0].ConnectionID != "conn_openai" || targets[1].ConnectionID != "conn_anthropic" {
		t.Fatalf("combo targets = %+v, want [conn_openai, conn_anthropic]", targets)
	}

	// Directly requesting a member target with allowed=[team-combo] returns ErrModelForbidden.
	if _, err := h.mod.Resolve(ctx, "openai/gpt-4o", []string{"team-combo"}); !errors.Is(err, contracts.ErrModelForbidden) {
		t.Fatalf("Resolve(openai/gpt-4o) with allowed=[team-combo] err = %v, want ErrModelForbidden", err)
	}

	// Prefix wildcard <provider>/* matches only exact provider prefix.
	if targets, err := h.mod.Resolve(ctx, "openai/gpt-4o", []string{"openai/*"}); err != nil || len(targets) != 1 {
		t.Fatalf("prefix wildcard allowed: targets=%v err=%v", targets, err)
	}
	if _, err := h.mod.Resolve(ctx, "openai-compat/gpt-4o", []string{"openai/*"}); !errors.Is(err, contracts.ErrModelForbidden) {
		t.Fatalf("prefix wildcard should not match openai-compat: err = %v, want ErrModelForbidden", err)
	}
	if _, err := h.mod.Resolve(ctx, "anthropic/claude-sonnet-4", []string{"openai/*"}); !errors.Is(err, contracts.ErrModelForbidden) {
		t.Fatalf("prefix wildcard should reject anthropic: err = %v, want ErrModelForbidden", err)
	}
}

func TestDisabledModelsAndCooldownStateMachine(t *testing.T) {
	t.Parallel()

	conns := []contracts.ProviderConnection{
		{
			ID:         "conn_openai_1",
			ProviderID: "openai",
			Name:       "OpenAI 1",
			Protocol:   "openai-chat",
			BaseURL:    "https://api.openai.com/v1",
			AuthType:   "api_key",
			Priority:   10,
			Enabled:    true,
			Models:     []string{"gpt-4o"},
		},
		{
			ID:         "conn_azure_1",
			ProviderID: "azure",
			Name:       "Azure 1",
			Protocol:   "openai-chat",
			BaseURL:    "https://azure.example.test/v1",
			AuthType:   "api_key",
			Priority:   20,
			Enabled:    true,
			Models:     []string{"gpt-4o"},
		},
	}

	h := newTestHarness(t, nil, conns...)
	ctx := context.Background()

	// Disable exact provider-qualified model openai/gpt-4o -> skips openai, keeps azure on bare lookup.
	if rec := h.doJSON(t, http.MethodPut, "/api/models/disabled", map[string]any{
		"model": "openai/gpt-4o",
	}, false); rec.Code != http.StatusOK {
		t.Fatalf("disable openai/gpt-4o = %d", rec.Code)
	}
	if _, err := h.mod.Resolve(ctx, "openai/gpt-4o", nil); !errors.Is(err, contracts.ErrModelNotFound) {
		t.Fatalf("Resolve(disabled openai/gpt-4o) = %v, want ErrModelNotFound", err)
	}
	targets, err := h.mod.Resolve(ctx, "gpt-4o", nil)
	if err != nil || len(targets) != 1 || targets[0].ConnectionID != "conn_azure_1" {
		t.Fatalf("Resolve(gpt-4o with openai/gpt-4o disabled) = %+v, err=%v; want [conn_azure_1]", targets, err)
	}

	// Disable bare model gpt-4o -> skips both openai and azure!
	if rec := h.doJSON(t, http.MethodPut, "/api/models/disabled", map[string]any{
		"model": "gpt-4o",
	}, false); rec.Code != http.StatusOK {
		t.Fatalf("disable gpt-4o = %d", rec.Code)
	}
	if _, err := h.mod.Resolve(ctx, "azure/gpt-4o", nil); !errors.Is(err, contracts.ErrModelNotFound) {
		t.Fatalf("Resolve(azure/gpt-4o with bare gpt-4o disabled) = %v, want ErrModelNotFound", err)
	}
	if _, err := h.mod.Resolve(ctx, "gpt-4o", nil); !errors.Is(err, contracts.ErrModelNotFound) {
		t.Fatalf("Resolve(bare gpt-4o disabled) = %v, want ErrModelNotFound", err)
	}

	// Re-enable both models.
	_ = h.doJSON(t, http.MethodDelete, "/api/models/disabled?model=openai/gpt-4o", nil, false)
	_ = h.doJSON(t, http.MethodDelete, "/api/models/disabled?model=gpt-4o", nil, false)

	// Cooldown progression on 429: 15s -> 30s -> 60s -> 120s -> 300s -> 300s.
	expectedBackoffs := []time.Duration{
		15 * time.Second,
		30 * time.Second,
		60 * time.Second,
		120 * time.Second,
		300 * time.Second,
		300 * time.Second,
	}
	for i, wantDur := range expectedBackoffs {
		start := h.clock.Now()
		h.mod.RecordFailure("conn_openai_1", http.StatusTooManyRequests, 0)
		if !h.mod.IsInCooldown("conn_openai_1", start.Add(wantDur-time.Nanosecond)) {
			t.Fatalf("failure #%d: expected in cooldown at +%v", i+1, wantDur-time.Nanosecond)
		}
		if h.mod.IsInCooldown("conn_openai_1", start.Add(wantDur)) {
			t.Fatalf("failure #%d: expected cooldown expired at +%v", i+1, wantDur)
		}
		h.clock.Advance(wantDur)
	}

	// RecordSuccess resets failure count immediately so next failure starts at 15s again.
	h.mod.RecordFailure("conn_openai_1", http.StatusInternalServerError, 0)
	if !h.mod.IsInCooldown("conn_openai_1", h.clock.Now()) {
		t.Fatal("expected conn_openai_1 in cooldown")
	}
	h.mod.RecordSuccess("conn_openai_1")
	if h.mod.IsInCooldown("conn_openai_1", h.clock.Now()) {
		t.Fatal("RecordSuccess should clear cooldown immediately")
	}
	afterReset := h.clock.Now()
	h.mod.RecordFailure("conn_openai_1", http.StatusTooManyRequests, 0)
	if !h.mod.IsInCooldown("conn_openai_1", afterReset.Add(14*time.Second)) || h.mod.IsInCooldown("conn_openai_1", afterReset.Add(15*time.Second)) {
		t.Fatal("after RecordSuccess, first failure should use 15s base backoff")
	}

	// 401/403 starts at 60s -> 120s -> 300s.
	h.mod.RecordSuccess("conn_azure_1")
	tAuth := h.clock.Now()
	h.mod.RecordFailure("conn_azure_1", http.StatusUnauthorized, 0)
	if !h.mod.IsInCooldown("conn_azure_1", tAuth.Add(59*time.Second)) || h.mod.IsInCooldown("conn_azure_1", tAuth.Add(60*time.Second)) {
		t.Fatal("first 401 failure should start at 60s")
	}
	h.clock.Advance(60 * time.Second)
	tAuth2 := h.clock.Now()
	h.mod.RecordFailure("conn_azure_1", http.StatusForbidden, 0)
	if !h.mod.IsInCooldown("conn_azure_1", tAuth2.Add(119*time.Second)) || h.mod.IsInCooldown("conn_azure_1", tAuth2.Add(120*time.Second)) {
		t.Fatal("second 403 failure should use 120s")
	}

	// Subsequent shorter-backoff failure must NOT shrink an active longer Retry-After window.
	h.mod.RecordSuccess("conn_openai_1")
	tLong := h.clock.Now()
	h.mod.RecordFailure("conn_openai_1", http.StatusTooManyRequests, 120*time.Second)
	h.clock.Advance(5 * time.Second)
	// Second failure has 30s base backoff -> T+35s < T+120s, so T+120s must be preserved.
	h.mod.RecordFailure("conn_openai_1", http.StatusInternalServerError, 0)
	if !h.mod.IsInCooldown("conn_openai_1", tLong.Add(119*time.Second)) {
		t.Fatal("subsequent shorter-backoff failure must not shrink active Retry-After window")
	}
	if h.mod.IsInCooldown("conn_openai_1", tLong.Add(120*time.Second)) {
		t.Fatal("preserved Retry-After window should expire at +120s")
	}

	// Retry-After honored when larger than base backoff, ignored when smaller.
	h.mod.RecordSuccess("conn_openai_1")
	tRetry := h.clock.Now()
	h.mod.RecordFailure("conn_openai_1", http.StatusTooManyRequests, 95*time.Second)
	if !h.mod.IsInCooldown("conn_openai_1", tRetry.Add(94*time.Second)) || h.mod.IsInCooldown("conn_openai_1", tRetry.Add(95*time.Second)) {
		t.Fatal("Retry-After 95s should override 15s base backoff")
	}
	h.mod.RecordSuccess("conn_openai_1")
	tRetrySmall := h.clock.Now()
	h.mod.RecordFailure("conn_openai_1", http.StatusTooManyRequests, 5*time.Second)
	if !h.mod.IsInCooldown("conn_openai_1", tRetrySmall.Add(14*time.Second)) || h.mod.IsInCooldown("conn_openai_1", tRetrySmall.Add(15*time.Second)) {
		t.Fatal("Retry-After 5s smaller than 15s base should still use 15s")
	}

	// Resolver skips cooling-down candidates and returns ErrAllAccountsCoolingDown when all are cooling down.
	// Right now both conn_openai_1 (15s left) and conn_azure_1 (120s left) are cooling down.
	if _, err := h.mod.Resolve(ctx, "gpt-4o", nil); !errors.Is(err, contracts.ErrAllAccountsCoolingDown) {
		t.Fatalf("Resolve(all cooling down) = %v, want ErrAllAccountsCoolingDown", err)
	}
	// Non-existent model still returns ErrModelNotFound (not ErrAllAccountsCoolingDown).
	if _, err := h.mod.Resolve(ctx, "does-not-exist", nil); !errors.Is(err, contracts.ErrModelNotFound) {
		t.Fatalf("Resolve(nonexistent) = %v, want ErrModelNotFound", err)
	}

	// GET /api/models/cooldowns reports both active cooldowns.
	cdRec := h.doJSON(t, http.MethodGet, "/api/models/cooldowns", nil, false)
	if cdRec.Code != http.StatusOK {
		t.Fatalf("GET /api/models/cooldowns = %d: %s", cdRec.Code, cdRec.Body.String())
	}
	cdResp := decodeJSON[struct {
		Cooldowns []combos.CooldownEntry `json:"cooldowns"`
	}](t, cdRec)
	if len(cdResp.Cooldowns) != 2 {
		t.Fatalf("active cooldowns = %+v, want 2", cdResp.Cooldowns)
	}

	// DELETE /api/models/cooldowns/{connectionId} manually clears cooldown.
	if rec := h.doJSON(t, http.MethodDelete, "/api/models/cooldowns/conn_openai_1", nil, false); rec.Code != http.StatusNoContent {
		t.Fatalf("DELETE /api/models/cooldowns/conn_openai_1 = %d, want 204", rec.Code)
	}
	targets, err = h.mod.Resolve(ctx, "gpt-4o", nil)
	if err != nil || len(targets) != 1 || targets[0].ConnectionID != "conn_openai_1" {
		t.Fatalf("Resolve after manual clear = %+v, err=%v; want [conn_openai_1]", targets, err)
	}

	// Advance past conn_azure_1 cooldown -> GET /api/models/cooldowns is now empty and both candidates resolve.
	h.clock.Advance(121 * time.Second)
	cdRec2 := h.doJSON(t, http.MethodGet, "/api/models/cooldowns", nil, false)
	cdResp2 := decodeJSON[struct {
		Cooldowns []combos.CooldownEntry `json:"cooldowns"`
	}](t, cdRec2)
	if len(cdResp2.Cooldowns) != 0 {
		t.Fatalf("expired cooldowns should not be listed: %+v", cdResp2.Cooldowns)
	}
	targets, err = h.mod.Resolve(ctx, "gpt-4o", nil)
	if err != nil || len(targets) != 2 {
		t.Fatalf("Resolve after cooldown expiry = %+v, err=%v; want 2 targets", targets, err)
	}
}

func TestResolverCycleProtectionAndDeduplication(t *testing.T) {
	t.Parallel()

	conns := []contracts.ProviderConnection{
		{
			ID:         "conn_openai",
			ProviderID: "openai",
			Name:       "OpenAI",
			Protocol:   "openai-chat",
			BaseURL:    "https://api.openai.com/v1",
			AuthType:   "api_key",
			Priority:   10,
			Enabled:    true,
			Models:     []string{"gpt-4o"},
		},
		{
			ID:         "conn_anthropic",
			ProviderID: "anthropic",
			Name:       "Anthropic",
			Protocol:   "anthropic",
			BaseURL:    "https://api.anthropic.com",
			AuthType:   "api_key",
			Priority:   10,
			Enabled:    true,
			Models:     []string{"claude-sonnet-4"},
		},
	}

	h := newTestHarness(t, nil, conns...)
	ctx := context.Background()

	// 1. Self-referential alias seeded directly in SQLite (bypassing HTTP guard).
	nowStr := h.clock.Now().Format("2006-01-02T15:04:05.000000000Z")
	if _, err := h.db.Writer().ExecContext(
		ctx,
		`INSERT INTO model_aliases (alias, target_model, created_at, updated_at) VALUES (?, ?, ?, ?)`,
		"self-loop", "self-loop", nowStr, nowStr,
	); err != nil {
		t.Fatalf("seed self-loop alias: %v", err)
	}
	if _, err := h.mod.Resolve(ctx, "self-loop", nil); !errors.Is(err, contracts.ErrModelNotFound) {
		t.Fatalf("Resolve(self-loop) = %v, want ErrModelNotFound", err)
	}

	// 2. Mutual alias cycle: alpha -> beta -> alpha.
	_ = h.doJSON(t, http.MethodPut, "/api/aliases/alpha", map[string]any{"target_model": "beta"}, false)
	_ = h.doJSON(t, http.MethodPut, "/api/aliases/beta", map[string]any{"target_model": "alpha"}, false)
	if _, err := h.mod.Resolve(ctx, "alpha", nil); !errors.Is(err, contracts.ErrModelNotFound) {
		t.Fatalf("Resolve(alpha cycle) = %v, want ErrModelNotFound", err)
	}

	// 3. Combo cycle with healthy fallback targets: combo-a -> [combo-b, openai/gpt-4o], combo-b -> [combo-a, anthropic/claude-sonnet-4].
	_ = h.doJSON(t, http.MethodPost, "/api/combos", map[string]any{
		"name":    "combo-a",
		"targets": []string{"combo-b", "openai/gpt-4o"},
	}, false)
	_ = h.doJSON(t, http.MethodPost, "/api/combos", map[string]any{
		"name":    "combo-b",
		"targets": []string{"combo-a", "anthropic/claude-sonnet-4"},
	}, false)
	targets, err := h.mod.Resolve(ctx, "combo-a", nil)
	if err != nil {
		t.Fatalf("Resolve(combo-a cycle) err = %v", err)
	}
	if len(targets) != 2 || targets[0].ConnectionID != "conn_anthropic" || targets[1].ConnectionID != "conn_openai" {
		t.Fatalf("combo-a targets = %+v, want [conn_anthropic, conn_openai]", targets)
	}

	// 4. Combo -> Alias -> Combo cycle.
	_ = h.doJSON(t, http.MethodPost, "/api/combos", map[string]any{
		"name":    "combo-loop",
		"targets": []string{"alias-to-combo"},
	}, false)
	_ = h.doJSON(t, http.MethodPut, "/api/aliases/alias-to-combo", map[string]any{
		"target_model": "combo-loop",
	}, false)
	if _, err := h.mod.Resolve(ctx, "combo-loop", nil); !errors.Is(err, contracts.ErrModelNotFound) {
		t.Fatalf("Resolve(combo-loop) = %v, want ErrModelNotFound", err)
	}

	// 5. Deduplication when combo targets overlap on the same connection + upstream model.
	_ = h.doJSON(t, http.MethodPost, "/api/combos", map[string]any{
		"name":    "dedup-combo",
		"targets": []string{"openai/gpt-4o", "gpt-4o"},
	}, false)
	targets, err = h.mod.Resolve(ctx, "dedup-combo", nil)
	if err != nil {
		t.Fatalf("Resolve(dedup-combo) err = %v", err)
	}
	if len(targets) != 1 || targets[0].ConnectionID != "conn_openai" {
		t.Fatalf("dedup-combo targets = %+v, want 1 deduplicated target", targets)
	}
}

func TestListModelsDeduplicationAndFiltering(t *testing.T) {
	t.Parallel()

	conns := []contracts.ProviderConnection{
		{
			ID:         "conn_openai_1",
			ProviderID: "openai",
			Name:       "OpenAI 1",
			Protocol:   "openai-chat",
			BaseURL:    "https://api.openai.com/v1",
			AuthType:   "api_key",
			Priority:   10,
			Enabled:    true,
			Models:     []string{"gpt-4o", "gpt-4o-mini", "*"},
		},
		{
			ID:         "conn_openai_2",
			ProviderID: "openai",
			Name:       "OpenAI 2",
			Protocol:   "openai-chat",
			BaseURL:    "https://api.openai.com/v1",
			AuthType:   "api_key",
			Priority:   20,
			Enabled:    true,
			Models:     []string{"gpt-4o"},
		},
		{
			ID:         "conn_anthropic",
			ProviderID: "anthropic",
			Name:       "Anthropic",
			Protocol:   "anthropic",
			BaseURL:    "https://api.anthropic.com",
			AuthType:   "api_key",
			Priority:   15,
			Enabled:    true,
			Models:     []string{"claude-sonnet-4"},
		},
	}

	h := newTestHarness(t, nil, conns...)
	ctx := context.Background()

	// Active combo, disabled combo, alias, and disabled model.
	_ = h.doJSON(t, http.MethodPost, "/api/combos", map[string]any{
		"name":    "prod-combo",
		"targets": []string{"openai/gpt-4o", "anthropic/claude-sonnet-4"},
	}, false)
	_ = h.doJSON(t, http.MethodPost, "/api/combos", map[string]any{
		"name":    "inactive-combo",
		"targets": []string{"openai/gpt-4o"},
		"enabled": false,
	}, false)
	_ = h.doJSON(t, http.MethodPut, "/api/aliases/fast", map[string]any{
		"target_model": "openai/gpt-4o-mini",
	}, false)
	// Alias with same name as active combo should be shadowed by combo.
	_ = h.doJSON(t, http.MethodPut, "/api/aliases/prod-combo", map[string]any{
		"target_model": "openai/gpt-4o",
	}, false)
	_ = h.doJSON(t, http.MethodPut, "/api/models/disabled", map[string]any{
		"model": "gpt-4o-mini",
	}, false)

	// Unrestricted listing.
	all, err := h.mod.ListModels(ctx, nil)
	if err != nil {
		t.Fatalf("ListModels(nil) err = %v", err)
	}
	byID := make(map[string]contracts.ModelDescriptor, len(all))
	for _, m := range all {
		if _, dup := byID[m.ID]; dup {
			t.Fatalf("duplicate model ID %q in ListModels: %+v", m.ID, all)
		}
		byID[m.ID] = m
	}
	if byID["prod-combo"].Kind != "combo" {
		t.Fatalf("prod-combo descriptor = %+v, want Kind=combo", byID["prod-combo"])
	}
	if _, exists := byID["inactive-combo"]; exists {
		t.Fatal("disabled combo must not appear in ListModels")
	}
	if byID["fast"].Kind != "alias" {
		t.Fatalf("fast descriptor = %+v, want Kind=alias", byID["fast"])
	}
	if _, exists := byID["gpt-4o-mini"]; exists {
		t.Fatal("disabled bare model gpt-4o-mini must not appear")
	}
	if _, exists := byID["openai/gpt-4o-mini"]; exists {
		t.Fatal("qualified openai/gpt-4o-mini must not appear when bare gpt-4o-mini is disabled")
	}
	if _, exists := byID["*"]; exists {
		t.Fatal("literal wildcard '*' must not appear in ListModels")
	}
	if byID["openai/gpt-4o"].Kind != "model" || byID["gpt-4o"].Kind != "model" || byID["anthropic/claude-sonnet-4"].Kind != "model" {
		t.Fatalf("missing expected model descriptors in %+v", byID)
	}

	// Scoped to combo only.
	comboOnly, err := h.mod.ListModels(ctx, []string{"prod-combo"})
	if err != nil {
		t.Fatalf("ListModels([prod-combo]) err = %v", err)
	}
	if len(comboOnly) != 1 || comboOnly[0].ID != "prod-combo" || comboOnly[0].Kind != "combo" {
		t.Fatalf("comboOnly = %+v, want [prod-combo]", comboOnly)
	}

	// Scoped to provider prefix wildcard.
	openaiOnly, err := h.mod.ListModels(ctx, []string{"openai/*"})
	if err != nil {
		t.Fatalf("ListModels([openai/*]) err = %v", err)
	}
	if len(openaiOnly) != 1 || openaiOnly[0].ID != "openai/gpt-4o" {
		t.Fatalf("openaiOnly = %+v, want [openai/gpt-4o]", openaiOnly)
	}
}

func TestProviderQualifiedModelsListAndBareResolveSymmetry(t *testing.T) {
	t.Parallel()

	conns := []contracts.ProviderConnection{
		{
			ID:         "conn_openai_qualified",
			ProviderID: "openai",
			Name:       "OpenAI Qualified",
			Protocol:   "openai-chat",
			BaseURL:    "https://api.openai.com/v1",
			AuthType:   "api_key",
			Priority:   10,
			Enabled:    true,
			Models:     []string{" openai/gpt-4o ", "openai/*", "openai/"},
		},
	}

	h := newTestHarness(t, nil, conns...)
	ctx := context.Background()

	listed, err := h.mod.ListModels(ctx, nil)
	if err != nil {
		t.Fatalf("ListModels err = %v", err)
	}
	if len(listed) != 2 {
		t.Fatalf("ListModels = %+v, want exactly [gpt-4o, openai/gpt-4o] without bogus wildcard/empty entries", listed)
	}
	if listed[0].ID != "gpt-4o" || listed[1].ID != "openai/gpt-4o" {
		t.Fatalf("ListModels IDs = [%s, %s], want [gpt-4o, openai/gpt-4o]", listed[0].ID, listed[1].ID)
	}

	// Bare Tier 4 resolution must resolve models advertised as bare by ListModels.
	targets, err := h.mod.Resolve(ctx, "gpt-4o", nil)
	if err != nil {
		t.Fatalf("Resolve(bare gpt-4o) err = %v", err)
	}
	if len(targets) != 1 || targets[0].ConnectionID != "conn_openai_qualified" || targets[0].UpstreamModel != "gpt-4o" {
		t.Fatalf("Resolve(bare gpt-4o) = %+v, want [conn_openai_qualified/gpt-4o]", targets)
	}

	// Provider-prefixed wildcard openai/* must not match arbitrary unrelated bare models.
	if _, err := h.mod.Resolve(ctx, "claude-sonnet-4", nil); !errors.Is(err, contracts.ErrModelNotFound) {
		t.Fatalf("Resolve(unrelated bare claude-sonnet-4) = %v, want ErrModelNotFound", err)
	}

	// Provider-prefixed wildcard openai/* supports explicit Tier 3 lookup for that provider.
	explicitTargets, err := h.mod.Resolve(ctx, "openai/gpt-4.1", nil)
	if err != nil {
		t.Fatalf("Resolve(explicit openai/gpt-4.1 via openai/*) err = %v", err)
	}
	if len(explicitTargets) != 1 || explicitTargets[0].ConnectionID != "conn_openai_qualified" || explicitTargets[0].UpstreamModel != "gpt-4.1" {
		t.Fatalf("Resolve(explicit openai/gpt-4.1) = %+v, want [conn_openai_qualified/gpt-4.1]", explicitTargets)
	}
}

func TestCooldownConcurrency(t *testing.T) {
	t.Parallel()

	h := newTestHarness(t, nil)

	var wg sync.WaitGroup
	for i := range 16 {
		wg.Add(1)
		go func(n int) {
			defer wg.Done()
			id := "conn_race"
			for range 50 {
				h.mod.RecordFailure(id, http.StatusTooManyRequests, 0)
				_ = h.mod.IsInCooldown(id, h.clock.Now())
				if n%2 == 0 {
					h.mod.RecordSuccess(id)
				}
				_ = h.mod.ListCooldowns(h.clock.Now())
			}
		}(i)
	}
	wg.Wait()
}
