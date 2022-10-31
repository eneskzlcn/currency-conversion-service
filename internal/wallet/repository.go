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
	panic("implement me")
}
