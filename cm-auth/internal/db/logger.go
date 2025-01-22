package db

import (
	"context"
	"log/slog"

	"github.com/jackc/pgx/v5"
)

type myQueryTracer struct {
	l *slog.Logger
}

func (tracer *myQueryTracer) TraceQueryStart(
	ctx context.Context,
	_ *pgx.Conn,
	data pgx.TraceQueryStartData) context.Context {
	tracer.l.Info("DB", "sql", data.SQL, "args", data.Args)

	return ctx
}

func (tracer *myQueryTracer) TraceQueryEnd(ctx context.Context, conn *pgx.Conn, data pgx.TraceQueryEndData) {
}
