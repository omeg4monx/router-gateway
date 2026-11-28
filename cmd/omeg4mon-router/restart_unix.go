//go:build !windows

package main

import "syscall"

func execReplacement(exePath string, args, env []string) error {
	//nolint:gosec // Self-restart intentionally executes the verified replacement binary path.
	return syscall.Exec(exePath, args, env)
}
