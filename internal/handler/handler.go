package handler

import (
	"context"
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
	Config   config.Config
	Logger   *slog.Logger
	Database *pgxpool.Pool
}

func NewHandler(logger *slog.Logger, database *pgxpool.Pool, config config.Config) *Handler {
	return &Handler{
		Logger:   logger,
		Database: database,
		Config:   config,
	}
}

func (handler *Handler) InitRoutes() chi.Router {
	handler.validateHandler()
	repo := repository.NewRepository(handler.Database)

	worker := utils.NewAccrualWorker(handler.Config.BalanceSystemAddr, &repo, handler.Logger)
	worker.Start(context.Background())

	services := service.NewServices(
		handler.Logger,
		&repo,
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
		r.Use(AuthMiddleware(&repo))
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
