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
			if svc == nil {
				c.JSON(http.StatusServiceUnavailable, gin.H{"error": "service_unavailable"})
				return
			}
			page := max(1, atoi(c.DefaultQuery("page", "1")))
			size := clamp(atoi(c.DefaultQuery("size", "10")), 1, 100)
			q := repository.SearchQuery{
				Q:     c.DefaultQuery("q", "*:*"),
				Page:  page,
				Size:  size,
				Sort:  c.DefaultQuery("sort", ""),
				Order: c.DefaultQuery("order", "desc"),
				Filters: map[string]string{
					"meal_type":   c.Query("meal_type"),
					"is_available": c.Query("is_available"),
					"capacity":    c.Query("capacity"),
					"date":        c.Query("date"),
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
			if svc == nil {
				c.JSON(http.StatusServiceUnavailable, gin.H{"error": "service_unavailable"})
				return
			}
			id := c.Param("id")
			res, err := svc.GetByID(c.Request.Context(), id)
			if err != nil {
				c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
				return
			}
			c.JSON(http.StatusOK, res)
		})
		api.GET("/search/stats", func(c *gin.Context) {
			if svc == nil {
				c.JSON(http.StatusServiceUnavailable, gin.H{"error": "service_unavailable"})
				return
			}
			stats, err := svc.Stats(c.Request.Context())
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
				return
			}
			c.JSON(http.StatusOK, stats)
		})
		api.POST("/search/reindex", func(c *gin.Context) {
			if svc == nil {
				c.JSON(http.StatusServiceUnavailable, gin.H{"error": "service_unavailable"})
				return
			}
			go svc.Reindex(c.Request.Context())
			c.JSON(http.StatusAccepted, gin.H{"status": "reindex started"})
		})
	}

	// Cache introspection endpoints
	cache := r.Group("/__cache")
	{
		cache.GET("/stats", func(c *gin.Context) {
			if svc == nil {
				c.JSON(http.StatusServiceUnavailable, gin.H{"error": "service_unavailable"})
				return
			}
			stats, err := svc.Stats(c.Request.Context())
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
				return
			}
			c.JSON(http.StatusOK, stats.Cache)
		})
		cache.GET("/get", func(c *gin.Context) {
			if svc == nil {
				c.JSON(http.StatusServiceUnavailable, gin.H{"error": "service_unavailable"})
				return
			}
			key := c.Query("key")
			if key == "" {
				c.JSON(http.StatusBadRequest, gin.H{"error": "key parameter required"})
				return
			}
			val, ok := svc.GetCacheValue(key)
			if !ok {
				c.JSON(http.StatusNotFound, gin.H{"error": "key not found in cache"})
				return
			}
			c.JSON(http.StatusOK, gin.H{"key": key, "value": val})
		})
		cache.POST("/invalidate", func(c *gin.Context) {
			if svc == nil {
				c.JSON(http.StatusServiceUnavailable, gin.H{"error": "service_unavailable"})
				return
			}
			svc.InvalidateAll()
			c.JSON(http.StatusOK, gin.H{"status": "cache invalidated"})
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

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

func clamp(val, min, maxV int) int {
	if val < min {
		return min
	}
	if val > maxV {
		return maxV
	}
	return val
}
