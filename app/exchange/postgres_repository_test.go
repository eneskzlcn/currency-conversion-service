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
		FROM currency_exchange_values WHERE currency_from = $1 AND currency_to = $2`)
	t.Run("given existing currencies then it should return exchange", func(t *testing.T) {
		givenCurrencyFrom := "TRY"
		givenCurrencyTo := "USD"
		expectedExchange := entity.CurrencyExchangeValues{
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

func TestRepository_SetUserActiveExchangeRateOffer(t *testing.T) {
	db, sqlmock := postgres.NewMockPostgres()
	repository := exchange.NewRepository(db, zap.S())
	query := regexp.QuoteMeta(`INSERT INTO user_exchange_offers(user_id, 
	currency_from, currency_to, exchange_rate, offer_created_at, offer_expires_at)
	VALUES ($1, $2, $3, $4, $5, $6) 
	RETURNING id;`)

	t.Run("given existing user, and currencies then it should create if not exists or update active exchange offer", func(t *testing.T) {
		givenCreateExchangeRateOfferDTO := exchange.NewCreateExchangeRateOfferDTO(2,
			"TRY", "USD", 2.2, time.Now(), 12312412414)
		expectedCreatedOfferID := 2
		sqlmock.ExpectQuery(query).WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(expectedCreatedOfferID))
		createdID, err := repository.CreateExchangeRateOffer(context.TODO(), givenCreateExchangeRateOfferDTO)
		assert.Nil(t, err)
		assert.Equal(t, expectedCreatedOfferID, createdID)
	})
	t.Run("given not existing user or currencies then it should return false with error", func(t *testing.T) {
		givenCreateExchangeRateOfferDTO := exchange.NewCreateExchangeRateOfferDTO(-1,
			"", "", 2.2, time.Now(), 12312412414)
		sqlmock.ExpectQuery(query).WillReturnError(errors.New("not existing user or currency"))
		expectedCreatedOfferID := -1
		createdExchangeRateOfferID, err := repository.CreateExchangeRateOffer(context.TODO(), givenCreateExchangeRateOfferDTO)
		assert.NotNil(t, err)
		assert.Equal(t, expectedCreatedOfferID, createdExchangeRateOfferID)
	})
}
