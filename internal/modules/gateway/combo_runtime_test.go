package gateway_test

import (
	"bytes"
	"context"
	"net/http"
	"net/http/httptest"
	"sync/atomic"
	"testing"
	"time"

	"github.com/omeg4monx/router-gateway/internal/contracts"
	"github.com/omeg4monx/router-gateway/internal/modules/gateway"
)

func TestGatewayComboRuntimeRetriesTimeoutQuotaFallbackAndStickiness(t *testing.T) {
	// 1. Verify combo_max_retries + combo_retry_delay_ms retries transient 503 on a multi-candidate combo.
	var primaryAttempts atomic.Int32
	var backupAttempts atomic.Int32

	primarySrv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		if primaryAttempts.Add(1) < 3 {
			w.WriteHeader(http.StatusServiceUnavailable)
			_, _ = w.Write([]byte(`{"error":{"message":"temporary overload"}}`))
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{"id":"chatcmpl-retry-ok","object":"chat.completion","choices":[{"index":0,"message":{"role":"assistant","content":"recovered"},"finish_reason":"stop"}]}`))
	}))
	defer primarySrv.Close()

	backupSrv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		backupAttempts.Add(1)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{"id":"chatcmpl-backup","object":"chat.completion","choices":[{"index":0,"message":{"role":"assistant","content":"backup"},"finish_reason":"stop"}]}`))
	}))
	defer backupSrv.Close()

	env := setupGatewayTest(t, nil, 0, false)
	env.providers.creds["conn_primary"] = contracts.UpstreamCredentials{
		ConnectionID: "conn_primary", ProviderID: "openai", Protocol: gateway.ProtocolOpenAIChat,
		BaseURL: primarySrv.URL + "/v1", AuthType: "api_key", APIKey: "sk-1",
	}
	env.providers.creds["conn_backup"] = contracts.UpstreamCredentials{
		ConnectionID: "conn_backup", ProviderID: "openai", Protocol: gateway.ProtocolOpenAIChat,
		BaseURL: backupSrv.URL + "/v1", AuthType: "api_key", APIKey: "sk-2",
	}
	env.resolver.fn = func(context.Context, string, []string) ([]contracts.ResolvedTarget, error) {
		return []contracts.ResolvedTarget{
			{
				ProviderID:    "openai",
				ProviderType:  gateway.ProtocolOpenAIChat,
				ConnectionID:  "conn_primary",
				Model:         "retry-combo",
				UpstreamModel: "gpt-4o",
				Endpoint:      primarySrv.URL + "/v1",
				Metadata: map[string]string{
					"combo_max_retries":    "2",
					"combo_retry_delay_ms": "5",
				},
			},
			{
				ProviderID:    "openai",
				ProviderType:  gateway.ProtocolOpenAIChat,
				ConnectionID:  "conn_backup",
				Model:         "retry-combo",
				UpstreamModel: "gpt-4o-mini",
				Endpoint:      backupSrv.URL + "/v1",
			},
		}, nil
	}

	req := httptest.NewRequest(http.MethodPost, "/v1/chat/completions", bytes.NewBufferString(`{"model":"retry-combo","messages":[{"role":"user","content":"hi"}]}`))
	req.Header.Set("Authorization", "Bearer sk-valid")
	rec := httptest.NewRecorder()
	env.server.Handler().ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected 200 after 2 retries, got %d: %s", rec.Code, rec.Body.String())
	}
	if primaryAttempts.Load() != 3 {
		t.Fatalf("expected 3 primary attempts (1 + 2 retries), got %d", primaryAttempts.Load())
	}
	if backupAttempts.Load() != 0 {
		t.Fatalf("expected backup not to be called when primary recovered on retry, got %d", backupAttempts.Load())
	}

	// 2. Verify combo_fallback_only_on_quota stops cascade on 500 (non-429).
	quotaOnlyPrimary := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = w.Write([]byte(`{"error":{"message":"internal crash"}}`))
	}))
	defer quotaOnlyPrimary.Close()

	backupAttempts.Store(0)
	env.providers.creds["conn_q1"] = contracts.UpstreamCredentials{
		ConnectionID: "conn_q1", ProviderID: "openai", Protocol: gateway.ProtocolOpenAIChat,
		BaseURL: quotaOnlyPrimary.URL + "/v1", AuthType: "api_key", APIKey: "sk-1",
	}
	env.resolver.fn = func(context.Context, string, []string) ([]contracts.ResolvedTarget, error) {
		return []contracts.ResolvedTarget{
			{
				ProviderID:    "openai",
				ProviderType:  gateway.ProtocolOpenAIChat,
				ConnectionID:  "conn_q1",
				Model:         "quota-only-combo",
				UpstreamModel: "gpt-4o",
				Endpoint:      quotaOnlyPrimary.URL + "/v1",
				Metadata: map[string]string{
					"combo_max_retries":            "0",
					"combo_fallback_only_on_quota": "true",
				},
			},
			{
				ProviderID:    "openai",
				ProviderType:  gateway.ProtocolOpenAIChat,
				ConnectionID:  "conn_backup",
				Model:         "quota-only-combo",
				UpstreamModel: "gpt-4o-mini",
				Endpoint:      backupSrv.URL + "/v1",
			},
		}, nil
	}

	req2 := httptest.NewRequest(http.MethodPost, "/v1/chat/completions", bytes.NewBufferString(`{"model":"quota-only-combo","messages":[{"role":"user","content":"hi"}]}`))
	req2.Header.Set("Authorization", "Bearer sk-valid")
	rec2 := httptest.NewRecorder()
	env.server.Handler().ServeHTTP(rec2, req2)

	if rec2.Code != http.StatusInternalServerError {
		t.Fatalf("expected 500 when fallback_only_on_quota blocks non-429 cascade, got %d: %s", rec2.Code, rec2.Body.String())
	}
	if backupAttempts.Load() != 0 {
		t.Fatalf("expected backup not to be called on 500 when combo_fallback_only_on_quota=true, got %d", backupAttempts.Load())
	}

	// 3. Verify combo_timeout_ms triggers failover to backup when primary stalls.
	slowSrv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		select {
		case <-time.After(3 * time.Second):
			w.WriteHeader(http.StatusOK)
		case <-r.Context().Done():
		}
	}))
	defer slowSrv.Close()

	backupAttempts.Store(0)
	env.providers.creds["conn_slow"] = contracts.UpstreamCredentials{
		ConnectionID: "conn_slow", ProviderID: "openai", Protocol: gateway.ProtocolOpenAIChat,
		BaseURL: slowSrv.URL + "/v1", AuthType: "api_key", APIKey: "sk-1",
	}
	env.resolver.fn = func(context.Context, string, []string) ([]contracts.ResolvedTarget, error) {
		return []contracts.ResolvedTarget{
			{
				ProviderID:    "openai",
				ProviderType:  gateway.ProtocolOpenAIChat,
				ConnectionID:  "conn_slow",
				Model:         "timeout-combo",
				UpstreamModel: "gpt-4o",
				Endpoint:      slowSrv.URL + "/v1",
				Metadata: map[string]string{
					"combo_max_retries": "0",
					"combo_timeout_ms":  "1000",
				},
			},
			{
				ProviderID:    "openai",
				ProviderType:  gateway.ProtocolOpenAIChat,
				ConnectionID:  "conn_backup",
				Model:         "timeout-combo",
				UpstreamModel: "gpt-4o-mini",
				Endpoint:      backupSrv.URL + "/v1",
			},
		}, nil
	}

	req3 := httptest.NewRequest(http.MethodPost, "/v1/chat/completions", bytes.NewBufferString(`{"model":"timeout-combo","messages":[{"role":"user","content":"hi"}]}`))
	req3.Header.Set("Authorization", "Bearer sk-valid")
	rec3 := httptest.NewRecorder()
	env.server.Handler().ServeHTTP(rec3, req3)

	if rec3.Code != http.StatusOK {
		t.Fatalf("expected 200 via backup after primary combo_timeout_ms elapsed, got %d: %s", rec3.Code, rec3.Body.String())
	}
	if backupAttempts.Load() != 1 {
		t.Fatalf("expected backup to handle failover after timeout, got %d", backupAttempts.Load())
	}
}
