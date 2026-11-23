package usage

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/omeg4monx/router-gateway/internal/platform/db"
)

const rfc3339NanoFixed = "2006-01-02T15:04:05.000000000Z"

// ErrPricingNotFound indicates no custom pricing override exists for the requested model.
var ErrPricingNotFound = errors.New("usage: model pricing override not found")

// UsageLog represents a single persisted request telemetry record.
type UsageLog struct {
	ID           int64   `json:"id"`
	RequestID    string  `json:"request_id"`
	KeyID        string  `json:"key_id"`
	Model        string  `json:"model"`
	ProviderID   string  `json:"provider_id"`
	ConnectionID string  `json:"connection_id"`
	Protocol     string  `json:"protocol"`
	InputTokens  int64   `json:"input_tokens"`
	OutputTokens int64   `json:"output_tokens"`
	CachedTokens int64   `json:"cached_tokens"`
	CostUSD      float64 `json:"cost_usd"`
	LatencyMS    int64   `json:"latency_ms"`
	TTFTMS       int64   `json:"ttft_ms"`
	StatusCode   int     `json:"status_code"`
	ErrorSummary string  `json:"error_summary"`
	Stream       bool    `json:"stream"`
	CreatedAt    string  `json:"created_at"`
}

// LogsFilter constrains queries against persisted usage_logs.
type LogsFilter struct {
	Model      string
	ProviderID string
	KeyID      string
	Status     string
	Limit      int
	Offset     int
}

// LogsResponse is the paginated response payload for GET /api/usage/logs.
type LogsResponse struct {
	Logs   []UsageLog `json:"logs"`
	Total  int64      `json:"total"`
	Limit  int        `json:"limit"`
	Offset int        `json:"offset"`
}

// SummaryTotals aggregates overall request volume, tokens, cost, latency, and errors.
type SummaryTotals struct {
	Requests     int64   `json:"requests"`
	InputTokens  int64   `json:"input_tokens"`
	OutputTokens int64   `json:"output_tokens"`
	CachedTokens int64   `json:"cached_tokens"`
	CostUSD      float64 `json:"cost_usd"`
	AvgLatencyMS float64 `json:"avg_latency_ms"`
	AvgTTFTMS    float64 `json:"avg_ttft_ms"`
	ErrorCount   int64   `json:"error_count"`
	ErrorRate    float64 `json:"error_rate"`
}

// ModelBreakdown summarizes usage metrics for a single model.
type ModelBreakdown struct {
	Model        string  `json:"model"`
	Requests     int64   `json:"requests"`
	InputTokens  int64   `json:"input_tokens"`
	OutputTokens int64   `json:"output_tokens"`
	CachedTokens int64   `json:"cached_tokens"`
	CostUSD      float64 `json:"cost_usd"`
	AvgLatencyMS float64 `json:"avg_latency_ms"`
	ErrorCount   int64   `json:"error_count"`
}

// ProviderBreakdown summarizes usage metrics for a single provider.
type ProviderBreakdown struct {
	ProviderID   string  `json:"provider_id"`
	Requests     int64   `json:"requests"`
	InputTokens  int64   `json:"input_tokens"`
	OutputTokens int64   `json:"output_tokens"`
	CachedTokens int64   `json:"cached_tokens"`
	CostUSD      float64 `json:"cost_usd"`
	AvgLatencyMS float64 `json:"avg_latency_ms"`
	ErrorCount   int64   `json:"error_count"`
}

// SeriesBucket aggregates usage metrics for an hourly or daily time bucket.
type SeriesBucket struct {
	Bucket       string  `json:"bucket"`
	Requests     int64   `json:"requests"`
	InputTokens  int64   `json:"input_tokens"`
	OutputTokens int64   `json:"output_tokens"`
	CachedTokens int64   `json:"cached_tokens"`
	CostUSD      float64 `json:"cost_usd"`
	ErrorCount   int64   `json:"error_count"`
}

// SummaryResponse is the payload returned by GET /api/usage/summary.
type SummaryResponse struct {
	Range      string              `json:"range"`
	Totals     SummaryTotals       `json:"totals"`
	ByModel    []ModelBreakdown    `json:"by_model"`
	ByProvider []ProviderBreakdown `json:"by_provider"`
	Series     []SeriesBucket      `json:"series"`
}

type repository struct {
	db  *db.DB
	now func() time.Time
}

func newRepository(database *db.DB, now func() time.Time) *repository {
	return &repository{
		db:  database,
		now: now,
	}
}

const errorPredicateSQL = `(status_code >= 400 OR status_code = 0 OR error_summary != '')`
const successPredicateSQL = `(status_code >= 200 AND status_code < 400 AND error_summary = '')`

func (r *repository) InsertBatch(ctx context.Context, entries []UsageLog) ([]UsageLog, error) {
	if r == nil || r.db == nil {
		return nil, errors.New("usage: database is not initialized")
	}
	if len(entries) == 0 {
		return nil, nil
	}

	tx, err := r.db.Writer().BeginTx(ctx, nil)
	if err != nil {
		return nil, fmt.Errorf("usage: begin batch insert transaction: %w", err)
	}
	defer func() {
		_ = tx.Rollback()
	}()

	stmt, err := tx.PrepareContext(ctx, `INSERT INTO usage_logs (
		request_id, key_id, model, provider_id, connection_id, protocol,
		input_tokens, output_tokens, cached_tokens, cost_usd,
		latency_ms, ttft_ms, status_code, error_summary, stream, created_at
	) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`)
	if err != nil {
		return nil, fmt.Errorf("usage: prepare batch insert: %w", err)
	}
	defer func() {
		_ = stmt.Close()
	}()

	inserted := make([]UsageLog, len(entries))
	for i, entry := range entries {
		streamInt := 0
		if entry.Stream {
			streamInt = 1
		}
		res, err := stmt.ExecContext(
			ctx,
			entry.RequestID,
			entry.KeyID,
			entry.Model,
			entry.ProviderID,
			entry.ConnectionID,
			entry.Protocol,
			entry.InputTokens,
			entry.OutputTokens,
			entry.CachedTokens,
			entry.CostUSD,
			entry.LatencyMS,
			entry.TTFTMS,
			entry.StatusCode,
			entry.ErrorSummary,
			streamInt,
			entry.CreatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("usage: insert log %q: %w", entry.RequestID, err)
		}
		id, err := res.LastInsertId()
		if err != nil {
			return nil, fmt.Errorf("usage: read inserted log id: %w", err)
		}
		entry.ID = id
		inserted[i] = entry
	}

	if err := tx.Commit(); err != nil {
		return nil, fmt.Errorf("usage: commit batch insert: %w", err)
	}
	return inserted, nil
}

func (r *repository) ListLogs(ctx context.Context, filter LogsFilter) (LogsResponse, error) {
	if r == nil || r.db == nil {
		return LogsResponse{}, errors.New("usage: database is not initialized")
	}

	whereSQL, args, err := buildLogsWhere(filter)
	if err != nil {
		return LogsResponse{}, err
	}

	var total int64
	countQuery := `SELECT COUNT(*) FROM usage_logs` + whereSQL
	if err := r.db.Reader().QueryRowContext(ctx, countQuery, args...).Scan(&total); err != nil {
		return LogsResponse{}, fmt.Errorf("usage: count logs: %w", err)
	}

	queryArgs := append(append([]any{}, args...), filter.Limit, filter.Offset)
	rows, err := r.db.Reader().QueryContext(
		ctx,
		`SELECT id, request_id, key_id, model, provider_id, connection_id, protocol,
		        input_tokens, output_tokens, cached_tokens, cost_usd,
		        latency_ms, ttft_ms, status_code, error_summary, stream, created_at
		 FROM usage_logs`+whereSQL+`
		 ORDER BY created_at DESC, id DESC
		 LIMIT ? OFFSET ?`,
		queryArgs...,
	)
	if err != nil {
		return LogsResponse{}, fmt.Errorf("usage: query logs: %w", err)
	}
	defer func() {
		_ = rows.Close()
	}()

	logs := make([]UsageLog, 0)
	for rows.Next() {
		var (
			item      UsageLog
			streamInt int
		)
		if err := rows.Scan(
			&item.ID,
			&item.RequestID,
			&item.KeyID,
			&item.Model,
			&item.ProviderID,
			&item.ConnectionID,
			&item.Protocol,
			&item.InputTokens,
			&item.OutputTokens,
			&item.CachedTokens,
			&item.CostUSD,
			&item.LatencyMS,
			&item.TTFTMS,
			&item.StatusCode,
			&item.ErrorSummary,
			&streamInt,
			&item.CreatedAt,
		); err != nil {
			return LogsResponse{}, fmt.Errorf("usage: scan log row: %w", err)
		}
		item.Stream = streamInt != 0
		logs = append(logs, item)
	}
	if err := rows.Err(); err != nil {
		return LogsResponse{}, fmt.Errorf("usage: iterate log rows: %w", err)
	}

	return LogsResponse{
		Logs:   logs,
		Total:  total,
		Limit:  filter.Limit,
		Offset: filter.Offset,
	}, nil
}

func buildLogsWhere(filter LogsFilter) (whereSQL string, args []any, err error) {
	clauses := make([]string, 0, 4)
	args = make([]any, 0, 4)

	if filter.Model != "" {
		clauses = append(clauses, `model = ?`)
		args = append(args, filter.Model)
	}
	if filter.ProviderID != "" {
		clauses = append(clauses, `provider_id = ?`)
		args = append(args, filter.ProviderID)
	}
	if filter.KeyID != "" {
		clauses = append(clauses, `key_id = ?`)
		args = append(args, filter.KeyID)
	}
	if filter.Status != "" {
		switch strings.ToLower(filter.Status) {
		case "error":
			clauses = append(clauses, errorPredicateSQL)
		case "success", "ok":
			clauses = append(clauses, successPredicateSQL)
		default:
			code, err := strconv.Atoi(filter.Status)
			if err != nil || code < 100 || code > 599 {
				return "", nil, errors.New("status must be 'success', 'error', or an HTTP status code")
			}
			clauses = append(clauses, `status_code = ?`)
			args = append(args, code)
		}
	}

	if len(clauses) == 0 {
		return "", nil, nil
	}
	return " WHERE " + strings.Join(clauses, " AND "), args, nil
}

func (r *repository) Summary(ctx context.Context, rng string) (SummaryResponse, error) {
	if r == nil || r.db == nil {
		return SummaryResponse{}, errors.New("usage: database is not initialized")
	}

	whereSQL, args, bucketExpr, err := r.summaryWindow(rng)
	if err != nil {
		return SummaryResponse{}, err
	}

	var totals SummaryTotals
	totalsQuery := `SELECT
		COUNT(*),
		COALESCE(SUM(input_tokens), 0),
		COALESCE(SUM(output_tokens), 0),
		COALESCE(SUM(cached_tokens), 0),
		COALESCE(SUM(cost_usd), 0.0),
		COALESCE(AVG(latency_ms), 0.0),
		COALESCE(AVG(ttft_ms), 0.0),
		COALESCE(SUM(CASE WHEN ` + errorPredicateSQL + ` THEN 1 ELSE 0 END), 0)
	FROM usage_logs` + whereSQL
	if err := r.db.Reader().QueryRowContext(ctx, totalsQuery, args...).Scan(
		&totals.Requests,
		&totals.InputTokens,
		&totals.OutputTokens,
		&totals.CachedTokens,
		&totals.CostUSD,
		&totals.AvgLatencyMS,
		&totals.AvgTTFTMS,
		&totals.ErrorCount,
	); err != nil {
		return SummaryResponse{}, fmt.Errorf("usage: query summary totals: %w", err)
	}
	if totals.Requests > 0 {
		totals.ErrorRate = float64(totals.ErrorCount) / float64(totals.Requests)
	}

	byModel, err := r.queryModelBreakdown(ctx, whereSQL, args)
	if err != nil {
		return SummaryResponse{}, err
	}
	byProvider, err := r.queryProviderBreakdown(ctx, whereSQL, args)
	if err != nil {
		return SummaryResponse{}, err
	}
	series, err := r.querySeries(ctx, whereSQL, args, bucketExpr)
	if err != nil {
		return SummaryResponse{}, err
	}

	return SummaryResponse{
		Range:      rng,
		Totals:     totals,
		ByModel:    byModel,
		ByProvider: byProvider,
		Series:     series,
	}, nil
}

func (r *repository) summaryWindow(rng string) (whereSQL string, args []any, bucketExpr string, err error) {
	now := r.now().UTC()
	switch rng {
	case "24h":
		cutoff := now.Add(-24 * time.Hour).Format(rfc3339NanoFixed)
		return ` WHERE created_at >= ?`, []any{cutoff}, `substr(created_at, 1, 13) || ':00:00Z'`, nil
	case "7d":
		cutoff := now.Add(-7 * 24 * time.Hour).Format(rfc3339NanoFixed)
		return ` WHERE created_at >= ?`, []any{cutoff}, `substr(created_at, 1, 10)`, nil
	case "30d":
		cutoff := now.Add(-30 * 24 * time.Hour).Format(rfc3339NanoFixed)
		return ` WHERE created_at >= ?`, []any{cutoff}, `substr(created_at, 1, 10)`, nil
	case "all":
		return "", nil, `substr(created_at, 1, 10)`, nil
	default:
		return "", nil, "", errors.New("range must be one of: 24h, 7d, 30d, all")
	}
}

func (r *repository) queryModelBreakdown(ctx context.Context, whereSQL string, args []any) ([]ModelBreakdown, error) {
	rows, err := r.db.Reader().QueryContext(
		ctx,
		`SELECT
			model,
			COUNT(*),
			COALESCE(SUM(input_tokens), 0),
			COALESCE(SUM(output_tokens), 0),
			COALESCE(SUM(cached_tokens), 0),
			COALESCE(SUM(cost_usd), 0.0),
			COALESCE(AVG(latency_ms), 0.0),
			COALESCE(SUM(CASE WHEN `+errorPredicateSQL+` THEN 1 ELSE 0 END), 0)
		FROM usage_logs`+whereSQL+`
		GROUP BY model
		ORDER BY SUM(cost_usd) DESC, COUNT(*) DESC, model ASC`,
		args...,
	)
	if err != nil {
		return nil, fmt.Errorf("usage: query model breakdown: %w", err)
	}
	defer func() {
		_ = rows.Close()
	}()

	out := make([]ModelBreakdown, 0)
	for rows.Next() {
		var item ModelBreakdown
		if err := rows.Scan(
			&item.Model,
			&item.Requests,
			&item.InputTokens,
			&item.OutputTokens,
			&item.CachedTokens,
			&item.CostUSD,
			&item.AvgLatencyMS,
			&item.ErrorCount,
		); err != nil {
			return nil, fmt.Errorf("usage: scan model breakdown: %w", err)
		}
		out = append(out, item)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("usage: iterate model breakdown: %w", err)
	}
	return out, nil
}

func (r *repository) queryProviderBreakdown(ctx context.Context, whereSQL string, args []any) ([]ProviderBreakdown, error) {
	rows, err := r.db.Reader().QueryContext(
		ctx,
		`SELECT
			provider_id,
			COUNT(*),
			COALESCE(SUM(input_tokens), 0),
			COALESCE(SUM(output_tokens), 0),
			COALESCE(SUM(cached_tokens), 0),
			COALESCE(SUM(cost_usd), 0.0),
			COALESCE(AVG(latency_ms), 0.0),
			COALESCE(SUM(CASE WHEN `+errorPredicateSQL+` THEN 1 ELSE 0 END), 0)
		FROM usage_logs`+whereSQL+`
		GROUP BY provider_id
		ORDER BY SUM(cost_usd) DESC, COUNT(*) DESC, provider_id ASC`,
		args...,
	)
	if err != nil {
		return nil, fmt.Errorf("usage: query provider breakdown: %w", err)
	}
	defer func() {
		_ = rows.Close()
	}()

	out := make([]ProviderBreakdown, 0)
	for rows.Next() {
		var item ProviderBreakdown
		if err := rows.Scan(
			&item.ProviderID,
			&item.Requests,
			&item.InputTokens,
			&item.OutputTokens,
			&item.CachedTokens,
			&item.CostUSD,
			&item.AvgLatencyMS,
			&item.ErrorCount,
		); err != nil {
			return nil, fmt.Errorf("usage: scan provider breakdown: %w", err)
		}
		out = append(out, item)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("usage: iterate provider breakdown: %w", err)
	}
	return out, nil
}

func (r *repository) querySeries(ctx context.Context, whereSQL string, args []any, bucketExpr string) ([]SeriesBucket, error) {
	rows, err := r.db.Reader().QueryContext(
		ctx,
		`SELECT
			`+bucketExpr+` AS bucket,
			COUNT(*),
			COALESCE(SUM(input_tokens), 0),
			COALESCE(SUM(output_tokens), 0),
			COALESCE(SUM(cached_tokens), 0),
			COALESCE(SUM(cost_usd), 0.0),
			COALESCE(SUM(CASE WHEN `+errorPredicateSQL+` THEN 1 ELSE 0 END), 0)
		FROM usage_logs`+whereSQL+`
		GROUP BY bucket
		ORDER BY bucket ASC`,
		args...,
	)
	if err != nil {
		return nil, fmt.Errorf("usage: query usage series: %w", err)
	}
	defer func() {
		_ = rows.Close()
	}()

	out := make([]SeriesBucket, 0)
	for rows.Next() {
		var item SeriesBucket
		if err := rows.Scan(
			&item.Bucket,
			&item.Requests,
			&item.InputTokens,
			&item.OutputTokens,
			&item.CachedTokens,
			&item.CostUSD,
			&item.ErrorCount,
		); err != nil {
			return nil, fmt.Errorf("usage: scan usage series: %w", err)
		}
		out = append(out, item)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("usage: iterate usage series: %w", err)
	}
	return out, nil
}

func (r *repository) PruneBefore(ctx context.Context, cutoff string) (int64, error) {
	if r == nil || r.db == nil {
		return 0, errors.New("usage: database is not initialized")
	}

	res, err := r.db.Writer().ExecContext(ctx, `DELETE FROM usage_logs WHERE created_at < ?`, cutoff)
	if err != nil {
		return 0, fmt.Errorf("usage: prune expired logs: %w", err)
	}
	deleted, err := res.RowsAffected()
	if err != nil {
		return 0, fmt.Errorf("usage: count pruned logs: %w", err)
	}
	return deleted, nil
}

func (r *repository) ListPricingOverrides(ctx context.Context) ([]ModelPricing, error) {
	if r == nil || r.db == nil {
		return nil, errors.New("usage: database is not initialized")
	}

	rows, err := r.db.Reader().QueryContext(
		ctx,
		`SELECT model, input_per_million, output_per_million, cached_per_million, updated_at
		 FROM model_pricing
		 ORDER BY model ASC`,
	)
	if err != nil {
		return nil, fmt.Errorf("usage: list model pricing overrides: %w", err)
	}
	defer func() {
		_ = rows.Close()
	}()

	out := make([]ModelPricing, 0)
	for rows.Next() {
		var item ModelPricing
		if err := rows.Scan(
			&item.Model,
			&item.InputPerMillion,
			&item.OutputPerMillion,
			&item.CachedPerMillion,
			&item.UpdatedAt,
		); err != nil {
			return nil, fmt.Errorf("usage: scan model pricing override: %w", err)
		}
		item.Source = pricingSourceCustom
		out = append(out, item)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("usage: iterate model pricing overrides: %w", err)
	}
	return out, nil
}

func (r *repository) UpsertPricingOverride(ctx context.Context, model string, inputRate, outputRate, cachedRate float64) (ModelPricing, error) {
	if r == nil || r.db == nil {
		return ModelPricing{}, errors.New("usage: database is not initialized")
	}

	nowStr := r.now().UTC().Format(rfc3339NanoFixed)
	if _, err := r.db.Writer().ExecContext(
		ctx,
		`INSERT INTO model_pricing (model, input_per_million, output_per_million, cached_per_million, updated_at)
		 VALUES (?, ?, ?, ?, ?)
		 ON CONFLICT(model) DO UPDATE SET
			input_per_million = excluded.input_per_million,
			output_per_million = excluded.output_per_million,
			cached_per_million = excluded.cached_per_million,
			updated_at = excluded.updated_at`,
		model,
		inputRate,
		outputRate,
		cachedRate,
		nowStr,
	); err != nil {
		return ModelPricing{}, fmt.Errorf("usage: upsert model pricing %q: %w", model, err)
	}

	return ModelPricing{
		Model:            model,
		InputPerMillion:  inputRate,
		OutputPerMillion: outputRate,
		CachedPerMillion: cachedRate,
		Source:           pricingSourceCustom,
		UpdatedAt:        nowStr,
	}, nil
}

func (r *repository) DeletePricingOverride(ctx context.Context, model string) error {
	if r == nil || r.db == nil {
		return errors.New("usage: database is not initialized")
	}

	res, err := r.db.Writer().ExecContext(ctx, `DELETE FROM model_pricing WHERE model = ?`, model)
	if err != nil {
		return fmt.Errorf("usage: delete model pricing %q: %w", model, err)
	}
	affected, err := res.RowsAffected()
	if err != nil {
		return fmt.Errorf("usage: read deleted pricing rows for %q: %w", model, err)
	}
	if affected == 0 {
		return ErrPricingNotFound
	}
	return nil
}

// ensure sql package reference stays explicit for future QueryRow checks.
var _ = sql.ErrNoRows
