package exchange

import (
	"context"
	"errors"
	"github.com/eneskzlcn/currency-conversion-service/app/entity"
	"go.uber.org/zap"
	"time"
)

type ExchangeRepository interface {
	IsCurrencyExists(ctx context.Context, currency string) (bool, error)
	GetExchangeValuesForGivenCurrencies(ctx context.Context, fromCurrency, toCurrency string) (entity.Exchange, error)
	SetUserActiveExchangeRateOffer(ctx context.Context, offer entity.UserActiveExchangeOffer) (bool, error)
}
type Service struct {
	exchangeRepository ExchangeRepository
	logger             *zap.SugaredLogger
}

func NewService(repository ExchangeRepository, logger *zap.SugaredLogger) *Service {
	return &Service{exchangeRepository: repository, logger: logger}
}
func (s *Service) PrepareExchangeRateOffer(ctx context.Context, userID int, request ExchangeRateRequest) (ExchangeRateResponse, error) {
	if err := s.checkAreCurrenciesExists(ctx, request); err != nil {
		s.logger.Debug(err)
		return ExchangeRateResponse{}, err
	}
	exchange, err := s.exchangeRepository.
		GetExchangeValuesForGivenCurrencies(ctx, request.FromCurrency, request.ToCurrency)
	if err != nil {
		s.logger.Debug(err)
		return ExchangeRateResponse{}, err
	}
	exchangeRateOfferFeedWithMarkupRate := exchange.ExchangeRate + exchange.MarkupRate
	createdAt := time.Now()
	expiresAt := createdAt.Add(ExchangeRateExpirationMinutes * time.Minute).Unix()

	exchangeRateOfferSet, err := s.exchangeRepository.SetUserActiveExchangeRateOffer(ctx, entity.UserActiveExchangeOffer{
		UserID:         userID,
		FromCurrency:   request.FromCurrency,
		ToCurrency:     request.ToCurrency,
		ExchangeRate:   exchangeRateOfferFeedWithMarkupRate,
		OfferCreatedAt: createdAt,
		OfferExpiresAt: expiresAt,
	})
	if err != nil || !exchangeRateOfferSet {
		return ExchangeRateResponse{}, err
	}
	return ExchangeRateResponse{
		FromCurrency: exchange.FromCurrency,
		ToCurrency:   exchange.ToCurrency,
		ExchangeRate: exchangeRateOfferFeedWithMarkupRate,
		CreatedAt:    createdAt,
		ExpiresAt:    expiresAt,
	}, nil
}
func (s *Service) checkAreCurrenciesExists(ctx context.Context, request ExchangeRateRequest) error {
	exists, err := s.exchangeRepository.IsCurrencyExists(ctx, request.FromCurrency)
	if err != nil {
		return err
	}
	if !exists {
		return errors.New(NotValidCurrency)
	}
	exists, err = s.exchangeRepository.IsCurrencyExists(ctx, request.ToCurrency)
	if err != nil {
		return err
	}
	if !exists {
		return errors.New(NotValidCurrency)
	}
	return nil
}
