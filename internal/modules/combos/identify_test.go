package combos_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/omeg4monx/router-gateway/internal/contracts"
)

func TestIdentifyModelEndpoint(t *testing.T) {
	upstream := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var reqBody struct {
			Model string `json:"model"`
		}
		_ = json.NewDecoder(r.Body).Decode(&reqBody)

		switch reqBody.Model {
		case "claude-opus-4-6":
			// Proxy claims claude-opus-4-6, but upstream JSON + reply reveal DeepSeek-V3.
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			_, _ = w.Write([]byte(`{
				"id": "chatcmpl-ds-123",
				"model": "deepseek-chat",
				"choices": [{"message": {"role": "assistant", "content": "I am DeepSeek-V3, developed by DeepSeek."}}]
			}`))
		case "gpt-5":
			// Proxy rewrites response.model to gpt-5, but model self-identifies as Qwen.
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			_, _ = w.Write([]byte(`{
				"id": "chatcmpl-qw-456",
				"model": "gpt-5",
				"choices": [{"message": {"role": "assistant", "content": "I am Qwen, a large language model created by Alibaba Cloud."}}]
			}`))
		case "gpt-4o":
			// Genuine OpenAI GPT-4o with versioned snapshot & vendor header.
			w.Header().Set("Content-Type", "application/json")
			w.Header().Set("openai-model", "gpt-4o-2024-08-06")
			w.WriteHeader(http.StatusOK)
			_, _ = w.Write([]byte(`{
				"id": "chatcmpl-oa-789",
				"model": "gpt-4o-2024-08-06",
				"system_fingerprint": "fp_test123",
				"choices": [{"message": {"role": "assistant", "content": "I am GPT-4o, developed by OpenAI."}}]
			}`))
		case "broken-model":
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusBadRequest)
			_, _ = w.Write([]byte(`{"error":{"message":"invalid key sk-identify-secret"}}`))
		default:
			w.WriteHeader(http.StatusNotFound)
		}
	}))
	defer upstream.Close()

	store := &catalogFakeProviderStore{
		connections: []contracts.ProviderConnection{
			{
				ID:         "conn_id_1",
				ProviderID: "custom-proxy",
				Name:       "Proxy",
				Protocol:   "openai-chat",
				BaseURL:    upstream.URL + "/v1",
				Priority:   1,
				Enabled:    true,
			},
		},
		creds: map[string]contracts.UpstreamCredentials{
			"conn_id_1": {
				ConnectionID: "conn_id_1",
				ProviderID:   "custom-proxy",
				Protocol:     "openai-chat",
				BaseURL:      upstream.URL + "/v1",
				AuthType:     "api_key",
				APIKey:       "sk-identify-secret",
			},
		},
	}

	_, router := setupCatalogTestModule(t, store, upstream.Client())

	type identifyResp struct {
		OK                bool              `json:"ok"`
		ClaimedModel      string            `json:"claimedModel"`
		ClaimedFamily     string            `json:"claimedFamily"`
		Verdict           string            `json:"verdict"`
		Confidence        string            `json:"confidence"`
		DetectedModel     string            `json:"detectedModel"`
		DetectedFamily    string            `json:"detectedFamily"`
		ReportedModel     string            `json:"reportedModel"`
		ResponseIDPrefix  string            `json:"responseIdPrefix"`
		SystemFingerprint string            `json:"systemFingerprint"`
		VendorHeaders     map[string]string `json:"vendorHeaders"`
		SelfReport        string            `json:"selfReport"`
		Signals           []string          `json:"signals"`
		LimitationNote    string            `json:"limitationNote"`
		StatusCode        int               `json:"statusCode"`
		Error             string            `json:"error"`
	}

	// 1. Metadata + self-report mismatch -> suspected_masked (high)
	req1 := httptest.NewRequest(http.MethodPost, "/api/models/identify", bytes.NewBufferString(`{"providerId":"custom-proxy","modelId":"claude-opus-4-6"}`))
	rec1 := httptest.NewRecorder()
	router.ServeHTTP(rec1, req1)
	if rec1.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d: %s", rec1.Code, rec1.Body.String())
	}
	var res1 identifyResp
	if err := json.Unmarshal(rec1.Body.Bytes(), &res1); err != nil {
		t.Fatalf("unmarshal res1: %v", err)
	}
	if !res1.OK || res1.Verdict != "suspected_masked" || res1.Confidence != "high" {
		t.Errorf("expected suspected_masked high, got %+v", res1)
	}
	if res1.DetectedModel != "deepseek-chat" || res1.DetectedFamily != "deepseek" {
		t.Errorf("expected detected deepseek-chat/deepseek, got %q / %q", res1.DetectedModel, res1.DetectedFamily)
	}
	if len(res1.Signals) == 0 || strings.TrimSpace(res1.LimitationNote) == "" {
		t.Errorf("expected non-empty signals and limitationNote, got %+v", res1)
	}

	// 2. Rewritten model field, but self-report exposes Qwen -> suspected_masked (medium)
	req2 := httptest.NewRequest(http.MethodPost, "/api/models/identify", bytes.NewBufferString(`{"providerId":"custom-proxy","modelId":"gpt-5"}`))
	rec2 := httptest.NewRecorder()
	router.ServeHTTP(rec2, req2)
	var res2 identifyResp
	_ = json.Unmarshal(rec2.Body.Bytes(), &res2)
	if res2.Verdict != "suspected_masked" || res2.DetectedFamily != "qwen" {
		t.Errorf("expected suspected_masked with detectedFamily=qwen, got %+v", res2)
	}

	// 3. Genuine OpenAI gpt-4o -> verified (high)
	req3 := httptest.NewRequest(http.MethodPost, "/api/models/identify", bytes.NewBufferString(`{"providerId":"custom-proxy","modelId":"gpt-4o"}`))
	rec3 := httptest.NewRecorder()
	router.ServeHTTP(rec3, req3)
	var res3 identifyResp
	_ = json.Unmarshal(rec3.Body.Bytes(), &res3)
	if res3.Verdict != "verified" || res3.Confidence != "high" || res3.SystemFingerprint != "fp_test123" {
		t.Errorf("expected verified high with fingerprint, got %+v", res3)
	}
	if res3.VendorHeaders["openai-model"] != "gpt-4o-2024-08-06" {
		t.Errorf("expected openai-model vendor header, got %+v", res3.VendorHeaders)
	}

	// 4. Error response must sanitize secrets and return verdict=error
	req4 := httptest.NewRequest(http.MethodPost, "/api/models/identify", bytes.NewBufferString(`{"providerId":"custom-proxy","modelId":"broken-model"}`))
	rec4 := httptest.NewRecorder()
	router.ServeHTTP(rec4, req4)
	if strings.Contains(rec4.Body.String(), "sk-identify-secret") {
		t.Fatalf("secret leaked in identify error response: %s", rec4.Body.String())
	}
	var res4 identifyResp
	_ = json.Unmarshal(rec4.Body.Bytes(), &res4)
	if res4.OK || res4.Verdict != "error" || strings.TrimSpace(res4.LimitationNote) == "" {
		t.Errorf("expected error verdict with limitationNote, got %+v", res4)
	}

	// 5. Validation error on empty fields -> 400
	req5 := httptest.NewRequest(http.MethodPost, "/api/models/identify", bytes.NewBufferString(`{"providerId":"","modelId":""}`))
	rec5 := httptest.NewRecorder()
	router.ServeHTTP(rec5, req5)
	if rec5.Code != http.StatusBadRequest {
		t.Errorf("expected 400 for empty payload, got %d", rec5.Code)
	}
}
