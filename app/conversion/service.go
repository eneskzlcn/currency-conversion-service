package conversion

import (
	"context"
	"errors"
	"github.com/eneskzlcn/currency-conversion-service/app/entity"
	"go.uber.org/zap"
	"time"
)

type WalletService interface {
	GetUserBalanceOnGivenCurrency(ctx context.Context, userID int, currency string) (float32, error)
	AdjustUserBalanceOnGivenCurrency(ctx context.Context, userID int, currency string, balance float32) (bool, error)
}
type Repository interface {
	GetUserActiveExchangeOffer(ctx context.Context, dto UserActiveExchangeOfferDTO) (entity.UserActiveExchangeOffer, error)
}
type service struct {
	walletService WalletService
	logger        *zap.SugaredLogger
	repository    Repository
}

func NewService(walletService WalletService, logger *zap.SugaredLogger, repository Repository) *service {
	if walletService == nil {
		return nil
	}
	return &service{walletService: walletService, logger: logger, repository: repository}
}
func (s *service) ConvertCurrencies(ctx context.Context, userID int, request CurrencyConversionOfferRequest) (bool, error) {
	if err := s.checkExchangeOfferSameWithTheUserActiveExchangeOffer(ctx, userID, request); err != nil {
		s.logger.Debug(err)
		return false, err
	}
	if s.isCurrencyConversionOfferExpired(request.ExpiresAt) {
		s.logger.Debug(CurrencyConversionOfferExpired)
		return false, errors.New(CurrencyConversionOfferExpired)
	}
	isUserHasEnoughBalance, err := s.isUserHasEnoughBalanceToMakeConversion(ctx,
		userID, request.FromCurrency, request.Balance)
	if err != nil {
		s.logger.Debug(err)
		return false, err
	}
	if !isUserHasEnoughBalance {
		s.logger.Debug(NotEnoughBalanceForConversionOffer)
		return false, errors.New(NotEnoughBalanceForConversionOffer)
	}
	err = s.updateUserWalletBalancesByConversion(ctx, userID, request)
	if err != nil {
		s.logger.Debug(err)
		return false, err
	}
	return true, nil
}
func (s *service) checkExchangeOfferSameWithTheUserActiveExchangeOffer(ctx context.Context, userID int, request CurrencyConversionOfferRequest) error {
	userActiveExchangeOffer, err := s.repository.GetUserActiveExchangeOffer(ctx, UserActiveExchangeOfferDTO{
		UserID:       userID,
		FromCurrency: request.FromCurrency,
		ToCurrency:   request.ToCurrency,
	})
	if err != nil {
		s.logger.Debug(err)
		return err
	}
	if userActiveExchangeOffer.ExchangeRate != request.ExchangeRate ||
		userActiveExchangeOffer.OfferExpiresAt != request.ExpiresAt {
		return errors.New(ExchangeOfferNotSameWithSavedUserOffer)
	}
	return nil
}
func (s *service) isCurrencyConversionOfferExpired(expiresAtUnix int64) bool {
	return expiresAtUnix < time.Now().Local().Unix()
}
func (s *service) isUserHasEnoughBalanceToMakeConversion(ctx context.Context, userID int, currency string, conversionBalance float32) (bool, error) {
	userBalanceInCurrencyFrom, err := s.walletService.
		GetUserBalanceOnGivenCurrency(ctx, userID, currency)
	if err != nil {
		return false, err
	}
	if userBalanceInCurrencyFrom < conversionBalance {
		return false, errors.New(NotEnoughBalanceForConversionOffer)
	}
	return true, nil
}
func (s *service) updateUserWalletBalancesByConversion(ctx context.Context, userID int, request CurrencyConversionOfferRequest) error {
	fromCurrencyBalanceAdjustAmount := -1 * request.Balance
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
