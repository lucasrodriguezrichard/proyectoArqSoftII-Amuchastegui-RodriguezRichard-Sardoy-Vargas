// internal/transport/http/users_endpoints.go
package http

import (
	"net/http"

	"github.com/blassardoy/restaurant-reservas/users-api/internal/domain"
	"github.com/gin-gonic/gin"
)

// UsersHandler groups all user-related HTTP endpoints.
type UsersHandler struct {
	service domain.UserService
}

func NewUsersHandler(service domain.UserService) *UsersHandler {
	return &UsersHandler{service: service}
}

// RegisterUserRoutes attaches user endpoints under the given router/group.
func RegisterUserRoutes(r gin.IRouter, service domain.UserService) {
	h := NewUsersHandler(service)
	r.POST("/users/register", h.Register)
	r.POST("/users/login", h.Login)
}

type registerRequest struct {
	Username  string `json:"username"`
	Email     string `json:"email"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Password  string `json:"password"`
}

func (h *UsersHandler) Register(c *gin.Context) {
	var req registerRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid_request"})
		return
	}

	in := domain.RegisterInput{
		Username:  req.Username,
		Email:     req.Email,
		FirstName: req.FirstName,
		LastName:  req.LastName,
		Password:  req.Password,
	}

	user, err := h.service.Register(c.Request.Context(), in)
	if err != nil {
		switch err {
		case domain.ErrInvalidInput:
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid_input"})
			return
		case domain.ErrUserExists:
			c.JSON(http.StatusConflict, gin.H{"error": "user_already_exists"})
			return
		default:
			c.JSON(http.StatusInternalServerError, gin.H{"error": "internal_error"})
			return
		}
	}

	// Never include PasswordHash
	c.JSON(http.StatusCreated, gin.H{
		"user": gin.H{
			"id":         user.ID,
			"username":   user.Username,
			"email":      user.Email,
			"first_name": user.FirstName,
			"last_name":  user.LastName,
			"role":       user.Role,
			"created_at": user.CreatedAt,
			"updated_at": user.UpdatedAt,
		},
	})
}

type loginRequest struct {
	Identifier string `json:"identifier"`
	Password   string `json:"password"`
}

func (h *UsersHandler) Login(c *gin.Context) {
	var req loginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid_request"})
		return
	}

	in := domain.LoginInput{
		Identifier: req.Identifier,
		Password:   req.Password,
	}

	tokens, user, err := h.service.Login(c.Request.Context(), in)
	if err != nil {
		switch err {
		case domain.ErrInvalidInput, domain.ErrInvalidCredentials:
			c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid_credentials"})
			return
		default:
			c.JSON(http.StatusInternalServerError, gin.H{"error": "internal_error"})
			return
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"tokens": gin.H{
			"access_token":  tokens.AccessToken,
			"refresh_token": tokens.RefreshToken,
			"expires_at":    tokens.ExpiresAt,
		},
		"user": gin.H{
			"id":         user.ID,
			"username":   user.Username,
			"email":      user.Email,
			"first_name": user.FirstName,
			"last_name":  user.LastName,
			"role":       user.Role,
			"created_at": user.CreatedAt,
			"updated_at": user.UpdatedAt,
		},
	})
}
