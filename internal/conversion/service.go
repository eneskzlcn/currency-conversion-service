package conversion

import (
	"context"
	"errors"
	"go.uber.org/zap"
	"time"
)

type WalletService interface {
	GetUserBalanceOnGivenCurrency(ctx context.Context, userID int, currency string) (float32, error)
	AdjustUserBalanceOnGivenCurrency(ctx context.Context, userID int, currency string, balance float32) (bool, error)
}
type Service struct {
	walletService WalletService
	logger        *zap.SugaredLogger
}

func NewService(walletService WalletService, logger *zap.SugaredLogger) *Service {
	if walletService == nil {
		return nil
	}
	return &Service{walletService: walletService, logger: logger}
}
func (s *Service) ConvertCurrencies(ctx context.Context, userID int, request CurrencyConversionOfferRequest) (bool, error) {
	if !s.isValidConversionOfferExchangeRate(request.ExpiresAt) {
		return false, CurrencyConversionOfferExpiredErr
	}
	isUserHasEnoughBalance, err := s.isUserHasEnoughBalanceToMakeConversion(ctx,
		userID, request.FromCurrency, request.Balance)
	if err != nil {
		return false, err
	}
	if !isUserHasEnoughBalance {
		return false, NotEnoughBalanceForConversionOfferErr
	}
	err = s.updateUserWalletBalancesByConversion(ctx, userID, request)
	if err != nil {
		return false, err
	}
	return true, nil
}
func (s *Service) isValidConversionOfferExchangeRate(expiresAtUnix int64) bool {
	return expiresAtUnix >= time.Now().Local().Unix()
}
func (s *Service) isUserHasEnoughBalanceToMakeConversion(ctx context.Context, userID int, currency string, conversionBalance float32) (bool, error) {
	userBalanceInCurrencyFrom, err := s.walletService.
		GetUserBalanceOnGivenCurrency(ctx, userID, currency)
	if err != nil {
		return false, err
	}
	if userBalanceInCurrencyFrom < conversionBalance {
		return false, NotEnoughBalanceForConversionOfferErr
	}
	return true, nil
}
func (s *Service) updateUserWalletBalancesByConversion(ctx context.Context, userID int, request CurrencyConversionOfferRequest) error {
	fromCurrencyBalanceAdjustAmount := request.Balance
	toCurrencyBalanceAdjustAmount := request.Balance * request.ExchangeRate

	success, err := s.walletService.AdjustUserBalanceOnGivenCurrency(ctx, userID,
		request.FromCurrency, fromCurrencyBalanceAdjustAmount)

	if err != nil {
		return err
	}
	if !success {
		return errors.New("error occurred on database")
	}
	success, err = s.walletService.AdjustUserBalanceOnGivenCurrency(ctx, userID,
		request.ToCurrency, toCurrencyBalanceAdjustAmount)

	//some retry or backup must be done here.
	if err != nil {
		return err
	}
	if !success {
		return errors.New("error occurred on database")
	}
	return nil
}
