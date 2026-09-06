package ai

import (
	"testing"
)

func TestNewRequiresModel(t *testing.T) {
	_, err := New(Config{})
	if err == nil {
		t.Fatal("expected error for empty model")
	}
}

func TestConfigWithTemperature(t *testing.T) {
	cfg := Config{Model: "test"}.WithTemperature(0.7)
	if cfg.Options["temperature"] != 0.7 {
		t.Fatalf("expected temperature 0.7, got %v", cfg.Options["temperature"])
	}
}

func TestNewInvalidHost(t *testing.T) {
	_, err := New(Config{Model: "test", Host: "://bad"})
	if err == nil {
		t.Fatal("expected error for invalid host")
	}
}
