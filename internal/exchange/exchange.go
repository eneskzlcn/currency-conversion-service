package exchange

import "time"

const (
	ExchangeRateExpirationMinutes = 3
)

type ExchangeRateRequest struct {
	FromCurrency string `json:"from_currency"`
	ToCurrency   string `json:"to_currency"`
}

type ExchangeRateResponse struct {
	FromCurrency string    `json:"from_currency"`
	ToCurrency   string    `json:"to_currency"`
	ExchangeRate float32   `json:"exchange_rate"`
	CreatedAt    time.Time `json:"created_at"`
	ExpiresAt    int64     `json:"expires_at"`
}
