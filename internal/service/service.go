package service

import (
	"context"
	"log/slog"

	"github.com/lxmp7p/ya-dipl1/internal/repository"
)

type Services struct {
	logger  *slog.Logger
	Auth    AuthServiceInterface
	Order   OrderServiceInterface
	Balance BalanceServiceInterface
}

func NewServices(logger *slog.Logger, repo *repository.Repository) *Services {
	return &Services{
		logger:  logger,
		Auth:    NewAuthService(repo),
		Order:   NewOrderService(logger, repo, repo),
		Balance: NewBalanceService(logger, repo, repo),
	}
}

type AuthServiceInterface interface {
	Registration(ctx context.Context, login, password string) (string, error)
	Login(ctx context.Context, login, password string) (string, error)
	ValidateToken(ctx context.Context, login, hash string) (bool, error)
}

type OrderServiceInterface interface {
	UploadOrder(ctx context.Context, orderNumber string, userID string) error
	List(ctx context.Context, userID string) ([]OrderResponse, error)
}

type BalanceServiceInterface interface {
	Balance(ctx context.Context, userID string) (BalanceResponse, error)
	Withdrawn(ctx context.Context, orderNumber string, sum float64, userID string) error
	ListWithdrawn(ctx context.Context, userID string) ([]repository.Withdrawal, error)
}
