package service

import (
	"context"
	"crypto/sha1"
	"encoding/hex"
	"fmt"
	"time"

	"github.com/blassardoy/restaurant-reservas/search-api/internal/cache"
	"github.com/blassardoy/restaurant-reservas/search-api/internal/domain"
	"github.com/blassardoy/restaurant-reservas/search-api/internal/repository"
)

type SearchService interface {
	Search(ctx context.Context, q repository.SearchQuery) (*repository.SearchResult, error)
	GetByID(ctx context.Context, id string) (*repository.SearchResult, error)
}

type searchService struct {
	repo  repository.SearchRepository
	cache *cache.DualCache
}

func NewSearchService(repo repository.SearchRepository, cacheTTL time.Duration) SearchService {
	return &searchService{repo: repo, cache: cache.NewDual(cacheTTL)}
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

func (s *searchService) GetByID(ctx context.Context, id string) (*repository.SearchResult, error) {
	key := "id:" + id
	if v, ok := s.cache.Get(key); ok {
		if res, ok2 := v.(*repository.SearchResult); ok2 {
			return res, nil
		}
	}
	doc, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	res := &repository.SearchResult{Results: []domain.ReservationDocument{*doc}, Total: 1, Page: 1, Size: 1, Pages: 1}
	s.cache.Set(key, res)
	return res, nil
}

func cacheKey(q repository.SearchQuery) string {
	// Deterministic key from query
	s := fmt.Sprintf("q=%s|p=%d|s=%d|sort=%s|order=%s|f=%v", q.Q, q.Page, q.Size, q.Sort, q.Order, q.Filters)
	h := sha1.Sum([]byte(s))
	return hex.EncodeToString(h[:])
}

// no additional helpers yet
