package store

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	sql_driver "database/sql/driver"
	"github.com/go-sql-driver/mysql"
	"github.com/shoet/blog/config"
	"github.com/shoet/blog/logging"

	"github.com/jmoiron/sqlx"
	"github.com/qustavo/sqlhooks/v2"
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

	// true/falseの文字列指定のためboolから変換
	tlsConfigString := "false"
	if cfg.DBTlsEnabled {
		tlsConfigString = "true"
	}

	// register sql query logger
	withHooksDriverName, err := InitSQLDriverWithLogs("mysql", &mysql.MySQLDriver{})
	if err != nil {
		return nil, fmt.Errorf("failed init sql driver: %w", err)
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
		TLSConfig:            tlsConfigString,
	}
	db, err := sql.Open(withHooksDriverName, config.FormatDSN())
	if err != nil {
		return nil, fmt.Errorf("failed open mysql: %w", err)
	}
	if err := db.PingContext(ctx); err != nil {
		return nil, fmt.Errorf("failed connect mysql: %w", err)
	}
	xdb := sqlx.NewDb(db, "mysql")
	return xdb, nil
}

func InitSQLDriverWithLogs(driverName string, driver sql_driver.Driver) (string, error) {
	hooksDriverName := fmt.Sprintf("%sWithLoggerHooks", driverName)
	for _, existingDriverName := range sql.Drivers() {
		if existingDriverName == hooksDriverName {
			return hooksDriverName, fmt.Errorf("driver %s already exists", hooksDriverName)
		}
	}
	sql.Register(hooksDriverName, sqlhooks.Wrap(driver, &SQLQueryLoggerHooks{}))
	return hooksDriverName, nil
}

type SQLQueryLoggerHooks struct{}

func (s *SQLQueryLoggerHooks) Before(
	ctx context.Context, query string, args ...interface{},
) (context.Context, error) {
	logger := logging.GetLogger(ctx)
	logger.Info().Msgf("query: %s", query)
	return ctx, nil
}

func (s *SQLQueryLoggerHooks) After(
	ctx context.Context, query string, args ...interface{},
) (context.Context, error) {
	return ctx, nil
}

var _ sqlhooks.Hooks = (*SQLQueryLoggerHooks)(nil)
