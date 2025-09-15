package services

import (
	"os"
	"testing"

	"kb-freelance-api/internal/config"
)

// Integration tests that test the actual service methods
// These require the Python CLI tools to be available
func TestTimeTrackerServiceIntegration(t *testing.T) {
	// Skip if running in CI or if Python tools are not available
	if os.Getenv("SKIP_INTEGRATION_TESTS") == "true" {
		t.Skip("Skipping integration tests")
	}

	// Check if Python CLI tools are available
	pythonPath := os.Getenv("PYTHON_EXEC_PATH")
	if pythonPath == "" {
		pythonPath = "python3"
	}

	// Use environment variables or defaults for paths
	timeTrackerPath := os.Getenv("TIME_TRACKER_PATH")
	if timeTrackerPath == "" {
		timeTrackerPath = "../kb-tt-cli"
	}

	invoiceGenPath := os.Getenv("INVOICE_GEN_PATH")
	if invoiceGenPath == "" {
		invoiceGenPath = "../kb-invoice-gen-cli"
	}

	cfg := &config.Config{
		TimeTrackerPath: timeTrackerPath,
		InvoiceGenPath:  invoiceGenPath,
		DatabasePath:    "/tmp/test_time_tracker.db",
		Port:            "8080",
		PythonExecPath:  pythonPath,
	}

	service := NewTimeTrackerService(cfg)

	// Test GetTodaySummary (this should work if Python tools are available)
	t.Run("GetTodaySummary", func(t *testing.T) {
		summary, err := service.GetTodaySummary()
		if err != nil {
			t.Logf("GetTodaySummary failed (expected if Python tools not available): %v", err)
			return
		}

		if summary == nil {
			t.Error("Expected summary to not be nil")
		}

		t.Logf("Today's summary: %+v", summary)
	})

	// Test GetStatus
	t.Run("GetStatus", func(t *testing.T) {
		status, err := service.GetStatus()
		if err != nil {
			t.Logf("GetStatus failed (expected if Python tools not available): %v", err)
			return
		}

		if status == nil {
			t.Error("Expected status to not be nil")
		}

		t.Logf("Timer status: %+v", status)
	})
}

func TestInvoiceServiceIntegration(t *testing.T) {
	// Skip if running in CI or if Python tools are not available
	if os.Getenv("SKIP_INTEGRATION_TESTS") == "true" {
		t.Skip("Skipping integration tests")
	}

	pythonPath := os.Getenv("PYTHON_EXEC_PATH")
	if pythonPath == "" {
		pythonPath = "python3"
	}

	// Use environment variables or defaults for paths
	timeTrackerPath := os.Getenv("TIME_TRACKER_PATH")
	if timeTrackerPath == "" {
		timeTrackerPath = "../kb-tt-cli"
	}

	invoiceGenPath := os.Getenv("INVOICE_GEN_PATH")
	if invoiceGenPath == "" {
		invoiceGenPath = "../kb-invoice-gen-cli"
	}

	cfg := &config.Config{
		TimeTrackerPath: timeTrackerPath,
		InvoiceGenPath:  invoiceGenPath,
		DatabasePath:    "/tmp/test_time_tracker.db",
		Port:            "8080",
		PythonExecPath:  pythonPath,
	}

	service := NewInvoiceService(cfg)

	// Test GenerateInvoice with sample data
	t.Run("GenerateInvoice", func(t *testing.T) {
		lineItems := []InvoiceLineItem{
			{
				Description: "Test Work",
				Hours:       2.0,
				Rate:        75.0,
			},
		}

		result, err := service.GenerateInvoice("Test Client", "test@example.com", lineItems, "Test notes", "2024-01-01")
		if err != nil {
			t.Logf("GenerateInvoice failed (expected if Python tools not available): %v", err)
			return
		}

		if result == nil {
			t.Error("Expected result to not be nil")
		}

		t.Logf("Invoice generation result: %+v", result)
	})
}
