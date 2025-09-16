package config

import (
	"fmt"
	"os"
	"path/filepath"
)

type Config struct {
	TimeTrackerPath string
	InvoiceGenPath  string
	DatabasePath    string
	Port            string
	PythonExecPath  string
}

func Load() *Config {
	// Get the absolute path to the freelance_tools directory
	// This assumes the API is in freelance_tools/kb-freelance-api/
	currentDir, _ := os.Getwd()
	freelanceToolsDir := filepath.Dir(currentDir) // Go up one level from kb-freelance-api

	config := &Config{
		TimeTrackerPath: getEnv("TIME_TRACKER_PATH", filepath.Join(freelanceToolsDir, "kb-tt-cli")),
		InvoiceGenPath:  getEnv("INVOICE_GEN_PATH", filepath.Join(freelanceToolsDir, "kb-invoice-gen-cli")),
		DatabasePath:    getEnv("DATABASE_PATH", filepath.Join(os.Getenv("HOME"), ".kb-tt-cli", "time_tracker.db")),
		Port:            getEnv("PORT", "8080"),
		PythonExecPath:  getEnv("PYTHON_EXEC_PATH", "/Users/kevinbinder/anaconda3/envs/kb-freelance/bin/python"),
	}

	// Debug: log the paths
	fmt.Printf("DEBUG: Current directory: %s\n", currentDir)
	fmt.Printf("DEBUG: Freelance tools directory: %s\n", freelanceToolsDir)
	fmt.Printf("DEBUG: Time tracker path: %s\n", config.TimeTrackerPath)
	fmt.Printf("DEBUG: Invoice gen path: %s\n", config.InvoiceGenPath)
	fmt.Printf("DEBUG: Python executable: %s\n", config.PythonExecPath)

	return config
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
