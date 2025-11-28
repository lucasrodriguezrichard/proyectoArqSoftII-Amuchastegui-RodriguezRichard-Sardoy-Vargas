package main

import (
	"context"
	"log"
	"net/http"
	"strings"

	"github.com/blassardoy/restaurant-reservas/users-api/internal/auth"
	"github.com/blassardoy/restaurant-reservas/users-api/internal/config"
	"github.com/blassardoy/restaurant-reservas/users-api/internal/crypto"
	"github.com/blassardoy/restaurant-reservas/users-api/internal/db"
	"github.com/blassardoy/restaurant-reservas/users-api/internal/domain"
	"github.com/blassardoy/restaurant-reservas/users-api/internal/repository"
	servicepkg "github.com/blassardoy/restaurant-reservas/users-api/internal/service"
	httptransport "github.com/blassardoy/restaurant-reservas/users-api/internal/transport/http"
	"github.com/joho/godotenv"
)

func main() {
	// Load .env file (ignore error in production where env vars are set directly)
	_ = godotenv.Load()

	// Load configuration
	cfg := config.FromEnv()

	// DB connection
	gdb, err := db.New(db.Config{
		Host:     cfg.DBHost,
		Port:     cfg.DBPort,
		User:     cfg.DBUser,
		Password: cfg.DBPass,
		DBName:   cfg.DBName,
	})
	if err != nil {
		log.Fatalf("db connection error: %v", err)
	}

	// Auto-migrate user model (optional, safe on start)
	if err := gdb.AutoMigrate(&domain.User{}); err != nil {
		log.Fatalf("auto-migrate error: %v", err)
	}

	// Wire dependencies
	userRepo := repository.NewUserRepository(gdb)
	hasher := crypto.NewBcryptHasher(0)
	issuer := auth.NewJWTIssuer(cfg.JWTSecret, cfg.JWTAccessTTL, cfg.JWTRefreshTTL)
	service := servicepkg.NewUserService(userRepo, hasher, issuer)

	// Seed default admin so teammates can log into the admin panel without manual DB edits.
	seedDefaultAdmin(context.Background(), cfg, userRepo, service)

	// HTTP router with JWT secret for middleware
	r := httptransport.NewRouterWithConfig(service, cfg.JWTSecret)

	log.Println("users-api escuchando en :8080")
	log.Fatal(http.ListenAndServe(":8080", r))
}

func seedDefaultAdmin(ctx context.Context, cfg config.AppConfig, repo domain.UserRepository, svc *servicepkg.UserService) {
	if !cfg.AdminSeedEnabled {
		log.Println("admin seed disabled via config")
		return
	}

	if cfg.AdminUsername == "" || cfg.AdminEmail == "" || cfg.AdminPassword == "" {
		log.Println("admin seed skipped: missing credentials")
		return
	}

	exists, err := repo.ExistsByEmailOrUsername(ctx, strings.ToLower(cfg.AdminEmail), cfg.AdminUsername)
	if err != nil {
		log.Printf("admin seed: could not check existing user: %v", err)
		return
	}
	if exists {
		log.Printf("admin seed: user %s already exists, skipping", cfg.AdminUsername)
		return
	}

	input := domain.RegisterInput{
		Username:  cfg.AdminUsername,
		Email:     cfg.AdminEmail,
		FirstName: cfg.AdminFirstName,
		LastName:  cfg.AdminLastName,
		Password:  cfg.AdminPassword,
	}

	if _, err := svc.CreateAdmin(ctx, input); err != nil {
		log.Printf("admin seed: failed to create default admin: %v", err)
		return
	}

	log.Printf("admin seed: created default admin %s (%s)", cfg.AdminUsername, cfg.AdminEmail)
}
