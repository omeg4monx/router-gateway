package main

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"log/slog"
	"net/http"
	"net/url"
	"os"
	"os/signal"
	"path/filepath"
	"slices"
	"sync"
	"sync/atomic"
	"syscall"
	"time"

	"github.com/omeg4monx/router-gateway/internal/contracts"
	"github.com/omeg4monx/router-gateway/internal/modules/auth"
	"github.com/omeg4monx/router-gateway/internal/modules/clitools"
	"github.com/omeg4monx/router-gateway/internal/modules/combos"
	"github.com/omeg4monx/router-gateway/internal/modules/gateway"
	"github.com/omeg4monx/router-gateway/internal/modules/keys"
	"github.com/omeg4monx/router-gateway/internal/modules/modeleditor"
	"github.com/omeg4monx/router-gateway/internal/modules/notifications"
	"github.com/omeg4monx/router-gateway/internal/modules/promptinjection"
	"github.com/omeg4monx/router-gateway/internal/modules/providers"
	"github.com/omeg4monx/router-gateway/internal/modules/proxypools"
	"github.com/omeg4monx/router-gateway/internal/modules/security"
	"github.com/omeg4monx/router-gateway/internal/modules/settings"
	"github.com/omeg4monx/router-gateway/internal/modules/system"
	"github.com/omeg4monx/router-gateway/internal/modules/tokensaver"
	"github.com/omeg4monx/router-gateway/internal/modules/translator"
	"github.com/omeg4monx/router-gateway/internal/modules/usage"
	"github.com/omeg4monx/router-gateway/internal/platform/config"
	"github.com/omeg4monx/router-gateway/internal/platform/crypto"
	"github.com/omeg4monx/router-gateway/internal/platform/db"
	"github.com/omeg4monx/router-gateway/internal/platform/eventbus"
	"github.com/omeg4monx/router-gateway/internal/platform/httpclient"
	"github.com/omeg4monx/router-gateway/internal/platform/httpserver"
	"github.com/omeg4monx/router-gateway/internal/platform/jobs"
	"github.com/omeg4monx/router-gateway/internal/platform/lifecycle"
	"github.com/omeg4monx/router-gateway/internal/platform/logging"
)

const (
	defaultLogRingCap      = 500
	defaultShutdownTimeout = 10 * time.Second
)

var (
	version   = "dev"
	commit    = "none"
	buildTime = "unknown"
)

type dynamicProxyResolver struct {
	target atomic.Pointer[contracts.ProxyDialerResolver]
}

func (d *dynamicProxyResolver) Set(r contracts.ProxyDialerResolver) {
	if r == nil {
		d.target.Store(nil)
		return
	}
	d.target.Store(&r)
}

func (d *dynamicProxyResolver) ProxyURLForProvider(ctx context.Context, providerID string) (*url.URL, error) {
	p := d.target.Load()
	if p == nil || *p == nil {
		return nil, nil
	}
	return (*p).ProxyURLForProvider(ctx, providerID)
}

type jobInspectorAdapter struct {
	runner *jobs.Runner
}

func (a *jobInspectorAdapter) Snapshot() []contracts.JobStatus {
	if a == nil || a.runner == nil {
		return nil
	}
	raw := a.runner.Snapshot()
	out := make([]contracts.JobStatus, 0, len(raw))
	for _, st := range raw {
		out = append(out, contracts.JobStatus{
			Name:           st.Name,
			IntervalMs:     st.IntervalMs,
			Running:        st.Running,
			RunCount:       st.RunCount,
			ErrorCount:     st.ErrorCount,
			LastStatus:     st.LastStatus,
			LastError:      st.LastError,
			LastDurationMs: st.LastDurationMs,
			LastStartedAt:  st.LastStartedAt,
			LastFinishedAt: st.LastFinishedAt,
			NextRunAt:      st.NextRunAt,
		})
	}
	return out
}

func (a *jobInspectorAdapter) Trigger(ctx context.Context, name string) error {
	if a == nil || a.runner == nil {
		return contracts.ErrJobNotFound
	}
	err := a.runner.Trigger(ctx, name)
	switch {
	case errors.Is(err, jobs.ErrJobNotFound):
		return contracts.ErrJobNotFound
	case errors.Is(err, jobs.ErrJobAlreadyRunning):
		return contracts.ErrJobAlreadyRunning
	default:
		return err
	}
}

// App holds the wired platform primitives, module lifecycle registry, and HTTP
// server for a running omeg4mon-router instance.
type App struct {
	cfg        config.Config
	logger     *slog.Logger
	logRing    *logging.RingBuffer
	bus        *eventbus.Bus
	vault      *crypto.Vault
	sessions   *crypto.SessionSigner
	httpClient *http.Client
	db         *db.DB
	jobs       *jobs.Runner
	registry   *lifecycle.Registry
	server     *httpserver.Server

	runMu             sync.Mutex
	runCancel         context.CancelFunc
	shutdownRequested atomic.Bool

	closeOnce sync.Once
	closeErr  error
}

// buildApp wires all platform primitives, opens SQLite, runs module migrations,
// mounts HTTP routes, and starts active modules and background jobs. On any
// error partway through initialization, already-opened resources are cleaned up.
func buildApp(ctx context.Context, cfg config.Config, build system.BuildInfo) (*App, error) {
	if ctx == nil {
		ctx = context.Background()
	}

	logger, logRing := logging.New(cfg.LogLevel, os.Stderr, defaultLogRingCap)
	bus := eventbus.New()

	masterKey, err := crypto.LoadOrCreateMasterKey(cfg.DataDir, cfg.MasterKeyHex)
	if err != nil {
		return nil, fmt.Errorf("load master key: %w", err)
	}
	defer clear(masterKey)

	vault, err := crypto.NewVault(masterKey)
	if err != nil {
		return nil, fmt.Errorf("initialize vault: %w", err)
	}

	sessions, err := crypto.NewSessionSigner(masterKey)
	if err != nil {
		return nil, fmt.Errorf("initialize session signer: %w", err)
	}

	proxyResolver := &dynamicProxyResolver{}
	httpClient := httpclient.NewClient(httpclient.Options{
		ProxyFunc: httpclient.NewProviderProxyFunc(proxyResolver),
	})

	var (
		database       *db.DB
		registry       *lifecycle.Registry
		jobRunner      *jobs.Runner
		modulesStarted bool
		jobsStarted    bool
		ok             bool
	)
	defer func() {
		if ok {
			return
		}
		cleanupCtx, cancel := context.WithTimeout(context.Background(), defaultShutdownTimeout)
		defer cancel()
		if jobsStarted && jobRunner != nil {
			_ = jobRunner.Stop(cleanupCtx)
		}
		if modulesStarted && registry != nil {
			_ = registry.StopAll(cleanupCtx)
		}
		httpClient.CloseIdleConnections()
		if database != nil {
			_ = database.Close()
		}
	}()

	dbPath := filepath.Join(cfg.DataDir, "data.sqlite")
	database, err = db.Open(dbPath)
	if err != nil {
		return nil, fmt.Errorf("open database: %w", err)
	}

	jobRunner = jobs.NewRunner(logger)
	registry = lifecycle.NewRegistry(cfg.DisabledModules)

	uiFS, devUIHandler, err := loadUI()
	if err != nil {
		return nil, fmt.Errorf("load UI assets: %w", err)
	}

	app := &App{
		cfg:        cfg,
		logger:     logger,
		logRing:    logRing,
		bus:        bus,
		vault:      vault,
		sessions:   sessions,
		httpClient: httpClient,
		db:         database,
		jobs:       jobRunner,
		registry:   registry,
	}

	sysMod := system.New(system.Options{
		DB:            database,
		Config:        cfg,
		Build:         build,
		ActiveModules: registry.ActiveNames,
		Shutdown:      app.requestShutdown,
	})
	authMod := auth.New(auth.Options{
		Store:    sysMod,
		Sessions: sessions,
		Config:   cfg,
		Bus:      bus,
	})

	srv := httpserver.New(httpserver.Options{
		Config:       cfg,
		Logger:       logger,
		Auth:         authMod,
		UI:           uiFS,
		DevUIHandler: devUIHandler,
	})
	app.server = srv

	settingsMod := settings.New(settings.Options{
		DB:  database,
		Bus: bus,
	})
	keysMod := keys.New(keys.Options{
		DB:  database,
		Bus: bus,
	})
	providersMod := providers.New(providers.Options{
		DB:         database,
		Vault:      vault,
		HTTPClient: httpClient,
		Bus:        bus,
	})
	combosMod := combos.New(combos.Options{
		DB:        database,
		Providers: providersMod,
		Bus:       bus,
	})
	usageMod := usage.New(usage.Options{
		DB:     database,
		Bus:    bus,
		Logger: logger,
	})
	promptGuardMod := promptinjection.New(promptinjection.Options{
		DB:     database,
		Bus:    bus,
		Logger: logger,
	})
	modelEditorMod := modeleditor.New(modeleditor.Options{
		DB:     database,
		Bus:    bus,
		Logger: logger,
	})
	tokenSaverMod := tokensaver.New(tokensaver.Options{
		DB:     database,
		Bus:    bus,
		Logger: logger,
	})
	interceptors := make([]contracts.GatewayInterceptor, 0, 3)
	if !cfg.DisabledModules[promptGuardMod.Name()] {
		interceptors = append(interceptors, promptGuardMod)
	}
	if !cfg.DisabledModules[modelEditorMod.Name()] {
		interceptors = append(interceptors, modelEditorMod)
	}
	if !cfg.DisabledModules[tokenSaverMod.Name()] {
		interceptors = append(interceptors, tokenSaverMod)
	}
	gatewayMod := gateway.New(gateway.Options{
		DB:           database,
		HTTPClient:   httpClient,
		Keys:         keysMod,
		Resolver:     combosMod,
		Models:       combosMod,
		Providers:    providersMod,
		Health:       combosMod,
		Interceptors: interceptors,
		Usage:        usageMod,
		Logger:       logger,
	})
	clitoolsMod := clitools.New(clitools.Options{
		Config: cfg,
	})
	proxypoolsMod := proxypools.New(proxypools.Options{
		DB:     database,
		Vault:  vault,
		Bus:    bus,
		Logger: logger,
	})
	translatorMod := translator.New(translator.Options{
		Translator: gatewayMod,
		LogRing:    logRing,
	})
	securityMod := security.New(security.Options{
		DB:            database,
		Bus:           bus,
		Logger:        logger,
		Config:        cfg,
		Store:         sysMod,
		ActiveModules: registry.ActiveNames,
	})
	notificationsMod := notifications.New(notifications.Options{
		DB:     database,
		Bus:    bus,
		Logger: logger,
		Jobs:   &jobInspectorAdapter{runner: jobRunner},
	})

	if err := registry.Register(sysMod); err != nil {
		return nil, fmt.Errorf("register system module: %w", err)
	}
	if err := registry.Register(authMod); err != nil {
		return nil, fmt.Errorf("register auth module: %w", err)
	}
	if err := registry.Register(settingsMod); err != nil {
		return nil, fmt.Errorf("register settings module: %w", err)
	}
	if err := registry.Register(keysMod); err != nil {
		return nil, fmt.Errorf("register keys module: %w", err)
	}
	if err := registry.Register(providersMod); err != nil {
		return nil, fmt.Errorf("register providers module: %w", err)
	}
	if err := registry.Register(combosMod); err != nil {
		return nil, fmt.Errorf("register combos module: %w", err)
	}
	if err := registry.Register(usageMod); err != nil {
		return nil, fmt.Errorf("register usage module: %w", err)
	}
	if err := registry.Register(gatewayMod); err != nil {
		return nil, fmt.Errorf("register gateway module: %w", err)
	}
	if err := registry.Register(clitoolsMod); err != nil {
		return nil, fmt.Errorf("register clitools module: %w", err)
	}
	if err := registry.Register(proxypoolsMod); err != nil {
		return nil, fmt.Errorf("register proxypools module: %w", err)
	}
	if err := registry.Register(translatorMod); err != nil {
		return nil, fmt.Errorf("register translator module: %w", err)
	}
	if err := registry.Register(promptGuardMod); err != nil {
		return nil, fmt.Errorf("register promptinjection module: %w", err)
	}
	if err := registry.Register(modelEditorMod); err != nil {
		return nil, fmt.Errorf("register modeleditor module: %w", err)
	}
	if err := registry.Register(tokenSaverMod); err != nil {
		return nil, fmt.Errorf("register tokensaver module: %w", err)
	}
	if err := registry.Register(securityMod); err != nil {
		return nil, fmt.Errorf("register security module: %w", err)
	}
	if err := registry.Register(notificationsMod); err != nil {
		return nil, fmt.Errorf("register notifications module: %w", err)
	}

	if err := registry.MigrateAll(ctx, database); err != nil {
		return nil, fmt.Errorf("run module migrations: %w", err)
	}

	if err := registry.MountAll(srv.Mounts()); err != nil {
		return nil, fmt.Errorf("mount module routes: %w", err)
	}

	if err := registry.StartAll(ctx); err != nil {
		return nil, fmt.Errorf("start modules: %w", err)
	}
	modulesStarted = true
	activeNames := registry.ActiveNames()
	if slices.Contains(activeNames, proxypoolsMod.Name()) {
		proxyResolver.Set(proxypoolsMod)
	}

	jobRunner.SetObserver(func(st jobs.JobStatus) {
		if st.LastStatus == "error" || st.LastStatus == "panic" {
			bus.Publish(notifications.TopicJobFailed, notifications.EmitPayload{
				Title:    fmt.Sprintf("Background job %s failed", st.Name),
				Message:  st.LastError,
				Severity: "error",
				Category: "jobs",
				Source:   "jobs.runner",
				Metadata: map[string]any{
					"job":         st.Name,
					"last_status": st.LastStatus,
					"error_count": st.ErrorCount,
				},
			})
		}
	})
	jobRunner.Every("usage-prune", time.Hour, usageMod.PruneExpired)
	if slices.Contains(activeNames, securityMod.Name()) {
		jobRunner.Every("security-audit-prune", 6*time.Hour, securityMod.PruneExpired)
	}
	if slices.Contains(activeNames, notificationsMod.Name()) {
		jobRunner.Every("notifications-prune", 6*time.Hour, notificationsMod.PruneExpired)
	}
	jobRunner.Start(context.Background())
	jobsStarted = true

	ok = true
	return app, nil
}

// Handler returns the root HTTP handler for testing or embedding.
func (a *App) Handler() http.Handler {
	return a.server.Handler()
}

// Addr returns the active TCP listener address once Run binds a socket, or the
// configured host:port prior to startup.
func (a *App) Addr() string {
	return a.server.Addr()
}

// Run starts the HTTP server and blocks until ctx is canceled, the shutdown
// endpoint triggers shutdown, or the listener fails, then performs graceful
// teardown via Close.
func (a *App) Run(ctx context.Context) error {
	if ctx == nil {
		ctx = context.Background()
	}

	runCtx, cancel := context.WithCancel(ctx)
	a.runMu.Lock()
	a.runCancel = cancel
	alreadyRequested := a.shutdownRequested.Load()
	a.runMu.Unlock()

	defer func() {
		a.runMu.Lock()
		a.runCancel = nil
		a.runMu.Unlock()
		cancel()
	}()

	var serveErr error
	if !alreadyRequested && ctx.Err() == nil {
		serveErr = a.server.Start(runCtx)
		if errors.Is(serveErr, context.Canceled) && (ctx.Err() != nil || a.shutdownRequested.Load()) {
			serveErr = nil
		}
	}

	closeCtx, closeCancel := context.WithTimeout(context.WithoutCancel(ctx), defaultShutdownTimeout)
	defer closeCancel()

	closeErr := a.Close(closeCtx)
	return errors.Join(serveErr, closeErr)
}

// Close gracefully stops the HTTP server, background jobs, active modules,
// idle HTTP connections, and SQLite connection pools in strict LIFO order.
// Close is idempotent and safe to call multiple times.
func (a *App) Close(ctx context.Context) error {
	if a == nil {
		return nil
	}
	if ctx == nil {
		ctx = context.Background()
	}

	a.closeOnce.Do(func() {
		a.requestShutdown()

		var errs []error
		if a.server != nil {
			if err := a.server.Shutdown(ctx); err != nil {
				errs = append(errs, fmt.Errorf("shutdown http server: %w", err))
			}
		}
		if a.jobs != nil {
			if err := a.jobs.Stop(ctx); err != nil {
				errs = append(errs, fmt.Errorf("stop jobs runner: %w", err))
			}
		}
		if a.registry != nil {
			if err := a.registry.StopAll(ctx); err != nil {
				errs = append(errs, fmt.Errorf("stop modules: %w", err))
			}
		}
		if a.httpClient != nil {
			a.httpClient.CloseIdleConnections()
		}
		if a.db != nil {
			if err := a.db.Close(); err != nil {
				errs = append(errs, fmt.Errorf("close database: %w", err))
			}
		}
		a.closeErr = errors.Join(errs...)
	})

	return a.closeErr
}

func (a *App) requestShutdown() {
	a.shutdownRequested.Store(true)
	a.runMu.Lock()
	cancel := a.runCancel
	a.runMu.Unlock()
	if cancel != nil {
		cancel()
	}
}

func main() {
	os.Exit(runMain(os.Args[1:], os.Getenv))
}

func runMain(args []string, getenv func(string) string) int {
	cfg, err := config.Load(args, getenv)
	if err != nil {
		if errors.Is(err, flag.ErrHelp) {
			fmt.Fprintln(os.Stderr, "Usage: omeg4mon-router [--host 127.0.0.1] [--port 20128] [--data-dir PATH] [--log-level info]")
			return 0
		}
		fmt.Fprintf(os.Stderr, "omeg4mon-router: configuration error: %v\n", err)
		return 1
	}

	sigCtx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	build := system.BuildInfo{
		Version:   version,
		Commit:    commit,
		BuildTime: buildTime,
	}

	app, err := buildApp(sigCtx, cfg, build)
	if err != nil {
		fmt.Fprintf(os.Stderr, "omeg4mon-router: startup error: %v\n", err)
		return 1
	}

	app.logger.Info("omeg4mon-router starting",
		"version", build.Version,
		"commit", build.Commit,
		"addr", cfg.Addr(),
		"data_dir", cfg.DataDir,
	)

	if err := app.Run(sigCtx); err != nil {
		app.logger.Error("omeg4mon-router exited with error", "error", err)
		return 1
	}

	app.logger.Info("omeg4mon-router stopped")
	return 0
}
