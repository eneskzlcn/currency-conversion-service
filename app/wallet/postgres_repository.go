package wallet

import (
	"context"
	"database/sql"
	"github.com/eneskzlcn/currency-conversion-service/app/common/logutil"
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

func (r *postgresRepository) IsUserWithUserIDExists(ctx context.Context, userID int) (bool, error) {
	query := `SELECT EXISTS (SELECT 1 FROM users WHERE id = $1)`
	row := r.db.QueryRowContext(ctx, query, userID)
	var exists bool
	if err := row.Scan(&exists); err != nil {
		return false, err
	}
	return exists, nil
}

func (r *postgresRepository) GetUserWalletAccounts(ctx context.Context, userID int) ([]model.UserWallet, error) {
	query := `
		SELECT user_id, currency_code, balance, created_at, updated_at
		FROM user_wallets uw 
		WHERE uw.user_id = $1`
	rows, err := r.db.QueryContext(ctx, query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var userWallets []model.UserWallet
	for rows.Next() {
		var wallet model.UserWallet
		err = rows.Scan(
			&wallet.UserID,
			&wallet.Currency,
			&wallet.Balance,
			&wallet.CreatedAt,
			&wallet.UpdatedAt)
		if err != nil {
			return nil, err
		}
		userWallets = append(userWallets, wallet)
	}
	if err = rows.Close(); err != nil {
		return nil, err
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}
	return userWallets, nil
}

func (r *postgresRepository) GetUserBalanceOnGivenCurrency(ctx context.Context, userID int, currency string) (float32, error) {
	query := `SELECT balance FROM user_wallets WHERE user_id = $1 AND currency_code = $2`
	row := r.db.QueryRowContext(ctx, query, userID, currency)
	var balance float32
	if err := row.Scan(&balance); err != nil {
		return -1, err
	}
	return balance, nil
}
func (r *postgresRepository) AdjustUserBalanceOnGivenCurrency(ctx context.Context, userID int, currency string, balance float32) (bool, error) {
	panic("IMPLEMENT ME")
}

func (r *postgresRepository) TransferBalancesBetweenUserWallets(ctx context.Context,
	dto TransferBalanceBetweenUserWalletsDTO) error {
	tx, err := r.db.BeginTx(ctx, &sql.TxOptions{})
	if err != nil {
		return logutil.LogThenReturn(r.logger, err)
	}
	if err = r.adjustUserBalanceOnGivenCurrency(ctx, tx, dto.UserID(),
		dto.FromCurrency(), dto.SenderBalanceDecAmount()); err != nil {
		return logutil.LogThenReturn(r.logger, err)
	}
	if err = r.adjustUserBalanceOnGivenCurrency(ctx, tx, dto.UserID(),
		dto.ToCurrency(), dto.ReceiverBalanceIncAmount()); err != nil {
		return logutil.LogThenReturn(r.logger, err)
	}
	return tx.Commit()
}

func (r *postgresRepository) adjustUserBalanceOnGivenCurrency(ctx context.Context, tx *sql.Tx, userID int, currency string, balance float32) error {
	query := `
	UPDATE user_wallets
	SET balance = balance + $1
	WHERE user_id = $2 AND currency_code = $3`
	_, err := tx.ExecContext(ctx, query, balance, userID, currency)

	if err != nil {
		if rollBackErr := tx.Rollback(); rollBackErr != nil {
			return logutil.LogThenReturn(r.logger, err)
		}
		return logutil.LogThenReturn(r.logger, err)
	}
	return nil
}
