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
	GetCurrencyExchangeValuesByCurrency(ctx context.Context, dto ExchangeCurrencyDTO) (entity.CurrencyExchangeValues, error)
	CreateExchangeRateOffer(ctx context.Context, dto CreateExchangeRateOfferDTO) (int, error)
}
type Service struct {
	exchangeRepository ExchangeRepository
	logger             *zap.SugaredLogger
}

func NewService(repository ExchangeRepository, logger *zap.SugaredLogger) *Service {
	return &Service{exchangeRepository: repository, logger: logger}
}
func (s *Service) PrepareExchangeRateOffer(ctx context.Context, userID int, request ExchangeRateRequest) (ExchangeRateOfferResponse, error) {
	if err := s.checkExchangeValuesBetweenCurrenciesExists(ctx, request.FromCurrency,
		request.ToCurrency); err != nil {
		return ExchangeRateOfferResponse{}, err
	}
	currencyExchangeValues, err := s.exchangeRepository.
		GetCurrencyExchangeValuesByCurrency(ctx, ExchangeCurrencyDTO{
			FromCurrency: request.FromCurrency,
			ToCurrency:   request.ToCurrency,
		})
	if err != nil {
		return ExchangeRateOfferResponse{}, err
	}
	exchangeRateOfferID, err := s.createExchangeRateOffer(ctx, userID, currencyExchangeValues)
	if err != nil {
		return ExchangeRateOfferResponse{}, err
	}
	return ExchangeRateOfferResponse{
		ExchangeRateOfferID: exchangeRateOfferID,
	}, nil
}
func (s *Service) checkExchangeValuesBetweenCurrenciesExists(ctx context.Context, fromCurrency, toCurrency string) error {
	exists, err := s.exchangeRepository.IsCurrencyExists(ctx, fromCurrency)
	if err != nil {
		return err
	}
	if !exists {
		return errors.New(NotValidCurrency)
	}
	exists, err = s.exchangeRepository.IsCurrencyExists(ctx, toCurrency)
	if err != nil {
		return err
	}
	if !exists {
		return errors.New(NotValidCurrency)
	}
	return nil
}
func (s *Service) createExchangeRateOffer(ctx context.Context, userID int, values entity.CurrencyExchangeValues) (int, error) {
	exchangeRateOfferFeedWithMarkupRate := values.ExchangeRate + values.MarkupRate
	createdAt := time.Now()
	expiresAt := createdAt.Add(ExchangeRateExpirationMinutes * time.Minute).Unix()

	createExchangeRateOfferDTO := CreateExchangeRateOfferDTO{
		UserID:         userID,
		FromCurrency:   values.FromCurrency,
		ToCurrency:     values.ToCurrency,
		ExchangeRate:   exchangeRateOfferFeedWithMarkupRate,
		OfferCreatedAt: createdAt,
		OfferExpiresAt: expiresAt,
	}
	return s.exchangeRepository.CreateExchangeRateOffer(ctx, createExchangeRateOfferDTO)
}
