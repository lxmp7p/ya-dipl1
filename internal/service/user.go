package service

import (
	"context"
	"errors"
	"log/slog"
	"net/http"

	"github.com/lxmp7p/ya-dipl1/internal/repository"
)

type UserService struct {
	logger    *slog.Logger
	repo      repository.User
	orderRepo repository.Order
}

type BalanceResponse struct {
	Current   float64 `json:"current,omitempty"`
	Withdrawn float64 `json:"withdrawn,omitempty"`
}

func NewUserService(
	logger *slog.Logger,
	repo repository.User,
	orderRepo repository.Order,
) *UserService {
	return &UserService{
		logger:    logger,
		repo:      repo,
		orderRepo: orderRepo,
	}
}

func (ors *UserService) Balance(ctx context.Context, userID string) (BalanceResponse, error) {
	balance, err := ors.repo.Balance(ctx, userID)
	if err != nil {
		ors.logger.Error(err.Error())
		return BalanceResponse{}, err
	}

	return BalanceResponse{
		Current:   balance.Balance,
		Withdrawn: balance.Withdrawn,
	}, nil
}

func (ors *UserService) Withdrawn(ctx context.Context, orderNumber string, money float64) error {
	order, exist, err := ors.orderRepo.FindOrderWithUserByNumber(ctx, orderNumber)

	if err != nil {
		return errors.New(http.StatusText(http.StatusBadRequest))
	}

	if !exist {
		return ErrOrderExists
	}

	err = ors.repo.Withdrawn(ctx, money, order.UserID)
	if err != nil {
		ors.logger.Error(err.Error())
		return err
	}

	return nil
}
