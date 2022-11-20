package wallet_test

import (
	"errors"
	"github.com/eneskzlcn/currency-conversion-service/app/common"
	"github.com/eneskzlcn/currency-conversion-service/app/common/testhttp"
	mocks "github.com/eneskzlcn/currency-conversion-service/app/mocks/wallet"
	"github.com/eneskzlcn/currency-conversion-service/app/wallet"
	"github.com/gofiber/fiber/v2"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
	"testing"
)

type HttpHandler interface {
	RegisterRoutes(app *fiber.App)
	GetUserWalletAccounts(ctx *fiber.Ctx) error
}

func TestHandler_GetUserWallets(t *testing.T) {
	handler, mockWalletService, _ := createHandlerWithMockWalletServiceAndAuthGuard(t)
	route := "/wallets"
	userID := 2
	mockAuthMiddleware := func(handl fiber.Handler) fiber.Handler {
		return func(ctx *fiber.Ctx) error {
			ctx.Locals(common.USER_ID_CTX_KEY, userID)
			return handl(ctx)
		}
	}

	t.Run("not given userID or invalid userID from context then it should return status bad request", func(t *testing.T) {
		app := fiber.New()
		app.Get(route, handler.GetUserWalletAccounts)
		req := testhttp.MakeRequest(fiber.MethodGet, route)
		resp, err := app.Test(req)
		assert.Nil(t, err)
		assert.Equal(t, fiber.StatusBadRequest, resp.StatusCode)
	})
	t.Run("given valid userID but error occurred on service then it should return status internal server error", func(t *testing.T) {
		app := fiber.New()
		app.Get(route, mockAuthMiddleware(handler.GetUserWalletAccounts))
		req := testhttp.MakeRequestWithBody(fiber.MethodGet, route, nil)

		mockWalletService.EXPECT().GetUserWalletAccounts(gomock.Any(), userID).
			Return(wallet.UserWalletAccountsResponse{}, errors.New("error occurred on service"))
		resp, err := app.Test(req)
		assert.Nil(t, err)
		assert.Equal(t, fiber.StatusInternalServerError, resp.StatusCode)
	})
	t.Run("given valid userID then it should return wallet accounts with status ok", func(t *testing.T) {
		app := fiber.New()
		app.Get(route, mockAuthMiddleware(handler.GetUserWalletAccounts))

		req := testhttp.MakeRequestWithBody(fiber.MethodGet, route, nil)
		expectedWalletAccountResp := wallet.UserWalletAccountsResponse{Accounts: []wallet.UserWalletAccount{
			{
				Currency: "TRY",
				Balance:  200,
			},
			{
				Currency: "USD",
				Balance:  300,
			},
		}}

		mockWalletService.EXPECT().GetUserWalletAccounts(gomock.Any(), userID).Return(expectedWalletAccountResp, nil)
		resp, err := app.Test(req)
		assert.Nil(t, err)
		assert.Equal(t, fiber.StatusOK, resp.StatusCode)
		testhttp.AssertBodyEqual(t, resp.Body, expectedWalletAccountResp)
	})

}
func TestHandler_RegisterRoutes(t *testing.T) {
	app := fiber.New()
	handler, _, mockAuthGuard := createHandlerWithMockWalletServiceAndAuthGuard(t)
	mockAuthGuard.EXPECT().ProtectWithJWT(gomock.Any()).Return(func(ctx *fiber.Ctx) error { return nil })
	handler.RegisterRoutes(app)
	testhttp.AssertRouteRegistered(t, app, fiber.MethodGet, "/wallets")
}

func createHandlerWithMockWalletServiceAndAuthGuard(t *testing.T) (HttpHandler, *mocks.MockService, *mocks.MockAuthGuard) {
	ctrl := gomock.NewController(t)
	mockWalletService := mocks.NewMockService(ctrl)
	mockAuthGuard := mocks.NewMockAuthGuard(ctrl)
	return wallet.NewHttpHandler(mockWalletService, mockAuthGuard, zap.S()), mockWalletService, mockAuthGuard
}
