package providers

import (
	"context"
	"encoding/json"
	"fmt"
	"html"
	"io"
	"net"
	"net/http"
	"strings"
	"sync"
	"time"
)

const loopbackSessionTTL = 15 * time.Minute

type loopbackSessionEntry struct {
	code      string
	errMsg    string
	done      bool
	createdAt time.Time
}

type loopbackProviderServer struct {
	providerID   string
	callbackPath string
	port         int
	server       *http.Server
	sessions     map[string]*loopbackSessionEntry
}

// oauthLoopbackProxy manages per-provider local loopback OAuth callback servers
// (fixed-port such as Codex :1455/auth/callback, or dynamic-port when port == 0),
// matching xrouter's unified callback proxy pattern.
type oauthLoopbackProxy struct {
	mu      sync.Mutex
	servers map[string]*loopbackProviderServer
}

func newOAuthLoopbackProxy() *oauthLoopbackProxy {
	return &oauthLoopbackProxy{
		servers: make(map[string]*loopbackProviderServer),
	}
}

func (p *oauthLoopbackProxy) getOrCreateLocked(providerID, callbackPath string) *loopbackProviderServer {
	if callbackPath == "" {
		callbackPath = "/callback"
	}
	if !strings.HasPrefix(callbackPath, "/") {
		callbackPath = "/" + callbackPath
	}
	ps, ok := p.servers[providerID]
	if !ok {
		ps = &loopbackProviderServer{
			providerID:   providerID,
			callbackPath: callbackPath,
			sessions:     make(map[string]*loopbackSessionEntry),
		}
		p.servers[providerID] = ps
	}
	return ps
}

func (p *oauthLoopbackProxy) registerSession(providerID, callbackPath, state string) {
	p.mu.Lock()
	defer p.mu.Unlock()
	ps := p.getOrCreateLocked(providerID, callbackPath)
	now := time.Now()
	for k, v := range ps.sessions {
		if now.Sub(v.createdAt) > loopbackSessionTTL {
			delete(ps.sessions, k)
		}
	}
	ps.sessions[state] = &loopbackSessionEntry{createdAt: now}
}

// start binds a loopback HTTP listener for providerID (port == 0 picks an ephemeral port)
// and registers the OAuth state for callback capture. If a listener is already running
// for providerID, it is reused.
func (p *oauthLoopbackProxy) start(providerID string, port int, callbackPath, state string) (int, error) {
	p.registerSession(providerID, callbackPath, state)

	p.mu.Lock()
	ps := p.getOrCreateLocked(providerID, callbackPath)
	if ps.server != nil {
		boundPort := ps.port
		p.mu.Unlock()
		return boundPort, nil
	}

	addr := fmt.Sprintf("127.0.0.1:%d", port)
	ln, err := net.Listen("tcp", addr)
	if err != nil {
		p.mu.Unlock()
		return 0, err
	}
	tcpAddr, ok := ln.Addr().(*net.TCPAddr)
	if !ok {
		_ = ln.Close()
		p.mu.Unlock()
		return 0, fmt.Errorf("providers: unexpected loopback listener address %T", ln.Addr())
	}
	actualPort := tcpAddr.Port

	mux := http.NewServeMux()
	mux.HandleFunc(ps.callbackPath, func(w http.ResponseWriter, r *http.Request) {
		p.handleCallback(providerID, w, r)
	})
	srv := &http.Server{
		Handler:           mux,
		ReadHeaderTimeout: 5 * time.Second,
	}
	ps.server = srv
	ps.port = actualPort
	p.mu.Unlock()

	go func() {
		_ = srv.Serve(ln)
	}()
	return actualPort, nil
}

func (p *oauthLoopbackProxy) stop(providerID string) {
	p.mu.Lock()
	ps, ok := p.servers[providerID]
	var srv *http.Server
	if ok {
		srv = ps.server
		ps.server = nil
		ps.port = 0
	}
	p.mu.Unlock()

	if srv != nil {
		ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
		defer cancel()
		_ = srv.Shutdown(ctx)
	}
}

func (p *oauthLoopbackProxy) stopAll() {
	p.mu.Lock()
	var servers []*http.Server
	for _, ps := range p.servers {
		if ps.server != nil {
			servers = append(servers, ps.server)
			ps.server = nil
			ps.port = 0
		}
	}
	p.mu.Unlock()

	for _, srv := range servers {
		ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
		_ = srv.Shutdown(ctx)
		cancel()
	}
}

func (p *oauthLoopbackProxy) status(providerID, state string) (status, errMsg string) {
	p.mu.Lock()
	defer p.mu.Unlock()
	ps, ok := p.servers[providerID]
	if !ok {
		return "waiting", ""
	}
	entry, ok := ps.sessions[state]
	if !ok {
		return "waiting", ""
	}
	if entry.errMsg != "" {
		return "error", entry.errMsg
	}
	if entry.done && entry.code != "" {
		return "done", ""
	}
	return "waiting", ""
}

func (p *oauthLoopbackProxy) consumeCode(providerID, state string) (string, bool) {
	p.mu.Lock()
	defer p.mu.Unlock()
	ps, ok := p.servers[providerID]
	if !ok {
		return "", false
	}
	entry, ok := ps.sessions[state]
	if !ok || !entry.done || entry.code == "" {
		return "", false
	}
	code := entry.code
	delete(ps.sessions, state)
	return code, true
}

func (p *oauthLoopbackProxy) handleCallback(providerID string, w http.ResponseWriter, r *http.Request) {
	q := r.URL.Query()
	code := strings.TrimSpace(q.Get("code"))
	state := strings.TrimSpace(q.Get("state"))
	oauthErr := strings.TrimSpace(q.Get("error"))
	oauthErrDesc := strings.TrimSpace(q.Get("error_description"))

	p.mu.Lock()
	ps := p.getOrCreateLocked(providerID, r.URL.Path)
	entry := ps.sessions[state]
	if entry == nil && state != "" {
		entry = &loopbackSessionEntry{createdAt: time.Now()}
		ps.sessions[state] = entry
	}
	if entry != nil {
		if oauthErr != "" {
			msg := oauthErr
			if oauthErrDesc != "" {
				msg = msg + ": " + oauthErrDesc
			}
			entry.errMsg = msg
		} else if code != "" {
			entry.code = code
			entry.done = true
		}
	}
	p.mu.Unlock()

	payloadBytes, _ := json.Marshal(map[string]any{
		"code":             code,
		"state":            state,
		"error":            oauthErr,
		"errorDescription": oauthErrDesc,
		"providerId":       providerID,
		"timestamp":        time.Now().UnixMilli(),
	})
	safePayload := strings.ReplaceAll(string(payloadBytes), "</", `<\/`)

	w.Header().Set("Cache-Control", "no-store")
	w.Header().Set("Content-Security-Policy", "default-src 'none'; style-src 'unsafe-inline'; script-src 'unsafe-inline'")
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.WriteHeader(http.StatusOK)

	title := "Authorization Received"
	msg := "Returning to OMEG4MON Router dashboard… You can close this window."
	if oauthErr != "" {
		title = "Authorization Error"
		msg = oauthErr
		if oauthErrDesc != "" {
			msg = msg + " — " + oauthErrDesc
		}
	}

	var b strings.Builder
	b.WriteString("<!DOCTYPE html><html lang=\"en\"><head><meta charset=\"utf-8\"><title>")
	b.WriteString(html.EscapeString(title))
	b.WriteString("</title><style>body{font-family:system-ui,sans-serif;max-width:32rem;margin:3rem auto;padding:0 1rem;color:#111827;background:#f8fafc}.card{background:#fff;border:1px solid #e2e8f0;border-radius:16px;padding:1.5rem}</style></head><body><div class=\"card\"><h1>")
	b.WriteString(html.EscapeString(title))
	b.WriteString("</h1><p>")
	b.WriteString(html.EscapeString(msg))
	b.WriteString("</p></div><script>(function(){var d=")
	b.WriteString(safePayload)
	b.WriteString(";try{if(window.opener&&!window.opener.closed){window.opener.postMessage({type:\"oauth_callback\",data:d},\"*\");}}catch(e){}try{var bc=new BroadcastChannel(\"oauth_callback\");bc.postMessage(d);bc.close();}catch(e){}try{localStorage.setItem(\"oauth_callback\",JSON.stringify(d));}catch(e){}if(!d.error&&window.opener){setTimeout(function(){try{window.close();}catch(e){}},600);}})();</script></body></html>")
	_, _ = io.WriteString(w, b.String())
}
