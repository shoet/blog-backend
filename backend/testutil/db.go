package testutil

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/go-sql-driver/mysql"
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

func NewDBMySQLForTest(ctx context.Context) (*sqlx.DB, error) {
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
	if err := db.PingContext(ctx); err != nil {
		return nil, fmt.Errorf("failed connect mysql: %w", err)
	}
	xdb := sqlx.NewDb(db, "mysql")
	return xdb, nil
}
