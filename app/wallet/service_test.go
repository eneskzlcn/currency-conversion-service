package wallet_test

import (
	"context"
	"errors"
	"github.com/eneskzlcn/currency-conversion-service/app/message"
	mocks "github.com/eneskzlcn/currency-conversion-service/app/mocks/wallet"
	"github.com/eneskzlcn/currency-conversion-service/app/model"
	"github.com/eneskzlcn/currency-conversion-service/app/wallet"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
	"testing"
)

type Service interface {
	GetUserWalletAccounts(ctx context.Context, userID int) (wallet.UserWalletAccountsResponse, error)
	GetUserBalanceOnGivenCurrency(ctx context.Context, userID int, currency string) (float32, error)
	TransferBalancesBetweenUserWallets(ctx context.Context, msg message.CurrencyConvertedMessage) error
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
		assert.Equal(t, err.Error(), wallet.UserWithUserIDNotExists)
	})
	t.Run("given existing user id then it should return users wallet accounts", func(t *testing.T) {
		userID := 2
		mockUserWalletAccounts := []model.UserWallet{
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

func createServiceWithMockWalletRepository(t *testing.T) (Service, *mocks.MockRepository) {
	ctrl := gomock.NewController(t)
	mockWalletRepo := mocks.NewMockRepository(ctrl)
	return wallet.NewService(mockWalletRepo, zap.S()), mockWalletRepo
}

func TestService_TransferBalancesBetweenUserWallets(t *testing.T) {
	service, mockRepo := createServiceWithMockWalletRepository(t)
	givenTransferMessage := message.CurrencyConvertedMessage{
		UserID:                   2,
		FromCurrency:             "TRY",
		ToCurrency:               "USD",
		SenderBalanceDecAmount:   200,
		ReceiverBalanceIncAmount: 10,
	}
	givenTransferDTO := wallet.NewTransferBalanceBetweenUserWalletsDTO(givenTransferMessage.UserID,
		givenTransferMessage.FromCurrency, givenTransferMessage.ToCurrency,
		givenTransferMessage.SenderBalanceDecAmount, givenTransferMessage.ReceiverBalanceIncAmount)
	t.Run("given transfer balance message but error occurred in repository then it should return error", func(t *testing.T) {
		mockRepo.EXPECT().TransferBalancesBetweenUserWallets(gomock.Any(), givenTransferDTO).
			Return(errors.New(""))
		err := service.TransferBalancesBetweenUserWallets(context.TODO(), givenTransferMessage)
		assert.NotNil(t, err)
	})
	t.Run("given transfer balance message then it should return nil", func(t *testing.T) {
		mockRepo.EXPECT().TransferBalancesBetweenUserWallets(gomock.Any(), givenTransferDTO).
			Return(nil)
		err := service.TransferBalancesBetweenUserWallets(context.TODO(), givenTransferMessage)
		assert.Nil(t, err)
	})
}
