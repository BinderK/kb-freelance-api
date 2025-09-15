package config

import (
	"os"
	"path/filepath"
	"testing"
)

func TestLoad(t *testing.T) {
	// Test that Load returns a valid config
	cfg := Load()

	if cfg == nil {
		t.Fatal("Config should not be nil")
	}

	// Test that paths are set
	if cfg.TimeTrackerPath == "" {
		t.Error("TimeTrackerPath should not be empty")
	}

	if cfg.InvoiceGenPath == "" {
		t.Error("InvoiceGenPath should not be empty")
	}

	if cfg.DatabasePath == "" {
		t.Error("DatabasePath should not be empty")
	}

	// Test that paths exist
	if !filepath.IsAbs(cfg.TimeTrackerPath) {
		t.Error("TimeTrackerPath should be absolute")
	}

	if !filepath.IsAbs(cfg.InvoiceGenPath) {
		t.Error("InvoiceGenPath should be absolute")
	}

	if !filepath.IsAbs(cfg.DatabasePath) {
		t.Error("DatabasePath should be absolute")
	}

	// Test default port
	if cfg.Port == "" {
		t.Error("Port should have a default value")
	}
}

func TestGetEnv(t *testing.T) {
	// Test with existing environment variable
	os.Setenv("TEST_VAR", "test_value")
	defer os.Unsetenv("TEST_VAR")

	result := getEnv("TEST_VAR", "default")
	if result != "test_value" {
		t.Errorf("Expected 'test_value', got '%s'", result)
	}

	// Test with non-existing environment variable
	result = getEnv("NON_EXISTING_VAR", "default_value")
	if result != "default_value" {
		t.Errorf("Expected 'default_value', got '%s'", result)
	}
}

func TestConfigPaths(t *testing.T) {
	cfg := Load()

	// Test that the paths point to the correct directories
	expectedTimeTrackerPath := filepath.Join(filepath.Dir(cfg.TimeTrackerPath), "kb-tt-cli")
	if cfg.TimeTrackerPath != expectedTimeTrackerPath {
		t.Errorf("Expected TimeTrackerPath to be %s, got %s", expectedTimeTrackerPath, cfg.TimeTrackerPath)
	}

	expectedInvoiceGenPath := filepath.Join(filepath.Dir(cfg.InvoiceGenPath), "kb-invoice-gen-cli")
	if cfg.InvoiceGenPath != expectedInvoiceGenPath {
		t.Errorf("Expected InvoiceGenPath to be %s, got %s", expectedInvoiceGenPath, cfg.InvoiceGenPath)
	}
}
