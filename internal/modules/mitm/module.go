package mitm

import (
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
	"runtime"
	"slices"
	"strconv"
	"strings"
	"sync"

	"github.com/go-chi/chi/v5"

	"github.com/omeg4monx/router-gateway/internal/platform/config"
	"github.com/omeg4monx/router-gateway/internal/platform/db"
	"github.com/omeg4monx/router-gateway/internal/platform/eventbus"
	"github.com/omeg4monx/router-gateway/internal/platform/httpserver"
	"github.com/omeg4monx/router-gateway/internal/platform/lifecycle"
)

const (
	moduleName       = "mitm"
	maxJSONBodyBytes = 1 << 20
	eventTopicStatus = "mitm.status"
)

const (
	StateStopped = "stopped"
	StateRunning = "running"
	StateError   = "error"
)

//go:embed migrations/*.sql
var migrationsFS embed.FS

// InstallGuide provides OS-specific Root CA trust and hosts-file instructions.
type InstallGuide struct {
	OS           string `json:"os"`
	CACommand    string `json:"ca_command"`
	HostsPath    string `json:"hosts_path"`
	HostsSnippet string `json:"hosts_snippet"`
}

// Status is returned by GET /api/mitm/status, PATCH /api/mitm/config, POST /api/mitm/start, and POST /api/mitm/stop.
type Status struct {
	Status           string           `json:"status"`
	ListenAddr       string           `json:"listen_addr"`
	BoundAddr        string           `json:"bound_addr"`
	TargetGatewayURL string           `json:"target_gateway_url"`
	EnabledTools     []string         `json:"enabled_tools"`
	AutoStart        bool             `json:"auto_start"`
	StartedAt        string           `json:"started_at"`
	LastError        string           `json:"last_error"`
	CA               CAInfo           `json:"ca"`
	Tools            []ToolDescriptor `json:"tools"`
	InstallGuide     InstallGuide     `json:"install_guide"`
}

// AliasListResponse is returned by GET /api/mitm/aliases.
type AliasListResponse struct {
	Aliases []Alias `json:"aliases"`
}

// EventListResponse is returned by GET /api/mitm/events.
type EventListResponse struct {
	Events []Event `json:"events"`
}

// Options configures a mitm Module instance.
type Options struct {
	DB     *db.DB
	Config config.Config
	Bus    *eventbus.Bus
	Logger *slog.Logger
}

// Module implements lifecycle.Module for local TLS MITM interception.
type Module struct {
	repo   *repository
	cfg    config.Config
	bus    *eventbus.Bus
	logger *slog.Logger
	certs  *certManager

	opMu     sync.Mutex
	configMu sync.Mutex

	mu          sync.RWMutex
	state       string
	boundAddr   string
	startedAt   string
	lastError   string
	srv         *http.Server
	listener    net.Listener
	proxyClient *http.Client
	serveDone   chan struct{}
}

var _ lifecycle.Module = (*Module)(nil)

// New constructs a mitm Module.
func New(opts Options) *Module {
	logger := opts.Logger
	if logger == nil {
		logger = slog.Default()
	}
	return &Module{
		repo:   newRepository(opts.DB),
		cfg:    opts.Config,
		bus:    opts.Bus,
		logger: logger,
		certs:  newCertManager(opts.Config.DataDir),
		state:  StateStopped,
	}
}

func (m *Module) Name() string           { return moduleName }
func (m *Module) Dependencies() []string { return []string{"system"} }

func (m *Module) Migrations() fs.FS {
	sub, err := fs.Sub(migrationsFS, "migrations")
	if err != nil {
		panic(fmt.Sprintf("mitm: sub migrations fs: %v", err))
	}
	return sub
}

func (m *Module) defaultGatewayURL() string {
	host := strings.Trim(strings.TrimSpace(m.cfg.Host), "[]")
	if !isLoopbackHost(host) {
		host = "127.0.0.1"
	}
	port := m.cfg.Port
	if port <= 0 {
		port = 20128
	}
	return fmt.Sprintf("http://%s", net.JoinHostPort(host, strconv.Itoa(port)))
}

func (m *Module) Start(ctx context.Context) error {
	if err := m.repo.ensureDefaults(ctx, m.defaultGatewayURL()); err != nil {
		return err
	}
	if err := m.certs.loadFromDisk(); err != nil {
		m.logger.Warn("mitm: failed to load existing root CA from disk", "error", err)
	}
	cfg, err := m.repo.getConfig(ctx)
	if err != nil {
		return err
	}
	if cfg.AutoStart {
		if _, startErr := m.startProxyServer(cfg); startErr != nil {
			m.mu.Lock()
			m.state = StateError
			m.lastError = startErr.Error()
			m.mu.Unlock()
			m.logger.Warn("mitm: auto-start failed", "error", startErr)
		}
	}
	return nil
}

func (m *Module) Stop(ctx context.Context) error {
	return m.stopProxyServer(ctx)
}

func (m *Module) RegisterRoutes(mounts lifecycle.RouteMounts) {
	if mounts.AdminAPI == nil {
		return
	}
	mounts.AdminAPI.Route("/mitm", func(r chi.Router) {
		r.Get("/status", m.handleGetStatus)
		r.Patch("/config", m.handlePatchConfig)
		r.Post("/start", m.handleStart)
		r.Post("/stop", m.handleStop)
		r.Post("/ca/generate", m.handleGenerateCA)
		r.Get("/ca/download", m.handleDownloadCA)
		r.Get("/aliases", m.handleListAliases)
		r.Post("/aliases", m.handleCreateAlias)
		r.Delete("/aliases/{id}", m.handleDeleteAlias)
		r.Get("/events", m.handleListEvents)
		r.Delete("/events", m.handleClearEvents)
	})
}

func (m *Module) buildStatus(ctx context.Context) (Status, error) {
	cfg, err := m.repo.getConfig(ctx)
	if err != nil {
		return Status{}, err
	}
	m.mu.RLock()
	state := m.state
	boundAddr := m.boundAddr
	startedAt := m.startedAt
	lastErr := m.lastError
	m.mu.RUnlock()

	caInfo := m.certs.info()
	tools := make([]ToolDescriptor, 0, len(toolCatalog))
	var hostLines []string
	for _, t := range toolCatalog {
		enabled := slices.Contains(cfg.EnabledTools, t.id)
		hostsCopy := make([]string, len(t.hosts))
		copy(hostsCopy, t.hosts)
		tools = append(tools, ToolDescriptor{
			ID:      t.id,
			Name:    t.name,
			Hosts:   hostsCopy,
			Enabled: enabled,
		})
		if enabled {
			for _, h := range t.hosts {
				hostLines = append(hostLines, "127.0.0.1 "+h)
			}
		}
	}
	if len(hostLines) == 0 {
		hostLines = append(hostLines, "# Enable one or more IDE tools above to generate hosts entries")
	}

	return Status{
		Status:           state,
		ListenAddr:       cfg.ListenAddr,
		BoundAddr:        boundAddr,
		TargetGatewayURL: cfg.TargetGatewayURL,
		EnabledTools:     cfg.EnabledTools,
		AutoStart:        cfg.AutoStart,
		StartedAt:        startedAt,
		LastError:        lastErr,
		CA:               caInfo,
		Tools:            tools,
		InstallGuide:     buildInstallGuide(caInfo.CertPath, strings.Join(hostLines, "\n")),
	}, nil
}

func buildInstallGuide(certPath, hostsSnippet string) InstallGuide {
	switch runtime.GOOS {
	case "darwin":
		return InstallGuide{
			OS:           "darwin",
			CACommand:    fmt.Sprintf("sudo security add-trusted-cert -d -r trustRoot -k /Library/Keychains/System.keychain %q", certPath),
			HostsPath:    "/etc/hosts",
			HostsSnippet: hostsSnippet,
		}
	case "windows":
		return InstallGuide{
			OS:           "windows",
			CACommand:    fmt.Sprintf("certutil -addstore -f Root %q", certPath),
			HostsPath:    `C:\Windows\System32\drivers\etc\hosts`,
			HostsSnippet: hostsSnippet,
		}
	default:
		return InstallGuide{
			OS:           "linux",
			CACommand:    fmt.Sprintf("sudo cp %q /usr/local/share/ca-certificates/omeg4mon-mitm.crt && sudo update-ca-certificates", certPath),
			HostsPath:    "/etc/hosts",
			HostsSnippet: hostsSnippet,
		}
	}
}

func (m *Module) handleGetStatus(w http.ResponseWriter, r *http.Request) {
	st, err := m.buildStatus(r.Context())
	if err != nil {
		httpserver.WriteError(w, r, http.StatusInternalServerError, "mitm_status_error", err.Error())
		return
	}
	httpserver.WriteJSON(w, http.StatusOK, st)
}

type updateConfigRequest struct {
	ListenAddr       *string   `json:"listen_addr"`
	TargetGatewayURL *string   `json:"target_gateway_url"`
	EnabledTools     *[]string `json:"enabled_tools"`
	AutoStart        *bool     `json:"auto_start"`
}

func (m *Module) handlePatchConfig(w http.ResponseWriter, r *http.Request) {
	var req updateConfigRequest
	if err := decodeStrictJSON(r, &req); err != nil {
		httpserver.WriteError(w, r, http.StatusBadRequest, "invalid_request", err.Error())
		return
	}
	m.configMu.Lock()
	cfg, err := m.repo.getConfig(r.Context())
	if err != nil {
		m.configMu.Unlock()
		httpserver.WriteError(w, r, http.StatusInternalServerError, "mitm_config_error", err.Error())
		return
	}
	if req.ListenAddr != nil {
		if err := validateLoopbackListenAddr(*req.ListenAddr); err != nil {
			m.configMu.Unlock()
			httpserver.WriteError(w, r, http.StatusBadRequest, "invalid_listen_addr", err.Error())
			return
		}
		cfg.ListenAddr = strings.TrimSpace(*req.ListenAddr)
	}
	if req.TargetGatewayURL != nil {
		if err := validateLoopbackURL(*req.TargetGatewayURL); err != nil {
			m.configMu.Unlock()
			httpserver.WriteError(w, r, http.StatusBadRequest, "invalid_target_url", err.Error())
			return
		}
		cfg.TargetGatewayURL = strings.TrimSpace(*req.TargetGatewayURL)
	}
	if req.EnabledTools != nil {
		deduped := make([]string, 0, len(*req.EnabledTools))
		for _, t := range *req.EnabledTools {
			clean := strings.ToLower(strings.TrimSpace(t))
			if !validToolID(clean) {
				m.configMu.Unlock()
				httpserver.WriteError(w, r, http.StatusBadRequest, "invalid_tool", fmt.Sprintf("unknown MITM tool %q", t))
				return
			}
			if !slices.Contains(deduped, clean) {
				deduped = append(deduped, clean)
			}
		}
		cfg.EnabledTools = deduped
	}
	if req.AutoStart != nil {
		cfg.AutoStart = *req.AutoStart
	}
	if _, err := m.repo.updateConfig(r.Context(), cfg); err != nil {
		m.configMu.Unlock()
		httpserver.WriteError(w, r, http.StatusInternalServerError, "mitm_config_error", err.Error())
		return
	}
	m.configMu.Unlock()
	st, err := m.buildStatus(r.Context())
	if err != nil {
		httpserver.WriteError(w, r, http.StatusInternalServerError, "mitm_status_error", err.Error())
		return
	}
	httpserver.WriteJSON(w, http.StatusOK, st)
}

func (m *Module) handleStart(w http.ResponseWriter, r *http.Request) {
	cfg, err := m.repo.getConfig(r.Context())
	if err != nil {
		httpserver.WriteError(w, r, http.StatusInternalServerError, "mitm_config_error", err.Error())
		return
	}
	if err := validateLoopbackListenAddr(cfg.ListenAddr); err != nil {
		httpserver.WriteError(w, r, http.StatusBadRequest, "invalid_listen_addr", err.Error())
		return
	}
	if err := validateLoopbackURL(cfg.TargetGatewayURL); err != nil {
		httpserver.WriteError(w, r, http.StatusBadRequest, "invalid_target_url", err.Error())
		return
	}
	if _, err := m.startProxyServer(cfg); err != nil {
		if errors.Is(err, ErrAlreadyRunning) {
			httpserver.WriteError(w, r, http.StatusConflict, "mitm_already_running", "MITM proxy is already running")
			return
		}
		m.mu.Lock()
		if m.state != StateRunning {
			m.state = StateError
			m.lastError = err.Error()
		}
		m.mu.Unlock()
		httpserver.WriteError(w, r, http.StatusInternalServerError, "mitm_start_failed", err.Error())
		return
	}
	if m.bus != nil {
		m.bus.Publish(eventTopicStatus, map[string]any{"status": StateRunning})
	}
	st, err := m.buildStatus(r.Context())
	if err != nil {
		httpserver.WriteError(w, r, http.StatusInternalServerError, "mitm_status_error", err.Error())
		return
	}
	httpserver.WriteJSON(w, http.StatusOK, st)
}

func (m *Module) handleStop(w http.ResponseWriter, r *http.Request) {
	if err := m.stopProxyServer(r.Context()); err != nil {
		httpserver.WriteError(w, r, http.StatusInternalServerError, "mitm_stop_failed", err.Error())
		return
	}
	if m.bus != nil {
		m.bus.Publish(eventTopicStatus, map[string]any{"status": StateStopped})
	}
	st, err := m.buildStatus(r.Context())
	if err != nil {
		httpserver.WriteError(w, r, http.StatusInternalServerError, "mitm_status_error", err.Error())
		return
	}
	httpserver.WriteJSON(w, http.StatusOK, st)
}

type generateCARequest struct {
	Force bool `json:"force"`
}

func (m *Module) handleGenerateCA(w http.ResponseWriter, r *http.Request) {
	var req generateCARequest
	if r.Body != nil {
		raw, err := io.ReadAll(io.LimitReader(r.Body, maxJSONBodyBytes))
		_ = r.Body.Close()
		if err != nil {
			httpserver.WriteError(w, r, http.StatusBadRequest, "invalid_request", err.Error())
			return
		}
		if strings.TrimSpace(string(raw)) != "" {
			dec := json.NewDecoder(strings.NewReader(string(raw)))
			dec.DisallowUnknownFields()
			if err := dec.Decode(&req); err != nil {
				httpserver.WriteError(w, r, http.StatusBadRequest, "invalid_request", err.Error())
				return
			}
			if err := dec.Decode(&struct{}{}); !errors.Is(err, io.EOF) {
				httpserver.WriteError(w, r, http.StatusBadRequest, "invalid_request", "request body must contain a single JSON object")
				return
			}
		}
	}
	info, err := m.certs.ensureCA(req.Force)
	if err != nil {
		httpserver.WriteError(w, r, http.StatusInternalServerError, "mitm_ca_error", err.Error())
		return
	}
	httpserver.WriteJSON(w, http.StatusOK, info)
}

func (m *Module) handleDownloadCA(w http.ResponseWriter, r *http.Request) {
	pemBytes, err := m.certs.downloadPEM()
	if err != nil {
		if errors.Is(err, ErrCANotFound) {
			httpserver.WriteError(w, r, http.StatusNotFound, "ca_not_found", "Root CA has not been generated yet")
			return
		}
		httpserver.WriteError(w, r, http.StatusInternalServerError, "mitm_ca_error", err.Error())
		return
	}
	w.Header().Set("Content-Type", "application/x-pem-file")
	w.Header().Set("Content-Disposition", `attachment; filename="omeg4mon-mitm-rootCA.crt"`)
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write(pemBytes)
}

func (m *Module) handleListAliases(w http.ResponseWriter, r *http.Request) {
	aliases, err := m.repo.listAliases(r.Context())
	if err != nil {
		httpserver.WriteError(w, r, http.StatusInternalServerError, "mitm_alias_error", err.Error())
		return
	}
	httpserver.WriteJSON(w, http.StatusOK, AliasListResponse{Aliases: aliases})
}

type createAliasRequest struct {
	Tool        string `json:"tool"`
	SourceModel string `json:"source_model"`
	TargetModel string `json:"target_model"`
	Enabled     *bool  `json:"enabled"`
}

func (m *Module) handleCreateAlias(w http.ResponseWriter, r *http.Request) {
	var req createAliasRequest
	if err := decodeStrictJSON(r, &req); err != nil {
		httpserver.WriteError(w, r, http.StatusBadRequest, "invalid_request", err.Error())
		return
	}
	tool := strings.ToLower(strings.TrimSpace(req.Tool))
	if tool == "" {
		tool = "*"
	}
	if tool != "*" && !validToolID(tool) {
		httpserver.WriteError(w, r, http.StatusBadRequest, "invalid_tool", fmt.Sprintf("unknown tool %q", req.Tool))
		return
	}
	src := strings.TrimSpace(req.SourceModel)
	tgt := strings.TrimSpace(req.TargetModel)
	if src == "" || tgt == "" {
		httpserver.WriteError(w, r, http.StatusBadRequest, "invalid_alias", "source_model and target_model are required")
		return
	}
	enabled := true
	if req.Enabled != nil {
		enabled = *req.Enabled
	}
	created, err := m.repo.createAlias(r.Context(), Alias{
		ID:          newID("mal"),
		Tool:        tool,
		SourceModel: src,
		TargetModel: tgt,
		Enabled:     enabled,
	})
	if err != nil {
		if errors.Is(err, ErrAliasConflict) {
			httpserver.WriteError(w, r, http.StatusConflict, "alias_conflict", err.Error())
			return
		}
		httpserver.WriteError(w, r, http.StatusInternalServerError, "mitm_alias_error", err.Error())
		return
	}
	httpserver.WriteJSON(w, http.StatusCreated, created)
}

func (m *Module) handleDeleteAlias(w http.ResponseWriter, r *http.Request) {
	id := strings.TrimSpace(chi.URLParam(r, "id"))
	if id == "" {
		httpserver.WriteError(w, r, http.StatusBadRequest, "invalid_id", "alias id is required")
		return
	}
	if err := m.repo.deleteAlias(r.Context(), id); err != nil {
		if errors.Is(err, ErrAliasNotFound) {
			httpserver.WriteError(w, r, http.StatusNotFound, "alias_not_found", err.Error())
			return
		}
		httpserver.WriteError(w, r, http.StatusInternalServerError, "mitm_alias_error", err.Error())
		return
	}
	httpserver.WriteJSON(w, http.StatusOK, map[string]any{"deleted": true, "id": id})
}

func (m *Module) handleListEvents(w http.ResponseWriter, r *http.Request) {
	limit := 50
	if raw := strings.TrimSpace(r.URL.Query().Get("limit")); raw != "" {
		if parsed, err := strconv.Atoi(raw); err == nil && parsed > 0 {
			limit = parsed
		}
	}
	events, err := m.repo.listEvents(r.Context(), limit)
	if err != nil {
		httpserver.WriteError(w, r, http.StatusInternalServerError, "mitm_events_error", err.Error())
		return
	}
	httpserver.WriteJSON(w, http.StatusOK, EventListResponse{Events: events})
}

func (m *Module) handleClearEvents(w http.ResponseWriter, r *http.Request) {
	deleted, err := m.repo.clearEvents(r.Context())
	if err != nil {
		httpserver.WriteError(w, r, http.StatusInternalServerError, "mitm_events_error", err.Error())
		return
	}
	httpserver.WriteJSON(w, http.StatusOK, map[string]any{"deleted": deleted})
}

func decodeStrictJSON(r *http.Request, dst any) error {
	defer func() { _ = r.Body.Close() }()
	dec := json.NewDecoder(io.LimitReader(r.Body, maxJSONBodyBytes))
	dec.DisallowUnknownFields()
	if err := dec.Decode(dst); err != nil {
		return err
	}
	if err := dec.Decode(&struct{}{}); !errors.Is(err, io.EOF) {
		return errors.New("request body must contain a single JSON object")
	}
	return nil
}
