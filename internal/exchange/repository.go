package exchange

import (
	"context"
	"database/sql"
	"github.com/eneskzlcn/currency-conversion-service/internal/entity"
)

type Repository struct {
	db *sql.DB
}

func NewRepository(db *sql.DB) *Repository {
	if db == nil {
		return nil
	}
	return &Repository{db: db}
}
func (r *Repository) IsCurrencyExists(ctx context.Context, currency string) (bool, error) {
	query := `SELECT EXISTS ( SELECT 1 FROM currencies c WHERE c.code = $1)`
	row := r.db.QueryRowContext(ctx, query, currency)
	var isExists bool
	if err := row.Scan(&isExists); err != nil {
		return false, err
	}
	return isExists, nil
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
