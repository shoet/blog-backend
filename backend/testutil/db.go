package testutil

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/jmoiron/sqlx"
)

func NewDBSQLite3ForTest(ctx context.Context) (*sqlx.DB, error) {
	db, err := sql.Open("sqlite3", ":memory:")
	if err != nil {
		return nil, fmt.Errorf("failed open sqlite3: %w", err)
	}
	if err := db.PingContext(ctx); err != nil {
		return nil, fmt.Errorf("failed connect sqlite3: %w", err)
	}
	xdb := sqlx.NewDb(db, "sqlite3")
	if _, err := xdb.Exec("PRAGMA foreign_keys= true;"); err != nil {
		return nil, fmt.Errorf("failed set foreign_keys: %w", err)
	}
	return xdb, nil
}
