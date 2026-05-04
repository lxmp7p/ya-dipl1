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

func NewUserService(logger *slog.Logger, repo repository.User) *UserService {
	return &UserService{
		logger: logger,
		repo:   repo,
	}
}

func (ors *UserService) Balance(ctx context.Context, userID string) ([]OrderResponse, error) {
	_, err := ors.repo.Balance(ctx, userID)
	if err != nil {
		ors.logger.Error(err.Error())
		return []OrderResponse{}, nil
	}
	var ordersResult []OrderResponse

	return ordersResult, nil
}
