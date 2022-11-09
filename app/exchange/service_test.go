package exchange_test

import (
	"context"
	"errors"
	"github.com/eneskzlcn/currency-conversion-service/app/entity"
	"github.com/eneskzlcn/currency-conversion-service/app/exchange"
	mocks "github.com/eneskzlcn/currency-conversion-service/app/mocks/exchange"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
	"testing"
	"time"
)

func TestService_CreateExchangeRate(t *testing.T) {
	mockExchangeRepository := mocks.NewMockRepository(gomock.NewController(t))
	service := exchange.NewService(mockExchangeRepository, zap.S())

	t.Run("given not supported currency from then it should return not valid currency error", func(t *testing.T) {
		givenRequest := exchange.ExchangeRateRequest{
			FromCurrency: "",
			ToCurrency:   "TRY",
		}
		userID := 1
		mockExchangeRepository.EXPECT().IsCurrencyExists(gomock.Any(), givenRequest.FromCurrency).
			Return(false, nil)
		exchangeDateResp, err := service.PrepareExchangeRateOffer(context.Background(), userID, givenRequest)
		assert.Equal(t, err.Error(), exchange.NotValidCurrency)
		assert.Empty(t, exchangeDateResp)
	})
	t.Run("given not supported currency to then it should return not valid currency error", func(t *testing.T) {
		givenRequest := exchange.ExchangeRateRequest{
			FromCurrency: "TRY",
			ToCurrency:   "",
		}
		userID := 1
		mockExchangeRepository.EXPECT().IsCurrencyExists(gomock.Any(), givenRequest.FromCurrency).
			Return(true, nil)
		mockExchangeRepository.EXPECT().IsCurrencyExists(gomock.Any(), givenRequest.ToCurrency).
			Return(false, nil)
		exchangeDateResp, err := service.PrepareExchangeRateOffer(context.Background(), userID, givenRequest)
		assert.Equal(t, err.Error(), exchange.NotValidCurrency)
		assert.Empty(t, exchangeDateResp)
	})
	t.Run("given supported currency values but exchange values can not be taken then it should return error", func(t *testing.T) {
		givenRequest := exchange.ExchangeRateRequest{
			FromCurrency: "TRY",
			ToCurrency:   "USD",
		}
		mockExchangeRepository.EXPECT().IsCurrencyExists(gomock.Any(), givenRequest.FromCurrency).
			Return(true, nil)
		mockExchangeRepository.EXPECT().IsCurrencyExists(gomock.Any(), givenRequest.ToCurrency).
			Return(true, nil)

		mockExchangeRepository.EXPECT().
			GetExchangeValuesForGivenCurrencies(gomock.Any(),
				givenRequest.FromCurrency, givenRequest.ToCurrency).
			Return(entity.Exchange{}, errors.New("exchange not found"))

		exchangeDateResp, err := service.PrepareExchangeRateOffer(context.Background(), 1, givenRequest)
		assert.NotNil(t, err)
		assert.Empty(t, exchangeDateResp)
	})
	t.Run("given supported currency values and exchange values taken then it should return exchange rate response", func(t *testing.T) {
		givenRequest := exchange.ExchangeRateRequest{
			FromCurrency: "TRY",
			ToCurrency:   "USD",
		}
		givenExchange := entity.Exchange{
			FromCurrency: givenRequest.FromCurrency,
			ToCurrency:   givenRequest.ToCurrency,
			ExchangeRate: 2.3,
			MarkupRate:   1.1,
			CreatedAt:    time.Now(),
			UpdatedAt:    time.Now(),
		}
		userID := 1
		expectedResponse := exchange.ExchangeRateResponse{
			FromCurrency: givenRequest.FromCurrency,
			ToCurrency:   givenRequest.ToCurrency,
			ExchangeRate: givenExchange.ExchangeRate + givenExchange.MarkupRate,
			CreatedAt:    givenExchange.CreatedAt,
			ExpiresAt:    givenExchange.UpdatedAt.Unix(),
		}

		mockExchangeRepository.EXPECT().IsCurrencyExists(gomock.Any(), givenRequest.FromCurrency).
			Return(true, nil)
		mockExchangeRepository.EXPECT().IsCurrencyExists(gomock.Any(), givenRequest.ToCurrency).
			Return(true, nil)

		mockExchangeRepository.EXPECT().
			GetExchangeValuesForGivenCurrencies(gomock.Any(),
				givenRequest.FromCurrency, givenRequest.ToCurrency).
			Return(givenExchange, nil)

		mockExchangeRepository.EXPECT().
			SetUserActiveExchangeRateOffer(gomock.Any(), gomock.Any()).Return(true, nil)
		exchangeRateResp, err := service.PrepareExchangeRateOffer(context.Background(), userID, givenRequest)
		assert.Nil(t, err)
		assert.Equal(t, expectedResponse.ExchangeRate, exchangeRateResp.ExchangeRate)
	})
}
