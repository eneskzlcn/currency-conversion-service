package wallet

import (
	"context"
	"database/sql"
)

type Repository struct {
	db *sql.DB
}

func NewRepository(db *sql.DB) *Repository {
	if db == nil {
		return nil
	}
	return &Repository{db: db}
}
func (r *Repository) IsUserWithUserIDExists(ctx context.Context, userID int) (bool, error) {
	query := `SELECT EXISTS (SELECT 1 FROM users WHERE id = $1`
	row := r.db.QueryRowContext(ctx, query, userID)
	var exists bool
	if err := row.Scan(&exists); err != nil {
		return false, err
	}
	return exists, nil
}
