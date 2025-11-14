package repository

import (
	"context"
	"fmt"

	"github.com/blassardoy/restaurant-reservas/search-api/internal/domain"
)

// SearchQuery represents incoming search parameters
type SearchQuery struct {
	Q       string
	Page    int
	Size    int
	Sort    string
	Order   string
	Filters map[string]string
}

// SearchResult represents table availability search results (MAIN ENTITY)
type SearchResult struct {
	Results []domain.TableAvailability
	Total   int
	Page    int
	Size    int
	Pages   int
}

// SearchRepository defines interactions with the search backend (Solr)
// It indexes and searches TableAvailability (the main entity)
type SearchRepository interface {
	Search(ctx context.Context, q SearchQuery) (*SearchResult, error)
	GetByID(ctx context.Context, id string) (*domain.TableAvailability, error)
	Index(ctx context.Context, doc domain.TableAvailability) error
	Update(ctx context.Context, doc domain.TableAvailability) error
	Delete(ctx context.Context, id string) error
}

// NoopRepository is a placeholder that returns empty results (to be replaced by Solr client)
type NoopRepository struct{}

func NewNoopRepository() *NoopRepository { return &NoopRepository{} }

func (r *NoopRepository) Search(ctx context.Context, q SearchQuery) (*SearchResult, error) {
	return &SearchResult{Results: []domain.TableAvailability{}, Total: 0, Page: q.Page, Size: q.Size, Pages: 0}, nil
}
func (r *NoopRepository) GetByID(ctx context.Context, id string) (*domain.TableAvailability, error) {
	return nil, fmt.Errorf("not found")
}
func (r *NoopRepository) Index(ctx context.Context, doc domain.TableAvailability) error { return nil }
func (r *NoopRepository) Update(ctx context.Context, doc domain.TableAvailability) error {
	return nil
}
func (r *NoopRepository) Delete(ctx context.Context, id string) error { return nil }
