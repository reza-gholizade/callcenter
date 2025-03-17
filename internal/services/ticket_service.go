package services

import (
	"context"
	"fmt"
	"time"

	"gorm.io/gorm"

	"callcenter/internal/models"
)

// TicketService handles ticket-related operations
type TicketService struct {
	db *gorm.DB
}

// NewTicketService creates a new instance of TicketService
func NewTicketService(db *gorm.DB) *TicketService {
	return &TicketService{
		db: db,
	}
}

// GetTicket retrieves ticket details by ticket number
func (s *TicketService) GetTicket(ctx context.Context, ticketNumber string) (*models.Ticket, error) {
	var ticket models.Ticket
	if err := s.db.WithContext(ctx).
		Where("ticket_number = ?", ticketNumber).
		First(&ticket).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("ticket not found: %s", ticketNumber)
		}
		return nil, fmt.Errorf("failed to get ticket: %w", err)
	}
	return &ticket, nil
}

// GetTicketByPhone retrieves ticket details by passenger phone number
func (s *TicketService) GetTicketByPhone(ctx context.Context, phoneNumber string) ([]models.Ticket, error) {
	var tickets []models.Ticket
	if err := s.db.WithContext(ctx).
		Where("passenger_phone = ?", phoneNumber).
		Find(&tickets).Error; err != nil {
		return nil, fmt.Errorf("failed to get tickets: %w", err)
	}
	return tickets, nil
}

// CancelTicket cancels a ticket and initiates refund process
func (s *TicketService) CancelTicket(ctx context.Context, ticketNumber string, reason string) error {
	// Start transaction
	tx := s.db.WithContext(ctx).Begin()
	if tx.Error != nil {
		return fmt.Errorf("failed to start transaction: %w", tx.Error)
	}
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	// Get ticket
	var ticket models.Ticket
	if err := tx.Where("ticket_number = ?", ticketNumber).First(&ticket).Error; err != nil {
		tx.Rollback()
		if err == gorm.ErrRecordNotFound {
			return fmt.Errorf("ticket not found: %s", ticketNumber)
		}
		return fmt.Errorf("failed to get ticket: %w", err)
	}

	// Check if ticket can be cancelled
	if ticket.Status != "active" {
		tx.Rollback()
		return fmt.Errorf("ticket cannot be cancelled: current status is %s", ticket.Status)
	}

	// Calculate refund amount based on ticket type and cancellation policy
	refundAmount, err := s.calculateRefundAmount(ctx, &ticket)
	if err != nil {
		tx.Rollback()
		return fmt.Errorf("failed to calculate refund amount: %w", err)
	}

	// Create refund request
	refundRequest := &models.RefundRequest{
		TicketNumber: ticketNumber,
		RequestedBy:  "system", // TODO: Replace with actual user ID
		Reason:       reason,
		Status:       "pending",
		Amount:       refundAmount,
		Currency:     "USD", // Hardcoded for now since Ticket model doesn't have Currency field
	}

	if err := tx.Create(refundRequest).Error; err != nil {
		tx.Rollback()
		return fmt.Errorf("failed to create refund request: %w", err)
	}

	// Update ticket status
	ticket.Status = "cancelled"

	// Update ticket in database
	if err := tx.Model(&ticket).Updates(map[string]interface{}{
		"status": "cancelled",
	}).Error; err != nil {
		tx.Rollback()
		return fmt.Errorf("failed to update ticket: %w", err)
	}

	// Commit transaction
	if err := tx.Commit().Error; err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	return nil
}

// GetRefundStatus retrieves the refund status for a ticket
func (s *TicketService) GetRefundStatus(ctx context.Context, ticketNumber string) (*models.RefundRequest, error) {
	var refundRequest models.RefundRequest
	if err := s.db.WithContext(ctx).
		Where("ticket_number = ?", ticketNumber).
		Order("created_at DESC").
		First(&refundRequest).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("no refund request found for ticket: %s", ticketNumber)
		}
		return nil, fmt.Errorf("failed to get refund status: %w", err)
	}
	return &refundRequest, nil
}

// calculateRefundAmount calculates the refund amount based on ticket type and cancellation policy
func (s *TicketService) calculateRefundAmount(ctx context.Context, ticket *models.Ticket) (float64, error) {
	// TODO: Implement actual refund calculation logic
	// This should consider:
	// 1. Ticket type (charter/systematic)
	// 2. Time until departure
	// 3. Airline cancellation policy
	// 4. Any applicable fees or penalties

	// Placeholder implementation
	switch ticket.TicketType {
	case "charter":
		return ticket.Price * 0.5, nil // 50% refund for charter tickets
	case "systematic":
		return ticket.Price * 0.8, nil // 80% refund for systematic tickets
	default:
		return 0, fmt.Errorf("invalid ticket type: %s", ticket.TicketType)
	}
}

// UpdateRefundStatus updates the refund status for a ticket
func (s *TicketService) UpdateRefundStatus(ctx context.Context, ticketNumber string, status string, processedBy string) error {
	// Start transaction
	tx := s.db.WithContext(ctx).Begin()
	if tx.Error != nil {
		return fmt.Errorf("failed to start transaction: %w", tx.Error)
	}
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	// Update refund request
	var refundRequest models.RefundRequest
	if err := tx.Where("ticket_number = ?", ticketNumber).
		Order("created_at DESC").
		First(&refundRequest).Error; err != nil {
		tx.Rollback()
		if err == gorm.ErrRecordNotFound {
			return fmt.Errorf("no refund request found for ticket: %s", ticketNumber)
		}
		return fmt.Errorf("failed to get refund request: %w", err)
	}

	refundRequest.Status = status
	refundRequest.ProcessedBy = processedBy
	refundRequest.ProcessedAt = time.Now()

	if err := tx.Save(&refundRequest).Error; err != nil {
		tx.Rollback()
		return fmt.Errorf("failed to update refund request: %w", err)
	}

	// Update ticket refund status
	var ticket models.Ticket
	if err := tx.Where("ticket_number = ?", ticketNumber).First(&ticket).Error; err != nil {
		tx.Rollback()
		return fmt.Errorf("failed to get ticket: %w", err)
	}

	ticket.RefundStatus = status
	if status == "processed" {
		ticket.RefundProcessed = time.Now()
	}

	if err := tx.Save(&ticket).Error; err != nil {
		tx.Rollback()
		return fmt.Errorf("failed to update ticket: %w", err)
	}

	// Commit transaction
	if err := tx.Commit().Error; err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	return nil
}
