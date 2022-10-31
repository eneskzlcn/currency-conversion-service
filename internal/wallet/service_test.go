package wallet_test

import (
	mocks "github.com/eneskzlcn/currency-conversion-service/internal/mocks/wallet"
	"github.com/eneskzlcn/currency-conversion-service/internal/wallet"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestNewService(t *testing.T) {
	t.Run("given not valid arguments then it should return nil", func(t *testing.T) {
		assert.Nil(t, wallet.NewService(nil))
	})
	t.Run("given valid arguments then it should return new Service", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		mockWalletRepo := mocks.NewMockWalletRepository(ctrl)
		assert.NotNil(t, wallet.NewService(mockWalletRepo))
	})
}
