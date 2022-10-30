package entity

import "time"

type Exchange struct {
	FromCurrency string
	ToCurrency   string
	ExchangeRate float32
	MarkupRate   float32
	CreatedAt    time.Time
	UpdatedAt    time.Time
}
