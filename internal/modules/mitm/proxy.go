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
	"slices"
	"strconv"
	"strings"
	"time"
)

const (
	maxProxyBodyBytes = 8 << 20 // 8 MiB
	streamCopyBufSize = 32 << 10
)

// ErrAlreadyRunning indicates the MITM proxy is already active.
var ErrAlreadyRunning = errors.New("mitm: proxy is already running")

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

var hopByHopHeaders = map[string]struct{}{
	"Connection":          {},
	"Proxy-Connection":    {},
	"Keep-Alive":          {},
	"Proxy-Authenticate":  {},
	"Proxy-Authorization": {},
	"Te":                  {},
	"Trailer":             {},
	"Transfer-Encoding":   {},
	"Upgrade":             {},
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
	host, portStr, err := net.SplitHostPort(trimmed)
	if err != nil {
		return fmt.Errorf("listen_addr must be host:port: %w", err)
	}
	port, err := strconv.Atoi(portStr)
	if err != nil || port < 0 || port > 65535 {
		return fmt.Errorf("listen_addr port %q must be an integer between 0 and 65535", portStr)
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
	if portStr := parsed.Port(); portStr != "" {
		port, err := strconv.Atoi(portStr)
		if err != nil || port < 1 || port > 65535 {
			return fmt.Errorf("target_gateway_url port %q must be between 1 and 65535", portStr)
		}
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
	m.opMu.Lock()
	defer m.opMu.Unlock()

	m.mu.RLock()
	alreadyRunning := m.state == StateRunning && m.srv != nil
	m.mu.RUnlock()
	if alreadyRunning {
		return "", ErrAlreadyRunning
	}

	if err := validateLoopbackListenAddr(cfg.ListenAddr); err != nil {
		return "", err
	}
	if err := validateLoopbackURL(cfg.TargetGatewayURL); err != nil {
		return "", err
	}
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
		CheckRedirect: func(*http.Request, []*http.Request) error {
			return http.ErrUseLastResponse
		},
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
			proxyClient.CloseIdleConnections()
			m.mu.Lock()
			if m.srv == srv {
				m.srv = nil
				m.listener = nil
				m.proxyClient = nil
				m.serveDone = nil
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
	m.opMu.Lock()
	defer m.opMu.Unlock()

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
	if err != nil {
		_ = srv.Close()
	}
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

		toolActive := tool == "generic" || m.isToolEnabled(r.Context(), tool)
		sourceModel, targetModel, rewrittenBody := "", "", rawBody
		if toolActive {
			sourceModel, targetModel, rewrittenBody = m.rewriteModelIfNeeded(r.Context(), tool, rawBody)
		}

		outURL := *targetBase
		outURL.Path = singleJoiningSlash(targetBase.Path, r.URL.Path)
		outURL.RawQuery = r.URL.RawQuery

		var bodyReader io.Reader
		if len(rewrittenBody) > 0 {
			bodyReader = bytes.NewReader(rewrittenBody)
		}
		outReq, err := http.NewRequestWithContext(r.Context(), r.Method, outURL.String(), bodyReader)
		if err != nil {
			http.Error(w, "failed to build upstream request", http.StatusBadGateway)
			return
		}
		copyNonHopHeaders(outReq.Header, r.Header)
		outReq.Header.Del("Host")
		outReq.Header.Set("X-Mitm-Tool", tool)
		outReq.Header.Set("X-Mitm-Original-Host", r.Host)
		outReq.ContentLength = int64(len(rewrittenBody))

		resp, err := client.Do(outReq) //nolint:gosec // targetBase is strictly validated as a loopback gateway URL and CheckRedirect refuses redirects
		if err != nil {
			m.recordInterceptEvent(r.Context(), tool, r.Host, r.Method, r.URL.Path, sourceModel, targetModel, http.StatusBadGateway, time.Since(start))
			http.Error(w, "upstream gateway unreachable", http.StatusBadGateway)
			return
		}
		defer func() { _ = resp.Body.Close() }()

		m.recordInterceptEvent(r.Context(), tool, r.Host, r.Method, r.URL.Path, sourceModel, targetModel, resp.StatusCode, time.Since(start))

		copyNonHopHeaders(w.Header(), resp.Header)
		w.WriteHeader(resp.StatusCode)
		streamResponseBody(w, resp.Body)
	})
}

func (m *Module) isToolEnabled(ctx context.Context, tool string) bool {
	cfg, err := m.repo.getConfig(ctx)
	if err != nil {
		return true
	}
	return slices.Contains(cfg.EnabledTools, tool)
}

func streamResponseBody(w http.ResponseWriter, body io.Reader) {
	rc := http.NewResponseController(w)
	buf := make([]byte, streamCopyBufSize)
	for {
		n, readErr := body.Read(buf)
		if n > 0 {
			if _, writeErr := w.Write(buf[:n]); writeErr != nil {
				return
			}
			_ = rc.Flush()
		}
		if readErr != nil {
			return
		}
	}
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
	evCtx, cancel := context.WithTimeout(context.WithoutCancel(ctx), 2*time.Second)
	defer cancel()

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
	if err := m.repo.recordEvent(evCtx, ev); err != nil {
		m.logger.Warn("mitm: failed to record intercept event", "error", err)
	}
}

func copyNonHopHeaders(dst, src http.Header) {
	extraHop := make(map[string]struct{})
	for _, connVal := range src.Values("Connection") {
		for token := range strings.SplitSeq(connVal, ",") {
			if name := http.CanonicalHeaderKey(strings.TrimSpace(token)); name != "" {
				extraHop[name] = struct{}{}
			}
		}
	}
	for k, vals := range src {
		ck := http.CanonicalHeaderKey(k)
		if _, isHop := hopByHopHeaders[ck]; isHop {
			continue
		}
		if _, isExtra := extraHop[ck]; isExtra {
			continue
		}
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
