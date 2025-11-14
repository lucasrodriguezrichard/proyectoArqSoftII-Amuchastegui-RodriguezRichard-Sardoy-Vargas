package repository

import (
	"context"
	"fmt"
	"time"

	"github.com/blassardoy/restaurant-reservas/search-api/internal/domain"
	"github.com/blassardoy/restaurant-reservas/search-api/internal/solr"
)

type SolrRepository struct {
	c *solr.Client
}

func NewSolrRepository(c *solr.Client) *SolrRepository { return &SolrRepository{c: c} }

// Search for table availability (MAIN ENTITY)
func (r *SolrRepository) Search(ctx context.Context, q SearchQuery) (*SearchResult, error) {
	if q.Page <= 0 {
		q.Page = 1
	}
	if q.Size <= 0 {
		q.Size = 10
	}
	start := (q.Page - 1) * q.Size

	sort := ""
	if q.Sort != "" {
		order := q.Order
		if order == "" {
			order = "asc"
		}
		sort = fmt.Sprintf("%s %s", q.Sort, order)
	}

	sr, err := r.c.Search(q.Q, start, q.Size, q.Filters, sort)
	if err != nil {
		return nil, err
	}

	docs := make([]domain.TableAvailability, 0, len(sr.Response.Docs))
	for _, d := range sr.Response.Docs {
		doc := mapToTableAvailability(d)
		docs = append(docs, doc)
	}

	pages := 0
	if q.Size > 0 {
		pages = (sr.Response.NumFound + q.Size - 1) / q.Size
	}
	return &SearchResult{Results: docs, Total: sr.Response.NumFound, Page: q.Page, Size: q.Size, Pages: pages}, nil
}

func (r *SolrRepository) GetByID(ctx context.Context, id string) (*domain.TableAvailability, error) {
	q := SearchQuery{Q: fmt.Sprintf("id:%s", id), Page: 1, Size: 1}
	res, err := r.Search(ctx, q)
	if err != nil {
		return nil, err
	}
	if len(res.Results) == 0 {
		return nil, fmt.Errorf("not found")
	}
	return &res.Results[0], nil
}

func (r *SolrRepository) Index(ctx context.Context, doc domain.TableAvailability) error {
	return r.c.Index(doc)
}
func (r *SolrRepository) Update(ctx context.Context, doc domain.TableAvailability) error {
	return r.c.Update(doc)
}
func (r *SolrRepository) Delete(ctx context.Context, id string) error { return r.c.Delete(id) }

// mapToTableAvailability converts Solr document to TableAvailability
// Solr returns multivalued fields as arrays, so we need to handle both cases
func mapToTableAvailability(m map[string]any) domain.TableAvailability {
	var doc domain.TableAvailability

	// ID
	if v, ok := m[solr.FieldID].(string); ok {
		doc.ID = v
	}

	// TableNumber - can be float64 or array
	if v, ok := m[solr.FieldTableNumber].(float64); ok {
		doc.TableNumber = int(v)
	} else if arr, ok := m[solr.FieldTableNumber].([]interface{}); ok && len(arr) > 0 {
		if num, ok := arr[0].(float64); ok {
			doc.TableNumber = int(num)
		}
	}

	// Capacity - can be float64 or array
	if v, ok := m[solr.FieldCapacity].(float64); ok {
		doc.Capacity = int(v)
	} else if arr, ok := m[solr.FieldCapacity].([]interface{}); ok && len(arr) > 0 {
		if num, ok := arr[0].(float64); ok {
			doc.Capacity = int(num)
		}
	}

	// MealType - can be string or array
	if v, ok := m[solr.FieldMealType].(string); ok {
		doc.MealType = v
	} else if arr, ok := m[solr.FieldMealType].([]interface{}); ok && len(arr) > 0 {
		if str, ok := arr[0].(string); ok {
			doc.MealType = str
		}
	}

	// Date - can be string or array
	if v, ok := m[solr.FieldDate].(string); ok {
		doc.Date = v
	} else if arr, ok := m[solr.FieldDate].([]interface{}); ok && len(arr) > 0 {
		if str, ok := arr[0].(string); ok {
			// Extract just the date part (YYYY-MM-DD)
			if len(str) >= 10 {
				doc.Date = str[:10]
			} else {
				doc.Date = str
			}
		}
	}

	// IsAvailable - can be bool or array
	if v, ok := m[solr.FieldIsAvailable].(bool); ok {
		doc.IsAvailable = v
	} else if arr, ok := m[solr.FieldIsAvailable].([]interface{}); ok && len(arr) > 0 {
		if b, ok := arr[0].(bool); ok {
			doc.IsAvailable = b
		}
	}

	// ReservationID - can be string or array
	if v, ok := m[solr.FieldReservationID].(string); ok {
		doc.ReservationID = v
	} else if arr, ok := m[solr.FieldReservationID].([]interface{}); ok && len(arr) > 0 {
		if str, ok := arr[0].(string); ok {
			doc.ReservationID = str
		}
	}

	// CreatedAt - can be string or array
	var createdAtStr string
	if v, ok := m[solr.FieldCreatedAt].(string); ok {
		createdAtStr = v
	} else if arr, ok := m[solr.FieldCreatedAt].([]interface{}); ok && len(arr) > 0 {
		if str, ok := arr[0].(string); ok {
			createdAtStr = str
		}
	}
	if createdAtStr != "" {
		if t, err := time.Parse(time.RFC3339, createdAtStr); err == nil {
			doc.CreatedAt = t
		}
	}

	// UpdatedAt - can be string or array
	var updatedAtStr string
	if v, ok := m[solr.FieldUpdatedAt].(string); ok {
		updatedAtStr = v
	} else if arr, ok := m[solr.FieldUpdatedAt].([]interface{}); ok && len(arr) > 0 {
		if str, ok := arr[0].(string); ok {
			updatedAtStr = str
		}
	}
	if updatedAtStr != "" {
		if t, err := time.Parse(time.RFC3339, updatedAtStr); err == nil {
			doc.UpdatedAt = t
		}
	}

	return doc
}

