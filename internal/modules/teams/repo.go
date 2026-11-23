package teams

import (
	"context"
	"crypto/rand"
	"database/sql"
	"encoding/hex"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/omeg4monx/router-gateway/internal/platform/db"
)

const timestampLayout = "2006-01-02T15:04:05.000000000Z"

var (
	errTeamNotFound = errors.New("teams: team not found")
	errTeamExists   = errors.New("teams: team name already exists")
	errKeyNotInTeam = errors.New("teams: key not assigned to team")
)

// Team represents an API-key grouping with shared USD budget and RPM quota.
type Team struct {
	ID           string   `json:"id"`
	Name         string   `json:"name"`
	Description  string   `json:"description"`
	BudgetUSD    float64  `json:"budget_usd"`
	SpentUSD     float64  `json:"spent_usd"`
	RateLimitRPM int      `json:"rate_limit_rpm"`
	Enabled      bool     `json:"enabled"`
	KeyIDs       []string `json:"key_ids"`
	CreatedAt    string   `json:"created_at"`
	UpdatedAt    string   `json:"updated_at"`
}

type createTeamInput struct {
	Name         string
	Description  string
	BudgetUSD    float64
	RateLimitRPM int
	Enabled      bool
}

type patchTeamInput struct {
	Name         *string
	Description  *string
	BudgetUSD    *float64
	RateLimitRPM *int
	Enabled      *bool
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

func newTeamID() string {
	var buf [10]byte
	if _, err := rand.Read(buf[:]); err != nil {
		return fmt.Sprintf("team_%d", time.Now().UTC().UnixNano())
	}
	return "team_" + hex.EncodeToString(buf[:])
}

func (r *repo) list(ctx context.Context) ([]Team, error) {
	rows, err := r.db.Reader().QueryContext(ctx,
		`SELECT id, name, description, budget_usd, spent_usd, rate_limit_rpm, enabled, created_at, updated_at
		 FROM teams ORDER BY name ASC`,
	)
	if err != nil {
		return nil, fmt.Errorf("list teams: %w", err)
	}
	defer func() { _ = rows.Close() }()

	items := make([]Team, 0)
	byID := make(map[string]int)
	for rows.Next() {
		var t Team
		var enabled int
		if err := rows.Scan(
			&t.ID, &t.Name, &t.Description, &t.BudgetUSD, &t.SpentUSD,
			&t.RateLimitRPM, &enabled, &t.CreatedAt, &t.UpdatedAt,
		); err != nil {
			return nil, fmt.Errorf("scan team: %w", err)
		}
		t.Enabled = enabled != 0
		t.KeyIDs = []string{}
		byID[t.ID] = len(items)
		items = append(items, t)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate teams: %w", err)
	}

	kRows, err := r.db.Reader().QueryContext(ctx, `SELECT team_id, key_id FROM team_keys ORDER BY key_id ASC`)
	if err != nil {
		return nil, fmt.Errorf("list team keys: %w", err)
	}
	defer func() { _ = kRows.Close() }()

	for kRows.Next() {
		var teamID, keyID string
		if err := kRows.Scan(&teamID, &keyID); err != nil {
			return nil, fmt.Errorf("scan team key: %w", err)
		}
		if idx, ok := byID[teamID]; ok {
			items[idx].KeyIDs = append(items[idx].KeyIDs, keyID)
		}
	}
	if err := kRows.Err(); err != nil {
		return nil, fmt.Errorf("iterate team keys: %w", err)
	}

	return items, nil
}

func (r *repo) getByID(ctx context.Context, id string) (Team, error) {
	var t Team
	var enabled int
	err := r.db.Reader().QueryRowContext(ctx,
		`SELECT id, name, description, budget_usd, spent_usd, rate_limit_rpm, enabled, created_at, updated_at
		 FROM teams WHERE id = ?`, id,
	).Scan(
		&t.ID, &t.Name, &t.Description, &t.BudgetUSD, &t.SpentUSD,
		&t.RateLimitRPM, &enabled, &t.CreatedAt, &t.UpdatedAt,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return Team{}, errTeamNotFound
		}
		return Team{}, fmt.Errorf("get team: %w", err)
	}
	t.Enabled = enabled != 0
	t.KeyIDs = []string{}

	rows, err := r.db.Reader().QueryContext(ctx, `SELECT key_id FROM team_keys WHERE team_id = ? ORDER BY key_id ASC`, id)
	if err != nil {
		return Team{}, fmt.Errorf("list keys for team: %w", err)
	}
	defer func() { _ = rows.Close() }()
	for rows.Next() {
		var keyID string
		if err := rows.Scan(&keyID); err != nil {
			return Team{}, fmt.Errorf("scan team key: %w", err)
		}
		t.KeyIDs = append(t.KeyIDs, keyID)
	}
	return t, rows.Err()
}

func (r *repo) create(ctx context.Context, in createTeamInput) (Team, error) {
	now := formatUTC(time.Now())
	t := Team{
		ID:           newTeamID(),
		Name:         strings.TrimSpace(in.Name),
		Description:  strings.TrimSpace(in.Description),
		BudgetUSD:    in.BudgetUSD,
		SpentUSD:     0,
		RateLimitRPM: in.RateLimitRPM,
		Enabled:      in.Enabled,
		KeyIDs:       []string{},
		CreatedAt:    now,
		UpdatedAt:    now,
	}
	enabledInt := 0
	if t.Enabled {
		enabledInt = 1
	}

	_, err := r.db.Writer().ExecContext(ctx,
		`INSERT INTO teams (id, name, description, budget_usd, spent_usd, rate_limit_rpm, enabled, created_at, updated_at)
		 VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)`,
		t.ID, t.Name, t.Description, t.BudgetUSD, t.SpentUSD, t.RateLimitRPM, enabledInt, t.CreatedAt, t.UpdatedAt,
	)
	if err != nil {
		if strings.Contains(strings.ToLower(err.Error()), "unique") {
			return Team{}, errTeamExists
		}
		return Team{}, fmt.Errorf("insert team: %w", err)
	}
	return t, nil
}

func (r *repo) patch(ctx context.Context, id string, in patchTeamInput) (Team, error) {
	cur, err := r.getByID(ctx, id)
	if err != nil {
		return Team{}, err
	}
	if in.Name != nil {
		cur.Name = strings.TrimSpace(*in.Name)
	}
	if in.Description != nil {
		cur.Description = strings.TrimSpace(*in.Description)
	}
	if in.BudgetUSD != nil {
		cur.BudgetUSD = *in.BudgetUSD
	}
	if in.RateLimitRPM != nil {
		cur.RateLimitRPM = *in.RateLimitRPM
	}
	if in.Enabled != nil {
		cur.Enabled = *in.Enabled
	}
	cur.UpdatedAt = formatUTC(time.Now())

	enabledInt := 0
	if cur.Enabled {
		enabledInt = 1
	}
	_, err = r.db.Writer().ExecContext(ctx,
		`UPDATE teams SET name = ?, description = ?, budget_usd = ?, rate_limit_rpm = ?, enabled = ?, updated_at = ?
		 WHERE id = ?`,
		cur.Name, cur.Description, cur.BudgetUSD, cur.RateLimitRPM, enabledInt, cur.UpdatedAt, id,
	)
	if err != nil {
		if strings.Contains(strings.ToLower(err.Error()), "unique") {
			return Team{}, errTeamExists
		}
		return Team{}, fmt.Errorf("update team: %w", err)
	}
	return r.getByID(ctx, id)
}

func (r *repo) deleteByID(ctx context.Context, id string) error {
	res, err := r.db.Writer().ExecContext(ctx, `DELETE FROM teams WHERE id = ?`, id)
	if err != nil {
		return fmt.Errorf("delete team: %w", err)
	}
	n, _ := res.RowsAffected()
	if n == 0 {
		return errTeamNotFound
	}
	return nil
}

func (r *repo) assignKey(ctx context.Context, teamID, keyID string) (Team, error) {
	if _, err := r.getByID(ctx, teamID); err != nil {
		return Team{}, err
	}
	now := formatUTC(time.Now())
	_, err := r.db.Writer().ExecContext(ctx,
		`INSERT INTO team_keys (team_id, key_id, added_at) VALUES (?, ?, ?)
		 ON CONFLICT(key_id) DO UPDATE SET team_id = excluded.team_id, added_at = excluded.added_at`,
		teamID, strings.TrimSpace(keyID), now,
	)
	if err != nil {
		return Team{}, fmt.Errorf("assign key to team: %w", err)
	}
	return r.getByID(ctx, teamID)
}

func (r *repo) removeKey(ctx context.Context, teamID, keyID string) (Team, error) {
	if _, err := r.getByID(ctx, teamID); err != nil {
		return Team{}, err
	}
	res, err := r.db.Writer().ExecContext(ctx,
		`DELETE FROM team_keys WHERE team_id = ? AND key_id = ?`,
		teamID, strings.TrimSpace(keyID),
	)
	if err != nil {
		return Team{}, fmt.Errorf("remove key from team: %w", err)
	}
	n, _ := res.RowsAffected()
	if n == 0 {
		return Team{}, errKeyNotInTeam
	}
	return r.getByID(ctx, teamID)
}

func (r *repo) addSpendForKey(ctx context.Context, keyID string, costUSD float64) (string, error) {
	if costUSD <= 0 || strings.TrimSpace(keyID) == "" {
		return "", nil
	}
	var teamID string
	err := r.db.Reader().QueryRowContext(ctx, `SELECT team_id FROM team_keys WHERE key_id = ?`, keyID).Scan(&teamID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return "", nil
		}
		return "", fmt.Errorf("lookup team for key: %w", err)
	}
	now := formatUTC(time.Now())
	_, err = r.db.Writer().ExecContext(ctx,
		`UPDATE teams SET spent_usd = spent_usd + ?, updated_at = ? WHERE id = ?`,
		costUSD, now, teamID,
	)
	if err != nil {
		return "", fmt.Errorf("increment team spend: %w", err)
	}
	return teamID, nil
}

func (r *repo) resetSpend(ctx context.Context, id string) (Team, error) {
	now := formatUTC(time.Now())
	res, err := r.db.Writer().ExecContext(ctx,
		`UPDATE teams SET spent_usd = 0, updated_at = ? WHERE id = ?`,
		now, id,
	)
	if err != nil {
		return Team{}, fmt.Errorf("reset team spend: %w", err)
	}
	n, _ := res.RowsAffected()
	if n == 0 {
		return Team{}, errTeamNotFound
	}
	return r.getByID(ctx, id)
}
