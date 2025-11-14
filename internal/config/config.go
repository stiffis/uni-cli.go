package config

import (
	"os"
	"path/filepath"
)

type Config struct {
	DatabasePath string
	DataDir      string
	Theme        Theme
}

type Theme struct {
	Primary   string
	Secondary string
	Success   string
	Warning   string
	Danger    string
	Info      string
	Muted     string
}

func DefaultTheme() Theme {
	return Theme{
		Primary:   "#7C3AED",
		Secondary: "#06B6D4",
		Success:   "#10B981",
		Warning:   "#F59E0B",
		Danger:    "#EF4444",
		Info:      "#3B82F6",
		Muted:     "#6B7280",
	}
}

func Load() (*Config, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return nil, err
	}

	dataDir := filepath.Join(homeDir, ".unicli")

	if err := os.MkdirAll(dataDir, 0755); err != nil {
		return nil, err
	}

	cfg := &Config{
		DatabasePath: filepath.Join(dataDir, "unicli.db"),
		DataDir:      dataDir,
		Theme:        DefaultTheme(),
	}

	return cfg, nil
}
