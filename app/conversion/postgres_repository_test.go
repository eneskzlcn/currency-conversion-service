package conversion_test

import (
	"context"
	"errors"
	"github.com/eneskzlcn/currency-conversion-service/app/conversion"
	"github.com/eneskzlcn/currency-conversion-service/app/model"
	"github.com/eneskzlcn/currency-conversion-service/postgres"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
	"regexp"
	"testing"
	"time"
)

func TestPostgresRepository_GetExchangeOfferByID(t *testing.T) {
	db, sqlmock := postgres.NewMockPostgres()
	repository := conversion.NewPostgresRepository(db, zap.S())
	query := regexp.QuoteMeta(`
		SELECT id, user_id, currency_from, currency_to, exchange_rate, offer_created_at, offer_expires_at
		FROM user_exchange_offers
		WHERE id = $1`)

	t.Run("given not existing exchange rate offer id then it should return error", func(t *testing.T) {
		givenExchangeRateOfferID := 2
		givenGetExchangeRateOfferDTO := conversion.NewGetExchangeRateOfferDTO(givenExchangeRateOfferID)
		sqlmock.ExpectQuery(query).WithArgs(givenExchangeRateOfferID).WillReturnError(errors.New(""))

		exchangeRateOffer, err := repository.GetExchangeOfferByID(context.TODO(),
			givenGetExchangeRateOfferDTO)

		assert.Nil(t, sqlmock.ExpectationsWereMet())
		assert.NotNil(t, err)
		assert.Empty(t, exchangeRateOffer)
	})
	t.Run("given existing user id, currency from and currency to query arguments then it should return exchange rate offer", func(t *testing.T) {
		givenExchangeRateOfferID := 2
		givenExchangeRateOfferDTO := conversion.NewGetExchangeRateOfferDTO(givenExchangeRateOfferID)
		expectedOffer := model.ExchangeRateOffer{
			ID:             givenExchangeRateOfferID,
			UserID:         2,
			FromCurrency:   "TRY",
			ToCurrency:     "USD",
			ExchangeRate:   2.3,
			OfferCreatedAt: time.Now().UTC(),
			OfferExpiresAt: time.Now().Add(123123 * time.Minute).Unix(),
		}

		expectedRows := sqlmock.NewRows([]string{"id", "user_id", "currency_from",
			"currency_to", "exchange_rate", "offer_created_at", "offer_expires_at"}).
			AddRow(givenExchangeRateOfferID, expectedOffer.UserID, expectedOffer.FromCurrency,
				expectedOffer.ToCurrency, expectedOffer.ExchangeRate,
				expectedOffer.OfferCreatedAt, expectedOffer.OfferExpiresAt)

		sqlmock.ExpectQuery(query).WithArgs(givenExchangeRateOfferID).WillReturnRows(expectedRows)

		exchangeRateOffer, err := repository.GetExchangeOfferByID(context.TODO(),
			givenExchangeRateOfferDTO)

		assert.Nil(t, sqlmock.ExpectationsWereMet())
		assert.Nil(t, err)
		assert.Equal(t, expectedOffer, exchangeRateOffer)
	})
}
func TestPostgresRepository_CreateUserConversion(t *testing.T) {
	db, sqlmock := postgres.NewMockPostgres()
	repository := conversion.NewPostgresRepository(db, zap.S())
	query := regexp.QuoteMeta(`
		INSERT INTO user_currency_conversions (user_id, currency_from, currency_to, 
		sender_balance_dec_amount, receiver_balance_inc_amount)
		VALUES($1, $2, $3, $4, $5)
		RETURNING created_at, id;`)

	t.Run("given create conversion dto then it should return created user currency conversion", func(t *testing.T) {
		givenCreateConversionDTO := conversion.NewCreateUserConversionDTO(2, "TRY",
			"USD", 23.4, 155.5)
		expectedCreatedAt := time.Now().Local()
		expectedConversionId := 1
		expectedRows := sqlmock.NewRows([]string{"created_at", "id"}).AddRow(expectedCreatedAt, expectedConversionId)

		sqlmock.ExpectQuery(query).WithArgs(givenCreateConversionDTO.UserID(),
			givenCreateConversionDTO.FromCurrency(), givenCreateConversionDTO.ToCurrency(),
			givenCreateConversionDTO.SenderBalanceDecAmount(), givenCreateConversionDTO.ReceiverBalanceIncAmount()).
			WillReturnRows(expectedRows)

		expectedCurrencyConversion := model.UserCurrencyConversion{
			ID:                       expectedConversionId,
			UserID:                   givenCreateConversionDTO.UserID(),
			FromCurrency:             givenCreateConversionDTO.FromCurrency(),
			ToCurrency:               givenCreateConversionDTO.ToCurrency(),
			SenderBalanceDecAmount:   givenCreateConversionDTO.SenderBalanceDecAmount(),
			ReceiverBalanceIncAmount: givenCreateConversionDTO.ReceiverBalanceIncAmount(),
			CreatedAt:                expectedCreatedAt,
		}
		createdCurrencyConversion, err := repository.CreateUserConversion(context.TODO(), givenCreateConversionDTO)
		assert.Nil(t, err)
		assert.Equal(t, expectedCurrencyConversion, createdCurrencyConversion)
	})
}
