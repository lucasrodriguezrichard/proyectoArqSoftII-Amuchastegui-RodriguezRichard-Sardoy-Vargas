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
	createFunc                  func(ctx context.Context, u *domain.User) error
	getByIDFunc                 func(ctx context.Context, id uint64) (domain.User, error)
	getByEmailFunc              func(ctx context.Context, email string) (domain.User, error)
	getByUsernameFunc           func(ctx context.Context, username string) (domain.User, error)
	existsByEmailOrUsernameFunc func(ctx context.Context, email, username string) (bool, error)
	updatePasswordHashFunc      func(ctx context.Context, id uint64, newHash string) error
}

func (m *mockUserRepository) Create(ctx context.Context, u *domain.User) error {
	if m.createFunc != nil {
		return m.createFunc(ctx, u)
	}
	u.ID = 1
	return nil
}

func (m *mockUserRepository) GetByID(ctx context.Context, id uint64) (domain.User, error) {
	if m.getByIDFunc != nil {
		return m.getByIDFunc(ctx, id)
	}
	return domain.User{}, errors.New("not found")
}

func (m *mockUserRepository) GetByEmail(ctx context.Context, email string) (domain.User, error) {
	if m.getByEmailFunc != nil {
		return m.getByEmailFunc(ctx, email)
	}
	return domain.User{}, errors.New("not found")
}

func (m *mockUserRepository) GetByUsername(ctx context.Context, username string) (domain.User, error) {
	if m.getByUsernameFunc != nil {
		return m.getByUsernameFunc(ctx, username)
	}
	return domain.User{}, errors.New("not found")
}

func (m *mockUserRepository) ExistsByEmailOrUsername(ctx context.Context, email, username string) (bool, error) {
	if m.existsByEmailOrUsernameFunc != nil {
		return m.existsByEmailOrUsernameFunc(ctx, email, username)
	}
	return false, nil
}

func (m *mockUserRepository) UpdatePasswordHash(ctx context.Context, id uint64, newHash string) error {
	if m.updatePasswordHashFunc != nil {
		return m.updatePasswordHashFunc(ctx, id, newHash)
	}
	return nil
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
	issueAccessTokenFunc  func(u domain.User) (token string, exp time.Time, err error)
	issueRefreshTokenFunc func(u domain.User) (token string, err error)
}

func (m *mockTokenIssuer) IssueAccessToken(u domain.User) (token string, exp time.Time, err error) {
	if m.issueAccessTokenFunc != nil {
		return m.issueAccessTokenFunc(u)
	}
	return "access_token", time.Now().Add(15 * time.Minute), nil
}

func (m *mockTokenIssuer) IssueRefreshToken(u domain.User) (token string, err error) {
	if m.issueRefreshTokenFunc != nil {
		return m.issueRefreshTokenFunc(u)
	}
	return "refresh_token", nil
}

// Tests for Register

func TestUserService_Register_Success(t *testing.T) {
	mockRepo := &mockUserRepository{
		existsByEmailOrUsernameFunc: func(ctx context.Context, email, username string) (bool, error) {
			return false, nil
		},
		createFunc: func(ctx context.Context, u *domain.User) error {
			u.ID = 1
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

	if user.Role != domain.RoleUser {
		t.Errorf("expected role 'user', got %s", user.Role)
	}
}

func TestUserService_Register_InvalidInput(t *testing.T) {
	mockRepo := &mockUserRepository{}
	mockHasher := &mockPasswordHasher{}
	mockIssuer := &mockTokenIssuer{}

	service := NewUserService(mockRepo, mockHasher, mockIssuer)

	input := domain.RegisterInput{
		Username:  "ab", // Too short
		Email:     "invalid-email",
		FirstName: "Test",
		LastName:  "User",
		Password:  "123", // Too short
	}

	_, err := service.Register(context.Background(), input)

	if err != domain.ErrInvalidInput {
		t.Fatalf("expected ErrInvalidInput, got %v", err)
	}
}

func TestUserService_Register_UserExists(t *testing.T) {
	mockRepo := &mockUserRepository{
		existsByEmailOrUsernameFunc: func(ctx context.Context, email, username string) (bool, error) {
			return true, nil
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

func TestUserService_Register_HashError(t *testing.T) {
	mockRepo := &mockUserRepository{
		existsByEmailOrUsernameFunc: func(ctx context.Context, email, username string) (bool, error) {
			return false, nil
		},
	}
	mockHasher := &mockPasswordHasher{
		hashFunc: func(plain string) (string, error) {
			return "", errors.New("hash error")
		},
	}
	mockIssuer := &mockTokenIssuer{}

	service := NewUserService(mockRepo, mockHasher, mockIssuer)

	input := domain.RegisterInput{
		Username:  "testuser",
		Email:     "test@example.com",
		FirstName: "Test",
		LastName:  "User",
		Password:  "password123",
	}

	_, err := service.Register(context.Background(), input)

	if err == nil {
		t.Fatal("expected error, got nil")
	}
}

// Tests for Login

func TestUserService_Login_SuccessWithEmail(t *testing.T) {
	mockUser := domain.User{
		ID:           1,
		Username:     "testuser",
		Email:        "test@example.com",
		PasswordHash: "hashed_password123",
		Role:         domain.RoleUser,
	}

	mockRepo := &mockUserRepository{
		getByEmailFunc: func(ctx context.Context, email string) (domain.User, error) {
			if email == "test@example.com" {
				return mockUser, nil
			}
			return domain.User{}, errors.New("not found")
		},
	}
	mockHasher := &mockPasswordHasher{}
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

	if user.ID != 1 {
		t.Errorf("expected user ID 1, got %d", user.ID)
	}
}

func TestUserService_Login_SuccessWithUsername(t *testing.T) {
	mockUser := domain.User{
		ID:           1,
		Username:     "testuser",
		Email:        "test@example.com",
		PasswordHash: "hashed_password123",
		Role:         domain.RoleUser,
	}

	mockRepo := &mockUserRepository{
		getByUsernameFunc: func(ctx context.Context, username string) (domain.User, error) {
			if username == "testuser" {
				return mockUser, nil
			}
			return domain.User{}, errors.New("not found")
		},
	}
	mockHasher := &mockPasswordHasher{}
	mockIssuer := &mockTokenIssuer{}

	service := NewUserService(mockRepo, mockHasher, mockIssuer)

	input := domain.LoginInput{
		Identifier: "testuser",
		Password:   "password123",
	}

	tokens, user, err := service.Login(context.Background(), input)

	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if tokens.AccessToken == "" {
		t.Error("expected access token to be set")
	}

	if user.ID != 1 {
		t.Errorf("expected user ID 1, got %d", user.ID)
	}
}

func TestUserService_Login_InvalidInput(t *testing.T) {
	mockRepo := &mockUserRepository{}
	mockHasher := &mockPasswordHasher{}
	mockIssuer := &mockTokenIssuer{}

	service := NewUserService(mockRepo, mockHasher, mockIssuer)

	input := domain.LoginInput{
		Identifier: "",
		Password:   "",
	}

	_, _, err := service.Login(context.Background(), input)

	if err != domain.ErrInvalidInput {
		t.Fatalf("expected ErrInvalidInput, got %v", err)
	}
}

func TestUserService_Login_UserNotFound(t *testing.T) {
	mockRepo := &mockUserRepository{
		getByEmailFunc: func(ctx context.Context, email string) (domain.User, error) {
			return domain.User{}, errors.New("not found")
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

func TestUserService_Login_WrongPassword(t *testing.T) {
	mockUser := domain.User{
		ID:           1,
		Username:     "testuser",
		Email:        "test@example.com",
		PasswordHash: "hashed_correctpassword",
		Role:         domain.RoleUser,
	}

	mockRepo := &mockUserRepository{
		getByEmailFunc: func(ctx context.Context, email string) (domain.User, error) {
			return mockUser, nil
		},
	}
	mockHasher := &mockPasswordHasher{}
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

// Tests for GetByID

func TestUserService_GetByID_Success(t *testing.T) {
	mockUser := domain.User{
		ID:       1,
		Username: "testuser",
		Email:    "test@example.com",
		Role:     domain.RoleUser,
	}

	mockRepo := &mockUserRepository{
		getByIDFunc: func(ctx context.Context, id uint64) (domain.User, error) {
			if id == 1 {
				return mockUser, nil
			}
			return domain.User{}, errors.New("not found")
		},
	}
	mockHasher := &mockPasswordHasher{}
	mockIssuer := &mockTokenIssuer{}

	service := NewUserService(mockRepo, mockHasher, mockIssuer)

	user, err := service.GetByID(context.Background(), 1)

	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if user.ID != 1 {
		t.Errorf("expected user ID 1, got %d", user.ID)
	}
}

func TestUserService_GetByID_NotFound(t *testing.T) {
	mockRepo := &mockUserRepository{
		getByIDFunc: func(ctx context.Context, id uint64) (domain.User, error) {
			return domain.User{}, errors.New("not found")
		},
	}
	mockHasher := &mockPasswordHasher{}
	mockIssuer := &mockTokenIssuer{}

	service := NewUserService(mockRepo, mockHasher, mockIssuer)

	_, err := service.GetByID(context.Background(), 999)

	if err == nil {
		t.Fatal("expected error, got nil")
	}
}

// Tests for CreateAdmin

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

	if user.ID == 0 {
		t.Error("expected user ID to be set")
	}
}

func TestUserService_CreateAdmin_InvalidInput(t *testing.T) {
	mockRepo := &mockUserRepository{}
	mockHasher := &mockPasswordHasher{}
	mockIssuer := &mockTokenIssuer{}

	service := NewUserService(mockRepo, mockHasher, mockIssuer)

	input := domain.RegisterInput{
		Username:  "a", // Too short
		Email:     "invalid",
		FirstName: "Admin",
		LastName:  "User",
		Password:  "short",
	}

	_, err := service.CreateAdmin(context.Background(), input)

	if err != domain.ErrInvalidInput {
		t.Fatalf("expected ErrInvalidInput, got %v", err)
	}
}

func TestUserService_CreateAdmin_UserExists(t *testing.T) {
	mockRepo := &mockUserRepository{
		existsByEmailOrUsernameFunc: func(ctx context.Context, email, username string) (bool, error) {
			return true, nil
		},
	}
	mockHasher := &mockPasswordHasher{}
	mockIssuer := &mockTokenIssuer{}

	service := NewUserService(mockRepo, mockHasher, mockIssuer)

	input := domain.RegisterInput{
		Username:  "existingadmin",
		Email:     "existing@example.com",
		FirstName: "Admin",
		LastName:  "User",
		Password:  "adminpass123",
	}

	_, err := service.CreateAdmin(context.Background(), input)

	if err != domain.ErrUserExists {
		t.Fatalf("expected ErrUserExists, got %v", err)
	}
}
