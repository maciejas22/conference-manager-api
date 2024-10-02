package db

import (
	"context"
	"errors"
	"log"
	"log/slog"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/stdlib"
	"github.com/jackc/pgx/v5/tracelog"
	"github.com/jmoiron/sqlx"
	"github.com/maciejas22/conference-manager/api/internal/config"
	_ "github.com/mattn/go-sqlite3"
)

type DB struct {
	Logger *slog.Logger
	Conn   *sqlx.DB
}

func GetDriverConfig(l *slog.Logger) *pgx.ConnConfig {
	connConfig, _ := pgx.ParseConfig(config.AppConfig.DatabaseURL)
	l.Info("Connecting to database", "url", config.AppConfig.DatabaseURL)
	adapterLogger := NewLogger(l)
	m := MultiQueryTracer{
		Tracers: []pgx.QueryTracer{
			&tracelog.TraceLog{
				Logger:   adapterLogger,
				LogLevel: tracelog.LogLevelTrace,
			},
		},
	}
	connConfig.Tracer = &m
	return connConfig
}

func Connect(ctx context.Context, logger *slog.Logger) (*DB, error) {
	dc := GetDriverConfig(logger)
	connStr := stdlib.RegisterConnConfig(dc)

	conn, err := sqlx.Connect("pgx", connStr)
	if err != nil {
		logger.Error("failed to connect to db", "error", err)
		return nil, errors.New("failed to connect to db")
	}

	return &DB{
		Conn:   conn,
		Logger: logger,
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

func Transaction(ctx context.Context, db *sqlx.DB, fn func(*sqlx.Tx) error) error {
	tx, err := db.BeginTxx(ctx, nil)
	if err != nil {
		log.Println("Error in transaction ", err)
		return err
	}

	err = fn(tx)

	if err != nil {
		log.Println("Error in transaction2 ", err)
		if rbErr := tx.Rollback(); rbErr != nil {
			log.Println("Error in transaction rollback ", rbErr)
			return rbErr
		}
		return err
	}

	if err := tx.Commit(); err != nil {
		log.Println("Error in transaction commit ", err)
		return err
	}

	return nil
}
