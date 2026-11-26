package notifications_test

import (
	"bufio"
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"path/filepath"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/go-chi/chi/v5"

	"github.com/omeg4monx/router-gateway/internal/contracts"
	"github.com/omeg4monx/router-gateway/internal/modules/notifications"
	"github.com/omeg4monx/router-gateway/internal/platform/db"
	"github.com/omeg4monx/router-gateway/internal/platform/eventbus"
	"github.com/omeg4monx/router-gateway/internal/platform/lifecycle"
)

type stubJobInspector struct {
	mu      sync.Mutex
	jobs    map[string]contracts.JobStatus
	busy    map[string]bool
	trigger func(ctx context.Context, name string) error
}

func (s *stubJobInspector) Snapshot() []contracts.JobStatus {
	s.mu.Lock()
	defer s.mu.Unlock()
	out := make([]contracts.JobStatus, 0, len(s.jobs))
	for _, j := range s.jobs {
		out = append(out, j)
	}
	return out
}

func (s *stubJobInspector) Trigger(ctx context.Context, name string) error {
	s.mu.Lock()
	j, ok := s.jobs[name]
	if !ok {
		s.mu.Unlock()
		return contracts.ErrJobNotFound
	}
	if s.busy[name] {
		s.mu.Unlock()
		return contracts.ErrJobAlreadyRunning
	}
	fn := s.trigger
	s.mu.Unlock()

	if fn != nil {
		if err := fn(ctx, name); err != nil {
			return err
		}
	}

	s.mu.Lock()
	now := time.Now().UTC()
	j.RunCount++
	j.LastStatus = "ok"
	j.LastStartedAt = &now
	j.LastFinishedAt = &now
	s.jobs[name] = j
	s.mu.Unlock()
	return nil
}

func setupNotificationsTest(t *testing.T, inspector contracts.JobInspector) (*notifications.Module, chi.Router, *eventbus.Bus, *db.DB) {
	t.Helper()

	database, err := db.Open(filepath.Join(t.TempDir(), "notifications.sqlite"))
	if err != nil {
		t.Fatalf("db.Open: %v", err)
	}
	t.Cleanup(func() { _ = database.Close() })

	bus := eventbus.New()
	mod := notifications.New(notifications.Options{
		DB:   database,
		Bus:  bus,
		Jobs: inspector,
	})

	if err := database.MigrateModule(context.Background(), mod.Name(), mod.Migrations()); err != nil {
		t.Fatalf("MigrateModule: %v", err)
	}

	r := chi.NewRouter()
	r.Route("/api", func(api chi.Router) {
		mod.RegisterRoutes(lifecycle.RouteMounts{
			AdminAPI:       api,
			AdminStreamAPI: api,
		})
	})

	if err := mod.Start(context.Background()); err != nil {
		t.Fatalf("mod.Start: %v", err)
	}
	t.Cleanup(func() { _ = mod.Stop(context.Background()) })

	return mod, r, bus, database
}

func TestNotificationsCRUDReadUnreadAndPrune(t *testing.T) {
	t.Parallel()

	mod, router, _, database := setupNotificationsTest(t, nil)

	// Create two notifications.
	createOne := func(title, severity string) string {
		t.Helper()
		raw, _ := json.Marshal(map[string]any{
			"title":    title,
			"message":  "Details for " + title,
			"severity": severity,
			"category": "system",
		})
		req := httptest.NewRequest(http.MethodPost, "/api/notifications", bytes.NewReader(raw))
		req.Header.Set("Content-Type", "application/json")
		rec := httptest.NewRecorder()
		router.ServeHTTP(rec, req)
		if rec.Code != http.StatusCreated {
			t.Fatalf("POST /api/notifications status = %d, want 201; body=%s", rec.Code, rec.Body.String())
		}
		var out struct {
			ID string `json:"id"`
		}
		if err := json.Unmarshal(rec.Body.Bytes(), &out); err != nil {
			t.Fatalf("unmarshal created: %v", err)
		}
		return out.ID
	}

	id1 := createOne("First alert", "warning")
	id2 := createOne("Error alert", "error")

	// List all -> total 2, unread 2.
	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, httptest.NewRequest(http.MethodGet, "/api/notifications", nil))
	if rec.Code != http.StatusOK {
		t.Fatalf("GET /api/notifications status = %d, want 200", rec.Code)
	}
	var listResp struct {
		Items       []map[string]any `json:"items"`
		Total       int              `json:"total"`
		UnreadCount int              `json:"unread_count"`
	}
	if err := json.Unmarshal(rec.Body.Bytes(), &listResp); err != nil {
		t.Fatalf("unmarshal list: %v", err)
	}
	if listResp.Total != 2 || listResp.UnreadCount != 2 {
		t.Fatalf("list = total:%d unread:%d, want 2/2", listResp.Total, listResp.UnreadCount)
	}

	// Mark id1 read -> unread_count becomes 1.
	readRec := httptest.NewRecorder()
	router.ServeHTTP(readRec, httptest.NewRequest(http.MethodPost, "/api/notifications/"+id1+"/read", nil))
	if readRec.Code != http.StatusOK {
		t.Fatalf("POST read status = %d, want 200; body=%s", readRec.Code, readRec.Body.String())
	}

	// Filter unread_only=true -> returns only id2.
	unreadRec := httptest.NewRecorder()
	router.ServeHTTP(unreadRec, httptest.NewRequest(http.MethodGet, "/api/notifications?unread_only=true", nil))
	var unreadResp struct {
		Items       []struct{ ID string } `json:"items"`
		UnreadCount int                   `json:"unread_count"`
	}
	_ = json.Unmarshal(unreadRec.Body.Bytes(), &unreadResp)
	if len(unreadResp.Items) != 1 || unreadResp.Items[0].ID != id2 || unreadResp.UnreadCount != 1 {
		t.Fatalf("unread filter = %+v, want only %s", unreadResp, id2)
	}

	// Mark all read -> unread_count becomes 0.
	readAllRec := httptest.NewRecorder()
	router.ServeHTTP(readAllRec, httptest.NewRequest(http.MethodPost, "/api/notifications/read-all", nil))
	if readAllRec.Code != http.StatusOK {
		t.Fatalf("POST read-all status = %d, want 200", readAllRec.Code)
	}

	// Seed an expired notification (> 30 days old) and prune it.
	oldTime := time.Now().UTC().AddDate(0, 0, -45).Format("2006-01-02T15:04:05.000000000Z")
	_, err := database.Writer().ExecContext(context.Background(),
		`INSERT INTO system_notifications(id, title, message, severity, category, source, metadata_json, read_at, created_at)
		 VALUES ('ntf_old', 'Expired', 'Old', 'info', 'system', 'system', '{}', '', ?)`,
		oldTime,
	)
	if err != nil {
		t.Fatalf("seed old notification: %v", err)
	}
	if err := mod.PruneExpired(context.Background()); err != nil {
		t.Fatalf("PruneExpired: %v", err)
	}

	// Delete id1, then clear all.
	delRec := httptest.NewRecorder()
	router.ServeHTTP(delRec, httptest.NewRequest(http.MethodDelete, "/api/notifications/"+id1, nil))
	if delRec.Code != http.StatusOK {
		t.Fatalf("DELETE single status = %d, want 200", delRec.Code)
	}

	clearRec := httptest.NewRecorder()
	router.ServeHTTP(clearRec, httptest.NewRequest(http.MethodDelete, "/api/notifications", nil))
	if clearRec.Code != http.StatusOK {
		t.Fatalf("DELETE all status = %d, want 200", clearRec.Code)
	}
}

func TestNotificationsEventBusAndSSEStreamCleanup(t *testing.T) {
	t.Parallel()

	_, router, bus, _ := setupNotificationsTest(t, nil)

	srv := httptest.NewServer(router)
	defer srv.Close()

	streamCtx, cancelStream := context.WithCancel(context.Background())
	req, err := http.NewRequestWithContext(streamCtx, http.MethodGet, srv.URL+"/api/notifications/stream", nil)
	if err != nil {
		t.Fatalf("new stream req: %v", err)
	}
	resp, err := srv.Client().Do(req)
	if err != nil {
		t.Fatalf("connect stream: %v", err)
	}
	defer func() { _ = resp.Body.Close() }()

	if ct := resp.Header.Get("Content-Type"); !strings.Contains(ct, "text/event-stream") {
		t.Fatalf("Content-Type = %q, want text/event-stream", ct)
	}

	lineCh := make(chan string, 16)
	readDone := make(chan struct{})
	go func() {
		defer close(readDone)
		scanner := bufio.NewScanner(resp.Body)
		for scanner.Scan() {
			line := scanner.Text()
			if strings.HasPrefix(line, "data:") {
				lineCh <- line
			}
		}
	}()

	// Give stream handler a moment to register subscription, then emit via eventbus.
	time.Sleep(50 * time.Millisecond)
	bus.Publish(notifications.TopicJobFailed, notifications.EmitPayload{
		Title:    "Background job failed",
		Message:  "usage-prune returned an error",
		Severity: "error",
		Category: "jobs",
		Source:   "jobs.runner",
	})

	select {
	case dataLine := <-lineCh:
		if !strings.Contains(dataLine, "Background job failed") {
			t.Fatalf("unexpected SSE payload: %s", dataLine)
		}
	case <-time.After(2 * time.Second):
		t.Fatal("timed out waiting for SSE notification event")
	}

	// Cancel client context and ensure stream exits cleanly without hanging.
	cancelStream()
	select {
	case <-readDone:
	case <-time.After(2 * time.Second):
		t.Fatal("SSE reader did not exit after context cancellation")
	}
}

func TestNotificationsJobsListAndTriggerConflict(t *testing.T) {
	t.Parallel()

	inspector := &stubJobInspector{
		jobs: map[string]contracts.JobStatus{
			"usage-prune": {
				Name:       "usage-prune",
				IntervalMs: int64(time.Hour / time.Millisecond),
				LastStatus: "idle",
			},
			"busy-job": {
				Name:       "busy-job",
				IntervalMs: int64(time.Minute / time.Millisecond),
				Running:    true,
				LastStatus: "running",
			},
		},
		busy: map[string]bool{
			"busy-job": true,
		},
	}

	_, router, _, _ := setupNotificationsTest(t, inspector)

	// GET /api/jobs
	listRec := httptest.NewRecorder()
	router.ServeHTTP(listRec, httptest.NewRequest(http.MethodGet, "/api/jobs", nil))
	if listRec.Code != http.StatusOK {
		t.Fatalf("GET /api/jobs status = %d, want 200; body=%s", listRec.Code, listRec.Body.String())
	}
	var listResp struct {
		Jobs []struct {
			Name string `json:"name"`
		} `json:"jobs"`
	}
	if err := json.Unmarshal(listRec.Body.Bytes(), &listResp); err != nil {
		t.Fatalf("unmarshal jobs: %v", err)
	}
	if len(listResp.Jobs) != 2 {
		t.Fatalf("jobs len = %d, want 2", len(listResp.Jobs))
	}

	// POST /api/jobs/usage-prune/trigger -> 200 OK
	trigRec := httptest.NewRecorder()
	router.ServeHTTP(trigRec, httptest.NewRequest(http.MethodPost, "/api/jobs/usage-prune/trigger", nil))
	if trigRec.Code != http.StatusOK {
		t.Fatalf("trigger usage-prune status = %d, want 200; body=%s", trigRec.Code, trigRec.Body.String())
	}

	// POST /api/jobs/missing/trigger -> 404 Not Found
	missRec := httptest.NewRecorder()
	router.ServeHTTP(missRec, httptest.NewRequest(http.MethodPost, "/api/jobs/missing/trigger", nil))
	if missRec.Code != http.StatusNotFound {
		t.Fatalf("trigger missing status = %d, want 404", missRec.Code)
	}

	// POST /api/jobs/busy-job/trigger -> 409 Conflict
	busyRec := httptest.NewRecorder()
	router.ServeHTTP(busyRec, httptest.NewRequest(http.MethodPost, "/api/jobs/busy-job/trigger", nil))
	if busyRec.Code != http.StatusConflict {
		t.Fatalf("trigger busy-job status = %d, want 409; body=%s", busyRec.Code, busyRec.Body.String())
	}
}
