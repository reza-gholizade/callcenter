package handlers

import (
	"context"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"gorm.io/gorm"

	"callcenter/internal/models"
	"callcenter/internal/services"
)

// ChatHandler handles chat-related HTTP requests
type ChatHandler struct {
	db          *gorm.DB
	nlpService  *services.NLPService
	chatService *services.ChatService
}

// NewChatHandler creates a new instance of ChatHandler
func NewChatHandler(db *gorm.DB, nlpService *services.NLPService, chatService *services.ChatService) *ChatHandler {
	return &ChatHandler{
		db:          db,
		nlpService:  nlpService,
		chatService: chatService,
	}
}

// MessageRequest represents the request body for sending a message
type MessageRequest struct {
	Content  string `json:"content" binding:"required"`
	Platform string `json:"platform" binding:"required"`
	UserID   string `json:"user_id" binding:"required"`
}

// HandleMessage handles incoming chat messages
func (h *ChatHandler) HandleMessage(c *gin.Context) {
	var req MessageRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Start timing the request
	startTime := time.Now()

	// Detect intent and extract entities
	intent, err := h.nlpService.DetectIntent(c.Request.Context(), req.Content)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to process message"})
		return
	}

	// Create or get chat session
	session, err := h.chatService.CreateSession(c.Request.Context(), req.Platform, req.UserID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create chat session"})
		return
	}

	// Add user message
	_, err = h.chatService.AddMessage(c.Request.Context(), session.ID, req.Content, "user", intent.Name, intent.Entities)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save message"})
		return
	}

	// Generate bot response based on intent
	response, err := h.generateResponse(c.Request.Context(), intent, session.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate response"})
		return
	}

	// Add bot response
	_, err = h.chatService.AddMessage(c.Request.Context(), session.ID, response, "assistant", "", nil)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save bot response"})
		return
	}

	// Log the chat event
	processingTime := time.Since(startTime).Milliseconds()
	err = h.chatService.LogChatEvent(c.Request.Context(), session.ID, "message", map[string]interface{}{
		"intent":   intent.Name,
		"entities": intent.Entities,
	}, processingTime, true, nil)
	if err != nil {
		// Log error but don't fail the request
	}

	c.JSON(http.StatusOK, gin.H{
		"session_id": session.ID,
		"response":   response,
	})
}

// GetChatHistory retrieves the chat history for a session
func (h *ChatHandler) GetChatHistory(c *gin.Context) {
	sessionID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid session ID"})
		return
	}

	messages, err := h.chatService.GetChatHistory(c.Request.Context(), sessionID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get chat history"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"messages": messages,
	})
}

// generateResponse generates a bot response based on the detected intent
func (h *ChatHandler) generateResponse(ctx context.Context, intent *services.Intent, sessionID uuid.UUID) (string, error) {
	// TODO: Implement actual response generation logic
	// This should consider:
	// 1. The detected intent
	// 2. Extracted entities
	// 3. Conversation context
	// 4. Business rules and policies

	// Placeholder implementation
	switch intent.Name {
	case "ticket_lookup":
		return "I can help you find your ticket. Could you please provide your ticket number or booking reference?", nil
	case "ticket_cancellation":
		return "I can help you cancel your ticket. Could you please provide your ticket number?", nil
	case "refund_inquiry":
		return "I can help you check your refund status. Could you please provide your ticket number?", nil
	case "baggage_policy":
		return "I can help you with baggage policy information. Which airline are you flying with?", nil
	default:
		return "I'm not sure I understand. Could you please rephrase your question?", nil
	}
}

type CreateSessionRequest struct {
	Subject string `json:"subject" binding:"required"`
}

func (h *ChatHandler) CreateSession(c *gin.Context) {
	var req CreateSessionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	userID, _ := c.Get("userID")

	session := models.ChatSession{
		ID:           uuid.New(),
		UserID:       userID.(uuid.UUID),
		Platform:     "web",
		Status:       "active",
		LastActivity: time.Now(),
	}

	if err := h.db.Create(&session).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create session"})
		return
	}

	c.JSON(http.StatusCreated, session)
}

func (h *ChatHandler) ListSessions(c *gin.Context) {
	userID, _ := c.Get("userID")

	var sessions []models.ChatSession
	if err := h.db.Where("user_id = ?", userID).Find(&sessions).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch sessions"})
		return
	}

	c.JSON(http.StatusOK, sessions)
}

func (h *ChatHandler) GetSession(c *gin.Context) {
	userID, _ := c.Get("userID")
	sessionID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid session ID"})
		return
	}

	var session models.ChatSession
	if err := h.db.Where("id = ? AND user_id = ?", sessionID, userID).First(&session).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Session not found"})
		return
	}

	c.JSON(http.StatusOK, session)
}

type SendMessageRequest struct {
	Content string `json:"content" binding:"required"`
}

func (h *ChatHandler) SendMessage(c *gin.Context) {
	userID, _ := c.Get("userID")
	sessionID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid session ID"})
		return
	}

	var req SendMessageRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Verify session exists and belongs to user
	var session models.ChatSession
	if err := h.db.Where("id = ? AND user_id = ?", sessionID, userID).First(&session).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Session not found"})
		return
	}

	// Create user message
	userMessage := models.ChatMessage{
		ID:        uuid.New(),
		SessionID: sessionID,
		Content:   req.Content,
		Role:      "user",
	}

	if err := h.db.Create(&userMessage).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save message"})
		return
	}

	// TODO: Process message with NLP service and generate response
	// For now, just echo the message
	assistantMessage := models.ChatMessage{
		ID:        uuid.New(),
		SessionID: sessionID,
		Content:   "I received your message: " + req.Content,
		Role:      "assistant",
	}

	if err := h.db.Create(&assistantMessage).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save response"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"userMessage":      userMessage,
		"assistantMessage": assistantMessage,
	})
}

func (h *ChatHandler) GetMessages(c *gin.Context) {
	userID, _ := c.Get("userID")
	sessionID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid session ID"})
		return
	}

	// Verify session exists and belongs to user
	var session models.ChatSession
	if err := h.db.Where("id = ? AND user_id = ?", sessionID, userID).First(&session).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Session not found"})
		return
	}

	var messages []models.ChatMessage
	if err := h.db.Where("session_id = ?", sessionID).Order("created_at asc").Find(&messages).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch messages"})
		return
	}

	c.JSON(http.StatusOK, messages)
}
