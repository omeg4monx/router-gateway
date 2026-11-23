package notifications

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

const (
	timestampLayout         = "2006-01-02T15:04:05.000000000Z"
	defaultListLimit        = 50
	maxListLimit            = 200
	maxNotificationAge      = 30 * 24 * time.Hour
	maxRetainedNotification = 1000
)

var errNotificationNotFound = errors.New("notifications: not found")

// Notification represents a persisted system notification item.
type Notification struct {
	ID        string         `json:"id"`
	Title     string         `json:"title"`
	Message   string         `json:"message"`
	Severity  string         `json:"severity"`
	Category  string         `json:"category"`
	Source    string         `json:"source"`
	Metadata  map[string]any `json:"metadata"`
	ReadAt    string         `json:"read_at"`
	Read      bool           `json:"read"`
	CreatedAt string         `json:"created_at"`
}

// EmitPayload describes a new notification to persist and fan out.
type EmitPayload struct {
	Title    string         `json:"title"`
	Message  string         `json:"message"`
	Severity string         `json:"severity"`
	Category string         `json:"category"`
	Source   string         `json:"source"`
	Metadata map[string]any `json:"metadata"`
}

// ListFilter narrows notification queries.
type ListFilter struct {
	UnreadOnly bool
	Severity   string
	Limit      int
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

func newNotificationID() string {
	var buf [12]byte
	if _, err := rand.Read(buf[:]); err != nil {
		return fmt.Sprintf("ntf_%d", time.Now().UTC().UnixNano())
	}
	return "ntf_" + hex.EncodeToString(buf[:])
}

func (r *repo) insert(ctx context.Context, in EmitPayload) (Notification, error) {
	n := Notification{
		ID:        newNotificationID(),
		Title:     strings.TrimSpace(in.Title),
		Message:   strings.TrimSpace(in.Message),
		Severity:  normalizeSeverity(in.Severity),
		Category:  strings.TrimSpace(in.Category),
		Source:    strings.TrimSpace(in.Source),
		Metadata:  in.Metadata,
		ReadAt:    "",
		Read:      false,
		CreatedAt: formatUTC(time.Now()),
	}
	if n.Category == "" {
		n.Category = "system"
	}
	if n.Source == "" {
		n.Source = "system"
	}
	if n.Metadata == nil {
		n.Metadata = map[string]any{}
	}

	metaBytes, err := json.Marshal(n.Metadata)
	if err != nil {
		metaBytes = []byte("{}")
	}

	_, err = r.db.Writer().ExecContext(ctx,
		`INSERT INTO system_notifications
		 (id, title, message, severity, category, source, metadata_json, read_at, created_at)
		 VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)`,
		n.ID, n.Title, n.Message, n.Severity, n.Category, n.Source, string(metaBytes), n.ReadAt, n.CreatedAt,
	)
	if err != nil {
		return Notification{}, fmt.Errorf("insert notification: %w", err)
	}
	return n, nil
}

func (r *repo) list(ctx context.Context, filter ListFilter) (items []Notification, total, unread int, err error) {
	limit := filter.Limit
	if limit <= 0 {
		limit = defaultListLimit
	}
	if limit > maxListLimit {
		limit = maxListLimit
	}

	where, args := buildWhere(filter)

	if err = r.db.Reader().QueryRowContext(ctx, `SELECT COUNT(*) FROM system_notifications WHERE read_at = ''`).Scan(&unread); err != nil {
		return nil, 0, 0, fmt.Errorf("count unread notifications: %w", err)
	}
	if err = r.db.Reader().QueryRowContext(ctx, "SELECT COUNT(*) FROM system_notifications"+where, args...).Scan(&total); err != nil {
		return nil, 0, 0, fmt.Errorf("count notifications: %w", err)
	}

	query := `SELECT id, title, message, severity, category, source, metadata_json, read_at, created_at
		FROM system_notifications` + where + ` ORDER BY created_at DESC, id DESC LIMIT ?`
	rows, err := r.db.Reader().QueryContext(ctx, query, append(args, limit)...)
	if err != nil {
		return nil, 0, 0, fmt.Errorf("list notifications: %w", err)
	}
	defer func() { _ = rows.Close() }()

	items = make([]Notification, 0)
	for rows.Next() {
		item, scanErr := scanNotification(rows)
		if scanErr != nil {
			return nil, 0, 0, scanErr
		}
		items = append(items, item)
	}
	if err = rows.Err(); err != nil {
		return nil, 0, 0, fmt.Errorf("iterate notifications: %w", err)
	}

	return items, total, unread, nil
}

func (r *repo) markRead(ctx context.Context, id string) (Notification, error) {
	now := formatUTC(time.Now())
	res, err := r.db.Writer().ExecContext(ctx,
		`UPDATE system_notifications SET read_at = CASE WHEN read_at = '' THEN ? ELSE read_at END WHERE id = ?`,
		now, id,
	)
	if err != nil {
		return Notification{}, fmt.Errorf("mark notification read: %w", err)
	}
	affected, _ := res.RowsAffected()
	if affected == 0 {
		return Notification{}, errNotificationNotFound
	}

	row := r.db.Reader().QueryRowContext(ctx,
		`SELECT id, title, message, severity, category, source, metadata_json, read_at, created_at
		 FROM system_notifications WHERE id = ?`, id,
	)
	return scanNotificationRow(row)
}

func (r *repo) markAllRead(ctx context.Context) (int64, error) {
	now := formatUTC(time.Now())
	res, err := r.db.Writer().ExecContext(ctx,
		`UPDATE system_notifications SET read_at = ? WHERE read_at = ''`,
		now,
	)
	if err != nil {
		return 0, fmt.Errorf("mark all notifications read: %w", err)
	}
	n, _ := res.RowsAffected()
	return n, nil
}

func (r *repo) deleteByID(ctx context.Context, id string) error {
	res, err := r.db.Writer().ExecContext(ctx, `DELETE FROM system_notifications WHERE id = ?`, id)
	if err != nil {
		return fmt.Errorf("delete notification: %w", err)
	}
	affected, _ := res.RowsAffected()
	if affected == 0 {
		return errNotificationNotFound
	}
	return nil
}

func (r *repo) clearAll(ctx context.Context) (int64, error) {
	res, err := r.db.Writer().ExecContext(ctx, `DELETE FROM system_notifications`)
	if err != nil {
		return 0, fmt.Errorf("clear notifications: %w", err)
	}
	n, _ := res.RowsAffected()
	return n, nil
}

func (r *repo) prune(ctx context.Context, cutoff time.Time, maxRows int) (int64, error) {
	var totalDeleted int64

	res, err := r.db.Writer().ExecContext(ctx,
		`DELETE FROM system_notifications WHERE created_at < ?`,
		formatUTC(cutoff),
	)
	if err != nil {
		return 0, fmt.Errorf("prune expired notifications: %w", err)
	}
	if n, err := res.RowsAffected(); err == nil {
		totalDeleted += n
	}

	if maxRows > 0 {
		res2, err := r.db.Writer().ExecContext(ctx,
			`DELETE FROM system_notifications WHERE id NOT IN (
				SELECT id FROM system_notifications ORDER BY created_at DESC, id DESC LIMIT ?
			)`,
			maxRows,
		)
		if err != nil {
			return totalDeleted, fmt.Errorf("cap notifications: %w", err)
		}
		if n, err := res2.RowsAffected(); err == nil {
			totalDeleted += n
		}
	}

	return totalDeleted, nil
}

func buildWhere(filter ListFilter) (where string, args []any) {
	var clauses []string

	if filter.UnreadOnly {
		clauses = append(clauses, "read_at = ''")
	}
	if sev := strings.TrimSpace(filter.Severity); sev != "" {
		clauses = append(clauses, "severity = ?")
		args = append(args, sev)
	}

	if len(clauses) == 0 {
		return "", nil
	}
	return " WHERE " + strings.Join(clauses, " AND "), args
}

func scanNotification(rows *sql.Rows) (Notification, error) {
	var n Notification
	var metaJSON string
	if err := rows.Scan(
		&n.ID, &n.Title, &n.Message, &n.Severity, &n.Category,
		&n.Source, &metaJSON, &n.ReadAt, &n.CreatedAt,
	); err != nil {
		return Notification{}, fmt.Errorf("scan notification: %w", err)
	}
	n.Read = strings.TrimSpace(n.ReadAt) != ""
	n.Metadata = map[string]any{}
	if strings.TrimSpace(metaJSON) != "" {
		_ = json.Unmarshal([]byte(metaJSON), &n.Metadata)
	}
	return n, nil
}

func scanNotificationRow(row *sql.Row) (Notification, error) {
	var n Notification
	var metaJSON string
	if err := row.Scan(
		&n.ID, &n.Title, &n.Message, &n.Severity, &n.Category,
		&n.Source, &metaJSON, &n.ReadAt, &n.CreatedAt,
	); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return Notification{}, errNotificationNotFound
		}
		return Notification{}, fmt.Errorf("scan notification row: %w", err)
	}
	n.Read = strings.TrimSpace(n.ReadAt) != ""
	n.Metadata = map[string]any{}
	if strings.TrimSpace(metaJSON) != "" {
		_ = json.Unmarshal([]byte(metaJSON), &n.Metadata)
	}
	return n, nil
}

func normalizeSeverity(s string) string {
	switch strings.ToLower(strings.TrimSpace(s)) {
	case "success", "ok":
		return "success"
	case "warning", "warn":
		return "warning"
	case "error", "critical":
		return "error"
	default:
		return "info"
	}
}
