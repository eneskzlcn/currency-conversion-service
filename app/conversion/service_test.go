//go:build unit

package conversion_test

import (
	"context"
	"github.com/eneskzlcn/currency-conversion-service/app/conversion"
	mocks "github.com/eneskzlcn/currency-conversion-service/app/mocks/conversion"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
	"testing"
	"time"
)

func TestService_CreateCurrencyConversion(t *testing.T) {
	service, mockWalletService := createServiceWithMockConversionRepository(t)
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
		assert.Equal(t, err.Error(), conversion.CurrencyConversionOfferExpired)
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

		mockWalletService.EXPECT().
			GetUserBalanceOnGivenCurrency(gomock.Any(), userID, givenConversionOfferReq.FromCurrency).
			Return(float32(100), nil)
		success, err := service.ConvertCurrencies(context.TODO(), userID, givenConversionOfferReq)
		assert.False(t, success)
		assert.NotNil(t, err)
		assert.Equal(t, err.Error(), conversion.NotEnoughBalanceForConversionOffer)
	})
	t.Run("given valid currency conversion offer then it should return true", func(t *testing.T) {
		givenConversionOfferReq := conversion.CurrencyConversionOfferRequest{
			FromCurrency: "USD",
			ToCurrency:   "TRY",
			ExchangeRate: 2.30,
			CreatedAt:    time.Now(),
			ExpiresAt:    time.Now().Add(10 * time.Minute).Unix(),
			Balance:      float32(200),
		}
		userID := 2
		targetCurrencyBalanceAdjustAmount := givenConversionOfferReq.ExchangeRate * givenConversionOfferReq.Balance
		mockWalletService.EXPECT().
			GetUserBalanceOnGivenCurrency(gomock.Any(), userID, givenConversionOfferReq.FromCurrency).
			Return(float32(500), nil)

		mockWalletService.EXPECT().
			AdjustUserBalanceOnGivenCurrency(gomock.Any(), userID,
				givenConversionOfferReq.FromCurrency, -1*givenConversionOfferReq.Balance).
			Return(true, nil)

		mockWalletService.EXPECT().
			AdjustUserBalanceOnGivenCurrency(gomock.Any(), userID,
				givenConversionOfferReq.ToCurrency, targetCurrencyBalanceAdjustAmount).
			Return(true, nil)

		success, err := service.ConvertCurrencies(context.TODO(), userID, givenConversionOfferReq)
		assert.Nil(t, err)
		assert.True(t, success)
	})
}
func createServiceWithMockConversionRepository(t *testing.T) (*conversion.Service, *mocks.MockWalletService) {
	ctrl := gomock.NewController(t)
	mockWalletRepo := mocks.NewMockWalletService(ctrl)
	return conversion.NewService(mockWalletRepo, zap.S()), mockWalletRepo
}
