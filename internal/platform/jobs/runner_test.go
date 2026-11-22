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
