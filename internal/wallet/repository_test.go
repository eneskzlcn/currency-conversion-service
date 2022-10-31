package wallet_test

import (
	"context"
	"github.com/eneskzlcn/currency-conversion-service/internal/entity"
	"github.com/eneskzlcn/currency-conversion-service/internal/wallet"
	"github.com/eneskzlcn/currency-conversion-service/postgres"
	"github.com/stretchr/testify/assert"
	"regexp"
	"testing"
	"time"
)

func TestNewRepository(t *testing.T) {
	t.Run("given not valid arguments then it should return nil", func(t *testing.T) {
		assert.Nil(t, wallet.NewRepository(nil))
	})
	t.Run("given valid arguments then it should return new Repository", func(t *testing.T) {
		db, _ := postgres.NewMockPostgres()
		assert.NotNil(t, wallet.NewRepository(db))
	})
}
func TestRepository_IsUserWithUserIDExists(t *testing.T) {
	db, sqlmock := postgres.NewMockPostgres()
	repository := wallet.NewRepository(db)
	query := regexp.QuoteMeta(`SELECT EXISTS (SELECT 1 FROM users WHERE id = $1)`)
	t.Run("given existing user id then it should return true ", func(t *testing.T) {
		userID := 2
		expectedRows := sqlmock.NewRows([]string{"exists"}).AddRow(true)
		sqlmock.ExpectQuery(query).WithArgs(userID).
			WillReturnRows(expectedRows)
		exists, err := repository.IsUserWithUserIDExists(context.TODO(), userID)
		assert.Nil(t, sqlmock.ExpectationsWereMet())
		assert.Nil(t, err)
		assert.True(t, exists)
	})
	t.Run("given not existing user id then it should return false", func(t *testing.T) {
		userID := 4
		expectedRows := sqlmock.NewRows([]string{"exists"}).AddRow(false)
		sqlmock.ExpectQuery(query).WithArgs(userID).
			WillReturnRows(expectedRows)
		exists, err := repository.IsUserWithUserIDExists(context.TODO(), userID)
		assert.Nil(t, sqlmock.ExpectationsWereMet())
		assert.Nil(t, err)
		assert.False(t, exists)
	})
}
func TestRepository_GetUserWalletAccounts(t *testing.T) {
	db, sqlmock := postgres.NewMockPostgres()
	repository := wallet.NewRepository(db)
	query := regexp.QuoteMeta(`SELECT user_id, currency_code, balance, created_at, updated_at
			FROM user_wallets uw WHERE uw.user_id = $1`)

	t.Run("given user that not has wallet account then it should return empty wallet account array", func(t *testing.T) {
		userID := 3
		expectedRows := sqlmock.NewRows([]string{"user_id", "currency_code", "balance", "created_at", "updated_at"})
		sqlmock.ExpectQuery(query).WithArgs(userID).WillReturnRows(expectedRows)
		userWallets, err := repository.GetUserWalletAccounts(context.TODO(), userID)
		assert.Nil(t, err)
		assert.Empty(t, userWallets)
	})
	t.Run("given user that has wallet accounts then it should return wallet accounts array", func(t *testing.T) {
		userID := 4
		expectedRows := sqlmock.NewRows([]string{"user_id", "currency_code", "balance", "created_at", "updated_at"})
		expectedUserWallets := []entity.UserWallet{
			{
				UserID:    userID,
				Currency:  "TRY",
				Balance:   200,
				CreatedAt: time.Now(),
				UpdatedAt: time.Now(),
			},
			{
				UserID:    userID,
				Currency:  "USD",
				Balance:   300,
				CreatedAt: time.Now(),
				UpdatedAt: time.Now(),
			},
		}
		for _, userWallet := range expectedUserWallets {
			expectedRows.AddRow(userWallet.UserID, userWallet.Currency,
				userWallet.Balance, userWallet.CreatedAt, userWallet.UpdatedAt)
		}
		sqlmock.ExpectQuery(query).WithArgs(userID).WillReturnRows(expectedRows)
		userWallets, err := repository.GetUserWalletAccounts(context.TODO(), userID)
		assert.Nil(t, err)
		assert.Empty(t, userWallets)
	})
}
