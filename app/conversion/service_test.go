package conversion_test

import (
	"context"
	"errors"
	"github.com/eneskzlcn/currency-conversion-service/app/conversion"
	"github.com/eneskzlcn/currency-conversion-service/app/message"
	mocks "github.com/eneskzlcn/currency-conversion-service/app/mocks/conversion"
	"github.com/eneskzlcn/currency-conversion-service/app/model"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
	"testing"
	"time"
)

func TestService_CreateCurrencyConversion(t *testing.T) {
	service, mockWalletService,
		mockConversionRepo, mockRabbitmqProducer := createServiceWithMockWalletServiceAndConversionRepository(t)
	t.Run("given exchange rate offer id that not found in database then it should return false with error", func(t *testing.T) {
		givenConversionOfferReq := conversion.CurrencyConversionOfferRequest{
			ExchangeRateOfferID: 2,
			Balance:             200,
		}
		userID := 1
		givenUserActiveExchangeOfferDTO := conversion.NewGetExchangeRateOfferDTO(givenConversionOfferReq.ExchangeRateOfferID)
		mockConversionRepo.EXPECT().GetExchangeOfferByID(gomock.Any(), givenUserActiveExchangeOfferDTO).
			Return(model.ExchangeRateOffer{}, errors.New(""))
		success, err := service.ConvertCurrencies(context.TODO(), userID, givenConversionOfferReq)
		assert.NotNil(t, err)
		assert.False(t, success)
	})
	t.Run("given expired exchange offer request then it should return false with error", func(t *testing.T) {
		givenConversionOfferReq := conversion.CurrencyConversionOfferRequest{
			ExchangeRateOfferID: 2,
			Balance:             200,
		}
		userID := 2
		expectedExchangeRateOffer := model.ExchangeRateOffer{
			UserID:         userID,
			FromCurrency:   "TRY",
			ToCurrency:     "USD",
			ExchangeRate:   2.3,
			OfferCreatedAt: time.Now(),
			OfferExpiresAt: time.Now().Add(time.Hour * -3).Unix(),
		}
		givenGetExchangeRateOfferDTO := conversion.NewGetExchangeRateOfferDTO(givenConversionOfferReq.ExchangeRateOfferID)
		mockConversionRepo.EXPECT().GetExchangeOfferByID(gomock.Any(), givenGetExchangeRateOfferDTO).
			Return(expectedExchangeRateOffer, nil)
		success, err := service.ConvertCurrencies(context.TODO(), userID, givenConversionOfferReq)
		assert.False(t, success)
		assert.NotNil(t, err)
		assert.Equal(t, err.Error(), conversion.CurrencyConversionOfferExpired)
	})
	t.Run("given senderBalanceDecAmount amount that user not have in its from currency account then it should return false with error", func(t *testing.T) {
		givenConversionOfferReq := conversion.CurrencyConversionOfferRequest{
			ExchangeRateOfferID: 2,
			Balance:             200,
		}
		userID := 2
		expectedUserActiveExchangeOffer := model.ExchangeRateOffer{
			ID:             givenConversionOfferReq.ExchangeRateOfferID,
			UserID:         userID,
			FromCurrency:   "USD",
			ToCurrency:     "TRY",
			ExchangeRate:   2.30,
			OfferCreatedAt: time.Now(),
			OfferExpiresAt: time.Now().Add(10 * time.Minute).Unix(),
		}

		givenGetExchangeRateOfferDTO := conversion.NewGetExchangeRateOfferDTO(givenConversionOfferReq.ExchangeRateOfferID)
		mockConversionRepo.EXPECT().GetExchangeOfferByID(gomock.Any(), givenGetExchangeRateOfferDTO).
			Return(expectedUserActiveExchangeOffer, nil)

		mockWalletService.EXPECT().
			GetUserBalanceOnGivenCurrency(gomock.Any(), userID, expectedUserActiveExchangeOffer.FromCurrency).
			Return(float32(100), nil)

		success, err := service.ConvertCurrencies(context.TODO(), userID, givenConversionOfferReq)
		assert.False(t, success)
		assert.NotNil(t, err)
		assert.Equal(t, err.Error(), conversion.NotEnoughBalanceForConversionOffer)
	})
	t.Run("given error occurred when producing currency converted message on rabbitmq then it should return error", func(t *testing.T) {
		givenConversionOfferReq := conversion.CurrencyConversionOfferRequest{
			ExchangeRateOfferID: 2,
			Balance:             200,
		}
		userID := 2
		expectedUserActiveExchangeOffer := model.ExchangeRateOffer{
			ID:             givenConversionOfferReq.ExchangeRateOfferID,
			UserID:         userID,
			FromCurrency:   "USD",
			ToCurrency:     "TRY",
			ExchangeRate:   2.30,
			OfferCreatedAt: time.Now(),
			OfferExpiresAt: time.Now().Add(10 * time.Minute).Unix(),
		}

		givenGetExchangeRateOfferDTO := conversion.NewGetExchangeRateOfferDTO(givenConversionOfferReq.ExchangeRateOfferID)
		mockConversionRepo.EXPECT().GetExchangeOfferByID(gomock.Any(), givenGetExchangeRateOfferDTO).
			Return(expectedUserActiveExchangeOffer, nil)

		mockWalletService.EXPECT().
			GetUserBalanceOnGivenCurrency(gomock.Any(), userID, expectedUserActiveExchangeOffer.FromCurrency).
			Return(float32(500), nil)

		senderBalanceDecAmount := givenConversionOfferReq.Balance
		receiverBalanceIncAmount := givenConversionOfferReq.Balance * expectedUserActiveExchangeOffer.ExchangeRate
		currencyConvertedMessage := message.CurrencyConvertedMessage{
			UserID:                   userID,
			FromCurrency:             expectedUserActiveExchangeOffer.FromCurrency,
			ToCurrency:               expectedUserActiveExchangeOffer.ToCurrency,
			SenderBalanceDecAmount:   senderBalanceDecAmount,
			ReceiverBalanceIncAmount: receiverBalanceIncAmount,
		}
		createUserConversionDTO := conversion.NewCreateUserConversionDTO(userID,
			expectedUserActiveExchangeOffer.FromCurrency,
			expectedUserActiveExchangeOffer.ToCurrency, senderBalanceDecAmount, receiverBalanceIncAmount)

		expectedCurrencyConversion := model.UserCurrencyConversion{
			ID:                       2,
			UserID:                   userID,
			FromCurrency:             expectedUserActiveExchangeOffer.FromCurrency,
			ToCurrency:               expectedUserActiveExchangeOffer.ToCurrency,
			SenderBalanceDecAmount:   givenConversionOfferReq.Balance,
			ReceiverBalanceIncAmount: givenConversionOfferReq.Balance * expectedUserActiveExchangeOffer.ExchangeRate,
			CreatedAt:                time.Time{},
		}
		mockConversionRepo.EXPECT().CreateUserConversion(gomock.Any(), createUserConversionDTO).
			Return(expectedCurrencyConversion, nil)

		mockRabbitmqProducer.EXPECT().PushConversionCreatedMessage(currencyConvertedMessage).
			Return(errors.New(""))

		success, err := service.ConvertCurrencies(context.TODO(), userID,
			givenConversionOfferReq)
		assert.NotNil(t, err)
		assert.False(t, success)
	})
	t.Run("given valid currency conversion offer then it should return true", func(t *testing.T) {
		givenConversionOfferReq := conversion.CurrencyConversionOfferRequest{
			ExchangeRateOfferID: 2,
			Balance:             200,
		}
		userID := 2
		expectedUserActiveExchangeOffer := model.ExchangeRateOffer{
			ID:             givenConversionOfferReq.ExchangeRateOfferID,
			UserID:         userID,
			FromCurrency:   "USD",
			ToCurrency:     "TRY",
			ExchangeRate:   2.30,
			OfferCreatedAt: time.Now(),
			OfferExpiresAt: time.Now().Add(10 * time.Minute).Unix(),
		}

		givenGetExchangeRateOfferDTO := conversion.NewGetExchangeRateOfferDTO(givenConversionOfferReq.ExchangeRateOfferID)
		mockConversionRepo.EXPECT().GetExchangeOfferByID(gomock.Any(), givenGetExchangeRateOfferDTO).
			Return(expectedUserActiveExchangeOffer, nil)

		mockWalletService.EXPECT().
			GetUserBalanceOnGivenCurrency(gomock.Any(), userID, expectedUserActiveExchangeOffer.FromCurrency).
			Return(float32(500), nil)

		senderBalanceDecAmount := givenConversionOfferReq.Balance
		receiverBalanceIncAmount := givenConversionOfferReq.Balance * expectedUserActiveExchangeOffer.ExchangeRate

		currencyConvertedMessage := message.CurrencyConvertedMessage{
			UserID:                   userID,
			FromCurrency:             expectedUserActiveExchangeOffer.FromCurrency,
			ToCurrency:               expectedUserActiveExchangeOffer.ToCurrency,
			SenderBalanceDecAmount:   senderBalanceDecAmount,
			ReceiverBalanceIncAmount: receiverBalanceIncAmount,
		}
		createUserConversionDTO := conversion.NewCreateUserConversionDTO(userID,
			expectedUserActiveExchangeOffer.FromCurrency,
			expectedUserActiveExchangeOffer.ToCurrency, senderBalanceDecAmount, receiverBalanceIncAmount)
		expectedCurrencyConversion := model.UserCurrencyConversion{
			ID:                       2,
			UserID:                   userID,
			FromCurrency:             expectedUserActiveExchangeOffer.FromCurrency,
			ToCurrency:               expectedUserActiveExchangeOffer.ToCurrency,
			SenderBalanceDecAmount:   givenConversionOfferReq.Balance,
			ReceiverBalanceIncAmount: givenConversionOfferReq.Balance * expectedUserActiveExchangeOffer.ExchangeRate,
			CreatedAt:                time.Time{},
		}
		mockConversionRepo.EXPECT().CreateUserConversion(gomock.Any(), createUserConversionDTO).
			Return(expectedCurrencyConversion, nil)
		mockRabbitmqProducer.EXPECT().PushConversionCreatedMessage(currencyConvertedMessage).Return(nil)

		success, err := service.ConvertCurrencies(context.TODO(), userID, givenConversionOfferReq)
		assert.Nil(t, err)
		assert.True(t, success)
	})
}
func createServiceWithMockWalletServiceAndConversionRepository(t *testing.T) (conversion.Service, *mocks.MockWalletService,
	*mocks.MockRepository, *mocks.MockRabbitmqProducer) {
	ctrl := gomock.NewController(t)
	mockWalletService := mocks.NewMockWalletService(ctrl)
	mockConversionRepo := mocks.NewMockRepository(ctrl)
	mockRabbitmqProducer := mocks.NewMockRabbitmqProducer(ctrl)
	mockUserBalanceAdequacyPolicy := mocks.NewMockUserBalanceAdequacyPolicy(ctrl)
	return conversion.NewService(mockWalletService, zap.S(), mockConversionRepo,
		mockRabbitmqProducer, mockUserBalanceAdequacyPolicy), mockWalletService, mockConversionRepo, mockRabbitmqProducer
}
