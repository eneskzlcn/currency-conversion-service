package conversion

type CurrencyConversionOfferRequest struct {
	ExchangeRateOfferID int     `json:"exchange_rate_offer_id"`
	Balance             float32 `json:"senderBalanceDecAmount"`
}
