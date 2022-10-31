package conversion

import "time"

type CurrencyConversionOfferRequest struct {
	FromCurrency string    `json:"from_currency"`
	ToCurrency   string    `json:"to_currency"`
	ExchangeRate float32   `json:"exchange_rate"`
	CreatedAt    time.Time `json:"created_at"`
	ExpiresAt    int64     `json:"expires_at"`
	Balance      float32   `json:"balance"`
}
