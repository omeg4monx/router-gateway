package plugins

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/omeg4monx/router-gateway/internal/platform/db"
)

const rfc3339NanoFixed = "2006-01-02T15:04:05.000000000Z"

// ErrPluginNotFound indicates the requested built-in plugin key does not exist.
var ErrPluginNotFound = errors.New("plugins: plugin not found")

// PluginSpec describes metadata for a built-in request-enhancement plugin.
type PluginSpec struct {
	Key         string         `json:"key"`
	Name        string         `json:"name"`
	Description string         `json:"description"`
	Category    string         `json:"category"`
	DefaultCfg  map[string]any `json:"-"`
}

// PluginState combines built-in metadata with persisted runtime configuration.
type PluginState struct {
	Key         string         `json:"key"`
	Name        string         `json:"name"`
	Description string         `json:"description"`
	Category    string         `json:"category"`
	Enabled     bool           `json:"enabled"`
	Models      []string       `json:"models"`
	Config      map[string]any `json:"config"`
	UpdatedAt   string         `json:"updated_at"`
}

var builtinPlugins = []PluginSpec{
	{
		Key:         "imageVision",
		Name:        "Image Vision Helper",
		Description: "Preserves multimodal image blocks and injects visual-inspection instructions only when image blocks are present.",
		Category:    "multimodal",
		DefaultCfg:  map[string]any{},
	},
	{
		Key:         "thinkDeeper",
		Name:        "Think Deeper",
		Description: "Enables provider-native deep reasoning (Claude thinking budget or OpenAI reasoning_effort) for thorough multi-step analysis.",
		Category:    "reasoning",
		DefaultCfg:  map[string]any{"effort": "high", "budget_tokens": 10240},
	},
	{
		Key:         "speedMode",
		Name:        "Speed Mode",
		Description: "Disables provider-native reasoning/thinking overhead for low-latency direct answers with optional max_tokens cap.",
		Category:    "performance",
		DefaultCfg:  map[string]any{"max_tokens": 0},
	},
	{
		Key:         "jsonGuard",
		Name:        "JSON Guard",
		Description: "Requests strict JSON when prompted and repairs malformed/fenced/truncated JSON and tool-call arguments on responses.",
		Category:    "reliability",
		DefaultCfg:  map[string]any{"force": false},
	},
	{
		Key:         "contextSqueezer",
		Name:        "Context Squeezer",
		Description: "Fits long conversations into the target context window by shortening oversized messages and summarizing dropped oldest turns.",
		Category:    "optimization",
		DefaultCfg:  map[string]any{"context_window": 128000, "fill_ratio": 0.75, "min_turns": 4},
	},
	{
		Key:         "openaiToolBridge",
		Name:        "OpenAI Tool Bridge",
		Description: "Normalizes tool_choice/parallel_tool_calls on requests and extracts schema-validated tool_calls from prose responses.",
		Category:    "tools",
		DefaultCfg:  map[string]any{"parallel_tool_calls": true},
	},
	{
		Key:         "antiSlop",
		Name:        "Anti-Slop Filter",
		Description: "Injects tiered anti-slop craftsmanship rules (lite, standard, full, ultra) to eliminate generic AI UI, copy, and code.",
		Category:    "quality",
		DefaultCfg:  map[string]any{"level": "full"},
	},
}

func lookupPluginSpec(key string) (PluginSpec, bool) {
	for _, spec := range builtinPlugins {
		if spec.Key == key {
			return spec, true
		}
	}
	return PluginSpec{}, false
}

type repository struct {
	db  *db.DB
	now func() time.Time
}

func newRepository(database *db.DB, clock func() time.Time) *repository {
	if clock == nil {
		clock = func() time.Time { return time.Now().UTC() }
	}
	return &repository{db: database, now: clock}
}

func (r *repository) nowUTC() string {
	return r.now().UTC().Format(rfc3339NanoFixed)
}

func (r *repository) ListPlugins(ctx context.Context) ([]PluginState, error) {
	rows, err := r.db.Reader().QueryContext(
		ctx,
		`SELECT key, enabled, models_json, config_json, updated_at FROM plugin_configs`,
	)
	if err != nil {
		return nil, fmt.Errorf("plugins: query plugin_configs: %w", err)
	}
	defer rows.Close()

	type stored struct {
		enabled   bool
		models    []string
		cfg       map[string]any
		updatedAt string
	}
	persisted := make(map[string]stored)
	for rows.Next() {
		var (
			key        string
			enabledInt int
			modelsJSON string
			cfgJSON    string
			updatedAt  string
		)
		if err := rows.Scan(&key, &enabledInt, &modelsJSON, &cfgJSON, &updatedAt); err != nil {
			return nil, fmt.Errorf("plugins: scan plugin_configs: %w", err)
		}
		var models []string
		_ = json.Unmarshal([]byte(modelsJSON), &models)
		if models == nil {
			models = []string{}
		}
		var cfg map[string]any
		_ = json.Unmarshal([]byte(cfgJSON), &cfg)
		if cfg == nil {
			cfg = map[string]any{}
		}
		persisted[key] = stored{
			enabled:   enabledInt != 0,
			models:    models,
			cfg:       cfg,
			updatedAt: updatedAt,
		}
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("plugins: iterate plugin_configs: %w", err)
	}

	out := make([]PluginState, 0, len(builtinPlugins))
	for _, spec := range builtinPlugins {
		st := PluginState{
			Key:         spec.Key,
			Name:        spec.Name,
			Description: spec.Description,
			Category:    spec.Category,
			Enabled:     false,
			Models:      []string{},
			Config:      cloneMap(spec.DefaultCfg),
		}
		if row, ok := persisted[spec.Key]; ok {
			st.Enabled = row.enabled
			st.Models = row.models
			st.Config = mergeDefaults(spec.DefaultCfg, row.cfg)
			st.UpdatedAt = row.updatedAt
		}
		out = append(out, st)
	}
	return out, nil
}

func (r *repository) UpsertPlugin(ctx context.Context, key string, enabled bool, models []string, cfg map[string]any) (PluginState, error) {
	spec, ok := lookupPluginSpec(key)
	if !ok {
		return PluginState{}, ErrPluginNotFound
	}
	cleanModels := sanitizeStringSlice(models)
	mergedCfg := sanitizePluginConfig(key, mergeDefaults(spec.DefaultCfg, cfg))

	modelsBytes, err := json.Marshal(cleanModels)
	if err != nil {
		return PluginState{}, fmt.Errorf("plugins: marshal models: %w", err)
	}
	cfgBytes, err := json.Marshal(mergedCfg)
	if err != nil {
		return PluginState{}, fmt.Errorf("plugins: marshal config: %w", err)
	}
	now := r.nowUTC()
	_, err = r.db.Writer().ExecContext(
		ctx,
		`INSERT INTO plugin_configs (key, enabled, models_json, config_json, updated_at)
		 VALUES (?, ?, ?, ?, ?)
		 ON CONFLICT(key) DO UPDATE SET
		   enabled = excluded.enabled,
		   models_json = excluded.models_json,
		   config_json = excluded.config_json,
		   updated_at = excluded.updated_at`,
		key, boolToInt(enabled), string(modelsBytes), string(cfgBytes), now,
	)
	if err != nil {
		return PluginState{}, fmt.Errorf("plugins: upsert plugin %q: %w", key, err)
	}
	return PluginState{
		Key:         spec.Key,
		Name:        spec.Name,
		Description: spec.Description,
		Category:    spec.Category,
		Enabled:     enabled,
		Models:      cleanModels,
		Config:      mergedCfg,
		UpdatedAt:   now,
	}, nil
}

func boolToInt(b bool) int {
	if b {
		return 1
	}
	return 0
}

func sanitizeStringSlice(items []string) []string {
	out := make([]string, 0, len(items))
	seen := make(map[string]struct{}, len(items))
	for _, raw := range items {
		s := strings.TrimSpace(raw)
		if s == "" {
			continue
		}
		if _, ok := seen[s]; ok {
			continue
		}
		seen[s] = struct{}{}
		out = append(out, s)
	}
	return out
}

func cloneMap(in map[string]any) map[string]any {
	out := make(map[string]any, len(in))
	for k, v := range in {
		out[k] = v
	}
	return out
}

func mergeDefaults(defaults, override map[string]any) map[string]any {
	out := cloneMap(defaults)
	for k, v := range override {
		out[k] = v
	}
	return out
}

func sanitizePluginConfig(key string, cfg map[string]any) map[string]any {
	out := cloneMap(cfg)
	switch key {
	case "antiSlop":
		lvl, _ := out["level"].(string)
		out["level"] = normalizeAntiSlopLevel(lvl)
	case "thinkDeeper":
		effort, _ := out["effort"].(string)
		switch strings.ToLower(strings.TrimSpace(effort)) {
		case "low", "medium", "high":
			out["effort"] = strings.ToLower(strings.TrimSpace(effort))
		default:
			out["effort"] = "high"
		}
		budget := intFromAny(out["budget_tokens"])
		if budget < 1024 {
			budget = 10240
		}
		out["budget_tokens"] = budget
	case "speedMode":
		mt := intFromAny(out["max_tokens"])
		if mt < 0 {
			mt = 0
		}
		out["max_tokens"] = mt
	case "jsonGuard":
		force, _ := out["force"].(bool)
		out["force"] = force
	case "contextSqueezer":
		ctxWin := intFromAny(out["context_window"])
		if ctxWin < 1024 {
			ctxWin = defaultContextWindow
		}
		out["context_window"] = ctxWin
		ratio := floatFromAny(out["fill_ratio"])
		if ratio <= 0.1 || ratio > 0.95 {
			ratio = defaultFillRatio
		}
		out["fill_ratio"] = ratio
		minTurns := intFromAny(out["min_turns"])
		if minTurns < 1 {
			minTurns = defaultMinTurns
		}
		out["min_turns"] = minTurns
	case "openaiToolBridge":
		if ptc, ok := out["parallel_tool_calls"].(bool); ok {
			out["parallel_tool_calls"] = ptc
		} else {
			out["parallel_tool_calls"] = true
		}
	}
	return out
}
