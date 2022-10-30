package exchange_test

import (
	"github.com/eneskzlcn/currency-conversion-service/internal/exchange"
	mocks "github.com/eneskzlcn/currency-conversion-service/internal/mocks/exchange"
	"github.com/gofiber/fiber/v2"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"net/http/httptest"
	"testing"
)

func TestNewHandler(t *testing.T) {
	t.Run("given currency service then it should return new handler", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		mockExchangeService := mocks.NewMockExchangeService(ctrl)
		handler := exchange.NewHandler(mockExchangeService)
		assert.NotNil(t, handler)
	})
	t.Run("given currency service then it should return new handler", func(t *testing.T) {
		handler := exchange.NewHandler(nil)
		assert.Nil(t, handler)
	})
}
func TestHandler_RegisterRoutes(t *testing.T) {
	app := fiber.New()
	ctrl := gomock.NewController(t)
	mockExchangeService := mocks.NewMockExchangeService(ctrl)
	handler := exchange.NewHandler(mockExchangeService)
	handler.RegisterRoutes(app)

	assertRouteRegistered(t, app, fiber.MethodGet, "/exchange/rate")
}
func assertRouteRegistered(t *testing.T, app *fiber.App, method, route string) {
	resp, err := app.Test(httptest.NewRequest(fiber.MethodGet, "/exchange/rate", nil))
	assert.Nil(t, err)
	assert.NotEqual(t, fiber.StatusNotFound, resp.StatusCode)
}
