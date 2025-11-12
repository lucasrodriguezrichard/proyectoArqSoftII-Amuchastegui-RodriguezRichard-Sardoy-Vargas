package http

import (
	"net/http"

	"github.com/blassardoy/restaurant-reservas/search-api/internal/repository"
	"github.com/blassardoy/restaurant-reservas/search-api/internal/service"
	"github.com/gin-gonic/gin"
)

func NewRouterWithService(svc service.SearchService) *gin.Engine {
	r := gin.Default()

	// CORS
	r.Use(func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
		if c.Request.Method == http.MethodOptions {
			c.AbortWithStatus(http.StatusNoContent)
			return
		}
		c.Next()
	})

	r.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "healthy"})
	})

	api := r.Group("/api")
	{
		api.GET("/search", func(c *gin.Context) {
			q := repository.SearchQuery{
				Q:     c.DefaultQuery("q", "*:*"),
				Page:  atoi(c.DefaultQuery("page", "1")),
				Size:  atoi(c.DefaultQuery("size", "10")),
				Sort:  c.DefaultQuery("sort", "created_at"),
				Order: c.DefaultQuery("order", "desc"),
				Filters: map[string]string{
					"meal_type": c.Query("meal_type"),
					"status":    c.Query("status"),
					"guests":    c.Query("guests"),
				},
			}
			res, err := svc.Search(c.Request.Context(), q)
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
				return
			}
			c.JSON(http.StatusOK, res)
		})
		api.GET("/search/:id", func(c *gin.Context) {
			id := c.Param("id")
			res, err := svc.GetByID(c.Request.Context(), id)
			if err != nil {
				c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
				return
			}
			c.JSON(http.StatusOK, res)
		})
		api.GET("/search/stats", func(c *gin.Context) {
			c.JSON(http.StatusOK, gin.H{"documents": 0, "cache_hit_rate": 0})
		})
		api.POST("/search/reindex", func(c *gin.Context) {
			c.JSON(http.StatusAccepted, gin.H{"status": "reindex started"})
		})
	}

	return r
}

func NewRouter() *gin.Engine { return NewRouterWithService(nil) }

func atoi(s string) int {
	n := 0
	for _, ch := range s {
		if ch < '0' || ch > '9' {
			return n
		}
		n = n*10 + int(ch-'0')
	}
	return n
}
