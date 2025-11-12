package http

import (
	"context"
	"net/http"
	"strings"

	"github.com/blassardoy/restaurant-reservas/users-api/internal/domain"
	"github.com/golang-jwt/jwt/v5"
)

// contextKey is a custom type for context keys to avoid collisions
type contextKey string

const (
	userIDKey   contextKey = "user_id"
	usernameKey contextKey = "username"
	roleKey     contextKey = "role"
)

// JWTMiddleware validates the JWT token and adds user info to context
type JWTMiddleware struct {
	jwtSecret []byte
}

func NewJWTMiddleware(secret string) *JWTMiddleware {
	return &JWTMiddleware{jwtSecret: []byte(secret)}
}

// Authenticate is the middleware function that validates JWT
func (m *JWTMiddleware) Authenticate(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Get token from Authorization header
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			http.Error(w, `{"error":"missing_authorization_header"}`, http.StatusUnauthorized)
			return
		}

		// Expected format: "Bearer <token>"
		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || strings.ToLower(parts[0]) != "bearer" {
			http.Error(w, `{"error":"invalid_authorization_format"}`, http.StatusUnauthorized)
			return
		}

		tokenString := parts[1]

		// Parse and validate token
		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			// Validate signing method
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, jwt.ErrSignatureInvalid
			}
			return m.jwtSecret, nil
		})

		if err != nil || !token.Valid {
			http.Error(w, `{"error":"invalid_token"}`, http.StatusUnauthorized)
			return
		}

		// Extract claims
		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			http.Error(w, `{"error":"invalid_token_claims"}`, http.StatusUnauthorized)
			return
		}

		// Extract user info from claims
		userID, ok := claims["sub"].(float64) // JWT numbers are float64
		if !ok {
			http.Error(w, `{"error":"invalid_user_id_in_token"}`, http.StatusUnauthorized)
			return
		}

		username, _ := claims["username"].(string)
		role, _ := claims["role"].(string)

		// Add user info to context
		ctx := context.WithValue(r.Context(), userIDKey, uint64(userID))
		ctx = context.WithValue(ctx, usernameKey, username)
		ctx = context.WithValue(ctx, roleKey, domain.Role(role))

		// Call next handler with updated context
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// RequireAdmin is a middleware that checks if the user is an admin
func (m *JWTMiddleware) RequireAdmin(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		role, ok := r.Context().Value(roleKey).(domain.Role)
		if !ok {
			http.Error(w, `{"error":"missing_role_in_context"}`, http.StatusForbidden)
			return
		}

		if role != domain.RoleAdmin {
			http.Error(w, `{"error":"admin_access_required"}`, http.StatusForbidden)
			return
		}

		next.ServeHTTP(w, r)
	})
}

// Helper functions to extract user info from context

// GetUserIDFromContext extracts the user ID from the request context
func GetUserIDFromContext(ctx context.Context) (uint64, bool) {
	userID, ok := ctx.Value(userIDKey).(uint64)
	return userID, ok
}

// GetUsernameFromContext extracts the username from the request context
func GetUsernameFromContext(ctx context.Context) (string, bool) {
	username, ok := ctx.Value(usernameKey).(string)
	return username, ok
}

// GetRoleFromContext extracts the role from the request context
func GetRoleFromContext(ctx context.Context) (domain.Role, bool) {
	role, ok := ctx.Value(roleKey).(domain.Role)
	return role, ok
}
