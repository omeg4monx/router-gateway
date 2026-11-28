package system_test

import (
	"bytes"
	"context"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"sync/atomic"
	"testing"
	"time"

	"github.com/omeg4monx/router-gateway/internal/modules/system"
	"github.com/omeg4monx/router-gateway/internal/platform/config"
	"github.com/omeg4monx/router-gateway/internal/platform/db"
	"github.com/omeg4monx/router-gateway/internal/platform/httpserver"
	"github.com/omeg4monx/router-gateway/internal/platform/lifecycle"
)

type updateStatusPayload struct {
	CurrentVersion    string `json:"currentVersion"`
	CurrentCommit     string `json:"currentCommit"`
	Branch            string `json:"branch"`
	LatestCommit      string `json:"latestCommit"`
	LatestShortCommit string `json:"latestShortCommit"`
	LatestMessage     string `json:"latestMessage"`
	LatestPublishedAt string `json:"latestPublishedAt"`
	LatestURL         string `json:"latestUrl"`
	UpdateAvailable   bool   `json:"updateAvailable"`
	CanApply          bool   `json:"canApply"`
	Mode              string `json:"mode"`
	UnavailableReason string `json:"unavailableReason"`
	State             string `json:"state"`
	Message           string `json:"message"`
	Error             string `json:"error"`
	CheckedAt         string `json:"checkedAt"`
	UpdatedAt         string `json:"updatedAt"`
}

func newUpdaterTestHarness(t *testing.T, opts system.Options, auth ...*stubAdminAuth) *testHarness {
	t.Helper()

	if opts.Config.Host == "" {
		opts.Config.Host = "127.0.0.1"
	}
	if opts.Config.Port == 0 {
		opts.Config.Port = 20128
	}
	if opts.Config.APITimeout == 0 {
		opts.Config.APITimeout = 30 * time.Second
	}
	if opts.Config.UpdateRepo == "" {
		opts.Config.UpdateRepo = "omeg4monx/router-gateway"
	}
	if opts.Config.UpdateBranch == "" {
		opts.Config.UpdateBranch = "main"
	}

	dbPath := filepath.Join(t.TempDir(), "router.sqlite")
	database, err := db.Open(dbPath)
	if err != nil {
		t.Fatalf("db.Open: %v", err)
	}
	t.Cleanup(func() {
		_ = database.Close()
	})

	reg := lifecycle.NewRegistry(nil)
	opts.DB = database
	if opts.ActiveModules == nil {
		opts.ActiveModules = reg.ActiveNames
	}

	mod := system.New(opts)
	if err := reg.Register(mod); err != nil {
		t.Fatalf("reg.Register: %v", err)
	}
	if err := reg.MigrateAll(context.Background(), database); err != nil {
		t.Fatalf("reg.MigrateAll: %v", err)
	}

	srvOpts := httpserver.Options{
		Config: opts.Config,
		Logger: slog.New(slog.NewTextHandler(io.Discard, nil)),
	}
	if len(auth) > 0 && auth[0] != nil {
		srvOpts.Auth = auth[0]
	}
	srv := httpserver.New(srvOpts)
	if err := reg.MountAll(srv.Mounts()); err != nil {
		t.Fatalf("reg.MountAll: %v", err)
	}
	if err := reg.StartAll(context.Background()); err != nil {
		t.Fatalf("reg.StartAll: %v", err)
	}
	t.Cleanup(func() {
		_ = reg.StopAll(context.Background())
	})

	return &testHarness{
		db:     database,
		reg:    reg,
		server: srv,
		mod:    mod,
	}
}

func expectedPlatformAssetName() string {
	ext := ""
	if runtime.GOOS == "windows" {
		ext = ".exe"
	}
	return fmt.Sprintf("omeg4mon-router-%s-%s%s", runtime.GOOS, runtime.GOARCH, ext)
}

func makeFakeBinaryBytes() []byte {
	switch runtime.GOOS {
	case "windows":
		return append([]byte("MZ"), bytes.Repeat([]byte("A"), 256)...)
	case "darwin":
		return append([]byte{0xcf, 0xfa, 0xed, 0xfe}, bytes.Repeat([]byte("A"), 256)...)
	default:
		return append([]byte{0x7f, 'E', 'L', 'F'}, bytes.Repeat([]byte("A"), 256)...)
	}
}

func TestSystemUpdateCheckUpToDateAndAvailable(t *testing.T) {
	const fullSHA = "1234567890abcdef1234567890abcdef12345678"
	assetName := expectedPlatformAssetName()

	var ghServer *httptest.Server
	ghServer = httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case "/repos/omeg4monx/router-gateway/commits/main":
			_ = json.NewEncoder(w).Encode(map[string]any{
				"sha":      fullSHA,
				"html_url": "https://github.com/omeg4monx/router-gateway/commit/" + fullSHA,
				"commit": map[string]any{
					"message": "feat: new routing improvements\n\nmore details",
					"committer": map[string]any{
						"date": "2026-11-28T10:00:00Z",
					},
				},
			})
		case "/repos/omeg4monx/router-gateway/releases/tags/main-latest":
			_ = json.NewEncoder(w).Encode(map[string]any{
				"tag_name":         "main-latest",
				"target_commitish": fullSHA,
				"assets": []map[string]any{
					{
						"name":                 assetName,
						"browser_download_url": ghServer.URL + "/download/" + assetName,
					},
					{
						"name":                 "SHA256SUMS.txt",
						"browser_download_url": ghServer.URL + "/download/SHA256SUMS.txt",
					},
				},
			})
		default:
			http.NotFound(w, r)
		}
	}))
	defer ghServer.Close()

	standaloneExe := filepath.Join(t.TempDir(), "omeg4mon-router")
	if err := os.WriteFile(standaloneExe, makeFakeBinaryBytes(), 0o755); err != nil {
		t.Fatalf("write exe: %v", err)
	}

	t.Run("up to date when short commit matches full SHA prefix", func(t *testing.T) {
		h := newUpdaterTestHarness(t, system.Options{
			Build:            system.BuildInfo{Version: "v0.3.0", Commit: "1234567"},
			UpdateAPIBaseURL: ghServer.URL,
			ExecutablePath:   func() (string, error) { return standaloneExe, nil },
		})

		rec := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodGet, "/api/system/update", nil)
		h.server.Handler().ServeHTTP(rec, req)

		if rec.Code != http.StatusOK {
			t.Fatalf("status = %d, want 200 (body: %s)", rec.Code, rec.Body.String())
		}
		var payload updateStatusPayload
		if err := json.NewDecoder(rec.Body).Decode(&payload); err != nil {
			t.Fatalf("decode: %v", err)
		}
		if payload.UpdateAvailable {
			t.Fatalf("UpdateAvailable = true, want false")
		}
		if payload.LatestShortCommit != "1234567" {
			t.Fatalf("LatestShortCommit = %q, want 1234567", payload.LatestShortCommit)
		}
		if payload.LatestMessage != "feat: new routing improvements" {
			t.Fatalf("LatestMessage = %q, want first line only", payload.LatestMessage)
		}
	})

	t.Run("update available in binary mode when commit differs", func(t *testing.T) {
		h := newUpdaterTestHarness(t, system.Options{
			Build:            system.BuildInfo{Version: "v0.3.0", Commit: "aaaaaaa"},
			UpdateAPIBaseURL: ghServer.URL,
			ExecutablePath:   func() (string, error) { return standaloneExe, nil },
		})

		rec := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodGet, "/api/system/update", nil)
		h.server.Handler().ServeHTTP(rec, req)

		if rec.Code != http.StatusOK {
			t.Fatalf("status = %d, want 200", rec.Code)
		}
		var payload updateStatusPayload
		if err := json.NewDecoder(rec.Body).Decode(&payload); err != nil {
			t.Fatalf("decode: %v", err)
		}
		if !payload.UpdateAvailable {
			t.Fatal("UpdateAvailable = false, want true")
		}
		if !payload.CanApply {
			t.Fatalf("CanApply = false, want true (reason: %q)", payload.UnavailableReason)
		}
		if payload.Mode != "binary" {
			t.Fatalf("Mode = %q, want binary", payload.Mode)
		}
	})
}

func TestSystemUpdateBinaryApplyVerifiesSHA256SwapsAndRestarts(t *testing.T) {
	const fullSHA = "fedcba9876543210fedcba9876543210fedcba98"
	assetName := expectedPlatformAssetName()
	newBin := makeFakeBinaryBytes()
	newBin = append(newBin, []byte("NEW_VERSION_MARKER")...)
	sum := sha256.Sum256(newBin)
	sumHex := hex.EncodeToString(sum[:])

	var badChecksum atomic.Bool
	var ghServer *httptest.Server
	ghServer = httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case "/repos/omeg4monx/router-gateway/commits/main":
			_ = json.NewEncoder(w).Encode(map[string]any{
				"sha": fullSHA,
				"commit": map[string]any{
					"message":   "fix: binary update",
					"committer": map[string]any{"date": "2026-11-28T10:00:00Z"},
				},
			})
		case "/repos/omeg4monx/router-gateway/releases/tags/main-latest":
			_ = json.NewEncoder(w).Encode(map[string]any{
				"tag_name":         "main-latest",
				"target_commitish": fullSHA,
				"assets": []map[string]any{
					{"name": assetName, "browser_download_url": ghServer.URL + "/download/" + assetName},
					{"name": "SHA256SUMS.txt", "browser_download_url": ghServer.URL + "/download/SHA256SUMS.txt"},
				},
			})
		case "/download/" + assetName:
			_, _ = w.Write(newBin)
		case "/download/SHA256SUMS.txt":
			digest := sumHex
			if badChecksum.Load() {
				digest = strings.Repeat("0", 64)
			}
			_, _ = fmt.Fprintf(w, "%s  %s\n", digest, assetName)
		default:
			http.NotFound(w, r)
		}
	}))
	defer ghServer.Close()

	t.Run("checksum mismatch fails cleanly without touching original binary", func(t *testing.T) {
		badChecksum.Store(true)
		exeDir := t.TempDir()
		exePath := filepath.Join(exeDir, "omeg4mon-router")
		origBytes := []byte("ORIGINAL_BINARY")
		if err := os.WriteFile(exePath, origBytes, 0o755); err != nil {
			t.Fatalf("write exe: %v", err)
		}

		restarted := make(chan string, 1)
		h := newUpdaterTestHarness(t, system.Options{
			Build:            system.BuildInfo{Version: "v0.3.0", Commit: "1111111"},
			UpdateAPIBaseURL: ghServer.URL,
			ExecutablePath:   func() (string, error) { return exePath, nil },
			Restart:          func(p string) { restarted <- p },
		})

		rec := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodPost, "/api/system/update", bytes.NewBufferString(`{"mode":"binary"}`))
		h.server.Handler().ServeHTTP(rec, req)
		if rec.Code != http.StatusAccepted {
			t.Fatalf("POST status = %d, want 202 (body: %s)", rec.Code, rec.Body.String())
		}

		deadline := time.Now().Add(5 * time.Second)
		var finalPayload updateStatusPayload
		for time.Now().Before(deadline) {
			pollRec := httptest.NewRecorder()
			pollReq := httptest.NewRequest(http.MethodGet, "/api/system/update", nil)
			h.server.Handler().ServeHTTP(pollRec, pollReq)
			_ = json.NewDecoder(pollRec.Body).Decode(&finalPayload)
			if finalPayload.State == "failed" {
				break
			}
			time.Sleep(20 * time.Millisecond)
		}
		if finalPayload.State != "failed" {
			t.Fatalf("State = %q, want failed", finalPayload.State)
		}
		if !strings.Contains(strings.ToLower(finalPayload.Error), "checksum") {
			t.Fatalf("Error = %q, want checksum mismatch mention", finalPayload.Error)
		}
		gotBytes, err := os.ReadFile(exePath)
		if err != nil {
			t.Fatalf("read exe: %v", err)
		}
		if !bytes.Equal(gotBytes, origBytes) {
			t.Fatalf("original executable was modified on checksum mismatch")
		}
		select {
		case p := <-restarted:
			t.Fatalf("unexpected restart invoked with %q", p)
		default:
		}
	})

	t.Run("valid checksum swaps binary and triggers restart", func(t *testing.T) {
		badChecksum.Store(false)
		exeDir := t.TempDir()
		exePath := filepath.Join(exeDir, "omeg4mon-router")
		if err := os.WriteFile(exePath, []byte("ORIGINAL_BINARY"), 0o755); err != nil {
			t.Fatalf("write exe: %v", err)
		}

		restarted := make(chan string, 1)
		h := newUpdaterTestHarness(t, system.Options{
			Build:            system.BuildInfo{Version: "v0.3.0", Commit: "1111111"},
			UpdateAPIBaseURL: ghServer.URL,
			ExecutablePath:   func() (string, error) { return exePath, nil },
			Restart:          func(p string) { restarted <- p },
		})

		rec := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodPost, "/api/system/update", bytes.NewBufferString(`{"mode":"auto"}`))
		h.server.Handler().ServeHTTP(rec, req)
		if rec.Code != http.StatusAccepted {
			t.Fatalf("POST status = %d, want 202 (body: %s)", rec.Code, rec.Body.String())
		}

		select {
		case gotPath := <-restarted:
			if gotPath != exePath {
				t.Fatalf("restart path = %q, want %q", gotPath, exePath)
			}
		case <-time.After(5 * time.Second):
			t.Fatal("timed out waiting for restart callback")
		}

		gotBytes, err := os.ReadFile(exePath)
		if err != nil {
			t.Fatalf("read replaced exe: %v", err)
		}
		if !bytes.Equal(gotBytes, newBin) {
			t.Fatal("executable content was not replaced with verified binary")
		}
	})
}

func TestSystemUpdateLoopbackGuardAndSourceDirtyWorktree(t *testing.T) {
	const fullSHA = "abcdef1234567890abcdef1234567890abcdef12"
	ghServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case "/repos/omeg4monx/router-gateway/commits/main":
			_ = json.NewEncoder(w).Encode(map[string]any{
				"sha": fullSHA,
				"commit": map[string]any{
					"message":   "chore: bump",
					"committer": map[string]any{"date": "2026-11-28T10:00:00Z"},
				},
			})
		default:
			http.NotFound(w, r)
		}
	}))
	defer ghServer.Close()

	t.Run("POST rejected with 403 on non-loopback without AllowInsecureRemote", func(t *testing.T) {
		h := newUpdaterTestHarness(t, system.Options{
			Config: config.Config{
				Host:                "0.0.0.0",
				AllowInsecureRemote: false,
			},
			Build:            system.BuildInfo{Version: "v0.3.0", Commit: "0000000"},
			UpdateAPIBaseURL: ghServer.URL,
		}, &stubAdminAuth{requireLogin: true, validToken: "admin-tok"})

		rec := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodPost, "/api/system/update", nil)
		req.Header.Set("Authorization", "Bearer admin-tok")
		h.server.Handler().ServeHTTP(rec, req)
		if rec.Code != http.StatusForbidden {
			t.Fatalf("status = %d, want 403 (body: %s)", rec.Code, rec.Body.String())
		}
	})

	t.Run("local git checkout with uncommitted changes reports dirty_worktree", func(t *testing.T) {
		if _, err := exec.LookPath("git"); err != nil {
			t.Skip("git not available")
		}
		repoDir := t.TempDir()
		runGit := func(args ...string) {
			cmd := exec.Command("git", args...)
			cmd.Dir = repoDir
			cmd.Env = append(os.Environ(),
				"GIT_AUTHOR_NAME=test",
				"GIT_AUTHOR_EMAIL=test@example.com",
				"GIT_COMMITTER_NAME=test",
				"GIT_COMMITTER_EMAIL=test@example.com",
			)
			if out, err := cmd.CombinedOutput(); err != nil {
				t.Fatalf("git %v: %v (%s)", args, err, string(out))
			}
		}
		runGit("init", "-b", "main")
		if err := os.WriteFile(filepath.Join(repoDir, "go.mod"), []byte("module github.com/omeg4monx/router-gateway\n"), 0o644); err != nil {
			t.Fatalf("write go.mod: %v", err)
		}
		runGit("add", "go.mod")
		runGit("commit", "-m", "init")
		// Make working tree dirty
		if err := os.WriteFile(filepath.Join(repoDir, "dirty.txt"), []byte("uncommitted"), 0o644); err != nil {
			t.Fatalf("write dirty.txt: %v", err)
		}

		exePath := filepath.Join(repoDir, "bin", "omeg4mon-router")
		if err := os.MkdirAll(filepath.Dir(exePath), 0o755); err != nil {
			t.Fatalf("mkdir bin: %v", err)
		}
		if err := os.WriteFile(exePath, makeFakeBinaryBytes(), 0o755); err != nil {
			t.Fatalf("write exe: %v", err)
		}

		h := newUpdaterTestHarness(t, system.Options{
			Build:            system.BuildInfo{Version: "v0.3.0", Commit: "0000000"},
			UpdateAPIBaseURL: ghServer.URL,
			ExecutablePath:   func() (string, error) { return exePath, nil },
			LookupPath:       func(name string) (string, error) { return "/usr/bin/" + name, nil },
		})

		rec := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodGet, "/api/system/update", nil)
		h.server.Handler().ServeHTTP(rec, req)
		if rec.Code != http.StatusOK {
			t.Fatalf("status = %d, want 200", rec.Code)
		}
		var payload updateStatusPayload
		if err := json.NewDecoder(rec.Body).Decode(&payload); err != nil {
			t.Fatalf("decode: %v", err)
		}
		if !payload.UpdateAvailable {
			t.Fatal("UpdateAvailable = false, want true")
		}
		if payload.Mode != "source" {
			t.Fatalf("Mode = %q, want source", payload.Mode)
		}
		if payload.CanApply {
			t.Fatal("CanApply = true, want false on dirty worktree")
		}
		if payload.UnavailableReason != "dirty_worktree" {
			t.Fatalf("UnavailableReason = %q, want dirty_worktree", payload.UnavailableReason)
		}
	})

	t.Run("GitHub API failure degrades gracefully with 200", func(t *testing.T) {
		failServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
			http.Error(w, "rate limit exceeded", http.StatusForbidden)
		}))
		defer failServer.Close()

		h := newUpdaterTestHarness(t, system.Options{
			Build:            system.BuildInfo{Version: "v0.3.0", Commit: "1234567"},
			UpdateAPIBaseURL: failServer.URL,
			ExecutablePath:   func() (string, error) { return filepath.Join(t.TempDir(), "omeg4mon-router"), nil },
		})

		rec := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodGet, "/api/system/update", nil)
		h.server.Handler().ServeHTTP(rec, req)
		if rec.Code != http.StatusOK {
			t.Fatalf("status = %d, want 200 on graceful degradation", rec.Code)
		}
		var payload updateStatusPayload
		if err := json.NewDecoder(rec.Body).Decode(&payload); err != nil {
			t.Fatalf("decode: %v", err)
		}
		if payload.UpdateAvailable {
			t.Fatal("UpdateAvailable = true, want false on API error")
		}
		if payload.Error == "" {
			t.Fatal("Error = empty, want diagnostic error message")
		}
	})
}

func TestSystemUpdateCommitTxtETagAndWrongBranch(t *testing.T) {
	const fullSHA = "9999999888888887777777666666655555555444"
	assetName := expectedPlatformAssetName()

	t.Run("uses COMMIT.txt when target_commitish is branch name and honors ETag 304", func(t *testing.T) {
		var commitCalls atomic.Int32
		var ghServer *httptest.Server
		ghServer = httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			switch r.URL.Path {
			case "/repos/omeg4monx/router-gateway/commits/main":
				commitCalls.Add(1)
				if r.Header.Get("If-None-Match") == `"etag-v1"` {
					w.WriteHeader(http.StatusNotModified)
					return
				}
				w.Header().Set("ETag", `"etag-v1"`)
				_ = json.NewEncoder(w).Encode(map[string]any{
					"sha": fullSHA,
					"commit": map[string]any{
						"message":   "feat: rolling commit",
						"committer": map[string]any{"date": "2026-11-28T11:00:00Z"},
					},
				})
			case "/repos/omeg4monx/router-gateway/releases/tags/main-latest":
				_ = json.NewEncoder(w).Encode(map[string]any{
					"tag_name":         "main-latest",
					"target_commitish": "main",
					"assets": []map[string]any{
						{"name": assetName, "browser_download_url": ghServer.URL + "/download/" + assetName},
						{"name": "SHA256SUMS.txt", "browser_download_url": ghServer.URL + "/download/SHA256SUMS.txt"},
						{"name": "COMMIT.txt", "browser_download_url": ghServer.URL + "/download/COMMIT.txt"},
					},
				})
			case "/download/COMMIT.txt":
				_, _ = w.Write([]byte(fullSHA + "\n"))
			default:
				http.NotFound(w, r)
			}
		}))
		defer ghServer.Close()

		exeDir := t.TempDir()
		exePath := filepath.Join(exeDir, "omeg4mon-router")
		if err := os.WriteFile(exePath, makeFakeBinaryBytes(), 0o755); err != nil {
			t.Fatalf("write exe: %v", err)
		}
		bakPath := exePath + ".bak"
		if err := os.WriteFile(bakPath, []byte("stale-backup"), 0o755); err != nil {
			t.Fatalf("write bak: %v", err)
		}

		h := newUpdaterTestHarness(t, system.Options{
			Build:            system.BuildInfo{Version: "v0.3.0", Commit: "0000000"},
			UpdateAPIBaseURL: ghServer.URL,
			ExecutablePath:   func() (string, error) { return exePath, nil },
		})

		// Stale .bak should have been cleaned up on Start()
		if _, err := os.Stat(bakPath); !os.IsNotExist(err) {
			t.Fatalf("expected stale %s to be removed on Start()", bakPath)
		}

		rec1 := httptest.NewRecorder()
		h.server.Handler().ServeHTTP(rec1, httptest.NewRequest(http.MethodGet, "/api/system/update", nil))
		var p1 updateStatusPayload
		if err := json.NewDecoder(rec1.Body).Decode(&p1); err != nil {
			t.Fatalf("decode p1: %v", err)
		}
		if !p1.UpdateAvailable || !p1.CanApply || p1.Mode != "binary" {
			t.Fatalf("p1 = %+v, want UpdateAvailable=true CanApply=true Mode=binary", p1)
		}

		// Force check should send If-None-Match and gracefully handle 304 Not Modified
		rec2 := httptest.NewRecorder()
		h.server.Handler().ServeHTTP(rec2, httptest.NewRequest(http.MethodGet, "/api/system/update?force=1", nil))
		var p2 updateStatusPayload
		if err := json.NewDecoder(rec2.Body).Decode(&p2); err != nil {
			t.Fatalf("decode p2: %v", err)
		}
		if commitCalls.Load() != 2 {
			t.Fatalf("commitCalls = %d, want 2", commitCalls.Load())
		}
		if p2.LatestCommit != fullSHA || p2.Error != "" || !p2.CanApply {
			t.Fatalf("p2 after 304 = %+v, want preserved commit and CanApply=true", p2)
		}
	})

	t.Run("local git checkout on non-main branch reports wrong_branch", func(t *testing.T) {
		if _, err := exec.LookPath("git"); err != nil {
			t.Skip("git not available")
		}
		repoDir := t.TempDir()
		runGit := func(args ...string) {
			cmd := exec.Command("git", args...)
			cmd.Dir = repoDir
			cmd.Env = append(os.Environ(),
				"GIT_AUTHOR_NAME=test",
				"GIT_AUTHOR_EMAIL=test@example.com",
				"GIT_COMMITTER_NAME=test",
				"GIT_COMMITTER_EMAIL=test@example.com",
			)
			if out, err := cmd.CombinedOutput(); err != nil {
				t.Fatalf("git %v: %v (%s)", args, err, string(out))
			}
		}
		runGit("init", "-b", "dev")
		if err := os.WriteFile(filepath.Join(repoDir, "go.mod"), []byte("module github.com/omeg4monx/router-gateway\n"), 0o644); err != nil {
			t.Fatalf("write go.mod: %v", err)
		}
		runGit("add", "go.mod")
		runGit("commit", "-m", "init on dev")

		exePath := filepath.Join(repoDir, "bin", "omeg4mon-router")
		if err := os.MkdirAll(filepath.Dir(exePath), 0o755); err != nil {
			t.Fatalf("mkdir bin: %v", err)
		}
		if err := os.WriteFile(exePath, makeFakeBinaryBytes(), 0o755); err != nil {
			t.Fatalf("write exe: %v", err)
		}

		ghServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if r.URL.Path == "/repos/omeg4monx/router-gateway/commits/main" {
				_ = json.NewEncoder(w).Encode(map[string]any{
					"sha": fullSHA,
					"commit": map[string]any{
						"message":   "feat: main commit",
						"committer": map[string]any{"date": "2026-11-28T11:00:00Z"},
					},
				})
				return
			}
			http.NotFound(w, r)
		}))
		defer ghServer.Close()

		h := newUpdaterTestHarness(t, system.Options{
			Build:            system.BuildInfo{Version: "v0.3.0", Commit: "0000000"},
			UpdateAPIBaseURL: ghServer.URL,
			ExecutablePath:   func() (string, error) { return exePath, nil },
			LookupPath:       func(name string) (string, error) { return "/usr/bin/" + name, nil },
		})

		rec := httptest.NewRecorder()
		h.server.Handler().ServeHTTP(rec, httptest.NewRequest(http.MethodGet, "/api/system/update", nil))
		var payload updateStatusPayload
		if err := json.NewDecoder(rec.Body).Decode(&payload); err != nil {
			t.Fatalf("decode: %v", err)
		}
		if payload.Mode != "source" || payload.CanApply || payload.UnavailableReason != "wrong_branch" {
			t.Fatalf("payload = %+v, want Mode=source CanApply=false UnavailableReason=wrong_branch", payload)
		}
	})
}
