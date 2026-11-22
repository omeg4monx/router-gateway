package contracts_test

import (
	"context"
	"errors"
	"net/http"
	"net/url"
	"testing"
	"time"

	"github.com/omeg4monx/router-gateway/internal/contracts"
)

type stubAuthVerifier struct{}

func (stubAuthVerifier) AuthenticateAdmin(_ *http.Request) (string, error) { return "admin", nil }
func (stubAuthVerifier) RequireLogin(_ context.Context) bool               { return true }

type stubKeyValidator struct{}

func (stubKeyValidator) ValidateKey(_ context.Context, _ string) (contracts.KeyInfo, error) {
	return contracts.KeyInfo{
		ID:            "key_1",
		Name:          "primary",
		AllowedModels: []string{"*"},
	}, nil
}

type stubProviderResolver struct{}

func (stubProviderResolver) Resolve(_ context.Context, model string, _ []string) ([]contracts.ResolvedTarget, error) {
	return []contracts.ResolvedTarget{
		{
			ProviderID:    "openai",
			ProviderType:  "openai-compatible",
			ConnectionID:  "conn_1",
			Model:         model,
			UpstreamModel: model,
			Endpoint:      "https://api.openai.com/v1",
			Priority:      1,
		},
	}, nil
}

type stubInterceptor struct{}

func (stubInterceptor) Name() string  { return "noop" }
func (stubInterceptor) Priority() int { return 100 }
func (stubInterceptor) OnRequest(_ context.Context, _ *contracts.CanonicalRequest) error {
	return nil
}
func (stubInterceptor) OnResponseHeaders(_ context.Context, _ *contracts.ResponseMeta) error {
	return nil
}
func (stubInterceptor) OnStreamChunk(_ context.Context, _ *contracts.StreamEvent) error {
	return nil
}

type stubUsageRecorder struct {
	recorded []contracts.UsageRecord
}

func (s *stubUsageRecorder) Record(_ context.Context, rec contracts.UsageRecord) {
	s.recorded = append(s.recorded, rec)
}

type stubProxyDialerResolver struct{}

func (stubProxyDialerResolver) ProxyURLForProvider(_ context.Context, _ string) (*url.URL, error) {
	return url.Parse("http://127.0.0.1:8080")
}

type stubAdminCredentialStore struct {
	hash         string
	requireLogin bool
}

var _ contracts.AdminCredentialStore = (*stubAdminCredentialStore)(nil)

func (s *stubAdminCredentialStore) AdminPasswordHash(_ context.Context) (string, error) {
	if s.hash == "" {
		return "", contracts.ErrNotInitialized
	}
	return s.hash, nil
}

func (s *stubAdminCredentialStore) UpdateAdminPasswordHash(_ context.Context, hash string) error {
	if s.hash == "" {
		return contracts.ErrNotInitialized
	}
	s.hash = hash
	return nil
}

func (s *stubAdminCredentialStore) RequireLogin(_ context.Context) (bool, error) {
	return s.requireLogin, nil
}

func (s *stubAdminCredentialStore) SetRequireLogin(_ context.Context, required bool) error {
	s.requireLogin = required
	return nil
}

type stubProviderStore struct{}

var _ contracts.ProviderStore = stubProviderStore{}

func (stubProviderStore) ListEnabledConnections(_ context.Context) ([]contracts.ProviderConnection, error) {
	return []contracts.ProviderConnection{
		{
			ID:           "conn_1",
			ProviderID:   "openai",
			Name:         "Primary OpenAI",
			Protocol:     "openai-compatible",
			BaseURL:      "https://api.openai.com/v1",
			AuthType:     "api_key",
			Priority:     1,
			Enabled:      true,
			Models:       []string{"gpt-5"},
			ExtraHeaders: map[string]string{"X-Custom": "1"},
		},
	}, nil
}

func (stubProviderStore) ResolveCredentials(_ context.Context, connectionID string) (contracts.UpstreamCredentials, error) {
	if connectionID == "" {
		return contracts.UpstreamCredentials{}, contracts.ErrConnectionNotFound
	}
	return contracts.UpstreamCredentials{
		ConnectionID: connectionID,
		ProviderID:   "openai",
		Protocol:     "openai-compatible",
		BaseURL:      "https://api.openai.com/v1",
		AuthType:     "api_key",
		APIKey:       "sk-test",
		AccessToken:  "",
		ExtraHeaders: map[string]string{"X-Custom": "1"},
	}, nil
}

func (stubProviderStore) RefreshOAuthToken(_ context.Context, _ string) (contracts.UpstreamCredentials, error) {
	return contracts.UpstreamCredentials{}, contracts.ErrOAuthRefreshUnsupported
}

type stubAccountHealthTracker struct {
	cooling bool
}

var _ contracts.AccountHealthTracker = (*stubAccountHealthTracker)(nil)

func (s *stubAccountHealthTracker) IsInCooldown(_ string, _ time.Time) bool {
	return s.cooling
}

func (s *stubAccountHealthTracker) RecordFailure(_ string, _ int, _ time.Duration) {
	s.cooling = true
}

func (s *stubAccountHealthTracker) RecordSuccess(_ string) {
	s.cooling = false
}

type stubModelLister struct{}

var _ contracts.ModelLister = stubModelLister{}

func (stubModelLister) ListModels(_ context.Context, _ []string) ([]contracts.ModelDescriptor, error) {
	return []contracts.ModelDescriptor{
		{ID: "gpt-5", OwnedBy: "openai", Kind: "model"},
		{ID: "fast", OwnedBy: "router", Kind: "combo"},
		{ID: "latest", OwnedBy: "router", Kind: "alias"},
	}, nil
}

func TestContractInterfacesSatisfied(t *testing.T) {
	var (
		auth      contracts.AuthVerifier         = stubAuthVerifier{}
		keys      contracts.KeyValidator         = stubKeyValidator{}
		resolver  contracts.ProviderResolver     = stubProviderResolver{}
		hook      contracts.GatewayInterceptor   = stubInterceptor{}
		recorder  contracts.UsageRecorder        = &stubUsageRecorder{}
		proxy     contracts.ProxyDialerResolver  = stubProxyDialerResolver{}
		credStore contracts.AdminCredentialStore = &stubAdminCredentialStore{requireLogin: true}
		providers contracts.ProviderStore        = stubProviderStore{}
		health    contracts.AccountHealthTracker = &stubAccountHealthTracker{}
		models    contracts.ModelLister          = stubModelLister{}
	)

	ctx := context.Background()

	if !auth.RequireLogin(ctx) {
		t.Fatal("expected RequireLogin to return true")
	}
	if _, err := auth.AuthenticateAdmin(&http.Request{}); err != nil {
		t.Fatalf("AuthenticateAdmin error: %v", err)
	}

	info, err := keys.ValidateKey(ctx, "sk-test")
	if err != nil || info.ID != "key_1" {
		t.Fatalf("ValidateKey unexpected result: %+v, %v", info, err)
	}

	targets, err := resolver.Resolve(ctx, "gpt-5", info.AllowedModels)
	if err != nil || len(targets) != 1 {
		t.Fatalf("Resolve unexpected result: %+v, %v", targets, err)
	}

	req := &contracts.CanonicalRequest{RequestID: "req_1", Model: "gpt-5", Stream: true}
	meta := &contracts.ResponseMeta{RequestID: "req_1", StatusCode: http.StatusOK}
	ev := &contracts.StreamEvent{RequestID: "req_1", Event: "message", Data: []byte(`{}`)}
	if hook.Name() != "noop" || hook.Priority() != 100 {
		t.Fatalf("unexpected interceptor metadata: %s %d", hook.Name(), hook.Priority())
	}
	if err := hook.OnRequest(ctx, req); err != nil {
		t.Fatalf("OnRequest error: %v", err)
	}
	if err := hook.OnResponseHeaders(ctx, meta); err != nil {
		t.Fatalf("OnResponseHeaders error: %v", err)
	}
	if err := hook.OnStreamChunk(ctx, ev); err != nil {
		t.Fatalf("OnStreamChunk error: %v", err)
	}

	recorder.Record(ctx, contracts.UsageRecord{
		RequestID:    "req_1",
		KeyID:        "key_1",
		Model:        "gpt-5",
		ProviderID:   "openai",
		ConnectionID: "conn_1",
		InputTokens:  12,
		OutputTokens: 34,
		Latency:      150 * time.Millisecond,
		StatusCode:   http.StatusOK,
		Timestamp:    time.Now().UTC(),
	})

	u, err := proxy.ProxyURLForProvider(ctx, "openai")
	if err != nil || u == nil || u.Host != "127.0.0.1:8080" {
		t.Fatalf("ProxyURLForProvider unexpected result: %v, %v", u, err)
	}

	if _, err := credStore.AdminPasswordHash(ctx); !errors.Is(err, contracts.ErrNotInitialized) {
		t.Fatalf("AdminPasswordHash before init err = %v, want ErrNotInitialized", err)
	}
	if err := credStore.UpdateAdminPasswordHash(ctx, "$2a$10$hash"); !errors.Is(err, contracts.ErrNotInitialized) {
		t.Fatalf("UpdateAdminPasswordHash before init err = %v, want ErrNotInitialized", err)
	}
	if reqLogin, err := credStore.RequireLogin(ctx); err != nil || !reqLogin {
		t.Fatalf("RequireLogin = %v, %v, want true, nil", reqLogin, err)
	}
	if err := credStore.SetRequireLogin(ctx, false); err != nil {
		t.Fatalf("SetRequireLogin error: %v", err)
	}

	conns, err := providers.ListEnabledConnections(ctx)
	if err != nil || len(conns) != 1 || !conns[0].Enabled || conns[0].Priority != 1 || conns[0].Name != "Primary OpenAI" {
		t.Fatalf("ListEnabledConnections unexpected result: %+v, %v", conns, err)
	}
	creds, err := providers.ResolveCredentials(ctx, conns[0].ID)
	if err != nil || creds.APIKey != "sk-test" || creds.ConnectionID != "conn_1" {
		t.Fatalf("ResolveCredentials unexpected result: %+v, %v", creds, err)
	}
	if _, err := providers.RefreshOAuthToken(ctx, conns[0].ID); !errors.Is(err, contracts.ErrOAuthRefreshUnsupported) {
		t.Fatalf("RefreshOAuthToken err = %v, want ErrOAuthRefreshUnsupported", err)
	}

	now := time.Now().UTC()
	if health.IsInCooldown("conn_1", now) {
		t.Fatal("expected IsInCooldown to be false initially")
	}
	health.RecordFailure("conn_1", http.StatusTooManyRequests, 5*time.Second)
	if !health.IsInCooldown("conn_1", now) {
		t.Fatal("expected IsInCooldown to be true after RecordFailure")
	}
	health.RecordSuccess("conn_1")
	if health.IsInCooldown("conn_1", now) {
		t.Fatal("expected IsInCooldown to be false after RecordSuccess")
	}

	descs, err := models.ListModels(ctx, info.AllowedModels)
	if err != nil || len(descs) != 3 || descs[1].Kind != "combo" || descs[2].Kind != "alias" {
		t.Fatalf("ListModels unexpected result: %+v, %v", descs, err)
	}
}

func TestSentinelErrorsDistinct(t *testing.T) {
	sentinels := []struct {
		name string
		err  error
	}{
		{name: "ErrNotInitialized", err: contracts.ErrNotInitialized},
		{name: "ErrModelNotFound", err: contracts.ErrModelNotFound},
		{name: "ErrModelForbidden", err: contracts.ErrModelForbidden},
		{name: "ErrAllAccountsCoolingDown", err: contracts.ErrAllAccountsCoolingDown},
		{name: "ErrConnectionNotFound", err: contracts.ErrConnectionNotFound},
		{name: "ErrOAuthRefreshUnsupported", err: contracts.ErrOAuthRefreshUnsupported},
	}

	for i, a := range sentinels {
		if a.err == nil {
			t.Fatalf("sentinel %s is nil", a.name)
		}
		if a.err.Error() == "" {
			t.Fatalf("sentinel %s has empty error string", a.name)
		}
		for j, b := range sentinels {
			if i == j {
				if !errors.Is(a.err, b.err) {
					t.Fatalf("expected %s to match itself with errors.Is", a.name)
				}
				continue
			}
			if errors.Is(a.err, b.err) {
				t.Fatalf("expected %s and %s to be distinct sentinel errors", a.name, b.name)
			}
		}
	}
}
