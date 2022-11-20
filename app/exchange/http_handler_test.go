package exchange_test

import (
	"errors"
	"github.com/eneskzlcn/currency-conversion-service/app/common"
	"github.com/eneskzlcn/currency-conversion-service/app/common/testhttp"
	"github.com/eneskzlcn/currency-conversion-service/app/exchange"
	mocks "github.com/eneskzlcn/currency-conversion-service/app/mocks/exchange"
	"github.com/gofiber/fiber/v2"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
	"testing"
)

type HttpHandler interface {
	RegisterRoutes(app *fiber.App)
	GetExchangeRateOffer(ctx *fiber.Ctx) error
}

func TestHandler_RegisterRoutes(t *testing.T) {
	app := fiber.New()
	httpHandler, _, mockAuthGuard := createHandlerWithMockExchangeServiceAndAuthGuard(t)
	mockAuthGuard.EXPECT().ProtectWithJWT(gomock.Any()).Return(func(ctx *fiber.Ctx) error { return nil })
	httpHandler.RegisterRoutes(app)

	testhttp.AssertRouteRegistered(t, app, fiber.MethodGet, "/exchange/rate")
}
func TestHandler_GetExchangeRate(t *testing.T) {
	httpHandler, mockExchangeService, _ := createHandlerWithMockExchangeServiceAndAuthGuard(t)
	userID := 1
	mockAuthMiddleware := func(handl fiber.Handler) fiber.Handler {
		return func(ctx *fiber.Ctx) error {
			ctx.Locals(common.USER_ID_CTX_KEY, userID)
			return handl(ctx)
		}
	}
	t.Run("given not valid exchange rate request then it should return bad request", func(t *testing.T) {
		app := fiber.New()
		app.Get("/rate", httpHandler.GetExchangeRateOffer)
		givenRequest := "notvalidexchangerate"
		req := testhttp.MakeRequestWithBody(fiber.MethodGet, "/rate", givenRequest)
		resp, err := app.Test(req)
		assert.Nil(t, err)
		assert.Equal(t, fiber.StatusBadRequest, resp.StatusCode)
	})
	t.Run("given valid exchange rate request but unexpected error occurred on service then return status internal server error", func(t *testing.T) {
		app := fiber.New()
		app.Get("/rate", mockAuthMiddleware(httpHandler.GetExchangeRateOffer))
		givenRequest := exchange.ExchangeRateRequest{
			FromCurrency: "TRY",
			ToCurrency:   "USD",
		}

		mockExchangeService.EXPECT().PrepareExchangeRateOffer(gomock.Any(), userID, givenRequest).
			Return(exchange.ExchangeRateResponse{}, errors.New(""))

		req := testhttp.MakeRequestWithBody(fiber.MethodGet, "/rate", givenRequest)
		resp, err := app.Test(req)
		assert.Nil(t, err)
		assert.Equal(t, fiber.StatusInternalServerError, resp.StatusCode)
	})
	t.Run("given valid exchange rate then it should return exchange rate response with status created", func(t *testing.T) {
		app := fiber.New()
		app.Get("/rate", mockAuthMiddleware(httpHandler.GetExchangeRateOffer))
		givenRequest := exchange.ExchangeRateRequest{
			FromCurrency: "TRY",
			ToCurrency:   "USD",
		}
		expectedResponse := exchange.ExchangeRateResponse{
			ExchangeRateOfferID: 2,
		}
		mockExchangeService.EXPECT().PrepareExchangeRateOffer(gomock.Any(), gomock.Any(), givenRequest).
			Return(expectedResponse, nil)

		req := testhttp.MakeRequestWithBody(fiber.MethodGet, "/rate", givenRequest)
		resp, err := app.Test(req)
		assert.Nil(t, err)
		assert.Equal(t, fiber.StatusCreated, resp.StatusCode)
		testhttp.AssertBodyEqual(t, resp.Body, expectedResponse)
	})
}

func createHandlerWithMockExchangeServiceAndAuthGuard(t *testing.T) (HttpHandler, *mocks.MockService, *mocks.MockAuthGuard) {
	ctrl := gomock.NewController(t)
	mockExchangeService := mocks.NewMockService(ctrl)
	mockAuthGuard := mocks.NewMockAuthGuard(ctrl)
	return exchange.NewHttpHandler(mockExchangeService, mockAuthGuard, zap.S()), mockExchangeService, mockAuthGuard
}
