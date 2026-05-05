package handler

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"log/slog"
	"net/http"
	"strings"

	"github.com/lxmp7p/ya-dipl1/internal/service"

	"github.com/go-chi/chi/v5"
)

type OrderHandler struct {
	logger       *slog.Logger
	orderService *service.OrderService
}

func (orders *OrderHandler) OrdersRoutes() chi.Router {
	r := chi.NewRouter()

	r.Post("/", orders.UploadOrder)
	r.Get("/", orders.List)

	return r
}

func (orders *OrderHandler) UploadOrder(w http.ResponseWriter, r *http.Request) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Failed to read request body", http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	if r.Header.Get("Content-Type") != "text/plain" {
		http.Error(w, "Invalid Content-Type, expected text/plain", http.StatusBadRequest)
		return
	}

	orderNumber := strings.TrimSpace(string(body))
	userID, ok := r.Context().Value("user_id").(string)
	if !ok {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	if err := orders.orderService.UploadOrder(r.Context(), orderNumber, userID); err != nil {
		if errors.Is(err, service.ErrOrderInvalid) {
			http.Error(w, http.StatusText(http.StatusUnprocessableEntity), http.StatusUnprocessableEntity)
			return
		}
		if errors.Is(err, service.ErrOrderExists) {
			http.Error(w, http.StatusText(http.StatusConflict), http.StatusConflict)
			return
		}
		if errors.Is(err, service.ErrUserOrderExists) {
			w.WriteHeader(http.StatusOK)
			return
		}
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusAccepted)
}

func (orders *OrderHandler) List(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value("user_id").(string)
	if !ok {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
	result, err := orders.orderService.List(r.Context(), userID)
	if err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	if len(result) == 0 {
		w.WriteHeader(http.StatusNoContent)
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
