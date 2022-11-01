package entity

import "time"

type UserActiveExchangeOffer struct {
	UserID         int
	FromCurrency   string
	ToCurrency     string
	ExchangeRate   float32
	OfferCreatedAt time.Time
	OfferExpiresAt int64
}
