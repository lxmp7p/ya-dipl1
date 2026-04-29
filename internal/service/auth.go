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

func NewAuthService(repo repository.Auth) *AuthService {
	return &AuthService{
		repo: repo,
	}
}

func (auth *AuthService) Registration(ctx context.Context, login, password string) (string, error) {
	hash, err := HashPassword(password)
	if err != nil {
		return "", err
	}

	user, err := auth.repo.Registration(ctx, login, hash)
	if err != nil {
		if errors.Is(err, repository.ErrUserExists) {
			return "", ErrUserExists
		}
		return "", err
	}

	sessionID := uuid.NewString()
	err = auth.repo.CreateSession(ctx, user.Login, user.ID, user.ID)
	if err != nil {
		return "", err
	}

	return sessionID, nil
}

func (auth *AuthService) Login(ctx context.Context, login, password string) (string, error) {
	authData, err := auth.repo.GetAuthDataByLogin(ctx, login)
	if err != nil {
		if errors.Is(err, repository.ErrUserNotFound) {
			return "", ErrInvalidCredentials
		}
		return "", err
	}

	if ok := CheckPassword(password, authData.Password_hash); !ok {
		return "", ErrInvalidCredentials
	}

	sessionID := uuid.NewString()
	err = auth.repo.CreateSession(ctx, login, authData.UserId, sessionID)
	if err != nil {
		return "", err
	}

	return sessionID, nil
}

func (auth *AuthService) ValidateToken(ctx context.Context, login, hash string) (bool, error) {
	_, err := auth.repo.GetAuthDataByLogin(ctx, login)
	if err != nil {
		return false, err
	}
	return true, nil
}
