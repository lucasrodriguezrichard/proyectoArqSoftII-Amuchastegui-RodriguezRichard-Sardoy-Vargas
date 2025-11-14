package service

import (
	"context"
	"fmt"
	"log"

	"github.com/blassardoy/restaurant-reservas/search-api/internal/cache"
	"github.com/blassardoy/restaurant-reservas/search-api/internal/domain"
	"github.com/blassardoy/restaurant-reservas/search-api/internal/repository"
)

type SyncService struct {
	repo      repository.SearchRepository
	resClient *ReservationClient
	cache     *cache.DualCache
}

func NewSyncService(repo repository.SearchRepository, resClient *ReservationClient, cacheLayer *cache.DualCache) *SyncService {
	return &SyncService{repo: repo, resClient: resClient, cache: cacheLayer}
}

// HandleEvent processes reservation events and updates table availability in Solr
func (s *SyncService) HandleEvent(ctx context.Context, op string, reservationID string) error {
	log.Printf("HandleEvent: op=%s, reservationID=%s", op, reservationID)

	// Get reservation details from Reservations API
	reservation, err := s.resClient.GetReservationByID(reservationID)
	if err != nil {
		log.Printf("ERROR: Failed to get reservation %s: %v", reservationID, err)
		return fmt.Errorf("failed to get reservation: %w", err)
	}

	// Extract table info from reservation
	tableNumber := reservation.TableNumber
	mealType := reservation.MealType
	date := reservation.DateTime.Format("2006-01-02")

	// Get table capacity from predefined tables
	capacity, found := getTableCapacity(tableNumber, mealType)
	if !found {
		log.Printf("WARNING: Unknown table config for table %d, meal_type %s. Using default capacity 4", tableNumber, mealType)
		capacity = 4
	}

	// Generate TableAvailability ID
	tableAvailID := domain.GenerateTableAvailabilityID(mealType, tableNumber, date)

	var updateErr error
	switch op {
	case "create", "confirm":
		// Mark table as NOT available (reserved)
		tableAvail := domain.NewTableAvailability(tableNumber, capacity, mealType, date)
		tableAvail.IsAvailable = false
		tableAvail.ReservationID = reservationID

		log.Printf("Marking table as UNAVAILABLE: %s", tableAvailID)
		updateErr = s.repo.Index(ctx, *tableAvail)

	case "delete", "cancel":
		// Mark table as available again (reservation cancelled/deleted)
		tableAvail := domain.NewTableAvailability(tableNumber, capacity, mealType, date)
		tableAvail.IsAvailable = true
		tableAvail.ReservationID = ""

		log.Printf("Marking table as AVAILABLE: %s", tableAvailID)
		updateErr = s.repo.Update(ctx, *tableAvail)

	case "update":
		// For updates, check the reservation status
		// If status is cancelled, mark as available; otherwise keep as unavailable
		tableAvail := domain.NewTableAvailability(tableNumber, capacity, mealType, date)
		if reservation.Status == "cancelled" {
			tableAvail.IsAvailable = true
			tableAvail.ReservationID = ""
			log.Printf("Update: Marking table as AVAILABLE (cancelled): %s", tableAvailID)
		} else {
			tableAvail.IsAvailable = false
			tableAvail.ReservationID = reservationID
			log.Printf("Update: Keeping table as UNAVAILABLE: %s", tableAvailID)
		}
		updateErr = s.repo.Update(ctx, *tableAvail)

	default:
		log.Printf("Unknown operation: %s", op)
		return nil
	}

	if updateErr != nil {
		log.Printf("ERROR: Failed to update Solr for table %s: %v", tableAvailID, updateErr)
		return updateErr
	}

	// Clear cache on success
	if s.cache != nil {
		s.cache.Clear()
		log.Printf("Cache cleared after processing event")
	}

	log.Printf("Successfully processed event: op=%s, table=%s", op, tableAvailID)
	return nil
}

// getTableCapacity returns the capacity for a given table number and meal type
func getTableCapacity(tableNumber int, mealType string) (int, bool) {
	// Predefined table capacities (same as in reservations-api)
	capacities := map[string][]int{
		"breakfast": {2, 2, 4, 4, 4, 6, 6, 6, 8, 8},
		"lunch":     {2, 2, 4, 4, 4, 6, 6, 6, 8, 8},
		"dinner":    {2, 2, 4, 4, 4, 6, 6, 6, 8, 8},
		"event":     {8, 10, 10, 12, 12, 15, 15, 18, 20, 20},
	}

	if caps, ok := capacities[mealType]; ok {
		if tableNumber >= 1 && tableNumber <= len(caps) {
			return caps[tableNumber-1], true
		}
	}
	return 0, false
}
