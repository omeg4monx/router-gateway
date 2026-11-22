// Package config loads and validates runtime configuration from environment
// variables and command-line flags.
package config

import (
	"errors"
	"flag"
	"fmt"
	"io"
	"net"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"
)

const (
	defaultHost       = "127.0.0.1"
	defaultPort       = 20128
	defaultLogLevel   = "info"
	defaultAPITimeout = 30 * time.Second
)

// Config holds the resolved runtime settings for omeg4mon-router.
type Config struct {
	Host                string
	Port                int
	DataDir             string
	MasterKeyHex        string
	AdminPassword       string
	DisabledModules     map[string]bool
	AllowInsecureRemote bool
	LogLevel            string
	APITimeout          time.Duration
}

// CoreModules returns the names of core modules that cannot be disabled.
func CoreModules() []string {
	return []string{
		"system",
		"auth",
		"settings",
		"providers",
		"keys",
		"combos",
		"gateway",
		"usage",
	}
}

// Load parses configuration from environment variables (via getenv) and CLI
// arguments, applying defaults and validating constraints.
func Load(args []string, getenv func(string) string) (Config, error) {
	if getenv == nil {
		getenv = os.Getenv
	}

	cfg := Config{
		Host:       defaultHost,
		Port:       defaultPort,
		DataDir:    resolveDefaultDataDir(getenv),
		LogLevel:   defaultLogLevel,
		APITimeout: defaultAPITimeout,
	}

	if v := strings.TrimSpace(getenv("OMEG4MON_HOST")); v != "" {
		cfg.Host = v
	}
	if v := strings.TrimSpace(getenv("OMEG4MON_PORT")); v != "" {
		port, err := strconv.Atoi(v)
		if err != nil {
			return Config{}, fmt.Errorf("invalid OMEG4MON_PORT %q: %w", v, err)
		}
		cfg.Port = port
	}
	if v := strings.TrimSpace(getenv("OMEG4MON_DATA_DIR")); v != "" {
		cfg.DataDir = v
	}
	cfg.MasterKeyHex = strings.TrimSpace(getenv("OMEG4MON_MASTER_KEY"))
	cfg.AdminPassword = getenv("OMEG4MON_ADMIN_PASSWORD")
	if v := strings.TrimSpace(getenv("OMEG4MON_LOG_LEVEL")); v != "" {
		cfg.LogLevel = v
	}
	if v := strings.TrimSpace(getenv("OMEG4MON_ALLOW_INSECURE_REMOTE")); v != "" {
		allow, err := strconv.ParseBool(v)
		if err != nil {
			return Config{}, fmt.Errorf("invalid OMEG4MON_ALLOW_INSECURE_REMOTE %q: %w", v, err)
		}
		cfg.AllowInsecureRemote = allow
	}
	if v := strings.TrimSpace(getenv("OMEG4MON_API_TIMEOUT")); v != "" {
		d, err := time.ParseDuration(v)
		if err != nil || d <= 0 {
			return Config{}, fmt.Errorf("invalid OMEG4MON_API_TIMEOUT %q: must be a positive duration", v)
		}
		cfg.APITimeout = d
	}

	disabled, err := parseDisabledModules(getenv("OMEG4MON_DISABLED_MODULES"))
	if err != nil {
		return Config{}, err
	}
	cfg.DisabledModules = disabled

	fs := flag.NewFlagSet("omeg4mon-router", flag.ContinueOnError)
	fs.SetOutput(io.Discard)
	fs.StringVar(&cfg.Host, "host", cfg.Host, "listen host")
	fs.IntVar(&cfg.Port, "port", cfg.Port, "listen port")
	fs.StringVar(&cfg.DataDir, "data-dir", cfg.DataDir, "data directory path")
	fs.StringVar(&cfg.LogLevel, "log-level", cfg.LogLevel, "log level")
	if err := fs.Parse(args); err != nil {
		return Config{}, err
	}

	cfg.Host = strings.TrimSpace(cfg.Host)
	cfg.DataDir = strings.TrimSpace(cfg.DataDir)
	cfg.LogLevel = strings.TrimSpace(cfg.LogLevel)

	if cfg.Host == "" {
		return Config{}, errors.New("host must not be empty")
	}
	if cfg.Port < 1 || cfg.Port > 65535 {
		return Config{}, fmt.Errorf("invalid port %d: must be between 1 and 65535", cfg.Port)
	}
	if cfg.DataDir == "" {
		return Config{}, errors.New("data-dir must not be empty")
	}

	return cfg, nil
}

// Addr returns the formatted host:port listen address.
func (c Config) Addr() string {
	return net.JoinHostPort(c.Host, strconv.Itoa(c.Port))
}

// IsLoopback reports whether Host resolves to a loopback address or localhost.
func (c Config) IsLoopback() bool {
	host := strings.Trim(strings.TrimSpace(c.Host), "[]")
	if strings.EqualFold(host, "localhost") {
		return true
	}
	if ip := net.ParseIP(host); ip != nil {
		return ip.IsLoopback()
	}
	return false
}

// IsModuleDisabled reports whether the named module is in DisabledModules.
func (c Config) IsModuleDisabled(name string) bool {
	if len(c.DisabledModules) == 0 {
		return false
	}
	return c.DisabledModules[strings.ToLower(strings.TrimSpace(name))]
}

func resolveDefaultDataDir(getenv func(string) string) string {
	if home := strings.TrimSpace(getenv("HOME")); home != "" {
		return filepath.Join(home, ".omeg4mon-router")
	}
	if home, err := os.UserHomeDir(); err == nil && home != "" {
		return filepath.Join(home, ".omeg4mon-router")
	}
	return ".omeg4mon-router"
}

func parseDisabledModules(raw string) (map[string]bool, error) {
	disabled := make(map[string]bool)
	if strings.TrimSpace(raw) == "" {
		return disabled, nil
	}

	coreSet := make(map[string]bool, len(CoreModules()))
	for _, m := range CoreModules() {
		coreSet[m] = true
	}

	for _, part := range strings.Split(raw, ",") {
		name := strings.ToLower(strings.TrimSpace(part))
		if name == "" {
			continue
		}
		if coreSet[name] {
			return nil, fmt.Errorf("core module %q cannot be disabled", name)
		}
		disabled[name] = true
	}

	return disabled, nil
}
