package settings_test

import (
	"context"
	"encoding/json"
	"errors"
	"io"
	"io/fs"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"path/filepath"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/omeg4monx/router-gateway/internal/modules/settings"
	"github.com/omeg4monx/router-gateway/internal/platform/config"
	"github.com/omeg4monx/router-gateway/internal/platform/db"
	"github.com/omeg4monx/router-gateway/internal/platform/eventbus"
	"github.com/omeg4monx/router-gateway/internal/platform/httpserver"
	"github.com/omeg4monx/router-gateway/internal/platform/lifecycle"
)

type stubSystemModule struct{}

func (stubSystemModule) Name() string                         { return "system" }
func (stubSystemModule) Dependencies() []string               { return nil }
func (stubSystemModule) Migrations() fs.FS                    { return nil }
func (stubSystemModule) RegisterRoutes(lifecycle.RouteMounts) {}
func (stubSystemModule) Start(context.Context) error          { return nil }
func (stubSystemModule) Stop(context.Context) error           { return nil }

type stubAdminAuth struct {
	subject string
}

func (a *stubAdminAuth) AuthenticateAdmin(r *http.Request) (string, error) {
	if r.Header.Get("Authorization") == "Bearer test-admin-token" {
		return a.subject, nil
	}
	return "", errors.New("unauthorized")
}

func (a *stubAdminAuth) RequireLogin(context.Context) bool {
	return true
}

type testHarness struct {
	db     *db.DB
	bus    *eventbus.Bus
	reg    *lifecycle.Registry
	server *httpserver.Server
	mod    *settings.Module
	auth   *stubAdminAuth
}

func newTestHarness(t *testing.T, auth *stubAdminAuth) *testHarness {
	t.Helper()

	cfg := config.Config{
		Host:       "127.0.0.1",
		Port:       20128,
		APITimeout: 30 * time.Second,
	}

	dbPath := filepath.Join(t.TempDir(), "router.sqlite")
	database, err := db.Open(dbPath)
	if err != nil {
		t.Fatalf("db.Open: %v", err)
	}
	t.Cleanup(func() {
		_ = database.Close()
	})

	bus := eventbus.New()
	reg := lifecycle.NewRegistry(nil)
	mod := settings.New(settings.Options{
		DB:  database,
		Bus: bus,
	})

	if err := reg.Register(stubSystemModule{}); err != nil {
		t.Fatalf("reg.Register(system): %v", err)
	}
	if err := reg.Register(mod); err != nil {
		t.Fatalf("reg.Register(settings): %v", err)
	}

	ctx := context.Background()
	if err := reg.MigrateAll(ctx, database); err != nil {
		t.Fatalf("reg.MigrateAll: %v", err)
	}

	srvOpts := httpserver.Options{
		Config: cfg,
		Logger: slog.New(slog.NewTextHandler(io.Discard, nil)),
	}
	if auth != nil {
		srvOpts.Auth = auth
	}
	srv := httpserver.New(srvOpts)
	if err := reg.MountAll(srv.Mounts()); err != nil {
		t.Fatalf("reg.MountAll: %v", err)
	}
	if err := reg.StartAll(ctx); err != nil {
		t.Fatalf("reg.StartAll: %v", err)
	}
	t.Cleanup(func() {
		_ = reg.StopAll(context.Background())
	})

	return &testHarness{
		db:     database,
		bus:    bus,
		reg:    reg,
		server: srv,
		mod:    mod,
		auth:   auth,
	}
}

func (h *testHarness) doJSON(t *testing.T, method, path, body string) *httptest.ResponseRecorder {
	t.Helper()
	var reader io.Reader
	if body != "" {
		reader = strings.NewReader(body)
	}
	req := httptest.NewRequest(method, path, reader)
	if body != "" {
		req.Header.Set("Content-Type", "application/json")
	}
	if h.auth != nil {
		req.Header.Set("Authorization", "Bearer test-admin-token")
	}
	rec := httptest.NewRecorder()
	h.server.Handler().ServeHTTP(rec, req)
	return rec
}

type settingsListResp struct {
	Settings []settings.SettingItem `json:"settings"`
}

type historyListResp struct {
	Items []settings.HistoryItem `json:"items"`
}

func indexByKey(items []settings.SettingItem) map[string]settings.SettingItem {
	out := make(map[string]settings.SettingItem, len(items))
	for _, it := range items {
		out[it.Key] = it
	}
	return out
}

func TestSettingsDefaultsAndPatch(t *testing.T) {
	h := newTestHarness(t, &stubAdminAuth{subject: "operator-1"})

	// Verify migration version 1 is recorded for settings.
	versions, err := h.db.AppliedVersions(context.Background(), "settings")
	if err != nil {
		t.Fatalf("AppliedVersions: %v", err)
	}
	if len(versions) != 1 || versions[0] != 1 {
		t.Fatalf("applied versions = %v, want [1]", versions)
	}

	// 1) GET /api/settings returns all 7 built-in defaults.
	rec := h.doJSON(t, http.MethodGet, "/api/settings", "")
	if rec.Code != http.StatusOK {
		t.Fatalf("GET /api/settings status = %d, want 200 (body=%s)", rec.Code, rec.Body.String())
	}
	var initial settingsListResp
	if err := json.Unmarshal(rec.Body.Bytes(), &initial); err != nil {
		t.Fatalf("unmarshal GET /api/settings: %v", err)
	}
	if len(initial.Settings) != 7 {
		t.Fatalf("settings count = %d, want 7", len(initial.Settings))
	}

	byKey := indexByKey(initial.Settings)
	wantDefaults := map[string]any{
		"ui.theme":                              "dark",
		"ui.locale":                             "en",
		"usage.payload_logging_enabled":         false,
		"usage.payload_retention_days":          float64(7),
		"usage.metadata_retention_days":         float64(90),
		"gateway.default_timeout_seconds":       float64(120),
		"gateway.allow_response_chain_fallback": true,
	}
	for k, want := range wantDefaults {
		item, ok := byKey[k]
		if !ok {
			t.Fatalf("missing built-in setting %q", k)
		}
		if item.Value != want {
			t.Errorf("%s value = %#v (%T), want %#v", k, item.Value, item.Value, want)
		}
		if item.DefaultValue != want {
			t.Errorf("%s default_value = %#v, want %#v", k, item.DefaultValue, want)
		}
		if item.Revision != 0 {
			t.Errorf("%s initial revision = %d, want 0", k, item.Revision)
		}
		if item.Type == "" || item.Description == "" {
			t.Errorf("%s missing metadata type=%q description=%q", k, item.Type, item.Description)
		}
	}

	// 2) PATCH /api/settings updates ui.theme and usage.payload_retention_days.
	rec = h.doJSON(t, http.MethodPatch, "/api/settings", `{
		"updates": {
			"ui.theme": "light",
			"usage.payload_retention_days": 14
		},
		"reason": "increase retention and switch theme"
	}`)
	if rec.Code != http.StatusOK {
		t.Fatalf("PATCH /api/settings status = %d, want 200 (body=%s)", rec.Code, rec.Body.String())
	}
	var patched settingsListResp
	if err := json.Unmarshal(rec.Body.Bytes(), &patched); err != nil {
		t.Fatalf("unmarshal PATCH /api/settings: %v", err)
	}
	patchedByKey := indexByKey(patched.Settings)
	if patchedByKey["ui.theme"].Value != "light" || patchedByKey["ui.theme"].Revision != 1 {
		t.Errorf("PATCH ui.theme = %+v, want value=light revision=1", patchedByKey["ui.theme"])
	}
	if patchedByKey["ui.theme"].UpdatedAt == "" {
		t.Errorf("PATCH ui.theme updated_at should be non-empty")
	}
	if patchedByKey["ui.theme"].UpdatedBy != "operator-1" {
		t.Errorf("PATCH ui.theme updated_by = %q, want operator-1", patchedByKey["ui.theme"].UpdatedBy)
	}
	if patchedByKey["usage.payload_retention_days"].Value != float64(14) || patchedByKey["usage.payload_retention_days"].Revision != 1 {
		t.Errorf("PATCH usage.payload_retention_days = %+v, want value=14 revision=1", patchedByKey["usage.payload_retention_days"])
	}
	if patchedByKey["ui.locale"].Value != "en" || patchedByKey["ui.locale"].Revision != 0 {
		t.Errorf("untouched ui.locale = %+v, want value=en revision=0", patchedByKey["ui.locale"])
	}

	// 3) Subsequent GET reflects updates and incremented revision.
	rec = h.doJSON(t, http.MethodGet, "/api/settings", "")
	if rec.Code != http.StatusOK {
		t.Fatalf("GET /api/settings after PATCH status = %d, want 200", rec.Code)
	}
	var after settingsListResp
	if err := json.Unmarshal(rec.Body.Bytes(), &after); err != nil {
		t.Fatalf("unmarshal GET after PATCH: %v", err)
	}
	afterByKey := indexByKey(after.Settings)
	if afterByKey["ui.theme"].Value != "light" || afterByKey["ui.theme"].Revision != 1 {
		t.Errorf("GET ui.theme = %+v, want value=light revision=1", afterByKey["ui.theme"])
	}
	if afterByKey["usage.payload_retention_days"].Value != float64(14) || afterByKey["usage.payload_retention_days"].Revision != 1 {
		t.Errorf("GET usage.payload_retention_days = %+v, want value=14 revision=1", afterByKey["usage.payload_retention_days"])
	}

	// 4) Second PATCH increments revision to 2.
	rec = h.doJSON(t, http.MethodPatch, "/api/settings", `{"updates":{"ui.theme":"system"}}`)
	if rec.Code != http.StatusOK {
		t.Fatalf("second PATCH status = %d, want 200 (body=%s)", rec.Code, rec.Body.String())
	}
	var second settingsListResp
	if err := json.Unmarshal(rec.Body.Bytes(), &second); err != nil {
		t.Fatalf("unmarshal second PATCH: %v", err)
	}
	secondByKey := indexByKey(second.Settings)
	if secondByKey["ui.theme"].Value != "system" || secondByKey["ui.theme"].Revision != 2 {
		t.Errorf("second PATCH ui.theme = %+v, want value=system revision=2", secondByKey["ui.theme"])
	}
}

func TestSettingsValidationRejectsUnknownAndInvalid(t *testing.T) {
	h := newTestHarness(t, nil)

	cases := []struct {
		name     string
		body     string
		wantCode string
	}{
		{
			name:     "unknown key",
			body:     `{"updates":{"nonexistent.key":"value"}}`,
			wantCode: "invalid_setting_key",
		},
		{
			name:     "out of range int high",
			body:     `{"updates":{"usage.payload_retention_days":999}}`,
			wantCode: "invalid_setting_value",
		},
		{
			name:     "out of range int low",
			body:     `{"updates":{"usage.payload_retention_days":0}}`,
			wantCode: "invalid_setting_value",
		},
		{
			name:     "out of range timeout low",
			body:     `{"updates":{"gateway.default_timeout_seconds":4}}`,
			wantCode: "invalid_setting_value",
		},
		{
			name:     "out of range timeout high",
			body:     `{"updates":{"gateway.default_timeout_seconds":601}}`,
			wantCode: "invalid_setting_value",
		},
		{
			name:     "out of range metadata retention high",
			body:     `{"updates":{"usage.metadata_retention_days":3651}}`,
			wantCode: "invalid_setting_value",
		},
		{
			name:     "float rejected for integer setting",
			body:     `{"updates":{"usage.payload_retention_days":14.5}}`,
			wantCode: "invalid_setting_value",
		},
		{
			name:     "string rejected for integer setting",
			body:     `{"updates":{"usage.payload_retention_days":"14"}}`,
			wantCode: "invalid_setting_value",
		},
		{
			name:     "null rejected for integer setting",
			body:     `{"updates":{"usage.payload_retention_days":null}}`,
			wantCode: "invalid_setting_value",
		},
		{
			name:     "wrong type for enum setting",
			body:     `{"updates":{"ui.theme":123}}`,
			wantCode: "invalid_setting_value",
		},
		{
			name:     "invalid enum option",
			body:     `{"updates":{"ui.theme":"neon"}}`,
			wantCode: "invalid_setting_value",
		},
		{
			name:     "invalid locale option",
			body:     `{"updates":{"ui.locale":"fr"}}`,
			wantCode: "invalid_setting_value",
		},
		{
			name:     "wrong type for boolean setting",
			body:     `{"updates":{"usage.payload_logging_enabled":"true"}}`,
			wantCode: "invalid_setting_value",
		},
		{
			name:     "empty updates map",
			body:     `{"updates":{}}`,
			wantCode: "invalid_request",
		},
		{
			name:     "unknown top-level field",
			body:     `{"updates":{"ui.theme":"light"},"extra":true}`,
			wantCode: "invalid_request",
		},
		{
			name:     "malformed json",
			body:     `{"updates":`,
			wantCode: "invalid_request",
		},
		{
			name:     "trailing json object",
			body:     `{"updates":{"ui.theme":"light"}}{}`,
			wantCode: "invalid_request",
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			rec := h.doJSON(t, http.MethodPatch, "/api/settings", tc.body)
			if rec.Code != http.StatusBadRequest {
				t.Fatalf("PATCH (%s) status = %d, want 400 (body=%s)", tc.name, rec.Code, rec.Body.String())
			}
			var env httpserver.ErrorEnvelope
			if err := json.Unmarshal(rec.Body.Bytes(), &env); err != nil {
				t.Fatalf("unmarshal error envelope: %v (raw=%s)", err, rec.Body.String())
			}
			if env.Error.Code != tc.wantCode {
				t.Errorf("PATCH (%s) error.code = %q, want %q", tc.name, env.Error.Code, tc.wantCode)
			}
		})
	}

	t.Run("multi-key PATCH rolls back atomically on invalid value", func(t *testing.T) {
		rec := h.doJSON(t, http.MethodPatch, "/api/settings", `{
			"updates": {
				"ui.theme": "light",
				"usage.payload_retention_days": 999
			}
		}`)
		if rec.Code != http.StatusBadRequest {
			t.Fatalf("multi-key PATCH status = %d, want 400 (body=%s)", rec.Code, rec.Body.String())
		}
		var env httpserver.ErrorEnvelope
		if err := json.Unmarshal(rec.Body.Bytes(), &env); err != nil {
			t.Fatalf("unmarshal error envelope: %v", err)
		}
		if env.Error.Code != "invalid_setting_value" {
			t.Fatalf("error.code = %q, want invalid_setting_value", env.Error.Code)
		}

		// Verify ui.theme was NOT persisted and history remains empty.
		getRec := h.doJSON(t, http.MethodGet, "/api/settings", "")
		if getRec.Code != http.StatusOK {
			t.Fatalf("GET /api/settings status = %d, want 200", getRec.Code)
		}
		var resp settingsListResp
		if err := json.Unmarshal(getRec.Body.Bytes(), &resp); err != nil {
			t.Fatalf("unmarshal GET /api/settings: %v", err)
		}
		byKey := indexByKey(resp.Settings)
		if byKey["ui.theme"].Value != "dark" || byKey["ui.theme"].Revision != 0 {
			t.Fatalf("ui.theme mutated despite failed multi-key PATCH: %+v", byKey["ui.theme"])
		}

		histRec := h.doJSON(t, http.MethodGet, "/api/settings/history", "")
		if histRec.Code != http.StatusOK {
			t.Fatalf("GET /api/settings/history status = %d, want 200", histRec.Code)
		}
		var hist historyListResp
		if err := json.Unmarshal(histRec.Body.Bytes(), &hist); err != nil {
			t.Fatalf("unmarshal history: %v", err)
		}
		if len(hist.Items) != 0 {
			t.Fatalf("history len = %d, want 0 after failed atomic PATCH", len(hist.Items))
		}
	})

	t.Run("multi-key PATCH rolls back atomically on unknown key", func(t *testing.T) {
		rec := h.doJSON(t, http.MethodPatch, "/api/settings", `{
			"updates": {
				"ui.theme": "light",
				"unknown.key": true
			}
		}`)
		if rec.Code != http.StatusBadRequest {
			t.Fatalf("multi-key unknown PATCH status = %d, want 400 (body=%s)", rec.Code, rec.Body.String())
		}
		var env httpserver.ErrorEnvelope
		if err := json.Unmarshal(rec.Body.Bytes(), &env); err != nil {
			t.Fatalf("unmarshal error envelope: %v", err)
		}
		if env.Error.Code != "invalid_setting_key" {
			t.Fatalf("error.code = %q, want invalid_setting_key", env.Error.Code)
		}

		getRec := h.doJSON(t, http.MethodGet, "/api/settings", "")
		var resp settingsListResp
		if err := json.Unmarshal(getRec.Body.Bytes(), &resp); err != nil {
			t.Fatalf("unmarshal GET /api/settings: %v", err)
		}
		byKey := indexByKey(resp.Settings)
		if byKey["ui.theme"].Value != "dark" || byKey["ui.theme"].Revision != 0 {
			t.Fatalf("ui.theme mutated despite unknown key in PATCH: %+v", byKey["ui.theme"])
		}
	})
}

func TestSettingsHistoryAndRollback(t *testing.T) {
	h := newTestHarness(t, nil)

	// 1) Create revision history entries.
	rec := h.doJSON(t, http.MethodPatch, "/api/settings", `{"updates":{"ui.theme":"light"},"reason":"first theme override"}`)
	if rec.Code != http.StatusOK {
		t.Fatalf("PATCH #1 status = %d (body=%s)", rec.Code, rec.Body.String())
	}
	rec = h.doJSON(t, http.MethodPatch, "/api/settings", `{"updates":{"ui.theme":"system"},"reason":"second theme override"}`)
	if rec.Code != http.StatusOK {
		t.Fatalf("PATCH #2 status = %d (body=%s)", rec.Code, rec.Body.String())
	}
	rec = h.doJSON(t, http.MethodPatch, "/api/settings", `{"updates":{"usage.payload_logging_enabled":true},"reason":"enable payload logs"}`)
	if rec.Code != http.StatusOK {
		t.Fatalf("PATCH #3 status = %d (body=%s)", rec.Code, rec.Body.String())
	}

	// 2) GET /api/settings/history returns newest-first with old/new values.
	histRec := h.doJSON(t, http.MethodGet, "/api/settings/history", "")
	if histRec.Code != http.StatusOK {
		t.Fatalf("GET /api/settings/history status = %d (body=%s)", histRec.Code, histRec.Body.String())
	}
	var hist historyListResp
	if err := json.Unmarshal(histRec.Body.Bytes(), &hist); err != nil {
		t.Fatalf("unmarshal history: %v", err)
	}
	if len(hist.Items) != 3 {
		t.Fatalf("history len = %d, want 3", len(hist.Items))
	}
	if hist.Items[0].SettingKey != "usage.payload_logging_enabled" || hist.Items[0].OldValue != nil || hist.Items[0].NewValue != true {
		t.Errorf("history[0] = %+v, want usage.payload_logging_enabled nil->true", hist.Items[0])
	}
	if hist.Items[1].SettingKey != "ui.theme" || hist.Items[1].OldValue != "light" || hist.Items[1].NewValue != "system" {
		t.Errorf("history[1] = %+v, want ui.theme light->system", hist.Items[1])
	}
	if hist.Items[2].SettingKey != "ui.theme" || hist.Items[2].OldValue != nil || hist.Items[2].NewValue != "light" {
		t.Errorf("history[2] = %+v, want ui.theme nil->light", hist.Items[2])
	}
	if hist.Items[2].ChangedBy != "admin" {
		t.Errorf("default changed_by = %q, want admin", hist.Items[2].ChangedBy)
	}

	secondThemeHistID := hist.Items[1].ID
	firstThemeHistID := hist.Items[2].ID

	// 3) Query filters (key & limit) and invalid limit rejection.
	filteredRec := h.doJSON(t, http.MethodGet, "/api/settings/history?key=ui.theme&limit=1", "")
	if filteredRec.Code != http.StatusOK {
		t.Fatalf("filtered history status = %d (body=%s)", filteredRec.Code, filteredRec.Body.String())
	}
	var filtered historyListResp
	if err := json.Unmarshal(filteredRec.Body.Bytes(), &filtered); err != nil {
		t.Fatalf("unmarshal filtered history: %v", err)
	}
	if len(filtered.Items) != 1 || filtered.Items[0].ID != secondThemeHistID {
		t.Fatalf("filtered history = %+v, want single item id=%d", filtered.Items, secondThemeHistID)
	}

	badLimitRec := h.doJSON(t, http.MethodGet, "/api/settings/history?limit=abc", "")
	if badLimitRec.Code != http.StatusBadRequest {
		t.Fatalf("invalid limit status = %d, want 400", badLimitRec.Code)
	}

	// 4) Rollback second theme override (old_value = "light") -> restores "light" at revision 3.
	rbRec := h.doJSON(t, http.MethodPost, "/api/settings/rollback", `{"history_id":`+strconv.FormatInt(secondThemeHistID, 10)+`,"reason":"revert to light"}`)
	if rbRec.Code != http.StatusOK {
		t.Fatalf("POST /api/settings/rollback status = %d, want 200 (body=%s)", rbRec.Code, rbRec.Body.String())
	}
	var restored settings.SettingItem
	if err := json.Unmarshal(rbRec.Body.Bytes(), &restored); err != nil {
		t.Fatalf("unmarshal rollback response: %v", err)
	}
	if restored.Key != "ui.theme" || restored.Value != "light" || restored.Revision != 3 {
		t.Fatalf("restored item = %+v, want ui.theme=light rev=3", restored)
	}

	// 5) Rollback first theme override (old_value = NULL) -> restores built-in default "dark" at revision 4.
	rbDefaultRec := h.doJSON(t, http.MethodPost, "/api/settings/rollback", `{"history_id":`+strconv.FormatInt(firstThemeHistID, 10)+`}`)
	if rbDefaultRec.Code != http.StatusOK {
		t.Fatalf("POST /api/settings/rollback (default) status = %d, want 200 (body=%s)", rbDefaultRec.Code, rbDefaultRec.Body.String())
	}
	var restoredDefault settings.SettingItem
	if err := json.Unmarshal(rbDefaultRec.Body.Bytes(), &restoredDefault); err != nil {
		t.Fatalf("unmarshal rollback default response: %v", err)
	}
	if restoredDefault.Key != "ui.theme" || restoredDefault.Value != "dark" || restoredDefault.Revision != 4 {
		t.Fatalf("restored default item = %+v, want ui.theme=dark rev=4", restoredDefault)
	}

	// Confirm history now contains both rollback audit entries at the top.
	histAfterRec := h.doJSON(t, http.MethodGet, "/api/settings/history?key=ui.theme", "")
	if histAfterRec.Code != http.StatusOK {
		t.Fatalf("GET history after rollback status = %d", histAfterRec.Code)
	}
	var histAfter historyListResp
	if err := json.Unmarshal(histAfterRec.Body.Bytes(), &histAfter); err != nil {
		t.Fatalf("unmarshal history after rollback: %v", err)
	}
	if len(histAfter.Items) != 4 {
		t.Fatalf("ui.theme history len = %d, want 4", len(histAfter.Items))
	}
	if histAfter.Items[0].OldValue != "light" || histAfter.Items[0].NewValue != "dark" || histAfter.Items[0].Revision != 4 {
		t.Errorf("latest rollback audit entry = %+v, want light->dark rev=4", histAfter.Items[0])
	}
	if histAfter.Items[1].OldValue != "system" || histAfter.Items[1].NewValue != "light" || histAfter.Items[1].Revision != 3 || histAfter.Items[1].Reason != "revert to light" {
		t.Errorf("first rollback audit entry = %+v, want system->light rev=3 reason='revert to light'", histAfter.Items[1])
	}

	// 6) Rolling back nonexistent history_id returns 404 revision_not_found without mutating state.
	notFoundRec := h.doJSON(t, http.MethodPost, "/api/settings/rollback", `{"history_id":99999}`)
	if notFoundRec.Code != http.StatusNotFound {
		t.Fatalf("rollback nonexistent status = %d, want 404 (body=%s)", notFoundRec.Code, notFoundRec.Body.String())
	}
	var env httpserver.ErrorEnvelope
	if err := json.Unmarshal(notFoundRec.Body.Bytes(), &env); err != nil {
		t.Fatalf("unmarshal 404 envelope: %v", err)
	}
	if env.Error.Code != "revision_not_found" {
		t.Fatalf("404 error.code = %q, want revision_not_found", env.Error.Code)
	}

	// Verify state was untouched by the 404 rollback.
	checkRec := h.doJSON(t, http.MethodGet, "/api/settings/history?key=ui.theme", "")
	var histCheck historyListResp
	if err := json.Unmarshal(checkRec.Body.Bytes(), &histCheck); err != nil {
		t.Fatalf("unmarshal history check: %v", err)
	}
	if len(histCheck.Items) != 4 {
		t.Fatalf("history len changed after 404 rollback: got %d, want 4", len(histCheck.Items))
	}

	// 7) Strict JSON checks on rollback body.
	for _, tc := range []struct {
		name string
		body string
	}{
		{"missing history_id", `{"reason":"missing id"}`},
		{"unknown field", `{"history_id":1,"extra":true}`},
		{"float history_id", `{"history_id":1.5}`},
		{"string history_id", `{"history_id":"1"}`},
	} {
		t.Run(tc.name, func(t *testing.T) {
			rec := h.doJSON(t, http.MethodPost, "/api/settings/rollback", tc.body)
			if rec.Code != http.StatusBadRequest {
				t.Fatalf("rollback (%s) status = %d, want 400 (body=%s)", tc.name, rec.Code, rec.Body.String())
			}
		})
	}
}

func TestSettingsPublishesEventBus(t *testing.T) {
	h := newTestHarness(t, &stubAdminAuth{subject: "ops-admin"})

	events, unsub := h.bus.Subscribe(settings.TopicSettingsUpdated, 8)
	defer unsub()

	// Invalid PATCH must not emit an event.
	badRec := h.doJSON(t, http.MethodPatch, "/api/settings", `{"updates":{"unknown.key":1}}`)
	if badRec.Code != http.StatusBadRequest {
		t.Fatalf("bad PATCH status = %d, want 400", badRec.Code)
	}
	select {
	case ev := <-events:
		t.Fatalf("unexpected event after failed PATCH: %+v", ev)
	default:
	}

	// Valid PATCH emits settings.updated event.
	rec := h.doJSON(t, http.MethodPatch, "/api/settings", `{"updates":{"ui.theme":"light","gateway.allow_response_chain_fallback":false},"reason":"ops change"}`)
	if rec.Code != http.StatusOK {
		t.Fatalf("PATCH status = %d, want 200 (body=%s)", rec.Code, rec.Body.String())
	}

	select {
	case ev := <-events:
		if ev.Topic != settings.TopicSettingsUpdated {
			t.Fatalf("event topic = %q, want %q", ev.Topic, settings.TopicSettingsUpdated)
		}
		payload, ok := ev.Payload.(settings.UpdatedEvent)
		if !ok {
			t.Fatalf("event payload type = %T, want settings.UpdatedEvent", ev.Payload)
		}
		if len(payload.Keys) != 2 || payload.Keys[0] != "gateway.allow_response_chain_fallback" || payload.Keys[1] != "ui.theme" {
			t.Errorf("event keys = %v, want [gateway.allow_response_chain_fallback ui.theme]", payload.Keys)
		}
		if payload.Values["ui.theme"] != "light" || payload.Values["gateway.allow_response_chain_fallback"] != false {
			t.Errorf("event values = %#v", payload.Values)
		}
		if payload.ChangedBy != "ops-admin" || payload.Reason != "ops change" || payload.Rollback {
			t.Errorf("event metadata = %+v", payload)
		}
	case <-time.After(time.Second):
		t.Fatalf("timed out waiting for settings.updated event on PATCH")
	}

	// Rollback also emits settings.updated event.
	histRec := h.doJSON(t, http.MethodGet, "/api/settings/history?key=ui.theme", "")
	var hist historyListResp
	if err := json.Unmarshal(histRec.Body.Bytes(), &hist); err != nil {
		t.Fatalf("unmarshal history: %v", err)
	}
	rbRec := h.doJSON(t, http.MethodPost, "/api/settings/rollback", `{"history_id":`+strconv.FormatInt(hist.Items[0].ID, 10)+`}`)
	if rbRec.Code != http.StatusOK {
		t.Fatalf("rollback status = %d, want 200 (body=%s)", rbRec.Code, rbRec.Body.String())
	}

	select {
	case ev := <-events:
		payload, ok := ev.Payload.(settings.UpdatedEvent)
		if !ok {
			t.Fatalf("rollback event payload type = %T, want settings.UpdatedEvent", ev.Payload)
		}
		if !payload.Rollback || len(payload.Keys) != 1 || payload.Keys[0] != "ui.theme" || payload.Values["ui.theme"] != "dark" {
			t.Errorf("unexpected rollback event payload: %+v", payload)
		}
	case <-time.After(time.Second):
		t.Fatalf("timed out waiting for settings.updated event on rollback")
	}
}

func TestSettingsTimestampsFixedWidthUTC(t *testing.T) {
	const fixedLayout = "2006-01-02T15:04:05.000000000Z"

	assertFixedWidthUTC := func(t *testing.T, label, ts string) {
		t.Helper()
		if len(ts) != len(fixedLayout) {
			t.Fatalf("%s = %q (len=%d), want fixed width len=%d", label, ts, len(ts), len(fixedLayout))
		}
		if !strings.HasSuffix(ts, "Z") {
			t.Fatalf("%s = %q, want UTC suffix Z", label, ts)
		}
		if ts[19] != '.' {
			t.Fatalf("%s = %q, want fractional dot at index 19", label, ts)
		}
		if _, err := time.Parse(time.RFC3339Nano, ts); err != nil {
			t.Fatalf("%s = %q is not valid RFC3339: %v", label, ts, err)
		}
	}

	h := newTestHarness(t, nil)

	rec := h.doJSON(t, http.MethodPatch, "/api/settings", `{"updates":{"ui.theme":"light"},"reason":"check timestamp"}`)
	if rec.Code != http.StatusOK {
		t.Fatalf("PATCH status = %d (body=%s)", rec.Code, rec.Body.String())
	}
	var patched settingsListResp
	if err := json.Unmarshal(rec.Body.Bytes(), &patched); err != nil {
		t.Fatalf("unmarshal PATCH: %v", err)
	}
	byKey := indexByKey(patched.Settings)
	assertFixedWidthUTC(t, "settings.updated_at", byKey["ui.theme"].UpdatedAt)

	histRec := h.doJSON(t, http.MethodGet, "/api/settings/history?key=ui.theme", "")
	if histRec.Code != http.StatusOK {
		t.Fatalf("GET history status = %d (body=%s)", histRec.Code, histRec.Body.String())
	}
	var hist historyListResp
	if err := json.Unmarshal(histRec.Body.Bytes(), &hist); err != nil {
		t.Fatalf("unmarshal history: %v", err)
	}
	if len(hist.Items) != 1 {
		t.Fatalf("history len = %d, want 1", len(hist.Items))
	}
	assertFixedWidthUTC(t, "history.created_at (patch)", hist.Items[0].CreatedAt)

	rbRec := h.doJSON(t, http.MethodPost, "/api/settings/rollback", `{"history_id":`+strconv.FormatInt(hist.Items[0].ID, 10)+`}`)
	if rbRec.Code != http.StatusOK {
		t.Fatalf("rollback status = %d (body=%s)", rbRec.Code, rbRec.Body.String())
	}
	var rolled settings.SettingItem
	if err := json.Unmarshal(rbRec.Body.Bytes(), &rolled); err != nil {
		t.Fatalf("unmarshal rollback: %v", err)
	}
	assertFixedWidthUTC(t, "rollback.updated_at", rolled.UpdatedAt)

	histAfterRec := h.doJSON(t, http.MethodGet, "/api/settings/history?key=ui.theme", "")
	if histAfterRec.Code != http.StatusOK {
		t.Fatalf("GET history after rollback status = %d", histAfterRec.Code)
	}
	var histAfter historyListResp
	if err := json.Unmarshal(histAfterRec.Body.Bytes(), &histAfter); err != nil {
		t.Fatalf("unmarshal history after rollback: %v", err)
	}
	if len(histAfter.Items) != 2 {
		t.Fatalf("history len = %d, want 2", len(histAfter.Items))
	}
	assertFixedWidthUTC(t, "history.created_at (rollback)", histAfter.Items[0].CreatedAt)
	if histAfter.Items[0].CreatedAt < histAfter.Items[1].CreatedAt {
		t.Fatalf("expected newest-first lexicographical order: %q < %q", histAfter.Items[0].CreatedAt, histAfter.Items[1].CreatedAt)
	}
}
