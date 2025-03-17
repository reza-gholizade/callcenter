package main

import (
	"callcenter/internal/database"
	"callcenter/internal/routes"
	"log"
	"os"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	// Load environment variables
	if err := godotenv.Load(); err != nil {
		log.Printf("Warning: .env file not found")
	}

	// Initialize database
	db, err := database.InitDB()
	if err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}

	// Create Gin router
	r := gin.Default()

	// Configure CORS
	config := cors.DefaultConfig()
	config.AllowOrigins = []string{"http://localhost:3000"} // Add your frontend URL
	config.AllowMethods = []string{"GET", "POST", "PUT", "PATCH", "DELETE", "HEAD", "OPTIONS"}
	config.AllowHeaders = []string{"Origin", "Content-Length", "Content-Type", "Authorization"}
	config.AllowCredentials = true
	r.Use(cors.New(config))

	// Set JWT secret in context
	r.Use(func(c *gin.Context) {
		c.Set("JWT_SECRET", os.Getenv("JWT_SECRET"))
		c.Next()
	})

	// Setup routes
	routes.SetupAuthRoutes(r, db)
	routes.SetupChatRoutes(r, db)
	routes.SetupTicketRoutes(r, db)

	// Start server
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	if err := r.Run(":" + port); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}

func initializeServices(router *gin.Engine) {
	// TODO: Initialize database connection
	// TODO: Initialize chat service
	// TODO: Initialize NLP service
	// TODO: Initialize external API clients

	// Setup routes
	setupRoutes(router)
}

func setupRoutes(router *gin.Engine) {
	// Health check
	router.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status": "ok",
		})
	})

	// API routes group
	api := router.Group("/api/v1")
	{
		// Chat endpoints
		chat := api.Group("/chat")
		{
			chat.POST("/message", handleChatMessage)
			chat.GET("/history/:sessionId", getChatHistory)
		}

		// Ticket management endpoints
		tickets := api.Group("/tickets")
		{
			tickets.GET("/:ticketNumber", getTicketDetails)
			tickets.POST("/:ticketNumber/cancel", cancelTicket)
			tickets.GET("/:ticketNumber/refund-status", getRefundStatus)
		}
	}
}

func handleChatMessage(c *gin.Context) {
	// TODO: Implement chat message handling
	c.JSON(200, gin.H{
		"message": "Chat message endpoint",
	})
}

func getChatHistory(c *gin.Context) {
	// TODO: Implement chat history retrieval
	c.JSON(200, gin.H{
		"message": "Chat history endpoint",
	})
}

func getTicketDetails(c *gin.Context) {
	// TODO: Implement ticket details retrieval
	c.JSON(200, gin.H{
		"message": "Ticket details endpoint",
	})
}

func cancelTicket(c *gin.Context) {
	// TODO: Implement ticket cancellation
	c.JSON(200, gin.H{
		"message": "Ticket cancellation endpoint",
	})
}

func getRefundStatus(c *gin.Context) {
	// TODO: Implement refund status check
	c.JSON(200, gin.H{
		"message": "Refund status endpoint",
	})
}
