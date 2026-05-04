package repository

import (
	"context"
	"errors"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
)

type OrderData struct {
	OrderNumber string
	UserID      string
	Login       string
	Status      string
	Accrual     *float64
	UploadedAt  time.Time
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
            o.accrual,
			uploaded_at
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
		&result.UploadedAt,
	)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return OrderData{}, false, nil
		}
		return OrderData{}, false, err
	}

	return result, true, nil
}

func (rep *Repository) List(ctx context.Context, userId string) ([]OrderData, error) {
	query := `
        SELECT 
            o.order_number, 
            o.user_id, 
            u.login,
            o.status, 
            o.accrual,
			o.uploaded_at
        FROM orders o
        JOIN auth u ON o.user_id = u.id
        WHERE o.user_id = $1
		ORDER BY o.uploaded_at DESC
    `

	rows, err := rep.Db.Query(ctx, query, userId)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return []OrderData{}, ErrOrderNotFound
		}
		return []OrderData{}, err
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
			return []OrderData{}, err
		}
		orders = append(orders, order)
	}

	return orders, nil
}

func (rep *Repository) ListAllOrders(ctx context.Context) ([]OrderData, error) {
	query := `
        SELECT 
            o.order_number, 
            o.user_id, 
            u.login,
            o.status, 
            o.accrual,
			o.uploaded_at
        FROM orders o
        JOIN auth u ON o.user_id = u.id
		ORDER BY o.uploaded_at DESC
    `

	rows, err := rep.Db.Query(ctx, query)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return []OrderData{}, ErrOrderNotFound
		}
		return []OrderData{}, err
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
			return []OrderData{}, err
		}
		orders = append(orders, order)
	}

	return orders, nil
}

func (rep *Repository) Update(ctx context.Context, orderNumber, status string, accrual float64) error {
	query := `
		UPDATE orders 
		SET status = $1, accrual = $2 
		WHERE order_number = $3
	`
	_, err := rep.Db.Exec(ctx, query, status, accrual, orderNumber)

	if err != nil {
		return err
	}
	return err
}
