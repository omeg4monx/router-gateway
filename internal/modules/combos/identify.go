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
	"strings"
	"time"

	"github.com/omeg4monx/router-gateway/internal/contracts"
	"github.com/omeg4monx/router-gateway/internal/platform/httpclient"
	"github.com/omeg4monx/router-gateway/internal/platform/httpserver"
	"github.com/omeg4monx/router-gateway/internal/platform/venice"
)

const (
	maxIdentifyResponseBytes = 32 * 1024
	maxSelfReportRunes       = 280
	identityProbePrompt      = "Identify yourself concisely in one sentence: state your exact model name, model family, developer/organization, and knowledge cutoff date."
	identityLimitationNote   = "Detection combines upstream response metadata, vendor headers, and a neutral self-identification probe. Proxies that rewrite all metadata/headers and inject strict system prompts can still disguise the underlying model, so treat results as indicative evidence."
)

// ModelIdentityResult describes the outcome of an upstream model unmasking probe.
type ModelIdentityResult struct {
	OK                bool              `json:"ok"`
	ClaimedModel      string            `json:"claimedModel"`
	ClaimedFamily     string            `json:"claimedFamily,omitempty"`
	Verdict           string            `json:"verdict"`
	Confidence        string            `json:"confidence"`
	DetectedModel     string            `json:"detectedModel,omitempty"`
	DetectedFamily    string            `json:"detectedFamily,omitempty"`
	ReportedModel     string            `json:"reportedModel,omitempty"`
	ResponseIDPrefix  string            `json:"responseIdPrefix,omitempty"`
	SystemFingerprint string            `json:"systemFingerprint,omitempty"`
	VendorHeaders     map[string]string `json:"vendorHeaders,omitempty"`
	SelfReport        string            `json:"selfReport,omitempty"`
	Signals           []string          `json:"signals"`
	LimitationNote    string            `json:"limitationNote"`
	LatencyMS         int64             `json:"latencyMs"`
	StatusCode        int               `json:"statusCode,omitempty"`
	Error             string            `json:"error,omitempty"`
}

func (m *Module) handleIdentifyModel(w http.ResponseWriter, r *http.Request) {
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

	res := m.identifyProviderModel(r.Context(), providerID, modelID)
	httpserver.WriteJSON(w, http.StatusOK, res)
}

func (m *Module) identifyProviderModel(ctx context.Context, providerID, modelID string) ModelIdentityResult {
	claimedFamily := inferModelFamily(modelID)
	baseError := func(errMsg string, statusCode int, latency int64) ModelIdentityResult {
		return ModelIdentityResult{
			OK:             false,
			ClaimedModel:   modelID,
			ClaimedFamily:  claimedFamily,
			Verdict:        "error",
			Confidence:     "low",
			Signals:        []string{},
			LimitationNote: identityLimitationNote,
			LatencyMS:      latency,
			StatusCode:     statusCode,
			Error:          errMsg,
		}
	}

	conns, err := m.providers.ListEnabledConnections(ctx)
	if err != nil {
		return baseError("failed to load provider connections", 0, 0)
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
	var (
		connectionID string
		creds        contracts.UpstreamCredentials
	)
	if matchedConn == nil {
		keylessCreds, ok := resolveKeylessProbeCredentials(providerID)
		if !ok {
			return baseError(fmt.Sprintf("no enabled connection found for provider %q", providerID), 0, 0)
		}
		creds = keylessCreds
	} else {
		connectionID = matchedConn.ID
		creds, err = m.providers.ResolveCredentials(ctx, matchedConn.ID)
		if err != nil {
			return baseError("failed to resolve provider credentials", 0, 0)
		}
	}

	probeCtx, cancel := context.WithTimeout(ctx, defaultProbeTimeout)
	defer cancel()
	probeCtx = httpclient.WithProviderID(probeCtx, creds.ProviderID)

	client := selectProbeHTTPClient(m.httpClient, creds)

	executeIdentify := func(activeCreds contracts.UpstreamCredentials) (int, http.Header, []byte, int64, error) {
		var (
			req            *http.Request
			reqErr         error
			veniceProxyKey string
		)
		reqClient := client
		isVenice := venice.IsVeniceKeyless(activeCreds)
		if isVenice {
			targetURL := venice.ResolveInferenceURL(activeCreds.BaseURL)
			reqClient, veniceProxyKey = venice.ResolveSession(probeCtx, m.httpClient, targetURL)
			attToken, attErr := venice.DefaultAttestationManager().GetTokenForProxy(probeCtx, reqClient, targetURL, veniceProxyKey)
			if attErr != nil {
				return 0, nil, nil, 0, attErr
			}
			inBody, _ := json.Marshal(map[string]any{
				"messages":    []map[string]string{{"role": "user", "content": identityProbePrompt}},
				"temperature": 0,
			})
			probeBody, _, buildErr := venice.BuildChatPayload(modelID, inBody)
			if buildErr != nil {
				return 0, nil, nil, 0, buildErr
			}
			req, reqErr = http.NewRequestWithContext(probeCtx, http.MethodPost, targetURL, bytes.NewReader(probeBody))
			if reqErr == nil {
				venice.ApplyHeaders(req, attToken, true)
			}
		} else {
			req, reqErr = buildIdentityProbeRequest(probeCtx, activeCreds, modelID)
		}
		if reqErr != nil {
			return 0, nil, nil, 0, reqErr
		}
		start := time.Now()
		resp, doErr := reqClient.Do(req)
		latency := time.Since(start).Milliseconds()
		if doErr != nil {
			if isVenice {
				venice.DefaultAttestationManager().EvictForProxy(activeCreds.BaseURL, veniceProxyKey)
			}
			return 0, nil, nil, latency, doErr
		}
		defer func() { _ = resp.Body.Close() }()
		body, _ := io.ReadAll(io.LimitReader(resp.Body, maxIdentifyResponseBytes))
		if isVenice {
			if resp.StatusCode < http.StatusOK || resp.StatusCode >= http.StatusMultipleChoices {
				venice.DefaultAttestationManager().EvictForProxy(activeCreds.BaseURL, veniceProxyKey)
			} else if converted, _, convErr := venice.ParseNDJSONToOpenAI(body, modelID); convErr == nil {
				body = converted
			}
		}
		return resp.StatusCode, resp.Header.Clone(), body, latency, nil
	}

	statusCode, headers, body, latency, probeErr := executeIdentify(creds)
	if probeErr == nil && connectionID != "" && strings.EqualFold(creds.AuthType, "oauth") &&
		(statusCode == http.StatusUnauthorized || statusCode == http.StatusForbidden) {
		if refreshed, refreshErr := m.providers.RefreshOAuthToken(probeCtx, connectionID); refreshErr == nil &&
			strings.TrimSpace(refreshed.AccessToken) != "" {
			creds = refreshed
			statusCode, headers, body, latency, probeErr = executeIdentify(creds)
		}
	}

	if probeErr != nil {
		errMsg := sanitizeProbeError(probeErr.Error(), creds)
		m.emitUsage(ctx, contracts.UsageRecord{
			KeyID:        probeUsageKeyID,
			Model:        modelID,
			ProviderID:   creds.ProviderID,
			ConnectionID: connectionID,
			Protocol:     creds.Protocol,
			Latency:      time.Duration(latency) * time.Millisecond,
			TTFT:         time.Duration(latency) * time.Millisecond,
			StatusCode:   0,
			ErrorSummary: errMsg,
			Stream:       false,
		})
		return baseError(errMsg, 0, latency)
	}
	if statusCode < http.StatusOK || statusCode >= http.StatusMultipleChoices {
		errMsg := extractUpstreamProbeError(statusCode, body, creds)
		m.emitUsage(ctx, contracts.UsageRecord{
			KeyID:        probeUsageKeyID,
			Model:        modelID,
			ProviderID:   creds.ProviderID,
			ConnectionID: connectionID,
			Protocol:     creds.Protocol,
			Latency:      time.Duration(latency) * time.Millisecond,
			TTFT:         time.Duration(latency) * time.Millisecond,
			StatusCode:   statusCode,
			ErrorSummary: errMsg,
			Stream:       false,
		})
		return baseError(errMsg, statusCode, latency)
	}

	inTok, outTok, cachedTok := extractProbeTokenUsage(creds.Protocol, body)
	m.emitUsage(ctx, contracts.UsageRecord{
		KeyID:        probeUsageKeyID,
		Model:        modelID,
		ProviderID:   creds.ProviderID,
		ConnectionID: connectionID,
		Protocol:     creds.Protocol,
		InputTokens:  inTok,
		OutputTokens: outTok,
		CachedTokens: cachedTok,
		Latency:      time.Duration(latency) * time.Millisecond,
		TTFT:         time.Duration(latency) * time.Millisecond,
		StatusCode:   statusCode,
		Stream:       false,
	})

	return analyzeIdentityProbeResponse(modelID, statusCode, latency, headers, body)
}

func buildIdentityProbeRequest(ctx context.Context, creds contracts.UpstreamCredentials, modelID string) (*http.Request, error) {
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
			"model":       modelID,
			"max_tokens":  96,
			"temperature": 0,
			"messages":    []map[string]string{{"role": "user", "content": identityProbePrompt}},
		}
	case "gemini":
		targetURL = fmt.Sprintf("%s/models/%s:generateContent", baseURL, url.PathEscape(modelID))
		payload = map[string]any{
			"contents":         []map[string]any{{"role": "user", "parts": []map[string]string{{"text": identityProbePrompt}}}},
			"generationConfig": map[string]any{"maxOutputTokens": 96, "temperature": 0},
		}
	case "openai-responses":
		if strings.HasSuffix(baseURL, "/responses") {
			targetURL = baseURL
		} else {
			targetURL = baseURL + "/responses"
		}
		payload = map[string]any{
			"model":             modelID,
			"input":             identityProbePrompt,
			"max_output_tokens": 96,
		}
	default:
		if strings.HasSuffix(baseURL, "/chat/completions") ||
			strings.HasSuffix(baseURL, "/responses") ||
			strings.HasSuffix(baseURL, "/inference/chat") {
			targetURL = baseURL
		} else {
			targetURL = baseURL + "/chat/completions"
		}
		payload = map[string]any{
			"model":       modelID,
			"max_tokens":  96,
			"temperature": 0,
			"messages":    []map[string]string{{"role": "user", "content": identityProbePrompt}},
		}
		if isOpenCodeFreeCreds(creds) {
			applyOpenCodeFreePayload(payload, identityProbePrompt)
		}
	}

	rawBody, err := marshalProbeJSON(payload)
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

	if isOpenCodeFreeCreds(creds) {
		applyOpenCodeFreeHeaders(req)
	}

	for k, v := range creds.ExtraHeaders {
		if strings.TrimSpace(k) != "" {
			req.Header.Set(k, v)
		}
	}
	return req, nil
}

func analyzeIdentityProbeResponse(
	claimedModel string,
	statusCode int,
	latencyMS int64,
	headers http.Header,
	body []byte,
) ModelIdentityResult {
	claimedFamily := inferModelFamily(claimedModel)
	vendorHeaders := extractVendorHeaders(headers)

	reportedModel, responseIDPrefix, systemFingerprint, selfReport, hasReasoningField := parseIdentityPayload(body)
	if reportedModel == "" {
		if hdrModel := strings.TrimSpace(vendorHeaders["openai-model"]); hdrModel != "" {
			reportedModel = hdrModel
		}
	}

	signals := make([]string, 0, 5)
	mismatchCount := 0
	highConfidenceMismatch := false
	confirmCount := 0
	detectedModel := ""
	detectedFamily := ""

	// Layer 1: Upstream response JSON model field vs claimed model.
	if reportedModel != "" {
		reportedFamily := inferModelFamily(reportedModel)
		if !modelsLooselyMatch(claimedModel, reportedModel) {
			mismatchCount++
			highConfidenceMismatch = true
			detectedModel = reportedModel
			if reportedFamily != "" {
				detectedFamily = reportedFamily
			}
			signals = append(signals, fmt.Sprintf("Response metadata reported model %q instead of %q", reportedModel, claimedModel))
		} else {
			confirmCount++
			signals = append(signals, fmt.Sprintf("Response metadata matched %q", reportedModel))
		}
	}

	// Layer 2: Vendor HTTP header & response ID prefix signals.
	if hdrModel := strings.TrimSpace(vendorHeaders["openai-model"]); hdrModel != "" {
		if !modelsLooselyMatch(claimedModel, hdrModel) {
			mismatchCount++
			highConfidenceMismatch = true
			if detectedModel == "" {
				detectedModel = hdrModel
			}
			signals = append(signals, fmt.Sprintf("Header openai-model reported %q", hdrModel))
		} else {
			confirmCount++
			signals = append(signals, fmt.Sprintf("Header openai-model confirmed %q", hdrModel))
		}
	}
	if responseIDPrefix != "" {
		if responseIDPrefix == "msg_" && claimedFamily != "" && claimedFamily != "anthropic" {
			signals = append(signals, "Response ID uses Anthropic msg_ prefix")
		} else if responseIDPrefix == "msg_" && claimedFamily == "anthropic" {
			confirmCount++
		}
	}
	if hasReasoningField {
		signals = append(signals, "Response included reasoning_content field")
	}

	// Layer 3: Self-identification reply family check.
	if selfReport != "" {
		selfFamily := inferModelFamily(selfReport)
		if selfFamily != "" {
			switch {
			case claimedFamily != "" && selfFamily != claimedFamily:
				mismatchCount++
				if detectedFamily == "" {
					detectedFamily = selfFamily
				}
				signals = append(signals, fmt.Sprintf("Self-identification reply pointed to family %q (claimed %q)", selfFamily, claimedFamily))
			case claimedFamily != "" && selfFamily == claimedFamily:
				confirmCount++
				signals = append(signals, fmt.Sprintf("Self-identification reply matched family %q", selfFamily))
			default:
				detectedFamily = selfFamily
				signals = append(signals, fmt.Sprintf("Self-identification reply indicated family %q", selfFamily))
			}
		}
	}

	verdict := "inconclusive"
	confidence := "low"
	switch {
	case mismatchCount > 0:
		verdict = "suspected_masked"
		if highConfidenceMismatch || mismatchCount >= 2 {
			confidence = "high"
		} else {
			confidence = "medium"
		}
	case confirmCount >= 2:
		verdict = "verified"
		confidence = "high"
		if detectedModel == "" {
			detectedModel = reportedModel
		}
		if detectedFamily == "" {
			detectedFamily = claimedFamily
		}
	case confirmCount == 1:
		verdict = "verified"
		confidence = "medium"
		if detectedModel == "" {
			detectedModel = reportedModel
		}
		if detectedFamily == "" {
			detectedFamily = claimedFamily
		}
	}

	return ModelIdentityResult{
		OK:                true,
		ClaimedModel:      claimedModel,
		ClaimedFamily:     claimedFamily,
		Verdict:           verdict,
		Confidence:        confidence,
		DetectedModel:     detectedModel,
		DetectedFamily:    detectedFamily,
		ReportedModel:     reportedModel,
		ResponseIDPrefix:  responseIDPrefix,
		SystemFingerprint: systemFingerprint,
		VendorHeaders:     vendorHeaders,
		SelfReport:        truncateRunes(selfReport, maxSelfReportRunes),
		Signals:           signals,
		LimitationNote:    identityLimitationNote,
		LatencyMS:         latencyMS,
		StatusCode:        statusCode,
	}
}

func extractVendorHeaders(headers http.Header) map[string]string {
	if len(headers) == 0 {
		return nil
	}
	allowExact := map[string]bool{
		"openai-model":               true,
		"openai-version":             true,
		"openai-processing-ms":       true,
		"anthropic-organization-id":  true,
		"x-goog-api-version":         true,
		"x-openrouter-provider-name": true,
	}
	out := make(map[string]string)
	for k, vals := range headers {
		lower := strings.ToLower(strings.TrimSpace(k))
		if len(vals) == 0 {
			continue
		}
		if allowExact[lower] || strings.HasPrefix(lower, "x-deepseek-") || strings.HasPrefix(lower, "x-qwen-") {
			v := strings.TrimSpace(vals[0])
			if v != "" {
				out[lower] = v
			}
		}
	}
	if len(out) == 0 {
		return nil
	}
	return out
}

func parseIdentityPayload(body []byte) (reportedModel, idPrefix, fingerprint, selfReport string, hasReasoning bool) {
	var raw map[string]any
	if err := json.Unmarshal(body, &raw); err != nil {
		if bytes.Contains(body, []byte("data:")) {
			return parseSSEIdentityPayload(body)
		}
		return "", "", "", "", false
	}

	if m, ok := raw["model"].(string); ok {
		reportedModel = strings.TrimSpace(m)
	} else if mv, ok := raw["modelVersion"].(string); ok {
		reportedModel = strings.TrimSpace(mv)
	}
	if fp, ok := raw["system_fingerprint"].(string); ok {
		fingerprint = strings.TrimSpace(fp)
	}
	if id, ok := raw["id"].(string); ok {
		idPrefix = extractIDPrefix(strings.TrimSpace(id))
	}

	// OpenAI Chat Completions
	if choices, ok := raw["choices"].([]any); ok && len(choices) > 0 {
		if first, ok := choices[0].(map[string]any); ok {
			if msg, ok := first["message"].(map[string]any); ok {
				if rc, ok := msg["reasoning_content"].(string); ok && strings.TrimSpace(rc) != "" {
					hasReasoning = true
				}
				if c, ok := msg["content"].(string); ok {
					selfReport = strings.TrimSpace(c)
				}
			}
		}
	}

	// Anthropic Messages
	if selfReport == "" {
		if content, ok := raw["content"].([]any); ok {
			var parts []string
			for _, item := range content {
				block, ok := item.(map[string]any)
				if !ok {
					continue
				}
				if block["type"] == "thinking" {
					hasReasoning = true
				}
				if txt, ok := block["text"].(string); ok && strings.TrimSpace(txt) != "" {
					parts = append(parts, strings.TrimSpace(txt))
				}
			}
			selfReport = strings.Join(parts, " ")
		}
	}

	// OpenAI Responses
	if selfReport == "" {
		if output, ok := raw["output"].([]any); ok {
			var parts []string
			for _, item := range output {
				entry, ok := item.(map[string]any)
				if !ok {
					continue
				}
				if content, ok := entry["content"].([]any); ok {
					for _, sub := range content {
						block, ok := sub.(map[string]any)
						if !ok {
							continue
						}
						if txt, ok := block["text"].(string); ok && strings.TrimSpace(txt) != "" {
							parts = append(parts, strings.TrimSpace(txt))
						}
					}
				}
			}
			selfReport = strings.Join(parts, " ")
		}
	}

	// Gemini GenerateContent
	if selfReport == "" {
		if candidates, ok := raw["candidates"].([]any); ok && len(candidates) > 0 {
			if first, ok := candidates[0].(map[string]any); ok {
				if content, ok := first["content"].(map[string]any); ok {
					if partsArr, ok := content["parts"].([]any); ok {
						var parts []string
						for _, p := range partsArr {
							pm, ok := p.(map[string]any)
							if !ok {
								continue
							}
							if txt, ok := pm["text"].(string); ok && strings.TrimSpace(txt) != "" {
								parts = append(parts, strings.TrimSpace(txt))
							}
						}
						selfReport = strings.Join(parts, " ")
					}
				}
			}
		}
	}

	return reportedModel, idPrefix, fingerprint, strings.TrimSpace(selfReport), hasReasoning
}

func extractIDPrefix(id string) string {
	for _, prefix := range []string{"chatcmpl-", "msg_", "resp_", "cmpl-"} {
		if strings.HasPrefix(id, prefix) {
			return prefix
		}
	}
	return ""
}

func normalizeModelForCompare(raw string) string {
	s := strings.ToLower(strings.TrimSpace(raw))
	if idx := strings.LastIndex(s, "/"); idx >= 0 {
		s = s[idx+1:]
	}
	if idx := strings.Index(s, ":"); idx >= 0 {
		s = s[:idx]
	}
	return s
}

func modelsLooselyMatch(claimed, reported string) bool {
	c := normalizeModelForCompare(claimed)
	r := normalizeModelForCompare(reported)
	if c == "" || r == "" {
		return true
	}
	if c == r || strings.HasPrefix(r, c+"-") || strings.HasPrefix(c, r+"-") ||
		strings.HasPrefix(r, c+"@") || strings.HasPrefix(c, r+"@") {
		return true
	}
	return false
}

func inferModelFamily(text string) string {
	s := strings.ToLower(strings.TrimSpace(text))
	if s == "" {
		return ""
	}
	type familyMatcher struct {
		family   string
		keywords []string
	}
	matchers := []familyMatcher{
		{family: "anthropic", keywords: []string{"claude", "anthropic"}},
		{family: "deepseek", keywords: []string{"deepseek"}},
		{family: "qwen", keywords: []string{"qwen", "qwq", "alibaba cloud", "tongyi"}},
		{family: "google", keywords: []string{"gemini", "gemma", "google deepmind"}},
		{family: "xai", keywords: []string{"grok", "xai", "x.ai"}},
		{family: "meta", keywords: []string{"llama", "meta ai"}},
		{family: "mistral", keywords: []string{"mistral", "ministral", "codestral", "mixtral"}},
		{family: "moonshot", keywords: []string{"kimi", "moonshot"}},
		{family: "zhipu", keywords: []string{"glm-", "chatglm", "zhipu"}},
		{family: "minimax", keywords: []string{"minimax", "abab"}},
		{family: "cohere", keywords: []string{"command-r", "cohere"}},
		{family: "openai", keywords: []string{"gpt-", "chatgpt", "openai", "o1-", "o3-", "o4-"}},
	}
	for _, matcher := range matchers {
		for _, kw := range matcher.keywords {
			if strings.Contains(s, kw) {
				return matcher.family
			}
		}
	}
	switch s {
	case "o1", "o3", "o4":
		return "openai"
	}
	return ""
}

func truncateRunes(s string, maxRunes int) string {
	runes := []rune(strings.TrimSpace(s))
	if len(runes) <= maxRunes {
		return string(runes)
	}
	return string(runes[:maxRunes]) + "..."
}
