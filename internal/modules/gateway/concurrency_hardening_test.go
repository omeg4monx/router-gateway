package gateway

import (
	"context"
	"errors"
	"io"
	"net/http"
	"path/filepath"
	"strings"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"github.com/omeg4monx/router-gateway/internal/platform/db"
)

type roundTripFunc func(*http.Request) (*http.Response, error)

func (f roundTripFunc) RoundTrip(r *http.Request) (*http.Response, error) {
	return f(r)
}

func TestResponseCoalescerLeaderCancelDoesNotAbortWaiters(t *testing.T) {
	t.Parallel()

	coalescer := newResponseCoalescer(16, 30*time.Second, nil)
	leaderStarted := make(chan struct{})
	releaseUpstream := make(chan struct{})

	leaderCtx, cancelLeader := context.WithCancel(context.Background())
	defer cancelLeader()

	var (
		wg         sync.WaitGroup
		waiterResp cachedGatewayResponse
		waiterErr  error
		leaderErr  error
		calls      atomic.Int32
	)

	wg.Add(1)
	go func() {
		defer wg.Done()
		_, _, leaderErr = coalescer.do(leaderCtx, "k-leader-cancel", func(ctx context.Context) (cachedGatewayResponse, error) {
			calls.Add(1)
			close(leaderStarted)
			select {
			case <-releaseUpstream:
				return cachedGatewayResponse{
					statusCode: http.StatusOK,
					body:       []byte(`{"ok":true}`),
				}, nil
			case <-ctx.Done():
				return cachedGatewayResponse{}, ctx.Err()
			}
		})
	}()

	<-leaderStarted

	waiterReady := make(chan struct{})
	wg.Add(1)
	go func() {
		defer wg.Done()
		close(waiterReady)
		waiterResp, _, waiterErr = coalescer.do(context.Background(), "k-leader-cancel", func(context.Context) (cachedGatewayResponse, error) {
			calls.Add(1)
			return cachedGatewayResponse{}, errors.New("waiter should not execute upstream")
		})
	}()
	<-waiterReady
	time.Sleep(25 * time.Millisecond)

	cancelLeader()
	time.Sleep(25 * time.Millisecond)
	close(releaseUpstream)
	wg.Wait()

	if !errors.Is(leaderErr, context.Canceled) {
		t.Fatalf("leaderErr = %v, want context.Canceled", leaderErr)
	}
	if waiterErr != nil {
		t.Fatalf("waiterErr = %v, want nil (leader cancel must not abort waiters)", waiterErr)
	}
	if string(waiterResp.body) != `{"ok":true}` {
		t.Fatalf("waiterResp.body = %q, want %q", string(waiterResp.body), `{"ok":true}`)
	}
	if calls.Load() != 1 {
		t.Fatalf("upstream calls = %d, want 1", calls.Load())
	}
}

func TestResponseCoalescerPanicDoesNotDeadlockWaiters(t *testing.T) {
	t.Parallel()

	coalescer := newResponseCoalescer(16, 30*time.Second, nil)
	leaderStarted := make(chan struct{})
	releasePanic := make(chan struct{})

	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()
		defer func() { _ = recover() }()
		_, _, _ = coalescer.do(context.Background(), "k-panic", func(context.Context) (cachedGatewayResponse, error) {
			close(leaderStarted)
			<-releasePanic
			panic("boom")
		})
	}()

	<-leaderStarted

	waiterDone := make(chan error, 1)
	wg.Add(1)
	go func() {
		defer wg.Done()
		ctx, cancel := context.WithTimeout(context.Background(), 500*time.Millisecond)
		defer cancel()
		_, _, err := coalescer.do(ctx, "k-panic", func(context.Context) (cachedGatewayResponse, error) {
			return cachedGatewayResponse{statusCode: http.StatusOK, body: []byte(`{"ok":true}`)}, nil
		})
		waiterDone <- err
	}()

	time.Sleep(25 * time.Millisecond)
	close(releasePanic)

	select {
	case err := <-waiterDone:
		if errors.Is(err, context.DeadlineExceeded) {
			t.Fatal("waiter deadlocked until timeout after leader panic")
		}
		if err == nil {
			t.Fatal("expected non-nil error after leader panic")
		}
	case <-time.After(1 * time.Second):
		t.Fatal("timed out waiting for waiter to unblock after leader panic")
	}

	wg.Wait()

	// Subsequent request on the same key must succeed cleanly (no leaked inflight entry).
	resp, status, err := coalescer.do(context.Background(), "k-panic", func(context.Context) (cachedGatewayResponse, error) {
		return cachedGatewayResponse{statusCode: http.StatusOK, body: []byte(`{"recovered":true}`)}, nil
	})
	if err != nil {
		t.Fatalf("subsequent do err = %v, want nil", err)
	}
	if status != cacheStatusMiss || string(resp.body) != `{"recovered":true}` {
		t.Fatalf("subsequent do = (%s, %q), want (MISS, recovered)", status, string(resp.body))
	}
}

func TestAntigravityQuotaMonitorLeaderCancelDoesNotAbortWaiters(t *testing.T) {
	t.Parallel()

	monitor := newAntigravityQuotaMonitor()
	leaderStarted := make(chan struct{})
	releaseHTTP := make(chan struct{})

	client := &http.Client{
		Transport: roundTripFunc(func(req *http.Request) (*http.Response, error) {
			close(leaderStarted)
			select {
			case <-releaseHTTP:
				body := `{"groups":[{"displayName":"Gemini","buckets":[{"window":"5h","remainingFraction":0.0,"resetTime":"2099-01-01T00:00:00Z"}]}]}`
				return &http.Response{
					StatusCode: http.StatusOK,
					Header:     http.Header{"Content-Type": []string{"application/json"}},
					Body:       io.NopCloser(strings.NewReader(body)),
				}, nil
			case <-req.Context().Done():
				return nil, req.Context().Err()
			}
		}),
	}

	leaderCtx, cancelLeader := context.WithCancel(context.Background())
	defer cancelLeader()

	var (
		wg            sync.WaitGroup
		waiterDrained bool
	)

	wg.Add(1)
	go func() {
		defer wg.Done()
		_, _ = monitor.checkDrained(leaderCtx, client, "https://ag.example", "tok-1", "gemini-2.5-pro", "conn-1")
	}()

	<-leaderStarted

	wg.Add(1)
	go func() {
		defer wg.Done()
		waiterDrained, _ = monitor.checkDrained(context.Background(), client, "https://ag.example", "tok-1", "gemini-2.5-pro", "conn-1")
	}()

	time.Sleep(25 * time.Millisecond)
	cancelLeader()
	time.Sleep(25 * time.Millisecond)
	close(releaseHTTP)
	wg.Wait()

	if !waiterDrained {
		t.Fatal("expected waiter to observe drained=true even when leader canceled")
	}
}

func TestChainStoreThrottlesExpiredPrune(t *testing.T) {
	t.Parallel()

	database, err := db.Open(filepath.Join(t.TempDir(), "chains_throttle.db"))
	if err != nil {
		t.Fatalf("db.Open: %v", err)
	}
	t.Cleanup(func() { _ = database.Close() })

	mod := New(Options{DB: database})
	ctx := context.Background()
	if err := database.MigrateModule(ctx, mod.Name(), mod.Migrations()); err != nil {
		t.Fatalf("MigrateModule: %v", err)
	}
	store := newChainStore(database, time.Hour)

	t0 := time.Date(2026, 11, 28, 12, 0, 0, 0, time.UTC)
	expired := t0.Add(-5 * time.Minute).Format(rfc3339NanoFixed)
	created := t0.Add(-10 * time.Minute).Format(rfc3339NanoFixed)

	// Insert expired row #1; first saveChain must prune it.
	if _, err := database.Writer().ExecContext(ctx,
		`INSERT INTO response_chains (response_id, key_id, model, messages_json, expires_at, created_at) VALUES (?, ?, ?, ?, ?, ?)`,
		"exp-1", "key-1", "gpt-4o", `[]`, expired, created,
	); err != nil {
		t.Fatalf("seed exp-1: %v", err)
	}
	if err := store.saveChain(ctx, "live-1", "key-1", "gpt-4o", []CanonicalMessage{{Role: "user"}}, t0); err != nil {
		t.Fatalf("saveChain #1: %v", err)
	}
	var count int
	if err := database.Reader().QueryRowContext(ctx, `SELECT COUNT(*) FROM response_chains WHERE response_id = 'exp-1'`).Scan(&count); err != nil {
		t.Fatalf("count exp-1: %v", err)
	}
	if count != 0 {
		t.Fatalf("exp-1 count = %d, want 0 after first saveChain", count)
	}

	// Insert expired row #2; second saveChain 5s later must skip redundant DELETE.
	if _, err := database.Writer().ExecContext(ctx,
		`INSERT INTO response_chains (response_id, key_id, model, messages_json, expires_at, created_at) VALUES (?, ?, ?, ?, ?, ?)`,
		"exp-2", "key-1", "gpt-4o", `[]`, expired, created,
	); err != nil {
		t.Fatalf("seed exp-2: %v", err)
	}
	if err := store.saveChain(ctx, "live-2", "key-1", "gpt-4o", []CanonicalMessage{{Role: "user"}}, t0.Add(5*time.Second)); err != nil {
		t.Fatalf("saveChain #2: %v", err)
	}
	if err := database.Reader().QueryRowContext(ctx, `SELECT COUNT(*) FROM response_chains WHERE response_id = 'exp-2'`).Scan(&count); err != nil {
		t.Fatalf("count exp-2: %v", err)
	}
	if count != 1 {
		t.Fatalf("exp-2 count = %d, want 1 (prune should be throttled within pruneInterval)", count)
	}

	// Third saveChain >1m later must prune exp-2.
	if err := store.saveChain(ctx, "live-3", "key-1", "gpt-4o", []CanonicalMessage{{Role: "user"}}, t0.Add(2*time.Minute)); err != nil {
		t.Fatalf("saveChain #3: %v", err)
	}
	if err := database.Reader().QueryRowContext(ctx, `SELECT COUNT(*) FROM response_chains WHERE response_id = 'exp-2'`).Scan(&count); err != nil {
		t.Fatalf("count exp-2 after interval: %v", err)
	}
	if count != 0 {
		t.Fatalf("exp-2 count = %d, want 0 after pruneInterval elapsed", count)
	}
}
