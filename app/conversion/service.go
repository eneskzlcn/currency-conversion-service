package conversion

import (
	"context"
	"errors"
	"github.com/eneskzlcn/currency-conversion-service/app/message"
	"github.com/eneskzlcn/currency-conversion-service/app/model"
	"go.uber.org/zap"
	"time"
)

type WalletService interface {
	GetUserBalanceOnGivenCurrency(ctx context.Context, userID int, currency string) (float32, error)
	AdjustUserBalanceOnGivenCurrency(ctx context.Context, userID int, currency string, balance float32) (bool, error)
}
type Repository interface {
	GetExchangeOfferByID(ctx context.Context, dto GetExchangeRateOfferDTO) (model.ExchangeRateOffer, error)
	CreateUserConversion(ctx context.Context, dto CreateUserConversionDTO) (model.UserCurrencyConversion, error)
}
type RabbitmqProducer interface {
	PushConversionCreatedMessage(message message.CurrencyConvertedMessage) error
}

type service struct {
	walletService    WalletService
	logger           *zap.SugaredLogger
	repository       Repository
	rabbitmqProducer RabbitmqProducer
}

func NewService(walletService WalletService, logger *zap.SugaredLogger, repository Repository,
	rabbitmqProducer RabbitmqProducer) *service {
	if walletService == nil {
		return nil
	}
	return &service{walletService: walletService, logger: logger, repository: repository, rabbitmqProducer: rabbitmqProducer}
}
func (s *service) ConvertCurrencies(ctx context.Context, userID int, request CurrencyConversionOfferRequest) (bool, error) {
	exchangeRateOffer, err := s.repository.GetExchangeOfferByID(ctx, GetExchangeRateOfferDTO{
		exchangeRateOfferID: request.ExchangeRateOfferID})
	if err != nil {
		s.logger.Debug(err)
		return false, err
	}
	if s.isCurrencyConversionOfferExpired(exchangeRateOffer.OfferExpiresAt) {
		s.logger.Debug(CurrencyConversionOfferExpired)
		return false, errors.New(CurrencyConversionOfferExpired)
	}
	isUserHasEnoughBalance, err := s.isUserHasEnoughBalanceToMakeConversion(ctx,
		userID, exchangeRateOffer.FromCurrency, request.Balance)
	if err != nil {
		s.logger.Debug(err)
		return false, err
	}
	if !isUserHasEnoughBalance {
		s.logger.Debug(NotEnoughBalanceForConversionOffer)
		return false, errors.New(NotEnoughBalanceForConversionOffer)
	}
	currencyConversion, err := s.repository.CreateUserConversion(ctx, CreateUserConversionDTO{
		userID:                 userID,
		fromCurrency:           exchangeRateOffer.FromCurrency,
		toCurrency:             exchangeRateOffer.ToCurrency,
		senderBalanceDecAmount: request.Balance,
	})
	if err != nil {
		s.logger.Error(err)
		return false, err
	}
	err = s.rabbitmqProducer.PushConversionCreatedMessage(message.CurrencyConvertedMessage{
		UserID:                   userID,
		FromCurrency:             currencyConversion.FromCurrency,
		ToCurrency:               currencyConversion.ToCurrency,
		SenderBalanceDecAmount:   currencyConversion.SenderBalanceDecAmount,
		ReceiverBalanceIncAmount: currencyConversion.ReceiverBalanceIncAmount,
	})
	if err != nil {
		s.logger.Debug(err)
		return false, err
	}
	return true, nil
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
