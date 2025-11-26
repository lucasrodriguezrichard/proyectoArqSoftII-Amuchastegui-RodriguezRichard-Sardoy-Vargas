package controller

import (
	"net/http"

	"github.com/blassardoy/restaurant-reservas/reservations-api/internal/service"
	"github.com/gin-gonic/gin"
)

type TableController struct {
	service *service.TableService
}

func NewTableController(service *service.TableService) *TableController {
	return &TableController{service: service}
}

// CreateTable handles POST /api/tables
func (c *TableController) CreateTable(ctx *gin.Context) {
	var input service.CreateTableInput
	if err := ctx.ShouldBindJSON(&input); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	table, err := c.service.CreateTable(ctx.Request.Context(), input)
	if err != nil {
		ctx.JSON(http.StatusUnprocessableEntity, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusCreated, table)
}

// GetTable handles GET /api/tables/:id
func (c *TableController) GetTable(ctx *gin.Context) {
	id := ctx.Param("id")

	table, err := c.service.GetTableByID(ctx.Request.Context(), id)
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, table)
}

// GetAllTables handles GET /api/tables
func (c *TableController) GetAllTables(ctx *gin.Context) {
	mealType := ctx.Query("meal_type")

	var tables interface{}
	var err error

	if mealType != "" {
		tables, err = c.service.GetTablesByMealType(ctx.Request.Context(), mealType)
	} else {
		tables, err = c.service.GetAllTables(ctx.Request.Context())
	}

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, tables)
}

// UpdateTable handles PUT /api/tables/:id
func (c *TableController) UpdateTable(ctx *gin.Context) {
	id := ctx.Param("id")

	var input service.UpdateTableInput
	if err := ctx.ShouldBindJSON(&input); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	table, err := c.service.UpdateTable(ctx.Request.Context(), id, input)
	if err != nil {
		ctx.JSON(http.StatusUnprocessableEntity, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, table)
}

// DeleteTable handles DELETE /api/tables/:id
func (c *TableController) DeleteTable(ctx *gin.Context) {
	id := ctx.Param("id")

	if err := c.service.DeleteTable(ctx.Request.Context(), id); err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusNoContent, nil)
}
