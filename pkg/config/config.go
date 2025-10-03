// Package config provides configuration management for the Zero HTTP server.
// It supports loading configuration from environment variables and provides
// sensible defaults for server settings.
package config

import (
	"os"
	"strconv"
	"time"
)

// Config holds the server configuration
type Config struct {
	Host           string
	Port           int
	ReadTimeout    time.Duration
	WriteTimeout   time.Duration
	IdleTimeout    time.Duration
	AllowedOrigins []string
	RateLimit      int
	EnableLogging  bool
	EnableCORS     bool
	EnableRecovery bool
}

// Default returns a default configuration
func Default() *Config {
	return &Config{
		Host:           "localhost",
		Port:           8000,
		ReadTimeout:    5 * time.Second,
		WriteTimeout:   15 * time.Second,
		IdleTimeout:    60 * time.Second,
		AllowedOrigins: []string{"*"},
		RateLimit:      100,
		EnableLogging:  true,
		EnableCORS:     true,
		EnableRecovery: true,
	}
}

// FromEnv loads configuration from environment variables
func FromEnv() *Config {
	config := Default()

	if host := os.Getenv("ZERO_HOST"); host != "" {
		config.Host = host
	}

	if portStr := os.Getenv("ZERO_PORT"); portStr != "" {
		if port, err := strconv.Atoi(portStr); err == nil {
			config.Port = port
		}
	}

	if readTimeoutStr := os.Getenv("ZERO_READ_TIMEOUT"); readTimeoutStr != "" {
		if readTimeout, err := time.ParseDuration(readTimeoutStr); err == nil {
			config.ReadTimeout = readTimeout
		}
	}

	if writeTimeoutStr := os.Getenv("ZERO_WRITE_TIMEOUT"); writeTimeoutStr != "" {
		if writeTimeout, err := time.ParseDuration(writeTimeoutStr); err == nil {
			config.WriteTimeout = writeTimeout
		}
	}

	if idleTimeoutStr := os.Getenv("ZERO_IDLE_TIMEOUT"); idleTimeoutStr != "" {
		if idleTimeout, err := time.ParseDuration(idleTimeoutStr); err == nil {
			config.IdleTimeout = idleTimeout
		}
	}

	if rateLimitStr := os.Getenv("ZERO_RATE_LIMIT"); rateLimitStr != "" {
		if rateLimit, err := strconv.Atoi(rateLimitStr); err == nil {
			config.RateLimit = rateLimit
		}
	}

	if enableLoggingStr := os.Getenv("ZERO_ENABLE_LOGGING"); enableLoggingStr != "" {
		config.EnableLogging = enableLoggingStr == "true"
	}

	if enableCORSStr := os.Getenv("ZERO_ENABLE_CORS"); enableCORSStr != "" {
		config.EnableCORS = enableCORSStr == "true"
	}

	if enableRecoveryStr := os.Getenv("ZERO_ENABLE_RECOVERY"); enableRecoveryStr != "" {
		config.EnableRecovery = enableRecoveryStr == "true"
	}

	return config
}

// GetAddr returns the server address as a string
func (c *Config) GetAddr() string {
	return c.Host + ":" + strconv.Itoa(c.Port)
}

// SetHost sets the host
func (c *Config) SetHost(host string) {
	c.Host = host
}

// SetPort sets the port
func (c *Config) SetPort(port int) {
	c.Port = port
}

// SetReadTimeout sets the read timeout
func (c *Config) SetReadTimeout(timeout time.Duration) {
	c.ReadTimeout = timeout
}

// SetWriteTimeout sets the write timeout
func (c *Config) SetWriteTimeout(timeout time.Duration) {
	c.WriteTimeout = timeout
}

// SetIdleTimeout sets the idle timeout
func (c *Config) SetIdleTimeout(timeout time.Duration) {
	c.IdleTimeout = timeout
}

// SetAllowedOrigins sets the allowed CORS origins
func (c *Config) SetAllowedOrigins(origins []string) {
	c.AllowedOrigins = origins
}

// SetRateLimit sets the rate limit
func (c *Config) SetRateLimit(limit int) {
	c.RateLimit = limit
}
