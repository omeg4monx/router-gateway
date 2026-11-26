package importer

import (
	"bytes"
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"net/url"
	"os"
	"path/filepath"
	"sort"
	"strings"

	_ "modernc.org/sqlite"

	"github.com/omeg4monx/router-gateway/internal/contracts"
)

var sqliteMagic = []byte("SQLite format 3\x00")

type legacyKeyRecord struct {
	Name          string
	AllowedModels []string
	RateLimitRPM  int
}

type legacyExtract struct {
	TablesFound    []string
	Settings       map[string]any
	Combos         []contracts.ComboSnapshot
	Keys           []legacyKeyRecord
	ProvidersCount int
	Warnings       []string
}

func validateSQLiteSource(rawPath string) (string, error) {
	trimmed := strings.TrimSpace(rawPath)
	if trimmed == "" {
		return "", errors.New("source_path is required")
	}
	cleaned := filepath.Clean(trimmed)
	info, err := os.Stat(cleaned)
	if err != nil {
		return "", fmt.Errorf("cannot access source_path: %w", err)
	}
	if !info.Mode().IsRegular() {
		return "", errors.New("source_path must be a regular SQLite file")
	}
	if info.Size() < int64(len(sqliteMagic)) {
		return "", errors.New("source_path is too small to be a valid SQLite database")
	}
	f, err := os.Open(cleaned)
	if err != nil {
		return "", fmt.Errorf("open source_path: %w", err)
	}
	defer f.Close()

	header := make([]byte, len(sqliteMagic))
	if _, err := f.Read(header); err != nil {
		return "", fmt.Errorf("read SQLite header: %w", err)
	}
	if !bytes.Equal(header, sqliteMagic) {
		return "", errors.New("source_path is not a valid SQLite 3 database")
	}
	return cleaned, nil
}

func readLegacySQLite(ctx context.Context, path string) (legacyExtract, error) {
	dsn := "file:" + url.PathEscape(path) + "?mode=ro&_query_only=1"
	conn, err := sql.Open("sqlite", dsn)
	if err != nil {
		return legacyExtract{}, fmt.Errorf("open legacy sqlite read-only: %w", err)
	}
	defer conn.Close()
	conn.SetMaxOpenConns(1)

	tables, err := listTables(ctx, conn)
	if err != nil {
		return legacyExtract{}, err
	}
	tableSet := make(map[string]struct{}, len(tables))
	for _, t := range tables {
		tableSet[t] = struct{}{}
	}

	out := legacyExtract{
		TablesFound: tables,
		Settings:    make(map[string]any),
		Combos:      make([]contracts.ComboSnapshot, 0),
		Keys:        make([]legacyKeyRecord, 0),
		Warnings:    make([]string, 0),
	}

	if _, ok := tableSet["settings"]; ok {
		if err := readLegacySettings(ctx, conn, out.Settings); err != nil {
			out.Warnings = append(out.Warnings, "settings read skipped: "+err.Error())
		}
	}
	if _, ok := tableSet["combos"]; ok {
		combos, err := readLegacyCombos(ctx, conn)
		if err != nil {
			out.Warnings = append(out.Warnings, "combos read skipped: "+err.Error())
		} else {
			out.Combos = combos
		}
	}
	if _, ok := tableSet["apiKeys"]; ok {
		keys, err := readLegacyKeys(ctx, conn)
		if err != nil {
			out.Warnings = append(out.Warnings, "apiKeys read skipped: "+err.Error())
		} else {
			out.Keys = keys
		}
	}
	if _, ok := tableSet["providerConnections"]; ok {
		var cnt int
		if err := conn.QueryRowContext(ctx, `SELECT COUNT(*) FROM providerConnections`).Scan(&cnt); err == nil {
			out.ProvidersCount = cnt
			if cnt > 0 {
				out.Warnings = append(out.Warnings, "Legacy providerConnections are reported for audit only; re-link provider credentials via /dashboard/providers/ to encrypt with the new AES-256-GCM vault.")
			}
		}
	}
	return out, nil
}

func listTables(ctx context.Context, conn *sql.DB) ([]string, error) {
	rows, err := conn.QueryContext(ctx, `SELECT name FROM sqlite_master WHERE type = 'table' AND name NOT LIKE 'sqlite_%'`)
	if err != nil {
		return nil, fmt.Errorf("list legacy tables: %w", err)
	}
	defer rows.Close()

	var tables []string
	for rows.Next() {
		var name string
		if err := rows.Scan(&name); err != nil {
			return nil, err
		}
		tables = append(tables, name)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	sort.Strings(tables)
	return tables, nil
}

var legacySettingsKeyMap = map[string]string{
	"theme":                "ui.theme",
	"ui.theme":             "ui.theme",
	"locale":               "ui.locale",
	"ui.locale":            "ui.locale",
	"payloadLogging":       "usage.payload_logging_enabled",
	"payloadRetentionDays": "usage.payload_retention_days",
}

func readLegacySettings(ctx context.Context, conn *sql.DB, dst map[string]any) error {
	rows, err := conn.QueryContext(ctx, `SELECT key, value FROM settings`)
	if err != nil {
		return err
	}
	defer rows.Close()

	for rows.Next() {
		var k, rawVal string
		if err := rows.Scan(&k, &rawVal); err != nil {
			return err
		}
		mapped, ok := legacySettingsKeyMap[strings.TrimSpace(k)]
		if !ok {
			continue
		}
		var decoded any
		if err := json.Unmarshal([]byte(rawVal), &decoded); err != nil {
			decoded = strings.TrimSpace(rawVal)
		}
		dst[mapped] = decoded
	}
	return rows.Err()
}

func readLegacyCombos(ctx context.Context, conn *sql.DB) ([]contracts.ComboSnapshot, error) {
	rows, err := conn.QueryContext(ctx, `SELECT name, data FROM combos`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	out := make([]contracts.ComboSnapshot, 0)
	for rows.Next() {
		var name, dataJSON string
		if err := rows.Scan(&name, &dataJSON); err != nil {
			return nil, err
		}
		var payload struct {
			Description string   `json:"description"`
			Models      []string `json:"models"`
			Targets     []string `json:"targets"`
			Strategy    string   `json:"strategy"`
			IsActive    *bool    `json:"isActive"`
			Enabled     *bool    `json:"enabled"`
		}
		_ = json.Unmarshal([]byte(dataJSON), &payload)
		targets := payload.Models
		if len(targets) == 0 {
			targets = payload.Targets
		}
		if strings.TrimSpace(name) == "" || len(targets) == 0 {
			continue
		}
		enabled := true
		if payload.IsActive != nil {
			enabled = *payload.IsActive
		} else if payload.Enabled != nil {
			enabled = *payload.Enabled
		}
		strategy := strings.TrimSpace(payload.Strategy)
		if strategy == "" {
			strategy = "ordered"
		}
		out = append(out, contracts.ComboSnapshot{
			Name:        strings.TrimSpace(name),
			Description: strings.TrimSpace(payload.Description),
			Strategy:    strategy,
			Targets:     targets,
			Enabled:     enabled,
		})
	}
	return out, rows.Err()
}

func readLegacyKeys(ctx context.Context, conn *sql.DB) ([]legacyKeyRecord, error) {
	rows, err := conn.QueryContext(ctx, `SELECT name, data FROM apiKeys`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	out := make([]legacyKeyRecord, 0)
	for rows.Next() {
		var name, dataJSON string
		if err := rows.Scan(&name, &dataJSON); err != nil {
			return nil, err
		}
		var payload struct {
			AllowedModels []string `json:"allowedModels"`
			RateLimitRPM  int      `json:"rateLimitRpm"`
		}
		_ = json.Unmarshal([]byte(dataJSON), &payload)
		trimmedName := strings.TrimSpace(name)
		if trimmedName == "" {
			trimmedName = "legacy-imported-key"
		}
		out = append(out, legacyKeyRecord{
			Name:          trimmedName,
			AllowedModels: payload.AllowedModels,
			RateLimitRPM:  payload.RateLimitRPM,
		})
	}
	return out, rows.Err()
}
