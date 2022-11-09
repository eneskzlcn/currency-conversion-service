package exchange

import (
	"context"
	"database/sql"
	"github.com/eneskzlcn/currency-conversion-service/app/entity"
	"go.uber.org/zap"
)

type postgresRepository struct {
	db     *sql.DB
	logger *zap.SugaredLogger
}

func NewRepository(db *sql.DB, logger *zap.SugaredLogger) *postgresRepository {
	return &postgresRepository{db: db, logger: logger}
}
func (r *postgresRepository) IsCurrencyExists(ctx context.Context, currency string) (bool, error) {
	query := `SELECT EXISTS ( SELECT 1 FROM currencies WHERE code = $1)`
	row := r.db.QueryRowContext(ctx, query, currency)
	var isExists bool
	err := row.Scan(&isExists)
	return isExists, err
}
func (r *postgresRepository) GetExchangeValuesForGivenCurrencies(ctx context.Context,
	fromCurrency, toCurrency string) (entity.Exchange, error) {
	query := `
		SELECT currency_from, currency_to, exchange_rate, markup_rate, created_at, updated_at
		FROM exchanges e WHERE currency_from = $1 AND currency_to = $2`

	row := r.db.QueryRowContext(ctx, query, fromCurrency, toCurrency)
	var exchange entity.Exchange
	err := row.Scan(&exchange.FromCurrency,
		&exchange.ToCurrency,
		&exchange.ExchangeRate,
		&exchange.MarkupRate,
		&exchange.CreatedAt,
		&exchange.UpdatedAt,
	)
	if err != nil {
		return entity.Exchange{}, err
	}
	return exchange, nil
}
func (r *postgresRepository) SetUserActiveExchangeRateOffer(ctx context.Context, offer entity.UserActiveExchangeOffer) (bool, error) {
	query := `
	INSERT INTO user_active_exchange_offers(user_id, 
	currency_from, currency_to, exchange_rate, offer_created_at, offer_expires_at)
	VALUES ($1, $2, $3, $4, $5, $6) 
	ON CONFLICT(user_id, currency_from, currency_To) DO
	UPDATE SET exchange_rate = $4, offer_created_at = $5, offer_expires_at = $6
	`
	row := r.db.QueryRowContext(ctx, query, offer.UserID,
		offer.FromCurrency, offer.ToCurrency, offer.ExchangeRate, offer.OfferCreatedAt, offer.OfferExpiresAt)
	if err := row.Err(); err != nil {
		return false, err
	}
	return true, nil
}
