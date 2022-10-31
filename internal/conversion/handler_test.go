package conversion_test

import (
	"errors"
	"github.com/eneskzlcn/currency-conversion-service/internal/common/testutil"
	"github.com/eneskzlcn/currency-conversion-service/internal/conversion"
	mocks "github.com/eneskzlcn/currency-conversion-service/internal/mocks/conversion"
	"github.com/gofiber/fiber/v2"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
	"strconv"
	"testing"
	"time"
)

func TestHandler_ConvertCurrencies(t *testing.T) {
	handler, mockConversionService, _ := createHandlerWithMockConversionServiceAndAuthGuard(t)
	app := fiber.New()
	route := "/offer"
	app.Post(route, handler.CurrencyConversionOffer)
	t.Run("given invalid conversion offer request then it should return status bad request", func(t *testing.T) {
		givenOfferRequest := "invalidRequest"
		req := testutil.MakeTestRequestWithBody(fiber.MethodPost, route, givenOfferRequest)
		resp, err := app.Test(req)
		assert.Nil(t, err)
		assert.Equal(t, fiber.StatusBadRequest, resp.StatusCode)
	})
	t.Run("given conversion offer request but userID not in context then it should return status internal server error", func(t *testing.T) {
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
	t.Run("given conversion offer request but error returned from service then it should return status internal server error ", func(t *testing.T) {
		givenOfferRequest := conversion.CurrencyConversionOfferRequest{
			FromCurrency: "TRY",
			ToCurrency:   "USD",
			ExchangeRate: 2.30,
			CreatedAt:    time.Now().Local(),
			ExpiresAt:    time.Now().Add(3 * time.Minute).Unix(),
			Balance:      200,
		}
		userID := 3
		req := testutil.MakeTestRequestWithBody(fiber.MethodPost, route, givenOfferRequest)
		req.Header.Set("userID", strconv.FormatInt(int64(userID), 10))
		mockConversionService.EXPECT().ConvertCurrencies(gomock.Any(), userID, givenOfferRequest).
			Return(false, errors.New(""))

		resp, err := app.Test(req)
		assert.Nil(t, err)
		assert.Equal(t, fiber.StatusInternalServerError, resp.StatusCode)
	})
	t.Run("given conversion offer request then it should return status ok and ", func(t *testing.T) {
		givenOfferRequest := conversion.CurrencyConversionOfferRequest{
			FromCurrency: "TRY",
			ToCurrency:   "USD",
			ExchangeRate: 2.30,
			CreatedAt:    time.Now().Local(),
			ExpiresAt:    time.Now().Add(3 * time.Minute).Unix(),
			Balance:      200,
		}
		userID := 3
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
	handler, _, mockAuthGuard := createHandlerWithMockConversionServiceAndAuthGuard(t)
	mockAuthGuard.EXPECT().ProtectWithJWT(gomock.Any()).Return(func(ctx *fiber.Ctx) error { return nil })
	handler.RegisterRoutes(app)
	testutil.AssertRouteRegistered(t, app, fiber.MethodPost, "/conversion/offer")

}
func createHandlerWithMockConversionServiceAndAuthGuard(t *testing.T) (*conversion.Handler, *mocks.MockConversionService, *mocks.MockAuthGuard) {
	ctrl := gomock.NewController(t)
	mockConversionService := mocks.NewMockConversionService(ctrl)
	mockAuthGuard := mocks.NewMockAuthGuard(ctrl)
	return conversion.NewHandler(mockConversionService, mockAuthGuard, zap.L().Sugar()),
		mockConversionService, mockAuthGuard
}
