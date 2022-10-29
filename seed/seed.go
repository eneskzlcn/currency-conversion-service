package seed

import (
	"context"
	"database/sql"
	_ "embed"
)

//go:embed create_seed.sql
var createSeedSql string

//go:embed drop_seed.sql
var dropSeedSql string

type DB interface {
	ExecContext(context.Context, string, ...interface{}) (sql.Result, error)
}

func MigrateTables(ctx context.Context, db DB) error {
	_, err := db.ExecContext(ctx, createSeedSql)
	return err
}

func DropTables(ctx context.Context, db DB) error {
	_, err := db.ExecContext(ctx, dropSeedSql)
	return err
}
