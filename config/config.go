package config

import (
	"fmt"
	"os"
	"time"
)

type Config struct {
	Database struct {
		URL string
	}
	RefreshInterval time.Duration
}

// LoadConfig loads the configuration from environment variables
func LoadConfig() (*Config, error) {
	config := &Config{}

	// Load database URL from environment
	if dbURL := os.Getenv("RIVER_DB_URL"); dbURL != "" {
		config.Database.URL = dbURL
	}

	// Load timeout from environment
	if timeoutStr := os.Getenv("RIVER_CLI_TIMEOUT"); timeoutStr != "" {
		timeout, err := time.ParseDuration(timeoutStr)
		if err != nil {
			return nil, fmt.Errorf("invalid RIVER_CLI_TIMEOUT value: %w", err)
		}
		config.RefreshInterval = timeout
	} else {
		config.RefreshInterval = 1 * time.Second
	}

	return config, nil
}

// UpdateConfigFromFlags updates the configuration with values from command-line flags
func UpdateConfigFromFlags(config *Config, dbURL string, refreshInterval time.Duration) {
	if dbURL != "" {
		config.Database.URL = dbURL
	}
	if refreshInterval != 0 {
		config.RefreshInterval = refreshInterval
	}
}
