package media

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"io"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"net/netip"
	"path/filepath"
	"strings"
	"sync"
	"sync/atomic"
	"testing"

	"github.com/go-chi/chi/v5"

	"github.com/omeg4monx/router-gateway/internal/contracts"
	"github.com/omeg4monx/router-gateway/internal/platform/crypto"
	"github.com/omeg4monx/router-gateway/internal/platform/db"
	"github.com/omeg4monx/router-gateway/internal/platform/eventbus"
	"github.com/omeg4monx/router-gateway/internal/platform/lifecycle"
)

type stubKeyValidator struct {
	validKey string
}

func (s *stubKeyValidator) ValidateKey(_ context.Context, rawKey string) (contracts.KeyInfo, error) {
	if rawKey != s.validKey {
		return contracts.KeyInfo{}, errors.New("invalid api key")
	}
	return contracts.KeyInfo{
		ID:   "key_test_1",
		Name: "test-key",
	}, nil
}

type recordingUsage struct {
	mu      sync.Mutex
	records []contracts.UsageRecord
}

func (r *recordingUsage) Record(_ context.Context, rec contracts.UsageRecord) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.records = append(r.records, rec)
}

func (r *recordingUsage) len() int {
	r.mu.Lock()
	defer r.mu.Unlock()
	return len(r.records)
}

func setupMediaTest(t *testing.T, allowPrivateFetch bool) (*Module, *chi.Mux, *chi.Mux, *recordingUsage) {
	t.Helper()
	dir := t.TempDir()
	database, err := db.Open(filepath.Join(dir, "media.sqlite"))
	if err != nil {
		t.Fatalf("open db: %v", err)
	}
	t.Cleanup(func() { _ = database.Close() })

	masterKey, err := crypto.LoadOrCreateMasterKey(dir, "")
	if err != nil {
		t.Fatalf("master key: %v", err)
	}
	vault, err := crypto.NewVault(masterKey)
	if err != nil {
		t.Fatalf("vault: %v", err)
	}

	usage := &recordingUsage{}
	mod := New(Options{
		DB:                database,
		Vault:             vault,
		HTTPClient:        &http.Client{},
		Keys:              &stubKeyValidator{validKey: "om_valid_key"},
		Usage:             usage,
		Bus:               eventbus.New(),
		AllowPrivateFetch: allowPrivateFetch,
	})

	if err := database.MigrateModule(context.Background(), mod.Name(), mod.Migrations()); err != nil {
		t.Fatalf("run migrations: %v", err)
	}
	if err := mod.Start(context.Background()); err != nil {
		t.Fatalf("start module: %v", err)
	}
	t.Cleanup(func() { _ = mod.Stop(context.Background()) })

	adminRouter := chi.NewRouter()
	gatewayRouter := chi.NewRouter()
	mod.RegisterRoutes(lifecycle.RouteMounts{
		AdminAPI: adminRouter,
		Gateway:  gatewayRouter,
	})
	return mod, adminRouter, gatewayRouter, usage
}

func doJSONReq(t *testing.T, h http.Handler, method, path string, body any, headers map[string]string) *httptest.ResponseRecorder {
	t.Helper()
	var r io.Reader
	if body != nil {
		raw, err := json.Marshal(body)
		if err != nil {
			t.Fatalf("marshal body: %v", err)
		}
		r = bytes.NewReader(raw)
	}
	req := httptest.NewRequest(method, path, r)
	if body != nil {
		req.Header.Set("Content-Type", "application/json")
	}
	for k, v := range headers {
		req.Header.Set(k, v)
	}
	rec := httptest.NewRecorder()
	h.ServeHTTP(rec, req)
	return rec
}

func TestMediaProvidersAndCombosCRUD(t *testing.T) {
	_, admin, _, _ := setupMediaTest(t, false)

	// Reject invalid capability.
	rec := doJSONReq(t, admin, http.MethodPost, "/media/providers", map[string]any{
		"name":       "bad-cap",
		"capability": "unknown",
		"driver":     "openai_compatible",
		"base_url":   "https://api.example.com",
	}, nil)
	if rec.Code != http.StatusBadRequest {
		t.Fatalf("expected 400 for invalid capability, got %d: %s", rec.Code, rec.Body.String())
	}

	// Create valid provider with secret API key.
	rec = doJSONReq(t, admin, http.MethodPost, "/media/providers", map[string]any{
		"name":          "openai-img",
		"capability":    "image",
		"driver":        "openai_compatible",
		"base_url":      "https://api.openai.com",
		"api_key":       "sk-super-secret-media-key",
		"default_model": "gpt-image-1",
		"priority":      10,
	}, nil)
	if rec.Code != http.StatusCreated {
		t.Fatalf("expected 201, got %d: %s", rec.Code, rec.Body.String())
	}
	if strings.Contains(rec.Body.String(), "sk-super-secret-media-key") {
		t.Fatalf("response leaked plaintext api_key: %s", rec.Body.String())
	}
	var created ProviderSummary
	if err := json.Unmarshal(rec.Body.Bytes(), &created); err != nil {
		t.Fatalf("decode created: %v", err)
	}
	if !created.HasAPIKey || created.ID == "" {
		t.Fatalf("expected has_api_key=true and non-empty id, got %+v", created)
	}

	// Duplicate name returns 409.
	rec = doJSONReq(t, admin, http.MethodPost, "/media/providers", map[string]any{
		"name":       "openai-img",
		"capability": "image",
		"driver":     "openai_compatible",
		"base_url":   "https://api.openai.com",
	}, nil)
	if rec.Code != http.StatusConflict {
		t.Fatalf("expected 409 duplicate name, got %d: %s", rec.Code, rec.Body.String())
	}

	// Patch clears API key and updates priority.
	rec = doJSONReq(t, admin, http.MethodPatch, "/media/providers/"+created.ID, map[string]any{
		"priority":      5,
		"clear_api_key": true,
	}, nil)
	if rec.Code != http.StatusOK {
		t.Fatalf("expected 200 patch, got %d: %s", rec.Code, rec.Body.String())
	}
	var patched ProviderSummary
	if err := json.Unmarshal(rec.Body.Bytes(), &patched); err != nil {
		t.Fatalf("decode patched: %v", err)
	}
	if patched.HasAPIKey || patched.Priority != 5 {
		t.Fatalf("expected has_api_key=false priority=5, got %+v", patched)
	}

	// Create media combo.
	rec = doJSONReq(t, admin, http.MethodPost, "/media/combos", map[string]any{
		"name":       "image-default",
		"capability": "image",
		"targets":    []string{created.ID},
	}, nil)
	if rec.Code != http.StatusCreated {
		t.Fatalf("expected 201 combo, got %d: %s", rec.Code, rec.Body.String())
	}
	var combo Combo
	if err := json.Unmarshal(rec.Body.Bytes(), &combo); err != nil {
		t.Fatalf("decode combo: %v", err)
	}
	if combo.ID == "" || len(combo.Targets) != 1 {
		t.Fatalf("unexpected combo: %+v", combo)
	}

	// List combos.
	rec = doJSONReq(t, admin, http.MethodGet, "/media/combos", nil, nil)
	if rec.Code != http.StatusOK {
		t.Fatalf("expected 200 list combos, got %d", rec.Code)
	}

	// Delete combo & provider.
	rec = doJSONReq(t, admin, http.MethodDelete, "/media/combos/"+combo.ID, nil, nil)
	if rec.Code != http.StatusOK {
		t.Fatalf("expected 200 delete combo, got %d", rec.Code)
	}
	rec = doJSONReq(t, admin, http.MethodDelete, "/media/providers/"+created.ID, nil, nil)
	if rec.Code != http.StatusOK {
		t.Fatalf("expected 200 delete provider, got %d", rec.Code)
	}
}

func TestMediaGatewayImagesFallbackAndClientErrorNoFallback(t *testing.T) {
	_, admin, gw, usage := setupMediaTest(t, false)

	var primaryCalls atomic.Int32
	var primaryStatus atomic.Int32
	primaryStatus.Store(http.StatusServiceUnavailable)
	primarySrv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		primaryCalls.Add(1)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(int(primaryStatus.Load()))
		_, _ = w.Write([]byte(`{"error":{"message":"primary failed"}}`))
	}))
	defer primarySrv.Close()

	var secondaryCalls atomic.Int32
	var gotAuth string
	secondarySrv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		secondaryCalls.Add(1)
		gotAuth = r.Header.Get("Authorization")
		if r.URL.Path != "/v1/images/generations" {
			w.WriteHeader(http.StatusNotFound)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{"created":1700000000,"data":[{"url":"https://cdn.example.com/img.png"}]}`))
	}))
	defer secondarySrv.Close()

	// Create primary + secondary providers.
	rec1 := doJSONReq(t, admin, http.MethodPost, "/media/providers", map[string]any{
		"name": "img-primary", "capability": "image", "driver": "openai_compatible",
		"base_url": primarySrv.URL, "api_key": "sk-p1", "priority": 1,
	}, nil)
	if rec1.Code != http.StatusCreated {
		t.Fatalf("create primary: %d %s", rec1.Code, rec1.Body.String())
	}
	var p1 ProviderSummary
	_ = json.Unmarshal(rec1.Body.Bytes(), &p1)

	rec2 := doJSONReq(t, admin, http.MethodPost, "/media/providers", map[string]any{
		"name": "img-secondary", "capability": "image", "driver": "openai_compatible",
		"base_url": secondarySrv.URL, "api_key": "sk-p2", "priority": 2,
	}, nil)
	if rec2.Code != http.StatusCreated {
		t.Fatalf("create secondary: %d %s", rec2.Code, rec2.Body.String())
	}
	var p2 ProviderSummary
	_ = json.Unmarshal(rec2.Body.Bytes(), &p2)

	// Create combo.
	recCombo := doJSONReq(t, admin, http.MethodPost, "/media/combos", map[string]any{
		"name": "img-combo", "capability": "image", "targets": []string{p1.ID, p2.ID},
	}, nil)
	if recCombo.Code != http.StatusCreated {
		t.Fatalf("create combo: %d %s", recCombo.Code, recCombo.Body.String())
	}

	// Missing API key returns 401.
	unauth := doJSONReq(t, gw, http.MethodPost, "/v1/images/generations", map[string]any{
		"model": "img-combo", "prompt": "ember",
	}, nil)
	if unauth.Code != http.StatusUnauthorized {
		t.Fatalf("expected 401 without key, got %d", unauth.Code)
	}

	// 503 on primary fails over to secondary and succeeds with 200.
	okRec := doJSONReq(t, gw, http.MethodPost, "/v1/images/generations", map[string]any{
		"model": "img-combo", "prompt": "ember",
	}, map[string]string{"Authorization": "Bearer om_valid_key"})
	if okRec.Code != http.StatusOK {
		t.Fatalf("expected 200 after fallback, got %d: %s", okRec.Code, okRec.Body.String())
	}
	if primaryCalls.Load() != 1 || secondaryCalls.Load() != 1 {
		t.Fatalf("expected 1 primary + 1 secondary call, got %d / %d", primaryCalls.Load(), secondaryCalls.Load())
	}
	if gotAuth != "Bearer sk-p2" {
		t.Fatalf("expected decrypted secondary bearer token, got %q", gotAuth)
	}
	if usage.len() != 1 {
		t.Fatalf("expected 1 usage record, got %d", usage.len())
	}

	// 400 client error on primary must NOT fail over to secondary.
	primaryStatus.Store(http.StatusBadRequest)
	badRec := doJSONReq(t, gw, http.MethodPost, "/v1/images/generations", map[string]any{
		"model": "img-combo", "prompt": "",
	}, map[string]string{"Authorization": "Bearer om_valid_key"})
	if badRec.Code != http.StatusBadRequest {
		t.Fatalf("expected 400 without fallback, got %d", badRec.Code)
	}
	if secondaryCalls.Load() != 1 {
		t.Fatalf("expected secondary NOT to be called on 400, got %d calls", secondaryCalls.Load())
	}
}

func TestMediaGatewayAudioVideoSearchEndpoints(t *testing.T) {
	_, admin, gw, _ := setupMediaTest(t, false)

	upstream := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case "/v1/audio/speech":
			w.Header().Set("Content-Type", "audio/mpeg")
			w.WriteHeader(http.StatusOK)
			_, _ = w.Write([]byte("ID3-FAKE-MP3"))
		case "/v1/audio/transcriptions":
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			_, _ = w.Write([]byte(`{"text":"hello router"}`))
		case "/v1/videos/generations":
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			_, _ = w.Write([]byte(`{"id":"vid_123","status":"queued"}`))
		case "/v1/videos/vid_123":
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			_, _ = w.Write([]byte(`{"id":"vid_123","status":"completed"}`))
		case "/search":
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			_, _ = w.Write([]byte(`{"results":[{"title":"Go","url":"https://go.dev"}]}`))
		default:
			w.WriteHeader(http.StatusNotFound)
		}
	}))
	defer upstream.Close()

	for _, p := range []map[string]any{
		{"name": "audio-p", "capability": "audio", "driver": "openai_compatible", "base_url": upstream.URL},
		{"name": "video-p", "capability": "video", "driver": "openai_compatible", "base_url": upstream.URL},
		{"name": "search-p", "capability": "search", "driver": "generic_http", "base_url": upstream.URL + "/search"},
	} {
		rec := doJSONReq(t, admin, http.MethodPost, "/media/providers", p, nil)
		if rec.Code != http.StatusCreated {
			t.Fatalf("create provider %v: %d %s", p["name"], rec.Code, rec.Body.String())
		}
	}

	headers := map[string]string{"Authorization": "Bearer om_valid_key"}

	// Audio speech.
	speechRec := doJSONReq(t, gw, http.MethodPost, "/v1/audio/speech", map[string]any{
		"model": "tts-1", "input": "hi", "voice": "alloy",
	}, headers)
	if speechRec.Code != http.StatusOK || speechRec.Body.String() != "ID3-FAKE-MP3" {
		t.Fatalf("speech failed: %d %q", speechRec.Code, speechRec.Body.String())
	}

	// Audio transcription (multipart).
	var mp bytes.Buffer
	mw := multipart.NewWriter(&mp)
	_ = mw.WriteField("model", "whisper-1")
	part, _ := mw.CreateFormFile("file", "clip.mp3")
	_, _ = part.Write([]byte("fake-audio"))
	_ = mw.Close()
	mpReq := httptest.NewRequest(http.MethodPost, "/v1/audio/transcriptions", &mp)
	mpReq.Header.Set("Content-Type", mw.FormDataContentType())
	mpReq.Header.Set("Authorization", "Bearer om_valid_key")
	mpRec := httptest.NewRecorder()
	gw.ServeHTTP(mpRec, mpReq)
	if mpRec.Code != http.StatusOK || !strings.Contains(mpRec.Body.String(), "hello router") {
		t.Fatalf("transcription failed: %d %s", mpRec.Code, mpRec.Body.String())
	}

	// Video generation & status.
	vidCreate := doJSONReq(t, gw, http.MethodPost, "/v1/videos/generations", map[string]any{
		"model": "sora-2", "prompt": "sunset",
	}, headers)
	if vidCreate.Code != http.StatusOK || !strings.Contains(vidCreate.Body.String(), "vid_123") {
		t.Fatalf("video create failed: %d %s", vidCreate.Code, vidCreate.Body.String())
	}
	vidGet := doJSONReq(t, gw, http.MethodGet, "/v1/videos/vid_123", nil, headers)
	if vidGet.Code != http.StatusOK || !strings.Contains(vidGet.Body.String(), "completed") {
		t.Fatalf("video get failed: %d %s", vidGet.Code, vidGet.Body.String())
	}

	// Search.
	searchRec := doJSONReq(t, gw, http.MethodPost, "/v1/search", map[string]any{
		"query": "golang",
	}, headers)
	if searchRec.Code != http.StatusOK || !strings.Contains(searchRec.Body.String(), "https://go.dev") {
		t.Fatalf("search failed: %d %s", searchRec.Code, searchRec.Body.String())
	}
}

func TestMediaWebFetchSSRFGuardAndExtraction(t *testing.T) {
	headers := map[string]string{"Authorization": "Bearer om_valid_key"}

	// SSRF-strict instance.
	_, _, gwStrict, _ := setupMediaTest(t, false)
	for _, blocked := range []string{
		"http://127.0.0.1:20128/api/health",
		"http://localhost/secret",
		"http://169.254.169.254/latest/meta-data/",
		"http://10.0.0.1/internal",
		"file:///etc/passwd",
	} {
		rec := doJSONReq(t, gwStrict, http.MethodPost, "/v1/web/fetch", map[string]any{
			"url": blocked,
		}, headers)
		if rec.Code != http.StatusForbidden && rec.Code != http.StatusBadRequest {
			t.Fatalf("expected 403/400 for SSRF URL %q, got %d: %s", blocked, rec.Code, rec.Body.String())
		}
	}

	// Private-fetch-allowed instance for HTML extraction & truncation test.
	pageSrv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		_, _ = w.Write([]byte(`<html><head><title>Ember Docs</title><style>body{color:red}</style></head><body><script>alert(1)</script><h1>Hello World</h1><p>Clean router text.</p></body></html>`))
	}))
	defer pageSrv.Close()

	_, _, gwAllow, _ := setupMediaTest(t, true)
	rec := doJSONReq(t, gwAllow, http.MethodPost, "/v1/web/fetch", map[string]any{
		"url": pageSrv.URL,
	}, headers)
	if rec.Code != http.StatusOK {
		t.Fatalf("expected 200 web fetch, got %d: %s", rec.Code, rec.Body.String())
	}
	var fetched WebFetchResponse
	if err := json.Unmarshal(rec.Body.Bytes(), &fetched); err != nil {
		t.Fatalf("decode fetch: %v", err)
	}
	if fetched.Title != "Ember Docs" {
		t.Fatalf("expected title 'Ember Docs', got %q", fetched.Title)
	}
	if strings.Contains(fetched.Text, "alert(1)") || strings.Contains(fetched.Text, "color:red") {
		t.Fatalf("expected script/style stripped, got %q", fetched.Text)
	}
	if !strings.Contains(fetched.Text, "Hello World") || !strings.Contains(fetched.Text, "Clean router text.") {
		t.Fatalf("expected readable body text, got %q", fetched.Text)
	}
}

func TestMediaSSRFReservedRangesAndRedirectBlock(t *testing.T) {
	blockedIPs := []string{
		"0.0.0.1",
		"100.64.0.1",
		"192.0.0.1",
		"198.18.0.1",
		"240.0.0.1",
		"64:ff9b::7f00:1",
		"::127.0.0.1",
		"fe80::1%eth0",
	}
	for _, raw := range blockedIPs {
		addr, err := netip.ParseAddr(raw)
		if err != nil {
			t.Fatalf("parse %q: %v", raw, err)
		}
		if !isBlockedAddr(addr) {
			t.Fatalf("expected %q to be blocked by SSRF policy", raw)
		}
	}
	if err := validateRedirectTarget(context.Background(), "http://127.0.0.1:20128/api/health", false); !errors.Is(err, errSSRFBlocked) {
		t.Fatalf("expected redirect to loopback to be blocked with errSSRFBlocked, got %v", err)
	}
	if err := validateRedirectTarget(context.Background(), "file:///etc/passwd", false); !errors.Is(err, errInvalidFetchURL) {
		t.Fatalf("expected file:// redirect to be blocked with errInvalidFetchURL, got %v", err)
	}
}

func TestMediaKeyAllowedModelsEnforcement(t *testing.T) {
	if !isModelAllowed(nil, "img-combo", CapabilityImage) {
		t.Fatalf("empty allowed list must allow all")
	}
	if !isModelAllowed([]string{"*"}, "img-combo", CapabilityImage) {
		t.Fatalf("wildcard must allow all")
	}
	if isModelAllowed([]string{"gpt-4o"}, "img-combo", CapabilityImage) {
		t.Fatalf("restricted key must NOT allow unrelated media capability")
	}
	if !isModelAllowed([]string{"media/image"}, "", CapabilityImage) {
		t.Fatalf("media/<capability> must be permitted")
	}
}
