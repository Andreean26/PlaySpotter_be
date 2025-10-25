package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"playspotter/internal/config"
	"playspotter/internal/db"
	"playspotter/internal/handlers"
	"playspotter/internal/repositories"
	"playspotter/internal/routes"
	"playspotter/internal/services"
	"playspotter/pkg/jwt"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

// @title PlaySpotter API
// @version 1.0
// @description Backend API for PlaySpotter - Tinder for sports events
// @termsOfService http://swagger.io/terms/

// @contact.name API Support
// @contact.email support@playspotter.com

// @license.name MIT
// @license.url https://opensource.org/licenses/MIT

// @host localhost:8080
// @BasePath /

// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
// @description Type "Bearer" followed by a space and JWT token.

// @securityDefinitions.apikey SetupToken
// @in header
// @name X-Setup-Token
// @description Setup token for bootstrapping admin user.

func main() {
	// Load configuration
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	// Set Gin mode
	if cfg.Env == "production" {
		gin.SetMode(gin.ReleaseMode)
	}

	// Connect to database
	database, err := db.Connect(cfg.DatabaseURL)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	// Initialize repositories
	userRepo := repositories.NewUserRepository(database)
	eventRepo := repositories.NewEventRepository(database)
	participantRepo := repositories.NewParticipantRepository(database)
	swipeRepo := repositories.NewSwipeRepository(database)
	tokenRepo := repositories.NewTokenRepository(database)

	// Initialize JWT manager
	jwtManager := jwt.NewManager(
		cfg.JWTAccessSecret,
		cfg.JWTRefreshSecret,
		cfg.AccessTTL,
		cfg.RefreshTTL,
	)

	// Initialize services
	authService := services.NewAuthService(userRepo, tokenRepo, jwtManager)
	userService := services.NewUserService(userRepo)
	eventService := services.NewEventService(eventRepo, participantRepo)
	swipeService := services.NewSwipeService(swipeRepo)

	// Initialize handlers
	authHandler := handlers.NewAuthHandler(authService)
	meHandler := handlers.NewMeHandler(userService)
	eventHandler := handlers.NewEventHandler(eventService, swipeService)
	adminHandler := handlers.NewAdminHandler(userService, eventService)

	// Setup router
	router := gin.Default()

	// Swagger documentation
	router.GET("/docs/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// Setup routes
	apiRouter := routes.NewRouter(
		authHandler,
		meHandler,
		eventHandler,
		adminHandler,
		jwtManager,
		cfg,
	)
	apiRouter.Setup(router)

	// Create server
	srv := &http.Server{
		Addr:    ":" + cfg.Port,
		Handler: router,
	}

	// Start server in a goroutine
	go func() {
		log.Printf("Starting server on port %s", cfg.Port)
		log.Printf("Swagger docs available at http://localhost:%s/docs/index.html", cfg.Port)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Failed to start server: %v", err)
		}
	}()

	// Wait for interrupt signal for graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("Shutting down server...")

	// Graceful shutdown with 5 second timeout
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Fatalf("Server forced to shutdown: %v", err)
	}

	log.Println("Server exited")
}
