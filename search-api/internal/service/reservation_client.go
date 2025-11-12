package service

import (
    "encoding/json"
    "fmt"
    "net/http"
    "time"

    "github.com/blassardoy/restaurant-reservas/search-api/internal/domain"
)

type ReservationClient struct {
    baseURL string
    httpc   *http.Client
}

func NewReservationClient(baseURL string) *ReservationClient {
    return &ReservationClient{baseURL: baseURL, httpc: &http.Client{Timeout: 5 * time.Second}}
}

func (c *ReservationClient) GetReservationByID(id string) (*domain.ReservationDocument, error) {
    url := fmt.Sprintf("%s/api/reservations/%s", c.baseURL, id)
    resp, err := c.httpc.Get(url)
    if err != nil { return nil, err }
    defer resp.Body.Close()
    if resp.StatusCode != http.StatusOK {
        return nil, fmt.Errorf("reservations api status %d", resp.StatusCode)
    }
    var doc domain.ReservationDocument
    if err := json.NewDecoder(resp.Body).Decode(&doc); err != nil { return nil, err }
    return &doc, nil
}

