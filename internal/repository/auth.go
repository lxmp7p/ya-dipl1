package repository

import "context"

func (rep *Repository) Registration(ctx context.Context, login, passwordHash string) error {
	query := `
		INSERT INTO auth (login, password_hash)
		VALUES ($1, $2)
	`

	_, err := rep.db.Exec(ctx, query, login, passwordHash)
	return err
}

func (rep *Repository) CreateSession(ctx context.Context, userLogin, sessionID string) error {
	query := `
		INSERT INTO sessions (session_id, login)
		VALUES ($1, $2)
	`
	_, err := rep.db.Exec(ctx, query, sessionID, userLogin)
	return err
}
