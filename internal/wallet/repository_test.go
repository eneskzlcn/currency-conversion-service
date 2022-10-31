package wallet_test

import (
	"context"
	"github.com/eneskzlcn/currency-conversion-service/internal/wallet"
	"github.com/eneskzlcn/currency-conversion-service/postgres"
	"github.com/stretchr/testify/assert"
	"regexp"
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
	db, sqlmock := postgres.NewMockPostgres()
	repository := wallet.NewRepository(db)
	query := regexp.QuoteMeta(`SELECT EXISTS (SELECT 1 FROM users WHERE id = $1`)
	t.Run("given existing user id then it should return true ", func(t *testing.T) {
		userID := 2
		expectedRows := sqlmock.NewRows([]string{"exists"}).AddRow(true)
		sqlmock.ExpectQuery(query).WithArgs(userID).
			WillReturnRows(expectedRows)
		exists, err := repository.IsUserWithUserIDExists(context.TODO(), userID)
		assert.Nil(t, sqlmock.ExpectationsWereMet())
		assert.Nil(t, err)
		assert.True(t, exists)
	})
	t.Run("given not existing user id then it should return false", func(t *testing.T) {
		userID := 4
		expectedRows := sqlmock.NewRows([]string{"exists"}).AddRow(false)
		sqlmock.ExpectQuery(query).WithArgs(userID).
			WillReturnRows(expectedRows)
		exists, err := repository.IsUserWithUserIDExists(context.TODO(), userID)
		assert.Nil(t, sqlmock.ExpectationsWereMet())
		assert.Nil(t, err)
		assert.False(t, exists)
	})
}
