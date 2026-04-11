package app

import (
	"encoding/json"
	"os"
	"path/filepath"
)

// Config holds persistent user configuration.
type Config struct {
	DefaultOutputDir string   `json:"default_output_dir,omitempty"`
	Theme            string   `json:"theme"`
	RecentFiles      []string `json:"recent_files"`
	FavoritePresets   []string `json:"favorite_presets,omitempty"`
	HWAccel          string   `json:"hw_accel"`
	FFmpegPath       string   `json:"ffmpeg_path,omitempty"`
}

// DefaultConfig returns the default configuration.
func DefaultConfig() *Config {
	return &Config{
		Theme:       "dark",
		RecentFiles: []string{},
		HWAccel:     "auto",
	}
}

// LoadConfig reads config from disk, or returns defaults.
func LoadConfig() *Config {
	path := configFilePath()
	data, err := os.ReadFile(path)
	if err != nil {
		return DefaultConfig()
	}

	cfg := DefaultConfig()
	if err := json.Unmarshal(data, cfg); err != nil {
		return DefaultConfig()
	}
	return cfg
}

// Save writes config to disk.
func (c *Config) Save() error {
	dir := configDirPath()
	if err := os.MkdirAll(dir, 0755); err != nil {
		return err
	}

	data, err := json.MarshalIndent(c, "", "  ")
	if err != nil {
		return err
	}

	return os.WriteFile(configFilePath(), data, 0644)
}

// AddRecentFile adds a file to the recent files list (max 10, no dupes).
func (c *Config) AddRecentFile(path string) {
	// Remove if already present
	filtered := make([]string, 0, len(c.RecentFiles))
	for _, f := range c.RecentFiles {
		if f != path {
			filtered = append(filtered, f)
		}
	}

	// Prepend
	c.RecentFiles = append([]string{path}, filtered...)

	// Cap at 10
	if len(c.RecentFiles) > 10 {
		c.RecentFiles = c.RecentFiles[:10]
	}
}

func configDirPath() string {
	home, _ := os.UserHomeDir()
	return filepath.Join(home, ".config", "nano-ffmpeg")
}

func configFilePath() string {
	return filepath.Join(configDirPath(), "config.json")
}
