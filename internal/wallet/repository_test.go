package wallet_test

import (
	"github.com/eneskzlcn/currency-conversion-service/internal/wallet"
	"github.com/eneskzlcn/currency-conversion-service/postgres"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestNewRepository(t *testing.T) {
	t.Run("given not valid arguments then it should return nil", func(t *testing.T) {
		assert.Nil(t, wallet.NewRepository(nil))
	})
	t.Run("given valid arguments then it should return new Repository", func(t *testing.T) {
		db, _ := postgres.NewMockPostgres()
		assert.NotNil(t, wallet.NewRepository(db))
	})
}
func TestRepository_IsUserWithUserIDExists(t *testing.T) {

}
