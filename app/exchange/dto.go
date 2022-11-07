package exchange

import "time"

type CreateExchangeRateOfferDTO struct {
	UserID         int
	FromCurrency   string
	ToCurrency     string
	ExchangeRate   float32
	OfferCreatedAt time.Time
	OfferExpiresAt int64
}

type ExchangeCurrencyDTO struct {
	FromCurrency string
	ToCurrency   string
}
