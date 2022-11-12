package conversion

type GetExchangeRateOfferDTO struct {
	exchangeRateOfferID int
}

func NewGetExchangeRateOfferDTO(exchangeRateOfferID int) GetExchangeRateOfferDTO {
	return GetExchangeRateOfferDTO{exchangeRateOfferID: exchangeRateOfferID}
}
func (g GetExchangeRateOfferDTO) ExchangeRateOfferID() int {
	return g.exchangeRateOfferID
}

type CreateUserConversionDTO struct {
	userID                   int
	fromCurrency             string
	toCurrency               string
	senderBalanceDecAmount   float32
	receiverBalanceIncAmount float32
}

func NewCreateUserConversionDTO(userID int, fromCurrency string, toCurrency string, senderBalanceDecAmount float32, receiverBalanceIncAmount float32) CreateUserConversionDTO {
	return CreateUserConversionDTO{userID: userID, fromCurrency: fromCurrency,
		toCurrency: toCurrency, senderBalanceDecAmount: senderBalanceDecAmount,
		receiverBalanceIncAmount: receiverBalanceIncAmount}
}

func (c CreateUserConversionDTO) SenderBalanceDecAmount() float32 {
	return c.senderBalanceDecAmount
}

func (c CreateUserConversionDTO) ReceiverBalanceIncAmount() float32 {
	return c.receiverBalanceIncAmount
}

func (c CreateUserConversionDTO) UserID() int {
	return c.userID
}

func (c CreateUserConversionDTO) FromCurrency() string {
	return c.fromCurrency
}

func (c CreateUserConversionDTO) ToCurrency() string {
	return c.toCurrency
}
