//go:build windows

package main

import (
	"os"
	"os/exec"
)

func execReplacement(exePath string, args, env []string) error {
	var childArgs []string
	if len(args) > 1 {
		childArgs = args[1:]
	}
	//nolint:gosec // Self-restart intentionally spawns the verified replacement binary path.
	cmd := exec.Command(exePath, childArgs...)
	cmd.Env = env
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Start()
}
