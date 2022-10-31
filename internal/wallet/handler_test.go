package wallet_test

import (
	"errors"
	"github.com/eneskzlcn/currency-conversion-service/internal/common/testutil"
	mocks "github.com/eneskzlcn/currency-conversion-service/internal/mocks/wallet"
	"github.com/eneskzlcn/currency-conversion-service/internal/wallet"
	"github.com/gofiber/fiber/v2"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"strconv"
	"testing"
)

func TestNewHandler(t *testing.T) {
	ctrl := gomock.NewController(t)
	mockWalletService := mocks.NewMockWalletService(ctrl)
	mockAuthGuard := mocks.NewMockAuthGuard(ctrl)
	t.Run("given invalid arguments then it should return nil", func(t *testing.T) {
		assert.Nil(t, wallet.NewHandler(nil, nil))
		assert.Nil(t, wallet.NewHandler(mockWalletService, nil))
		assert.Nil(t, wallet.NewHandler(nil, mockAuthGuard))
	})
	t.Run("given valid arguments then it should return new Handler", func(t *testing.T) {
		handler := wallet.NewHandler(mockWalletService, mockAuthGuard)
		assert.NotNil(t, handler)
	})
}
func TestHandler_GetUserWallets(t *testing.T) {
	handler, mockWalletService, _ := createHandlerWithMockWalletServiceAndAuthGuard(t)
	app := fiber.New()
	route := "/wallets"
	app.Get(route, handler.GetUserWalletAccounts)

	t.Run("not given userID or invalid userID from context then it should return status bad request", func(t *testing.T) {
		req := testutil.MakeTestRequestWithoutBody(fiber.MethodGet, route)
		resp, err := app.Test(req)
		assert.Nil(t, err)
		assert.Equal(t, fiber.StatusBadRequest, resp.StatusCode)
	})
	t.Run("given valid userID but error occurred on service then it should return status internal server error", func(t *testing.T) {
		userID := 2
		req := testutil.MakeTestRequestWithBody(fiber.MethodGet, route, nil)
		req.Header.Set("userID", strconv.FormatInt(int64(userID), 10))
		mockWalletService.EXPECT().GetUserWalletAccounts(gomock.Any(), userID).
			Return(wallet.UserWalletAccountsResponse{}, errors.New("error occurred on service"))
		resp, err := app.Test(req)
		assert.Nil(t, err)
		assert.Equal(t, fiber.StatusInternalServerError, resp.StatusCode)
	})
	t.Run("given valid userID then it should return wallet accounts with status ok", func(t *testing.T) {
		userID := 2
		req := testutil.MakeTestRequestWithBody(fiber.MethodGet, route, nil)
		req.Header.Set("userID", strconv.FormatInt(int64(userID), 10))
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
		testutil.AssertBodyEqual(t, resp.Body, expectedWalletAccountResp)
	})

}
func TestHandler_RegisterRoutes(t *testing.T) {
	app := fiber.New()
	handler, _, mockAuthGuard := createHandlerWithMockWalletServiceAndAuthGuard(t)
	mockAuthGuard.EXPECT().ProtectWithJWT(gomock.Any()).Return(func(ctx *fiber.Ctx) error { return nil })
	handler.RegisterRoutes(app)
	testutil.AssertRouteRegistered(t, app, fiber.MethodGet, "/wallets")
}

func createHandlerWithMockWalletServiceAndAuthGuard(t *testing.T) (*wallet.Handler, *mocks.MockWalletService, *mocks.MockAuthGuard) {
	ctrl := gomock.NewController(t)
	mockWalletService := mocks.NewMockWalletService(ctrl)
	mockAuthGuard := mocks.NewMockAuthGuard(ctrl)
	return wallet.NewHandler(mockWalletService, mockAuthGuard), mockWalletService, mockAuthGuard
}
