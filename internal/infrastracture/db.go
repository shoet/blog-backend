package infrastracture

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	sql_driver "database/sql/driver"

	"github.com/go-sql-driver/mysql"
	_ "github.com/jackc/pgx/stdlib"
	"github.com/shoet/blog/internal/config"
	"github.com/shoet/blog/internal/logging"

	"github.com/jmoiron/sqlx"
	"github.com/qustavo/sqlhooks/v2"
)

// type DB = *sqlx.DB

type DB interface {
	BeginTxx(ctx context.Context, opts *sql.TxOptions) (*sqlx.Tx, error)
	TX
}

type TX interface {
	QueryRowxContext(ctx context.Context, query string, args ...interface{}) *sqlx.Row
	SelectContext(ctx context.Context, dest interface{}, query string, args ...interface{}) error
	ExecContext(ctx context.Context, query string, args ...interface{}) (sql.Result, error)
}

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
	xdb := sqlx.NewDb(db, withHooksDriverName)
	return xdb, nil
}

func NewDBPostgres(ctx context.Context, cfg *config.Config) (*sqlx.DB, error) {
	dbDsn := fmt.Sprintf(
		"postgres://%s:%s@%s:%d/%s",
		cfg.DBUser, cfg.DBPass, cfg.DBHost, cfg.DBPort, cfg.DBName)

	if cfg.DBSSLMode != "" {
		dbDsn += fmt.Sprintf("?sslmode=%s", cfg.DBSSLMode)
	}

	db, err := sql.Open("pgx", dbDsn)
	if err != nil {
		return nil, fmt.Errorf("failed open postgres: %w", err)
	}

	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("failed connect postgres: %w", err)
	}

	return sqlx.NewDb(db, "pgx"), nil
}

func InitSQLDriverWithLogs(driverName string, driver sql_driver.Driver) (string, error) {
	hooksDriverName := fmt.Sprintf("%sWithLoggerHooks", driverName)
	for _, existingDriverName := range sql.Drivers() {
		if existingDriverName == hooksDriverName {
			fmt.Printf("driver %s already exists\n", hooksDriverName)
			return hooksDriverName, nil
		}
	}
	sql.Register(hooksDriverName, sqlhooks.Wrap(driver, &SQLQueryLoggerHooks{}))
	return hooksDriverName, nil
}

type TransactionProvider struct {
	db DB
}

func NewTransactionProvider(db DB) *TransactionProvider {
	return &TransactionProvider{db: db}
}

var DBTxKey = struct{}{}
var ErrNoTransaction = fmt.Errorf("no transaction")

func (t *TransactionProvider) DoInTx(
	ctx context.Context, f func(tx TX) (interface{}, error),
) (interface{}, error) {
	tx, err := t.db.BeginTxx(ctx, nil)
	if err != nil {
		return nil, fmt.Errorf("failed begin transaction: %w", err)
	}
	v, err := f(tx)
	if err != nil {
		if err := tx.Rollback(); err != nil {
			return nil, fmt.Errorf("failed rollback transaction: %w", err)
		}
		return nil, fmt.Errorf("failed transaction: %w", err)
	}
	if err := tx.Commit(); err != nil {
		if err := tx.Rollback(); err != nil {
			return nil, fmt.Errorf("failed rollback transaction: %w", err)
		}
		return nil, fmt.Errorf("failed commit transaction: %w", err)
	}
	return v, nil
}

type SQLQueryLoggerHooks struct{}

func (s *SQLQueryLoggerHooks) Before(
	ctx context.Context, query string, args ...interface{},
) (context.Context, error) {
	logger := logging.GetLogger(ctx)
	logger.Debug(fmt.Sprintf("query: %s\nargs: %v", query, args))
	return ctx, nil
}

func (s *SQLQueryLoggerHooks) After(
	ctx context.Context, query string, args ...interface{},
) (context.Context, error) {
	return ctx, nil
}

var _ sqlhooks.Hooks = (*SQLQueryLoggerHooks)(nil)
