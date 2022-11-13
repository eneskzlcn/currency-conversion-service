package conversion

import "errors"

type userBalanceAdequacyPolicy struct {
}

func NewUserBalanceAdequacyPolicy() *userBalanceAdequacyPolicy {
	return &userBalanceAdequacyPolicy{}
}
func (u *userBalanceAdequacyPolicy) IsAllowed(userBalance, conversionBalance float32) error {
	if userBalance < conversionBalance {
		return errors.New(NotEnoughBalanceForConversionOffer)
	}
	return nil
}
