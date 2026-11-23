package jobs_test

import (
	"bytes"
	"context"
	"errors"
	"log/slog"
	"strings"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"github.com/omeg4monx/router-gateway/internal/platform/jobs"
)

type syncBuffer struct {
	mu  sync.Mutex
	buf bytes.Buffer
}

func (b *syncBuffer) Write(p []byte) (int, error) {
	b.mu.Lock()
	defer b.mu.Unlock()
	return b.buf.Write(p)
}

func (b *syncBuffer) String() string {
	b.mu.Lock()
	defer b.mu.Unlock()
	return b.buf.String()
}

func TestRunnerRecoversFromPanicAndStopsCleanly(t *testing.T) {
	t.Parallel()

	var logBuf syncBuffer
	logger := slog.New(slog.NewTextHandler(&logBuf, &slog.HandlerOptions{Level: slog.LevelDebug}))
	runner := jobs.NewRunner(logger)

	var runs atomic.Int32
	secondRunDone := make(chan struct{})

	runner.Every("flaky-job", 10*time.Millisecond, func(_ context.Context) error {
		n := runs.Add(1)
		if n == 1 {
			panic("boom on run 1")
		}
		if n == 2 {
			close(secondRunDone)
		}
		return nil
	})

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	runner.Start(ctx)

	select {
	case <-secondRunDone:
	case <-time.After(2 * time.Second):
		t.Fatalf("timed out waiting for job run 2 after panic; runs=%d", runs.Load())
	}

	stopCtx, stopCancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer stopCancel()

	if err := runner.Stop(stopCtx); err != nil {
		t.Fatalf("Stop returned unexpected error: %v", err)
	}

	logs := logBuf.String()
	if !strings.Contains(logs, "flaky-job") || !strings.Contains(logs, "boom on run 1") {
		t.Fatalf("expected panic log to include job name and panic value, got: %s", logs)
	}
	if !strings.Contains(logs, "stack=") {
		t.Fatalf("expected panic log to include stack trace, got: %s", logs)
	}
}

func TestRunnerIdempotentStartStopAndGuards(t *testing.T) {
	t.Parallel()

	// Nil logger should not panic.
	runner := jobs.NewRunner(nil)

	// Invalid intervals or nil fn should be ignored safely.
	runner.Every("zero-interval", 0, func(context.Context) error { return nil })
	runner.Every("negative-interval", -time.Second, func(context.Context) error { return nil })
	runner.Every("nil-fn", 10*time.Millisecond, nil)

	// Stop before Start is a safe no-op.
	if err := runner.Stop(context.Background()); err != nil {
		t.Fatalf("Stop before Start: %v", err)
	}

	var count atomic.Int32
	ran := make(chan struct{})
	var once sync.Once

	runner.Every("healthy-job", 10*time.Millisecond, func(context.Context) error {
		count.Add(1)
		once.Do(func() { close(ran) })
		return errors.New("non-fatal error")
	})

	ctx := context.Background()
	runner.Start(ctx)
	// Second Start must be idempotent.
	runner.Start(ctx)

	select {
	case <-ran:
	case <-time.After(2 * time.Second):
		t.Fatal("timed out waiting for healthy-job to run")
	}

	stopCtx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	if err := runner.Stop(stopCtx); err != nil {
		t.Fatalf("first Stop: %v", err)
	}
	// Second Stop must be idempotent.
	if err := runner.Stop(stopCtx); err != nil {
		t.Fatalf("second Stop: %v", err)
	}

	if count.Load() < 1 {
		t.Fatalf("expected at least 1 run, got %d", count.Load())
	}
}

func TestRunnerStopRespectsContextDeadline(t *testing.T) {
	t.Parallel()

	runner := jobs.NewRunner(slog.New(slog.NewTextHandler(&syncBuffer{}, nil)))

	started := make(chan struct{})
	release := make(chan struct{})
	var once sync.Once

	runner.Every("blocking-job", 5*time.Millisecond, func(_ context.Context) error {
		once.Do(func() { close(started) })
		<-release
		return nil
	})

	runner.Start(context.Background())

	select {
	case <-started:
	case <-time.After(2 * time.Second):
		t.Fatal("timed out waiting for blocking-job to start")
	}

	shortCtx, shortCancel := context.WithTimeout(context.Background(), 25*time.Millisecond)
	defer shortCancel()

	err := runner.Stop(shortCtx)
	if !errors.Is(err, context.DeadlineExceeded) {
		close(release)
		t.Fatalf("Stop err = %v, want context.DeadlineExceeded", err)
	}

	close(release)

	cleanCtx, cleanCancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cleanCancel()
	if err := runner.Stop(cleanCtx); err != nil {
		t.Fatalf("Stop after worker unblocked: %v", err)
	}
}

func TestRunnerSnapshotAndManualTrigger(t *testing.T) {
	t.Parallel()

	runner := jobs.NewRunner(slog.New(slog.NewTextHandler(&syncBuffer{}, nil)))

	var betaRuns atomic.Int32
	runner.Every("beta-job", time.Hour, func(_ context.Context) error {
		betaRuns.Add(1)
		return nil
	})
	runner.Every("alpha-job", 30*time.Minute, func(_ context.Context) error {
		return errors.New("alpha failed")
	})

	snap := runner.Snapshot()
	if len(snap) != 2 {
		t.Fatalf("Snapshot len = %d, want 2", len(snap))
	}
	if snap[0].Name != "alpha-job" || snap[1].Name != "beta-job" {
		t.Fatalf("Snapshot order = [%s, %s], want [alpha-job, beta-job]", snap[0].Name, snap[1].Name)
	}
	if snap[0].LastStatus != "idle" || snap[1].LastStatus != "idle" {
		t.Fatalf("initial LastStatus = [%s, %s], want idle", snap[0].LastStatus, snap[1].LastStatus)
	}

	if err := runner.Trigger(context.Background(), "missing-job"); !errors.Is(err, jobs.ErrJobNotFound) {
		t.Fatalf("Trigger missing err = %v, want ErrJobNotFound", err)
	}

	if err := runner.Trigger(context.Background(), "beta-job"); err != nil {
		t.Fatalf("Trigger beta-job err = %v, want nil", err)
	}
	if betaRuns.Load() != 1 {
		t.Fatalf("betaRuns = %d, want 1", betaRuns.Load())
	}

	if err := runner.Trigger(context.Background(), "alpha-job"); err == nil || err.Error() != "alpha failed" {
		t.Fatalf("Trigger alpha-job err = %v, want alpha failed", err)
	}

	after := runner.Snapshot()
	if after[0].RunCount != 1 || after[0].ErrorCount != 1 || after[0].LastStatus != "error" || after[0].LastError != "alpha failed" {
		t.Fatalf("alpha status = %+v, want RunCount=1 ErrorCount=1 LastStatus=error", after[0])
	}
	if after[0].LastStartedAt == nil || after[0].LastFinishedAt == nil {
		t.Fatalf("alpha timestamps missing: %+v", after[0])
	}
	if after[1].RunCount != 1 || after[1].ErrorCount != 0 || after[1].LastStatus != "ok" || after[1].LastError != "" {
		t.Fatalf("beta status = %+v, want RunCount=1 ErrorCount=0 LastStatus=ok", after[1])
	}
}

func TestRunnerOverlapGuardReturnsAlreadyRunning(t *testing.T) {
	t.Parallel()

	runner := jobs.NewRunner(slog.New(slog.NewTextHandler(&syncBuffer{}, nil)))

	started := make(chan struct{})
	release := make(chan struct{})
	runner.Every("slow-job", time.Hour, func(_ context.Context) error {
		close(started)
		<-release
		return nil
	})

	done := make(chan error, 1)
	go func() {
		done <- runner.Trigger(context.Background(), "slow-job")
	}()

	select {
	case <-started:
	case <-time.After(2 * time.Second):
		t.Fatal("timed out waiting for slow-job start")
	}

	if err := runner.Trigger(context.Background(), "slow-job"); !errors.Is(err, jobs.ErrJobAlreadyRunning) {
		close(release)
		t.Fatalf("overlapping Trigger err = %v, want ErrJobAlreadyRunning", err)
	}

	close(release)
	if err := <-done; err != nil {
		t.Fatalf("first Trigger err = %v, want nil", err)
	}
}

func TestRunnerPanicRecoveryAndObserverCallback(t *testing.T) {
	t.Parallel()

	runner := jobs.NewRunner(slog.New(slog.NewTextHandler(&syncBuffer{}, nil)))

	var mu sync.Mutex
	var observed []jobs.JobStatus
	runner.SetObserver(func(st jobs.JobStatus) {
		mu.Lock()
		observed = append(observed, st)
		mu.Unlock()
	})

	runner.Every("panic-job", time.Hour, func(_ context.Context) error {
		panic("kaboom")
	})

	err := runner.Trigger(context.Background(), "panic-job")
	if err == nil || !strings.Contains(err.Error(), "kaboom") {
		t.Fatalf("Trigger panic-job err = %v, want panic error containing kaboom", err)
	}

	snap := runner.Snapshot()
	if len(snap) != 1 || snap[0].LastStatus != "panic" || snap[0].ErrorCount != 1 || snap[0].RunCount != 1 {
		t.Fatalf("panic-job snapshot = %+v, want LastStatus=panic ErrorCount=1 RunCount=1", snap)
	}

	mu.Lock()
	defer mu.Unlock()
	if len(observed) != 1 || observed[0].LastStatus != "panic" || !strings.Contains(observed[0].LastError, "kaboom") {
		t.Fatalf("observed = %+v, want 1 panic event", observed)
	}
}
