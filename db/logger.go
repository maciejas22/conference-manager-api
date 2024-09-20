package db

import (
	"context"
	"database/sql"
	"fmt"
	"log/slog"

	"github.com/jmoiron/sqlx"
)

type QueryExecutor struct {
	queryer sqlx.Queryer
	execer  sqlx.Execer
	logger  *slog.Logger
}

func printQuery(l *slog.Logger, query string, args ...interface{}) {
	argsStr := fmt.Sprintf("%v", args...)
	l.Debug("", slog.String("query", query), slog.String("args", argsStr))
}

func (p *QueryExecutor) Query(query string, args ...interface{}) (*sql.Rows, error) {
	printQuery(p.logger, query, args...)
	res, err := p.queryer.Query(query, args...)
	if err != nil {
		p.logger.Error("failed to execute query", "error", err)
	}
	return res, err
}

func (p *QueryExecutor) Queryx(query string, args ...interface{}) (*sqlx.Rows, error) {
	printQuery(p.logger, query, args...)
	res, err := p.queryer.Queryx(query, args...)
	if err != nil {
		p.logger.Error("failed to execute query", "error", err)
	}
	return res, err
}

func (p *QueryExecutor) QueryRowx(query string, args ...interface{}) *sqlx.Row {
	printQuery(p.logger, query, args...)
	row := p.queryer.QueryRowx(query, args...)
	if row.Err() != nil {
		p.logger.Error("failed to execute query", "error", row.Err())
	}
	return row
}

func (p *QueryExecutor) Exec(query string, args ...interface{}) (sql.Result, error) {
	printQuery(p.logger, query, args...)
	res, err := p.execer.Exec(query, args...)
	if err != nil {
		p.logger.Error("failed to execute query", "error", err)
	}
	return res, err
}

func (ql *QueryExecutor) BeginTxx(ctx context.Context, opts *sql.TxOptions) (*sqlx.Tx, error) {
	tx, err := ql.queryer.(*sqlx.DB).BeginTxx(ctx, opts)
	if err != nil {
		ql.logger.Error("failed to begin transaction", "error", err)
		return nil, err
	}
	return tx, err
}

func (ql *QueryExecutor) WithTx(tx *sqlx.Tx) *QueryExecutor {
	return &QueryExecutor{
		queryer: tx,
		execer:  tx,
		logger:  ql.logger,
	}
}
