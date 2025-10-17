package main

import (
	"log"
	"restaurant-system/internal/config"
	"restaurant-system/internal/controllers"
	"restaurant-system/internal/dao"
	"restaurant-system/internal/middleware"
	"restaurant-system/internal/repository"
	"restaurant-system/internal/services"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	// Load environment variables
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, using environment variables")
	}

	// Load configuration
	cfg := config.LoadConfig()

	// Initialize database
	db, err := dao.NewDB(&cfg.Database)
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}
	defer db.Close()

	log.Println("Successfully connected to database")

	// Run migrations
	if err := dao.RunMigrations(db); err != nil {
		log.Fatal("Failed to run migrations:", err)
	}
	log.Println("Database migrations completed successfully")

	// Seed database with initial data (only in development)
	if cfg.Database.Host != "production-host" {
		if err := dao.SeedDatabase(db); err != nil {
			log.Println("Warning: Failed to seed database:", err)
		}
	}

	// Initialize DAOs
	reservationDAO := dao.NewReservationDAO(db)
	orderDAO := dao.NewOrderDAO(db)
	reviewDAO := dao.NewReviewDAO(db)
	ticketDAO := dao.NewTicketDAO(db)

	// Initialize repositories
	reservationRepo := repository.NewReservationRepository(reservationDAO)
	orderRepo := repository.NewOrderRepository(orderDAO)
	reviewRepo := repository.NewReviewRepository(reviewDAO)
	ticketRepo := repository.NewTicketRepository(ticketDAO)

	// Initialize services
	reservationService := services.NewReservationService(reservationRepo)
	orderService := services.NewOrderService(orderRepo)
	reviewService := services.NewReviewService(reviewRepo)
	ticketService := services.NewTicketService(ticketRepo, orderRepo)

	// Initialize controllers
	reservationController := controllers.NewReservationController(reservationService)
	orderController := controllers.NewOrderController(orderService)
	reviewController := controllers.NewReviewController(reviewService)
	ticketController := controllers.NewTicketController(ticketService)

	// Initialize Gin router
	router := gin.Default()

	// Add middleware
	router.Use(middleware.CORS())
	router.Use(middleware.Logger())

	// Health check endpoint
	router.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status":  "ok",
			"message": "Restaurant System API is running",
		})
	})

	// API routes
	api := router.Group("/api/v1")
	{
		// Reservation routes
		reservations := api.Group("/reservations")
		{
			reservations.POST("", reservationController.CreateReservation)
			reservations.GET("/:id", reservationController.GetReservation)
			reservations.GET("", reservationController.GetReservations)
			reservations.PUT("/:id/confirm", reservationController.ConfirmReservation)
			reservations.PUT("/:id/cancel", reservationController.CancelReservation)
			reservations.PUT("/:id/complete", reservationController.CompleteReservation)
		}

		// Order routes
		orders := api.Group("/orders")
		{
			orders.POST("", orderController.CreateOrder)
			orders.GET("/:id", orderController.GetOrder)
			orders.GET("/reservation/:reservation_id", orderController.GetOrdersByReservation)
			orders.PUT("/:id/status", orderController.UpdateOrderStatus)
		}

		// Review routes
		reviews := api.Group("/reviews")
		{
			reviews.POST("", reviewController.CreateReview)
			reviews.GET("/:id", reviewController.GetReview)
			reviews.GET("", reviewController.GetReviews)
			reviews.GET("/reservation/:reservation_id", reviewController.GetReviewsByReservation)
			reviews.GET("/stats/average", reviewController.GetAverageRatings)
			reviews.PUT("/:id", reviewController.UpdateReview)
			reviews.DELETE("/:id", reviewController.DeleteReview)
		}

		// Ticket routes
		tickets := api.Group("/tickets")
		{
			tickets.POST("", ticketController.GenerateTicket)
			tickets.GET("/:id", ticketController.GetTicket)
			tickets.GET("/order/:order_id", ticketController.GetTicketByOrder)
			tickets.GET("/reservation/:reservation_id", ticketController.GetTicketsByReservation)
			tickets.GET("/reports/sales", ticketController.GetSalesReport)
		}
	}

	// Admin routes (protected with auth middleware)
	admin := api.Group("/admin")
	admin.Use(middleware.AuthMiddleware(cfg.JWT.Secret))
	{
		// Admin reservation management
		admin.GET("/reservations", reservationController.GetReservations)

		// Admin order management
		admin.GET("/orders", orderController.GetOrdersByReservation)

		// Admin review management
		admin.GET("/reviews", reviewController.GetReviews)
		admin.DELETE("/reviews/:id", reviewController.DeleteReview)

		// Admin reports
		admin.GET("/reports/sales", ticketController.GetSalesReport)
		admin.GET("/reports/reviews", reviewController.GetAverageRatings)
	}

	// Start server
	log.Printf("🚀 Server starting on port %s", cfg.Server.Port)
	log.Printf("📚 API documentation available at http://localhost:%s/health", cfg.Server.Port)
	if err := router.Run(":" + cfg.Server.Port); err != nil {
		log.Fatal("Failed to start server:", err)
	}
}
