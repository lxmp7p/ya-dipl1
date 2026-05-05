package handler

import (
	"log/slog"

	"github.com/lxmp7p/ya-dipl1/internal/config"
	"github.com/lxmp7p/ya-dipl1/internal/repository"
	"github.com/lxmp7p/ya-dipl1/internal/service"
	"github.com/lxmp7p/ya-dipl1/internal/utils.go"

	"github.com/go-chi/chi/v5"
)

var (
	DefaultApiRoute = "/api"
)

type Handler struct {
	Config     config.Config
	Logger     *slog.Logger
	Repository repository.Repository
}

func NewHandler(logger *slog.Logger, repository repository.Repository, config config.Config) *Handler {
	return &Handler{
		Logger:     logger,
		Repository: repository,
		Config:     config,
	}
}

func (handler *Handler) InitRoutes() chi.Router {
	handler.validateHandler()

	services := service.NewServices(
		handler.Logger,
		&handler.Repository,
	)

	apiRouter := chi.NewRouter()
	authHandler := AuthHandler{
		logger:      handler.Logger,
		authService: services.Auth,
	}

	orderHandler := OrderHandler{
		logger:       handler.Logger,
		orderService: services.Order,
	}

	userHandler := BalanceHandler{
		logger:         handler.Logger,
		balanceService: services.Balance,
	}

	apiRouter.Mount(DefaultApiRoute+"/user", authHandler.AuthRoutes())
	apiRouter.Group(func(r chi.Router) {
		r.Use(AuthMiddleware(&handler.Repository))
		r.Mount(DefaultApiRoute+"/user/orders", orderHandler.OrdersRoutes())
		r.Mount(DefaultApiRoute+"/user/balance", userHandler.BalanceRoutes())
		r.Get(DefaultApiRoute+"/user/withdrawals", userHandler.ListWithdrawn)
	})
	return apiRouter
}

func (handler *Handler) validateHandler() {
	if handler.Logger == nil {
		handler.Logger = utils.CreateLogger()
	}
}
