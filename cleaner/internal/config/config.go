package config

import (
	"fmt"
	"os"
	"time"

	"github.com/llmariner/common/pkg/db"
	"gopkg.in/yaml.v3"
)

// Config is the configuration for the cleaner.
type Config struct {
	// RetentionPeriod is the duration for which records are kept.
	RetentionPeriod time.Duration `yaml:"retentionPeriod"`
	// PollInterval is the interval at which cleaner runs.
	PollInterval time.Duration `yaml:"pollInterval"`
	// Database is the database configuration.
	Database db.Config `yaml:"database"`
}

// Validate validates the configuration.
func (c *Config) Validate() error {
	if c.RetentionPeriod <= 0 {
		return fmt.Errorf("retentionPeriod must be greater than 0")
	}
	if c.PollInterval <= 0 {
		return fmt.Errorf("pollInterval must be greater than 0")
	}
	if err := c.Database.Validate(); err != nil {
		return fmt.Errorf("database: %s", err)
	}
	return nil
}

// Parse parses the configuration file at the given path, returning a new
// Config struct.
func Parse(path string) (*Config, error) {
	b, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("config: read: %s", err)
	}
	var config Config
	if err = yaml.Unmarshal(b, &config); err != nil {
		return nil, fmt.Errorf("config: unmarshal: %s", err)
	}
	return &config, nil
}
