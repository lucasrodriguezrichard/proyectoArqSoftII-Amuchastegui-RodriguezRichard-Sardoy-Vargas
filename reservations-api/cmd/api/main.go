package main

import (
	"context"
	"log"
	"time"

	"github.com/blassardoy/restaurant-reservas/reservations-api/internal/config"
	"github.com/blassardoy/restaurant-reservas/reservations-api/internal/controller"
	"github.com/blassardoy/restaurant-reservas/reservations-api/internal/db"
	"github.com/blassardoy/restaurant-reservas/reservations-api/internal/repository"
	"github.com/blassardoy/restaurant-reservas/reservations-api/internal/service"
	httptransport "github.com/blassardoy/restaurant-reservas/reservations-api/internal/transport/http"
	"github.com/joho/godotenv"
)

func main() {
	// Load .env file
	_ = godotenv.Load()

	// Load configuration
	cfg := config.FromEnv()

	// Connect to MongoDB
	log.Println("Connecting to MongoDB...")
	client, collection, err := db.Connect(db.Config{
		URI:        cfg.MongoURI,
		Database:   cfg.MongoDB,
		Collection: cfg.MongoCollection,
	})
	if err != nil {
		log.Fatalf("MongoDB connection error: %v", err)
	}
	defer func() {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		if err := client.Disconnect(ctx); err != nil {
			log.Printf("Error disconnecting MongoDB: %v", err)
		}
	}()
	log.Println("Connected to MongoDB successfully")

	// Connect to RabbitMQ
	log.Println("Connecting to RabbitMQ...")
	rmqPublisher, err := service.NewRabbitMQPublisher(
		cfg.RabbitMQURI,
		cfg.RabbitMQExchange,
		cfg.RabbitMQQueue,
	)
	if err != nil {
		log.Fatalf("RabbitMQ connection error: %v", err)
	}
	defer rmqPublisher.Close()
	log.Println("Connected to RabbitMQ successfully")

	// Initialize layers
	repo := repository.NewMongoReservationRepository(collection)
	userClient := service.NewUserClient(cfg.UsersAPIURL)
	svc := service.NewReservationService(repo, userClient, rmqPublisher)
	ctrl := controller.NewReservationController(svc)

	// Setup HTTP router
	router := httptransport.NewRouter(ctrl)

	// Start server
	addr := ":" + cfg.Port
	log.Printf("reservations-api listening on %s", addr)
	if err := router.Run(addr); err != nil {
		log.Fatalf("Server error: %v", err)
	}
}
