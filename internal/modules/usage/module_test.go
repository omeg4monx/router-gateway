package usage_test

import (
	"bufio"
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"io"
	"io/fs"
	"log/slog"
	"math"
	"net/http"
	"net/http/httptest"
	"path/filepath"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/go-chi/chi/v5"

	"github.com/omeg4monx/router-gateway/internal/contracts"
	"github.com/omeg4monx/router-gateway/internal/modules/usage"
	"github.com/omeg4monx/router-gateway/internal/platform/config"
	"github.com/omeg4monx/router-gateway/internal/platform/db"
	"github.com/omeg4monx/router-gateway/internal/platform/eventbus"
	"github.com/omeg4monx/router-gateway/internal/platform/httpserver"
	"github.com/omeg4monx/router-gateway/internal/platform/lifecycle"
)

var (
	_ lifecycle.Module        = (*usage.Module)(nil)
	_ contracts.UsageRecorder = (*usage.Module)(nil)
)

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

type fakeClock struct {
	mu  sync.Mutex
	now time.Time
}

func newFakeClock(t time.Time) *fakeClock {
	return &fakeClock{now: t.UTC()}
}

func (c *fakeClock) Now() time.Time {
	c.mu.Lock()
	defer c.mu.Unlock()
	return c.now
}

func (c *fakeClock) Set(t time.Time) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.now = t.UTC()
}

type testHarness struct {
	db     *db.DB
	bus    *eventbus.Bus
	reg    *lifecycle.Registry
	server *httpserver.Server
	mod    *usage.Module
	clock  *fakeClock
}

func newTestHarness(t *testing.T, opts usage.Options, withAuth bool) *testHarness {
	t.Helper()

	clock := newFakeClock(time.Date(2026, 11, 23, 12, 0, 0, 0, time.UTC))
	dbPath := filepath.Join(t.TempDir(), "router.sqlite")
	database, err := db.Open(dbPath)
	if err != nil {
		t.Fatalf("db.Open: %v", err)
	}
	t.Cleanup(func() {
		_ = database.Close()
	})

	bus := eventbus.New()
	opts.DB = database
	opts.Bus = bus
	if opts.Logger == nil {
		opts.Logger = slog.New(slog.NewTextHandler(io.Discard, nil))
	}
	if opts.Now == nil {
		opts.Now = clock.Now
	}

	mod := usage.New(opts)
	reg := lifecycle.NewRegistry(nil)
	if err := reg.Register(stubModule{name: "system"}); err != nil {
		t.Fatalf("reg.Register(system): %v", err)
	}
	if err := reg.Register(stubModule{name: "settings"}); err != nil {
		t.Fatalf("reg.Register(settings): %v", err)
	}
	if err := reg.Register(mod); err != nil {
		t.Fatalf("reg.Register(usage): %v", err)
	}

	ctx := context.Background()
	if err := reg.MigrateAll(ctx, database); err != nil {
		t.Fatalf("reg.MigrateAll: %v", err)
	}

	srvOpts := httpserver.Options{
		Config: config.Config{
			Host:       "127.0.0.1",
			Port:       20128,
			APITimeout: 30 * time.Second,
		},
		Logger: slog.New(slog.NewTextHandler(io.Discard, nil)),
	}
	if withAuth {
		srvOpts.Auth = &stubAdminAuth{subject: "admin"}
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
		clock:  clock,
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

func waitForLogCount(t *testing.T, database *db.DB, want int64) {
	t.Helper()
	deadline := time.Now().Add(3 * time.Second)
	for time.Now().Before(deadline) {
		var count int64
		err := database.Reader().QueryRow(`SELECT COUNT(*) FROM usage_logs`).Scan(&count)
		if err == nil && count >= want {
			return
		}
		time.Sleep(15 * time.Millisecond)
	}
	var got int64
	_ = database.Reader().QueryRow(`SELECT COUNT(*) FROM usage_logs`).Scan(&got)
	t.Fatalf("timed out waiting for %d usage_logs rows, got %d", want, got)
}

func TestUsageAsyncBatcherAndGracefulDrain(t *testing.T) {
	h := newTestHarness(t, usage.Options{
		BatchSize:     20,
		FlushInterval: 5 * time.Second, // long interval forces Stop() drain to flush remainder
		BufferSize:    256,
	}, false)

	const totalEvents = 75
	var wg sync.WaitGroup
	for i := range totalEvents {
		wg.Add(1)
		go func(idx int) {
			defer wg.Done()
			h.mod.Record(context.Background(), contracts.UsageRecord{
				RequestID:    "req_" + string(rune('a'+idx%26)),
				KeyID:        "key_1",
				Model:        "gpt-4o",
				ProviderID:   "openai",
				ConnectionID: "conn_1",
				Protocol:     "openai_chat",
				InputTokens:  100,
				OutputTokens: 50,
				Latency:      120 * time.Millisecond,
				TTFT:         30 * time.Millisecond,
				StatusCode:   200,
				Timestamp:    h.clock.Now(),
			})
		}(i)
	}
	wg.Wait()

	if err := h.reg.StopAll(context.Background()); err != nil {
		t.Fatalf("reg.StopAll: %v", err)
	}

	var count int64
	if err := h.db.Reader().QueryRow(`SELECT COUNT(*) FROM usage_logs`).Scan(&count); err != nil {
		t.Fatalf("count usage_logs: %v", err)
	}
	if count != totalEvents {
		t.Fatalf("persisted rows = %d, want %d", count, totalEvents)
	}

	// Late Record calls after Stop must be safe non-panicking no-ops.
	var lateWG sync.WaitGroup
	for range 10 {
		lateWG.Add(1)
		go func() {
			defer lateWG.Done()
			h.mod.Record(context.Background(), contracts.UsageRecord{
				RequestID:  "req_after_stop",
				Model:      "gpt-4o",
				StatusCode: 200,
			})
		}()
	}
	lateWG.Wait()

	var afterCount int64
	if err := h.db.Reader().QueryRow(`SELECT COUNT(*) FROM usage_logs`).Scan(&afterCount); err != nil {
		t.Fatalf("count after stop: %v", err)
	}
	if afterCount != totalEvents {
		t.Fatalf("rows after stop = %d, want %d", afterCount, totalEvents)
	}
}

func TestUsageCostCalculationAndPricingCRUD(t *testing.T) {
	h := newTestHarness(t, usage.Options{
		BatchSize:     1,
		FlushInterval: 20 * time.Millisecond,
	}, true)

	// Unauthenticated pricing call must return 401.
	if rec := h.doJSON(t, http.MethodGet, "/api/usage/pricing", nil, false); rec.Code != http.StatusUnauthorized {
		t.Fatalf("unauthenticated GET /api/usage/pricing = %d, want 401", rec.Code)
	}

	// Built-in pricing list includes popular models.
	listRec := h.doJSON(t, http.MethodGet, "/api/usage/pricing", nil, true)
	if listRec.Code != http.StatusOK {
		t.Fatalf("GET /api/usage/pricing = %d: %s", listRec.Code, listRec.Body.String())
	}
	pricingResp := decodeJSON[struct {
		Pricing []usage.ModelPricing `json:"pricing"`
	}](t, listRec)

	byModel := make(map[string]usage.ModelPricing, len(pricingResp.Pricing))
	for _, p := range pricingResp.Pricing {
		byModel[p.Model] = p
	}
	for _, wantModel := range []string{"gpt-4o", "claude-sonnet-4", "gemini-2.5-pro"} {
		if _, ok := byModel[wantModel]; !ok {
			t.Fatalf("built-in pricing missing %q", wantModel)
		}
	}

	// Record events with CostUSD == 0 to trigger built-in calculation (including <provider>/ prefix strip).
	h.mod.Record(context.Background(), contracts.UsageRecord{
		RequestID:    "req_openai",
		Model:        "openai/gpt-4o",
		ProviderID:   "openai",
		InputTokens:  1_000_000,
		OutputTokens: 500_000,
		CachedTokens: 200_000,
		StatusCode:   200,
		Timestamp:    h.clock.Now(),
	})
	h.mod.Record(context.Background(), contracts.UsageRecord{
		RequestID:    "req_claude",
		Model:        "claude-sonnet-4",
		ProviderID:   "anthropic",
		InputTokens:  1_000_000,
		OutputTokens: 1_000_000,
		StatusCode:   200,
		Timestamp:    h.clock.Now(),
	})
	h.mod.Record(context.Background(), contracts.UsageRecord{
		RequestID:    "req_gemini",
		Model:        "gemini-2.5-pro",
		ProviderID:   "gemini",
		InputTokens:  1_000_000,
		OutputTokens: 1_000_000,
		StatusCode:   200,
		Timestamp:    h.clock.Now(),
	})
	h.mod.Record(context.Background(), contracts.UsageRecord{
		RequestID:    "req_explicit",
		Model:        "gpt-4o",
		ProviderID:   "openai",
		InputTokens:  1000,
		OutputTokens: 1000,
		CostUSD:      0.42,
		StatusCode:   200,
		Timestamp:    h.clock.Now(),
	})

	waitForLogCount(t, h.db, 4)

	logsRec := h.doJSON(t, http.MethodGet, "/api/usage/logs?limit=10", nil, true)
	if logsRec.Code != http.StatusOK {
		t.Fatalf("GET /api/usage/logs = %d: %s", logsRec.Code, logsRec.Body.String())
	}
	logsResp := decodeJSON[struct {
		Logs []usage.UsageLog `json:"logs"`
	}](t, logsRec)

	costByReq := make(map[string]float64, len(logsResp.Logs))
	for _, l := range logsResp.Logs {
		costByReq[l.RequestID] = l.CostUSD
	}

	gpt4o := byModel["gpt-4o"]
	wantOpenAI := gpt4o.InputPerMillion + 0.5*gpt4o.OutputPerMillion + 0.2*gpt4o.CachedPerMillion
	if math.Abs(costByReq["req_openai"]-wantOpenAI) > 1e-6 {
		t.Fatalf("req_openai cost = %v, want %v", costByReq["req_openai"], wantOpenAI)
	}
	claude := byModel["claude-sonnet-4"]
	wantClaude := claude.InputPerMillion + claude.OutputPerMillion
	if math.Abs(costByReq["req_claude"]-wantClaude) > 1e-6 {
		t.Fatalf("req_claude cost = %v, want %v", costByReq["req_claude"], wantClaude)
	}
	gemini := byModel["gemini-2.5-pro"]
	wantGemini := gemini.InputPerMillion + gemini.OutputPerMillion
	if math.Abs(costByReq["req_gemini"]-wantGemini) > 1e-6 {
		t.Fatalf("req_gemini cost = %v, want %v", costByReq["req_gemini"], wantGemini)
	}
	if math.Abs(costByReq["req_explicit"]-0.42) > 1e-6 {
		t.Fatalf("req_explicit cost = %v, want 0.42", costByReq["req_explicit"])
	}

	// Custom PUT /api/usage/pricing/{model} overrides calculation.
	putRec := h.doJSON(t, http.MethodPut, "/api/usage/pricing/openai%2Fgpt-4o", map[string]any{
		"input_per_million":  10.0,
		"output_per_million": 20.0,
		"cached_per_million": 2.0,
	}, true)
	if putRec.Code != http.StatusOK {
		t.Fatalf("PUT /api/usage/pricing/openai%%2Fgpt-4o = %d: %s", putRec.Code, putRec.Body.String())
	}
	updatedPricing := decodeJSON[usage.ModelPricing](t, putRec)
	if updatedPricing.Model != "openai/gpt-4o" || updatedPricing.Source != "custom" || updatedPricing.InputPerMillion != 10.0 {
		t.Fatalf("unexpected updated pricing: %+v", updatedPricing)
	}

	// Negative price or unknown JSON field returns 400.
	if rec := h.doJSON(t, http.MethodPut, "/api/usage/pricing/gpt-4o", map[string]any{
		"input_per_million":  -1.0,
		"output_per_million": 1.0,
	}, true); rec.Code != http.StatusBadRequest {
		t.Fatalf("negative price PUT = %d, want 400", rec.Code)
	}
	if rec := h.doJSON(t, http.MethodPut, "/api/usage/pricing/gpt-4o", `{"input_per_million":1,"output_per_million":1,"extra":true}`, true); rec.Code != http.StatusBadRequest {
		t.Fatalf("unknown field PUT = %d, want 400", rec.Code)
	}

	h.mod.Record(context.Background(), contracts.UsageRecord{
		RequestID:    "req_custom_override",
		Model:        "openai/gpt-4o",
		ProviderID:   "openai",
		InputTokens:  1_000_000,
		OutputTokens: 500_000,
		StatusCode:   200,
		Timestamp:    h.clock.Now(),
	})
	waitForLogCount(t, h.db, 5)

	logsRec2 := h.doJSON(t, http.MethodGet, "/api/usage/logs?limit=1", nil, true)
	logsResp2 := decodeJSON[struct {
		Logs []usage.UsageLog `json:"logs"`
	}](t, logsRec2)
	if len(logsResp2.Logs) != 1 || math.Abs(logsResp2.Logs[0].CostUSD-20.0) > 1e-6 {
		t.Fatalf("custom override cost = %+v, want 20.0", logsResp2.Logs)
	}

	// DELETE resets override; second DELETE returns 404.
	if rec := h.doJSON(t, http.MethodDelete, "/api/usage/pricing/openai%2Fgpt-4o", nil, true); rec.Code != http.StatusNoContent {
		t.Fatalf("DELETE pricing = %d, want 204: %s", rec.Code, rec.Body.String())
	}
	if rec := h.doJSON(t, http.MethodDelete, "/api/usage/pricing/openai%2Fgpt-4o", nil, true); rec.Code != http.StatusNotFound {
		t.Fatalf("second DELETE pricing = %d, want 404", rec.Code)
	}
}

func TestUsageSummaryLogsAndRetentionPrune(t *testing.T) {
	h := newTestHarness(t, usage.Options{
		BatchSize:             1,
		FlushInterval:         20 * time.Millisecond,
		MetadataRetentionDays: 90,
	}, false)

	now := h.clock.Now()
	// 1. Recent success (2h ago)
	h.mod.Record(context.Background(), contracts.UsageRecord{
		RequestID:    "req_1",
		KeyID:        "key_alpha",
		Model:        "gpt-4o",
		ProviderID:   "openai",
		ConnectionID: "conn_1",
		Protocol:     "openai_chat",
		InputTokens:  1000,
		OutputTokens: 500,
		CachedTokens: 100,
		CostUSD:      0.50,
		Latency:      200 * time.Millisecond,
		TTFT:         50 * time.Millisecond,
		StatusCode:   200,
		Timestamp:    now.Add(-2 * time.Hour),
	})
	// 2. Recent error (1h ago)
	h.mod.Record(context.Background(), contracts.UsageRecord{
		RequestID:    "req_2",
		KeyID:        "key_beta",
		Model:        "claude-sonnet-4",
		ProviderID:   "anthropic",
		ConnectionID: "conn_2",
		Protocol:     "anthropic",
		InputTokens:  400,
		OutputTokens: 0,
		CostUSD:      0.10,
		Latency:      400 * time.Millisecond,
		TTFT:         0,
		StatusCode:   502,
		ErrorSummary: "upstream bad gateway",
		Timestamp:    now.Add(-1 * time.Hour),
	})
	// 3. 3-day-old success (included in 7d & all, excluded from 24h)
	h.mod.Record(context.Background(), contracts.UsageRecord{
		RequestID:    "req_3",
		KeyID:        "key_alpha",
		Model:        "gpt-4o",
		ProviderID:   "openai",
		ConnectionID: "conn_1",
		Protocol:     "openai_chat",
		InputTokens:  2000,
		OutputTokens: 1000,
		CostUSD:      1.00,
		Latency:      300 * time.Millisecond,
		TTFT:         100 * time.Millisecond,
		StatusCode:   200,
		Timestamp:    now.Add(-72 * time.Hour),
	})
	// 4. 100-day-old row (expired under 90-day retention)
	h.mod.Record(context.Background(), contracts.UsageRecord{
		RequestID:    "req_expired",
		KeyID:        "key_old",
		Model:        "gpt-4o-mini",
		ProviderID:   "openai",
		ConnectionID: "conn_1",
		Protocol:     "openai_chat",
		InputTokens:  500,
		OutputTokens: 200,
		CostUSD:      0.01,
		Latency:      100 * time.Millisecond,
		StatusCode:   200,
		Timestamp:    now.Add(-100 * 24 * time.Hour),
	})

	waitForLogCount(t, h.db, 4)

	// Summary 24h
	sum24 := h.doJSON(t, http.MethodGet, "/api/usage/summary?range=24h", nil, false)
	if sum24.Code != http.StatusOK {
		t.Fatalf("GET /api/usage/summary?range=24h = %d: %s", sum24.Code, sum24.Body.String())
	}
	s24 := decodeJSON[usage.SummaryResponse](t, sum24)
	if s24.Range != "24h" || s24.Totals.Requests != 2 || s24.Totals.ErrorCount != 1 {
		t.Fatalf("unexpected 24h summary totals: %+v", s24)
	}
	if math.Abs(s24.Totals.CostUSD-0.60) > 1e-6 {
		t.Fatalf("24h cost = %v, want 0.60", s24.Totals.CostUSD)
	}
	if math.Abs(s24.Totals.ErrorRate-0.5) > 1e-6 {
		t.Fatalf("24h error_rate = %v, want 0.5", s24.Totals.ErrorRate)
	}
	if len(s24.ByModel) != 2 || len(s24.ByProvider) != 2 || len(s24.Series) != 2 {
		t.Fatalf("unexpected 24h breakdowns: models=%d providers=%d series=%d", len(s24.ByModel), len(s24.ByProvider), len(s24.Series))
	}

	// Summary 7d
	sum7d := decodeJSON[usage.SummaryResponse](t, h.doJSON(t, http.MethodGet, "/api/usage/summary?range=7d", nil, false))
	if sum7d.Totals.Requests != 3 || math.Abs(sum7d.Totals.CostUSD-1.60) > 1e-6 {
		t.Fatalf("unexpected 7d summary: %+v", sum7d.Totals)
	}

	// Invalid range
	if rec := h.doJSON(t, http.MethodGet, "/api/usage/summary?range=99y", nil, false); rec.Code != http.StatusBadRequest {
		t.Fatalf("invalid range = %d, want 400", rec.Code)
	}

	// Logs filtering: model, provider, status
	byModel := decodeJSON[usage.LogsResponse](t, h.doJSON(t, http.MethodGet, "/api/usage/logs?model=gpt-4o&provider=openai", nil, false))
	if byModel.Total != 2 || len(byModel.Logs) != 2 {
		t.Fatalf("filtered logs by model/provider = total %d len %d, want 2", byModel.Total, len(byModel.Logs))
	}
	byStatusErr := decodeJSON[usage.LogsResponse](t, h.doJSON(t, http.MethodGet, "/api/usage/logs?status=error", nil, false))
	if byStatusErr.Total != 1 || byStatusErr.Logs[0].RequestID != "req_2" {
		t.Fatalf("filtered logs by status=error = %+v", byStatusErr)
	}
	byStatusCode := decodeJSON[usage.LogsResponse](t, h.doJSON(t, http.MethodGet, "/api/usage/logs?status=502", nil, false))
	if byStatusCode.Total != 1 || byStatusCode.Logs[0].RequestID != "req_2" {
		t.Fatalf("filtered logs by status=502 = %+v", byStatusCode)
	}

	// Retention prune deletes only expired rows (> 90 days old).
	if err := h.mod.PruneExpired(context.Background()); err != nil {
		t.Fatalf("PruneExpired: %v", err)
	}

	afterPrune := decodeJSON[usage.LogsResponse](t, h.doJSON(t, http.MethodGet, "/api/usage/logs?limit=50", nil, false))
	if afterPrune.Total != 3 {
		t.Fatalf("logs after PruneExpired = %d, want 3", afterPrune.Total)
	}
	for _, l := range afterPrune.Logs {
		if l.RequestID == "req_expired" {
			t.Fatalf("expired log was not pruned")
		}
	}
}

func TestUsageLiveSSEStream(t *testing.T) {
	h := newTestHarness(t, usage.Options{
		BatchSize:     1,
		FlushInterval: 20 * time.Millisecond,
	}, true)

	ts := httptest.NewServer(h.server.Handler())
	defer ts.Close()

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, ts.URL+"/api/usage/stream", nil)
	if err != nil {
		t.Fatalf("NewRequest: %v", err)
	}
	req.Header.Set("Authorization", "Bearer test-admin-token")

	resp, err := ts.Client().Do(req)
	if err != nil {
		t.Fatalf("GET /api/usage/stream: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Fatalf("stream status = %d, want 200", resp.StatusCode)
	}
	if ct := resp.Header.Get("Content-Type"); !strings.HasPrefix(ct, "text/event-stream") {
		t.Fatalf("Content-Type = %q, want text/event-stream", ct)
	}
	if cc := resp.Header.Get("Cache-Control"); cc != "no-cache" {
		t.Fatalf("Cache-Control = %q, want no-cache", cc)
	}

	lineCh := make(chan string, 32)
	errCh := make(chan error, 1)
	go func() {
		scanner := bufio.NewScanner(resp.Body)
		for scanner.Scan() {
			lineCh <- scanner.Text()
		}
		errCh <- scanner.Err()
	}()

	// Wait for initial `: connected` comment.
	select {
	case line := <-lineCh:
		if line != ": connected" {
			t.Fatalf("first SSE line = %q, want ': connected'", line)
		}
	case <-ctx.Done():
		t.Fatalf("timed out waiting for initial SSE comment")
	}

	// Emit a record and expect `event: usage` + `data: {...}`.
	h.mod.Record(context.Background(), contracts.UsageRecord{
		RequestID:    "req_sse_1",
		KeyID:        "key_1",
		Model:        "gpt-4o-mini",
		ProviderID:   "openai",
		Protocol:     "openai_chat",
		InputTokens:  1000,
		OutputTokens: 500,
		Latency:      45 * time.Millisecond,
		StatusCode:   200,
		Stream:       true,
		Timestamp:    h.clock.Now(),
	})

	var gotEvent, gotData string
	for gotEvent == "" || gotData == "" {
		select {
		case line := <-lineCh:
			if strings.HasPrefix(line, "event: ") {
				gotEvent = strings.TrimPrefix(line, "event: ")
			}
			if strings.HasPrefix(line, "data: ") {
				gotData = strings.TrimPrefix(line, "data: ")
			}
		case <-ctx.Done():
			t.Fatalf("timed out waiting for SSE usage frame (gotEvent=%q gotData=%q)", gotEvent, gotData)
		}
	}

	if gotEvent != "usage" {
		t.Fatalf("SSE event = %q, want 'usage'", gotEvent)
	}
	var payload usage.UsageLog
	if err := json.Unmarshal([]byte(gotData), &payload); err != nil {
		t.Fatalf("unmarshal SSE data %q: %v", gotData, err)
	}
	if payload.RequestID != "req_sse_1" || payload.Model != "gpt-4o-mini" || !payload.Stream {
		t.Fatalf("unexpected SSE payload: %+v", payload)
	}

	// Cancel client context and verify clean exit.
	cancel()
	select {
	case <-errCh:
	case <-time.After(2 * time.Second):
		t.Fatalf("SSE reader did not exit after context cancel")
	}
}

func TestUsageBufferFullNeverBlocks(t *testing.T) {
	h := newTestHarness(t, usage.Options{
		BatchSize:     100,
		FlushInterval: 10 * time.Second,
		BufferSize:    2,
	}, false)

	start := time.Now()
	for i := range 200 {
		h.mod.Record(context.Background(), contracts.UsageRecord{
			RequestID:  "req_flood_" + string(rune('0'+i%10)),
			Model:      "gpt-4o",
			StatusCode: 200,
			Timestamp:  h.clock.Now(),
		})
	}
	if elapsed := time.Since(start); elapsed > 500*time.Millisecond {
		t.Fatalf("Record blocked on full buffer: took %v", elapsed)
	}
	if h.mod.DroppedCount() == 0 {
		t.Fatalf("expected non-zero DroppedCount under buffer flood with BufferSize=2")
	}
}

func TestUsageStreamMountsOnAdminStreamAPI(t *testing.T) {
	database, err := db.Open(t.TempDir() + "/usage-mount-test.sqlite")
	if err != nil {
		t.Fatalf("open db: %v", err)
	}
	t.Cleanup(func() { _ = database.Close() })

	mod := usage.New(usage.Options{DB: database})
	adminRouter := chi.NewRouter()
	streamRouter := chi.NewRouter()
	mod.RegisterRoutes(lifecycle.RouteMounts{
		AdminAPI:       adminRouter,
		AdminStreamAPI: streamRouter,
	})

	summaryRec := httptest.NewRecorder()
	adminRouter.ServeHTTP(summaryRec, httptest.NewRequest(http.MethodGet, "/usage/summary", nil))
	if summaryRec.Code == http.StatusNotFound {
		t.Fatalf("expected /usage/summary mounted on AdminAPI, got 404")
	}

	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	adminStreamRec := httptest.NewRecorder()
	adminRouter.ServeHTTP(adminStreamRec, httptest.NewRequest(http.MethodGet, "/usage/stream", nil).WithContext(ctx))
	if adminStreamRec.Code != http.StatusNotFound {
		t.Fatalf("expected /usage/stream NOT on AdminAPI when AdminStreamAPI is provided, got %d", adminStreamRec.Code)
	}

	streamRec := httptest.NewRecorder()
	streamRouter.ServeHTTP(streamRec, httptest.NewRequest(http.MethodGet, "/usage/stream", nil).WithContext(ctx))
	if streamRec.Code != http.StatusOK {
		t.Fatalf("expected /usage/stream mounted on AdminStreamAPI to return 200, got %d", streamRec.Code)
	}
}

func TestUsageStreamForwardsInventoryInvalidations(t *testing.T) {
	h := newTestHarness(t, usage.Options{}, true)
	srv := httptest.NewServer(h.server.Handler())
	t.Cleanup(srv.Close)

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, srv.URL+"/api/usage/stream", nil)
	if err != nil {
		t.Fatalf("new request: %v", err)
	}
	req.Header.Set("Authorization", "Bearer test-admin-token")
	resp, err := srv.Client().Do(req)
	if err != nil {
		t.Fatalf("stream GET: %v", err)
	}
	defer func() { _ = resp.Body.Close() }()

	if resp.StatusCode != http.StatusOK {
		t.Fatalf("expected 200 OK, got %d", resp.StatusCode)
	}

	// Publish a providers.changed event on the shared event bus.
	go func() {
		time.Sleep(50 * time.Millisecond)
		h.bus.Publish("providers.changed", map[string]any{"action": "updated"})
	}()

	scanner := bufio.NewScanner(resp.Body)
	gotEvent := false
	gotTopic := false
	for scanner.Scan() {
		line := scanner.Text()
		if line == "event: invalidation" {
			gotEvent = true
		}
		if gotEvent && strings.HasPrefix(line, "data: ") && strings.Contains(line, `"providers.changed"`) {
			gotTopic = true
			break
		}
	}
	if !gotEvent || !gotTopic {
		t.Fatalf("expected SSE invalidation event for providers.changed, gotEvent=%v gotTopic=%v", gotEvent, gotTopic)
	}
}
