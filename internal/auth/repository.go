package auth

import (
	"context"
	"database/sql"
	"github.com/eneskzlcn/currency-conversion-service/internal/entity"
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
func (r *Repository) GetUserByUsernameAndPassword(ctx context.Context, username string, password string) (entity.User, error) {
	query := `
		SELECT id, username, password, created_at, updated_at
		FROM users WHERE username = $1 AND password = $2`
	row := r.db.QueryRowContext(ctx, query, username, password)
	var user entity.User
	err := row.Scan(&user.ID, &user.Username, &user.Password, &user.CreatedAt, &user.UpdatedAt)
	if err != nil {
		return entity.User{}, err
	}
	if err = row.Err(); err != nil {
		return entity.User{}, err
	}
	return user, nil
}