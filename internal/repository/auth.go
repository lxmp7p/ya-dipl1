package repository

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5"
)

func (rep *Repository) Registration(ctx context.Context, login, passwordHash string) error {
	query := `
		INSERT INTO auth (login, password_hash)
		VALUES ($1, $2)
	`

	_, err := rep.Db.Exec(ctx, query, login, passwordHash)
	return err
}

func (rep *Repository) ValidatePasswordHash(ctx context.Context, login, hash string) (string, error) {
	query := `
		SELECT password_hash FROM auth WHERE login = $1 AND password_hash = $2
	`
	var passwordHash string
	err := rep.Db.QueryRow(ctx, query, login, hash).Scan(&passwordHash)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return "", errors.New("UserNotFound")
		}
		return "", err
	}
	return passwordHash, nil
}

func (rep *Repository) CreateSession(ctx context.Context, userLogin, sessionID string) error {
	query := `
		INSERT INTO sessions (session_id, login)
		VALUES ($1, $2)
	`
	_, err := rep.Db.Exec(ctx, query, sessionID, userLogin)
	return err
}
