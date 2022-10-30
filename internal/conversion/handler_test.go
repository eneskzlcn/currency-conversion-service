package conversion_test

import (
	"github.com/eneskzlcn/currency-conversion-service/internal/common/testutil"
	"github.com/eneskzlcn/currency-conversion-service/internal/conversion"
	mocks "github.com/eneskzlcn/currency-conversion-service/internal/mocks/conversion"
	"github.com/gofiber/fiber/v2"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestNewHandler(t *testing.T) {
	ctrl := gomock.NewController(t)
	mockConversionService := mocks.NewMockConversionService(ctrl)
	mockAuthGuard := mocks.NewMockAuthGuard(ctrl)

	t.Run("given empty service or auth guard then it should return nil", func(t *testing.T) {
		handler := conversion.NewHandler(nil, nil)
		assert.Nil(t, handler)
		handler = conversion.NewHandler(mockConversionService, nil)
		assert.Nil(t, handler)
		handler = conversion.NewHandler(nil, mockAuthGuard)
		assert.Nil(t, handler)
	})
	t.Run("given conversion service then it should return handler", func(t *testing.T) {
		handler := conversion.NewHandler(mockConversionService, mockAuthGuard)
		assert.NotNil(t, handler)
	})
}
func TestHandler_ConvertCurrencies(t *testing.T) {
	handler, _, _ := createHandlerWithMockConversionServiceAndAuthGuard(t)
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
	t.Run("given conversion offer request but error returned from service then it should return status internal server error ", func(t *testing.T) {

	})
	t.Run("given conversion offer request then it should return status ok and ", func(t *testing.T) {

	})
}
func TestHandler_RegisterRoutes(t *testing.T) {
	app := fiber.New()
	handler, _, mockAuthGuard := createHandlerWithMockConversionServiceAndAuthGuard(t)
	mockAuthGuard.EXPECT().ProtectWithJWT(gomock.Any()).Return(func(ctx *fiber.Ctx) error { return nil })
	handler.RegisterRoutes(app)
	testutil.AssertRouteRegistered(t, app, fiber.MethodPost, "/conversion/convert")

}
func createHandlerWithMockConversionServiceAndAuthGuard(t *testing.T) (*conversion.Handler, *mocks.MockConversionService, *mocks.MockAuthGuard) {
	ctrl := gomock.NewController(t)
	mockConversionService := mocks.NewMockConversionService(ctrl)
	mockAuthGuard := mocks.NewMockAuthGuard(ctrl)
	return conversion.NewHandler(mockConversionService, mockAuthGuard),
		mockConversionService, mockAuthGuard
}
