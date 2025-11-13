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
}

type searchService struct {
	repo  repository.SearchRepository
	cache *cache.DualCache
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

func NewSearchService(repo repository.SearchRepository, cacheLayer *cache.DualCache) SearchService {
	return &searchService{repo: repo, cache: cacheLayer}
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

func (s *searchService) Reindex(ctx context.Context) error {
	s.InvalidateAll()
	// TODO: trigger background job to re-sync Solr.
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
