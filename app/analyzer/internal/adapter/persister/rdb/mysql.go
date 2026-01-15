package rdb

import (
	"database/sql"

	_ "github.com/go-sql-driver/mysql"
)

func newDB() (*sql.DB, error) {
	return sql.Open("mysql", "user:password@/dbname")
}
