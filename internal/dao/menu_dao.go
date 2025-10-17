package dao

import (
	"database/sql"
	"fmt"
	"restaurant-system/internal/domain"
)

// MenuDAO handles database operations for menu items
type MenuDAO struct {
	db *DB
}

// NewMenuDAO creates a new MenuDAO
func NewMenuDAO(db *DB) *MenuDAO {
	return &MenuDAO{db: db}
}

// Create creates a new menu item
func (dao *MenuDAO) Create(item *domain.MenuItem) error {
	query := `
		INSERT INTO menu_items (id, restaurant_id, name, description, price, category, is_available)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
	`

	_, err := dao.db.Exec(query, item.ID, item.RestaurantID, item.Name,
		item.Description, item.Price, item.Category, item.IsAvailable)

	return err
}

// GetByID retrieves a menu item by ID
func (dao *MenuDAO) GetByID(id string) (*domain.MenuItem, error) {
	query := `
		SELECT id, restaurant_id, name, description, price, category, is_available
		FROM menu_items WHERE id = $1
	`

	item := &domain.MenuItem{}
	err := dao.db.QueryRow(query, id).Scan(
		&item.ID, &item.RestaurantID, &item.Name, &item.Description,
		&item.Price, &item.Category, &item.IsAvailable)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("menu item not found")
		}
		return nil, err
	}

	return item, nil
}

// GetByRestaurantID retrieves all menu items for a restaurant
func (dao *MenuDAO) GetByRestaurantID(restaurantID string) ([]*domain.MenuItem, error) {
	query := `
		SELECT id, restaurant_id, name, description, price, category, is_available
		FROM menu_items WHERE restaurant_id = $1
		ORDER BY category, name ASC
	`

	rows, err := dao.db.Query(query, restaurantID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var items []*domain.MenuItem
	for rows.Next() {
		item := &domain.MenuItem{}
		err := rows.Scan(&item.ID, &item.RestaurantID, &item.Name,
			&item.Description, &item.Price, &item.Category, &item.IsAvailable)
		if err != nil {
			return nil, err
		}
		items = append(items, item)
	}

	return items, nil
}

// GetByCategory retrieves menu items by category
func (dao *MenuDAO) GetByCategory(restaurantID, category string) ([]*domain.MenuItem, error) {
	query := `
		SELECT id, restaurant_id, name, description, price, category, is_available
		FROM menu_items 
		WHERE restaurant_id = $1 AND category = $2
		ORDER BY name ASC
	`

	rows, err := dao.db.Query(query, restaurantID, category)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var items []*domain.MenuItem
	for rows.Next() {
		item := &domain.MenuItem{}
		err := rows.Scan(&item.ID, &item.RestaurantID, &item.Name,
			&item.Description, &item.Price, &item.Category, &item.IsAvailable)
		if err != nil {
			return nil, err
		}
		items = append(items, item)
	}

	return items, nil
}

// GetAvailableItems retrieves all available menu items
func (dao *MenuDAO) GetAvailableItems(restaurantID string) ([]*domain.MenuItem, error) {
	query := `
		SELECT id, restaurant_id, name, description, price, category, is_available
		FROM menu_items 
		WHERE restaurant_id = $1 AND is_available = true
		ORDER BY category, name ASC
	`

	rows, err := dao.db.Query(query, restaurantID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var items []*domain.MenuItem
	for rows.Next() {
		item := &domain.MenuItem{}
		err := rows.Scan(&item.ID, &item.RestaurantID, &item.Name,
			&item.Description, &item.Price, &item.Category, &item.IsAvailable)
		if err != nil {
			return nil, err
		}
		items = append(items, item)
	}

	return items, nil
}

// Update updates a menu item
func (dao *MenuDAO) Update(item *domain.MenuItem) error {
	query := `
		UPDATE menu_items 
		SET name = $1, description = $2, price = $3, category = $4, is_available = $5
		WHERE id = $6
	`

	_, err := dao.db.Exec(query, item.Name, item.Description, item.Price,
		item.Category, item.IsAvailable, item.ID)

	return err
}

// UpdateAvailability updates the availability status of a menu item
func (dao *MenuDAO) UpdateAvailability(id string, isAvailable bool) error {
	query := `UPDATE menu_items SET is_available = $1 WHERE id = $2`
	_, err := dao.db.Exec(query, isAvailable, id)
	return err
}

// Delete deletes a menu item
func (dao *MenuDAO) Delete(id string) error {
	query := `DELETE FROM menu_items WHERE id = $1`
	_, err := dao.db.Exec(query, id)
	return err
}