package api

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"kb-freelance-api/internal/services"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// Mock services for testing
type MockTimeTrackerService struct {
	mock.Mock
}

func (m *MockTimeTrackerService) StartTimer(client, project, description string) (map[string]interface{}, error) {
	args := m.Called(client, project, description)
	return args.Get(0).(map[string]interface{}), args.Error(1)
}

func (m *MockTimeTrackerService) StopTimer() (map[string]interface{}, error) {
	args := m.Called()
	return args.Get(0).(map[string]interface{}), args.Error(1)
}

func (m *MockTimeTrackerService) GetStatus() (*services.TimerStatus, error) {
	args := m.Called()
	return args.Get(0).(*services.TimerStatus), args.Error(1)
}

func (m *MockTimeTrackerService) GetRecentEntries(limit int) ([]services.TimeEntry, error) {
	args := m.Called(limit)
	return args.Get(0).([]services.TimeEntry), args.Error(1)
}

func (m *MockTimeTrackerService) GetTodaySummary() (*services.TodaySummary, error) {
	args := m.Called()
	return args.Get(0).(*services.TodaySummary), args.Error(1)
}

type MockInvoiceService struct {
	mock.Mock
}

func (m *MockInvoiceService) GenerateInvoice(clientName, clientEmail string, lineItems []services.InvoiceLineItem, notes, date string) (map[string]interface{}, error) {
	args := m.Called(clientName, clientEmail, lineItems, notes, date)
	return args.Get(0).(map[string]interface{}), args.Error(1)
}

func setupTestRouter() (*gin.Engine, *MockTimeTrackerService, *MockInvoiceService) {
	gin.SetMode(gin.TestMode)

	mockTimeTracker := &MockTimeTrackerService{}
	mockInvoice := &MockInvoiceService{}

	router := gin.New()

	// Add middleware
	router.Use(gin.Logger())
	router.Use(gin.Recovery())

	// Health check endpoint
	router.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok", "service": "kb-freelance-api"})
	})

	// API routes
	api := router.Group("/api")
	{
		// Time tracking routes
		time := api.Group("/time")
		{
			time.POST("/start", func(c *gin.Context) {
				var req StartTimerRequest
				if err := c.ShouldBindJSON(&req); err != nil {
					c.JSON(400, gin.H{"error": err.Error()})
					return
				}

				result, err := mockTimeTracker.StartTimer(req.Client, req.Project, req.Description)
				if err != nil {
					c.JSON(500, gin.H{"error": err.Error()})
					return
				}

				c.JSON(200, result)
			})

			time.POST("/stop", func(c *gin.Context) {
				result, err := mockTimeTracker.StopTimer()
				if err != nil {
					c.JSON(500, gin.H{"error": err.Error()})
					return
				}

				c.JSON(200, result)
			})

			time.GET("/status", func(c *gin.Context) {
				status, err := mockTimeTracker.GetStatus()
				if err != nil {
					c.JSON(500, gin.H{"error": err.Error()})
					return
				}

				c.JSON(200, status)
			})

			time.GET("/entries", func(c *gin.Context) {
				limit := 10
				if limitStr := c.Query("limit"); limitStr != "" {
					// Parse limit parameter
					limit = 10 // Simplified for test
				}

				entries, err := mockTimeTracker.GetRecentEntries(limit)
				if err != nil {
					c.JSON(500, gin.H{"error": err.Error()})
					return
				}

				c.JSON(200, gin.H{"entries": entries})
			})

			time.GET("/today", func(c *gin.Context) {
				summary, err := mockTimeTracker.GetTodaySummary()
				if err != nil {
					c.JSON(500, gin.H{"error": err.Error()})
					return
				}

				c.JSON(200, summary)
			})
		}

		// Invoice routes
		invoice := api.Group("/invoice")
		{
			invoice.POST("/generate", func(c *gin.Context) {
				var req GenerateInvoiceRequest
				if err := c.ShouldBindJSON(&req); err != nil {
					c.JSON(400, gin.H{"error": err.Error()})
					return
				}

				// Convert request line items to service line items
				lineItems := make([]services.InvoiceLineItem, len(req.LineItems))
				for i, item := range req.LineItems {
					lineItems[i] = services.InvoiceLineItem{
						Description: item.Description,
						Hours:       item.Hours,
						Rate:        item.Rate,
					}
				}

				result, err := mockInvoice.GenerateInvoice(req.ClientName, req.ClientEmail, lineItems, req.Notes, req.Date)
				if err != nil {
					c.JSON(500, gin.H{"error": err.Error()})
					return
				}

				c.JSON(200, result)
			})
		}
	}

	return router, mockTimeTracker, mockInvoice
}

func TestHealthEndpoint(t *testing.T) {
	router, _, _ := setupTestRouter()

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/health", nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, 200, w.Code)

	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, "ok", response["status"])
	assert.Equal(t, "kb-freelance-api", response["service"])
}

func TestStartTimerEndpoint(t *testing.T) {
	router, mockTimeTracker, _ := setupTestRouter()

	// Mock the service call
	expectedResult := map[string]interface{}{
		"status":      "success",
		"message":     "Started tracking time for Test Client/Test Project",
		"client":      "Test Client",
		"project":     "Test Project",
		"description": "Test Description",
	}
	mockTimeTracker.On("StartTimer", "Test Client", "Test Project", "Test Description").Return(expectedResult, nil)

	// Create request
	requestBody := StartTimerRequest{
		Client:      "Test Client",
		Project:     "Test Project",
		Description: "Test Description",
	}
	jsonBody, _ := json.Marshal(requestBody)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/api/time/start", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)

	assert.Equal(t, 200, w.Code)

	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, "success", response["status"])

	mockTimeTracker.AssertExpectations(t)
}

func TestStopTimerEndpoint(t *testing.T) {
	router, mockTimeTracker, _ := setupTestRouter()

	// Mock the service call
	expectedResult := map[string]interface{}{
		"status":  "success",
		"message": "Timer stopped successfully",
	}
	mockTimeTracker.On("StopTimer").Return(expectedResult, nil)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/api/time/stop", nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, 200, w.Code)

	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, "success", response["status"])

	mockTimeTracker.AssertExpectations(t)
}

func TestGetStatusEndpoint(t *testing.T) {
	router, mockTimeTracker, _ := setupTestRouter()

	// Mock the service call
	expectedStatus := &services.TimerStatus{
		IsRunning:       true,
		Client:          "Test Client",
		Project:         "Test Project",
		Description:     "Test Description",
		DurationMinutes: 30,
	}
	mockTimeTracker.On("GetStatus").Return(expectedStatus, nil)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/api/time/status", nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, 200, w.Code)

	var response services.TimerStatus
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.True(t, response.IsRunning)
	assert.Equal(t, "Test Client", response.Client)
	assert.Equal(t, "Test Project", response.Project)

	mockTimeTracker.AssertExpectations(t)
}

func TestGetTodaySummaryEndpoint(t *testing.T) {
	router, mockTimeTracker, _ := setupTestRouter()

	// Mock the service call
	expectedSummary := &services.TodaySummary{
		TotalHours:   2.5,
		TotalMinutes: 150,
		EntryCount:   2,
		Breakdown: []services.Breakdown{
			{
				ClientProject: "Client A/Project A",
				Hours:         1.5,
				Minutes:       90,
			},
			{
				ClientProject: "Client B/Project B",
				Hours:         1.0,
				Minutes:       60,
			},
		},
	}
	mockTimeTracker.On("GetTodaySummary").Return(expectedSummary, nil)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/api/time/today", nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, 200, w.Code)

	var response services.TodaySummary
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, 2.5, response.TotalHours)
	assert.Equal(t, 150, response.TotalMinutes)
	assert.Equal(t, 2, response.EntryCount)
	assert.Len(t, response.Breakdown, 2)

	mockTimeTracker.AssertExpectations(t)
}

func TestGenerateInvoiceEndpoint(t *testing.T) {
	router, _, mockInvoice := setupTestRouter()

	// Mock the service call
	expectedResult := map[string]interface{}{
		"status":  "success",
		"message": "Invoice generated successfully",
	}

	lineItems := []services.InvoiceLineItem{
		{
			Description: "Test Work",
			Hours:       2.0,
			Rate:        75.0,
		},
	}

	mockInvoice.On("GenerateInvoice", "Test Client", "test@example.com", lineItems, "Test notes", "2024-01-01").Return(expectedResult, nil)

	// Create request
	requestBody := GenerateInvoiceRequest{
		ClientName:  "Test Client",
		ClientEmail: "test@example.com",
		LineItems: []InvoiceLineItemRequest{
			{
				Description: "Test Work",
				Hours:       2.0,
				Rate:        75.0,
			},
		},
		Notes: "Test notes",
		Date:  "2024-01-01",
	}
	jsonBody, _ := json.Marshal(requestBody)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/api/invoice/generate", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)

	assert.Equal(t, 200, w.Code)

	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, "success", response["status"])

	mockInvoice.AssertExpectations(t)
}

func TestInvalidJSONRequest(t *testing.T) {
	router, _, _ := setupTestRouter()

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/api/time/start", bytes.NewBufferString("invalid json"))
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)

	assert.Equal(t, 400, w.Code)

	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Contains(t, response["error"], "invalid character")
}
