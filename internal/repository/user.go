package repository

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5"
)

type BalanceInfo struct {
	Balance   float64
	Withdrawn float64
}

func (rep *Repository) Balance(ctx context.Context, userId string) (BalanceInfo, error) {
	query := `
        SELECT balance, withdrawn FROM users WHERE user_id = $1
    `

	var balance BalanceInfo
	err := rep.Db.QueryRow(ctx, query, userId).Scan(&balance.Balance, &balance.Withdrawn)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return BalanceInfo{Balance: 0, Withdrawn: 0}, nil
		}
		return BalanceInfo{}, err
	}

	return balance, nil
}
