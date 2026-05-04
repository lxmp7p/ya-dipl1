package service

import (
	"context"
	"errors"
	"log/slog"
	"net/http"
	"time"

	"github.com/lxmp7p/ya-dipl1/internal/repository"
)

type OrderService struct {
	logger   *slog.Logger
	repo     repository.Order
	userRepo repository.User
}

func NewOrderService(
	logger *slog.Logger,
	repo repository.Order,
	userRepo repository.User,
) *OrderService {
	return &OrderService{
		logger:   logger,
		repo:     repo,
		userRepo: userRepo,
	}
}

type OrderResponse struct {
	Number     string   `json:"number"`
	Status     string   `json:"status"`
	Accrual    *float64 `json:"accrual,omitempty"`
	UploadedAt string   `json:"uploaded_at"`
}

func (ors *OrderService) UploadOrder(ctx context.Context, orderNumber string, userID string) error {
	if !isDigitsOnly(orderNumber) {
		return ErrOrderInvalid
	}

	if !isValidLuhn(orderNumber) {
		return ErrOrderInvalid
	}

	order, exist, err := ors.repo.FindOrderWithUserByNumber(ctx, orderNumber)
	if err != nil {
		return errors.New(http.StatusText(http.StatusBadRequest))
	}

	if exist {
		if order.UserID != userID {
			return ErrOrderExists
		} else {
			return nil
		}
	}

	if err = ors.repo.Create(ctx, orderNumber, userID); err != nil {
		return err
	}

	return nil
}

func (ors *OrderService) List(ctx context.Context, userID string) ([]OrderResponse, error) {
	orders, err := ors.repo.List(ctx, userID)
	if err != nil {
		ors.logger.Error(err.Error())
		return []OrderResponse{}, err
	}
	var ordersResult []OrderResponse
	for _, order := range orders {
		ordersResult = append(ordersResult, OrderResponse{
			Number:     order.OrderNumber,
			Status:     order.Status,
			Accrual:    order.Accrual,
			UploadedAt: order.UploadedAt.Format(time.RFC3339),
		})
	}

	return ordersResult, nil
}
