package repository

import (
	"context"
	"log/slog"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/lxmp7p/ya-dipl1/internal/config"
)

var (
	migrationsPath = "file://migrations"
)

func InitDB(config config.Config, logger *slog.Logger) (*pgxpool.Pool, error) {
	if config.DatabaseDsn != "" {
		logger.Info("Running migrations...")
		if err := runMigrations(config.DatabaseDsn, logger); err != nil {
			return nil, err
		}
		logger.Info("Migrations done")
	}

	pool, err := pgxpool.New(context.Background(), config.DatabaseDsn)
	if err != nil {
		return nil, err
	}

	return pool, nil
}

func runMigrations(DSN string, logger *slog.Logger) error {
	m, err := migrate.New(migrationsPath, DSN)
	if err != nil {
		logger.Error("Failed to initialize migrate", "err", err)
		return err
	}
	if err := m.Up(); err != nil && err != migrate.ErrNoChange {
		logger.Error("Migration failed:", "err", err)
		return err
	}

	logger.Info("Migrations applied successfully!")
	return nil
}
