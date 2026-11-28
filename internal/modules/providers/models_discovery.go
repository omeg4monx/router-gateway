package providers

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"math"
	"net/http"
	"net/url"
	"sort"
	"strconv"
	"strings"
	"time"
	"unicode/utf8"

	"github.com/go-chi/chi/v5"

	"github.com/omeg4monx/router-gateway/internal/contracts"
	"github.com/omeg4monx/router-gateway/internal/platform/httpclient"
	"github.com/omeg4monx/router-gateway/internal/platform/httpserver"
)

const (
	maxDiscoveryBodyBytes = 1 << 20 // 1 MiB
	codexModelsClientVer  = "0.167.1"
	geminiCLIModelsURL    = "https://cloudcode-pa.googleapis.com/v1internal:fetchAvailableModels"
)

type discoveredModel struct {
	ID     string `json:"id"`
	Name   string `json:"name,omitempty"`
	IsFree *bool  `json:"is_free,omitempty"`
}

type fetchConnectionModelsResponse struct {
	ConnectionID string            `json:"connection_id"`
	ProviderID   string            `json:"provider_id"`
	Models       []discoveredModel `json:"models"`
	Source       string            `json:"source"`
	LatencyMS    int64             `json:"latency_ms"`
	Warning      string            `json:"warning,omitempty"`
}

func (m *Module) handleFetchConnectionModels(w http.ResponseWriter, r *http.Request) {
	id := strings.TrimSpace(chi.URLParam(r, "id"))
	if id == "" {
		httpserver.WriteError(w, r, http.StatusBadRequest, "invalid_request", "connection id is required")
		return
	}

	creds, err := m.ResolveCredentials(r.Context(), id)
	if err != nil {
		if errors.Is(err, contracts.ErrConnectionNotFound) {
			httpserver.WriteError(w, r, http.StatusNotFound, "connection_not_found", "provider connection not found")
			return
		}
		httpserver.WriteError(w, r, http.StatusInternalServerError, "internal_error", "failed to load provider connection credentials")
		return
	}

	fallbackModels := make([]discoveredModel, 0)
	if summary, sumErr := m.repo.GetConnectionSummary(r.Context(), id); sumErr == nil {
		for _, mid := range summary.Models {
			trimmed := strings.TrimSpace(mid)
			if trimmed != "" {
				fallbackModels = append(fallbackModels, discoveredModel{ID: trimmed})
			}
		}
	}
	if len(fallbackModels) == 0 {
		if adapter, ok := m.GetAdapter(creds.ProviderID); ok {
			for _, mid := range adapter.DefaultModels() {
				trimmed := strings.TrimSpace(mid)
				if trimmed != "" {
					fallbackModels = append(fallbackModels, discoveredModel{ID: trimmed})
				}
			}
		}
	}

	reqCtx := httpclient.WithProviderID(r.Context(), creds.ProviderID)
	requests, reqErr := m.buildModelsDiscoveryRequests(reqCtx, creds)
	if reqErr != nil || len(requests) == 0 {
		warn := "unable to build model discovery request"
		if reqErr != nil {
			warn = truncateError(sanitizeUpstreamError(reqErr.Error(), creds))
		}
		httpserver.WriteJSON(w, http.StatusOK, fetchConnectionModelsResponse{
			ConnectionID: id,
			ProviderID:   creds.ProviderID,
			Models:       fallbackModels,
			Source:       "fallback",
			Warning:      warn,
		})
		return
	}

	start := time.Now()
	var lastWarning string
	for _, req := range requests {
		resp, doErr := m.httpClient.Do(req) //nolint:gosec // G704: upstream provider base_url is configured by operator and validated http/https
		if doErr != nil {
			lastWarning = truncateError(sanitizeUpstreamError(doErr.Error(), creds))
			continue
		}
		body, _ := io.ReadAll(io.LimitReader(resp.Body, maxDiscoveryBodyBytes))
		_ = resp.Body.Close()

		if resp.StatusCode < http.StatusOK || resp.StatusCode >= http.StatusMultipleChoices {
			msg := fmt.Sprintf("upstream returned HTTP %d", resp.StatusCode)
			if trimmed := strings.TrimSpace(string(body)); trimmed != "" {
				msg = msg + ": " + sanitizeUpstreamError(trimmed, creds)
			}
			lastWarning = truncateError(msg)
			continue
		}

		discovered := parseDiscoveredModels(body)
		if len(discovered) > 0 {
			m.persistDiscoveredModels(r.Context(), id, discovered)
			httpserver.WriteJSON(w, http.StatusOK, fetchConnectionModelsResponse{
				ConnectionID: id,
				ProviderID:   creds.ProviderID,
				Models:       discovered,
				Source:       "upstream",
				LatencyMS:    time.Since(start).Milliseconds(),
			})
			return
		}
		lastWarning = "upstream returned an empty model list"
	}

	httpserver.WriteJSON(w, http.StatusOK, fetchConnectionModelsResponse{
		ConnectionID: id,
		ProviderID:   creds.ProviderID,
		Models:       fallbackModels,
		Source:       "fallback",
		LatencyMS:    time.Since(start).Milliseconds(),
		Warning:      lastWarning,
	})
}

var errConnectionModelsUnchanged = errors.New("providers: connection models unchanged")

func mergeConnectionModels(existing []string, discovered []discoveredModel) ([]string, bool) {
	seen := make(map[string]struct{}, len(existing)+len(discovered))
	merged := make([]string, 0, len(existing)+len(discovered))

	appendIfValid := func(raw string) {
		if len(merged) >= maxSliceItems {
			return
		}
		trimmed := strings.TrimSpace(raw)
		if trimmed == "" || utf8.RuneCountInString(trimmed) > maxSliceItemRunes {
			return
		}
		if _, dup := seen[trimmed]; dup {
			return
		}
		seen[trimmed] = struct{}{}
		merged = append(merged, trimmed)
	}

	for _, id := range existing {
		appendIfValid(id)
	}
	for _, item := range discovered {
		appendIfValid(item.ID)
	}

	if len(merged) == len(existing) {
		same := true
		for i := range merged {
			if merged[i] != existing[i] {
				same = false
				break
			}
		}
		if same {
			return merged, false
		}
	}
	return merged, true
}

func (m *Module) persistDiscoveredModels(ctx context.Context, connectionID string, discovered []discoveredModel) {
	if m == nil || m.repo == nil || len(discovered) == 0 {
		return
	}
	summary, err := m.repo.UpdateConnection(ctx, connectionID, func(cur connectionRecord) (updateConnectionValues, error) {
		merged, changed := mergeConnectionModels(cur.Summary.Models, discovered)
		if !changed {
			return updateConnectionValues{}, errConnectionModelsUnchanged
		}
		return updateConnectionValues{
			Name:           cur.Summary.Name,
			BaseURL:        cur.Summary.BaseURL,
			CredentialsEnc: cur.CredentialsEnc,
			Priority:       cur.Summary.Priority,
			Enabled:        cur.Summary.Enabled,
			Models:         merged,
			Headers:        cur.Summary.Headers,
			TokenExpiresAt: cur.Summary.TokenExpiresAt,
		}, nil
	})
	if err != nil {
		return
	}
	m.publishChanged("updated", "connection", summary.ID, summary.ProviderID, actorFromContext(ctx), summary.UpdatedAt)
}

func (m *Module) buildModelsDiscoveryRequests(ctx context.Context, creds contracts.UpstreamCredentials) ([]*http.Request, error) {
	// Codex OAuth uses the Codex models endpoint with client_version.
	if creds.ProviderID == "codex" && creds.AuthType == authTypeOAuth && strings.Contains(creds.BaseURL, "chatgpt.com") {
		targetURL := strings.TrimRight(creds.BaseURL, "/") + "/models?client_version=" + codexModelsClientVer
		req, err := newAuthenticatedDiscoveryRequest(ctx, http.MethodGet, targetURL, nil, creds)
		if err != nil {
			return nil, err
		}
		return []*http.Request{req}, nil
	}

	// Antigravity / Gemini OAuth cloudcode-pa uses POST v1internal:fetchAvailableModels.
	if (creds.ProviderID == "antigravity" || creds.ProviderID == "gemini-oauth") &&
		creds.AuthType == authTypeOAuth &&
		strings.Contains(creds.BaseURL, "cloudcode-pa.googleapis.com") {
		req, err := newAuthenticatedDiscoveryRequest(ctx, http.MethodPost, geminiCLIModelsURL, []byte(`{}`), creds)
		if err != nil {
			return nil, err
		}
		req.Header.Set("Content-Type", "application/json")
		return []*http.Request{req}, nil
	}

	// JetBrains AI (Grazie) authenticates /models via Grazie-Authenticate-JWT + Grazie-Agent.
	if creds.ProviderID == providerIDJetBrains {
		targetURL := resolveJetBrainsModelsURL(creds.BaseURL)
		req, err := http.NewRequestWithContext(ctx, http.MethodGet, targetURL, http.NoBody)
		if err != nil {
			return nil, err
		}
		req.Header.Set("Accept", "application/json")
		req.Header.Set("Grazie-Agent", jetbrainsGrazieAgentHeader)
		jwt := strings.TrimSpace(creds.ExtraHeaders["Grazie-Authenticate-JWT"])
		if jwt == "" {
			jwt = strings.TrimSpace(creds.AccessToken)
		}
		if jwt == "" {
			jwt = strings.TrimSpace(creds.APIKey)
		}
		if jwt != "" {
			req.Header.Set("Grazie-Authenticate-JWT", jwt)
		}
		return []*http.Request{req}, nil
	}

	baseURL, err := validateHTTPURL(creds.BaseURL)
	if err != nil {
		return nil, fmt.Errorf("providers: invalid base_url: %w", err)
	}

	candidates := buildDiscoveryCandidateURLs(creds.ProviderID, creds.Protocol, baseURL)
	reqs := make([]*http.Request, 0, len(candidates))
	for _, targetURL := range candidates {
		req, buildErr := newAuthenticatedDiscoveryRequest(ctx, http.MethodGet, targetURL, nil, creds)
		if buildErr != nil {
			continue
		}
		reqs = append(reqs, req)
	}
	if len(reqs) == 0 {
		return nil, errors.New("providers: no valid discovery URL candidates")
	}
	return reqs, nil
}

func buildDiscoveryCandidateURLs(providerID, protocol, rawBaseURL string) []string {
	trimmed := strings.TrimRight(strings.TrimSpace(rawBaseURL), "/")
	if trimmed == "" {
		return nil
	}

	// If caller already configured a /models endpoint directly, use it as-is.
	if strings.HasSuffix(trimmed, "/models") {
		return []string{trimmed}
	}

	// Strip common inference suffixes so full endpoint URLs (e.g. /v1/chat/completions, /v1/messages)
	// resolve cleanly to their sibling /models discovery endpoint.
	for _, suffix := range []string{
		"/chat/completions",
		"/responses",
		"/messages",
		"/completions",
		"/inference/chat",
	} {
		if strings.HasSuffix(trimmed, suffix) {
			trimmed = strings.TrimSuffix(trimmed, suffix)
			break
		}
	}
	trimmed = strings.TrimRight(trimmed, "/")

	out := make([]string, 0, 3)
	add := func(u string) {
		for _, existing := range out {
			if existing == u {
				return
			}
		}
		out = append(out, u)
	}

	switch protocol {
	case protocolAnthropic:
		if strings.HasSuffix(trimmed, "/v1") {
			add(trimmed + "/models")
		} else {
			add(trimmed + "/v1/models")
			add(trimmed + "/models")
		}
	case protocolGemini:
		if strings.HasSuffix(trimmed, "/v1beta") || strings.HasSuffix(trimmed, "/v1") {
			add(trimmed + "/models")
		} else {
			add(trimmed + "/v1beta/models")
			add(trimmed + "/v1/models")
		}
	default:
		add(trimmed + "/models")
		if !strings.HasSuffix(trimmed, "/v1") && !strings.HasSuffix(trimmed, "/v1beta") {
			add(trimmed + "/v1/models")
		}
	}

	// Auto-detect Ollama native /api/tags fallback.
	if strings.Contains(strings.ToLower(providerID), "ollama") || strings.Contains(trimmed, ":11434") {
		if parsed, err := url.Parse(trimmed); err == nil && parsed.Scheme != "" && parsed.Host != "" {
			add(parsed.Scheme + "://" + parsed.Host + "/api/tags")
		}
	}

	return out
}

func newAuthenticatedDiscoveryRequest(
	ctx context.Context,
	method string,
	targetURL string,
	body []byte,
	creds contracts.UpstreamCredentials,
) (*http.Request, error) {
	var bodyReader io.Reader = http.NoBody
	if len(body) > 0 {
		bodyReader = bytes.NewReader(body)
	}
	req, err := http.NewRequestWithContext(ctx, method, targetURL, bodyReader)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Accept", "application/json")
	if creds.Protocol == protocolAnthropic {
		req.Header.Set("anthropic-version", defaultAnthropicVersion)
	}
	if creds.ProviderID == "github-copilot" {
		req.Header.Set("Editor-Version", "vscode/1.99.0")
		req.Header.Set("Copilot-Integration-Id", "vscode-chat")
	}
	for k, v := range creds.ExtraHeaders {
		cleanKey := strings.TrimSpace(k)
		if cleanKey == "" {
			continue
		}
		if creds.ProviderID == "antigravity" && strings.EqualFold(cleanKey, "x-goog-user-project") && isBlockedAntigravityProjectID(v) {
			continue
		}
		req.Header.Set(cleanKey, v)
	}

	switch creds.Protocol {
	case protocolAnthropic:
		if creds.AuthType == authTypeOAuth || creds.AccessToken != "" {
			if creds.AccessToken != "" {
				req.Header.Set("Authorization", "Bearer "+creds.AccessToken)
			}
		} else if creds.APIKey != "" {
			req.Header.Set("x-api-key", creds.APIKey)
		}
	case protocolGemini:
		if creds.AuthType == authTypeOAuth || creds.AccessToken != "" {
			if creds.AccessToken != "" {
				req.Header.Set("Authorization", "Bearer "+creds.AccessToken)
			}
		} else if creds.APIKey != "" {
			req.Header.Set("x-goog-api-key", creds.APIKey)
		}
	default:
		token := creds.APIKey
		if creds.AuthType == authTypeOAuth || creds.AccessToken != "" {
			token = creds.AccessToken
		}
		if token != "" {
			req.Header.Set("Authorization", "Bearer "+token)
		}
	}

	return req, nil
}

func parseDiscoveredModels(body []byte) []discoveredModel {
	trimmed := bytes.TrimSpace(body)
	if len(trimmed) == 0 {
		return nil
	}

	var parsed any
	if err := json.Unmarshal(trimmed, &parsed); err != nil {
		return nil
	}

	seen := make(map[string]struct{})
	out := make([]discoveredModel, 0, 16)
	add := func(rawID, rawName string, isFree *bool) {
		id := strings.TrimSpace(rawID)
		id = strings.TrimPrefix(id, "models/")
		if id == "" {
			return
		}
		if _, dup := seen[id]; dup {
			return
		}
		seen[id] = struct{}{}
		name := strings.TrimSpace(rawName)
		if name == id {
			name = ""
		}
		out = append(out, discoveredModel{ID: id, Name: name, IsFree: isFree})
	}

	extractList := func(items []any) {
		for _, item := range items {
			switch v := item.(type) {
			case string:
				add(v, "", nil)
			case map[string]any:
				if internal, ok := v["isInternal"].(bool); ok && internal {
					continue
				}
				id := firstStringField(v, "id", "slug", "model", "name")
				name := firstStringField(v, "displayName", "display_name", "name")
				add(id, name, extractDiscoveredModelIsFree(v))
			}
		}
	}

	switch root := parsed.(type) {
	case []any:
		extractList(root)
	case map[string]any:
		if dataList, ok := root["data"].([]any); ok {
			extractList(dataList)
		} else if modelsList, ok := root["models"].([]any); ok {
			extractList(modelsList)
		} else if resultsList, ok := root["results"].([]any); ok {
			extractList(resultsList)
		} else if modelsMap, ok := root["models"].(map[string]any); ok {
			keys := make([]string, 0, len(modelsMap))
			for k := range modelsMap {
				keys = append(keys, k)
			}
			sort.Strings(keys)
			for _, k := range keys {
				info, _ := modelsMap[k].(map[string]any)
				if info != nil {
					if internal, ok := info["isInternal"].(bool); ok && internal {
						continue
					}
					add(k, firstStringField(info, "displayName", "display_name", "name"), extractDiscoveredModelIsFree(info))
				} else {
					add(k, "", nil)
				}
			}
		}
	}

	return out
}

func extractDiscoveredModelIsFree(m map[string]any) *bool {
	for _, key := range []string{"is_free", "isFree", "free"} {
		if b, ok := m[key].(bool); ok {
			return &b
		}
	}
	if tier, ok := m["tier"].(string); ok && strings.EqualFold(strings.TrimSpace(tier), "free") {
		v := true
		return &v
	}
	if pricing, ok := m["pricing"].(map[string]any); ok && pricing != nil {
		prompt, hasPrompt := parseOptionalPrice(pricing["prompt"])
		completion, hasCompletion := parseOptionalPrice(pricing["completion"])
		if hasPrompt && hasCompletion {
			v := prompt <= 0 && completion <= 0
			return &v
		}
	}
	return nil
}

func parseOptionalPrice(raw any) (price float64, ok bool) {
	switch v := raw.(type) {
	case float64:
		if !math.IsNaN(v) && !math.IsInf(v, 0) {
			return v, true
		}
	case string:
		trimmed := strings.TrimSpace(v)
		if trimmed == "" {
			return 0, false
		}
		if f, err := strconv.ParseFloat(trimmed, 64); err == nil && !math.IsNaN(f) && !math.IsInf(f, 0) {
			return f, true
		}
	}
	return 0, false
}

func firstStringField(m map[string]any, keys ...string) string {
	for _, k := range keys {
		if v, ok := m[k].(string); ok {
			if trimmed := strings.TrimSpace(v); trimmed != "" {
				return trimmed
			}
		}
	}
	return ""
}
