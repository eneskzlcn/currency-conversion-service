package exchange_test

import (
	"context"
	"errors"
	"github.com/eneskzlcn/currency-conversion-service/internal/entity"
	"github.com/eneskzlcn/currency-conversion-service/internal/exchange"
	mocks "github.com/eneskzlcn/currency-conversion-service/internal/mocks/exchange"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestNewService(t *testing.T) {
	t.Run("given empty repository then it should return nil", func(t *testing.T) {
		service := exchange.NewService(nil)
		assert.Nil(t, service)
	})
	t.Run("given empty repository then it should return nil", func(t *testing.T) {
		mockExchangeRepository := mocks.NewMockExchangeRepository(gomock.NewController(t))
		service := exchange.NewService(mockExchangeRepository)
		assert.NotNil(t, service)
	})
}

func TestService_CreateExchangeRate(t *testing.T) {
	mockExchangeRepository := mocks.NewMockExchangeRepository(gomock.NewController(t))
	service := exchange.NewService(mockExchangeRepository)

	t.Run("given not supported currency from then it should return not valid currency error", func(t *testing.T) {
		givenRequest := exchange.ExchangeRateRequest{
			FromCurrency: "",
			ToCurrency:   "TRY",
		}
		mockExchangeRepository.EXPECT().IsCurrencyExists(gomock.Any(), givenRequest.FromCurrency).
			Return(false)
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
			Return(true)
		mockExchangeRepository.EXPECT().IsCurrencyExists(gomock.Any(), givenRequest.ToCurrency).
			Return(false)
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
			Return(true)
		mockExchangeRepository.EXPECT().IsCurrencyExists(gomock.Any(), givenRequest.ToCurrency).
			Return(true)

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
			ExpiresAt:    givenExchange.UpdatedAt,
		}
		mockExchangeRepository.EXPECT().IsCurrencyExists(gomock.Any(), givenRequest.FromCurrency).
			Return(true)
		mockExchangeRepository.EXPECT().IsCurrencyExists(gomock.Any(), givenRequest.ToCurrency).
			Return(true)

		mockExchangeRepository.EXPECT().
			GetExchangeValuesForGivenCurrencies(gomock.Any(),
				givenRequest.FromCurrency, givenRequest.ToCurrency).
			Return(givenExchange, nil)

		exchangeRateResp, err := service.CreateExchangeRate(context.Background(), givenRequest)
		assert.Nil(t, err)
		assert.Equal(t, expectedResponse.ExchangeRate, exchangeRateResp.ExchangeRate)
	})
}
