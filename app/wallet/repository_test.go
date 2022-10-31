package wallet_test

import (
	"context"
	"errors"
	"github.com/eneskzlcn/currency-conversion-service/app/entity"
	"github.com/eneskzlcn/currency-conversion-service/app/wallet"
	"github.com/eneskzlcn/currency-conversion-service/postgres"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
	"regexp"
	"testing"
	"time"
)

func TestRepository_IsUserWithUserIDExists(t *testing.T) {
	db, sqlmock := postgres.NewMockPostgres()
	repository := wallet.NewRepository(db, zap.S())
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
	repository := wallet.NewRepository(db, zap.S())
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
func TestRepository_GetUserBalanceOnGivenCurrency(t *testing.T) {
	db, sqlmock := postgres.NewMockPostgres()
	repository := wallet.NewRepository(db, zap.S())
	query := regexp.QuoteMeta(`SELECT balance FROM user_wallets WHERE user_id = $1 AND currency_code = $2`)
	t.Run("given existing wallet with user id and currency then it should return balance", func(t *testing.T) {
		userID := 2
		currency := "TRY"
		expectedBalance := float32(200)
		sqlmock.ExpectQuery(query).WithArgs(userID, currency).
			WillReturnRows(sqlmock.NewRows([]string{"balance"}).
				AddRow(expectedBalance))
		balance, err := repository.GetUserBalanceOnGivenCurrency(context.TODO(), userID, currency)
		assert.Nil(t, err)
		assert.Equal(t, expectedBalance, balance)
	})
	t.Run("given not existing wallet with user id and currency then it should return -1 with error", func(t *testing.T) {
		userID := 2
		currency := "TRY"
		expectedBalance := float32(-1)
		sqlmock.ExpectQuery(query).WithArgs(userID, currency).WillReturnRows(sqlmock.NewRows([]string{"balance"}))
		balance, err := repository.GetUserBalanceOnGivenCurrency(context.TODO(), userID, currency)
		assert.NotNil(t, err)
		assert.Equal(t, expectedBalance, balance)
	})
}
func TestRepository_AdjustUserBalanceOnGivenCurrency(t *testing.T) {
	db, sqlmock := postgres.NewMockPostgres()
	repository := wallet.NewRepository(db, zap.S())
	query := regexp.QuoteMeta(`
	UPDATE user_wallets 
	SET balance = balance + $1 
	WHERE user_id = $2 AND currency_code = $3`)

	t.Run("given existing wallet with user_id and currency_code then it should return true", func(t *testing.T) {
		userID := 2
		currency := "TRY"
		balance := float32(200)
		sqlmock.ExpectQuery(query).WithArgs(balance, userID, currency).WillReturnRows(sqlmock.NewRows([]string{}))
		success, err := repository.AdjustUserBalanceOnGivenCurrency(context.TODO(),
			userID, currency, balance)
		assert.Nil(t, sqlmock.ExpectationsWereMet())
		assert.Nil(t, err)
		assert.True(t, success)
	})
	t.Run("given non existing wallet with user_id and currency_code then  it should return false", func(t *testing.T) {
		userID := -1
		currency := "TsxY"
		balance := float32(500)
		sqlmock.ExpectQuery(query).WithArgs(balance, userID, currency).WillReturnError(errors.New("wallet not exist"))
		success, err := repository.AdjustUserBalanceOnGivenCurrency(context.TODO(),
			userID, currency, balance)
		assert.NotNil(t, err)
		assert.False(t, success)
	})
}
