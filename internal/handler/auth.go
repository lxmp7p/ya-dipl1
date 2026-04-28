package handler

import (
	"encoding/json"
	"errors"
	"log/slog"
	"net/http"

	"github.com/lxmp7p/ya-dipl1/internal/service"

	"github.com/go-chi/chi/v5"
)

type RegisterRequest struct {
	Login    string `json:"login"`
	Password string `json:"password"`
}

type RegisterResponse struct {
	SessionID string `json:"session_id"`
}

type AuthHandler struct {
	logger      *slog.Logger
	authService service.AuthService
}

func (auth *AuthHandler) AuthRoutes() chi.Router {
	r := chi.NewRouter()

	//r.Use(AuthMiddleware())

	r.Post(DefaultApiRoute+"/user/register", auth.Registration)

	return r
}

func (auth *AuthHandler) Registration(w http.ResponseWriter, r *http.Request) {
	var request RegisterRequest

	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	if request.Login == "" || request.Password == "" {
		http.Error(w, "empty fields", http.StatusBadRequest)
		return
	}

	sessionID, err := auth.authService.Registration(r.Context(), request.Login, request.Password)

	if err != nil {
		switch {
		case errors.Is(err, service.ErrUserExists):
			http.Error(w, http.StatusText(http.StatusConflict), http.StatusConflict)

		default:
			auth.logger.Debug(err.Error())
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

	json.NewEncoder(w).Encode(RegisterResponse{
		SessionID: sessionID,
	})
}
