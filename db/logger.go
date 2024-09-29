package db

import (
	"context"
	"log/slog"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/tracelog"
)

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
