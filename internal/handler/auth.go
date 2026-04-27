package handler

import (
	"github.com/lxmp7p/ya-dipl1/internal/service"

	"github.com/go-chi/chi/v5"
)

func AuthRoutes(authService *service.AuthService) chi.Router {
	r := chi.NewRouter()

	r.Use(AuthMiddleware())

	r.Post(DefaultApiRoute+"/user/register", authService.Registration)
	
	return r
}
