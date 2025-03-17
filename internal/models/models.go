package models

import (
	"time"

	"gorm.io/gorm"
)

// RefundRequest represents a ticket refund request
type RefundRequest struct {
	gorm.Model
	TicketNumber string  `gorm:"not null"`
	RequestedBy  string  `gorm:"not null"` // user ID or system
	Reason       string  `gorm:"type:text"`
	Status       string  `gorm:"not null"` // pending, approved, rejected, processed
	Amount       float64 `gorm:"not null"`
	Currency     string  `gorm:"not null"`
	ProcessedAt  time.Time
	ProcessedBy  string // agent ID or system
	Notes        string `gorm:"type:text"`
}
