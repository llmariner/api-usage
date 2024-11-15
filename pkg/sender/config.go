package sender

import (
	"fmt"
	"time"
)

const (
	defaultInitialDelay = 5 * time.Second
	defaultInterval     = 10 * time.Second

	defaultUsageChannelSize = 300

	// defaultMaxMessageSize is the default maximum size of a message.
	// This value considers the default Kafka message size limit.
	defaultMaxMessageSize = 768 * 1024
)

// Config is the configuration for the sender.
type Config struct {
	// Enable is the flag to enable the sender.
	Enable bool `yaml:"enable"`
	// APIUsageInternalServerAddr is the address of the server to send usage data to.
	APIUsageInternalServerAddr string `yaml:"apiUsageInternalServerAddr"`
	// InitialDelay is the time to wait before starting the sender.
	InitialDelay time.Duration `yaml:"initialDelay"`
	// Interval is the interval at which the sender sends usage data to the server.
	Interval time.Duration `yaml:"interval"`
	// UsageChannelSize is the size of the channel that stores usage data.
	UsageChannelSize int `yaml:"usageChannelSize"`
	// MaxMessageSize is the maximum size of a message that the sender can send one gRPC call.
	MaxMessageSize int `yaml:"maxMessageSize"`
}

// Validate validates the configuration.
func (c *Config) Validate() error {
	if !c.Enable {
		return nil
	}

	if c.APIUsageInternalServerAddr == "" {
		return fmt.Errorf("server address is required")
	}

	if c.InitialDelay == 0 {
		c.InitialDelay = defaultInitialDelay
	} else if c.InitialDelay < 0 {
		return fmt.Errorf("initial delay must be greater than 0")
	}
	if c.Interval == 0 {
		c.Interval = defaultInterval
	} else if c.Interval < 0 {
		return fmt.Errorf("interval must be greater than 0")
	}

	if c.UsageChannelSize == 0 {
		c.UsageChannelSize = defaultUsageChannelSize
	} else if c.UsageChannelSize < 0 {
		return fmt.Errorf("usage channel size must be greater than 0")
	}

	if c.MaxMessageSize == 0 {
		c.MaxMessageSize = defaultMaxMessageSize
	} else if c.MaxMessageSize < 0 {
		return fmt.Errorf("max message size must be greater than 0")
	}
	return nil
}
