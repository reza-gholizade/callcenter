package services

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"

	"callcenter/internal/models"
)

// ChatService handles the core chatbot functionality
type ChatService struct {
	db *gorm.DB
}

// NewChatService creates a new instance of ChatService
func NewChatService(db *gorm.DB) *ChatService {
	return &ChatService{
		db: db,
	}
}

// CreateSession creates a new chat session
func (s *ChatService) CreateSession(ctx context.Context, platform, userID string) (*models.ChatSession, error) {
	uid, err := uuid.Parse(userID)
	if err != nil {
		return nil, fmt.Errorf("invalid user ID: %w", err)
	}

	session := &models.ChatSession{
		ID:           uuid.New(),
		UserID:       uid,
		Platform:     platform,
		Status:       "active",
		LastActivity: time.Now(),
	}

	if err := s.db.WithContext(ctx).Create(session).Error; err != nil {
		return nil, fmt.Errorf("failed to create chat session: %w", err)
	}

	return session, nil
}

// AddMessage adds a new message to a chat session
func (s *ChatService) AddMessage(ctx context.Context, sessionID uuid.UUID, content string, role string, intent string, entities map[string]interface{}) (*models.ChatMessage, error) {
	entitiesJSON, err := json.Marshal(entities)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal entities: %w", err)
	}

	message := &models.ChatMessage{
		ID:        uuid.New(),
		SessionID: sessionID,
		Content:   content,
		Role:      role,
		Intent:    intent,
		Entities:  string(entitiesJSON),
	}

	if err := s.db.WithContext(ctx).Create(message).Error; err != nil {
		return nil, fmt.Errorf("failed to create message: %w", err)
	}

	// Update session last activity
	if err := s.db.WithContext(ctx).Model(&models.ChatSession{}).
		Where("id = ?", sessionID).
		Update("last_activity", time.Now()).Error; err != nil {
		return nil, fmt.Errorf("failed to update session activity: %w", err)
	}

	return message, nil
}

// GetChatHistory retrieves the message history for a session
func (s *ChatService) GetChatHistory(ctx context.Context, sessionID uuid.UUID) ([]models.ChatMessage, error) {
	var messages []models.ChatMessage
	if err := s.db.WithContext(ctx).
		Where("session_id = ?", sessionID).
		Order("created_at ASC").
		Find(&messages).Error; err != nil {
		return nil, fmt.Errorf("failed to get chat history: %w", err)
	}
	return messages, nil
}

// LogChatEvent logs a chat-related event for monitoring and analytics
func (s *ChatService) LogChatEvent(ctx context.Context, sessionID uuid.UUID, eventType string, eventData map[string]interface{}, processingTime int64, success bool, err error) error {
	var errorMsg string
	if err != nil {
		errorMsg = err.Error()
	}

	eventDataJSON, err := json.Marshal(eventData)
	if err != nil {
		return fmt.Errorf("failed to marshal event data: %w", err)
	}

	log := &models.ChatLog{
		ID:             uuid.New(),
		SessionID:      sessionID,
		EventType:      eventType,
		EventData:      string(eventDataJSON),
		Timestamp:      time.Now(),
		ProcessingTime: processingTime,
		Success:        success,
		Error:          errorMsg,
	}

	if err := s.db.WithContext(ctx).Create(log).Error; err != nil {
		return fmt.Errorf("failed to create chat log: %w", err)
	}

	return nil
}

// EscalateSession marks a chat session as escalated to human support
func (s *ChatService) EscalateSession(ctx context.Context, sessionID uuid.UUID) error {
	if err := s.db.WithContext(ctx).Model(&models.ChatSession{}).
		Where("id = ?", sessionID).
		Update("status", "escalated").Error; err != nil {
		return fmt.Errorf("failed to escalate session: %w", err)
	}
	return nil
}

// CloseSession marks a chat session as closed
func (s *ChatService) CloseSession(ctx context.Context, sessionID uuid.UUID) error {
	if err := s.db.WithContext(ctx).Model(&models.ChatSession{}).
		Where("id = ?", sessionID).
		Update("status", "closed").Error; err != nil {
		return fmt.Errorf("failed to close session: %w", err)
	}
	return nil
}
