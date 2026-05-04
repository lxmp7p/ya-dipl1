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
	Registration(ctx context.Context, login, passwordHash string) (UserInfo, error)
	GetAuthDataByLogin(ctx context.Context, login string) (AuthData, error)
	CreateSession(ctx context.Context, userLogin, userID, sessionID string) error
	CheckAuthDataBySession(ctx context.Context, sessionID string) (SessionData, error)
}

type Order interface {
	Create(ctx context.Context, orderNumber string, userId string) error
	FindOrderWithUserByNumber(ctx context.Context, orderNumber string) (OrderData, bool, error)
	List(ctx context.Context, userId string) ([]OrderData, error)
	ListAllOrders(ctx context.Context) ([]OrderData, error)
}

type User interface {
	Balance(ctx context.Context, userId string) (BalanceInfo, error)
	AddBalance(ctx context.Context, userId string, amount float64) (BalanceInfo, error)
}
