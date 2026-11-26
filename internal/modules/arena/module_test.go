package arena_test

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"io"
	"io/fs"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"path/filepath"
	"sync"
	"testing"

	"github.com/go-chi/chi/v5"

	"github.com/omeg4monx/router-gateway/internal/contracts"
	"github.com/omeg4monx/router-gateway/internal/modules/arena"
	"github.com/omeg4monx/router-gateway/internal/platform/db"
	"github.com/omeg4monx/router-gateway/internal/platform/eventbus"
	"github.com/omeg4monx/router-gateway/internal/platform/lifecycle"
)

type stubDep struct{ name string }

func (s stubDep) Name() string                         { return s.name }
func (s stubDep) Dependencies() []string               { return nil }
func (s stubDep) Migrations() fs.FS                    { return nil }
func (s stubDep) RegisterRoutes(lifecycle.RouteMounts) {}
func (s stubDep) Start(context.Context) error          { return nil }
func (s stubDep) Stop(context.Context) error           { return nil }

type stubModelLister struct {
	models []contracts.ModelDescriptor
}

func (s *stubModelLister) ListModels(context.Context, []string) ([]contracts.ModelDescriptor, error) {
	return s.models, nil
}

type stubCompleter struct {
	mu       sync.Mutex
	requests []contracts.ChatCompletionRequest
	fn       func(ctx context.Context, req contracts.ChatCompletionRequest) (contracts.ChatCompletionResult, error)
}

func (s *stubCompleter) CompleteChat(ctx context.Context, req contracts.ChatCompletionRequest) (contracts.ChatCompletionResult, error) {
	s.mu.Lock()
	s.requests = append(s.requests, req)
	fn := s.fn
	s.mu.Unlock()
	return fn(ctx, req)
}

func (s *stubCompleter) lastRequest() contracts.ChatCompletionRequest {
	s.mu.Lock()
	defer s.mu.Unlock()
	if len(s.requests) == 0 {
		return contracts.ChatCompletionRequest{}
	}
	return s.requests[len(s.requests)-1]
}

type arenaTestEnv struct {
	srv       *httptest.Server
	completer *stubCompleter
}

func setupArenaTest(t *testing.T) *arenaTestEnv {
	t.Helper()

	database, err := db.Open(filepath.Join(t.TempDir(), "arena.sqlite"))
	if err != nil {
		t.Fatalf("db.Open: %v", err)
	}
	t.Cleanup(func() { _ = database.Close() })

	comp := &stubCompleter{
		fn: func(_ context.Context, req contracts.ChatCompletionRequest) (contracts.ChatCompletionResult, error) {
			if req.Model == "broken-model" {
				return contracts.ChatCompletionResult{Model: req.Model, LatencyMs: 12}, errors.New("upstream 500")
			}
			return contracts.ChatCompletionResult{
				Model:            req.Model,
				Content:          "Reply from " + req.Model,
				PromptTokens:     10,
				CompletionTokens: 5,
				TotalTokens:      15,
				LatencyMs:        25,
			}, nil
		},
	}

	mod := arena.New(arena.Options{
		DB:        database,
		Completer: comp,
		Models: &stubModelLister{models: []contracts.ModelDescriptor{
			{ID: "gpt-4o", OwnedBy: "openai", Kind: "model"},
			{ID: "claude-sonnet", OwnedBy: "anthropic", Kind: "model"},
		}},
		Bus:    eventbus.New(),
		Logger: slog.New(slog.NewTextHandler(io.Discard, nil)),
	})

	reg := lifecycle.NewRegistry(nil)
	for _, dep := range []lifecycle.Module{stubDep{name: "system"}, stubDep{name: "gateway"}, mod} {
		if err := reg.Register(dep); err != nil {
			t.Fatalf("register %s: %v", dep.Name(), err)
		}
	}
	if err := reg.MigrateAll(context.Background(), database); err != nil {
		t.Fatalf("MigrateAll: %v", err)
	}

	r := chi.NewRouter()
	r.Route("/api", func(ar chi.Router) {
		mod.RegisterRoutes(lifecycle.RouteMounts{AdminAPI: ar})
	})

	if err := reg.StartAll(context.Background()); err != nil {
		t.Fatalf("StartAll: %v", err)
	}
	t.Cleanup(func() { _ = reg.StopAll(context.Background()) })

	srv := httptest.NewServer(r)
	t.Cleanup(srv.Close)

	return &arenaTestEnv{srv: srv, completer: comp}
}

func doJSON(t *testing.T, method, url string, body any) (*http.Response, []byte) {
	t.Helper()
	var rdr io.Reader
	if body != nil {
		b, err := json.Marshal(body)
		if err != nil {
			t.Fatalf("marshal: %v", err)
		}
		rdr = bytes.NewReader(b)
	}
	req, err := http.NewRequest(method, url, rdr)
	if err != nil {
		t.Fatalf("NewRequest: %v", err)
	}
	if body != nil {
		req.Header.Set("Content-Type", "application/json")
	}
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Fatalf("Do %s %s: %v", method, url, err)
	}
	defer resp.Body.Close()
	raw, _ := io.ReadAll(resp.Body)
	return resp, raw
}

func TestArenaBattlePartialFailureVoteAndLeaderboard(t *testing.T) {
	env := setupArenaTest(t)

	// 1. Models endpoint.
	resp, raw := doJSON(t, http.MethodGet, env.srv.URL+"/api/arena/models", nil)
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("GET models status=%d body=%s", resp.StatusCode, raw)
	}

	// 2. Reject < 2 models.
	resp, _ = doJSON(t, http.MethodPost, env.srv.URL+"/api/arena/battles", map[string]any{
		"prompt": "hi",
		"models": []string{"gpt-4o"},
	})
	if resp.StatusCode != http.StatusBadRequest {
		t.Fatalf("expected 400 for <2 models, got %d", resp.StatusCode)
	}

	// 3. Run battle with partial failure (gpt-4o ok, broken-model fails).
	resp, raw = doJSON(t, http.MethodPost, env.srv.URL+"/api/arena/battles", map[string]any{
		"prompt":        "Compare Go and Rust",
		"system_prompt": "Be concise",
		"models":        []string{"gpt-4o", "broken-model"},
	})
	if resp.StatusCode != http.StatusCreated {
		t.Fatalf("POST battle status=%d body=%s", resp.StatusCode, raw)
	}
	var battle struct {
		ID         string `json:"id"`
		Contenders []struct {
			Model   string `json:"model"`
			Status  string `json:"status"`
			Content string `json:"content"`
			Error   string `json:"error"`
		} `json:"contenders"`
	}
	if err := json.Unmarshal(raw, &battle); err != nil {
		t.Fatalf("unmarshal battle: %v", err)
	}
	if len(battle.Contenders) != 2 {
		t.Fatalf("expected 2 contenders, got %d", len(battle.Contenders))
	}
	if battle.Contenders[0].Status != "ok" || battle.Contenders[1].Status != "error" {
		t.Fatalf("expected [ok, error], got %+v", battle.Contenders)
	}

	// 4. Reject invalid winner vote.
	resp, _ = doJSON(t, http.MethodPost, env.srv.URL+"/api/arena/battles/"+battle.ID+"/vote", map[string]any{
		"winner_model": "non-participant",
	})
	if resp.StatusCode != http.StatusBadRequest {
		t.Fatalf("expected 400 for non-participant vote, got %d", resp.StatusCode)
	}

	// 5. Vote gpt-4o winner and verify leaderboard.
	resp, raw = doJSON(t, http.MethodPost, env.srv.URL+"/api/arena/battles/"+battle.ID+"/vote", map[string]any{
		"winner_model": "gpt-4o",
	})
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("vote status=%d body=%s", resp.StatusCode, raw)
	}

	resp, raw = doJSON(t, http.MethodGet, env.srv.URL+"/api/arena/leaderboard", nil)
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("leaderboard status=%d body=%s", resp.StatusCode, raw)
	}
	var lb struct {
		Entries []struct {
			Model string `json:"model"`
			Wins  int    `json:"wins"`
		} `json:"entries"`
	}
	if err := json.Unmarshal(raw, &lb); err != nil || len(lb.Entries) < 2 {
		t.Fatalf("unexpected leaderboard: %s", raw)
	}
	if lb.Entries[0].Model != "gpt-4o" || lb.Entries[0].Wins != 1 {
		t.Fatalf("expected gpt-4o at #1 with 1 win, got %+v", lb.Entries[0])
	}

	// 6. Delete battle.
	resp, _ = doJSON(t, http.MethodDelete, env.srv.URL+"/api/arena/battles/"+battle.ID, nil)
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("delete battle status=%d", resp.StatusCode)
	}
}

func TestBasicChatSessionsHistoryAndValidation(t *testing.T) {
	env := setupArenaTest(t)

	// 1. Nonexistent session -> 404.
	resp, _ := doJSON(t, http.MethodPost, env.srv.URL+"/api/chat/sessions/missing-id/messages", map[string]any{
		"content": "hello",
	})
	if resp.StatusCode != http.StatusNotFound {
		t.Fatalf("expected 404 for missing session, got %d", resp.StatusCode)
	}

	// 2. Create session.
	resp, raw := doJSON(t, http.MethodPost, env.srv.URL+"/api/chat/sessions", map[string]any{
		"model":         "gpt-4o",
		"system_prompt": "Always reply briefly.",
	})
	if resp.StatusCode != http.StatusCreated {
		t.Fatalf("create session status=%d body=%s", resp.StatusCode, raw)
	}
	var session struct {
		ID    string `json:"id"`
		Title string `json:"title"`
	}
	if err := json.Unmarshal(raw, &session); err != nil || session.ID == "" {
		t.Fatalf("invalid session response: %s", raw)
	}

	// 3. Empty message -> 400.
	resp, _ = doJSON(t, http.MethodPost, env.srv.URL+"/api/chat/sessions/"+session.ID+"/messages", map[string]any{
		"content": "   ",
	})
	if resp.StatusCode != http.StatusBadRequest {
		t.Fatalf("expected 400 for empty message, got %d", resp.StatusCode)
	}

	// 4. Turn 1 + Turn 2 verify chronological history passed to ChatCompleter.
	resp, raw = doJSON(t, http.MethodPost, env.srv.URL+"/api/chat/sessions/"+session.ID+"/messages", map[string]any{
		"content": "First question about Go",
	})
	if resp.StatusCode != http.StatusCreated {
		t.Fatalf("turn 1 status=%d body=%s", resp.StatusCode, raw)
	}

	resp, raw = doJSON(t, http.MethodPost, env.srv.URL+"/api/chat/sessions/"+session.ID+"/messages", map[string]any{
		"content": "Follow-up question",
	})
	if resp.StatusCode != http.StatusCreated {
		t.Fatalf("turn 2 status=%d body=%s", resp.StatusCode, raw)
	}

	lastReq := env.completer.lastRequest()
	// Expect: system + user1 + assistant1 + user2 = 4 messages in chronological order.
	if len(lastReq.Messages) != 4 {
		t.Fatalf("expected 4 chronological messages, got %d: %+v", len(lastReq.Messages), lastReq.Messages)
	}
	if lastReq.Messages[0].Role != "system" || lastReq.Messages[1].Role != "user" || lastReq.Messages[2].Role != "assistant" || lastReq.Messages[3].Role != "user" {
		t.Fatalf("unexpected message order: %+v", lastReq.Messages)
	}

	// 5. Fetch session detail and verify auto-title + 4 persisted messages.
	resp, raw = doJSON(t, http.MethodGet, env.srv.URL+"/api/chat/sessions/"+session.ID, nil)
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("GET session status=%d body=%s", resp.StatusCode, raw)
	}
	var detail struct {
		Session struct {
			Title string `json:"title"`
		} `json:"session"`
		Messages []map[string]any `json:"messages"`
	}
	if err := json.Unmarshal(raw, &detail); err != nil {
		t.Fatalf("unmarshal detail: %v", err)
	}
	if detail.Session.Title == "New chat" || detail.Session.Title == "" {
		t.Fatalf("expected auto-generated title, got %q", detail.Session.Title)
	}
	if len(detail.Messages) != 4 {
		t.Fatalf("expected 4 persisted chat messages, got %d", len(detail.Messages))
	}

	// 6. Delete session.
	resp, _ = doJSON(t, http.MethodDelete, env.srv.URL+"/api/chat/sessions/"+session.ID, nil)
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("DELETE session status=%d", resp.StatusCode)
	}
}
