package repository

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
)

type AuthData struct {
	Login         string
	Password_hash string
}

func (rep *Repository) Registration(ctx context.Context, login, passwordHash string) error {
	query := `
		INSERT INTO auth (login, password_hash)
		VALUES ($1, $2)
	`

	_, err := rep.Db.Exec(ctx, query, login, passwordHash)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == "23505" {
			return ErrUserExists
		}
	}
	return err
}

func (rep *Repository) GetAuthDataByLogin(ctx context.Context, login string) (AuthData, error) {
	query := `
		SELECT login, password_hash FROM auth WHERE login = $1
	`
	var authData AuthData
	err := rep.Db.QueryRow(ctx, query, login).Scan(&authData.Login, &authData.Password_hash)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return AuthData{}, ErrUserNotFound
		}
		return AuthData{}, err
	}
	return authData, nil
}

func (rep *Repository) CreateSession(ctx context.Context, userLogin, sessionID string) error {
	query := `
		INSERT INTO sessions (session_id, login)
		VALUES ($1, $2)
	`
	_, err := rep.Db.Exec(ctx, query, sessionID, userLogin)
	return err
}
