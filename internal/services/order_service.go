package services

import (
	"fmt"
	"restaurant-system/internal/domain"
	"restaurant-system/internal/repository"
	"time"

	"github.com/google/uuid"
)

// OrderService handles order business logic
type OrderService struct {
	orderRepo repository.OrderRepository
}

// NewOrderService creates a new order service
func NewOrderService(orderRepo repository.OrderRepository) *OrderService {
	return &OrderService{
		orderRepo: orderRepo,
	}
}

// CreateOrder creates a new order
func (s *OrderService) CreateOrder(req CreateOrderRequest) (*domain.Order, error) {
	// Validate order items
	if len(req.Items) == 0 {
		return nil, fmt.Errorf("order must have at least one item")
	}

	// Calculate totals
	var subtotal float64
	for _, item := range req.Items {
		subtotal += item.Price * float64(item.Quantity)
	}

	tax := subtotal * 0.1 // 10% tax
	tip := req.Tip
	total := subtotal + tax + tip

	order := &domain.Order{
		ID:            uuid.New().String(),
		ReservationID: req.ReservationID,
		TableID:       req.TableID,
		CustomerName:  req.CustomerName,
		Status:        domain.OrderStatusPending,
		Items:         req.Items,
		Subtotal:      subtotal,
		Tax:           tax,
		Tip:           tip,
		Total:         total,
		CreatedAt:     time.Now(),
		UpdatedAt:     time.Now(),
	}

	err := s.orderRepo.Create(order)
	if err != nil {
		return nil, fmt.Errorf("failed to create order: %w", err)
	}

	return order, nil
}

// GetOrder retrieves an order by ID
func (s *OrderService) GetOrder(id string) (*domain.Order, error) {
	return s.orderRepo.GetByID(id)
}

// GetOrdersByReservation retrieves orders by reservation ID
func (s *OrderService) GetOrdersByReservation(reservationID string) ([]*domain.Order, error) {
	return s.orderRepo.GetByReservationID(reservationID)
}

// UpdateOrderStatus updates the status of an order
func (s *OrderService) UpdateOrderStatus(id string, status string) error {
	return s.orderRepo.UpdateStatus(id, domain.OrderStatus(status))
}

// CreateOrderRequest represents the request to create an order
type CreateOrderRequest struct {
	ReservationID string             `json:"reservation_id" binding:"required"`
	TableID       string             `json:"table_id" binding:"required"`
	CustomerName  string             `json:"customer_name" binding:"required"`
	Items         []domain.OrderItem `json:"items" binding:"required"`
	Tip           float64            `json:"tip"`
}
