package exchange

const (
	ExchangeRateExpirationMinutes = 3
)

type ExchangeRateRequest struct {
	FromCurrency string `json:"from_currency"`
	ToCurrency   string `json:"to_currency"`
}

type ExchangeRateResponse struct {
	ExchangeRateOfferID int `json:"exchange_rate_offer_id"`
}
