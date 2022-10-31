package wallet

type UserWalletAccountsResponse struct {
	Accounts []UserWalletAccount `json:"accounts"`
}
type UserWalletAccount struct {
	Currency string  `json:"currency"`
	Balance  float32 `json:"balance"`
}
