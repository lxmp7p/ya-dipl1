package handler

import (
	"log/slog"

	"github.com/lxmp7p/ya-dipl1/internal/config"
	"github.com/lxmp7p/ya-dipl1/internal/service"
	"github.com/lxmp7p/ya-dipl1/internal/utils.go"

	"github.com/go-chi/chi/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

var (
	DefaultApiRoute = "/api"
)

type Handler struct {
	Services *service.Services
	Config   config.Config
	Logger   *slog.Logger
	Database *pgxpool.Pool
}

func NewHandler(logger *slog.Logger, database *pgxpool.Pool) *Handler {
	return &Handler{
		Logger:   logger,
		Database: database,
	}
}

func (handler *Handler) InitRoutes() chi.Router {
	handler.validateHandler()

	services := service.NewServices(service.AuthService{})
	handler.Services = services

	apiRouter := chi.NewRouter()

	apiRouter.Mount("/", AuthRoutes(&services.Auth))
	return apiRouter
}

func (handler *Handler) validateHandler() {
	if handler.Logger == nil {
		handler.Logger = utils.CreateLogger()
	}
}
