package wallet

import (
	"context"
	"github.com/eneskzlcn/currency-conversion-service/internal/entity"
)

type WalletRepository interface {
	GetUserWalletAccounts(ctx context.Context, userID int) ([]entity.UserWallet, error)
	IsUserWithUserIDExists(ctx context.Context, userID int) (bool, error)
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
func (s *Service) GetUserWalletAccounts(ctx context.Context, userID int) (UserWalletAccountsResponse, error) {
	exists, err := s.walletRepository.IsUserWithUserIDExists(ctx, userID)
	if err != nil {
		return UserWalletAccountsResponse{}, err
	}
	if !exists {
		return UserWalletAccountsResponse{}, UserWithUserIDNotExistsErr
	}
	userWalletsAccounts, err := s.walletRepository.GetUserWalletAccounts(ctx, userID)
	if err != nil {
		return UserWalletAccountsResponse{}, err
	}
	userWalletAccountsResponse := UserWalletAccountResponseFromUserWallets(userWalletsAccounts)
	return userWalletAccountsResponse, nil
}
