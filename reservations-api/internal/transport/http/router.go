package http

import (
	"github.com/blassardoy/restaurant-reservas/reservations-api/internal/controller"
	"github.com/blassardoy/restaurant-reservas/reservations-api/internal/middleware"
	"github.com/gin-gonic/gin"
)

func NewRouter(reservationCtrl *controller.ReservationController, tableCtrl *controller.TableController) *gin.Engine {
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
			reservations.POST("", reservationCtrl.CreateReservation)
			reservations.GET("", reservationCtrl.GetAllReservations)
			reservations.GET("/:id", reservationCtrl.GetReservation)
			reservations.GET("/user/:user_id", reservationCtrl.GetUserReservations)
			reservations.PUT("/:id", reservationCtrl.UpdateReservation)
			reservations.DELETE("/:id", reservationCtrl.DeleteReservation)
			reservations.POST("/:id/confirm", reservationCtrl.ConfirmReservation)
		}

		// Public table routes
		tables := api.Group("/tables")
		{
			tables.GET("/available", reservationCtrl.GetAvailableTables)
			tables.GET("", tableCtrl.GetAllTables) // Public: list all tables
			tables.GET("/:id", tableCtrl.GetTable) // Public: get table by ID
		}

		// Admin-only table management routes
		adminTables := api.Group("/admin/tables")
		adminTables.Use(middleware.AuthMiddleware(), middleware.RequireAdmin())
		{
			adminTables.POST("", tableCtrl.CreateTable)
			adminTables.PUT("/:id", tableCtrl.UpdateTable)
			adminTables.DELETE("/:id", tableCtrl.DeleteTable)
		}
	}

	return r
}
