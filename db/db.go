package db

import (
	"context"
	"log/slog"

	"github.com/jmoiron/sqlx"
	"github.com/maciejas22/conference-manager/api/internal/config"
	_ "github.com/mattn/go-sqlite3"
)

type DB struct {
	Logger        *slog.Logger
	QueryExecutor *QueryExecutor
	Conn          *sqlx.DB
}

func Connect(ctx context.Context, logger *slog.Logger) (*DB, error) {
	connStr := config.AppConfig.DatabaseURL

	db, err := sqlx.Connect("sqlite3", connStr)
	if err != nil {
		logger.Error("failed to connect to db", "error", err)
		return nil, err
	}

	return &DB{
		Conn:   db,
		Logger: logger,
		QueryExecutor: &QueryExecutor{
			queryer: db,
			execer:  db,
			logger:  logger,
		},
	}, nil
}

func (db *DB) Close() (err error) {
	if sqlErr := db.Conn.Close(); sqlErr != nil {
		db.Logger.Error("failed to close db connection", "error", sqlErr)
		return sqlErr
	}
	db.Logger.Info("closed connection to db")

	return nil
}

func Transaction(ctx context.Context, ql *QueryExecutor, fn func(*QueryExecutor) error) error {
	tx, err := ql.BeginTxx(ctx, nil)
	if err != nil {
		ql.logger.Error("failed to begin transaction", "error", err)
		return err
	}

	txLogger := ql.WithTx(tx)

	err = fn(txLogger)

	if err != nil {
		if rbErr := tx.Rollback(); rbErr != nil {
			ql.logger.Error("failed to rollback transaction", "error", rbErr)
			return rbErr
		}
		ql.logger.Error("transaction failed", "error", err)
		return err
	}

	if err := tx.Commit(); err != nil {
		ql.logger.Error("failed to commit transaction", "error", err)
		return err
	}

	return nil
}
