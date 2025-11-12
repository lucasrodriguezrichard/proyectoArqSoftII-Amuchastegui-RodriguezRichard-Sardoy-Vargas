package config

import (
	"os"
)

type AppConfig struct {
	// MongoDB
	MongoURI        string
	MongoDB         string
	MongoCollection string

	// RabbitMQ
	RabbitMQURI      string
	RabbitMQExchange string
	RabbitMQQueue    string

	// Users API
	UsersAPIURL string

	// Server
	Port   string
	AppEnv string
}

func getenv(key, def string) string {
	v := os.Getenv(key)
	if v == "" {
		return def
	}
	return v
}

func FromEnv() AppConfig {
	return AppConfig{
		MongoURI:         getenv("MONGO_URI", "mongodb://localhost:27017"),
		MongoDB:          getenv("MONGO_DB", "reservations_db"),
		MongoCollection:  getenv("MONGO_COLLECTION", "reservations"),
		RabbitMQURI:      getenv("RABBITMQ_URI", "amqp://guest:guest@localhost:5672/"),
		RabbitMQExchange: getenv("RABBITMQ_EXCHANGE", "restaurant_events"),
		RabbitMQQueue:    getenv("RABBITMQ_QUEUE", "reservations_updates"),
		UsersAPIURL:      getenv("USERS_API_URL", "http://localhost:8080"),
		Port:             getenv("PORT", "8081"),
		AppEnv:           getenv("APP_ENV", "development"),
	}
}
