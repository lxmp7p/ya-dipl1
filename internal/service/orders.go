package service

import (
	"context"
	"errors"
	"strings"

	"github.com/lxmp7p/ya-dipl1/internal/repository"
)

type OrderService struct {
	repo repository.Auth
}

func NewOrderService(repo repository.Auth) *AuthService {
	return &AuthService{
		repo: repo,
	}
}

func (h *OrderService) UploadOrder(ctx context.Context, orderString string) error {
	orderNumber := strings.TrimSpace(orderString)

	if !isDigitsOnly(orderNumber) {
		return errors.New("bad request")
	}

	if !isValidLuhn(orderNumber) {
		return errors.New("bad request")
	}
	return nil
}
