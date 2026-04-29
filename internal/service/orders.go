package service

import (
	"context"
	"errors"

	"github.com/lxmp7p/ya-dipl1/internal/repository"
)

type OrderService struct {
	repo repository.Order
}

func NewOrderService(repo repository.Order) *OrderService {
	return &OrderService{
		repo: repo,
	}
}

func (ors *OrderService) UploadOrder(ctx context.Context, orderNumber string, userID string) error {
	if !isDigitsOnly(orderNumber) {
		return errors.New("bad request")
	}

	if !isValidLuhn(orderNumber) {
		return errors.New("bad request")
	}

	order, exist, err := ors.repo.FindOrderWithUserByNumber(ctx, orderNumber)
	if err != nil {
		return errors.New("bad request")
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
	// проверка заказа в базе
	return nil
}
