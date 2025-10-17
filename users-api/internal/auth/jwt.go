package auth

import (
	"time"

	"github.com/blassardoy/restaurant-reservas/users-api/internal/domain"
	"github.com/golang-jwt/jwt/v5"
)

// JWTIssuer issues JWT access and refresh tokens.
type JWTIssuer struct {
	secret     []byte
	accessTTL  time.Duration
	refreshTTL time.Duration
}

func NewJWTIssuer(secret string, accessTTL, refreshTTL time.Duration) *JWTIssuer {
	return &JWTIssuer{secret: []byte(secret), accessTTL: accessTTL, refreshTTL: refreshTTL}
}

func (i *JWTIssuer) IssueAccessToken(u domain.User) (string, time.Time, error) {
	exp := time.Now().Add(i.accessTTL)
	claims := jwt.MapClaims{
		"sub":      u.ID,
		"username": u.Username,
		"role":     string(u.Role),
		"exp":      exp.Unix(),
		"iat":      time.Now().Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signed, err := token.SignedString(i.secret)
	if err != nil {
		return "", time.Time{}, err
	}
	return signed, exp, nil
}

func (i *JWTIssuer) IssueRefreshToken(u domain.User) (string, error) {
	exp := time.Now().Add(i.refreshTTL)
	claims := jwt.MapClaims{
		"sub":  u.ID,
		"type": "refresh",
		"exp":  exp.Unix(),
		"iat":  time.Now().Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(i.secret)
}

var _ domain.TokenIssuer = (*JWTIssuer)(nil)
