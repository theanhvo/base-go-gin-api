package routes

import (
	"time"

	"baseApi/dto"
	"baseApi/handlers"
	"baseApi/middleware"

	"github.com/gin-gonic/gin"
)

/* SetupRoutes configures all API routes */
func SetupRoutes() *gin.Engine {
	// Set Gin to release mode in production
	// gin.SetMode(gin.ReleaseMode)

	router := gin.New()

	// Apply global middleware
	router.Use(middleware.RecoveryWithSentry()) // Custom recovery with Sentry
	router.Use(middleware.CORSMiddleware())
	router.Use(middleware.SentryMiddleware())        // Sentry error tracking and performance
	router.Use(middleware.LoggingMiddleware())       // Request logging
	router.Use(middleware.CaptureErrorMiddleware())  // Capture Gin errors

	// Health check endpoint with standardized response
	router.GET("/health", func(c *gin.Context) {
		healthData := gin.H{
			"status":    "healthy",
			"timestamp": time.Now().Format(time.RFC3339),
			"uptime":    "running",
			"version":   "v1.0.0",
			"services": gin.H{
				"database": "connected",
				"redis":    "connected",
			},
		}
		
		response := dto.SuccessResponse(
		dto.StatusOK,
		"Service is healthy",
		healthData,
	)
		c.JSON(response.StatusCode, response)
	})

	// API v1 routes
	v1 := router.Group("/v1")
	{
		setupUserRoutes(v1)
	}

	return router
}

/* setupUserRoutes configures user-related routes */
func setupUserRoutes(rg *gin.RouterGroup) {
	userHandler := handlers.NewUserHandler()

	users := rg.Group("/users")
	{
		users.POST("", userHandler.CreateUser)                          // POST /api/v1/users
		users.GET("", userHandler.GetAllUsers)                          // GET /api/v1/users?page=1&limit=10
		users.GET("/:id", userHandler.GetUser)                          // GET /api/v1/users/1
		users.GET("/username/:username", userHandler.GetUserByUsername) // GET /api/v1/users/username/john
		users.PUT("/:id", userHandler.UpdateUser)                       // PUT /api/v1/users/1
		users.DELETE("/:id", userHandler.DeleteUser)                    // DELETE /api/v1/users/1
	}
}