package repository

import (
	"context"
	"errors"
	"time"

	"github.com/jackc/pgx/v5"
)

type BalanceInfo struct {
	Balance   float64
	Withdrawn float64
}

type Withdrawal struct {
	Order       string    `json:"order"`
	Sum         float64   `json:"sum"`
	ProcessedAt time.Time `json:"processed_at"`
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

func (rep *Repository) ListWithdrawn(ctx context.Context, userId string) ([]Withdrawal, error) {
	query := `
        SELECT order_number, sum, processed_at 
        FROM withdrawn 
        WHERE user_id = $1
        ORDER BY processed_at DESC
    `

	rows, err := rep.Db.Query(ctx, query, userId)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return []Withdrawal{}, nil
		}
		return nil, err
	}
	defer rows.Close()

	var withdrawals []Withdrawal
	for rows.Next() {
		var w Withdrawal
		err := rows.Scan(
			&w.Order,
			&w.Sum,
			&w.ProcessedAt,
		)
		if err != nil {
			return nil, err
		}
		withdrawals = append(withdrawals, w)
	}

	return withdrawals, nil
}

func (rep *Repository) AddBalance(ctx context.Context, userId string, amount float64) (BalanceInfo, error) {
	query := `
        INSERT INTO users (user_id, balance, withdrawn) 
        VALUES ($1, $2, $3)
        ON CONFLICT (user_id) DO UPDATE 
        SET balance = users.balance + $2
        RETURNING balance, withdrawn
    `

	var balance BalanceInfo
	err := rep.Db.QueryRow(ctx, query, userId, amount, 0).Scan(&balance.Balance, &balance.Withdrawn)

	if err != nil {
		return BalanceInfo{}, err
	}

	return balance, nil
}

func (rep *Repository) Withdrawn(ctx context.Context, money float64, userID string, orderNumber string) error {
	query := `
		UPDATE users 
		SET balance = balance - $1 withdrawn = withdrawn + $1
		WHERE user_id = $2
	`
	_, err := rep.Db.Exec(ctx, query, money, userID)

	if err != nil {
		return err
	}

	query = `
        INSERT INTO withdrawn (order_number, sum, user_id, processed_at)
        VALUES ($1, $2, $3, NOW())
    `
	_, err = rep.Db.Exec(ctx, query, orderNumber, money, userID)

	return err
}
