package service

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/blassardoy/restaurant-reservas/users-api/internal/domain"
)

// Mock implementations

type mockUserRepository struct {
	getUserByIDFunc              func(ctx context.Context, id uint64) (domain.User, error)
	getUserByEmailFunc           func(ctx context.Context, email string) (domain.User, error)
	getUserByUsernameFunc        func(ctx context.Context, username string) (domain.User, error)
	createFunc                   func(ctx context.Context, u *domain.User) error
	existsByEmailOrUsernameFunc  func(ctx context.Context, email, username string) (bool, error)
}

func (m *mockUserRepository) GetByID(ctx context.Context, id uint64) (domain.User, error) {
	if m.getUserByIDFunc != nil {
		return m.getUserByIDFunc(ctx, id)
	}
	return domain.User{}, errors.New("not implemented")
}

func (m *mockUserRepository) GetByEmail(ctx context.Context, email string) (domain.User, error) {
	if m.getUserByEmailFunc != nil {
		return m.getUserByEmailFunc(ctx, email)
	}
	return domain.User{}, errors.New("not implemented")
}

func (m *mockUserRepository) GetByUsername(ctx context.Context, username string) (domain.User, error) {
	if m.getUserByUsernameFunc != nil {
		return m.getUserByUsernameFunc(ctx, username)
	}
	return domain.User{}, errors.New("not implemented")
}

func (m *mockUserRepository) Create(ctx context.Context, u *domain.User) error {
	if m.createFunc != nil {
		return m.createFunc(ctx, u)
	}
	return nil
}

func (m *mockUserRepository) UpdatePasswordHash(ctx context.Context, id uint64, newHash string) error {
	return nil
}

func (m *mockUserRepository) ExistsByEmailOrUsername(ctx context.Context, email, username string) (bool, error) {
	if m.existsByEmailOrUsernameFunc != nil {
		return m.existsByEmailOrUsernameFunc(ctx, email, username)
	}
	return false, nil
}

type mockPasswordHasher struct {
	hashFunc    func(plain string) (string, error)
	compareFunc func(hash string, plain string) bool
}

func (m *mockPasswordHasher) Hash(plain string) (string, error) {
	if m.hashFunc != nil {
		return m.hashFunc(plain)
	}
	return "hashed_" + plain, nil
}

func (m *mockPasswordHasher) Compare(hash string, plain string) bool {
	if m.compareFunc != nil {
		return m.compareFunc(hash, plain)
	}
	return hash == "hashed_"+plain
}

type mockTokenIssuer struct {
	issueAccessTokenFunc  func(u domain.User) (string, time.Time, error)
	issueRefreshTokenFunc func(u domain.User) (string, error)
}

func (m *mockTokenIssuer) IssueAccessToken(u domain.User) (string, time.Time, error) {
	if m.issueAccessTokenFunc != nil {
		return m.issueAccessTokenFunc(u)
	}
	exp := time.Now().Add(time.Hour)
	return "access_token_" + u.Username, exp, nil
}

func (m *mockTokenIssuer) IssueRefreshToken(u domain.User) (string, error) {
	if m.issueRefreshTokenFunc != nil {
		return m.issueRefreshTokenFunc(u)
	}
	return "refresh_token_" + u.Username, nil
}

// Tests

func TestUserService_Register_Success(t *testing.T) {
	mockRepo := &mockUserRepository{
		existsByEmailOrUsernameFunc: func(ctx context.Context, email, username string) (bool, error) {
			return false, nil // User doesn't exist
		},
		createFunc: func(ctx context.Context, u *domain.User) error {
			u.ID = 1 // Simulate DB auto-increment
			return nil
		},
	}

	mockHasher := &mockPasswordHasher{}
	mockIssuer := &mockTokenIssuer{}

	service := NewUserService(mockRepo, mockHasher, mockIssuer)

	input := domain.RegisterInput{
		Username:  "testuser",
		Email:     "test@example.com",
		FirstName: "Test",
		LastName:  "User",
		Password:  "password123",
	}

	user, err := service.Register(context.Background(), input)

	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if user.ID == 0 {
		t.Error("expected user ID to be set")
	}

	if user.Username != "testuser" {
		t.Errorf("expected username 'testuser', got %s", user.Username)
	}

	if user.Email != "test@example.com" {
		t.Errorf("expected email 'test@example.com', got %s", user.Email)
	}

	if user.Role != domain.RoleUser {
		t.Errorf("expected role 'user', got %s", user.Role)
	}

	if user.PasswordHash != "hashed_password123" {
		t.Errorf("expected password to be hashed, got %s", user.PasswordHash)
	}
}

func TestUserService_Register_UserExists(t *testing.T) {
	mockRepo := &mockUserRepository{
		existsByEmailOrUsernameFunc: func(ctx context.Context, email, username string) (bool, error) {
			return true, nil // User already exists
		},
	}

	mockHasher := &mockPasswordHasher{}
	mockIssuer := &mockTokenIssuer{}

	service := NewUserService(mockRepo, mockHasher, mockIssuer)

	input := domain.RegisterInput{
		Username:  "existinguser",
		Email:     "existing@example.com",
		FirstName: "Test",
		LastName:  "User",
		Password:  "password123",
	}

	_, err := service.Register(context.Background(), input)

	if err != domain.ErrUserExists {
		t.Fatalf("expected ErrUserExists, got %v", err)
	}
}

func TestUserService_Register_InvalidInput(t *testing.T) {
	mockRepo := &mockUserRepository{}
	mockHasher := &mockPasswordHasher{}
	mockIssuer := &mockTokenIssuer{}

	service := NewUserService(mockRepo, mockHasher, mockIssuer)

	testCases := []struct {
		name  string
		input domain.RegisterInput
	}{
		{
			name: "short password",
			input: domain.RegisterInput{
				Username:  "testuser",
				Email:     "test@example.com",
				FirstName: "Test",
				LastName:  "User",
				Password:  "short",
			},
		},
		{
			name: "invalid email",
			input: domain.RegisterInput{
				Username:  "testuser",
				Email:     "invalid-email",
				FirstName: "Test",
				LastName:  "User",
				Password:  "password123",
			},
		},
		{
			name: "invalid username",
			input: domain.RegisterInput{
				Username:  "ab", // too short
				Email:     "test@example.com",
				FirstName: "Test",
				LastName:  "User",
				Password:  "password123",
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			_, err := service.Register(context.Background(), tc.input)
			if err != domain.ErrInvalidInput {
				t.Errorf("expected ErrInvalidInput, got %v", err)
			}
		})
	}
}

func TestUserService_Login_Success_WithEmail(t *testing.T) {
	expectedUser := domain.User{
		ID:           1,
		Username:     "testuser",
		Email:        "test@example.com",
		PasswordHash: "hashed_password123",
		Role:         domain.RoleUser,
	}

	mockRepo := &mockUserRepository{
		getUserByEmailFunc: func(ctx context.Context, email string) (domain.User, error) {
			return expectedUser, nil
		},
	}

	mockHasher := &mockPasswordHasher{
		compareFunc: func(hash string, plain string) bool {
			return hash == "hashed_password123" && plain == "password123"
		},
	}

	mockIssuer := &mockTokenIssuer{}

	service := NewUserService(mockRepo, mockHasher, mockIssuer)

	input := domain.LoginInput{
		Identifier: "test@example.com",
		Password:   "password123",
	}

	tokens, user, err := service.Login(context.Background(), input)

	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if tokens.AccessToken == "" {
		t.Error("expected access token to be set")
	}

	if user.ID != expectedUser.ID {
		t.Errorf("expected user ID %d, got %d", expectedUser.ID, user.ID)
	}

	if user.Username != expectedUser.Username {
		t.Errorf("expected username %s, got %s", expectedUser.Username, user.Username)
	}
}

func TestUserService_Login_Success_WithUsername(t *testing.T) {
	expectedUser := domain.User{
		ID:           1,
		Username:     "testuser",
		Email:        "test@example.com",
		PasswordHash: "hashed_password123",
		Role:         domain.RoleUser,
	}

	mockRepo := &mockUserRepository{
		getUserByUsernameFunc: func(ctx context.Context, username string) (domain.User, error) {
			return expectedUser, nil
		},
	}

	mockHasher := &mockPasswordHasher{
		compareFunc: func(hash string, plain string) bool {
			return hash == "hashed_password123" && plain == "password123"
		},
	}

	mockIssuer := &mockTokenIssuer{}

	service := NewUserService(mockRepo, mockHasher, mockIssuer)

	input := domain.LoginInput{
		Identifier: "testuser", // No @ sign, so treated as username
		Password:   "password123",
	}

	tokens, user, err := service.Login(context.Background(), input)

	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if tokens.AccessToken == "" {
		t.Error("expected access token to be set")
	}

	if user.Username != expectedUser.Username {
		t.Errorf("expected username %s, got %s", expectedUser.Username, user.Username)
	}
}

func TestUserService_Login_InvalidCredentials_WrongPassword(t *testing.T) {
	expectedUser := domain.User{
		ID:           1,
		Username:     "testuser",
		Email:        "test@example.com",
		PasswordHash: "hashed_password123",
		Role:         domain.RoleUser,
	}

	mockRepo := &mockUserRepository{
		getUserByEmailFunc: func(ctx context.Context, email string) (domain.User, error) {
			return expectedUser, nil
		},
	}

	mockHasher := &mockPasswordHasher{
		compareFunc: func(hash string, plain string) bool {
			return false // Password doesn't match
		},
	}

	mockIssuer := &mockTokenIssuer{}

	service := NewUserService(mockRepo, mockHasher, mockIssuer)

	input := domain.LoginInput{
		Identifier: "test@example.com",
		Password:   "wrongpassword",
	}

	_, _, err := service.Login(context.Background(), input)

	if err != domain.ErrInvalidCredentials {
		t.Fatalf("expected ErrInvalidCredentials, got %v", err)
	}
}

func TestUserService_Login_InvalidCredentials_UserNotFound(t *testing.T) {
	mockRepo := &mockUserRepository{
		getUserByEmailFunc: func(ctx context.Context, email string) (domain.User, error) {
			return domain.User{}, domain.ErrUserNotFound
		},
	}

	mockHasher := &mockPasswordHasher{}
	mockIssuer := &mockTokenIssuer{}

	service := NewUserService(mockRepo, mockHasher, mockIssuer)

	input := domain.LoginInput{
		Identifier: "nonexistent@example.com",
		Password:   "password123",
	}

	_, _, err := service.Login(context.Background(), input)

	if err != domain.ErrInvalidCredentials {
		t.Fatalf("expected ErrInvalidCredentials, got %v", err)
	}
}

func TestUserService_GetByID_Success(t *testing.T) {
	expectedUser := domain.User{
		ID:       1,
		Username: "testuser",
		Email:    "test@example.com",
		Role:     domain.RoleUser,
	}

	mockRepo := &mockUserRepository{
		getUserByIDFunc: func(ctx context.Context, id uint64) (domain.User, error) {
			if id == 1 {
				return expectedUser, nil
			}
			return domain.User{}, domain.ErrUserNotFound
		},
	}

	mockHasher := &mockPasswordHasher{}
	mockIssuer := &mockTokenIssuer{}

	service := NewUserService(mockRepo, mockHasher, mockIssuer)

	user, err := service.GetByID(context.Background(), 1)

	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if user.ID != expectedUser.ID {
		t.Errorf("expected user ID %d, got %d", expectedUser.ID, user.ID)
	}

	if user.Username != expectedUser.Username {
		t.Errorf("expected username %s, got %s", expectedUser.Username, user.Username)
	}
}

func TestUserService_GetByID_NotFound(t *testing.T) {
	mockRepo := &mockUserRepository{
		getUserByIDFunc: func(ctx context.Context, id uint64) (domain.User, error) {
			return domain.User{}, domain.ErrUserNotFound
		},
	}

	mockHasher := &mockPasswordHasher{}
	mockIssuer := &mockTokenIssuer{}

	service := NewUserService(mockRepo, mockHasher, mockIssuer)

	_, err := service.GetByID(context.Background(), 999)

	if err != domain.ErrUserNotFound {
		t.Fatalf("expected ErrUserNotFound, got %v", err)
	}
}

func TestUserService_CreateAdmin_Success(t *testing.T) {
	mockRepo := &mockUserRepository{
		existsByEmailOrUsernameFunc: func(ctx context.Context, email, username string) (bool, error) {
			return false, nil
		},
		createFunc: func(ctx context.Context, u *domain.User) error {
			u.ID = 2
			return nil
		},
	}

	mockHasher := &mockPasswordHasher{}
	mockIssuer := &mockTokenIssuer{}

	service := NewUserService(mockRepo, mockHasher, mockIssuer)

	input := domain.RegisterInput{
		Username:  "adminuser",
		Email:     "admin@example.com",
		FirstName: "Admin",
		LastName:  "User",
		Password:  "adminpass123",
	}

	user, err := service.CreateAdmin(context.Background(), input)

	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if user.Role != domain.RoleAdmin {
		t.Errorf("expected role 'admin', got %s", user.Role)
	}

	if user.Username != "adminuser" {
		t.Errorf("expected username 'adminuser', got %s", user.Username)
	}
}
