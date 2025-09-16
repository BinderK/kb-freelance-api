package api

import (
	"log"

	"kb-freelance-api/internal/config"
	"kb-freelance-api/internal/services"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

type Server struct {
	config             *config.Config
	timeTrackerService *services.TimeTrackerService
	invoiceService     *services.InvoiceService
}

func NewServer(cfg *config.Config) *Server {
	return &Server{
		config:             cfg,
		timeTrackerService: services.NewTimeTrackerService(cfg),
		invoiceService:     services.NewInvoiceService(cfg),
	}
}

func (s *Server) Start(addr string) error {
	// Set Gin mode
	gin.SetMode(gin.ReleaseMode)

	// Create router
	router := gin.New()

	// Add middleware
	router.Use(gin.Logger())
	router.Use(gin.Recovery())
	router.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:3000", "http://127.0.0.1:3000", "http://localhost:5173", "http://127.0.0.1:5173"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Accept", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
	}))

	// Health check endpoints
	router.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok", "service": "kb-freelance-api"})
	})

	// API health check endpoint
	router.GET("/api/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"success": true, "data": gin.H{"status": "ok", "service": "kb-freelance-api", "timestamp": "2024-01-01T00:00:00Z"}})
	})

	// Static file serving for PDFs
	router.Static("/files", s.config.InvoiceGenPath+"/output")

	// API routes
	api := router.Group("/api")
	{
		// Time tracking routes
		time := api.Group("/time")
		{
			time.POST("/start", s.startTimer)
			time.POST("/stop", s.stopTimer)
			time.GET("/current", s.getTimerStatus) // Changed from /status to /current
			time.GET("/entries", s.getTimeEntries)
			time.GET("/today", s.getTodaySummary)
		}

		// Invoice routes
		invoice := api.Group("/invoice")
		{
			invoice.POST("/generate", s.generateInvoice)
			invoice.GET("/preview", s.previewInvoice)
		}
	}

	log.Printf("Server starting on %s", addr)
	return router.Run(addr)
}
