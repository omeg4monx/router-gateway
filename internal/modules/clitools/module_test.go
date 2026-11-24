package clitools_test

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"io"
	"io/fs"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"testing"
	"time"

	"github.com/omeg4monx/router-gateway/internal/modules/clitools"
	"github.com/omeg4monx/router-gateway/internal/platform/config"
	"github.com/omeg4monx/router-gateway/internal/platform/httpserver"
	"github.com/omeg4monx/router-gateway/internal/platform/lifecycle"
)

var _ lifecycle.Module = (*clitools.Module)(nil)

type stubModule struct {
	name string
}

func (s stubModule) Name() string                         { return s.name }
func (s stubModule) Dependencies() []string               { return nil }
func (s stubModule) Migrations() fs.FS                    { return nil }
func (s stubModule) RegisterRoutes(lifecycle.RouteMounts) {}
func (s stubModule) Start(context.Context) error          { return nil }
func (s stubModule) Stop(context.Context) error           { return nil }

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
	homeDir string
	reg     *lifecycle.Registry
	server  *httpserver.Server
	mod     *clitools.Module
}

func newTestHarness(t *testing.T, host string, withAuth bool) *testHarness {
	t.Helper()

	homeDir := t.TempDir()
	cfg := config.Config{
		Host:       host,
		Port:       20128,
		APITimeout: 30 * time.Second,
	}

	mod := clitools.New(clitools.Options{
		Config:  cfg,
		HomeDir: homeDir,
	})

	reg := lifecycle.NewRegistry(nil)
	if err := reg.Register(stubModule{name: "system"}); err != nil {
		t.Fatalf("reg.Register(system): %v", err)
	}
	if err := reg.Register(mod); err != nil {
		t.Fatalf("reg.Register(clitools): %v", err)
	}

	srvOpts := httpserver.Options{
		Config: cfg,
		Logger: slog.New(slog.NewTextHandler(io.Discard, nil)),
	}
	if withAuth {
		srvOpts.Auth = &stubAdminAuth{subject: "admin"}
	}
	srv := httpserver.New(srvOpts)
	if err := reg.MountAll(srv.Mounts()); err != nil {
		t.Fatalf("reg.MountAll: %v", err)
	}
	if err := reg.StartAll(context.Background()); err != nil {
		t.Fatalf("reg.StartAll: %v", err)
	}
	t.Cleanup(func() {
		_ = reg.StopAll(context.Background())
	})

	return &testHarness{
		homeDir: homeDir,
		reg:     reg,
		server:  srv,
		mod:     mod,
	}
}

func (h *testHarness) doJSON(t *testing.T, method, path string, body any, authenticated bool) *httptest.ResponseRecorder {
	t.Helper()

	var reader io.Reader
	if body != nil {
		switch v := body.(type) {
		case string:
			reader = strings.NewReader(v)
		default:
			raw, err := json.Marshal(v)
			if err != nil {
				t.Fatalf("json.Marshal: %v", err)
			}
			reader = bytes.NewReader(raw)
		}
	}

	req := httptest.NewRequest(method, path, reader)
	if body != nil {
		req.Header.Set("Content-Type", "application/json")
	}
	if authenticated {
		req.Header.Set("Authorization", "Bearer test-admin-token")
	}
	rec := httptest.NewRecorder()
	h.server.Handler().ServeHTTP(rec, req)
	return rec
}

func decodeJSON[T any](t *testing.T, rec *httptest.ResponseRecorder) T {
	t.Helper()
	var out T
	if err := json.Unmarshal(rec.Body.Bytes(), &out); err != nil {
		t.Fatalf("json.Unmarshal(%s): %v", rec.Body.String(), err)
	}
	return out
}

type catalogResponse struct {
	LocalApplyEnabled bool `json:"local_apply_enabled"`
	Tools             []struct {
		Tool        string `json:"tool"`
		Name        string `json:"name"`
		Format      string `json:"format"`
		Filename    string `json:"filename"`
		TargetPath  string `json:"target_path"`
		Description string `json:"description"`
	} `json:"tools"`
}

type generateResponse struct {
	Tool         string `json:"tool"`
	Name         string `json:"name"`
	Format       string `json:"format"`
	Filename     string `json:"filename"`
	TargetPath   string `json:"target_path"`
	Content      string `json:"content"`
	Instructions string `json:"instructions"`
}

type applyResponse struct {
	Applied    bool   `json:"applied"`
	Tool       string `json:"tool"`
	TargetPath string `json:"target_path"`
	BackupPath string `json:"backup_path"`
}

type rollbackResponse struct {
	Restored   bool   `json:"restored"`
	Tool       string `json:"tool"`
	TargetPath string `json:"target_path"`
	BackupPath string `json:"backup_path"`
}

func TestCLIToolsModuleContract(t *testing.T) {
	h := newTestHarness(t, "127.0.0.1", true)

	if got := h.mod.Name(); got != "clitools" {
		t.Fatalf("Name() = %q, want clitools", got)
	}
	if got := h.mod.Dependencies(); !reflect.DeepEqual(got, []string{"system"}) {
		t.Fatalf("Dependencies() = %v, want [system]", got)
	}
	if got := h.mod.Migrations(); got != nil {
		t.Fatalf("Migrations() = %v, want nil", got)
	}

	// Admin endpoints must enforce admin authentication when configured.
	unauthRec := h.doJSON(t, http.MethodGet, "/api/cli-tools", nil, false)
	if unauthRec.Code != http.StatusUnauthorized {
		t.Fatalf("unauthenticated GET /api/cli-tools status = %d, want 401", unauthRec.Code)
	}
}

func TestCLIToolsCatalogAndGenerateAllTools(t *testing.T) {
	h := newTestHarness(t, "127.0.0.1", false)

	catRec := h.doJSON(t, http.MethodGet, "/api/cli-tools", nil, false)
	if catRec.Code != http.StatusOK {
		t.Fatalf("GET /api/cli-tools status = %d, body = %s", catRec.Code, catRec.Body.String())
	}
	cat := decodeJSON[catalogResponse](t, catRec)
	if !cat.LocalApplyEnabled {
		t.Fatalf("local_apply_enabled = false on loopback, want true")
	}

	expectedTargets := map[string]struct {
		format string
		rel    string
	}{
		"claude-code": {format: "json", rel: filepath.Join(".claude", "settings.json")},
		"codex-cli":   {format: "toml", rel: filepath.Join(".codex", "config.toml")},
		"opencode":    {format: "json", rel: filepath.Join(".config", "opencode", "opencode.json")},
		"cline":       {format: "json", rel: filepath.Join(".cline", "omeg4mon-router.json")},
		"roo-code":    {format: "json", rel: filepath.Join(".roo-code", "omeg4mon-router.json")},
		"droid":       {format: "json", rel: filepath.Join(".factory", "config.json")},
	}

	if len(cat.Tools) != len(expectedTargets) {
		t.Fatalf("len(tools) = %d, want %d", len(cat.Tools), len(expectedTargets))
	}

	seen := make(map[string]bool, len(cat.Tools))
	for _, tool := range cat.Tools {
		seen[tool.Tool] = true
		exp, ok := expectedTargets[tool.Tool]
		if !ok {
			t.Fatalf("unexpected tool in catalog: %q", tool.Tool)
		}
		if tool.Format != exp.format {
			t.Fatalf("tool %q format = %q, want %q", tool.Tool, tool.Format, exp.format)
		}
		wantPath := filepath.Join(h.homeDir, exp.rel)
		if tool.TargetPath != wantPath {
			t.Fatalf("tool %q target_path = %q, want %q", tool.Tool, tool.TargetPath, wantPath)
		}
		if tool.Name == "" || tool.Filename == "" || tool.Description == "" {
			t.Fatalf("tool %q has empty metadata: %+v", tool.Tool, tool)
		}
	}
	for id := range expectedTargets {
		if !seen[id] {
			t.Fatalf("missing tool %q from catalog", id)
		}
	}

	// Generate config for all 6 tools.
	for id, exp := range expectedTargets {
		rec := h.doJSON(t, http.MethodPost, "/api/cli-tools/"+id+"/generate", map[string]string{
			"base_url": "http://127.0.0.1:20128/",
			"api_key":  `om4_secret_"quoted"\key`,
			"model":    "claude-sonnet-4",
		}, false)
		if rec.Code != http.StatusOK {
			t.Fatalf("POST /api/cli-tools/%s/generate status = %d, body = %s", id, rec.Code, rec.Body.String())
		}
		resp := decodeJSON[generateResponse](t, rec)
		if resp.Tool != id {
			t.Fatalf("generate %q tool = %q", id, resp.Tool)
		}
		if resp.Format != exp.format {
			t.Fatalf("generate %q format = %q, want %q", id, resp.Format, exp.format)
		}
		if resp.TargetPath != filepath.Join(h.homeDir, exp.rel) {
			t.Fatalf("generate %q target_path = %q", id, resp.TargetPath)
		}
		if resp.Instructions == "" {
			t.Fatalf("generate %q instructions must not be empty", id)
		}

		switch exp.format {
		case "json":
			if !json.Valid([]byte(resp.Content)) {
				t.Fatalf("tool %q produced invalid JSON:\n%s", id, resp.Content)
			}
			var parsed map[string]any
			if err := json.Unmarshal([]byte(resp.Content), &parsed); err != nil {
				t.Fatalf("tool %q json.Unmarshal: %v", id, err)
			}
			if id == "claude-code" {
				env, ok := parsed["env"].(map[string]any)
				if !ok {
					t.Fatalf("claude-code missing env object: %s", resp.Content)
				}
				if env["ANTHROPIC_BASE_URL"] != "http://127.0.0.1:20128" {
					t.Fatalf("ANTHROPIC_BASE_URL = %v", env["ANTHROPIC_BASE_URL"])
				}
				if env["ANTHROPIC_AUTH_TOKEN"] != `om4_secret_"quoted"\key` {
					t.Fatalf("ANTHROPIC_AUTH_TOKEN = %v", env["ANTHROPIC_AUTH_TOKEN"])
				}
				if env["ANTHROPIC_MODEL"] != "claude-sonnet-4" {
					t.Fatalf("ANTHROPIC_MODEL = %v", env["ANTHROPIC_MODEL"])
				}
			} else {
				if !strings.Contains(resp.Content, "http://127.0.0.1:20128/v1") {
					t.Fatalf("tool %q expected /v1 base URL in content:\n%s", id, resp.Content)
				}
			}
		case "toml":
			if !strings.Contains(resp.Content, `[model_providers.omeg4mon-router]`) {
				t.Fatalf("codex-cli TOML missing provider table:\n%s", resp.Content)
			}
			if !strings.Contains(resp.Content, `base_url = "http://127.0.0.1:20128/v1"`) {
				t.Fatalf("codex-cli TOML missing /v1 base_url:\n%s", resp.Content)
			}
			if !strings.Contains(resp.Content, `"om4_secret_\"quoted\"\\key"`) {
				t.Fatalf("codex-cli TOML did not escape string cleanly with strconv.Quote:\n%s", resp.Content)
			}
		}
	}

	// Unknown tool -> 404.
	unknownRec := h.doJSON(t, http.MethodPost, "/api/cli-tools/nonexistent/generate", map[string]string{
		"base_url": "http://127.0.0.1:20128",
	}, false)
	if unknownRec.Code != http.StatusNotFound {
		t.Fatalf("unknown tool generate status = %d, want 404", unknownRec.Code)
	}
	unknownErr := decodeJSON[httpserver.ErrorEnvelope](t, unknownRec)
	if unknownErr.Error.Code != "tool_not_found" {
		t.Fatalf("unknown tool error code = %q, want tool_not_found", unknownErr.Error.Code)
	}

	// Non-HTTP base_url -> 400 invalid_base_url.
	for _, badURL := range []string{"", "ftp://127.0.0.1:20128", "javascript:alert(1)", "http://", "/relative/path"} {
		badRec := h.doJSON(t, http.MethodPost, "/api/cli-tools/claude-code/generate", map[string]string{
			"base_url": badURL,
		}, false)
		if badRec.Code != http.StatusBadRequest {
			t.Fatalf("bad base_url %q status = %d, want 400 (body=%s)", badURL, badRec.Code, badRec.Body.String())
		}
		badErr := decodeJSON[httpserver.ErrorEnvelope](t, badRec)
		if badErr.Error.Code != "invalid_base_url" {
			t.Fatalf("bad base_url %q code = %q, want invalid_base_url", badURL, badErr.Error.Code)
		}
	}

	// Strict JSON decoding: unknown field rejected with 400.
	unknownFieldRec := h.doJSON(t, http.MethodPost, "/api/cli-tools/claude-code/generate", `{"base_url":"http://127.0.0.1:20128","extra":1}`, false)
	if unknownFieldRec.Code != http.StatusBadRequest {
		t.Fatalf("unknown field status = %d, want 400", unknownFieldRec.Code)
	}

	// Base URL already ending in /v1 should not duplicate /v1/v1 for OpenAI-compatible tools.
	v1Rec := h.doJSON(t, http.MethodPost, "/api/cli-tools/codex-cli/generate", map[string]string{
		"base_url": "http://127.0.0.1:20128/v1/",
	}, false)
	if v1Rec.Code != http.StatusOK {
		t.Fatalf("codex-cli /v1 generate status = %d", v1Rec.Code)
	}
	v1Resp := decodeJSON[generateResponse](t, v1Rec)
	if strings.Contains(v1Resp.Content, "/v1/v1") {
		t.Fatalf("codex-cli duplicated /v1 suffix:\n%s", v1Resp.Content)
	}
	if !strings.Contains(v1Resp.Content, `base_url = "http://127.0.0.1:20128/v1"`) {
		t.Fatalf("codex-cli missing expected /v1 base_url:\n%s", v1Resp.Content)
	}

	// Optional fields omitted: claude-code omits ANTHROPIC_MODEL and uses placeholder token.
	omitRec := h.doJSON(t, http.MethodPost, "/api/cli-tools/claude-code/generate", map[string]string{
		"base_url": "http://127.0.0.1:20128",
	}, false)
	if omitRec.Code != http.StatusOK {
		t.Fatalf("claude-code omitted optional fields status = %d", omitRec.Code)
	}
	omitResp := decodeJSON[generateResponse](t, omitRec)
	if strings.Contains(omitResp.Content, "ANTHROPIC_MODEL") {
		t.Fatalf("expected ANTHROPIC_MODEL omitted when model empty:\n%s", omitResp.Content)
	}
}

func TestCLIToolsLoopbackAtomicApplyAndRollback(t *testing.T) {
	h := newTestHarness(t, "127.0.0.1", false)

	targetPath := filepath.Join(h.homeDir, ".claude", "settings.json")
	if err := os.MkdirAll(filepath.Dir(targetPath), 0o700); err != nil {
		t.Fatalf("os.MkdirAll: %v", err)
	}
	originalBytes := []byte("{\n  \"existing\": true,\n  \"theme\": \"ember\"\n}\n")
	if err := os.WriteFile(targetPath, originalBytes, 0o600); err != nil {
		t.Fatalf("os.WriteFile: %v", err)
	}

	// Apply on loopback writes new config atomically and creates timestamped backup.
	applyRec := h.doJSON(t, http.MethodPost, "/api/cli-tools/claude-code/apply", map[string]string{
		"base_url": "http://127.0.0.1:20128",
		"api_key":  "om4_test_key_123",
		"model":    "claude-sonnet-4",
	}, false)
	if applyRec.Code != http.StatusOK {
		t.Fatalf("apply status = %d, body = %s", applyRec.Code, applyRec.Body.String())
	}
	applied := decodeJSON[applyResponse](t, applyRec)
	if !applied.Applied || applied.Tool != "claude-code" || applied.TargetPath != targetPath {
		t.Fatalf("unexpected apply response: %+v", applied)
	}
	if !strings.HasPrefix(applied.BackupPath, targetPath+".om4.bak.") {
		t.Fatalf("backup_path = %q, want prefix %q", applied.BackupPath, targetPath+".om4.bak.")
	}

	newBytes, err := os.ReadFile(targetPath)
	if err != nil {
		t.Fatalf("os.ReadFile(target): %v", err)
	}
	if !json.Valid(newBytes) || !strings.Contains(string(newBytes), "om4_test_key_123") {
		t.Fatalf("target file not updated with generated JSON:\n%s", string(newBytes))
	}

	backupBytes, err := os.ReadFile(applied.BackupPath)
	if err != nil {
		t.Fatalf("os.ReadFile(backup): %v", err)
	}
	if !bytes.Equal(backupBytes, originalBytes) {
		t.Fatalf("backup bytes = %q, want %q", string(backupBytes), string(originalBytes))
	}

	// Ensure no leftover temp files remain in target directory.
	entries, err := os.ReadDir(filepath.Dir(targetPath))
	if err != nil {
		t.Fatalf("os.ReadDir: %v", err)
	}
	if len(entries) != 2 {
		t.Fatalf("expected exactly 2 files (target + backup), got %d", len(entries))
	}

	// Rollback restores exact original bytes atomically.
	rbRec := h.doJSON(t, http.MethodPost, "/api/cli-tools/claude-code/rollback", nil, false)
	if rbRec.Code != http.StatusOK {
		t.Fatalf("rollback status = %d, body = %s", rbRec.Code, rbRec.Body.String())
	}
	rb := decodeJSON[rollbackResponse](t, rbRec)
	if !rb.Restored || rb.Tool != "claude-code" || rb.TargetPath != targetPath || rb.BackupPath != applied.BackupPath {
		t.Fatalf("unexpected rollback response: %+v", rb)
	}

	restoredBytes, err := os.ReadFile(targetPath)
	if err != nil {
		t.Fatalf("os.ReadFile(restored): %v", err)
	}
	if !bytes.Equal(restoredBytes, originalBytes) {
		t.Fatalf("restored bytes = %q, want %q", string(restoredBytes), string(originalBytes))
	}

	// Apply to a tool without a pre-existing file creates parent dir and leaves backup_path empty.
	freshRec := h.doJSON(t, http.MethodPost, "/api/cli-tools/codex-cli/apply", map[string]string{
		"base_url": "http://127.0.0.1:20128",
	}, false)
	if freshRec.Code != http.StatusOK {
		t.Fatalf("fresh codex-cli apply status = %d, body = %s", freshRec.Code, freshRec.Body.String())
	}
	fresh := decodeJSON[applyResponse](t, freshRec)
	if fresh.BackupPath != "" {
		t.Fatalf("fresh apply backup_path = %q, want empty", fresh.BackupPath)
	}

	// Rollback without any backup returns 404 backup_not_found.
	noBackupRec := h.doJSON(t, http.MethodPost, "/api/cli-tools/codex-cli/rollback", nil, false)
	if noBackupRec.Code != http.StatusNotFound {
		t.Fatalf("rollback without backup status = %d, want 404", noBackupRec.Code)
	}
	noBackupErr := decodeJSON[httpserver.ErrorEnvelope](t, noBackupRec)
	if noBackupErr.Error.Code != "backup_not_found" {
		t.Fatalf("rollback without backup code = %q, want backup_not_found", noBackupErr.Error.Code)
	}

	// Multiple applies create ordered backups; rollback restores the lexicographically newest backup.
	apply2Rec := h.doJSON(t, http.MethodPost, "/api/cli-tools/codex-cli/apply", map[string]string{
		"base_url": "http://127.0.0.1:20128",
		"api_key":  "second_apply_key",
	}, false)
	if apply2Rec.Code != http.StatusOK {
		t.Fatalf("second codex-cli apply status = %d, body = %s", apply2Rec.Code, apply2Rec.Body.String())
	}
	apply2 := decodeJSON[applyResponse](t, apply2Rec)
	if apply2.BackupPath == "" {
		t.Fatalf("expected backup_path on second apply")
	}

	apply3Rec := h.doJSON(t, http.MethodPost, "/api/cli-tools/codex-cli/apply", map[string]string{
		"base_url": "http://127.0.0.1:20128",
		"api_key":  "third_apply_key",
	}, false)
	if apply3Rec.Code != http.StatusOK {
		t.Fatalf("third codex-cli apply status = %d, body = %s", apply3Rec.Code, apply3Rec.Body.String())
	}
	apply3 := decodeJSON[applyResponse](t, apply3Rec)
	if apply3.BackupPath <= apply2.BackupPath {
		t.Fatalf("expected strictly newer backup path: %q <= %q", apply3.BackupPath, apply2.BackupPath)
	}

	rbNewestRec := h.doJSON(t, http.MethodPost, "/api/cli-tools/codex-cli/rollback", nil, false)
	if rbNewestRec.Code != http.StatusOK {
		t.Fatalf("rollback newest status = %d, body = %s", rbNewestRec.Code, rbNewestRec.Body.String())
	}
	rbNewest := decodeJSON[rollbackResponse](t, rbNewestRec)
	if rbNewest.BackupPath != apply3.BackupPath {
		t.Fatalf("rollback chose %q, want newest %q", rbNewest.BackupPath, apply3.BackupPath)
	}
	codexBytes, err := os.ReadFile(fresh.TargetPath)
	if err != nil {
		t.Fatalf("os.ReadFile(codex): %v", err)
	}
	if !strings.Contains(string(codexBytes), "second_apply_key") || strings.Contains(string(codexBytes), "third_apply_key") {
		t.Fatalf("expected rollback to restore second_apply_key state, got:\n%s", string(codexBytes))
	}
}

func TestCLIToolsNonLoopbackRejectsApplyAndRollback(t *testing.T) {
	h := newTestHarness(t, "0.0.0.0", true)

	// Catalog reports local_apply_enabled = false on non-loopback bind.
	catRec := h.doJSON(t, http.MethodGet, "/api/cli-tools", nil, true)
	if catRec.Code != http.StatusOK {
		t.Fatalf("GET /api/cli-tools status = %d", catRec.Code)
	}
	cat := decodeJSON[catalogResponse](t, catRec)
	if cat.LocalApplyEnabled {
		t.Fatalf("local_apply_enabled = true on 0.0.0.0, want false")
	}

	// Generate still works on non-loopback (stateless).
	genRec := h.doJSON(t, http.MethodPost, "/api/cli-tools/claude-code/generate", map[string]string{
		"base_url": "https://router.example.com",
	}, true)
	if genRec.Code != http.StatusOK {
		t.Fatalf("generate on non-loopback status = %d, want 200", genRec.Code)
	}

	// Seed a file and a backup to prove neither apply nor rollback touches HomeDir.
	targetPath := filepath.Join(h.homeDir, ".claude", "settings.json")
	if err := os.MkdirAll(filepath.Dir(targetPath), 0o700); err != nil {
		t.Fatalf("os.MkdirAll: %v", err)
	}
	orig := []byte(`{"untouched":true}`)
	if err := os.WriteFile(targetPath, orig, 0o600); err != nil {
		t.Fatalf("os.WriteFile: %v", err)
	}
	bakPath := targetPath + ".om4.bak.20261123T000000.000000000Z"
	if err := os.WriteFile(bakPath, []byte(`{"backup":true}`), 0o600); err != nil {
		t.Fatalf("os.WriteFile(bak): %v", err)
	}

	applyRec := h.doJSON(t, http.MethodPost, "/api/cli-tools/claude-code/apply", map[string]string{
		"base_url": "http://127.0.0.1:20128",
	}, true)
	if applyRec.Code != http.StatusForbidden {
		t.Fatalf("non-loopback apply status = %d, want 403 (body=%s)", applyRec.Code, applyRec.Body.String())
	}
	applyErr := decodeJSON[httpserver.ErrorEnvelope](t, applyRec)
	if applyErr.Error.Code != "local_apply_disabled" {
		t.Fatalf("non-loopback apply code = %q, want local_apply_disabled", applyErr.Error.Code)
	}

	rbRec := h.doJSON(t, http.MethodPost, "/api/cli-tools/claude-code/rollback", nil, true)
	if rbRec.Code != http.StatusForbidden {
		t.Fatalf("non-loopback rollback status = %d, want 403 (body=%s)", rbRec.Code, rbRec.Body.String())
	}
	rbErr := decodeJSON[httpserver.ErrorEnvelope](t, rbRec)
	if rbErr.Error.Code != "local_apply_disabled" {
		t.Fatalf("non-loopback rollback code = %q, want local_apply_disabled", rbErr.Error.Code)
	}

	// Also verify fresh HomeDir directory (e.g. .codex) is never created on non-loopback.
	codexApplyRec := h.doJSON(t, http.MethodPost, "/api/cli-tools/codex-cli/apply", map[string]string{
		"base_url": "http://127.0.0.1:20128",
	}, true)
	if codexApplyRec.Code != http.StatusForbidden {
		t.Fatalf("non-loopback codex-cli apply status = %d, want 403", codexApplyRec.Code)
	}
	if _, err := os.Stat(filepath.Join(h.homeDir, ".codex")); !os.IsNotExist(err) {
		t.Fatalf("expected .codex dir not to exist on non-loopback apply, stat err = %v", err)
	}

	// Original file and directory contents must remain completely untouched.
	gotBytes, err := os.ReadFile(targetPath)
	if err != nil {
		t.Fatalf("os.ReadFile: %v", err)
	}
	if !bytes.Equal(gotBytes, orig) {
		t.Fatalf("target file modified on non-loopback: got %q, want %q", string(gotBytes), string(orig))
	}
	entries, err := os.ReadDir(filepath.Dir(targetPath))
	if err != nil {
		t.Fatalf("os.ReadDir: %v", err)
	}
	if len(entries) != 2 {
		t.Fatalf("expected .claude dir untouched with 2 entries, got %d", len(entries))
	}
}
