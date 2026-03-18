package nri

import (
	"fmt"
	"os"

	"github.com/pelletier/go-toml/v2"
)

// Config holds the nono-nri plugin configuration loaded from a TOML file.
type Config struct {
	RuntimeClasses []string `toml:"runtime_classes"`
	DefaultProfile string   `toml:"default_profile"`
	NonoBinPath    string   `toml:"nono_bin_path"`
	SocketPath     string   `toml:"socket_path"`
}

// LoadConfig reads and parses a TOML config file at the given path.
// Returns an error if the file cannot be read, fails to parse, or required fields are invalid.
// Unknown TOML keys are silently ignored (go-toml/v2 default behaviour — intentional).
func LoadConfig(path string) (*Config, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("reading config: %w", err)
	}
	var cfg Config
	if err := toml.Unmarshal(data, &cfg); err != nil {
		return nil, fmt.Errorf("parsing config: %w", err)
	}
	if len(cfg.RuntimeClasses) == 0 {
		return nil, fmt.Errorf("config: runtime_classes must not be empty")
	}
	if cfg.NonoBinPath == "" {
		return nil, fmt.Errorf("config: nono_bin_path must not be empty")
	}
	return &cfg, nil
}

// DefaultConfig returns a Config with sensible defaults.
func DefaultConfig() *Config {
	return &Config{
		DefaultProfile: "default",
		SocketPath:     "/var/run/nri/nri.sock",
	}
}
