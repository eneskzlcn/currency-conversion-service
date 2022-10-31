package wallet_test

import (
	"github.com/eneskzlcn/currency-conversion-service/internal/common/testutil"
	mocks "github.com/eneskzlcn/currency-conversion-service/internal/mocks/wallet"
	"github.com/eneskzlcn/currency-conversion-service/internal/wallet"
	"github.com/gofiber/fiber/v2"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
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
