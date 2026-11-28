package venice

import (
	"bufio"
	"bytes"
	"context"
	"crypto/rand"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/omeg4monx/router-gateway/internal/contracts"
	"github.com/omeg4monx/router-gateway/internal/platform/httpclient"
)

const (
	// DefaultOuterfaceOrigin is the base origin for Venice keyless webchat APIs.
	DefaultOuterfaceOrigin = "https://outerface.venice.ai"
	// DefaultInferenceURL is the Venice keyless webchat chat inference endpoint.
	DefaultInferenceURL = DefaultOuterfaceOrigin + "/api/inference/chat"
	// MintPath is the relative path for Venice client attestation token minting.
	MintPath = "/api/app/user/client_attestation/mint"
	// MiddlefaceVersion matches the official Venice webchat middleface header version.
	MiddlefaceVersion = "0.1.1272"
	// DefaultModel is used when no explicit Venice model ID is supplied.
	DefaultModel = "deepseek-v4-1-flash-uncensored"
	// VirtualConnectionID identifies the virtual keyless connection used when no DB connection exists.
	VirtualConnectionID = "noauth:venice"
	// DirectProxyKey is the session/token bucket key used for direct (non-proxied) Venice requests.
	DirectProxyKey = "direct"

	mintBatchCount        = 5
	tokenExpirationGrace  = 60 * time.Second
	defaultTokenTTL       = 5 * time.Minute
	sessionTTL            = 10 * time.Minute
	transientRetryDelay   = 400 * time.Millisecond
	maxAttestationBuckets = 32
	maxVeniceSessions     = 32
	maxMintResponseBytes  = 64 << 10
	maxNDJSONLineBytes    = 1 << 20
	estimatedPromptTokens = 15
	secCHUAHeader         = `"Not(A:Brand";v="99", "Google Chrome";v="133", "Chromium";v="133"`
	secCHUAMobileHeader   = "?0"
	secCHUAPlatformHeader = `"Windows"`
)

// DefaultModels returns the canonical list of Venice keyless webchat model IDs.
func DefaultModels() []string {
	return []string{
		"deepseek-v4-1-flash-uncensored",
		"olafangensan-glm-4.7-flash-heretic",
		"venice-uncensored-1-2",
		"qwen3-coder-480b-a35b-instruct-turbo",
		"nvidia-nemotron-3-nano-30b-a3b",
		"nvidia-nemotron-3-ultra-550b-a55b",
	}
}

// IsVeniceProviderID returns true for canonical Venice provider IDs and aliases.
func IsVeniceProviderID(providerID string) bool {
	switch strings.ToLower(strings.TrimSpace(providerID)) {
	case "venice", "venice-ai", "veniceai", "venice-free":
		return true
	default:
		return false
	}
}

// IsVirtualConnectionID reports whether connectionID is the virtual keyless Venice connection.
func IsVirtualConnectionID(connectionID string) bool {
	return strings.TrimSpace(connectionID) == VirtualConnectionID
}

// DefaultCredentials returns keyless upstream credentials for Venice webchat inference.
func DefaultCredentials() contracts.UpstreamCredentials {
	return contracts.UpstreamCredentials{
		ConnectionID: VirtualConnectionID,
		ProviderID:   "venice",
		Protocol:     "openai-chat",
		BaseURL:      DefaultInferenceURL,
		AuthType:     "none",
	}
}

// IsVeniceKeyless reports whether creds should execute via the keyless outerface
// attestation + NDJSON webchat protocol rather than paid OpenAI-compatible api.venice.ai.
func IsVeniceKeyless(creds contracts.UpstreamCredentials) bool {
	if !IsVeniceProviderID(creds.ProviderID) {
		return false
	}
	baseLower := strings.ToLower(strings.TrimSpace(creds.BaseURL))
	if strings.Contains(baseLower, "outerface.venice.ai") || strings.HasSuffix(baseLower, "/inference/chat") {
		return true
	}
	if strings.EqualFold(strings.TrimSpace(creds.AuthType), "none") {
		return true
	}
	return strings.TrimSpace(creds.APIKey) == "" && strings.TrimSpace(creds.AccessToken) == ""
}

// ResolveInferenceURL returns the full Venice chat inference URL for baseURL.
func ResolveInferenceURL(baseURL string) string {
	trimmed := strings.TrimRight(strings.TrimSpace(baseURL), "/")
	if trimmed == "" {
		return DefaultInferenceURL
	}
	if strings.HasSuffix(trimmed, "/inference/chat") {
		return trimmed
	}
	return trimmed + "/api/inference/chat"
}

// ResolveMintURL derives the attestation mint URL from an inference or origin baseURL.
func ResolveMintURL(baseURL string) string {
	trimmed := strings.TrimSpace(baseURL)
	if trimmed == "" {
		return DefaultOuterfaceOrigin + MintPath
	}
	parsed, err := url.Parse(trimmed)
	if err != nil || parsed.Scheme == "" || parsed.Host == "" {
		return DefaultOuterfaceOrigin + MintPath
	}
	return parsed.Scheme + "://" + parsed.Host + MintPath
}

// StripModelPrefix removes venice/ aliases from model and falls back to DefaultModel.
func StripModelPrefix(model string) string {
	clean := strings.TrimSpace(model)
	for _, prefix := range []string{"venice/", "venice-ai/", "veniceai/", "venice-free/"} {
		if strings.HasPrefix(strings.ToLower(clean), prefix) {
			clean = strings.TrimSpace(clean[len(prefix):])
			break
		}
	}
	if clean == "" {
		return DefaultModel
	}
	return clean
}

// ApplyHeaders sets the Venice middleface headers required by outerface.venice.ai.
func ApplyHeaders(req *http.Request, attestationToken string, isInference bool) {
	req.Header.Del("Authorization")
	req.Header.Set("Content-Type", "application/json")
	if isInference {
		req.Header.Set("Accept", "text/event-stream, application/json, */*")
	} else {
		req.Header.Set("Accept", "application/json")
	}
	req.Header.Set("Origin", "https://venice.ai")
	req.Header.Set("Referer", "https://venice.ai/chat/v2")
	req.Header.Set("User-Agent", "venice-middleface/"+MiddlefaceVersion)
	req.Header.Set("x-venice-middleface-version", MiddlefaceVersion)
	req.Header.Set("x-venice-request-timestamp-ms", strconv.FormatInt(time.Now().UnixMilli(), 10))
	req.Header.Set("Sec-CH-UA", secCHUAHeader)
	req.Header.Set("Sec-CH-UA-Mobile", secCHUAMobileHeader)
	req.Header.Set("Sec-CH-UA-Platform", secCHUAPlatformHeader)
	req.Header.Set("Sec-Fetch-Dest", "empty")
	req.Header.Set("Sec-Fetch-Mode", "cors")
	req.Header.Set("Sec-Fetch-Site", "same-site")
	if isInference && strings.TrimSpace(attestationToken) != "" {
		req.Header.Set("x-venice-client-attestation", attestationToken)
	}
}

type pooledVeniceSession struct {
	client    *http.Client
	createdAt time.Time
	lastUsed  time.Time
}

type veniceSessionPool struct {
	mu       sync.Mutex
	sessions map[string]*pooledVeniceSession
}

var (
	sharedVeniceClientOnce sync.Once
	sharedVeniceClient     *http.Client
	sharedVeniceJarOnce    sync.Once
	sharedVeniceJar        http.CookieJar
	defaultSessionPool     = &veniceSessionPool{
		sessions: make(map[string]*pooledVeniceSession),
	}
)

func getSharedCookieJar() http.CookieJar {
	sharedVeniceJarOnce.Do(func() {
		sharedVeniceJar, _ = cookiejar.New(nil)
	})
	return sharedVeniceJar
}

func normalizeProxyKey(proxyKey string) string {
	trimmed := strings.TrimSpace(proxyKey)
	if trimmed == "" {
		return DirectProxyKey
	}
	return trimmed
}

func (p *veniceSessionPool) getOrCreate(proxyKey string, proxyURL *url.URL, timeout time.Duration) *http.Client {
	key := normalizeProxyKey(proxyKey)
	now := time.Now()

	p.mu.Lock()
	defer p.mu.Unlock()

	if entry, ok := p.sessions[key]; ok {
		if now.Sub(entry.createdAt) < sessionTTL {
			entry.lastUsed = now
			if entry.client.Timeout == timeout {
				return entry.client
			}
			cloned := *entry.client
			cloned.Timeout = timeout
			return &cloned
		}
		delete(p.sessions, key)
	}

	if len(p.sessions) >= maxVeniceSessions {
		var oldestKey string
		var oldestTime time.Time
		for k, s := range p.sessions {
			if oldestKey == "" || s.lastUsed.Before(oldestTime) {
				oldestKey = k
				oldestTime = s.lastUsed
			}
		}
		if oldestKey != "" {
			delete(p.sessions, oldestKey)
		}
	}

	var proxyFn func(*http.Request) (*url.URL, error)
	if proxyURL != nil {
		fixedURL := *proxyURL
		proxyFn = func(*http.Request) (*url.URL, error) {
			u := fixedURL
			return &u, nil
		}
	} else {
		proxyFn = func(*http.Request) (*url.URL, error) {
			return nil, nil
		}
	}

	client := httpclient.NewChromeTLSClientWithProxy(0, proxyFn)
	jar, _ := cookiejar.New(nil)
	client.Jar = jar

	p.sessions[key] = &pooledVeniceSession{
		client:    client,
		createdAt: now,
		lastUsed:  now,
	}

	if timeout == 0 {
		return client
	}
	cloned := *client
	cloned.Timeout = timeout
	return &cloned
}

func (p *veniceSessionPool) evict(proxyKey string) {
	key := normalizeProxyKey(proxyKey)
	p.mu.Lock()
	delete(p.sessions, key)
	p.mu.Unlock()
}

// EvictSession drops the pooled Chrome uTLS client and CookieJar for proxyKey.
func EvictSession(proxyKey string) {
	defaultSessionPool.evict(proxyKey)
}

// SharedHTTPClient returns a process-wide pooled Chrome uTLS client with an
// in-memory CookieJar so Cloudflare cookies minted on attestation carry over to
// inference requests.
func SharedHTTPClient() *http.Client {
	sharedVeniceClientOnce.Do(func() {
		base := httpclient.NewChromeTLSClient(0)
		base.Jar = getSharedCookieJar()
		sharedVeniceClient = base
	})
	return sharedVeniceClient
}

// SelectHTTPClient preserves configured.Transport (including omeg4mon-router's
// provider-bound Proxy Pool transport and custom test RoundTrippers) while
// ensuring a shared CookieJar is attached for Venice attestation affinity.
func SelectHTTPClient(configured *http.Client) *http.Client {
	if configured == nil {
		return SharedHTTPClient()
	}
	if configured.Jar != nil {
		return configured
	}
	cloned := *configured
	cloned.Jar = getSharedCookieJar()
	return &cloned
}

// ResolveSession returns a Chrome uTLS *http.Client and its canonical proxyKey
// for a Venice request attempt. When configured uses an *http.Transport with a
// Proxy resolver (such as omeg4mon-router's Proxy Pool), the proxy URL is
// resolved once so both attestation mint and inference share the same outbound
// proxy, Chrome TLS fingerprint, and dedicated CookieJar. Custom non-*http.Transport
// RoundTrippers (such as unit-test stubs) are preserved verbatim.
func ResolveSession(ctx context.Context, configured *http.Client, targetURL string) (client *http.Client, proxyKey string) {
	if configured == nil {
		return defaultSessionPool.getOrCreate(DirectProxyKey, nil, 0), DirectProxyKey
	}

	tr, isStdTransport := configured.Transport.(*http.Transport)
	if !isStdTransport {
		if configured.Jar != nil {
			return configured, DirectProxyKey
		}
		cloned := *configured
		cloned.Jar = getSharedCookieJar()
		return &cloned, DirectProxyKey
	}

	var proxyURL *url.URL
	if tr.Proxy != nil {
		probeTarget := ResolveInferenceURL(targetURL)
		if req, err := http.NewRequestWithContext(ctx, http.MethodPost, probeTarget, http.NoBody); err == nil {
			if resolved, proxyErr := tr.Proxy(req); proxyErr == nil && resolved != nil {
				proxyURL = resolved
			}
		}
	}

	proxyKey = DirectProxyKey
	if proxyURL != nil {
		proxyKey = proxyURL.String()
	}
	return defaultSessionPool.getOrCreate(proxyKey, proxyURL, configured.Timeout), proxyKey
}

type tokenBucket struct {
	queue     []string
	expiresAt time.Time
	lastUsed  time.Time
}

type inFlightMint struct {
	done chan struct{}
	err  error
}

// AttestationManager caches batched Venice client attestation tokens per origin
// with preemptive expiry refresh and singleflight deduplication.
type AttestationManager struct {
	mu       sync.Mutex
	buckets  map[string]*tokenBucket
	inFlight map[string]*inFlightMint
	nowFn    func() time.Time
}

// NewAttestationManager creates a new AttestationManager. Passing nil for nowFn uses time.Now.
func NewAttestationManager(nowFn func() time.Time) *AttestationManager {
	if nowFn == nil {
		nowFn = time.Now
	}
	return &AttestationManager{
		buckets:  make(map[string]*tokenBucket),
		inFlight: make(map[string]*inFlightMint),
		nowFn:    nowFn,
	}
}

var (
	defaultMgrOnce sync.Once
	defaultMgr     *AttestationManager
)

// DefaultAttestationManager returns the process-wide shared AttestationManager.
func DefaultAttestationManager() *AttestationManager {
	defaultMgrOnce.Do(func() {
		defaultMgr = NewAttestationManager(nil)
	})
	return defaultMgr
}

func attestationBucketKey(mintURL, proxyKey string) string {
	return mintURL + "|" + normalizeProxyKey(proxyKey)
}

// Evict removes cached attestation tokens for the direct origin of baseURL.
func (m *AttestationManager) Evict(baseURL string) {
	m.EvictForProxy(baseURL, DirectProxyKey)
}

// EvictForProxy removes cached attestation tokens for (baseURL, proxyKey) and
// drops the corresponding pooled Venice session/CookieJar.
func (m *AttestationManager) EvictForProxy(baseURL, proxyKey string) {
	key := attestationBucketKey(ResolveMintURL(baseURL), proxyKey)
	m.mu.Lock()
	delete(m.buckets, key)
	m.mu.Unlock()
	EvictSession(proxyKey)
}

// GetToken returns a fresh Venice client attestation token for direct connections.
func (m *AttestationManager) GetToken(ctx context.Context, client *http.Client, baseURL string) (string, error) {
	return m.GetTokenForProxy(ctx, client, baseURL, DirectProxyKey)
}

// GetTokenForProxy returns a fresh Venice client attestation token isolated to
// (baseURL, proxyKey), batch-minting and deduplicating concurrent requests when
// the local queue is empty or near expiry.
func (m *AttestationManager) GetTokenForProxy(ctx context.Context, client *http.Client, baseURL, proxyKey string) (string, error) {
	mintURL := ResolveMintURL(baseURL)
	bucketKey := attestationBucketKey(mintURL, proxyKey)
	httpClient := SelectHTTPClient(client)

	for {
		m.mu.Lock()
		now := m.nowFn()
		if bucket, ok := m.buckets[bucketKey]; ok {
			if len(bucket.queue) > 0 && bucket.expiresAt.Sub(now) > tokenExpirationGrace {
				tok := bucket.queue[0]
				bucket.queue = bucket.queue[1:]
				bucket.lastUsed = now
				m.mu.Unlock()
				return tok, nil
			}
			delete(m.buckets, bucketKey)
		}

		if call, ok := m.inFlight[bucketKey]; ok {
			m.mu.Unlock()
			select {
			case <-call.done:
				if call.err != nil {
					return "", call.err
				}
				continue
			case <-ctx.Done():
				return "", ctx.Err()
			}
		}

		call := &inFlightMint{done: make(chan struct{})}
		m.inFlight[bucketKey] = call
		m.mu.Unlock()

		go func() {
			mintCtx, cancel := context.WithTimeout(context.WithoutCancel(ctx), 30*time.Second)
			defer cancel()

			var (
				tokens    []string
				expiresAt time.Time
				err       error
			)
			defer func() {
				if rec := recover(); rec != nil {
					err = fmt.Errorf("venice: attestation mint panicked: %v", rec)
				}
				m.mu.Lock()
				call.err = err
				delete(m.inFlight, bucketKey)
				if err == nil && len(tokens) > 0 {
					if len(m.buckets) >= maxAttestationBuckets {
						m.evictOldestLocked()
					}
					m.buckets[bucketKey] = &tokenBucket{
						queue:     tokens,
						expiresAt: expiresAt,
						lastUsed:  m.nowFn(),
					}
				}
				close(call.done)
				m.mu.Unlock()
			}()

			tokens, expiresAt, err = m.mintBatch(mintCtx, httpClient, mintURL)
		}()

		select {
		case <-call.done:
			if call.err != nil {
				return "", call.err
			}
			continue
		case <-ctx.Done():
			return "", ctx.Err()
		}
	}
}

func (m *AttestationManager) evictOldestLocked() {
	var oldestKey string
	var oldestTime time.Time
	for k, b := range m.buckets {
		if oldestKey == "" || b.lastUsed.Before(oldestTime) {
			oldestKey = k
			oldestTime = b.lastUsed
		}
	}
	if oldestKey != "" {
		delete(m.buckets, oldestKey)
	}
}

func (m *AttestationManager) mintBatch(ctx context.Context, client *http.Client, mintURL string) ([]string, time.Time, error) {
	payload := fmt.Appendf(nil,
		`{"count":%d,"env":{"nativeIntegrityOk":true,"webdriver":false,"unexpectedGlobalsCount":0}}`,
		mintBatchCount,
	)

	doOnce := func() (*http.Response, error) {
		req, err := http.NewRequestWithContext(ctx, http.MethodPost, mintURL, bytes.NewReader(payload))
		if err != nil {
			return nil, err
		}
		ApplyHeaders(req, "", false)
		return client.Do(req)
	}

	resp, err := doOnce()
	if err != nil {
		if ctx.Err() != nil {
			return nil, time.Time{}, fmt.Errorf("venice: attestation mint request failed: %w", err)
		}
		timer := time.NewTimer(transientRetryDelay)
		select {
		case <-ctx.Done():
			timer.Stop()
			return nil, time.Time{}, ctx.Err()
		case <-timer.C:
		}
		resp, err = doOnce()
		if err != nil {
			return nil, time.Time{}, fmt.Errorf("venice: attestation mint retry failed: %w", err)
		}
	}

	if isTransientMintStatus(resp.StatusCode) {
		_, _ = io.Copy(io.Discard, io.LimitReader(resp.Body, maxMintResponseBytes))
		_ = resp.Body.Close()

		timer := time.NewTimer(transientRetryDelay)
		select {
		case <-ctx.Done():
			timer.Stop()
			return nil, time.Time{}, ctx.Err()
		case <-timer.C:
		}

		resp, err = doOnce()
		if err != nil {
			return nil, time.Time{}, fmt.Errorf("venice: attestation mint retry failed: %w", err)
		}
	}
	defer resp.Body.Close()

	raw, readErr := io.ReadAll(io.LimitReader(resp.Body, maxMintResponseBytes))
	if readErr != nil {
		return nil, time.Time{}, fmt.Errorf("venice: read attestation response: %w", readErr)
	}
	if resp.StatusCode != http.StatusOK {
		return nil, time.Time{}, fmt.Errorf("venice: attestation mint failed with HTTP %d: %s", resp.StatusCode, strings.TrimSpace(string(raw)))
	}

	var parsed struct {
		ExpiresAt int64             `json:"expiresAt"`
		Token     string            `json:"token"`
		Tokens    []json.RawMessage `json:"tokens"`
	}
	if err := json.Unmarshal(raw, &parsed); err != nil {
		return nil, time.Time{}, fmt.Errorf("venice: decode attestation response: %w", err)
	}

	tokens := make([]string, 0, len(parsed.Tokens)+1)
	for _, entry := range parsed.Tokens {
		var s string
		if json.Unmarshal(entry, &s) == nil && strings.TrimSpace(s) != "" {
			tokens = append(tokens, strings.TrimSpace(s))
			continue
		}
		var obj struct {
			Token string `json:"token"`
		}
		if json.Unmarshal(entry, &obj) == nil && strings.TrimSpace(obj.Token) != "" {
			tokens = append(tokens, strings.TrimSpace(obj.Token))
		}
	}
	if len(tokens) == 0 && strings.TrimSpace(parsed.Token) != "" {
		tokens = append(tokens, strings.TrimSpace(parsed.Token))
	}
	if len(tokens) == 0 {
		return nil, time.Time{}, errors.New("venice: no tokens returned from attestation mint")
	}

	expiresAt := m.nowFn().Add(defaultTokenTTL)
	if parsed.ExpiresAt > 0 {
		expiresAt = time.Unix(parsed.ExpiresAt, 0)
	}
	return tokens, expiresAt, nil
}

func isTransientMintStatus(code int) bool {
	switch code {
	case http.StatusTooManyRequests,
		http.StatusInternalServerError,
		http.StatusBadGateway,
		http.StatusServiceUnavailable,
		http.StatusGatewayTimeout:
		return true
	default:
		return false
	}
}

type openAIChatInput struct {
	Model       string            `json:"model"`
	Messages    []openAIMessage   `json:"messages"`
	Tools       []json.RawMessage `json:"tools"`
	Temperature *float64          `json:"temperature"`
	TopP        *float64          `json:"top_p"`
}

type openAIMessage struct {
	Role      string           `json:"role"`
	Name      string           `json:"name"`
	Content   json.RawMessage  `json:"content"`
	ToolCalls []openAIToolCall `json:"tool_calls"`
}

type openAIToolCall struct {
	Function struct {
		Name      string          `json:"name"`
		Arguments json.RawMessage `json:"arguments"`
	} `json:"function"`
}

type venicePromptMessage struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type veniceChatPayload struct {
	RequestID    string                `json:"requestId"`
	ModelID      string                `json:"modelId"`
	Prompt       []venicePromptMessage `json:"prompt"`
	SystemPrompt string                `json:"systemPrompt,omitempty"`
	Temperature  float64               `json:"temperature"`
	TopP         float64               `json:"topP"`
	SimpleMode   bool                  `json:"simpleMode"`
}

// BuildChatPayload translates an OpenAI Chat Completions JSON payload into the
// Venice outerface.venice.ai/api/inference/chat schema.
func BuildChatPayload(model string, openAIBody []byte) (encoded []byte, resolvedModel string, err error) {
	var in openAIChatInput
	if len(bytes.TrimSpace(openAIBody)) > 0 {
		if err := json.Unmarshal(openAIBody, &in); err != nil {
			return nil, "", fmt.Errorf("venice: decode openai request: %w", err)
		}
	}

	targetModel := model
	if strings.TrimSpace(targetModel) == "" {
		targetModel = in.Model
	}
	resolvedModel = StripModelPrefix(targetModel)

	var systemParts []string
	if len(in.Tools) > 0 {
		if toolBlock := formatToolsInstruction(in.Tools); toolBlock != "" {
			systemParts = append(systemParts, toolBlock)
		}
	}

	prompt := make([]venicePromptMessage, 0, len(in.Messages)+1)
	hasUser := false

	for _, msg := range in.Messages {
		contentStr := extractContentText(msg.Content)
		switch strings.ToLower(strings.TrimSpace(msg.Role)) {
		case "system", "developer":
			if contentStr != "" {
				systemParts = append(systemParts, contentStr)
			}
		case "tool":
			toolLabel := ""
			if strings.TrimSpace(msg.Name) != "" {
				toolLabel = " (" + strings.TrimSpace(msg.Name) + ")"
			}
			prompt = append(prompt, venicePromptMessage{
				Role:    "user",
				Content: fmt.Sprintf("[Tool Result%s]:\n%s", toolLabel, contentStr),
			})
			hasUser = true
		case "assistant":
			var sb strings.Builder
			if contentStr != "" {
				sb.WriteString(contentStr)
			}
			for _, tc := range msg.ToolCalls {
				if sb.Len() > 0 {
					sb.WriteByte('\n')
				}
				argsStr := extractToolArguments(tc.Function.Arguments)
				fmt.Fprintf(&sb, `<tool_call>{"name":%q,"arguments":%s}</tool_call>`, tc.Function.Name, argsStr)
			}
			prompt = append(prompt, venicePromptMessage{
				Role:    "assistant",
				Content: strings.TrimSpace(sb.String()),
			})
		default:
			prompt = append(prompt, venicePromptMessage{
				Role:    "user",
				Content: contentStr,
			})
			hasUser = true
		}
	}

	if !hasUser {
		prompt = append(prompt, venicePromptMessage{
			Role:    "user",
			Content: "Hello",
		})
	}

	temperature := 0.7
	if in.Temperature != nil {
		temperature = *in.Temperature
	}
	topP := 0.9
	if in.TopP != nil {
		topP = *in.TopP
	}

	payload := veniceChatPayload{
		RequestID:    randomUUIDv4(),
		ModelID:      resolvedModel,
		Prompt:       prompt,
		SystemPrompt: strings.Join(systemParts, "\n\n"),
		Temperature:  temperature,
		TopP:         topP,
		SimpleMode:   false,
	}

	encoded, err = json.Marshal(payload)
	if err != nil {
		return nil, "", fmt.Errorf("venice: marshal chat payload: %w", err)
	}
	return encoded, resolvedModel, nil
}

func formatToolsInstruction(tools []json.RawMessage) string {
	lines := make([]string, 0, len(tools))
	for _, rawTool := range tools {
		var wrapper struct {
			Function *struct {
				Name        string          `json:"name"`
				Description string          `json:"description"`
				Parameters  json.RawMessage `json:"parameters"`
			} `json:"function"`
			Name        string          `json:"name"`
			Description string          `json:"description"`
			Parameters  json.RawMessage `json:"parameters"`
		}
		if json.Unmarshal(rawTool, &wrapper) != nil {
			continue
		}
		name, desc, params := wrapper.Name, wrapper.Description, wrapper.Parameters
		if wrapper.Function != nil {
			name = wrapper.Function.Name
			desc = wrapper.Function.Description
			params = wrapper.Function.Parameters
		}
		if strings.TrimSpace(name) == "" {
			continue
		}
		var b strings.Builder
		b.WriteString("- ")
		b.WriteString(name)
		if strings.TrimSpace(desc) != "" {
			b.WriteString(": ")
			b.WriteString(strings.TrimSpace(desc))
		}
		if len(bytes.TrimSpace(params)) > 0 {
			b.WriteString("\n  parameters: ")
			b.Write(bytes.TrimSpace(params))
		}
		lines = append(lines, b.String())
	}
	if len(lines) == 0 {
		return ""
	}
	return "[Tools Available]\nYou have access to the following tools:\n" +
		strings.Join(lines, "\n") +
		"\n\nTo invoke a tool, output a JSON tool_call block:\n" +
		`<tool_call>{"name": "tool_name", "arguments": { ... }}</tool_call>` +
		"\nOnly call a tool when necessary to fulfill the request."
}

func extractContentText(raw json.RawMessage) string {
	trimmed := bytes.TrimSpace(raw)
	if len(trimmed) == 0 || bytes.Equal(trimmed, []byte("null")) {
		return ""
	}
	var s string
	if json.Unmarshal(trimmed, &s) == nil {
		return s
	}
	var parts []struct {
		Type string `json:"type"`
		Text string `json:"text"`
	}
	if json.Unmarshal(trimmed, &parts) == nil {
		var sb strings.Builder
		for _, p := range parts {
			if p.Text != "" {
				sb.WriteString(p.Text)
			}
		}
		if sb.Len() > 0 {
			return sb.String()
		}
	}
	return string(trimmed)
}

func extractToolArguments(raw json.RawMessage) string {
	trimmed := bytes.TrimSpace(raw)
	if len(trimmed) == 0 {
		return "{}"
	}
	var s string
	if json.Unmarshal(trimmed, &s) == nil {
		s = strings.TrimSpace(s)
		if s == "" {
			return "{}"
		}
		if json.Valid([]byte(s)) {
			return s
		}
		quoted, _ := json.Marshal(s)
		return string(quoted)
	}
	if json.Valid(trimmed) {
		return string(trimmed)
	}
	return "{}"
}

func randomUUIDv4() string {
	var b [16]byte
	_, _ = rand.Read(b[:])
	b[6] = (b[6] & 0x0f) | 0x40
	b[8] = (b[8] & 0x3f) | 0x80
	return fmt.Sprintf("%08x-%04x-%04x-%04x-%012x", b[0:4], b[4:6], b[6:8], b[8:10], b[10:16])
}

type veniceNDJSONLine struct {
	Kind         string `json:"kind"`
	Content      string `json:"content"`
	CompletionID string `json:"completion_id"`
}

// ParseNDJSONToOpenAI converts a Venice NDJSON response body into an OpenAI
// chat.completion JSON payload. If raw is already a valid OpenAI chat.completion
// JSON object (with choices), it is returned as-is.
func ParseNDJSONToOpenAI(raw []byte, model string) (openAIJSON []byte, fullText string, err error) {
	return ParseNDJSONToOpenAIWithTools(raw, model, nil)
}

// ParseNDJSONToOpenAIWithTools converts a Venice NDJSON response body into an
// OpenAI chat.completion JSON payload and, when requestTools is non-empty,
// bridges prose <tool_call> blocks into structured OpenAI message.tool_calls.
func ParseNDJSONToOpenAIWithTools(raw []byte, model string, requestTools []json.RawMessage) (openAIJSON []byte, fullText string, err error) {
	trimmed := bytes.TrimSpace(raw)
	if len(trimmed) == 0 {
		return nil, "", errors.New("venice: empty upstream response")
	}

	// Pass through already-formatted OpenAI JSON responses (e.g. mocks or compatible proxies).
	if trimmed[0] == '{' && json.Valid(trimmed) {
		var probe struct {
			Choices []struct {
				Message struct {
					Content string `json:"content"`
				} `json:"message"`
			} `json:"choices"`
		}
		if json.Unmarshal(trimmed, &probe) == nil && len(probe.Choices) > 0 {
			return bytes.Clone(trimmed), probe.Choices[0].Message.Content, nil
		}
	}

	var (
		completionID string
		textBuilder  strings.Builder
		parsedAny    bool
	)
	scanner := bufio.NewScanner(bytes.NewReader(trimmed))
	scanner.Buffer(make([]byte, 0, 16*1024), maxNDJSONLineBytes)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" {
			continue
		}
		var item veniceNDJSONLine
		if json.Unmarshal([]byte(line), &item) != nil {
			continue
		}
		if item.Kind == "" {
			continue
		}
		parsedAny = true
		if item.CompletionID != "" && completionID == "" {
			completionID = item.CompletionID
		}
		if item.Kind == "content" && item.Content != "" {
			textBuilder.WriteString(item.Content)
		}
	}
	if err := scanner.Err(); err != nil {
		return nil, "", fmt.Errorf("venice: scan ndjson response: %w", err)
	}
	if !parsedAny {
		return nil, "", errors.New("venice: upstream response contained no valid NDJSON frames")
	}

	fullText = textBuilder.String()
	if completionID == "" {
		completionID = fmt.Sprintf("chatcmpl-vnc-%d", time.Now().UnixMilli())
	}
	completionTokens := (len(fullText) + 3) / 4

	msgObj := map[string]any{
		"role":    "assistant",
		"content": fullText,
	}
	finishReason := "stop"
	if len(requestTools) > 0 && fullText != "" {
		if cleaned, calls, applied := ApplyToolBridge(fullText, requestTools); applied {
			msgObj["content"] = cleaned
			msgObj["tool_calls"] = calls
			finishReason = "tool_calls"
		}
	}

	respObj := map[string]any{
		"id":      completionID,
		"object":  "chat.completion",
		"created": time.Now().Unix(),
		"model":   model,
		"choices": []map[string]any{
			{
				"index":         0,
				"message":       msgObj,
				"finish_reason": finishReason,
			},
		},
		"usage": map[string]any{
			"prompt_tokens":     estimatedPromptTokens,
			"completion_tokens": completionTokens,
			"total_tokens":      estimatedPromptTokens + completionTokens,
		},
	}

	out, err := json.Marshal(respObj)
	if err != nil {
		return nil, "", fmt.Errorf("venice: marshal openai response: %w", err)
	}
	return out, fullText, nil
}

type ssePipeReader struct {
	pr  *io.PipeReader
	src io.ReadCloser
}

func (r *ssePipeReader) Read(p []byte) (int, error) {
	return r.pr.Read(p)
}

func (r *ssePipeReader) Close() error {
	_ = r.pr.Close()
	return r.src.Close()
}

// NewNDJSONToOpenAISSEReader wraps a Venice NDJSON stream and emits OpenAI
// chat.completion.chunk SSE frames in real time so downstream gateway SSE
// translators can stream to OpenAI, Anthropic, Responses, or Gemini clients.
func NewNDJSONToOpenAISSEReader(ctx context.Context, src io.ReadCloser, model string) io.ReadCloser {
	return NewNDJSONToOpenAISSEReaderWithTools(ctx, src, model, nil)
}

// NewNDJSONToOpenAISSEReaderWithTools wraps a Venice NDJSON stream into OpenAI
// SSE chunks and, when requestTools is non-empty, bridges prose <tool_call>
// blocks into structured delta.tool_calls chunks with finish_reason "tool_calls".
func NewNDJSONToOpenAISSEReaderWithTools(ctx context.Context, src io.ReadCloser, model string, requestTools []json.RawMessage) io.ReadCloser {
	pr, pw := io.Pipe()
	go func() {
		defer func() {
			_ = src.Close()
		}()

		completionID := fmt.Sprintf("chatcmpl-vnc-%d", time.Now().UnixMilli())
		created := time.Now().Unix()

		writeChunk := func(delta map[string]any, finishReason any) error {
			chunk := map[string]any{
				"id":      completionID,
				"object":  "chat.completion.chunk",
				"created": created,
				"model":   model,
				"choices": []map[string]any{
					{
						"index":         0,
						"delta":         delta,
						"finish_reason": finishReason,
					},
				},
			}
			raw, err := json.Marshal(chunk)
			if err != nil {
				return err
			}
			_, err = fmt.Fprintf(pw, "data: %s\n\n", raw)
			return err
		}

		scanner := bufio.NewScanner(src)
		scanner.Buffer(make([]byte, 0, 16*1024), maxNDJSONLineBytes)

		var toolBuf strings.Builder
		bufferingForTools := false
		hasTools := len(requestTools) > 0

		for scanner.Scan() {
			if ctx.Err() != nil {
				_ = pw.CloseWithError(ctx.Err())
				return
			}
			line := strings.TrimSpace(scanner.Text())
			if line == "" {
				continue
			}
			var item veniceNDJSONLine
			if json.Unmarshal([]byte(line), &item) != nil {
				continue
			}
			if item.CompletionID != "" {
				completionID = item.CompletionID
			}
			if item.Kind != "content" || item.Content == "" {
				continue
			}

			if !hasTools {
				if err := writeChunk(map[string]any{"content": item.Content}, nil); err != nil {
					_ = pw.CloseWithError(err)
					return
				}
				continue
			}

			if !bufferingForTools {
				idx := strings.IndexAny(item.Content, "<`")
				if idx < 0 {
					if err := writeChunk(map[string]any{"content": item.Content}, nil); err != nil {
						_ = pw.CloseWithError(err)
						return
					}
					continue
				}
				if idx > 0 {
					if err := writeChunk(map[string]any{"content": item.Content[:idx]}, nil); err != nil {
						_ = pw.CloseWithError(err)
						return
					}
				}
				tail := item.Content[idx:]
				if couldBeToolMarkerPrefix(tail) {
					bufferingForTools = true
					toolBuf.WriteString(tail)
				} else if err := writeChunk(map[string]any{"content": tail}, nil); err != nil {
					_ = pw.CloseWithError(err)
					return
				}
			} else {
				toolBuf.WriteString(item.Content)
				if !couldBeToolMarkerPrefix(toolBuf.String()) {
					if err := writeChunk(map[string]any{"content": toolBuf.String()}, nil); err != nil {
						_ = pw.CloseWithError(err)
						return
					}
					toolBuf.Reset()
					bufferingForTools = false
				}
			}
		}

		if err := scanner.Err(); err != nil && ctx.Err() == nil {
			_ = pw.CloseWithError(err)
			return
		}
		if ctx.Err() != nil {
			_ = pw.CloseWithError(ctx.Err())
			return
		}

		finishReason := "stop"
		if bufferingForTools && toolBuf.Len() > 0 {
			buffered := toolBuf.String()
			if mayContainToolMarker(buffered) {
				if cleaned, calls, applied := ApplyToolBridge(buffered, requestTools); applied {
					if cleaned != "" {
						if err := writeChunk(map[string]any{"content": cleaned}, nil); err != nil {
							_ = pw.CloseWithError(err)
							return
						}
					}
					for idx, call := range calls {
						delta := map[string]any{
							"tool_calls": []map[string]any{
								{
									"index": idx,
									"id":    call.ID,
									"type":  call.Type,
									"function": map[string]any{
										"name":      call.Function.Name,
										"arguments": call.Function.Arguments,
									},
								},
							},
						}
						if err := writeChunk(delta, nil); err != nil {
							_ = pw.CloseWithError(err)
							return
						}
					}
					finishReason = "tool_calls"
				} else if err := writeChunk(map[string]any{"content": buffered}, nil); err != nil {
					_ = pw.CloseWithError(err)
					return
				}
			} else if err := writeChunk(map[string]any{"content": buffered}, nil); err != nil {
				_ = pw.CloseWithError(err)
				return
			}
		}

		if err := writeChunk(map[string]any{}, finishReason); err == nil {
			_, _ = io.WriteString(pw, "data: [DONE]\n\n")
		}
		_ = pw.Close()
	}()

	return &ssePipeReader{pr: pr, src: src}
}
