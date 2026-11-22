package lifecycle_test

import (
	"context"
	"errors"
	"io/fs"
	"net/http"
	"net/http/httptest"
	"path/filepath"
	"reflect"
	"sync"
	"testing"
	"testing/fstest"

	"github.com/go-chi/chi/v5"

	"github.com/omeg4monx/router-gateway/internal/platform/db"
	"github.com/omeg4monx/router-gateway/internal/platform/lifecycle"
)

type stubModule struct {
	name       string
	deps       []string
	migrations fs.FS
	mountFn    func(mounts lifecycle.RouteMounts)
	startFn    func(ctx context.Context) error
	stopFn     func(ctx context.Context) error
}

func (m *stubModule) Name() string           { return m.name }
func (m *stubModule) Dependencies() []string { return m.deps }
func (m *stubModule) Migrations() fs.FS      { return m.migrations }

func (m *stubModule) RegisterRoutes(mounts lifecycle.RouteMounts) {
	if m.mountFn != nil {
		m.mountFn(mounts)
	}
}

func (m *stubModule) Start(ctx context.Context) error {
	if m.startFn != nil {
		return m.startFn(ctx)
	}
	return nil
}

func (m *stubModule) Stop(ctx context.Context) error {
	if m.stopFn != nil {
		return m.stopFn(ctx)
	}
	return nil
}

func TestRegistryTopologicalOrderAndReverseStop(t *testing.T) {
	var (
		mu         sync.Mutex
		startOrder []string
		stopOrder  []string
	)

	newTracked := func(name string, deps ...string) *stubModule {
		return &stubModule{
			name: name,
			deps: deps,
			startFn: func(_ context.Context) error {
				mu.Lock()
				startOrder = append(startOrder, name)
				mu.Unlock()
				return nil
			},
			stopFn: func(_ context.Context) error {
				mu.Lock()
				stopOrder = append(stopOrder, name)
				mu.Unlock()
				return nil
			},
		}
	}

	reg := lifecycle.NewRegistry(nil)
	// Register in reverse dependency order: C -> B -> A.
	if err := reg.Register(newTracked("c", "b")); err != nil {
		t.Fatalf("register c: %v", err)
	}
	if err := reg.Register(newTracked("b", "a")); err != nil {
		t.Fatalf("register b: %v", err)
	}
	if err := reg.Register(newTracked("a")); err != nil {
		t.Fatalf("register a: %v", err)
	}

	ordered, err := reg.Ordered()
	if err != nil {
		t.Fatalf("Ordered() error: %v", err)
	}
	gotNames := make([]string, len(ordered))
	for i, m := range ordered {
		gotNames[i] = m.Name()
	}
	wantStart := []string{"a", "b", "c"}
	if !reflect.DeepEqual(gotNames, wantStart) {
		t.Fatalf("Ordered() = %v, want %v", gotNames, wantStart)
	}

	ctx := context.Background()
	if err := reg.StartAll(ctx); err != nil {
		t.Fatalf("StartAll() error: %v", err)
	}
	if !reflect.DeepEqual(startOrder, wantStart) {
		t.Fatalf("StartAll order = %v, want %v", startOrder, wantStart)
	}

	if err := reg.StopAll(ctx); err != nil {
		t.Fatalf("StopAll() error: %v", err)
	}
	wantStop := []string{"c", "b", "a"}
	if !reflect.DeepEqual(stopOrder, wantStop) {
		t.Fatalf("StopAll order = %v, want %v", stopOrder, wantStop)
	}

	// Second StopAll must be idempotent and not invoke Stop again.
	if err := reg.StopAll(ctx); err != nil {
		t.Fatalf("second StopAll() error: %v", err)
	}
	if !reflect.DeepEqual(stopOrder, wantStop) {
		t.Fatalf("StopAll was not idempotent, stopOrder = %v, want %v", stopOrder, wantStop)
	}
}

func TestRegistryDeterministicLexicographicalTieBreak(t *testing.T) {
	// Independent and branching modules registered in different orders must
	// always resolve to the exact same lexicographical tie-break order.
	want := []string{"alpha", "beta", "delta", "gamma", "zeta"}

	orders := [][]struct {
		name string
		deps []string
	}{
		{
			{name: "zeta", deps: []string{"beta"}},
			{name: "gamma", deps: []string{"alpha"}},
			{name: "delta", deps: nil},
			{name: "beta", deps: nil},
			{name: "alpha", deps: nil},
		},
		{
			{name: "alpha", deps: nil},
			{name: "beta", deps: nil},
			{name: "delta", deps: nil},
			{name: "gamma", deps: []string{"alpha"}},
			{name: "zeta", deps: []string{"beta"}},
		},
		{
			{name: "delta", deps: nil},
			{name: "zeta", deps: []string{"beta"}},
			{name: "alpha", deps: nil},
			{name: "gamma", deps: []string{"alpha"}},
			{name: "beta", deps: nil},
		},
	}

	for idx, mods := range orders {
		reg := lifecycle.NewRegistry(nil)
		for _, spec := range mods {
			if err := reg.Register(&stubModule{name: spec.name, deps: spec.deps}); err != nil {
				t.Fatalf("case %d register %s: %v", idx, spec.name, err)
			}
		}
		ordered, err := reg.Ordered()
		if err != nil {
			t.Fatalf("case %d Ordered(): %v", idx, err)
		}
		got := make([]string, len(ordered))
		for i, m := range ordered {
			got[i] = m.Name()
		}
		if !reflect.DeepEqual(got, want) {
			t.Fatalf("case %d Ordered() = %v, want %v", idx, got, want)
		}
	}
}

func TestRegistryRejectsCyclesAndMissingDeps(t *testing.T) {
	t.Run("two-node cycle", func(t *testing.T) {
		reg := lifecycle.NewRegistry(nil)
		if err := reg.Register(&stubModule{name: "a", deps: []string{"b"}}); err != nil {
			t.Fatalf("register a: %v", err)
		}
		if err := reg.Register(&stubModule{name: "b", deps: []string{"a"}}); err != nil {
			t.Fatalf("register b: %v", err)
		}
		_, err := reg.Ordered()
		if !errors.Is(err, lifecycle.ErrDependencyCycle) {
			t.Fatalf("Ordered() err = %v, want ErrDependencyCycle", err)
		}
	})

	t.Run("self cycle", func(t *testing.T) {
		reg := lifecycle.NewRegistry(nil)
		if err := reg.Register(&stubModule{name: "self", deps: []string{"self"}}); err != nil {
			t.Fatalf("register self: %v", err)
		}
		_, err := reg.Ordered()
		if !errors.Is(err, lifecycle.ErrDependencyCycle) {
			t.Fatalf("Ordered() err = %v, want ErrDependencyCycle", err)
		}
	})

	t.Run("missing dependency", func(t *testing.T) {
		reg := lifecycle.NewRegistry(nil)
		if err := reg.Register(&stubModule{name: "b", deps: []string{"x"}}); err != nil {
			t.Fatalf("register b: %v", err)
		}
		_, err := reg.Ordered()
		if !errors.Is(err, lifecycle.ErrMissingDependency) {
			t.Fatalf("Ordered() err = %v, want ErrMissingDependency", err)
		}
	})

	t.Run("duplicate active module name", func(t *testing.T) {
		reg := lifecycle.NewRegistry(nil)
		if err := reg.Register(&stubModule{name: "auth"}); err != nil {
			t.Fatalf("register first auth: %v", err)
		}
		err := reg.Register(&stubModule{name: "AUTH"})
		if !errors.Is(err, lifecycle.ErrDuplicateModule) {
			t.Fatalf("duplicate Register() err = %v, want ErrDuplicateModule", err)
		}
	})

	t.Run("duplicate disabled module name", func(t *testing.T) {
		reg := lifecycle.NewRegistry(map[string]bool{"mitm": true})
		if err := reg.Register(&stubModule{name: "mitm"}); err != nil {
			t.Fatalf("register first mitm: %v", err)
		}
		err := reg.Register(&stubModule{name: "mitm"})
		if !errors.Is(err, lifecycle.ErrDuplicateModule) {
			t.Fatalf("duplicate disabled Register() err = %v, want ErrDuplicateModule", err)
		}
	})

	t.Run("nil module and empty name", func(t *testing.T) {
		reg := lifecycle.NewRegistry(nil)
		if err := reg.Register(nil); !errors.Is(err, lifecycle.ErrNilModule) {
			t.Fatalf("Register(nil) err = %v, want ErrNilModule", err)
		}
		if err := reg.Register(&stubModule{name: "   "}); !errors.Is(err, lifecycle.ErrEmptyModuleName) {
			t.Fatalf("Register(blank) err = %v, want ErrEmptyModuleName", err)
		}
	})
}

func TestRegistrySkipsDisabledOptionalModule(t *testing.T) {
	disabled := map[string]bool{" MITM ": true, "ignored": false}
	reg := lifecycle.NewRegistry(disabled)
	// Mutating caller map after NewRegistry must not affect registry state.
	disabled["system"] = true

	mitmStarted := false
	systemStarted := false

	if err := reg.Register(&stubModule{
		name: "mitm",
		startFn: func(_ context.Context) error {
			mitmStarted = true
			return nil
		},
	}); err != nil {
		t.Fatalf("register mitm: %v", err)
	}
	if err := reg.Register(&stubModule{
		name: "system",
		startFn: func(_ context.Context) error {
			systemStarted = true
			return nil
		},
	}); err != nil {
		t.Fatalf("register system: %v", err)
	}

	wantActive := []string{"system"}
	if got := reg.ActiveNames(); !reflect.DeepEqual(got, wantActive) {
		t.Fatalf("ActiveNames() = %v, want %v", got, wantActive)
	}

	if err := reg.StartAll(context.Background()); err != nil {
		t.Fatalf("StartAll() error: %v", err)
	}
	if mitmStarted {
		t.Fatal("disabled module mitm must not be started")
	}
	if !systemStarted {
		t.Fatal("active module system must be started")
	}

	t.Run("active module depending on disabled module returns error", func(t *testing.T) {
		r := lifecycle.NewRegistry(map[string]bool{"mitm": true})
		if err := r.Register(&stubModule{name: "mitm"}); err != nil {
			t.Fatalf("register mitm: %v", err)
		}
		if err := r.Register(&stubModule{name: "proxy", deps: []string{"mitm"}}); err != nil {
			t.Fatalf("register proxy: %v", err)
		}
		_, err := r.Ordered()
		if !errors.Is(err, lifecycle.ErrDisabledDependency) {
			t.Fatalf("Ordered() err = %v, want ErrDisabledDependency", err)
		}
	})
}

func TestRegistryStartAllRollsBackOnPartialFailure(t *testing.T) {
	errBoom := errors.New("module c failed to start")
	errRollback := errors.New("module a rollback warning")

	var (
		started    []string
		stopped    []string
		stopCtxErr []error
	)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	reg := lifecycle.NewRegistry(nil)
	mustReg := func(m lifecycle.Module) {
		t.Helper()
		if err := reg.Register(m); err != nil {
			t.Fatalf("register %s: %v", m.Name(), err)
		}
	}

	mustReg(&stubModule{
		name: "a",
		startFn: func(_ context.Context) error {
			started = append(started, "a")
			return nil
		},
		stopFn: func(sctx context.Context) error {
			stopped = append(stopped, "a")
			stopCtxErr = append(stopCtxErr, sctx.Err())
			return errRollback
		},
	})
	mustReg(&stubModule{
		name: "b",
		deps: []string{"a"},
		startFn: func(_ context.Context) error {
			started = append(started, "b")
			return nil
		},
		stopFn: func(sctx context.Context) error {
			stopped = append(stopped, "b")
			stopCtxErr = append(stopCtxErr, sctx.Err())
			return nil
		},
	})
	mustReg(&stubModule{
		name: "c",
		deps: []string{"b"},
		startFn: func(_ context.Context) error {
			started = append(started, "c")
			// Simulate parent context cancellation happening alongside failure.
			cancel()
			return errBoom
		},
		stopFn: func(_ context.Context) error {
			stopped = append(stopped, "c")
			return nil
		},
	})
	mustReg(&stubModule{
		name: "d",
		deps: []string{"c"},
		startFn: func(_ context.Context) error {
			started = append(started, "d")
			return nil
		},
		stopFn: func(_ context.Context) error {
			stopped = append(stopped, "d")
			return nil
		},
	})

	err := reg.StartAll(ctx)
	if !errors.Is(err, errBoom) {
		t.Fatalf("StartAll() err = %v, want wrapped %v", err, errBoom)
	}
	if !errors.Is(err, errRollback) {
		t.Fatalf("StartAll() err = %v, want wrapped rollback error %v", err, errRollback)
	}

	wantStarted := []string{"a", "b", "c"}
	if !reflect.DeepEqual(started, wantStarted) {
		t.Fatalf("started = %v, want %v", started, wantStarted)
	}
	wantStopped := []string{"b", "a"}
	if !reflect.DeepEqual(stopped, wantStopped) {
		t.Fatalf("rollback stopped = %v, want %v", stopped, wantStopped)
	}
	for i, cerr := range stopCtxErr {
		if cerr != nil {
			t.Fatalf("rollback stop %d received cancelled context: %v", i, cerr)
		}
	}

	// StopAll after failed StartAll must be a no-op since rollback already stopped A and B.
	if err := reg.StopAll(context.Background()); err != nil {
		t.Fatalf("StopAll() after rollback error: %v", err)
	}
	if !reflect.DeepEqual(stopped, wantStopped) {
		t.Fatalf("StopAll() double-stopped modules: stopped = %v, want %v", stopped, wantStopped)
	}
}

func TestRegistryStopAllJoinsErrors(t *testing.T) {
	errStopA := errors.New("stop a failed")
	errStopC := errors.New("stop c failed")

	var stopped []string
	reg := lifecycle.NewRegistry(nil)
	for _, tc := range []struct {
		name string
		deps []string
		err  error
	}{
		{"a", nil, errStopA},
		{"b", []string{"a"}, nil},
		{"c", []string{"b"}, errStopC},
	} {
		tc := tc
		if err := reg.Register(&stubModule{
			name: tc.name,
			deps: tc.deps,
			stopFn: func(_ context.Context) error {
				stopped = append(stopped, tc.name)
				return tc.err
			},
		}); err != nil {
			t.Fatalf("register %s: %v", tc.name, err)
		}
	}

	ctx := context.Background()
	if err := reg.StartAll(ctx); err != nil {
		t.Fatalf("StartAll() error: %v", err)
	}

	err := reg.StopAll(ctx)
	if !errors.Is(err, errStopA) || !errors.Is(err, errStopC) {
		t.Fatalf("StopAll() err = %v, want both errStopA and errStopC", err)
	}
	wantStopped := []string{"c", "b", "a"}
	if !reflect.DeepEqual(stopped, wantStopped) {
		t.Fatalf("stopped = %v, want %v", stopped, wantStopped)
	}
}

func TestRegistryMigrateAllRunsInDependencyOrder(t *testing.T) {
	database, err := db.Open(filepath.Join(t.TempDir(), "registry.db"))
	if err != nil {
		t.Fatalf("db.Open: %v", err)
	}
	defer func() { _ = database.Close() }()

	parentFS := fstest.MapFS{
		"0001_init.sql": &fstest.MapFile{
			Data: []byte(`CREATE TABLE parent_items (id INTEGER PRIMARY KEY);`),
		},
	}
	childFS := fstest.MapFS{
		"0001_init.sql": &fstest.MapFile{
			Data: []byte(`CREATE TABLE child_items (id INTEGER PRIMARY KEY, parent_id INTEGER NOT NULL REFERENCES parent_items(id));`),
		},
	}
	disabledFS := fstest.MapFS{
		"0001_init.sql": &fstest.MapFile{
			Data: []byte(`CREATE TABLE should_not_exist (id INTEGER PRIMARY KEY);`),
		},
	}

	reg := lifecycle.NewRegistry(map[string]bool{"mitm": true})
	// Register child before parent to verify topological migration order.
	if err := reg.Register(&stubModule{name: "child", deps: []string{"parent"}, migrations: childFS}); err != nil {
		t.Fatalf("register child: %v", err)
	}
	if err := reg.Register(&stubModule{name: "parent", migrations: parentFS}); err != nil {
		t.Fatalf("register parent: %v", err)
	}
	if err := reg.Register(&stubModule{name: "stateless", deps: []string{"parent"}, migrations: nil}); err != nil {
		t.Fatalf("register stateless: %v", err)
	}
	if err := reg.Register(&stubModule{name: "mitm", migrations: disabledFS}); err != nil {
		t.Fatalf("register mitm: %v", err)
	}

	ctx := context.Background()
	if err := reg.MigrateAll(ctx, nil); err == nil {
		t.Fatal("MigrateAll(ctx, nil) expected error, got nil")
	}

	if err := reg.MigrateAll(ctx, database); err != nil {
		t.Fatalf("MigrateAll() error: %v", err)
	}

	for _, mod := range []string{"parent", "child"} {
		versions, err := database.AppliedVersions(ctx, mod)
		if err != nil {
			t.Fatalf("AppliedVersions(%s): %v", mod, err)
		}
		if !reflect.DeepEqual(versions, []int{1}) {
			t.Fatalf("AppliedVersions(%s) = %v, want [1]", mod, versions)
		}
	}

	mitmVersions, err := database.AppliedVersions(ctx, "mitm")
	if err != nil {
		t.Fatalf("AppliedVersions(mitm): %v", err)
	}
	if len(mitmVersions) != 0 {
		t.Fatalf("disabled module mitm was migrated: %v", mitmVersions)
	}
}

func TestRegistryMountAllMountsInDependencyOrder(t *testing.T) {
	var mountOrder []string

	reg := lifecycle.NewRegistry(nil)
	if err := reg.Register(&stubModule{
		name: "gateway",
		deps: []string{"system"},
		mountFn: func(mounts lifecycle.RouteMounts) {
			mountOrder = append(mountOrder, "gateway")
			mounts.Gateway.Get("/v1/models", func(w http.ResponseWriter, _ *http.Request) {
				w.WriteHeader(http.StatusOK)
			})
		},
	}); err != nil {
		t.Fatalf("register gateway: %v", err)
	}
	if err := reg.Register(&stubModule{
		name: "system",
		mountFn: func(mounts lifecycle.RouteMounts) {
			mountOrder = append(mountOrder, "system")
			mounts.PublicAPI.Get("/api/health", func(w http.ResponseWriter, _ *http.Request) {
				w.WriteHeader(http.StatusOK)
			})
			mounts.AdminAPI.Get("/api/admin/ping", func(w http.ResponseWriter, _ *http.Request) {
				w.WriteHeader(http.StatusAccepted)
			})
		},
	}); err != nil {
		t.Fatalf("register system: %v", err)
	}

	mounts := lifecycle.RouteMounts{
		PublicAPI: chi.NewRouter(),
		AdminAPI:  chi.NewRouter(),
		Gateway:   chi.NewRouter(),
	}

	if err := reg.MountAll(mounts); err != nil {
		t.Fatalf("MountAll() error: %v", err)
	}
	wantOrder := []string{"system", "gateway"}
	if !reflect.DeepEqual(mountOrder, wantOrder) {
		t.Fatalf("mount order = %v, want %v", mountOrder, wantOrder)
	}

	rec := httptest.NewRecorder()
	mounts.PublicAPI.ServeHTTP(rec, httptest.NewRequest(http.MethodGet, "/api/health", nil))
	if rec.Code != http.StatusOK {
		t.Fatalf("PublicAPI /api/health status = %d, want 200", rec.Code)
	}

	rec = httptest.NewRecorder()
	mounts.AdminAPI.ServeHTTP(rec, httptest.NewRequest(http.MethodGet, "/api/admin/ping", nil))
	if rec.Code != http.StatusAccepted {
		t.Fatalf("AdminAPI /api/admin/ping status = %d, want 202", rec.Code)
	}

	rec = httptest.NewRecorder()
	mounts.Gateway.ServeHTTP(rec, httptest.NewRequest(http.MethodGet, "/v1/models", nil))
	if rec.Code != http.StatusOK {
		t.Fatalf("Gateway /v1/models status = %d, want 200", rec.Code)
	}
}
