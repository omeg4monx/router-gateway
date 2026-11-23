package vouchers_test

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"path/filepath"
	"sync"
	"sync/atomic"
	"testing"

	"github.com/go-chi/chi/v5"

	"github.com/omeg4monx/router-gateway/internal/contracts"
	"github.com/omeg4monx/router-gateway/internal/modules/vouchers"
	"github.com/omeg4monx/router-gateway/internal/platform/db"
	"github.com/omeg4monx/router-gateway/internal/platform/eventbus"
	"github.com/omeg4monx/router-gateway/internal/platform/lifecycle"
)

type stubKeyProvisioner struct {
	counter atomic.Int64
}

func (s *stubKeyProvisioner) ProvisionKey(_ context.Context, req contracts.ProvisionKeyRequest) (contracts.ProvisionedKey, error) {
	n := s.counter.Add(1)
	return contracts.ProvisionedKey{
		ID:        fmt.Sprintf("key_v_%d", n),
		Name:      req.Name,
		KeyPrefix: "sk-omg-vouch",
		Secret:    fmt.Sprintf("sk-omg-vouch-secret-%d", n),
	}, nil
}

type stubTeamAssigner struct {
	mu          sync.Mutex
	assignments []string
}

func (s *stubTeamAssigner) AssignKeyToTeam(_ context.Context, teamID, keyID string) error {
	s.mu.Lock()
	s.assignments = append(s.assignments, teamID+":"+keyID)
	s.mu.Unlock()
	return nil
}

func setupVouchersTest(t *testing.T, maxClaimsPerIP int) (*vouchers.Module, chi.Router, *stubKeyProvisioner, *stubTeamAssigner) {
	t.Helper()

	database, err := db.Open(filepath.Join(t.TempDir(), "vouchers.sqlite"))
	if err != nil {
		t.Fatalf("db.Open: %v", err)
	}
	t.Cleanup(func() { _ = database.Close() })

	prov := &stubKeyProvisioner{}
	assigner := &stubTeamAssigner{}
	mod := vouchers.New(vouchers.Options{
		DB:             database,
		Bus:            eventbus.New(),
		KeyProvisioner: prov,
		TeamAssigner:   assigner,
		MaxClaimsPerIP: maxClaimsPerIP,
	})
	if err := database.MigrateModule(context.Background(), mod.Name(), mod.Migrations()); err != nil {
		t.Fatalf("MigrateModule: %v", err)
	}

	r := chi.NewRouter()
	r.Route("/api", func(api chi.Router) {
		mod.RegisterRoutes(lifecycle.RouteMounts{
			PublicAPI: api,
			AdminAPI:  api,
		})
	})

	if err := mod.Start(context.Background()); err != nil {
		t.Fatalf("mod.Start: %v", err)
	}
	t.Cleanup(func() { _ = mod.Stop(context.Background()) })

	return mod, r, prov, assigner
}

func TestVouchersCRUDClaimTeamBindingAndRedemptions(t *testing.T) {
	t.Parallel()

	_, router, prov, assigner := setupVouchersTest(t, 20)

	createBody, _ := json.Marshal(map[string]any{
		"code":            "LAUNCH2026",
		"name":            "Launch Promo",
		"quota_usd":       10.0,
		"rate_limit_rpm":  30,
		"allowed_models":  []string{"gpt-4o"},
		"team_id":         "team_promo",
		"max_redemptions": 2,
	})
	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/api/vouchers", bytes.NewReader(createBody))
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(rec, req)
	if rec.Code != http.StatusCreated {
		t.Fatalf("POST /api/vouchers status = %d, want 201; body=%s", rec.Code, rec.Body.String())
	}
	var created struct {
		ID   string `json:"id"`
		Code string `json:"code"`
	}
	_ = json.Unmarshal(rec.Body.Bytes(), &created)

	// Claim by alice -> 201 with provisioned key secret and team assignment.
	claim1Body, _ := json.Marshal(map[string]any{"code": "launch2026", "claimer_id": "alice"})
	claim1Rec := httptest.NewRecorder()
	claim1Req := httptest.NewRequest(http.MethodPost, "/api/vouchers/claim", bytes.NewReader(claim1Body))
	claim1Req.Header.Set("Content-Type", "application/json")
	claim1Req.RemoteAddr = "10.0.0.1:1234"
	router.ServeHTTP(claim1Rec, claim1Req)
	if claim1Rec.Code != http.StatusOK {
		t.Fatalf("POST /api/vouchers/claim status = %d, want 200; body=%s", claim1Rec.Code, claim1Rec.Body.String())
	}
	var claimResp struct {
		KeyID  string `json:"key_id"`
		Secret string `json:"secret"`
	}
	_ = json.Unmarshal(claim1Rec.Body.Bytes(), &claimResp)
	if claimResp.KeyID == "" || claimResp.Secret == "" {
		t.Fatalf("expected provisioned key in claim response, got %+v", claimResp)
	}
	if prov.counter.Load() != 1 {
		t.Fatalf("provisioned keys = %d, want 1", prov.counter.Load())
	}
	assigner.mu.Lock()
	gotAssign := len(assigner.assignments)
	assigner.mu.Unlock()
	if gotAssign != 1 {
		t.Fatalf("team assignments = %d, want 1", gotAssign)
	}

	// Duplicate claim by alice -> 409 Conflict (anti-abuse per-actor uniqueness).
	dupRec := httptest.NewRecorder()
	dupReq := httptest.NewRequest(http.MethodPost, "/api/vouchers/claim", bytes.NewReader(claim1Body))
	dupReq.Header.Set("Content-Type", "application/json")
	dupReq.RemoteAddr = "10.0.0.2:1234"
	router.ServeHTTP(dupRec, dupReq)
	if dupRec.Code != http.StatusConflict {
		t.Fatalf("duplicate claim by alice status = %d, want 409; body=%s", dupRec.Code, dupRec.Body.String())
	}

	// List redemptions -> 1 entry.
	redRec := httptest.NewRecorder()
	router.ServeHTTP(redRec, httptest.NewRequest(http.MethodGet, "/api/vouchers/"+created.ID+"/redemptions", nil))
	if redRec.Code != http.StatusOK {
		t.Fatalf("GET redemptions status = %d, want 200", redRec.Code)
	}
}

func TestVouchersConcurrentSingleUseAtomicityAndIPRateLimit(t *testing.T) {
	t.Parallel()

	_, router, prov, _ := setupVouchersTest(t, 2)

	createBody, _ := json.Marshal(map[string]any{
		"code":            "SINGLEUSE",
		"name":            "One-shot",
		"quota_usd":       5.0,
		"max_redemptions": 1,
	})
	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/api/vouchers", bytes.NewReader(createBody))
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(rec, req)
	if rec.Code != http.StatusCreated {
		t.Fatalf("create voucher = %d: %s", rec.Code, rec.Body.String())
	}

	// 8 concurrent callers from distinct IPs & claimer_ids -> exactly 1 succeeds, 7 get 409.
	var okCount, conflictCount atomic.Int32
	var wg sync.WaitGroup
	for i := range 8 {
		wg.Add(1)
		go func(idx int) {
			defer wg.Done()
			body, _ := json.Marshal(map[string]any{
				"code":       "SINGLEUSE",
				"claimer_id": fmt.Sprintf("user_%d", idx),
			})
			r := httptest.NewRequest(http.MethodPost, "/api/vouchers/claim", bytes.NewReader(body))
			r.Header.Set("Content-Type", "application/json")
			r.RemoteAddr = fmt.Sprintf("10.1.0.%d:5555", idx+1)
			w := httptest.NewRecorder()
			router.ServeHTTP(w, r)
			switch w.Code {
			case http.StatusOK:
				okCount.Add(1)
			case http.StatusConflict:
				conflictCount.Add(1)
			}
		}(i)
	}
	wg.Wait()

	if okCount.Load() != 1 || conflictCount.Load() != 7 || prov.counter.Load() != 1 {
		t.Fatalf("concurrent claim results: ok=%d conflict=%d keys=%d, want 1/7/1",
			okCount.Load(), conflictCount.Load(), prov.counter.Load())
	}

	// IP rate limit test: 3rd attempt from same IP (limit=2) returns 429 Too Many Requests.
	for i := range 3 {
		body, _ := json.Marshal(map[string]any{"code": "UNKNOWN", "claimer_id": fmt.Sprintf("spam_%d", i)})
		r := httptest.NewRequest(http.MethodPost, "/api/vouchers/claim", bytes.NewReader(body))
		r.Header.Set("Content-Type", "application/json")
		r.RemoteAddr = "192.168.1.99:9999"
		w := httptest.NewRecorder()
		router.ServeHTTP(w, r)
		if i == 2 && w.Code != http.StatusTooManyRequests {
			t.Fatalf("expected 429 on 3rd rapid claim from same IP, got %d: %s", w.Code, w.Body.String())
		}
	}
}
