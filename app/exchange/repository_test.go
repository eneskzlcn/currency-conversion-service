//go:build unit

package exchange_test

import (
	"context"
	"errors"
	"github.com/eneskzlcn/currency-conversion-service/app/entity"
	"github.com/eneskzlcn/currency-conversion-service/app/exchange"
	"github.com/eneskzlcn/currency-conversion-service/postgres"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
	"regexp"
	"testing"
	"time"
)

func TestRepository_IsCurrencyExists(t *testing.T) {
	db, sqlmock := postgres.NewMockPostgres()
	repository := exchange.NewRepository(db, zap.S())
	query := regexp.QuoteMeta(`SELECT EXISTS ( SELECT 1 FROM currencies WHERE code = $1)`)
	t.Run("given not existing currency code then it should return false", func(t *testing.T) {
		givenCurrency := "asf"
		sqlmock.ExpectQuery(query).WithArgs(givenCurrency).
			WillReturnRows(sqlmock.NewRows([]string{"exists"}).AddRow(false))
		exists, err := repository.IsCurrencyExists(context.Background(), givenCurrency)
		assert.Nil(t, err)
		assert.Nil(t, sqlmock.ExpectationsWereMet())
		assert.False(t, exists)
	})
	t.Run("given existing currency code then it should return true", func(t *testing.T) {
		givenCurrency := "TRY"
		sqlmock.ExpectQuery(query).WithArgs(givenCurrency).
			WillReturnRows(sqlmock.NewRows([]string{"exists"}).AddRow(true))
		exists, err := repository.IsCurrencyExists(context.Background(), givenCurrency)
		assert.Nil(t, err)
		assert.Nil(t, sqlmock.ExpectationsWereMet())
		assert.True(t, exists)
	})
}
func TestRepository_GetExchangeValuesForGivenCurrencies(t *testing.T) {
	db, sqlmock := postgres.NewMockPostgres()
	repository := exchange.NewRepository(db, zap.S())
	query := regexp.QuoteMeta(`
		SELECT currency_from, currency_to, exchange_rate, markup_rate, created_at, updated_at
		FROM currency_exchange_values e WHERE currency_from = $1 AND currency_to = $2`)
	t.Run("given existing currencies then it should return exchange", func(t *testing.T) {
		givenCurrencyExchangeDTO := exchange.ExchangeCurrencyDTO{
			FromCurrency: "TRY",
			ToCurrency:   "USD",
		}
		expectedExchange := entity.CurrencyExchangeValues{
			FromCurrency: givenCurrencyExchangeDTO.FromCurrency,
			ToCurrency:   givenCurrencyExchangeDTO.ToCurrency,
			ExchangeRate: 12.3,
			MarkupRate:   3.2,
			CreatedAt:    time.Now(),
			UpdatedAt:    time.Now(),
		}
		expectedRows := sqlmock.NewRows([]string{"currency_from",
			"currency_to", "exchange_rate", "markup_rate", "created_at", "updated_at"})
		expectedRows.AddRow(expectedExchange.FromCurrency,
			expectedExchange.ToCurrency, expectedExchange.ExchangeRate,
			expectedExchange.MarkupRate, expectedExchange.CreatedAt, expectedExchange.UpdatedAt)

		sqlmock.ExpectQuery(query).WithArgs(givenCurrencyExchangeDTO.FromCurrency,
			givenCurrencyExchangeDTO.ToCurrency).
			WillReturnRows(expectedRows)

		exchange, err := repository.GetCurrencyExchangeValuesByCurrency(context.Background(),
			givenCurrencyExchangeDTO)
		assert.Nil(t, err)
		assert.Nil(t, sqlmock.ExpectationsWereMet())
		assert.Equal(t, expectedExchange, exchange)
	})
	t.Run("given not existing currencies then it should return error", func(t *testing.T) {
		givenCurrencyExchangeDTO := exchange.ExchangeCurrencyDTO{
			FromCurrency: "TRY",
			ToCurrency:   "USD",
		}
		sqlmock.ExpectQuery(query).WithArgs(givenCurrencyExchangeDTO.FromCurrency,
			givenCurrencyExchangeDTO.ToCurrency).
			WillReturnError(errors.New("exchange not found"))
		exchange, err := repository.GetCurrencyExchangeValuesByCurrency(context.Background(),
			givenCurrencyExchangeDTO)
		assert.NotNil(t, err)
		assert.Empty(t, exchange)
	})
}

func TestRepository_CreateExchangeRateOffer(t *testing.T) {
	db, sqlmock := postgres.NewMockPostgres()
	repository := exchange.NewRepository(db, zap.S())
	query := regexp.QuoteMeta(`INSERT INTO exchange_rate_offers(user_id, 
	currency_from, currency_to, exchange_rate, offer_created_at, offer_expires_at)
	VALUES ($1, $2, $3, $4, $5, $6) 
	RETURNING id;`)

	t.Run("given existing user, and currencies then it should create if not exists or update active exchange offer", func(t *testing.T) {
		offerDTO := exchange.CreateExchangeRateOfferDTO{
			UserID:         2,
			FromCurrency:   "TRY",
			ToCurrency:     "USD",
			ExchangeRate:   2.2,
			OfferCreatedAt: time.Now(),
			OfferExpiresAt: 12312412414,
		}
		expectedExchangeRateOfferID := 2
		sqlmock.ExpectQuery(query).WillReturnRows(sqlmock.NewRows([]string{"id"}).
			AddRow(expectedExchangeRateOfferID))
		exchangeRateOfferID, err := repository.CreateExchangeRateOffer(context.TODO(), offerDTO)
		assert.Nil(t, err)
		assert.Equal(t, exchangeRateOfferID, expectedExchangeRateOfferID)
	})
	t.Run("given not existing user or currencies then it should return false with error", func(t *testing.T) {
		offer := exchange.CreateExchangeRateOfferDTO{
			UserID:         -1,
			FromCurrency:   "",
			ToCurrency:     "",
			ExchangeRate:   2.2,
			OfferCreatedAt: time.Now(),
			OfferExpiresAt: 12312412414,
		}
		expectedExchangeRateOfferID := -1
		sqlmock.ExpectQuery(query).WillReturnError(errors.New("not existing user or currency"))
		exchangeRateOfferID, err := repository.CreateExchangeRateOffer(context.TODO(), offer)
		assert.NotNil(t, err)
		assert.Equal(t, exchangeRateOfferID, expectedExchangeRateOfferID)
	})
}
