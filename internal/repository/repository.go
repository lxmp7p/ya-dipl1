package repository

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
)

type Repository struct {
	db *pgxpool.Pool
}

type Auth interface {
	Registration(ctx context.Context, login, password string) error
	CreateSession(ctx context.Context, userLogin, sessionID string) error
}
