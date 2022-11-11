package auth

import (
	"context"
	"database/sql"
	"github.com/eneskzlcn/currency-conversion-service/app/model"
	"go.uber.org/zap"
)

type postgresRepository struct {
	db     *sql.DB
	logger *zap.SugaredLogger
}

func NewPostgresRepository(db *sql.DB, logger *zap.SugaredLogger) *postgresRepository {
	return &postgresRepository{db: db, logger: logger}
}
func (r *postgresRepository) GetUserByUsernameAndPassword(ctx context.Context, username string, password string) (model.User, error) {
	query := `
		SELECT id, username, password, created_at, updated_at
		FROM users WHERE username = $1 AND password = $2`
	row := r.db.QueryRowContext(ctx, query, username, password)
	var user model.User
	err := row.Scan(&user.ID, &user.Username, &user.Password, &user.CreatedAt, &user.UpdatedAt)
	if err != nil {
		return model.User{}, err
	}
	if err = row.Err(); err != nil {
		return model.User{}, err
	}
	return user, nil
}
