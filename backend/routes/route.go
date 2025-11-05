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
	router := gin.Default()

	corsConfig := cors.Config{
		AllowOrigins:     getAllowedOrigins(),
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS", "PATCH"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}

	router.Use(cors.New(corsConfig))

	router.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status":  "healthy",
			"service": "SlotSwapper API",
		})
	})

	auth := router.Group("/api/auth")
	{
		auth.POST("/signup", controllers.Signup)
		auth.POST("/login", controllers.Login)
	}
	protected := router.Group("/api")
	protected.Use(middleware.AuthMiddleware())
	{
		// Event routes
		protected.GET("/events", controllers.GetEvents)
		protected.POST("/events", controllers.CreateEvent)
		protected.PATCH("/events/:id/status", controllers.UpdateEventStatus)
		protected.DELETE("/events/:id", controllers.DeleteEvent)

		// Swap routes
		protected.GET("/swappable-slots", controllers.GetSwappableSlots)
		protected.POST("/swap-request", controllers.InitiateSwap)
		protected.GET("/swap-requests", controllers.GetMySwapRequests)
		protected.POST("/swap-request/:id/accept", controllers.AcceptSwap)
		protected.POST("/swap-request/:id/reject", controllers.RejectSwap)
		protected.DELETE("/swap-request/:id", controllers.CancelSwap)
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
