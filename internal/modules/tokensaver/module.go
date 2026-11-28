package tokensaver

import (
	"context"
	"embed"
	"encoding/json"
	"fmt"
	"io"
	"io/fs"
	"log/slog"
	"net/http"
	"strconv"
	"strings"
	"sync"
	"sync/atomic"
	"time"

	"github.com/go-chi/chi/v5"

	"github.com/omeg4monx/router-gateway/internal/contracts"
	"github.com/omeg4monx/router-gateway/internal/platform/db"
	"github.com/omeg4monx/router-gateway/internal/platform/eventbus"
	"github.com/omeg4monx/router-gateway/internal/platform/httpserver"
	"github.com/omeg4monx/router-gateway/internal/platform/lifecycle"
)

const (
	moduleName          = "tokensaver"
	interceptorPriority = 30
	maxBodyBytes        = 1 << 20 // 1 MiB
	eventConfigChanged  = "tokensaver.changed"
)

//go:embed migrations/*.sql
var migrationsFS embed.FS

// CompressionStats aggregates Go-native token compression metrics.
type CompressionStats struct {
	RequestsProcessed    int64 `json:"requests_processed"`
	RequestsCompressed   int64 `json:"requests_compressed"`
	BytesSaved           int64 `json:"bytes_saved"`
	EstimatedTokensSaved int64 `json:"estimated_tokens_saved"`
}

// HeadroomStats aggregates prefix-cache affinity metrics.
type HeadroomStats struct {
	CacheHits     int64               `json:"cache_hits"`
	CacheMisses   int64               `json:"cache_misses"`
	HitRate       float64             `json:"hit_rate"`
	ActiveEntries int                 `json:"active_entries"`
	MaxEntries    int                 `json:"max_entries"`
	WarmPrefixes  []WarmPrefixSummary `json:"warm_prefixes"`
}

// StatsResponse is returned by GET /api/token-saver/stats.
type StatsResponse struct {
	Compression CompressionStats `json:"compression"`
	Headroom    HeadroomStats    `json:"headroom"`
}

// PreviewResponse is returned by POST /api/token-saver/preview.
type PreviewResponse struct {
	OriginalBytes      int            `json:"original_bytes"`
	CompressedBytes    int            `json:"compressed_bytes"`
	BytesSaved         int64          `json:"bytes_saved"`
	EstimatedSavedToks int64          `json:"estimated_tokens_saved"`
	CacheKey           string         `json:"cache_key"`
	Body               map[string]any `json:"body"`
}

// Options configures a tokensaver Module instance.
type Options struct {
	DB     *db.DB
	Bus    *eventbus.Bus
	Logger *slog.Logger
}

// Module implements lifecycle.Module and contracts.GatewayInterceptor.
type Module struct {
	repo   *repository
	bus    *eventbus.Bus
	logger *slog.Logger

	mu         sync.RWMutex
	cfgPatchMu sync.Mutex
	cfg        Config
	headroom   *headroomCache

	reqProcessed  atomic.Int64
	reqCompressed atomic.Int64
	bytesSaved    atomic.Int64
	tokensSaved   atomic.Int64
	cacheHits     atomic.Int64
	cacheMisses   atomic.Int64
}

var (
	_ lifecycle.Module             = (*Module)(nil)
	_ contracts.GatewayInterceptor = (*Module)(nil)
)

// New constructs a tokensaver Module.
func New(opts Options) *Module {
	logger := opts.Logger
	if logger == nil {
		logger = slog.Default()
	}
	return &Module{
		repo:     newRepository(opts.DB),
		bus:      opts.Bus,
		logger:   logger,
		headroom: newHeadroomCache(1024, 30*time.Minute),
	}
}

func (m *Module) Name() string           { return moduleName }
func (m *Module) Dependencies() []string { return []string{"system"} }
func (m *Module) Priority() int          { return interceptorPriority }

func (m *Module) Migrations() fs.FS {
	sub, err := fs.Sub(migrationsFS, "migrations")
	if err != nil {
		panic(fmt.Sprintf("tokensaver: sub migrations fs: %v", err))
	}
	return sub
}

func (m *Module) Start(ctx context.Context) error {
	if err := m.repo.ensureDefaults(ctx); err != nil {
		return err
	}
	return m.reloadConfig(ctx)
}

func (m *Module) Stop(context.Context) error { return nil }

func (m *Module) reloadConfig(ctx context.Context) error {
	cfg, err := m.repo.getConfig(ctx)
	if err != nil {
		return err
	}
	m.mu.Lock()
	m.cfg = cfg
	m.mu.Unlock()
	m.headroom.configure(cfg.MaxCacheEntries, time.Duration(cfg.CacheTTLSeconds)*time.Second)
	return nil
}

func (m *Module) snapshotConfig() Config {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return m.cfg
}

// OnRequest applies token compression and Headroom prefix cache lookup
// using a single JSON unmarshal pass over req.Body.
func (m *Module) OnRequest(_ context.Context, req *contracts.CanonicalRequest) error {
	if req == nil || len(req.Body) == 0 {
		return nil
	}
	cfg := m.snapshotConfig()
	if !cfg.Enabled && !cfg.HeadroomEnabled && !cfg.CavemanEnabled && !cfg.PonytailEnabled && !cfg.TrimHistoryEnabled {
		return nil
	}
	if req.Metadata == nil {
		req.Metadata = make(map[string]any)
	}

	bypassSaver, _ := req.Metadata["bypass_tokensaver"].(bool)
	bypassHeadroom, _ := req.Metadata["bypass_headroom"].(bool)
	runSaver := !bypassSaver && (cfg.Enabled || cfg.TrimHistoryEnabled || cfg.CavemanEnabled || cfg.PonytailEnabled)
	runHeadroom := cfg.HeadroomEnabled && !bypassHeadroom
	if !runSaver && !runHeadroom {
		return nil
	}

	var root map[string]any
	if err := json.Unmarshal(req.Body, &root); err != nil || root == nil {
		return nil
	}

	rootDirty := false

	if !bypassSaver {
		if cfg.Enabled {
			m.reqProcessed.Add(1)
			out, savedBytes, savedTokens, mutated := compressRoot(root, req.Body, cfg)
			if mutated {
				req.Body = out
				m.reqCompressed.Add(1)
				m.bytesSaved.Add(savedBytes)
				m.tokensSaved.Add(savedTokens)
				req.Metadata["tokens_saved"] = savedTokens
			} else if err := json.Unmarshal(req.Body, &root); err != nil || root == nil {
				// Restore unmutated root if compression did not shrink payload.
				return nil
			}
		} else if cfg.TrimHistoryEnabled {
			m.reqProcessed.Add(1)
			out, savedBytes, savedTokens, mutated := trimRootHistory(root, req.Body, cfg.MaxHistoryMessages)
			if mutated {
				req.Body = out
				m.reqCompressed.Add(1)
				m.bytesSaved.Add(savedBytes)
				m.tokensSaved.Add(savedTokens)
				req.Metadata["tokens_saved"] = savedTokens
			}
		}
		if cfg.CavemanEnabled || cfg.PonytailEnabled {
			if injectOutputDirectivesRoot(req.Protocol, root, cfg) {
				rootDirty = true
			}
		}
	}

	if runHeadroom {
		if cfg.InjectCacheControl && req.Protocol == "anthropic" {
			if injectAnthropicCacheControlRoot(root) {
				rootDirty = true
			}
		}
		key := computePrefixFingerprintFromRoot(req.Protocol, req.Model, root, cfg.PrefixDepth)
		if key != "" {
			hit, provID, connID := m.headroom.lookup(key, time.Now().UTC())
			req.Metadata["headroom_cache_key"] = key
			req.Metadata["headroom_cache_hit"] = hit
			req.Metadata["headroom_protocol"] = req.Protocol
			req.Metadata["headroom_model"] = req.Model
			if hit {
				m.cacheHits.Add(1)
				if provID != "" {
					req.Metadata["headroom_affinity_provider"] = provID
				}
				if connID != "" {
					req.Metadata["headroom_affinity_connection"] = connID
				}
			} else {
				m.cacheMisses.Add(1)
			}
		}
	}

	if rootDirty {
		if encoded, err := json.Marshal(root); err == nil {
			req.Body = encoded
		}
	}

	return nil
}

// OnResponseHeaders records provider/connection affinity and emits observability headers.
func (m *Module) OnResponseHeaders(_ context.Context, meta *contracts.ResponseMeta) error {
	if meta == nil || meta.Metadata == nil {
		return nil
	}
	if saved, ok := meta.Metadata["tokens_saved"].(int64); ok && saved > 0 && meta.Headers != nil {
		meta.Headers.Set("X-Omeg4mon-Tokens-Saved", strconv.FormatInt(saved, 10))
	}
	if key, ok := meta.Metadata["headroom_cache_key"].(string); ok && key != "" {
		proto, _ := meta.Metadata["headroom_protocol"].(string)
		model, _ := meta.Metadata["headroom_model"].(string)
		if model == "" {
			model = meta.Model
		}
		m.headroom.recordAffinity(key, proto, model, meta.ProviderID, meta.ConnectionID, time.Now().UTC())
		if meta.Headers != nil {
			if hit, _ := meta.Metadata["headroom_cache_hit"].(bool); hit {
				meta.Headers.Set("X-Omeg4mon-Headroom-Cache", "hit")
			} else {
				meta.Headers.Set("X-Omeg4mon-Headroom-Cache", "miss")
			}
		}
	}
	return nil
}

// OnStreamChunk is a no-op for tokensaver.
func (m *Module) OnStreamChunk(context.Context, *contracts.StreamEvent) error { return nil }

func (m *Module) RegisterRoutes(mounts lifecycle.RouteMounts) {
	if mounts.AdminAPI == nil {
		return
	}
	mounts.AdminAPI.Route("/token-saver", func(r chi.Router) {
		r.Get("/config", m.handleGetConfig)
		r.Patch("/config", m.handlePatchConfig)
		r.Get("/stats", m.handleGetStats)
		r.Delete("/stats", m.handleResetStats)
		r.Delete("/cache", m.handleClearCache)
		r.Post("/preview", m.handlePreview)
	})
}

func (m *Module) handleGetConfig(w http.ResponseWriter, r *http.Request) {
	cfg, err := m.repo.getConfig(r.Context())
	if err != nil {
		httpserver.WriteError(w, r, http.StatusInternalServerError, "db_error", "failed to load tokensaver config")
		return
	}
	httpserver.WriteJSON(w, http.StatusOK, cfg)
}

func (m *Module) handlePatchConfig(w http.ResponseWriter, r *http.Request) {
	var patch struct {
		Enabled            *bool   `json:"enabled"`
		Strategy           *string `json:"strategy"`
		MaxHistoryMessages *int    `json:"max_history_messages"`
		MaxToolOutputChars *int    `json:"max_tool_output_chars"`
		HeadroomEnabled    *bool   `json:"headroom_enabled"`
		PrefixDepth        *int    `json:"prefix_depth"`
		MaxCacheEntries    *int    `json:"max_cache_entries"`
		CacheTTLSeconds    *int    `json:"cache_ttl_seconds"`
		InjectCacheControl *bool   `json:"inject_cache_control"`
		CavemanEnabled     *bool   `json:"caveman_enabled"`
		CavemanLevel       *string `json:"caveman_level"`
		PonytailEnabled    *bool   `json:"ponytail_enabled"`
		PonytailLevel      *string `json:"ponytail_level"`
		TrimHistoryEnabled *bool   `json:"trim_history_enabled"`
	}
	if !decodeJSONBody(w, r, &patch) {
		return
	}

	m.cfgPatchMu.Lock()
	defer m.cfgPatchMu.Unlock()

	cfg, err := m.repo.getConfig(r.Context())
	if err != nil {
		httpserver.WriteError(w, r, http.StatusInternalServerError, "db_error", "failed to load tokensaver config")
		return
	}
	if patch.Enabled != nil {
		cfg.Enabled = *patch.Enabled
	}
	if patch.Strategy != nil {
		s := strings.TrimSpace(*patch.Strategy)
		if s != StrategyOff && s != StrategyWhitespace && s != StrategyDedupe && s != StrategyAggressive {
			httpserver.WriteError(w, r, http.StatusBadRequest, "invalid_strategy", "strategy must be off, whitespace, dedupe, or aggressive")
			return
		}
		cfg.Strategy = s
	}
	if patch.MaxHistoryMessages != nil {
		if *patch.MaxHistoryMessages < 4 || *patch.MaxHistoryMessages > 500 {
			httpserver.WriteError(w, r, http.StatusBadRequest, "invalid_max_history", "max_history_messages must be 4..500")
			return
		}
		cfg.MaxHistoryMessages = *patch.MaxHistoryMessages
	}
	if patch.MaxToolOutputChars != nil {
		if *patch.MaxToolOutputChars < 32 || *patch.MaxToolOutputChars > 200000 {
			httpserver.WriteError(w, r, http.StatusBadRequest, "invalid_max_tool_chars", "max_tool_output_chars must be 32..200000")
			return
		}
		cfg.MaxToolOutputChars = *patch.MaxToolOutputChars
	}
	if patch.HeadroomEnabled != nil {
		cfg.HeadroomEnabled = *patch.HeadroomEnabled
	}
	if patch.PrefixDepth != nil {
		if *patch.PrefixDepth < 1 || *patch.PrefixDepth > 10 {
			httpserver.WriteError(w, r, http.StatusBadRequest, "invalid_prefix_depth", "prefix_depth must be 1..10")
			return
		}
		cfg.PrefixDepth = *patch.PrefixDepth
	}
	if patch.MaxCacheEntries != nil {
		if *patch.MaxCacheEntries < 8 || *patch.MaxCacheEntries > 10000 {
			httpserver.WriteError(w, r, http.StatusBadRequest, "invalid_max_cache_entries", "max_cache_entries must be 8..10000")
			return
		}
		cfg.MaxCacheEntries = *patch.MaxCacheEntries
	}
	if patch.CacheTTLSeconds != nil {
		if *patch.CacheTTLSeconds < 60 || *patch.CacheTTLSeconds > 86400 {
			httpserver.WriteError(w, r, http.StatusBadRequest, "invalid_cache_ttl", "cache_ttl_seconds must be 60..86400")
			return
		}
		cfg.CacheTTLSeconds = *patch.CacheTTLSeconds
	}
	if patch.InjectCacheControl != nil {
		cfg.InjectCacheControl = *patch.InjectCacheControl
	}
	if patch.CavemanEnabled != nil {
		cfg.CavemanEnabled = *patch.CavemanEnabled
	}
	if patch.CavemanLevel != nil {
		lvl := strings.TrimSpace(*patch.CavemanLevel)
		if !isValidLevel(lvl) {
			httpserver.WriteError(w, r, http.StatusBadRequest, "invalid_caveman_level", "caveman_level must be lite, full, or ultra")
			return
		}
		cfg.CavemanLevel = lvl
	}
	if patch.PonytailEnabled != nil {
		cfg.PonytailEnabled = *patch.PonytailEnabled
	}
	if patch.PonytailLevel != nil {
		lvl := strings.TrimSpace(*patch.PonytailLevel)
		if !isValidLevel(lvl) {
			httpserver.WriteError(w, r, http.StatusBadRequest, "invalid_ponytail_level", "ponytail_level must be lite, full, or ultra")
			return
		}
		cfg.PonytailLevel = lvl
	}
	if patch.TrimHistoryEnabled != nil {
		cfg.TrimHistoryEnabled = *patch.TrimHistoryEnabled
	}

	saved, err := m.repo.saveConfig(r.Context(), cfg)
	if err != nil {
		httpserver.WriteError(w, r, http.StatusInternalServerError, "db_error", "failed to save tokensaver config")
		return
	}
	m.mu.Lock()
	m.cfg = saved
	m.mu.Unlock()
	m.headroom.configure(saved.MaxCacheEntries, time.Duration(saved.CacheTTLSeconds)*time.Second)
	if m.bus != nil {
		m.bus.Publish(eventConfigChanged, nil)
	}
	httpserver.WriteJSON(w, http.StatusOK, saved)
}

func (m *Module) handleGetStats(w http.ResponseWriter, _ *http.Request) {
	hits := m.cacheHits.Load()
	misses := m.cacheMisses.Load()
	hitRate := 0.0
	if total := hits + misses; total > 0 {
		hitRate = float64(hits) / float64(total)
	}
	active, maxEntries, warm := m.headroom.snapshot(12)
	httpserver.WriteJSON(w, http.StatusOK, StatsResponse{
		Compression: CompressionStats{
			RequestsProcessed:    m.reqProcessed.Load(),
			RequestsCompressed:   m.reqCompressed.Load(),
			BytesSaved:           m.bytesSaved.Load(),
			EstimatedTokensSaved: m.tokensSaved.Load(),
		},
		Headroom: HeadroomStats{
			CacheHits:     hits,
			CacheMisses:   misses,
			HitRate:       hitRate,
			ActiveEntries: active,
			MaxEntries:    maxEntries,
			WarmPrefixes:  warm,
		},
	})
}

func (m *Module) handleResetStats(w http.ResponseWriter, _ *http.Request) {
	m.reqProcessed.Store(0)
	m.reqCompressed.Store(0)
	m.bytesSaved.Store(0)
	m.tokensSaved.Store(0)
	m.cacheHits.Store(0)
	m.cacheMisses.Store(0)
	w.WriteHeader(http.StatusNoContent)
}

func (m *Module) handleClearCache(w http.ResponseWriter, _ *http.Request) {
	m.headroom.clear()
	w.WriteHeader(http.StatusNoContent)
}

func (m *Module) handlePreview(w http.ResponseWriter, r *http.Request) {
	var body struct {
		Protocol string         `json:"protocol"`
		Model    string         `json:"model"`
		Strategy string         `json:"strategy"`
		Body     map[string]any `json:"body"`
	}
	if !decodeJSONBody(w, r, &body) {
		return
	}
	if body.Body == nil {
		httpserver.WriteError(w, r, http.StatusBadRequest, "invalid_body", "body object is required")
		return
	}
	proto := strings.TrimSpace(body.Protocol)
	if proto == "" {
		proto = "openai-chat"
	}
	model := strings.TrimSpace(body.Model)
	if model == "" {
		model = "gpt-4o"
	}
	raw, err := json.Marshal(body.Body)
	if err != nil {
		httpserver.WriteError(w, r, http.StatusBadRequest, "invalid_body", "failed to encode preview body")
		return
	}
	cfg := m.snapshotConfig()
	if body.Strategy != "" {
		cfg.Strategy = body.Strategy
	}
	if cfg.Strategy == StrategyOff {
		cfg.Strategy = StrategyWhitespace
	}
	outRaw, savedBytes, savedTokens, _ := compressPayload(raw, cfg)
	cacheKey := computePrefixFingerprint(proto, model, outRaw, cfg.PrefixDepth)
	outMap := make(map[string]any)
	_ = json.Unmarshal(outRaw, &outMap)
	httpserver.WriteJSON(w, http.StatusOK, PreviewResponse{
		OriginalBytes:      len(raw),
		CompressedBytes:    len(outRaw),
		BytesSaved:         savedBytes,
		EstimatedSavedToks: savedTokens,
		CacheKey:           cacheKey,
		Body:               outMap,
	})
}

func decodeJSONBody(w http.ResponseWriter, r *http.Request, dst any) bool {
	limited := io.LimitReader(r.Body, maxBodyBytes+1)
	raw, err := io.ReadAll(limited)
	if err != nil || int64(len(raw)) > maxBodyBytes {
		httpserver.WriteError(w, r, http.StatusBadRequest, "invalid_body", "request body unreadable or exceeds 1 MiB")
		return false
	}
	if err := json.Unmarshal(raw, dst); err != nil {
		httpserver.WriteError(w, r, http.StatusBadRequest, "invalid_json", "malformed JSON payload")
		return false
	}
	return true
}
