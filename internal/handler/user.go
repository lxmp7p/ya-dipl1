package handler

import (
	"bytes"
	"encoding/json"
	"log/slog"
	"net/http"

	"github.com/lxmp7p/ya-dipl1/internal/service"

	"github.com/go-chi/chi/v5"
)

type UserHandler struct {
	logger      *slog.Logger
	userService *service.UserService
}

func (orders *UserHandler) UsersRoutes() chi.Router {
	r := chi.NewRouter()

	r.Get("/", orders.Balance)

	return r
}

func (orders *UserHandler) Balance(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value("user_id").(string)
	if !ok {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
	result, err := orders.userService.Balance(r.Context(), userID)
	if err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")

	buf := &bytes.Buffer{}
	if err := json.NewEncoder(buf).Encode(result); err != nil {
		orders.logger.Debug("Error encoding orders:", "err", err)
		http.Error(w, "internal error", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	_, _ = w.Write(buf.Bytes())
}
