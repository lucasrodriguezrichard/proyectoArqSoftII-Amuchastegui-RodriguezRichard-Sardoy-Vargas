package http

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"

	"github.com/blassardoy/restaurant-reservas/users-api/internal/domain"
)

type Router struct {
	svc domain.UserService
}

func NewRouter(svc domain.UserService) http.Handler {
	rt := &Router{svc: svc}

	r := chi.NewRouter()
	r.Use(middleware.RequestID, middleware.RealIP, middleware.Recoverer)

	// health
	r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	// grupo RESTAURANT
	r.Route("/restaurant", func(r chi.Router) {
		r.Post("/users/register", rt.register)
		r.Post("/users/login", rt.login)
	})

	return r
}

// request structs are defined in users_endpoints.go (Gin handlers)

func (rt *Router) register(w http.ResponseWriter, r *http.Request) {
	var in registerRequest
	if err := json.NewDecoder(r.Body).Decode(&in); err != nil {
		http.Error(w, "bad json", http.StatusBadRequest)
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
		http.Error(w, err.Error(), status)
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
		http.Error(w, "bad json", http.StatusBadRequest)
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
		http.Error(w, err.Error(), status)
		return
	}
	u.PasswordHash = ""
	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(struct {
		User   domain.User       `json:"user"`
		Tokens domain.AuthTokens `json:"tokens"`
	}{User: u, Tokens: tok})
}
