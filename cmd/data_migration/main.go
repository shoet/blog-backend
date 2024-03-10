package main

import (
	"context"
	"database/sql"
	"fmt"
	"os"
	"time"

	"github.com/caarlos0/env/v10"
	"github.com/jmoiron/sqlx"
	"github.com/shoet/blog/internal/infrastracture"
	"github.com/shoet/blog/internal/infrastracture/models"
)

type MigrateConfig struct {
	SrcDSN string `env:"SRC_DSN,required"`
	DstDSN string `env:"DST_DSN,required"`
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
	tableName string
	scanQuery string
	writeFunc func(ctx context.Context, tx infrastracture.TX, rows *sqlx.Rows) error
}

func Migration(input *MigrationInput, tableName string, dryrun bool) error {
	ctx := context.Background()
	dstTx, err := input.dst.BeginTxx(ctx, nil)
	if err != nil {
		panic(err)
	}
	defer dstTx.Rollback()

	_, err = dstTx.ExecContext(ctx, fmt.Sprintf("TRUNCATE TABLE %s RESTART IDENTITY;", tableName))
	if err != nil {
		return fmt.Errorf("failed to truncate blogs: %w", err)
	}

	rows, err := input.src.QueryxContext(ctx, input.scanQuery)
	if err != nil {
		return fmt.Errorf("failed to read from source database: %w", err)
	}

	cnt := 0
	for rows.Next() {
		if err := input.writeFunc(ctx, dstTx, rows); err != nil {
			return fmt.Errorf("failed to write to destination database: %w", err)
		}
		cnt++
	}

	fmt.Println("migrated: ", cnt)
	if !dryrun {
		if err := dstTx.Commit(); err != nil {
			return fmt.Errorf("failed to commit to destination database: %w", err)
		}
	}
	return nil
}

type BeforeBlog struct {
	models.Blog
	Created  time.Time `db:"created"`
	Modified time.Time `db:"modified"`
}

func main() {
	cfg, err := ReadConfig()
	if err != nil {
		panic(err)
	}

	dryrun := false

	if len(os.Args) > 1 && os.Args[1] == "dryrun" {
		dryrun = true
	}

	fmt.Printf("srcDSN: %s\n", cfg.SrcDSN)
	fmt.Printf("dstDSN: %s\n", cfg.DstDSN)

	src, dst, err := SetupDB(cfg.SrcDSN, cfg.DstDSN)
	if err != nil {
		panic(err)
	}

	srcx := sqlx.NewDb(src, "mysql")
	dstx := sqlx.NewDb(dst, "pgx")
	defer src.Close()
	defer dst.Close()

	for _, m := range []MigrationInput{
		{
			tableName: "blogs",
			src:       srcx,
			dst:       dstx,
			scanQuery: `SELECT * FROM blogs ORDER BY id;`,
			writeFunc: func(ctx context.Context, tx infrastracture.TX, rows *sqlx.Rows) error {
				var blog BeforeBlog
				if err := rows.StructScan(&blog); err != nil {
					return fmt.Errorf("failed to scan to destination database: %w", err)
				}
				created := blog.Created.Unix()
				modified := blog.Modified.Unix()

				dstQuery := `
				INSERT INTO
					blogs
				(
					id, title, description, content, author_id, 
					thumbnail_image_file_name, is_public,
					created, modified
				)
				VALUES
					($1, $2, $3, $4, $5, $6, $7, $8, $9)
				`
				if _, err := tx.ExecContext(
					ctx,
					dstQuery,
					blog.Id, blog.Title, blog.Description, blog.Content, blog.AuthorId,
					blog.ThumbnailImageFileName, blog.IsPublic, created, modified,
				); err != nil {
					return fmt.Errorf("failed to insert to destination database: %w", err)
				}
				return nil
			},
		},
		{
			tableName: "tags",
			src:       srcx,
			dst:       dstx,
			scanQuery: `SELECT * FROM tags ORDER BY id;`,
			writeFunc: func(ctx context.Context, tx infrastracture.TX, rows *sqlx.Rows) error {
				var tag models.Tag
				if err := rows.StructScan(&tag); err != nil {
					return fmt.Errorf("failed to scan to destination database: %w", err)
				}

				dstQuery := `
				INSERT INTO
					tags
				(
					id, name
				)
				VALUES
					($1, $2)
				`
				if _, err := tx.ExecContext(
					ctx,
					dstQuery,
					tag.Id, tag.Name,
				); err != nil {
					return fmt.Errorf("failed to insert to destination database: %w", err)
				}
				return nil
			},
		},
		{
			tableName: "blogs_tags",
			src:       srcx,
			dst:       dstx,
			scanQuery: `SELECT blog_id, tag_id FROM blogs_tags ORDER BY id;`,
			writeFunc: func(ctx context.Context, tx infrastracture.TX, rows *sqlx.Rows) error {
				var blogsTags models.BlogsTags
				if err := rows.StructScan(&blogsTags); err != nil {
					return fmt.Errorf("failed to scan to destination database: %w", err)
				}

				dstQuery := `
				INSERT INTO
					blogs_tags
				(
					blog_id, tag_id
				)
				VALUES
					($1, $2)
				`
				if _, err := tx.ExecContext(
					ctx,
					dstQuery,
					blogsTags.BlogId, blogsTags.TagId,
				); err != nil {
					return fmt.Errorf("failed to insert to destination database: %w", err)
				}
				return nil
			},
		},
	} {
		fmt.Println("start migration: ", m.tableName)
		if err := Migration(&m, m.tableName, dryrun); err != nil {
			panic(err)
		}
	}
}
