package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	Port            string
	MongoURI        string
	MongoDB         string
	MongoCollection string
	UsersAPIURL     string
	JWTSecret       string
	RabbitURL       string
	RabbitExchange  string
	RabbitKind      string
}

func Load() *Config {
	_ = godotenv.Load()

	cfg := &Config{
		Port:            getenv("PORT", "8081"),
		MongoURI:        getenv("MONGO_URI", "mongodb://localhost:27017"),
		MongoDB:         getenv("MONGO_DB", "restaurant"),
		MongoCollection: getenv("MONGO_COLLECTION", "reservations"),
		UsersAPIURL:     getenv("USERS_API_URL", "http://localhost:8080"),
		JWTSecret:       getenv("JWT_SECRET", "changeme"),
		RabbitURL:       getenv("RABBIT_URL", ""),
		RabbitExchange:  getenv("RABBIT_EXCHANGE", "reservations"),
		RabbitKind:      getenv("RABBIT_KIND", "topic"),
	}

	if cfg.UsersAPIURL == "" || cfg.MongoURI == "" {
		log.Fatal("config inválida: faltan USERS_API_URL o MONGO_URI")
	}
	return cfg
}

func getenv(k, def string) string {
	if v := os.Getenv(k); v != "" {
		return v
	}
	return def
}
