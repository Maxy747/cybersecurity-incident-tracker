// Package main is the entry point for the CyberGuard API server.
//
// @title           CyberGuard Security Incident Management API
// @version         1.0
// @description     A cybersecurity incident management platform for SOC teams. Provides incident tracking, investigation notes, role-based access control, audit logging, and dashboard analytics.
// @contact.name    CyberGuard Team
// @host            localhost:8080
// @BasePath        /
// @securityDefinitions.apikey BearerAuth
// @in              header
// @name            Authorization
// @description     Enter "Bearer <token>" — obtain token via POST /api/auth/login
package main

import (
	"log"

	"github.com/gin-gonic/gin"

	_ "github.com/marwans200/cyberguard/docs" // swag-generated docs

	"github.com/marwans200/cyberguard/internal/config"
	"github.com/marwans200/cyberguard/internal/database"
	"github.com/marwans200/cyberguard/internal/middleware"
	"github.com/marwans200/cyberguard/internal/routes"
)

func main() {
	// Load configuration
	cfg := config.Load()

	// Wire JWT secret into middleware package
	middleware.JWTSecret = []byte(cfg.JWTSecret)

	// Connect to PostgreSQL (retries until ready)
	database.Connect(cfg)

	// Run idempotent schema migrations
	database.Migrate()

	// Configure Gin
	r := gin.Default()
	r.SetTrustedProxies(nil) // Disable proxy trust for security

	// Register all routes
	routes.Setup(r)

	addr := ":" + cfg.ServerPort
	log.Printf("[server] CyberGuard listening on %s", addr)
	if err := r.Run(addr); err != nil {
		log.Fatalf("[server] failed to start: %v", err)
	}
}
