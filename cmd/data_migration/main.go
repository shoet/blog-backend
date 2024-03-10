package main

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/caarlos0/env/v10"
	"github.com/jmoiron/sqlx"
	"github.com/shoet/blog/internal/infrastracture"
	"github.com/shoet/blog/internal/infrastracture/models"
)

type MigrateConfig struct {
	SrcDSN string `env:"SRC_DSN,required"`
	DstDSN string `env:"DST_DSN,required"`
}

func StreamMigrationSource(
	ctx context.Context,
	src infrastracture.DB,
	dst infrastracture.TX,
	scanQuery string,
	writeFunc func(ctx context.Context, tx infrastracture.TX, rows *sqlx.Rows) error,
) error {
	rows, err := src.QueryxContext(ctx, scanQuery)
	if err != nil {
		return fmt.Errorf("failed to read from source database: %w", err)
	}

	cnt := 0
	for rows.Next() {
		if err := writeFunc(ctx, dst, rows); err != nil {
			return fmt.Errorf("failed to write to destination database: %w", err)
		}
		cnt++
	}
	fmt.Println("migrated: ", cnt)
	return nil
}

func SetupDB(srcDSN, dstDSN string) (src, dst *sql.DB, err error) {
	src, err = ConnectDB("mysql", srcDSN)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to connect to source database: %w", err)
	}
	dst, err = ConnectDB("pgx", dstDSN)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to connect to destination database: %w", err)
	}
	return src, dst, nil
}

func ConnectDB(driverName string, dsn string) (*sql.DB, error) {
	db, err := sql.Open(driverName, dsn)
	if err != nil {
		return nil, fmt.Errorf("failed to open source database: %w", err)
	}
	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("failed to ping source database: %w", err)
	}
	return db, nil
}

func ReadConfig() (*MigrateConfig, error) {
	var config MigrateConfig
	if err := env.Parse(&config); err != nil {
		return nil, fmt.Errorf("failed to read config: %w", err)
	}
	return &config, nil
}

type MigrationInput struct {
	src       *sqlx.DB
	dst       *sqlx.DB
	name      string
	scanQuery string
	writeFunc func(ctx context.Context, tx infrastracture.TX, rows *sqlx.Rows) error
}

func Migration(input *MigrationInput) error {
	ctx := context.Background()
	dstTx, err := input.src.BeginTxx(ctx, nil)
	if err != nil {
		panic(err)
	}
	defer dstTx.Rollback()
	if err := StreamMigrationSource(
		ctx, input.src, dstTx, input.scanQuery, input.writeFunc,
	); err != nil {
		return fmt.Errorf("failed to migrate: %w", err)
	}
	if err := dstTx.Commit(); err != nil {
		return fmt.Errorf("failed to commit to destination database: %w", err)
	}
	return nil
}

func main() {
	cfg, err := ReadConfig()
	if err != nil {
		panic(err)
	}

	src, dst, err := SetupDB(cfg.SrcDSN, cfg.DstDSN)
	if err != nil {
		panic(err)
	}

	srcx := sqlx.NewDb(src, "mysql")
	dstx := sqlx.NewDb(dst, "pgx")
	defer src.Close()
	defer dst.Close()

	for _, m := range []*MigrationInput{
		{
			name:      "blogs",
			src:       srcx,
			dst:       dstx,
			scanQuery: `SELECT * FROM blogs ORDER BY id;`,
			writeFunc: func(ctx context.Context, tx infrastracture.TX, rows *sqlx.Rows) error {
				dstQuery := `
				INSERT INTO
					blogs
				(
					id, title, description, content, author_id, 
					thumbnail_image_file_name, is_public, tags,
					created, modified
				)
				VALUES
					($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)
				`
				rows, err := tx.Queryx(dstQuery)
				if err != nil {
					return fmt.Errorf("failed to insert to destination database: %w", err)
				}
				for rows.Next() {
					var blog models.Blog
					if err := rows.StructScan(&blog); err != nil {
						return fmt.Errorf("failed to scan to destination database: %w", err)
					}
					if _, err := tx.ExecContext(
						ctx, dstQuery,
						blog.Id, blog.Title, blog.Description, blog.Content, blog.AuthorId,
						blog.ThumbnailImageFileName, blog.IsPublic, blog.Tags,
						blog.Created, blog.Modified,
					); err != nil {
						return fmt.Errorf("failed to insert to destination database: %w", err)
					}
				}
				return nil
			},
		},
	} {
		fmt.Println("start migration: ", m.name)
		if err := Migration(m); err != nil {
			panic(err)
		}
	}

}
