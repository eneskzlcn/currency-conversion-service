package auth

import (
	"context"
	"database/sql"
	"github.com/eneskzlcn/currency-conversion-service/app/entity"
	"go.uber.org/zap"
)

type Repository struct {
	db     *sql.DB
	logger *zap.SugaredLogger
}

func NewRepository(db *sql.DB, logger *zap.SugaredLogger) *Repository {
	return &Repository{db: db, logger: logger}
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
