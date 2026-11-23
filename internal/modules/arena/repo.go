package arena

import (
	"context"
	"crypto/rand"
	"database/sql"
	"encoding/hex"
	"errors"
	"fmt"
	"math"
	"slices"
	"strings"
	"time"

	"github.com/omeg4monx/router-gateway/internal/platform/db"
)

const (
	timestampFormat = "2006-01-02T15:04:05.000000000Z"
	maxListLimit    = 100
)

var (
	errBattleNotFound  = errors.New("arena battle not found")
	errSessionNotFound = errors.New("chat session not found")
)

// ContenderResult stores one model's output and metrics within an arena battle.
type ContenderResult struct {
	ID               string `json:"id"`
	BattleID         string `json:"battle_id"`
	SlotIndex        int    `json:"slot_index"`
	Model            string `json:"model"`
	Content          string `json:"content"`
	Status           string `json:"status"`
	Error            string `json:"error"`
	LatencyMs        int64  `json:"latency_ms"`
	PromptTokens     int64  `json:"prompt_tokens"`
	CompletionTokens int64  `json:"completion_tokens"`
	TotalTokens      int64  `json:"total_tokens"`
}

// Battle represents a multi-model side-by-side comparison run.
type Battle struct {
	ID           string            `json:"id"`
	Prompt       string            `json:"prompt"`
	SystemPrompt string            `json:"system_prompt"`
	Temperature  *float64          `json:"temperature,omitempty"`
	MaxTokens    int               `json:"max_tokens"`
	WinnerModel  string            `json:"winner_model"`
	CreatedAt    string            `json:"created_at"`
	Contenders   []ContenderResult `json:"contenders"`
}

// LeaderboardEntry summarizes battle outcomes and latency for a model.
type LeaderboardEntry struct {
	Model        string  `json:"model"`
	Battles      int     `json:"battles"`
	Wins         int     `json:"wins"`
	Ties         int     `json:"ties"`
	Losses       int     `json:"losses"`
	WinRate      float64 `json:"win_rate"`
	AvgLatencyMs int64   `json:"avg_latency_ms"`
}

// ChatSession represents a persistent Basic Chat conversation.
type ChatSession struct {
	ID           string   `json:"id"`
	Title        string   `json:"title"`
	Model        string   `json:"model"`
	SystemPrompt string   `json:"system_prompt"`
	Temperature  *float64 `json:"temperature,omitempty"`
	MaxTokens    int      `json:"max_tokens"`
	MessageCount int      `json:"message_count"`
	CreatedAt    string   `json:"created_at"`
	UpdatedAt    string   `json:"updated_at"`
}

// ChatMessage represents a single persisted turn inside a Basic Chat session.
type ChatMessage struct {
	ID               string `json:"id"`
	SessionID        string `json:"session_id"`
	Role             string `json:"role"`
	Content          string `json:"content"`
	Model            string `json:"model"`
	LatencyMs        int64  `json:"latency_ms"`
	PromptTokens     int64  `json:"prompt_tokens"`
	CompletionTokens int64  `json:"completion_tokens"`
	TotalTokens      int64  `json:"total_tokens"`
	CreatedAt        string `json:"created_at"`
}

type repo struct {
	db *db.DB
}

func newRepo(database *db.DB) *repo {
	return &repo{db: database}
}

func nowTimestamp() string {
	return time.Now().UTC().Format(timestampFormat)
}

func randomHexID(prefix string) string {
	var b [8]byte
	_, _ = rand.Read(b[:])
	return prefix + "_" + hex.EncodeToString(b[:])
}

func (r *repo) createBattle(ctx context.Context, battle Battle) error {
	tx, err := r.db.Writer().BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("begin battle tx: %w", err)
	}
	defer func() { _ = tx.Rollback() }()

	var tempVal any
	if battle.Temperature != nil {
		tempVal = *battle.Temperature
	}
	if _, err := tx.ExecContext(ctx, `
		INSERT INTO arena_battles (id, prompt, system_prompt, temperature, max_tokens, winner_model, created_at)
		VALUES (?, ?, ?, ?, ?, ?, ?)
	`, battle.ID, battle.Prompt, battle.SystemPrompt, tempVal, battle.MaxTokens, battle.WinnerModel, battle.CreatedAt); err != nil {
		return fmt.Errorf("insert arena battle: %w", err)
	}

	for _, c := range battle.Contenders {
		if _, err := tx.ExecContext(ctx, `
			INSERT INTO arena_contenders (
				id, battle_id, slot_index, model, content, status, error,
				latency_ms, prompt_tokens, completion_tokens, total_tokens
			) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
		`, c.ID, battle.ID, c.SlotIndex, c.Model, c.Content, c.Status, c.Error,
			c.LatencyMs, c.PromptTokens, c.CompletionTokens, c.TotalTokens); err != nil {
			return fmt.Errorf("insert arena contender: %w", err)
		}
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("commit battle tx: %w", err)
	}
	return nil
}

func (r *repo) getBattle(ctx context.Context, id string) (Battle, error) {
	var (
		b    Battle
		temp sql.NullFloat64
	)
	err := r.db.Reader().QueryRowContext(ctx, `
		SELECT id, prompt, system_prompt, temperature, max_tokens, winner_model, created_at
		FROM arena_battles
		WHERE id = ?
	`, id).Scan(&b.ID, &b.Prompt, &b.SystemPrompt, &temp, &b.MaxTokens, &b.WinnerModel, &b.CreatedAt)
	if errors.Is(err, sql.ErrNoRows) {
		return Battle{}, errBattleNotFound
	}
	if err != nil {
		return Battle{}, fmt.Errorf("query arena battle: %w", err)
	}
	if temp.Valid {
		v := temp.Float64
		b.Temperature = &v
	}

	contenders, err := r.listContenders(ctx, b.ID)
	if err != nil {
		return Battle{}, err
	}
	b.Contenders = contenders
	return b, nil
}

func (r *repo) listContenders(ctx context.Context, battleID string) ([]ContenderResult, error) {
	rows, err := r.db.Reader().QueryContext(ctx, `
		SELECT id, battle_id, slot_index, model, content, status, error,
		       latency_ms, prompt_tokens, completion_tokens, total_tokens
		FROM arena_contenders
		WHERE battle_id = ?
		ORDER BY slot_index ASC
	`, battleID)
	if err != nil {
		return nil, fmt.Errorf("query arena contenders: %w", err)
	}
	defer rows.Close()

	out := make([]ContenderResult, 0, 4)
	for rows.Next() {
		var c ContenderResult
		if err := rows.Scan(
			&c.ID, &c.BattleID, &c.SlotIndex, &c.Model, &c.Content, &c.Status, &c.Error,
			&c.LatencyMs, &c.PromptTokens, &c.CompletionTokens, &c.TotalTokens,
		); err != nil {
			return nil, fmt.Errorf("scan arena contender: %w", err)
		}
		out = append(out, c)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate arena contenders: %w", err)
	}
	return out, nil
}

func (r *repo) listBattles(ctx context.Context, limit int) ([]Battle, error) {
	if limit <= 0 || limit > maxListLimit {
		limit = 30
	}
	rows, err := r.db.Reader().QueryContext(ctx, `
		SELECT id, prompt, system_prompt, temperature, max_tokens, winner_model, created_at
		FROM arena_battles
		ORDER BY created_at DESC
		LIMIT ?
	`, limit)
	if err != nil {
		return nil, fmt.Errorf("list arena battles: %w", err)
	}
	defer rows.Close()

	battles := make([]Battle, 0, limit)
	for rows.Next() {
		var (
			b    Battle
			temp sql.NullFloat64
		)
		if err := rows.Scan(&b.ID, &b.Prompt, &b.SystemPrompt, &temp, &b.MaxTokens, &b.WinnerModel, &b.CreatedAt); err != nil {
			return nil, fmt.Errorf("scan arena battle: %w", err)
		}
		if temp.Valid {
			v := temp.Float64
			b.Temperature = &v
		}
		battles = append(battles, b)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate arena battles: %w", err)
	}

	for i := range battles {
		contenders, err := r.listContenders(ctx, battles[i].ID)
		if err != nil {
			return nil, err
		}
		battles[i].Contenders = contenders
	}
	return battles, nil
}

func (r *repo) setBattleWinner(ctx context.Context, id, winner string) error {
	res, err := r.db.Writer().ExecContext(ctx, `
		UPDATE arena_battles SET winner_model = ? WHERE id = ?
	`, winner, id)
	if err != nil {
		return fmt.Errorf("update arena winner: %w", err)
	}
	n, _ := res.RowsAffected()
	if n == 0 {
		return errBattleNotFound
	}
	return nil
}

func (r *repo) deleteBattle(ctx context.Context, id string) error {
	res, err := r.db.Writer().ExecContext(ctx, `DELETE FROM arena_battles WHERE id = ?`, id)
	if err != nil {
		return fmt.Errorf("delete arena battle: %w", err)
	}
	n, _ := res.RowsAffected()
	if n == 0 {
		return errBattleNotFound
	}
	return nil
}

func (r *repo) leaderboard(ctx context.Context) ([]LeaderboardEntry, error) {
	rows, err := r.db.Reader().QueryContext(ctx, `
		SELECT c.model, c.latency_ms, b.winner_model
		FROM arena_contenders c
		JOIN arena_battles b ON b.id = c.battle_id
	`)
	if err != nil {
		return nil, fmt.Errorf("query leaderboard: %w", err)
	}
	defer rows.Close()

	type acc struct {
		battles      int
		wins         int
		ties         int
		losses       int
		totalLatency int64
	}
	stats := make(map[string]*acc)
	for rows.Next() {
		var (
			model   string
			latency int64
			winner  string
		)
		if err := rows.Scan(&model, &latency, &winner); err != nil {
			return nil, fmt.Errorf("scan leaderboard row: %w", err)
		}
		entry := stats[model]
		if entry == nil {
			entry = &acc{}
			stats[model] = entry
		}
		entry.battles++
		entry.totalLatency += latency
		switch winner {
		case "":
			// Unvoted battle: still counted in total battles and avg latency.
		case "tie":
			entry.ties++
		case model:
			entry.wins++
		default:
			entry.losses++
		}
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate leaderboard rows: %w", err)
	}

	entries := make([]LeaderboardEntry, 0, len(stats))
	for model, a := range stats {
		decided := a.wins + a.losses + a.ties
		winRate := 0.0
		if decided > 0 {
			winRate = math.Round((float64(a.wins)/float64(decided))*1000) / 10
		}
		avgLatency := int64(0)
		if a.battles > 0 {
			avgLatency = a.totalLatency / int64(a.battles)
		}
		entries = append(entries, LeaderboardEntry{
			Model:        model,
			Battles:      a.battles,
			Wins:         a.wins,
			Ties:         a.ties,
			Losses:       a.losses,
			WinRate:      winRate,
			AvgLatencyMs: avgLatency,
		})
	}

	slices.SortFunc(entries, func(a, b LeaderboardEntry) int {
		if a.WinRate != b.WinRate {
			if a.WinRate > b.WinRate {
				return -1
			}
			return 1
		}
		if a.Wins != b.Wins {
			if a.Wins > b.Wins {
				return -1
			}
			return 1
		}
		if a.AvgLatencyMs != b.AvgLatencyMs {
			if a.AvgLatencyMs < b.AvgLatencyMs {
				return -1
			}
			return 1
		}
		return strings.Compare(a.Model, b.Model)
	})
	return entries, nil
}

func (r *repo) createSession(ctx context.Context, s ChatSession) error {
	var tempVal any
	if s.Temperature != nil {
		tempVal = *s.Temperature
	}
	_, err := r.db.Writer().ExecContext(ctx, `
		INSERT INTO chat_sessions (id, title, model, system_prompt, temperature, max_tokens, created_at, updated_at)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?)
	`, s.ID, s.Title, s.Model, s.SystemPrompt, tempVal, s.MaxTokens, s.CreatedAt, s.UpdatedAt)
	if err != nil {
		return fmt.Errorf("insert chat session: %w", err)
	}
	return nil
}

func (r *repo) listSessions(ctx context.Context, limit int) ([]ChatSession, error) {
	if limit <= 0 || limit > maxListLimit {
		limit = 50
	}
	rows, err := r.db.Reader().QueryContext(ctx, `
		SELECT s.id, s.title, s.model, s.system_prompt, s.temperature, s.max_tokens,
		       (SELECT COUNT(*) FROM chat_messages m WHERE m.session_id = s.id) AS message_count,
		       s.created_at, s.updated_at
		FROM chat_sessions s
		ORDER BY s.updated_at DESC
		LIMIT ?
	`, limit)
	if err != nil {
		return nil, fmt.Errorf("list chat sessions: %w", err)
	}
	defer rows.Close()

	out := make([]ChatSession, 0, limit)
	for rows.Next() {
		var (
			s    ChatSession
			temp sql.NullFloat64
		)
		if err := rows.Scan(
			&s.ID, &s.Title, &s.Model, &s.SystemPrompt, &temp, &s.MaxTokens,
			&s.MessageCount, &s.CreatedAt, &s.UpdatedAt,
		); err != nil {
			return nil, fmt.Errorf("scan chat session: %w", err)
		}
		if temp.Valid {
			v := temp.Float64
			s.Temperature = &v
		}
		out = append(out, s)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate chat sessions: %w", err)
	}
	return out, nil
}

func (r *repo) getSession(ctx context.Context, id string) (ChatSession, error) {
	var (
		s    ChatSession
		temp sql.NullFloat64
	)
	err := r.db.Reader().QueryRowContext(ctx, `
		SELECT s.id, s.title, s.model, s.system_prompt, s.temperature, s.max_tokens,
		       (SELECT COUNT(*) FROM chat_messages m WHERE m.session_id = s.id) AS message_count,
		       s.created_at, s.updated_at
		FROM chat_sessions s
		WHERE s.id = ?
	`, id).Scan(
		&s.ID, &s.Title, &s.Model, &s.SystemPrompt, &temp, &s.MaxTokens,
		&s.MessageCount, &s.CreatedAt, &s.UpdatedAt,
	)
	if errors.Is(err, sql.ErrNoRows) {
		return ChatSession{}, errSessionNotFound
	}
	if err != nil {
		return ChatSession{}, fmt.Errorf("query chat session: %w", err)
	}
	if temp.Valid {
		v := temp.Float64
		s.Temperature = &v
	}
	return s, nil
}

func (r *repo) updateSession(ctx context.Context, s ChatSession) error {
	var tempVal any
	if s.Temperature != nil {
		tempVal = *s.Temperature
	}
	res, err := r.db.Writer().ExecContext(ctx, `
		UPDATE chat_sessions
		SET title = ?, model = ?, system_prompt = ?, temperature = ?, max_tokens = ?, updated_at = ?
		WHERE id = ?
	`, s.Title, s.Model, s.SystemPrompt, tempVal, s.MaxTokens, s.UpdatedAt, s.ID)
	if err != nil {
		return fmt.Errorf("update chat session: %w", err)
	}
	n, _ := res.RowsAffected()
	if n == 0 {
		return errSessionNotFound
	}
	return nil
}

func (r *repo) deleteSession(ctx context.Context, id string) error {
	res, err := r.db.Writer().ExecContext(ctx, `DELETE FROM chat_sessions WHERE id = ?`, id)
	if err != nil {
		return fmt.Errorf("delete chat session: %w", err)
	}
	n, _ := res.RowsAffected()
	if n == 0 {
		return errSessionNotFound
	}
	return nil
}

func (r *repo) listMessages(ctx context.Context, sessionID string) ([]ChatMessage, error) {
	rows, err := r.db.Reader().QueryContext(ctx, `
		SELECT id, session_id, role, content, model, latency_ms,
		       prompt_tokens, completion_tokens, total_tokens, created_at
		FROM chat_messages
		WHERE session_id = ?
		ORDER BY created_at ASC
	`, sessionID)
	if err != nil {
		return nil, fmt.Errorf("list chat messages: %w", err)
	}
	defer rows.Close()

	out := make([]ChatMessage, 0, 16)
	for rows.Next() {
		var m ChatMessage
		if err := rows.Scan(
			&m.ID, &m.SessionID, &m.Role, &m.Content, &m.Model, &m.LatencyMs,
			&m.PromptTokens, &m.CompletionTokens, &m.TotalTokens, &m.CreatedAt,
		); err != nil {
			return nil, fmt.Errorf("scan chat message: %w", err)
		}
		out = append(out, m)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate chat messages: %w", err)
	}
	return out, nil
}

func (r *repo) appendTurn(ctx context.Context, session ChatSession, userMsg, assistantMsg ChatMessage) error {
	tx, err := r.db.Writer().BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("begin chat turn tx: %w", err)
	}
	defer func() { _ = tx.Rollback() }()

	for _, msg := range []ChatMessage{userMsg, assistantMsg} {
		if _, err := tx.ExecContext(ctx, `
			INSERT INTO chat_messages (
				id, session_id, role, content, model, latency_ms,
				prompt_tokens, completion_tokens, total_tokens, created_at
			) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
		`, msg.ID, msg.SessionID, msg.Role, msg.Content, msg.Model, msg.LatencyMs,
			msg.PromptTokens, msg.CompletionTokens, msg.TotalTokens, msg.CreatedAt); err != nil {
			return fmt.Errorf("insert chat message: %w", err)
		}
	}

	if _, err := tx.ExecContext(ctx, `
		UPDATE chat_sessions SET title = ?, model = ?, updated_at = ? WHERE id = ?
	`, session.Title, session.Model, session.UpdatedAt, session.ID); err != nil {
		return fmt.Errorf("touch chat session: %w", err)
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("commit chat turn tx: %w", err)
	}
	return nil
}
