package httpserver

import (
	"context"
	"errors"
	"io/fs"
	"log/slog"
	"net"
	"net/http"
	"runtime/debug"
	"sync"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"

	"github.com/omeg4monx/router-gateway/internal/contracts"
	"github.com/omeg4monx/router-gateway/internal/platform/config"
	"github.com/omeg4monx/router-gateway/internal/platform/lifecycle"
)

const (
	defaultReadHeaderTimeout = 10 * time.Second
	defaultIdleTimeout       = 120 * time.Second
	cancelShutdownGrace      = 5 * time.Second
)

var gatewayPrefixes = []string{
	"/v1",
	"/v1beta",
	"/codex",
	"/responses",
	"/models",
}

type adminSubjectKey struct{}

// AdminSubject returns the authenticated admin subject attached to ctx by the
// AdminAPI authentication middleware, or an empty string if unauthenticated.
func AdminSubject(ctx context.Context) string {
	if ctx == nil {
		return ""
	}
	v, _ := ctx.Value(adminSubjectKey{}).(string)
	return v
}

// Options configures a new HTTP Server instance.
type Options struct {
	Config       config.Config
	Logger       *slog.Logger
	Auth         contracts.AuthVerifier
	UI           fs.FS
	DevUIHandler http.Handler
}

// Server wraps net/http and chi with isolated middleware stacks for control
// plane APIs (/api/*), streaming gateway routes (/v1/*, /v1beta/*, and
// compatibility aliases), and static/dev UI fallback routes.
type Server struct {
	cfg        config.Config
	logger     *slog.Logger
	auth       contracts.AuthVerifier
	root       *chi.Mux
	mounts     lifecycle.RouteMounts
	httpServer *http.Server

	mu        sync.Mutex
	started   bool
	closed    bool
	boundAddr string
}

// New constructs a Server with pre-wired PublicAPI, AdminAPI, and Gateway
// route mounts and safe HTTP server timeout defaults (WriteTimeout: 0 so SSE
// streaming routes are never cut off by the server).
func New(opts Options) *Server {
	logger := opts.Logger
	if logger == nil {
		logger = slog.Default()
	}

	s := &Server{
		cfg:    opts.Config,
		logger: logger,
		auth:   opts.Auth,
	}

	root := chi.NewRouter()
	root.Use(middleware.RequestID)

	// Control plane router mounted under /api with panic recovery. Standard
	// PublicAPI and AdminAPI mounts enforce APITimeout, while AdminStreamAPI
	// omits APITimeout so authenticated SSE streams survive long sessions.
	apiRouter := chi.NewRouter()
	apiRouter.Use(s.recoverJSON)
	apiRouter.NotFound(jsonNotFoundHandler)
	apiRouter.MethodNotAllowed(jsonMethodNotAllowedHandler)

	var publicAPI, adminAPI chi.Router
	if opts.Config.APITimeout > 0 {
		timeoutMW := apiTimeoutMiddleware(opts.Config.APITimeout)
		publicAPI = apiRouter.With(timeoutMW)
		adminAPI = apiRouter.With(timeoutMW, s.adminAuthMiddleware)
	} else {
		publicAPI = apiRouter.With()
		adminAPI = apiRouter.With(s.adminAuthMiddleware)
	}
	adminStreamAPI := apiRouter.With(s.adminAuthMiddleware)
	root.Mount("/api", apiRouter)

	// Streaming gateway router: panic recovery without APITimeout so SSE streams
	// cleanly past APITimeout. Routes are registered with full paths (e.g.
	// /v1/chat/completions, /v1beta/models/..., /codex/..., /responses, /models/...).
	gatewayRouter := chi.NewRouter()
	gatewayRouter.Use(s.recoverJSON)
	gatewayRouter.NotFound(jsonNotFoundHandler)
	gatewayRouter.MethodNotAllowed(jsonMethodNotAllowedHandler)

	gwDispatch := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		rctx := chi.NewRouteContext()
		gatewayRouter.ServeHTTP(w, r.WithContext(context.WithValue(r.Context(), chi.RouteCtxKey, rctx)))
	})
	for _, prefix := range gatewayPrefixes {
		root.Handle(prefix, gwDispatch)
		root.Handle(prefix+"/*", gwDispatch)
	}

	// Static or dev UI fallback for non-API, non-gateway paths.
	var uiHandler http.Handler
	switch {
	case opts.UI != nil:
		uiHandler = NewStaticUIHandler(opts.UI)
	case opts.DevUIHandler != nil:
		uiHandler = opts.DevUIHandler
	default:
		uiHandler = http.HandlerFunc(jsonNotFoundHandler)
	}
	wrappedUI := s.recoverUI(uiHandler)
	root.Handle("/", wrappedUI)
	root.Handle("/*", wrappedUI)
	root.NotFound(wrappedUI.ServeHTTP)

	s.root = root
	s.mounts = lifecycle.RouteMounts{
		PublicAPI:      publicAPI,
		AdminAPI:       adminAPI,
		AdminStreamAPI: adminStreamAPI,
		Gateway:        gatewayRouter,
	}
	s.httpServer = &http.Server{
		Addr:              opts.Config.Addr(),
		Handler:           root,
		ReadHeaderTimeout: defaultReadHeaderTimeout,
		IdleTimeout:       defaultIdleTimeout,
		WriteTimeout:      0,
		ErrorLog:          slog.NewLogLogger(logger.Handler(), slog.LevelError),
	}

	return s
}

// Mounts returns the three chi.Router tiers onto which modules register their
// endpoints during application bootstrap.
func (s *Server) Mounts() lifecycle.RouteMounts {
	return s.mounts
}

// Handler returns the root http.Handler for testing or embedding.
func (s *Server) Handler() http.Handler {
	return s.root
}

// Addr returns the bound listener address once Start has bound a socket, or the
// configured listen address before Start is called.
func (s *Server) Addr() string {
	s.mu.Lock()
	defer s.mu.Unlock()
	if s.boundAddr != "" {
		return s.boundAddr
	}
	return s.httpServer.Addr
}

// Start binds the configured TCP address and serves HTTP requests until
// Shutdown is called or ctx is canceled. It returns nil on graceful shutdown.
func (s *Server) Start(ctx context.Context) error {
	if ctx == nil {
		ctx = context.Background()
	}
	if err := ctx.Err(); err != nil {
		return err
	}

	s.mu.Lock()
	if s.closed {
		s.mu.Unlock()
		return nil
	}
	if s.started {
		s.mu.Unlock()
		return errors.New("http server already started")
	}
	s.started = true
	addr := s.httpServer.Addr
	s.mu.Unlock()

	var lc net.ListenConfig
	ln, err := lc.Listen(ctx, "tcp", addr)
	if err != nil {
		s.mu.Lock()
		s.started = false
		s.mu.Unlock()
		return err
	}

	s.mu.Lock()
	if s.closed {
		s.mu.Unlock()
		_ = ln.Close()
		return nil
	}
	s.boundAddr = ln.Addr().String()
	s.mu.Unlock()

	s.logger.Info("http server listening", "addr", ln.Addr().String())

	serveDone := make(chan struct{})
	watcherDone := make(chan struct{})

	go func() {
		defer close(watcherDone)
		select {
		case <-ctx.Done():
			shutCtx, cancel := context.WithTimeout(context.WithoutCancel(ctx), cancelShutdownGrace)
			defer cancel()
			_ = s.Shutdown(shutCtx)
		case <-serveDone:
		}
	}()

	err = s.httpServer.Serve(ln)
	close(serveDone)
	<-watcherDone
	if errors.Is(err, http.ErrServerClosed) {
		return nil
	}
	return err
}

// Shutdown gracefully drains active connections until ctx expires.
func (s *Server) Shutdown(ctx context.Context) error {
	if ctx == nil {
		ctx = context.Background()
	}

	s.mu.Lock()
	s.closed = true
	s.mu.Unlock()

	return s.httpServer.Shutdown(ctx)
}

func (s *Server) adminAuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		remoteEnforced := !s.cfg.IsLoopback() && !s.cfg.AllowInsecureRemote
		if s.auth == nil {
			if remoteEnforced {
				WriteError(w, r, http.StatusUnauthorized, "unauthorized", "authentication required")
				return
			}
			next.ServeHTTP(w, r)
			return
		}

		enforced := s.auth.RequireLogin(r.Context()) || remoteEnforced
		if !enforced {
			next.ServeHTTP(w, r)
			return
		}

		subject, err := s.auth.AuthenticateAdmin(r)
		if err != nil {
			WriteError(w, r, http.StatusUnauthorized, "unauthorized", "authentication required")
			return
		}

		ctx := context.WithValue(r.Context(), adminSubjectKey{}, subject)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func (s *Server) recoverJSON(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		tw := wrapTrackedResponseWriter(w)
		defer func() {
			if rec := recover(); rec != nil {
				if err, ok := rec.(error); ok && errors.Is(err, http.ErrAbortHandler) {
					panic(rec)
				}
				s.logger.Error("http handler panicked",
					"request_id", middleware.GetReqID(r.Context()),
					"method", r.Method,
					"path", r.URL.Path,
					"panic", rec,
					"stack", string(debug.Stack()),
				)
				if !tw.wroteHeader {
					WriteError(tw, r, http.StatusInternalServerError, "internal_error", "internal server error")
				}
			}
		}()
		next.ServeHTTP(tw, r)
	})
}

func (s *Server) recoverUI(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		tw := wrapTrackedResponseWriter(w)
		defer func() {
			if rec := recover(); rec != nil {
				if err, ok := rec.(error); ok && errors.Is(err, http.ErrAbortHandler) {
					panic(rec)
				}
				s.logger.Error("ui handler panicked",
					"request_id", middleware.GetReqID(r.Context()),
					"method", r.Method,
					"path", r.URL.Path,
					"panic", rec,
					"stack", string(debug.Stack()),
				)
				if !tw.wroteHeader {
					http.Error(tw, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
					return
				}
			}
		}()
		next.ServeHTTP(tw, r)
	})
}

func apiTimeoutMiddleware(timeout time.Duration) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			tw := wrapTrackedResponseWriter(w)
			ctx, cancel := context.WithTimeout(r.Context(), timeout)
			defer func() {
				err := ctx.Err()
				cancel()
				if errors.Is(err, context.DeadlineExceeded) && !tw.wroteHeader {
					WriteError(tw, r, http.StatusGatewayTimeout, "timeout", "request timed out")
				}
			}()
			next.ServeHTTP(tw, r.WithContext(ctx))
		})
	}
}

func jsonNotFoundHandler(w http.ResponseWriter, r *http.Request) {
	WriteError(w, r, http.StatusNotFound, "not_found", "route not found")
}

func jsonMethodNotAllowedHandler(w http.ResponseWriter, r *http.Request) {
	WriteError(w, r, http.StatusMethodNotAllowed, "method_not_allowed", "method not allowed")
}

type trackedResponseWriter struct {
	http.ResponseWriter
	wroteHeader bool
}

func wrapTrackedResponseWriter(w http.ResponseWriter) *trackedResponseWriter {
	if existing, ok := w.(*trackedResponseWriter); ok {
		return existing
	}
	return &trackedResponseWriter{ResponseWriter: w}
}

func (w *trackedResponseWriter) WriteHeader(code int) {
	if w.wroteHeader {
		return
	}
	if code >= 200 {
		w.wroteHeader = true
	}
	w.ResponseWriter.WriteHeader(code)
}

func (w *trackedResponseWriter) Write(p []byte) (int, error) {
	if !w.wroteHeader {
		w.wroteHeader = true
	}
	return w.ResponseWriter.Write(p)
}

func (w *trackedResponseWriter) Flush() {
	w.wroteHeader = true
	if f, ok := w.ResponseWriter.(http.Flusher); ok {
		f.Flush()
	}
}

func (w *trackedResponseWriter) Unwrap() http.ResponseWriter {
	return w.ResponseWriter
}
