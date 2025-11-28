package service

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestUserClient_ValidateUser_Success(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/api/users/1" {
			t.Errorf("expected path '/api/users/1', got %s", r.URL.Path)
		}
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"id":1,"username":"testuser","email":"test@example.com"}`))
	}))
	defer server.Close()

	client := NewUserClient(server.URL)
	err := client.ValidateUser("1")

	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
}

func TestUserClient_ValidateUser_NotFound(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
	}))
	defer server.Close()

	client := NewUserClient(server.URL)
	err := client.ValidateUser("999")

	if err == nil {
		t.Fatal("expected error, got nil")
	}

	if err.Error() != "user not found" {
		t.Errorf("expected 'user not found' error, got %v", err)
	}
}

func TestUserClient_ValidateUser_ServerError(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	}))
	defer server.Close()

	client := NewUserClient(server.URL)
	err := client.ValidateUser("1")

	if err == nil {
		t.Fatal("expected error, got nil")
	}
}

func TestUserClient_GetUser_Success(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/api/users/1" {
			t.Errorf("expected path '/api/users/1', got %s", r.URL.Path)
		}
		w.WriteHeader(http.StatusOK)
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`{"id":1,"username":"testuser","email":"test@example.com","first_name":"Test","last_name":"User","role":"user"}`))
	}))
	defer server.Close()

	client := NewUserClient(server.URL)
	user, err := client.GetUser("1")

	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if user.ID != 1 {
		t.Errorf("expected user ID 1, got %d", user.ID)
	}

	if user.Username != "testuser" {
		t.Errorf("expected username 'testuser', got %s", user.Username)
	}

	if user.Email != "test@example.com" {
		t.Errorf("expected email 'test@example.com', got %s", user.Email)
	}
}

func TestUserClient_GetUser_NotFound(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
	}))
	defer server.Close()

	client := NewUserClient(server.URL)
	user, err := client.GetUser("999")

	if err == nil {
		t.Fatal("expected error, got nil")
	}

	if user != nil {
		t.Error("expected nil user")
	}

	if err.Error() != "user not found" {
		t.Errorf("expected 'user not found' error, got %v", err)
	}
}

func TestUserClient_GetUser_InvalidJSON(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`invalid json`))
	}))
	defer server.Close()

	client := NewUserClient(server.URL)
	user, err := client.GetUser("1")

	if err == nil {
		t.Fatal("expected error, got nil")
	}

	if user != nil {
		t.Error("expected nil user")
	}
}
