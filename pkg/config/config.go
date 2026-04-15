// Package config provides configuration management for the Zero HTTP server.
// It supports loading configuration from environment variables with sensible defaults.
// Options (WithConfig/WithPort/etc) override everything.
package config

import (
	"log/slog"
	"os"
	"strconv"
	"time"

	"github.com/aleal/zero/pkg/parser"
)

// Config holds the server configuration
type Config struct {
	Host                string
	Port                int
	ReadTimeout         time.Duration
	WriteTimeout        time.Duration
	IdleTimeout         time.Duration
	MaxJSONBodySize     int64
	MaxUploadedFileSize int64
}

// Load loads configuration from environment variables with sensible defaults.
// Unparseable env values are logged at WARN level and the default is used.
func Load() *Config {
	config := &Config{
		Host:            "localhost",
		Port:            8000,
		ReadTimeout:     5 * time.Second,
		WriteTimeout:    15 * time.Second,
		IdleTimeout:     60 * time.Second,
		MaxJSONBodySize:     1 << 20, // 1 MB
		MaxUploadedFileSize: 10 << 20, // 10 MB
	}
	parseEnv(config)
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

// SetMaxJSONBodySize sets the maximum JSON body size in bytes
func (c *Config) SetMaxJSONBodySize(size int64) {
	c.MaxJSONBodySize = size
}

// SetMaxUploadedFileSize sets the maximum uploaded file size in bytes
func (c *Config) SetMaxUploadedFileSize(size int64) {
	c.MaxUploadedFileSize = size
}

func parseEnv(config *Config) {
	config.Host = getValueFromEnvOrDefault("ZERO_HOST", config.Host)
	config.Port = getValueFromEnvOrDefault("ZERO_PORT", config.Port)
	config.ReadTimeout = getValueFromEnvOrDefault("ZERO_READ_TIMEOUT", config.ReadTimeout)
	config.WriteTimeout = getValueFromEnvOrDefault("ZERO_WRITE_TIMEOUT", config.WriteTimeout)
	config.IdleTimeout = getValueFromEnvOrDefault("ZERO_IDLE_TIMEOUT", config.IdleTimeout)
	config.MaxJSONBodySize = getValueFromEnvOrDefault("ZERO_MAX_JSON_REQUEST_BODY_SIZE", config.MaxJSONBodySize)
	config.MaxUploadedFileSize = getValueFromEnvOrDefault("ZERO_MAX_UPLOADED_FILE_SIZE", config.MaxUploadedFileSize)
}

func getValueFromEnvOrDefault[T parser.ParseType](key string, defaultValue T) T {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	var parsedValue T
	if err := parser.ParseString(value, &parsedValue); err != nil {
		slog.Warn("invalid env value, using default",
			slog.String("key", key),
			slog.String("value", value),
			slog.Any("default", defaultValue),
			slog.Any("error", err),
		)
		return defaultValue
	}
	return parsedValue
}
