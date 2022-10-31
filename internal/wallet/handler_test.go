package wallet_test

import (
	mocks "github.com/eneskzlcn/currency-conversion-service/internal/mocks/wallet"
	"github.com/eneskzlcn/currency-conversion-service/internal/wallet"
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
		assert.Nil(t, handler)
	})
}
