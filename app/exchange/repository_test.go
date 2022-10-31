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
	query := regexp.QuoteMeta(`SELECT EXISTS ( SELECT 1 FROM currencies c WHERE c.code = $1)`)
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
		FROM exchanges e WHERE currency_from = $1 AND currency_to = $2`)
	t.Run("given existing currencies then it should return exchange", func(t *testing.T) {
		givenCurrencyFrom := "TRY"
		givenCurrencyTo := "USD"
		expectedExchange := entity.Exchange{
			FromCurrency: givenCurrencyFrom,
			ToCurrency:   givenCurrencyTo,
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

		sqlmock.ExpectQuery(query).WithArgs(givenCurrencyFrom, givenCurrencyTo).
			WillReturnRows(expectedRows)

		exchange, err := repository.GetExchangeValuesForGivenCurrencies(context.Background(),
			givenCurrencyFrom, givenCurrencyTo)
		assert.Nil(t, err)
		assert.Nil(t, sqlmock.ExpectationsWereMet())
		assert.Equal(t, expectedExchange, exchange)
	})
	t.Run("given not existing currencies then it should return error", func(t *testing.T) {
		givenCurrencyFrom := "vxcxcg"
		givenCurrencyTo := "sfasf"
		sqlmock.ExpectQuery(query).WithArgs(givenCurrencyFrom, givenCurrencyTo).
			WillReturnError(errors.New("exchange not found"))
		exchange, err := repository.GetExchangeValuesForGivenCurrencies(context.Background(),
			givenCurrencyFrom, givenCurrencyTo)
		assert.NotNil(t, err)
		assert.Empty(t, exchange)
	})
}
