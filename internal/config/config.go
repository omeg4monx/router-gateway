package config

import (
	"os"
	"path/filepath"
)

type Config struct {
	Host            string
	Port            string
	DataDir         string
	DBPath          string
	InitialPassword string
	PublicBaseUrl   string
	Version         string
}

func Load() *Config {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	// Loopback by default: this server has no auth; set HOST=0.0.0.0 only behind one.
	host := os.Getenv("HOST")
	if host == "" {
		host = "127.0.0.1"
	}

	dataDir := os.Getenv("DATA_DIR")
	if dataDir == "" {
		home, err := os.UserHomeDir()
		if err != nil {
			home = "/root"
		}
		// Default to shared xrouter data directory if present
		candidate := filepath.Join(home, "xrouter-remake", "data")
		if _, err := os.Stat(candidate); err == nil {
			dataDir = candidate
		} else {
			dataDir = filepath.Join(home, ".xrouter", "data")
		}
	}

	dbPath := filepath.Join(dataDir, "db", "data.sqlite")

	initialPassword := os.Getenv("INITIAL_PASSWORD")
	if initialPassword == "" {
		initialPassword = "xrouter123"
	}

	publicBaseUrl := os.Getenv("PUBLIC_BASE_URL")
	if publicBaseUrl == "" {
		publicBaseUrl = os.Getenv("BASE_URL")
	}

	return &Config{
		Host:            host,
		Port:            port,
		DataDir:         dataDir,
		DBPath:          dbPath,
		InitialPassword: initialPassword,
		PublicBaseUrl:   publicBaseUrl,
		Version:         "0.0.1",
	}
}
