package model

import "time"

type ExchangeRateOffer struct {
	ID             int
	UserID         int
	FromCurrency   string
	ToCurrency     string
	ExchangeRate   float32
	OfferCreatedAt time.Time
	OfferExpiresAt int64
}
