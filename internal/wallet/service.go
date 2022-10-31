package wallet

import (
	"context"
	"github.com/eneskzlcn/currency-conversion-service/internal/entity"
)

type WalletRepository interface {
	GetUserWalletAccounts(ctx context.Context, userID int) ([]entity.UserWallet, error)
}

type Service struct {
	walletRepository WalletRepository
}

func NewService(repository WalletRepository) *Service {
	if repository == nil {
		return nil
	}
	return &Service{walletRepository: repository}
}
