package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Ticket struct {
	gorm.Model
	ID              uuid.UUID `gorm:"type:uuid;primary_key"`
	UserID          uuid.UUID `gorm:"type:uuid;not null"`
	Number          string    `gorm:"uniqueIndex;not null"`
	Status          string    `gorm:"not null;default:'open'"`
	Subject         string    `gorm:"not null"`
	Description     string    `gorm:"not null"`
	Priority        string    `gorm:"not null;default:'medium'"`
	TicketType      string    `gorm:"not null"` // charter, systematic
	Price           float64   `gorm:"not null"`
	Currency        string    `gorm:"not null"`
	RefundStatus    string
	RefundAmount    float64
	RefundProcessed time.Time
	User            User `gorm:"foreignKey:UserID"`
	History         []TicketHistory
}

type TicketStatus struct {
	gorm.Model
	ID          uuid.UUID `gorm:"type:uuid;primary_key"`
	TicketID    uuid.UUID `gorm:"type:uuid;not null"`
	Status      string    `gorm:"not null"`
	Description string
	Ticket      Ticket `gorm:"foreignKey:TicketID"`
}

type TicketHistory struct {
	gorm.Model
	ID          uuid.UUID `gorm:"type:uuid;primary_key"`
	TicketID    uuid.UUID `gorm:"type:uuid;not null"`
	Action      string    `gorm:"not null"`
	Description string
	UserID      uuid.UUID `gorm:"type:uuid;not null"`
	Ticket      Ticket    `gorm:"foreignKey:TicketID"`
	User        User      `gorm:"foreignKey:UserID"`
}
