package service

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"time"
)

type TableDocument struct {
	ID          string    `json:"id"`
	TableNumber int       `json:"table_number"`
	Capacity    int       `json:"capacity"`
	MealType    string    `json:"meal_type"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

type TableClient struct {
	baseURL string
	httpc   *http.Client
}

func NewTableClient(baseURL string) *TableClient {
	return &TableClient{baseURL: baseURL, httpc: &http.Client{Timeout: 5 * time.Second}}
}

func (c *TableClient) GetTableByID(id string) (*TableDocument, error) {
	url := fmt.Sprintf("%s/api/tables/%s", c.baseURL, id)
	resp, err := c.httpc.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("tables api status %d", resp.StatusCode)
	}
	var doc TableDocument
	if err := json.NewDecoder(resp.Body).Decode(&doc); err != nil {
		return nil, err
	}
	return &doc, nil
}

func (c *TableClient) GetAllTables() ([]TableDocument, error) {
	url := fmt.Sprintf("%s/api/tables", c.baseURL)
	resp, err := c.httpc.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("tables api status %d", resp.StatusCode)
	}
	var docs []TableDocument
	if err := json.NewDecoder(resp.Body).Decode(&docs); err != nil {
		return nil, err
	}
	return docs, nil
}

// GetTablesByMealType fetches tables filtered by meal type.
func (c *TableClient) GetTablesByMealType(mealType string) ([]TableDocument, error) {
	u := fmt.Sprintf("%s/api/tables", c.baseURL)
	if mealType != "" {
		u = fmt.Sprintf("%s?meal_type=%s", u, url.QueryEscape(mealType))
	}
	resp, err := c.httpc.Get(u)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("tables api status %d", resp.StatusCode)
	}
	var docs []TableDocument
	if err := json.NewDecoder(resp.Body).Decode(&docs); err != nil {
		return nil, err
	}
	return docs, nil
}

// FindTable searches for a table by meal type and table number.
func (c *TableClient) FindTable(mealType string, tableNumber int) (*TableDocument, error) {
	tables, err := c.GetTablesByMealType(mealType)
	if err != nil {
		return nil, err
	}
	for _, table := range tables {
		if table.TableNumber == tableNumber {
			t := table
			return &t, nil
		}
	}
	return nil, fmt.Errorf("table %d (%s) not found", tableNumber, mealType)
}
