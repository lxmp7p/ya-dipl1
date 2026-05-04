package service

import (
	"context"
	"log/slog"

	"github.com/lxmp7p/ya-dipl1/internal/repository"
)

type UserService struct {
	logger *slog.Logger
	repo   repository.User
}

type BalanceResponse struct {
	Accrual   *float64 `json:"current,omitempty"`
	Withdrawn *float64 `json:"withdrawn,omitempty"`
}

func NewUserService(logger *slog.Logger, repo repository.User) *UserService {
	return &UserService{
		logger: logger,
		repo:   repo,
	}
}

func (ors *UserService) Balance(ctx context.Context, userID string) (BalanceResponse, error) {
	balance, err := ors.repo.Balance(ctx, userID)
	if err != nil {
		ors.logger.Error(err.Error())
		return BalanceResponse{}, nil
	}

	return BalanceResponse{
		Accrual:   &balance.Balance,
		Withdrawn: &balance.Withdrawn,
	}, nil
}
