package service

import (
	"encoding/json"
	"net/http"

	"github.com/google/uuid"
	"github.com/lxmp7p/ya-dipl1/internal/repository"
)

type AuthData struct {
	Login    string `json:"login"`
	Password string `json:"password"`
}

type AuthService struct {
	repo repository.Auth
}

func (auth *AuthService) Registration(w http.ResponseWriter, r *http.Request) {
	var request AuthData

	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	if request.Login == "" || request.Password == "" {
		http.Error(w, "empty fields", http.StatusBadRequest)
		return
	}

	hash, err := HashPassword(request.Password)
	if err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	err = auth.repo.Registration(r.Context(), request.Login, hash)
	if err != nil {
		if isUniqueViolation(err) {
			http.Error(w, http.StatusText(http.StatusConflict), http.StatusConflict)
			return
		}

		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	sessionID := uuid.NewString()
	err = auth.repo.CreateSession(r.Context(), request.Login, sessionID)
	if err != nil {
		http.Error(w, "internal error", http.StatusInternalServerError)
		return
	}

	http.SetCookie(w, &http.Cookie{
		Name:     "session_id",
		Value:    sessionID,
		HttpOnly: true,
		Path:     "/",
	})

	w.WriteHeader(http.StatusOK)
}
