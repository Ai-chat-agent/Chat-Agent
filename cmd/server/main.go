package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"

	"github.com/Ai-chat-agent/Chat-Agent.git/internal/config"
	"github.com/Ai-chat-agent/Chat-Agent.git/internal/database"
	"github.com/Ai-chat-agent/Chat-Agent.git/pkg/handlers"
	"github.com/Ai-chat-agent/Chat-Agent.git/pkg/logger"
	"github.com/Ai-chat-agent/Chat-Agent.git/pkg/middleware"
)

func main() {
	// Load environment variables
	if err := godotenv.Load(); err != nil {
		// It's okay if .env file doesn't exist
		fmt.Println("No .env file found, using environment variables")
	}

	// Load configuration
	cfg, err := config.Load()
	if err != nil {
		fmt.Printf("Failed to load configuration: %v\n", err)
		os.Exit(1)
	}

	// Initialize logger
	log, err := logger.NewZapLogger(cfg.Log.Level, cfg.Log.Format)
	if err != nil {
		fmt.Printf("Failed to initialize logger: %v\n", err)
		os.Exit(1)
	}

	log.Info("Starting Chat Agent Server",
		logger.F("version", cfg.App.Version),
		logger.F("environment", cfg.App.Environment),
	)

	// Initialize database
	db, err := database.New(&cfg.Database, log)
	if err != nil {
		log.Fatal("Failed to connect to database", logger.F("error", err.Error()))
	}
	defer db.Close()

	// Run database migrations
	if err := db.AutoMigrate(); err != nil {
		log.Fatal("Failed to run database migrations", logger.F("error", err.Error()))
	}

	// Initialize repositories
	userRepo := database.NewUserRepository(db.DB, log)
	chatRepo := database.NewChatRepository(db.DB, log)

	// Initialize handlers
	healthHandler := handlers.NewHealthHandler()
	chatHandler := handlers.NewChatHandler(log)

	// Initialize Gin router
	if cfg.App.Environment == "production" {
		gin.SetMode(gin.ReleaseMode)
	}

	router := gin.New()

	// Add middleware
	router.Use(middleware.LoggerMiddleware(log))
	router.Use(middleware.RecoveryMiddleware(log))
	router.Use(middleware.CORSMiddleware())
	router.Use(middleware.SecurityHeadersMiddleware())
	router.Use(middleware.RequestIDMiddleware())

	// Health check routes
	router.GET("/health", healthHandler.Health)
	router.GET("/health/ready", healthHandler.Ready)
	router.GET("/health/live", healthHandler.Live)

	// API routes
	v1 := router.Group("/api/v1")
	{
		// Chat endpoints
		chat := v1.Group("/chat")
		{
			chat.POST("/message", chatHandler.SendMessage)
			chat.GET("/history/:userID", chatHandler.GetChatHistory)
			chat.DELETE("/message/:messageID", chatHandler.DeleteMessage)
		}

		// User endpoints (placeholder for future implementation)
		users := v1.Group("/users")
		{
			users.GET("/", func(c *gin.Context) {
				c.JSON(http.StatusOK, gin.H{"message": "Users endpoint - Coming Soon"})
			})
		}
	}

	// Welcome route
	router.GET("/", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message":     "Welcome to Chat Agent API",
			"version":     cfg.App.Version,
			"environment": cfg.App.Environment,
			"status":      "running",
			"endpoints": map[string]string{
				"health":   "/health",
				"api_docs": "/api/v1",
				"chat":     "/api/v1/chat",
			},
		})
	})

	// Create HTTP server
	server := &http.Server{
		Addr:         fmt.Sprintf("%s:%s", cfg.Server.Host, cfg.Server.Port),
		Handler:      router,
		ReadTimeout:  time.Duration(cfg.Server.ReadTimeout) * time.Second,
		WriteTimeout: time.Duration(cfg.Server.WriteTimeout) * time.Second,
		IdleTimeout:  time.Duration(cfg.Server.IdleTimeout) * time.Second,
	}

	// Start server in a goroutine
	go func() {
		log.Info("Server starting",
			logger.F("address", server.Addr),
			logger.F("environment", cfg.App.Environment),
		)

		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatal("Server failed to start", logger.F("error", err.Error()))
		}
	}()

	// Wait for interrupt signal to gracefully shutdown the server
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Info("Server shutting down...")

	// Give outstanding requests a deadline for completion
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// Attempt graceful shutdown
	if err := server.Shutdown(ctx); err != nil {
		log.Error("Server forced to shutdown", logger.F("error", err.Error()))
		os.Exit(1)
	}

	log.Info("Server exited gracefully")

	// Suppress unused variable warnings (remove these when implementing actual functionality)
	_ = userRepo
	_ = chatRepo
}
