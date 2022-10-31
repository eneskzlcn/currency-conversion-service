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
	mockExchangeRepository := mocks.NewMockExchangeRepository(gomock.NewController(t))
	service := exchange.NewService(mockExchangeRepository, zap.S())

	t.Run("given not supported currency from then it should return not valid currency error", func(t *testing.T) {
		givenRequest := exchange.ExchangeRateRequest{
			FromCurrency: "",
			ToCurrency:   "TRY",
		}
		mockExchangeRepository.EXPECT().IsCurrencyExists(gomock.Any(), givenRequest.FromCurrency).
			Return(false, nil)
		exchangeDateResp, err := service.CreateExchangeRate(context.Background(), givenRequest)
		assert.True(t, errors.Is(err, exchange.NotValidCurrencyErr))
		assert.Empty(t, exchangeDateResp)
	})
	t.Run("given not supported currency to then it should return not valid currency error", func(t *testing.T) {
		givenRequest := exchange.ExchangeRateRequest{
			FromCurrency: "TRY",
			ToCurrency:   "",
		}
		mockExchangeRepository.EXPECT().IsCurrencyExists(gomock.Any(), givenRequest.FromCurrency).
			Return(true, nil)
		mockExchangeRepository.EXPECT().IsCurrencyExists(gomock.Any(), givenRequest.ToCurrency).
			Return(false, nil)
		exchangeDateResp, err := service.CreateExchangeRate(context.Background(), givenRequest)
		assert.True(t, errors.Is(err, exchange.NotValidCurrencyErr))
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

		exchangeDateResp, err := service.CreateExchangeRate(context.Background(), givenRequest)
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

		exchangeRateResp, err := service.CreateExchangeRate(context.Background(), givenRequest)
		assert.Nil(t, err)
		assert.Equal(t, expectedResponse.ExchangeRate, exchangeRateResp.ExchangeRate)
	})
}
