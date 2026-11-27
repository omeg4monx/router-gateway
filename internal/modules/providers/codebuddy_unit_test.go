package providers

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/omeg4monx/router-gateway/internal/contracts"
)

func TestCodeBuddyAdaptersAndOAuthRegistry(t *testing.T) {
	m := New(Options{})
	adapters := make(map[string]ProviderAdapter)
	for _, a := range m.ListAdapters() {
		adapters[a.ID()] = a
	}

	for _, id := range []string{providerIDCodeBuddyCN, providerIDCodeBuddyIntl} {
		a, ok := adapters[id]
		if !ok {
			t.Fatalf("missing built-in adapter %q", id)
		}
		if a.AuthType() != authTypeOAuth {
			t.Errorf("%s AuthType = %q, want %q", id, a.AuthType(), authTypeOAuth)
		}
		if a.OAuthConfig() == nil {
			t.Errorf("%s OAuthConfig is nil", id)
		}
		strategy := m.getOAuthStrategy(id)
		if strategy == nil || strategy.FlowType != oauthFlowDeviceCode {
			t.Fatalf("%s strategy flowType = %+v, want %q", id, strategy, oauthFlowDeviceCode)
		}
		if !strategy.ClientIDOptional {
			t.Errorf("%s expected ClientIDOptional=true", id)
		}
		if strategy.RefreshToken == nil {
			t.Errorf("%s expected RefreshToken handler", id)
		}
	}
}

func TestCodeBuddyTokenRefresh(t *testing.T) {
	mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/v2/plugin/auth/token/refresh" {
			w.WriteHeader(http.StatusNotFound)
			return
		}
		if r.Header.Get("X-Refresh-Token") != "old-refresh" {
			t.Errorf("X-Refresh-Token = %q, want old-refresh", r.Header.Get("X-Refresh-Token"))
		}
		if r.Header.Get("X-Auth-Refresh-Source") != "plugin" {
			t.Errorf("X-Auth-Refresh-Source = %q, want plugin", r.Header.Get("X-Auth-Refresh-Source"))
		}
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(map[string]any{
			"code": 0,
			"data": map[string]any{
				"accessToken":  "new-access-tok",
				"refreshToken": "new-refresh-tok",
				"expiresIn":    43200,
			},
		})
	}))
	defer mockServer.Close()

	m := New(Options{HTTPClient: mockServer.Client()})
	tok, err := refreshCodeBuddyToken(
		context.Background(),
		m,
		providerIDCodeBuddyCN,
		mockServer.URL+"/v2/plugin/auth/token/refresh",
		storedCredentials{RefreshToken: "old-refresh"},
	)
	if err != nil {
		t.Fatalf("refreshCodeBuddyToken: %v", err)
	}
	if tok.AccessToken != "new-access-tok" || tok.RefreshToken != "new-refresh-tok" || tok.ExpiresIn != 43200 {
		t.Fatalf("unexpected refreshed token: %+v", tok)
	}
}

func TestCodeBuddyAndFreePIProtocolTestRequests(t *testing.T) {
	fpReq, err := buildProtocolTestRequest(context.Background(), contracts.UpstreamCredentials{
		ConnectionID: "c-fp",
		ProviderID:   "freepi",
		Protocol:     protocolOpenAIChat,
		AuthType:     authTypeOAuth,
		BaseURL:      "https://api.freepi.ai/v1",
		AccessToken:  "fp-tok",
	})
	if err != nil {
		t.Fatalf("freepi buildProtocolTestRequest: %v", err)
	}
	if fpReq.Method != http.MethodPost || fpReq.URL.String() != "https://api.freepi.ai/v1/chat/completions" {
		t.Fatalf("freepi test req = %s %s, want POST .../chat/completions", fpReq.Method, fpReq.URL)
	}
	if fpReq.Header.Get("x-client-version") == "" {
		t.Fatalf("freepi test req missing x-client-version")
	}

	cbReq, err := buildProtocolTestRequest(context.Background(), contracts.UpstreamCredentials{
		ConnectionID: "c-cb",
		ProviderID:   "codebuddy-cn",
		Protocol:     protocolOpenAIChat,
		AuthType:     authTypeOAuth,
		BaseURL:      "https://copilot.tencent.com/v2",
		AccessToken:  "cb-tok",
	})
	if err != nil {
		t.Fatalf("codebuddy buildProtocolTestRequest: %v", err)
	}
	if cbReq.Method != http.MethodPost || cbReq.URL.String() != "https://copilot.tencent.com/v2/billing/meter/get-user-resource" {
		t.Fatalf("codebuddy test req = %s %s", cbReq.Method, cbReq.URL)
	}
	if cbReq.Header.Get("x-codebuddy-request") != "1" {
		t.Fatalf("codebuddy test req missing x-codebuddy-request header")
	}
}

func TestCodeBuddyQuotaFetcherRefillAndBonusPacks(t *testing.T) {
	now := time.Now().UTC().Truncate(time.Second)
	cycleStart := now.Add(-5 * 24 * time.Hour)
	cycleEnd := now.Add(25 * 24 * time.Hour)
	deductionEndFar := now.Add(365 * 24 * time.Hour)
	bonusEnd := now.Add(10 * 24 * time.Hour)

	mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/v2/billing/meter/get-user-resource" {
			w.WriteHeader(http.StatusNotFound)
			return
		}
		if r.Method != http.MethodPost {
			t.Errorf("quota method = %s, want POST", r.Method)
		}
		if r.Header.Get("Authorization") != "Bearer cb-secret-tok" {
			t.Errorf("Authorization = %q", r.Header.Get("Authorization"))
		}
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(map[string]any{
			"code": 0,
			"data": map[string]any{
				"Response": map[string]any{
					"Data": map[string]any{
						"Accounts": []map[string]any{
							{
								"PackageName":              "CodeBuddy Pro",
								"CycleStartTime":           cycleStart.UnixMilli(),
								"CycleEndTime":             cycleEnd.UnixMilli(),
								"DeductionEndTime":         deductionEndFar.UnixMilli(),
								"CycleCapacityUsedPrecise": "50.0",
								"CycleCapacitySizePrecise": "500.0",
							},
							{
								"PackageName":         "Promo Bonus",
								"CycleStartTime":      cycleStart.UnixMilli(),
								"CycleEndTime":        bonusEnd.UnixMilli(),
								"DeductionEndTime":    bonusEnd.UnixMilli(),
								"CapacityUsedPrecise": "25",
								"CapacitySizePrecise": "100",
							},
						},
					},
				},
			},
		})
	}))
	defer mockServer.Close()

	fetcher := codebuddyQuotaFetcher{providerID: providerIDCodeBuddyCN}
	res, err := fetcher.Fetch(
		context.Background(),
		mockServer.Client(),
		ConnectionSummary{
			ID:         "conn-cb-1",
			ProviderID: providerIDCodeBuddyCN,
			BaseURL:    mockServer.URL + "/v2",
		},
		contracts.UpstreamCredentials{
			ProviderID:  providerIDCodeBuddyCN,
			BaseURL:     mockServer.URL + "/v2",
			AccessToken: "cb-secret-tok",
		},
	)
	if err != nil {
		t.Fatalf("Fetch: %v", err)
	}
	if res.PlanName == nil || *res.PlanName != "CodeBuddy Pro" {
		t.Fatalf("PlanName = %v, want CodeBuddy Pro", res.PlanName)
	}
	if !res.OnPlan || !res.Activated || res.Status != "ok" {
		t.Fatalf("unexpected flags: onPlan=%v activated=%v status=%q", res.OnPlan, res.Activated, res.Status)
	}
	if len(res.Buckets) != 2 {
		t.Fatalf("len(Buckets) = %d, want 2: %+v", len(res.Buckets), res.Buckets)
	}
	if res.Buckets[0].Label != "Monthly" || res.Buckets[0].UsedPct != 10 || res.Buckets[0].RemainingPct != 90 {
		t.Fatalf("bucket[0] = %+v, want Monthly 10%%/90%%", res.Buckets[0])
	}
	if res.Buckets[1].Label != "Bonus Pack 1" || res.Buckets[1].UsedPct != 25 || res.Buckets[1].RemainingPct != 75 {
		t.Fatalf("bucket[1] = %+v, want Bonus Pack 1 25%%/75%%", res.Buckets[1])
	}
}
