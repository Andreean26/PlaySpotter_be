package routes

import (
	"playspotter/internal/config"
	"playspotter/internal/handlers"
	"playspotter/internal/middlewares"
	"playspotter/pkg/jwt"

	"github.com/gin-gonic/gin"
)

type Router struct {
	authHandler  *handlers.AuthHandler
	meHandler    *handlers.MeHandler
	eventHandler *handlers.EventHandler
	adminHandler *handlers.AdminHandler
	jwtManager   *jwt.Manager
	cfg          *config.Config
}

func NewRouter(
	authHandler *handlers.AuthHandler,
	meHandler *handlers.MeHandler,
	eventHandler *handlers.EventHandler,
	adminHandler *handlers.AdminHandler,
	jwtManager *jwt.Manager,
	cfg *config.Config,
) *Router {
	return &Router{
		authHandler:  authHandler,
		meHandler:    meHandler,
		eventHandler: eventHandler,
		adminHandler: adminHandler,
		jwtManager:   jwtManager,
		cfg:          cfg,
	}
}

func (r *Router) Setup(router *gin.Engine) {
	// CORS middleware
	router.Use(middlewares.CORS(r.cfg.AllowedOrigins))

	// Health check
	router.GET("/health", handlers.HealthCheck)

	// Auth routes (with rate limiting)
	authLimiter := middlewares.NewAuthRateLimiter()
	auth := router.Group("/auth")
	auth.Use(middlewares.RateLimitMiddleware(authLimiter))
	{
		auth.POST("/register", r.authHandler.Register)
		auth.POST("/login", r.authHandler.Login)
		auth.POST("/refresh", r.authHandler.RefreshToken)
		auth.POST("/logout", r.authHandler.Logout)
	}

	// Internal routes
	internal := router.Group("/internal")
	{
		internal.POST("/bootstrap-admin", func(c *gin.Context) {
			r.authHandler.BootstrapAdmin(c, r.cfg.AdminEmail, r.cfg.AdminPassword, r.cfg.AdminBootstrapToken)
		})
	}

	// Protected routes (require JWT)
	jwtAuth := middlewares.JWTAuth(r.jwtManager)

	// Me routes
	router.GET("/me", jwtAuth, r.meHandler.GetMe)
	router.PUT("/me", jwtAuth, r.meHandler.UpdateMe)

	// Event routes
	events := router.Group("/events")
	{
		// Public routes
		events.GET("", r.eventHandler.ListEvents)
		events.GET("/:id", r.eventHandler.GetEvent)

		// Protected routes
		events.POST("", jwtAuth, r.eventHandler.CreateEvent)
		events.PUT("/:id", jwtAuth, r.eventHandler.UpdateEvent)
		events.DELETE("/:id", jwtAuth, r.eventHandler.DeleteEvent)
		events.POST("/:id/join", jwtAuth, r.eventHandler.JoinEvent)
		events.POST("/:id/leave", jwtAuth, r.eventHandler.LeaveEvent)
		events.POST("/:id/swipe", jwtAuth, r.eventHandler.SwipeEvent)
	}

	// Admin routes (require admin role)
	admin := router.Group("/admin")
	admin.Use(jwtAuth, middlewares.RequireRole("admin"))
	{
		admin.GET("/users", r.adminHandler.ListUsers)
		admin.PUT("/users/:id/role", r.adminHandler.UpdateUserRole)
		admin.GET("/events", r.adminHandler.ListAllEvents)
		admin.PUT("/events/:id/status", r.adminHandler.UpdateEventStatus)
	}
}
