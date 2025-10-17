package services

import (
	"fmt"
	"restaurant-system/internal/domain"
	"time"

	"github.com/google/uuid"
)

// TicketService handles ticket/bill business logic
type TicketService struct {
	// Add repository when implemented
}

// NewTicketService creates a new ticket service
func NewTicketService() *TicketService {
	return &TicketService{}
}

// GenerateTicket generates a ticket for an order
func (s *TicketService) GenerateTicket(req GenerateTicketRequest) (*domain.Ticket, error) {
	// Validate payment method
	validPaymentMethods := []string{"cash", "card", "digital_wallet", "other"}
	isValidPaymentMethod := false
	for _, method := range validPaymentMethods {
		if req.PaymentMethod == method {
			isValidPaymentMethod = true
			break
		}
	}
	if !isValidPaymentMethod {
		return nil, fmt.Errorf("invalid payment method")
	}

	ticket := &domain.Ticket{
		ID:            uuid.New().String(),
		OrderID:       req.OrderID,
		ReservationID: req.ReservationID,
		TableNumber:   req.TableNumber,
		CustomerName:  req.CustomerName,
		Items:         req.Items,
		Subtotal:      req.Subtotal,
		Tax:           req.Tax,
		Tip:           req.Tip,
		Total:         req.Total,
		PaymentMethod: req.PaymentMethod,
		PaidAt:        time.Now(),
		CreatedAt:     time.Now(),
	}

	// TODO: Save to database
	// err := s.ticketRepo.Create(ticket)
	// if err != nil {
	//     return nil, fmt.Errorf("failed to create ticket: %w", err)
	// }

	return ticket, nil
}

// GenerateTicketRequest represents the request to generate a ticket
type GenerateTicketRequest struct {
	OrderID       string             `json:"order_id" binding:"required"`
	ReservationID string             `json:"reservation_id" binding:"required"`
	TableNumber   int                `json:"table_number" binding:"required"`
	CustomerName  string             `json:"customer_name" binding:"required"`
	Items         []domain.OrderItem `json:"items" binding:"required"`
	Subtotal      float64            `json:"subtotal" binding:"required"`
	Tax           float64            `json:"tax" binding:"required"`
	Tip           float64            `json:"tip"`
	Total         float64            `json:"total" binding:"required"`
	PaymentMethod string             `json:"payment_method" binding:"required"`
}
