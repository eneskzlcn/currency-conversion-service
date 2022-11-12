package wallet

import "github.com/eneskzlcn/currency-conversion-service/app/message"

type TransferBalanceBetweenUserWalletsDTO struct {
	userID                   int
	fromCurrency             string
	toCurrency               string
	senderBalanceDecAmount   float32
	receiverBalanceIncAmount float32
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

func (t TransferBalanceBetweenUserWalletsDTO) SenderBalanceDecAmount() float32 {
	return t.senderBalanceDecAmount
}

func (t TransferBalanceBetweenUserWalletsDTO) ReceiverBalanceIncAmount() float32 {
	return t.receiverBalanceIncAmount
}

func NewTransferBalanceBetweenUserWalletsDTO(userID int, fromCurrency string, toCurrency string,
	senderBalanceDecAmount float32, receiverBalanceIncAmount float32) TransferBalanceBetweenUserWalletsDTO {
	return TransferBalanceBetweenUserWalletsDTO{userID: userID, fromCurrency: fromCurrency,
		toCurrency: toCurrency, senderBalanceDecAmount: senderBalanceDecAmount,
		receiverBalanceIncAmount: receiverBalanceIncAmount}
}

func ToTransferBalanceBetweenUserWalletsDTO(msg message.CurrencyConvertedMessage) TransferBalanceBetweenUserWalletsDTO {
	return TransferBalanceBetweenUserWalletsDTO{
		userID:                   msg.UserID,
		fromCurrency:             msg.FromCurrency,
		toCurrency:               msg.ToCurrency,
		senderBalanceDecAmount:   msg.SenderBalanceDecAmount,
		receiverBalanceIncAmount: msg.ReceiverBalanceIncAmount,
	}
}
