package service

import "log/slog"

type Services struct {
	Logger *slog.Logger
	Auth   AuthService
}

func NewServices(logger *slog.Logger, authService AuthService) *Services {
	return &Services{
		Logger: logger,
		Auth:   authService,
	}
}
