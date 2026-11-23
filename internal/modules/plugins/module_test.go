package plugins_test

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"path/filepath"
	"strings"
	"testing"

	"github.com/go-chi/chi/v5"

	"github.com/omeg4monx/router-gateway/internal/contracts"
	"github.com/omeg4monx/router-gateway/internal/modules/plugins"
	"github.com/omeg4monx/router-gateway/internal/platform/db"
	"github.com/omeg4monx/router-gateway/internal/platform/eventbus"
	"github.com/omeg4monx/router-gateway/internal/platform/lifecycle"
)

func setupPluginsTest(t *testing.T) (*plugins.Module, http.Handler) {
	t.Helper()
	database, err := db.Open(filepath.Join(t.TempDir(), "plugins.sqlite"))
	if err != nil {
		t.Fatalf("db.Open: %v", err)
	}
	t.Cleanup(func() { _ = database.Close() })

	mod := plugins.New(plugins.Options{
		DB:  database,
		Bus: eventbus.New(),
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
	apiRouter := chi.NewRouter()
	mod.RegisterRoutes(lifecycle.RouteMounts{AdminAPI: apiRouter})
	r.Mount("/api", apiRouter)
	return mod, r
}

func TestPluginsBuiltinsUpdateAndInterceptorScoping(t *testing.T) {
	mod, handler := setupPluginsTest(t)

	// List built-in 7 plugins.
	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, httptest.NewRequest(http.MethodGet, "/api/plugins", nil))
	if rec.Code != http.StatusOK {
		t.Fatalf("GET /api/plugins status = %d, body = %s", rec.Code, rec.Body.String())
	}
	var listResp struct {
		Plugins []struct {
			Key     string   `json:"key"`
			Enabled bool     `json:"enabled"`
			Models  []string `json:"models"`
		} `json:"plugins"`
	}
	if err := json.Unmarshal(rec.Body.Bytes(), &listResp); err != nil {
		t.Fatalf("decode list: %v", err)
	}
	if len(listResp.Plugins) != 7 {
		t.Fatalf("expected 7 built-in plugins, got %d", len(listResp.Plugins))
	}

	// Enable antiSlop (level=full) and jsonGuard scoped to gpt-4o.
	putBody := `{"enabled":true,"models":["gpt-4o"],"config":{"level":"full"}}`
	rec = httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPut, "/api/plugins/antiSlop", strings.NewReader(putBody))
	handler.ServeHTTP(rec, req)
	if rec.Code != http.StatusOK {
		t.Fatalf("PUT /api/plugins/antiSlop status = %d, body = %s", rec.Code, rec.Body.String())
	}

	rec = httptest.NewRecorder()
	req = httptest.NewRequest(http.MethodPut, "/api/plugins/jsonGuard", strings.NewReader(`{"enabled":true,"models":["gpt-4o"]}`))
	handler.ServeHTTP(rec, req)
	if rec.Code != http.StatusOK {
		t.Fatalf("PUT /api/plugins/jsonGuard status = %d, body = %s", rec.Code, rec.Body.String())
	}

	// Unknown plugin returns 404.
	rec = httptest.NewRecorder()
	handler.ServeHTTP(rec, httptest.NewRequest(http.MethodPut, "/api/plugins/doesNotExist", strings.NewReader(`{"enabled":true}`)))
	if rec.Code != http.StatusNotFound {
		t.Fatalf("expected 404 for unknown plugin, got %d", rec.Code)
	}

	// Interceptor mutates matching model request.
	canonReq := &contracts.CanonicalRequest{
		RequestID: "req_1",
		Protocol:  "openai-chat",
		Model:     "gpt-4o",
		Body:      []byte(`{"model":"gpt-4o","messages":[{"role":"user","content":"Return JSON"}]}`),
		Metadata:  map[string]any{},
	}
	if err := mod.OnRequest(context.Background(), canonReq); err != nil {
		t.Fatalf("OnRequest matching model: %v", err)
	}
	var mutated map[string]any
	if err := json.Unmarshal(canonReq.Body, &mutated); err != nil {
		t.Fatalf("unmarshal mutated body: %v", err)
	}
	msgs, _ := mutated["messages"].([]any)
	if len(msgs) < 2 {
		t.Fatalf("expected injected system directive message, got %v", mutated["messages"])
	}
	if rf, _ := mutated["response_format"].(map[string]any); rf["type"] != "json_object" {
		t.Fatalf("expected jsonGuard to set response_format.type=json_object, got %v", mutated["response_format"])
	}

	// Non-matching model passes through untouched.
	nonMatchBody := []byte(`{"model":"claude-3-5-sonnet","messages":[{"role":"user","content":"Hello"}]}`)
	otherReq := &contracts.CanonicalRequest{
		RequestID: "req_2",
		Protocol:  "openai-chat",
		Model:     "claude-3-5-sonnet",
		Body:      append([]byte(nil), nonMatchBody...),
		Metadata:  map[string]any{},
	}
	if err := mod.OnRequest(context.Background(), otherReq); err != nil {
		t.Fatalf("OnRequest non-matching model: %v", err)
	}
	if !bytes.Equal(otherReq.Body, nonMatchBody) {
		t.Fatalf("expected non-matching model body untouched, got %s", string(otherReq.Body))
	}

	// Malformed JSON body must not panic or error.
	badReq := &contracts.CanonicalRequest{
		RequestID: "req_3",
		Protocol:  "openai-chat",
		Model:     "gpt-4o",
		Body:      []byte(`{not-json`),
		Metadata:  map[string]any{},
	}
	if err := mod.OnRequest(context.Background(), badReq); err != nil {
		t.Fatalf("expected malformed JSON to pass through safely, got %v", err)
	}

	// Dry-run preview endpoint reports applied plugins.
	previewBody := `{"model":"gpt-4o","payload":{"model":"gpt-4o","messages":[{"role":"user","content":"hello"}]}}`
	rec = httptest.NewRecorder()
	handler.ServeHTTP(rec, httptest.NewRequest(http.MethodPost, "/api/plugins/preview", strings.NewReader(previewBody)))
	if rec.Code != http.StatusOK {
		t.Fatalf("POST /api/plugins/preview status = %d, body = %s", rec.Code, rec.Body.String())
	}
	var prevResp struct {
		AppliedPlugins []string `json:"applied_plugins"`
	}
	if err := json.Unmarshal(rec.Body.Bytes(), &prevResp); err != nil {
		t.Fatalf("decode preview: %v", err)
	}
	if len(prevResp.AppliedPlugins) != 2 {
		t.Fatalf("expected 2 applied plugins in preview, got %v", prevResp.AppliedPlugins)
	}
}

func TestSkillsCatalogCRUDAndExport(t *testing.T) {
	_, handler := setupPluginsTest(t)

	// Built-in skills seeded on start.
	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, httptest.NewRequest(http.MethodGet, "/api/skills", nil))
	if rec.Code != http.StatusOK {
		t.Fatalf("GET /api/skills status = %d, body = %s", rec.Code, rec.Body.String())
	}
	var listResp struct {
		Skills []struct {
			ID      string `json:"id"`
			Slug    string `json:"slug"`
			Builtin bool   `json:"builtin"`
		} `json:"skills"`
	}
	if err := json.Unmarshal(rec.Body.Bytes(), &listResp); err != nil {
		t.Fatalf("decode skills: %v", err)
	}
	if len(listResp.Skills) < 3 {
		t.Fatalf("expected at least 3 built-in skills, got %d", len(listResp.Skills))
	}

	// Create custom skill.
	createPayload := `{"slug":"custom-routing-audit","name":"Routing Audit","description":"Audit fallback routing","category":"operations","content":"# Routing Audit\n\nCheck combos and cooldowns."}`
	rec = httptest.NewRecorder()
	handler.ServeHTTP(rec, httptest.NewRequest(http.MethodPost, "/api/skills", strings.NewReader(createPayload)))
	if rec.Code != http.StatusCreated {
		t.Fatalf("POST /api/skills status = %d, body = %s", rec.Code, rec.Body.String())
	}
	var created struct {
		ID   string `json:"id"`
		Slug string `json:"slug"`
	}
	if err := json.Unmarshal(rec.Body.Bytes(), &created); err != nil {
		t.Fatalf("decode created skill: %v", err)
	}

	// Export SKILL.md with frontmatter.
	rec = httptest.NewRecorder()
	handler.ServeHTTP(rec, httptest.NewRequest(http.MethodGet, "/api/skills/"+created.ID+"/export", nil))
	if rec.Code != http.StatusOK {
		t.Fatalf("GET export status = %d, body = %s", rec.Code, rec.Body.String())
	}
	if !strings.Contains(rec.Body.String(), "name: custom-routing-audit") || !strings.Contains(rec.Body.String(), "# Routing Audit") {
		t.Fatalf("unexpected SKILL.md export: %s", rec.Body.String())
	}

	// Update custom skill.
	rec = httptest.NewRecorder()
	handler.ServeHTTP(rec, httptest.NewRequest(http.MethodPatch, "/api/skills/"+created.ID, strings.NewReader(`{"name":"Routing Audit v2","enabled":false}`)))
	if rec.Code != http.StatusOK {
		t.Fatalf("PATCH /api/skills status = %d, body = %s", rec.Code, rec.Body.String())
	}

	// Delete custom skill succeeds; deleting built-in skill returns 400.
	rec = httptest.NewRecorder()
	handler.ServeHTTP(rec, httptest.NewRequest(http.MethodDelete, "/api/skills/"+created.ID, nil))
	if rec.Code != http.StatusOK {
		t.Fatalf("DELETE custom skill status = %d, body = %s", rec.Code, rec.Body.String())
	}

	rec = httptest.NewRecorder()
	handler.ServeHTTP(rec, httptest.NewRequest(http.MethodDelete, "/api/skills/"+listResp.Skills[0].ID, nil))
	if rec.Code != http.StatusBadRequest {
		t.Fatalf("expected 400 deleting built-in skill, got %d", rec.Code)
	}
}
