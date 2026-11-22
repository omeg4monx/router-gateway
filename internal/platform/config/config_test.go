package config

import (
	"path/filepath"
	"strings"
	"testing"
	"time"
)

func TestLoadDefaults(t *testing.T) {
	cfg, err := Load(nil, func(string) string { return "" })
	if err != nil {
		t.Fatalf("Load() unexpected error: %v", err)
	}
	if cfg.Host != "127.0.0.1" {
		t.Fatalf("Host = %q, want %q", cfg.Host, "127.0.0.1")
	}
	if cfg.Port != 20128 {
		t.Fatalf("Port = %d, want %d", cfg.Port, 20128)
	}
	if cfg.Addr() != "127.0.0.1:20128" {
		t.Fatalf("Addr() = %q, want %q", cfg.Addr(), "127.0.0.1:20128")
	}
	if !cfg.IsLoopback() {
		t.Fatalf("IsLoopback() = false, want true")
	}
	if cfg.LogLevel != "info" {
		t.Fatalf("LogLevel = %q, want %q", cfg.LogLevel, "info")
	}
	if cfg.AllowInsecureRemote {
		t.Fatalf("AllowInsecureRemote = true, want false")
	}
	if cfg.APITimeout <= 0 {
		t.Fatalf("APITimeout = %v, want > 0", cfg.APITimeout)
	}
	if cfg.IsModuleDisabled("mitm") {
		t.Fatalf("IsModuleDisabled(mitm) = true, want false")
	}

	withHome, err := Load(nil, func(key string) string {
		if key == "HOME" {
			return "/home/testuser"
		}
		return ""
	})
	if err != nil {
		t.Fatalf("Load(with HOME) unexpected error: %v", err)
	}
	wantDataDir := filepath.Join("/home/testuser", ".omeg4mon-router")
	if withHome.DataDir != wantDataDir {
		t.Fatalf("DataDir = %q, want %q", withHome.DataDir, wantDataDir)
	}
}

func TestLoadRejectsInvalidPortAndCoreModuleDisable(t *testing.T) {
	if _, err := Load(nil, func(key string) string {
		if key == "OMEG4MON_PORT" {
			return "99999"
		}
		return ""
	}); err == nil {
		t.Fatalf("expected error for OMEG4MON_PORT=99999, got nil")
	}

	if _, err := Load(nil, func(key string) string {
		if key == "OMEG4MON_PORT" {
			return "0"
		}
		return ""
	}); err == nil {
		t.Fatalf("expected error for OMEG4MON_PORT=0, got nil")
	}

	if _, err := Load([]string{"--port", "70000"}, func(string) string { return "" }); err == nil {
		t.Fatalf("expected error for --port 70000, got nil")
	}

	_, err := Load(nil, func(key string) string {
		if key == "OMEG4MON_DISABLED_MODULES" {
			return "mitm,gateway"
		}
		return ""
	})
	if err == nil {
		t.Fatalf("expected error when disabling core module gateway, got nil")
	}
	if !strings.Contains(err.Error(), "gateway") {
		t.Fatalf("expected error to mention core module \"gateway\", got %q", err.Error())
	}

	for _, core := range CoreModules() {
		_, err := Load(nil, func(key string) string {
			if key == "OMEG4MON_DISABLED_MODULES" {
				return core
			}
			return ""
		})
		if err == nil {
			t.Fatalf("expected error when disabling core module %q, got nil", core)
		}
	}

	cfg, err := Load(nil, func(key string) string {
		if key == "OMEG4MON_DISABLED_MODULES" {
			return "mitm, telegram"
		}
		return ""
	})
	if err != nil {
		t.Fatalf("unexpected error for optional disabled modules: %v", err)
	}
	if !cfg.IsModuleDisabled("mitm") {
		t.Fatalf("IsModuleDisabled(\"mitm\") = false, want true")
	}
	if !cfg.IsModuleDisabled("telegram") {
		t.Fatalf("IsModuleDisabled(\"telegram\") = false, want true")
	}
	if cfg.IsModuleDisabled("gateway") {
		t.Fatalf("IsModuleDisabled(\"gateway\") = true, want false")
	}
}

func TestIsLoopbackClassification(t *testing.T) {
	cases := []struct {
		host string
		want bool
	}{
		{host: "127.0.0.1", want: true},
		{host: "::1", want: true},
		{host: "localhost", want: true},
		{host: "LOCALHOST", want: true},
		{host: "0.0.0.0", want: false},
		{host: "::", want: false},
		{host: "192.168.1.10", want: false},
		{host: "example.com", want: false},
		{host: "", want: false},
	}

	for _, tc := range cases {
		cfg := Config{Host: tc.host}
		if got := cfg.IsLoopback(); got != tc.want {
			t.Errorf("Config{Host: %q}.IsLoopback() = %v, want %v", tc.host, got, tc.want)
		}
	}
}

func TestLoadEnvAndFlagOverrides(t *testing.T) {
	env := map[string]string{
		"OMEG4MON_HOST":                  "0.0.0.0",
		"OMEG4MON_PORT":                  "3000",
		"OMEG4MON_DATA_DIR":              "/var/lib/omeg4mon",
		"OMEG4MON_MASTER_KEY":            "0123456789abcdef",
		"OMEG4MON_ADMIN_PASSWORD":        "s3cr3t",
		"OMEG4MON_ALLOW_INSECURE_REMOTE": "true",
		"OMEG4MON_LOG_LEVEL":             "debug",
		"OMEG4MON_API_TIMEOUT":           "45s",
	}
	getenv := func(k string) string { return env[k] }

	cfg, err := Load([]string{"--host", "::1", "--port", "20200", "--data-dir", "/tmp/data", "--log-level", "warn"}, getenv)
	if err != nil {
		t.Fatalf("Load() unexpected error: %v", err)
	}
	if cfg.Host != "::1" {
		t.Fatalf("Host = %q, want %q", cfg.Host, "::1")
	}
	if cfg.Port != 20200 {
		t.Fatalf("Port = %d, want %d", cfg.Port, 20200)
	}
	if cfg.Addr() != "[::1]:20200" {
		t.Fatalf("Addr() = %q, want %q", cfg.Addr(), "[::1]:20200")
	}
	if cfg.DataDir != "/tmp/data" {
		t.Fatalf("DataDir = %q, want %q", cfg.DataDir, "/tmp/data")
	}
	if cfg.MasterKeyHex != "0123456789abcdef" {
		t.Fatalf("MasterKeyHex = %q, want %q", cfg.MasterKeyHex, "0123456789abcdef")
	}
	if cfg.AdminPassword != "s3cr3t" {
		t.Fatalf("AdminPassword = %q, want %q", cfg.AdminPassword, "s3cr3t")
	}
	if !cfg.AllowInsecureRemote {
		t.Fatalf("AllowInsecureRemote = false, want true")
	}
	if cfg.LogLevel != "warn" {
		t.Fatalf("LogLevel = %q, want %q", cfg.LogLevel, "warn")
	}
	if cfg.APITimeout != 45*time.Second {
		t.Fatalf("APITimeout = %v, want %v", cfg.APITimeout, 45*time.Second)
	}
}
