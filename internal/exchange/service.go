package exchange

import (
	"context"
	"github.com/eneskzlcn/currency-conversion-service/internal/entity"
	"time"
)

type ExchangeRepository interface {
	IsCurrencyExists(ctx context.Context, currency string) (bool, error)
	GetExchangeValuesForGivenCurrencies(ctx context.Context, fromCurrency, toCurrency string) (entity.Exchange, error)
}
type Service struct {
	exchangeRepository ExchangeRepository
}

func NewService(repository ExchangeRepository) *Service {
	if repository == nil {
		return nil
	}
	return &Service{exchangeRepository: repository}
}
func (s *Service) CreateExchangeRate(ctx context.Context, request ExchangeRateRequest) (ExchangeRateResponse, error) {
	exists, err := s.exchangeRepository.IsCurrencyExists(ctx, request.FromCurrency)
	if err != nil {
		return ExchangeRateResponse{}, err
	}
	if !exists {
		return ExchangeRateResponse{}, NotValidCurrencyErr
	}
	exists, err = s.exchangeRepository.IsCurrencyExists(ctx, request.ToCurrency)
	if err != nil {
		return ExchangeRateResponse{}, err
	}
	if !exists {
		return ExchangeRateResponse{}, NotValidCurrencyErr
	}
	exchange, err := s.exchangeRepository.
		GetExchangeValuesForGivenCurrencies(ctx, request.FromCurrency, request.ToCurrency)
	if err != nil {
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
