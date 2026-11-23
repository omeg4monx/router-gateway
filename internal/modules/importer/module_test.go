package importer_test

import (
	"context"
	"database/sql"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"testing"

	"github.com/go-chi/chi/v5"
	_ "modernc.org/sqlite"

	"github.com/omeg4monx/router-gateway/internal/contracts"
	"github.com/omeg4monx/router-gateway/internal/modules/importer"
	"github.com/omeg4monx/router-gateway/internal/platform/db"
	"github.com/omeg4monx/router-gateway/internal/platform/eventbus"
	"github.com/omeg4monx/router-gateway/internal/platform/lifecycle"
)

type stubSettings struct {
	mu     sync.Mutex
	values map[string]any
}

func (s *stubSettings) ExportSettings(context.Context) (map[string]any, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.values, nil
}

func (s *stubSettings) ImportSettings(_ context.Context, values map[string]any, _, _ string) (int, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	if s.values == nil {
		s.values = make(map[string]any)
	}
	for k, v := range values {
		s.values[k] = v
	}
	return len(values), nil
}

type stubCombos struct {
	mu      sync.Mutex
	combos  []contracts.ComboSnapshot
	aliases []contracts.AliasSnapshot
}

func (s *stubCombos) ExportCombos(context.Context) ([]contracts.ComboSnapshot, []contracts.AliasSnapshot, error) {
	return s.combos, s.aliases, nil
}

func (s *stubCombos) ImportCombos(_ context.Context, combos []contracts.ComboSnapshot, aliases []contracts.AliasSnapshot, _ bool, _ string) (int, int, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.combos = append(s.combos, combos...)
	s.aliases = append(s.aliases, aliases...)
	return len(combos), len(aliases), nil
}

type stubKeyProvisioner struct {
	mu   sync.Mutex
	keys []contracts.ProvisionKeyRequest
}

func (s *stubKeyProvisioner) ProvisionKey(_ context.Context, req contracts.ProvisionKeyRequest) (contracts.ProvisionedKey, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.keys = append(s.keys, req)
	return contracts.ProvisionedKey{
		ID:        "key_imported",
		Name:      req.Name,
		KeyPrefix: "om_imp_",
		Secret:    "om_imp_secret",
	}, nil
}

func (s *stubKeyProvisioner) RevokeKey(context.Context, string) error { return nil }

func createLegacySQLite(t *testing.T) string {
	t.Helper()
	path := filepath.Join(t.TempDir(), "legacy-data.sqlite")
	conn, err := sql.Open("sqlite", path)
	if err != nil {
		t.Fatalf("open legacy sqlite: %v", err)
	}
	defer conn.Close()

	stmts := []string{
		`CREATE TABLE settings (key TEXT PRIMARY KEY, value TEXT)`,
		`INSERT INTO settings (key, value) VALUES ('theme', '"dark"'), ('locale', '"en"')`,
		`CREATE TABLE combos (id TEXT PRIMARY KEY, name TEXT, data TEXT)`,
		`INSERT INTO combos (id, name, data) VALUES ('c1', 'legacy-primary', '{"description":"Migrated","models":["openai/gpt-4o","anthropic/claude-3-5-sonnet"],"isActive":true}')`,
		`CREATE TABLE apiKeys (id TEXT PRIMARY KEY, name TEXT, key TEXT, data TEXT)`,
		`INSERT INTO apiKeys (id, name, key, data) VALUES ('k1', 'legacy-dev', 'sk-legacy-1', '{"allowedModels":["gpt-4o"],"rateLimitRpm":60}')`,
		`CREATE TABLE providerConnections (id TEXT PRIMARY KEY, provider TEXT, name TEXT, data TEXT)`,
		`INSERT INTO providerConnections (id, provider, name, data) VALUES ('p1', 'openai', 'Work OpenAI', '{}')`,
	}
	for _, stmt := range stmts {
		if _, err := conn.Exec(stmt); err != nil {
			t.Fatalf("exec %q: %v", stmt, err)
		}
	}
	return path
}

func setupImporterTest(t *testing.T) (http.Handler, *stubSettings, *stubCombos, *stubKeyProvisioner) {
	t.Helper()
	database, err := db.Open(filepath.Join(t.TempDir(), "importer.sqlite"))
	if err != nil {
		t.Fatalf("db.Open: %v", err)
	}
	t.Cleanup(func() { _ = database.Close() })

	st := &stubSettings{}
	cb := &stubCombos{}
	kp := &stubKeyProvisioner{}
	mod := importer.New(importer.Options{
		DB:             database,
		Bus:            eventbus.New(),
		Settings:       st,
		Combos:         cb,
		KeyProvisioner: kp,
	})

	ctx := context.Background()
	if err := database.MigrateModule(ctx, mod.Name(), mod.Migrations()); err != nil {
		t.Fatalf("MigrateModule: %v", err)
	}
	if err := mod.Start(ctx); err != nil {
		t.Fatalf("Start: %v", err)
	}
	t.Cleanup(func() { _ = mod.Stop(context.Background()) })

	r := chi.NewRouter()
	api := chi.NewRouter()
	mod.RegisterRoutes(lifecycle.RouteMounts{AdminAPI: api})
	r.Mount("/api", api)
	return r, st, cb, kp
}

func TestImporterRejectsInvalidFilesAndPreservesReadOnlySource(t *testing.T) {
	handler, st, cb, kp := setupImporterTest(t)

	// Missing file -> 400.
	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, httptest.NewRequest(http.MethodPost, "/api/importer/inspect", strings.NewReader(`{"source_path":"/nonexistent/data.sqlite"}`)))
	if rec.Code != http.StatusBadRequest {
		t.Fatalf("expected 400 for missing file, got %d (%s)", rec.Code, rec.Body.String())
	}

	// Non-SQLite file -> 400.
	txtFile := filepath.Join(t.TempDir(), "not-sqlite.db")
	if err := os.WriteFile(txtFile, []byte("plain text not sqlite"), 0o600); err != nil {
		t.Fatalf("write txt file: %v", err)
	}
	rec = httptest.NewRecorder()
	handler.ServeHTTP(rec, httptest.NewRequest(http.MethodPost, "/api/importer/inspect", strings.NewReader(`{"source_path":"`+txtFile+`"}`)))
	if rec.Code != http.StatusBadRequest {
		t.Fatalf("expected 400 for non-SQLite file, got %d (%s)", rec.Code, rec.Body.String())
	}

	// Valid legacy SQLite inspection must leave file unmodified.
	legacyPath := createLegacySQLite(t)
	beforeInfo, err := os.Stat(legacyPath)
	if err != nil {
		t.Fatalf("stat before: %v", err)
	}

	rec = httptest.NewRecorder()
	handler.ServeHTTP(rec, httptest.NewRequest(http.MethodPost, "/api/importer/inspect", strings.NewReader(`{"source_path":"`+legacyPath+`"}`)))
	if rec.Code != http.StatusOK {
		t.Fatalf("POST /api/importer/inspect status = %d, body = %s", rec.Code, rec.Body.String())
	}
	var inspectResp struct {
		SettingsCount  int      `json:"settings_count"`
		CombosCount    int      `json:"combos_count"`
		KeysCount      int      `json:"keys_count"`
		ProvidersCount int      `json:"providers_count"`
		TablesFound    []string `json:"tables_found"`
	}
	if err := json.Unmarshal(rec.Body.Bytes(), &inspectResp); err != nil {
		t.Fatalf("decode inspect: %v", err)
	}
	if inspectResp.SettingsCount != 2 || inspectResp.CombosCount != 1 || inspectResp.KeysCount != 1 || inspectResp.ProvidersCount != 1 {
		t.Fatalf("unexpected inspect counts: %+v", inspectResp)
	}

	// Dry-run performs zero mutations.
	dryBody := `{"source_path":"` + legacyPath + `","dry_run":true,"import_settings":true,"import_combos":true,"import_keys":true}`
	rec = httptest.NewRecorder()
	handler.ServeHTTP(rec, httptest.NewRequest(http.MethodPost, "/api/importer/run", strings.NewReader(dryBody)))
	if rec.Code != http.StatusOK {
		t.Fatalf("POST /api/importer/run (dry_run) status = %d, body = %s", rec.Code, rec.Body.String())
	}
	if len(st.values) != 0 || len(cb.combos) != 0 || len(kp.keys) != 0 {
		t.Fatalf("dry_run must not mutate target stores")
	}

	// Live import applies settings, combos, and provisions keys.
	liveBody := `{"source_path":"` + legacyPath + `","dry_run":false,"import_settings":true,"import_combos":true,"import_keys":true}`
	rec = httptest.NewRecorder()
	handler.ServeHTTP(rec, httptest.NewRequest(http.MethodPost, "/api/importer/run", strings.NewReader(liveBody)))
	if rec.Code != http.StatusOK {
		t.Fatalf("POST /api/importer/run (live) status = %d, body = %s", rec.Code, rec.Body.String())
	}
	if st.values["ui.theme"] != "dark" || len(cb.combos) != 1 || len(kp.keys) != 1 {
		t.Fatalf("expected live import to populate settings/combos/keys: settings=%v combos=%d keys=%d", st.values, len(cb.combos), len(kp.keys))
	}

	afterInfo, err := os.Stat(legacyPath)
	if err != nil {
		t.Fatalf("stat after: %v", err)
	}
	if !afterInfo.ModTime().Equal(beforeInfo.ModTime()) || afterInfo.Size() != beforeInfo.Size() {
		t.Fatalf("legacy source file was modified during read-only import")
	}

	// History records runs.
	rec = httptest.NewRecorder()
	handler.ServeHTTP(rec, httptest.NewRequest(http.MethodGet, "/api/importer/history", nil))
	if rec.Code != http.StatusOK {
		t.Fatalf("GET /api/importer/history status = %d", rec.Code)
	}
	var hist struct {
		Runs []map[string]any `json:"runs"`
	}
	if err := json.Unmarshal(rec.Body.Bytes(), &hist); err != nil {
		t.Fatalf("decode history: %v", err)
	}
	if len(hist.Runs) < 2 {
		t.Fatalf("expected >= 2 runs in history, got %d", len(hist.Runs))
	}
}
