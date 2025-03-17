package routes

import (
	"callcenter/internal/handlers"
	"callcenter/internal/middleware"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func SetupAuthRoutes(r *gin.Engine, db *gorm.DB) {
	authHandler := handlers.NewAuthHandler(db)

	auth := r.Group("/api/v1/auth")
	{
		auth.POST("/login", authHandler.Login)
		auth.POST("/register", authHandler.Register)

		// Protected routes
		auth.Use(middleware.AuthMiddleware())
		auth.GET("/me", authHandler.GetCurrentUser)
	}
}
