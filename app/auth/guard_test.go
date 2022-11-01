//go:build unit

package auth_test

import (
	"errors"
	"github.com/eneskzlcn/currency-conversion-service/app/auth"
	"github.com/eneskzlcn/currency-conversion-service/app/common/testutil"
	"github.com/eneskzlcn/currency-conversion-service/app/entity"
	mocks "github.com/eneskzlcn/currency-conversion-service/app/mocks/auth"
	"github.com/eneskzlcn/currency-conversion-service/config"
	"github.com/gofiber/fiber/v2"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
	"testing"
)

func TestProtectWithJWTProtectsTheGivenHandlerWithJWTWhenItAppliedToAHandlerAsMiddleware(t *testing.T) {
	ctrl := gomock.NewController(t)
	mockAuthService := mocks.NewMockAuthService(ctrl)
	guard := auth.NewGuard(mockAuthService, zap.S())
	app := createFiberAppWithAMockProtectedEndpoint(guard)

	t.Run("given valid token then it should call next handler without status unauthorized", func(t *testing.T) {
		givenConfig := config.Jwt{
			ATPrivateKey:        "private",
			ATExpirationMinutes: 10,
		}
		givenUser := entity.User{
			ID:       1,
			Username: "user",
			Password: "user",
		}
		token, err := createMockValidToken(givenConfig, givenUser)
		assert.Nil(t, err)

		mockAuthService.EXPECT().ValidateToken(gomock.Any(), token).Return(nil)
		mockAuthService.EXPECT().ExtractUserIDFromToken(token).Return(givenUser.ID, nil)
		req := testutil.MakeTestRequestWithoutBodyToProtectedEndpoint(fiber.MethodGet, "/test", token)
		resp, err := app.Test(req)
		assert.Nil(t, err)
		assert.Equal(t, fiber.StatusOK, resp.StatusCode)
	})

	t.Run("given invalid token then it should return status unauthorized", func(t *testing.T) {
		invalidToken := "invalidToken"
		mockAuthService.EXPECT().ValidateToken(gomock.Any(), invalidToken).Return(errors.New("authentication error"))
		req := testutil.MakeTestRequestWithoutBodyToProtectedEndpoint(fiber.MethodGet, "/test", invalidToken)
		resp, err := app.Test(req)
		assert.Nil(t, err)
		assert.Equal(t, fiber.StatusUnauthorized, resp.StatusCode)
	})
}

func createFiberAppWithAMockProtectedEndpoint(guard *auth.Guard) *fiber.App {
	app := fiber.New()
	endpointToProtect := func(ctx *fiber.Ctx) error {
		return ctx.Status(fiber.StatusOK).SendString("Hello")
	}
	app.Get("/test", guard.ProtectWithJWT(endpointToProtect))
	return app
}
