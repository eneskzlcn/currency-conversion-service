package exchange

import (
	"context"
	"errors"
	"github.com/eneskzlcn/currency-conversion-service/app/entity"
	"go.uber.org/zap"
	"time"
)

type Repository interface {
	IsCurrencyExists(ctx context.Context, currency string) (bool, error)
	GetExchangeValuesForGivenCurrencies(ctx context.Context, fromCurrency, toCurrency string) (entity.CurrencyExchangeValues, error)
	CreateExchangeRateOffer(ctx context.Context, dto CreateExchangeRateOfferDTO) (int, error)
}
type service struct {
	repository Repository
	logger     *zap.SugaredLogger
}

func NewService(repository Repository, logger *zap.SugaredLogger) *service {
	return &service{repository: repository, logger: logger}
}
func (s *service) PrepareExchangeRateOffer(ctx context.Context, userID int, request ExchangeRateRequest) (ExchangeRateResponse, error) {
	if err := s.checkAreCurrenciesExists(ctx, request); err != nil {
		s.logger.Debug(err)
		return ExchangeRateResponse{}, err
	}
	exchangeValues, err := s.repository.
		GetExchangeValuesForGivenCurrencies(ctx, request.FromCurrency, request.ToCurrency)
	if err != nil {
		s.logger.Debug(err)
		return ExchangeRateResponse{}, err
	}
	createdExchangeRateOfferID, err := s.createExchangeRateOffer(ctx, userID, exchangeValues)
	if err != nil {
		return ExchangeRateResponse{}, err
	}
	return ExchangeRateResponse{ExchangeRateOfferID: createdExchangeRateOfferID}, nil
}
func (s *service) checkAreCurrenciesExists(ctx context.Context, request ExchangeRateRequest) error {
	exists, err := s.repository.IsCurrencyExists(ctx, request.FromCurrency)
	if err != nil {
		return err
	}
	if !exists {
		return errors.New(NotValidCurrency)
	}
	exists, err = s.repository.IsCurrencyExists(ctx, request.ToCurrency)
	if err != nil {
		return err
	}
	if !exists {
		return errors.New(NotValidCurrency)
	}
	return nil
}
func (s *service) createExchangeRateOffer(ctx context.Context, userID int, exchangeValues entity.CurrencyExchangeValues) (int, error) {
	exchangeRateOfferFeedWithMarkupRate := exchangeValues.ExchangeRate + exchangeValues.MarkupRate
	createdAt := time.Now()
	expiresAt := createdAt.Add(ExchangeRateExpirationMinutes * time.Minute).Unix()
	createExchangeRateOfferDTO := NewCreateExchangeRateOfferDTO(userID, exchangeValues.FromCurrency,
		exchangeValues.ToCurrency, exchangeRateOfferFeedWithMarkupRate, createdAt, expiresAt)

	return s.repository.CreateExchangeRateOffer(ctx, createExchangeRateOfferDTO)
}
