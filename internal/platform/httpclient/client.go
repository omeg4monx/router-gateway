// Package httpclient constructs pooled HTTP/2 clients suited for both
// request/response calls and long-lived Server-Sent Events (SSE) streams.
package httpclient

import (
	"net"
	"net/http"
	"net/url"
	"time"
)

const (
	defaultDialTimeout           = 10 * time.Second
	defaultResponseHeaderTimeout = 30 * time.Second
	defaultIdleConnTimeout       = 90 * time.Second
	defaultMaxIdleConns          = 100
	defaultMaxIdleConnsPerHost   = 20
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
