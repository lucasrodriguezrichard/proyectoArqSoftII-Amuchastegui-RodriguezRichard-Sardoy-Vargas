package service

import (
	"context"
	"crypto/sha1"
	"encoding/hex"
	"fmt"

	"github.com/blassardoy/restaurant-reservas/search-api/internal/cache"
	"github.com/blassardoy/restaurant-reservas/search-api/internal/domain"
	"github.com/blassardoy/restaurant-reservas/search-api/internal/repository"
)

type SearchService interface {
	Search(ctx context.Context, q repository.SearchQuery) (*repository.SearchResult, error)
	GetByID(ctx context.Context, id string) (*domain.ReservationDocument, error)
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
	key := cacheKey(q)
	if v, ok := s.cache.Get(key); ok {
		if res, ok2 := v.(*repository.SearchResult); ok2 {
			return res, nil
		}
	}
	res, err := s.repo.Search(ctx, q)
	if err != nil {
		return nil, err
	}
	s.cache.Set(key, res)
	return res, nil
}

func (s *searchService) GetByID(ctx context.Context, id string) (*domain.ReservationDocument, error) {
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
	res := &repository.SearchResult{Results: []domain.ReservationDocument{*doc}, Total: 1, Page: 1, Size: 1, Pages: 1}
	s.cache.Set(key, res)
	return doc, nil
}

func (s *searchService) Stats(ctx context.Context) (Stats, error) {
	local, hits, misses := s.cache.Stats()
	total := 0
	res, err := s.repo.Search(ctx, repository.SearchQuery{Q: "*:*", Page: 1, Size: 1})
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

	// Get all reservations from Reservations API
	if s.resClient == nil {
		return fmt.Errorf("reservation client not configured")
	}

	reservations, err := s.resClient.GetAllReservations()
	if err != nil {
		return fmt.Errorf("failed to fetch reservations: %w", err)
	}

	// Re-index all documents in Solr
	for _, doc := range reservations {
		if err := s.repo.Index(ctx, doc); err != nil {
			// Log error but continue with other docs
			fmt.Printf("Warning: failed to index doc %s: %v\n", doc.ID, err)
		}
	}

	return nil
}

func cacheKey(q repository.SearchQuery) string {
	// Deterministic key from query
	s := fmt.Sprintf("q=%s|p=%d|s=%d|sort=%s|order=%s|f=%v", q.Q, q.Page, q.Size, q.Sort, q.Order, q.Filters)
	h := sha1.Sum([]byte(s))
	return hex.EncodeToString(h[:])
}

func docCacheKey(id string) string {
	return "doc:" + id
}
