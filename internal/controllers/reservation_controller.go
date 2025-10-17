package controllers

import (
	"net/http"
	"restaurant-system/internal/services"
	"time"

	"github.com/gin-gonic/gin"
)

// ReservationController handles reservation HTTP requests
type ReservationController struct {
	reservationService *services.ReservationService
}

// NewReservationController creates a new reservation controller
func NewReservationController(reservationService *services.ReservationService) *ReservationController {
	return &ReservationController{
		reservationService: reservationService,
	}
}

// CreateReservation handles POST /reservations
func (c *ReservationController) CreateReservation(ctx *gin.Context) {
	var req services.CreateReservationRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	reservation, err := c.reservationService.CreateReservation(req)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusCreated, reservation)
}

// GetReservation handles GET /reservations/:id
func (c *ReservationController) GetReservation(ctx *gin.Context) {
	id := ctx.Param("id")

	reservation, err := c.reservationService.GetReservation(id)
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, reservation)
}

// GetReservations handles GET /reservations
func (c *ReservationController) GetReservations(ctx *gin.Context) {
	startDateStr := ctx.Query("start_date")
	endDateStr := ctx.Query("end_date")

	if startDateStr == "" || endDateStr == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "start_date and end_date are required"})
		return
	}

	startDate, err := time.Parse("2006-01-02", startDateStr)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid start_date format"})
		return
	}

	endDate, err := time.Parse("2006-01-02", endDateStr)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid end_date format"})
		return
	}

	reservations, err := c.reservationService.GetReservationsByDateRange(startDate, endDate)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, reservations)
}

// ConfirmReservation handles PUT /reservations/:id/confirm
func (c *ReservationController) ConfirmReservation(ctx *gin.Context) {
	id := ctx.Param("id")

	err := c.reservationService.ConfirmReservation(id)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "Reservation confirmed"})
}

// CancelReservation handles PUT /reservations/:id/cancel
func (c *ReservationController) CancelReservation(ctx *gin.Context) {
	id := ctx.Param("id")

	err := c.reservationService.CancelReservation(id)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "Reservation cancelled"})
}

// CompleteReservation handles PUT /reservations/:id/complete
func (c *ReservationController) CompleteReservation(ctx *gin.Context) {
	id := ctx.Param("id")

	err := c.reservationService.CompleteReservation(id)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "Reservation completed"})
}
