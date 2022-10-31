package wallet

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

func (r *Repository) IsUserWithUserIDExists(ctx context.Context, userID int) (bool, error) {
	query := `SELECT EXISTS (SELECT 1 FROM users WHERE id = $1)`
	row := r.db.QueryRowContext(ctx, query, userID)
	var exists bool
	if err := row.Scan(&exists); err != nil {
		return false, err
	}
	return exists, nil
}

func (r *Repository) GetUserWalletAccounts(ctx context.Context, userID int) ([]entity.UserWallet, error) {
	query := `
		SELECT user_id, currency_code, balance, created_at, updated_at
		FROM user_wallets uw 
		WHERE uw.user_id = $1`
	rows, err := r.db.QueryContext(ctx, query, userID)
	if err != nil {
		return nil, err
	}
	userWallets := make([]entity.UserWallet, 0)
	for rows.Next() {
		var wallet entity.UserWallet
		err = rows.Scan(&wallet.UserID,
			&wallet.Currency, &wallet.Balance, &wallet.CreatedAt,
			&wallet.UpdatedAt)
		if err != nil {
			return nil, err
		}
	}
	if err = rows.Close(); err != nil {
		return nil, err
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}
	return userWallets, nil
}

func (r *Repository) GetUserBalanceOnGivenCurrency(ctx context.Context, userID int, currency string) (float32, error) {
	query := `SELECT balance FROM user_wallets WHERE user_id = $2 AND currency_code = $2`
	row := r.db.QueryRowContext(ctx, query, userID, currency)
	var balance float32
	if err := row.Scan(&balance); err != nil {
		return -1, err
	}
	return balance, nil
}
func (r *Repository) AdjustUserBalanceOnGivenCurrency(ctx context.Context, userID int, currency string, balance float32) (bool, error) {
	panic("implement me")
}
