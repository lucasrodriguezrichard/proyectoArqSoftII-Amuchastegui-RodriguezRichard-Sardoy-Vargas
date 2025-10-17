package dao

import (
	"database/sql"
	"fmt"
	"restaurant-system/internal/domain"
	"time"
)

// ReservationDAO handles database operations for reservations
type ReservationDAO struct {
	db *DB
}

// NewReservationDAO creates a new ReservationDAO
func NewReservationDAO(db *DB) *ReservationDAO {
	return &ReservationDAO{db: db}
}

// Create creates a new reservation
func (dao *ReservationDAO) Create(reservation *domain.Reservation) error {
	query := `
		INSERT INTO reservations (id, restaurant_id, table_id, customer_name, customer_phone, 
			customer_email, reservation_date, meal_type, party_size, status, special_requests, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13)
	`

	_, err := dao.db.Exec(query, reservation.ID, reservation.RestaurantID, reservation.TableID,
		reservation.CustomerName, reservation.CustomerPhone, reservation.CustomerEmail,
		reservation.ReservationDate, reservation.MealType, reservation.PartySize,
		reservation.Status, reservation.SpecialRequests, reservation.CreatedAt, reservation.UpdatedAt)

	return err
}

// GetByID retrieves a reservation by ID
func (dao *ReservationDAO) GetByID(id string) (*domain.Reservation, error) {
	query := `
		SELECT id, restaurant_id, table_id, customer_name, customer_phone, customer_email,
			reservation_date, meal_type, party_size, status, special_requests, created_at, updated_at
		FROM reservations WHERE id = $1
	`

	reservation := &domain.Reservation{}
	err := dao.db.QueryRow(query, id).Scan(
		&reservation.ID, &reservation.RestaurantID, &reservation.TableID,
		&reservation.CustomerName, &reservation.CustomerPhone, &reservation.CustomerEmail,
		&reservation.ReservationDate, &reservation.MealType, &reservation.PartySize,
		&reservation.Status, &reservation.SpecialRequests, &reservation.CreatedAt, &reservation.UpdatedAt)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("reservation not found")
		}
		return nil, err
	}

	return reservation, nil
}

// GetByDateRange retrieves reservations within a date range
func (dao *ReservationDAO) GetByDateRange(startDate, endDate time.Time) ([]*domain.Reservation, error) {
	query := `
		SELECT id, restaurant_id, table_id, customer_name, customer_phone, customer_email,
			reservation_date, meal_type, party_size, status, special_requests, created_at, updated_at
		FROM reservations 
		WHERE reservation_date BETWEEN $1 AND $2
		ORDER BY reservation_date ASC
	`

	rows, err := dao.db.Query(query, startDate, endDate)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var reservations []*domain.Reservation
	for rows.Next() {
		reservation := &domain.Reservation{}
		err := rows.Scan(
			&reservation.ID, &reservation.RestaurantID, &reservation.TableID,
			&reservation.CustomerName, &reservation.CustomerPhone, &reservation.CustomerEmail,
			&reservation.ReservationDate, &reservation.MealType, &reservation.PartySize,
			&reservation.Status, &reservation.SpecialRequests, &reservation.CreatedAt, &reservation.UpdatedAt)
		if err != nil {
			return nil, err
		}
		reservations = append(reservations, reservation)
	}

	return reservations, nil
}

// UpdateStatus updates the status of a reservation
func (dao *ReservationDAO) UpdateStatus(id string, status domain.ReservationStatus) error {
	query := `UPDATE reservations SET status = $1, updated_at = $2 WHERE id = $3`
	_, err := dao.db.Exec(query, status, time.Now(), id)
	return err
}

// Delete deletes a reservation
func (dao *ReservationDAO) Delete(id string) error {
	query := `DELETE FROM reservations WHERE id = $1`
	_, err := dao.db.Exec(query, id)
	return err
}
