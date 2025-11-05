package routes

import (
	"os"
	"strings"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"

	"SlotSwapper/controllers"
	"SlotSwapper/middleware"
)

func SetupRoutes() *gin.Engine {
	// Set Gin mode based on environment
	if os.Getenv("APP_ENV") == "production" {
		gin.SetMode(gin.ReleaseMode)
	}

	router := gin.Default()

	// CORS configuration
	corsConfig := cors.Config{
		AllowOrigins:     getAllowedOrigins(),
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS", "PATCH"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		ExposeHeaders:    []string{"Content-Length", "X-RateLimit-Limit", "X-RateLimit-Remaining"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}
	router.Use(cors.New(corsConfig))

	// Global rate limiter (100 req/min per IP)
	router.Use(middleware.GlobalRateLimiter())

	// Health check endpoints (no auth, no additional rate limit)
	health := router.Group("/health")
	{
		health.GET("", controllers.HealthCheck)
		health.GET("/detailed", controllers.DetailedHealthCheck)
		health.GET("/ready", controllers.ReadinessCheck)
		health.GET("/live", controllers.LivenessCheck)
	}

	// Auth routes with strict rate limiting (5 req/min per IP)
	auth := router.Group("/api/auth")
	auth.Use(middleware.AuthRateLimiter())
	{
		auth.POST("/signup", controllers.Signup)
		auth.POST("/login", controllers.Login)
	}

	// Protected routes
	protected := router.Group("/api")
	protected.Use(middleware.AuthMiddleware())
	{
		// User routes
		protected.GET("/users/me", controllers.GetCurrentUser)
		protected.PATCH("/users/me", controllers.UpdateUserProfile)

		// Event routes
		protected.GET("/events", controllers.GetEvents)
		protected.POST("/events", controllers.CreateEvent)
		protected.PATCH("/events/:id/status", controllers.UpdateEventStatus)
		protected.DELETE("/events/:id", controllers.DeleteEvent)

		// Swap routes with moderate rate limiting (20 req/min per IP)
		swaps := protected.Group("/")
		swaps.Use(middleware.SwapRateLimiter())
		{
			swaps.GET("/swappable-slots", controllers.GetSwappableSlots)
			swaps.GET("/swap-requests", controllers.GetMySwapRequests)
			swaps.POST("/swap-request", controllers.InitiateSwap)
			swaps.POST("/swap-request/:id/accept", controllers.AcceptSwap)
			swaps.POST("/swap-request/:id/reject", controllers.RejectSwap)
			swaps.DELETE("/swap-request/:id", controllers.CancelSwap)
		}
	}

	return router
}

func getAllowedOrigins() []string {
	origins := os.Getenv("CORS_ALLOWED_ORIGINS")
	if origins == "" {
		return []string{"http://localhost:5173", "http://localhost:3000"}
	}
	return strings.Split(origins, ",")
}
