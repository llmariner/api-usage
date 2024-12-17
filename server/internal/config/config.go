package config

import (
	"fmt"
	"os"
	"time"

	"github.com/llmariner/common/pkg/db"
	"gopkg.in/yaml.v3"
)

// CleanerConfig is the configuration for the cleaner.
type CleanerConfig struct {
	RetentionPeriod time.Duration `yaml:"retentionPeriod"`
	PollInterval    time.Duration `yaml:"pollInterval"`
}

// Validate validates the configuration.
func (c CleanerConfig) Validate() error {
	if c.RetentionPeriod <= 0 {
		return fmt.Errorf("retentionPeriod must be greater than 0")
	}
	if c.PollInterval <= 0 {
		return fmt.Errorf("pollInterval must be greater than 0")
	}
	return nil
}

// AuthConfig is the authentication configuration.
type AuthConfig struct {
	Enable                 bool   `yaml:"enable"`
	RBACInternalServerAddr string `yaml:"rbacInternalServerAddr"`
}

// Validate validates the configuration.
func (c *AuthConfig) validate() error {
	if !c.Enable {
		return nil
	}
	if c.RBACInternalServerAddr == "" {
		return fmt.Errorf("rbacInternalServerAddr must be set")
	}
	return nil
}

// Config is the configuration.
type Config struct {
	AdminGRPCPort    int           `yaml:"adminGrpcPort"`
	GRPCPort         int           `yaml:"grpcPort"`
	HTTPPort         int           `yaml:"httpPort"`
	InternalGRPCPort int           `yaml:"internalGrpcPort"`
	Cleaner          CleanerConfig `yaml:"cleaner"`

	AuthConfig AuthConfig `yaml:"auth"`

	Database db.Config `yaml:"database"`
}

// Validate validates the configuration.
func (c *Config) Validate() error {
	if c.AdminGRPCPort <= 0 {
		return fmt.Errorf("adminGrpcPort must be greater than 0")
	}
	if c.GRPCPort <= 0 {
		return fmt.Errorf("grpcPort must be greater than 0")
	}
	if c.HTTPPort <= 0 {
		return fmt.Errorf("httpPort must be greater than 0")
	}
	if c.InternalGRPCPort <= 0 {
		return fmt.Errorf("grpcPort must be greater than 0")
	}
	if err := c.Cleaner.Validate(); err != nil {
		return fmt.Errorf("cleaner: %s", err)
	}
	if err := c.AuthConfig.validate(); err != nil {
		return fmt.Errorf("auth: %s", err)
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
