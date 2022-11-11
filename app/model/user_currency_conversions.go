package model

import "time"

type UserCurrencyConversion struct {
	ID                       int
	UserID                   int
	FromCurrency             string
	ToCurrency               string
	SenderBalanceDecAmount   float32
	ReceiverBalanceIncAmount float32
	CreatedAt                time.Time
}
