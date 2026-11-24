package combos

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"sort"
	"strings"
	"time"

	"github.com/omeg4monx/router-gateway/internal/contracts"
	"github.com/omeg4monx/router-gateway/internal/platform/httpclient"
	"github.com/omeg4monx/router-gateway/internal/platform/httpserver"
)

const (
	defaultProbeTimeout    = 20 * time.Second
	maxProbeResponseBytes  = 64 << 10 // 64 KiB
	maxBatchTestModels     = 50
	slowModelThresholdMS   = 5000
	maxProbeErrorSummaryCh = 280
)

// CatalogModel describes a single model entry in OmniRoute's /api/models/catalog response.
type CatalogModel struct {
	ID                 string          `json:"id"`
	Name               string          `json:"name"`
	Type               string          `json:"type"`
	Subtype            string          `json:"subtype,omitempty"`
	Custom             bool            `json:"custom,omitempty"`
	Free               bool            `json:"free,omitempty"`
	Capabilities       map[string]bool `json:"capabilities,omitempty"`
	ContextLength      int             `json:"context_length,omitempty"`
	MaxOutputTokens    int             `json:"max_output_tokens,omitempty"`
	InputModalities    []string        `json:"input_modalities,omitempty"`
	OutputModalities   []string        `json:"output_modalities,omitempty"`
	SupportedEndpoints []string        `json:"supported_endpoints,omitempty"`
	Hidden             bool            `json:"hidden,omitempty"`
}

// CatalogProviderGroup groups models under a single provider ID.
type CatalogProviderGroup struct {
	Provider string         `json:"provider"`
	Models   []CatalogModel `json:"models"`
}

type modelCatalogResponse struct {
	Catalog map[string]CatalogProviderGroup `json:"catalog"`
}

type testModelRequest struct {
	ProviderID string `json:"providerId"`
	ModelID    string `json:"modelId"`
}

// ModelTestResult matches OmniRoute's /api/models/test and /api/models/test-all item payload.
type ModelTestResult struct {
	OK          bool   `json:"ok"`
	Status      string `json:"status"` // "ok" | "slow" | "error"
	LatencyMS   int64  `json:"latencyMs"`
	StatusCode  int    `json:"statusCode,omitempty"`
	Error       string `json:"error,omitempty"`
	RateLimited bool   `json:"rateLimited,omitempty"`
	IsQuota     bool   `json:"isQuota,omitempty"`
	IsTimeout   bool   `json:"isTimeout,omitempty"`
}

type testAllModelsRequest struct {
	ProviderID       string   `json:"providerId"`
	ModelIDs         []string `json:"modelIds"`
	RespectRateLimit bool     `json:"respectRateLimit"`
}

type testAllModelsResponse struct {
	Results map[string]ModelTestResult `json:"results"`
}

type testComboRequest struct {
	ComboName string `json:"comboName"`
}

// ComboTargetTestResult reports the probe outcome of a single target in a combo chain.
type ComboTargetTestResult struct {
	Target     string `json:"target"`
	Status     string `json:"status"` // "ok" | "error"
	LatencyMS  int64  `json:"latencyMs,omitempty"`
	StatusCode int    `json:"statusCode,omitempty"`
	Error      string `json:"error,omitempty"`
}

type testComboResponse struct {
	ResolvedBy string                  `json:"resolvedBy,omitempty"`
	Results    []ComboTargetTestResult `json:"results"`
	Error      string                  `json:"error,omitempty"`
}

// ProviderHealthMatrixEntry summarizes live connection health for a provider.
type ProviderHealthMatrixEntry struct {
	Provider               string `json:"provider"`
	State                  string `json:"state"` // "healthy" | "degraded" | "down" | "untested"
	TotalConnections       int    `json:"totalConnections"`
	HealthyConnections     int    `json:"healthyConnections"`
	CoolingDownConnections int    `json:"coolingDownConnections"`
}

type providerHealthMatrixResponse struct {
	Providers []ProviderHealthMatrixEntry `json:"providers"`
}

func (m *Module) handleModelCatalog(w http.ResponseWriter, r *http.Request) {
	includeHiddenParam := strings.TrimSpace(strings.ToLower(r.URL.Query().Get("includeHidden")))
	includeHidden := includeHiddenParam == "1" || includeHiddenParam == "true"

	conns, err := m.providers.ListEnabledConnections(r.Context())
	if err != nil {
		httpserver.WriteError(w, r, http.StatusInternalServerError, "internal_error", "failed to load provider connections")
		return
	}

	disabledList, err := m.repo.ListDisabledModels(r.Context())
	if err != nil {
		httpserver.WriteError(w, r, http.StatusInternalServerError, "internal_error", "failed to load disabled models")
		return
	}
	disabledSet := make(map[string]struct{}, len(disabledList))
	for _, d := range disabledList {
		disabledSet[strings.TrimSpace(d.Model)] = struct{}{}
	}

	type provAcc struct {
		name     string
		seen     map[string]struct{}
		models   []CatalogModel
		protocol string
		isCustom bool
	}
	groups := make(map[string]*provAcc)

	for _, conn := range conns {
		pid := strings.TrimSpace(conn.ProviderID)
		if pid == "" {
			continue
		}
		acc, ok := groups[pid]
		if !ok {
			acc = &provAcc{
				name:     providerDisplayName(pid, conn.Name),
				seen:     make(map[string]struct{}),
				protocol: conn.Protocol,
				isCustom: strings.HasPrefix(pid, "custom-") || strings.HasPrefix(pid, "node_"),
			}
			groups[pid] = acc
		}

		for _, rawModel := range conn.Models {
			trimmed := strings.TrimSpace(rawModel)
			if trimmed == "" {
				continue
			}
			modelID := trimmed
			if strings.HasPrefix(modelID, pid+"/") {
				modelID = strings.TrimPrefix(modelID, pid+"/")
			}
			if modelID == "" {
				continue
			}
			if _, dup := acc.seen[modelID]; dup {
				continue
			}

			_, bareHidden := disabledSet[modelID]
			_, qualHidden := disabledSet[pid+"/"+modelID]
			isHidden := bareHidden || qualHidden
			if isHidden && !includeHidden {
				continue
			}

			acc.seen[modelID] = struct{}{}
			entry := inferCatalogModel(pid, modelID, conn.Protocol, acc.isCustom)
			entry.Hidden = isHidden
			acc.models = append(acc.models, entry)
		}
	}

	catalog := make(map[string]CatalogProviderGroup, len(groups))
	for pid, acc := range groups {
		sort.Slice(acc.models, func(i, j int) bool {
			return acc.models[i].ID < acc.models[j].ID
		})
		catalog[pid] = CatalogProviderGroup{
			Provider: acc.name,
			Models:   acc.models,
		}
	}

	httpserver.WriteJSON(w, http.StatusOK, modelCatalogResponse{Catalog: catalog})
}

func (m *Module) handleProviderHealthMatrix(w http.ResponseWriter, r *http.Request) {
	conns, err := m.providers.ListEnabledConnections(r.Context())
	if err != nil {
		httpserver.WriteError(w, r, http.StatusInternalServerError, "internal_error", "failed to load provider connections")
		return
	}

	now := m.now().UTC()
	type agg struct {
		total   int
		healthy int
		cooling int
	}
	byProvider := make(map[string]*agg)
	for _, conn := range conns {
		pid := strings.TrimSpace(conn.ProviderID)
		if pid == "" {
			continue
		}
		a, ok := byProvider[pid]
		if !ok {
			a = &agg{}
			byProvider[pid] = a
		}
		a.total++
		if m.cooldown.IsInCooldown(conn.ID, now) {
			a.cooling++
		} else {
			a.healthy++
		}
	}

	keys := make([]string, 0, len(byProvider))
	for pid := range byProvider {
		keys = append(keys, pid)
	}
	sort.Strings(keys)

	entries := make([]ProviderHealthMatrixEntry, 0, len(keys))
	for _, pid := range keys {
		a := byProvider[pid]
		state := "untested"
		switch {
		case a.total == 0:
			state = "untested"
		case a.cooling == 0:
			state = "healthy"
		case a.healthy > 0:
			state = "degraded"
		default:
			state = "down"
		}
		entries = append(entries, ProviderHealthMatrixEntry{
			Provider:               pid,
			State:                  state,
			TotalConnections:       a.total,
			HealthyConnections:     a.healthy,
			CoolingDownConnections: a.cooling,
		})
	}

	httpserver.WriteJSON(w, http.StatusOK, providerHealthMatrixResponse{Providers: entries})
}

func (m *Module) handleTestModel(w http.ResponseWriter, r *http.Request) {
	var req testModelRequest
	if !decodeJSONBody(w, r, &req) {
		return
	}
	providerID := strings.TrimSpace(req.ProviderID)
	modelID := strings.TrimSpace(req.ModelID)
	if providerID == "" || modelID == "" {
		httpserver.WriteError(w, r, http.StatusBadRequest, "invalid_request", "providerId and modelId are required")
		return
	}

	res := m.probeProviderModel(r.Context(), providerID, modelID)
	httpserver.WriteJSON(w, http.StatusOK, res)
}

func (m *Module) handleTestAllModels(w http.ResponseWriter, r *http.Request) {
	var req testAllModelsRequest
	if !decodeJSONBody(w, r, &req) {
		return
	}
	providerID := strings.TrimSpace(req.ProviderID)
	if providerID == "" || len(req.ModelIDs) == 0 {
		httpserver.WriteError(w, r, http.StatusBadRequest, "invalid_request", "providerId and non-empty modelIds are required")
		return
	}

	limit := len(req.ModelIDs)
	if limit > maxBatchTestModels {
		limit = maxBatchTestModels
	}

	results := make(map[string]ModelTestResult, limit)
	rateLimited := false
	for i := 0; i < limit; i++ {
		mid := strings.TrimSpace(req.ModelIDs[i])
		if mid == "" {
			continue
		}
		if r.Context().Err() != nil {
			break
		}
		if rateLimited && req.RespectRateLimit {
			results[mid] = ModelTestResult{
				OK:          false,
				Status:      "error",
				StatusCode:  http.StatusTooManyRequests,
				Error:       "skipped due to upstream rate limit",
				RateLimited: true,
			}
			continue
		}
		res := m.probeProviderModel(r.Context(), providerID, mid)
		results[mid] = res
		if res.RateLimited {
			rateLimited = true
		}
	}

	httpserver.WriteJSON(w, http.StatusOK, testAllModelsResponse{Results: results})
}

func (m *Module) handleTestCombo(w http.ResponseWriter, r *http.Request) {
	var req testComboRequest
	if !decodeJSONBody(w, r, &req) {
		return
	}
	comboName := strings.TrimSpace(req.ComboName)
	if comboName == "" {
		httpserver.WriteError(w, r, http.StatusBadRequest, "invalid_request", "comboName is required")
		return
	}

	allCombos, err := m.repo.ListCombos(r.Context())
	if err != nil {
		httpserver.WriteError(w, r, http.StatusInternalServerError, "internal_error", "failed to load combos")
		return
	}
	var targetCombo *Combo
	for i := range allCombos {
		if allCombos[i].Name == comboName {
			targetCombo = &allCombos[i]
			break
		}
	}
	if targetCombo == nil {
		httpserver.WriteError(w, r, http.StatusNotFound, "combo_not_found", "combo not found")
		return
	}

	resp := testComboResponse{
		Results: make([]ComboTargetTestResult, 0, len(targetCombo.Targets)),
	}

	for _, target := range targetCombo.Targets {
		if r.Context().Err() != nil {
			break
		}
		providerID, modelID := splitTargetIdentifier(target)
		var probeRes ModelTestResult
		if providerID != "" {
			probeRes = m.probeProviderModel(r.Context(), providerID, modelID)
		} else {
			// Resolve bare model via resolver to pick best provider connection.
			resolved, resolveErr := m.Resolve(r.Context(), modelID, nil)
			if resolveErr != nil || len(resolved) == 0 {
				probeRes = ModelTestResult{
					OK:     false,
					Status: "error",
					Error:  "no enabled provider connection for target",
				}
			} else {
				probeRes = m.probeConnectionModel(r.Context(), resolved[0].ConnectionID, resolved[0].UpstreamModel)
			}
		}

		item := ComboTargetTestResult{
			Target:     target,
			Status:     "error",
			LatencyMS:  probeRes.LatencyMS,
			StatusCode: probeRes.StatusCode,
			Error:      probeRes.Error,
		}
		if probeRes.OK {
			item.Status = "ok"
			resp.Results = append(resp.Results, item)
			resp.ResolvedBy = target
			break
		}
		resp.Results = append(resp.Results, item)
	}

	if resp.ResolvedBy == "" && len(resp.Results) > 0 {
		resp.Error = resp.Results[len(resp.Results)-1].Error
	}

	httpserver.WriteJSON(w, http.StatusOK, resp)
}

func (m *Module) probeProviderModel(ctx context.Context, providerID, modelID string) ModelTestResult {
	conns, err := m.providers.ListEnabledConnections(ctx)
	if err != nil {
		return ModelTestResult{
			OK:     false,
			Status: "error",
			Error:  "failed to load provider connections",
		}
	}

	var matchedConn *contracts.ProviderConnection
	for i := range conns {
		if conns[i].ProviderID != providerID {
			continue
		}
		if matchedConn == nil || conns[i].Priority < matchedConn.Priority {
			c := conns[i]
			matchedConn = &c
		}
	}
	if matchedConn == nil {
		return ModelTestResult{
			OK:     false,
			Status: "error",
			Error:  fmt.Sprintf("no enabled connection found for provider %q", providerID),
		}
	}

	return m.probeConnectionModel(ctx, matchedConn.ID, modelID)
}

func (m *Module) probeConnectionModel(ctx context.Context, connectionID, modelID string) ModelTestResult {
	creds, err := m.providers.ResolveCredentials(ctx, connectionID)
	if err != nil {
		return ModelTestResult{
			OK:     false,
			Status: "error",
			Error:  "failed to resolve provider credentials",
		}
	}

	probeCtx, cancel := context.WithTimeout(ctx, defaultProbeTimeout)
	defer cancel()
	probeCtx = httpclient.WithProviderID(probeCtx, creds.ProviderID)

	req, err := buildModelProbeRequest(probeCtx, creds, modelID)
	if err != nil {
		return ModelTestResult{
			OK:     false,
			Status: "error",
			Error:  sanitizeProbeError(err.Error(), creds),
		}
	}

	client := m.httpClient
	if client == nil {
		client = &http.Client{Timeout: defaultProbeTimeout}
	}

	start := time.Now()
	resp, doErr := client.Do(req) //nolint:gosec // G704: upstream provider base_url is configured by operator
	latency := time.Since(start).Milliseconds()

	if doErr != nil {
		isTimeout := errors.Is(doErr, context.DeadlineExceeded) || strings.Contains(strings.ToLower(doErr.Error()), "timeout")
		return ModelTestResult{
			OK:        false,
			Status:    "error",
			LatencyMS: latency,
			Error:     sanitizeProbeError(doErr.Error(), creds),
			IsTimeout: isTimeout,
		}
	}
	defer func() { _ = resp.Body.Close() }()

	body, _ := io.ReadAll(io.LimitReader(resp.Body, maxProbeResponseBytes))

	if resp.StatusCode >= http.StatusOK && resp.StatusCode < http.StatusMultipleChoices {
		m.cooldown.RecordSuccess(connectionID)
		status := "ok"
		if latency >= slowModelThresholdMS {
			status = "slow"
		}
		return ModelTestResult{
			OK:         true,
			Status:     status,
			LatencyMS:  latency,
			StatusCode: resp.StatusCode,
		}
	}

	errMsg := extractUpstreamProbeError(resp.StatusCode, body, creds)
	lowerErr := strings.ToLower(errMsg)
	rateLimited := resp.StatusCode == http.StatusTooManyRequests || strings.Contains(lowerErr, "rate limit") || strings.Contains(lowerErr, "rate_limit")
	isQuota := strings.Contains(lowerErr, "quota") || strings.Contains(lowerErr, "insufficient") || strings.Contains(lowerErr, "billing") || strings.Contains(lowerErr, "credit")
	isTimeout := resp.StatusCode == http.StatusRequestTimeout || resp.StatusCode == http.StatusGatewayTimeout || strings.Contains(lowerErr, "timeout")

	return ModelTestResult{
		OK:          false,
		Status:      "error",
		LatencyMS:   latency,
		StatusCode:  resp.StatusCode,
		Error:       errMsg,
		RateLimited: rateLimited,
		IsQuota:     isQuota,
		IsTimeout:   isTimeout,
	}
}

func buildModelProbeRequest(ctx context.Context, creds contracts.UpstreamCredentials, modelID string) (*http.Request, error) {
	baseURL := strings.TrimRight(strings.TrimSpace(creds.BaseURL), "/")
	if baseURL == "" {
		return nil, errors.New("provider base_url is empty")
	}
	parsed, err := url.Parse(baseURL)
	if err != nil || (parsed.Scheme != "http" && parsed.Scheme != "https") || parsed.Host == "" {
		return nil, errors.New("invalid provider base_url")
	}

	token := strings.TrimSpace(creds.AccessToken)
	if token == "" {
		token = strings.TrimSpace(creds.APIKey)
	}

	var targetURL string
	var payload map[string]any

	switch strings.ToLower(strings.TrimSpace(creds.Protocol)) {
	case "anthropic":
		if strings.HasSuffix(baseURL, "/messages") {
			targetURL = baseURL
		} else {
			targetURL = baseURL + "/messages"
		}
		payload = map[string]any{
			"model":      modelID,
			"max_tokens": 1,
			"messages":   []map[string]string{{"role": "user", "content": "hi"}},
		}
	case "gemini":
		targetURL = fmt.Sprintf("%s/models/%s:generateContent", baseURL, url.PathEscape(modelID))
		payload = map[string]any{
			"contents":         []map[string]any{{"role": "user", "parts": []map[string]string{{"text": "hi"}}}},
			"generationConfig": map[string]int{"maxOutputTokens": 1},
		}
	default:
		if strings.HasSuffix(baseURL, "/chat/completions") || strings.HasSuffix(baseURL, "/responses") {
			targetURL = baseURL
		} else {
			targetURL = baseURL + "/chat/completions"
		}
		payload = map[string]any{
			"model":      modelID,
			"max_tokens": 1,
			"messages":   []map[string]string{{"role": "user", "content": "hi"}},
		}
	}

	rawBody, err := json.Marshal(payload)
	if err != nil {
		return nil, err
	}
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, targetURL, bytes.NewReader(rawBody))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")

	switch strings.ToLower(strings.TrimSpace(creds.Protocol)) {
	case "anthropic":
		if token != "" {
			if strings.EqualFold(creds.AuthType, "oauth") {
				req.Header.Set("Authorization", "Bearer "+token)
			} else {
				req.Header.Set("x-api-key", token)
			}
		}
		req.Header.Set("anthropic-version", "2023-06-01")
	case "gemini":
		if token != "" {
			if strings.EqualFold(creds.AuthType, "oauth") {
				req.Header.Set("Authorization", "Bearer "+token)
			} else {
				req.Header.Set("x-goog-api-key", token)
			}
		}
	default:
		if token != "" {
			req.Header.Set("Authorization", "Bearer "+token)
		}
	}

	for k, v := range creds.ExtraHeaders {
		if strings.TrimSpace(k) != "" {
			req.Header.Set(k, v)
		}
	}
	return req, nil
}

func extractUpstreamProbeError(statusCode int, body []byte, creds contracts.UpstreamCredentials) string {
	var parsed struct {
		Error any    `json:"error"`
		Msg   string `json:"message"`
	}
	var msg string
	if json.Unmarshal(body, &parsed) == nil {
		switch e := parsed.Error.(type) {
		case string:
			msg = e
		case map[string]any:
			if m, ok := e["message"].(string); ok {
				msg = m
			}
		}
		if msg == "" {
			msg = parsed.Msg
		}
	}
	if strings.TrimSpace(msg) == "" {
		msg = strings.TrimSpace(string(body))
	}
	if msg == "" {
		msg = fmt.Sprintf("HTTP %d %s", statusCode, http.StatusText(statusCode))
	} else {
		msg = fmt.Sprintf("HTTP %d: %s", statusCode, msg)
	}
	return sanitizeProbeError(msg, creds)
}

func sanitizeProbeError(msg string, creds contracts.UpstreamCredentials) string {
	out := msg
	if creds.APIKey != "" {
		out = strings.ReplaceAll(out, creds.APIKey, "[REDACTED]")
	}
	if creds.AccessToken != "" {
		out = strings.ReplaceAll(out, creds.AccessToken, "[REDACTED]")
	}
	out = strings.TrimSpace(out)
	if len(out) > maxProbeErrorSummaryCh {
		out = out[:maxProbeErrorSummaryCh] + "..."
	}
	return out
}

func splitTargetIdentifier(target string) (providerID, modelID string) {
	trimmed := strings.TrimSpace(target)
	if idx := strings.Index(trimmed, "/"); idx > 0 {
		return trimmed[:idx], trimmed[idx+1:]
	}
	return "", trimmed
}

func providerDisplayName(providerID, connName string) string {
	if trimmed := strings.TrimSpace(providerID); trimmed != "" {
		return trimmed
	}
	return strings.TrimSpace(connName)
}

func inferCatalogModel(providerID, modelID, protocol string, isCustom bool) CatalogModel {
	lower := strings.ToLower(modelID)
	lowerProv := strings.ToLower(providerID)

	modelType := "chat"
	subtype := "general"
	ctxLen := 128000
	maxOut := 8192
	caps := map[string]bool{
		"streaming":    true,
		"tool_calling": true,
	}
	inMods := []string{"text"}
	outMods := []string{"text"}
	endpoints := []string{"/v1/chat/completions"}
	if strings.EqualFold(protocol, "anthropic") {
		endpoints = []string{"/v1/messages", "/v1/chat/completions"}
	} else if strings.EqualFold(protocol, "gemini") {
		endpoints = []string{"/v1beta/models", "/v1/chat/completions"}
	}

	// Type inference
	switch {
	case strings.Contains(lower, "embed") || strings.Contains(lower, "bge-") || strings.Contains(lower, "e5-"):
		modelType = "embedding"
		subtype = "embedding"
		ctxLen = 8192
		maxOut = 0
		caps = map[string]bool{"embeddings": true}
		endpoints = []string{"/v1/embeddings"}
	case strings.Contains(lower, "rerank"):
		modelType = "rerank"
		subtype = "rerank"
		ctxLen = 8192
		maxOut = 0
		caps = map[string]bool{"rerank": true}
		endpoints = []string{"/v1/rerank"}
	case strings.Contains(lower, "dall-e") || strings.Contains(lower, "imagen") || strings.Contains(lower, "flux") || strings.Contains(lower, "stable-diffusion") || strings.Contains(lower, "gpt-image"):
		modelType = "image"
		subtype = "image-generation"
		ctxLen = 4096
		maxOut = 0
		caps = map[string]bool{"image_generation": true}
		outMods = []string{"image"}
		endpoints = []string{"/v1/images/generations"}
	case strings.Contains(lower, "whisper") || strings.Contains(lower, "tts-") || strings.Contains(lower, "audio"):
		modelType = "audio"
		subtype = "audio"
		ctxLen = 16384
		maxOut = 4096
		caps = map[string]bool{"audio": true}
		inMods = []string{"audio", "text"}
		outMods = []string{"audio", "text"}
	default:
		// Chat model capability & context inference
		if strings.Contains(lower, "o1") || strings.Contains(lower, "o3") || strings.Contains(lower, "o4") ||
			strings.Contains(lower, "r1") || strings.Contains(lower, "thinking") || strings.Contains(lower, "reason") ||
			strings.Contains(lower, "qwq") || strings.Contains(lower, "claude-3-7") || strings.Contains(lower, "claude-sonnet-4") ||
			strings.Contains(lower, "claude-opus-4") || strings.Contains(lower, "gemini-2.5") || strings.Contains(lower, "gpt-5") {
			caps["reasoning"] = true
			subtype = "reasoning"
		}
		if strings.Contains(lower, "coder") || strings.Contains(lower, "codex") || strings.Contains(lower, "codestral") {
			caps["code"] = true
			if subtype == "general" {
				subtype = "code"
			}
		}
		if strings.Contains(lower, "4o") || strings.Contains(lower, "gpt-5") || strings.Contains(lower, "claude-3") ||
			strings.Contains(lower, "claude-sonnet-4") || strings.Contains(lower, "claude-opus-4") ||
			strings.Contains(lower, "gemini-2") || strings.Contains(lower, "gemini-1.5") || strings.Contains(lower, "vision") ||
			strings.Contains(lower, "-vl") {
			caps["vision"] = true
			inMods = []string{"text", "image"}
		}
		if strings.Contains(lower, "mini") || strings.Contains(lower, "flash") || strings.Contains(lower, "haiku") || strings.Contains(lower, "nano") {
			if subtype == "general" {
				subtype = "fast"
			}
		}

		// Context / max output heuristics
		switch {
		case strings.Contains(lower, "gemini-2") || strings.Contains(lower, "gemini-1.5"):
			ctxLen = 1048576
			maxOut = 65536
		case strings.Contains(lower, "claude-"):
			ctxLen = 200000
			maxOut = 16384
			if strings.Contains(lower, "sonnet-4") || strings.Contains(lower, "opus-4") || strings.Contains(lower, "3-7") {
				maxOut = 64000
			}
		case strings.Contains(lower, "gpt-5") || strings.Contains(lower, "gpt-4.1"):
			ctxLen = 1000000
			maxOut = 32768
		case strings.Contains(lower, "o3") || strings.Contains(lower, "o4"):
			ctxLen = 200000
			maxOut = 100000
		case strings.Contains(lower, "gpt-4o"):
			ctxLen = 128000
			maxOut = 16384
		case strings.Contains(lower, "deepseek"):
			ctxLen = 128000
			maxOut = 16384
		}
	}

	isFree := strings.Contains(lower, ":free") || strings.Contains(lower, "-free") ||
		strings.Contains(lowerProv, "ollama") || strings.Contains(lowerProv, "lmstudio") ||
		strings.Contains(lowerProv, "groq") || strings.Contains(lowerProv, "cerebras")

	return CatalogModel{
		ID:                 modelID,
		Name:               modelID,
		Type:               modelType,
		Subtype:            subtype,
		Custom:             isCustom,
		Free:               isFree,
		Capabilities:       caps,
		ContextLength:      ctxLen,
		MaxOutputTokens:    maxOut,
		InputModalities:    inMods,
		OutputModalities:   outMods,
		SupportedEndpoints: endpoints,
	}
}
