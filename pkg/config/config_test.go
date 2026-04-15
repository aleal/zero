package config

import (
	"testing"
	"time"
)

func TestLoad(t *testing.T) {
	cfg := Load()

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
}

func TestFromEnv(t *testing.T) {
	t.Setenv("ZERO_HOST", "0.0.0.0")
	t.Setenv("ZERO_PORT", "9000")
	t.Setenv("ZERO_READ_TIMEOUT", "10s")
	t.Setenv("ZERO_WRITE_TIMEOUT", "20s")
	t.Setenv("ZERO_IDLE_TIMEOUT", "120s")

	cfg := Load()

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
}

func TestLoadWithInvalidEnv(t *testing.T) {
	t.Setenv("ZERO_PORT", "not-a-number")

	// Should not panic — logs warning and uses default
	cfg := Load()
	if cfg.Port != 8000 {
		t.Errorf("Expected default port 8000 for invalid env, got %d", cfg.Port)
	}
}

func TestGetAddr(t *testing.T) {
	cfg := &Config{Host: "localhost", Port: 8080}

	addr := cfg.GetAddr()
	if addr != "localhost:8080" {
		t.Errorf("GetAddr() = %s, want localhost:8080", addr)
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

func BenchmarkLoad(b *testing.B) {
	for i := 0; i < b.N; i++ {
		Load()
	}
}
