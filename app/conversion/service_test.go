//go:build unit

package conversion_test

import (
	"context"
	"errors"
	"github.com/eneskzlcn/currency-conversion-service/app/conversion"
	"github.com/eneskzlcn/currency-conversion-service/app/entity"
	mocks "github.com/eneskzlcn/currency-conversion-service/app/mocks/conversion"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
	"testing"
	"time"
)

func TestService_CreateCurrencyConversion(t *testing.T) {
	service, mockWalletService, mockConversionRepo := createServiceWithMockWalletServiceAndConversionRepository(t)
	t.Run("given offer not same with saved user exchange rate offer then it should return false with error", func(t *testing.T) {
		givenConversionOfferReq := conversion.CurrencyConversionOfferRequest{
			FromCurrency: "USD",
			ToCurrency:   "TRY",
			ExchangeRate: 2.30,
			CreatedAt:    time.Now(),
			ExpiresAt:    time.Now().Add(-10 * time.Minute).Unix(),
			Balance:      200,
		}
		userID := 1
		givenUserActiveExchangeOfferDTO := conversion.UserActiveExchangeOfferDTO{
			UserID:       userID,
			FromCurrency: givenConversionOfferReq.FromCurrency,
			ToCurrency:   givenConversionOfferReq.ToCurrency,
		}

		mockConversionRepo.EXPECT().GetUserActiveExchangeOffer(gomock.Any(), givenUserActiveExchangeOfferDTO).
			Return(entity.UserActiveExchangeOffer{}, errors.New(""))
		success, err := service.ConvertCurrencies(context.TODO(), userID, givenConversionOfferReq)
		assert.NotNil(t, err)
		assert.False(t, success)
	})
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
		expectedUserActiveExchangeOffer := entity.UserActiveExchangeOffer{
			UserID:         userID,
			FromCurrency:   givenConversionOfferReq.FromCurrency,
			ToCurrency:     givenConversionOfferReq.ToCurrency,
			ExchangeRate:   givenConversionOfferReq.ExchangeRate,
			OfferCreatedAt: givenConversionOfferReq.CreatedAt,
			OfferExpiresAt: givenConversionOfferReq.ExpiresAt,
		}
		mockConversionRepo.EXPECT().GetUserActiveExchangeOffer(gomock.Any(), gomock.Any()).
			Return(expectedUserActiveExchangeOffer, nil)
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
		expectedUserActiveExchangeOffer := entity.UserActiveExchangeOffer{
			UserID:         userID,
			FromCurrency:   givenConversionOfferReq.FromCurrency,
			ToCurrency:     givenConversionOfferReq.ToCurrency,
			ExchangeRate:   givenConversionOfferReq.ExchangeRate,
			OfferCreatedAt: givenConversionOfferReq.CreatedAt,
			OfferExpiresAt: givenConversionOfferReq.ExpiresAt,
		}
		mockConversionRepo.EXPECT().GetUserActiveExchangeOffer(gomock.Any(), gomock.Any()).
			Return(expectedUserActiveExchangeOffer, nil)
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
		expectedUserActiveExchangeOffer := entity.UserActiveExchangeOffer{
			UserID:         userID,
			FromCurrency:   givenConversionOfferReq.FromCurrency,
			ToCurrency:     givenConversionOfferReq.ToCurrency,
			ExchangeRate:   givenConversionOfferReq.ExchangeRate,
			OfferCreatedAt: givenConversionOfferReq.CreatedAt,
			OfferExpiresAt: givenConversionOfferReq.ExpiresAt,
		}
		mockConversionRepo.EXPECT().GetUserActiveExchangeOffer(gomock.Any(), gomock.Any()).
			Return(expectedUserActiveExchangeOffer, nil)
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
func createServiceWithMockWalletServiceAndConversionRepository(t *testing.T) (*conversion.Service, *mocks.MockWalletService, *mocks.MockConversionRepository) {
	ctrl := gomock.NewController(t)
	mockWalletService := mocks.NewMockWalletService(ctrl)
	mockConversionRepo := mocks.NewMockConversionRepository(ctrl)
	return conversion.NewService(mockWalletService, zap.S(), mockConversionRepo), mockWalletService, mockConversionRepo
}
