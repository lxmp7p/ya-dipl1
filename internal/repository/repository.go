package repository

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
)

type Repository struct {
	Db *pgxpool.Pool
}

func NewRepository(database *pgxpool.Pool) Repository {
	return Repository{
		Db: database,
	}
}

type Auth interface {
	Registration(ctx context.Context, login, password string) error
	GetAuthDataByLogin(ctx context.Context, login string) (AuthData, error)
	CreateSession(ctx context.Context, userLogin, sessionID string) error
}
