// Package main is the entrypoint for The Url Mentalist API.
//
//	@title			The Url Mentalist API
//	@version		1.0
//	@description	URL safety analysis API. Checks phishing signals, redirects, TLS, DNS, brand mismatches, and more.
//	@license.name	AGPL-3.0
//
//	@host		localhost:8080
//	@BasePath	/api/v1
//
//	@tag.name			Analysis
//	@tag.description	Full URL analysis
//	@tag.name			URL
//	@tag.description	URL structure checks
//	@tag.name			DNS
//	@tag.description	DNS and infrastructure checks
//	@tag.name			Security
//	@tag.description	Security protocol checks
//	@tag.name			Threat Intel
//	@tag.description	Threat intelligence lookups
//	@tag.name			Utility
//	@tag.description	Utility and administrative endpoints
package main

import (
	"os"
	"os/signal"
	"syscall"

	"github.com/AARYA-001/the-url-mentalist/internal/admintoken"
	_ "github.com/AARYA-001/the-url-mentalist/internal/docs" // swagger docs registration
	"github.com/AARYA-001/the-url-mentalist/internal/handler"
	"github.com/AARYA-001/the-url-mentalist/internal/logger"
	"github.com/AARYA-001/the-url-mentalist/internal/service/rank"
	"github.com/AARYA-001/the-url-mentalist/internal/service/screenshot"
	"github.com/AARYA-001/the-url-mentalist/internal/service/typosquat"
	"github.com/joho/godotenv"
)

func main() {
	// Load .env file if it exists (check multiple paths for local & docker compatibility)
	var loadedEnv bool
	for _, envPath := range []string{".env", "server/.env", "../.env", "/app/.env"} {
		if err := godotenv.Load(envPath); err == nil {
			loadedEnv = true
			break
		}
	}

	logger.Init()

	if !loadedEnv {
		logger.Info("no .env file found, using environment variables or defaults")
	}

	// Fail fast if required secrets are missing, before serving any requests.
	admintoken.Preload()

	// Initialize screenshot service (shared browser allocator)
	_, err := screenshot.GetService()
	if err != nil {
		logger.Warn("screenshot service init failed, continuing without it", "err", err)
	} else {
		logger.Info("screenshot service initialized")
		// Ensure cleanup on shutdown
		defer func() {
			service, _ := screenshot.GetService()
			if service != nil {
				service.Close()
			}
		}()
	}

	r := handler.SetupRouter()

	err = rank.LoadDomainRanks()
	if err != nil {
		logger.Fatal("failed to load domain ranks", "err", err)
	}

	if err := typosquat.LoadTopDomains(); err != nil {
		logger.Fatal("failed to load top domains", "err", err)
	}

	// Get port from environment variable, default to 8080
	port := getEnv("PORT", "8080")
	addr := ":" + port

	// Setup graceful shutdown
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)

	go func() {
		logger.Info("starting server", "addr", addr)
		if err := r.Run(addr); err != nil {
			logger.Fatal("server error", "err", err)
		}
	}()

	// Wait for interrupt signal
	<-sigChan
	logger.Info("shutting down server")

	// Cleanup screenshot service
	service, _ := screenshot.GetService()
	if service != nil {
		service.Close()
	}
}

// getEnv gets an environment variable or returns a default value
func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
