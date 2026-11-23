package gateway

import (
	"bytes"
	"context"
	"crypto/rand"
	"embed"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/fs"
	"log/slog"
	"net/http"
	"net/url"
	"slices"
	"strings"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"

	"github.com/omeg4monx/router-gateway/internal/contracts"
	"github.com/omeg4monx/router-gateway/internal/platform/db"
	"github.com/omeg4monx/router-gateway/internal/platform/httpclient"
	"github.com/omeg4monx/router-gateway/internal/platform/httpserver"
	"github.com/omeg4monx/router-gateway/internal/platform/lifecycle"
)

const (
	moduleName = "gateway"

	maxInboundBodyBytes  = 32 << 20 // 32 MiB
	maxResponseBodyBytes = 32 << 20 // 32 MiB
	maxErrorBodyBytes    = 64 << 10 // 64 KiB
)

//go:embed migrations/*.sql
var migrationsFS embed.FS

// Options configures a new gateway Module instance.
type Options struct {
	DB           *db.DB
	HTTPClient   *http.Client
	Keys         contracts.KeyValidator
	Resolver     contracts.ProviderResolver
	Models       contracts.ModelLister
	Providers    contracts.ProviderStore
	Health       contracts.AccountHealthTracker
	Interceptors []contracts.GatewayInterceptor
	Usage        contracts.UsageRecorder
	Logger       *slog.Logger
	ChainTTL     time.Duration
}

// Module implements lifecycle.Module and exposes multi-protocol AI gateway
// endpoints onto the streaming Gateway router mount.
type Module struct {
	db           *db.DB
	httpClient   *http.Client
	keys         contracts.KeyValidator
	resolver     contracts.ProviderResolver
	models       contracts.ModelLister
	providers    contracts.ProviderStore
	health       contracts.AccountHealthTracker
	interceptors []contracts.GatewayInterceptor
	usage        contracts.UsageRecorder
	logger       *slog.Logger
	chains       *chainStore
}

var _ lifecycle.Module = (*Module)(nil)

// New constructs a gateway Module with sorted interceptors and bounded
// response chain persistence.
func New(opts Options) *Module {
	client := opts.HTTPClient
	if client == nil {
		client = httpclient.NewClient(httpclient.Options{})
	}
	logger := opts.Logger
	if logger == nil {
		logger = slog.Default()
	}

	sorted := make([]contracts.GatewayInterceptor, 0, len(opts.Interceptors))
	for _, ic := range opts.Interceptors {
		if ic != nil {
			sorted = append(sorted, ic)
		}
	}
	slices.SortStableFunc(sorted, func(a, b contracts.GatewayInterceptor) int {
		if a.Priority() != b.Priority() {
			if a.Priority() < b.Priority() {
				return -1
			}
			return 1
		}
		return strings.Compare(a.Name(), b.Name())
	})

	return &Module{
		db:           opts.DB,
		httpClient:   client,
		keys:         opts.Keys,
		resolver:     opts.Resolver,
		models:       opts.Models,
		providers:    opts.Providers,
		health:       opts.Health,
		interceptors: sorted,
		usage:        opts.Usage,
		logger:       logger,
		chains:       newChainStore(opts.DB, opts.ChainTTL),
	}
}

// Name returns the canonical module name.
func (m *Module) Name() string {
	return moduleName
}

// Dependencies declares the upstream modules required by the gateway.
func (m *Module) Dependencies() []string {
	return []string{"system", "keys", "providers", "combos"}
}

// Migrations returns the embedded SQL migration filesystem rooted at migrations/.
func (m *Module) Migrations() fs.FS {
	sub, err := fs.Sub(migrationsFS, "migrations")
	if err != nil {
		panic(fmt.Sprintf("gateway: load embedded migrations: %v", err))
	}
	return sub
}

// RegisterRoutes mounts all OpenAI, Anthropic, Gemini, and compatibility
// gateway routes onto mounts.Gateway.
func (m *Module) RegisterRoutes(mounts lifecycle.RouteMounts) {
	if mounts.Gateway == nil {
		return
	}

	mounts.Gateway.Post("/v1/chat/completions", m.handleOpenAIChat)
	mounts.Gateway.Post("/v1/messages", m.handleAnthropicMessages)
	mounts.Gateway.Post("/v1/messages/count_tokens", m.handleAnthropicCountTokens)
	mounts.Gateway.Post("/v1/responses", m.handleOpenAIResponses)
	mounts.Gateway.Post("/v1/responses/compact", m.handleOpenAIResponsesCompact)
	mounts.Gateway.Get("/v1/models", m.handleListModels)
	mounts.Gateway.Get("/v1/models/{id}", m.handleGetModel)
	mounts.Gateway.Get("/v1/models/*", m.handleGetModel)

	mounts.Gateway.Get("/v1beta/models", m.handleGeminiListModels)
	mounts.Gateway.Post("/v1beta/models/{modelAction}", m.handleGeminiGenerate)

	// Compatibility aliases.
	mounts.Gateway.Post("/codex/responses", m.handleOpenAIResponses)
	mounts.Gateway.Post("/responses", m.handleOpenAIResponses)
	mounts.Gateway.Get("/models", m.handleListModels)
	mounts.Gateway.Get("/models/{id}", m.handleGetModel)
	mounts.Gateway.Get("/models/*", m.handleGetModel)
}

// Start validates that required dependencies are wired.
func (m *Module) Start(context.Context) error {
	if m.db == nil {
		return errors.New("gateway: database is required")
	}
	if m.keys == nil {
		return errors.New("gateway: key validator is required")
	}
	if m.resolver == nil {
		return errors.New("gateway: provider resolver is required")
	}
	if m.models == nil {
		return errors.New("gateway: model lister is required")
	}
	if m.providers == nil {
		return errors.New("gateway: provider store is required")
	}
	if m.health == nil {
		return errors.New("gateway: account health tracker is required")
	}
	return nil
}

// Stop performs graceful module teardown.
func (m *Module) Stop(context.Context) error {
	return nil
}

func (m *Module) handleOpenAIChat(w http.ResponseWriter, r *http.Request) {
	m.handleInference(w, r, ProtocolOpenAIChat, "", false, false)
}

func (m *Module) handleAnthropicMessages(w http.ResponseWriter, r *http.Request) {
	m.handleInference(w, r, ProtocolAnthropic, "", false, false)
}

func (m *Module) handleOpenAIResponses(w http.ResponseWriter, r *http.Request) {
	m.handleInference(w, r, ProtocolOpenAIResponses, "", false, false)
}

func (m *Module) handleOpenAIResponsesCompact(w http.ResponseWriter, r *http.Request) {
	m.handleInference(w, r, ProtocolOpenAIResponses, "", false, true)
}

func (m *Module) handleGeminiGenerate(w http.ResponseWriter, r *http.Request) {
	rawAction := chi.URLParam(r, "modelAction")
	if unescaped, err := url.PathUnescape(rawAction); err == nil {
		rawAction = unescaped
	}
	colon := strings.LastIndexByte(rawAction, ':')
	if colon <= 0 || colon == len(rawAction)-1 {
		if _, ok := m.authenticateRequest(w, r, ProtocolGemini); !ok {
			return
		}
		writeGatewayError(w, r, ProtocolGemini, http.StatusBadRequest, "invalid_request_error", "invalid_model_action", "expected /v1beta/models/{model}:generateContent or :streamGenerateContent")
		return
	}
	modelName := rawAction[:colon]
	action := rawAction[colon+1:]

	var stream bool
	switch action {
	case "generateContent":
		stream = false
	case "streamGenerateContent":
		stream = true
	default:
		if _, ok := m.authenticateRequest(w, r, ProtocolGemini); !ok {
			return
		}
		writeGatewayError(w, r, ProtocolGemini, http.StatusBadRequest, "invalid_request_error", "unsupported_action", fmt.Sprintf("unsupported Gemini action %q", action))
		return
	}

	m.handleInference(w, r, ProtocolGemini, modelName, stream, false)
}

func (m *Module) handleListModels(w http.ResponseWriter, r *http.Request) {
	keyInfo, ok := m.authenticateRequest(w, r, ProtocolOpenAIChat)
	if !ok {
		return
	}
	list, err := m.models.ListModels(r.Context(), keyInfo.AllowedModels)
	if err != nil {
		writeGatewayError(w, r, ProtocolOpenAIChat, http.StatusInternalServerError, "server_error", "model_list_failed", "failed to list models")
		return
	}
	items := make([]map[string]any, 0, len(list))
	for _, d := range list {
		items = append(items, map[string]any{
			"id":       d.ID,
			"object":   "model",
			"created":  0,
			"owned_by": d.OwnedBy,
		})
	}
	httpserver.WriteJSON(w, http.StatusOK, map[string]any{
		"object": "list",
		"data":   items,
	})
}

func (m *Module) handleGetModel(w http.ResponseWriter, r *http.Request) {
	keyInfo, ok := m.authenticateRequest(w, r, ProtocolOpenAIChat)
	if !ok {
		return
	}
	modelID := chi.URLParam(r, "id")
	if modelID == "" {
		modelID = chi.URLParam(r, "*")
	}
	if unescaped, err := url.PathUnescape(modelID); err == nil {
		modelID = unescaped
	}
	list, err := m.models.ListModels(r.Context(), keyInfo.AllowedModels)
	if err != nil {
		writeGatewayError(w, r, ProtocolOpenAIChat, http.StatusInternalServerError, "server_error", "model_lookup_failed", "failed to list models")
		return
	}
	for _, d := range list {
		if d.ID == modelID {
			httpserver.WriteJSON(w, http.StatusOK, map[string]any{
				"id":       d.ID,
				"object":   "model",
				"created":  0,
				"owned_by": d.OwnedBy,
			})
			return
		}
	}
	writeGatewayError(w, r, ProtocolOpenAIChat, http.StatusNotFound, "invalid_request_error", "model_not_found", fmt.Sprintf("model %q not found", modelID))
}

func (m *Module) handleGeminiListModels(w http.ResponseWriter, r *http.Request) {
	keyInfo, ok := m.authenticateRequest(w, r, ProtocolGemini)
	if !ok {
		return
	}
	list, err := m.models.ListModels(r.Context(), keyInfo.AllowedModels)
	if err != nil {
		writeGatewayError(w, r, ProtocolGemini, http.StatusInternalServerError, "server_error", "model_list_failed", "failed to list models")
		return
	}
	items := make([]map[string]any, 0, len(list))
	for _, d := range list {
		items = append(items, map[string]any{
			"name":                       "models/" + d.ID,
			"displayName":                d.ID,
			"supportedGenerationMethods": []string{"generateContent", "streamGenerateContent"},
		})
	}
	httpserver.WriteJSON(w, http.StatusOK, map[string]any{
		"models": items,
	})
}

func (m *Module) handleAnthropicCountTokens(w http.ResponseWriter, r *http.Request) {
	keyInfo, ok := m.authenticateRequest(w, r, ProtocolAnthropic)
	if !ok {
		return
	}
	body, err := readLimitedBody(w, r, maxInboundBodyBytes)
	if err != nil {
		writeGatewayError(w, r, ProtocolAnthropic, http.StatusBadRequest, "invalid_request_error", "invalid_body", "failed to read request body")
		return
	}

	canonReq, err := DecodeRequest(ProtocolAnthropic, body, "")
	if err != nil {
		writeGatewayError(w, r, ProtocolAnthropic, http.StatusBadRequest, "invalid_request_error", "invalid_payload", err.Error())
		return
	}
	if strings.TrimSpace(canonReq.Model) == "" {
		writeGatewayError(w, r, ProtocolAnthropic, http.StatusBadRequest, "invalid_request_error", "missing_model", "model is required")
		return
	}

	candidates, err := m.resolver.Resolve(r.Context(), canonReq.Model, keyInfo.AllowedModels)
	if err != nil && !errors.Is(err, contracts.ErrAllAccountsCoolingDown) {
		writeResolutionError(w, r, ProtocolAnthropic, err)
		return
	}

	if err == nil && len(candidates) > 0 && candidates[0].ProviderType == ProtocolAnthropic {
		if tokens, ok := m.tryUpstreamAnthropicCountTokens(r.Context(), candidates[0], body); ok {
			httpserver.WriteJSON(w, http.StatusOK, map[string]any{"input_tokens": tokens})
			return
		}
	}

	httpserver.WriteJSON(w, http.StatusOK, map[string]any{
		"input_tokens": EstimateTokens(canonReq),
	})
}

func (m *Module) tryUpstreamAnthropicCountTokens(ctx context.Context, cand contracts.ResolvedTarget, body []byte) (int64, bool) {
	creds, err := m.providers.ResolveCredentials(ctx, cand.ConnectionID)
	if err != nil {
		return 0, false
	}
	baseURL := cand.Endpoint
	if baseURL == "" {
		baseURL = creds.BaseURL
	}
	targetURL, err := buildUpstreamURL(baseURL, ProtocolAnthropic, cand.UpstreamModel, false, true, false)
	if err != nil {
		return 0, false
	}
	rewritten, err := RewriteModelField(body, ProtocolAnthropic, cand.UpstreamModel)
	if err != nil {
		return 0, false
	}
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, targetURL, bytes.NewReader(rewritten))
	if err != nil {
		return 0, false
	}
	applyUpstreamHeaders(req, creds, false)
	resp, err := m.httpClient.Do(req)
	if err != nil {
		return 0, false
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return 0, false
	}
	raw, err := io.ReadAll(io.LimitReader(resp.Body, maxErrorBodyBytes))
	if err != nil {
		return 0, false
	}
	var out struct {
		InputTokens int64 `json:"input_tokens"`
	}
	if err := json.Unmarshal(raw, &out); err != nil || out.InputTokens <= 0 {
		return 0, false
	}
	return out.InputTokens, true
}

type statusError interface {
	StatusCode() int
}

func (m *Module) handleInference(w http.ResponseWriter, r *http.Request, inProto, pathModel string, forceStream, isCompact bool) {
	startTime := time.Now()
	reqID := ensureRequestID(r)

	keyInfo, ok := m.authenticateRequest(w, r, inProto)
	if !ok {
		return
	}

	body, err := readLimitedBody(w, r, maxInboundBodyBytes)
	if err != nil {
		writeGatewayError(w, r, inProto, http.StatusBadRequest, "invalid_request_error", "invalid_body", "failed to read request body")
		return
	}

	reqModel, reqStream, err := parseInboundMeta(inProto, body, pathModel, forceStream)
	if err != nil {
		writeGatewayError(w, r, inProto, http.StatusBadRequest, "invalid_request_error", "invalid_json", err.Error())
		return
	}
	if strings.TrimSpace(reqModel) == "" {
		writeGatewayError(w, r, inProto, http.StatusBadRequest, "invalid_request_error", "missing_model", "model is required")
		return
	}

	icReq := &contracts.CanonicalRequest{
		RequestID: reqID,
		Protocol:  inProto,
		Model:     reqModel,
		Stream:    reqStream,
		KeyID:     keyInfo.ID,
		Headers:   r.Header.Clone(),
		Body:      body,
		Metadata:  make(map[string]any),
	}
	for _, ic := range m.interceptors {
		icErr := ic.OnRequest(r.Context(), icReq)
		if icErr == nil {
			continue
		}
		status := http.StatusBadRequest
		var se statusError
		if errors.As(icErr, &se) && se.StatusCode() >= 400 {
			status = se.StatusCode()
		}
		writeGatewayError(w, r, inProto, status, "invalid_request_error", "interceptor_rejected", icErr.Error())
		m.emitUsage(r.Context(), contracts.UsageRecord{
			RequestID:    reqID,
			KeyID:        keyInfo.ID,
			Model:        icReq.Model,
			Protocol:     inProto,
			Latency:      time.Since(startTime),
			StatusCode:   status,
			ErrorSummary: icErr.Error(),
			Stream:       icReq.Stream,
			Timestamp:    time.Now().UTC(),
		})
		return
	}

	candidates, err := m.resolver.Resolve(r.Context(), icReq.Model, keyInfo.AllowedModels)
	if err != nil {
		writeResolutionError(w, r, inProto, err)
		return
	}

	var (
		lastStatus = http.StatusBadGateway
		lastMsg    = "upstream request failed"
		lastCand   contracts.ResolvedTarget
	)

	for _, cand := range candidates {
		if r.Context().Err() != nil {
			return
		}
		lastCand = cand

		creds, credErr := m.providers.ResolveCredentials(r.Context(), cand.ConnectionID)
		if credErr != nil {
			m.health.RecordFailure(cand.ConnectionID, http.StatusBadGateway, 0)
			lastStatus = http.StatusBadGateway
			lastMsg = "failed to resolve upstream credentials"
			continue
		}

		outProto := cand.ProviderType
		if outProto == "" {
			outProto = creds.Protocol
		}
		if creds.Protocol == "" {
			creds.Protocol = outProto
		}
		baseURL := cand.Endpoint
		if baseURL == "" {
			baseURL = creds.BaseURL
		}
		upstreamModel := cand.UpstreamModel
		if upstreamModel == "" {
			upstreamModel = cand.Model
		}

		upBody, chainInputMsgs, prepErr := m.prepareUpstreamBody(r.Context(), inProto, outProto, keyInfo.ID, upstreamModel, icReq.Stream, isCompact, icReq.Body)
		if prepErr != nil {
			writeGatewayError(w, r, inProto, http.StatusBadRequest, "invalid_request_error", "invalid_payload", prepErr.Error())
			return
		}

		refreshed := false
		for {
			outcome := m.executeCandidateAttempt(
				r.Context(),
				w,
				r,
				reqID,
				keyInfo.ID,
				icReq.Model,
				inProto,
				outProto,
				icReq.Stream,
				isCompact,
				cand,
				creds,
				baseURL,
				upstreamModel,
				upBody,
				chainInputMsgs,
				startTime,
				!refreshed,
			)
			if outcome.handled {
				return
			}
			if outcome.needsRefresh {
				refreshed = true
				newCreds, refErr := m.providers.RefreshOAuthToken(r.Context(), cand.ConnectionID)
				if refErr != nil {
					if r.Context().Err() != nil || errors.Is(refErr, context.Canceled) {
						return
					}
					m.health.RecordFailure(cand.ConnectionID, http.StatusUnauthorized, 0)
					lastStatus = http.StatusUnauthorized
					lastMsg = "oauth token refresh failed"
					break
				}
				if newCreds.Protocol == "" {
					newCreds.Protocol = outProto
				}
				creds = newCreds
				continue
			}
			if outcome.clientCanceled {
				return
			}
			lastStatus = outcome.statusCode
			lastMsg = outcome.errMsg
			break
		}
	}

	m.emitUsage(r.Context(), contracts.UsageRecord{
		RequestID:    reqID,
		KeyID:        keyInfo.ID,
		Model:        icReq.Model,
		ProviderID:   lastCand.ProviderID,
		ConnectionID: lastCand.ConnectionID,
		Protocol:     inProto,
		Latency:      time.Since(startTime),
		StatusCode:   lastStatus,
		ErrorSummary: lastMsg,
		Stream:       icReq.Stream,
		Timestamp:    time.Now().UTC(),
	})
	writeGatewayError(w, r, inProto, lastStatus, "server_error", "upstream_failed", lastMsg)
}

type attemptOutcome struct {
	handled        bool
	needsRefresh   bool
	clientCanceled bool
	statusCode     int
	errMsg         string
}

func (m *Module) prepareUpstreamBody(
	ctx context.Context,
	inProto string,
	outProto string,
	keyID string,
	upstreamModel string,
	stream bool,
	isCompact bool,
	body []byte,
) ([]byte, []CanonicalMessage, error) {
	if isCompact && CanFastPath(inProto, outProto) {
		rewritten, err := RewriteModelField(body, inProto, upstreamModel)
		return rewritten, nil, err
	}

	if inProto == ProtocolOpenAIResponses && outProto != ProtocolOpenAIResponses {
		canonReq, err := DecodeRequest(inProto, body, upstreamModel)
		if err != nil {
			return nil, nil, err
		}
		canonReq.Stream = stream
		if canonReq.PreviousResponseID != "" {
			prior, found, chainErr := m.chains.loadChain(ctx, canonReq.PreviousResponseID, keyID, time.Now().UTC())
			if chainErr != nil {
				m.logger.Warn("load response chain failed", "err", chainErr)
			} else if found && len(prior) > 0 {
				combined := make([]CanonicalMessage, 0, len(prior)+len(canonReq.Messages))
				combined = append(combined, prior...)
				combined = append(combined, canonReq.Messages...)
				canonReq.Messages = combined
			}
			canonReq.PreviousResponseID = ""
		}
		if isCompact && len(canonReq.System) == 0 {
			canonReq.System = []SystemBlock{{Text: "Compact and summarize the conversation state concisely."}}
		}
		encoded, err := EncodeRequest(outProto, canonReq)
		if err != nil {
			return nil, nil, err
		}
		return encoded, slices.Clone(canonReq.Messages), nil
	}

	if CanFastPath(inProto, outProto) {
		rewritten, err := RewriteModelField(body, inProto, upstreamModel)
		return rewritten, nil, err
	}

	canonReq, err := DecodeRequest(inProto, body, upstreamModel)
	if err != nil {
		return nil, nil, err
	}
	canonReq.Stream = stream
	encoded, err := EncodeRequest(outProto, canonReq)
	if err != nil {
		return nil, nil, err
	}
	return encoded, nil, nil
}

func (m *Module) executeCandidateAttempt(
	ctx context.Context,
	w http.ResponseWriter,
	r *http.Request,
	reqID string,
	keyID string,
	requestedModel string,
	inProto string,
	outProto string,
	stream bool,
	isCompact bool,
	cand contracts.ResolvedTarget,
	creds contracts.UpstreamCredentials,
	baseURL string,
	upstreamModel string,
	upBody []byte,
	chainInputMsgs []CanonicalMessage,
	startTime time.Time,
	allowOAuthRefresh bool,
) attemptOutcome {
	upURL, err := buildUpstreamURL(baseURL, outProto, upstreamModel, stream, false, isCompact)
	if err != nil {
		m.health.RecordFailure(cand.ConnectionID, http.StatusBadGateway, 0)
		return attemptOutcome{statusCode: http.StatusBadGateway, errMsg: err.Error()}
	}
	upReq, err := http.NewRequestWithContext(ctx, http.MethodPost, upURL, bytes.NewReader(upBody))
	if err != nil {
		m.health.RecordFailure(cand.ConnectionID, http.StatusBadGateway, 0)
		return attemptOutcome{statusCode: http.StatusBadGateway, errMsg: err.Error()}
	}
	applyUpstreamHeaders(upReq, creds, stream)

	resp, err := m.httpClient.Do(upReq)
	if err != nil {
		if ctx.Err() != nil {
			return attemptOutcome{clientCanceled: true}
		}
		m.health.RecordFailure(cand.ConnectionID, http.StatusBadGateway, 0)
		return attemptOutcome{statusCode: http.StatusBadGateway, errMsg: err.Error()}
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 200 && resp.StatusCode < 300 {
		m.health.RecordSuccess(cand.ConnectionID)

		meta := &contracts.ResponseMeta{
			RequestID:    reqID,
			StatusCode:   resp.StatusCode,
			Headers:      resp.Header.Clone(),
			ProviderID:   cand.ProviderID,
			ConnectionID: cand.ConnectionID,
			Model:        requestedModel,
			Metadata:     make(map[string]any),
		}
		for _, ic := range m.interceptors {
			if icErr := ic.OnResponseHeaders(ctx, meta); icErr != nil {
				writeGatewayError(w, r, inProto, http.StatusBadGateway, "server_error", "interceptor_error", icErr.Error())
				return attemptOutcome{handled: true}
			}
		}
		copyResponseHeaders(w.Header(), meta.Headers)

		if stream {
			sres, streamErr := m.streamUpstreamSSE(ctx, w, resp.Body, reqID, inProto, outProto, requestedModel, startTime)
			if streamErr != nil && ctx.Err() == nil {
				m.logger.Warn("stream interrupted", "request_id", reqID, "err", streamErr)
			}
			if inProto == ProtocolOpenAIResponses && outProto != ProtocolOpenAIResponses && sres.responseID != "" {
				msgs := slices.Clone(chainInputMsgs)
				if sres.assistantMsg != nil {
					msgs = append(msgs, *sres.assistantMsg)
				}
				if saveErr := m.chains.saveChain(ctx, sres.responseID, keyID, requestedModel, msgs, time.Now().UTC()); saveErr != nil {
					m.logger.Warn("save streamed response chain failed", "err", saveErr)
				}
			}
			m.emitUsage(ctx, contracts.UsageRecord{
				RequestID:    reqID,
				KeyID:        keyID,
				Model:        requestedModel,
				ProviderID:   cand.ProviderID,
				ConnectionID: cand.ConnectionID,
				Protocol:     inProto,
				InputTokens:  sres.usage.InputTokens,
				OutputTokens: sres.usage.OutputTokens,
				CachedTokens: sres.usage.CachedTokens,
				Latency:      time.Since(startTime),
				TTFT:         sres.ttft,
				StatusCode:   resp.StatusCode,
				Stream:       true,
				Timestamp:    time.Now().UTC(),
			})
			return attemptOutcome{handled: true}
		}

		rawResp, readErr := io.ReadAll(io.LimitReader(resp.Body, maxResponseBodyBytes))
		if readErr != nil {
			writeGatewayError(w, r, inProto, http.StatusBadGateway, "server_error", "upstream_read_error", readErr.Error())
			return attemptOutcome{handled: true}
		}

		var usage CanonicalUsage
		outBytes := rawResp
		switch {
		case isCompact && CanFastPath(inProto, outProto):
			outBytes = rawResp
		case CanFastPath(inProto, outProto):
			if canon, decErr := DecodeResponse(outProto, rawResp); decErr == nil {
				usage = canon.Usage
			}
		default:
			canon, decErr := DecodeResponse(outProto, rawResp)
			if decErr != nil {
				writeGatewayError(w, r, inProto, http.StatusBadGateway, "server_error", "invalid_upstream_response", decErr.Error())
				return attemptOutcome{handled: true}
			}
			canon.Model = requestedModel
			if inProto == ProtocolOpenAIResponses && canon.ID == "" {
				canon.ID = randomResponseID()
			}
			usage = canon.Usage

			if isCompact && inProto == ProtocolOpenAIResponses {
				compactOut := buildCompactedResponse(canon)
				encoded, encErr := json.Marshal(compactOut)
				if encErr != nil {
					writeGatewayError(w, r, inProto, http.StatusBadGateway, "server_error", "encode_error", encErr.Error())
					return attemptOutcome{handled: true}
				}
				outBytes = encoded
			} else {
				encoded, encErr := EncodeResponse(inProto, canon)
				if encErr != nil {
					writeGatewayError(w, r, inProto, http.StatusBadGateway, "server_error", "encode_error", encErr.Error())
					return attemptOutcome{handled: true}
				}
				outBytes = encoded
			}

			if inProto == ProtocolOpenAIResponses && !isCompact && canon.ID != "" {
				msgs := slices.Clone(chainInputMsgs)
				if len(canon.Content) > 0 {
					msgs = append(msgs, CanonicalMessage{Role: "assistant", Content: canon.Content})
				}
				if saveErr := m.chains.saveChain(ctx, canon.ID, keyID, requestedModel, msgs, time.Now().UTC()); saveErr != nil {
					m.logger.Warn("save response chain failed", "err", saveErr)
				}
			}
		}

		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		w.WriteHeader(resp.StatusCode)
		_, _ = w.Write(outBytes)

		m.emitUsage(ctx, contracts.UsageRecord{
			RequestID:    reqID,
			KeyID:        keyID,
			Model:        requestedModel,
			ProviderID:   cand.ProviderID,
			ConnectionID: cand.ConnectionID,
			Protocol:     inProto,
			InputTokens:  usage.InputTokens,
			OutputTokens: usage.OutputTokens,
			CachedTokens: usage.CachedTokens,
			Latency:      time.Since(startTime),
			TTFT:         time.Since(startTime),
			StatusCode:   resp.StatusCode,
			Stream:       false,
			Timestamp:    time.Now().UTC(),
		})
		return attemptOutcome{handled: true}
	}

	if resp.StatusCode == http.StatusUnauthorized && creds.AuthType == "oauth" && allowOAuthRefresh {
		return attemptOutcome{needsRefresh: true}
	}

	errBody, _ := io.ReadAll(io.LimitReader(resp.Body, maxErrorBodyBytes))
	errMsg := extractUpstreamErrorMessage(errBody, resp.StatusCode)

	if classifyStatus(resp.StatusCode, errBody) == dispositionClientError {
		m.emitUsage(ctx, contracts.UsageRecord{
			RequestID:    reqID,
			KeyID:        keyID,
			Model:        requestedModel,
			ProviderID:   cand.ProviderID,
			ConnectionID: cand.ConnectionID,
			Protocol:     inProto,
			Latency:      time.Since(startTime),
			StatusCode:   resp.StatusCode,
			ErrorSummary: errMsg,
			Stream:       stream,
			Timestamp:    time.Now().UTC(),
		})
		if CanFastPath(inProto, outProto) && json.Valid(errBody) {
			w.Header().Set("Content-Type", "application/json; charset=utf-8")
			w.WriteHeader(resp.StatusCode)
			_, _ = w.Write(errBody)
			return attemptOutcome{handled: true}
		}
		writeGatewayError(w, r, inProto, resp.StatusCode, "invalid_request_error", "upstream_client_error", errMsg)
		return attemptOutcome{handled: true}
	}

	retryAfter := parseRetryAfter(resp.Header.Get("Retry-After"), time.Now().UTC())
	m.health.RecordFailure(cand.ConnectionID, resp.StatusCode, retryAfter)
	return attemptOutcome{
		statusCode: resp.StatusCode,
		errMsg:     errMsg,
	}
}

func (m *Module) emitUsage(ctx context.Context, rec contracts.UsageRecord) {
	if m.usage == nil {
		return
	}
	m.usage.Record(ctx, rec)
}

func readLimitedBody(w http.ResponseWriter, r *http.Request, maxBytes int64) ([]byte, error) {
	if r.Body == nil {
		return nil, errors.New("empty request body")
	}
	defer r.Body.Close()
	reader := http.MaxBytesReader(w, r.Body, maxBytes)
	return io.ReadAll(reader)
}

func parseInboundMeta(proto string, body []byte, pathModel string, forceStream bool) (model string, stream bool, err error) {
	trimmed := bytes.TrimSpace(body)
	if len(trimmed) == 0 {
		return "", false, errors.New("empty request body")
	}
	var meta struct {
		Model  string `json:"model"`
		Stream bool   `json:"stream"`
	}
	if err := json.Unmarshal(trimmed, &meta); err != nil {
		return "", false, fmt.Errorf("malformed JSON payload: %w", err)
	}
	model = meta.Model
	if proto == ProtocolGemini && pathModel != "" {
		model = pathModel
	}
	stream = meta.Stream
	if proto == ProtocolGemini {
		stream = forceStream
	}
	return model, stream, nil
}

func copyResponseHeaders(dst, src http.Header) {
	for k, vals := range src {
		switch strings.ToLower(k) {
		case "content-length", "content-type", "content-encoding", "connection", "transfer-encoding", "keep-alive":
			continue
		}
		for _, v := range vals {
			dst.Add(k, v)
		}
	}
}

func ensureRequestID(r *http.Request) string {
	if r != nil {
		if id := middleware.GetReqID(r.Context()); id != "" {
			return id
		}
	}
	var buf [8]byte
	if _, err := rand.Read(buf[:]); err == nil {
		return "req_" + hex.EncodeToString(buf[:])
	}
	return fmt.Sprintf("req_%d", time.Now().UnixNano())
}

func randomResponseID() string {
	var buf [8]byte
	if _, err := rand.Read(buf[:]); err == nil {
		return "resp_" + hex.EncodeToString(buf[:])
	}
	return fmt.Sprintf("resp_%d", time.Now().UnixNano())
}

func buildCompactedResponse(canon *CanonicalChatResponse) map[string]any {
	id := canon.ID
	if id == "" {
		id = randomResponseID()
	}
	var text strings.Builder
	for _, b := range canon.Content {
		if b.Type == BlockTypeText && b.Text != "" {
			text.WriteString(b.Text)
		}
	}
	return map[string]any{
		"id":     id,
		"object": "response.compaction",
		"output": []map[string]any{{
			"type": "message",
			"role": "assistant",
			"content": []map[string]any{{
				"type": "output_text",
				"text": text.String(),
			}},
		}},
	}
}
