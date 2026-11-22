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

func TestContractInterfacesSatisfied(t *testing.T) {
	var (
		auth      contracts.AuthVerifier         = stubAuthVerifier{}
		keys      contracts.KeyValidator         = stubKeyValidator{}
		resolver  contracts.ProviderResolver     = stubProviderResolver{}
		hook      contracts.GatewayInterceptor   = stubInterceptor{}
		recorder  contracts.UsageRecorder        = &stubUsageRecorder{}
		proxy     contracts.ProxyDialerResolver  = stubProxyDialerResolver{}
		credStore contracts.AdminCredentialStore = &stubAdminCredentialStore{requireLogin: true}
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
}
