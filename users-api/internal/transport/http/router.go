package http

import (
	"github.com/blassardoy/restaurant-reservas/users-api/internal/domain"
	"github.com/gin-gonic/gin"
)

// NewRouter builds and returns a gin.Engine with users routes mounted.
func NewRouter(service domain.UserService) *gin.Engine {
	r := gin.Default()

	// Health
	r.GET("/health", func(c *gin.Context) { c.Status(200) })

	// Users endpoints
	RegisterUserRoutes(r, service)

	return r
}
