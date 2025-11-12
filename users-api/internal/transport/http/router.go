package http

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"

	"github.com/blassardoy/restaurant-reservas/users-api/internal/domain"
)

type Router struct {
	svc            domain.UserService
	jwtSecret      string
	authMiddleware *JWTMiddleware
}

func NewRouter(svc domain.UserService) http.Handler {
	return NewRouterWithConfig(svc, "dev-secret")
}

func NewRouterWithConfig(svc domain.UserService, jwtSecret string) http.Handler {
	rt := &Router{
		svc:            svc,
		jwtSecret:      jwtSecret,
		authMiddleware: NewJWTMiddleware(jwtSecret),
	}

	r := chi.NewRouter()
	r.Use(middleware.RequestID, middleware.RealIP, middleware.Recoverer, middleware.Logger)

	// CORS middleware
	r.Use(corsMiddleware)

	// health
	r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_ = json.NewEncoder(w).Encode(map[string]string{"status": "healthy"})
	})

	// Public routes (no authentication required)
	r.Route("/api", func(r chi.Router) {
		r.Post("/users/register", rt.register)
		r.Post("/users/login", rt.login)
		r.Get("/users/{id}", rt.getUserByID) // Public for inter-service communication
	})

	// Protected routes (authentication required)
	r.Route("/api/admin", func(r chi.Router) {
		r.Use(rt.authMiddleware.Authenticate)
		r.Use(rt.authMiddleware.RequireAdmin)

		r.Post("/users", rt.createAdmin) // Create admin user
	})

	// Legacy routes (backward compatibility)
	r.Route("/restaurant", func(r chi.Router) {
		r.Post("/users/register", rt.register)
		r.Post("/users/login", rt.login)
	})

	return r
}

// CORS middleware
func corsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}

		next.ServeHTTP(w, r)
	})
}

// request structs are defined in users_endpoints.go (Gin handlers)

func (rt *Router) register(w http.ResponseWriter, r *http.Request) {
	var in registerRequest
	if err := json.NewDecoder(r.Body).Decode(&in); err != nil {
		respondError(w, http.StatusBadRequest, "invalid_request")
		return
	}
	req := domain.RegisterInput{
		Username:  in.Username,
		Email:     in.Email,
		FirstName: in.FirstName,
		LastName:  in.LastName,
		Password:  in.Password,
	}
	u, err := rt.svc.Register(r.Context(), req)
	if err != nil {
		status := http.StatusBadRequest
		switch err {
		case domain.ErrUserExists:
			status = http.StatusConflict
		case domain.ErrInvalidInput:
			status = http.StatusUnprocessableEntity
		default:
			status = http.StatusInternalServerError
		}
		respondError(w, status, err.Error())
		return
	}
	u.PasswordHash = "" // nunca exponer
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	_ = json.NewEncoder(w).Encode(u)
}

func (rt *Router) login(w http.ResponseWriter, r *http.Request) {
	var in loginRequest
	if err := json.NewDecoder(r.Body).Decode(&in); err != nil {
		respondError(w, http.StatusBadRequest, "invalid_request")
		return
	}
	tok, u, err := rt.svc.Login(r.Context(), domain.LoginInput{
		Identifier: in.Identifier,
		Password:   in.Password,
	})
	if err != nil {
		status := http.StatusUnauthorized
		if err == domain.ErrInvalidInput {
			status = http.StatusBadRequest
		}
		respondError(w, status, err.Error())
		return
	}
	u.PasswordHash = ""
	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(struct {
		User   domain.User       `json:"user"`
		Tokens domain.AuthTokens `json:"tokens"`
	}{User: u, Tokens: tok})
}

// getUserByID returns user information by ID (public endpoint for inter-service communication)
func (rt *Router) getUserByID(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		http.Error(w, `{"error":"invalid_user_id"}`, http.StatusBadRequest)
		return
	}

	user, err := rt.svc.GetByID(r.Context(), id)
	if err != nil {
		status := http.StatusInternalServerError
		if err == domain.ErrUserNotFound {
			status = http.StatusNotFound
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(status)
		_ = json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
		return
	}

	// Never expose password hash
	user.PasswordHash = ""
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(user)
}

// createAdmin creates a new admin user (requires admin authentication)
func (rt *Router) createAdmin(w http.ResponseWriter, r *http.Request) {
	var in registerRequest
	if err := json.NewDecoder(r.Body).Decode(&in); err != nil {
		respondError(w, http.StatusBadRequest, "invalid_request")
		return
	}

	req := domain.RegisterInput{
		Username:  in.Username,
		Email:     in.Email,
		FirstName: in.FirstName,
		LastName:  in.LastName,
		Password:  in.Password,
	}

	user, err := rt.svc.CreateAdmin(r.Context(), req)
	if err != nil {
		status := http.StatusBadRequest
		switch err {
		case domain.ErrUserExists:
			status = http.StatusConflict
		case domain.ErrInvalidInput:
			status = http.StatusUnprocessableEntity
		default:
			status = http.StatusInternalServerError
		}
		respondError(w, status, err.Error())
		return
	}

	// Never expose password hash
	user.PasswordHash = ""
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	_ = json.NewEncoder(w).Encode(user)
}

func respondError(w http.ResponseWriter, status int, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(map[string]string{"error": message})
}
