// Package usage implements async batched SQLite request telemetry logging,
// model pricing calculations and overrides, usage summary aggregations,
// real-time SSE telemetry streaming, and retention pruning.
package usage

import (
	"context"
	"embed"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/fs"
	"log/slog"
	"math"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"sync"
	"time"
	"unicode/utf8"

	"github.com/go-chi/chi/v5"

	"github.com/omeg4monx/router-gateway/internal/contracts"
	"github.com/omeg4monx/router-gateway/internal/platform/db"
	"github.com/omeg4monx/router-gateway/internal/platform/eventbus"
	"github.com/omeg4monx/router-gateway/internal/platform/httpserver"
	"github.com/omeg4monx/router-gateway/internal/platform/lifecycle"
)

const (
	moduleName = "usage"

	// TopicUsageRecorded is published on the event bus for each persisted
	// usage telemetry entry.
	TopicUsageRecorded = "usage.recorded"
	// TopicTrafficStarted is published immediately when a gateway request begins.
	TopicTrafficStarted = "traffic.started"
	// TopicTrafficFinished is published immediately when a gateway request finishes
	// (before async SQLite batch persistence) for 0ms live monitoring updates.
	TopicTrafficFinished = "traffic.finished"

	defaultBatchSize     = 50
	defaultFlushInterval = 250 * time.Millisecond
	defaultBufferSize    = 2048
	defaultRetentionDays = 90
	defaultLogsLimit     = 50
	maxLogsLimit         = 200
	maxUsageBodyBytes    = 64 << 10 // 64 KiB
	maxModelRunes        = 128
	sseSubscriberBufSize = 128
	sseKeepaliveInterval = 15 * time.Second
)

//go:embed migrations/*.sql
var migrationsFS embed.FS

// Options configures a new usage Module instance.
type Options struct {
	DB                    *db.DB
	Bus                   *eventbus.Bus
	Logger                *slog.Logger
	BatchSize             int
	FlushInterval         time.Duration
	BufferSize            int
	MetadataRetentionDays int
	Now                   func() time.Time
}

// Module implements lifecycle.Module and contracts.UsageRecorder.
type Module struct {
	db            *db.DB
	bus           *eventbus.Bus
	logger        *slog.Logger
	now           func() time.Time
	retentionDays int
	repo          *repository
	pricing       *pricingStore
	batcher       *batcher

	stopMu     sync.Mutex
	stopNotify chan struct{}
}

var (
	_ lifecycle.Module        = (*Module)(nil)
	_ contracts.UsageRecorder = (*Module)(nil)
)

// New constructs a usage Module with defaults applied for unset options.
func New(opts Options) *Module {
	clock := opts.Now
	if clock == nil {
		clock = func() time.Time { return time.Now().UTC() }
	}
	logger := opts.Logger
	if logger == nil {
		logger = slog.Default()
	}
	batchSize := opts.BatchSize
	if batchSize <= 0 {
		batchSize = defaultBatchSize
	}
	flushInterval := opts.FlushInterval
	if flushInterval <= 0 {
		flushInterval = defaultFlushInterval
	}
	bufferSize := opts.BufferSize
	if bufferSize <= 0 {
		bufferSize = defaultBufferSize
	}
	retentionDays := opts.MetadataRetentionDays
	if retentionDays <= 0 {
		retentionDays = defaultRetentionDays
	}

	repo := newRepository(opts.DB, clock)
	pricing := newPricingStore()

	return &Module{
		db:            opts.DB,
		bus:           opts.Bus,
		logger:        logger,
		now:           clock,
		retentionDays: retentionDays,
		repo:          repo,
		pricing:       pricing,
		batcher:       newBatcher(repo, pricing, opts.Bus, logger, clock, batchSize, flushInterval, bufferSize),
	}
}

// Name returns the canonical module name.
func (m *Module) Name() string {
	return moduleName
}

// Dependencies declares that usage depends on system and settings.
func (m *Module) Dependencies() []string {
	return []string{"system", "settings"}
}

// Migrations returns the embedded SQL migration filesystem rooted at migrations/.
func (m *Module) Migrations() fs.FS {
	sub, err := fs.Sub(migrationsFS, "migrations")
	if err != nil {
		panic(fmt.Sprintf("usage: load embedded migrations: %v", err))
	}
	return sub
}

// RegisterRoutes mounts the administrative usage endpoints onto AdminAPI and
// mounts the long-lived SSE stream onto AdminStreamAPI so APITimeout does not
// prematurely terminate live subscriptions.
func (m *Module) RegisterRoutes(mounts lifecycle.RouteMounts) {
	if mounts.AdminAPI != nil {
		mounts.AdminAPI.Get("/usage/summary", m.handleGetSummary)
		mounts.AdminAPI.Get("/usage/logs", m.handleListLogs)
		mounts.AdminAPI.Get("/usage/pricing", m.handleListPricing)
		mounts.AdminAPI.Put("/usage/pricing/{model}", m.handleUpsertPricing)
		mounts.AdminAPI.Delete("/usage/pricing/{model}", m.handleDeletePricing)
	}
	streamMount := mounts.AdminStreamAPI
	if streamMount == nil {
		streamMount = mounts.AdminAPI
	}
	if streamMount != nil {
		streamMount.Get("/usage/stream", m.handleStreamUsage)
	}
}

// Start loads custom model pricing overrides and starts the async batch writer.
func (m *Module) Start(ctx context.Context) error {
	if m.db == nil {
		return errors.New("usage: database is required")
	}

	overrides, err := m.repo.ListPricingOverrides(ctx)
	if err != nil {
		return fmt.Errorf("usage: load model pricing overrides: %w", err)
	}
	m.pricing.SetOverrides(overrides)

	m.stopMu.Lock()
	m.stopNotify = make(chan struct{})
	m.stopMu.Unlock()

	m.batcher.Start()
	return nil
}

// Stop signals SSE stream listeners to close, drains all buffered usage records
// to SQLite, and clears the in-memory pricing cache.
func (m *Module) Stop(ctx context.Context) error {
	m.stopMu.Lock()
	if m.stopNotify != nil {
		close(m.stopNotify)
		m.stopNotify = nil
	}
	m.stopMu.Unlock()

	err := m.batcher.Stop(ctx)
	m.pricing.Clear()
	return err
}

// Record enqueues a request telemetry record without blocking the caller and
// immediately publishes a live traffic.finished event on the event bus.
// Records emitted after Stop or when the buffer is full are safely dropped.
func (m *Module) Record(_ context.Context, rec contracts.UsageRecord) {
	if m.bus != nil {
		ts := rec.Timestamp
		if ts.IsZero() {
			ts = m.now()
		}
		inTok := max(rec.InputTokens, 0)
		outTok := max(rec.OutputTokens, 0)
		cachedTok := max(rec.CachedTokens, 0)
		cost := rec.CostUSD
		if cost < 0 {
			cost = 0
		} else if cost == 0 {
			cost = m.pricing.CalculateCost(rec.Model, inTok, outTok, cachedTok)
		}
		m.bus.Publish(TopicTrafficFinished, UsageLog{
			RequestID:    rec.RequestID,
			KeyID:        rec.KeyID,
			Model:        rec.Model,
			ProviderID:   rec.ProviderID,
			ConnectionID: rec.ConnectionID,
			Protocol:     rec.Protocol,
			InputTokens:  inTok,
			OutputTokens: outTok,
			CachedTokens: cachedTok,
			CostUSD:      cost,
			LatencyMS:    max(rec.Latency.Milliseconds(), 0),
			TTFTMS:       max(rec.TTFT.Milliseconds(), 0),
			StatusCode:   rec.StatusCode,
			ErrorSummary: rec.ErrorSummary,
			Stream:       rec.Stream,
			CreatedAt:    ts.UTC().Format(rfc3339NanoFixed),
		})
	}
	m.batcher.Enqueue(rec)
}

// DroppedCount returns the total number of telemetry records dropped due to
// a full async buffer.
func (m *Module) DroppedCount() int64 {
	return m.batcher.DroppedCount()
}

// PruneExpired deletes usage_logs rows older than MetadataRetentionDays.
func (m *Module) PruneExpired(ctx context.Context) error {
	if m.db == nil {
		return errors.New("usage: database is required")
	}
	cutoff := m.now().UTC().Add(-time.Duration(m.retentionDays) * 24 * time.Hour).Format(rfc3339NanoFixed)
	_, err := m.repo.PruneBefore(ctx, cutoff)
	return err
}

func (m *Module) currentStopNotify() <-chan struct{} {
	m.stopMu.Lock()
	defer m.stopMu.Unlock()
	return m.stopNotify
}

type pricingListResponse struct {
	Pricing []ModelPricing `json:"pricing"`
}

type upsertPricingRequest struct {
	InputPerMillion  *float64 `json:"input_per_million"`
	OutputPerMillion *float64 `json:"output_per_million"`
	CachedPerMillion *float64 `json:"cached_per_million"`
}

func (m *Module) handleGetSummary(w http.ResponseWriter, r *http.Request) {
	q := r.URL.Query()
	rng := strings.TrimSpace(q.Get("range"))
	if rng == "" {
		rng = "24h"
	}

	includeSeries := true
	if rawInclude := strings.TrimSpace(q.Get("include_series")); rawInclude != "" {
		parsed, err := strconv.ParseBool(rawInclude)
		if err != nil {
			httpserver.WriteError(w, r, http.StatusBadRequest, "invalid_request", "include_series must be a boolean")
			return
		}
		includeSeries = parsed
	}

	summary, err := m.repo.Summary(r.Context(), rng, includeSeries)
	if err != nil {
		if _, _, _, winErr := m.repo.summaryWindow(rng); winErr != nil {
			httpserver.WriteError(w, r, http.StatusBadRequest, "invalid_request", winErr.Error())
			return
		}
		httpserver.WriteError(w, r, http.StatusInternalServerError, "internal_error", "failed to load usage summary")
		return
	}
	httpserver.WriteJSON(w, http.StatusOK, summary)
}

func (m *Module) handleListLogs(w http.ResponseWriter, r *http.Request) {
	q := r.URL.Query()
	limit := defaultLogsLimit
	if rawLimit := strings.TrimSpace(q.Get("limit")); rawLimit != "" {
		parsed, err := strconv.Atoi(rawLimit)
		if err != nil || parsed <= 0 {
			httpserver.WriteError(w, r, http.StatusBadRequest, "invalid_request", "limit must be a positive integer")
			return
		}
		if parsed > maxLogsLimit {
			parsed = maxLogsLimit
		}
		limit = parsed
	}

	offset := 0
	if rawOffset := strings.TrimSpace(q.Get("offset")); rawOffset != "" {
		parsed, err := strconv.Atoi(rawOffset)
		if err != nil || parsed < 0 {
			httpserver.WriteError(w, r, http.StatusBadRequest, "invalid_request", "offset must be a non-negative integer")
			return
		}
		offset = parsed
	}

	filter := LogsFilter{
		Range:      strings.TrimSpace(q.Get("range")),
		Search:     strings.TrimSpace(q.Get("q")),
		Model:      strings.TrimSpace(q.Get("model")),
		ProviderID: strings.TrimSpace(q.Get("provider")),
		KeyID:      strings.TrimSpace(q.Get("key_id")),
		Status:     strings.TrimSpace(q.Get("status")),
		Limit:      limit,
		Offset:     offset,
	}
	if _, _, err := m.repo.buildLogsWhere(filter); err != nil {
		httpserver.WriteError(w, r, http.StatusBadRequest, "invalid_request", err.Error())
		return
	}

	resp, err := m.repo.ListLogs(r.Context(), filter)
	if err != nil {
		httpserver.WriteError(w, r, http.StatusInternalServerError, "internal_error", "failed to list usage logs")
		return
	}
	httpserver.WriteJSON(w, http.StatusOK, resp)
}

func (m *Module) handleStreamUsage(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")
	w.Header().Set("X-Accel-Buffering", "no")
	w.WriteHeader(http.StatusOK)

	includeLiveTraffic := r.URL.Query().Get("live") == "1" || strings.EqualFold(r.URL.Query().Get("live"), "true")

	var (
		events         <-chan eventbus.Event
		startedEvents  <-chan eventbus.Event
		finishedEvents <-chan eventbus.Event
		provEvents     <-chan eventbus.Event
		keyEvents      <-chan eventbus.Event
		cmbEvents      <-chan eventbus.Event
		poolEvents     <-chan eventbus.Event
	)
	if m.bus != nil {
		var unsubUsage, unsubStarted, unsubFinished, unsubProv, unsubKey, unsubCmb, unsubPool func()
		events, unsubUsage = m.bus.Subscribe(TopicUsageRecorded, sseSubscriberBufSize)
		if includeLiveTraffic {
			startedEvents, unsubStarted = m.bus.Subscribe(TopicTrafficStarted, sseSubscriberBufSize)
			finishedEvents, unsubFinished = m.bus.Subscribe(TopicTrafficFinished, sseSubscriberBufSize)
		}
		provEvents, unsubProv = m.bus.Subscribe("providers.changed", 16)
		keyEvents, unsubKey = m.bus.Subscribe("keys.changed", 16)
		cmbEvents, unsubCmb = m.bus.Subscribe("combos.changed", 16)
		poolEvents, unsubPool = m.bus.Subscribe("proxypools.changed", 16)
		defer func() {
			unsubUsage()
			if unsubStarted != nil {
				unsubStarted()
			}
			if unsubFinished != nil {
				unsubFinished()
			}
			unsubProv()
			unsubKey()
			unsubCmb()
			unsubPool()
		}()
	}

	rc := http.NewResponseController(w)
	if _, err := io.WriteString(w, ": connected\n\n"); err != nil {
		return
	}
	if err := rc.Flush(); err != nil {
		return
	}

	keepalive := time.NewTicker(sseKeepaliveInterval)
	defer keepalive.Stop()

	stopNotify := m.currentStopNotify()

	writeEvent := func(eventName string, payload any) bool {
		data, err := json.Marshal(payload)
		if err != nil {
			return true
		}
		if _, err := fmt.Fprintf(w, "event: %s\ndata: %s\n\n", eventName, data); err != nil {
			return false
		}
		return rc.Flush() == nil
	}

	writeInvalidation := func(topic string) bool {
		if _, err := fmt.Fprintf(w, "event: invalidation\ndata: {\"topic\":%q}\n\n", topic); err != nil {
			return false
		}
		return rc.Flush() == nil
	}

	for {
		select {
		case <-r.Context().Done():
			return
		case <-stopNotify:
			return
		case evt, ok := <-events:
			if !ok {
				return
			}
			if !writeEvent("usage", evt.Payload) {
				return
			}
		case evt, ok := <-startedEvents:
			if !ok {
				startedEvents = nil
				continue
			}
			if !writeEvent("traffic_started", evt.Payload) {
				return
			}
		case evt, ok := <-finishedEvents:
			if !ok {
				finishedEvents = nil
				continue
			}
			if !writeEvent("traffic_finished", evt.Payload) {
				return
			}
		case _, ok := <-provEvents:
			if !ok {
				provEvents = nil
				continue
			}
			if !writeInvalidation("providers.changed") {
				return
			}
		case _, ok := <-keyEvents:
			if !ok {
				keyEvents = nil
				continue
			}
			if !writeInvalidation("keys.changed") {
				return
			}
		case _, ok := <-cmbEvents:
			if !ok {
				cmbEvents = nil
				continue
			}
			if !writeInvalidation("combos.changed") {
				return
			}
		case _, ok := <-poolEvents:
			if !ok {
				poolEvents = nil
				continue
			}
			if !writeInvalidation("proxypools.changed") {
				return
			}
		case <-keepalive.C:
			if _, err := io.WriteString(w, ": keepalive\n\n"); err != nil {
				return
			}
			if err := rc.Flush(); err != nil {
				return
			}
		}
	}
}

func (m *Module) handleListPricing(w http.ResponseWriter, _ *http.Request) {
	httpserver.WriteJSON(w, http.StatusOK, pricingListResponse{
		Pricing: m.pricing.List(),
	})
}

func (m *Module) handleUpsertPricing(w http.ResponseWriter, r *http.Request) {
	model, err := decodePathModelParam(chi.URLParam(r, "model"))
	if err != nil {
		httpserver.WriteError(w, r, http.StatusBadRequest, "invalid_request", err.Error())
		return
	}

	var req upsertPricingRequest
	if !decodeJSONBody(w, r, &req) {
		return
	}

	if req.InputPerMillion == nil || req.OutputPerMillion == nil {
		httpserver.WriteError(w, r, http.StatusBadRequest, "invalid_request", "input_per_million and output_per_million are required")
		return
	}
	inRate := *req.InputPerMillion
	outRate := *req.OutputPerMillion
	cachedRate := 0.0
	if req.CachedPerMillion != nil {
		cachedRate = *req.CachedPerMillion
	}

	if !isValidRate(inRate) || !isValidRate(outRate) || !isValidRate(cachedRate) {
		httpserver.WriteError(w, r, http.StatusBadRequest, "invalid_request", "pricing rates must be non-negative finite numbers")
		return
	}

	saved, err := m.repo.UpsertPricingOverride(r.Context(), model, inRate, outRate, cachedRate)
	if err != nil {
		httpserver.WriteError(w, r, http.StatusInternalServerError, "internal_error", "failed to save model pricing")
		return
	}

	m.pricing.UpsertOverride(saved)
	httpserver.WriteJSON(w, http.StatusOK, saved)
}

func (m *Module) handleDeletePricing(w http.ResponseWriter, r *http.Request) {
	model, err := decodePathModelParam(chi.URLParam(r, "model"))
	if err != nil {
		httpserver.WriteError(w, r, http.StatusBadRequest, "invalid_request", err.Error())
		return
	}

	if err := m.repo.DeletePricingOverride(r.Context(), model); err != nil {
		if errors.Is(err, ErrPricingNotFound) {
			httpserver.WriteError(w, r, http.StatusNotFound, "pricing_not_found", "custom model pricing override not found")
			return
		}
		httpserver.WriteError(w, r, http.StatusInternalServerError, "internal_error", "failed to delete model pricing")
		return
	}

	m.pricing.DeleteOverride(model)
	w.WriteHeader(http.StatusNoContent)
}

func decodePathModelParam(raw string) (string, error) {
	unescaped, err := url.PathUnescape(raw)
	if err != nil {
		return "", errors.New("model contains invalid URL escape sequence")
	}
	trimmed := strings.TrimSpace(unescaped)
	if trimmed == "" {
		return "", errors.New("model must not be empty")
	}
	if trimmed == "*" {
		return "", errors.New("model must not be wildcard '*'")
	}
	if strings.ContainsAny(trimmed, "\r\n\t") {
		return "", errors.New("model must not contain control whitespace")
	}
	if utf8.RuneCountInString(trimmed) > maxModelRunes {
		return "", fmt.Errorf("model must not exceed %d characters", maxModelRunes)
	}
	return trimmed, nil
}

func isValidRate(v float64) bool {
	return v >= 0 && !math.IsNaN(v) && !math.IsInf(v, 0)
}

func decodeJSONBody(w http.ResponseWriter, r *http.Request, dst any) bool {
	r.Body = http.MaxBytesReader(w, r.Body, maxUsageBodyBytes)
	dec := json.NewDecoder(r.Body)
	dec.DisallowUnknownFields()

	if err := dec.Decode(dst); err != nil {
		var maxBytesErr *http.MaxBytesError
		if errors.As(err, &maxBytesErr) {
			httpserver.WriteError(w, r, http.StatusRequestEntityTooLarge, "request_too_large", "request body exceeds maximum allowed size")
			return false
		}
		httpserver.WriteError(w, r, http.StatusBadRequest, "invalid_request", "invalid JSON request body")
		return false
	}
	if dec.More() {
		httpserver.WriteError(w, r, http.StatusBadRequest, "invalid_request", "request body must contain a single JSON object")
		return false
	}
	if err := dec.Decode(&struct{}{}); !errors.Is(err, io.EOF) {
		httpserver.WriteError(w, r, http.StatusBadRequest, "invalid_request", "request body must contain a single JSON object")
		return false
	}
	return true
}
