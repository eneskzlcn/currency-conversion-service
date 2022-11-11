package exchange

import "time"

type CreateExchangeRateOfferDTO struct {
	userID         int
	fromCurrency   string
	toCurrency     string
	exchangeRate   float32
	offerCreatedAt time.Time
	offerExpiresAt int64
}

func NewCreateExchangeRateOfferDTO(userID int, fromCurrency, toCurrency string,
	exchangeRate float32, offerCreatedAt time.Time, offerExpiresAt int64) CreateExchangeRateOfferDTO {
	return CreateExchangeRateOfferDTO{
		userID:         userID,
		fromCurrency:   fromCurrency,
		toCurrency:     toCurrency,
		exchangeRate:   exchangeRate,
		offerCreatedAt: offerCreatedAt,
		offerExpiresAt: offerExpiresAt,
	}
}
func (c *CreateExchangeRateOfferDTO) UserID() int {
	return c.userID
}
func (c *CreateExchangeRateOfferDTO) FromCurrency() string {
	return c.fromCurrency
}
func (c *CreateExchangeRateOfferDTO) ToCurrency() string {
	return c.toCurrency
}
func (c *CreateExchangeRateOfferDTO) ExchangeRate() float32 {
	return c.exchangeRate
}
func (c *CreateExchangeRateOfferDTO) OfferCreatedAt() time.Time {
	return c.offerCreatedAt
}
func (c *CreateExchangeRateOfferDTO) OfferExpiresAt() int64 {
	return c.offerExpiresAt
}
