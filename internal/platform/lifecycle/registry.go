// Package lifecycle defines the modular monolith Module contract and a
// deterministic topological Registry that coordinates module migrations,
// route mounting, startup, rollback on partial start failure, and reverse-order
// graceful shutdown.
package lifecycle

import (
	"container/heap"
	"context"
	"errors"
	"fmt"
	"io/fs"
	"slices"
	"strings"
	"sync"
	"time"

	"github.com/go-chi/chi/v5"

	"github.com/omeg4monx/router-gateway/internal/platform/db"
)

const rollbackStopTimeout = 10 * time.Second

// Sentinel errors returned by Registry operations.
var (
	ErrNilModule          = errors.New("module must not be nil")
	ErrEmptyModuleName    = errors.New("module name must not be empty")
	ErrDuplicateModule    = errors.New("duplicate module name")
	ErrInvalidDependency  = errors.New("module dependency name must not be empty")
	ErrMissingDependency  = errors.New("missing module dependency")
	ErrDisabledDependency = errors.New("module depends on disabled module")
	ErrDependencyCycle    = errors.New("cyclic module dependency detected")
	ErrAlreadyStarted     = errors.New("registry modules already started")
)

// RouteMounts provides the HTTP router tiers onto which modules mount their
// endpoints during bootstrap. AdminStreamAPI omits APITimeout for long-lived
// authenticated SSE endpoints.
type RouteMounts struct {
	PublicAPI      chi.Router
	AdminAPI       chi.Router
	AdminStreamAPI chi.Router
	Gateway        chi.Router
}

// Module is implemented by every vertical-slice feature module registered with
// the application bootstrapper.
type Module interface {
	Name() string
	Dependencies() []string
	Migrations() fs.FS
	RegisterRoutes(mounts RouteMounts)
	Start(ctx context.Context) error
	Stop(ctx context.Context) error
}

type registryState int

const (
	stateIdle registryState = iota
	stateStarting
	stateRunning
	stateStopping
)

// Registry tracks registered modules, filters disabled optional modules,
// validates the dependency graph, and drives ordered lifecycle transitions.
type Registry struct {
	lifecycleMu sync.Mutex
	mu          sync.Mutex
	state       registryState
	disabled    map[string]bool
	registered  map[string]bool
	modules     map[string]Module
	started     []Module
}

// NewRegistry creates a Registry initialized with a normalized defensive copy
// of disabled module names.
func NewRegistry(disabled map[string]bool) *Registry {
	normDisabled := make(map[string]bool, len(disabled))
	for k, v := range disabled {
		if !v {
			continue
		}
		if name := normalizeName(k); name != "" {
			normDisabled[name] = true
		}
	}

	return &Registry{
		disabled:   normDisabled,
		registered: make(map[string]bool),
		modules:    make(map[string]Module),
	}
}

// Register validates and adds a module to the registry. Duplicate module names
// are always rejected (even if disabled). Modules marked disabled in the
// registry configuration are recorded as seen but excluded from active
// lifecycle execution.
func (r *Registry) Register(m Module) error {
	if m == nil {
		return ErrNilModule
	}

	name := normalizeName(m.Name())
	if name == "" {
		return ErrEmptyModuleName
	}

	r.mu.Lock()
	defer r.mu.Unlock()

	if r.state != stateIdle {
		return ErrAlreadyStarted
	}

	if r.registered[name] {
		return fmt.Errorf("%w: %q", ErrDuplicateModule, name)
	}
	r.registered[name] = true

	if r.disabled[name] {
		return nil
	}

	r.modules[name] = m
	return nil
}

// ActiveNames returns the normalized names of all enabled registered modules in
// deterministic lexicographical order.
func (r *Registry) ActiveNames() []string {
	r.mu.Lock()
	defer r.mu.Unlock()

	names := make([]string, 0, len(r.modules))
	for name := range r.modules {
		names = append(names, name)
	}
	slices.Sort(names)
	return names
}

// Ordered resolves the active module dependency graph using Kahn's topological
// sort with lexicographical tie-breaking by module name. It rejects missing
// dependencies, dependencies on disabled modules, and cycles.
func (r *Registry) Ordered() ([]Module, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	return r.orderedLocked()
}

func (r *Registry) orderedLocked() ([]Module, error) {
	inDegree := make(map[string]int, len(r.modules))
	dependents := make(map[string][]string, len(r.modules))

	activeNames := make([]string, 0, len(r.modules))
	for name := range r.modules {
		activeNames = append(activeNames, name)
		inDegree[name] = 0
	}
	// Sort active names before inspecting dependencies so error reporting is
	// deterministic even when multiple invalid edges exist.
	slices.Sort(activeNames)

	for _, name := range activeNames {
		mod := r.modules[name]
		seenDeps := make(map[string]bool)

		for _, rawDep := range mod.Dependencies() {
			dep := normalizeName(rawDep)
			if dep == "" {
				return nil, fmt.Errorf("module %q: %w", name, ErrInvalidDependency)
			}
			if seenDeps[dep] {
				continue
			}
			seenDeps[dep] = true

			if r.disabled[dep] {
				return nil, fmt.Errorf("module %q depends on %q: %w", name, dep, ErrDisabledDependency)
			}
			if _, ok := r.modules[dep]; !ok {
				return nil, fmt.Errorf("module %q depends on %q: %w", name, dep, ErrMissingDependency)
			}

			dependents[dep] = append(dependents[dep], name)
			inDegree[name]++
		}
	}

	ready := &nameMinHeap{}
	heap.Init(ready)
	for _, name := range activeNames {
		if inDegree[name] == 0 {
			heap.Push(ready, name)
		}
	}

	ordered := make([]Module, 0, len(r.modules))
	for ready.Len() > 0 {
		curr, ok := heap.Pop(ready).(string)
		if !ok {
			continue
		}
		ordered = append(ordered, r.modules[curr])

		for _, next := range dependents[curr] {
			inDegree[next]--
			if inDegree[next] == 0 {
				heap.Push(ready, next)
			}
		}
	}

	if len(ordered) != len(r.modules) {
		var cyclic []string
		for _, name := range activeNames {
			if inDegree[name] > 0 {
				cyclic = append(cyclic, name)
			}
		}
		return nil, fmt.Errorf("%w: involving [%s]", ErrDependencyCycle, strings.Join(cyclic, ", "))
	}

	return ordered, nil
}

// MigrateAll runs each active module's embedded SQL migrations in topological
// dependency order.
func (r *Registry) MigrateAll(ctx context.Context, database *db.DB) error {
	if database == nil {
		return errors.New("database must not be nil")
	}
	if ctx == nil {
		ctx = context.Background()
	}

	ordered, err := r.Ordered()
	if err != nil {
		return err
	}

	for _, mod := range ordered {
		if err := ctx.Err(); err != nil {
			return err
		}
		fsys := mod.Migrations()
		if fsys == nil {
			continue
		}
		name := normalizeName(mod.Name())
		if err := database.MigrateModule(ctx, name, fsys); err != nil {
			return fmt.Errorf("migrate module %q: %w", name, err)
		}
	}

	return nil
}

// MountAll invokes RegisterRoutes on every active module in topological
// dependency order.
func (r *Registry) MountAll(mounts RouteMounts) error {
	ordered, err := r.Ordered()
	if err != nil {
		return err
	}

	for _, mod := range ordered {
		mod.RegisterRoutes(mounts)
	}

	return nil
}

// StartAll starts all active modules in topological order. If ctx is already
// canceled before starting a module or if module i fails to start, StartAll
// automatically invokes Stop on already-started modules i-1..0 in reverse order
// with a fresh bounded rollback context before returning the joined error.
func (r *Registry) StartAll(ctx context.Context) error {
	if ctx == nil {
		ctx = context.Background()
	}

	r.lifecycleMu.Lock()
	defer r.lifecycleMu.Unlock()

	r.mu.Lock()
	if r.state != stateIdle {
		r.mu.Unlock()
		return ErrAlreadyStarted
	}
	ordered, err := r.orderedLocked()
	if err != nil {
		r.mu.Unlock()
		return err
	}
	r.state = stateStarting
	r.mu.Unlock()

	rollback := func(startErr error, started []Module) error {
		errs := []error{startErr}

		rollbackCtx, cancel := context.WithTimeout(context.WithoutCancel(ctx), rollbackStopTimeout)
		defer cancel()

		for i := len(started) - 1; i >= 0; i-- {
			prev := started[i]
			if stopErr := prev.Stop(rollbackCtx); stopErr != nil {
				errs = append(errs, fmt.Errorf("rollback stop module %q: %w", normalizeName(prev.Name()), stopErr))
			}
		}

		r.mu.Lock()
		r.started = nil
		r.state = stateIdle
		r.mu.Unlock()

		return errors.Join(errs...)
	}

	if err := ctx.Err(); err != nil {
		r.mu.Lock()
		r.state = stateIdle
		r.mu.Unlock()
		return err
	}

	started := make([]Module, 0, len(ordered))
	for _, mod := range ordered {
		name := normalizeName(mod.Name())
		if err := ctx.Err(); err != nil {
			return rollback(fmt.Errorf("start module %q: %w", name, err), started)
		}
		if err := mod.Start(ctx); err != nil {
			return rollback(fmt.Errorf("start module %q: %w", name, err), started)
		}
		started = append(started, mod)
	}

	r.mu.Lock()
	r.started = started
	r.state = stateRunning
	r.mu.Unlock()
	return nil
}

// StopAll stops all started modules in exact reverse startup order, joining
// any errors returned by individual modules. StopAll is idempotent.
func (r *Registry) StopAll(ctx context.Context) error {
	if ctx == nil {
		ctx = context.Background()
	}

	r.lifecycleMu.Lock()
	defer r.lifecycleMu.Unlock()

	r.mu.Lock()
	if r.state == stateIdle {
		r.mu.Unlock()
		return nil
	}
	started := r.started
	r.started = nil
	r.state = stateStopping
	r.mu.Unlock()

	var errs []error
	for i := len(started) - 1; i >= 0; i-- {
		mod := started[i]
		if err := mod.Stop(ctx); err != nil {
			errs = append(errs, fmt.Errorf("stop module %q: %w", normalizeName(mod.Name()), err))
		}
	}

	r.mu.Lock()
	r.state = stateIdle
	r.mu.Unlock()

	return errors.Join(errs...)
}

func normalizeName(name string) string {
	return strings.ToLower(strings.TrimSpace(name))
}

type nameMinHeap []string

func (h nameMinHeap) Len() int           { return len(h) }
func (h nameMinHeap) Less(i, j int) bool { return h[i] < h[j] }
func (h nameMinHeap) Swap(i, j int)      { h[i], h[j] = h[j], h[i] }

func (h *nameMinHeap) Push(x any) {
	if s, ok := x.(string); ok {
		*h = append(*h, s)
	}
}

func (h *nameMinHeap) Pop() any {
	old := *h
	n := len(old)
	item := old[n-1]
	*h = old[:n-1]
	return item
}
