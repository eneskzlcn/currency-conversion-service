package exchange

import (
	"context"
	"github.com/eneskzlcn/currency-conversion-service/app/entity"
	"go.uber.org/zap"
	"time"
)

type ExchangeRepository interface {
	IsCurrencyExists(ctx context.Context, currency string) (bool, error)
	GetExchangeValuesForGivenCurrencies(ctx context.Context, fromCurrency, toCurrency string) (entity.Exchange, error)
}
type Service struct {
	exchangeRepository ExchangeRepository
	logger             *zap.SugaredLogger
}

func NewService(repository ExchangeRepository, logger *zap.SugaredLogger) *Service {
	return &Service{exchangeRepository: repository, logger: logger}
}
func (s *Service) PrepareExchangeRateOffer(ctx context.Context, request ExchangeRateRequest) (ExchangeRateResponse, error) {
	if err := s.checkAreCurrenciesExists(ctx, request.FromCurrency, request.ToCurrency); err != nil {
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
	return ExchangeRateResponse{
		FromCurrency: exchange.FromCurrency,
		ToCurrency:   exchange.ToCurrency,
		ExchangeRate: exchangeRateOfferFeedWithMarkupRate,
		CreatedAt:    createdAt,
		ExpiresAt:    expiresAt,
	}, nil
}
func (s *Service) checkAreCurrenciesExists(ctx context.Context, currencyFrom, currencyTo string) error {
	exists, err := s.exchangeRepository.IsCurrencyExists(ctx, currencyFrom)
	if err != nil {
		return err
	}
	if !exists {
		return NotValidCurrencyErr
	}
	exists, err = s.exchangeRepository.IsCurrencyExists(ctx, currencyTo)
	if err != nil {
		return err
	}
	if !exists {
		return NotValidCurrencyErr
	}
	return nil
}
