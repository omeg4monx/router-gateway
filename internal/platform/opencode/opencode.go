package opencode

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"io"
	"net/http"
	"net/url"
	"os"
	"regexp"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/omeg4monx/router-gateway/internal/contracts"
	"github.com/omeg4monx/router-gateway/internal/platform/httpclient"
)

const (
	// DefaultFreeBaseURL is the OpenCode v2 free-tier OpenAI-compatible base URL.
	DefaultFreeBaseURL = "https://opencode.ai/inference/openai/v1"
	// DefaultZenBaseURL is the OpenCode Zen pay-as-you-go base URL.
	DefaultZenBaseURL = "https://opencode.ai/zen/v1"
	// DefaultGoBaseURL is the OpenCode Go subscription base URL.
	DefaultGoBaseURL = "https://opencode.ai/zen/go/v1"
	// DefaultGoGeminiBaseURL is the OpenCode Go Gemini-native models base URL.
	DefaultGoGeminiBaseURL = "https://opencode.ai/inference/go/v1/models"
	// DefaultSystemOneURL is the OpenCode System One (Jev) decision inference endpoint.
	DefaultSystemOneURL = "https://opencode.ai/inference/v1/systemone"
	// VirtualConnectionID identifies the virtual keyless connection used when no DB connection exists.
	VirtualConnectionID = "noauth:opencode"

	defaultVersionFloor = "2.4.1"
	defaultUpdaterURL   = "https://opencode.ai/update/api/latest/cli/npm"
	defaultChannel      = "latest"
	defaultClientHeader = "cli"
	anthropicAPIVersion = "2023-06-01"
	versionRefreshTTL   = 6 * time.Hour
	versionFetchTimeout = 5 * time.Second
	maxVersionRespBytes = 32 << 10
	maxPooledTLSClients = 32
	pooledTLSClientTTL  = 10 * time.Minute
	directProxyPoolKey  = "direct"
)

var (
	releaseVersionRe = regexp.MustCompile(`^\d{1,4}\.\d{1,4}\.\d{1,6}$`)
	validV2UARe      = regexp.MustCompile(`(?i)^opencode/[a-z0-9-]+/(\d+)\.\d+\.\d+/[a-z0-9-]+$`)
)

// DefaultFreeModels returns the curated live list of OpenCode Free model IDs.
func DefaultFreeModels() []string {
	return []string{
		"big-pickle",
		"mimo-v2.6-flash-free",
		"deepseek-v4.1-flash-free",
		"iceman-free",
		"beam-free",
		"flint-40b-free",
		"nemotron-3.5-super-vl-free",
		"nemotron-3.5-lightning-free",
		"hy4-flash-free",
		"tidepool-alpha-free",
	}
}

// IsFreeProviderID reports whether providerID targets OpenCode Free.
func IsFreeProviderID(providerID string) bool {
	switch strings.ToLower(strings.TrimSpace(providerID)) {
	case "opencode", "oc":
		return true
	default:
		return false
	}
}

// IsZenProviderID reports whether providerID targets OpenCode Zen (PAYG).
func IsZenProviderID(providerID string) bool {
	switch strings.ToLower(strings.TrimSpace(providerID)) {
	case "opencode-zen", "ocz":
		return true
	default:
		return false
	}
}

// IsGoProviderID reports whether providerID targets OpenCode Go (subscription).
func IsGoProviderID(providerID string) bool {
	switch strings.ToLower(strings.TrimSpace(providerID)) {
	case "opencode-go", "ocg":
		return true
	default:
		return false
	}
}

// IsAnyProviderID reports whether providerID belongs to the OpenCode provider family.
func IsAnyProviderID(providerID string) bool {
	return IsFreeProviderID(providerID) || IsZenProviderID(providerID) || IsGoProviderID(providerID)
}

// IsVirtualConnectionID reports whether connectionID is the virtual keyless OpenCode connection.
func IsVirtualConnectionID(connectionID string) bool {
	return strings.TrimSpace(connectionID) == VirtualConnectionID
}

// DefaultCredentials returns keyless upstream credentials for OpenCode Free.
func DefaultCredentials() contracts.UpstreamCredentials {
	return contracts.UpstreamCredentials{
		ConnectionID: VirtualConnectionID,
		ProviderID:   "opencode",
		Protocol:     "openai-chat",
		BaseURL:      DefaultFreeBaseURL,
		AuthType:     "none",
	}
}

// IsFreeCreds reports whether creds should execute via the keyless OpenCode v2 free-tier gate.
func IsFreeCreds(creds contracts.UpstreamCredentials) bool {
	if !IsFreeProviderID(creds.ProviderID) {
		return false
	}
	return strings.EqualFold(strings.TrimSpace(creds.AuthType), "none") ||
		(strings.TrimSpace(creds.APIKey) == "" && strings.TrimSpace(creds.AccessToken) == "")
}

// BaseModelID strips trailing thinking-level suffixes like "model(high)".
func BaseModelID(model string) string {
	trimmed := strings.TrimSpace(model)
	if idx := strings.LastIndex(trimmed, "("); idx > 0 && strings.HasSuffix(trimmed, ")") {
		return strings.TrimSpace(trimmed[:idx])
	}
	return trimmed
}

// IsResponsesModel reports whether model is served on an OpenAI Responses endpoint.
func IsResponsesModel(providerID, model, protocol string) bool {
	if strings.EqualFold(strings.TrimSpace(protocol), "openai-responses") {
		return true
	}
	base := strings.ToLower(BaseModelID(model))
	if strings.HasPrefix(base, "muse-spark") || strings.HasPrefix(base, "muse-wonder") ||
		base == "iceman-contributor-free" {
		return true
	}
	if IsZenProviderID(providerID) || IsGoProviderID(providerID) {
		if strings.HasPrefix(base, "gpt-5") || strings.HasPrefix(base, "gpt-6") ||
			strings.HasPrefix(base, "grok-4.") || strings.HasPrefix(base, "grok-build") {
			return true
		}
	}
	return false
}

// IsMessagesModel reports whether model is served on an Anthropic Messages endpoint.
func IsMessagesModel(providerID, model, protocol string) bool {
	if strings.EqualFold(strings.TrimSpace(protocol), "anthropic") {
		return true
	}
	base := strings.ToLower(BaseModelID(model))
	if base == "union-alpha" {
		return true
	}
	if IsZenProviderID(providerID) || IsGoProviderID(providerID) {
		if strings.HasPrefix(base, "claude-") {
			return true
		}
	}
	return false
}

// IsGeminiModel reports whether an OpenCode Go model uses the native Gemini endpoint.
func IsGeminiModel(providerID, model, protocol string) bool {
	if !IsGoProviderID(providerID) {
		return false
	}
	if strings.EqualFold(strings.TrimSpace(protocol), "gemini") {
		return true
	}
	base := strings.ToLower(BaseModelID(model))
	return strings.HasPrefix(base, "gemini-")
}

// ResolveTargetProtocol selects the exact wire protocol required by an OpenCode
// model before request translation runs, avoiding format/endpoint mismatches and
// enabling zero-translation passthrough when the model supports inProto natively.
func ResolveTargetProtocol(providerID, model, inProto, defaultProto string) string {
	base := strings.ToLower(BaseModelID(model))
	in := strings.ToLower(strings.TrimSpace(inProto))

	if IsGoProviderID(providerID) && strings.HasPrefix(base, "gemini-") {
		return "gemini"
	}
	if IsResponsesModel(providerID, model, "") {
		return "openai-responses"
	}
	if IsMessagesModel(providerID, model, "") {
		return "anthropic"
	}

	// OpenCode Go & Zen multi-format families that support direct Anthropic/Responses passthrough.
	if IsGoProviderID(providerID) || IsZenProviderID(providerID) {
		switch {
		case strings.HasPrefix(base, "deepseek-v4"):
			if in == "anthropic" || in == "openai-responses" || in == "openai-chat" {
				return in
			}
		case strings.HasPrefix(base, "minimax-") || strings.HasPrefix(base, "qwen3.") || base == "space-bunny-free":
			if in == "anthropic" || in == "openai-chat" {
				return in
			}
		}
	}

	if strings.TrimSpace(defaultProto) != "" {
		return defaultProto
	}
	return "openai-chat"
}

// ResolveUpstreamURL selects the exact OpenCode endpoint for providerID, model, and protocol.
func ResolveUpstreamURL(providerID, baseURL, model, protocol string, stream bool) string {
	trimmedBase := strings.TrimRight(strings.TrimSpace(baseURL), "/")

	if IsGoProviderID(providerID) && IsGeminiModel(providerID, model, protocol) {
		cleanModel := BaseModelID(model)
		action := "generateContent"
		if stream {
			action = "streamGenerateContent?alt=sse"
		}
		geminiBase := DefaultGoGeminiBaseURL
		if trimmedBase != "" && strings.Contains(trimmedBase, "/models") {
			geminiBase = trimmedBase
		}
		return geminiBase + "/" + cleanModel + ":" + action
	}

	if trimmedBase == "" {
		switch {
		case IsZenProviderID(providerID):
			trimmedBase = DefaultZenBaseURL
		case IsGoProviderID(providerID):
			trimmedBase = DefaultGoBaseURL
		default:
			trimmedBase = DefaultFreeBaseURL
		}
	}

	// Normalize if a caller passed a full endpoint URL instead of base URL.
	for _, suffix := range []string{"/chat/completions", "/responses", "/messages"} {
		if strings.HasSuffix(trimmedBase, suffix) {
			trimmedBase = strings.TrimSuffix(trimmedBase, suffix)
			break
		}
	}

	if IsResponsesModel(providerID, model, protocol) {
		return trimmedBase + "/responses"
	}
	if IsMessagesModel(providerID, model, protocol) {
		return trimmedBase + "/messages"
	}
	return trimmedBase + "/chat/completions"
}

// HasValidV2UserAgent reports whether ua matches opencode/<channel>/<v>=2.*/* .
func HasValidV2UserAgent(ua string) bool {
	m := validV2UARe.FindStringSubmatch(strings.TrimSpace(ua))
	if len(m) < 2 {
		return false
	}
	major, err := strconv.Atoi(m[1])
	return err == nil && major >= 2
}

// VersionResolver keeps the OpenCode CLI version fresh via non-blocking background
// polling of OpenCode's own npm updater endpoint.
type VersionResolver struct {
	mu         sync.Mutex
	updaterURL string
	client     *http.Client
	value      string
	fetchedAt  time.Time
	asking     bool
}

// NewVersionResolver creates a VersionResolver backed by updaterURL and client.
func NewVersionResolver(updaterURL string, client *http.Client) *VersionResolver {
	initial := strings.TrimSpace(os.Getenv("OPENCODE_VERSION"))
	if !releaseVersionRe.MatchString(initial) {
		initial = defaultVersionFloor
	}
	if strings.TrimSpace(updaterURL) == "" {
		updaterURL = defaultUpdaterURL
	}
	if client == nil {
		client = &http.Client{Timeout: versionFetchTimeout}
	}
	return &VersionResolver{
		updaterURL: updaterURL,
		client:     client,
		value:      initial,
	}
}

// Version returns the latest known OpenCode CLI version without blocking.
func (r *VersionResolver) Version() string {
	r.maybeRefresh()
	r.mu.Lock()
	defer r.mu.Unlock()
	return r.value
}

// UserAgent returns the canonical OpenCode v2 User-Agent header value.
func (r *VersionResolver) UserAgent() string {
	return "opencode/" + defaultChannel + "/" + r.Version() + "/" + defaultClientHeader
}

func (r *VersionResolver) maybeRefresh() {
	r.mu.Lock()
	if r.asking || (!r.fetchedAt.IsZero() && time.Since(r.fetchedAt) < versionRefreshTTL) {
		r.mu.Unlock()
		return
	}
	r.asking = true
	updaterURL := r.updaterURL
	client := r.client
	r.mu.Unlock()

	go func() {
		ctx, cancel := context.WithTimeout(context.Background(), versionFetchTimeout)
		defer cancel()

		var nextVersion string
		req, err := http.NewRequestWithContext(ctx, http.MethodGet, updaterURL, http.NoBody)
		if err == nil {
			resp, doErr := client.Do(req)
			if doErr == nil {
				body, readErr := io.ReadAll(io.LimitReader(resp.Body, maxVersionRespBytes))
				_ = resp.Body.Close()
				if readErr == nil && resp.StatusCode >= 200 && resp.StatusCode < 300 {
					var parsed struct {
						Version string `json:"version"`
					}
					if json.Unmarshal(body, &parsed) == nil {
						v := strings.TrimSpace(parsed.Version)
						if releaseVersionRe.MatchString(v) {
							nextVersion = v
						}
					}
				}
			}
		}

		r.mu.Lock()
		if nextVersion != "" {
			r.value = nextVersion
		}
		r.fetchedAt = time.Now()
		r.asking = false
		r.mu.Unlock()
	}()
}

var defaultVersionResolver = NewVersionResolver(defaultUpdaterURL, nil)

// DefaultUserAgent returns the process-wide dynamic OpenCode v2 User-Agent.
func DefaultUserAgent() string {
	return defaultVersionResolver.UserAgent()
}

// RandomHex returns a cryptographically random lowercase hex string of 2*numBytes chars.
func RandomHex(numBytes int) string {
	buf := make([]byte, numBytes)
	_, _ = rand.Read(buf)
	return hex.EncodeToString(buf)
}

// ApplyFreeHeaders applies all required OpenCode v2 free-tier headers onto req.
func ApplyFreeHeaders(req *http.Request, sessionID string, incoming http.Header, stream bool, targetURL string) {
	if strings.TrimSpace(sessionID) == "" {
		sessionID = GenerateSessionID(time.Now())
	}
	traceID := RandomHex(16)
	spanID := RandomHex(8)
	parentSpanID := RandomHex(8)

	ua := DefaultUserAgent()
	if incoming != nil {
		if downstreamUA := strings.TrimSpace(incoming.Get("User-Agent")); HasValidV2UserAgent(downstreamUA) {
			ua = downstreamUA
		}
	}

	clientHeader := defaultClientHeader
	projectID := ""
	if incoming != nil {
		if c := strings.TrimSpace(incoming.Get("x-opencode-client")); c != "" {
			clientHeader = c
		}
		if p := strings.TrimSpace(incoming.Get("x-opencode-project")); p != "" {
			projectID = p
		}
	}
	if projectID == "" {
		projectID = RandomHex(20)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer public")
	req.Header.Set("User-Agent", ua)
	if stream {
		req.Header.Set("Accept", "text/event-stream")
	} else {
		req.Header.Set("Accept", "*/*")
	}
	req.Header.Set("b3", traceID+"-"+spanID+"-1-"+parentSpanID)
	req.Header.Set("traceparent", "00-"+traceID+"-"+spanID+"-01")
	req.Header.Set("x-opencode-client", clientHeader)
	req.Header.Set("x-opencode-session", sessionID)
	req.Header.Set("x-opencode-session-id", sessionID)
	req.Header.Set("x-session-affinity", sessionID)
	req.Header.Set("x-session-id", sessionID)
	req.Header.Set("x-opencode-project", projectID)
	// Upstream rejects v1 x-opencode-request header with 403.
	req.Header.Del("x-opencode-request")
	if strings.HasSuffix(targetURL, "/messages") {
		req.Header.Set("anthropic-version", anthropicAPIVersion)
	}
}

// ApplyPaidHeaders applies OpenCode Zen / Go session & protocol headers while
// preserving the caller's API key auth.
func ApplyPaidHeaders(req *http.Request, creds contracts.UpstreamCredentials, sessionID, targetURL string) {
	if strings.TrimSpace(sessionID) != "" {
		req.Header.Set("x-opencode-session", sessionID)
		req.Header.Set("x-opencode-session-id", sessionID)
	}
	req.Header.Del("x-opencode-request")

	apiKey := strings.TrimSpace(creds.APIKey)
	if apiKey == "" {
		apiKey = strings.TrimSpace(creds.AccessToken)
	}
	switch {
	case strings.HasSuffix(targetURL, "/messages"):
		if apiKey != "" {
			req.Header.Del("Authorization")
			req.Header.Set("x-api-key", apiKey)
		}
		req.Header.Set("anthropic-version", anthropicAPIVersion)
	case strings.Contains(targetURL, ":generateContent") || strings.Contains(targetURL, ":streamGenerateContent"):
		if apiKey != "" {
			req.Header.Del("Authorization")
			req.Header.Set("x-goog-api-key", apiKey)
		}
	}
}

type pooledTLSClient struct {
	client    *http.Client
	createdAt time.Time
	lastUsed  time.Time
}

type tlsClientPool struct {
	mu      sync.Mutex
	clients map[string]*pooledTLSClient
}

var defaultTLSPool = &tlsClientPool{
	clients: make(map[string]*pooledTLSClient),
}

func (p *tlsClientPool) getOrCreate(proxyKey string, proxyURL *url.URL, timeout time.Duration) *http.Client {
	key := strings.TrimSpace(proxyKey)
	if key == "" {
		key = directProxyPoolKey
	}
	now := time.Now()

	p.mu.Lock()
	defer p.mu.Unlock()

	if entry, ok := p.clients[key]; ok {
		if now.Sub(entry.createdAt) < pooledTLSClientTTL {
			entry.lastUsed = now
			if entry.client.Timeout == timeout {
				return entry.client
			}
			cloned := *entry.client
			cloned.Timeout = timeout
			return &cloned
		}
		delete(p.clients, key)
	}

	if len(p.clients) >= maxPooledTLSClients {
		var oldestKey string
		var oldestTime time.Time
		for k, entry := range p.clients {
			if oldestKey == "" || entry.lastUsed.Before(oldestTime) {
				oldestKey = k
				oldestTime = entry.lastUsed
			}
		}
		if oldestKey != "" {
			delete(p.clients, oldestKey)
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
	p.clients[key] = &pooledTLSClient{
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

// ResolveClient returns a BoringSSL/Bun uTLS *http.Client preserving configured
// proxy-pool routing (when configured uses *http.Transport) and preserving
// custom unit-test RoundTrippers verbatim.
func ResolveClient(ctx context.Context, configured *http.Client, targetURL string) *http.Client {
	if configured == nil {
		return httpclient.SharedChromeTLSClient()
	}
	if configured.Transport == nil {
		return defaultTLSPool.getOrCreate(directProxyPoolKey, nil, configured.Timeout)
	}
	stdTransport, isStd := configured.Transport.(*http.Transport)
	if !isStd {
		// Preserve custom mock/stub RoundTrippers injected by tests.
		return configured
	}
	var resolvedProxy *url.URL
	if stdTransport.Proxy != nil {
		if parsed, err := url.Parse(targetURL); err == nil && parsed != nil {
			probeReq := &http.Request{URL: parsed}
			if ctx != nil {
				probeReq = probeReq.WithContext(ctx)
			}
			if u, err := stdTransport.Proxy(probeReq); err == nil && u != nil {
				resolvedProxy = u
			}
		}
	}
	proxyKey := directProxyPoolKey
	if resolvedProxy != nil {
		proxyKey = resolvedProxy.String()
	}
	return defaultTLSPool.getOrCreate(proxyKey, resolvedProxy, configured.Timeout)
}

// ResetClient closes idle connections and evicts the pooled uTLS client for
// targetURL's proxy key so the next request opens a fresh TCP/TLS socket
// (e.g. after a 429 rate-limit response).
func ResetClient(ctx context.Context, configured *http.Client, targetURL string) {
	proxyKey := directProxyPoolKey
	if configured != nil {
		if stdTransport, isStd := configured.Transport.(*http.Transport); isStd && stdTransport.Proxy != nil {
			if parsed, err := url.Parse(targetURL); err == nil && parsed != nil {
				probeReq := &http.Request{URL: parsed}
				if ctx != nil {
					probeReq = probeReq.WithContext(ctx)
				}
				if u, err := stdTransport.Proxy(probeReq); err == nil && u != nil {
					proxyKey = u.String()
				}
			}
		}
	}
	defaultTLSPool.mu.Lock()
	entry, ok := defaultTLSPool.clients[proxyKey]
	if ok {
		delete(defaultTLSPool.clients, proxyKey)
	}
	defaultTLSPool.mu.Unlock()
	if ok && entry.client != nil {
		entry.client.CloseIdleConnections()
	}
}

// ResolveSystemOneURL returns the OpenCode System One endpoint URL for baseURL.
func ResolveSystemOneURL(baseURL string) string {
	trimmed := strings.TrimRight(strings.TrimSpace(baseURL), "/")
	if trimmed == "" {
		return DefaultSystemOneURL
	}
	if strings.HasSuffix(trimmed, "/systemone") {
		return trimmed
	}
	if parsed, err := url.Parse(trimmed); err == nil && parsed.Scheme != "" && parsed.Host != "" {
		return parsed.Scheme + "://" + parsed.Host + "/inference/v1/systemone"
	}
	return DefaultSystemOneURL
}
