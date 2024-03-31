package repository

import (
	sq "github.com/Masterminds/squirrel"
)

/*
GetQueryBuilderPostgresはPostgreSQL用のQueryBuilderを返します。

PostgreSQLでのプレースホルダは$マークが使われるため。
*/
func GetQueryBuilderPostgres() sq.StatementBuilderType {
	psql := sq.StatementBuilder.PlaceholderFormat(sq.Dollar)
	return psql
}
