package security

import (
	"context"
	"crypto/rand"
	"database/sql"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/omeg4monx/router-gateway/internal/platform/db"
)

const (
	timestampLayout    = "2006-01-02T15:04:05.000000000Z"
	defaultListLimit   = 50
	maxListLimit       = 500
	maxAuditRetention  = 90 * 24 * time.Hour
	maxRetainedRecords = 5000
)

// AuditEvent represents a persisted, redacted security audit trail record.
type AuditEvent struct {
	ID         string         `json:"id"`
	Category   string         `json:"category"`
	Action     string         `json:"action"`
	Severity   string         `json:"severity"`
	Actor      string         `json:"actor"`
	ResourceID string         `json:"resource_id"`
	Summary    string         `json:"summary"`
	Metadata   map[string]any `json:"metadata"`
	IPAddress  string         `json:"ip_address"`
	CreatedAt  string         `json:"created_at"`
}

// AuditInput describes a new security audit entry prior to redaction and storage.
type AuditInput struct {
	Category   string         `json:"category"`
	Action     string         `json:"action"`
	Severity   string         `json:"severity"`
	Actor      string         `json:"actor"`
	ResourceID string         `json:"resource_id"`
	Summary    string         `json:"summary"`
	Metadata   map[string]any `json:"metadata"`
	IPAddress  string         `json:"ip_address"`
}

// ListFilter narrows security audit queries.
type ListFilter struct {
	Category string
	Severity string
	Query    string
	Limit    int
	Offset   int
}

type repo struct {
	db *db.DB
}

func newRepo(database *db.DB) *repo {
	return &repo{db: database}
}

func formatUTC(t time.Time) string {
	return t.UTC().Format(timestampLayout)
}

func newEventID() string {
	var buf [12]byte
	if _, err := rand.Read(buf[:]); err != nil {
		return fmt.Sprintf("aud_%d", time.Now().UTC().UnixNano())
	}
	return "aud_" + hex.EncodeToString(buf[:])
}

func (r *repo) insert(ctx context.Context, in AuditInput) (AuditEvent, error) {
	ev := AuditEvent{
		ID:         newEventID(),
		Category:   strings.TrimSpace(in.Category),
		Action:     strings.TrimSpace(in.Action),
		Severity:   normalizeSeverity(in.Severity),
		Actor:      strings.TrimSpace(in.Actor),
		ResourceID: strings.TrimSpace(in.ResourceID),
		Summary:    RedactText(in.Summary),
		Metadata:   RedactMetadata(in.Metadata),
		IPAddress:  strings.TrimSpace(in.IPAddress),
		CreatedAt:  formatUTC(time.Now()),
	}
	if ev.Actor == "" {
		ev.Actor = "system"
	}

	metaBytes, err := json.Marshal(ev.Metadata)
	if err != nil {
		metaBytes = []byte("{}")
	}

	_, err = r.db.Writer().ExecContext(ctx,
		`INSERT INTO security_audit_events
		 (id, category, action, severity, actor, resource_id, summary, metadata_json, ip_address, created_at)
		 VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`,
		ev.ID, ev.Category, ev.Action, ev.Severity, ev.Actor, ev.ResourceID,
		ev.Summary, string(metaBytes), ev.IPAddress, ev.CreatedAt,
	)
	if err != nil {
		return AuditEvent{}, fmt.Errorf("insert audit event: %w", err)
	}
	return ev, nil
}

func (r *repo) list(ctx context.Context, filter ListFilter) ([]AuditEvent, int, error) {
	limit := filter.Limit
	if limit <= 0 {
		limit = defaultListLimit
	}
	if limit > maxListLimit {
		limit = maxListLimit
	}

	offset := max(filter.Offset, 0)

	where, args := buildWhere(filter)

	var total int
	countQuery := "SELECT COUNT(*) FROM security_audit_events" + where
	if err := r.db.Reader().QueryRowContext(ctx, countQuery, args...).Scan(&total); err != nil {
		return nil, 0, fmt.Errorf("count audit events: %w", err)
	}

	query := `SELECT id, category, action, severity, actor, resource_id, summary, metadata_json, ip_address, created_at
		FROM security_audit_events` + where + ` ORDER BY created_at DESC, id DESC LIMIT ? OFFSET ?`
	rows, err := r.db.Reader().QueryContext(ctx, query, append(args, limit, offset)...)
	if err != nil {
		return nil, 0, fmt.Errorf("list audit events: %w", err)
	}
	defer func() { _ = rows.Close() }()

	items := make([]AuditEvent, 0)
	for rows.Next() {
		ev, err := scanAuditEvent(rows)
		if err != nil {
			return nil, 0, err
		}
		items = append(items, ev)
	}
	if err := rows.Err(); err != nil {
		return nil, 0, fmt.Errorf("iterate audit events: %w", err)
	}

	return items, total, nil
}

func (r *repo) listAllForExport(ctx context.Context) ([]AuditEvent, error) {
	rows, err := r.db.Reader().QueryContext(ctx,
		`SELECT id, category, action, severity, actor, resource_id, summary, metadata_json, ip_address, created_at
		 FROM security_audit_events ORDER BY created_at DESC, id DESC LIMIT ?`,
		maxRetainedRecords,
	)
	if err != nil {
		return nil, fmt.Errorf("export audit events: %w", err)
	}
	defer func() { _ = rows.Close() }()

	items := make([]AuditEvent, 0)
	for rows.Next() {
		ev, err := scanAuditEvent(rows)
		if err != nil {
			return nil, err
		}
		items = append(items, ev)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate export audit events: %w", err)
	}
	return items, nil
}

func (r *repo) countAll(ctx context.Context) (int, error) {
	var total int
	if err := r.db.Reader().QueryRowContext(ctx, `SELECT COUNT(*) FROM security_audit_events`).Scan(&total); err != nil {
		return 0, fmt.Errorf("count all audit events: %w", err)
	}
	return total, nil
}

// clearAndRecordMarker wipes existing audit records and atomically inserts a
// tamper-evident audit.cleared record inside the same transaction.
func (r *repo) clearAndRecordMarker(ctx context.Context, actor, ipAddress string) (int64, AuditEvent, error) {
	tx, err := r.db.Writer().BeginTx(ctx, nil)
	if err != nil {
		return 0, AuditEvent{}, fmt.Errorf("begin clear audit tx: %w", err)
	}
	defer func() { _ = tx.Rollback() }()

	res, err := tx.ExecContext(ctx, `DELETE FROM security_audit_events`)
	if err != nil {
		return 0, AuditEvent{}, fmt.Errorf("delete audit events: %w", err)
	}
	deleted, _ := res.RowsAffected()

	if strings.TrimSpace(actor) == "" {
		actor = "admin"
	}
	marker := AuditEvent{
		ID:        newEventID(),
		Category:  "security",
		Action:    "audit.cleared",
		Severity:  "warning",
		Actor:     actor,
		Summary:   fmt.Sprintf("Security audit trail cleared (%d prior events removed)", deleted),
		Metadata:  map[string]any{"cleared_count": deleted},
		IPAddress: strings.TrimSpace(ipAddress),
		CreatedAt: formatUTC(time.Now()),
	}
	metaBytes, _ := json.Marshal(marker.Metadata)

	_, err = tx.ExecContext(ctx,
		`INSERT INTO security_audit_events
		 (id, category, action, severity, actor, resource_id, summary, metadata_json, ip_address, created_at)
		 VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`,
		marker.ID, marker.Category, marker.Action, marker.Severity, marker.Actor,
		marker.ResourceID, marker.Summary, string(metaBytes), marker.IPAddress, marker.CreatedAt,
	)
	if err != nil {
		return 0, AuditEvent{}, fmt.Errorf("insert clear marker: %w", err)
	}

	if err := tx.Commit(); err != nil {
		return 0, AuditEvent{}, fmt.Errorf("commit clear audit tx: %w", err)
	}
	return deleted, marker, nil
}

func (r *repo) prune(ctx context.Context, cutoff time.Time, maxRows int) (int64, error) {
	var totalDeleted int64

	res, err := r.db.Writer().ExecContext(ctx,
		`DELETE FROM security_audit_events WHERE created_at < ?`,
		formatUTC(cutoff),
	)
	if err != nil {
		return 0, fmt.Errorf("prune expired audit events: %w", err)
	}
	if n, err := res.RowsAffected(); err == nil {
		totalDeleted += n
	}

	if maxRows > 0 {
		res2, err := r.db.Writer().ExecContext(ctx,
			`DELETE FROM security_audit_events WHERE id NOT IN (
				SELECT id FROM security_audit_events ORDER BY created_at DESC, id DESC LIMIT ?
			)`,
			maxRows,
		)
		if err != nil {
			return totalDeleted, fmt.Errorf("cap audit events: %w", err)
		}
		if n, err := res2.RowsAffected(); err == nil {
			totalDeleted += n
		}
	}

	return totalDeleted, nil
}

func buildWhere(filter ListFilter) (where string, args []any) {
	var clauses []string

	if cat := strings.TrimSpace(filter.Category); cat != "" {
		clauses = append(clauses, "category = ?")
		args = append(args, cat)
	}
	if sev := strings.TrimSpace(filter.Severity); sev != "" {
		clauses = append(clauses, "severity = ?")
		args = append(args, sev)
	}
	if q := strings.TrimSpace(filter.Query); q != "" {
		like := "%" + q + "%"
		clauses = append(clauses, "(summary LIKE ? OR action LIKE ? OR actor LIKE ? OR resource_id LIKE ?)")
		args = append(args, like, like, like, like)
	}

	if len(clauses) == 0 {
		return "", nil
	}
	return " WHERE " + strings.Join(clauses, " AND "), args
}

func scanAuditEvent(rows *sql.Rows) (AuditEvent, error) {
	var ev AuditEvent
	var metaJSON string
	if err := rows.Scan(
		&ev.ID, &ev.Category, &ev.Action, &ev.Severity, &ev.Actor,
		&ev.ResourceID, &ev.Summary, &metaJSON, &ev.IPAddress, &ev.CreatedAt,
	); err != nil {
		return AuditEvent{}, fmt.Errorf("scan audit event: %w", err)
	}
	ev.Metadata = map[string]any{}
	if strings.TrimSpace(metaJSON) != "" {
		_ = json.Unmarshal([]byte(metaJSON), &ev.Metadata)
	}
	return ev, nil
}

func normalizeSeverity(s string) string {
	switch strings.ToLower(strings.TrimSpace(s)) {
	case "warning", "warn":
		return "warning"
	case "critical", "error", "high":
		return "critical"
	default:
		return "info"
	}
}
