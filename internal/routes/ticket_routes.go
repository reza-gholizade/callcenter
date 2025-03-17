package routes

import (
	"callcenter/internal/handlers"
	"callcenter/internal/middleware"
	"callcenter/internal/services"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func SetupTicketRoutes(r *gin.Engine, db *gorm.DB) {
	ticketService := services.NewTicketService(db)
	ticketHandler := handlers.NewTicketHandler(db, ticketService)

	tickets := r.Group("/api/v1/tickets")
	tickets.Use(middleware.AuthMiddleware())
	{
		tickets.POST("", ticketHandler.CreateTicket)
		tickets.GET("", ticketHandler.ListTickets)
		tickets.GET("/:id", ticketHandler.GetTicket)
		tickets.PUT("/:id/status", ticketHandler.UpdateTicketStatus)
		tickets.GET("/:id/history", ticketHandler.GetTicketHistory)
	}
}
