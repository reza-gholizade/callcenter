package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"gorm.io/gorm"

	"callcenter/internal/models"
	"callcenter/internal/services"
)

// TicketHandler handles ticket-related HTTP requests
type TicketHandler struct {
	db            *gorm.DB
	ticketService *services.TicketService
}

// NewTicketHandler creates a new instance of TicketHandler
func NewTicketHandler(db *gorm.DB, ticketService *services.TicketService) *TicketHandler {
	return &TicketHandler{
		db:            db,
		ticketService: ticketService,
	}
}

type CreateTicketRequest struct {
	Subject     string `json:"subject" binding:"required"`
	Description string `json:"description" binding:"required"`
	Priority    string `json:"priority" binding:"required,oneof=low medium high"`
}

func (h *TicketHandler) CreateTicket(c *gin.Context) {
	var req CreateTicketRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	userID, _ := c.Get("userID")

	ticket := models.Ticket{
		ID:          uuid.New(),
		UserID:      userID.(uuid.UUID),
		Number:      generateTicketNumber(), // TODO: Implement ticket number generation
		Status:      "open",
		Subject:     req.Subject,
		Description: req.Description,
		Priority:    req.Priority,
	}

	if err := h.db.Create(&ticket).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create ticket"})
		return
	}

	// Create initial ticket history
	history := models.TicketHistory{
		ID:          uuid.New(),
		TicketID:    ticket.ID,
		Action:      "created",
		Description: "Ticket created",
		UserID:      userID.(uuid.UUID),
	}

	if err := h.db.Create(&history).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create ticket history"})
		return
	}

	c.JSON(http.StatusCreated, ticket)
}

func (h *TicketHandler) ListTickets(c *gin.Context) {
	userID, _ := c.Get("userID")

	var tickets []models.Ticket
	if err := h.db.Where("user_id = ?", userID).Find(&tickets).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch tickets"})
		return
	}

	c.JSON(http.StatusOK, tickets)
}

func (h *TicketHandler) GetTicket(c *gin.Context) {
	userID, _ := c.Get("userID")
	ticketID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ticket ID"})
		return
	}

	var ticket models.Ticket
	if err := h.db.Where("id = ? AND user_id = ?", ticketID, userID).First(&ticket).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Ticket not found"})
		return
	}

	c.JSON(http.StatusOK, ticket)
}

type UpdateStatusRequest struct {
	Status      string `json:"status" binding:"required,oneof=open in_progress resolved closed"`
	Description string `json:"description"`
}

func (h *TicketHandler) UpdateTicketStatus(c *gin.Context) {
	userID, _ := c.Get("userID")
	ticketID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ticket ID"})
		return
	}

	var req UpdateStatusRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Verify ticket exists and belongs to user
	var ticket models.Ticket
	if err := h.db.Where("id = ? AND user_id = ?", ticketID, userID).First(&ticket).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Ticket not found"})
		return
	}

	// Update ticket status
	ticket.Status = req.Status
	if err := h.db.Save(&ticket).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update ticket status"})
		return
	}

	// Create status history entry
	history := models.TicketHistory{
		ID:          uuid.New(),
		TicketID:    ticketID,
		Action:      "status_updated",
		Description: req.Description,
		UserID:      userID.(uuid.UUID),
	}

	if err := h.db.Create(&history).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create ticket history"})
		return
	}

	c.JSON(http.StatusOK, ticket)
}

func (h *TicketHandler) GetTicketHistory(c *gin.Context) {
	userID, _ := c.Get("userID")
	ticketID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ticket ID"})
		return
	}

	// Verify ticket exists and belongs to user
	var ticket models.Ticket
	if err := h.db.Where("id = ? AND user_id = ?", ticketID, userID).First(&ticket).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Ticket not found"})
		return
	}

	var history []models.TicketHistory
	if err := h.db.Where("ticket_id = ?", ticketID).Order("created_at asc").Find(&history).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch ticket history"})
		return
	}

	c.JSON(http.StatusOK, history)
}

// Helper function to generate ticket numbers
func generateTicketNumber() string {
	// TODO: Implement proper ticket number generation
	return "TKT-" + uuid.New().String()[:8]
}

// CancelTicketRequest represents the request body for cancelling a ticket
type CancelTicketRequest struct {
	Reason string `json:"reason" binding:"required"`
}

// CancelTicket handles ticket cancellation requests
func (h *TicketHandler) CancelTicket(c *gin.Context) {
	ticketNumber := c.Param("ticketNumber")
	if ticketNumber == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Ticket number is required"})
		return
	}

	var req CancelTicketRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Get ticket from database
	var ticket models.Ticket
	if err := h.db.Where("ticket_number = ?", ticketNumber).First(&ticket).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Ticket not found"})
		return
	}

	// Update ticket status to cancelled
	if err := h.db.Model(&ticket).Updates(map[string]interface{}{
		"status":              "cancelled",
		"cancellation_reason": req.Reason,
	}).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to cancel ticket"})
		return
	}

	// Create ticket history entry
	history := models.TicketHistory{
		TicketID:    ticket.ID,
		Action:      "cancelled",
		Description: req.Reason,
	}
	if err := h.db.Create(&history).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to record ticket history"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Ticket cancellation request submitted successfully",
	})
}

// GetRefundStatus retrieves the refund status for a ticket
func (h *TicketHandler) GetRefundStatus(c *gin.Context) {
	ticketNumber := c.Param("ticketNumber")
	if ticketNumber == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Ticket number is required"})
		return
	}

	// Get ticket from database
	var ticket models.Ticket
	if err := h.db.Where("ticket_number = ?", ticketNumber).First(&ticket).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Ticket not found"})
		return
	}

	// Get refund request from database
	var refundRequest models.RefundRequest
	if err := h.db.Where("ticket_id = ?", ticket.ID).First(&refundRequest).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "No refund request found for this ticket"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"refund_request": refundRequest,
	})
}

// UpdateRefundStatusRequest represents the request body for updating refund status
type UpdateRefundStatusRequest struct {
	Status      string `json:"status" binding:"required"`
	ProcessedBy string `json:"processed_by" binding:"required"`
}

// UpdateRefundStatus updates the refund status for a ticket
func (h *TicketHandler) UpdateRefundStatus(c *gin.Context) {
	ticketNumber := c.Param("ticketNumber")
	if ticketNumber == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Ticket number is required"})
		return
	}

	var req UpdateRefundStatusRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err := h.ticketService.UpdateRefundStatus(c.Request.Context(), ticketNumber, req.Status, req.ProcessedBy)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Refund status updated successfully",
	})
}
