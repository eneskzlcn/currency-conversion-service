package conversion

import "errors"

var (
	CurrencyConversionOfferExpiredErr     = errors.New("given conversion offer expired")
	NotEnoughBalanceForConversionOfferErr = errors.New("user not have enough balance to make the conversion")
)
