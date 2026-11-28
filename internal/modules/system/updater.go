package system

import (
	"bufio"
	"context"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"sync"
	"time"

	"github.com/omeg4monx/router-gateway/internal/platform/config"
	"github.com/omeg4monx/router-gateway/internal/platform/httpserver"
)

const (
	defaultUpdateAPIBaseURL = "https://api.github.com"
	defaultUpdateRepo       = "omeg4monx/router-gateway"
	defaultUpdateBranch     = "main"
	rollingReleaseTag       = "main-latest"
	updateCheckCacheTTL     = 5 * time.Minute
	maxMetadataBodyBytes    = 1 << 20   // 1 MiB
	maxBinaryDownloadBytes  = 256 << 20 // 256 MiB
	maxUpdateRequestBytes   = 16 << 10  // 16 KiB
)

type updateRequest struct {
	Mode string `json:"mode,omitempty"`
}

type updateStatusResponse struct {
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

type githubCommitResponse struct {
	SHA     string `json:"sha"`
	HTMLURL string `json:"html_url"`
	Commit  struct {
		Message   string `json:"message"`
		Committer struct {
			Date string `json:"date"`
		} `json:"committer"`
		Author struct {
			Date string `json:"date"`
		} `json:"author"`
	} `json:"commit"`
}

type githubReleaseResponse struct {
	TagName         string `json:"tag_name"`
	TargetCommitish string `json:"target_commitish"`
	Assets          []struct {
		Name               string `json:"name"`
		BrowserDownloadURL string `json:"browser_download_url"`
	} `json:"assets"`
}

type resolvedReleaseAssets struct {
	assetName      string
	assetURL       string
	checksumsURL   string
	commitURL      string
	targetCommit   string
	releaseMissing bool
}

type updater struct {
	cfg            config.Config
	build          BuildInfo
	httpClient     *http.Client
	apiBaseURL     string
	executablePath func() (string, error)
	lookupPath     func(string) (string, error)
	restart        func(exePath string)

	checkMu         sync.Mutex
	mu              sync.Mutex
	status          updateStatusResponse
	lastCheckedAt   time.Time
	commitETag      string
	cachedCommit    githubCommitResponse
	running         bool
	cachedAssets    resolvedReleaseAssets
	cachedRepoDir   string
	cachedExePath   string
	lifecycleCtx    context.Context
	lifecycleCancel context.CancelFunc
	wg              sync.WaitGroup
}

func newUpdater(opts Options) *updater {
	client := opts.HTTPClient
	if client == nil {
		client = &http.Client{Timeout: 60 * time.Second}
	}
	apiBase := strings.TrimRight(strings.TrimSpace(opts.UpdateAPIBaseURL), "/")
	if apiBase == "" {
		apiBase = defaultUpdateAPIBaseURL
	}
	exePathFn := opts.ExecutablePath
	if exePathFn == nil {
		exePathFn = os.Executable
	}
	lookupFn := opts.LookupPath
	if lookupFn == nil {
		lookupFn = exec.LookPath
	}
	restartFn := opts.Restart
	if restartFn == nil && opts.Shutdown != nil {
		restartFn = func(string) { opts.Shutdown() }
	}

	branch := strings.TrimSpace(opts.Config.UpdateBranch)
	if branch == "" {
		branch = defaultUpdateBranch
	}

	now := time.Now().UTC().Format(time.RFC3339)
	return &updater{
		cfg:            opts.Config,
		build:          opts.Build,
		httpClient:     client,
		apiBaseURL:     apiBase,
		executablePath: exePathFn,
		lookupPath:     lookupFn,
		restart:        restartFn,
		status: updateStatusResponse{
			CurrentVersion: opts.Build.Version,
			CurrentCommit:  opts.Build.Commit,
			Branch:         branch,
			Mode:           "none",
			State:          "idle",
			CheckedAt:      now,
			UpdatedAt:      now,
		},
	}
}

func (u *updater) start() {
	u.mu.Lock()
	if u.lifecycleCancel != nil {
		u.lifecycleCancel()
	}
	u.lifecycleCtx, u.lifecycleCancel = context.WithCancel(context.Background())
	u.mu.Unlock()

	if exePath, err := u.resolveExecutablePath(); err == nil && exePath != "" {
		_ = os.Remove(exePath + ".bak")
	}
}

func (u *updater) stop(ctx context.Context) {
	u.mu.Lock()
	cancel := u.lifecycleCancel
	u.lifecycleCancel = nil
	u.mu.Unlock()
	if cancel != nil {
		cancel()
	}

	done := make(chan struct{})
	go func() {
		u.wg.Wait()
		close(done)
	}()
	select {
	case <-done:
	case <-ctx.Done():
	}
}

func (u *updater) repoName() string {
	if r := strings.TrimSpace(u.cfg.UpdateRepo); r != "" {
		return r
	}
	return defaultUpdateRepo
}

func (u *updater) branchName() string {
	if b := strings.TrimSpace(u.cfg.UpdateBranch); b != "" {
		return b
	}
	return defaultUpdateBranch
}

func (u *updater) checkStatus(ctx context.Context, force bool) updateStatusResponse {
	u.mu.Lock()
	if u.running {
		snapshot := u.status
		u.mu.Unlock()
		return snapshot
	}
	if !force && !u.lastCheckedAt.IsZero() && time.Since(u.lastCheckedAt) < updateCheckCacheTTL {
		snapshot := u.status
		u.mu.Unlock()
		return snapshot
	}
	u.mu.Unlock()

	u.checkMu.Lock()
	defer u.checkMu.Unlock()

	u.mu.Lock()
	if u.running {
		snapshot := u.status
		u.mu.Unlock()
		return snapshot
	}
	if !force && !u.lastCheckedAt.IsZero() && time.Since(u.lastCheckedAt) < updateCheckCacheTTL {
		snapshot := u.status
		u.mu.Unlock()
		return snapshot
	}
	u.mu.Unlock()

	now := time.Now().UTC()
	nowStr := now.Format(time.RFC3339)

	exePath, exeErr := u.resolveExecutablePath()
	repoDir := detectLocalGitRepo(exePath)
	currentCommit := u.resolveCurrentCommit(repoDir)

	next := updateStatusResponse{
		CurrentVersion: u.build.Version,
		CurrentCommit:  currentCommit,
		Branch:         u.branchName(),
		Mode:           "none",
		State:          "idle",
		CheckedAt:      nowStr,
		UpdatedAt:      nowStr,
	}

	commitInfo, err := u.fetchLatestCommit(ctx)
	if err != nil {
		next.Error = fmt.Sprintf("failed to check %s branch: %v", u.branchName(), err)
		u.mu.Lock()
		if !u.running {
			u.status = next
			u.lastCheckedAt = now
		}
		snapshot := u.status
		u.mu.Unlock()
		return snapshot
	}

	next.LatestCommit = commitInfo.SHA
	next.LatestShortCommit = shortCommit(commitInfo.SHA)
	next.LatestMessage = firstLine(commitInfo.Commit.Message)
	next.LatestURL = commitInfo.HTMLURL
	if commitInfo.Commit.Committer.Date != "" {
		next.LatestPublishedAt = commitInfo.Commit.Committer.Date
	} else {
		next.LatestPublishedAt = commitInfo.Commit.Author.Date
	}

	next.UpdateAvailable = !commitsMatch(currentCommit, commitInfo.SHA)
	if next.UpdateAvailable && repoDir != "" && isCommitAncestorOfHead(ctx, repoDir, commitInfo.SHA) {
		next.UpdateAvailable = false
	}

	var resolvedAssets resolvedReleaseAssets
	if repoDir != "" && u.hasSourceToolchain() {
		next.Mode = "source"
		if next.UpdateAvailable {
			currentBranch := gitCurrentBranch(ctx, repoDir)
			switch {
			case exeErr != nil:
				next.UnavailableReason = "not_writable"
			case currentBranch != "" && currentBranch != u.branchName():
				next.UnavailableReason = "wrong_branch"
			case isGitWorktreeDirty(ctx, repoDir):
				next.UnavailableReason = "dirty_worktree"
			case !isDirWritable(filepath.Dir(exePath)):
				next.UnavailableReason = "not_writable"
			default:
				next.CanApply = true
			}
		}
	} else {
		next.Mode = "binary"
		if next.UpdateAvailable {
			assets, relErr := u.fetchRollingReleaseAssets(ctx)
			resolvedAssets = assets
			switch {
			case relErr != nil:
				next.UnavailableReason = "release_unavailable"
				next.Error = relErr.Error()
			case assets.releaseMissing || assets.assetURL == "" || assets.checksumsURL == "":
				next.UnavailableReason = "release_building"
			case isHexSHA(assets.targetCommit) && !commitsMatch(assets.targetCommit, commitInfo.SHA):
				next.UnavailableReason = "release_building"
			case exeErr != nil || !isDirWritable(filepath.Dir(exePath)):
				next.UnavailableReason = "not_writable"
			default:
				next.CanApply = true
			}
		}
	}

	u.mu.Lock()
	if !u.running {
		u.status = next
		u.lastCheckedAt = now
		u.cachedAssets = resolvedAssets
		u.cachedRepoDir = repoDir
		u.cachedExePath = exePath
	}
	snapshot := u.status
	u.mu.Unlock()
	return snapshot
}

func (u *updater) triggerUpdate(ctx context.Context, requestedMode string) (resp updateStatusResponse, statusCode int, errMsg string) {
	if requestedMode == "" {
		requestedMode = "auto"
	}
	if requestedMode != "auto" && requestedMode != "binary" && requestedMode != "source" {
		return updateStatusResponse{}, http.StatusBadRequest, "mode must be one of auto, binary, or source"
	}

	u.mu.Lock()
	if u.running {
		snapshot := u.status
		u.mu.Unlock()
		return snapshot, http.StatusConflict, "an update is already in progress"
	}
	u.mu.Unlock()

	status := u.checkStatus(ctx, true)
	if !status.UpdateAvailable {
		return status, http.StatusConflict, "router is already up to date"
	}

	mode := status.Mode
	if requestedMode != "auto" {
		mode = requestedMode
	}

	u.mu.Lock()
	if u.running {
		snapshot := u.status
		u.mu.Unlock()
		return snapshot, http.StatusConflict, "an update is already in progress"
	}
	exePath := u.cachedExePath
	repoDir := u.cachedRepoDir
	assets := u.cachedAssets
	baseCtx := u.lifecycleCtx
	if baseCtx == nil {
		baseCtx = context.Background()
	}

	if requestedMode == "auto" && !status.CanApply {
		u.mu.Unlock()
		reason := status.UnavailableReason
		if reason == "" {
			reason = "update cannot be applied in current environment"
		}
		return status, http.StatusConflict, fmt.Sprintf("cannot apply update: %s", reason)
	}
	if mode == "source" && repoDir == "" {
		u.mu.Unlock()
		return status, http.StatusConflict, "source update mode requires a local git checkout"
	}
	if mode == "binary" && (assets.assetURL == "" || assets.checksumsURL == "") {
		u.mu.Unlock()
		fetched, relErr := u.fetchRollingReleaseAssets(ctx)
		if relErr != nil || fetched.releaseMissing || fetched.assetURL == "" || fetched.checksumsURL == "" {
			return status, http.StatusConflict, "binary release assets are not available for the latest commit yet"
		}
		if isHexSHA(fetched.targetCommit) && !commitsMatch(fetched.targetCommit, status.LatestCommit) {
			return status, http.StatusConflict, "binary release assets are still building for the latest commit"
		}
		u.mu.Lock()
		if u.running {
			snapshot := u.status
			u.mu.Unlock()
			return snapshot, http.StatusConflict, "an update is already in progress"
		}
		assets = fetched
		u.cachedAssets = fetched
	}

	u.running = true
	nowStr := time.Now().UTC().Format(time.RFC3339)
	u.status.Mode = mode
	u.status.State = "checking"
	u.status.Message = "Preparing update"
	u.status.Error = ""
	u.status.UpdatedAt = nowStr
	snapshot := u.status
	u.wg.Add(1)
	u.mu.Unlock()

	go func() {
		defer u.wg.Done()
		runCtx, cancel := context.WithTimeout(baseCtx, 10*time.Minute)
		defer cancel()

		var err error
		if mode == "source" {
			err = u.applySourceUpdate(runCtx, repoDir, exePath)
		} else {
			err = u.applyBinaryUpdate(runCtx, assets, exePath)
		}

		if err != nil {
			u.setState("failed", "Update failed", err.Error(), false)
			return
		}

		u.setState("restarting", "Update applied; restarting router", "", false)
		time.Sleep(150 * time.Millisecond)
		if u.restart != nil {
			u.restart(exePath)
		}
	}()

	return snapshot, http.StatusAccepted, ""
}

func (u *updater) setState(state, message, errMsg string, keepRunning bool) {
	u.mu.Lock()
	defer u.mu.Unlock()
	u.status.State = state
	u.status.Message = message
	u.status.Error = errMsg
	u.status.UpdatedAt = time.Now().UTC().Format(time.RFC3339)
	if !keepRunning && (state == "failed" || state == "restarting" || state == "idle") {
		u.running = false
	}
}

func (u *updater) applyBinaryUpdate(ctx context.Context, assets resolvedReleaseAssets, exePath string) error {
	if exePath == "" {
		return errors.New("could not resolve current executable path")
	}

	u.setState("downloading", "Downloading release checksums", "", true)
	expectedSHA, err := u.fetchExpectedChecksum(ctx, assets.checksumsURL, assets.assetName)
	if err != nil {
		return err
	}

	u.setState("downloading", fmt.Sprintf("Downloading %s", assets.assetName), "", true)
	tmpPath := filepath.Join(filepath.Dir(exePath), fmt.Sprintf(".omeg4mon-router.update.%d.tmp", os.Getpid()))
	defer func() {
		_ = os.Remove(tmpPath)
	}()

	actualSHA, err := u.downloadFileAndHash(ctx, assets.assetURL, tmpPath)
	if err != nil {
		return err
	}

	u.setState("verifying", "Verifying binary SHA256 checksum", "", true)
	if !strings.EqualFold(actualSHA, expectedSHA) {
		return fmt.Errorf("checksum mismatch: expected %s, got %s", expectedSHA, actualSHA)
	}
	if err := verifyExecutableHeader(tmpPath); err != nil {
		return err
	}

	u.setState("applying", "Replacing router binary", "", true)
	if err := atomicReplaceExecutable(exePath, tmpPath); err != nil {
		return err
	}
	return nil
}

func (u *updater) applySourceUpdate(ctx context.Context, repoDir, exePath string) error {
	if exePath == "" {
		return errors.New("could not resolve current executable path")
	}

	u.setState("checking", "Verifying clean git working tree", "", true)
	if isGitWorktreeDirty(ctx, repoDir) {
		return errors.New("git working tree has uncommitted changes")
	}

	branch := u.branchName()
	u.setState("downloading", fmt.Sprintf("Fetching origin/%s", branch), "", true)
	if err := runCmd(ctx, repoDir, nil, "git", "fetch", "origin", branch); err != nil {
		return fmt.Errorf("git fetch failed: %w", err)
	}

	currentBranch, _ := runCmdOutput(ctx, repoDir, "git", "rev-parse", "--abbrev-ref", "HEAD")
	if strings.TrimSpace(currentBranch) != branch {
		if err := runCmd(ctx, repoDir, nil, "git", "checkout", branch); err != nil {
			return fmt.Errorf("git checkout %s failed: %w", branch, err)
		}
	}
	if err := runCmd(ctx, repoDir, nil, "git", "pull", "--ff-only", "origin", branch); err != nil {
		return fmt.Errorf("git pull --ff-only failed: %w", err)
	}

	u.setState("building", "Building embedded web UI", "", true)
	webDir := filepath.Join(repoDir, "web")
	if err := runCmd(ctx, webDir, nil, "bun", "run", "build"); err != nil {
		return fmt.Errorf("web build failed: %w", err)
	}
	uiDist := filepath.Join(repoDir, "cmd", "omeg4mon-router", "ui_dist")
	if err := replaceDirectory(filepath.Join(webDir, "out"), uiDist); err != nil {
		return fmt.Errorf("sync ui_dist failed: %w", err)
	}

	u.setState("building", "Compiling updated Go binary", "", true)
	newVersion, _ := runCmdOutput(ctx, repoDir, "git", "describe", "--tags", "--always")
	if strings.TrimSpace(newVersion) == "" {
		newVersion = "main"
	}
	newCommit, _ := runCmdOutput(ctx, repoDir, "git", "rev-parse", "--short", "HEAD")
	if strings.TrimSpace(newCommit) == "" {
		newCommit = "none"
	}
	buildTime := time.Now().UTC().Format(time.RFC3339)
	ldflags := fmt.Sprintf("-s -w -X main.version=%s -X main.commit=%s -X main.buildTime=%s",
		strings.TrimSpace(newVersion), strings.TrimSpace(newCommit), buildTime)

	tmpPath := filepath.Join(filepath.Dir(exePath), fmt.Sprintf(".omeg4mon-router.update.%d.tmp", os.Getpid()))
	defer func() {
		_ = os.Remove(tmpPath)
	}()

	buildEnv := append(os.Environ(), "CGO_ENABLED=0")
	if err := runCmd(ctx, repoDir, buildEnv, "go", "build", "-tags", "embedui", "-trimpath", "-ldflags", ldflags, "-o", tmpPath, "./cmd/omeg4mon-router"); err != nil {
		return fmt.Errorf("go build failed: %w", err)
	}

	u.setState("verifying", "Verifying compiled binary", "", true)
	if err := verifyExecutableHeader(tmpPath); err != nil {
		return err
	}

	u.setState("applying", "Replacing router binary", "", true)
	if err := atomicReplaceExecutable(exePath, tmpPath); err != nil {
		return err
	}
	return nil
}

func (u *updater) fetchLatestCommit(ctx context.Context) (githubCommitResponse, error) {
	url := fmt.Sprintf("%s/repos/%s/commits/%s", u.apiBaseURL, u.repoName(), u.branchName())
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, http.NoBody)
	if err != nil {
		return githubCommitResponse{}, err
	}
	req.Header.Set("Accept", "application/vnd.github+json")
	req.Header.Set("User-Agent", "omeg4mon-router-updater")
	if token := strings.TrimSpace(u.cfg.UpdateToken); token != "" {
		req.Header.Set("Authorization", "Bearer "+token)
	}

	u.mu.Lock()
	etag := u.commitETag
	cached := u.cachedCommit
	u.mu.Unlock()
	if etag != "" && cached.SHA != "" {
		req.Header.Set("If-None-Match", etag)
	}

	resp, err := u.httpClient.Do(req)
	if err != nil {
		return githubCommitResponse{}, err
	}
	defer func() { _ = resp.Body.Close() }()

	if resp.StatusCode == http.StatusNotModified && cached.SHA != "" {
		return cached, nil
	}
	if resp.StatusCode == http.StatusNotFound {
		return githubCommitResponse{}, errNotFound
	}
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		body, _ := io.ReadAll(io.LimitReader(resp.Body, 512))
		return githubCommitResponse{}, fmt.Errorf("HTTP %d: %s", resp.StatusCode, strings.TrimSpace(string(body)))
	}

	var out githubCommitResponse
	if err := json.NewDecoder(io.LimitReader(resp.Body, maxMetadataBodyBytes)).Decode(&out); err != nil {
		return githubCommitResponse{}, err
	}
	if strings.TrimSpace(out.SHA) == "" {
		return githubCommitResponse{}, errors.New("empty commit SHA returned by GitHub API")
	}

	u.mu.Lock()
	if newETag := strings.TrimSpace(resp.Header.Get("ETag")); newETag != "" {
		u.commitETag = newETag
	}
	u.cachedCommit = out
	u.mu.Unlock()
	return out, nil
}

func (u *updater) fetchRollingReleaseAssets(ctx context.Context) (resolvedReleaseAssets, error) {
	url := fmt.Sprintf("%s/repos/%s/releases/tags/%s", u.apiBaseURL, u.repoName(), rollingReleaseTag)
	var rel githubReleaseResponse
	err := u.getJSON(ctx, url, &rel)
	if err != nil {
		if errors.Is(err, errNotFound) {
			return resolvedReleaseAssets{releaseMissing: true}, nil
		}
		return resolvedReleaseAssets{}, err
	}

	wantAsset := platformAssetName(runtime.GOOS, runtime.GOARCH)
	res := resolvedReleaseAssets{
		assetName:    wantAsset,
		targetCommit: strings.TrimSpace(rel.TargetCommitish),
	}
	for _, a := range rel.Assets {
		switch a.Name {
		case wantAsset:
			res.assetURL = a.BrowserDownloadURL
		case "SHA256SUMS.txt":
			res.checksumsURL = a.BrowserDownloadURL
		case "COMMIT.txt":
			res.commitURL = a.BrowserDownloadURL
		}
	}
	if res.commitURL != "" {
		if sha, shaErr := u.fetchSmallText(ctx, res.commitURL); shaErr == nil && isHexSHA(sha) {
			res.targetCommit = sha
		}
	}
	return res, nil
}

func (u *updater) fetchSmallText(ctx context.Context, textURL string) (string, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, textURL, http.NoBody)
	if err != nil {
		return "", err
	}
	req.Header.Set("User-Agent", "omeg4mon-router-updater")
	if token := strings.TrimSpace(u.cfg.UpdateToken); token != "" {
		req.Header.Set("Authorization", "Bearer "+token)
	}
	resp, err := u.httpClient.Do(req)
	if err != nil {
		return "", err
	}
	defer func() { _ = resp.Body.Close() }()
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return "", fmt.Errorf("HTTP %d", resp.StatusCode)
	}
	raw, err := io.ReadAll(io.LimitReader(resp.Body, 4096))
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(string(raw)), nil
}

var errNotFound = errors.New("not found")

func (u *updater) getJSON(ctx context.Context, url string, dst any) error {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, http.NoBody)
	if err != nil {
		return err
	}
	req.Header.Set("Accept", "application/vnd.github+json")
	req.Header.Set("User-Agent", "omeg4mon-router-updater")
	if token := strings.TrimSpace(u.cfg.UpdateToken); token != "" {
		req.Header.Set("Authorization", "Bearer "+token)
	}

	resp, err := u.httpClient.Do(req)
	if err != nil {
		return err
	}
	defer func() { _ = resp.Body.Close() }()

	if resp.StatusCode == http.StatusNotFound {
		return errNotFound
	}
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		body, _ := io.ReadAll(io.LimitReader(resp.Body, 512))
		return fmt.Errorf("HTTP %d: %s", resp.StatusCode, strings.TrimSpace(string(body)))
	}

	dec := json.NewDecoder(io.LimitReader(resp.Body, maxMetadataBodyBytes))
	return dec.Decode(dst)
}

func (u *updater) fetchExpectedChecksum(ctx context.Context, checksumsURL, assetName string) (string, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, checksumsURL, http.NoBody)
	if err != nil {
		return "", err
	}
	req.Header.Set("User-Agent", "omeg4mon-router-updater")
	if token := strings.TrimSpace(u.cfg.UpdateToken); token != "" {
		req.Header.Set("Authorization", "Bearer "+token)
	}

	resp, err := u.httpClient.Do(req)
	if err != nil {
		return "", fmt.Errorf("download SHA256SUMS.txt: %w", err)
	}
	defer func() { _ = resp.Body.Close() }()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return "", fmt.Errorf("download SHA256SUMS.txt: HTTP %d", resp.StatusCode)
	}

	scanner := bufio.NewScanner(io.LimitReader(resp.Body, maxMetadataBodyBytes))
	for scanner.Scan() {
		fields := strings.Fields(scanner.Text())
		if len(fields) < 2 {
			continue
		}
		name := strings.TrimPrefix(fields[1], "*")
		if filepath.Base(name) == assetName {
			return strings.ToLower(fields[0]), nil
		}
	}
	if err := scanner.Err(); err != nil {
		return "", fmt.Errorf("read SHA256SUMS.txt: %w", err)
	}
	return "", fmt.Errorf("checksum entry not found for %s", assetName)
}

func (u *updater) downloadFileAndHash(ctx context.Context, downloadURL, destPath string) (string, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, downloadURL, http.NoBody)
	if err != nil {
		return "", err
	}
	req.Header.Set("User-Agent", "omeg4mon-router-updater")
	if token := strings.TrimSpace(u.cfg.UpdateToken); token != "" {
		req.Header.Set("Authorization", "Bearer "+token)
	}

	resp, err := u.httpClient.Do(req)
	if err != nil {
		return "", fmt.Errorf("download binary: %w", err)
	}
	defer func() { _ = resp.Body.Close() }()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return "", fmt.Errorf("download binary: HTTP %d", resp.StatusCode)
	}

	f, err := os.OpenFile(destPath, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0o600)
	if err != nil {
		return "", fmt.Errorf("create temp binary: %w", err)
	}

	hasher := sha256.New()
	written, copyErr := io.Copy(io.MultiWriter(f, hasher), io.LimitReader(resp.Body, maxBinaryDownloadBytes+1))
	closeErr := f.Close()
	if copyErr != nil {
		return "", fmt.Errorf("write temp binary: %w", copyErr)
	}
	if closeErr != nil {
		return "", fmt.Errorf("close temp binary: %w", closeErr)
	}
	if written == 0 {
		return "", errors.New("downloaded binary is empty")
	}
	if written > maxBinaryDownloadBytes {
		return "", errors.New("downloaded binary exceeds maximum allowed size")
	}

	return hex.EncodeToString(hasher.Sum(nil)), nil
}

func (u *updater) resolveExecutablePath() (string, error) {
	p, err := u.executablePath()
	if err != nil {
		return "", err
	}
	if resolved, evalErr := filepath.EvalSymlinks(p); evalErr == nil {
		return resolved, nil
	}
	return p, nil
}

func (u *updater) resolveCurrentCommit(repoDir string) string {
	c := strings.TrimSpace(u.build.Commit)
	if c != "" && !strings.EqualFold(c, "none") && !strings.EqualFold(c, "unknown") {
		return c
	}
	if repoDir != "" {
		ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
		defer cancel()
		if out, err := runCmdOutput(ctx, repoDir, "git", "rev-parse", "--short", "HEAD"); err == nil {
			if trimmed := strings.TrimSpace(out); trimmed != "" {
				return trimmed
			}
		}
	}
	return c
}

func (u *updater) hasSourceToolchain() bool {
	for _, bin := range []string{"git", "go", "bun"} {
		if _, err := u.lookupPath(bin); err != nil {
			return false
		}
	}
	return true
}

func platformAssetName(goos, goarch string) string {
	ext := ""
	if goos == "windows" {
		ext = ".exe"
	}
	return fmt.Sprintf("omeg4mon-router-%s-%s%s", goos, goarch, ext)
}

func commitsMatch(current, latest string) bool {
	c := strings.ToLower(strings.TrimSpace(current))
	l := strings.ToLower(strings.TrimSpace(latest))
	if c == "" || l == "" || c == "none" || c == "unknown" {
		return false
	}
	if len(c) < 7 {
		return c == l
	}
	return strings.HasPrefix(l, c) || strings.HasPrefix(c, l)
}

func shortCommit(sha string) string {
	s := strings.TrimSpace(sha)
	if len(s) > 7 {
		return s[:7]
	}
	return s
}

func firstLine(s string) string {
	line, _, _ := strings.Cut(strings.TrimSpace(s), "\n")
	return strings.TrimSpace(line)
}

func detectLocalGitRepo(exePath string) string {
	if exePath == "" {
		return ""
	}
	dir := filepath.Dir(exePath)
	for range 6 {
		gitPath := filepath.Join(dir, ".git")
		goModPath := filepath.Join(dir, "go.mod")
		if _, err := os.Stat(gitPath); err == nil {
			if raw, readErr := os.ReadFile(goModPath); readErr == nil {
				if strings.Contains(string(raw), "github.com/omeg4monx/router-gateway") {
					return dir
				}
			}
		}
		parent := filepath.Dir(dir)
		if parent == dir {
			break
		}
		dir = parent
	}
	return ""
}

func isHexSHA(s string) bool {
	trimmed := strings.TrimSpace(s)
	if len(trimmed) < 7 || len(trimmed) > 40 {
		return false
	}
	for i := range len(trimmed) {
		ch := trimmed[i]
		isHex := (ch >= '0' && ch <= '9') || (ch >= 'a' && ch <= 'f') || (ch >= 'A' && ch <= 'F')
		if !isHex {
			return false
		}
	}
	return true
}

func gitCurrentBranch(ctx context.Context, repoDir string) string {
	out, err := runCmdOutput(ctx, repoDir, "git", "rev-parse", "--abbrev-ref", "HEAD")
	if err != nil {
		return ""
	}
	return strings.TrimSpace(out)
}

func isCommitAncestorOfHead(ctx context.Context, repoDir, sha string) bool {
	if !isHexSHA(sha) {
		return false
	}
	return runCmd(ctx, repoDir, nil, "git", "merge-base", "--is-ancestor", sha, "HEAD") == nil
}

func isGitWorktreeDirty(ctx context.Context, repoDir string) bool {
	out, err := runCmdOutput(ctx, repoDir, "git", "status", "--porcelain")
	if err != nil {
		return true
	}
	return strings.TrimSpace(out) != ""
}

func isDirWritable(dir string) bool {
	f, err := os.CreateTemp(dir, ".omeg4mon-write-check-*")
	if err != nil {
		return false
	}
	name := f.Name()
	_ = f.Close()
	_ = os.Remove(name)
	return true
}

func verifyExecutableHeader(path string) error {
	f, err := os.Open(path)
	if err != nil {
		return fmt.Errorf("open downloaded binary: %w", err)
	}
	defer func() { _ = f.Close() }()

	var hdr [4]byte
	if _, err := io.ReadFull(f, hdr[:]); err != nil {
		return fmt.Errorf("invalid binary header: %w", err)
	}
	switch runtime.GOOS {
	case "windows":
		if hdr[0] != 'M' || hdr[1] != 'Z' {
			return errors.New("invalid PE executable header")
		}
	case "darwin":
		isMachO := (hdr[0] == 0xcf && hdr[1] == 0xfa && hdr[2] == 0xed && hdr[3] == 0xfe) ||
			(hdr[0] == 0xfe && hdr[1] == 0xed && hdr[2] == 0xfa && hdr[3] == 0xcf) ||
			(hdr[0] == 0xca && hdr[1] == 0xfe && hdr[2] == 0xba && hdr[3] == 0xbe)
		if !isMachO {
			return errors.New("invalid Mach-O executable header")
		}
	default:
		if hdr[0] != 0x7f || hdr[1] != 'E' || hdr[2] != 'L' || hdr[3] != 'F' {
			return errors.New("invalid ELF executable header")
		}
	}
	return nil
}

func atomicReplaceExecutable(targetPath, tmpPath string) error {
	//nolint:gosec // Replaced binary must be executable.
	if err := os.Chmod(tmpPath, 0o755); err != nil {
		return fmt.Errorf("chmod temp binary: %w", err)
	}
	backupPath := targetPath + ".bak"
	_ = os.Remove(backupPath)

	hadTarget := false
	if _, err := os.Stat(targetPath); err == nil {
		hadTarget = true
		if err := os.Rename(targetPath, backupPath); err != nil {
			return fmt.Errorf("backup existing binary: %w", err)
		}
	}

	if err := os.Rename(tmpPath, targetPath); err != nil {
		if hadTarget {
			_ = os.Rename(backupPath, targetPath)
		}
		return fmt.Errorf("install updated binary: %w", err)
	}
	return nil
}

func replaceDirectory(srcDir, dstDir string) error {
	if err := os.RemoveAll(dstDir); err != nil {
		return err
	}
	return copyDir(srcDir, dstDir)
}

func copyDir(src, dst string) error {
	if err := os.MkdirAll(dst, 0o750); err != nil {
		return err
	}
	srcRoot, err := os.OpenRoot(src)
	if err != nil {
		return err
	}
	defer func() { _ = srcRoot.Close() }()

	dstRoot, err := os.OpenRoot(dst)
	if err != nil {
		return err
	}
	defer func() { _ = dstRoot.Close() }()

	return filepath.WalkDir(src, func(path string, d os.DirEntry, walkErr error) error {
		if walkErr != nil {
			return walkErr
		}
		rel, err := filepath.Rel(src, path)
		if err != nil {
			return err
		}
		if rel == "." {
			return nil
		}
		if d.IsDir() {
			return dstRoot.MkdirAll(rel, 0o750)
		}
		data, err := srcRoot.ReadFile(rel)
		if err != nil {
			return err
		}
		return dstRoot.WriteFile(rel, data, 0o600)
	})
}

func runCmd(ctx context.Context, dir string, env []string, name string, args ...string) error {
	//nolint:gosec // Invokes vetted git/go/bun toolchain commands with fixed binary names.
	cmd := exec.CommandContext(ctx, name, args...)
	cmd.Dir = dir
	if env != nil {
		cmd.Env = env
	}
	out, err := cmd.CombinedOutput()
	if err != nil {
		trimmed := strings.TrimSpace(string(out))
		if len(trimmed) > 400 {
			trimmed = trimmed[len(trimmed)-400:]
		}
		return fmt.Errorf("%s %s: %w (%s)", name, strings.Join(args, " "), err, trimmed)
	}
	return nil
}

func runCmdOutput(ctx context.Context, dir, name string, args ...string) (string, error) {
	//nolint:gosec // Invokes vetted git toolchain commands with fixed binary names.
	cmd := exec.CommandContext(ctx, name, args...)
	cmd.Dir = dir
	out, err := cmd.Output()
	return string(out), err
}

func (m *Module) handleGetUpdate(w http.ResponseWriter, r *http.Request) {
	forceParam := strings.TrimSpace(r.URL.Query().Get("force"))
	force := forceParam == "1" || strings.EqualFold(forceParam, "true")
	status := m.updater.checkStatus(r.Context(), force)
	httpserver.WriteJSON(w, http.StatusOK, status)
}

func (m *Module) handlePostUpdate(w http.ResponseWriter, r *http.Request) {
	if !m.cfg.IsLoopback() && !m.cfg.AllowInsecureRemote {
		httpserver.WriteError(w, r, http.StatusForbidden, "forbidden", "update endpoint is disabled on non-loopback bind addresses")
		return
	}

	var req updateRequest
	if r.Body != nil {
		r.Body = http.MaxBytesReader(w, r.Body, maxUpdateRequestBytes)
		dec := json.NewDecoder(r.Body)
		dec.DisallowUnknownFields()
		if err := dec.Decode(&req); err != nil && !errors.Is(err, io.EOF) {
			httpserver.WriteError(w, r, http.StatusBadRequest, "invalid_request", "invalid JSON request body")
			return
		}
	}

	status, code, errMsg := m.updater.triggerUpdate(r.Context(), strings.TrimSpace(req.Mode))
	if code == http.StatusAccepted {
		httpserver.WriteJSON(w, http.StatusAccepted, status)
		return
	}
	errCode := "update_conflict"
	if code == http.StatusBadRequest {
		errCode = "invalid_request"
	}
	httpserver.WriteError(w, r, code, errCode, errMsg)
}
