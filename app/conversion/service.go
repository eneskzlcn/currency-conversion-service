package conversion

import (
	"context"
	"errors"
	"github.com/eneskzlcn/currency-conversion-service/app/common/logutil"
	"github.com/eneskzlcn/currency-conversion-service/app/message"
	"github.com/eneskzlcn/currency-conversion-service/app/model"
	"go.uber.org/zap"
	"time"
)

type WalletService interface {
	GetUserBalanceOnGivenCurrency(ctx context.Context, userID int, currency string) (float32, error)
}
type Repository interface {
	GetExchangeOfferByID(ctx context.Context, dto GetExchangeRateOfferDTO) (model.ExchangeRateOffer, error)
	CreateUserConversion(ctx context.Context, dto CreateUserConversionDTO) (model.UserCurrencyConversion, error)
}
type RabbitmqProducer interface {
	PushConversionCreatedMessage(message message.CurrencyConvertedMessage) error
}
type UserBalanceAdequacyPolicy interface {
	IsAllowed(userBalance, conversionBalance float32) error
}
type ServiceOptions struct {
	WalletService    WalletService
	Logger           *zap.SugaredLogger
	Repository       Repository
	RabbitmqProducer RabbitmqProducer
}
type service struct {
	walletService             WalletService
	logger                    *zap.SugaredLogger
	repository                Repository
	rabbitmqProducer          RabbitmqProducer
	userBalanceAdequacyPolicy UserBalanceAdequacyPolicy
}

func NewService(opts *ServiceOptions) *service {
	return &service{
		walletService:             opts.WalletService,
		logger:                    opts.Logger,
		repository:                opts.Repository,
		rabbitmqProducer:          opts.RabbitmqProducer,
		userBalanceAdequacyPolicy: NewUserBalanceAdequacyPolicy(),
	}
}
func (s *service) ConvertCurrencies(ctx context.Context, userID int, request CurrencyConversionOfferRequest) (bool, error) {
	exchangeRateOffer, err := s.repository.GetExchangeOfferByID(ctx, GetExchangeRateOfferDTO{
		exchangeRateOfferID: request.ExchangeRateOfferID})
	if err != nil {
		return false, logutil.LogThenReturn(s.logger, err)
	}
	if s.isCurrencyConversionOfferExpired(exchangeRateOffer.OfferExpiresAt) {
		return false, logutil.LogThenReturn(s.logger, errors.New(CurrencyConversionOfferExpired))
	}
	isUserHasEnoughBalance, err := s.isUserHasEnoughBalanceToMakeConversion(ctx,
		userID, exchangeRateOffer.FromCurrency, request.Balance)
	if err != nil {
		return false, logutil.LogThenReturn(s.logger, err)
	}
	if !isUserHasEnoughBalance {
		return false, logutil.LogThenReturn(s.logger, errors.New(NotEnoughBalanceForConversionOffer))
	}
	currencyConversion, err := s.repository.CreateUserConversion(ctx, CreateUserConversionDTO{
		userID:                   userID,
		fromCurrency:             exchangeRateOffer.FromCurrency,
		toCurrency:               exchangeRateOffer.ToCurrency,
		senderBalanceDecAmount:   request.Balance * -1,
		receiverBalanceIncAmount: request.Balance * exchangeRateOffer.ExchangeRate,
	})
	if err != nil {
		return false, logutil.LogThenReturn(s.logger, err)
	}
	err = s.rabbitmqProducer.PushConversionCreatedMessage(message.CurrencyConvertedMessage{
		UserID:                   userID,
		FromCurrency:             currencyConversion.FromCurrency,
		ToCurrency:               currencyConversion.ToCurrency,
		SenderBalanceDecAmount:   currencyConversion.SenderBalanceDecAmount,
		ReceiverBalanceIncAmount: currencyConversion.ReceiverBalanceIncAmount,
	})
	if err != nil {
		return false, logutil.LogThenReturn(s.logger, err)
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

	if err = s.userBalanceAdequacyPolicy.IsAllowed(userBalanceInCurrencyFrom, conversionBalance); err != nil {
		return false, errors.New(NotEnoughBalanceForConversionOffer)
	}
	return true, nil
}
