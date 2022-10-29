package auth_test

import (
	"github.com/eneskzlcn/currency-conversion-service/internal/auth"
	"github.com/eneskzlcn/currency-conversion-service/internal/config"
	"github.com/eneskzlcn/currency-conversion-service/internal/entity"
	mocks "github.com/eneskzlcn/currency-conversion-service/internal/mocks/auth"
	"github.com/gofiber/fiber/v2"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
)

func Test_NewGuard(t *testing.T) {
	t.Run("given valid auth service then it should return Guard when NewGuard called", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		mockAuthService := mocks.NewMockAuthService(ctrl)
		authGuard := auth.NewGuard(mockAuthService)
		assert.NotNil(t, authGuard)
	})
	t.Run("given nil auth service then it should return nil when NewGuard called", func(t *testing.T) {
		authGuard := auth.NewGuard(nil)
		assert.Nil(t, authGuard)
	})
}
func CreateFiberAppWithAMockProtectedEndpoint(guard *auth.Guard) *fiber.App {
	app := fiber.New()
	endpointToProtect := func(ctx *fiber.Ctx) error {
		return ctx.Status(fiber.StatusOK).SendString("Hello")
	}
	app.Get("/test", guard.ProtectWithJWT(endpointToProtect))
	return app
}
func makeTestRequestWithoutBodyToProtectedEndpoint(method string, route string, token string) *http.Request {
	req := httptest.NewRequest(method, route, nil)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Token", token)
	return req
}
func TestProtectWithJWTProtectsTheGivenHandlerWithJWTWhenItAppliedToAHandlerAsMiddleware(t *testing.T) {
	ctrl := gomock.NewController(t)
	mockAuthService := mocks.NewMockAuthService(ctrl)
	guard := auth.NewGuard(mockAuthService)
	app := CreateFiberAppWithAMockProtectedEndpoint(guard)

	t.Run("given valid token then it should return status unauthorized", func(t *testing.T) {
		givenConfig := config.Jwt{
			ATPrivateKey:        "private",
			ATExpirationSeconds: 10,
		}
		givenUser := entity.User{
			ID:       1,
			Username: "user",
			Password: "user",
		}
		token, err := CreateMockValidToken(givenConfig, givenUser)
		assert.Nil(t, err)

		req := makeTestRequestWithoutBodyToProtectedEndpoint(fiber.MethodGet, "/test", token)
		resp, err := app.Test(req)
		assert.Nil(t, err)
		assert.Equal(t, fiber.StatusOK, resp.StatusCode)
	})

	t.Run("given invalid token then it should return status unauthorized", func(t *testing.T) {
		req := makeTestRequestWithoutBodyToProtectedEndpoint(fiber.MethodGet, "/test", "invalidtoken")
		resp, err := app.Test(req)
		assert.Nil(t, err)
		assert.Equal(t, fiber.StatusUnauthorized, resp.StatusCode)
	})
}
