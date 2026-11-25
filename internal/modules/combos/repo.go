package combos

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/omeg4monx/router-gateway/internal/platform/db"
)

const rfc3339NanoFixed = "2006-01-02T15:04:05.000000000Z"

var (
	// ErrComboNotFound indicates no combo exists with the requested ID.
	ErrComboNotFound = errors.New("combos: combo not found")
	// ErrComboAlreadyExists indicates a combo with the requested name already exists.
	ErrComboAlreadyExists = errors.New("combos: combo name already exists")
	// ErrAliasNotFound indicates no model alias exists with the requested alias name.
	ErrAliasNotFound = errors.New("combos: model alias not found")
	// ErrDisabledModelNotFound indicates the requested model is not currently disabled.
	ErrDisabledModelNotFound = errors.New("combos: disabled model not found")
)

// Combo represents a multi-model routing combo following OmniRoute's schema.
type Combo struct {
	ID          string         `json:"id"`
	Name        string         `json:"name"`
	Description string         `json:"description"`
	Targets     []string       `json:"targets"`
	Steps       []ComboStep    `json:"steps"`
	Strategy    string         `json:"strategy"`
	Config      map[string]any `json:"config"`
	Enabled     bool           `json:"enabled"`
	CreatedAt   string         `json:"created_at"`
	UpdatedAt   string         `json:"updated_at"`
}

// ModelAlias maps a logical model alias to a target model, combo, or another alias.
type ModelAlias struct {
	Alias       string `json:"alias"`
	TargetModel string `json:"target_model"`
	CreatedAt   string `json:"created_at"`
	UpdatedAt   string `json:"updated_at"`
}

// DisabledModel represents a model identifier suppressed from routing and listing.
type DisabledModel struct {
	Model     string `json:"model"`
	Reason    string `json:"reason"`
	CreatedAt string `json:"created_at"`
}

type createComboParams struct {
	ID          string
	Name        string
	Description string
	Targets     []string
	Steps       []ComboStep
	Strategy    string
	Config      map[string]any
	Enabled     bool
}

type updateComboValues struct {
	Name        string
	Description string
	Targets     []string
	Steps       []ComboStep
	Strategy    string
	Config      map[string]any
	Enabled     bool
}

type repository struct {
	db  *db.DB
	now func() time.Time
}

func newRepository(database *db.DB, clock func() time.Time) *repository {
	if clock == nil {
		clock = func() time.Time { return time.Now().UTC() }
	}
	return &repository{
		db:  database,
		now: clock,
	}
}

type rowScanner interface {
	Scan(dest ...any) error
}

const selectComboColumns = `id, name, description, targets_json, steps_json, strategy, config_json, enabled, created_at, updated_at`

func scanCombo(sc rowScanner) (Combo, error) {
	var (
		c           Combo
		targetsJSON string
		stepsJSON   string
		configJSON  string
		enabledInt  int
	)
	if err := sc.Scan(
		&c.ID,
		&c.Name,
		&c.Description,
		&targetsJSON,
		&stepsJSON,
		&c.Strategy,
		&configJSON,
		&enabledInt,
		&c.CreatedAt,
		&c.UpdatedAt,
	); err != nil {
		return Combo{}, err
	}
	targets, err := decodeStringSlice(targetsJSON)
	if err != nil {
		return Combo{}, fmt.Errorf("combos: decode targets for %q: %w", c.ID, err)
	}
	c.Targets = targets
	c.Steps = decodeStepsJSON(stepsJSON, targets)
	c.Config = decodeConfigJSON(configJSON)
	if c.Strategy == "ordered" || c.Strategy == "" {
		c.Strategy = defaultComboStrategy
	}
	c.Enabled = enabledInt != 0
	return c, nil
}

func (r *repository) nowUTCString() string {
	return r.now().UTC().Format(rfc3339NanoFixed)
}

func (r *repository) ListCombos(ctx context.Context) ([]Combo, error) {
	if r == nil || r.db == nil {
		return nil, errors.New("combos: database is not initialized")
	}

	rows, err := r.db.Reader().QueryContext(
		ctx,
		`SELECT `+selectComboColumns+` FROM combos ORDER BY name ASC, id ASC`,
	)
	if err != nil {
		return nil, fmt.Errorf("combos: list combos: %w", err)
	}
	defer rows.Close()

	items := make([]Combo, 0)
	for rows.Next() {
		item, err := scanCombo(rows)
		if err != nil {
			return nil, fmt.Errorf("combos: scan combo: %w", err)
		}
		items = append(items, item)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("combos: iterate combos: %w", err)
	}
	return items, nil
}

func (r *repository) GetCombo(ctx context.Context, id string) (Combo, error) {
	if r == nil || r.db == nil {
		return Combo{}, errors.New("combos: database is not initialized")
	}

	row := r.db.Reader().QueryRowContext(
		ctx,
		`SELECT `+selectComboColumns+` FROM combos WHERE id = ?`,
		id,
	)
	item, err := scanCombo(row)
	if errors.Is(err, sql.ErrNoRows) {
		return Combo{}, ErrComboNotFound
	}
	if err != nil {
		return Combo{}, fmt.Errorf("combos: get combo %q: %w", id, err)
	}
	return item, nil
}

func (r *repository) CreateCombo(ctx context.Context, p createComboParams) (Combo, error) {
	if r == nil || r.db == nil {
		return Combo{}, errors.New("combos: database is not initialized")
	}

	targetsJSON, err := encodeStringSlice(p.Targets)
	if err != nil {
		return Combo{}, fmt.Errorf("combos: encode targets: %w", err)
	}
	steps := p.Steps
	if len(steps) == 0 {
		steps = stepsFromTargets(p.Targets)
	}
	stepsJSON, err := encodeStepsJSON(steps)
	if err != nil {
		return Combo{}, fmt.Errorf("combos: encode steps: %w", err)
	}
	configJSON, err := encodeConfigJSON(p.Config)
	if err != nil {
		return Combo{}, fmt.Errorf("combos: encode config: %w", err)
	}
	enabledInt := 0
	if p.Enabled {
		enabledInt = 1
	}

	now := r.nowUTCString()
	if _, err := r.db.Writer().ExecContext(
		ctx,
		`INSERT INTO combos (id, name, description, targets_json, steps_json, strategy, config_json, enabled, created_at, updated_at)
		 VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`,
		p.ID,
		p.Name,
		p.Description,
		targetsJSON,
		stepsJSON,
		p.Strategy,
		configJSON,
		enabledInt,
		now,
		now,
	); err != nil {
		if isUniqueConstraintError(err) {
			return Combo{}, ErrComboAlreadyExists
		}
		return Combo{}, fmt.Errorf("combos: insert combo: %w", err)
	}

	return Combo{
		ID:          p.ID,
		Name:        p.Name,
		Description: p.Description,
		Targets:     cloneStringSlice(p.Targets),
		Steps:       cloneSteps(steps),
		Strategy:    p.Strategy,
		Config:      cloneConfig(p.Config),
		Enabled:     p.Enabled,
		CreatedAt:   now,
		UpdatedAt:   now,
	}, nil
}

func (r *repository) UpdateCombo(ctx context.Context, id string, fn func(cur Combo) (updateComboValues, error)) (Combo, error) {
	if r == nil || r.db == nil {
		return Combo{}, errors.New("combos: database is not initialized")
	}

	tx, err := r.db.Writer().BeginTx(ctx, nil)
	if err != nil {
		return Combo{}, fmt.Errorf("combos: begin update combo tx: %w", err)
	}
	defer func() {
		_ = tx.Rollback()
	}()

	cur, err := scanCombo(tx.QueryRowContext(
		ctx,
		`SELECT `+selectComboColumns+` FROM combos WHERE id = ?`,
		id,
	))
	if errors.Is(err, sql.ErrNoRows) {
		return Combo{}, ErrComboNotFound
	}
	if err != nil {
		return Combo{}, fmt.Errorf("combos: read combo %q for update: %w", id, err)
	}

	next, err := fn(cur)
	if err != nil {
		return Combo{}, err
	}

	targetsJSON, err := encodeStringSlice(next.Targets)
	if err != nil {
		return Combo{}, fmt.Errorf("combos: encode updated targets: %w", err)
	}
	steps := next.Steps
	if len(steps) == 0 {
		steps = stepsFromTargets(next.Targets)
	}
	stepsJSON, err := encodeStepsJSON(steps)
	if err != nil {
		return Combo{}, fmt.Errorf("combos: encode updated steps: %w", err)
	}
	configJSON, err := encodeConfigJSON(next.Config)
	if err != nil {
		return Combo{}, fmt.Errorf("combos: encode updated config: %w", err)
	}
	enabledInt := 0
	if next.Enabled {
		enabledInt = 1
	}

	now := r.nowUTCString()
	if _, err := tx.ExecContext(
		ctx,
		`UPDATE combos
		 SET name = ?,
		     description = ?,
		     targets_json = ?,
		     steps_json = ?,
		     strategy = ?,
		     config_json = ?,
		     enabled = ?,
		     updated_at = ?
		 WHERE id = ?`,
		next.Name,
		next.Description,
		targetsJSON,
		stepsJSON,
		next.Strategy,
		configJSON,
		enabledInt,
		now,
		id,
	); err != nil {
		if isUniqueConstraintError(err) {
			return Combo{}, ErrComboAlreadyExists
		}
		return Combo{}, fmt.Errorf("combos: update combo %q: %w", id, err)
	}

	if err := tx.Commit(); err != nil {
		return Combo{}, fmt.Errorf("combos: commit update combo tx: %w", err)
	}

	cur.Name = next.Name
	cur.Description = next.Description
	cur.Targets = cloneStringSlice(next.Targets)
	cur.Steps = cloneSteps(steps)
	cur.Strategy = next.Strategy
	cur.Config = cloneConfig(next.Config)
	cur.Enabled = next.Enabled
	cur.UpdatedAt = now
	return cur, nil
}

func (r *repository) DeleteCombo(ctx context.Context, id string) (string, error) {
	if r == nil || r.db == nil {
		return "", errors.New("combos: database is not initialized")
	}

	res, err := r.db.Writer().ExecContext(ctx, `DELETE FROM combos WHERE id = ?`, id)
	if err != nil {
		return "", fmt.Errorf("combos: delete combo %q: %w", id, err)
	}
	affected, err := res.RowsAffected()
	if err != nil {
		return "", fmt.Errorf("combos: delete combo rows affected: %w", err)
	}
	if affected == 0 {
		return "", ErrComboNotFound
	}
	return r.nowUTCString(), nil
}

func (r *repository) ListAliases(ctx context.Context) ([]ModelAlias, error) {
	if r == nil || r.db == nil {
		return nil, errors.New("combos: database is not initialized")
	}

	rows, err := r.db.Reader().QueryContext(
		ctx,
		`SELECT alias, target_model, created_at, updated_at FROM model_aliases ORDER BY alias ASC`,
	)
	if err != nil {
		return nil, fmt.Errorf("combos: list model_aliases: %w", err)
	}
	defer rows.Close()

	items := make([]ModelAlias, 0)
	for rows.Next() {
		var item ModelAlias
		if err := rows.Scan(&item.Alias, &item.TargetModel, &item.CreatedAt, &item.UpdatedAt); err != nil {
			return nil, fmt.Errorf("combos: scan model_alias: %w", err)
		}
		items = append(items, item)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("combos: iterate model_aliases: %w", err)
	}
	return items, nil
}

func (r *repository) UpsertAlias(ctx context.Context, alias, targetModel string) (ModelAlias, error) {
	if r == nil || r.db == nil {
		return ModelAlias{}, errors.New("combos: database is not initialized")
	}

	tx, err := r.db.Writer().BeginTx(ctx, nil)
	if err != nil {
		return ModelAlias{}, fmt.Errorf("combos: begin upsert alias tx: %w", err)
	}
	defer func() {
		_ = tx.Rollback()
	}()

	now := r.nowUTCString()
	if _, err := tx.ExecContext(
		ctx,
		`INSERT INTO model_aliases (alias, target_model, created_at, updated_at)
		 VALUES (?, ?, ?, ?)
		 ON CONFLICT(alias) DO UPDATE SET
		     target_model = excluded.target_model,
		     updated_at = excluded.updated_at`,
		alias,
		targetModel,
		now,
		now,
	); err != nil {
		return ModelAlias{}, fmt.Errorf("combos: upsert model_alias %q: %w", alias, err)
	}

	var out ModelAlias
	if err := tx.QueryRowContext(
		ctx,
		`SELECT alias, target_model, created_at, updated_at FROM model_aliases WHERE alias = ?`,
		alias,
	).Scan(&out.Alias, &out.TargetModel, &out.CreatedAt, &out.UpdatedAt); err != nil {
		return ModelAlias{}, fmt.Errorf("combos: reload model_alias %q: %w", alias, err)
	}

	if err := tx.Commit(); err != nil {
		return ModelAlias{}, fmt.Errorf("combos: commit upsert alias tx: %w", err)
	}
	return out, nil
}

func (r *repository) DeleteAlias(ctx context.Context, alias string) (string, error) {
	if r == nil || r.db == nil {
		return "", errors.New("combos: database is not initialized")
	}

	res, err := r.db.Writer().ExecContext(ctx, `DELETE FROM model_aliases WHERE alias = ?`, alias)
	if err != nil {
		return "", fmt.Errorf("combos: delete model_alias %q: %w", alias, err)
	}
	affected, err := res.RowsAffected()
	if err != nil {
		return "", fmt.Errorf("combos: delete model_alias rows affected: %w", err)
	}
	if affected == 0 {
		return "", ErrAliasNotFound
	}
	return r.nowUTCString(), nil
}

func (r *repository) ListDisabledModels(ctx context.Context) ([]DisabledModel, error) {
	if r == nil || r.db == nil {
		return nil, errors.New("combos: database is not initialized")
	}

	rows, err := r.db.Reader().QueryContext(
		ctx,
		`SELECT model, reason, created_at FROM disabled_models ORDER BY model ASC`,
	)
	if err != nil {
		return nil, fmt.Errorf("combos: list disabled_models: %w", err)
	}
	defer rows.Close()

	items := make([]DisabledModel, 0)
	for rows.Next() {
		var item DisabledModel
		if err := rows.Scan(&item.Model, &item.Reason, &item.CreatedAt); err != nil {
			return nil, fmt.Errorf("combos: scan disabled_model: %w", err)
		}
		items = append(items, item)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("combos: iterate disabled_models: %w", err)
	}
	return items, nil
}

func (r *repository) UpsertDisabledModel(ctx context.Context, model, reason string) (DisabledModel, error) {
	if r == nil || r.db == nil {
		return DisabledModel{}, errors.New("combos: database is not initialized")
	}

	tx, err := r.db.Writer().BeginTx(ctx, nil)
	if err != nil {
		return DisabledModel{}, fmt.Errorf("combos: begin upsert disabled_model tx: %w", err)
	}
	defer func() {
		_ = tx.Rollback()
	}()

	now := r.nowUTCString()
	if _, err := tx.ExecContext(
		ctx,
		`INSERT INTO disabled_models (model, reason, created_at)
		 VALUES (?, ?, ?)
		 ON CONFLICT(model) DO UPDATE SET
		     reason = excluded.reason`,
		model,
		reason,
		now,
	); err != nil {
		return DisabledModel{}, fmt.Errorf("combos: upsert disabled_model %q: %w", model, err)
	}

	var out DisabledModel
	if err := tx.QueryRowContext(
		ctx,
		`SELECT model, reason, created_at FROM disabled_models WHERE model = ?`,
		model,
	).Scan(&out.Model, &out.Reason, &out.CreatedAt); err != nil {
		return DisabledModel{}, fmt.Errorf("combos: reload disabled_model %q: %w", model, err)
	}

	if err := tx.Commit(); err != nil {
		return DisabledModel{}, fmt.Errorf("combos: commit upsert disabled_model tx: %w", err)
	}
	return out, nil
}

func (r *repository) DeleteDisabledModel(ctx context.Context, model string) (string, error) {
	if r == nil || r.db == nil {
		return "", errors.New("combos: database is not initialized")
	}

	res, err := r.db.Writer().ExecContext(ctx, `DELETE FROM disabled_models WHERE model = ?`, model)
	if err != nil {
		return "", fmt.Errorf("combos: delete disabled_model %q: %w", model, err)
	}
	affected, err := res.RowsAffected()
	if err != nil {
		return "", fmt.Errorf("combos: delete disabled_model rows affected: %w", err)
	}
	if affected == 0 {
		return "", ErrDisabledModelNotFound
	}
	return r.nowUTCString(), nil
}

func (r *repository) ListProviderStrategies(ctx context.Context) ([]ProviderStrategyState, error) {
	if r == nil || r.db == nil {
		return nil, errors.New("combos: database is not initialized")
	}

	rows, err := r.db.Reader().QueryContext(
		ctx,
		`SELECT provider_id, strategy, updated_at FROM provider_strategies ORDER BY provider_id ASC`,
	)
	if err != nil {
		return nil, fmt.Errorf("combos: list provider_strategies: %w", err)
	}
	defer func() {
		_ = rows.Close()
	}()

	items := make([]ProviderStrategyState, 0)
	for rows.Next() {
		var item ProviderStrategyState
		if err := rows.Scan(&item.ProviderID, &item.Strategy, &item.UpdatedAt); err != nil {
			return nil, fmt.Errorf("combos: scan provider_strategy: %w", err)
		}
		item.DefaultStrategy = DefaultProviderStrategy
		item.IsDefault = false
		items = append(items, item)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("combos: iterate provider_strategies: %w", err)
	}
	return items, nil
}

func (r *repository) GetProviderStrategy(ctx context.Context, providerID string) (ProviderStrategyState, error) {
	if r == nil || r.db == nil {
		return ProviderStrategyState{}, errors.New("combos: database is not initialized")
	}

	var (
		strategy  string
		updatedAt string
	)
	err := r.db.Reader().QueryRowContext(
		ctx,
		`SELECT strategy, updated_at FROM provider_strategies WHERE provider_id = ?`,
		providerID,
	).Scan(&strategy, &updatedAt)
	if errors.Is(err, sql.ErrNoRows) {
		return ProviderStrategyState{
			ProviderID:      providerID,
			Strategy:        DefaultProviderStrategy,
			DefaultStrategy: DefaultProviderStrategy,
			IsDefault:       true,
		}, nil
	}
	if err != nil {
		return ProviderStrategyState{}, fmt.Errorf("combos: get provider_strategy %q: %w", providerID, err)
	}
	return ProviderStrategyState{
		ProviderID:      providerID,
		Strategy:        strategy,
		DefaultStrategy: DefaultProviderStrategy,
		IsDefault:       false,
		UpdatedAt:       updatedAt,
	}, nil
}

func (r *repository) UpsertProviderStrategy(ctx context.Context, providerID, strategy string) (ProviderStrategyState, error) {
	if r == nil || r.db == nil {
		return ProviderStrategyState{}, errors.New("combos: database is not initialized")
	}

	now := r.nowUTCString()
	if _, err := r.db.Writer().ExecContext(
		ctx,
		`INSERT INTO provider_strategies (provider_id, strategy, updated_at)
		 VALUES (?, ?, ?)
		 ON CONFLICT(provider_id) DO UPDATE SET
		     strategy = excluded.strategy,
		     updated_at = excluded.updated_at`,
		providerID,
		strategy,
		now,
	); err != nil {
		return ProviderStrategyState{}, fmt.Errorf("combos: upsert provider_strategy %q: %w", providerID, err)
	}

	return ProviderStrategyState{
		ProviderID:      providerID,
		Strategy:        strategy,
		DefaultStrategy: DefaultProviderStrategy,
		IsDefault:       false,
		UpdatedAt:       now,
	}, nil
}

func cloneStringSlice(items []string) []string {
	if len(items) == 0 {
		return []string{}
	}
	out := make([]string, len(items))
	copy(out, items)
	return out
}

func encodeStringSlice(items []string) (string, error) {
	data, err := json.Marshal(cloneStringSlice(items))
	if err != nil {
		return "", err
	}
	return string(data), nil
}

func decodeStringSlice(raw string) ([]string, error) {
	if strings.TrimSpace(raw) == "" {
		return []string{}, nil
	}
	var items []string
	if err := json.Unmarshal([]byte(raw), &items); err != nil {
		return nil, err
	}
	if items == nil {
		return []string{}, nil
	}
	return items, nil
}

func isUniqueConstraintError(err error) bool {
	if err == nil {
		return false
	}
	msg := strings.ToLower(err.Error())
	return strings.Contains(msg, "unique constraint failed") || strings.Contains(msg, "constraint_primarykey")
}
