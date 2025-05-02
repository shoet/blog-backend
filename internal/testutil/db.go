package testutil

import (
	"context"
	"database/sql"
	"fmt"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/go-sql-driver/mysql"
	_ "github.com/jackc/pgx/v5"
	"github.com/jmoiron/sqlx"
	migrate "github.com/rubenv/sql-migrate"
	"github.com/shoet/blog/internal/util"
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
		Addr:                 "127.0.0.1:33061",
		User:                 "blog",
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

func NewDBPostgreSQLForTest(t *testing.T, ctx context.Context) (*sqlx.DB, error) {
	t.Helper()
	dbDsn := fmt.Sprintf(
		"postgres://%s:%s@%s:%d/%s",
		"blog", "blog", "127.0.0.1", 54321, "blog")

	db, err := sql.Open("pgx", dbDsn)
	if err != nil {
		return nil, fmt.Errorf("failed open postgres: %w", err)
	}

	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("failed connect postgres: %w", err)
	}

	return sqlx.NewDb(db, "pgx"), nil
}

func RepositoryTestPrepare(t *testing.T, ctx context.Context, db *sqlx.DB) {
	t.Helper()

	curDir, err := os.Getwd()
	if err != nil {
		t.Fatalf("failed to get current directory: %v", err)
	}
	cwd, err := util.GetProjectRoot(curDir)
	if err != nil {
		t.Fatalf("failed to get project root: %v", err)
	}

	sqlDir := filepath.Join(cwd, "_tools/migrations/postgres")
	migrations := &migrate.FileMigrationSource{
		Dir: sqlDir,
	}

	_, err = migrate.ExecContext(ctx, db.DB, "postgres", migrations, migrate.Up)
	if err != nil {
		t.Fatalf("failed to migrate: %v", err)
	}
}
