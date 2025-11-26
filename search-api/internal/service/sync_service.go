package service

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/blassardoy/restaurant-reservas/search-api/internal/cache"
	"github.com/blassardoy/restaurant-reservas/search-api/internal/domain"
	"github.com/blassardoy/restaurant-reservas/search-api/internal/repository"
)

type SyncService struct {
	repo             repository.SearchRepository
	resClient        *ReservationClient
	tableClient      *TableClient
	cache            *cache.DualCache
	availabilityDays int
}

func NewSyncService(
	repo repository.SearchRepository,
	resClient *ReservationClient,
	tableClient *TableClient,
	cacheLayer *cache.DualCache,
	availabilityDays int,
) *SyncService {
	return &SyncService{
		repo:             repo,
		resClient:        resClient,
		tableClient:      tableClient,
		cache:            cacheLayer,
		availabilityDays: availabilityDays,
	}
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

	// Get table capacity from admin-managed tables
	capacity := 4
	if s.tableClient != nil {
		if table, err := s.tableClient.FindTable(mealType, tableNumber); err == nil {
			capacity = table.Capacity
		} else {
			log.Printf("WARNING: Failed to resolve table %d (%s): %v. Using default capacity 4", tableNumber, mealType, err)
		}
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

// HandleTableEvent processes table events and indexes available tables in Solr
func (s *SyncService) HandleTableEvent(ctx context.Context, op string, tableID string, metadata map[string]interface{}) error {
	log.Printf("HandleTableEvent: op=%s, tableID=%s", op, tableID)

	tableDoc, err := s.resolveTableDoc(tableID, metadata)
	if err != nil {
		if op == "delete" {
			log.Printf("WARNING: %v", err)
			return nil
		}
		return err
	}

	switch op {
	case "create":
		if err := s.indexTableWindow(ctx, tableDoc); err != nil {
			return err
		}
		log.Printf("Indexed table %d (%s) for %d days", tableDoc.TableNumber, tableDoc.MealType, s.windowDays())

	case "update":
		if err := s.indexTableWindow(ctx, tableDoc); err != nil {
			return err
		}
		log.Printf("Updated table %d (%s) for %d days", tableDoc.TableNumber, tableDoc.MealType, s.windowDays())

	case "delete":
		if err := s.deleteTableEntries(ctx, tableDoc); err != nil {
			return err
		}
		log.Printf("Deleted table %d (%s) entries", tableDoc.TableNumber, tableDoc.MealType)

	default:
		log.Printf("Unknown table operation: %s", op)
	}

	// Clear cache on success
	if s.cache != nil {
		s.cache.Clear()
		log.Printf("Cache cleared after processing table event")
	}

	return nil
}

func (s *SyncService) windowDays() int {
	if s.availabilityDays <= 0 {
		return 30
	}
	return s.availabilityDays
}

func (s *SyncService) indexTableWindow(ctx context.Context, table *TableDocument) error {
	now := time.Now()
	for i := 0; i < s.windowDays(); i++ {
		date := now.AddDate(0, 0, i).Format("2006-01-02")
		tableAvail := domain.NewTableAvailability(table.TableNumber, table.Capacity, table.MealType, date)
		tableAvail.IsAvailable = true
		if err := s.repo.Update(ctx, *tableAvail); err != nil {
			return fmt.Errorf("failed to index table availability for date %s: %w", date, err)
		}
	}
	return nil
}

func tableFromMetadata(metadata map[string]interface{}) *TableDocument {
	if metadata == nil {
		return nil
	}
	tableNumber := toInt(metadata["table_number"])
	mealType, _ := metadata["meal_type"].(string)
	capacity := toInt(metadata["capacity"])
	if tableNumber == 0 || mealType == "" {
		return nil
	}
	if capacity == 0 {
		capacity = 4
	}
	return &TableDocument{
		TableNumber: tableNumber,
		MealType:    mealType,
		Capacity:    capacity,
	}
}

func toInt(v interface{}) int {
	switch val := v.(type) {
	case float64:
		return int(val)
	case float32:
		return int(val)
	case int:
		return val
	case int64:
		return int(val)
	case json.Number:
		if i, err := val.Int64(); err == nil {
			return int(i)
		}
	}
	return 0
}

func (s *SyncService) resolveTableDoc(tableID string, metadata map[string]interface{}) (*TableDocument, error) {
	if doc := tableFromMetadata(metadata); doc != nil {
		return doc, nil
	}
	if s.tableClient == nil {
		return nil, fmt.Errorf("table metadata missing and table client not configured")
	}
	table, err := s.tableClient.GetTableByID(tableID)
	if err != nil {
		return nil, fmt.Errorf("failed to get table %s: %w", tableID, err)
	}
	return table, nil
}

func (s *SyncService) deleteTableEntries(ctx context.Context, table *TableDocument) error {
	if table == nil {
		return nil
	}
	query := fmt.Sprintf("table_number:%d AND meal_type:%q", table.TableNumber, table.MealType)
	if err := s.repo.DeleteByQuery(ctx, query); err != nil {
		return fmt.Errorf("failed to delete table %d (%s): %w", table.TableNumber, table.MealType, err)
	}
	return nil
}
