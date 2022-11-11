package model

import "time"

type CurrencyExchangeValues struct {
	FromCurrency string
	ToCurrency   string
	ExchangeRate float32
	MarkupRate   float32
	CreatedAt    time.Time
	UpdatedAt    time.Time
}
