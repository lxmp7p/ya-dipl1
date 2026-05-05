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
	UserId        string
}

type SessionData struct {
	SessionId string
	Login     string
	UserId    string
}

type UserInfo struct {
	ID           string `json:"id"`
	Login        string `json:"login"`
	PasswordHash string `json:"-"`
}

func (rep *Repository) Registration(ctx context.Context, login, passwordHash string) (UserInfo, error) {
	query := `
		INSERT INTO auth (login, password_hash)
		VALUES ($1, $2) RETURNING id, login, password_hash
	`
	var user UserInfo
	err := rep.Db.QueryRow(ctx, query, login, passwordHash).Scan(
		&user.ID,
		&user.Login,
		&user.PasswordHash,
	)

	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == "23505" {
			return UserInfo{}, ErrUserExists
		}
	}

	insertQuery := `
        INSERT INTO balance (user_id, balance, withdrawn)
        VALUES ($1, $2, $3)
    `

	_, err = rep.Db.Exec(ctx, insertQuery, user.ID, 0, 0)
	if err != nil {
		return UserInfo{}, err
	}

	return user, err
}

func (rep *Repository) GetAuthDataByLogin(ctx context.Context, login string) (AuthData, error) {
	query := `
		SELECT login, password_hash, id FROM auth WHERE login = $1
	`
	var authData AuthData
	err := rep.Db.QueryRow(ctx, query, login).Scan(&authData.Login, &authData.Password_hash, &authData.UserId)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return AuthData{}, ErrUserNotFound
		}
		return AuthData{}, err
	}
	return authData, nil
}

func (rep *Repository) CheckAuthDataBySession(ctx context.Context, sessionID string) (SessionData, error) {
	query := `
		SELECT session_id, login, user_id FROM sessions WHERE session_id = $1
	`

	var session SessionData
	err := rep.Db.QueryRow(ctx, query, sessionID).Scan(&session.SessionId, &session.Login, &session.UserId)
	if err != nil {
		return SessionData{}, err
	}
	return session, nil
}

func (rep *Repository) CreateSession(ctx context.Context, userLogin string, userID string, sessionID string) error {
	query := `
		INSERT INTO sessions (session_id, login, user_id)
		VALUES ($1, $2, $3)
	`
	_, err := rep.Db.Exec(ctx, query, sessionID, userLogin, userID)
	return err
}
