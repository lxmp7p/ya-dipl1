package service

import (
	"context"
	"log/slog"
)

type Services struct {
	logger *slog.Logger
	auth   AuthServiceInterface
	order  OrderServiceInterface
	user   UserServiceInterface
}

type AuthServiceInterface interface {
	Registration(ctx context.Context, login, password string) (string, error)
	Login(ctx context.Context, login, password string) (string, error)
	ValidateToken(ctx context.Context, login, hash string) (bool, error)
}

type OrderServiceInterface interface {
	UploadOrder(ctx context.Context, orderString string, userID int) error
}

type UserServiceInterface interface {
	Balance(ctx context.Context, userID string) (BalanceResponse, error)
}
