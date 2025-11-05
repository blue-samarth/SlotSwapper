package routes

import (
	"os"
	"time"
	"strings"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"

	"SlotSwapper/controllers"
	"SlotSwapper/middleware"
)

func SetupRoutes() *gin.Engine {
	router := gin.Default()

	corsConfig := cors.Config{
		AllowOrigins:     getAllowedOrigins(),
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
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
		// protected.GET("/events", controllers.GetEvents)
		// protected.POST("/events", controllers.CreateEvent)
		// protected.PUT("/events/:id", controllers.UpdateEvent)
		// protected.DELETE("/events/:id", controllers.DeleteEvent)
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