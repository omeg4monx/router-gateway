// Package tunnel implements a Cloudflare Quick Tunnel and Tailscale Serve/Funnel
// lifecycle manager with loopback-only target enforcement and event auditing.
package tunnel

import (
	"bufio"
	"context"
	"embed"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/fs"
	"log/slog"
	"net"
	"net/http"
	"net/url"
	"regexp"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/omeg4monx/router-gateway/internal/platform/config"
	"github.com/omeg4monx/router-gateway/internal/platform/db"
	"github.com/omeg4monx/router-gateway/internal/platform/eventbus"
	"github.com/omeg4monx/router-gateway/internal/platform/httpserver"
	"github.com/omeg4monx/router-gateway/internal/platform/lifecycle"
)

//go:embed migrations/*.sql
var migrationsFS embed.FS

const (
	moduleName          = "tunnel"
	maxBodyBytes        = 64 << 10
	TopicTunnelStatus   = "tunnel.status"
	ProviderCloudflare  = "cloudflare"
	ProviderTailscale   = "tailscale"
	TailscaleModeFunnel = "funnel"
	TailscaleModeServe  = "serve"
	StateStopped        = "stopped"
	StateStarting       = "starting"
	StateRunning        = "running"
	StateError          = "error"
	tailscaleCmdTimeout = 15 * time.Second
)

var (
	reCloudflareURL = regexp.MustCompile(`https://[a-zA-Z0-9-]+\.trycloudflare\.com`)
	reHTTPSURL      = regexp.MustCompile(`https://[^\s"'<>]+`)
)

// BinaryStatus reports whether an external tunnel CLI is available on PATH.
type BinaryStatus struct {
	Available bool   `json:"available"`
	Path      string `json:"path"`
}

// BinariesStatus groups availability checks for supported tunnel CLIs.
type BinariesStatus struct {
	Cloudflared BinaryStatus `json:"cloudflared"`
	Tailscale   BinaryStatus `json:"tailscale"`
}

// Status is the complete live and configured state returned by /api/tunnel/status.
type Status struct {
	Status        string         `json:"status"`
	Provider      string         `json:"provider"`
	PublicURL     string         `json:"public_url"`
	LocalURL      string         `json:"local_url"`
	TailscaleMode string         `json:"tailscale_mode"`
	AutoStart     bool           `json:"auto_start"`
	PID           int            `json:"pid"`
	StartedAt     string         `json:"started_at"`
	LastError     string         `json:"last_error"`
	Binaries      BinariesStatus `json:"binaries"`
}

// Options configures a tunnel Module instance.
type Options struct {
	DB     *db.DB
	Config config.Config
	Bus    *eventbus.Bus
	Logger *slog.Logger
	Runner ProcessRunner
}

// Module implements lifecycle.Module for tunnel management.
type Module struct {
	repo   *repo
	cfg    config.Config
	bus    *eventbus.Bus
	logger *slog.Logger
	runner ProcessRunner

	mu             sync.Mutex
	persisted      ConfigRecord
	state          string
	publicURL      string
	pid            int
	startedAt      string
	lastError      string
	activeHandle   ProcessHandle
	activeProv     string
	activeLocalURL string
	activeTSMode   string
	runWG          sync.WaitGroup
	stopRequested  bool
}

var _ lifecycle.Module = (*Module)(nil)

// New constructs a tunnel Module.
func New(opts Options) *Module {
	logger := opts.Logger
	if logger == nil {
		logger = slog.Default()
	}
	runner := opts.Runner
	if runner == nil {
		runner = osProcessRunner{}
	}
	return &Module{
		repo:   newRepo(opts.DB),
		cfg:    opts.Config,
		bus:    opts.Bus,
		logger: logger,
		runner: runner,
		state:  StateStopped,
	}
}

// Name returns the canonical module name.
func (m *Module) Name() string { return moduleName }

// Dependencies declares modules required before tunnel initializes.
func (m *Module) Dependencies() []string { return []string{"system"} }

// Migrations returns embedded SQLite migrations for tunnel.
func (m *Module) Migrations() fs.FS { return migrationsFS }

// RegisterRoutes mounts tunnel management routes on AdminAPI.
func (m *Module) RegisterRoutes(mounts lifecycle.RouteMounts) {
	if mounts.AdminAPI == nil {
		return
	}
	mounts.AdminAPI.Get("/tunnel/status", m.handleStatus)
	mounts.AdminAPI.Patch("/tunnel/config", m.handlePatchConfig)
	mounts.AdminAPI.Post("/tunnel/start", m.handleStart)
	mounts.AdminAPI.Post("/tunnel/stop", m.handleStop)
	mounts.AdminAPI.Get("/tunnel/events", m.handleEvents)
}

func (m *Module) defaultLocalURL() string {
	port := m.cfg.Port
	if port <= 0 {
		port = 20128
	}
	host := strings.TrimSpace(m.cfg.Host)
	if host == "" || host == "0.0.0.0" || host == "::" || host == "[::]" {
		host = "127.0.0.1"
	}
	return fmt.Sprintf("http://%s", net.JoinHostPort(host, strconv.Itoa(port)))
}

// Start loads persisted tunnel config and optionally auto-starts the tunnel.
func (m *Module) Start(ctx context.Context) error {
	cfg, err := m.repo.loadOrInitConfig(ctx, m.defaultLocalURL())
	if err != nil {
		return err
	}
	m.mu.Lock()
	m.persisted = cfg
	m.mu.Unlock()

	if cfg.AutoStart {
		if _, startErr := m.startTunnel(ctx, cfg.Provider, cfg.LocalURL, cfg.TailscaleMode); startErr != nil {
			m.logger.Warn("tunnel auto-start failed", "provider", cfg.Provider, "error", startErr)
		}
	}
	return nil
}

// Stop cleanly terminates any active tunnel process and waits for watcher
// goroutines to finish or until ctx is canceled.
func (m *Module) Stop(ctx context.Context) error {
	_, _ = m.stopTunnel(ctx)
	done := make(chan struct{})
	go func() {
		m.runWG.Wait()
		close(done)
	}()
	if ctx == nil {
		<-done
		return nil
	}
	select {
	case <-done:
		return nil
	case <-ctx.Done():
		return ctx.Err()
	}
}

func validateProvider(p string) error {
	if p != ProviderCloudflare && p != ProviderTailscale {
		return errors.New("provider must be 'cloudflare' or 'tailscale'")
	}
	return nil
}

func validateTailscaleMode(mode string) error {
	if mode != TailscaleModeFunnel && mode != TailscaleModeServe {
		return errors.New("tailscale_mode must be 'funnel' or 'serve'")
	}
	return nil
}

func validateLoopbackLocalURL(raw string) (normalized, port string, err error) {
	trimmed := strings.TrimRight(strings.TrimSpace(raw), "/")
	if trimmed == "" {
		return "", "", errors.New("local_url is required")
	}
	u, err := url.Parse(trimmed)
	if err != nil || (u.Scheme != "http" && u.Scheme != "https") {
		return "", "", errors.New("local_url must be an http:// or https:// URL")
	}
	host := strings.ToLower(strings.TrimSpace(u.Hostname()))
	isLoopback := host == "localhost"
	if !isLoopback {
		if ip := net.ParseIP(host); ip != nil && ip.IsLoopback() {
			isLoopback = true
		}
	}
	if !isLoopback {
		return "", "", errors.New("local_url host must be a loopback address (127.0.0.1, localhost, or ::1)")
	}
	p := u.Port()
	if p == "" {
		if u.Scheme == "https" {
			p = "443"
		} else {
			p = "80"
		}
	}
	portNum, convErr := strconv.Atoi(p)
	if convErr != nil || portNum < 1 || portNum > 65535 {
		return "", "", errors.New("local_url port must be between 1 and 65535")
	}
	return trimmed, p, nil
}

func (m *Module) detectBinaries() BinariesStatus {
	cfPath, cfErr := m.runner.LookPath("cloudflared")
	tsPath, tsErr := m.runner.LookPath("tailscale")
	return BinariesStatus{
		Cloudflared: BinaryStatus{Available: cfErr == nil && cfPath != "", Path: cfPath},
		Tailscale:   BinaryStatus{Available: tsErr == nil && tsPath != "", Path: tsPath},
	}
}

func (m *Module) snapshotLocked() Status {
	return Status{
		Status:        m.state,
		Provider:      m.persisted.Provider,
		PublicURL:     m.publicURL,
		LocalURL:      m.persisted.LocalURL,
		TailscaleMode: m.persisted.TailscaleMode,
		AutoStart:     m.persisted.AutoStart,
		PID:           m.pid,
		StartedAt:     m.startedAt,
		LastError:     m.lastError,
		Binaries:      m.detectBinaries(),
	}
}

func (m *Module) publishStatus(st Status) {
	if m.bus == nil {
		return
	}
	m.bus.Publish(TopicTunnelStatus, map[string]any{
		"status":     st.Status,
		"provider":   st.Provider,
		"public_url": st.PublicURL,
		"last_error": st.LastError,
	})
}

func (m *Module) handleStatus(w http.ResponseWriter, _ *http.Request) {
	m.mu.Lock()
	st := m.snapshotLocked()
	m.mu.Unlock()
	httpserver.WriteJSON(w, http.StatusOK, st)
}

type patchConfigReq struct {
	Provider      *string `json:"provider"`
	LocalURL      *string `json:"local_url"`
	TailscaleMode *string `json:"tailscale_mode"`
	AutoStart     *bool   `json:"auto_start"`
}

func (m *Module) handlePatchConfig(w http.ResponseWriter, r *http.Request) {
	var req patchConfigReq
	if err := decodeStrictJSON(r, maxBodyBytes, &req); err != nil {
		httpserver.WriteError(w, r, http.StatusBadRequest, "invalid_payload", err.Error())
		return
	}

	m.mu.Lock()
	next := m.persisted
	m.mu.Unlock()

	if req.Provider != nil {
		p := strings.TrimSpace(*req.Provider)
		if err := validateProvider(p); err != nil {
			httpserver.WriteError(w, r, http.StatusBadRequest, "invalid_provider", err.Error())
			return
		}
		next.Provider = p
	}
	if req.LocalURL != nil {
		norm, _, err := validateLoopbackLocalURL(*req.LocalURL)
		if err != nil {
			httpserver.WriteError(w, r, http.StatusBadRequest, "invalid_local_url", err.Error())
			return
		}
		next.LocalURL = norm
	}
	if req.TailscaleMode != nil {
		mode := strings.TrimSpace(*req.TailscaleMode)
		if err := validateTailscaleMode(mode); err != nil {
			httpserver.WriteError(w, r, http.StatusBadRequest, "invalid_tailscale_mode", err.Error())
			return
		}
		next.TailscaleMode = mode
	}
	if req.AutoStart != nil {
		next.AutoStart = *req.AutoStart
	}
	next.UpdatedAt = nowUTC()

	if err := m.repo.saveConfig(r.Context(), next); err != nil {
		httpserver.WriteError(w, r, http.StatusInternalServerError, "tunnel_config_save_failed", "failed to save tunnel configuration")
		return
	}

	m.mu.Lock()
	m.persisted = next
	st := m.snapshotLocked()
	m.mu.Unlock()

	_ = m.repo.appendEvent(r.Context(), next.Provider, "config_updated", "tunnel configuration updated", st.PublicURL)
	httpserver.WriteJSON(w, http.StatusOK, st)
}

type startTunnelReq struct {
	Provider      string `json:"provider"`
	LocalURL      string `json:"local_url"`
	TailscaleMode string `json:"tailscale_mode"`
}

var (
	errAlreadyRunning = errors.New("tunnel is already running")
	errBinaryNotFound = errors.New("tunnel CLI binary not found on PATH")
)

func (m *Module) handleStart(w http.ResponseWriter, r *http.Request) {
	var req startTunnelReq
	if r.Body != nil {
		raw, err := io.ReadAll(io.LimitReader(r.Body, maxBodyBytes))
		if err != nil {
			httpserver.WriteError(w, r, http.StatusBadRequest, "invalid_payload", "failed to read request body")
			return
		}
		if strings.TrimSpace(string(raw)) != "" {
			dec := json.NewDecoder(strings.NewReader(string(raw)))
			dec.DisallowUnknownFields()
			if err := dec.Decode(&req); err != nil {
				httpserver.WriteError(w, r, http.StatusBadRequest, "invalid_payload", err.Error())
				return
			}
		}
	}

	m.mu.Lock()
	prov := m.persisted.Provider
	localURL := m.persisted.LocalURL
	tsMode := m.persisted.TailscaleMode
	m.mu.Unlock()

	if strings.TrimSpace(req.Provider) != "" {
		prov = strings.TrimSpace(req.Provider)
	}
	if strings.TrimSpace(req.LocalURL) != "" {
		localURL = strings.TrimSpace(req.LocalURL)
	}
	if strings.TrimSpace(req.TailscaleMode) != "" {
		tsMode = strings.TrimSpace(req.TailscaleMode)
	}

	if err := validateProvider(prov); err != nil {
		httpserver.WriteError(w, r, http.StatusBadRequest, "invalid_provider", err.Error())
		return
	}
	if err := validateTailscaleMode(tsMode); err != nil {
		httpserver.WriteError(w, r, http.StatusBadRequest, "invalid_tailscale_mode", err.Error())
		return
	}
	normLocalURL, _, err := validateLoopbackLocalURL(localURL)
	if err != nil {
		httpserver.WriteError(w, r, http.StatusBadRequest, "invalid_local_url", err.Error())
		return
	}

	st, err := m.startTunnel(r.Context(), prov, normLocalURL, tsMode)
	if err != nil {
		switch {
		case errors.Is(err, errAlreadyRunning):
			httpserver.WriteError(w, r, http.StatusConflict, "tunnel_already_running", err.Error())
		case errors.Is(err, errBinaryNotFound):
			httpserver.WriteError(w, r, http.StatusUnprocessableEntity, "binary_not_found", err.Error())
		default:
			httpserver.WriteError(w, r, http.StatusInternalServerError, "tunnel_start_failed", err.Error())
		}
		return
	}
	httpserver.WriteJSON(w, http.StatusOK, st)
}

func (m *Module) startTunnel(ctx context.Context, provider, localURL, tsMode string) (Status, error) {
	normLocalURL, port, err := validateLoopbackLocalURL(localURL)
	if err != nil {
		return Status{}, err
	}

	m.mu.Lock()
	if m.state == StateRunning || m.state == StateStarting {
		m.mu.Unlock()
		return Status{}, errAlreadyRunning
	}

	if provider == ProviderCloudflare {
		if _, lpErr := m.runner.LookPath("cloudflared"); lpErr != nil {
			m.mu.Unlock()
			return Status{}, fmt.Errorf("%w: cloudflared is not installed", errBinaryNotFound)
		}
		handle, startErr := m.runner.Start(context.Background(), "cloudflared", "tunnel", "--url", normLocalURL, "--no-autoupdate")
		if startErr != nil {
			m.state = StateError
			m.lastError = startErr.Error()
			m.mu.Unlock()
			return Status{}, startErr
		}
		m.state = StateStarting
		m.publicURL = ""
		m.lastError = ""
		m.pid = handle.PID()
		m.startedAt = nowUTC()
		m.activeHandle = handle
		m.activeProv = provider
		m.activeLocalURL = normLocalURL
		m.stopRequested = false
		st := m.snapshotLocked()
		m.runWG.Add(1)
		m.mu.Unlock()

		_ = m.repo.appendEvent(ctx, provider, "started", "started cloudflared quick tunnel for "+normLocalURL, "")
		m.publishStatus(st)
		//nolint:gosec // G118: tunnel watcher intentionally outlives the HTTP start request and is stopped via Stop().
		go m.watchCloudflared(handle)
		return st, nil
	}

	// Tailscale serve / funnel: transition to StateStarting under lock BEFORE
	// executing the CLI so concurrent Start/Stop calls are serialized safely.
	if _, lpErr := m.runner.LookPath("tailscale"); lpErr != nil {
		m.mu.Unlock()
		return Status{}, fmt.Errorf("%w: tailscale is not installed", errBinaryNotFound)
	}
	m.state = StateStarting
	m.publicURL = ""
	m.lastError = ""
	m.pid = 0
	m.activeProv = provider
	m.activeLocalURL = normLocalURL
	m.activeTSMode = tsMode
	m.stopRequested = false
	m.mu.Unlock()

	cmdCtx, cancel := context.WithTimeout(ctx, tailscaleCmdTimeout)
	defer cancel()
	out, runErr := m.runner.Run(cmdCtx, "tailscale", tsMode, "--bg", port)
	if runErr != nil {
		msg := strings.TrimSpace(string(out))
		if msg == "" {
			msg = runErr.Error()
		}
		m.mu.Lock()
		m.state = StateError
		m.lastError = msg
		st := m.snapshotLocked()
		m.mu.Unlock()
		_ = m.repo.appendEvent(ctx, provider, "error", msg, "")
		m.publishStatus(st)
		return Status{}, errors.New(msg)
	}

	publicURL := reHTTPSURL.FindString(string(out))

	m.mu.Lock()
	if m.stopRequested {
		m.state = StateStopped
		m.publicURL = ""
		st := m.snapshotLocked()
		m.mu.Unlock()
		offCtx, offCancel := context.WithTimeout(context.Background(), tailscaleCmdTimeout)
		_, _ = m.runner.Run(offCtx, "tailscale", tsMode, "--bg", "off")
		offCancel()
		return st, nil
	}
	m.state = StateRunning
	m.publicURL = publicURL
	m.lastError = ""
	m.pid = 0
	m.startedAt = nowUTC()
	m.activeProv = provider
	m.activeLocalURL = normLocalURL
	m.activeTSMode = tsMode
	st := m.snapshotLocked()
	m.mu.Unlock()

	_ = m.repo.appendEvent(ctx, provider, "started", fmt.Sprintf("tailscale %s active on port %s", tsMode, port), publicURL)
	m.publishStatus(st)
	return st, nil
}

func (m *Module) watchCloudflared(handle ProcessHandle) {
	defer m.runWG.Done()
	if out := handle.Output(); out != nil {
		// Read with bufio.Reader instead of bufio.Scanner so lines >64 KiB
		// never abort draining and deadlock os/exec's pipe writer.
		reader := bufio.NewReader(out)
		for {
			line, err := reader.ReadString('\n')
			if found := reCloudflareURL.FindString(line); found != "" {
				m.mu.Lock()
				if m.activeHandle == handle && m.publicURL != found {
					m.publicURL = found
					m.state = StateRunning
					st := m.snapshotLocked()
					m.mu.Unlock()
					_ = m.repo.appendEvent(context.Background(), ProviderCloudflare, "url_ready", "public quick tunnel URL assigned", found)
					m.publishStatus(st)
				} else {
					m.mu.Unlock()
				}
			}
			if err != nil {
				break
			}
		}
		_, _ = io.Copy(io.Discard, reader)
	}

	waitErr := handle.Wait()
	m.mu.Lock()
	if m.activeHandle != handle {
		m.mu.Unlock()
		return
	}
	wasStopped := m.stopRequested
	m.activeHandle = nil
	m.pid = 0
	m.publicURL = ""
	if wasStopped {
		m.state = StateStopped
		m.lastError = ""
	} else {
		m.state = StateError
		if waitErr != nil {
			m.lastError = waitErr.Error()
		} else {
			m.lastError = "cloudflared process exited unexpectedly"
		}
	}
	st := m.snapshotLocked()
	m.mu.Unlock()

	if !wasStopped {
		_ = m.repo.appendEvent(context.Background(), ProviderCloudflare, "exited", st.LastError, "")
		m.publishStatus(st)
	}
}

func (m *Module) handleStop(w http.ResponseWriter, r *http.Request) {
	st, err := m.stopTunnel(r.Context())
	if err != nil {
		httpserver.WriteError(w, r, http.StatusInternalServerError, "tunnel_stop_failed", err.Error())
		return
	}
	httpserver.WriteJSON(w, http.StatusOK, st)
}

func (m *Module) stopTunnel(ctx context.Context) (Status, error) {
	m.mu.Lock()
	handle := m.activeHandle
	prov := m.activeProv
	tsMode := m.activeTSMode
	wasActive := m.state == StateRunning || m.state == StateStarting
	if prov == "" {
		prov = m.persisted.Provider
	}
	if tsMode == "" {
		tsMode = m.persisted.TailscaleMode
	}
	m.stopRequested = true
	m.mu.Unlock()

	if handle != nil {
		_ = handle.Kill()
	} else if prov == ProviderTailscale && wasActive {
		cmdCtx, cancel := context.WithTimeout(ctx, tailscaleCmdTimeout)
		out, runErr := m.runner.Run(cmdCtx, "tailscale", tsMode, "--bg", "off")
		cancel()
		if runErr != nil {
			msg := strings.TrimSpace(string(out))
			if msg == "" {
				msg = runErr.Error()
			}
			m.mu.Lock()
			m.state = StateError
			m.lastError = msg
			st := m.snapshotLocked()
			m.mu.Unlock()
			_ = m.repo.appendEvent(ctx, prov, "error", "tailscale stop failed: "+msg, st.PublicURL)
			m.publishStatus(st)
			return st, errors.New(msg)
		}
	}

	m.mu.Lock()
	m.activeHandle = nil
	m.activeProv = ""
	m.activeLocalURL = ""
	m.activeTSMode = ""
	m.state = StateStopped
	m.publicURL = ""
	m.pid = 0
	m.startedAt = ""
	m.lastError = ""
	st := m.snapshotLocked()
	m.mu.Unlock()

	_ = m.repo.appendEvent(ctx, prov, "stopped", "tunnel stopped", "")
	m.publishStatus(st)
	return st, nil
}

func (m *Module) handleEvents(w http.ResponseWriter, r *http.Request) {
	limit := defaultEventList
	if raw := strings.TrimSpace(r.URL.Query().Get("limit")); raw != "" {
		if n, err := strconv.Atoi(raw); err == nil && n > 0 {
			limit = n
		}
	}
	events, err := m.repo.listEvents(r.Context(), limit)
	if err != nil {
		httpserver.WriteError(w, r, http.StatusInternalServerError, "tunnel_events_failed", "failed to list tunnel events")
		return
	}
	httpserver.WriteJSON(w, http.StatusOK, map[string]any{"events": events})
}

func decodeStrictJSON(r *http.Request, limit int64, dst any) error {
	if r.Body == nil {
		return errors.New("request body is required")
	}
	dec := json.NewDecoder(io.LimitReader(r.Body, limit))
	dec.DisallowUnknownFields()
	if err := dec.Decode(dst); err != nil {
		return fmt.Errorf("invalid JSON body: %w", err)
	}
	return nil
}
