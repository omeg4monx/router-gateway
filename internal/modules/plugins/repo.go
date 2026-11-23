package plugins

import (
	"context"
	"crypto/rand"
	"database/sql"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/omeg4monx/router-gateway/internal/platform/db"
)

const rfc3339NanoFixed = "2006-01-02T15:04:05.000000000Z"

var (
	// ErrPluginNotFound indicates the requested built-in plugin key does not exist.
	ErrPluginNotFound = errors.New("plugins: plugin not found")
	// ErrSkillNotFound indicates no skill exists with the requested ID.
	ErrSkillNotFound = errors.New("plugins: skill not found")
	// ErrSkillSlugExists indicates a skill with the requested slug already exists.
	ErrSkillSlugExists = errors.New("plugins: skill slug already exists")
	// ErrBuiltinSkillProtected indicates a built-in skill cannot be deleted.
	ErrBuiltinSkillProtected = errors.New("plugins: built-in skill cannot be deleted")
)

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

// Skill represents a reusable agent/IDE skill markdown definition.
type Skill struct {
	ID          string `json:"id"`
	Slug        string `json:"slug"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Category    string `json:"category"`
	Content     string `json:"content"`
	Enabled     bool   `json:"enabled"`
	Builtin     bool   `json:"builtin"`
	CreatedAt   string `json:"created_at"`
	UpdatedAt   string `json:"updated_at"`
}

var builtinPlugins = []PluginSpec{
	{
		Key:         "imageVision",
		Name:        "Image Vision Helper",
		Description: "Normalizes image content blocks and injects concise visual analysis instructions when multimodal inputs are present.",
		Category:    "multimodal",
		DefaultCfg:  map[string]any{},
	},
	{
		Key:         "thinkDeeper",
		Name:        "Think Deeper",
		Description: "Injects structured step-by-step reasoning instructions and raises reasoning budget hints for complex prompts.",
		Category:    "reasoning",
		DefaultCfg:  map[string]any{"effort": "high"},
	},
	{
		Key:         "speedMode",
		Name:        "Speed Mode",
		Description: "Biases requests toward low-latency concise responses by capping verbosity and trimming redundant hints.",
		Category:    "performance",
		DefaultCfg:  map[string]any{"max_tokens": 1024},
	},
	{
		Key:         "jsonGuard",
		Name:        "JSON Guard",
		Description: "Enforces strict JSON output formatting via response_format and system-level schema discipline.",
		Category:    "reliability",
		DefaultCfg:  map[string]any{},
	},
	{
		Key:         "contextSqueezer",
		Name:        "Context Squeezer",
		Description: "Compacts redundant whitespace and blank lines in long user/system messages before upstream dispatch.",
		Category:    "optimization",
		DefaultCfg:  map[string]any{"min_chars": 256},
	},
	{
		Key:         "openaiToolBridge",
		Name:        "OpenAI Tool Bridge",
		Description: "Ensures tool_choice and parallel_tool_calls defaults are explicitly normalized for OpenAI-compatible upstreams.",
		Category:    "tools",
		DefaultCfg:  map[string]any{"parallel_tool_calls": true},
	},
	{
		Key:         "antiSlop",
		Name:        "Anti-Slop Filter",
		Description: "Prepends directness directives to eliminate filler preambles, hedging, and repetitive boilerplate.",
		Category:    "quality",
		DefaultCfg:  map[string]any{"level": "full"},
	},
}

var builtinSkills = []Skill{
	{
		ID:          "skill_builtin_gateway",
		Slug:        "omeg4mon-router-gateway",
		Name:        "OMEG4MON Router Gateway",
		Description: "Configure AI coding assistants and SDKs to route requests through the local OMEG4MON Router gateway.",
		Category:    "gateway",
		Content:     "# OMEG4MON Router Gateway\n\nPoint OpenAI, Anthropic, or Gemini clients at `http://127.0.0.1:20128/v1` with a gateway API key from `/dashboard/keys/`.\n\n## Verification\n- Run `curl -H \"Authorization: Bearer <key>\" http://127.0.0.1:20128/v1/models` to verify model visibility.\n- Prefer combo names for automatic multi-provider failover.",
		Enabled:     true,
		Builtin:     true,
	},
	{
		ID:          "skill_builtin_combos",
		Slug:        "omeg4mon-router-combos",
		Name:        "Fallback Combos & Aliases",
		Description: "Design resilient ordered fallback chains and logical aliases for zero-downtime model routing.",
		Category:    "routing",
		Content:     "# Fallback Combos & Aliases\n\n1. Create provider connections under `/dashboard/providers/`.\n2. Group primary and backup models into an ordered combo in `/dashboard/models/`.\n3. Scope API keys to the combo name so clients transparently fail over on 429/5xx.",
		Enabled:     true,
		Builtin:     true,
	},
	{
		ID:          "skill_builtin_cli_setup",
		Slug:        "omeg4mon-router-cli-setup",
		Name:        "CLI Coding Tools Setup",
		Description: "Generate and safely apply local configurations for Claude Code, Codex CLI, OpenCode, and Cline.",
		Category:    "clitools",
		Content:     "# CLI Coding Tools Setup\n\nOpen `/dashboard/cli-tools/` to generate copyable configs for supported coding tools. On loopback binds (`127.0.0.1`), one-click local apply writes atomic backups with instant rollback.",
		Enabled:     true,
		Builtin:     true,
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

func (r *repository) SeedBuiltins(ctx context.Context) error {
	now := r.nowUTC()
	for _, s := range builtinSkills {
		_, err := r.db.Writer().ExecContext(
			ctx,
			`INSERT OR IGNORE INTO plugin_skills
			 (id, slug, name, description, category, content, enabled, builtin, created_at, updated_at)
			 VALUES (?, ?, ?, ?, ?, ?, ?, 1, ?, ?)`,
			s.ID, s.Slug, s.Name, s.Description, s.Category, s.Content, boolToInt(s.Enabled), now, now,
		)
		if err != nil {
			return fmt.Errorf("plugins: seed builtin skill %q: %w", s.Slug, err)
		}
	}
	return nil
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
	mergedCfg := mergeDefaults(spec.DefaultCfg, cfg)
	if key == "antiSlop" {
		lvl, _ := mergedCfg["level"].(string)
		switch strings.ToLower(strings.TrimSpace(lvl)) {
		case "light", "standard", "full":
			mergedCfg["level"] = strings.ToLower(strings.TrimSpace(lvl))
		default:
			mergedCfg["level"] = "full"
		}
	}

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

func (r *repository) ListSkills(ctx context.Context) ([]Skill, error) {
	rows, err := r.db.Reader().QueryContext(
		ctx,
		`SELECT id, slug, name, description, category, content, enabled, builtin, created_at, updated_at
		 FROM plugin_skills
		 ORDER BY builtin DESC, name ASC`,
	)
	if err != nil {
		return nil, fmt.Errorf("plugins: query skills: %w", err)
	}
	defer rows.Close()

	skills := make([]Skill, 0)
	for rows.Next() {
		s, err := scanSkill(rows)
		if err != nil {
			return nil, err
		}
		skills = append(skills, s)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("plugins: iterate skills: %w", err)
	}
	return skills, nil
}

func (r *repository) GetSkill(ctx context.Context, id string) (Skill, error) {
	row := r.db.Reader().QueryRowContext(
		ctx,
		`SELECT id, slug, name, description, category, content, enabled, builtin, created_at, updated_at
		 FROM plugin_skills WHERE id = ?`,
		id,
	)
	s, err := scanSkill(row)
	if errors.Is(err, sql.ErrNoRows) {
		return Skill{}, ErrSkillNotFound
	}
	return s, err
}

func (r *repository) CreateSkill(ctx context.Context, slug, name, description, category, content string, enabled bool) (Skill, error) {
	raw := make([]byte, 12)
	if _, err := rand.Read(raw); err != nil {
		return Skill{}, fmt.Errorf("plugins: generate skill id: %w", err)
	}
	id := "skill_" + hex.EncodeToString(raw)
	now := r.nowUTC()
	_, err := r.db.Writer().ExecContext(
		ctx,
		`INSERT INTO plugin_skills
		 (id, slug, name, description, category, content, enabled, builtin, created_at, updated_at)
		 VALUES (?, ?, ?, ?, ?, ?, ?, 0, ?, ?)`,
		id, slug, name, description, category, content, boolToInt(enabled), now, now,
	)
	if err != nil {
		if strings.Contains(strings.ToLower(err.Error()), "unique") {
			return Skill{}, ErrSkillSlugExists
		}
		return Skill{}, fmt.Errorf("plugins: insert skill: %w", err)
	}
	return Skill{
		ID:          id,
		Slug:        slug,
		Name:        name,
		Description: description,
		Category:    category,
		Content:     content,
		Enabled:     enabled,
		Builtin:     false,
		CreatedAt:   now,
		UpdatedAt:   now,
	}, nil
}

type updateSkillPatch struct {
	Name        *string
	Description *string
	Category    *string
	Content     *string
	Enabled     *bool
}

func (r *repository) UpdateSkill(ctx context.Context, id string, patch updateSkillPatch) (Skill, error) {
	cur, err := r.GetSkill(ctx, id)
	if err != nil {
		return Skill{}, err
	}
	if patch.Name != nil {
		cur.Name = strings.TrimSpace(*patch.Name)
	}
	if patch.Description != nil {
		cur.Description = strings.TrimSpace(*patch.Description)
	}
	if patch.Category != nil {
		cur.Category = strings.TrimSpace(*patch.Category)
	}
	if patch.Content != nil {
		cur.Content = *patch.Content
	}
	if patch.Enabled != nil {
		cur.Enabled = *patch.Enabled
	}
	cur.UpdatedAt = r.nowUTC()

	res, err := r.db.Writer().ExecContext(
		ctx,
		`UPDATE plugin_skills
		 SET name = ?, description = ?, category = ?, content = ?, enabled = ?, updated_at = ?
		 WHERE id = ?`,
		cur.Name, cur.Description, cur.Category, cur.Content, boolToInt(cur.Enabled), cur.UpdatedAt, id,
	)
	if err != nil {
		return Skill{}, fmt.Errorf("plugins: update skill %q: %w", id, err)
	}
	n, _ := res.RowsAffected()
	if n == 0 {
		return Skill{}, ErrSkillNotFound
	}
	return cur, nil
}

func (r *repository) DeleteSkill(ctx context.Context, id string) error {
	cur, err := r.GetSkill(ctx, id)
	if err != nil {
		return err
	}
	if cur.Builtin {
		return ErrBuiltinSkillProtected
	}
	res, err := r.db.Writer().ExecContext(ctx, `DELETE FROM plugin_skills WHERE id = ? AND builtin = 0`, id)
	if err != nil {
		return fmt.Errorf("plugins: delete skill %q: %w", id, err)
	}
	n, _ := res.RowsAffected()
	if n == 0 {
		return ErrSkillNotFound
	}
	return nil
}

type scanner interface {
	Scan(dest ...any) error
}

func scanSkill(sc scanner) (Skill, error) {
	var (
		s          Skill
		enabledInt int
		builtinInt int
	)
	if err := sc.Scan(
		&s.ID,
		&s.Slug,
		&s.Name,
		&s.Description,
		&s.Category,
		&s.Content,
		&enabledInt,
		&builtinInt,
		&s.CreatedAt,
		&s.UpdatedAt,
	); err != nil {
		return Skill{}, err
	}
	s.Enabled = enabledInt != 0
	s.Builtin = builtinInt != 0
	return s, nil
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
