package main

import (
	"kb-freelance-api/internal/api"
	"kb-freelance-api/internal/config"
	"log"
	"os"
)

func main() {
	// Load configuration
	cfg := config.Load()

	// Create API server
	server := api.NewServer(cfg)

	// Start server
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("Starting server on port %s", port)
	if err := server.Start(":" + port); err != nil {
		log.Fatal("Failed to start server:", err)
	}
}
