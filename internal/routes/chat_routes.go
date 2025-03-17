package routes

import (
	"callcenter/internal/handlers"
	"callcenter/internal/middleware"
	"callcenter/internal/services"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func SetupChatRoutes(r *gin.Engine, db *gorm.DB) {
	nlpService := services.NewNLPService()
	chatService := services.NewChatService(db)
	chatHandler := handlers.NewChatHandler(db, nlpService, chatService)

	chat := r.Group("/api/v1/chat")
	chat.Use(middleware.AuthMiddleware())
	{
		chat.POST("/sessions", chatHandler.CreateSession)
		chat.GET("/sessions", chatHandler.ListSessions)
		chat.GET("/sessions/:id", chatHandler.GetSession)
		chat.POST("/sessions/:id/messages", chatHandler.SendMessage)
		chat.GET("/sessions/:id/messages", chatHandler.GetMessages)
		chat.POST("/message", chatHandler.HandleMessage)
	}
}
