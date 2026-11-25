package security_test

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/go-chi/chi/v5"

	"github.com/omeg4monx/router-gateway/internal/modules/security"
	"github.com/omeg4monx/router-gateway/internal/platform/config"
	"github.com/omeg4monx/router-gateway/internal/platform/db"
	"github.com/omeg4monx/router-gateway/internal/platform/eventbus"
	"github.com/omeg4monx/router-gateway/internal/platform/lifecycle"
)

type stubCredentialStore struct {
	requireLogin bool
}

func (s *stubCredentialStore) AdminPasswordHash(context.Context) (string, error) { return "hash", nil }
func (s *stubCredentialStore) UpdateAdminPasswordHash(context.Context, string) error {
	return nil
}
func (s *stubCredentialStore) RequireLogin(context.Context) (bool, error) {
	return s.requireLogin, nil
}
func (s *stubCredentialStore) SetRequireLogin(_ context.Context, required bool) error {
	s.requireLogin = required
	return nil
}

func setupSecurityTest(t *testing.T, cfg config.Config, store *stubCredentialStore) (*security.Module, chi.Router, *eventbus.Bus, *db.DB) {
	t.Helper()

	database, err := db.Open(filepath.Join(t.TempDir(), "security.sqlite"))
	if err != nil {
		t.Fatalf("db.Open: %v", err)
	}
	t.Cleanup(func() { _ = database.Close() })

	bus := eventbus.New()
	mod := security.New(security.Options{
		DB:     database,
		Bus:    bus,
		Config: cfg,
		Store:  store,
		ActiveModules: func() []string {
			return []string{"system", "auth", "security"}
		},
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

	return mod, r, bus, database
}

func TestSecurityRedactionAndTamperEvidentClear(t *testing.T) {
	t.Parallel()

	_, router, _, _ := setupSecurityTest(t, config.Config{
		Host: "127.0.0.1",
		Port: 20128,
	}, &stubCredentialStore{requireLogin: true})

	payload := map[string]any{
		"category": "keys",
		"action":   "key.tested",
		"severity": "warning",
		"summary":  "Failed auth with Authorization: Bearer top-secret-token-123 and key sk-proj-abc1234567890",
		"metadata": map[string]any{
			"api_key":  "sk-live-secret-9999",
			"password": "my-super-secret-password",
			"safe":     "visible-value",
		},
	}
	body, _ := json.Marshal(payload)
	req := httptest.NewRequest(http.MethodPost, "/api/security/audit", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)

	if rec.Code != http.StatusCreated {
		t.Fatalf("POST /api/security/audit status = %d, want 201; body=%s", rec.Code, rec.Body.String())
	}

	var created struct {
		ID       string         `json:"id"`
		Summary  string         `json:"summary"`
		Metadata map[string]any `json:"metadata"`
	}
	if err := json.Unmarshal(rec.Body.Bytes(), &created); err != nil {
		t.Fatalf("unmarshal created: %v", err)
	}
	if strings.Contains(created.Summary, "top-secret-token-123") || strings.Contains(created.Summary, "sk-proj-abc1234567890") {
		t.Fatalf("summary leaked secret: %q", created.Summary)
	}
	if !strings.Contains(created.Summary, "[REDACTED]") {
		t.Fatalf("expected [REDACTED] in summary, got %q", created.Summary)
	}
	if created.Metadata["api_key"] != "[REDACTED]" || created.Metadata["password"] != "[REDACTED]" {
		t.Fatalf("metadata leaked secret: %+v", created.Metadata)
	}
	if created.Metadata["safe"] != "visible-value" {
		t.Fatalf("non-secret metadata altered: %+v", created.Metadata)
	}

	// Clear must leave a tamper-evident security.audit_cleared marker event.
	delReq := httptest.NewRequest(http.MethodDelete, "/api/security/audit", nil)
	delRec := httptest.NewRecorder()
	router.ServeHTTP(delRec, delReq)
	if delRec.Code != http.StatusOK {
		t.Fatalf("DELETE /api/security/audit status = %d, want 200; body=%s", delRec.Code, delRec.Body.String())
	}

	listReq := httptest.NewRequest(http.MethodGet, "/api/security/audit", nil)
	listRec := httptest.NewRecorder()
	router.ServeHTTP(listRec, listReq)
	if listRec.Code != http.StatusOK {
		t.Fatalf("GET /api/security/audit status = %d, want 200", listRec.Code)
	}
	var listed struct {
		Items []struct {
			Category string `json:"category"`
			Action   string `json:"action"`
		} `json:"items"`
		Total int `json:"total"`
	}
	if err := json.Unmarshal(listRec.Body.Bytes(), &listed); err != nil {
		t.Fatalf("unmarshal list: %v", err)
	}
	if listed.Total != 1 || len(listed.Items) != 1 || listed.Items[0].Action != "audit.cleared" {
		t.Fatalf("expected 1 audit.cleared marker after clear, got total=%d items=%+v", listed.Total, listed.Items)
	}
}

func TestSecurityEventBusAutoAuditIngestion(t *testing.T) {
	t.Parallel()

	_, router, bus, _ := setupSecurityTest(t, config.Config{
		Host: "127.0.0.1",
		Port: 20128,
	}, &stubCredentialStore{requireLogin: true})

	bus.Publish("auth.policy_changed", map[string]any{"require_login": false})
	bus.Publish("keys.changed", map[string]any{"action": "created", "key_id": "key_1"})
	bus.Publish(security.TopicAudit, security.AuditInput{
		Category:   "gateway",
		Action:     "auth.rejected",
		Severity:   "critical",
		Actor:      "client",
		ResourceID: "req_42",
		Summary:    "Invalid gateway API key presented",
	})

	deadline := time.Now().Add(2 * time.Second)
	var total int
	for time.Now().Before(deadline) {
		req := httptest.NewRequest(http.MethodGet, "/api/security/audit", nil)
		rec := httptest.NewRecorder()
		router.ServeHTTP(rec, req)
		var resp struct {
			Total int `json:"total"`
		}
		_ = json.Unmarshal(rec.Body.Bytes(), &resp)
		total = resp.Total
		if total >= 3 {
			break
		}
		time.Sleep(15 * time.Millisecond)
	}
	if total < 3 {
		t.Fatalf("expected >= 3 auto-ingested audit events, got %d", total)
	}

	// Filter by severity=critical and q=req_42
	filterReq := httptest.NewRequest(http.MethodGet, "/api/security/audit?severity=critical&q=Invalid", nil)
	filterRec := httptest.NewRecorder()
	router.ServeHTTP(filterRec, filterReq)
	if filterRec.Code != http.StatusOK {
		t.Fatalf("GET filtered audit status = %d, want 200", filterRec.Code)
	}
	var filtered struct {
		Items []struct {
			Category   string `json:"category"`
			Action     string `json:"action"`
			ResourceID string `json:"resource_id"`
		} `json:"items"`
		Total int `json:"total"`
	}
	if err := json.Unmarshal(filterRec.Body.Bytes(), &filtered); err != nil {
		t.Fatalf("unmarshal filtered: %v", err)
	}
	if filtered.Total != 1 || len(filtered.Items) != 1 || filtered.Items[0].ResourceID != "req_42" {
		t.Fatalf("filtered audit = %+v, want 1 item with resource_id=req_42", filtered)
	}
}

func TestSecurityPostureAndExportAndPrune(t *testing.T) {
	t.Parallel()

	mod, router, _, database := setupSecurityTest(t, config.Config{
		Host:                "0.0.0.0",
		Port:                20128,
		AllowInsecureRemote: true,
		MasterKeyHex:        strings.Repeat("ab", 32),
		DisabledModules:     map[string]bool{"translator": true},
	}, &stubCredentialStore{requireLogin: false})

	// Seed an expired audit event directly (> 90 days old) and one recent event.
	oldTime := time.Now().UTC().AddDate(0, 0, -120).Format("2006-01-02T15:04:05.000000000Z")
	newTime := time.Now().UTC().Format("2006-01-02T15:04:05.000000000Z")
	_, err := database.Writer().ExecContext(context.Background(),
		`INSERT INTO security_audit_events(id, category, action, severity, actor, resource_id, summary, metadata_json, ip_address, created_at)
		 VALUES ('old_1', 'system', 'boot', 'info', 'system', '', 'old event', '{}', '', ?),
		        ('new_1', 'system', 'boot', 'info', 'system', '', 'new event', '{}', '', ?)`,
		oldTime, newTime,
	)
	if err != nil {
		t.Fatalf("seed audit events: %v", err)
	}

	if err := mod.PruneExpired(context.Background()); err != nil {
		t.Fatalf("PruneExpired: %v", err)
	}

	// Posture check should report non-loopback + allow_insecure_remote + require_login=false warnings and master_key_source=env.
	postureReq := httptest.NewRequest(http.MethodGet, "/api/security/posture", nil)
	postureRec := httptest.NewRecorder()
	router.ServeHTTP(postureRec, postureReq)
	if postureRec.Code != http.StatusOK {
		t.Fatalf("GET /api/security/posture status = %d, want 200; body=%s", postureRec.Code, postureRec.Body.String())
	}
	var posture struct {
		BindAddr            string   `json:"bind_addr"`
		IsLoopback          bool     `json:"is_loopback"`
		RequireLogin        bool     `json:"require_login"`
		AllowInsecureRemote bool     `json:"allow_insecure_remote"`
		MasterKeySource     string   `json:"master_key_source"`
		AuditEventCount     int      `json:"audit_event_count"`
		Warnings            []string `json:"warnings"`
	}
	if err := json.Unmarshal(postureRec.Body.Bytes(), &posture); err != nil {
		t.Fatalf("unmarshal posture: %v", err)
	}
	if posture.IsLoopback || posture.RequireLogin || !posture.AllowInsecureRemote {
		t.Fatalf("unexpected posture flags: %+v", posture)
	}
	if posture.MasterKeySource != "env" {
		t.Fatalf("MasterKeySource = %q, want env", posture.MasterKeySource)
	}
	if posture.AuditEventCount != 1 {
		t.Fatalf("AuditEventCount after prune = %d, want 1", posture.AuditEventCount)
	}
	if len(posture.Warnings) == 0 {
		t.Fatalf("expected security warnings for insecure remote bind, got none")
	}

	// Export check.
	exportReq := httptest.NewRequest(http.MethodGet, "/api/security/audit/export", nil)
	exportRec := httptest.NewRecorder()
	router.ServeHTTP(exportRec, exportReq)
	if exportRec.Code != http.StatusOK {
		t.Fatalf("GET /api/security/audit/export status = %d, want 200", exportRec.Code)
	}
	if !strings.Contains(exportRec.Header().Get("Content-Disposition"), "security-audit.json") {
		t.Fatalf("unexpected Content-Disposition: %q", exportRec.Header().Get("Content-Disposition"))
	}
	if !strings.Contains(exportRec.Body.String(), "new_1") || strings.Contains(exportRec.Body.String(), "old_1") {
		t.Fatalf("export body mismatch: %s", exportRec.Body.String())
	}
}

func TestSecurityWave2BTopicParityAndExtendedRedaction(t *testing.T) {
	t.Parallel()

	_, router, bus, _ := setupSecurityTest(t, config.Config{
		Host: "127.0.0.1",
		Port: 20128,
	}, &stubCredentialStore{requireLogin: true})

	// Pin Important #1: actual Wave 2B publisher topic names.
	bus.Publish("promptinjection.changed", nil)
	bus.Publish("modeleditor.changed", nil)
	bus.Publish("tokensaver.changed", nil)

	// Pin Important #2: quoted/JSON secrets, bare token=, and typed Go maps/slices.
	bus.Publish(security.TopicAudit, security.AuditInput{
		Category: "gateway",
		Action:   "debug.capture",
		Severity: "warning",
		Summary:  `Leak check token=bare-tok-999 password="quoted-pass-1" {"api_key":"json-key-2"}`,
		Metadata: map[string]any{
			"headers": map[string]string{
				"Authorization": "Bearer nested-bearer-secret",
				"X-Trace":       "token=nested-tok-888",
			},
			"lines": []string{`password='single-quoted-pw'`},
		},
	})

	deadline := time.Now().Add(2 * time.Second)
	var items []struct {
		Category string         `json:"category"`
		Summary  string         `json:"summary"`
		Metadata map[string]any `json:"metadata"`
	}
	for time.Now().Before(deadline) {
		req := httptest.NewRequest(http.MethodGet, "/api/security/audit?limit=20", nil)
		rec := httptest.NewRecorder()
		router.ServeHTTP(rec, req)
		var resp struct {
			Items []struct {
				Category string         `json:"category"`
				Summary  string         `json:"summary"`
				Metadata map[string]any `json:"metadata"`
			} `json:"items"`
		}
		_ = json.Unmarshal(rec.Body.Bytes(), &resp)
		items = resp.Items
		if len(items) >= 4 {
			break
		}
		time.Sleep(15 * time.Millisecond)
	}

	seen := map[string]bool{}
	for _, it := range items {
		seen[it.Category] = true
		if it.Category == "gateway" {
			rawMeta, _ := json.Marshal(it.Metadata)
			combined := it.Summary + " " + string(rawMeta)
			for _, leaked := range []string{"bare-tok-999", "quoted-pass-1", "json-key-2", "nested-bearer-secret", "nested-tok-888", "single-quoted-pw"} {
				if strings.Contains(combined, leaked) {
					t.Fatalf("audit record leaked secret %q: %s", leaked, combined)
				}
			}
		}
	}
	for _, wantCat := range []string{"promptinjection", "modeleditor", "tokensaver", "gateway"} {
		if !seen[wantCat] {
			t.Fatalf("missing auto-audited category %q; seen=%v", wantCat, seen)
		}
	}
}

func TestSecurityAuditPaginationOffset(t *testing.T) {
	t.Parallel()

	_, router, _, _ := setupSecurityTest(t, config.Config{
		Host: "127.0.0.1",
		Port: 20128,
	}, &stubCredentialStore{requireLogin: true})

	for i := 1; i <= 5; i++ {
		body := `{"category":"security","action":"page.item","severity":"info","summary":"event-` + string(rune('0'+i)) + `"}`
		req := httptest.NewRequest(http.MethodPost, "/api/security/audit", strings.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		rec := httptest.NewRecorder()
		router.ServeHTTP(rec, req)
		if rec.Code != http.StatusCreated {
			t.Fatalf("seed event %d status = %d: %s", i, rec.Code, rec.Body.String())
		}
	}

	// Page 1 (limit=2, offset=0)
	req1 := httptest.NewRequest(http.MethodGet, "/api/security/audit?limit=2&offset=0", nil)
	rec1 := httptest.NewRecorder()
	router.ServeHTTP(rec1, req1)
	if rec1.Code != http.StatusOK {
		t.Fatalf("page 1 status = %d: %s", rec1.Code, rec1.Body.String())
	}
	var page1 struct {
		Items []struct {
			ID string `json:"id"`
		} `json:"items"`
		Total int `json:"total"`
	}
	_ = json.Unmarshal(rec1.Body.Bytes(), &page1)
	if page1.Total != 5 || len(page1.Items) != 2 {
		t.Fatalf("page 1 = total %d len %d, want 5 and 2", page1.Total, len(page1.Items))
	}

	// Page 2 (limit=2, offset=2) must return distinct items.
	req2 := httptest.NewRequest(http.MethodGet, "/api/security/audit?limit=2&offset=2", nil)
	rec2 := httptest.NewRecorder()
	router.ServeHTTP(rec2, req2)
	if rec2.Code != http.StatusOK {
		t.Fatalf("page 2 status = %d: %s", rec2.Code, rec2.Body.String())
	}
	var page2 struct {
		Items []struct {
			ID string `json:"id"`
		} `json:"items"`
		Total int `json:"total"`
	}
	_ = json.Unmarshal(rec2.Body.Bytes(), &page2)
	if page2.Total != 5 || len(page2.Items) != 2 {
		t.Fatalf("page 2 = total %d len %d, want 5 and 2", page2.Total, len(page2.Items))
	}
	if page2.Items[0].ID == page1.Items[0].ID || page2.Items[0].ID == page1.Items[1].ID {
		t.Fatalf("page 2 overlapped with page 1: page1=%+v page2=%+v", page1.Items, page2.Items)
	}

	// Negative offset must return 400.
	badReq := httptest.NewRequest(http.MethodGet, "/api/security/audit?limit=2&offset=-1", nil)
	badRec := httptest.NewRecorder()
	router.ServeHTTP(badRec, badReq)
	if badRec.Code != http.StatusBadRequest {
		t.Fatalf("negative offset status = %d, want 400", badRec.Code)
	}
}
