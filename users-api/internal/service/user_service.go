// internal/service/user_service.go
package service

import (
	"context"
	"strings"

	"github.com/blassardoy/restaurant-reservas/users-api/internal/domain"
)

// UserService implements domain.UserService orchestrating repository, hashing and token issuance.
type UserService struct {
	userRepository domain.UserRepository
	passwordHasher domain.PasswordHasher
	tokenIssuer    domain.TokenIssuer
}

func NewUserService(repo domain.UserRepository, hasher domain.PasswordHasher, issuer domain.TokenIssuer) *UserService {
	return &UserService{
		userRepository: repo,
		passwordHasher: hasher,
		tokenIssuer:    issuer,
	}
}

// Register creates a new user with RoleUser. It does not issue tokens.
func (s *UserService) Register(ctx context.Context, in domain.RegisterInput) (domain.User, error) {
	in.Normalize()
	if err := in.Validate(); err != nil {
		return domain.User{}, domain.ErrInvalidInput
	}

	exists, err := s.userRepository.ExistsByEmailOrUsername(ctx, strings.ToLower(in.Email), in.Username)
	if err != nil {
		// Hide infra errors behind domain-level error when appropriate
		return domain.User{}, err
	}
	if exists {
		return domain.User{}, domain.ErrUserExists
	}

	hash, err := s.passwordHasher.Hash(in.Password)
	if err != nil {
		return domain.User{}, err
	}

	user := domain.NewUserFromRegister(in, hash)
	if err := s.userRepository.Create(ctx, &user); err != nil {
		return domain.User{}, err
	}
	return user, nil
}

// Login authenticates a user by email or username and returns tokens plus the user entity.
func (s *UserService) Login(ctx context.Context, in domain.LoginInput) (domain.AuthTokens, domain.User, error) {
	in.Normalize()
	if err := in.Validate(); err != nil {
		return domain.AuthTokens{}, domain.User{}, domain.ErrInvalidInput
	}

	var (
		user domain.User
		err  error
	)

	// Simple heuristic: identify email if it contains '@'; otherwise treat as username.
	if strings.Contains(in.Identifier, "@") {
		user, err = s.userRepository.GetByEmail(ctx, in.Identifier)
	} else {
		user, err = s.userRepository.GetByUsername(ctx, in.Identifier)
	}
	if err != nil {
		return domain.AuthTokens{}, domain.User{}, domain.ErrInvalidCredentials
	}

	if !s.passwordHasher.Compare(user.PasswordHash, in.Password) {
		return domain.AuthTokens{}, domain.User{}, domain.ErrInvalidCredentials
	}

	access, exp, err := s.tokenIssuer.IssueAccessToken(user)
	if err != nil {
		return domain.AuthTokens{}, domain.User{}, err
	}

	refresh, _ := s.tokenIssuer.IssueRefreshToken(user) // optional; ignore error if not implemented

	tokens := domain.AuthTokens{
		AccessToken:  access,
		RefreshToken: refresh,
		ExpiresAt:    exp,
	}
	return tokens, user, nil
}

// Ensure interface compliance at compile time.
var _ domain.UserService = (*UserService)(nil)
