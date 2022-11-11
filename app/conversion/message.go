package conversion

type CurrencyConvertedMessage struct {
	UserID                   int
	FromCurrency             string
	ToCurrency               string
	SenderBalanceDecAmount   float32
	ReceiverBalanceIncAmount float32
}
