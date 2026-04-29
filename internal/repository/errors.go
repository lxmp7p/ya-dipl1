package repository

import "errors"

var (
	ErrUserNotFound    = errors.New("user not found")
	ErrUserExists      = errors.New("user already exists")
	ErrSessionNotFound = errors.New("session not found")
	ErrInvalidSession  = errors.New("invalid session")
	ErrOrderNotFound   = errors.New("order not found")
	ErrOrderExists     = errors.New("order already exists")
)
