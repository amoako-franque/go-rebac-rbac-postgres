package main

import (
	"log"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
)

func main() {
	if err := ConnectDB(); err != nil {
		log.Fatal("Failed to connect to database:", err)
	}

	port := os.Getenv("PORT")
	if port == "" {
		port = "4000"
	}

	// Set Gin mode based on environment
	if os.Getenv("GIN_MODE") == "" {
		gin.SetMode(gin.ReleaseMode)
	}

	r := gin.Default()

	// CORS middleware
	r.Use(corsMiddleware())

	// Health check
	r.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok", "service": "rbac-rebac-api"})
	})

	// Auth routes
	r.POST("/auth/register", RegisterHandler)
	r.POST("/auth/login", LoginHandler)

	// Protected routes
	records := r.Group("/records")
	records.Use(RequireAuth())
	{
		records.GET("/rbac/:id", RBAC("record:read"), GetRecordRBAC)
		records.GET("/rebac/:id", ReBACResource("assigned_to"), GetRecordReBAC)
	}

	// Seed endpoint (should be protected in production)
	r.POST("/seed", func(c *gin.Context) {
		out := Seed()
		c.JSON(http.StatusOK, gin.H{"message": "database seeded successfully", "data": out})
	})

	log.Printf("Server listening on port %s", port)
	if err := r.Run(":" + port); err != nil {
		log.Fatal("Failed to start server:", err)
	}
}

func corsMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT, DELETE, PATCH")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(http.StatusNoContent)
			return
		}

		c.Next()
	}
}
