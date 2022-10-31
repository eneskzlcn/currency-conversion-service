package wallet

import (
	"context"
	"github.com/eneskzlcn/currency-conversion-service/internal/entity"
)

type WalletRepository interface {
	GetUserWalletAccounts(ctx context.Context, userID int) ([]entity.UserWallet, error)
	GetUserBalanceOnGivenCurrency(ctx context.Context, userID int, currency string) (float32, error)
	AdjustUserBalanceOnGivenCurrency(ctx context.Context, userID int, currency string, balance float32) (bool, error)
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
func (s *Service) GetUserBalanceOnGivenCurrency(ctx context.Context, userID int, currency string) (float32, error) {
	return s.walletRepository.GetUserBalanceOnGivenCurrency(ctx, userID, currency)
}
func (s *Service) AdjustUserBalanceOnGivenCurrency(ctx context.Context, userID int, currency string, balance float32) (bool, error) {
	panic("implement me")
}
