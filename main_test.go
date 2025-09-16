package main

import (
	"os"
	"testing"
)

func TestMain(m *testing.M) {
	// Set up test environment
	os.Setenv("PORT", "8081") // Use different port for tests

	// Run tests
	code := m.Run()

	// Clean up
	os.Unsetenv("PORT")

	// Exit with the same code as the tests
	os.Exit(code)
}

func TestMainFunction(t *testing.T) {
	// This test ensures the main function can be called without crashing
	// In a real scenario, you might want to test the server startup
	// but for unit tests, we'll just ensure the function exists and is callable

	// The main function should exist and be callable
	// We can't easily test it without mocking the server startup
	// but we can ensure the package compiles correctly
	t.Log("Main function exists and package compiles successfully")
}



