package repository

import (
	"context"
	"errors"
	"time"

	"github.com/jackc/pgx/v5"
)

type UserData struct {
	OrderNumber string
	UserID      string
	Login       string
	Status      string
	Accrual     *float64
	UploadedAt  time.Time
}

func (rep *Repository) Balance(ctx context.Context, userId string) (UserData, error) {
	query := `
        SELECT * FROM auth WHERE id = $1
    `

	rows, err := rep.Db.Query(ctx, query, userId)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return UserData{}, ErrOrderNotFound
		}
		return UserData{}, err
	}
	defer rows.Close()

	var orders []OrderData

	for rows.Next() {
		var order OrderData
		err := rows.Scan(
			&order.OrderNumber,
			&order.UserID,
			&order.Login,
			&order.Status,
			&order.Accrual,
			&order.UploadedAt,
		)
		if err != nil {
			return UserData{}, err
		}
		orders = append(orders, order)
	}

	return UserData{}, nil
}
