package providers

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"
	"time"

	"github.com/omeg4monx/router-gateway/internal/contracts"
)

func TestJetBrainsAdapterAndOAuthRegistry(t *testing.T) {
	m := New(Options{})

	adapter, ok := m.GetAdapter(providerIDJetBrains)
	if !ok {
		t.Fatalf("expected jetbrains adapter to be registered")
	}
	if adapter.AuthType() != authTypeOAuth || adapter.Protocol() != protocolOpenAIChat {
		t.Fatalf("unexpected jetbrains adapter: auth=%s proto=%s", adapter.AuthType(), adapter.Protocol())
	}
	models := adapter.DefaultModels()
	if len(models) != 1 || models[0] != "gpt-6-luna" {
		t.Fatalf("jetbrains default models = %v, want [gpt-6-luna]", models)
	}

	strategy := m.getOAuthStrategy(providerIDJetBrains)
	if strategy.FlowType != oauthFlowAuthCodePKCE || strategy.FixedPort != 62345 || strategy.CallbackPath != "/" {
		t.Fatalf("unexpected jetbrains oauth strategy: %+v", strategy)
	}

	authURL, err := strategy.BuildAuthURL(*adapter.OAuthConfig(), "http://localhost:20128/api/providers/oauth/callback", "st-1", "chal-1")
	if err != nil {
		t.Fatalf("BuildAuthURL: %v", err)
	}
	parsed, err := url.Parse(authURL)
	if err != nil {
		t.Fatalf("parse authURL: %v", err)
	}
	if got := parsed.Query().Get("redirect_uri"); got != "http://localhost:62345" {
		t.Fatalf("redirect_uri = %q, want http://localhost:62345", got)
	}
	if got := parsed.Query().Get("client_id"); got != "junie-cli" {
		t.Fatalf("client_id = %q, want junie-cli", got)
	}
}

func TestJetBrainsPostExchangeAndRefresh(t *testing.T) {
	var licenseCalls int
	var provideAccessCalls int

	idToken := makeTestJWT(t, map[string]any{
		"email": "dev@jetbrains.example",
		"name":  "Jet Dev",
		"sub":   "jb-sub-42",
	})
	grazieToken := makeTestJWT(t, map[string]any{
		"ownerEmail": "dev@jetbrains.example",
		"exp":        time.Now().Add(24 * time.Hour).Unix(),
	})

	mockAI := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case "/auth/jetbrains-jwt/license/obtain/grazie-lite":
			licenseCalls++
			if r.Header.Get("Authorization") != "Bearer "+idToken {
				w.WriteHeader(http.StatusUnauthorized)
				return
			}
			w.Header().Set("Content-Type", "application/json")
			_ = json.NewEncoder(w).Encode(map[string]any{
				"license": map[string]any{"licenseId": "GZL-999"},
			})
		case "/auth/jetbrains-jwt/provide-access/license/v2":
			provideAccessCalls++
			w.Header().Set("Content-Type", "application/json")
			_ = json.NewEncoder(w).Encode(map[string]any{
				"token": grazieToken,
			})
		case "/oauth2/token":
			w.Header().Set("Content-Type", "application/json")
			_ = json.NewEncoder(w).Encode(map[string]any{
				"access_token":  "jb-refreshed-access",
				"refresh_token": "jb-refreshed-refresh",
				"id_token":      idToken,
				"expires_in":    3600,
			})
		default:
			w.WriteHeader(http.StatusNotFound)
		}
	}))
	defer mockAI.Close()

	prevBase := jetbrainsAIBaseURL
	jetbrainsAIBaseURL = mockAI.URL
	t.Cleanup(func() { jetbrainsAIBaseURL = prevBase })

	m := New(Options{HTTPClient: mockAI.Client()})

	postRes := postExchangeJetBrains(context.Background(), m, mockAI.URL+"/oauth2/token", oauthTokenResponse{
		AccessToken: "jb-access",
		IDToken:     idToken,
	})
	if licenseCalls != 1 || provideAccessCalls != 1 {
		t.Fatalf("calls license=%d provide=%d, want 1/1", licenseCalls, provideAccessCalls)
	}
	if postRes.Email != "dev@jetbrains.example" {
		t.Fatalf("email = %q, want dev@jetbrains.example", postRes.Email)
	}
	if postRes.DisplayLabel != "Jet Dev (Grazie Lite)" {
		t.Fatalf("displayLabel = %q, want Jet Dev (Grazie Lite)", postRes.DisplayLabel)
	}
	if postRes.ExtraHeaders["Grazie-Authenticate-JWT"] != grazieToken {
		t.Fatalf("Grazie-Authenticate-JWT header mismatch")
	}
	if postRes.Metadata["license_id"] != "GZL-999" || postRes.Metadata["grazie_jwt"] != grazieToken {
		t.Fatalf("unexpected metadata: %+v", postRes.Metadata)
	}

	// Verify refresh re-provisions grazie_jwt.
	refreshed, err := refreshJetBrainsToken(context.Background(), m, mockAI.URL+"/oauth2/token", "junie-cli", "", storedCredentials{
		RefreshToken: "old-refresh",
	})
	if err != nil {
		t.Fatalf("refreshJetBrainsToken: %v", err)
	}
	if refreshed.AccessToken != "jb-refreshed-access" {
		t.Fatalf("refreshed access = %q", refreshed.AccessToken)
	}
	if refreshed.Metadata["grazie_jwt"] != grazieToken || refreshed.Metadata["license_id"] != "GZL-999" {
		t.Fatalf("refreshed metadata = %+v", refreshed.Metadata)
	}
}

func TestJetBrainsQuotaFetcher(t *testing.T) {
	exp := time.Now().UTC().Add(12 * time.Hour).Truncate(time.Second)
	grazieToken := makeTestJWT(t, map[string]any{
		"exp": exp.Unix(),
	})

	t.Run("returns active grazie lite quota on 200", func(t *testing.T) {
		srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if r.Header.Get("Grazie-Authenticate-JWT") != grazieToken {
				w.WriteHeader(http.StatusUnauthorized)
				return
			}
			w.Header().Set("Content-Type", "application/json")
			_, _ = w.Write([]byte(`{"data":[{"id":"gpt-6-luna"}]}`))
		}))
		defer srv.Close()

		f := jetbrainsQuotaFetcher{}
		conn := ConnectionSummary{ID: "conn-jb-1", ProviderID: providerIDJetBrains, BaseURL: srv.URL}
		if !f.Matches(conn) {
			t.Fatalf("expected jetbrainsQuotaFetcher to match jetbrains connection")
		}
		res, err := f.Fetch(context.Background(), srv.Client(), conn, contracts.UpstreamCredentials{
			ConnectionID: "conn-jb-1",
			ProviderID:   providerIDJetBrains,
			BaseURL:      srv.URL,
			ExtraHeaders: map[string]string{"Grazie-Authenticate-JWT": grazieToken},
		})
		if err != nil {
			t.Fatalf("Fetch: %v", err)
		}
		if res.Status != "ok" || !res.Activated || !res.OnPlan {
			t.Fatalf("unexpected quota result: %+v", res)
		}
		if res.PlanName == nil || *res.PlanName != "JetBrains AI (Grazie Lite)" {
			t.Fatalf("planName = %v", res.PlanName)
		}
		if len(res.Buckets) != 1 || res.Buckets[0].ID != "gpt-6-luna" || res.Buckets[0].RemainingPct != 100 {
			t.Fatalf("buckets = %+v", res.Buckets)
		}
		if res.Buckets[0].ResetAt == nil || !strings.HasPrefix(*res.Buckets[0].ResetAt, exp.Format("2006-01-02T15:04:05")) {
			t.Fatalf("resetAt = %v, want %s", res.Buckets[0].ResetAt, exp.Format(time.RFC3339))
		}
	})

	t.Run("returns unauthorized on 401", func(t *testing.T) {
		srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
			w.WriteHeader(http.StatusUnauthorized)
		}))
		defer srv.Close()

		f := jetbrainsQuotaFetcher{}
		_, err := f.Fetch(context.Background(), srv.Client(), ConnectionSummary{
			ID: "conn-jb-2", ProviderID: providerIDJetBrains, BaseURL: srv.URL,
		}, contracts.UpstreamCredentials{
			ConnectionID: "conn-jb-2",
			ProviderID:   providerIDJetBrains,
			BaseURL:      srv.URL,
			ExtraHeaders: map[string]string{"Grazie-Authenticate-JWT": "expired-jwt"},
		})
		if !errors.Is(err, errQuotaUnauthorized) {
			t.Fatalf("err = %v, want errQuotaUnauthorized", err)
		}
	})
}

func TestJetBrainsBuildTestRequest(t *testing.T) {
	req, err := buildProtocolTestRequest(context.Background(), contracts.UpstreamCredentials{
		ConnectionID: "conn-jb-1",
		ProviderID:   providerIDJetBrains,
		Protocol:     protocolOpenAIChat,
		BaseURL:      "https://api.jetbrains.ai/user/v5/llm/openai/v1",
		ExtraHeaders: map[string]string{
			"Grazie-Authenticate-JWT": "jwt-xyz",
		},
	})
	if err != nil {
		t.Fatalf("buildProtocolTestRequest: %v", err)
	}
	if req.Method != http.MethodPost {
		t.Fatalf("method = %s, want POST", req.Method)
	}
	if req.URL.String() != "https://api.jetbrains.ai/user/v5/llm/openai/v1/chat/completions" {
		t.Fatalf("url = %s", req.URL.String())
	}
	if req.Header.Get("Grazie-Authenticate-JWT") != "jwt-xyz" {
		t.Fatalf("missing Grazie-Authenticate-JWT header")
	}
	if req.Header.Get("X-Forced-License-Free") != "2" {
		t.Fatalf("missing X-Forced-License-Free header")
	}
}

func TestJetBrainsConnectivityTestAccepts400And429(t *testing.T) {
	for _, tc := range []struct {
		name   string
		status int
		wantOK bool
	}{
		{name: "200 ok", status: http.StatusOK, wantOK: true},
		{name: "400 validation ok", status: http.StatusBadRequest, wantOK: true},
		{name: "429 rate limited ok", status: http.StatusTooManyRequests, wantOK: true},
		{name: "401 unauthorized fails", status: http.StatusUnauthorized, wantOK: false},
	} {
		t.Run(tc.name, func(t *testing.T) {
			srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
				w.WriteHeader(tc.status)
			}))
			defer srv.Close()

			m := New(Options{HTTPClient: srv.Client()})
			res, statusText, err := m.executeConnectivityTest(context.Background(), contracts.UpstreamCredentials{
				ConnectionID: "conn-jb-test",
				ProviderID:   providerIDJetBrains,
				Protocol:     protocolOpenAIChat,
				BaseURL:      srv.URL,
				ExtraHeaders: map[string]string{"Grazie-Authenticate-JWT": "jwt-xyz"},
			})
			if err != nil {
				t.Fatalf("executeConnectivityTest: %v", err)
			}
			if res.OK != tc.wantOK {
				t.Fatalf("status=%d -> OK=%v (statusText=%q), want %v", tc.status, res.OK, statusText, tc.wantOK)
			}
		})
	}
}
