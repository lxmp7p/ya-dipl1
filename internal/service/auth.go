package service

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"github.com/lxmp7p/ya-dipl1/internal/repository"
)

type AuthService struct {
	repo repository.Auth
}

func NewAuthService(repo repository.Auth) AuthService {
	return AuthService{
		repo: repo,
	}
}

var (
	ErrUserExists = errors.New("user already exists")
)

func (auth *AuthService) Registration(ctx context.Context, login, password string) (string, error) {
	hash, err := HashPassword(password)
	if err != nil {
		return "", err
	}

	err = auth.repo.Registration(ctx, login, hash)
	if err != nil {
		if isUniqueViolation(err) {
			return "", ErrUserExists
		}
		return "", err
	}

	sessionID := uuid.NewString()
	err = auth.repo.CreateSession(ctx, login, sessionID)
	if err != nil {
		return "", err
	}

	return sessionID, nil
}
