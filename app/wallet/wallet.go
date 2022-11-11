package wallet

import "github.com/eneskzlcn/currency-conversion-service/app/model"

type UserWalletAccountsResponse struct {
	Accounts []UserWalletAccount `json:"accounts"`
}
type UserWalletAccount struct {
	Currency string  `json:"currency"`
	Balance  float32 `json:"balance"`
}

func UserWalletAccountResponseFromUserWallets(userWallets []model.UserWallet) UserWalletAccountsResponse {
	var response UserWalletAccountsResponse
	for _, userWallet := range userWallets {
		userWalletAccount := UserWalletAccount{
			Currency: userWallet.Currency,
			Balance:  userWallet.Balance,
		}
		response.Accounts = append(response.Accounts, userWalletAccount)
	}
	return response
}
