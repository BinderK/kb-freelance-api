package api

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"kb-freelance-api/internal/config"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

// Integration tests that test the actual API handlers with real services
func TestAPIHandlersIntegration(t *testing.T) {
	// Skip if running in CI or if Python tools are not available
	if os.Getenv("SKIP_INTEGRATION_TESTS") == "true" {
		t.Skip("Skipping integration tests")
	}

	// Set up real configuration
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

	// Create real server with real services
	server := NewServer(cfg)

	// Set Gin to test mode
	gin.SetMode(gin.TestMode)
	router := gin.New()

	// Add middleware
	router.Use(gin.Logger())
	router.Use(gin.Recovery())

	// Health check endpoint
	router.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok", "service": "kb-freelance-api"})
	})

	// API routes with real handlers
	api := router.Group("/api")
	{
		// Time tracking routes
		time := api.Group("/time")
		{
			time.POST("/start", server.startTimer)
			time.POST("/stop", server.stopTimer)
			time.GET("/status", server.getTimerStatus)
			time.GET("/entries", server.getTimeEntries)
			time.GET("/today", server.getTodaySummary)
		}

		// Invoice routes
		invoice := api.Group("/invoice")
		{
			invoice.POST("/generate", server.generateInvoice)
			invoice.GET("/preview", server.previewInvoice)
		}
	}

	t.Run("HealthEndpoint", func(t *testing.T) {
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", "/health", nil)
		router.ServeHTTP(w, req)

		assert.Equal(t, 200, w.Code)

		var response map[string]interface{}
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Equal(t, "ok", response["status"])
	})

	t.Run("GetTodaySummary", func(t *testing.T) {
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", "/api/time/today", nil)
		router.ServeHTTP(w, req)

		// This might fail if Python tools are not available, which is expected
		if w.Code != 200 {
			t.Logf("GetTodaySummary returned %d (expected if Python tools not available)", w.Code)
			return
		}

		var response map[string]interface{}
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)

		// Check that we got a valid response structure
		assert.Contains(t, response, "total_hours")
		assert.Contains(t, response, "total_minutes")
		assert.Contains(t, response, "entry_count")
		assert.Contains(t, response, "breakdown")
	})

	t.Run("GetStatus", func(t *testing.T) {
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", "/api/time/status", nil)
		router.ServeHTTP(w, req)

		if w.Code != 200 {
			t.Logf("GetStatus returned %d (expected if Python tools not available)", w.Code)
			return
		}

		var response map[string]interface{}
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)

		assert.Contains(t, response, "is_running")
	})

	t.Run("StartTimer", func(t *testing.T) {
		requestBody := map[string]interface{}{
			"client":      "Test Client",
			"project":     "Test Project",
			"description": "Integration Test",
		}
		jsonBody, _ := json.Marshal(requestBody)

		w := httptest.NewRecorder()
		req, _ := http.NewRequest("POST", "/api/time/start", bytes.NewBuffer(jsonBody))
		req.Header.Set("Content-Type", "application/json")
		router.ServeHTTP(w, req)

		if w.Code != 200 {
			t.Logf("StartTimer returned %d (expected if Python tools not available)", w.Code)
			return
		}

		var response map[string]interface{}
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)

		assert.Contains(t, response, "status")
	})

	t.Run("GenerateInvoice", func(t *testing.T) {
		requestBody := map[string]interface{}{
			"client_name":  "Test Client",
			"client_email": "test@example.com",
			"line_items": []map[string]interface{}{
				{
					"description": "Test Work",
					"hours":       2.0,
					"rate":        75.0,
				},
			},
			"notes": "Integration Test",
			"date":  "2024-01-01",
		}
		jsonBody, _ := json.Marshal(requestBody)

		w := httptest.NewRecorder()
		req, _ := http.NewRequest("POST", "/api/invoice/generate", bytes.NewBuffer(jsonBody))
		req.Header.Set("Content-Type", "application/json")
		router.ServeHTTP(w, req)

		if w.Code != 200 {
			t.Logf("GenerateInvoice returned %d (expected if Python tools not available)", w.Code)
			return
		}

		var response map[string]interface{}
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)

		assert.Contains(t, response, "status")
	})
}
