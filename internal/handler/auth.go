package handler

import (
	"encoding/json"
	"errors"
	"log/slog"
	"net/http"

	"github.com/lxmp7p/ya-dipl1/internal/service"

	"github.com/go-chi/chi/v5"
)

type AuthRequest struct {
	Login    string `json:"login"`
	Password string `json:"password"`
}

type AuthResponse struct {
	SessionID string `json:"session_id"`
}

type AuthHandler struct {
	logger      *slog.Logger
	authService *service.AuthService
}

func (r *AuthRequest) Validate() error {
	if r.Login == "" {
		return errors.New("login is required")
	}
	if r.Password == "" {
		return errors.New("password is required")
	}
	return nil
}

func (auth *AuthHandler) AuthRoutes() chi.Router {
	r := chi.NewRouter()

	//r.Use(AuthMiddleware())

	r.Post("/register", auth.Registration)
	r.Post("/login", auth.Login)

	return r
}

func (auth *AuthHandler) Registration(w http.ResponseWriter, r *http.Request) {
	var request AuthRequest

	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	if err := request.Validate(); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	sessionID, err := auth.authService.Registration(r.Context(), request.Login, request.Password)

	if err != nil {
		switch {
		case errors.Is(err, service.ErrUserExists):
			auth.logger.Error(err.Error())
			http.Error(w, http.StatusText(http.StatusConflict), http.StatusConflict)

		default:
			auth.logger.Error(err.Error())
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		}
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Authorization", sessionID)

	http.SetCookie(w, &http.Cookie{
		Name:     "session_id",
		Value:    sessionID,
		HttpOnly: true,
		Path:     "/",
	})

	w.WriteHeader(http.StatusOK)

	json.NewEncoder(w).Encode(AuthResponse{
		SessionID: sessionID,
	})
}

func (auth *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	var request AuthRequest

	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	if err := request.Validate(); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	sessionID, err := auth.authService.Login(r.Context(), request.Login, request.Password)

	if err != nil {
		if errors.Is(err, service.ErrInvalidCredentials) {
			http.Error(w, "Invalid login or password", http.StatusUnauthorized)
			return
		}
		auth.logger.Error("Login failed", "error", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Authorization", sessionID)

	http.SetCookie(w, &http.Cookie{
		Name:     "session_id",
		Value:    sessionID,
		HttpOnly: true,
		Path:     "/",
	})

	w.WriteHeader(http.StatusOK)

	json.NewEncoder(w).Encode(AuthResponse{
		SessionID: sessionID,
	})
}
