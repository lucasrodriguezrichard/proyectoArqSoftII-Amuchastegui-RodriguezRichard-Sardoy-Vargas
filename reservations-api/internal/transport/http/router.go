package http

import (
	"github.com/blassardoy/restaurant-reservas/reservations-api/internal/controller"
	"github.com/gin-gonic/gin"
)

func NewRouter(ctrl *controller.ReservationController) *gin.Engine {
	r := gin.Default()

	// CORS middleware
	r.Use(func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	})

	// Health check
	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "healthy"})
	})

	// API routes
	api := r.Group("/api")
	{
		reservations := api.Group("/reservations")
		{
			reservations.POST("", ctrl.CreateReservation)
			reservations.GET("", ctrl.GetAllReservations)
			reservations.GET("/:id", ctrl.GetReservation)
			reservations.GET("/user/:user_id", ctrl.GetUserReservations)
			reservations.PUT("/:id", ctrl.UpdateReservation)
			reservations.DELETE("/:id", ctrl.DeleteReservation)
			reservations.POST("/:id/confirm", ctrl.ConfirmReservation)
		}

		tables := api.Group("/tables")
		{
			tables.GET("/available", ctrl.GetAvailableTables)
		}
	}

	return r
}
