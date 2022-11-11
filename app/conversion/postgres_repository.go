package conversion

import (
	"context"
	"database/sql"
	"github.com/eneskzlcn/currency-conversion-service/app/model"
	"go.uber.org/zap"
	"time"
)

type postgresRepository struct {
	logger *zap.SugaredLogger
	db     *sql.DB
}

func NewPostgresRepository(db *sql.DB, logger *zap.SugaredLogger) *postgresRepository {
	return &postgresRepository{db: db, logger: logger}
}

func (r *postgresRepository) GetExchangeOfferByID(ctx context.Context,
	dto GetExchangeRateOfferDTO) (model.ExchangeRateOffer, error) {
	query := `
		SELECT id, user_id, currency_from, currency_to, exchange_rate, offer_created_at, offer_expires_at
		FROM user_exchange_offers
		WHERE id = $1;`
	row := r.db.QueryRowContext(ctx, query, dto.ExchangeRateOfferID())
	var offer model.ExchangeRateOffer
	err := row.Scan(
		&offer.ID,
		&offer.UserID,
		&offer.FromCurrency,
		&offer.ToCurrency,
		&offer.ExchangeRate,
		&offer.OfferCreatedAt,
		&offer.OfferExpiresAt)
	return offer, err
}
func (r *postgresRepository) CreateUserConversion(ctx context.Context,
	dto CreateUserConversionDTO) (model.UserCurrencyConversion, error) {
	query := `INSERT INTO user_currency_conversions (user_id, currency_from, currency_to, 
		sender_balance_dec_amount, receiver_balance_inc_amount)
		VALUES($1, $2, $3, $4, $5)
		RETURNING created_at, id;`
	row := r.db.QueryRowContext(ctx, query, dto.userID, dto.FromCurrency(), dto.ToCurrency(),
		dto.SenderBalanceDecAmount(), dto.ReceiverBalanceIncAmount())
	var conversionID int
	var createdAt time.Time
	err := row.Scan(&createdAt, &conversionID)
	return model.UserCurrencyConversion{
		ID:                       conversionID,
		UserID:                   dto.UserID(),
		FromCurrency:             dto.FromCurrency(),
		ToCurrency:               dto.ToCurrency(),
		SenderBalanceDecAmount:   dto.SenderBalanceDecAmount(),
		ReceiverBalanceIncAmount: dto.ReceiverBalanceIncAmount(),
		CreatedAt:                createdAt,
	}, err
}
