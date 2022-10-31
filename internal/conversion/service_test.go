package conversion_test

import (
	"github.com/eneskzlcn/currency-conversion-service/internal/conversion"
	mocks "github.com/eneskzlcn/currency-conversion-service/internal/mocks/conversion"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestNewService(t *testing.T) {
	t.Run("given empty repository then it should return nil", func(t *testing.T) {
		service := conversion.NewService(nil)
		assert.Nil(t, service)
	})
	t.Run("given valid arguments then it should return service", func(t *testing.T) {
		mockConversionRepository := mocks.NewMockConversionRepository(gomock.NewController(t))
		service := conversion.NewService(mockConversionRepository)
		assert.NotNil(t, service)
	})
}
