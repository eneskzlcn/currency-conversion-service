package exchange

import (
	"context"
	"database/sql"
	"github.com/eneskzlcn/currency-conversion-service/app/entity"
	"go.uber.org/zap"
)

type Repository struct {
	db     *sql.DB
	logger *zap.SugaredLogger
}

func NewRepository(db *sql.DB, logger *zap.SugaredLogger) *Repository {
	return &Repository{db: db, logger: logger}
}
func (r *Repository) IsCurrencyExists(ctx context.Context, currency string) (bool, error) {
	query := `SELECT EXISTS ( SELECT 1 FROM currencies WHERE code = $1)`
	row := r.db.QueryRowContext(ctx, query, currency)
	var isExists bool
	err := row.Scan(&isExists)
	return isExists, err
}
func (r *Repository) GetCurrencyExchangeValuesByCurrency(ctx context.Context, dto ExchangeCurrencyDTO) (entity.CurrencyExchangeValues, error) {
	query := `
		SELECT currency_from, currency_to, exchange_rate, markup_rate, created_at, updated_at
		FROM currency_exchange_values e WHERE currency_from = $1 AND currency_to = $2`

	row := r.db.QueryRowContext(ctx, query, dto.FromCurrency, dto.ToCurrency)
	var exchange entity.CurrencyExchangeValues
	err := row.Scan(&exchange.FromCurrency,
		&exchange.ToCurrency,
		&exchange.ExchangeRate,
		&exchange.MarkupRate,
		&exchange.CreatedAt,
		&exchange.UpdatedAt,
	)
	if err != nil {
		return entity.CurrencyExchangeValues{}, err
	}
	return exchange, nil
}
func (r *Repository) CreateExchangeRateOffer(ctx context.Context, dto CreateExchangeRateOfferDTO) (int, error) {
	query := `
	INSERT INTO exchange_rate_offers(user_id, 
	currency_from, currency_to, exchange_rate, offer_created_at, offer_expires_at)
	VALUES ($1, $2, $3, $4, $5, $6)
	RETURNING id;
	`
	row := r.db.QueryRowContext(ctx, query, dto.UserID,
		dto.FromCurrency, dto.ToCurrency, dto.ExchangeRate, dto.OfferCreatedAt, dto.OfferExpiresAt)

	var exchangeRateOfferID int
	if err := row.Scan(&exchangeRateOfferID); err != nil {
		return -1, err
	}
	return exchangeRateOfferID, nil
}
