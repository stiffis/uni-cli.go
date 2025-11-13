package config

import (
	"os"
	"path/filepath"
)

// Config holds application configuration
type Config struct {
	DatabasePath string
	DataDir      string
	Theme        Theme
}

// Theme defines the color scheme
type Theme struct {
	Primary   string
	Secondary string
	Success   string
	Warning   string
	Danger    string
	Info      string
	Muted     string
}

// DefaultTheme returns the default color theme
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

// Load loads or creates the default configuration
func Load() (*Config, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return nil, err
	}

	dataDir := filepath.Join(homeDir, ".unicli")

	// Create data directory if it doesn't exist
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
