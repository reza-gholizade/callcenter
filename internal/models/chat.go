package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type ChatSession struct {
	ID           uuid.UUID `gorm:"type:uuid;primary_key"`
	UserID       uuid.UUID `gorm:"type:uuid;not null"`
	Platform     string    `gorm:"not null"`
	Status       string    `gorm:"not null;default:'active'"`
	LastActivity time.Time
	Messages     []ChatMessage `gorm:"foreignKey:SessionID"`
	User         User          `gorm:"foreignKey:UserID"`
	CreatedAt    time.Time
	UpdatedAt    time.Time
	DeletedAt    gorm.DeletedAt `gorm:"index"`
}

type ChatMessage struct {
	ID        uuid.UUID `gorm:"type:uuid;primary_key"`
	SessionID uuid.UUID `gorm:"type:uuid;not null"`
	Content   string    `gorm:"not null"`
	Role      string    `gorm:"not null"` // 'user' or 'assistant'
	Intent    string
	Entities  string
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`
}

type ChatLog struct {
	ID             uuid.UUID `gorm:"type:uuid;primary_key"`
	SessionID      uuid.UUID `gorm:"type:uuid;not null"`
	EventType      string    `gorm:"not null"`
	EventData      string
	Timestamp      time.Time
	ProcessingTime int64
	Success        bool
	Error          string
	CreatedAt      time.Time
	UpdatedAt      time.Time
	DeletedAt      gorm.DeletedAt `gorm:"index"`
}
