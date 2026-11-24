package keys_test

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"io"
	"io/fs"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"path/filepath"
	"slices"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/omeg4monx/router-gateway/internal/contracts"
	"github.com/omeg4monx/router-gateway/internal/modules/keys"
	"github.com/omeg4monx/router-gateway/internal/platform/config"
	"github.com/omeg4monx/router-gateway/internal/platform/db"
	"github.com/omeg4monx/router-gateway/internal/platform/eventbus"
	"github.com/omeg4monx/router-gateway/internal/platform/httpserver"
	"github.com/omeg4monx/router-gateway/internal/platform/lifecycle"
)

var (
	_ lifecycle.Module       = (*keys.Module)(nil)
	_ contracts.KeyValidator = (*keys.Module)(nil)
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
	mod    *keys.Module
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
	mod := keys.New(keys.Options{
		DB:  database,
		Bus: bus,
	})

	if err := reg.Register(stubSystemModule{}); err != nil {
		t.Fatalf("reg.Register(system): %v", err)
	}
	if err := reg.Register(mod); err != nil {
		t.Fatalf("reg.Register(keys): %v", err)
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

type listKeysResp struct {
	Keys []keys.APIKeySummary `json:"keys"`
}

type keySecretResp struct {
	Key    keys.APIKeySummary `json:"key"`
	Secret string             `json:"secret"`
}

func sha256Hex(s string) string {
	sum := sha256.Sum256([]byte(s))
	return hex.EncodeToString(sum[:])
}

func assertFixedWidthUTC(t *testing.T, label, ts string) {
	t.Helper()
	const fixedLayout = "2006-01-02T15:04:05.000000000Z"
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

func TestKeysCRUDAndSecretHygiene(t *testing.T) {
	h := newTestHarness(t, &stubAdminAuth{subject: "operator-1"})

	// Verify migration version 1 is recorded for keys.
	versions, err := h.db.AppliedVersions(context.Background(), "keys")
	if err != nil {
		t.Fatalf("AppliedVersions: %v", err)
	}
	if len(versions) != 1 || versions[0] != 1 {
		t.Fatalf("applied versions = %v, want [1]", versions)
	}

	// 1) Initial GET /api/keys returns empty non-null array.
	rec := h.doJSON(t, http.MethodGet, "/api/keys", "")
	if rec.Code != http.StatusOK {
		t.Fatalf("GET /api/keys status = %d, want 200 (body=%s)", rec.Code, rec.Body.String())
	}
	if strings.TrimSpace(rec.Body.String()) != `{"keys":[]}` {
		t.Fatalf("empty GET /api/keys body = %s, want {\"keys\":[]}", rec.Body.String())
	}

	// 2) POST /api/keys returns 201 with secret starting with sk-om4- and full metadata.
	futureExpiry := time.Now().UTC().Add(24 * time.Hour).Truncate(time.Second).Format(time.RFC3339)
	rec = h.doJSON(t, http.MethodPost, "/api/keys", `{
		"name": "primary-service",
		"allowed_models": ["openai/gpt-4o", "anthropic/*"],
		"scopes": ["gateway", "models"],
		"tags": ["prod", "team-a"],
		"rate_limit_rpm": 120,
		"quota_usd": 50.25,
		"expires_at": "`+futureExpiry+`"
	}`)
	if rec.Code != http.StatusCreated {
		t.Fatalf("POST /api/keys status = %d, want 201 (body=%s)", rec.Code, rec.Body.String())
	}

	var created keySecretResp
	if err := json.Unmarshal(rec.Body.Bytes(), &created); err != nil {
		t.Fatalf("unmarshal POST /api/keys: %v", err)
	}
	if !strings.HasPrefix(created.Secret, "sk-om4-") {
		t.Fatalf("secret = %q, want prefix sk-om4-", created.Secret)
	}
	if len(created.Secret) != len("sk-om4-")+48 {
		t.Fatalf("secret len = %d, want %d", len(created.Secret), len("sk-om4-")+48)
	}
	if created.Key.ID == "" {
		t.Fatalf("created key id is empty")
	}
	if created.Key.Name != "primary-service" {
		t.Errorf("name = %q, want primary-service", created.Key.Name)
	}
	if created.Key.KeyPrefix != created.Secret[:14] {
		t.Errorf("key_prefix = %q, want %q", created.Key.KeyPrefix, created.Secret[:14])
	}
	if !slices.Equal(created.Key.AllowedModels, []string{"openai/gpt-4o", "anthropic/*"}) {
		t.Errorf("allowed_models = %v", created.Key.AllowedModels)
	}
	if !slices.Equal(created.Key.Scopes, []string{"gateway", "models"}) {
		t.Errorf("scopes = %v", created.Key.Scopes)
	}
	if !slices.Equal(created.Key.Tags, []string{"prod", "team-a"}) {
		t.Errorf("tags = %v", created.Key.Tags)
	}
	if created.Key.RateLimitRPM != 120 || created.Key.QuotaUSD != 50.25 || created.Key.Disabled {
		t.Errorf("limits/disabled = rpm:%d quota:%v disabled:%v", created.Key.RateLimitRPM, created.Key.QuotaUSD, created.Key.Disabled)
	}
	if created.Key.ExpiresAt == nil {
		t.Fatalf("expires_at is nil, want non-nil")
	}
	assertFixedWidthUTC(t, "expires_at", *created.Key.ExpiresAt)
	assertFixedWidthUTC(t, "created_at", created.Key.CreatedAt)
	assertFixedWidthUTC(t, "updated_at", created.Key.UpdatedAt)

	secretHash := sha256Hex(created.Secret)

	// 3) POST /api/keys with defaults only.
	rec2 := h.doJSON(t, http.MethodPost, "/api/keys", `{"name":"default-key"}`)
	if rec2.Code != http.StatusCreated {
		t.Fatalf("POST default key status = %d, want 201 (body=%s)", rec2.Code, rec2.Body.String())
	}
	var created2 keySecretResp
	if err := json.Unmarshal(rec2.Body.Bytes(), &created2); err != nil {
		t.Fatalf("unmarshal default key: %v", err)
	}
	if created2.Secret == created.Secret {
		t.Fatalf("generated duplicate secrets")
	}
	if !slices.Equal(created2.Key.Scopes, []string{"gateway"}) {
		t.Errorf("default scopes = %v, want [gateway]", created2.Key.Scopes)
	}
	if len(created2.Key.AllowedModels) != 0 || len(created2.Key.Tags) != 0 {
		t.Errorf("default slices not empty: allowed=%v tags=%v", created2.Key.AllowedModels, created2.Key.Tags)
	}
	if created2.Key.ExpiresAt != nil || created2.Key.LastUsedAt != nil {
		t.Errorf("default expires_at=%v last_used_at=%v, want nil", created2.Key.ExpiresAt, created2.Key.LastUsedAt)
	}

	// 4) GET /api/keys MUST NEVER leak raw secret or hash, and returns newest-first.
	listRec := h.doJSON(t, http.MethodGet, "/api/keys", "")
	if listRec.Code != http.StatusOK {
		t.Fatalf("GET /api/keys status = %d, want 200", listRec.Code)
	}
	listBody := listRec.Body.String()
	for _, forbidden := range []string{created.Secret, secretHash, created2.Secret, sha256Hex(created2.Secret), "key_hash", `"secret"`} {
		if strings.Contains(listBody, forbidden) {
			t.Fatalf("GET /api/keys leaked forbidden token/hash/field %q in body: %s", forbidden, listBody)
		}
	}
	var listed listKeysResp
	if err := json.Unmarshal(listRec.Body.Bytes(), &listed); err != nil {
		t.Fatalf("unmarshal GET /api/keys: %v", err)
	}
	if len(listed.Keys) != 2 {
		t.Fatalf("listed keys len = %d, want 2", len(listed.Keys))
	}
	if listed.Keys[0].ID != created2.Key.ID || listed.Keys[1].ID != created.Key.ID {
		t.Errorf("keys order = [%s, %s], want newest-first [%s, %s]", listed.Keys[0].ID, listed.Keys[1].ID, created2.Key.ID, created.Key.ID)
	}

	// 5) GET /api/keys/{id} MUST NEVER leak raw secret or hash.
	getRec := h.doJSON(t, http.MethodGet, "/api/keys/"+created.Key.ID, "")
	if getRec.Code != http.StatusOK {
		t.Fatalf("GET /api/keys/{id} status = %d, want 200 (body=%s)", getRec.Code, getRec.Body.String())
	}
	getBody := getRec.Body.String()
	for _, forbidden := range []string{created.Secret, secretHash, "key_hash", `"secret"`} {
		if strings.Contains(getBody, forbidden) {
			t.Fatalf("GET /api/keys/{id} leaked forbidden token/hash/field %q in body: %s", forbidden, getBody)
		}
	}
	var fetched keys.APIKeySummary
	if err := json.Unmarshal(getRec.Body.Bytes(), &fetched); err != nil {
		t.Fatalf("unmarshal GET /api/keys/{id}: %v", err)
	}
	if fetched.ID != created.Key.ID || fetched.KeyPrefix != created.Key.KeyPrefix {
		t.Errorf("fetched = %+v, want ID=%s prefix=%s", fetched, created.Key.ID, created.Key.KeyPrefix)
	}

	// 6) PATCH /api/keys/{id} MUST NEVER leak raw secret or hash.
	patchRec := h.doJSON(t, http.MethodPatch, "/api/keys/"+created.Key.ID, `{"name":"renamed-service"}`)
	if patchRec.Code != http.StatusOK {
		t.Fatalf("PATCH /api/keys/{id} status = %d, want 200 (body=%s)", patchRec.Code, patchRec.Body.String())
	}
	patchBody := patchRec.Body.String()
	for _, forbidden := range []string{created.Secret, secretHash, "key_hash", `"secret"`} {
		if strings.Contains(patchBody, forbidden) {
			t.Fatalf("PATCH /api/keys/{id} leaked forbidden token/hash/field %q in body: %s", forbidden, patchBody)
		}
	}

	// 7) Verify SQLite api_keys table stores SHA-256 hash and NEVER raw secret.
	var storedHash, storedPrefix, storedName string
	err = h.db.Reader().QueryRowContext(
		context.Background(),
		`SELECT key_hash, key_prefix, name FROM api_keys WHERE id = ?`,
		created.Key.ID,
	).Scan(&storedHash, &storedPrefix, &storedName)
	if err != nil {
		t.Fatalf("query api_keys directly: %v", err)
	}
	if storedHash != secretHash {
		t.Fatalf("stored key_hash = %q, want sha256 hex %q", storedHash, secretHash)
	}
	if strings.Contains(storedHash, created.Secret) || strings.Contains(storedPrefix, created.Secret) {
		t.Fatalf("raw secret found in SQLite row")
	}

	var rawLeakCount int
	err = h.db.Reader().QueryRowContext(
		context.Background(),
		`SELECT COUNT(*) FROM api_keys WHERE
			id LIKE '%' || ? || '%' OR
			name LIKE '%' || ? || '%' OR
			key_hash LIKE '%' || ? || '%' OR
			key_prefix LIKE '%' || ? || '%' OR
			allowed_models_json LIKE '%' || ? || '%' OR
			scopes_json LIKE '%' || ? || '%' OR
			tags_json LIKE '%' || ? || '%'`,
		created.Secret, created.Secret, created.Secret, created.Secret, created.Secret, created.Secret, created.Secret,
	).Scan(&rawLeakCount)
	if err != nil {
		t.Fatalf("scan raw leak count: %v", err)
	}
	if rawLeakCount != 0 {
		t.Fatalf("raw secret leaked into SQLite api_keys table (matches=%d)", rawLeakCount)
	}

	// 8) 404 responses for unknown IDs across GET/PATCH/rotate/DELETE.
	for _, tc := range []struct {
		method string
		path   string
		body   string
	}{
		{http.MethodGet, "/api/keys/key_doesnotexist", ""},
		{http.MethodPatch, "/api/keys/key_doesnotexist", `{"disabled":true}`},
		{http.MethodPost, "/api/keys/key_doesnotexist/rotate", ""},
		{http.MethodDelete, "/api/keys/key_doesnotexist", ""},
	} {
		rec := h.doJSON(t, tc.method, tc.path, tc.body)
		if rec.Code != http.StatusNotFound {
			t.Fatalf("%s %s status = %d, want 404 (body=%s)", tc.method, tc.path, rec.Code, rec.Body.String())
		}
		var env httpserver.ErrorEnvelope
		if err := json.Unmarshal(rec.Body.Bytes(), &env); err != nil {
			t.Fatalf("unmarshal 404 envelope: %v", err)
		}
		if env.Error.Code != "key_not_found" {
			t.Errorf("%s %s error.code = %q, want key_not_found", tc.method, tc.path, env.Error.Code)
		}
	}
}

func TestKeysValidateAndImmediateCacheInvalidation(t *testing.T) {
	h := newTestHarness(t, &stubAdminAuth{subject: "ops-admin"})
	ctx := context.Background()

	events, unsub := h.bus.Subscribe(keys.TopicKeysChanged, 16)
	defer unsub()

	expectEvent := func(t *testing.T, wantAction, wantKeyID string) {
		t.Helper()
		select {
		case ev := <-events:
			if ev.Topic != keys.TopicKeysChanged {
				t.Fatalf("event topic = %q, want %q", ev.Topic, keys.TopicKeysChanged)
			}
			payload, ok := ev.Payload.(keys.ChangedEvent)
			if !ok {
				t.Fatalf("event payload type = %T, want keys.ChangedEvent", ev.Payload)
			}
			if payload.Action != wantAction || payload.KeyID != wantKeyID {
				t.Fatalf("event = %+v, want action=%q key_id=%q", payload, wantAction, wantKeyID)
			}
			if payload.ChangedBy != "ops-admin" {
				t.Errorf("event changed_by = %q, want ops-admin", payload.ChangedBy)
			}
			assertFixedWidthUTC(t, "event.updated_at", payload.UpdatedAt)
		case <-time.After(time.Second):
			t.Fatalf("timed out waiting for %s event on %s", wantAction, keys.TopicKeysChanged)
		}
	}

	// 1) Create key.
	rec := h.doJSON(t, http.MethodPost, "/api/keys", `{
		"name": "gateway-client",
		"allowed_models": ["openai/gpt-4o", "anthropic/*"],
		"scopes": ["gateway"],
		"tags": ["edge"],
		"rate_limit_rpm": 60,
		"quota_usd": 15.5
	}`)
	if rec.Code != http.StatusCreated {
		t.Fatalf("POST /api/keys status = %d (body=%s)", rec.Code, rec.Body.String())
	}
	var created keySecretResp
	if err := json.Unmarshal(rec.Body.Bytes(), &created); err != nil {
		t.Fatalf("unmarshal create: %v", err)
	}
	expectEvent(t, "created", created.Key.ID)

	// 2) ValidateKey succeeds and populates cache.
	info, err := h.mod.ValidateKey(ctx, created.Secret)
	if err != nil {
		t.Fatalf("ValidateKey initial: %v", err)
	}
	if info.ID != created.Key.ID || info.Name != "gateway-client" || info.RateLimitRPM != 60 || info.QuotaUSD != 15.5 || info.Disabled {
		t.Fatalf("unexpected KeyInfo: %+v", info)
	}
	if !slices.Equal(info.AllowedModels, []string{"openai/gpt-4o", "anthropic/*"}) {
		t.Fatalf("AllowedModels = %v", info.AllowedModels)
	}

	// Mutating returned slices must not poison the in-memory cache.
	info.AllowedModels[0] = "tampered/model"
	info2, err := h.mod.ValidateKey(ctx, created.Secret)
	if err != nil {
		t.Fatalf("ValidateKey cached: %v", err)
	}
	if info2.AllowedModels[0] != "openai/gpt-4o" {
		t.Fatalf("cache poisoned by caller mutation: got %v", info2.AllowedModels)
	}

	// Prove ValidateKey is genuinely serving from the in-memory cache (not hitting SQLite):
	// a direct out-of-band SQLite update is invisible until PATCH invalidates the cache.
	if _, err := h.db.Writer().ExecContext(ctx, `UPDATE api_keys SET name = ? WHERE id = ?`, "out-of-band-edit", created.Key.ID); err != nil {
		t.Fatalf("direct SQLite update: %v", err)
	}
	cachedInfo, err := h.mod.ValidateKey(ctx, created.Secret)
	if err != nil {
		t.Fatalf("ValidateKey cache hit check: %v", err)
	}
	if cachedInfo.Name != "gateway-client" {
		t.Fatalf("expected cached name %q, got %q (cache was bypassed)", "gateway-client", cachedInfo.Name)
	}

	// 3) PATCH disabled: true immediately invalidates cache -> ErrKeyDisabled.
	rec = h.doJSON(t, http.MethodPatch, "/api/keys/"+created.Key.ID, `{"disabled":true}`)
	if rec.Code != http.StatusOK {
		t.Fatalf("PATCH disable status = %d (body=%s)", rec.Code, rec.Body.String())
	}
	expectEvent(t, "updated", created.Key.ID)

	if _, err := h.mod.ValidateKey(ctx, created.Secret); !errors.Is(err, keys.ErrKeyDisabled) {
		t.Fatalf("ValidateKey after disable err = %v, want ErrKeyDisabled", err)
	}

	// 4) PATCH re-enable + update allowed_models/limits immediately reflects in ValidateKey.
	rec = h.doJSON(t, http.MethodPatch, "/api/keys/"+created.Key.ID, `{
		"disabled": false,
		"allowed_models": ["openai/o3-mini"],
		"rate_limit_rpm": 240,
		"quota_usd": 99.0,
		"tags": ["vip"]
	}`)
	if rec.Code != http.StatusOK {
		t.Fatalf("PATCH re-enable status = %d (body=%s)", rec.Code, rec.Body.String())
	}
	expectEvent(t, "updated", created.Key.ID)

	info3, err := h.mod.ValidateKey(ctx, created.Secret)
	if err != nil {
		t.Fatalf("ValidateKey after re-enable: %v", err)
	}
	if !slices.Equal(info3.AllowedModels, []string{"openai/o3-mini"}) || info3.RateLimitRPM != 240 || info3.QuotaUSD != 99.0 || !slices.Equal(info3.Tags, []string{"vip"}) {
		t.Fatalf("stale metadata after PATCH: %+v", info3)
	}

	// 5) POST /api/keys/{id}/rotate immediately invalidates old secret and accepts new secret.
	oldSecret := created.Secret
	rotRec := h.doJSON(t, http.MethodPost, "/api/keys/"+created.Key.ID+"/rotate", "")
	if rotRec.Code != http.StatusOK {
		t.Fatalf("POST rotate status = %d (body=%s)", rotRec.Code, rotRec.Body.String())
	}
	expectEvent(t, "rotated", created.Key.ID)

	var rotated keySecretResp
	if err := json.Unmarshal(rotRec.Body.Bytes(), &rotated); err != nil {
		t.Fatalf("unmarshal rotate: %v", err)
	}
	if rotated.Secret == oldSecret {
		t.Fatalf("rotated secret equals old secret")
	}
	if !strings.HasPrefix(rotated.Secret, "sk-om4-") || len(rotated.Secret) != len("sk-om4-")+48 {
		t.Fatalf("rotated secret malformed: %q", rotated.Secret)
	}
	if rotated.Key.KeyPrefix != rotated.Secret[:14] {
		t.Fatalf("rotated key_prefix = %q, want %q", rotated.Key.KeyPrefix, rotated.Secret[:14])
	}

	if _, err := h.mod.ValidateKey(ctx, oldSecret); !errors.Is(err, keys.ErrKeyInvalid) {
		t.Fatalf("ValidateKey(oldSecret) after rotate err = %v, want ErrKeyInvalid", err)
	}
	infoRot, err := h.mod.ValidateKey(ctx, rotated.Secret)
	if err != nil {
		t.Fatalf("ValidateKey(newSecret) after rotate: %v", err)
	}
	if infoRot.ID != created.Key.ID || !slices.Equal(infoRot.AllowedModels, []string{"openai/o3-mini"}) {
		t.Fatalf("rotated KeyInfo mismatch: %+v", infoRot)
	}

	// Concurrent reads + PATCH exercise under -race.
	var wg sync.WaitGroup
	for i := 0; i < 8; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for j := 0; j < 20; j++ {
				_, _ = h.mod.ValidateKey(ctx, rotated.Secret)
			}
		}()
	}
	for i := 0; i < 4; i++ {
		patchRec := h.doJSON(t, http.MethodPatch, "/api/keys/"+created.Key.ID, `{"rate_limit_rpm":300}`)
		if patchRec.Code != http.StatusOK {
			t.Errorf("concurrent PATCH status = %d", patchRec.Code)
		}
		expectEvent(t, "updated", created.Key.ID)
	}
	wg.Wait()

	// 6) DELETE /api/keys/{id} immediately rejects secret and returns 204.
	delRec := h.doJSON(t, http.MethodDelete, "/api/keys/"+created.Key.ID, "")
	if delRec.Code != http.StatusNoContent {
		t.Fatalf("DELETE status = %d, want 204 (body=%s)", delRec.Code, delRec.Body.String())
	}
	expectEvent(t, "deleted", created.Key.ID)

	if _, err := h.mod.ValidateKey(ctx, rotated.Secret); !errors.Is(err, keys.ErrKeyInvalid) {
		t.Fatalf("ValidateKey after DELETE err = %v, want ErrKeyInvalid", err)
	}
	getAfterDel := h.doJSON(t, http.MethodGet, "/api/keys/"+created.Key.ID, "")
	if getAfterDel.Code != http.StatusNotFound {
		t.Fatalf("GET after DELETE status = %d, want 404", getAfterDel.Code)
	}
}

func TestKeysExpirationAndValidationErrors(t *testing.T) {
	h := newTestHarness(t, nil)
	ctx := context.Background()

	// 1) Expired key fails ValidateKey with ErrKeyExpired.
	pastExpiry := time.Now().UTC().Add(-2 * time.Minute).Format(time.RFC3339)
	rec := h.doJSON(t, http.MethodPost, "/api/keys", `{
		"name": "expired-key",
		"expires_at": "`+pastExpiry+`"
	}`)
	if rec.Code != http.StatusCreated {
		t.Fatalf("POST expired key status = %d, want 201 (body=%s)", rec.Code, rec.Body.String())
	}
	var expired keySecretResp
	if err := json.Unmarshal(rec.Body.Bytes(), &expired); err != nil {
		t.Fatalf("unmarshal expired key: %v", err)
	}
	if _, err := h.mod.ValidateKey(ctx, expired.Secret); !errors.Is(err, keys.ErrKeyExpired) {
		t.Fatalf("ValidateKey(expired) err = %v, want ErrKeyExpired", err)
	}

	// Clearing expires_at via PATCH null makes key valid immediately.
	clearRec := h.doJSON(t, http.MethodPatch, "/api/keys/"+expired.Key.ID, `{"expires_at":null}`)
	if clearRec.Code != http.StatusOK {
		t.Fatalf("PATCH clear expires_at status = %d, want 200 (body=%s)", clearRec.Code, clearRec.Body.String())
	}
	if _, err := h.mod.ValidateKey(ctx, expired.Secret); err != nil {
		t.Fatalf("ValidateKey after clearing expires_at: %v", err)
	}

	// Re-expiring a cached key via PATCH immediately fails with ErrKeyExpired.
	reExpireRec := h.doJSON(t, http.MethodPatch, "/api/keys/"+expired.Key.ID, `{"expires_at":"`+pastExpiry+`"}`)
	if reExpireRec.Code != http.StatusOK {
		t.Fatalf("PATCH re-expire status = %d, want 200 (body=%s)", reExpireRec.Code, reExpireRec.Body.String())
	}
	if _, err := h.mod.ValidateKey(ctx, expired.Secret); !errors.Is(err, keys.ErrKeyExpired) {
		t.Fatalf("ValidateKey after re-expiring err = %v, want ErrKeyExpired", err)
	}

	// A cached key whose expires_at passes in wall-clock time must also fail with ErrKeyExpired.
	soonExpiry := time.Now().UTC().Add(250 * time.Millisecond).Format(time.RFC3339Nano)
	soonRec := h.doJSON(t, http.MethodPost, "/api/keys", `{
		"name": "short-lived-key",
		"expires_at": "`+soonExpiry+`"
	}`)
	if soonRec.Code != http.StatusCreated {
		t.Fatalf("POST short-lived key status = %d, want 201 (body=%s)", soonRec.Code, soonRec.Body.String())
	}
	var soonKey keySecretResp
	if err := json.Unmarshal(soonRec.Body.Bytes(), &soonKey); err != nil {
		t.Fatalf("unmarshal short-lived key: %v", err)
	}
	if _, err := h.mod.ValidateKey(ctx, soonKey.Secret); err != nil {
		t.Fatalf("ValidateKey before short expiry: %v", err)
	}
	time.Sleep(300 * time.Millisecond)
	if _, err := h.mod.ValidateKey(ctx, soonKey.Secret); !errors.Is(err, keys.ErrKeyExpired) {
		t.Fatalf("ValidateKey after cached wall-clock expiry err = %v, want ErrKeyExpired", err)
	}

	// 2) Malformed and unknown raw keys fail with ErrKeyInvalid.
	for _, raw := range []string{
		"",
		"   ",
		"sk-wrong-" + strings.Repeat("a", 48),
		"sk-om4-short",
		"sk-om4-" + strings.Repeat("!", 48),
		"sk-om4-" + strings.Repeat("A", 48), // valid format, unknown key
	} {
		if _, err := h.mod.ValidateKey(ctx, raw); !errors.Is(err, keys.ErrKeyInvalid) {
			t.Errorf("ValidateKey(%q) err = %v, want ErrKeyInvalid", raw, err)
		}
	}

	// 3) POST /api/keys validation errors return 400.
	postBadCases := []struct {
		name string
		body string
	}{
		{"empty name", `{"name":""}`},
		{"whitespace name", `{"name":"   "}`},
		{"missing name", `{}`},
		{"negative rate_limit_rpm", `{"name":"k1","rate_limit_rpm":-1}`},
		{"negative quota_usd", `{"name":"k1","quota_usd":-0.01}`},
		{"invalid expires_at", `{"name":"k1","expires_at":"not-a-date"}`},
		{"empty expires_at string", `{"name":"k1","expires_at":""}`},
		{"empty allowed_models entry", `{"name":"k1","allowed_models":[""]}`},
		{"unknown field", `{"name":"k1","extra":true}`},
		{"malformed json", `{"name":`},
		{"trailing json", `{"name":"k1"}{}`},
	}
	for _, tc := range postBadCases {
		t.Run("POST_"+tc.name, func(t *testing.T) {
			rec := h.doJSON(t, http.MethodPost, "/api/keys", tc.body)
			if rec.Code != http.StatusBadRequest {
				t.Fatalf("POST (%s) status = %d, want 400 (body=%s)", tc.name, rec.Code, rec.Body.String())
			}
			var env httpserver.ErrorEnvelope
			if err := json.Unmarshal(rec.Body.Bytes(), &env); err != nil {
				t.Fatalf("unmarshal error envelope: %v", err)
			}
			if env.Error.Code != "invalid_request" {
				t.Errorf("POST (%s) error.code = %q, want invalid_request", tc.name, env.Error.Code)
			}
		})
	}

	// 4) PATCH /api/keys/{id} validation errors return 400.
	patchBadCases := []struct {
		name string
		body string
	}{
		{"empty patch body", `{}`},
		{"empty name", `{"name":""}`},
		{"whitespace name", `{"name":"   "}`},
		{"negative rate_limit_rpm", `{"rate_limit_rpm":-10}`},
		{"negative quota_usd", `{"quota_usd":-5}`},
		{"invalid expires_at", `{"expires_at":"tomorrow"}`},
		{"empty expires_at string", `{"expires_at":""}`},
		{"empty scopes entry", `{"scopes":["  "]}`},
		{"unknown field", `{"disabled":true,"unknown":1}`},
	}
	for _, tc := range patchBadCases {
		t.Run("PATCH_"+tc.name, func(t *testing.T) {
			rec := h.doJSON(t, http.MethodPatch, "/api/keys/"+expired.Key.ID, tc.body)
			if rec.Code != http.StatusBadRequest {
				t.Fatalf("PATCH (%s) status = %d, want 400 (body=%s)", tc.name, rec.Code, rec.Body.String())
			}
			var env httpserver.ErrorEnvelope
			if err := json.Unmarshal(rec.Body.Bytes(), &env); err != nil {
				t.Fatalf("unmarshal error envelope: %v", err)
			}
			if env.Error.Code != "invalid_request" {
				t.Errorf("PATCH (%s) error.code = %q, want invalid_request", tc.name, env.Error.Code)
			}
		})
	}
}

type httpStatusError interface {
	StatusCode() int
}

func TestValidateKeyUpdatesLastUsedAt(t *testing.T) {
	h := newTestHarness(t, nil)
	ctx := context.Background()

	rec := h.doJSON(t, http.MethodPost, "/api/keys", `{"name":"last-used-check"}`)
	if rec.Code != http.StatusCreated {
		t.Fatalf("POST /api/keys status = %d: %s", rec.Code, rec.Body.String())
	}
	var created keySecretResp
	if err := json.Unmarshal(rec.Body.Bytes(), &created); err != nil {
		t.Fatalf("unmarshal create: %v", err)
	}
	if created.Key.LastUsedAt != nil {
		t.Fatalf("initial last_used_at = %v, want nil", created.Key.LastUsedAt)
	}

	if _, err := h.mod.ValidateKey(ctx, created.Secret); err != nil {
		t.Fatalf("ValidateKey: %v", err)
	}

	getRec := h.doJSON(t, http.MethodGet, "/api/keys/"+created.Key.ID, "")
	if getRec.Code != http.StatusOK {
		t.Fatalf("GET /api/keys/{id} status = %d: %s", getRec.Code, getRec.Body.String())
	}
	var loaded keys.APIKeySummary
	if err := json.Unmarshal(getRec.Body.Bytes(), &loaded); err != nil {
		t.Fatalf("unmarshal get: %v", err)
	}
	if loaded.LastUsedAt == nil || strings.TrimSpace(*loaded.LastUsedAt) == "" {
		t.Fatalf("expected last_used_at to be populated after ValidateKey, got nil")
	}
	if _, err := time.Parse(time.RFC3339Nano, *loaded.LastUsedAt); err != nil {
		t.Fatalf("invalid RFC3339 last_used_at %q: %v", *loaded.LastUsedAt, err)
	}
}

func TestKeyInterceptorEnforcesRPMAndQuota(t *testing.T) {
	var _ contracts.GatewayInterceptor = (*keys.Module)(nil)

	h := newTestHarness(t, nil)
	ctx := context.Background()

	rec := h.doJSON(t, http.MethodPost, "/api/keys", `{
		"name": "limited-client",
		"rate_limit_rpm": 2,
		"quota_usd": 0.50
	}`)
	if rec.Code != http.StatusCreated {
		t.Fatalf("POST /api/keys status = %d: %s", rec.Code, rec.Body.String())
	}
	var created keySecretResp
	if err := json.Unmarshal(rec.Body.Bytes(), &created); err != nil {
		t.Fatalf("unmarshal create: %v", err)
	}

	// Internal studio/probe keys must never be blocked by user key limits.
	for _, internalID := range []string{"", "admin-internal", "admin-probe"} {
		if err := h.mod.OnRequest(ctx, &contracts.CanonicalRequest{KeyID: internalID}); err != nil {
			t.Fatalf("OnRequest(%q) unexpected error: %v", internalID, err)
		}
	}

	// 1) RPM enforcement: 2 allowed, 3rd rejected with HTTP 429.
	for i := range 2 {
		if err := h.mod.OnRequest(ctx, &contracts.CanonicalRequest{KeyID: created.Key.ID}); err != nil {
			t.Fatalf("OnRequest attempt %d unexpected error: %v", i+1, err)
		}
	}
	rpmErr := h.mod.OnRequest(ctx, &contracts.CanonicalRequest{KeyID: created.Key.ID})
	if rpmErr == nil {
		t.Fatalf("3rd OnRequest within 1 minute succeeded, want 429 rate limit error")
	}
	var statusErr httpStatusError
	if !errors.As(rpmErr, &statusErr) || statusErr.StatusCode() != http.StatusTooManyRequests {
		t.Fatalf("RPM error = %v, want HTTP 429 statusError", rpmErr)
	}
	if !strings.Contains(strings.ToLower(rpmErr.Error()), "rate limit") {
		t.Errorf("RPM error message %q should mention rate limit", rpmErr.Error())
	}

	// 2) Lift RPM cap via PATCH -> request passes immediately.
	patchRec := h.doJSON(t, http.MethodPatch, "/api/keys/"+created.Key.ID, `{"rate_limit_rpm": 0}`)
	if patchRec.Code != http.StatusOK {
		t.Fatalf("PATCH rpm=0 status = %d: %s", patchRec.Code, patchRec.Body.String())
	}
	if err := h.mod.OnRequest(ctx, &contracts.CanonicalRequest{KeyID: created.Key.ID}); err != nil {
		t.Fatalf("OnRequest after lifting RPM failed: %v", err)
	}

	// 3) Live usage.recorded event exceeds $0.50 quota -> next request blocked with HTTP 429.
	h.bus.Publish("usage.recorded", map[string]any{
		"key_id":   created.Key.ID,
		"cost_usd": 0.55,
	})
	deadline := time.Now().Add(2 * time.Second)
	var quotaErr error
	for time.Now().Before(deadline) {
		quotaErr = h.mod.OnRequest(ctx, &contracts.CanonicalRequest{KeyID: created.Key.ID})
		if quotaErr != nil {
			break
		}
		time.Sleep(15 * time.Millisecond)
	}
	if quotaErr == nil {
		t.Fatalf("expected quota exhaustion error after usage.recorded exceeded $0.50")
	}
	if !errors.As(quotaErr, &statusErr) || statusErr.StatusCode() != http.StatusTooManyRequests {
		t.Fatalf("quota error = %v, want HTTP 429 statusError", quotaErr)
	}
	if !strings.Contains(strings.ToLower(quotaErr.Error()), "quota") {
		t.Errorf("quota error message %q should mention quota", quotaErr.Error())
	}

	// 4) Raise quota via PATCH -> requests pass again.
	raiseRec := h.doJSON(t, http.MethodPatch, "/api/keys/"+created.Key.ID, `{"quota_usd": 1.00}`)
	if raiseRec.Code != http.StatusOK {
		t.Fatalf("PATCH quota_usd=1.00 status = %d: %s", raiseRec.Code, raiseRec.Body.String())
	}
	if err := h.mod.OnRequest(ctx, &contracts.CanonicalRequest{KeyID: created.Key.ID}); err != nil {
		t.Fatalf("OnRequest after raising quota failed: %v", err)
	}
}

func TestKeyInterceptorSeedsQuotaFromUsageLogsOnRestart(t *testing.T) {
	h := newTestHarness(t, nil)
	ctx := context.Background()

	// Create usage_logs table as persisted from a previous run.
	if _, err := h.db.Writer().ExecContext(ctx, `
		CREATE TABLE IF NOT EXISTS usage_logs (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			key_id TEXT NOT NULL DEFAULT '',
			cost_usd REAL NOT NULL DEFAULT 0
		)
	`); err != nil {
		t.Fatalf("create usage_logs: %v", err)
	}

	rec := h.doJSON(t, http.MethodPost, "/api/keys", `{
		"name": "pre-spent-key",
		"quota_usd": 0.50
	}`)
	if rec.Code != http.StatusCreated {
		t.Fatalf("POST /api/keys status = %d: %s", rec.Code, rec.Body.String())
	}
	var created keySecretResp
	if err := json.Unmarshal(rec.Body.Bytes(), &created); err != nil {
		t.Fatalf("unmarshal create: %v", err)
	}

	if _, err := h.db.Writer().ExecContext(ctx, `INSERT INTO usage_logs (key_id, cost_usd) VALUES (?, ?)`, created.Key.ID, 0.75); err != nil {
		t.Fatalf("insert prior spend: %v", err)
	}

	err := h.mod.OnRequest(ctx, &contracts.CanonicalRequest{KeyID: created.Key.ID})
	if err == nil {
		t.Fatalf("expected 429 quota error seeded from historical usage_logs spend ($0.75 >= $0.50)")
	}
	var statusErr httpStatusError
	if !errors.As(err, &statusErr) || statusErr.StatusCode() != http.StatusTooManyRequests {
		t.Fatalf("seeded quota err = %v, want HTTP 429", err)
	}
}
