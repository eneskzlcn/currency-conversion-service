package conversion_test

import (
	"errors"
	"github.com/eneskzlcn/currency-conversion-service/app/common"
	"github.com/eneskzlcn/currency-conversion-service/app/common/testutil"
	"github.com/eneskzlcn/currency-conversion-service/app/conversion"
	mocks "github.com/eneskzlcn/currency-conversion-service/app/mocks/conversion"
	"github.com/gofiber/fiber/v2"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
	"strconv"
	"testing"
	"time"
)

type HttpHandler interface {
	CurrencyConversionOffer(ctx *fiber.Ctx) error
	RegisterRoutes(app *fiber.App)
}

func TestHandler_ConvertCurrencies(t *testing.T) {
	httpHandler, mockConversionService, _ := createHandlerWithMockConversionServiceAndAuthGuard(t)
	route := "/offer"
	userID := 2
	mockAuthMiddleware := func(handl fiber.Handler) fiber.Handler {
		return func(ctx *fiber.Ctx) error {
			ctx.Locals(common.USER_ID_CTX_KEY, userID)
			return handl(ctx)
		}
	}
	t.Run("given conversion offer request but userID not in context then it should return status bad request", func(t *testing.T) {
		app := fiber.New()
		app.Post(route, httpHandler.CurrencyConversionOffer)
		givenOfferRequest := conversion.CurrencyConversionOfferRequest{
			FromCurrency: "TRY",
			ToCurrency:   "USD",
			ExchangeRate: 2.30,
			CreatedAt:    time.Now(),
			ExpiresAt:    time.Now().Add(3 * time.Minute).Unix(),
			Balance:      200,
		}
		req := testutil.MakeTestRequestWithBody(fiber.MethodPost, route, givenOfferRequest)
		resp, err := app.Test(req)
		assert.Nil(t, err)
		assert.Equal(t, fiber.StatusBadRequest, resp.StatusCode)
	})
	t.Run("given invalid conversion offer request then it should return status bad request", func(t *testing.T) {
		app := fiber.New()
		app.Post(route, mockAuthMiddleware(httpHandler.CurrencyConversionOffer))
		givenOfferRequest := "invalidRequest"
		req := testutil.MakeTestRequestWithBody(fiber.MethodPost, route, givenOfferRequest)
		resp, err := app.Test(req)
		assert.Nil(t, err)
		assert.Equal(t, fiber.StatusBadRequest, resp.StatusCode)
	})

	t.Run("given conversion offer request but error returned from service then it should return status internal server error ", func(t *testing.T) {
		app := fiber.New()
		app.Post(route, mockAuthMiddleware(httpHandler.CurrencyConversionOffer))
		givenOfferRequest := conversion.CurrencyConversionOfferRequest{
			FromCurrency: "TRY",
			ToCurrency:   "USD",
			ExchangeRate: 2.30,
			CreatedAt:    time.Now().Local(),
			ExpiresAt:    time.Now().Add(3 * time.Minute).Unix(),
			Balance:      200,
		}
		req := testutil.MakeTestRequestWithBody(fiber.MethodPost, route, givenOfferRequest)
		req.Header.Set("userID", strconv.FormatInt(int64(userID), 10))
		mockConversionService.EXPECT().ConvertCurrencies(gomock.Any(), userID, givenOfferRequest).
			Return(false, errors.New(""))

		resp, err := app.Test(req)
		assert.Nil(t, err)
		assert.Equal(t, fiber.StatusInternalServerError, resp.StatusCode)
	})
	t.Run("given conversion offer request then it should return status ok and ", func(t *testing.T) {
		app := fiber.New()
		app.Post(route, mockAuthMiddleware(httpHandler.CurrencyConversionOffer))
		givenOfferRequest := conversion.CurrencyConversionOfferRequest{
			FromCurrency: "TRY",
			ToCurrency:   "USD",
			ExchangeRate: 2.30,
			CreatedAt:    time.Now().Local(),
			ExpiresAt:    time.Now().Add(3 * time.Minute).Unix(),
			Balance:      200,
		}
		req := testutil.MakeTestRequestWithBody(fiber.MethodPost, route, givenOfferRequest)
		req.Header.Set("userID", strconv.FormatInt(int64(userID), 10))
		mockConversionService.EXPECT().ConvertCurrencies(gomock.Any(), userID, givenOfferRequest).
			Return(true, nil)

		resp, err := app.Test(req)
		assert.Nil(t, err)
		assert.Equal(t, fiber.StatusOK, resp.StatusCode)
	})
}
func TestHandler_RegisterRoutes(t *testing.T) {
	app := fiber.New()
	httpHandler, _, mockAuthGuard := createHandlerWithMockConversionServiceAndAuthGuard(t)
	mockAuthGuard.EXPECT().ProtectWithJWT(gomock.Any()).Return(func(ctx *fiber.Ctx) error { return nil })
	httpHandler.RegisterRoutes(app)
	testutil.AssertRouteRegistered(t, app, fiber.MethodPost, "/conversion/offer")

}
func createHandlerWithMockConversionServiceAndAuthGuard(t *testing.T) (HttpHandler, *mocks.MockService, *mocks.MockAuthGuard) {
	ctrl := gomock.NewController(t)
	mockConversionService := mocks.NewMockService(ctrl)
	mockAuthGuard := mocks.NewMockAuthGuard(ctrl)
	return conversion.NewHttpHandler(mockConversionService, mockAuthGuard, zap.S()),
		mockConversionService, mockAuthGuard
}
