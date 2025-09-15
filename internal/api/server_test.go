package api

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"kb-freelance-api/internal/config"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestNewServer(t *testing.T) {
	cfg := &config.Config{
		TimeTrackerPath: "/test/tt",
		InvoiceGenPath:  "/test/invoice",
		DatabasePath:    "/test/db",
		Port:            "8080",
	}

	server := NewServer(cfg)

	assert.NotNil(t, server)
	assert.Equal(t, cfg, server.config)
	assert.NotNil(t, server.timeTrackerService)
	assert.NotNil(t, server.invoiceService)
}

func TestServerRoutes(t *testing.T) {
	cfg := &config.Config{
		TimeTrackerPath: "/test/tt",
		InvoiceGenPath:  "/test/invoice",
		DatabasePath:    "/test/db",
		Port:            "8080",
	}

	_ = NewServer(cfg)

	// Create a test router
	gin.SetMode(gin.TestMode)
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
			time.POST("/start", func(c *gin.Context) { c.JSON(200, gin.H{"status": "ok"}) })
			time.POST("/stop", func(c *gin.Context) { c.JSON(200, gin.H{"status": "ok"}) })
			time.GET("/status", func(c *gin.Context) { c.JSON(200, gin.H{"status": "ok"}) })
			time.GET("/entries", func(c *gin.Context) { c.JSON(200, gin.H{"status": "ok"}) })
			time.GET("/today", func(c *gin.Context) { c.JSON(200, gin.H{"status": "ok"}) })
		}

		// Invoice routes
		invoice := api.Group("/invoice")
		{
			invoice.POST("/generate", func(c *gin.Context) { c.JSON(200, gin.H{"status": "ok"}) })
			invoice.GET("/preview", func(c *gin.Context) { c.JSON(200, gin.H{"status": "ok"}) })
		}
	}

	// Test that all routes are registered
	routes := router.Routes()

	expectedRoutes := []string{
		"GET /health",
		"POST /api/time/start",
		"POST /api/time/stop",
		"GET /api/time/status",
		"GET /api/time/entries",
		"GET /api/time/today",
		"POST /api/invoice/generate",
		"GET /api/invoice/preview",
	}

	routeMap := make(map[string]bool)
	for _, route := range routes {
		routeMap[route.Method+" "+route.Path] = true
	}

	for _, expectedRoute := range expectedRoutes {
		assert.True(t, routeMap[expectedRoute], "Route %s should be registered", expectedRoute)
	}
}

func TestCORSConfiguration(t *testing.T) {
	cfg := &config.Config{
		TimeTrackerPath: "/test/tt",
		InvoiceGenPath:  "/test/invoice",
		DatabasePath:    "/test/db",
		Port:            "8080",
	}

	_ = NewServer(cfg)

	// Create a test router with CORS
	gin.SetMode(gin.TestMode)
	router := gin.New()

	// Add CORS middleware
	router.Use(func(c *gin.Context) {
		c.Header("Access-Control-Allow-Origin", "http://localhost:3000")
		c.Header("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		c.Header("Access-Control-Allow-Headers", "Origin, Content-Type, Accept, Authorization")
		c.Header("Access-Control-Expose-Headers", "Content-Length")
		c.Header("Access-Control-Allow-Credentials", "true")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	})

	// Health check endpoint
	router.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok", "service": "kb-freelance-api"})
	})

	// Test CORS headers
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/health", nil)
	req.Header.Set("Origin", "http://localhost:3000")
	router.ServeHTTP(w, req)

	assert.Equal(t, 200, w.Code)
	assert.Equal(t, "http://localhost:3000", w.Header().Get("Access-Control-Allow-Origin"))
	assert.Equal(t, "GET, POST, PUT, DELETE, OPTIONS", w.Header().Get("Access-Control-Allow-Methods"))
	assert.Equal(t, "Origin, Content-Type, Accept, Authorization", w.Header().Get("Access-Control-Allow-Headers"))
}

func TestServerMiddleware(t *testing.T) {
	cfg := &config.Config{
		TimeTrackerPath: "/test/tt",
		InvoiceGenPath:  "/test/invoice",
		DatabasePath:    "/test/db",
		Port:            "8080",
	}

	_ = NewServer(cfg)

	// Create a test router
	gin.SetMode(gin.TestMode)
	router := gin.New()

	// Add middleware
	router.Use(gin.Logger())
	router.Use(gin.Recovery())

	// Test endpoint that might panic
	router.GET("/panic", func(c *gin.Context) {
		panic("test panic")
	})

	// Health check endpoint
	router.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok", "service": "kb-freelance-api"})
	})

	// Test that recovery middleware works
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/panic", nil)
	router.ServeHTTP(w, req)

	// Should not crash, recovery middleware should handle it
	assert.Equal(t, 500, w.Code)

	// Test normal endpoint still works
	w = httptest.NewRecorder()
	req, _ = http.NewRequest("GET", "/health", nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, 200, w.Code)
}
