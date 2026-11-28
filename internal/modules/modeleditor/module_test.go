package modeleditor

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"path/filepath"
	"strings"
	"sync"
	"testing"

	"github.com/go-chi/chi/v5"

	"github.com/omeg4monx/router-gateway/internal/contracts"
	"github.com/omeg4monx/router-gateway/internal/platform/db"
	"github.com/omeg4monx/router-gateway/internal/platform/eventbus"
	"github.com/omeg4monx/router-gateway/internal/platform/lifecycle"
)

type testEnv struct {
	mod    *Module
	router chi.Router
}

func setupModelEditorTest(t *testing.T) *testEnv {
	t.Helper()
	database, err := db.Open(filepath.Join(t.TempDir(), "test.sqlite"))
	if err != nil {
		t.Fatalf("db.Open: %v", err)
	}
	t.Cleanup(func() { _ = database.Close() })

	mod := New(Options{
		DB:     database,
		Bus:    eventbus.New(),
		Logger: slog.New(slog.NewTextHandler(io.Discard, nil)),
	})

	ctx := context.Background()
	if err := database.MigrateModule(ctx, mod.Name(), mod.Migrations()); err != nil {
		t.Fatalf("MigrateModule: %v", err)
	}
	if err := mod.Start(ctx); err != nil {
		t.Fatalf("Start: %v", err)
	}
	t.Cleanup(func() { _ = mod.Stop(context.Background()) })

	r := chi.NewRouter()
	r.Route("/api", func(ar chi.Router) {
		mod.RegisterRoutes(lifecycle.RouteMounts{AdminAPI: ar})
	})
	return &testEnv{mod: mod, router: r}
}

func doJSON(t *testing.T, h http.Handler, method, path string, payload any) *httptest.ResponseRecorder {
	t.Helper()
	var body io.Reader
	if payload != nil {
		raw, err := json.Marshal(payload)
		if err != nil {
			t.Fatalf("marshal: %v", err)
		}
		body = bytes.NewReader(raw)
	}
	req := httptest.NewRequest(method, path, body)
	if payload != nil {
		req.Header.Set("Content-Type", "application/json")
	}
	rec := httptest.NewRecorder()
	h.ServeHTTP(rec, req)
	return rec
}

func TestModelEditorOverridesAcrossAllFourProtocols(t *testing.T) {
	env := setupModelEditorTest(t)

	if env.mod.Priority() != 20 {
		t.Fatalf("expected priority 20, got %d", env.mod.Priority())
	}

	temp := 0.25
	maxTokens := 2048
	rec := doJSON(t, env.router, http.MethodPost, "/api/model-editor/overrides", map[string]any{
		"name":                 "Strict Coding Profile",
		"model_pattern":        "coding-*",
		"target_model":         "claude-sonnet-4",
		"system_prompt_prefix": "[PREFIX]",
		"system_prompt_suffix": "[SUFFIX]",
		"temperature":          temp,
		"max_tokens":           maxTokens,
		"extra_params":         map[string]any{"seed": 42},
		"strip_params":         []string{"presence_penalty"},
		"priority":             10,
	})
	if rec.Code != http.StatusCreated {
		t.Fatalf("create override: status %d body %s", rec.Code, rec.Body.String())
	}

	cases := []struct {
		name     string
		protocol string
		body     string
		check    func(t *testing.T, out map[string]any)
	}{
		{
			name:     "openai-chat",
			protocol: "openai-chat",
			body:     `{"model":"coding-fast","presence_penalty":0.5,"custom_keep":true,"messages":[{"role":"system","content":"Base sys"},{"role":"user","content":"hi"}]}`,
			check: func(t *testing.T, out map[string]any) {
				if out["model"] != "claude-sonnet-4" {
					t.Fatalf("expected model rewrite, got %v", out["model"])
				}
				if _, exists := out["presence_penalty"]; exists {
					t.Fatal("expected presence_penalty to be stripped")
				}
				if out["custom_keep"] != true || out["seed"] != float64(42) || out["temperature"] != 0.25 || out["max_tokens"] != float64(2048) {
					t.Fatalf("unexpected params: %+v", out)
				}
				msgs := out["messages"].([]any)
				sysMsg := msgs[0].(map[string]any)["content"].(string)
				if !strings.HasPrefix(sysMsg, "[PREFIX]") || !strings.HasSuffix(sysMsg, "[SUFFIX]") {
					t.Fatalf("unexpected system message: %q", sysMsg)
				}
			},
		},
		{
			name:     "anthropic",
			protocol: "anthropic",
			body:     `{"model":"coding-fast","presence_penalty":0.5,"system":"Base sys","messages":[{"role":"user","content":"hi"}]}`,
			check: func(t *testing.T, out map[string]any) {
				if out["model"] != "claude-sonnet-4" {
					t.Fatalf("expected model rewrite, got %v", out["model"])
				}
				sys := out["system"].(string)
				if !strings.HasPrefix(sys, "[PREFIX]") || !strings.HasSuffix(sys, "[SUFFIX]") {
					t.Fatalf("unexpected anthropic system: %q", sys)
				}
			},
		},
		{
			name:     "openai-responses",
			protocol: "openai-responses",
			body:     `{"model":"coding-fast","instructions":"Base inst","input":"hi"}`,
			check: func(t *testing.T, out map[string]any) {
				inst := out["instructions"].(string)
				if !strings.HasPrefix(inst, "[PREFIX]") || !strings.HasSuffix(inst, "[SUFFIX]") {
					t.Fatalf("unexpected responses instructions: %q", inst)
				}
				if out["max_output_tokens"] != float64(2048) {
					t.Fatalf("expected max_output_tokens=2048, got %v", out["max_output_tokens"])
				}
			},
		},
		{
			name:     "gemini",
			protocol: "gemini",
			body:     `{"contents":[{"role":"user","parts":[{"text":"hi"}]}]}`,
			check: func(t *testing.T, out map[string]any) {
				si := out["systemInstruction"].(map[string]any)
				parts := si["parts"].([]any)
				text := parts[0].(map[string]any)["text"].(string)
				if !strings.Contains(text, "[PREFIX]") || !strings.Contains(text, "[SUFFIX]") {
					t.Fatalf("unexpected gemini systemInstruction: %q", text)
				}
				gc := out["generationConfig"].(map[string]any)
				if gc["temperature"] != 0.25 || gc["maxOutputTokens"] != float64(2048) {
					t.Fatalf("unexpected gemini generationConfig: %+v", gc)
				}
			},
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			req := &contracts.CanonicalRequest{
				RequestID: "req-" + tc.name,
				Protocol:  tc.protocol,
				Model:     "coding-fast",
				Headers:   make(http.Header),
				Body:      []byte(tc.body),
				Metadata:  make(map[string]any),
			}
			if err := env.mod.OnRequest(context.Background(), req); err != nil {
				t.Fatalf("OnRequest: %v", err)
			}
			if req.Model != "claude-sonnet-4" {
				t.Fatalf("expected req.Model=claude-sonnet-4, got %q", req.Model)
			}
			var out map[string]any
			if err := json.Unmarshal(req.Body, &out); err != nil {
				t.Fatalf("invalid output JSON: %v", err)
			}
			tc.check(t, out)
		})
	}
}

func TestModelBypassRulesPriorityAndMetadataFlags(t *testing.T) {
	env := setupModelEditorTest(t)

	// Add an override for gpt-*.
	rec := doJSON(t, env.router, http.MethodPost, "/api/model-editor/overrides", map[string]any{
		"name":          "GPT Temp Override",
		"model_pattern": "gpt-*",
		"temperature":   0.1,
	})
	if rec.Code != http.StatusCreated {
		t.Fatalf("create override: status %d", rec.Code)
	}

	// Add bypass rules for header X-Raw-Mode: true -> bypass_overrides & bypass_tokensaver.
	rec = doJSON(t, env.router, http.MethodPost, "/api/model-editor/bypass", map[string]any{
		"name":          "Raw Mode Bypass Overrides",
		"model_pattern": "gpt-*",
		"header_name":   "X-Raw-Mode",
		"header_value":  "true",
		"action":        "bypass_overrides",
		"priority":      10,
	})
	if rec.Code != http.StatusCreated {
		t.Fatalf("create bypass 1: status %d body %s", rec.Code, rec.Body.String())
	}
	rec = doJSON(t, env.router, http.MethodPost, "/api/model-editor/bypass", map[string]any{
		"name":          "Raw Mode Bypass TokenSaver",
		"model_pattern": "gpt-*",
		"header_name":   "X-Raw-Mode",
		"header_value":  "true",
		"action":        "bypass_tokensaver",
		"priority":      20,
	})
	if rec.Code != http.StatusCreated {
		t.Fatalf("create bypass 2: status %d", rec.Code)
	}

	hdr := make(http.Header)
	hdr.Set("X-Raw-Mode", "true")
	req := &contracts.CanonicalRequest{
		RequestID: "req-bypass",
		Protocol:  "openai-chat",
		Model:     "gpt-4o",
		Headers:   hdr,
		Body:      []byte(`{"model":"gpt-4o","temperature":0.9,"messages":[{"role":"user","content":"hi"}]}`),
		Metadata:  make(map[string]any),
	}
	if err := env.mod.OnRequest(context.Background(), req); err != nil {
		t.Fatalf("OnRequest: %v", err)
	}
	if req.Metadata["bypass_tokensaver"] != true {
		t.Fatalf("expected bypass_tokensaver=true, got %v", req.Metadata["bypass_tokensaver"])
	}
	var out map[string]any
	_ = json.Unmarshal(req.Body, &out)
	if out["temperature"] != 0.9 {
		t.Fatalf("expected temperature to stay 0.9 due to bypass_overrides, got %v", out["temperature"])
	}

	meta := &contracts.ResponseMeta{Headers: make(http.Header), Metadata: req.Metadata}
	_ = env.mod.OnResponseHeaders(context.Background(), meta)
	if meta.Headers.Get("X-Omeg4mon-Model-Bypass") == "" {
		t.Fatal("expected X-Omeg4mon-Model-Bypass header to be set")
	}
}

func TestModelEditorCRUDValidationAndDryRunPreview(t *testing.T) {
	env := setupModelEditorTest(t)

	// Invalid bypass action returns 400.
	rec := doJSON(t, env.router, http.MethodPost, "/api/model-editor/bypass", map[string]any{
		"name":          "Invalid",
		"model_pattern": "*",
		"action":        "unknown_action",
	})
	if rec.Code != http.StatusBadRequest {
		t.Fatalf("expected 400 for invalid bypass action, got %d", rec.Code)
	}

	// Valid rewrite_model bypass rule.
	rec = doJSON(t, env.router, http.MethodPost, "/api/model-editor/bypass", map[string]any{
		"name":          "Redirect Legacy",
		"model_pattern": "legacy-*",
		"action":        "rewrite_model",
		"target_model":  "gpt-4o-mini",
		"priority":      5,
	})
	if rec.Code != http.StatusCreated {
		t.Fatalf("expected 201, got %d body %s", rec.Code, rec.Body.String())
	}

	// Preview endpoint applies rewrite_model cleanly.
	prevRec := doJSON(t, env.router, http.MethodPost, "/api/model-editor/preview", map[string]any{
		"protocol": "openai-chat",
		"model":    "legacy-v1",
		"body":     map[string]any{"model": "legacy-v1", "messages": []any{map[string]any{"role": "user", "content": "hi"}}},
	})
	if prevRec.Code != http.StatusOK {
		t.Fatalf("preview: status %d body %s", prevRec.Code, prevRec.Body.String())
	}
	var prev PreviewResult
	if err := json.Unmarshal(prevRec.Body.Bytes(), &prev); err != nil {
		t.Fatalf("unmarshal preview: %v", err)
	}
	if prev.Model != "gpt-4o-mini" || len(prev.AppliedBypassRules) != 1 {
		t.Fatalf("unexpected preview result: %+v", prev)
	}
}

func TestModelEditorFinalReviewFindings(t *testing.T) {
	env := setupModelEditorTest(t)

	// Finding #4: Invalid glob pattern must be rejected; valid slash-agnostic glob matches provider/model.
	badRec := doJSON(t, env.router, http.MethodPost, "/api/model-editor/overrides", map[string]any{
		"name":          "Bad Glob",
		"model_pattern": "gpt-[4",
	})
	if badRec.Code != http.StatusBadRequest {
		t.Fatalf("expected 400 for malformed glob pattern, got %d", badRec.Code)
	}

	temp := 0.3
	okRec := doJSON(t, env.router, http.MethodPost, "/api/model-editor/overrides", map[string]any{
		"name":                 "Claude Family",
		"model_pattern":        "*claude*",
		"system_prompt_prefix": "[P]",
		"temperature":          temp,
		"extra_params":         map[string]any{"generationConfig": map[string]any{"seed": 7}},
	})
	if okRec.Code != http.StatusCreated {
		t.Fatalf("create override: status %d body %s", okRec.Code, okRec.Body.String())
	}

	// Slash-containing model name matches *claude*.
	// Finding #5: OpenAI developer role with array content wraps in-place instead of prepending duplicate system message.
	devReq := &contracts.CanonicalRequest{
		RequestID: "req-dev",
		Protocol:  "openai-chat",
		Model:     "anthropic/claude-sonnet-4",
		Body:      []byte(`{"model":"anthropic/claude-sonnet-4","messages":[{"role":"developer","content":[{"type":"text","text":"Dev rules"}]},{"role":"user","content":"hi"}]}`),
		Metadata:  make(map[string]any),
	}
	if err := env.mod.OnRequest(context.Background(), devReq); err != nil {
		t.Fatalf("OnRequest dev: %v", err)
	}
	var devOut map[string]any
	_ = json.Unmarshal(devReq.Body, &devOut)
	msgs := devOut["messages"].([]any)
	if len(msgs) != 2 {
		t.Fatalf("expected developer message wrapped in-place (len 2), got len %d: %s", len(msgs), devReq.Body)
	}

	// Finding #5: Gemini system_instruction & generationConfig deep-merge with extra_params.
	gemReq := &contracts.CanonicalRequest{
		RequestID: "req-gem",
		Protocol:  "gemini",
		Model:     "google/claude-bridge",
		Body:      []byte(`{"system_instruction":{"parts":[{"text":"Base"}]},"contents":[{"role":"user","parts":[{"text":"hi"}]}]}`),
		Metadata:  make(map[string]any),
	}
	if err := env.mod.OnRequest(context.Background(), gemReq); err != nil {
		t.Fatalf("OnRequest gem: %v", err)
	}
	var gemOut map[string]any
	_ = json.Unmarshal(gemReq.Body, &gemOut)
	if _, dup := gemOut["systemInstruction"]; dup {
		t.Fatalf("expected existing system_instruction to be updated in place without duplicate systemInstruction: %s", gemReq.Body)
	}
	gc := gemOut["generationConfig"].(map[string]any)
	if gc["temperature"] != 0.3 || gc["seed"] != float64(7) {
		t.Fatalf("expected generationConfig deep-merged (temperature=0.3 & seed=7), got %+v", gc)
	}
}

func TestListVirtualModels(t *testing.T) {
	t.Parallel()
	m := &Module{
		overrides: []ModelOverride{
			{ID: "1", Name: "claude-neon", ModelPattern: "claude-neon", TargetModel: "claude-sonnet-4", Enabled: true},
			{ID: "2", Name: "disabled-one", ModelPattern: "disabled-one", TargetModel: "gpt-4o", Enabled: false},
			{ID: "3", Name: "wildcard-rule", ModelPattern: "coding-*", TargetModel: "gpt-4o", Enabled: true},
			{ID: "4", Name: "no-target", ModelPattern: "gpt-4o", TargetModel: "", Enabled: true},
			{ID: "5", Name: "dup-alias", ModelPattern: "claude-neon", TargetModel: "claude-sonnet-4", Enabled: true},
			{ID: "6", Name: "qmark-glob", ModelPattern: "gpt-4?", TargetModel: "gpt-4o", Enabled: true},
			{ID: "7", Name: "bracket-glob", ModelPattern: "claude-[34]", TargetModel: "claude-sonnet-4", Enabled: true},
		},
	}

	got := m.ListVirtualModels()
	if len(got) != 1 {
		t.Fatalf("ListVirtualModels() len = %d, want 1 (%+v)", len(got), got)
	}
	if got[0].ID != "claude-neon" || got[0].OwnedBy != "model-editor" || got[0].Kind != "override" {
		t.Fatalf("ListVirtualModels()[0] = %+v, want claude-neon/model-editor/override", got[0])
	}
}

func TestConcurrentExtraParamsNoDataRaceOrCacheMutation(t *testing.T) {
	env := setupModelEditorTest(t)

	// Override 1 introduces a nested map ("generationConfig": {"seed": 7}).
	rec1 := doJSON(t, env.router, http.MethodPost, "/api/model-editor/overrides", map[string]any{
		"name":          "Base GenConfig",
		"model_pattern": "race-*",
		"extra_params":  map[string]any{"generationConfig": map[string]any{"seed": 7}},
		"priority":      10,
	})
	if rec1.Code != http.StatusCreated {
		t.Fatalf("create override 1: status %d body %s", rec1.Code, rec1.Body.String())
	}

	// Override 2 merges another key ("topK": 40) into the same nested map.
	rec2 := doJSON(t, env.router, http.MethodPost, "/api/model-editor/overrides", map[string]any{
		"name":          "Patch GenConfig",
		"model_pattern": "race-*",
		"extra_params":  map[string]any{"generationConfig": map[string]any{"topK": 40}},
		"priority":      20,
	})
	if rec2.Code != http.StatusCreated {
		t.Fatalf("create override 2: status %d body %s", rec2.Code, rec2.Body.String())
	}

	var wg sync.WaitGroup
	for i := range 16 {
		wg.Add(1)
		go func(idx int) {
			defer wg.Done()
			req := &contracts.CanonicalRequest{
				RequestID: "race-req",
				Protocol:  "openai-chat",
				Model:     "race-model",
				Body:      []byte(`{"model":"race-model","messages":[{"role":"user","content":"hi"}]}`),
				Metadata:  make(map[string]any),
			}
			if err := env.mod.OnRequest(context.Background(), req); err != nil {
				t.Errorf("goroutine %d OnRequest: %v", idx, err)
				return
			}
			var out map[string]any
			if err := json.Unmarshal(req.Body, &out); err != nil {
				t.Errorf("goroutine %d unmarshal: %v", idx, err)
				return
			}
			gc, ok := out["generationConfig"].(map[string]any)
			if !ok || gc["seed"] != float64(7) || gc["topK"] != float64(40) {
				t.Errorf("goroutine %d unexpected generationConfig: %+v", idx, gc)
			}
		}(i)
	}
	wg.Wait()

	// Cached override 1's ExtraParams must remain untouched (must NOT contain topK).
	_, overrides := env.mod.snapshot()
	baseGC, ok := overrides[0].ExtraParams["generationConfig"].(map[string]any)
	if !ok {
		t.Fatalf("expected generationConfig map on cached override 0")
	}
	if _, leaked := baseGC["topK"]; leaked {
		t.Fatalf("cached override ExtraParams was mutated in-place: %+v", baseGC)
	}
	if overrides[0].compiledPattern == nil {
		t.Fatalf("expected precompiled pattern regex on cached override 0")
	}
}

func TestDualInjectionAndIdempotency(t *testing.T) {
	t.Parallel()
	ov := ModelOverride{
		ID:                 "ov-dual",
		Name:               "dual-test",
		ModelPattern:       "dual-test",
		TargetModel:        "claude-sonnet-4",
		SystemPromptSuffix: "ANCHOR_PROMPT",
		ExtraParams:        map[string]any{"_prompt_mode": "dual", "seed": 7},
		Enabled:            true,
	}

	rawBody := []byte(`{
		"model": "dual-test",
		"messages": [
			{"role": "system", "content": "BASE_SYS"},
			{"role": "user", "content": "TURN_1"},
			{"role": "assistant", "content": "RESP_1"},
			{"role": "user", "content": "TURN_2"}
		]
	}`)

	outModel, outBody, _, _ := applyPipeline("openai-chat", "dual-test", "", nil, rawBody, map[string]any{}, nil, []ModelOverride{ov})
	if outModel != "claude-sonnet-4" {
		t.Fatalf("outModel = %q, want claude-sonnet-4", outModel)
	}
	if strings.Contains(string(outBody), "_prompt_mode") {
		t.Fatalf("internal _prompt_mode leaked into outgoing body: %s", string(outBody))
	}

	var parsed struct {
		Seed     int `json:"seed"`
		Messages []struct {
			Role    string `json:"role"`
			Content string `json:"content"`
		} `json:"messages"`
	}
	if err := json.Unmarshal(outBody, &parsed); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}
	if parsed.Seed != 7 {
		t.Fatalf("seed = %d, want 7", parsed.Seed)
	}
	if parsed.Messages[0].Content != "BASE_SYS\n\nANCHOR_PROMPT" {
		t.Fatalf("system = %q, want cache-safe append", parsed.Messages[0].Content)
	}
	if parsed.Messages[1].Content != "TURN_1" {
		t.Fatalf("turn 1 should remain untouched, got %q", parsed.Messages[1].Content)
	}
	if parsed.Messages[3].Content != "ANCHOR_PROMPT\n\nTURN_2" {
		t.Fatalf("latest user turn = %q, want anchored", parsed.Messages[3].Content)
	}

	// Re-applying the same override (e.g. retry/fallback) must be strictly idempotent.
	_, secondBody, _, _ := applyPipeline("openai-chat", "dual-test", "", nil, outBody, map[string]any{}, nil, []ModelOverride{ov})
	if strings.Count(string(secondBody), "ANCHOR_PROMPT") != 2 {
		t.Fatalf("expected idempotent 2 occurrences (system + latest user), got %d in %s", strings.Count(string(secondBody), "ANCHOR_PROMPT"), string(secondBody))
	}
}
