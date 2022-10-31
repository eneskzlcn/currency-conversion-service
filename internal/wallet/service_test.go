package wallet_test

import (
	"context"
	"errors"
	"github.com/eneskzlcn/currency-conversion-service/internal/entity"
	mocks "github.com/eneskzlcn/currency-conversion-service/internal/mocks/wallet"
	"github.com/eneskzlcn/currency-conversion-service/internal/wallet"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
	"testing"
)

func TestService_GetUserWalletAccounts(t *testing.T) {
	service, mockWalletRepo := createServiceWithMockWalletRepository(t)
	t.Run("given not existing user id then it should return nil and error", func(t *testing.T) {
		userID := -34
		mockWalletRepo.EXPECT().IsUserWithUserIDExists(gomock.Any(), userID).
			Return(false, nil)
		userWalletAccountsResp, err := service.GetUserWalletAccounts(context.TODO(), userID)
		assert.Empty(t, userWalletAccountsResp)
		assert.NotNil(t, err)
		assert.True(t, errors.Is(err, wallet.UserWithUserIDNotExistsErr))
	})
	t.Run("given existing user id then it should return users wallet accounts", func(t *testing.T) {
		userID := 2
		mockUserWalletAccounts := []entity.UserWallet{
			{
				UserID:   userID,
				Currency: "TRY",
				Balance:  200,
			},
			{
				UserID:   userID,
				Currency: "USD",
				Balance:  2000,
			},
		}
		expectedUserWalletAccountsResp := wallet.UserWalletAccountsResponse{Accounts: []wallet.UserWalletAccount{
			{
				Currency: mockUserWalletAccounts[0].Currency,
				Balance:  mockUserWalletAccounts[0].Balance,
			},
			{
				Currency: mockUserWalletAccounts[1].Currency,
				Balance:  mockUserWalletAccounts[1].Balance,
			},
		}}
		mockWalletRepo.EXPECT().IsUserWithUserIDExists(gomock.Any(), userID).
			Return(true, nil)
		mockWalletRepo.EXPECT().GetUserWalletAccounts(gomock.Any(), userID).
			Return(mockUserWalletAccounts, nil)

		userWalletAccountsResp, err := service.GetUserWalletAccounts(context.TODO(), userID)

		assert.Nil(t, err)
		assert.Equal(t, expectedUserWalletAccountsResp, userWalletAccountsResp)

	})
}

func TestService_GetUserBalanceOnGivenCurrency(t *testing.T) {
	service, mockWalletRepo := createServiceWithMockWalletRepository(t)
	t.Run("given not existing user id then it should return -1 and error", func(t *testing.T) {
		userID := 2
		currency := "TRY"
		expectedBalance := float32(-1)
		mockWalletRepo.EXPECT().GetUserBalanceOnGivenCurrency(gomock.Any(), userID, currency).
			Return(expectedBalance, errors.New("user not found"))
		balance, err := service.GetUserBalanceOnGivenCurrency(context.TODO(), userID, currency)
		assert.NotNil(t, err)
		assert.LessOrEqual(t, balance, expectedBalance)
	})
	t.Run("given existing user id then it should return balance without error", func(t *testing.T) {
		userID := 2
		currency := "TRY"
		expectedBalance := float32(200)
		mockWalletRepo.EXPECT().GetUserBalanceOnGivenCurrency(gomock.Any(), userID, currency).
			Return(expectedBalance, nil)
		balance, err := service.GetUserBalanceOnGivenCurrency(context.TODO(), userID, currency)
		assert.Nil(t, err)
		assert.Equal(t, balance, expectedBalance)
	})
}

func TestService_AdjustUserBalanceOnGivenCurrency(t *testing.T) {
	service, mockWalletRepo := createServiceWithMockWalletRepository(t)

	t.Run("given not existing wallet with user id and currency then it should return false", func(t *testing.T) {
		userID := -1
		currency := "T"
		balance := float32(200)
		mockWalletRepo.EXPECT().AdjustUserBalanceOnGivenCurrency(gomock.Any(), userID,
			currency, balance).Return(false, errors.New(""))
		success, err := service.AdjustUserBalanceOnGivenCurrency(context.TODO(),
			userID, currency, balance)
		assert.NotNil(t, err)
		assert.False(t, success)
	})
	t.Run("given existing wallet with user id and currency then it should return true", func(t *testing.T) {
		userID := -1
		currency := "T"
		balance := float32(200)
		mockWalletRepo.EXPECT().AdjustUserBalanceOnGivenCurrency(gomock.Any(), userID,
			currency, balance).Return(true, nil)
		success, err := service.AdjustUserBalanceOnGivenCurrency(context.TODO(),
			userID, currency, balance)
		assert.Nil(t, err)
		assert.True(t, success)
	})
}

func createServiceWithMockWalletRepository(t *testing.T) (*wallet.Service, *mocks.MockWalletRepository) {
	ctrl := gomock.NewController(t)
	mockWalletRepo := mocks.NewMockWalletRepository(ctrl)
	return wallet.NewService(mockWalletRepo, zap.S()), mockWalletRepo
}
