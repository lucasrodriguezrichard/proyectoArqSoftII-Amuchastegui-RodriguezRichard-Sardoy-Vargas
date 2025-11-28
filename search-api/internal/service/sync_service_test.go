package service

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/blassardoy/restaurant-reservas/search-api/internal/cache"
	"github.com/blassardoy/restaurant-reservas/search-api/internal/domain"
)

// Tests for SyncService

func TestSyncService_HandleEvent_Create(t *testing.T) {
	indexedDoc := domain.TableAvailability{}

	mockRepo := &mockSearchRepository{
		indexFunc: func(ctx context.Context, doc domain.TableAvailability) error {
			indexedDoc = doc
			return nil
		},
	}

	// Mock reservation server
	futureDate := time.Now().Add(24 * time.Hour)
	resServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		response := `{"id":"res1","owner_id":"1","table_number":1,"guests":4,"date_time":"` + futureDate.Format(time.RFC3339) + `","meal_type":"lunch","status":"confirmed","total_price":100.0}`
		w.Write([]byte(response))
	}))
	defer resServer.Close()

	tableServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`[{"id":"1","table_number":1,"capacity":4,"meal_type":"lunch"}]`))
	}))
	defer tableServer.Close()

	cacheLayer := cache.NewDual(0, nil, nil)
	resClient := NewReservationClient(resServer.URL)
	tableClient := NewTableClient(tableServer.URL)
	service := NewSyncService(mockRepo, resClient, tableClient, cacheLayer, 30)

	err := service.HandleEvent(context.Background(), "create", "res1")

	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if indexedDoc.IsAvailable {
		t.Error("expected table to be unavailable after reservation created")
	}

	if indexedDoc.ReservationID != "res1" {
		t.Errorf("expected reservation ID 'res1', got %s", indexedDoc.ReservationID)
	}
}

func TestSyncService_HandleEvent_Delete(t *testing.T) {
	updatedDoc := domain.TableAvailability{}

	mockRepo := &mockSearchRepository{
		updateFunc: func(ctx context.Context, doc domain.TableAvailability) error {
			updatedDoc = doc
			return nil
		},
	}

	futureDate := time.Now().Add(24 * time.Hour)
	resServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		response := `{"id":"res1","owner_id":"1","table_number":1,"guests":4,"date_time":"` + futureDate.Format(time.RFC3339) + `","meal_type":"lunch","status":"confirmed","total_price":100.0}`
		w.Write([]byte(response))
	}))
	defer resServer.Close()

	tableServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`[{"id":"1","table_number":1,"capacity":4,"meal_type":"lunch"}]`))
	}))
	defer tableServer.Close()

	cacheLayer := cache.NewDual(0, nil, nil)
	resClient := NewReservationClient(resServer.URL)
	tableClient := NewTableClient(tableServer.URL)
	service := NewSyncService(mockRepo, resClient, tableClient, cacheLayer, 30)

	err := service.HandleEvent(context.Background(), "delete", "res1")

	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if !updatedDoc.IsAvailable {
		t.Error("expected table to be available after reservation deleted")
	}

	if updatedDoc.ReservationID != "" {
		t.Errorf("expected empty reservation ID, got %s", updatedDoc.ReservationID)
	}
}

func TestSyncService_HandleEvent_Cancel(t *testing.T) {
	updatedDoc := domain.TableAvailability{}

	mockRepo := &mockSearchRepository{
		updateFunc: func(ctx context.Context, doc domain.TableAvailability) error {
			updatedDoc = doc
			return nil
		},
	}

	futureDate := time.Now().Add(24 * time.Hour)
	resServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		response := `{"id":"res1","owner_id":"1","table_number":1,"guests":4,"date_time":"` + futureDate.Format(time.RFC3339) + `","meal_type":"lunch","status":"cancelled","total_price":100.0}`
		w.Write([]byte(response))
	}))
	defer resServer.Close()

	cacheLayer := cache.NewDual(0, nil, nil)
	resClient := NewReservationClient(resServer.URL)
	service := NewSyncService(mockRepo, resClient, nil, cacheLayer, 30)

	err := service.HandleEvent(context.Background(), "cancel", "res1")

	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if !updatedDoc.IsAvailable {
		t.Error("expected table to be available after reservation cancelled")
	}
}

func TestSyncService_HandleEvent_Update_Cancelled(t *testing.T) {
	updatedDoc := domain.TableAvailability{}

	mockRepo := &mockSearchRepository{
		updateFunc: func(ctx context.Context, doc domain.TableAvailability) error {
			updatedDoc = doc
			return nil
		},
	}

	futureDate := time.Now().Add(24 * time.Hour)
	resServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		response := `{"id":"res1","owner_id":"1","table_number":1,"guests":4,"date_time":"` + futureDate.Format(time.RFC3339) + `","meal_type":"lunch","status":"cancelled","total_price":100.0}`
		w.Write([]byte(response))
	}))
	defer resServer.Close()

	cacheLayer := cache.NewDual(0, nil, nil)
	resClient := NewReservationClient(resServer.URL)
	service := NewSyncService(mockRepo, resClient, nil, cacheLayer, 30)

	err := service.HandleEvent(context.Background(), "update", "res1")

	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if !updatedDoc.IsAvailable {
		t.Error("expected table to be available after reservation status changed to cancelled")
	}
}

func TestSyncService_HandleEvent_Update_NotCancelled(t *testing.T) {
	updatedDoc := domain.TableAvailability{}

	mockRepo := &mockSearchRepository{
		updateFunc: func(ctx context.Context, doc domain.TableAvailability) error {
			updatedDoc = doc
			return nil
		},
	}

	futureDate := time.Now().Add(24 * time.Hour)
	resServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		response := `{"id":"res1","owner_id":"1","table_number":1,"guests":4,"date_time":"` + futureDate.Format(time.RFC3339) + `","meal_type":"lunch","status":"confirmed","total_price":100.0}`
		w.Write([]byte(response))
	}))
	defer resServer.Close()

	cacheLayer := cache.NewDual(0, nil, nil)
	resClient := NewReservationClient(resServer.URL)
	service := NewSyncService(mockRepo, resClient, nil, cacheLayer, 30)

	err := service.HandleEvent(context.Background(), "update", "res1")

	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if updatedDoc.IsAvailable {
		t.Error("expected table to be unavailable for confirmed reservation")
	}
}

func TestSyncService_HandleEvent_UnknownOperation(t *testing.T) {
	mockRepo := &mockSearchRepository{}

	futureDate := time.Now().Add(24 * time.Hour)
	resServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		response := `{"id":"res1","owner_id":"1","table_number":1,"guests":4,"date_time":"` + futureDate.Format(time.RFC3339) + `","meal_type":"lunch","status":"confirmed","total_price":100.0}`
		w.Write([]byte(response))
	}))
	defer resServer.Close()

	cacheLayer := cache.NewDual(0, nil, nil)
	resClient := NewReservationClient(resServer.URL)
	service := NewSyncService(mockRepo, resClient, nil, cacheLayer, 30)

	err := service.HandleEvent(context.Background(), "unknown", "res1")

	if err != nil {
		t.Fatalf("expected no error for unknown operation, got %v", err)
	}
}

func TestSyncService_HandleEvent_ReservationNotFound(t *testing.T) {
	mockRepo := &mockSearchRepository{}

	resServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
	}))
	defer resServer.Close()

	cacheLayer := cache.NewDual(0, nil, nil)
	resClient := NewReservationClient(resServer.URL)
	service := NewSyncService(mockRepo, resClient, nil, cacheLayer, 30)

	err := service.HandleEvent(context.Background(), "create", "nonexistent")

	if err == nil {
		t.Fatal("expected error for nonexistent reservation, got nil")
	}
}

func TestSyncService_HandleTableEvent_Create(t *testing.T) {
	indexedDocs := []domain.TableAvailability{}

	mockRepo := &mockSearchRepository{
		updateFunc: func(ctx context.Context, doc domain.TableAvailability) error {
			indexedDocs = append(indexedDocs, doc)
			return nil
		},
	}

	metadata := map[string]interface{}{
		"table_number": float64(5),
		"meal_type":    "dinner",
		"capacity":     float64(8),
	}

	cacheLayer := cache.NewDual(0, nil, nil)
	service := NewSyncService(mockRepo, nil, nil, cacheLayer, 7)

	err := service.HandleTableEvent(context.Background(), "create", "table1", metadata)

	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	// Should create 7 days of availability (default window)
	if len(indexedDocs) != 7 {
		t.Errorf("expected 7 indexed documents, got %d", len(indexedDocs))
	}

	// Check first document
	if len(indexedDocs) > 0 {
		if indexedDocs[0].TableNumber != 5 {
			t.Errorf("expected table number 5, got %d", indexedDocs[0].TableNumber)
		}
		if indexedDocs[0].Capacity != 8 {
			t.Errorf("expected capacity 8, got %d", indexedDocs[0].Capacity)
		}
		if indexedDocs[0].MealType != "dinner" {
			t.Errorf("expected meal type 'dinner', got %s", indexedDocs[0].MealType)
		}
		if !indexedDocs[0].IsAvailable {
			t.Error("expected new table to be available")
		}
	}
}

func TestSyncService_HandleTableEvent_Update(t *testing.T) {
	updatedDocs := []domain.TableAvailability{}

	mockRepo := &mockSearchRepository{
		updateFunc: func(ctx context.Context, doc domain.TableAvailability) error {
			updatedDocs = append(updatedDocs, doc)
			return nil
		},
	}

	metadata := map[string]interface{}{
		"table_number": float64(3),
		"meal_type":    "lunch",
		"capacity":     float64(6),
	}

	cacheLayer := cache.NewDual(0, nil, nil)
	service := NewSyncService(mockRepo, nil, nil, cacheLayer, 5)

	err := service.HandleTableEvent(context.Background(), "update", "table1", metadata)

	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if len(updatedDocs) != 5 {
		t.Errorf("expected 5 updated documents, got %d", len(updatedDocs))
	}
}

func TestSyncService_HandleTableEvent_Delete(t *testing.T) {
	deletedQuery := ""

	mockRepo := &mockSearchRepository{
		deleteByQueryFunc: func(ctx context.Context, query string) error {
			deletedQuery = query
			return nil
		},
	}

	metadata := map[string]interface{}{
		"table_number": float64(2),
		"meal_type":    "breakfast",
		"capacity":     float64(4),
	}

	cacheLayer := cache.NewDual(0, nil, nil)
	service := NewSyncService(mockRepo, nil, nil, cacheLayer, 30)

	err := service.HandleTableEvent(context.Background(), "delete", "table1", metadata)

	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	expectedQuery := `table_number:2 AND meal_type:"breakfast"`
	if deletedQuery != expectedQuery {
		t.Errorf("expected query '%s', got '%s'", expectedQuery, deletedQuery)
	}
}

func TestSyncService_HandleTableEvent_NoMetadata_WithClient(t *testing.T) {
	indexedDocs := []domain.TableAvailability{}

	mockRepo := &mockSearchRepository{
		updateFunc: func(ctx context.Context, doc domain.TableAvailability) error {
			indexedDocs = append(indexedDocs, doc)
			return nil
		},
	}

	tableServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"id":"table1","table_number":7,"capacity":10,"meal_type":"event"}`))
	}))
	defer tableServer.Close()

	cacheLayer := cache.NewDual(0, nil, nil)
	tableClient := NewTableClient(tableServer.URL)
	service := NewSyncService(mockRepo, nil, tableClient, cacheLayer, 3)

	err := service.HandleTableEvent(context.Background(), "create", "table1", nil)

	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if len(indexedDocs) != 3 {
		t.Errorf("expected 3 indexed documents, got %d", len(indexedDocs))
	}

	if len(indexedDocs) > 0 {
		if indexedDocs[0].TableNumber != 7 {
			t.Errorf("expected table number 7, got %d", indexedDocs[0].TableNumber)
		}
	}
}

func TestSyncService_HandleTableEvent_NoMetadata_NoClient(t *testing.T) {
	mockRepo := &mockSearchRepository{}
	cacheLayer := cache.NewDual(0, nil, nil)
	service := NewSyncService(mockRepo, nil, nil, cacheLayer, 30)

	err := service.HandleTableEvent(context.Background(), "create", "table1", nil)

	if err == nil {
		t.Fatal("expected error when no metadata and no client, got nil")
	}

	if !contains(err.Error(), "table metadata missing") {
		t.Errorf("expected 'table metadata missing' error, got %v", err)
	}
}

func TestSyncService_HandleTableEvent_Delete_NoMetadata(t *testing.T) {
	mockRepo := &mockSearchRepository{}
	cacheLayer := cache.NewDual(0, nil, nil)
	service := NewSyncService(mockRepo, nil, nil, cacheLayer, 30)

	// Delete should not fail even without metadata
	err := service.HandleTableEvent(context.Background(), "delete", "table1", nil)

	if err != nil {
		t.Fatalf("expected no error for delete without metadata, got %v", err)
	}
}

func TestSyncService_WindowDays_Default(t *testing.T) {
	service := &SyncService{
		availabilityDays: 0,
	}

	days := service.windowDays()

	if days != 30 {
		t.Errorf("expected default 30 days, got %d", days)
	}
}

func TestSyncService_WindowDays_Custom(t *testing.T) {
	service := &SyncService{
		availabilityDays: 90,
	}

	days := service.windowDays()

	if days != 90 {
		t.Errorf("expected 90 days, got %d", days)
	}
}

func TestSyncService_HandleEvent_ClearCache(t *testing.T) {
	mockRepo := &mockSearchRepository{
		indexFunc: func(ctx context.Context, doc domain.TableAvailability) error {
			return nil
		},
	}

	futureDate := time.Now().Add(24 * time.Hour)
	resServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		response := `{"id":"res1","owner_id":"1","table_number":1,"guests":4,"date_time":"` + futureDate.Format(time.RFC3339) + `","meal_type":"lunch","status":"confirmed","total_price":100.0}`
		w.Write([]byte(response))
	}))
	defer resServer.Close()

	cacheLayer := cache.NewDual(0, nil, nil)
	cacheLayer.Set("test-key", "test-value")

	resClient := NewReservationClient(resServer.URL)
	service := NewSyncService(mockRepo, resClient, nil, cacheLayer, 30)

	err := service.HandleEvent(context.Background(), "create", "res1")

	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	// Cache should be cleared
	if _, ok := cacheLayer.Get("test-key"); ok {
		t.Error("expected cache to be cleared after event")
	}
}

func TestToInt_Success(t *testing.T) {
	tests := []struct {
		input    interface{}
		expected int
	}{
		{float64(42), 42},
		{float32(10), 10},
		{int(5), 5},
		{int64(100), 100},
		{"invalid", 0},
		{nil, 0},
	}

	for _, test := range tests {
		result := toInt(test.input)
		if result != test.expected {
			t.Errorf("toInt(%v) = %d, expected %d", test.input, result, test.expected)
		}
	}
}

// Helper function
func contains(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr || len(s) > len(substr) && containsHelper(s, substr))
}

func containsHelper(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}
