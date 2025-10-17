package controllers

import (
	"net/http"
	"restaurant-system/internal/services"

	"github.com/gin-gonic/gin"
)

// OrderController handles order HTTP requests
type OrderController struct {
	orderService *services.OrderService
}

// NewOrderController creates a new order controller
func NewOrderController(orderService *services.OrderService) *OrderController {
	return &OrderController{
		orderService: orderService,
	}
}

// CreateOrder handles POST /orders
func (c *OrderController) CreateOrder(ctx *gin.Context) {
	var req services.CreateOrderRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	order, err := c.orderService.CreateOrder(req)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusCreated, order)
}

// GetOrder handles GET /orders/:id
func (c *OrderController) GetOrder(ctx *gin.Context) {
	id := ctx.Param("id")

	order, err := c.orderService.GetOrder(id)
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, order)
}

// GetOrdersByReservation handles GET /orders/reservation/:reservation_id
func (c *OrderController) GetOrdersByReservation(ctx *gin.Context) {
	reservationID := ctx.Param("reservation_id")

	orders, err := c.orderService.GetOrdersByReservation(reservationID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, orders)
}

// UpdateOrderStatus handles PUT /orders/:id/status
func (c *OrderController) UpdateOrderStatus(ctx *gin.Context) {
	id := ctx.Param("id")

	var req struct {
		Status string `json:"status" binding:"required"`
	}

	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err := c.orderService.UpdateOrderStatus(id, req.Status)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "Order status updated"})
}
