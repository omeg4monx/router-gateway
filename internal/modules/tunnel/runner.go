package tunnel

import (
	"context"
	"io"
	"os/exec"
	"sync"
)

// ProcessHandle abstracts a long-running child tunnel process (such as cloudflared).
type ProcessHandle interface {
	PID() int
	Output() io.Reader
	Wait() error
	Kill() error
}

// ProcessRunner abstracts binary lookup and command execution so tunnel
// lifecycle behavior can be verified deterministically in tests.
type ProcessRunner interface {
	LookPath(name string) (string, error)
	Start(ctx context.Context, name string, args ...string) (ProcessHandle, error)
	Run(ctx context.Context, name string, args ...string) ([]byte, error)
}

type osProcessHandle struct {
	cmd    *exec.Cmd
	pr     *io.PipeReader
	pw     *io.PipeWriter
	cancel context.CancelFunc

	waitOnce sync.Once
	waitErr  error
}

func (h *osProcessHandle) PID() int {
	if h.cmd == nil || h.cmd.Process == nil {
		return 0
	}
	return h.cmd.Process.Pid
}

func (h *osProcessHandle) Output() io.Reader {
	return h.pr
}

func (h *osProcessHandle) Wait() error {
	h.waitOnce.Do(func() {
		h.waitErr = h.cmd.Wait()
		_ = h.pw.Close()
	})
	return h.waitErr
}

func (h *osProcessHandle) Kill() error {
	if h.cancel != nil {
		h.cancel()
	}
	if h.cmd != nil && h.cmd.Process != nil {
		_ = h.cmd.Process.Kill()
	}
	return nil
}

type osProcessRunner struct{}

func (osProcessRunner) LookPath(name string) (string, error) {
	return exec.LookPath(name)
}

func (osProcessRunner) Start(ctx context.Context, name string, args ...string) (ProcessHandle, error) {
	cmdCtx, cancel := context.WithCancel(ctx)
	//nolint:gosec // G204: binary name is strictly allowlisted ("cloudflared"/"tailscale") and args are validated.
	cmd := exec.CommandContext(cmdCtx, name, args...)
	pr, pw := io.Pipe()
	cmd.Stdout = pw
	cmd.Stderr = pw
	if err := cmd.Start(); err != nil {
		cancel()
		_ = pr.Close()
		_ = pw.Close()
		return nil, err
	}
	return &osProcessHandle{
		cmd:    cmd,
		pr:     pr,
		pw:     pw,
		cancel: cancel,
	}, nil
}

func (osProcessRunner) Run(ctx context.Context, name string, args ...string) ([]byte, error) {
	//nolint:gosec // G204: binary name is strictly allowlisted ("cloudflared"/"tailscale") and args are validated.
	cmd := exec.CommandContext(ctx, name, args...)
	return cmd.CombinedOutput()
}
