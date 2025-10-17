package repository

import (
	"restaurant-system/internal/dao"
	"restaurant-system/internal/domain"
)

// OrderRepository defines the interface for order operations
type OrderRepository interface {
	Create(order *domain.Order) error
	GetByID(id string) (*domain.Order, error)
	GetByReservationID(reservationID string) ([]*domain.Order, error)
	UpdateStatus(id string, status domain.OrderStatus) error
}

// orderRepository implements OrderRepository
type orderRepository struct {
	dao *dao.OrderDAO
}

// NewOrderRepository creates a new order repository
func NewOrderRepository(dao *dao.OrderDAO) OrderRepository {
	return &orderRepository{dao: dao}
}

func (r *orderRepository) Create(order *domain.Order) error {
	return r.dao.Create(order)
}

func (r *orderRepository) GetByID(id string) (*domain.Order, error) {
	return r.dao.GetByID(id)
}

func (r *orderRepository) GetByReservationID(reservationID string) ([]*domain.Order, error) {
	return r.dao.GetByReservationID(reservationID)
}

func (r *orderRepository) UpdateStatus(id string, status domain.OrderStatus) error {
	return r.dao.UpdateStatus(id, status)
}
