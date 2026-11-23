package cloudsync_test

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"path/filepath"
	"strings"
	"sync"
	"testing"

	"github.com/go-chi/chi/v5"

	"github.com/omeg4monx/router-gateway/internal/contracts"
	"github.com/omeg4monx/router-gateway/internal/modules/cloudsync"
	"github.com/omeg4monx/router-gateway/internal/platform/crypto"
	"github.com/omeg4monx/router-gateway/internal/platform/db"
	"github.com/omeg4monx/router-gateway/internal/platform/eventbus"
	"github.com/omeg4monx/router-gateway/internal/platform/lifecycle"
)

type fakeSettingsStore struct {
	mu     sync.Mutex
	values map[string]any
}

func (f *fakeSettingsStore) ExportSettings(context.Context) (map[string]any, error) {
	f.mu.Lock()
	defer f.mu.Unlock()
	out := make(map[string]any, len(f.values))
	for k, v := range f.values {
		out[k] = v
	}
	return out, nil
}

func (f *fakeSettingsStore) ImportSettings(_ context.Context, values map[string]any, _, _ string) (int, error) {
	f.mu.Lock()
	defer f.mu.Unlock()
	if f.values == nil {
		f.values = make(map[string]any)
	}
	count := 0
	for k, v := range values {
		f.values[k] = v
		count++
	}
	return count, nil
}

type fakeComboStore struct {
	mu      sync.Mutex
	combos  []contracts.ComboSnapshot
	aliases []contracts.AliasSnapshot
}

func (f *fakeComboStore) ExportCombos(context.Context) ([]contracts.ComboSnapshot, []contracts.AliasSnapshot, error) {
	f.mu.Lock()
	defer f.mu.Unlock()
	c := append([]contracts.ComboSnapshot(nil), f.combos...)
	a := append([]contracts.AliasSnapshot(nil), f.aliases...)
	return c, a, nil
}

func (f *fakeComboStore) ImportCombos(_ context.Context, combos []contracts.ComboSnapshot, aliases []contracts.AliasSnapshot, overwrite bool, _ string) (int, int, error) {
	f.mu.Lock()
	defer f.mu.Unlock()
	existingC := make(map[string]int)
	for i, c := range f.combos {
		existingC[c.Name] = i
	}
	cImported := 0
	for _, c := range combos {
		if idx, ok := existingC[c.Name]; ok {
			if overwrite {
				f.combos[idx] = c
				cImported++
			}
			continue
		}
		f.combos = append(f.combos, c)
		existingC[c.Name] = len(f.combos) - 1
		cImported++
	}

	existingA := make(map[string]int)
	for i, a := range f.aliases {
		existingA[a.Alias] = i
	}
	aImported := 0
	for _, a := range aliases {
		if idx, ok := existingA[a.Alias]; ok {
			if overwrite {
				f.aliases[idx] = a
				aImported++
			}
			continue
		}
		f.aliases = append(f.aliases, a)
		existingA[a.Alias] = len(f.aliases) - 1
		aImported++
	}
	return cImported, aImported, nil
}

func setupCloudSyncTest(t *testing.T, allowPrivate bool) (*cloudsync.Module, http.Handler, *fakeSettingsStore, *fakeComboStore) {
	t.Helper()
	database, err := db.Open(filepath.Join(t.TempDir(), "cloudsync.sqlite"))
	if err != nil {
		t.Fatalf("db.Open: %v", err)
	}
	t.Cleanup(func() { _ = database.Close() })

	key, err := crypto.LoadOrCreateMasterKey(t.TempDir(), "")
	if err != nil {
		t.Fatalf("LoadOrCreateMasterKey: %v", err)
	}
	vault, err := crypto.NewVault(key)
	if err != nil {
		t.Fatalf("NewVault: %v", err)
	}

	settingsStore := &fakeSettingsStore{
		values: map[string]any{
			"ui.theme":  "dark",
			"ui.locale": "en",
		},
	}
	comboStore := &fakeComboStore{
		combos: []contracts.ComboSnapshot{
			{Name: "primary", Strategy: "ordered", Targets: []string{"openai/gpt-4o"}, Enabled: true},
		},
		aliases: []contracts.AliasSnapshot{
			{Alias: "fast", TargetModel: "openai/gpt-4o-mini"},
		},
	}

	mod := cloudsync.New(cloudsync.Options{
		DB:                    database,
		Vault:                 vault,
		Bus:                   eventbus.New(),
		Settings:              settingsStore,
		Combos:                comboStore,
		AllowPrivateEndpoints: allowPrivate,
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
	return mod, r, settingsStore, comboStore
}

func TestCloudSyncSSRFGuardAndVaultedTokenIsolation(t *testing.T) {
	_, handler, _, _ := setupCloudSyncTest(t, false)

	// Private/loopback/metadata endpoints must be rejected by default.
	for _, badURL := range []string{
		"http://127.0.0.1:8080/sync",
		"http://localhost:9000/sync",
		"http://169.254.169.254/latest/meta-data/",
		"ftp://example.com/sync",
	} {
		body := `{"endpoint_url":"` + badURL + `"}`
		rec := httptest.NewRecorder()
		handler.ServeHTTP(rec, httptest.NewRequest(http.MethodPatch, "/api/cloudsync/config", strings.NewReader(body)))
		if rec.Code != http.StatusBadRequest {
			t.Fatalf("expected 400 for SSRF URL %q, got %d (%s)", badURL, rec.Code, rec.Body.String())
		}
	}

	// Valid public HTTPS endpoint + secret token never returns raw token.
	rawSecret := "sync_super_secret_token_98765"
	patch := `{"enabled":true,"endpoint_url":"https://sync.example.com/v1/state","auth_token":"` + rawSecret + `","conflict_policy":"merge"}`
	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, httptest.NewRequest(http.MethodPatch, "/api/cloudsync/config", strings.NewReader(patch)))
	if rec.Code != http.StatusOK {
		t.Fatalf("PATCH config status = %d, body = %s", rec.Code, rec.Body.String())
	}
	if strings.Contains(rec.Body.String(), rawSecret) {
		t.Fatalf("PATCH response leaked raw auth_token: %s", rec.Body.String())
	}

	rec = httptest.NewRecorder()
	handler.ServeHTTP(rec, httptest.NewRequest(http.MethodGet, "/api/cloudsync/config", nil))
	if rec.Code != http.StatusOK {
		t.Fatalf("GET config status = %d", rec.Code)
	}
	if strings.Contains(rec.Body.String(), rawSecret) {
		t.Fatalf("GET response leaked raw auth_token: %s", rec.Body.String())
	}
}

func TestCloudSyncPushPullAndSnapshotConflictPolicies(t *testing.T) {
	var (
		mu         sync.Mutex
		storedBody []byte
		gotAuth    string
	)
	remote := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		mu.Lock()
		defer mu.Unlock()
		gotAuth = r.Header.Get("Authorization")
		switch r.Method {
		case http.MethodPut:
			b, _ := io.ReadAll(r.Body)
			storedBody = b
			w.WriteHeader(http.StatusOK)
			_, _ = w.Write([]byte(`{"ok":true}`))
		case http.MethodGet:
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			_, _ = w.Write(storedBody)
		default:
			w.WriteHeader(http.StatusMethodNotAllowed)
		}
	}))
	defer remote.Close()

	_, handler, settingsStore, comboStore := setupCloudSyncTest(t, true)

	// Configure remote pointing to test server.
	patch := `{"enabled":true,"endpoint_url":"` + remote.URL + `","auth_token":"tok_abc12345","conflict_policy":"remote_wins"}`
	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, httptest.NewRequest(http.MethodPatch, "/api/cloudsync/config", strings.NewReader(patch)))
	if rec.Code != http.StatusOK {
		t.Fatalf("PATCH config status = %d, body = %s", rec.Code, rec.Body.String())
	}

	// Push local snapshot.
	rec = httptest.NewRecorder()
	handler.ServeHTTP(rec, httptest.NewRequest(http.MethodPost, "/api/cloudsync/push", nil))
	if rec.Code != http.StatusOK {
		t.Fatalf("POST push status = %d, body = %s", rec.Code, rec.Body.String())
	}
	mu.Lock()
	authHeader := gotAuth
	pushedLen := len(storedBody)
	mu.Unlock()
	if authHeader != "Bearer tok_abc12345" {
		t.Fatalf("expected Bearer auth header, got %q", authHeader)
	}
	if pushedLen == 0 {
		t.Fatalf("expected non-empty pushed snapshot")
	}

	// Mutate local state and verify local_wins preserves local theme while remote_wins restores remote.
	settingsStore.mu.Lock()
	settingsStore.values["ui.theme"] = "light"
	settingsStore.mu.Unlock()
	comboStore.mu.Lock()
	comboStore.combos = nil
	comboStore.mu.Unlock()

	// Pull with local_wins -> keeps ui.theme=light, restores missing combo "primary".
	rec = httptest.NewRecorder()
	handler.ServeHTTP(rec, httptest.NewRequest(http.MethodPatch, "/api/cloudsync/config", strings.NewReader(`{"conflict_policy":"local_wins"}`)))
	if rec.Code != http.StatusOK {
		t.Fatalf("PATCH local_wins status = %d", rec.Code)
	}
	rec = httptest.NewRecorder()
	handler.ServeHTTP(rec, httptest.NewRequest(http.MethodPost, "/api/cloudsync/pull", nil))
	if rec.Code != http.StatusOK {
		t.Fatalf("POST pull (local_wins) status = %d, body = %s", rec.Code, rec.Body.String())
	}
	if got := settingsStore.values["ui.theme"]; got != "light" {
		t.Fatalf("local_wins should keep ui.theme=light, got %v", got)
	}
	if len(comboStore.combos) != 1 {
		t.Fatalf("local_wins should still add missing combo, got %d", len(comboStore.combos))
	}

	// Pull with remote_wins -> overwrites ui.theme back to dark.
	rec = httptest.NewRecorder()
	handler.ServeHTTP(rec, httptest.NewRequest(http.MethodPatch, "/api/cloudsync/config", strings.NewReader(`{"conflict_policy":"remote_wins"}`)))
	if rec.Code != http.StatusOK {
		t.Fatalf("PATCH remote_wins status = %d", rec.Code)
	}
	rec = httptest.NewRecorder()
	handler.ServeHTTP(rec, httptest.NewRequest(http.MethodPost, "/api/cloudsync/pull", nil))
	if rec.Code != http.StatusOK {
		t.Fatalf("POST pull (remote_wins) status = %d, body = %s", rec.Code, rec.Body.String())
	}
	if got := settingsStore.values["ui.theme"]; got != "dark" {
		t.Fatalf("remote_wins should restore ui.theme=dark, got %v", got)
	}

	// History should have recorded push + 2 pulls.
	rec = httptest.NewRecorder()
	handler.ServeHTTP(rec, httptest.NewRequest(http.MethodGet, "/api/cloudsync/history", nil))
	if rec.Code != http.StatusOK {
		t.Fatalf("GET history status = %d", rec.Code)
	}
	var hist struct {
		Events []map[string]any `json:"events"`
	}
	if err := json.Unmarshal(rec.Body.Bytes(), &hist); err != nil {
		t.Fatalf("decode history: %v", err)
	}
	if len(hist.Events) < 3 {
		t.Fatalf("expected >= 3 sync history events, got %d", len(hist.Events))
	}
}
