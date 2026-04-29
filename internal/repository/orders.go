package repository

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
)

type OrderData struct {
	OrderNumber string
	UserID      string
	Login       string // username из таблицы users
	Status      string
	Accrual     int
}

func (rep *Repository) Create(ctx context.Context, orderNumber, userId string) error {
	query := `
		INSERT INTO orders (order_number, user_id)
		VALUES ($1, $2)
	`
	_, err := rep.Db.Exec(ctx, query, orderNumber, userId)

	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == "23505" {
			return ErrOrderExists
		}
		return err
	}
	return err
}

func (rep *Repository) FindOrderWithUserByNumber(ctx context.Context, orderNumber string) (OrderData, bool, error) {
	query := `
        SELECT 
            o.order_number, 
            o.user_id, 
            u.login,
            o.status, 
            o.accrual
        FROM orders o
        JOIN auth u ON o.user_id = u.id
        WHERE o.order_number = $1
    `

	var result OrderData
	err := rep.Db.QueryRow(ctx, query, orderNumber).Scan(
		&result.OrderNumber,
		&result.UserID,
		&result.Login,
		&result.Status,
		&result.Accrual,
	)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return OrderData{}, false, nil
		}
		return OrderData{}, false, err
	}

	return result, true, nil
}
