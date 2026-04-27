package handler

import (
	"github.com/lxmp7p/ya-dipl1/internal/service"

	"github.com/go-chi/chi/v5"
)

func AuthRoutes(authService *service.AuthService) chi.Router {
	r := chi.NewRouter()
	r.Post(DEFAULT_API_ROUTE+"/user/register", authService.Registration)

	return r
}
