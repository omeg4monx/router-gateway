package combos_test

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/omeg4monx/router-gateway/internal/contracts"
	"github.com/omeg4monx/router-gateway/internal/modules/combos"
)

func TestAllOmniRouteComboStrategiesCRUDAndPersistence(t *testing.T) {
	store := &catalogFakeProviderStore{
		connections: []contracts.ProviderConnection{
			{
				ID:         "conn_openai",
				ProviderID: "openai",
				Protocol:   "openai",
				Priority:   10,
				Enabled:    true,
				Models:     []string{"gpt-4o"},
			},
			{
				ID:         "conn_ollama",
				ProviderID: "ollama",
				Protocol:   "openai",
				Priority:   20,
				Enabled:    true,
				Models:     []string{"llama3.3"},
			},
		},
	}

	_, router := setupCatalogTestModule(t, store, http.DefaultClient)

	// Verify all 22 OmniRoute strategies can be created and retrieved cleanly.
	for i, strat := range combos.OmniRouteRoutingStrategies {
		payload := map[string]any{
			"name":        "combo-" + strat,
			"description": "OmniRoute strategy " + strat,
			"strategy":    strat,
			"steps": []map[string]any{
				{"kind": "model", "model": "openai/gpt-4o", "weight": 70, "label": "Primary"},
				{"kind": "model", "model": "ollama/llama3.3", "weight": 30, "label": "Local"},
			},
			"config": map[string]any{
				"maxRetries":     2,
				"routerStrategy": "cost",
			},
		}
		raw, _ := json.Marshal(payload)
		req := httptest.NewRequest(http.MethodPost, "/api/combos", bytes.NewReader(raw))
		rec := httptest.NewRecorder()
		router.ServeHTTP(rec, req)
		if rec.Code != http.StatusCreated {
			t.Fatalf("strategy[%d]=%q create failed: HTTP %d: %s", i, strat, rec.Code, rec.Body.String())
		}

		var created combos.Combo
		if err := json.Unmarshal(rec.Body.Bytes(), &created); err != nil {
			t.Fatalf("unmarshal created combo: %v", err)
		}
		if created.Strategy != strat {
			t.Errorf("expected strategy %q, got %q", strat, created.Strategy)
		}
		if len(created.Steps) != 2 || created.Steps[0].Weight != 70 || created.Steps[0].Label != "Primary" {
			t.Errorf("expected structured steps preserved for %q, got %+v", strat, created.Steps)
		}
		if len(created.Targets) != 2 || created.Targets[0] != "openai/gpt-4o" {
			t.Errorf("expected targets synced with steps for %q, got %+v", strat, created.Targets)
		}
		if created.Config["routerStrategy"] != "cost" {
			t.Errorf("expected runtime config persisted for %q, got %+v", strat, created.Config)
		}
	}

	// Verify combo-ref step nesting works
	nestPayload := `{"name":"nested-combo","strategy":"priority","steps":[{"kind":"combo-ref","comboName":"combo-priority","fallbackOnlyOnQuotaExhaustion":true},{"kind":"model","model":"ollama/llama3.3"}],"config":{"nestedComboMode":"execute"}}`
	nestReq := httptest.NewRequest(http.MethodPost, "/api/combos", bytes.NewBufferString(nestPayload))
	nestRec := httptest.NewRecorder()
	router.ServeHTTP(nestRec, nestReq)
	if nestRec.Code != http.StatusCreated {
		t.Fatalf("expected 201 for combo-ref step, got %d: %s", nestRec.Code, nestRec.Body.String())
	}
}

func TestOmniRouteResolverStrategiesExecution(t *testing.T) {
	store := &catalogFakeProviderStore{
		connections: []contracts.ProviderConnection{
			{
				ID:         "conn_openai",
				ProviderID: "openai",
				Protocol:   "openai",
				Priority:   10,
				Enabled:    true,
				Models:     []string{"gpt-4o"}, // 128K context, paid
			},
			{
				ID:         "conn_gemini",
				ProviderID: "gemini",
				Protocol:   "gemini",
				Priority:   20,
				Enabled:    true,
				Models:     []string{"gemini-2.5-pro"}, // 1M context, paid
			},
			{
				ID:         "conn_ollama",
				ProviderID: "ollama",
				Protocol:   "openai",
				Priority:   30,
				Enabled:    true,
				Models:     []string{"llama3.3"}, // free
			},
		},
	}

	mod, router := setupCatalogTestModule(t, store, http.DefaultClient)
	ctx := context.Background()

	createCombo := func(body string) {
		t.Helper()
		req := httptest.NewRequest(http.MethodPost, "/api/combos", bytes.NewBufferString(body))
		rec := httptest.NewRecorder()
		router.ServeHTTP(rec, req)
		if rec.Code != http.StatusCreated {
			t.Fatalf("create combo failed: %d %s", rec.Code, rec.Body.String())
		}
	}

	// 1. Round-robin should rotate primary target across consecutive Resolve() calls
	createCombo(`{"name":"rr-combo","strategy":"round-robin","targets":["openai/gpt-4o","gemini/gemini-2.5-pro","ollama/llama3.3"]}`)
	r1, err := mod.Resolve(ctx, "rr-combo", nil)
	if err != nil {
		t.Fatalf("resolve rr #1: %v", err)
	}
	r2, err := mod.Resolve(ctx, "rr-combo", nil)
	if err != nil {
		t.Fatalf("resolve rr #2: %v", err)
	}
	r3, err := mod.Resolve(ctx, "rr-combo", nil)
	if err != nil {
		t.Fatalf("resolve rr #3: %v", err)
	}
	if r1[0].ConnectionID == r2[0].ConnectionID || r2[0].ConnectionID == r3[0].ConnectionID {
		t.Errorf("expected round-robin to rotate primaries, got %q, %q, %q", r1[0].ConnectionID, r2[0].ConnectionID, r3[0].ConnectionID)
	}

	// 2. cost-optimized should put free provider (ollama) first
	createCombo(`{"name":"cost-combo","strategy":"cost-optimized","targets":["openai/gpt-4o","ollama/llama3.3"]}`)
	costRes, err := mod.Resolve(ctx, "cost-combo", nil)
	if err != nil {
		t.Fatalf("resolve cost-combo: %v", err)
	}
	if costRes[0].ProviderID != "ollama" {
		t.Errorf("expected cost-optimized to prefer free ollama first, got %q", costRes[0].ProviderID)
	}

	// 3. context-optimized should put largest context window (gemini-2.5-pro = 1M) first
	createCombo(`{"name":"ctx-combo","strategy":"context-optimized","targets":["openai/gpt-4o","gemini/gemini-2.5-pro"]}`)
	ctxRes, err := mod.Resolve(ctx, "ctx-combo", nil)
	if err != nil {
		t.Fatalf("resolve ctx-combo: %v", err)
	}
	if ctxRes[0].ProviderID != "gemini" {
		t.Errorf("expected context-optimized to prefer 1M-context gemini first, got %q", ctxRes[0].ProviderID)
	}

	// 4. auto with routerStrategy=cost should also prefer free ollama first
	createCombo(`{"name":"auto-cost","strategy":"auto","targets":["openai/gpt-4o","ollama/llama3.3"],"config":{"routerStrategy":"cost"}}`)
	autoRes, err := mod.Resolve(ctx, "auto-cost", nil)
	if err != nil {
		t.Fatalf("resolve auto-cost: %v", err)
	}
	if autoRes[0].ProviderID != "ollama" {
		t.Errorf("expected auto(cost) to prefer ollama first, got %q", autoRes[0].ProviderID)
	}
}
