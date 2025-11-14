package service

import (
	"context"
	"crypto/sha1"
	"encoding/hex"
	"fmt"
	"sort"
	"strings"
	"time"

	"github.com/blassardoy/restaurant-reservas/search-api/internal/cache"
	"github.com/blassardoy/restaurant-reservas/search-api/internal/domain"
	"github.com/blassardoy/restaurant-reservas/search-api/internal/repository"
)

type SearchService interface {
	Search(ctx context.Context, q repository.SearchQuery) (*repository.SearchResult, error)
	GetByID(ctx context.Context, id string) (*domain.TableAvailability, error)
	Stats(ctx context.Context) (Stats, error)
	Reindex(ctx context.Context) error
	InvalidateAll()
	GetCacheValue(key string) (any, bool)
}

type searchService struct {
	repo      repository.SearchRepository
	cache     *cache.DualCache
	resClient *ReservationClient
}

type Stats struct {
	Documents int        `json:"documents"`
	Cache     CacheStats `json:"cache"`
}

type CacheStats struct {
	LocalEntries      int    `json:"local_entries"`
	DistributedHits   uint64 `json:"distributed_hits"`
	DistributedMisses uint64 `json:"distributed_misses"`
}

func NewSearchService(repo repository.SearchRepository, cacheLayer *cache.DualCache, resClient *ReservationClient) SearchService {
	return &searchService{repo: repo, cache: cacheLayer, resClient: resClient}
}

func (s *searchService) Search(ctx context.Context, q repository.SearchQuery) (*repository.SearchResult, error) {
	normalized := normalizeQuery(q)
	key := cacheKey(normalized)
	if v, ok := s.cache.Get(key); ok {
		if res, ok2 := v.(*repository.SearchResult); ok2 {
			return res, nil
		}
	}
	res, err := s.repo.Search(ctx, normalized)
	if err != nil {
		return nil, err
	}
	s.cache.Set(key, res)
	return res, nil
}

func (s *searchService) GetByID(ctx context.Context, id string) (*domain.TableAvailability, error) {
	key := docCacheKey(id)
	if v, ok := s.cache.Get(key); ok {
		if res, ok2 := v.(*repository.SearchResult); ok2 && len(res.Results) > 0 {
			doc := res.Results[0]
			return &doc, nil
		}
	}
	doc, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	res := &repository.SearchResult{Results: []domain.TableAvailability{*doc}, Total: 1, Page: 1, Size: 1, Pages: 1}
	s.cache.Set(key, res)
	return doc, nil
}

func (s *searchService) Stats(ctx context.Context) (Stats, error) {
	local, hits, misses := s.cache.Stats()
	total := 0
	res, err := s.repo.Search(ctx, normalizeQuery(repository.SearchQuery{Q: "*:*", Page: 1, Size: 1}))
	if err == nil && res != nil {
		total = res.Total
	}
	if err != nil {
		return Stats{}, err
	}
	return Stats{
		Documents: total,
		Cache: CacheStats{
			LocalEntries:      local,
			DistributedHits:   hits,
			DistributedMisses: misses,
		},
	}, nil
}

func (s *searchService) InvalidateAll() {
	if s.cache != nil {
		s.cache.Clear()
	}
}

func (s *searchService) GetCacheValue(key string) (any, bool) {
	if s.cache == nil {
		return nil, false
	}
	return s.cache.Get(key)
}

func (s *searchService) Reindex(ctx context.Context) error {
	// Clear cache
	s.InvalidateAll()

	// Get all existing reservations to mark tables as unavailable
	reservations, err := s.resClient.GetAllReservations()
	if err != nil {
		return fmt.Errorf("failed to get reservations: %w", err)
	}

	// Create a map of reserved tables: "date-mealtype-tablenumber" -> reservationID
	reservedTables := make(map[string]string)
	for _, res := range reservations {
		if res.Status != "cancelled" {
			date := res.DateTime.Format("2006-01-02")
			key := fmt.Sprintf("%s-%s-%d", date, res.MealType, res.TableNumber)
			reservedTables[key] = res.ID
		}
	}

	// Index all tables for the next 30 days
	now := time.Now()
	mealTypes := []string{"breakfast", "lunch", "dinner", "event"}
	capacities := map[string][]int{
		"breakfast": {2, 2, 4, 4, 4, 6, 6, 6, 8, 8},
		"lunch":     {2, 2, 4, 4, 4, 6, 6, 6, 8, 8},
		"dinner":    {2, 2, 4, 4, 4, 6, 6, 6, 8, 8},
		"event":     {8, 10, 10, 12, 12, 15, 15, 18, 20, 20},
	}

	indexed := 0
	for day := 0; day < 30; day++ {
		date := now.AddDate(0, 0, day).Format("2006-01-02")

		for _, mealType := range mealTypes {
			caps := capacities[mealType]

			for tableNum := 1; tableNum <= len(caps); tableNum++ {
				capacity := caps[tableNum-1]

				// Check if this table is reserved
				key := fmt.Sprintf("%s-%s-%d", date, mealType, tableNum)
				reservationID, isReserved := reservedTables[key]

				// Create TableAvailability document
				tableAvail := domain.NewTableAvailability(tableNum, capacity, mealType, date)
				tableAvail.IsAvailable = !isReserved
				if isReserved {
					tableAvail.ReservationID = reservationID
				}

				// Index in Solr
				if err := s.repo.Index(ctx, *tableAvail); err != nil {
					return fmt.Errorf("failed to index table %s: %w", tableAvail.ID, err)
				}
				indexed++
			}
		}
	}

	return nil
}

func cacheKey(q repository.SearchQuery) string {
	// Deterministic key from query
	s := fmt.Sprintf("q=%s|p=%d|s=%d|sort=%s|order=%s|f=%s", q.Q, q.Page, q.Size, q.Sort, q.Order, canonicalFilters(q.Filters))
	h := sha1.Sum([]byte(s))
	return hex.EncodeToString(h[:])
}

func docCacheKey(id string) string {
	return "doc:" + id
}

func canonicalFilters(filters map[string]string) string {
	if len(filters) == 0 {
		return ""
	}
	keys := make([]string, 0, len(filters))
	for k := range filters {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	parts := make([]string, 0, len(keys))
	for _, k := range keys {
		parts = append(parts, fmt.Sprintf("%s=%s", k, filters[k]))
	}
	return strings.Join(parts, ",")
}
