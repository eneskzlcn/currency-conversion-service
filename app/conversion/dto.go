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

type TransferBalanceBetweenUserWalletsDTO struct {
	userID       int
	fromCurrency string
	toCurrency   string
	balance      float32
}

func NewTransferBalanceBetweenUserWalletsDTO(userID int, fromCurrency string, toCurrency string,
	balance float32) TransferBalanceBetweenUserWalletsDTO {
	return TransferBalanceBetweenUserWalletsDTO{userID: userID, fromCurrency: fromCurrency,
		toCurrency: toCurrency, balance: balance}
}

func (t TransferBalanceBetweenUserWalletsDTO) UserID() int {
	return t.userID
}

func (t TransferBalanceBetweenUserWalletsDTO) FromCurrency() string {
	return t.fromCurrency
}

func (t TransferBalanceBetweenUserWalletsDTO) ToCurrency() string {
	return t.toCurrency
}

func (t TransferBalanceBetweenUserWalletsDTO) Balance() float32 {
	return t.balance
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
