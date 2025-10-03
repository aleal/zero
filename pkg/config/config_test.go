package config

import (
	"os"
	"testing"
	"time"
)

func TestDefault(t *testing.T) {
	cfg := Default()

	if cfg.Host != "localhost" {
		t.Errorf("Default host = %s, want localhost", cfg.Host)
	}

	if cfg.Port != 8000 {
		t.Errorf("Default port = %d, want 8000", cfg.Port)
	}

	if cfg.ReadTimeout != 5*time.Second {
		t.Errorf("Default read timeout = %v, want 5s", cfg.ReadTimeout)
	}

	if cfg.WriteTimeout != 15*time.Second {
		t.Errorf("Default write timeout = %v, want 15s", cfg.WriteTimeout)
	}

	if cfg.IdleTimeout != 60*time.Second {
		t.Errorf("Default idle timeout = %v, want 60s", cfg.IdleTimeout)
	}

	if cfg.RateLimit != 100 {
		t.Errorf("Default rate limit = %d, want 100", cfg.RateLimit)
	}

	if !cfg.EnableLogging {
		t.Error("Default enable logging = false, want true")
	}

	if !cfg.EnableCORS {
		t.Error("Default enable CORS = false, want true")
	}

	if !cfg.EnableRecovery {
		t.Error("Default enable recovery = false, want true")
	}

	if len(cfg.AllowedOrigins) != 1 || cfg.AllowedOrigins[0] != "*" {
		t.Errorf("Default allowed origins = %v, want [*]", cfg.AllowedOrigins)
	}
}

func TestFromEnv(t *testing.T) {
	// Set environment variables
	os.Setenv("ZERO_HOST", "0.0.0.0")
	os.Setenv("ZERO_PORT", "9000")
	os.Setenv("ZERO_READ_TIMEOUT", "10s")
	os.Setenv("ZERO_WRITE_TIMEOUT", "20s")
	os.Setenv("ZERO_IDLE_TIMEOUT", "120s")
	os.Setenv("ZERO_RATE_LIMIT", "200")
	os.Setenv("ZERO_ENABLE_LOGGING", "false")
	os.Setenv("ZERO_ENABLE_CORS", "false")
	os.Setenv("ZERO_ENABLE_RECOVERY", "false")

	defer func() {
		os.Unsetenv("ZERO_HOST")
		os.Unsetenv("ZERO_PORT")
		os.Unsetenv("ZERO_READ_TIMEOUT")
		os.Unsetenv("ZERO_WRITE_TIMEOUT")
		os.Unsetenv("ZERO_IDLE_TIMEOUT")
		os.Unsetenv("ZERO_RATE_LIMIT")
		os.Unsetenv("ZERO_ENABLE_LOGGING")
		os.Unsetenv("ZERO_ENABLE_CORS")
		os.Unsetenv("ZERO_ENABLE_RECOVERY")
	}()

	cfg := FromEnv()

	if cfg.Host != "0.0.0.0" {
		t.Errorf("FromEnv host = %s, want 0.0.0.0", cfg.Host)
	}

	if cfg.Port != 9000 {
		t.Errorf("FromEnv port = %d, want 9000", cfg.Port)
	}

	if cfg.ReadTimeout != 10*time.Second {
		t.Errorf("FromEnv read timeout = %v, want 10s", cfg.ReadTimeout)
	}

	if cfg.WriteTimeout != 20*time.Second {
		t.Errorf("FromEnv write timeout = %v, want 20s", cfg.WriteTimeout)
	}

	if cfg.IdleTimeout != 120*time.Second {
		t.Errorf("FromEnv idle timeout = %v, want 120s", cfg.IdleTimeout)
	}

	if cfg.RateLimit != 200 {
		t.Errorf("FromEnv rate limit = %d, want 200", cfg.RateLimit)
	}

	if cfg.EnableLogging {
		t.Error("FromEnv enable logging = true, want false")
	}

	if cfg.EnableCORS {
		t.Error("FromEnv enable CORS = true, want false")
	}

	if cfg.EnableRecovery {
		t.Error("FromEnv enable recovery = true, want false")
	}
}

func TestGetAddr(t *testing.T) {
	cfg := &Config{
		Host: "localhost",
		Port: 8080,
	}

	addr := cfg.GetAddr()
	expected := "localhost:8080"

	if addr != expected {
		t.Errorf("GetAddr() = %s, want %s", addr, expected)
	}
}

func TestSetHost(t *testing.T) {
	cfg := &Config{}
	cfg.SetHost("0.0.0.0")

	if cfg.Host != "0.0.0.0" {
		t.Errorf("SetHost() = %s, want 0.0.0.0", cfg.Host)
	}
}

func TestSetPort(t *testing.T) {
	cfg := &Config{}
	cfg.SetPort(9000)

	if cfg.Port != 9000 {
		t.Errorf("SetPort() = %d, want 9000", cfg.Port)
	}
}

func TestSetReadTimeout(t *testing.T) {
	cfg := &Config{}
	timeout := 10 * time.Second
	cfg.SetReadTimeout(timeout)

	if cfg.ReadTimeout != timeout {
		t.Errorf("SetReadTimeout() = %v, want %v", cfg.ReadTimeout, timeout)
	}
}

func TestSetWriteTimeout(t *testing.T) {
	cfg := &Config{}
	timeout := 20 * time.Second
	cfg.SetWriteTimeout(timeout)

	if cfg.WriteTimeout != timeout {
		t.Errorf("SetWriteTimeout() = %v, want %v", cfg.WriteTimeout, timeout)
	}
}

func TestSetIdleTimeout(t *testing.T) {
	cfg := &Config{}
	timeout := 120 * time.Second
	cfg.SetIdleTimeout(timeout)

	if cfg.IdleTimeout != timeout {
		t.Errorf("SetIdleTimeout() = %v, want %v", cfg.IdleTimeout, timeout)
	}
}

func TestSetAllowedOrigins(t *testing.T) {
	cfg := &Config{}
	origins := []string{"https://example.com", "https://test.com"}
	cfg.SetAllowedOrigins(origins)

	if len(cfg.AllowedOrigins) != len(origins) {
		t.Errorf("SetAllowedOrigins() length = %d, want %d", len(cfg.AllowedOrigins), len(origins))
	}

	for i, origin := range origins {
		if cfg.AllowedOrigins[i] != origin {
			t.Errorf("SetAllowedOrigins()[%d] = %s, want %s", i, cfg.AllowedOrigins[i], origin)
		}
	}
}

func TestSetRateLimit(t *testing.T) {
	cfg := &Config{}
	cfg.SetRateLimit(200)

	if cfg.RateLimit != 200 {
		t.Errorf("SetRateLimit() = %d, want 200", cfg.RateLimit)
	}
}

func BenchmarkDefault(b *testing.B) {
	for i := 0; i < b.N; i++ {
		Default()
	}
}

func BenchmarkFromEnv(b *testing.B) {
	for i := 0; i < b.N; i++ {
		FromEnv()
	}
}
