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
	"strconv"
	"strings"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"

	"github.com/omeg4monx/router-gateway/internal/contracts"
	"github.com/omeg4monx/router-gateway/internal/platform/db"
	"github.com/omeg4monx/router-gateway/internal/platform/eventbus"
	"github.com/omeg4monx/router-gateway/internal/platform/httpclient"
	"github.com/omeg4monx/router-gateway/internal/platform/httpserver"
	"github.com/omeg4monx/router-gateway/internal/platform/lifecycle"
	"github.com/omeg4monx/router-gateway/internal/platform/venice"
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
	DB                      *db.DB
	HTTPClient              *http.Client
	Keys                    contracts.KeyValidator
	Resolver                contracts.ProviderResolver
	Models                  contracts.ModelLister
	Providers               contracts.ProviderStore
	Health                  contracts.AccountHealthTracker
	Interceptors            []contracts.GatewayInterceptor
	Usage                   contracts.UsageRecorder
	Bus                     *eventbus.Bus
	Logger                  *slog.Logger
	ChainTTL                time.Duration
	MaxGlobalInFlight       int
	MaxConnInFlight         int
	QueueTimeout            time.Duration
	ResponseCacheMaxEntries int
	ResponseCacheTTL        time.Duration
	DisableResponseCache    bool
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
	bus          *eventbus.Bus
	logger       *slog.Logger
	chains       *chainStore
	agQuota      *antigravityQuotaMonitor
	veniceMgr    *venice.AttestationManager
	shaper       *trafficShaper
	coalescer    *responseCoalescer
}

var (
	_ lifecycle.Module             = (*Module)(nil)
	_ contracts.ProtocolTranslator = (*Module)(nil)
	_ contracts.ChatCompleter      = (*Module)(nil)
)

// SupportedProtocols returns the ordered list of gateway wire protocols.
func (m *Module) SupportedProtocols() []string {
	return []string{
		ProtocolOpenAIChat,
		ProtocolAnthropic,
		ProtocolOpenAIResponses,
		ProtocolGemini,
	}
}

// TranslateRequest converts a wire-format request payload from fromProto to
// toProto using the gateway's fast-path validator or Canonical IR pipeline.
func (m *Module) TranslateRequest(fromProto, toProto string, body []byte) ([]byte, error) {
	if !IsSupportedProtocol(fromProto) {
		return nil, fmt.Errorf("%w: %q", ErrUnsupportedProtocol, fromProto)
	}
	if !IsSupportedProtocol(toProto) {
		return nil, fmt.Errorf("%w: %q", ErrUnsupportedProtocol, toProto)
	}
	if CanFastPath(fromProto, toProto) {
		return RewriteModelField(body, fromProto, "")
	}
	canon, err := DecodeRequest(fromProto, body, "")
	if err != nil {
		return nil, err
	}
	return EncodeRequest(toProto, canon)
}

// TranslateResponse converts a non-streaming wire-format response payload from
// fromProto to toProto using the gateway's Canonical IR pipeline.
func (m *Module) TranslateResponse(fromProto, toProto string, body []byte) ([]byte, error) {
	if !IsSupportedProtocol(fromProto) {
		return nil, fmt.Errorf("%w: %q", ErrUnsupportedProtocol, fromProto)
	}
	if !IsSupportedProtocol(toProto) {
		return nil, fmt.Errorf("%w: %q", ErrUnsupportedProtocol, toProto)
	}
	if CanFastPath(fromProto, toProto) {
		trimmed := bytes.TrimSpace(body)
		if len(trimmed) == 0 || trimmed[0] != '{' || !json.Valid(trimmed) {
			return nil, fmt.Errorf("%w: response body must be a JSON object", ErrInvalidPayload)
		}
		return bytes.Clone(trimmed), nil
	}
	canon, err := DecodeResponse(fromProto, body)
	if err != nil {
		return nil, err
	}
	return EncodeResponse(toProto, canon)
}

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

	var coalescer *responseCoalescer
	if !opts.DisableResponseCache {
		coalescer = newResponseCoalescer(opts.ResponseCacheMaxEntries, opts.ResponseCacheTTL, nil)
	}

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
		bus:          opts.Bus,
		logger:       logger,
		chains:       newChainStore(opts.DB, opts.ChainTTL),
		agQuota:      newAntigravityQuotaMonitor(),
		veniceMgr:    venice.NewAttestationManager(nil),
		shaper:       newTrafficShaper(opts.MaxGlobalInFlight, opts.MaxConnInFlight, opts.QueueTimeout, nil),
		coalescer:    coalescer,
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
	mounts.Gateway.Get("/v1/models/*", m.handleGetModel)

	mounts.Gateway.Get("/v1beta/models", m.handleGeminiListModels)
	mounts.Gateway.Post("/v1beta/models/{modelAction}", m.handleGeminiGenerate)

	// Compatibility aliases.
	mounts.Gateway.Post("/codex/responses", m.handleOpenAIResponses)
	mounts.Gateway.Post("/responses", m.handleOpenAIResponses)
	mounts.Gateway.Get("/models", m.handleListModels)
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
	upCtx := httpclient.WithProviderID(ctx, cand.ProviderID)
	req, err := http.NewRequestWithContext(upCtx, http.MethodPost, targetURL, bytes.NewReader(rewritten))
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
	keyInfo, ok := m.authenticateRequest(w, r, inProto)
	if !ok {
		return
	}
	m.handleInferenceWithKey(w, r, keyInfo, inProto, pathModel, forceStream, isCompact)
}

func (m *Module) handleInferenceWithKey(w http.ResponseWriter, r *http.Request, keyInfo contracts.KeyInfo, inProto, pathModel string, forceStream, isCompact bool) {
	startTime := time.Now()
	reqID := ensureRequestID(r)

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

	releaseGlobal, globalErr := m.shaper.acquireGlobal(r.Context())
	if globalErr != nil {
		if r.Context().Err() != nil {
			return
		}
		w.Header().Set("Retry-After", "1")
		writeGatewayError(w, r, inProto, http.StatusServiceUnavailable, "server_error", "overloaded", globalErr.Error())
		m.emitUsage(r.Context(), contracts.UsageRecord{
			RequestID:    reqID,
			KeyID:        keyInfo.ID,
			Model:        icReq.Model,
			Protocol:     inProto,
			Latency:      time.Since(startTime),
			StatusCode:   http.StatusServiceUnavailable,
			ErrorSummary: globalErr.Error(),
			Stream:       icReq.Stream,
			Timestamp:    time.Now().UTC(),
		})
		return
	}
	defer releaseGlobal()

	if m.coalescer != nil && m.coalescer.isCacheableRequest(inProto, icReq.Stream, isCompact, icReq.Body) {
		cKey := m.coalescer.cacheKey(keyInfo.ID, inProto, icReq.Model, icReq.Body)
		cachedResp, cacheStatus, cErr := m.coalescer.do(r.Context(), cKey, func(callCtx context.Context) (cachedGatewayResponse, error) {
			return m.executeCoalescedInference(callCtx, r, reqID, keyInfo, icReq, inProto, isCompact, startTime)
		})
		if cErr != nil {
			if r.Context().Err() != nil {
				return
			}
			var alreadyHandled *handledGatewayError
			if errors.As(cErr, &alreadyHandled) {
				switch {
				case alreadyHandled.resolveErr != nil:
					writeResolutionError(w, r, inProto, alreadyHandled.resolveErr)
				case len(alreadyHandled.rawBody) > 0:
					w.Header().Set("Content-Type", "application/json; charset=utf-8")
					w.WriteHeader(alreadyHandled.statusCode)
					_, _ = w.Write(alreadyHandled.rawBody)
				default:
					writeGatewayError(w, r, inProto, alreadyHandled.statusCode, alreadyHandled.errType, alreadyHandled.errCode, alreadyHandled.errMsg)
				}
				return
			}
			writeGatewayError(w, r, inProto, http.StatusBadGateway, "server_error", "upstream_failed", cErr.Error())
			return
		}

		w.Header().Set("X-Router-Cache", cacheStatus)
		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		w.WriteHeader(cachedResp.statusCode)
		_, _ = w.Write(cachedResp.body)

		usageRec := contracts.UsageRecord{
			RequestID:    reqID,
			KeyID:        keyInfo.ID,
			Model:        icReq.Model,
			ProviderID:   cachedResp.providerID,
			ConnectionID: cachedResp.connectionID,
			Protocol:     inProto,
			Latency:      time.Since(startTime),
			TTFT:         time.Since(startTime),
			StatusCode:   cachedResp.statusCode,
			Stream:       false,
			Timestamp:    time.Now().UTC(),
		}
		if cacheStatus == cacheStatusMiss {
			usageRec.InputTokens = cachedResp.usage.InputTokens
			usageRec.OutputTokens = cachedResp.usage.OutputTokens
			usageRec.CachedTokens = cachedResp.usage.CachedTokens
		} else {
			usageRec.CachedTokens = cachedResp.usage.InputTokens
		}
		m.emitUsage(r.Context(), usageRec)
		return
	}

	_ = m.executeInferenceCandidates(r.Context(), w, r, reqID, keyInfo, icReq, inProto, isCompact, startTime, nil)
}

type handledGatewayError struct {
	statusCode int
	errType    string
	errCode    string
	errMsg     string
	rawBody    []byte
	resolveErr error
}

func (e *handledGatewayError) Error() string { return e.errMsg }

func (m *Module) executeCoalescedInference(
	ctx context.Context,
	r *http.Request,
	reqID string,
	keyInfo contracts.KeyInfo,
	icReq *contracts.CanonicalRequest,
	inProto string,
	isCompact bool,
	startTime time.Time,
) (cachedGatewayResponse, error) {
	var captured cachedGatewayResponse
	err := m.executeInferenceCandidates(ctx, nil, r, reqID, keyInfo, icReq, inProto, isCompact, startTime, &captured)
	if err != nil {
		return cachedGatewayResponse{}, err
	}
	return captured, nil
}

func (m *Module) executeInferenceCandidates(
	ctx context.Context,
	w http.ResponseWriter,
	r *http.Request,
	reqID string,
	keyInfo contracts.KeyInfo,
	icReq *contracts.CanonicalRequest,
	inProto string,
	isCompact bool,
	startTime time.Time,
	capture *cachedGatewayResponse,
) error {
	candidates, err := m.resolver.Resolve(ctx, icReq.Model, keyInfo.AllowedModels)
	if err != nil {
		if capture == nil {
			writeResolutionError(w, r, inProto, err)
		}
		return &handledGatewayError{
			statusCode: http.StatusNotFound,
			errType:    "invalid_request_error",
			errCode:    "model_not_found",
			errMsg:     err.Error(),
			resolveErr: err,
		}
	}
	candidates = prioritizeHeadroomAffinity(candidates, icReq.Metadata)

	var (
		lastStatus      = http.StatusBadGateway
		lastMsg         = "upstream request failed"
		lastCand        contracts.ResolvedTarget
		singleCandidate = len(candidates) == 1
	)

	for _, cand := range candidates {
		if ctx.Err() != nil {
			return ctx.Err()
		}
		lastCand = cand

		var (
			creds   contracts.UpstreamCredentials
			credErr error
		)
		if venice.IsVirtualConnectionID(cand.ConnectionID) {
			creds = venice.DefaultCredentials()
			if strings.TrimSpace(cand.Endpoint) != "" {
				creds.BaseURL = cand.Endpoint
			}
		} else {
			creds, credErr = m.providers.ResolveCredentials(ctx, cand.ConnectionID)
		}
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

		upBody, chainInputMsgs, prepErr := m.prepareUpstreamBody(ctx, inProto, outProto, keyInfo.ID, upstreamModel, icReq.Stream, isCompact, icReq.Body)
		if prepErr != nil {
			if capture == nil {
				writeGatewayError(w, r, inProto, http.StatusBadRequest, "invalid_request_error", "invalid_payload", prepErr.Error())
			}
			return &handledGatewayError{
				statusCode: http.StatusBadRequest,
				errType:    "invalid_request_error",
				errCode:    "invalid_payload",
				errMsg:     prepErr.Error(),
			}
		}

		refreshed := false
		transientRetries := 0
		maxRetries, hasMaxRetries := parseComboMetadataInt(cand.Metadata, "combo_max_retries")
		retryDelayMs, hasRetryDelay := parseComboMetadataInt(cand.Metadata, "combo_retry_delay_ms")
		timeoutMs, hasTimeout := parseComboMetadataInt(cand.Metadata, "combo_timeout_ms")
		stopCascade := false

		for {
			releaseConn, connErr := m.shaper.acquireConnection(ctx, cand.ConnectionID)
			if connErr != nil {
				if ctx.Err() != nil {
					return ctx.Err()
				}
				lastStatus = http.StatusServiceUnavailable
				lastMsg = connErr.Error()
				break
			}

			canRetryTransient := singleCandidate && transientRetries < singleCandidateRetryMax
			if hasMaxRetries {
				canRetryTransient = transientRetries < maxRetries
			}

			attemptCtx := ctx
			var attemptCancel context.CancelFunc
			if hasTimeout && timeoutMs >= 1000 {
				attemptCtx, attemptCancel = context.WithTimeout(ctx, time.Duration(timeoutMs)*time.Millisecond)
			}

			if m.bus != nil {
				m.bus.Publish("traffic.started", map[string]any{
					"request_id":    reqID,
					"provider_id":   cand.ProviderID,
					"connection_id": cand.ConnectionID,
					"model":         upstreamModel,
					"protocol":      inProto,
					"stream":        icReq.Stream,
					"started_at":    startTime.UTC().Format(time.RFC3339Nano),
				})
			}
			if m.logger != nil {
				m.logger.Info("▶ gateway request",
					"request_id", reqID,
					"provider", cand.ProviderID,
					"model", upstreamModel,
					"protocol", inProto,
					"stream", icReq.Stream,
				)
			}

			outcome := m.executeCandidateAttempt(
				attemptCtx,
				w,
				r,
				reqID,
				keyInfo.ID,
				icReq.Model,
				icReq.Metadata,
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
				canRetryTransient,
				capture,
			)
			attemptTimedOut := attemptCtx.Err() == context.DeadlineExceeded && ctx.Err() == nil
			if attemptCancel != nil {
				attemptCancel()
			}
			releaseConn()

			if outcome.clientCanceled && attemptTimedOut {
				outcome = attemptOutcome{
					statusCode:         http.StatusGatewayTimeout,
					errMsg:             "combo upstream attempt timed out",
					deferredHealthFail: true,
				}
			}

			if outcome.handled {
				if capture != nil && outcome.statusCode >= 400 {
					errType := outcome.errType
					if errType == "" {
						errType = "invalid_request_error"
					}
					errCode := outcome.errCode
					if errCode == "" {
						errCode = "upstream_client_error"
					}
					return &handledGatewayError{
						statusCode: outcome.statusCode,
						errType:    errType,
						errCode:    errCode,
						errMsg:     outcome.errMsg,
						rawBody:    outcome.rawBody,
					}
				}
				return nil
			}
			if outcome.needsRefresh {
				refreshed = true
				newCreds, refErr := m.providers.RefreshOAuthToken(ctx, cand.ConnectionID)
				if refErr != nil {
					if ctx.Err() != nil || errors.Is(refErr, context.Canceled) {
						return ctx.Err()
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
				if m.bus != nil {
					m.bus.Publish("traffic.finished", map[string]any{
						"request_id":    reqID,
						"key_id":        keyInfo.ID,
						"provider_id":   cand.ProviderID,
						"connection_id": cand.ConnectionID,
						"model":         upstreamModel,
						"protocol":      inProto,
						"status_code":   499,
						"latency_ms":    max(time.Since(startTime).Milliseconds(), 0),
						"error_summary": "client canceled request",
						"stream":        icReq.Stream,
						"created_at":    time.Now().UTC().Format(time.RFC3339Nano),
					})
				}
				if m.logger != nil {
					m.logger.Warn("✗ client canceled",
						"request_id", reqID,
						"provider", cand.ProviderID,
						"model", upstreamModel,
						"status", 499,
						"latency_ms", max(time.Since(startTime).Milliseconds(), 0),
					)
				}
				return context.Canceled
			}
			if canRetryTransient && isTransientSingleCandidateStatus(outcome.statusCode) {
				var delay time.Duration
				if hasRetryDelay {
					delay = max(time.Duration(retryDelayMs)*time.Millisecond, outcome.retryAfter)
				} else {
					delay = computeJitteredRetryDelay(transientRetries, outcome.retryAfter)
				}
				transientRetries++
				if m.logger != nil {
					m.logger.Warn("↻ retrying upstream",
						"request_id", reqID,
						"provider", cand.ProviderID,
						"model", upstreamModel,
						"status", outcome.statusCode,
						"attempt", transientRetries,
					)
				}
				if delay > 0 {
					if err := sleepContext(ctx, delay); err != nil {
						return err
					}
				}
				continue
			}
			if outcome.deferredHealthFail {
				m.health.RecordFailure(cand.ConnectionID, outcome.statusCode, outcome.retryAfter)
			}
			lastStatus = outcome.statusCode
			lastMsg = outcome.errMsg
			if m.logger != nil {
				m.logger.Warn("✗ upstream failed",
					"request_id", reqID,
					"provider", cand.ProviderID,
					"model", upstreamModel,
					"status", outcome.statusCode,
					"error", outcome.errMsg,
				)
			}
			if cand.Metadata["combo_fallback_only_on_quota"] == "true" && outcome.statusCode != http.StatusTooManyRequests {
				stopCascade = true
			}
			break
		}
		if stopCascade {
			break
		}
	}

	m.emitUsage(ctx, contracts.UsageRecord{
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
	if capture == nil {
		writeGatewayError(w, r, inProto, lastStatus, "server_error", "upstream_failed", lastMsg)
	}
	return &handledGatewayError{
		statusCode: lastStatus,
		errType:    "server_error",
		errCode:    "upstream_failed",
		errMsg:     lastMsg,
	}
}

type attemptOutcome struct {
	handled            bool
	needsRefresh       bool
	clientCanceled     bool
	deferredHealthFail bool
	statusCode         int
	retryAfter         time.Duration
	errType            string
	errCode            string
	errMsg             string
	rawBody            []byte
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
	reqMeta map[string]any,
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
	deferHealthFailure bool,
	capture *cachedGatewayResponse,
) attemptOutcome {
	isVeniceKeyless := venice.IsVeniceKeyless(creds)
	var (
		upURL string
		err   error
	)
	if isVeniceKeyless {
		upURL = venice.ResolveInferenceURL(baseURL)
	} else {
		upURL, err = buildUpstreamURL(baseURL, outProto, upstreamModel, stream, false, isCompact)
		if err != nil {
			m.health.RecordFailure(cand.ConnectionID, http.StatusBadGateway, 0)
			return attemptOutcome{statusCode: http.StatusBadGateway, errMsg: err.Error()}
		}
	}

	// Proactive Antigravity quota pre-check: skip accounts whose family bucket
	// is already <= 2% remaining without burning a 429 on generateContent.
	if cand.ProviderID == "antigravity" && creds.AccessToken != "" && m.agQuota != nil {
		if drained, resetAfter := m.agQuota.checkDrained(ctx, m.httpClient, baseURL, creds.AccessToken, upstreamModel, cand.ConnectionID); drained {
			m.health.RecordFailure(cand.ConnectionID, http.StatusTooManyRequests, resetAfter)
			return attemptOutcome{
				statusCode: http.StatusTooManyRequests,
				errMsg:     "antigravity quota proactively drained",
			}
		}
	}

	upCtx := httpclient.WithProviderID(ctx, cand.ProviderID)
	reqClient := m.httpClient
	var (
		veniceAttestationToken string
		veniceProxyKey         string
		veniceRequestTools     []json.RawMessage
	)
	curBody := upBody
	if cand.ProviderID == "antigravity" {
		curBody = rewriteAntigravitySystemInstruction(curBody)
	}
	if cand.ProviderID == providerIDFreePI {
		curBody, _ = clampFreePIMaxTokens(curBody, resolveFreePIMinMaxTokens())
	}
	if isCodeBuddyProvider(cand.ProviderID) {
		curBody = transformCodeBuddyRequest(curBody)
	}
	if cand.ProviderID == providerIDJetBrains {
		if exceedsJetBrainsInputLimit(curBody) {
			m.emitUsage(ctx, contracts.UsageRecord{
				RequestID:    reqID,
				KeyID:        keyID,
				Model:        requestedModel,
				ProviderID:   cand.ProviderID,
				ConnectionID: cand.ConnectionID,
				Protocol:     inProto,
				Latency:      time.Since(startTime),
				StatusCode:   http.StatusBadRequest,
				ErrorSummary: jetbrainsContextLimitMessage,
				Stream:       stream,
				Timestamp:    time.Now().UTC(),
			})
			if capture == nil {
				writeGatewayError(w, r, inProto, http.StatusBadRequest, "invalid_request_error", "context_length_exceeded", jetbrainsContextLimitMessage)
			}
			return attemptOutcome{
				handled:    true,
				statusCode: http.StatusBadRequest,
				errType:    "invalid_request_error",
				errCode:    "context_length_exceeded",
				errMsg:     jetbrainsContextLimitMessage,
			}
		}
		curBody = transformJetBrainsRequest(curBody, stream)
	}
	if isVeniceKeyless {
		veniceRequestTools = venice.ExtractRequestTools(curBody)
		var payloadErr error
		curBody, _, payloadErr = venice.BuildChatPayload(upstreamModel, curBody)
		if payloadErr != nil {
			m.health.RecordFailure(cand.ConnectionID, http.StatusBadGateway, 0)
			return attemptOutcome{statusCode: http.StatusBadGateway, errMsg: payloadErr.Error()}
		}
		const maxVeniceMintAttempts = 3
		var attErr error
		for attempt := 0; attempt < maxVeniceMintAttempts; attempt++ {
			reqClient, veniceProxyKey = venice.ResolveSession(upCtx, m.httpClient, upURL)
			veniceAttestationToken, attErr = m.veniceMgr.GetTokenForProxy(upCtx, reqClient, upURL, veniceProxyKey)
			if attErr == nil {
				break
			}
			if ctx.Err() != nil {
				return attemptOutcome{clientCanceled: true}
			}
			m.veniceMgr.EvictForProxy(upURL, veniceProxyKey)
			if veniceProxyKey == venice.DirectProxyKey {
				break
			}
		}
		if attErr != nil {
			if !deferHealthFailure {
				m.health.RecordFailure(cand.ConnectionID, http.StatusBadGateway, 0)
			}
			return attemptOutcome{
				deferredHealthFail: deferHealthFailure,
				statusCode:         http.StatusBadGateway,
				errMsg:             attErr.Error(),
			}
		}
	}
	retriedSystemToUser := false
	retriedFreePIUpgrade := false
	retriedVeniceAuth := false

	var resp *http.Response
	var errBody []byte
	for {
		upReq, reqErr := http.NewRequestWithContext(upCtx, http.MethodPost, upURL, bytes.NewReader(curBody))
		if reqErr != nil {
			m.health.RecordFailure(cand.ConnectionID, http.StatusBadGateway, 0)
			return attemptOutcome{statusCode: http.StatusBadGateway, errMsg: reqErr.Error()}
		}
		if isVeniceKeyless {
			venice.ApplyHeaders(upReq, veniceAttestationToken, true)
		} else {
			applyUpstreamHeaders(upReq, creds, stream)
			switch {
			case cand.ProviderID == providerIDFreePI:
				applyFreePIHeaders(upReq)
			case isCodeBuddyProvider(cand.ProviderID):
				applyCodeBuddyHeaders(upReq, r, cand.ProviderID)
			case cand.ProviderID == providerIDJetBrains:
				applyJetBrainsHeaders(upReq, creds, stream)
			}
		}

		var doErr error
		resp, doErr = reqClient.Do(upReq)
		if doErr != nil {
			if isVeniceKeyless {
				m.veniceMgr.EvictForProxy(upURL, veniceProxyKey)
			}
			if ctx.Err() != nil {
				return attemptOutcome{clientCanceled: true}
			}
			if !deferHealthFailure {
				m.health.RecordFailure(cand.ConnectionID, http.StatusBadGateway, 0)
			}
			return attemptOutcome{
				deferredHealthFail: deferHealthFailure,
				statusCode:         http.StatusBadGateway,
				errMsg:             doErr.Error(),
			}
		}

		m.shaper.observeRateLimitHeaders(cand.ConnectionID, resp.Header)

		if isVeniceKeyless && !retriedVeniceAuth &&
			(resp.StatusCode == http.StatusUnauthorized || resp.StatusCode == http.StatusForbidden) {
			_, _ = io.Copy(io.Discard, io.LimitReader(resp.Body, maxErrorBodyBytes))
			_ = resp.Body.Close()
			m.veniceMgr.EvictForProxy(upURL, veniceProxyKey)
			reqClient, veniceProxyKey = venice.ResolveSession(upCtx, m.httpClient, upURL)
			freshToken, attErr := m.veniceMgr.GetTokenForProxy(upCtx, reqClient, upURL, veniceProxyKey)
			if attErr == nil && freshToken != "" {
				veniceAttestationToken = freshToken
				retriedVeniceAuth = true
				continue
			}
			if ctx.Err() != nil {
				return attemptOutcome{clientCanceled: true}
			}
			errMsg := "venice: empty attestation token"
			if attErr != nil {
				errMsg = attErr.Error()
			}
			if !deferHealthFailure {
				m.health.RecordFailure(cand.ConnectionID, http.StatusBadGateway, 0)
			}
			return attemptOutcome{
				deferredHealthFail: deferHealthFailure,
				statusCode:         http.StatusBadGateway,
				errMsg:             errMsg,
			}
		}

		if resp.StatusCode == http.StatusUpgradeRequired && cand.ProviderID == providerIDFreePI && !retriedFreePIUpgrade {
			errBody, _ = io.ReadAll(io.LimitReader(resp.Body, maxErrorBodyBytes))
			if _, _, lifted := maybeLiftFreePIClientVersion(errBody); lifted {
				_ = resp.Body.Close()
				retriedFreePIUpgrade = true
				errBody = nil
				continue
			}
		}

		if resp.StatusCode == http.StatusTooManyRequests && cand.ProviderID == "antigravity" {
			errBody, _ = io.ReadAll(io.LimitReader(resp.Body, maxErrorBodyBytes))
			if isAntigravityPromptRejection(resp.StatusCode, errBody) {
				if !retriedSystemToUser {
					if rewritten, ok := moveGeminiSystemInstructionToContents(curBody); ok {
						_ = resp.Body.Close()
						curBody = rewritten
						retriedSystemToUser = true
						errBody = nil
						continue
					}
				}
				_ = resp.Body.Close()
				m.emitUsage(ctx, contracts.UsageRecord{
					RequestID:    reqID,
					KeyID:        keyID,
					Model:        requestedModel,
					ProviderID:   cand.ProviderID,
					ConnectionID: cand.ConnectionID,
					Protocol:     inProto,
					Latency:      time.Since(startTime),
					StatusCode:   http.StatusBadRequest,
					ErrorSummary: antigravityPromptRejectedMessage,
					Stream:       stream,
					Timestamp:    time.Now().UTC(),
				})
				if capture == nil {
					writeGatewayError(w, r, inProto, http.StatusBadRequest, "invalid_request_error", "antigravity_prompt_rejected", antigravityPromptRejectedMessage)
				}
				return attemptOutcome{
					handled:    true,
					statusCode: http.StatusBadRequest,
					errType:    "invalid_request_error",
					errCode:    "antigravity_prompt_rejected",
					errMsg:     antigravityPromptRejectedMessage,
				}
			}
		}
		break
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 200 && resp.StatusCode < 300 {
		if cand.ProviderID == providerIDFreePI {
			replayedBody, pinMismatch := peekFreePIPinMismatch(resp.Body, freePIPinPeekLimit)
			if pinMismatch {
				if m.logger != nil {
					m.logger.Warn("freepi upstream opened with conversation-pin sentence; rotating account",
						"request_id", reqID,
						"connection_id", cand.ConnectionID,
					)
				}
				return attemptOutcome{
					statusCode: http.StatusForbidden,
					errMsg:     freePIPinMismatchErrorMessage,
				}
			}
			resp.Body = replayedBody
		}
		m.health.RecordSuccess(cand.ConnectionID)

		metaMap := make(map[string]any, len(reqMeta))
		for k, v := range reqMeta {
			metaMap[k] = v
		}
		meta := &contracts.ResponseMeta{
			RequestID:    reqID,
			StatusCode:   resp.StatusCode,
			Headers:      resp.Header.Clone(),
			ProviderID:   cand.ProviderID,
			ConnectionID: cand.ConnectionID,
			Model:        requestedModel,
			Metadata:     metaMap,
		}
		for _, ic := range m.interceptors {
			if icErr := ic.OnResponseHeaders(ctx, meta); icErr != nil {
				if capture == nil {
					writeGatewayError(w, r, inProto, http.StatusBadGateway, "server_error", "interceptor_error", icErr.Error())
				}
				return attemptOutcome{
					handled:    true,
					statusCode: http.StatusBadGateway,
					errType:    "server_error",
					errCode:    "interceptor_error",
					errMsg:     icErr.Error(),
				}
			}
		}
		if capture == nil {
			copyResponseHeaders(w.Header(), meta.Headers)
		}

		if stream {
			streamSrc := resp.Body
			if isVeniceKeyless {
				streamSrc = venice.NewNDJSONToOpenAISSEReaderWithTools(ctx, resp.Body, requestedModel, veniceRequestTools)
			}
			sres, streamErr := m.streamUpstreamSSE(ctx, w, streamSrc, reqID, inProto, outProto, requestedModel, startTime)
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
			if capture == nil {
				writeGatewayError(w, r, inProto, http.StatusBadGateway, "server_error", "upstream_read_error", readErr.Error())
			}
			return attemptOutcome{handled: true, statusCode: http.StatusBadGateway, errType: "server_error", errCode: "upstream_read_error", errMsg: readErr.Error()}
		}
		if isVeniceKeyless {
			converted, _, convErr := venice.ParseNDJSONToOpenAIWithTools(rawResp, requestedModel, veniceRequestTools)
			if convErr != nil {
				m.veniceMgr.EvictForProxy(upURL, veniceProxyKey)
				if capture == nil {
					writeGatewayError(w, r, inProto, http.StatusBadGateway, "server_error", "invalid_upstream_response", convErr.Error())
				}
				return attemptOutcome{handled: true, statusCode: http.StatusBadGateway, errType: "server_error", errCode: "invalid_upstream_response", errMsg: convErr.Error()}
			}
			rawResp = converted
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
				if capture == nil {
					writeGatewayError(w, r, inProto, http.StatusBadGateway, "server_error", "invalid_upstream_response", decErr.Error())
				}
				return attemptOutcome{handled: true, statusCode: http.StatusBadGateway, errType: "server_error", errCode: "invalid_upstream_response", errMsg: decErr.Error()}
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
					if capture == nil {
						writeGatewayError(w, r, inProto, http.StatusBadGateway, "server_error", "encode_error", encErr.Error())
					}
					return attemptOutcome{handled: true, statusCode: http.StatusBadGateway, errType: "server_error", errCode: "encode_error", errMsg: encErr.Error()}
				}
				outBytes = encoded
			} else {
				encoded, encErr := EncodeResponse(inProto, canon)
				if encErr != nil {
					if capture == nil {
						writeGatewayError(w, r, inProto, http.StatusBadGateway, "server_error", "encode_error", encErr.Error())
					}
					return attemptOutcome{handled: true, statusCode: http.StatusBadGateway, errType: "server_error", errCode: "encode_error", errMsg: encErr.Error()}
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

		for _, ic := range m.interceptors {
			rbi, ok := ic.(contracts.ResponseBodyInterceptor)
			if !ok {
				continue
			}
			nextBody, icErr := rbi.OnResponseBody(ctx, meta, outBytes)
			if icErr != nil {
				m.logger.Warn("response body interceptor failed", "interceptor", ic.Name(), "request_id", reqID, "err", icErr)
				continue
			}
			if len(nextBody) > 0 {
				outBytes = nextBody
			}
		}

		if capture != nil {
			*capture = cachedGatewayResponse{
				statusCode:   resp.StatusCode,
				body:         bytes.Clone(outBytes),
				usage:        usage,
				providerID:   cand.ProviderID,
				connectionID: cand.ConnectionID,
			}
			return attemptOutcome{handled: true, statusCode: resp.StatusCode}
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

	if isVeniceKeyless {
		m.veniceMgr.EvictForProxy(upURL, veniceProxyKey)
	}

	if resp.StatusCode == http.StatusUnauthorized && creds.AuthType == "oauth" && allowOAuthRefresh {
		return attemptOutcome{needsRefresh: true}
	}

	if errBody == nil {
		errBody, _ = io.ReadAll(io.LimitReader(resp.Body, maxErrorBodyBytes))
	}
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
			if capture == nil {
				w.Header().Set("Content-Type", "application/json; charset=utf-8")
				w.WriteHeader(resp.StatusCode)
				_, _ = w.Write(errBody)
			}
			return attemptOutcome{handled: true, statusCode: resp.StatusCode, errMsg: errMsg, rawBody: bytes.Clone(errBody)}
		}
		if capture == nil {
			writeGatewayError(w, r, inProto, resp.StatusCode, "invalid_request_error", "upstream_client_error", errMsg)
		}
		return attemptOutcome{handled: true, statusCode: resp.StatusCode, errType: "invalid_request_error", errCode: "upstream_client_error", errMsg: errMsg}
	}

	retryAfter := parseRetryAfter(resp.Header.Get("Retry-After"), time.Now().UTC())
	shouldDefer := deferHealthFailure && isTransientSingleCandidateStatus(resp.StatusCode)
	if !shouldDefer {
		m.health.RecordFailure(cand.ConnectionID, resp.StatusCode, retryAfter)
	}
	return attemptOutcome{
		deferredHealthFail: shouldDefer,
		statusCode:         resp.StatusCode,
		retryAfter:         retryAfter,
		errMsg:             errMsg,
	}
}

func (m *Module) emitUsage(ctx context.Context, rec contracts.UsageRecord) {
	if m.logger != nil {
		if rec.StatusCode >= http.StatusBadRequest || rec.ErrorSummary != "" {
			m.logger.Error("✗ gateway response",
				"request_id", rec.RequestID,
				"provider", rec.ProviderID,
				"model", rec.Model,
				"status", rec.StatusCode,
				"latency_ms", rec.Latency.Milliseconds(),
				"error", rec.ErrorSummary,
			)
		} else {
			m.logger.Info("✓ gateway response",
				"request_id", rec.RequestID,
				"provider", rec.ProviderID,
				"model", rec.Model,
				"status", rec.StatusCode,
				"latency_ms", rec.Latency.Milliseconds(),
				"in_tok", rec.InputTokens,
				"out_tok", rec.OutputTokens,
			)
		}
	}
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

// prioritizeHeadroomAffinity promotes the sticky Headroom prefix-cache candidate
// to index 0 while preserving fallback order for all remaining healthy candidates.
func prioritizeHeadroomAffinity(candidates []contracts.ResolvedTarget, metadata map[string]any) []contracts.ResolvedTarget {
	if len(candidates) <= 1 || len(metadata) == 0 {
		return candidates
	}
	if candidates[0].Metadata["combo_disable_stickiness"] == "true" {
		return candidates
	}
	connID, _ := metadata["headroom_affinity_connection"].(string)
	if connID == "" {
		return candidates
	}
	provID, _ := metadata["headroom_affinity_provider"].(string)

	matchIdx := -1
	for i, cand := range candidates {
		if cand.ConnectionID != connID {
			continue
		}
		if provID != "" && cand.ProviderID != provID {
			continue
		}
		matchIdx = i
		break
	}
	if matchIdx <= 0 {
		return candidates
	}

	ordered := make([]contracts.ResolvedTarget, 0, len(candidates))
	ordered = append(ordered, candidates[matchIdx])
	ordered = append(ordered, candidates[:matchIdx]...)
	ordered = append(ordered, candidates[matchIdx+1:]...)
	return ordered
}

func parseComboMetadataInt(meta map[string]string, key string) (int, bool) {
	if len(meta) == 0 {
		return 0, false
	}
	raw := strings.TrimSpace(meta[key])
	if raw == "" {
		return 0, false
	}
	val, err := strconv.Atoi(raw)
	if err != nil || val < 0 {
		return 0, false
	}
	return val, true
}
