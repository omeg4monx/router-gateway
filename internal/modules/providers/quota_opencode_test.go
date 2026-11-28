package providers

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/omeg4monx/router-gateway/internal/contracts"
)

func TestOpenCodeQuotaFetcherParsesUsageAndHandlesAuthErrors(t *testing.T) {
	t.Parallel()

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/zen/v1/usage" {
			w.WriteHeader(http.StatusNotFound)
			return
		}
		switch r.Header.Get("Authorization") {
		case "Bearer sk-good":
			w.Header().Set("Content-Type", "application/json")
			_, _ = w.Write([]byte(`{
				"usage": {
					"rolling": {"percent": 25, "resetsAt": "2026-12-01T00:00:00Z"},
					"weekly": {"percent": "60.5", "resetsAt": 1800000000},
					"monthly": {"percent": 100}
				}
			}`))
		default:
			w.WriteHeader(http.StatusUnauthorized)
		}
	}))
	defer srv.Close()

	fetcher := opencodeQuotaFetcher{
		providerID: "opencode-zen",
		planName:   "OpenCode Zen",
		usageURL:   srv.URL + "/zen/v1/usage",
	}

	conn := ConnectionSummary{
		ID:         "conn-zen-1",
		ProviderID: "opencode-zen",
		BaseURL:    srv.URL + "/zen/v1",
	}
	if !fetcher.Matches(conn) {
		t.Fatalf("expected fetcher to match opencode-zen connection")
	}

	res, err := fetcher.Fetch(context.Background(), srv.Client(), conn, contracts.UpstreamCredentials{
		ProviderID: "opencode-zen",
		BaseURL:    srv.URL + "/zen/v1",
		APIKey:     "sk-good",
	})
	if err != nil {
		t.Fatalf("Fetch good key error: %v", err)
	}
	if res.Status != "ok" || len(res.Buckets) != 3 {
		t.Fatalf("unexpected quota result: %+v", res)
	}
	if res.Buckets[0].ID != "rolling" || res.Buckets[0].UsedPct != 25 || res.Buckets[0].RemainingPct != 75 {
		t.Fatalf("unexpected rolling bucket: %+v", res.Buckets[0])
	}
	if res.Buckets[2].ID != "monthly" || !res.Buckets[2].Exhausted {
		t.Fatalf("expected monthly bucket exhausted: %+v", res.Buckets[2])
	}

	_, err = fetcher.Fetch(context.Background(), srv.Client(), conn, contracts.UpstreamCredentials{
		ProviderID: "opencode-zen",
		BaseURL:    srv.URL + "/zen/v1",
		APIKey:     "sk-bad",
	})
	if !errors.Is(err, errQuotaUnauthorized) {
		t.Fatalf("expected errQuotaUnauthorized, got %v", err)
	}
}

func TestOpenCodeDiscoveryURLsAndFreeTestRequest(t *testing.T) {
	t.Parallel()

	freeURLs := buildDiscoveryCandidateURLs("opencode", protocolOpenAIChat, "https://opencode.ai/inference/openai/v1")
	if len(freeURLs) == 0 || freeURLs[0] != "https://opencode.ai/inference/v1/models" {
		t.Fatalf("opencode free discovery URLs = %v, want first https://opencode.ai/inference/v1/models", freeURLs)
	}

	goURLs := buildDiscoveryCandidateURLs("opencode-go", protocolOpenAIChat, "https://opencode.ai/zen/go/v1")
	if len(goURLs) == 0 || goURLs[0] != "https://opencode.ai/inference/go/v1/models" {
		t.Fatalf("opencode-go discovery URLs = %v, want first https://opencode.ai/inference/go/v1/models", goURLs)
	}

	req, err := buildProtocolTestRequest(context.Background(), contracts.UpstreamCredentials{
		ProviderID: "opencode",
		Protocol:   protocolOpenAIChat,
		BaseURL:    "https://opencode.ai/inference/openai/v1",
		AuthType:   "none",
	})
	if err != nil {
		t.Fatalf("buildProtocolTestRequest(opencode) error: %v", err)
	}
	if req.URL.String() != "https://opencode.ai/inference/v1/models" {
		t.Fatalf("test req URL = %q, want https://opencode.ai/inference/v1/models", req.URL.String())
	}
	if req.Header.Get("Authorization") != "Bearer public" || req.Header.Get("x-opencode-client") != "cli" {
		t.Fatalf("unexpected opencode test headers: %v", req.Header)
	}
}
