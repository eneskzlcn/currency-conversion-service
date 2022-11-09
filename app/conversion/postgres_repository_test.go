package conversion_test

import (
	"context"
	"github.com/eneskzlcn/currency-conversion-service/app/conversion"
	"github.com/eneskzlcn/currency-conversion-service/app/entity"
	"github.com/eneskzlcn/currency-conversion-service/postgres"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
	"regexp"
	"testing"
	"time"
)

func TestPostgresRepository_GetUserActiveExchangeOffer(t *testing.T) {
	db, sqlmock := postgres.NewMockPostgres()
	repository := conversion.NewRepository(db, zap.S())
	query := regexp.QuoteMeta(`
		SELECT user_id, currency_from, currency_to, exchange_rate, offer_created_at, offer_expires_at
		FROM user_active_exchange_offers
		WHERE user_id = $1 AND currency_from = $2 AND currency_to = $3`)

	t.Run("given not existing user id, currency from and currency to query arguments then it should return error", func(t *testing.T) {
		expectedRows := sqlmock.NewRows([]string{"user_id", "currency_from",
			"currency_to", "exchange_rate", "offer_created_at", "offer_expires_at"})
		sqlmock.ExpectQuery(query).WillReturnRows(expectedRows)
		givenUserActiveExchangeDTO := conversion.UserActiveExchangeOfferDTO{
			UserID:       2,
			FromCurrency: "Tx",
			ToCurrency:   "yr",
		}
		userActiveExchangeOffer, err := repository.GetUserActiveExchangeOffer(context.TODO(),
			givenUserActiveExchangeDTO)
		assert.NotNil(t, err)
		assert.Empty(t, userActiveExchangeOffer)
	})
	t.Run("given existing user id, currency from and currency to query arguments then it should return user active exchange offer", func(t *testing.T) {
		givenUserActiveExchangeDTO := conversion.UserActiveExchangeOfferDTO{
			UserID:       2,
			FromCurrency: "TRY",
			ToCurrency:   "USD",
		}
		expectedOffer := entity.UserActiveExchangeOffer{
			UserID:         givenUserActiveExchangeDTO.UserID,
			FromCurrency:   givenUserActiveExchangeDTO.FromCurrency,
			ToCurrency:     givenUserActiveExchangeDTO.ToCurrency,
			ExchangeRate:   2.3,
			OfferCreatedAt: time.Now(),
			OfferExpiresAt: time.Now().Add(123123 * time.Minute).Unix(),
		}

		expectedRows := sqlmock.NewRows([]string{"user_id", "currency_from",
			"currency_to", "exchange_rate", "offer_created_at", "offer_expires_at"}).
			AddRow(expectedOffer.UserID, expectedOffer.FromCurrency,
				expectedOffer.ToCurrency, expectedOffer.ExchangeRate, expectedOffer.OfferCreatedAt, expectedOffer.OfferExpiresAt)
		sqlmock.ExpectQuery(query).WillReturnRows(expectedRows)

		userActiveExchangeOffer, err := repository.GetUserActiveExchangeOffer(context.TODO(),
			givenUserActiveExchangeDTO)
		assert.Nil(t, err)
		assert.Equal(t, userActiveExchangeOffer, expectedOffer)
	})
}
