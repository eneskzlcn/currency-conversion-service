//go:generate mockgen -destination=./mock_postgres.go -source=./postgres.go -package=postgres
package postgres

import (
	_ "context"
	"database/sql"
	_ "errors"
	"fmt"
	"github.com/eneskzlcn/currency-conversion-service/internal/config"
	_ "github.com/lib/pq"
)

func New(config config.DB) (*sql.DB, error) {
	db, err := sql.Open("postgres", createDSN(config))
	if err != nil {
		return nil, err
	}
	if err = db.Ping(); err != nil {
		return nil, err
	}
	return db, nil
}

func createDSN(config config.DB) string {
	return fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		config.Host, config.Port, config.Username, config.Password, config.DBName)
}
