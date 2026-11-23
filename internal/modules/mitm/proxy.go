package mitm

import (
	"bytes"
	"context"
	"crypto/rand"
	"crypto/tls"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net"
	"net/http"
	"net/url"
	"strings"
	"time"
)

const maxProxyBodyBytes = 8 << 20 // 8 MiB

// ToolDescriptor describes an IDE traffic interception profile and its hosts.
type ToolDescriptor struct {
	ID      string   `json:"id"`
	Name    string   `json:"name"`
	Hosts   []string `json:"hosts"`
	Enabled bool     `json:"enabled"`
}

var toolCatalog = []struct {
	id    string
	name  string
	hosts []string
}{
	{
		id:    "antigravity",
		name:  "Antigravity IDE",
		hosts: []string{"daily-cloudcode-pa.googleapis.com", "cloudcode-pa.googleapis.com"},
	},
	{
		id:    "copilot",
		name:  "GitHub Copilot",
		hosts: []string{"api.individual.githubcopilot.com"},
	},
	{
		id:    "kiro",
		name:  "Kiro IDE",
		hosts: []string{"runtime.us-east-1.kiro.dev", "q.us-east-1.amazonaws.com", "codewhisperer.us-east-1.amazonaws.com"},
	},
	{
		id:    "cursor",
		name:  "Cursor IDE",
		hosts: []string{"api2.cursor.sh"},
	},
}

func validToolID(id string) bool {
	for _, t := range toolCatalog {
		if t.id == id {
			return true
		}
	}
	return false
}

func matchToolForHost(host string) string {
	clean := strings.ToLower(strings.TrimSpace(host))
	if h, _, err := net.SplitHostPort(clean); err == nil && h != "" {
		clean = h
	}
	for _, t := range toolCatalog {
		for _, candidate := range t.hosts {
			if clean == candidate {
				return t.id
			}
		}
	}
	return "generic"
}

func validateLoopbackListenAddr(addr string) error {
	trimmed := strings.TrimSpace(addr)
	if trimmed == "" {
		return errors.New("listen_addr must not be empty")
	}
	host, port, err := net.SplitHostPort(trimmed)
	if err != nil {
		return fmt.Errorf("listen_addr must be host:port: %w", err)
	}
	if port == "" {
		return errors.New("listen_addr port must not be empty")
	}
	if !isLoopbackHost(host) {
		return fmt.Errorf("listen_addr host %q must be loopback (127.0.0.1, localhost, or ::1)", host)
	}
	return nil
}

func validateLoopbackURL(raw string) error {
	trimmed := strings.TrimSpace(raw)
	if trimmed == "" {
		return errors.New("target_gateway_url must not be empty")
	}
	parsed, err := url.Parse(trimmed)
	if err != nil {
		return fmt.Errorf("invalid target_gateway_url: %w", err)
	}
	if parsed.Scheme != "http" && parsed.Scheme != "https" {
		return fmt.Errorf("target_gateway_url scheme must be http or https, got %q", parsed.Scheme)
	}
	if !isLoopbackHost(parsed.Hostname()) {
		return fmt.Errorf("target_gateway_url host %q must be loopback (127.0.0.1, localhost, or ::1)", parsed.Hostname())
	}
	return nil
}

func isLoopbackHost(host string) bool {
	h := strings.ToLower(strings.TrimSpace(strings.Trim(host, "[]")))
	if h == "localhost" {
		return true
	}
	ip := net.ParseIP(h)
	return ip != nil && ip.IsLoopback()
}

func (m *Module) startProxyServer(cfg Config) (string, error) {
	if _, err := m.certs.ensureCA(false); err != nil {
		return "", fmt.Errorf("ensure root CA: %w", err)
	}
	targetURL, err := url.Parse(cfg.TargetGatewayURL)
	if err != nil {
		return "", fmt.Errorf("parse target gateway URL: %w", err)
	}

	rawListener, err := net.Listen("tcp", cfg.ListenAddr)
	if err != nil {
		return "", fmt.Errorf("listen on %s: %w", cfg.ListenAddr, err)
	}
	tlsListener := tls.NewListener(rawListener, &tls.Config{
		GetCertificate: m.certs.getCertificate,
		MinVersion:     tls.VersionTLS12,
	})

	proxyClient := &http.Client{
		Timeout: 0,
		Transport: &http.Transport{
			MaxIdleConns:          32,
			IdleConnTimeout:       60 * time.Second,
			ResponseHeaderTimeout: 60 * time.Second,
		},
	}

	srv := &http.Server{
		Handler:           m.buildProxyHandler(targetURL, proxyClient),
		ReadHeaderTimeout: 10 * time.Second,
		IdleTimeout:       120 * time.Second,
	}

	boundAddr := rawListener.Addr().String()
	done := make(chan struct{})

	m.mu.Lock()
	m.srv = srv
	m.listener = tlsListener
	m.proxyClient = proxyClient
	m.serveDone = done
	m.state = StateRunning
	m.boundAddr = boundAddr
	m.startedAt = time.Now().UTC().Format(time.RFC3339)
	m.lastError = ""
	m.mu.Unlock()

	go func() {
		defer close(done)
		if serveErr := srv.Serve(tlsListener); serveErr != nil && !errors.Is(serveErr, http.ErrServerClosed) {
			m.mu.Lock()
			if m.srv == srv {
				m.state = StateError
				m.lastError = serveErr.Error()
				m.boundAddr = ""
			}
			m.mu.Unlock()
			m.logger.Error("mitm proxy server exited with error", "error", serveErr)
		}
	}()

	return boundAddr, nil
}

func (m *Module) stopProxyServer(ctx context.Context) error {
	m.mu.Lock()
	srv := m.srv
	done := m.serveDone
	client := m.proxyClient
	m.srv = nil
	m.listener = nil
	m.proxyClient = nil
	m.serveDone = nil
	m.state = StateStopped
	m.boundAddr = ""
	m.startedAt = ""
	m.mu.Unlock()

	if srv == nil {
		return nil
	}

	shutdownCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()
	err := srv.Shutdown(shutdownCtx)
	if done != nil {
		<-done
	}
	if client != nil {
		client.CloseIdleConnections()
	}
	return err
}

func (m *Module) buildProxyHandler(targetBase *url.URL, client *http.Client) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		tool := matchToolForHost(r.Host)

		var rawBody []byte
		if r.Body != nil {
			limited := io.LimitReader(r.Body, maxProxyBodyBytes+1)
			b, err := io.ReadAll(limited)
			if err != nil {
				http.Error(w, "failed to read request body", http.StatusBadRequest)
				return
			}
			if len(b) > maxProxyBodyBytes {
				http.Error(w, "request body too large", http.StatusRequestEntityTooLarge)
				return
			}
			rawBody = b
		}

		sourceModel, targetModel, rewrittenBody := m.rewriteModelIfNeeded(r.Context(), tool, rawBody)

		outURL := *targetBase
		outURL.Path = singleJoiningSlash(targetBase.Path, r.URL.Path)
		outURL.RawQuery = r.URL.RawQuery

		outReq, err := http.NewRequestWithContext(r.Context(), r.Method, outURL.String(), bytes.NewReader(rewrittenBody))
		if err != nil {
			http.Error(w, "failed to build upstream request", http.StatusBadGateway)
			return
		}
		copyHeaders(outReq.Header, r.Header)
		outReq.Header.Del("Host")
		outReq.Header.Set("X-Mitm-Tool", tool)
		outReq.Header.Set("X-Mitm-Original-Host", r.Host)
		outReq.ContentLength = int64(len(rewrittenBody))

		resp, err := client.Do(outReq) //nolint:gosec // targetBase is strictly validated as a loopback gateway URL
		if err != nil {
			m.recordInterceptEvent(r.Context(), tool, r.Host, r.Method, r.URL.Path, sourceModel, targetModel, http.StatusBadGateway, time.Since(start))
			http.Error(w, "upstream gateway unreachable", http.StatusBadGateway)
			return
		}
		defer func() { _ = resp.Body.Close() }()

		copyHeaders(w.Header(), resp.Header)
		w.WriteHeader(resp.StatusCode)
		_, _ = io.Copy(w, resp.Body)

		m.recordInterceptEvent(r.Context(), tool, r.Host, r.Method, r.URL.Path, sourceModel, targetModel, resp.StatusCode, time.Since(start))
	})
}

func (m *Module) rewriteModelIfNeeded(ctx context.Context, tool string, body []byte) (sourceModel, targetModel string, out []byte) {
	if len(body) == 0 {
		return "", "", body
	}
	var payload map[string]any
	if err := json.Unmarshal(body, &payload); err != nil {
		return "", "", body
	}
	src, ok := payload["model"].(string)
	if !ok || strings.TrimSpace(src) == "" {
		return "", "", body
	}
	target, matched := m.repo.resolveAlias(ctx, tool, src)
	if !matched || target == "" {
		return src, src, body
	}
	payload["model"] = target
	encoded, err := json.Marshal(payload)
	if err != nil {
		return src, src, body
	}
	return src, target, encoded
}

func (m *Module) recordInterceptEvent(ctx context.Context, tool, host, method, path, srcModel, tgtModel string, statusCode int, latency time.Duration) {
	ev := Event{
		ID:          newID("mev"),
		Tool:        tool,
		Host:        host,
		Method:      method,
		Path:        path,
		SourceModel: srcModel,
		TargetModel: tgtModel,
		StatusCode:  statusCode,
		LatencyMs:   latency.Milliseconds(),
		CreatedAt:   time.Now().UTC().Format(timestampLayout),
	}
	_ = m.repo.recordEvent(ctx, ev)
}

func copyHeaders(dst, src http.Header) {
	for k, vals := range src {
		for _, v := range vals {
			dst.Add(k, v)
		}
	}
}

func singleJoiningSlash(a, b string) string {
	aslash := strings.HasSuffix(a, "/")
	bslash := strings.HasPrefix(b, "/")
	switch {
	case aslash && bslash:
		return a + b[1:]
	case !aslash && !bslash:
		return a + "/" + b
	default:
		return a + b
	}
}

func newID(prefix string) string {
	var b [8]byte
	_, _ = rand.Read(b[:])
	return prefix + "_" + hex.EncodeToString(b[:])
}
