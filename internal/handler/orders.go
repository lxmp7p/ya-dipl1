package handler

import (
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

	r.Post(DefaultApiRoute+"/user/orders", orders.UploadOrder)

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

	orders.orderService.UploadOrder(r.Context(), orderNumber)

	w.WriteHeader(http.StatusAccepted)
}
