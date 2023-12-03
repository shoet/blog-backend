package testutil

import (
	"context"
	"database/sql"
	"fmt"
	"testing"
	"time"

	"github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
)

func NewDBSQLite3ForTest(t *testing.T, ctx context.Context) (*sqlx.DB, error) {
	t.Helper()
	db, err := sql.Open("sqlite3", ":memory:")
	if err != nil {
		return nil, fmt.Errorf("failed open sqlite3: %w", err)
	}
	t.Cleanup(func() {
		db.Close()
	})
	if err := db.PingContext(ctx); err != nil {
		return nil, fmt.Errorf("failed connect sqlite3: %w", err)
	}
	xdb := sqlx.NewDb(db, "sqlite3")
	if _, err := xdb.Exec("PRAGMA foreign_keys= true;"); err != nil {
		return nil, fmt.Errorf("failed set foreign_keys: %w", err)
	}
	return xdb, nil
}

func NewDBMySQLForTest(t *testing.T, ctx context.Context) (*sqlx.DB, error) {
	t.Helper()
	jst, err := time.LoadLocation("Asia/Tokyo")
	if err != nil {
		return nil, fmt.Errorf("failed load location: %w", err)
	}
	config := mysql.Config{
		Addr:                 "127.0.0.1:3306",
		User:                 "blog_user",
		Passwd:               "blog",
		DBName:               "blog",
		Net:                  "tcp",
		ParseTime:            true,
		Loc:                  jst,
		AllowNativePasswords: true,
	}
	db, err := sql.Open("mysql", config.FormatDSN())
	if err != nil {
		return nil, fmt.Errorf("failed open mysql: %w", err)
	}
	t.Cleanup(func() {
		db.Close()
	})
	if err := db.PingContext(ctx); err != nil {
		return nil, fmt.Errorf("failed connect mysql: %w", err)
	}
	xdb := sqlx.NewDb(db, "mysql")
	return xdb, nil
}
