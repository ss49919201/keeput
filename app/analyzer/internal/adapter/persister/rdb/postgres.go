package rdb

import (
	"database/sql"
	"fmt"

	_ "github.com/lib/pq"
	"github.com/ss49919201/keeput/app/analyzer/internal/config"
)

func NewDB() (*sql.DB, error) {
	dsn := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		config.PostgresHost(),
		config.PostgresPort(),
		config.PostgresUser(),
		config.PostgresPassword(),
		config.PostgresDBName(),
	)

	return sql.Open("postgres", dsn)
}
