package store

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/go-sql-driver/mysql"
	_ "github.com/go-sql-driver/mysql"
	_ "github.com/mattn/go-sqlite3"
	"github.com/shoet/blog/config"

	"github.com/jmoiron/sqlx"
)

func NewDBSQLite3(ctx context.Context) (*sqlx.DB, error) {
	db, err := sql.Open("sqlite3", "database.sqlite")
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

func NewDBMySQL(ctx context.Context, cfg *config.Config) (*sqlx.DB, error) {
	jst, err := time.LoadLocation("Asia/Tokyo")
	if err != nil {
		return nil, fmt.Errorf("failed load location: %w", err)
	}
	config := mysql.Config{
		Addr:                 fmt.Sprintf("%s:%d", cfg.DBHost, cfg.DBPort),
		User:                 cfg.DBUser,
		Passwd:               cfg.DBPass,
		DBName:               cfg.DBName,
		Net:                  "tcp",
		ParseTime:            true,
		Loc:                  jst,
		AllowNativePasswords: true,
		TLSConfig:            "true",
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
