package integration

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

type UsersClient struct {
	BaseURL string
	http    *http.Client
}

func NewUsersClient(base string) *UsersClient {
	return &UsersClient{
		BaseURL: base,
		http:    &http.Client{Timeout: 4 * time.Second},
	}
}

type userDTO struct {
	ID uint `json:"id"`
}

func (c *UsersClient) Exists(userID uint) (bool, error) {
	req, _ := http.NewRequest(http.MethodGet, fmt.Sprintf("%s/users/%d", c.BaseURL, userID), nil)
	resp, err := c.http.Do(req)
	if err != nil {
		return false, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return false, nil
	}

	var u userDTO
	if err := json.NewDecoder(resp.Body).Decode(&u); err != nil {
		return false, err
	}
	return u.ID == userID, nil
}
