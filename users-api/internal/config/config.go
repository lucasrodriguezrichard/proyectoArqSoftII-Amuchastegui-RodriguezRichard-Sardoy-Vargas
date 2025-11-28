package config

import (
	"os"
	"strconv"
	"time"
)

// AppConfig centralizes runtime configuration for the users microservice.
type AppConfig struct {
	// MySQL
	DBHost string
	DBPort int
	DBUser string
	DBPass string
	DBName string

	// Auth
	JWTSecret     string
	JWTAccessTTL  time.Duration
	JWTRefreshTTL time.Duration

	// Admin seed
	AdminSeedEnabled bool
	AdminUsername    string
	AdminEmail       string
	AdminPassword    string
	AdminFirstName   string
	AdminLastName    string
}

func getenv(key, def string) string {
	v := os.Getenv(key)
	if v == "" {
		return def
	}
	return v
}

func getenvBool(key string, def bool) bool {
	v := os.Getenv(key)
	if v == "" {
		return def
	}
	b, err := strconv.ParseBool(v)
	if err != nil {
		return def
	}
	return b
}

func FromEnv() AppConfig {
	portStr := getenv("DB_PORT", "3306")
	port, err := strconv.Atoi(portStr)
	if err != nil {
		port = 3306
	}

	accessTTLStr := getenv("JWT_ACCESS_TTL", "1h")
	accessTTL, err := time.ParseDuration(accessTTLStr)
	if err != nil {
		accessTTL = time.Hour
	}
	refreshTTLStr := getenv("JWT_REFRESH_TTL", "720h") // 30d
	refreshTTL, err := time.ParseDuration(refreshTTLStr)
	if err != nil {
		refreshTTL = 720 * time.Hour
	}

	return AppConfig{
		DBHost:        getenv("DB_HOST", "localhost"),
		DBPort:        port,
		DBUser:        getenv("DB_USER", "root"),
		DBPass:        getenv("DB_PASS", "250498La"),
		DBName:        getenv("DB_NAME", "users"),
		JWTSecret:     getenv("JWT_SECRET", "dev-secret"),
		JWTAccessTTL:  accessTTL,
		JWTRefreshTTL: refreshTTL,

		AdminSeedEnabled: getenvBool("SEED_ADMIN", true),
		AdminUsername:    getenv("ADMIN_USERNAME", "admin"),
		AdminEmail:       getenv("ADMIN_EMAIL", "admin@restaurant.local"),
		AdminPassword:    getenv("ADMIN_PASSWORD", "Admin1234!"),
		AdminFirstName:   getenv("ADMIN_FIRST_NAME", "Default"),
		AdminLastName:    getenv("ADMIN_LAST_NAME", "Admin"),
	}
}
