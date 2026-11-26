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
	for i, strategy := range combos.OmniRouteRoutingStrategies {
		payload := map[string]any{
			"name":        "combo-" + strategy,
			"description": "OmniRoute strategy " + strategy,
			"strategy":    strategy,
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
			t.Fatalf("strategy[%d]=%q create failed: HTTP %d: %s", i, strategy, rec.Code, rec.Body.String())
		}

		var created combos.Combo
		if err := json.Unmarshal(rec.Body.Bytes(), &created); err != nil {
			t.Fatalf("unmarshal created combo: %v", err)
		}
		if created.Strategy != strategy {
			t.Errorf("expected strategy %q, got %q", strategy, created.Strategy)
		}
		if len(created.Steps) != 2 || created.Steps[0].Weight != 70 || created.Steps[0].Label != "Primary" {
			t.Errorf("expected structured steps preserved for %q, got %+v", strategy, created.Steps)
		}
		if len(created.Targets) != 2 || created.Targets[0] != "openai/gpt-4o" {
			t.Errorf("expected targets synced with steps for %q, got %+v", strategy, created.Targets)
		}
		if created.Config["routerStrategy"] != "cost" {
			t.Errorf("expected runtime config persisted for %q, got %+v", strategy, created.Config)
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

	// 5. lkgp should pin to the most recently succeeded connection and clear on failure
	createCombo(`{"name":"lkgp-combo","strategy":"lkgp","targets":["openai/gpt-4o","gemini/gemini-2.5-pro","ollama/llama3.3"]}`)
	mod.RecordSuccess("conn_gemini")
	lkgpRes, err := mod.Resolve(ctx, "lkgp-combo", nil)
	if err != nil {
		t.Fatalf("resolve lkgp-combo: %v", err)
	}
	if lkgpRes[0].ConnectionID != "conn_gemini" {
		t.Errorf("expected lkgp to pick last-succeeded conn_gemini first, got %q", lkgpRes[0].ConnectionID)
	}
	mod.RecordFailure("conn_gemini", http.StatusTooManyRequests, 0)
	lkgpAfterFail, err := mod.Resolve(ctx, "lkgp-combo", nil)
	if err != nil {
		t.Fatalf("resolve lkgp-combo after failure: %v", err)
	}
	if lkgpAfterFail[0].ConnectionID == "conn_gemini" {
		t.Errorf("expected lkgp to stop pinning cooling-down conn_gemini, got %q", lkgpAfterFail[0].ConnectionID)
	}
	mod.ClearCooldown("conn_gemini")

	// 6. least-used should balance across candidates
	createCombo(`{"name":"lu-combo","strategy":"least-used","targets":["openai/gpt-4o","gemini/gemini-2.5-pro","ollama/llama3.3"]}`)
	seenLU := make(map[string]int)
	for range 6 {
		res, err := mod.Resolve(ctx, "lu-combo", nil)
		if err != nil {
			t.Fatalf("resolve lu-combo: %v", err)
		}
		seenLU[res[0].ConnectionID]++
	}
	if len(seenLU) < 3 {
		t.Errorf("expected least-used to distribute across all 3 connections, got %v", seenLU)
	}

	// 7. cache-optimized should return deterministic order
	createCombo(`{"name":"cache-combo","strategy":"cache-optimized","targets":["openai/gpt-4o","gemini/gemini-2.5-pro","ollama/llama3.3"]}`)
	c1, err := mod.Resolve(ctx, "cache-combo", nil)
	if err != nil {
		t.Fatalf("resolve cache-combo #1: %v", err)
	}
	c2, err := mod.Resolve(ctx, "cache-combo", nil)
	if err != nil {
		t.Fatalf("resolve cache-combo #2: %v", err)
	}
	if c1[0].ConnectionID != c2[0].ConnectionID {
		t.Errorf("expected cache-optimized to be deterministic, got %q vs %q", c1[0].ConnectionID, c2[0].ConnectionID)
	}

	// 8. p2c, random, strict-random, weighted, quota-weighted, headroom, reset-aware, reset-window, expiry-first all resolve cleanly
	for _, strategy := range []string{
		"priority",
		"fill-first",
		"p2c",
		"random",
		"strict-random",
		"weighted",
		"quota-weighted",
		"headroom",
		"reset-aware",
		"reset-window",
		"expiry-first",
		"context-relay",
		"fusion",
		"pipeline",
		"jev",
	} {
		comboName := "exec-" + strategy
		createCombo(`{"name":"` + comboName + `","strategy":"` + strategy + `","steps":[{"kind":"model","model":"openai/gpt-4o","weight":60},{"kind":"model","model":"gemini/gemini-2.5-pro","weight":40}]}`)
		res, err := mod.Resolve(ctx, comboName, nil)
		if err != nil {
			t.Fatalf("resolve %s: %v", comboName, err)
		}
		if len(res) != 2 {
			t.Errorf("expected 2 targets for %s, got %d", strategy, len(res))
		}
	}
}

func TestProviderStrategyAPIAndResolverIntegration(t *testing.T) {
	store := &catalogFakeProviderStore{
		connections: []contracts.ProviderConnection{
			{
				ID:         "conn_oai_1",
				ProviderID: "openai",
				Protocol:   "openai",
				Priority:   1,
				Enabled:    true,
				Models:     []string{"gpt-4o"},
			},
			{
				ID:         "conn_oai_2",
				ProviderID: "openai",
				Protocol:   "openai",
				Priority:   2,
				Enabled:    true,
				Models:     []string{"gpt-4o"},
			},
			{
				ID:         "conn_oai_3",
				ProviderID: "openai",
				Protocol:   "openai",
				Priority:   3,
				Enabled:    true,
				Models:     []string{"gpt-4o"},
			},
		},
	}

	mod, router := setupCatalogTestModule(t, store, http.DefaultClient)
	ctx := context.Background()

	// 1. Default provider strategy is immediately active ("fill-first") without any prior write.
	getReq := httptest.NewRequest(http.MethodGet, "/api/providers/openai/strategy", nil)
	getRec := httptest.NewRecorder()
	router.ServeHTTP(getRec, getReq)
	if getRec.Code != http.StatusOK {
		t.Fatalf("GET provider strategy expected 200, got %d: %s", getRec.Code, getRec.Body.String())
	}
	var initial combos.ProviderStrategyState
	if err := json.Unmarshal(getRec.Body.Bytes(), &initial); err != nil {
		t.Fatalf("unmarshal initial provider strategy: %v", err)
	}
	if initial.Strategy != combos.DefaultProviderStrategy || !initial.IsDefault || initial.DefaultStrategy != combos.DefaultProviderStrategy {
		t.Fatalf("expected active default %q, got %+v", combos.DefaultProviderStrategy, initial)
	}

	// Default fill-first preserves priority order across repeated resolves.
	for range 3 {
		targets, err := mod.Resolve(ctx, "openai/gpt-4o", nil)
		if err != nil {
			t.Fatalf("resolve default fill-first: %v", err)
		}
		if targets[0].ConnectionID != "conn_oai_1" {
			t.Fatalf("expected fill-first default to keep conn_oai_1 first, got %q", targets[0].ConnectionID)
		}
	}

	// 2. Reject unsupported / combo-only strategy on provider endpoint.
	badReq := httptest.NewRequest(http.MethodPut, "/api/providers/openai/strategy", bytes.NewBufferString(`{"strategy":"fusion"}`))
	badRec := httptest.NewRecorder()
	router.ServeHTTP(badRec, badReq)
	if badRec.Code != http.StatusBadRequest {
		t.Fatalf("expected 400 for combo-only strategy on provider, got %d: %s", badRec.Code, badRec.Body.String())
	}

	// 3. Update provider strategy to round-robin and verify live resolution rotates connections.
	putReq := httptest.NewRequest(http.MethodPut, "/api/providers/openai/strategy", bytes.NewBufferString(`{"strategy":"round-robin"}`))
	putRec := httptest.NewRecorder()
	router.ServeHTTP(putRec, putReq)
	if putRec.Code != http.StatusOK {
		t.Fatalf("PUT provider strategy expected 200, got %d: %s", putRec.Code, putRec.Body.String())
	}
	var updated combos.ProviderStrategyState
	if err := json.Unmarshal(putRec.Body.Bytes(), &updated); err != nil {
		t.Fatalf("unmarshal updated provider strategy: %v", err)
	}
	if updated.Strategy != "round-robin" || updated.IsDefault {
		t.Fatalf("expected persisted round-robin override, got %+v", updated)
	}

	r1, err := mod.Resolve(ctx, "openai/gpt-4o", nil)
	if err != nil {
		t.Fatalf("resolve provider rr #1: %v", err)
	}
	r2, err := mod.Resolve(ctx, "openai/gpt-4o", nil)
	if err != nil {
		t.Fatalf("resolve provider rr #2: %v", err)
	}
	r3, err := mod.Resolve(ctx, "openai/gpt-4o", nil)
	if err != nil {
		t.Fatalf("resolve provider rr #3: %v", err)
	}
	if r1[0].ConnectionID == r2[0].ConnectionID || r2[0].ConnectionID == r3[0].ConnectionID {
		t.Errorf("expected provider round-robin to rotate connections, got %q, %q, %q", r1[0].ConnectionID, r2[0].ConnectionID, r3[0].ConnectionID)
	}

	// 4. Bare model resolution also honors provider strategy (lkgp).
	putLKGP := httptest.NewRequest(http.MethodPut, "/api/providers/openai/strategy", bytes.NewBufferString(`{"strategy":"lkgp"}`))
	putLKGPRec := httptest.NewRecorder()
	router.ServeHTTP(putLKGPRec, putLKGP)
	if putLKGPRec.Code != http.StatusOK {
		t.Fatalf("PUT lkgp expected 200, got %d", putLKGPRec.Code)
	}
	mod.RecordSuccess("conn_oai_3")
	bareRes, err := mod.Resolve(ctx, "gpt-4o", nil)
	if err != nil {
		t.Fatalf("resolve bare model with provider lkgp: %v", err)
	}
	if bareRes[0].ConnectionID != "conn_oai_3" {
		t.Errorf("expected bare model to honor provider lkgp conn_oai_3, got %q", bareRes[0].ConnectionID)
	}

	// 5. List all configured provider strategies.
	listReq := httptest.NewRequest(http.MethodGet, "/api/providers/strategies", nil)
	listRec := httptest.NewRecorder()
	router.ServeHTTP(listRec, listReq)
	if listRec.Code != http.StatusOK {
		t.Fatalf("GET /api/providers/strategies expected 200, got %d: %s", listRec.Code, listRec.Body.String())
	}
	var listResp struct {
		DefaultStrategy string                         `json:"default_strategy"`
		Strategies      []combos.ProviderStrategyState `json:"strategies"`
	}
	if err := json.Unmarshal(listRec.Body.Bytes(), &listResp); err != nil {
		t.Fatalf("unmarshal provider strategies list: %v", err)
	}
	if listResp.DefaultStrategy != combos.DefaultProviderStrategy || len(listResp.Strategies) != 1 || listResp.Strategies[0].ProviderID != "openai" {
		t.Fatalf("unexpected strategies list: %+v", listResp)
	}
}

func BenchmarkResolveStrategies(b *testing.B) {
	store := &catalogFakeProviderStore{
		connections: []contracts.ProviderConnection{
			{ID: "c1", ProviderID: "openai", Protocol: "openai", Priority: 1, Enabled: true, Models: []string{"gpt-4o"}},
			{ID: "c2", ProviderID: "openai", Protocol: "openai", Priority: 2, Enabled: true, Models: []string{"gpt-4o"}},
			{ID: "c3", ProviderID: "openai", Protocol: "openai", Priority: 3, Enabled: true, Models: []string{"gpt-4o"}},
			{ID: "c4", ProviderID: "openai", Protocol: "openai", Priority: 4, Enabled: true, Models: []string{"gpt-4o"}},
		},
	}
	mod, router := setupCatalogTestModule(b, store, http.DefaultClient)
	ctx := context.Background()

	for _, strategy := range []string{"fill-first", "round-robin", "p2c", "least-used", "random", "lkgp"} {
		b.Run(strategy, func(b *testing.B) {
			req := httptest.NewRequest(http.MethodPut, "/api/providers/openai/strategy", bytes.NewBufferString(`{"strategy":"`+strategy+`"}`))
			rec := httptest.NewRecorder()
			router.ServeHTTP(rec, req)
			if rec.Code != http.StatusOK {
				b.Fatalf("set strategy %s failed: %d", strategy, rec.Code)
			}
			mod.RecordSuccess("c2")
			b.ResetTimer()
			for range b.N {
				if _, err := mod.Resolve(ctx, "openai/gpt-4o", nil); err != nil {
					b.Fatalf("resolve failed: %v", err)
				}
			}
		})
	}
}
