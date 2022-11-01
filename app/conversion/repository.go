package conversion

import (
	"context"
	"database/sql"
	"github.com/eneskzlcn/currency-conversion-service/app/entity"
	"go.uber.org/zap"
)

type Repository struct {
	logger *zap.SugaredLogger
	db     *sql.DB
}

func NewRepository(db *sql.DB, logger *zap.SugaredLogger) *Repository {
	return &Repository{db: db, logger: logger}
}

func (r *Repository) GetUserActiveExchangeOffer(ctx context.Context,
	dto UserActiveExchangeOfferDTO) (entity.UserActiveExchangeOffer, error) {
	query := `
		SELECT user_id, currency_from, currency_to, exchange_rate, offer_created_at, offer_expires_at
		FROM user_active_exchange_offers
		WHERE user_id = $1 AND currency_from = $2 AND currency_to = $3`
	row := r.db.QueryRowContext(ctx, query, dto.UserID, dto.FromCurrency, dto.ToCurrency)
	var offer entity.UserActiveExchangeOffer
	err := row.Scan(&offer.UserID,
		&offer.FromCurrency,
		&offer.ToCurrency,
		&offer.ExchangeRate,
		&offer.OfferCreatedAt,
		&offer.OfferExpiresAt)
	if err != nil {
		return entity.UserActiveExchangeOffer{}, err
	}
	return offer, nil
}
