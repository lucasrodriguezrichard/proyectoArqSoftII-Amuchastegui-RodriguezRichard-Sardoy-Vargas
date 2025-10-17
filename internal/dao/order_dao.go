package dao

import (
	"database/sql"
	"fmt"
	"restaurant-system/internal/domain"
	"time"
)

// OrderDAO handles database operations for orders
type OrderDAO struct {
	db *DB
}

// NewOrderDAO creates a new OrderDAO
func NewOrderDAO(db *DB) *OrderDAO {
	return &OrderDAO{db: db}
}

// Create creates a new order
func (dao *OrderDAO) Create(order *domain.Order) error {
	tx, err := dao.db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	// Insert order
	orderQuery := `
		INSERT INTO orders (id, reservation_id, table_id, customer_name, status, 
			subtotal, tax, tip, total, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)
	`

	_, err = tx.Exec(orderQuery, order.ID, order.ReservationID, order.TableID,
		order.CustomerName, order.Status, order.Subtotal, order.Tax, order.Tip,
		order.Total, order.CreatedAt, order.UpdatedAt)

	if err != nil {
		return err
	}

	// Insert order items
	for _, item := range order.Items {
		itemQuery := `
			INSERT INTO order_items (id, order_id, menu_item_id, quantity, price, notes)
			VALUES ($1, $2, $3, $4, $5, $6)
		`
		_, err = tx.Exec(itemQuery, item.ID, item.OrderID, item.MenuItemID,
			item.Quantity, item.Price, item.Notes)
		if err != nil {
			return err
		}
	}

	return tx.Commit()
}

// GetByID retrieves an order by ID with items
func (dao *OrderDAO) GetByID(id string) (*domain.Order, error) {
	// Get order
	orderQuery := `
		SELECT id, reservation_id, table_id, customer_name, status, 
			subtotal, tax, tip, total, created_at, updated_at
		FROM orders WHERE id = $1
	`

	order := &domain.Order{}
	err := dao.db.QueryRow(orderQuery, id).Scan(
		&order.ID, &order.ReservationID, &order.TableID, &order.CustomerName,
		&order.Status, &order.Subtotal, &order.Tax, &order.Tip, &order.Total,
		&order.CreatedAt, &order.UpdatedAt)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("order not found")
		}
		return nil, err
	}

	// Get order items
	itemsQuery := `
		SELECT id, order_id, menu_item_id, quantity, price, notes
		FROM order_items WHERE order_id = $1
	`

	rows, err := dao.db.Query(itemsQuery, id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var items []domain.OrderItem
	for rows.Next() {
		item := domain.OrderItem{}
		err := rows.Scan(&item.ID, &item.OrderID, &item.MenuItemID,
			&item.Quantity, &item.Price, &item.Notes)
		if err != nil {
			return nil, err
		}
		items = append(items, item)
	}

	order.Items = items
	return order, nil
}

// GetByReservationID retrieves orders by reservation ID
func (dao *OrderDAO) GetByReservationID(reservationID string) ([]*domain.Order, error) {
	query := `
		SELECT id, reservation_id, table_id, customer_name, status, 
			subtotal, tax, tip, total, created_at, updated_at
		FROM orders WHERE reservation_id = $1
		ORDER BY created_at ASC
	`

	rows, err := dao.db.Query(query, reservationID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var orders []*domain.Order
	for rows.Next() {
		order := &domain.Order{}
		err := rows.Scan(&order.ID, &order.ReservationID, &order.TableID,
			&order.CustomerName, &order.Status, &order.Subtotal, &order.Tax,
			&order.Tip, &order.Total, &order.CreatedAt, &order.UpdatedAt)
		if err != nil {
			return nil, err
		}
		orders = append(orders, order)
	}

	return orders, nil
}

// UpdateStatus updates the status of an order
func (dao *OrderDAO) UpdateStatus(id string, status domain.OrderStatus) error {
	query := `UPDATE orders SET status = $1, updated_at = $2 WHERE id = $3`
	_, err := dao.db.Exec(query, string(status), time.Now(), id)
	return err
}
