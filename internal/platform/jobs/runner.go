// Package jobs provides a panic-safe periodic background task runner with
// graceful shutdown semantics.
package jobs

import (
	"context"
	"errors"
	"log/slog"
	"runtime/debug"
	"sync"
	"time"
)

type job struct {
	name     string
	interval time.Duration
	fn       func(context.Context) error
}

// Runner manages periodic background jobs, recovering from panics on every
// tick and coordinating clean shutdown via Stop.
type Runner struct {
	logger  *slog.Logger
	mu      sync.Mutex
	jobs    []job
	started bool
	cancel  context.CancelFunc
	wg      sync.WaitGroup
}

// NewRunner creates a Runner that logs job errors and panics to logger.
// If logger is nil, slog.Default() is used.
func NewRunner(logger *slog.Logger) *Runner {
	if logger == nil {
		logger = slog.Default()
	}
	return &Runner{
		logger: logger,
	}
}

// Every registers a periodic job to run at the given interval.
// Registrations with non-positive intervals, nil functions, or made after Start
// are ignored.
func (r *Runner) Every(name string, interval time.Duration, fn func(context.Context) error) {
	if interval <= 0 || fn == nil {
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
		r.logger.Warn("ignoring periodic job registration after runner started", "job", name)
		return
	}

	r.jobs = append(r.jobs, job{
		name:     name,
		interval: interval,
		fn:       fn,
	})
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

	snapshot := make([]job, len(r.jobs))
	copy(snapshot, r.jobs)
	r.wg.Add(len(snapshot))
	r.mu.Unlock()

	for _, j := range snapshot {
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

func (r *Runner) loop(ctx context.Context, j job) {
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
			r.runOnce(ctx, j)
		}
	}
}

func (r *Runner) runOnce(ctx context.Context, j job) {
	defer func() {
		if rec := recover(); rec != nil {
			r.logger.Error("periodic job panicked",
				"job", j.name,
				"panic", rec,
				"stack", string(debug.Stack()),
			)
		}
	}()

	if err := j.fn(ctx); err != nil && !(ctx.Err() != nil && errors.Is(err, context.Canceled)) {
		r.logger.Error("periodic job failed",
			"job", j.name,
			"error", err,
		)
	}
}
