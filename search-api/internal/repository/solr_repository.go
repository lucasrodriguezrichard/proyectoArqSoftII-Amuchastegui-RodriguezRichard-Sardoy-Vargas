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

func (r *SolrRepository) Search(ctx context.Context, q SearchQuery) (*SearchResult, error) {
    if q.Page <= 0 { q.Page = 1 }
    if q.Size <= 0 { q.Size = 10 }
    start := (q.Page - 1) * q.Size

    sort := ""
    if q.Sort != "" {
        order := q.Order
        if order == "" { order = "desc" }
        sort = fmt.Sprintf("%s %s", q.Sort, order)
    }

    sr, err := r.c.Search(q.Q, start, q.Size, q.Filters, sort)
    if err != nil { return nil, err }

    docs := make([]domain.ReservationDocument, 0, len(sr.Response.Docs))
    for _, d := range sr.Response.Docs {
        doc := mapToReservation(d)
        docs = append(docs, doc)
    }

    pages := 0
    if q.Size > 0 { pages = (sr.Response.NumFound + q.Size - 1) / q.Size }
    return &SearchResult{Results: docs, Total: sr.Response.NumFound, Page: q.Page, Size: q.Size, Pages: pages}, nil
}

func (r *SolrRepository) GetByID(ctx context.Context, id string) (*domain.ReservationDocument, error) {
    q := SearchQuery{Q: fmt.Sprintf("id:%s", id), Page: 1, Size: 1}
    res, err := r.Search(ctx, q)
    if err != nil { return nil, err }
    if len(res.Results) == 0 { return nil, fmt.Errorf("not found") }
    return &res.Results[0], nil
}

func (r *SolrRepository) Index(ctx context.Context, doc domain.ReservationDocument) error  { return r.c.Index(doc) }
func (r *SolrRepository) Update(ctx context.Context, doc domain.ReservationDocument) error { return r.c.Update(doc) }
func (r *SolrRepository) Delete(ctx context.Context, id string) error                      { return r.c.Delete(id) }

func mapToReservation(m map[string]any) domain.ReservationDocument {
    var doc domain.ReservationDocument
    if v, ok := m[solr.FieldID].(string); ok { doc.ID = v }
    if v, ok := m[solr.FieldOwnerID].(string); ok { doc.OwnerID = v }
    if v, ok := m[solr.FieldTableNumber].(float64); ok { doc.TableNumber = int(v) }
    if v, ok := m[solr.FieldGuests].(float64); ok { doc.Guests = int(v) }
    if v, ok := m[solr.FieldMealType].(string); ok { doc.MealType = v }
    if v, ok := m[solr.FieldStatus].(string); ok { doc.Status = v }
    if v, ok := m[solr.FieldTotalPrice].(float64); ok { doc.TotalPrice = v }
    // parse RFC3339 times if present
    if v, ok := m[solr.FieldDateTime].(string); ok {
        if t, err := time.Parse(time.RFC3339, v); err == nil { doc.DateTime = t }
    }
    if v, ok := m[solr.FieldCreatedAt].(string); ok {
        if t, err := time.Parse(time.RFC3339, v); err == nil { doc.CreatedAt = t }
    }
    if v, ok := m[solr.FieldUpdatedAt].(string); ok {
        if t, err := time.Parse(time.RFC3339, v); err == nil { doc.UpdatedAt = t }
    }
    return doc
}

