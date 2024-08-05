package db

import (
	"context"
	"errors"
	"log/slog"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/stdlib"
	"github.com/jackc/pgx/v5/tracelog"
	"github.com/jmoiron/sqlx"
	"github.com/maciejas22/conference-manager/api/internal/config"
)

type DB struct {
	Conn   *sqlx.DB
	logger *slog.Logger
}

type Logger struct {
	l *slog.Logger
}

func NewLogger(l *slog.Logger) *Logger {
	return &Logger{l: l}
}

func (l *Logger) Log(ctx context.Context, level tracelog.LogLevel, msg string, data map[string]interface{}) {
	attrs := make([]slog.Attr, 0, len(data))
	for k, v := range data {
		attrs = append(attrs, slog.Any(k, v))
	}

	var lvl slog.Level
	switch level {
	case tracelog.LogLevelTrace:
		lvl = slog.LevelDebug - 1
		attrs = append(attrs, slog.Any("PGX_LOG_LEVEL", level))
	case tracelog.LogLevelDebug:
		lvl = slog.LevelDebug
	case tracelog.LogLevelInfo:
		lvl = slog.LevelInfo
	case tracelog.LogLevelWarn:
		lvl = slog.LevelWarn
	case tracelog.LogLevelError:
		lvl = slog.LevelError
	default:
		lvl = slog.LevelError
		attrs = append(attrs, slog.Any("INVALID_PGX_LOG_LEVEL", level))
	}
	l.l.LogAttrs(ctx, lvl, msg, attrs...)
}

type MultiQueryTracer struct {
	Tracers []pgx.QueryTracer
}

func (m *MultiQueryTracer) TraceQueryStart(ctx context.Context, conn *pgx.Conn, data pgx.TraceQueryStartData) context.Context {
	for _, t := range m.Tracers {
		ctx = t.TraceQueryStart(ctx, conn, data)
	}

	return ctx
}

func (m *MultiQueryTracer) TraceQueryEnd(ctx context.Context, conn *pgx.Conn, data pgx.TraceQueryEndData) {
	for _, t := range m.Tracers {
		t.TraceQueryEnd(ctx, conn, data)
	}
}

func GetDriverConfig(l *slog.Logger) *pgx.ConnConfig {
	connConfig, _ := pgx.ParseConfig(config.AppConfig.DatabaseURL)
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

func SqlConnect(c context.Context, l *slog.Logger) (*sqlx.DB, error) {
	connConfig, _ := pgx.ParseConfig(config.AppConfig.DatabaseURL)
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

	connStr := stdlib.RegisterConnConfig(connConfig)
	conn, err := sqlx.Connect("pgx", connStr)
	if err != nil {
		return nil, err
	}

	return conn, nil
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
		logger: logger,
	}, nil
}

func (db *DB) Close() (err error) {
	if sqlErr := db.Conn.Close(); sqlErr != nil {
		db.logger.Error("failed to close db connection", "error", sqlErr)
		return sqlErr
	}
	db.logger.Info("closed connection to db")

	return nil
}

func Transaction(ctx context.Context, db *sqlx.DB, fn func(*sqlx.Tx) error) error {
	tx, err := db.BeginTxx(ctx, nil)
	if err != nil {
		return err
	}

	err = fn(tx)

	if err != nil {
		if rbErr := tx.Rollback(); rbErr != nil {
			return rbErr
		}
		return err
	}

	if err := tx.Commit(); err != nil {
		return err
	}

	return nil
}
