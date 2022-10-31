package conversion_test

import (
	"context"
	"errors"
	"github.com/eneskzlcn/currency-conversion-service/internal/conversion"
	mocks "github.com/eneskzlcn/currency-conversion-service/internal/mocks/conversion"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestNewService(t *testing.T) {
	t.Run("given empty repository then it should return nil", func(t *testing.T) {
		service := conversion.NewService(nil)
		assert.Nil(t, service)
	})
	t.Run("given valid arguments then it should return service", func(t *testing.T) {
		mockWalletRepo := mocks.NewMockWalletRepository(gomock.NewController(t))
		service := conversion.NewService(mockWalletRepo)
		assert.NotNil(t, service)
	})
}

func TestService_CreateCurrencyConversion(t *testing.T) {
	service, mockWalletRepository := createServiceWithMockConversionRepository(t)
	t.Run("given expired exchange offer request then it should return false with error", func(t *testing.T) {
		givenConversionOfferReq := conversion.CurrencyConversionOfferRequest{
			FromCurrency: "USD",
			ToCurrency:   "TRY",
			ExchangeRate: 2.30,
			CreatedAt:    time.Now(),
			ExpiresAt:    time.Now().Add(-10 * time.Minute).Unix(),
			Balance:      200,
		}
		userID := 2
		success, err := service.ConvertCurrencies(context.TODO(), userID, givenConversionOfferReq)
		assert.False(t, success)
		assert.NotNil(t, err)
		assert.True(t, errors.Is(err, conversion.CurrencyConversionOfferExpiredErr))
	})
	t.Run("given balance amount that user not have in its from currency account then it should return false with error", func(t *testing.T) {
		givenConversionOfferReq := conversion.CurrencyConversionOfferRequest{
			FromCurrency: "USD",
			ToCurrency:   "TRY",
			ExchangeRate: 2.30,
			CreatedAt:    time.Now(),
			ExpiresAt:    time.Now().Add(10 * time.Minute).Unix(),
			Balance:      200,
		}
		userID := 2

		mockWalletRepository.EXPECT().
			GetUserBalanceOnGivenCurrency(gomock.Any(), userID, givenConversionOfferReq.FromCurrency).
			Return(float32(100), nil)
		success, err := service.ConvertCurrencies(context.TODO(), userID, givenConversionOfferReq)
		assert.False(t, success)
		assert.NotNil(t, err)
		assert.True(t, errors.Is(err, conversion.NotEnoughBalanceForConversionOfferErr))
	})
	t.Run("given valid currency conversion offer then it should return true", func(t *testing.T) {
		givenConversionOfferReq := conversion.CurrencyConversionOfferRequest{
			FromCurrency: "USD",
			ToCurrency:   "TRY",
			ExchangeRate: 2.30,
			CreatedAt:    time.Now(),
			ExpiresAt:    time.Now().Add(10 * time.Minute).Unix(),
			Balance:      200,
		}
		userID := 2
		targetCurrencyBalanceAdjustAmount := givenConversionOfferReq.ExchangeRate * givenConversionOfferReq.Balance
		mockWalletRepository.EXPECT().
			GetUserBalanceOnGivenCurrency(gomock.Any(), userID, givenConversionOfferReq.FromCurrency).
			Return(float32(500), nil)

		mockWalletRepository.EXPECT().
			AdjustUserBalanceOnGivenCurrency(gomock.Any(), userID,
				givenConversionOfferReq.FromCurrency, givenConversionOfferReq.Balance).
			Return(true, nil)

		mockWalletRepository.EXPECT().
			AdjustUserBalanceOnGivenCurrency(gomock.Any(), userID,
				givenConversionOfferReq.ToCurrency, targetCurrencyBalanceAdjustAmount).
			Return(true, nil)

		success, err := service.ConvertCurrencies(context.TODO(), userID, givenConversionOfferReq)
		assert.Nil(t, err)
		assert.True(t, success)
	})
}
func createServiceWithMockConversionRepository(t *testing.T) (*conversion.Service, *mocks.MockWalletRepository) {
	ctrl := gomock.NewController(t)
	mockWalletRepo := mocks.NewMockWalletRepository(ctrl)
	return conversion.NewService(mockWalletRepo), mockWalletRepo
}
