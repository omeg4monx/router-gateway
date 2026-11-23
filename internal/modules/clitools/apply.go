package clitools

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"
)

const (
	backupSuffix     = ".om4.bak."
	backupTimeFormat = "20060102T150405.000000000Z"
	dirPermissions   = 0o700
	filePermissions  = 0o600
)

// ErrBackupNotFound indicates no backup file exists for the requested tool config.
var ErrBackupNotFound = errors.New("clitools: no backup found")

func applyConfigFile(targetPath, content string) (string, error) {
	if targetPath == "" {
		return "", errors.New("clitools: target path is empty")
	}
	parentDir := filepath.Dir(targetPath)
	if err := os.MkdirAll(parentDir, dirPermissions); err != nil {
		return "", fmt.Errorf("create target directory: %w", err)
	}

	var backupPath string
	info, err := os.Lstat(targetPath)
	if err == nil {
		if info.IsDir() {
			return "", fmt.Errorf("target path %q is a directory", targetPath)
		}
		origBytes, readErr := os.ReadFile(targetPath)
		if readErr != nil {
			return "", fmt.Errorf("read existing config: %w", readErr)
		}
		backupPath = nextBackupPath(targetPath, time.Now().UTC())
		if bakErr := writeFileAtomic(backupPath, origBytes, filePermissions); bakErr != nil {
			return "", fmt.Errorf("write backup file: %w", bakErr)
		}
	} else if !errors.Is(err, os.ErrNotExist) {
		return "", fmt.Errorf("inspect target config: %w", err)
	}

	if err := writeFileAtomic(targetPath, []byte(content), filePermissions); err != nil {
		return "", fmt.Errorf("write config file: %w", err)
	}
	return backupPath, nil
}

func rollbackConfigFile(targetPath string) (string, error) {
	if targetPath == "" {
		return "", errors.New("clitools: target path is empty")
	}
	parentDir := filepath.Dir(targetPath)
	entries, err := os.ReadDir(parentDir)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return "", ErrBackupNotFound
		}
		return "", fmt.Errorf("read config directory: %w", err)
	}

	prefix := filepath.Base(targetPath) + backupSuffix
	var newestName string
	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}
		name := entry.Name()
		if strings.HasPrefix(name, prefix) && filepath.Base(name) == name && name > newestName {
			newestName = name
		}
	}
	if newestName == "" {
		return "", ErrBackupNotFound
	}

	backupPath := filepath.Join(parentDir, newestName)
	backupBytes, err := os.ReadFile(backupPath) //nolint:gosec // G703: backupPath is derived from fixed catalog path and validated ReadDir prefix.
	if err != nil {
		return "", fmt.Errorf("read backup file: %w", err)
	}

	if err := writeFileAtomic(targetPath, backupBytes, filePermissions); err != nil {
		return "", fmt.Errorf("restore config file: %w", err)
	}
	return backupPath, nil
}

func nextBackupPath(targetPath string, now time.Time) string {
	ts := now.UTC()
	for {
		candidate := targetPath + backupSuffix + ts.Format(backupTimeFormat)
		if _, err := os.Lstat(candidate); errors.Is(err, os.ErrNotExist) {
			return candidate
		}
		ts = ts.Add(time.Nanosecond)
	}
}

func writeFileAtomic(targetPath string, data []byte, mode os.FileMode) error {
	dir := filepath.Dir(targetPath)
	tmpFile, err := os.CreateTemp(dir, "."+filepath.Base(targetPath)+".tmp-*")
	if err != nil {
		return err
	}
	tmpPath := tmpFile.Name()
	cleanup := true
	defer func() {
		if cleanup {
			_ = tmpFile.Close()
			_ = os.Remove(tmpPath) //nolint:gosec // G703: tmpPath was created by os.CreateTemp inside target parentDir.
		}
	}()

	if err := tmpFile.Chmod(mode); err != nil {
		return err
	}
	if _, err := tmpFile.Write(data); err != nil {
		return err
	}
	if err := tmpFile.Sync(); err != nil {
		return err
	}
	if err := tmpFile.Close(); err != nil {
		return err
	}
	if err := os.Rename(tmpPath, targetPath); err != nil { //nolint:gosec // G703: tmpPath is in same parentDir and targetPath comes from fixed catalog RelParts.
		return err
	}
	cleanup = false
	return nil
}
