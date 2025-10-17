package dao

import (
	"database/sql"
	"fmt"
	"restaurant-system/internal/domain"
)

// TableDAO handles database operations for tables
type TableDAO struct {
	db *DB
}

// NewTableDAO creates a new TableDAO
func NewTableDAO(db *DB) *TableDAO {
	return &TableDAO{db: db}
}

// Create creates a new table
func (dao *TableDAO) Create(table *domain.Table) error {
	query := `
		INSERT INTO tables (id, restaurant_id, number, capacity, is_available)
		VALUES ($1, $2, $3, $4, $5)
	`

	_, err := dao.db.Exec(query, table.ID, table.RestaurantID, table.Number,
		table.Capacity, table.IsAvailable)

	return err
}

// GetByID retrieves a table by ID
func (dao *TableDAO) GetByID(id string) (*domain.Table, error) {
	query := `
		SELECT id, restaurant_id, number, capacity, is_available
		FROM tables WHERE id = $1
	`

	table := &domain.Table{}
	err := dao.db.QueryRow(query, id).Scan(
		&table.ID, &table.RestaurantID, &table.Number,
		&table.Capacity, &table.IsAvailable)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("table not found")
		}
		return nil, err
	}

	return table, nil
}

// GetByRestaurantID retrieves all tables for a restaurant
func (dao *TableDAO) GetByRestaurantID(restaurantID string) ([]*domain.Table, error) {
	query := `
		SELECT id, restaurant_id, number, capacity, is_available
		FROM tables WHERE restaurant_id = $1
		ORDER BY number ASC
	`

	rows, err := dao.db.Query(query, restaurantID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var tables []*domain.Table
	for rows.Next() {
		table := &domain.Table{}
		err := rows.Scan(&table.ID, &table.RestaurantID, &table.Number,
			&table.Capacity, &table.IsAvailable)
		if err != nil {
			return nil, err
		}
		tables = append(tables, table)
	}

	return tables, nil
}

// GetAvailableTables retrieves all available tables for a restaurant
func (dao *TableDAO) GetAvailableTables(restaurantID string) ([]*domain.Table, error) {
	query := `
		SELECT id, restaurant_id, number, capacity, is_available
		FROM tables 
		WHERE restaurant_id = $1 AND is_available = true
		ORDER BY number ASC
	`

	rows, err := dao.db.Query(query, restaurantID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var tables []*domain.Table
	for rows.Next() {
		table := &domain.Table{}
		err := rows.Scan(&table.ID, &table.RestaurantID, &table.Number,
			&table.Capacity, &table.IsAvailable)
		if err != nil {
			return nil, err
		}
		tables = append(tables, table)
	}

	return tables, nil
}

// Update updates a table
func (dao *TableDAO) Update(table *domain.Table) error {
	query := `
		UPDATE tables 
		SET number = $1, capacity = $2, is_available = $3
		WHERE id = $4
	`

	_, err := dao.db.Exec(query, table.Number, table.Capacity,
		table.IsAvailable, table.ID)

	return err
}

// UpdateAvailability updates the availability status of a table
func (dao *TableDAO) UpdateAvailability(id string, isAvailable bool) error {
	query := `UPDATE tables SET is_available = $1 WHERE id = $2`
	_, err := dao.db.Exec(query, isAvailable, id)
	return err
}

// Delete deletes a table
func (dao *TableDAO) Delete(id string) error {
	query := `DELETE FROM tables WHERE id = $1`
	_, err := dao.db.Exec(query, id)
	return err
}
