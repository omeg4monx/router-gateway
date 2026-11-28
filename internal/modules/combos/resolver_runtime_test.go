package combos_test

import (
	"bytes"
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/omeg4monx/router-gateway/internal/contracts"
)

func TestResolveAttachesComboRuntimeMetadata(t *testing.T) {
	store := &catalogFakeProviderStore{
		connections: []contracts.ProviderConnection{
			{
				ID:         "conn_1",
				ProviderID: "openai",
				Name:       "Primary",
				Protocol:   "openai-chat",
				BaseURL:    "https://api.openai.com/v1",
				Priority:   1,
				Enabled:    true,
				Models:     []string{"gpt-4o"},
			},
			{
				ID:         "conn_2",
				ProviderID: "anthropic",
				Name:       "Backup",
				Protocol:   "anthropic",
				BaseURL:    "https://api.anthropic.com/v1",
				Priority:   2,
				Enabled:    true,
				Models:     []string{"claude-sonnet-4"},
			},
		},
	}

	mod, router := setupCatalogTestModule(t, store, http.DefaultClient)

	// Child combo
	recChild := httptest.NewRecorder()
	router.ServeHTTP(recChild, httptest.NewRequest(http.MethodPost, "/api/combos", bytes.NewBufferString(`{
		"name":"child-combo",
		"strategy":"priority",
		"targets":["openai/gpt-4o"]
	}`)))
	if recChild.Code != http.StatusCreated {
		t.Fatalf("create child combo: %d %s", recChild.Code, recChild.Body.String())
	}

	// Parent combo with runtime config + quota-only nested combo-ref
	recParent := httptest.NewRecorder()
	router.ServeHTTP(recParent, httptest.NewRequest(http.MethodPost, "/api/combos", bytes.NewBufferString(`{
		"name":"parent-runtime",
		"strategy":"strict-random",
		"steps":[
			{"kind":"combo-ref","comboName":"child-combo","fallbackOnlyOnQuotaExhaustion":true},
			{"kind":"model","model":"anthropic/claude-sonnet-4"}
		],
		"config":{
			"maxRetries":3,
			"retryDelayMs":250,
			"timeoutMs":5000,
			"nestedComboMode":"execute",
			"disableSessionStickiness":true
		}
	}`)))
	if recParent.Code != http.StatusCreated {
		t.Fatalf("create parent combo: %d %s", recParent.Code, recParent.Body.String())
	}

	resolved, err := mod.Resolve(context.Background(), "parent-runtime", nil)
	if err != nil {
		t.Fatalf("Resolve failed: %v", err)
	}
	if len(resolved) != 2 {
		t.Fatalf("expected 2 resolved targets, got %d", len(resolved))
	}

	// In nestedComboMode=execute, child-combo stays first and keeps quota-only tag.
	first := resolved[0]
	if first.ConnectionID != "conn_1" {
		t.Fatalf("expected first target conn_1 in execute mode, got %s", first.ConnectionID)
	}
	if first.Metadata["combo_fallback_only_on_quota"] != "true" {
		t.Fatalf("expected combo_fallback_only_on_quota=true on nested ref, got %+v", first.Metadata)
	}
	if first.Metadata["combo_max_retries"] != "3" ||
		first.Metadata["combo_retry_delay_ms"] != "250" ||
		first.Metadata["combo_timeout_ms"] != "5000" ||
		first.Metadata["combo_disable_stickiness"] != "true" {
		t.Fatalf("unexpected runtime metadata on first target: %+v", first.Metadata)
	}
}
