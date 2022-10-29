package auth_test

import (
	"github.com/eneskzlcn/currency-conversion-service/internal/auth"
	mocks "github.com/eneskzlcn/currency-conversion-service/internal/mocks/auth"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
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
