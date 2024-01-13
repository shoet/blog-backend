package repository

import (
	"context"
	"database/sql"

	"github.com/jmoiron/sqlx"
)

type Queryer interface {
	QueryRowxContext(ctx context.Context, query string, args ...interface{}) *sqlx.Row
	SelectContext(ctx context.Context, dest interface{}, query string, args ...interface{}) error
}

type Execer interface {
	Queryer
	ExecContext(ctx context.Context, query string, args ...interface{}) (sql.Result, error)
}

var _ Execer = (*sqlx.Tx)(nil)
var _ Execer = (*sqlx.DB)(nil)
var _ Queryer = (*sqlx.Tx)(nil)
var _ Queryer = (*sqlx.DB)(nil)
