package translator

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/fs"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/omeg4monx/router-gateway/internal/contracts"
	"github.com/omeg4monx/router-gateway/internal/platform/httpserver"
	"github.com/omeg4monx/router-gateway/internal/platform/lifecycle"
	"github.com/omeg4monx/router-gateway/internal/platform/logging"
)

const (
	moduleName         = "translator"
	maxBodyBytes       = 4 << 20 // 4 MiB
	defaultLogsLimit   = 100
	maxLogsLimit       = 500
	defaultStreamTail  = 50
	sseKeepalivePeriod = 15 * time.Second
)

// Options configures a translator Module instance.
type Options struct {
	Translator contracts.ProtocolTranslator
	LogRing    *logging.RingBuffer
}

// Module implements lifecycle.Module and provides interactive wire-protocol
// translation plus live structured console log streaming.
type Module struct {
	translator contracts.ProtocolTranslator
	logRing    *logging.RingBuffer
}

var _ lifecycle.Module = (*Module)(nil)

// New constructs a translator Module.
func New(opts Options) *Module {
	return &Module{
		translator: opts.Translator,
		logRing:    opts.LogRing,
	}
}

// Name returns the canonical module name.
func (m *Module) Name() string { return moduleName }

// Dependencies declares the upstream modules required by translator.
func (m *Module) Dependencies() []string { return []string{"system", "gateway"} }

// Migrations returns nil because translator is stateless.
func (m *Module) Migrations() fs.FS { return nil }

// RegisterRoutes mounts protocol translator and console log routes onto mounts.AdminAPI
// (and mounts the long-lived SSE log stream onto mounts.AdminStreamAPI when available).
func (m *Module) RegisterRoutes(mounts lifecycle.RouteMounts) {
	if mounts.AdminAPI != nil {
		mounts.AdminAPI.Get("/translator/protocols", m.handleProtocols)
		mounts.AdminAPI.Post("/translator/translate", m.handleTranslate)
		mounts.AdminAPI.Get("/translator/logs", m.handleLogs)
		mounts.AdminAPI.Delete("/translator/logs", m.handleClearLogs)
	}
	streamMount := mounts.AdminStreamAPI
	if streamMount == nil {
		streamMount = mounts.AdminAPI
	}
	if streamMount != nil {
		streamMount.Get("/translator/logs/stream", m.handleLogsStream)
	}
}

// Start is a no-op for the stateless translator module.
func (m *Module) Start(_ context.Context) error { return nil }

// Stop is a no-op for the stateless translator module.
func (m *Module) Stop(_ context.Context) error { return nil }

func (m *Module) handleProtocols(w http.ResponseWriter, _ *http.Request) {
	protos := []string{}
	if m.translator != nil {
		protos = m.translator.SupportedProtocols()
	}
	httpserver.WriteJSON(w, http.StatusOK, map[string]any{
		"protocols": protos,
		"modes":     []string{"request", "response"},
	})
}

func (m *Module) handleTranslate(w http.ResponseWriter, r *http.Request) {
	if m.translator == nil {
		httpserver.WriteError(w, r, http.StatusServiceUnavailable, "translator_unavailable", "protocol translator is not configured")
		return
	}

	var req struct {
		Mode         string          `json:"mode"`
		FromProtocol string          `json:"from_protocol"`
		ToProtocol   string          `json:"to_protocol"`
		Payload      json.RawMessage `json:"payload"`
	}
	if err := decodeStrictJSON(w, r, &req); err != nil {
		httpserver.WriteError(w, r, http.StatusBadRequest, "invalid_request", err.Error())
		return
	}

	mode := strings.ToLower(strings.TrimSpace(req.Mode))
	if mode != "request" && mode != "response" {
		httpserver.WriteError(w, r, http.StatusBadRequest, "invalid_mode", "mode must be 'request' or 'response'")
		return
	}

	fromProto := strings.TrimSpace(req.FromProtocol)
	toProto := strings.TrimSpace(req.ToProtocol)
	payload := bytes.TrimSpace(req.Payload)
	if len(payload) == 0 || payload[0] != '{' {
		httpserver.WriteError(w, r, http.StatusBadRequest, "invalid_payload", "payload must be a JSON object")
		return
	}

	start := time.Now()
	var (
		out []byte
		err error
	)
	if mode == "request" {
		out, err = m.translator.TranslateRequest(fromProto, toProto, payload)
	} else {
		out, err = m.translator.TranslateResponse(fromProto, toProto, payload)
	}
	if err != nil {
		httpserver.WriteError(w, r, http.StatusBadRequest, "translation_failed", err.Error())
		return
	}

	httpserver.WriteJSON(w, http.StatusOK, map[string]any{
		"mode":          mode,
		"from_protocol": fromProto,
		"to_protocol":   toProto,
		"fast_path":     fromProto == toProto,
		"duration_us":   time.Since(start).Microseconds(),
		"output":        json.RawMessage(out),
	})
}

func (m *Module) handleLogs(w http.ResponseWriter, r *http.Request) {
	limit := parseIntParam(r.URL.Query().Get("limit"), defaultLogsLimit, 1, maxLogsLimit)
	minLevel := normalizeLevelFilter(r.URL.Query().Get("level"))

	entries := make([]logging.LogEntry, 0)
	if m.logRing != nil {
		raw := m.logRing.Recent(maxLogsLimit)
		for _, entry := range raw {
			if matchesLevel(entry.Level, minLevel) {
				entries = append(entries, entry)
			}
		}
		if len(entries) > limit {
			entries = entries[len(entries)-limit:]
		}
	}

	httpserver.WriteJSON(w, http.StatusOK, map[string]any{
		"entries": entries,
		"total":   len(entries),
	})
}

func (m *Module) handleClearLogs(w http.ResponseWriter, _ *http.Request) {
	if m.logRing != nil {
		m.logRing.Clear()
	}
	httpserver.WriteJSON(w, http.StatusOK, map[string]any{
		"cleared": true,
	})
}

func (m *Module) handleLogsStream(w http.ResponseWriter, r *http.Request) {
	tail := parseIntParam(r.URL.Query().Get("tail"), defaultStreamTail, 0, maxLogsLimit)
	minLevel := normalizeLevelFilter(r.URL.Query().Get("level"))

	rc := http.NewResponseController(w)
	_ = rc.SetWriteDeadline(time.Time{})

	w.Header().Set("Content-Type", "text/event-stream; charset=utf-8")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")
	w.Header().Set("X-Accel-Buffering", "no")
	w.WriteHeader(http.StatusOK)

	writeEntry := func(entry logging.LogEntry) bool {
		if !matchesLevel(entry.Level, minLevel) {
			return true
		}
		data, err := json.Marshal(entry)
		if err != nil {
			return true
		}
		if _, err := fmt.Fprintf(w, "event: log\ndata: %s\n\n", data); err != nil {
			return false
		}
		_ = rc.Flush()
		return true
	}

	var (
		subCh      <-chan logging.LogEntry
		unsub      func()
		clearCh    <-chan struct{}
		unsubClear func()
	)
	if m.logRing != nil {
		subCh, unsub = m.logRing.Subscribe()
		defer unsub()
		clearCh, unsubClear = m.logRing.SubscribeClear()
		defer unsubClear()

		if tail > 0 {
			recent := m.logRing.Recent(tail)
			for _, entry := range recent {
				if !writeEntry(entry) {
					return
				}
			}
		}
	}
	_ = rc.Flush()

	ticker := time.NewTicker(sseKeepalivePeriod)
	defer ticker.Stop()

	for {
		select {
		case <-r.Context().Done():
			return
		case entry, ok := <-subCh:
			if !ok {
				return
			}
			if !writeEntry(entry) {
				return
			}
		case _, ok := <-clearCh:
			if !ok {
				return
			}
			if _, err := fmt.Fprint(w, "event: clear\ndata: {}\n\n"); err != nil {
				return
			}
			_ = rc.Flush()
		case <-ticker.C:
			if _, err := fmt.Fprint(w, ": keepalive\n\n"); err != nil {
				return
			}
			_ = rc.Flush()
		}
	}
}

func decodeStrictJSON(w http.ResponseWriter, r *http.Request, dst any) error {
	if r.Body == nil {
		return errors.New("request body is empty")
	}
	defer r.Body.Close()
	dec := json.NewDecoder(http.MaxBytesReader(w, r.Body, maxBodyBytes))
	dec.DisallowUnknownFields()
	if err := dec.Decode(dst); err != nil {
		return fmt.Errorf("malformed JSON: %w", err)
	}
	if err := dec.Decode(&struct{}{}); !errors.Is(err, io.EOF) {
		return errors.New("unexpected trailing JSON content")
	}
	return nil
}

func parseIntParam(raw string, def, minVal, maxVal int) int {
	if strings.TrimSpace(raw) == "" {
		return def
	}
	v, err := strconv.Atoi(strings.TrimSpace(raw))
	if err != nil {
		return def
	}
	if v < minVal {
		return minVal
	}
	if v > maxVal {
		return maxVal
	}
	return v
}

func normalizeLevelFilter(raw string) int {
	switch strings.ToUpper(strings.TrimSpace(raw)) {
	case "DEBUG":
		return 1
	case "INFO":
		return 2
	case "WARN", "WARNING":
		return 3
	case "ERROR":
		return 4
	default:
		return 0
	}
}

func matchesLevel(entryLevel string, minSeverity int) bool {
	if minSeverity == 0 {
		return true
	}
	return normalizeLevelFilter(entryLevel) >= minSeverity
}
