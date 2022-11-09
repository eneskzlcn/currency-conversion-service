package wallet

import (
	"context"
	"errors"
	"github.com/eneskzlcn/currency-conversion-service/app/entity"
	"go.uber.org/zap"
)

type Repository interface {
	GetUserWalletAccounts(ctx context.Context, userID int) ([]entity.UserWallet, error)
	GetUserBalanceOnGivenCurrency(ctx context.Context, userID int, currency string) (float32, error)
	AdjustUserBalanceOnGivenCurrency(ctx context.Context, userID int, currency string, balance float32) (bool, error)
	IsUserWithUserIDExists(ctx context.Context, userID int) (bool, error)
}

type service struct {
	repository Repository
	logger     *zap.SugaredLogger
}

func NewService(repository Repository, logger *zap.SugaredLogger) *service {
	return &service{repository: repository, logger: logger}
}

func (s *service) GetUserWalletAccounts(ctx context.Context, userID int) (UserWalletAccountsResponse, error) {
	exists, err := s.repository.IsUserWithUserIDExists(ctx, userID)
	if err != nil {
		return UserWalletAccountsResponse{}, err
	}
	if !exists {
		return UserWalletAccountsResponse{}, errors.New(UserWithUserIDNotExists)
	}
	userWalletsAccounts, err := s.repository.GetUserWalletAccounts(ctx, userID)
	s.logger.Debug(userWalletsAccounts)
	if err != nil {
		return UserWalletAccountsResponse{}, err
	}
	userWalletAccountsResponse := UserWalletAccountResponseFromUserWallets(userWalletsAccounts)
	s.logger.Debug(userWalletAccountsResponse)
	return userWalletAccountsResponse, nil
}

func (s *service) GetUserBalanceOnGivenCurrency(ctx context.Context, userID int,
	currency string) (float32, error) {
	return s.repository.GetUserBalanceOnGivenCurrency(ctx, userID, currency)
}

func (s *service) AdjustUserBalanceOnGivenCurrency(ctx context.Context,
	userID int, currency string, balance float32) (bool, error) {
	return s.repository.AdjustUserBalanceOnGivenCurrency(ctx, userID, currency, balance)
}
