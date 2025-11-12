package controller

import (
	"net/http"
	"strconv"

	"github.com/blassardoy/restaurant-reservas/reservations-api/internal/domain"
	"github.com/blassardoy/restaurant-reservas/reservations-api/internal/service"
	"github.com/gin-gonic/gin"
)

type ReservationController struct {
	service service.ReservationService
}

func NewReservationController(service service.ReservationService) *ReservationController {
	return &ReservationController{service: service}
}

// CreateReservation handles POST /api/reservations
func (c *ReservationController) CreateReservation(ctx *gin.Context) {
	var req domain.CreateReservationRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	reservation, err := c.service.CreateReservation(ctx.Request.Context(), req)
	if err != nil {
		ctx.JSON(http.StatusUnprocessableEntity, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusCreated, reservation)
}

// GetReservation handles GET /api/reservations/:id
func (c *ReservationController) GetReservation(ctx *gin.Context) {
	id := ctx.Param("id")

	reservation, err := c.service.GetReservation(ctx.Request.Context(), id)
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, reservation)
}

// GetAllReservations handles GET /api/reservations
func (c *ReservationController) GetAllReservations(ctx *gin.Context) {
	limit, _ := strconv.Atoi(ctx.DefaultQuery("limit", "50"))
	offset, _ := strconv.Atoi(ctx.DefaultQuery("offset", "0"))

	reservations, err := c.service.GetAllReservations(ctx.Request.Context(), limit, offset)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, reservations)
}

// GetUserReservations handles GET /api/reservations/user/:user_id
func (c *ReservationController) GetUserReservations(ctx *gin.Context) {
	userID := ctx.Param("user_id")

	reservations, err := c.service.GetUserReservations(ctx.Request.Context(), userID)
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, reservations)
}

// UpdateReservation handles PUT /api/reservations/:id
func (c *ReservationController) UpdateReservation(ctx *gin.Context) {
	id := ctx.Param("id")

	var req domain.UpdateReservationRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	reservation, err := c.service.UpdateReservation(ctx.Request.Context(), id, req)
	if err != nil {
		ctx.JSON(http.StatusUnprocessableEntity, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, reservation)
}

// DeleteReservation handles DELETE /api/reservations/:id
func (c *ReservationController) DeleteReservation(ctx *gin.Context) {
	id := ctx.Param("id")

	if err := c.service.DeleteReservation(ctx.Request.Context(), id); err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	ctx.Status(http.StatusNoContent)
}

// ConfirmReservation handles POST /api/reservations/:id/confirm
func (c *ReservationController) ConfirmReservation(ctx *gin.Context) {
	id := ctx.Param("id")

	var req domain.ConfirmReservationRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	reservation, err := c.service.ConfirmReservation(ctx.Request.Context(), id, req)
	if err != nil {
		ctx.JSON(http.StatusUnprocessableEntity, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, reservation)
}
