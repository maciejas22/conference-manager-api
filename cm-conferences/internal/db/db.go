package db

import (
	"context"
	"errors"
	"log/slog"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/stdlib"
	"github.com/jmoiron/sqlx"
	"github.com/maciejas22/conference-manager-api/cm-conferences/config"
)

type DB struct {
	Logger *slog.Logger
	Conn   *sqlx.DB
}

func GetDriverConfig(l *slog.Logger) *pgx.ConnConfig {
	connConfig, _ := pgx.ParseConfig(config.AppConfig.DatabaseURL)
	l.Info("DB", "url", config.AppConfig.DatabaseURL)
	connConfig.Tracer = &myQueryTracer{l}
	return connConfig
}

func Connect(ctx context.Context, logger *slog.Logger) (*DB, error) {
	dc := GetDriverConfig(logger)
	connStr := stdlib.RegisterConnConfig(dc)

	conn, err := sqlx.Connect("pgx", connStr)
	if err != nil {
		logger.Error("DB", "error", err)
		return nil, errors.New("Failed to connect to db")
	}

	return &DB{
		Conn:   conn,
		Logger: logger,
	}, nil
}

func (db *DB) Close() (err error) {
	if sqlErr := db.Conn.Close(); sqlErr != nil {
		db.Logger.Error("DB", "error", sqlErr)
		return sqlErr
	}
	db.Logger.Info("DB", "message", "closed connection to db")

	return nil
}
