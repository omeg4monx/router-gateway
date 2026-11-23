// Package jobs provides a panic-safe periodic background task runner with
// graceful shutdown semantics, live status snapshots, and overlap-guarded
// manual triggers.
package jobs

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"runtime/debug"
	"slices"
	"strings"
	"sync"
	"time"
)

// ErrJobNotFound indicates that no registered job matches the requested name.
var ErrJobNotFound = errors.New("jobs: job not found")

// ErrJobAlreadyRunning indicates that the requested job is already executing.
var ErrJobAlreadyRunning = errors.New("jobs: job already running")

// JobStatus is a point-in-time snapshot of a registered periodic job's
// schedule and execution telemetry.
type JobStatus struct {
	Name           string
	IntervalMs     int64
	Running        bool
	RunCount       int64
	ErrorCount     int64
	LastStatus     string
	LastError      string
	LastDurationMs int64
	LastStartedAt  *time.Time
	LastFinishedAt *time.Time
	NextRunAt      *time.Time
}

type job struct {
	name     string
	interval time.Duration
	fn       func(context.Context) error

	mu             sync.Mutex
	running        bool
	runCount       int64
	errorCount     int64
	lastStatus     string
	lastError      string
	lastDurationMs int64
	lastStartedAt  *time.Time
	lastFinishedAt *time.Time
	nextRunAt      *time.Time
}

func (j *job) snapshot() JobStatus {
	j.mu.Lock()
	defer j.mu.Unlock()

	st := JobStatus{
		Name:           j.name,
		IntervalMs:     j.interval.Milliseconds(),
		Running:        j.running,
		RunCount:       j.runCount,
		ErrorCount:     j.errorCount,
		LastStatus:     j.lastStatus,
		LastError:      j.lastError,
		LastDurationMs: j.lastDurationMs,
	}
	if j.lastStartedAt != nil {
		t := *j.lastStartedAt
		st.LastStartedAt = &t
	}
	if j.lastFinishedAt != nil {
		t := *j.lastFinishedAt
		st.LastFinishedAt = &t
	}
	if j.nextRunAt != nil {
		t := *j.nextRunAt
		st.NextRunAt = &t
	}
	return st
}

// Runner manages periodic background jobs, recovering from panics on every
// tick and coordinating clean shutdown via Stop.
type Runner struct {
	logger   *slog.Logger
	mu       sync.Mutex
	jobs     []*job
	byName   map[string]*job
	observer func(JobStatus)
	started  bool
	cancel   context.CancelFunc
	wg       sync.WaitGroup
}

// NewRunner creates a Runner that logs job errors and panics to logger.
// If logger is nil, slog.Default() is used.
func NewRunner(logger *slog.Logger) *Runner {
	if logger == nil {
		logger = slog.Default()
	}
	return &Runner{
		logger: logger,
		byName: make(map[string]*job),
	}
}

// SetObserver registers an optional callback invoked after each job execution
// completes (whether ok, error, or panic). Pass nil to clear.
func (r *Runner) SetObserver(fn func(JobStatus)) {
	r.mu.Lock()
	r.observer = fn
	r.mu.Unlock()
}

// Every registers a periodic job to run at the given interval.
// Registrations with empty names, non-positive intervals, nil functions,
// duplicate names, or made after Start are ignored.
func (r *Runner) Every(name string, interval time.Duration, fn func(context.Context) error) {
	trimmed := strings.TrimSpace(name)
	if trimmed == "" || interval <= 0 || fn == nil {
		r.logger.Warn("ignoring invalid periodic job registration",
			"job", name,
			"interval", interval,
			"nil_fn", fn == nil,
		)
		return
	}

	r.mu.Lock()
	defer r.mu.Unlock()

	if r.started {
		r.logger.Warn("ignoring periodic job registration after runner started", "job", trimmed)
		return
	}
	if _, exists := r.byName[trimmed]; exists {
		r.logger.Warn("ignoring duplicate periodic job registration", "job", trimmed)
		return
	}

	j := &job{
		name:       trimmed,
		interval:   interval,
		fn:         fn,
		lastStatus: "idle",
	}
	r.jobs = append(r.jobs, j)
	r.byName[trimmed] = j
}

// Snapshot returns a sorted point-in-time status slice for all registered jobs.
func (r *Runner) Snapshot() []JobStatus {
	r.mu.Lock()
	list := make([]*job, len(r.jobs))
	copy(list, r.jobs)
	r.mu.Unlock()

	out := make([]JobStatus, 0, len(list))
	for _, j := range list {
		out = append(out, j.snapshot())
	}
	slices.SortFunc(out, func(a, b JobStatus) int {
		return strings.Compare(a.Name, b.Name)
	})
	return out
}

// Trigger runs a registered job immediately on demand. It returns ErrJobNotFound
// if name is unknown, or ErrJobAlreadyRunning if the job is currently executing.
func (r *Runner) Trigger(ctx context.Context, name string) error {
	if ctx == nil {
		ctx = context.Background()
	}

	r.mu.Lock()
	j := r.byName[strings.TrimSpace(name)]
	r.mu.Unlock()
	if j == nil {
		return ErrJobNotFound
	}

	return r.executeJob(ctx, j, false)
}

// Start launches all registered periodic jobs in dedicated goroutines.
// Calling Start more than once is a safe no-op.
func (r *Runner) Start(ctx context.Context) {
	if ctx == nil {
		ctx = context.Background()
	}

	r.mu.Lock()
	if r.started {
		r.mu.Unlock()
		return
	}

	runCtx, cancel := context.WithCancel(ctx)
	r.cancel = cancel
	r.started = true

	snapshot := make([]*job, len(r.jobs))
	copy(snapshot, r.jobs)
	r.wg.Add(len(snapshot))
	r.mu.Unlock()

	now := time.Now().UTC()
	for _, j := range snapshot {
		next := now.Add(j.interval)
		j.mu.Lock()
		j.nextRunAt = &next
		j.mu.Unlock()
		go r.loop(runCtx, j)
	}
}

// Stop cancels all running jobs and waits for their goroutines to exit or for
// ctx to expire. Calling Stop before Start or multiple times is safe.
func (r *Runner) Stop(ctx context.Context) error {
	if ctx == nil {
		ctx = context.Background()
	}

	r.mu.Lock()
	if !r.started {
		r.mu.Unlock()
		return nil
	}
	cancel := r.cancel
	r.mu.Unlock()

	if cancel != nil {
		cancel()
	}

	done := make(chan struct{})
	go func() {
		r.wg.Wait()
		close(done)
	}()

	select {
	case <-done:
		return nil
	case <-ctx.Done():
		return ctx.Err()
	}
}

func (r *Runner) loop(ctx context.Context, j *job) {
	defer r.wg.Done()

	ticker := time.NewTicker(j.interval)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			if ctx.Err() != nil {
				return
			}
			_ = r.executeJob(ctx, j, true)
			next := time.Now().UTC().Add(j.interval)
			j.mu.Lock()
			j.nextRunAt = &next
			j.mu.Unlock()
		}
	}
}

func (r *Runner) executeJob(ctx context.Context, j *job, scheduled bool) (runErr error) {
	startedAt := time.Now().UTC()

	j.mu.Lock()
	if j.running {
		j.mu.Unlock()
		return ErrJobAlreadyRunning
	}
	j.running = true
	j.lastStatus = "running"
	j.lastStartedAt = &startedAt
	j.mu.Unlock()

	panicked := false
	defer func() {
		if rec := recover(); rec != nil {
			panicked = true
			runErr = fmt.Errorf("panic: %v", rec)
			r.logger.Error("periodic job panicked",
				"job", j.name,
				"panic", rec,
				"stack", string(debug.Stack()),
			)
		}

		finishedAt := time.Now().UTC()
		durMs := max(finishedAt.Sub(startedAt).Milliseconds(), 0)

		j.mu.Lock()
		j.running = false
		j.runCount++
		j.lastDurationMs = durMs
		j.lastFinishedAt = &finishedAt
		switch {
		case panicked:
			j.errorCount++
			j.lastStatus = "panic"
			j.lastError = runErr.Error()
		case runErr != nil && (ctx.Err() == nil || !errors.Is(runErr, context.Canceled)):
			j.errorCount++
			j.lastStatus = "error"
			j.lastError = runErr.Error()
			if scheduled {
				r.logger.Error("periodic job failed",
					"job", j.name,
					"error", runErr,
				)
			}
		default:
			j.lastStatus = "ok"
			j.lastError = ""
		}
		j.mu.Unlock()

		r.mu.Lock()
		obs := r.observer
		r.mu.Unlock()
		if obs != nil {
			obs(j.snapshot())
		}
	}()

	runErr = j.fn(ctx)
	return runErr
}
