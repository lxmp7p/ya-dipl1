package handler

import (
	"log/slog"

	"github.com/lxmp7p/ya-dipl1/internal/config"
	"github.com/lxmp7p/ya-dipl1/internal/repository"
	"github.com/lxmp7p/ya-dipl1/internal/service"
	"github.com/lxmp7p/ya-dipl1/internal/utils.go"

	"github.com/go-chi/chi/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

var (
	DefaultApiRoute = "/api"
)

type Handler struct {
	Services service.Service
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
	repo := repository.NewRepository(handler.Database)

	authService := service.NewAuthService(&repo)
	orderService := service.NewOrderService(handler.Logger, &repo, &repo)
	userService := service.NewUserService(handler.Logger, &repo)

	apiRouter := chi.NewRouter()
	authHandler := AuthHandler{
		logger:      handler.Logger,
		authService: authService,
	}

	orderHandler := OrderHandler{
		logger:       handler.Logger,
		orderService: orderService,
	}

	userHandler := UserHandler{
		logger: handler.Logger, userService: userService,
	}

	apiRouter.Mount(DefaultApiRoute+"/user", authHandler.AuthRoutes())
	apiRouter.Group(func(r chi.Router) {
		r.Use(AuthMiddleware(&repo))
		r.Mount(DefaultApiRoute+"/user/orders", orderHandler.OrdersRoutes())
		r.Mount(DefaultApiRoute+"/user/balance", userHandler.UsersRoutes())
	})
	return apiRouter
}

func (handler *Handler) validateHandler() {
	if handler.Logger == nil {
		handler.Logger = utils.CreateLogger()
	}
}
