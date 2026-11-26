// Package httpclient constructs pooled HTTP/2 clients suited for both
// request/response calls and long-lived Server-Sent Events (SSE) streams.
package httpclient

import (
	"context"
	"net"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/omeg4monx/router-gateway/internal/contracts"
)

type providerIDContextKey struct{}

// WithProviderID tags ctx with the target upstream provider identifier so
// NewProviderProxyFunc can route requests through provider-bound proxy pools.
func WithProviderID(ctx context.Context, providerID string) context.Context {
	if ctx == nil {
		ctx = context.Background()
	}
	trimmed := strings.TrimSpace(providerID)
	if trimmed == "" {
		return ctx
	}
	return context.WithValue(ctx, providerIDContextKey{}, trimmed)
}

// ProviderIDFromContext extracts the upstream provider identifier previously
// attached via WithProviderID, or returns an empty string if unset.
func ProviderIDFromContext(ctx context.Context) string {
	if ctx == nil {
		return ""
	}
	v, _ := ctx.Value(providerIDContextKey{}).(string)
	return v
}

// NewProviderProxyFunc returns an http.Transport Proxy function that resolves
// provider-bound outbound proxies via resolver when a request context carries a
// provider ID, falling back to http.ProxyFromEnvironment when unbound.
func NewProviderProxyFunc(resolver contracts.ProxyDialerResolver) func(*http.Request) (*url.URL, error) {
	return func(req *http.Request) (*url.URL, error) {
		if req != nil && resolver != nil {
			if providerID := ProviderIDFromContext(req.Context()); providerID != "" {
				proxyURL, err := resolver.ProxyURLForProvider(req.Context(), providerID)
				if err != nil {
					return nil, err
				}
				if proxyURL != nil {
					return proxyURL, nil
				}
			}
		}
		return http.ProxyFromEnvironment(req)
	}
}

const (
	defaultDialTimeout           = 10 * time.Second
	defaultResponseHeaderTimeout = 30 * time.Second
	defaultIdleConnTimeout       = 90 * time.Second
	defaultMaxIdleConns          = 256
	defaultMaxIdleConnsPerHost   = 64
	defaultMaxConnsPerHost       = 128
	defaultKeepAlive             = 30 * time.Second
	defaultTLSHandshakeTimeout   = 10 * time.Second
	defaultExpectContinueTimeout = 1 * time.Second
)

// Options configures connection pooling and header timeouts for NewClient.
// Zero or negative values fall back to production defaults.
type Options struct {
	DialTimeout           time.Duration
	ResponseHeaderTimeout time.Duration
	IdleConnTimeout       time.Duration
	MaxIdleConns          int
	MaxIdleConnsPerHost   int
	MaxConnsPerHost       int
	ProxyFunc             func(*http.Request) (*url.URL, error)
}

func applyDefaults(opts Options) Options {
	if opts.DialTimeout <= 0 {
		opts.DialTimeout = defaultDialTimeout
	}
	if opts.ResponseHeaderTimeout <= 0 {
		opts.ResponseHeaderTimeout = defaultResponseHeaderTimeout
	}
	if opts.IdleConnTimeout <= 0 {
		opts.IdleConnTimeout = defaultIdleConnTimeout
	}
	if opts.MaxIdleConns <= 0 {
		opts.MaxIdleConns = defaultMaxIdleConns
	}
	if opts.MaxIdleConnsPerHost <= 0 {
		opts.MaxIdleConnsPerHost = defaultMaxIdleConnsPerHost
	}
	if opts.MaxConnsPerHost <= 0 {
		opts.MaxConnsPerHost = defaultMaxConnsPerHost
	}
	if opts.ProxyFunc == nil {
		opts.ProxyFunc = http.ProxyFromEnvironment
	}
	return opts
}

// NewClient returns an *http.Client backed by a pooled HTTP/2-capable
// *http.Transport. Client.Timeout is intentionally left at 0 so long-lived
// streaming responses are bounded only by request context cancellation and
// ResponseHeaderTimeout.
func NewClient(opts Options) *http.Client {
	cfg := applyDefaults(opts)

	dialer := &net.Dialer{
		Timeout:   cfg.DialTimeout,
		KeepAlive: defaultKeepAlive,
	}

	transport := &http.Transport{
		Proxy:                 cfg.ProxyFunc,
		DialContext:           dialer.DialContext,
		ForceAttemptHTTP2:     true,
		MaxIdleConns:          cfg.MaxIdleConns,
		MaxIdleConnsPerHost:   cfg.MaxIdleConnsPerHost,
		MaxConnsPerHost:       cfg.MaxConnsPerHost,
		IdleConnTimeout:       cfg.IdleConnTimeout,
		TLSHandshakeTimeout:   defaultTLSHandshakeTimeout,
		ExpectContinueTimeout: defaultExpectContinueTimeout,
		ResponseHeaderTimeout: cfg.ResponseHeaderTimeout,
	}

	return &http.Client{
		Transport: transport,
		Timeout:   0,
	}
}
