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
func (r *Repository) GetExchangeValuesForGivenCurrencies(ctx context.Context, fromCurrency, toCurrency string) (entity.Exchange, error) {
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
