package service

import (
	"context"
	"errors"
	"log/slog"

	"github.com/lxmp7p/ya-dipl1/internal/repository"
)

type BalanceService struct {
	logger    *slog.Logger
	repo      repository.User
	orderRepo repository.Order
}

type BalanceResponse struct {
	Current   float64 `json:"current,omitempty"`
	Withdrawn float64 `json:"withdrawn,omitempty"`
}

type WithdrawnListResponse struct {
	Order        string  `json:"order"`
	Sum          float64 `json:"sum"`
	Processed_at string  `json:"processed_at"`
}

func NewBalanceService(
	logger *slog.Logger,
	repo repository.User,
	orderRepo repository.Order,
) *BalanceService {
	return &BalanceService{
		logger:    logger,
		repo:      repo,
		orderRepo: orderRepo,
	}
}

func (ors *BalanceService) Balance(ctx context.Context, userID string) (BalanceResponse, error) {
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

func (ors *BalanceService) Withdrawn(ctx context.Context, orderNumber string, money float64, userID string) error {
	if !isValidLuhn(orderNumber) {
		return ErrOrderInvalid
	}

	err := ors.repo.Withdrawn(ctx, money, userID, orderNumber)
	if err != nil {
		ors.logger.Error(err.Error())
		if errors.Is(err, repository.ErrNoEnoughMoney) {
			return ErrNoEnoughMoney
		}
		return err
	}

	return nil
}

func (ors *BalanceService) ListWithdrawn(ctx context.Context, userID string) ([]repository.Withdrawal, error) {
	withdrawals, err := ors.repo.ListWithdrawn(ctx, userID)
	if err != nil {
		ors.logger.Error(err.Error())
		return []repository.Withdrawal{}, err
	}

	return withdrawals, nil
}
