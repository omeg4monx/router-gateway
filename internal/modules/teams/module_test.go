package teams_test

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"path/filepath"
	"testing"
	"time"

	"github.com/go-chi/chi/v5"

	"github.com/omeg4monx/router-gateway/internal/contracts"
	"github.com/omeg4monx/router-gateway/internal/modules/teams"
	"github.com/omeg4monx/router-gateway/internal/platform/db"
	"github.com/omeg4monx/router-gateway/internal/platform/eventbus"
	"github.com/omeg4monx/router-gateway/internal/platform/lifecycle"
)

func setupTeamsTest(t *testing.T) (*teams.Module, chi.Router, *eventbus.Bus) {
	t.Helper()

	database, err := db.Open(filepath.Join(t.TempDir(), "teams.sqlite"))
	if err != nil {
		t.Fatalf("db.Open: %v", err)
	}
	t.Cleanup(func() { _ = database.Close() })

	bus := eventbus.New()
	mod := teams.New(teams.Options{
		DB:  database,
		Bus: bus,
	})
	if err := database.MigrateModule(context.Background(), mod.Name(), mod.Migrations()); err != nil {
		t.Fatalf("MigrateModule: %v", err)
	}

	r := chi.NewRouter()
	r.Route("/api", func(api chi.Router) {
		mod.RegisterRoutes(lifecycle.RouteMounts{AdminAPI: api})
	})

	if err := mod.Start(context.Background()); err != nil {
		t.Fatalf("mod.Start: %v", err)
	}
	t.Cleanup(func() { _ = mod.Stop(context.Background()) })

	return mod, r, bus
}

func TestTeamsCRUDMembershipAndSpendReset(t *testing.T) {
	t.Parallel()

	mod, router, _ := setupTeamsTest(t)

	// Create team.
	createBody, _ := json.Marshal(map[string]any{
		"name":           "Engineering",
		"description":    "Core engineering keys",
		"budget_usd":     25.5,
		"rate_limit_rpm": 60,
	})
	req := httptest.NewRequest(http.MethodPost, "/api/teams", bytes.NewReader(createBody))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)
	if rec.Code != http.StatusCreated {
		t.Fatalf("POST /api/teams status = %d, want 201; body=%s", rec.Code, rec.Body.String())
	}

	var created struct {
		ID      string   `json:"id"`
		Name    string   `json:"name"`
		KeyIDs  []string `json:"key_ids"`
		Enabled bool     `json:"enabled"`
	}
	if err := json.Unmarshal(rec.Body.Bytes(), &created); err != nil {
		t.Fatalf("unmarshal team: %v", err)
	}
	if created.ID == "" || created.Name != "Engineering" || !created.Enabled {
		t.Fatalf("unexpected created team: %+v", created)
	}

	// Assign key via contract method and via HTTP.
	if err := mod.AssignKeyToTeam(context.Background(), created.ID, "key_1"); err != nil {
		t.Fatalf("AssignKeyToTeam: %v", err)
	}
	addKeyBody, _ := json.Marshal(map[string]string{"key_id": "key_2"})
	addRec := httptest.NewRecorder()
	addReq := httptest.NewRequest(http.MethodPost, "/api/teams/"+created.ID+"/keys", bytes.NewReader(addKeyBody))
	addReq.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(addRec, addReq)
	if addRec.Code != http.StatusOK {
		t.Fatalf("POST /api/teams/{id}/keys status = %d, want 200; body=%s", addRec.Code, addRec.Body.String())
	}

	// Get team -> has 2 keys.
	getRec := httptest.NewRecorder()
	router.ServeHTTP(getRec, httptest.NewRequest(http.MethodGet, "/api/teams/"+created.ID, nil))
	if getRec.Code != http.StatusOK {
		t.Fatalf("GET /api/teams/{id} status = %d, want 200", getRec.Code)
	}
	var fetched struct {
		KeyIDs []string `json:"key_ids"`
	}
	_ = json.Unmarshal(getRec.Body.Bytes(), &fetched)
	if len(fetched.KeyIDs) != 2 {
		t.Fatalf("expected 2 keys in team, got %v", fetched.KeyIDs)
	}

	// Remove key_2.
	rmRec := httptest.NewRecorder()
	router.ServeHTTP(rmRec, httptest.NewRequest(http.MethodDelete, "/api/teams/"+created.ID+"/keys/key_2", nil))
	if rmRec.Code != http.StatusOK {
		t.Fatalf("DELETE key status = %d, want 200", rmRec.Code)
	}

	// Delete team.
	delRec := httptest.NewRecorder()
	router.ServeHTTP(delRec, httptest.NewRequest(http.MethodDelete, "/api/teams/"+created.ID, nil))
	if delRec.Code != http.StatusOK {
		t.Fatalf("DELETE team status = %d, want 200", delRec.Code)
	}
}

func TestTeamsBudgetAndRPMEnforcementViaInterceptor(t *testing.T) {
	t.Parallel()

	mod, router, bus := setupTeamsTest(t)

	// Create team with tight $1.00 budget and 2 RPM limit.
	createBody, _ := json.Marshal(map[string]any{
		"name":           "Limited Team",
		"budget_usd":     1.0,
		"rate_limit_rpm": 2,
	})
	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/api/teams", bytes.NewReader(createBody))
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(rec, req)
	if rec.Code != http.StatusCreated {
		t.Fatalf("create team = %d: %s", rec.Code, rec.Body.String())
	}
	var team struct {
		ID string `json:"id"`
	}
	_ = json.Unmarshal(rec.Body.Bytes(), &team)

	if err := mod.AssignKeyToTeam(context.Background(), team.ID, "key_limited"); err != nil {
		t.Fatalf("AssignKeyToTeam: %v", err)
	}

	// First 2 requests pass RPM.
	for i := range 2 {
		if err := mod.OnRequest(context.Background(), &contracts.CanonicalRequest{
			RequestID: "req",
			KeyID:     "key_limited",
			Metadata:  map[string]any{},
		}); err != nil {
			t.Fatalf("request %d unexpectedly rejected: %v", i+1, err)
		}
	}

	// 3rd request within the same minute must be rejected with 429.
	rpmErr := mod.OnRequest(context.Background(), &contracts.CanonicalRequest{
		RequestID: "req_3",
		KeyID:     "key_limited",
		Metadata:  map[string]any{},
	})
	if rpmErr == nil {
		t.Fatal("expected 3rd request to exceed team RPM limit")
	}
	type statusCoder interface{ StatusCode() int }
	if sc, ok := rpmErr.(statusCoder); !ok || sc.StatusCode() != http.StatusTooManyRequests {
		t.Fatalf("expected 429 status error for RPM, got %T: %v", rpmErr, rpmErr)
	}

	// Remove RPM limit so we can test USD budget enforcement cleanly.
	patchBody, _ := json.Marshal(map[string]any{"rate_limit_rpm": 0})
	patchRec := httptest.NewRecorder()
	patchReq := httptest.NewRequest(http.MethodPatch, "/api/teams/"+team.ID, bytes.NewReader(patchBody))
	patchReq.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(patchRec, patchReq)
	if patchRec.Code != http.StatusOK {
		t.Fatalf("PATCH team = %d: %s", patchRec.Code, patchRec.Body.String())
	}

	// Emit usage.recorded totaling $1.25 -> exceeds $1.00 budget.
	bus.Publish("usage.recorded", map[string]any{
		"key_id":   "key_limited",
		"cost_usd": 1.25,
	})

	deadline := time.Now().Add(2 * time.Second)
	var budgetErr error
	for time.Now().Before(deadline) {
		budgetErr = mod.OnRequest(context.Background(), &contracts.CanonicalRequest{
			RequestID: "req_budget",
			KeyID:     "key_limited",
			Metadata:  map[string]any{},
		})
		if budgetErr != nil {
			break
		}
		time.Sleep(15 * time.Millisecond)
	}
	if budgetErr == nil {
		t.Fatal("expected request to be rejected after team budget exceeded")
	}
	if sc, ok := budgetErr.(statusCoder); !ok || sc.StatusCode() != http.StatusTooManyRequests {
		t.Fatalf("expected 429 status error for budget, got %T: %v", budgetErr, budgetErr)
	}

	// Reset spend -> requests allowed again.
	resetRec := httptest.NewRecorder()
	router.ServeHTTP(resetRec, httptest.NewRequest(http.MethodPost, "/api/teams/"+team.ID+"/reset-spend", nil))
	if resetRec.Code != http.StatusOK {
		t.Fatalf("POST reset-spend = %d: %s", resetRec.Code, resetRec.Body.String())
	}
	if err := mod.OnRequest(context.Background(), &contracts.CanonicalRequest{
		RequestID: "req_after_reset",
		KeyID:     "key_limited",
		Metadata:  map[string]any{},
	}); err != nil {
		t.Fatalf("request after reset-spend rejected: %v", err)
	}
}
