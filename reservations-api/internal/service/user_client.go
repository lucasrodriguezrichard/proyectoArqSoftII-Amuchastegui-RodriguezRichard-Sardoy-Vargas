package service

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

// UserClient handles communication with the Users API
type UserClient struct {
	baseURL    string
	httpClient *http.Client
}

// UserResponse represents the response from Users API
type UserResponse struct {
	ID        uint64    `json:"id"`
	Username  string    `json:"username"`
	Email     string    `json:"email"`
	FirstName string    `json:"first_name"`
	LastName  string    `json:"last_name"`
	Role      string    `json:"role"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// NewUserClient creates a new Users API client
func NewUserClient(baseURL string) *UserClient {
	return &UserClient{
		baseURL: baseURL,
		httpClient: &http.Client{
			Timeout: 5 * time.Second,
		},
	}
}

// ValidateUser checks if a user exists by calling the Users API
func (c *UserClient) ValidateUser(userID string) error {
	url := fmt.Sprintf("%s/api/users/%s", c.baseURL, userID)

	resp, err := c.httpClient.Get(url)
	if err != nil {
		return fmt.Errorf("failed to call users API: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusNotFound {
		return fmt.Errorf("user not found")
	}

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("users API returned status %d", resp.StatusCode)
	}

	return nil
}

// GetUser retrieves user information from the Users API
func (c *UserClient) GetUser(userID string) (*UserResponse, error) {
	url := fmt.Sprintf("%s/api/users/%s", c.baseURL, userID)

	resp, err := c.httpClient.Get(url)
	if err != nil {
		return nil, fmt.Errorf("failed to call users API: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusNotFound {
		return nil, fmt.Errorf("user not found")
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("users API returned status %d", resp.StatusCode)
	}

	var user UserResponse
	if err := json.NewDecoder(resp.Body).Decode(&user); err != nil {
		return nil, fmt.Errorf("failed to decode user response: %w", err)
	}

	return &user, nil
}
