package wallet_test

import (
	"context"
	"errors"
	"github.com/eneskzlcn/currency-conversion-service/internal/entity"
	mocks "github.com/eneskzlcn/currency-conversion-service/internal/mocks/wallet"
	"github.com/eneskzlcn/currency-conversion-service/internal/wallet"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestNewService(t *testing.T) {
	t.Run("given not valid arguments then it should return nil", func(t *testing.T) {
		assert.Nil(t, wallet.NewService(nil))
	})
	t.Run("given valid arguments then it should return new Service", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		mockWalletRepo := mocks.NewMockWalletRepository(ctrl)
		assert.NotNil(t, wallet.NewService(mockWalletRepo))
	})
}

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
func createServiceWithMockWalletRepository(t *testing.T) (*wallet.Service, *mocks.MockWalletRepository) {
	ctrl := gomock.NewController(t)
	mockWalletRepo := mocks.NewMockWalletRepository(ctrl)
	return wallet.NewService(mockWalletRepo), mockWalletRepo
}
