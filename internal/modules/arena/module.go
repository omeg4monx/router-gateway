// Package arena implements multi-model side-by-side comparison battles,
// leaderboard scoring, and persistent multi-session Basic Chat backed by
// contracts.ChatCompleter.
package arena

import (
	"context"
	"embed"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/fs"
	"log/slog"
	"net/http"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/go-chi/chi/v5"

	"github.com/omeg4monx/router-gateway/internal/contracts"
	"github.com/omeg4monx/router-gateway/internal/platform/db"
	"github.com/omeg4monx/router-gateway/internal/platform/eventbus"
	"github.com/omeg4monx/router-gateway/internal/platform/httpserver"
	"github.com/omeg4monx/router-gateway/internal/platform/lifecycle"
)

//go:embed migrations/*.sql
var migrationsFS embed.FS

const (
	moduleName          = "arena"
	maxJSONBodyBytes    = 1 << 20
	minBattleModels     = 2
	maxBattleModels     = 4
	maxChatHistoryTurns = 50
	perModelTimeout     = 60 * time.Second
	defaultSessionTitle = "New chat"
)

// Options configures an arena Module instance.
type Options struct {
	DB        *db.DB
	Completer contracts.ChatCompleter
	Models    contracts.ModelLister
	Bus       *eventbus.Bus
	Logger    *slog.Logger
}

// ModelOption is the wire representation of an available model in Arena & Chat.
type ModelOption struct {
	ID      string `json:"id"`
	OwnedBy string `json:"owned_by"`
	Kind    string `json:"kind"`
}

// Module implements lifecycle.Module for Model Arena and Basic Chat.
type Module struct {
	repo      *repo
	completer contracts.ChatCompleter
	models    contracts.ModelLister
	bus       *eventbus.Bus
	logger    *slog.Logger
}

var _ lifecycle.Module = (*Module)(nil)

// New constructs an arena Module.
func New(opts Options) *Module {
	logger := opts.Logger
	if logger == nil {
		logger = slog.Default()
	}
	return &Module{
		repo:      newRepo(opts.DB),
		completer: opts.Completer,
		models:    opts.Models,
		bus:       opts.Bus,
		logger:    logger,
	}
}

// Name returns the canonical module name.
func (m *Module) Name() string { return moduleName }

// Dependencies declares modules required before arena initializes.
func (m *Module) Dependencies() []string { return []string{"system", "gateway"} }

// Migrations returns the embedded SQLite migrations for arena and chat.
func (m *Module) Migrations() fs.FS {
	sub, err := fs.Sub(migrationsFS, "migrations")
	if err != nil {
		panic(fmt.Sprintf("arena: load embedded migrations: %v", err))
	}
	return sub
}

// RegisterRoutes mounts Arena and Basic Chat AdminAPI endpoints.
// Long-running inference routes mount on AdminStreamAPI when available so
// per-model timeouts (60s) are not truncated by the default 30s APITimeout.
func (m *Module) RegisterRoutes(mounts lifecycle.RouteMounts) {
	if mounts.AdminAPI == nil {
		return
	}
	mounts.AdminAPI.Get("/arena/models", m.handleListModels)
	mounts.AdminAPI.Get("/arena/battles", m.handleListBattles)
	mounts.AdminAPI.Get("/arena/battles/{id}", m.handleGetBattle)
	mounts.AdminAPI.Post("/arena/battles/{id}/vote", m.handleVoteBattle)
	mounts.AdminAPI.Delete("/arena/battles/{id}", m.handleDeleteBattle)
	mounts.AdminAPI.Get("/arena/leaderboard", m.handleLeaderboard)

	mounts.AdminAPI.Get("/chat/sessions", m.handleListSessions)
	mounts.AdminAPI.Post("/chat/sessions", m.handleCreateSession)
	mounts.AdminAPI.Get("/chat/sessions/{id}", m.handleGetSession)
	mounts.AdminAPI.Patch("/chat/sessions/{id}", m.handleUpdateSession)
	mounts.AdminAPI.Delete("/chat/sessions/{id}", m.handleDeleteSession)

	inferenceRouter := mounts.AdminStreamAPI
	if inferenceRouter == nil {
		inferenceRouter = mounts.AdminAPI
	}
	inferenceRouter.Post("/arena/battles", m.handleRunBattle)
	inferenceRouter.Post("/chat/sessions/{id}/messages", m.handleSendMessage)
}

// Start validates module wiring.
func (m *Module) Start(context.Context) error {
	if m.repo == nil || m.repo.db == nil {
		return errors.New("arena: database is required")
	}
	return nil
}

// Stop performs graceful teardown.
func (m *Module) Stop(context.Context) error {
	return nil
}

func (m *Module) handleListModels(w http.ResponseWriter, r *http.Request) {
	if m.models == nil {
		httpserver.WriteJSON(w, http.StatusOK, map[string]any{"models": []ModelOption{}})
		return
	}
	list, err := m.models.ListModels(r.Context(), nil)
	if err != nil {
		httpserver.WriteError(w, r, http.StatusInternalServerError, "arena_models_failed", "failed to list models")
		return
	}
	items := make([]ModelOption, 0, len(list))
	for _, d := range list {
		items = append(items, ModelOption{
			ID:      d.ID,
			OwnedBy: d.OwnedBy,
			Kind:    d.Kind,
		})
	}
	httpserver.WriteJSON(w, http.StatusOK, map[string]any{"models": items})
}

func (m *Module) handleListBattles(w http.ResponseWriter, r *http.Request) {
	limit := 30
	if raw := strings.TrimSpace(r.URL.Query().Get("limit")); raw != "" {
		if parsed, err := strconv.Atoi(raw); err == nil && parsed > 0 {
			limit = parsed
		}
	}
	battles, err := m.repo.listBattles(r.Context(), limit)
	if err != nil {
		httpserver.WriteError(w, r, http.StatusInternalServerError, "arena_list_failed", "failed to list arena battles")
		return
	}
	httpserver.WriteJSON(w, http.StatusOK, map[string]any{"battles": battles})
}

func (m *Module) handleRunBattle(w http.ResponseWriter, r *http.Request) {
	if m.completer == nil {
		httpserver.WriteError(w, r, http.StatusServiceUnavailable, "completer_unavailable", "gateway chat completer is not wired")
		return
	}
	var body struct {
		Prompt       string   `json:"prompt"`
		SystemPrompt string   `json:"system_prompt"`
		Models       []string `json:"models"`
		Temperature  *float64 `json:"temperature"`
		MaxTokens    int      `json:"max_tokens"`
	}
	if err := decodeJSON(r, &body); err != nil {
		httpserver.WriteError(w, r, http.StatusBadRequest, "invalid_json", err.Error())
		return
	}
	prompt := strings.TrimSpace(body.Prompt)
	if prompt == "" {
		httpserver.WriteError(w, r, http.StatusBadRequest, "prompt_required", "prompt is required")
		return
	}

	cleanedModels := make([]string, 0, len(body.Models))
	seen := make(map[string]struct{}, len(body.Models))
	for _, rawModel := range body.Models {
		model := strings.TrimSpace(rawModel)
		if model == "" {
			continue
		}
		if _, dup := seen[model]; dup {
			continue
		}
		seen[model] = struct{}{}
		cleanedModels = append(cleanedModels, model)
	}
	if len(cleanedModels) < minBattleModels || len(cleanedModels) > maxBattleModels {
		httpserver.WriteError(w, r, http.StatusBadRequest, "invalid_models", "select between 2 and 4 distinct models")
		return
	}
	if body.Temperature != nil && (*body.Temperature < 0 || *body.Temperature > 2) {
		httpserver.WriteError(w, r, http.StatusBadRequest, "invalid_temperature", "temperature must be between 0 and 2")
		return
	}
	if body.MaxTokens < 0 || body.MaxTokens > 128000 {
		httpserver.WriteError(w, r, http.StatusBadRequest, "invalid_max_tokens", "max_tokens must be between 0 and 128000")
		return
	}

	messages := make([]contracts.ChatMessage, 0, 2)
	sysPrompt := strings.TrimSpace(body.SystemPrompt)
	if sysPrompt != "" {
		messages = append(messages, contracts.ChatMessage{Role: "system", Content: sysPrompt})
	}
	messages = append(messages, contracts.ChatMessage{Role: "user", Content: prompt})

	battleID := randomHexID("battle")
	contenders := make([]ContenderResult, len(cleanedModels))

	var wg sync.WaitGroup
	for idx, modelName := range cleanedModels {
		wg.Add(1)
		go func(slot int, model string) {
			defer wg.Done()
			callCtx, cancel := context.WithTimeout(r.Context(), perModelTimeout)
			defer cancel()

			start := time.Now()
			res, err := m.completer.CompleteChat(callCtx, contracts.ChatCompletionRequest{
				Model:       model,
				Messages:    messages,
				Temperature: body.Temperature,
				MaxTokens:   body.MaxTokens,
			})
			latencyMs := res.LatencyMs
			if latencyMs <= 0 {
				latencyMs = time.Since(start).Milliseconds()
			}
			item := ContenderResult{
				ID:               randomHexID("cont"),
				BattleID:         battleID,
				SlotIndex:        slot,
				Model:            model,
				Content:          res.Content,
				Status:           "ok",
				LatencyMs:        latencyMs,
				PromptTokens:     res.PromptTokens,
				CompletionTokens: res.CompletionTokens,
				TotalTokens:      res.TotalTokens,
			}
			if err != nil {
				item.Status = "error"
				item.Error = err.Error()
			}
			contenders[slot] = item
		}(idx, modelName)
	}
	wg.Wait()

	battle := Battle{
		ID:           battleID,
		Prompt:       prompt,
		SystemPrompt: sysPrompt,
		Temperature:  body.Temperature,
		MaxTokens:    body.MaxTokens,
		CreatedAt:    nowTimestamp(),
		Contenders:   contenders,
	}
	saveCtx, saveCancel := context.WithTimeout(context.WithoutCancel(r.Context()), 5*time.Second)
	defer saveCancel()
	if err := m.repo.createBattle(saveCtx, battle); err != nil {
		httpserver.WriteError(w, r, http.StatusInternalServerError, "arena_save_failed", "failed to save battle")
		return
	}
	httpserver.WriteJSON(w, http.StatusCreated, battle)
}

func (m *Module) handleGetBattle(w http.ResponseWriter, r *http.Request) {
	id := strings.TrimSpace(chi.URLParam(r, "id"))
	battle, err := m.repo.getBattle(r.Context(), id)
	if errors.Is(err, errBattleNotFound) {
		httpserver.WriteError(w, r, http.StatusNotFound, "battle_not_found", "arena battle not found")
		return
	}
	if err != nil {
		httpserver.WriteError(w, r, http.StatusInternalServerError, "arena_read_failed", "failed to load arena battle")
		return
	}
	httpserver.WriteJSON(w, http.StatusOK, battle)
}

func (m *Module) handleVoteBattle(w http.ResponseWriter, r *http.Request) {
	id := strings.TrimSpace(chi.URLParam(r, "id"))
	var body struct {
		WinnerModel string `json:"winner_model"`
	}
	if err := decodeJSON(r, &body); err != nil {
		httpserver.WriteError(w, r, http.StatusBadRequest, "invalid_json", err.Error())
		return
	}
	winner := strings.TrimSpace(body.WinnerModel)
	if winner == "" {
		httpserver.WriteError(w, r, http.StatusBadRequest, "winner_required", "winner_model is required (model ID or 'tie')")
		return
	}

	battle, err := m.repo.getBattle(r.Context(), id)
	if errors.Is(err, errBattleNotFound) {
		httpserver.WriteError(w, r, http.StatusNotFound, "battle_not_found", "arena battle not found")
		return
	}
	if err != nil {
		httpserver.WriteError(w, r, http.StatusInternalServerError, "arena_read_failed", "failed to load arena battle")
		return
	}

	if winner != "tie" {
		matched := false
		for _, c := range battle.Contenders {
			if c.Model == winner {
				matched = true
				break
			}
		}
		if !matched {
			httpserver.WriteError(w, r, http.StatusBadRequest, "invalid_winner", "winner_model must match a participating model or 'tie'")
			return
		}
	}

	if err := m.repo.setBattleWinner(r.Context(), id, winner); err != nil {
		httpserver.WriteError(w, r, http.StatusInternalServerError, "arena_vote_failed", "failed to record battle vote")
		return
	}
	battle.WinnerModel = winner
	httpserver.WriteJSON(w, http.StatusOK, battle)
}

func (m *Module) handleDeleteBattle(w http.ResponseWriter, r *http.Request) {
	id := strings.TrimSpace(chi.URLParam(r, "id"))
	err := m.repo.deleteBattle(r.Context(), id)
	if errors.Is(err, errBattleNotFound) {
		httpserver.WriteError(w, r, http.StatusNotFound, "battle_not_found", "arena battle not found")
		return
	}
	if err != nil {
		httpserver.WriteError(w, r, http.StatusInternalServerError, "arena_delete_failed", "failed to delete arena battle")
		return
	}
	httpserver.WriteJSON(w, http.StatusOK, map[string]any{"deleted": true, "id": id})
}

func (m *Module) handleLeaderboard(w http.ResponseWriter, r *http.Request) {
	entries, err := m.repo.leaderboard(r.Context())
	if err != nil {
		httpserver.WriteError(w, r, http.StatusInternalServerError, "arena_leaderboard_failed", "failed to compute leaderboard")
		return
	}
	httpserver.WriteJSON(w, http.StatusOK, map[string]any{"entries": entries})
}

func (m *Module) handleListSessions(w http.ResponseWriter, r *http.Request) {
	sessions, err := m.repo.listSessions(r.Context(), 50)
	if err != nil {
		httpserver.WriteError(w, r, http.StatusInternalServerError, "chat_sessions_list_failed", "failed to list chat sessions")
		return
	}
	httpserver.WriteJSON(w, http.StatusOK, map[string]any{"sessions": sessions})
}

func (m *Module) handleCreateSession(w http.ResponseWriter, r *http.Request) {
	var body struct {
		Title        string   `json:"title"`
		Model        string   `json:"model"`
		SystemPrompt string   `json:"system_prompt"`
		Temperature  *float64 `json:"temperature"`
		MaxTokens    int      `json:"max_tokens"`
	}
	if err := decodeJSON(r, &body); err != nil {
		httpserver.WriteError(w, r, http.StatusBadRequest, "invalid_json", err.Error())
		return
	}
	model := strings.TrimSpace(body.Model)
	if model == "" {
		httpserver.WriteError(w, r, http.StatusBadRequest, "model_required", "model is required")
		return
	}
	if body.Temperature != nil && (*body.Temperature < 0 || *body.Temperature > 2) {
		httpserver.WriteError(w, r, http.StatusBadRequest, "invalid_temperature", "temperature must be between 0 and 2")
		return
	}
	if body.MaxTokens < 0 || body.MaxTokens > 128000 {
		httpserver.WriteError(w, r, http.StatusBadRequest, "invalid_max_tokens", "max_tokens must be between 0 and 128000")
		return
	}

	title := strings.TrimSpace(body.Title)
	if title == "" {
		title = defaultSessionTitle
	}
	now := nowTimestamp()
	session := ChatSession{
		ID:           randomHexID("chat"),
		Title:        title,
		Model:        model,
		SystemPrompt: strings.TrimSpace(body.SystemPrompt),
		Temperature:  body.Temperature,
		MaxTokens:    body.MaxTokens,
		CreatedAt:    now,
		UpdatedAt:    now,
	}
	if err := m.repo.createSession(r.Context(), session); err != nil {
		httpserver.WriteError(w, r, http.StatusInternalServerError, "chat_session_create_failed", "failed to create chat session")
		return
	}
	httpserver.WriteJSON(w, http.StatusCreated, session)
}

func (m *Module) handleGetSession(w http.ResponseWriter, r *http.Request) {
	id := strings.TrimSpace(chi.URLParam(r, "id"))
	session, err := m.repo.getSession(r.Context(), id)
	if errors.Is(err, errSessionNotFound) {
		httpserver.WriteError(w, r, http.StatusNotFound, "session_not_found", "chat session not found")
		return
	}
	if err != nil {
		httpserver.WriteError(w, r, http.StatusInternalServerError, "chat_session_read_failed", "failed to load chat session")
		return
	}
	messages, err := m.repo.listMessages(r.Context(), id)
	if err != nil {
		httpserver.WriteError(w, r, http.StatusInternalServerError, "chat_messages_read_failed", "failed to load chat messages")
		return
	}
	httpserver.WriteJSON(w, http.StatusOK, map[string]any{
		"session":  session,
		"messages": messages,
	})
}

func (m *Module) handleUpdateSession(w http.ResponseWriter, r *http.Request) {
	id := strings.TrimSpace(chi.URLParam(r, "id"))
	var body struct {
		Title        *string  `json:"title"`
		Model        *string  `json:"model"`
		SystemPrompt *string  `json:"system_prompt"`
		Temperature  *float64 `json:"temperature"`
		MaxTokens    *int     `json:"max_tokens"`
	}
	if err := decodeJSON(r, &body); err != nil {
		httpserver.WriteError(w, r, http.StatusBadRequest, "invalid_json", err.Error())
		return
	}

	session, err := m.repo.getSession(r.Context(), id)
	if errors.Is(err, errSessionNotFound) {
		httpserver.WriteError(w, r, http.StatusNotFound, "session_not_found", "chat session not found")
		return
	}
	if err != nil {
		httpserver.WriteError(w, r, http.StatusInternalServerError, "chat_session_read_failed", "failed to load chat session")
		return
	}

	if body.Title != nil && strings.TrimSpace(*body.Title) != "" {
		session.Title = strings.TrimSpace(*body.Title)
	}
	if body.Model != nil && strings.TrimSpace(*body.Model) != "" {
		session.Model = strings.TrimSpace(*body.Model)
	}
	if body.SystemPrompt != nil {
		session.SystemPrompt = strings.TrimSpace(*body.SystemPrompt)
	}
	if body.Temperature != nil {
		if *body.Temperature < 0 || *body.Temperature > 2 {
			httpserver.WriteError(w, r, http.StatusBadRequest, "invalid_temperature", "temperature must be between 0 and 2")
			return
		}
		session.Temperature = body.Temperature
	}
	if body.MaxTokens != nil {
		if *body.MaxTokens < 0 || *body.MaxTokens > 128000 {
			httpserver.WriteError(w, r, http.StatusBadRequest, "invalid_max_tokens", "max_tokens must be between 0 and 128000")
			return
		}
		session.MaxTokens = *body.MaxTokens
	}
	session.UpdatedAt = nowTimestamp()

	if err := m.repo.updateSession(r.Context(), session); err != nil {
		httpserver.WriteError(w, r, http.StatusInternalServerError, "chat_session_update_failed", "failed to update chat session")
		return
	}
	httpserver.WriteJSON(w, http.StatusOK, session)
}

func (m *Module) handleDeleteSession(w http.ResponseWriter, r *http.Request) {
	id := strings.TrimSpace(chi.URLParam(r, "id"))
	err := m.repo.deleteSession(r.Context(), id)
	if errors.Is(err, errSessionNotFound) {
		httpserver.WriteError(w, r, http.StatusNotFound, "session_not_found", "chat session not found")
		return
	}
	if err != nil {
		httpserver.WriteError(w, r, http.StatusInternalServerError, "chat_session_delete_failed", "failed to delete chat session")
		return
	}
	httpserver.WriteJSON(w, http.StatusOK, map[string]any{"deleted": true, "id": id})
}

func (m *Module) handleSendMessage(w http.ResponseWriter, r *http.Request) {
	if m.completer == nil {
		httpserver.WriteError(w, r, http.StatusServiceUnavailable, "completer_unavailable", "gateway chat completer is not wired")
		return
	}
	id := strings.TrimSpace(chi.URLParam(r, "id"))
	var body struct {
		Content string `json:"content"`
		Model   string `json:"model"`
	}
	if err := decodeJSON(r, &body); err != nil {
		httpserver.WriteError(w, r, http.StatusBadRequest, "invalid_json", err.Error())
		return
	}
	content := strings.TrimSpace(body.Content)
	if content == "" {
		httpserver.WriteError(w, r, http.StatusBadRequest, "content_required", "message content is required")
		return
	}

	session, err := m.repo.getSession(r.Context(), id)
	if errors.Is(err, errSessionNotFound) {
		httpserver.WriteError(w, r, http.StatusNotFound, "session_not_found", "chat session not found")
		return
	}
	if err != nil {
		httpserver.WriteError(w, r, http.StatusInternalServerError, "chat_session_read_failed", "failed to load chat session")
		return
	}

	model := session.Model
	if override := strings.TrimSpace(body.Model); override != "" {
		model = override
		session.Model = override
	}

	history, err := m.repo.listMessages(r.Context(), id)
	if err != nil {
		httpserver.WriteError(w, r, http.StatusInternalServerError, "chat_messages_read_failed", "failed to load chat history")
		return
	}
	if len(history) > maxChatHistoryTurns {
		history = history[len(history)-maxChatHistoryTurns:]
	}

	reqMessages := make([]contracts.ChatMessage, 0, len(history)+2)
	if session.SystemPrompt != "" {
		reqMessages = append(reqMessages, contracts.ChatMessage{Role: "system", Content: session.SystemPrompt})
	}
	for _, h := range history {
		reqMessages = append(reqMessages, contracts.ChatMessage{Role: h.Role, Content: h.Content})
	}
	reqMessages = append(reqMessages, contracts.ChatMessage{Role: "user", Content: content})

	callCtx, cancel := context.WithTimeout(r.Context(), perModelTimeout)
	defer cancel()

	userTime := time.Now().UTC()
	userTS := userTime.Format(timestampFormat)
	start := time.Now()
	res, err := m.completer.CompleteChat(callCtx, contracts.ChatCompletionRequest{
		Model:       model,
		Messages:    reqMessages,
		Temperature: session.Temperature,
		MaxTokens:   session.MaxTokens,
	})
	if err != nil {
		status := http.StatusBadGateway
		code := "chat_completion_failed"
		if errors.Is(err, contracts.ErrModelNotFound) {
			status = http.StatusBadRequest
			code = "model_not_found"
		}
		httpserver.WriteError(w, r, status, code, err.Error())
		return
	}

	latencyMs := res.LatencyMs
	if latencyMs <= 0 {
		latencyMs = time.Since(start).Milliseconds()
	}
	assistantTime := time.Now().UTC()
	if !assistantTime.After(userTime) {
		assistantTime = userTime.Add(time.Microsecond)
	}
	assistantTS := assistantTime.Format(timestampFormat)

	userMsg := ChatMessage{
		ID:        randomHexID("msg"),
		SessionID: session.ID,
		Role:      "user",
		Content:   content,
		Model:     model,
		CreatedAt: userTS,
	}
	assistantMsg := ChatMessage{
		ID:               randomHexID("msg"),
		SessionID:        session.ID,
		Role:             "assistant",
		Content:          res.Content,
		Model:            res.Model,
		LatencyMs:        latencyMs,
		PromptTokens:     res.PromptTokens,
		CompletionTokens: res.CompletionTokens,
		TotalTokens:      res.TotalTokens,
		CreatedAt:        assistantTS,
	}

	if session.Title == "" || session.Title == defaultSessionTitle {
		session.Title = deriveSessionTitle(content)
	}
	session.UpdatedAt = assistantTS
	session.MessageCount += 2

	if err := m.repo.appendTurn(r.Context(), session, userMsg, assistantMsg); err != nil {
		httpserver.WriteError(w, r, http.StatusInternalServerError, "chat_turn_save_failed", "failed to save chat turn")
		return
	}

	httpserver.WriteJSON(w, http.StatusCreated, map[string]any{
		"session":           session,
		"user_message":      userMsg,
		"assistant_message": assistantMsg,
	})
}

func deriveSessionTitle(prompt string) string {
	collapsed := strings.Join(strings.Fields(prompt), " ")
	if collapsed == "" {
		return defaultSessionTitle
	}
	runes := []rune(collapsed)
	if len(runes) > 52 {
		return strings.TrimSpace(string(runes[:52])) + "…"
	}
	return collapsed
}

func decodeJSON(r *http.Request, dst any) error {
	raw, err := io.ReadAll(io.LimitReader(r.Body, maxJSONBodyBytes+1))
	if err != nil {
		return errors.New("failed to read request body")
	}
	if len(raw) > maxJSONBodyBytes {
		return errors.New("request body exceeds limit")
	}
	if strings.TrimSpace(string(raw)) == "" {
		return errors.New("request body is empty")
	}
	if err := json.Unmarshal(raw, dst); err != nil {
		return fmt.Errorf("invalid JSON body: %w", err)
	}
	return nil
}
