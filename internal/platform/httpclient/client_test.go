package httpclient_test

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"io"
	"net"
	"net/http"
	"net/http/httptest"
	"net/url"
	"sync/atomic"
	"testing"
	"time"

	"github.com/omeg4monx/router-gateway/internal/platform/httpclient"
)

func TestNewClientAllowsLongStreamButEnforcesResponseHeaderTimeout(t *testing.T) {
	t.Parallel()

	client := httpclient.NewClient(httpclient.Options{
		DialTimeout:           2 * time.Second,
		ResponseHeaderTimeout: 50 * time.Millisecond,
		IdleConnTimeout:       5 * time.Second,
		MaxIdleConns:          10,
		MaxIdleConnsPerHost:   5,
	})

	if client.Timeout != 0 {
		t.Fatalf("expected client.Timeout == 0 for streaming, got %v", client.Timeout)
	}

	t.Run("enforces ResponseHeaderTimeout when headers are delayed", func(t *testing.T) {
		t.Parallel()

		srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			select {
			case <-time.After(150 * time.Millisecond):
				w.WriteHeader(http.StatusOK)
			case <-r.Context().Done():
			}
		}))
		defer srv.Close()

		req, err := http.NewRequestWithContext(context.Background(), http.MethodGet, srv.URL, nil)
		if err != nil {
			t.Fatalf("new request: %v", err)
		}

		resp, err := client.Do(req)
		if err == nil {
			resp.Body.Close()
			t.Fatal("expected ResponseHeaderTimeout error, got nil")
		}

		var netErr net.Error
		if !errors.As(err, &netErr) || !netErr.Timeout() {
			t.Fatalf("expected timeout net.Error, got %T: %v", err, err)
		}
	})

	t.Run("allows slow chunked stream longer than ResponseHeaderTimeout", func(t *testing.T) {
		t.Parallel()

		srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			flusher, ok := w.(http.Flusher)
			if !ok {
				http.Error(w, "streaming unsupported", http.StatusInternalServerError)
				return
			}
			w.Header().Set("Content-Type", "text/event-stream")
			w.WriteHeader(http.StatusOK)
			flusher.Flush()

			fmt.Fprint(w, "data: chunk1\n\n")
			flusher.Flush()

			select {
			case <-time.After(80 * time.Millisecond):
			case <-r.Context().Done():
				return
			}

			fmt.Fprint(w, "data: chunk2\n\n")
			flusher.Flush()
		}))
		defer srv.Close()

		req, err := http.NewRequestWithContext(context.Background(), http.MethodGet, srv.URL, nil)
		if err != nil {
			t.Fatalf("new request: %v", err)
		}

		resp, err := client.Do(req)
		if err != nil {
			t.Fatalf("expected streaming request to succeed, got %v", err)
		}
		defer resp.Body.Close()

		body, err := io.ReadAll(resp.Body)
		if err != nil {
			t.Fatalf("read stream body: %v", err)
		}

		want := "data: chunk1\n\ndata: chunk2\n\n"
		if !bytes.Equal(body, []byte(want)) {
			t.Fatalf("body = %q, want %q", string(body), want)
		}
	})
}

func TestNewClientAppliesDefaults(t *testing.T) {
	t.Parallel()

	client := httpclient.NewClient(httpclient.Options{})
	if client.Timeout != 0 {
		t.Fatalf("Timeout = %v, want 0", client.Timeout)
	}

	tr, ok := client.Transport.(*http.Transport)
	if !ok {
		t.Fatalf("Transport type = %T, want *http.Transport", client.Transport)
	}

	if !tr.ForceAttemptHTTP2 {
		t.Fatal("ForceAttemptHTTP2 = false, want true")
	}
	if tr.ResponseHeaderTimeout != 30*time.Second {
		t.Fatalf("ResponseHeaderTimeout = %v, want 30s", tr.ResponseHeaderTimeout)
	}
	if tr.IdleConnTimeout != 90*time.Second {
		t.Fatalf("IdleConnTimeout = %v, want 90s", tr.IdleConnTimeout)
	}
	if tr.MaxIdleConns != 256 {
		t.Fatalf("MaxIdleConns = %d, want 256", tr.MaxIdleConns)
	}
	if tr.MaxIdleConnsPerHost != 64 {
		t.Fatalf("MaxIdleConnsPerHost = %d, want 64", tr.MaxIdleConnsPerHost)
	}
	if tr.MaxConnsPerHost != 128 {
		t.Fatalf("MaxConnsPerHost = %d, want 128", tr.MaxConnsPerHost)
	}
	if tr.Proxy == nil {
		t.Fatal("Proxy = nil, want default http.ProxyFromEnvironment")
	}
	if tr.DialContext == nil {
		t.Fatal("DialContext = nil, want configured dialer")
	}
}

func TestNewClientUsesCustomProxyFunc(t *testing.T) {
	t.Parallel()

	proxyURL, err := url.Parse("http://proxy.internal:8080")
	if err != nil {
		t.Fatalf("parse proxy url: %v", err)
	}

	var called atomic.Bool
	client := httpclient.NewClient(httpclient.Options{
		ProxyFunc: func(_ *http.Request) (*url.URL, error) {
			called.Store(true)
			return proxyURL, nil
		},
	})

	tr, ok := client.Transport.(*http.Transport)
	if !ok {
		t.Fatalf("Transport type = %T, want *http.Transport", client.Transport)
	}

	req, err := http.NewRequest(http.MethodGet, "https://example.com", nil)
	if err != nil {
		t.Fatalf("new request: %v", err)
	}

	got, err := tr.Proxy(req)
	if err != nil {
		t.Fatalf("Proxy returned error: %v", err)
	}
	if !called.Load() {
		t.Fatal("expected custom ProxyFunc to be called")
	}
	if got.String() != proxyURL.String() {
		t.Fatalf("proxy URL = %q, want %q", got.String(), proxyURL.String())
	}
}

func TestNewClientRequestContextCancellationAbortsStream(t *testing.T) {
	t.Parallel()

	started := make(chan struct{})
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		flusher, _ := w.(http.Flusher)
		w.WriteHeader(http.StatusOK)
		if flusher != nil {
			flusher.Flush()
		}
		close(started)
		<-r.Context().Done()
	}))
	defer srv.Close()

	client := httpclient.NewClient(httpclient.Options{
		ResponseHeaderTimeout: 2 * time.Second,
	})

	ctx, cancel := context.WithCancel(context.Background())
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, srv.URL, nil)
	if err != nil {
		t.Fatalf("new request: %v", err)
	}

	resp, err := client.Do(req)
	if err != nil {
		t.Fatalf("Do: %v", err)
	}
	defer resp.Body.Close()

	<-started
	cancel()

	_, err = io.ReadAll(resp.Body)
	if !errors.Is(err, context.Canceled) {
		t.Fatalf("expected context.Canceled on stream read, got %v", err)
	}
}

type stubProxyResolver struct {
	byProvider map[string]*url.URL
	err        error
}

func (s *stubProxyResolver) ProxyURLForProvider(_ context.Context, providerID string) (*url.URL, error) {
	if s.err != nil {
		return nil, s.err
	}
	return s.byProvider[providerID], nil
}

func TestProviderProxyFuncContextRouting(t *testing.T) {
	t.Parallel()

	wantProxy, err := url.Parse("http://proxy.bound.local:8080")
	if err != nil {
		t.Fatalf("parse proxy url: %v", err)
	}

	resolver := &stubProxyResolver{
		byProvider: map[string]*url.URL{
			"openai": wantProxy,
		},
	}
	proxyFn := httpclient.NewProviderProxyFunc(resolver)

	// Bound provider resolves to bound proxy URL.
	ctx := httpclient.WithProviderID(context.Background(), "openai")
	if got := httpclient.ProviderIDFromContext(ctx); got != "openai" {
		t.Fatalf("ProviderIDFromContext = %q, want %q", got, "openai")
	}
	reqBound, err := http.NewRequestWithContext(ctx, http.MethodGet, "https://api.openai.com/v1/models", nil)
	if err != nil {
		t.Fatalf("new bound req: %v", err)
	}
	gotURL, err := proxyFn(reqBound)
	if err != nil {
		t.Fatalf("proxyFn(bound) err: %v", err)
	}
	if gotURL == nil || gotURL.String() != wantProxy.String() {
		t.Fatalf("proxyFn(bound) = %v, want %v", gotURL, wantProxy)
	}

	// Unbound provider returns nil (no error).
	ctxUnbound := httpclient.WithProviderID(context.Background(), "anthropic")
	reqUnbound, err := http.NewRequestWithContext(ctxUnbound, http.MethodGet, "https://api.anthropic.com/v1/models", nil)
	if err != nil {
		t.Fatalf("new unbound req: %v", err)
	}
	gotUnbound, err := proxyFn(reqUnbound)
	if err != nil {
		t.Fatalf("proxyFn(unbound) err: %v", err)
	}
	if gotUnbound != nil {
		t.Fatalf("proxyFn(unbound) = %v, want nil", gotUnbound)
	}

	// Untagged request falls back cleanly without error.
	reqPlain, err := http.NewRequestWithContext(context.Background(), http.MethodGet, "https://example.com", nil)
	if err != nil {
		t.Fatalf("new plain req: %v", err)
	}
	if _, err := proxyFn(reqPlain); err != nil {
		t.Fatalf("proxyFn(plain) err: %v", err)
	}
}
