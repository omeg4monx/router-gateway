// Package contracts defines narrow, domain-agnostic interfaces and lightweight
// value types for cross-module communication so feature modules never import
// each other directly.
package contracts

import (
	"context"
	"net/http"
	"net/url"
	"time"
)

// AuthVerifier validates administrative dashboard requests and reports whether
// dashboard authentication is currently enforced.
type AuthVerifier interface {
	AuthenticateAdmin(r *http.Request) (string, error)
	RequireLogin(ctx context.Context) bool
}

// KeyInfo describes an authenticated gateway API key and its access limits.
type KeyInfo struct {
	ID            string
	Name          string
	AllowedModels []string
	Scopes        []string
	Tags          []string
	RateLimitRPM  int
	QuotaUSD      float64
	Disabled      bool
}

// KeyValidator authenticates raw gateway API keys and returns their metadata.
type KeyValidator interface {
	ValidateKey(ctx context.Context, rawKey string) (KeyInfo, error)
}

// ResolvedTarget identifies a concrete upstream provider connection candidate
// returned by ProviderResolver in fallback priority order.
type ResolvedTarget struct {
	ProviderID    string
	ProviderType  string
	ConnectionID  string
	Model         string
	UpstreamModel string
	Endpoint      string
	Priority      int
	Metadata      map[string]string
}

// ProviderResolver resolves a requested model identifier (combo, alias,
// provider/model, or bare model) into an ordered list of upstream targets
// permitted by the caller's allowed model rules.
type ProviderResolver interface {
	Resolve(ctx context.Context, model string, allowed []string) ([]ResolvedTarget, error)
}

// CanonicalRequest carries normalized inbound gateway request state through
// registered GatewayInterceptor hooks.
type CanonicalRequest struct {
	RequestID string
	Protocol  string
	Model     string
	Stream    bool
	KeyID     string
	Headers   http.Header
	Body      []byte
	Metadata  map[string]any
}

// ResponseMeta carries upstream response header metadata through registered
// GatewayInterceptor hooks before the response body is streamed to the client.
type ResponseMeta struct {
	RequestID    string
	StatusCode   int
	Headers      http.Header
	ProviderID   string
	ConnectionID string
	Model        string
	Metadata     map[string]any
}

// StreamEvent represents a single streaming chunk passed through registered
// GatewayInterceptor hooks during SSE forwarding.
type StreamEvent struct {
	RequestID string
	Event     string
	Data      []byte
	Done      bool
	Metadata  map[string]any
}

// GatewayInterceptor allows extension modules to inspect or mutate gateway
// requests, response headers, and stream chunks in deterministic priority order.
type GatewayInterceptor interface {
	Name() string
	Priority() int
	OnRequest(ctx context.Context, req *CanonicalRequest) error
	OnResponseHeaders(ctx context.Context, meta *ResponseMeta) error
	OnStreamChunk(ctx context.Context, ev *StreamEvent) error
}

// UsageRecord captures per-request telemetry emitted asynchronously by the
// gateway to the usage module.
type UsageRecord struct {
	RequestID    string
	KeyID        string
	Model        string
	ProviderID   string
	ConnectionID string
	Protocol     string
	InputTokens  int64
	OutputTokens int64
	CachedTokens int64
	CostUSD      float64
	Latency      time.Duration
	TTFT         time.Duration
	StatusCode   int
	ErrorSummary string
	Stream       bool
	Timestamp    time.Time
}

// UsageRecorder records gateway request usage asynchronously without blocking
// response streaming.
type UsageRecorder interface {
	Record(ctx context.Context, rec UsageRecord)
}

// ProxyDialerResolver resolves an optional outbound proxy URL for a given
// upstream provider. A nil *url.URL with a nil error indicates direct dial.
type ProxyDialerResolver interface {
	ProxyURLForProvider(ctx context.Context, providerID string) (*url.URL, error)
}
