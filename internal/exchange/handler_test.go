package exchange_test

import (
	"errors"
	"github.com/eneskzlcn/currency-conversion-service/internal/common/testutil"
	"github.com/eneskzlcn/currency-conversion-service/internal/exchange"
	mocks "github.com/eneskzlcn/currency-conversion-service/internal/mocks/exchange"
	"github.com/gofiber/fiber/v2"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestNewHandler(t *testing.T) {
	t.Run("given currency service and auth guard then it should return new handler", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		mockExchangeService := mocks.NewMockExchangeService(ctrl)
		mockAuthGuard := mocks.NewMockAuthGuard(ctrl)
		handler := exchange.NewHandler(mockExchangeService, mockAuthGuard)
		assert.NotNil(t, handler)
	})
	t.Run("given empty service or auth guard then it should return nil", func(t *testing.T) {
		handler := exchange.NewHandler(nil, nil)
		assert.Nil(t, handler)
		handler = exchange.NewHandler(nil, mocks.NewMockAuthGuard(gomock.NewController(t)))
		assert.Nil(t, handler)
		handler = exchange.NewHandler(mocks.NewMockExchangeService(gomock.NewController(t)), nil)
		assert.Nil(t, handler)
	})
}
func TestHandler_RegisterRoutes(t *testing.T) {
	app := fiber.New()
	handler, _, mockAuthGuard := createHandlerWithMockExchangeServiceAndAuthGuard(t)
	mockAuthGuard.EXPECT().ProtectWithJWT(gomock.Any()).Return(func(ctx *fiber.Ctx) error { return nil })
	handler.RegisterRoutes(app)

	testutil.AssertRouteRegistered(t, app, fiber.MethodGet, "/exchange/rate")
}
func TestHandler_GetExchangeRate(t *testing.T) {
	app := fiber.New()
	handler, mockExchangeService, _ := createHandlerWithMockExchangeServiceAndAuthGuard(t)
	app.Get("/rate", handler.GetExchangeRate)
	t.Run("given not valid exchange rate request then it should return bad request", func(t *testing.T) {
		givenRequest := "notvalidexchangerate"
		req := testutil.MakeTestRequestWithBody(fiber.MethodGet, "/rate", givenRequest)
		resp, err := app.Test(req)
		assert.Nil(t, err)
		assert.Equal(t, fiber.StatusBadRequest, resp.StatusCode)
	})
	t.Run("given valid exchange rate request but unexpected error occurred on service then return status internal server error", func(t *testing.T) {
		givenRequest := exchange.ExchangeRateRequest{
			FromCurrency: "TRY",
			ToCurrency:   "USD",
		}
		mockExchangeService.EXPECT().CreateExchangeRate(gomock.Any(), givenRequest).
			Return(exchange.ExchangeRateResponse{}, errors.New(""))

		req := testutil.MakeTestRequestWithBody(fiber.MethodGet, "/rate", givenRequest)
		resp, err := app.Test(req)
		assert.Nil(t, err)
		assert.Equal(t, fiber.StatusInternalServerError, resp.StatusCode)
	})
	t.Run("given valid exchange rate then it should return exchange rate response with status created", func(t *testing.T) {
		givenRequest := exchange.ExchangeRateRequest{
			FromCurrency: "TRY",
			ToCurrency:   "USD",
		}
		expectedResponse := exchange.ExchangeRateResponse{
			FromCurrency: "TRY",
			ToCurrency:   "USD",
			ExchangeRate: 0.23,
			CreatedAt:    time.Now(),
			ExpiresAt:    time.Now().Add(exchange.ExchangeRateExpirationMinutes * time.Minute).Unix(),
		}
		mockExchangeService.EXPECT().CreateExchangeRate(gomock.Any(), givenRequest).
			Return(expectedResponse, nil)

		req := testutil.MakeTestRequestWithBody(fiber.MethodGet, "/rate", givenRequest)
		resp, err := app.Test(req)
		assert.Nil(t, err)
		assert.Equal(t, fiber.StatusCreated, resp.StatusCode)
		testutil.AssertBodyEqual(t, resp.Body, expectedResponse)
	})
}

func createHandlerWithMockExchangeServiceAndAuthGuard(t *testing.T) (*exchange.Handler, *mocks.MockExchangeService, *mocks.MockAuthGuard) {
	ctrl := gomock.NewController(t)
	mockExchangeService := mocks.NewMockExchangeService(ctrl)
	mockAuthGuard := mocks.NewMockAuthGuard(ctrl)
	return exchange.NewHandler(mockExchangeService, mockAuthGuard), mockExchangeService, mockAuthGuard
}
