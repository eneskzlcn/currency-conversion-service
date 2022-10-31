package auth_test

import (
	"context"
	"errors"
	"github.com/eneskzlcn/currency-conversion-service/internal/auth"
	"github.com/eneskzlcn/currency-conversion-service/internal/entity"
	"github.com/eneskzlcn/currency-conversion-service/postgres"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
	"regexp"
	"testing"
	"time"
)

func TestRepository_GetUserByUsernameAndPassword(t *testing.T) {
	db, sqlmock := postgres.NewMockPostgres()
	repository := auth.NewRepository(db, zap.S())
	query := regexp.QuoteMeta(`
		SELECT id, username, password, created_at, updated_at
		FROM users WHERE username = $1 AND password = $2`)

	t.Run("given not existing username and password then it should return error", func(t *testing.T) {
		givenUsername := "notexistinguser"
		givenPassword := "notexistingpass"

		sqlmock.ExpectQuery(query).
			WithArgs(givenUsername, givenPassword).
			WillReturnError(errors.New("user not exist"))
		user, err := repository.GetUserByUsernameAndPassword(context.Background(), givenUsername, givenPassword)
		assert.NotNil(t, err)
		assert.Empty(t, user)
		assert.Nil(t, sqlmock.ExpectationsWereMet())
	})
	t.Run("given existing username and password then it should return user without error", func(t *testing.T) {
		givenUsername := "existinguser"
		givenPassword := "existingpass"
		expectedUser := entity.User{
			ID:        1,
			Username:  givenUsername,
			Password:  givenPassword,
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		}
		expectedRows := sqlmock.NewRows([]string{"id", "username", "password", "created_at", "updated_at"}).
			AddRow(expectedUser.ID, expectedUser.Username,
				expectedUser.Password, expectedUser.CreatedAt, expectedUser.UpdatedAt)

		sqlmock.ExpectQuery(query).
			WithArgs(givenUsername, givenPassword).
			WillReturnRows(expectedRows)

		user, err := repository.GetUserByUsernameAndPassword(context.Background(), givenUsername, givenPassword)
		assert.Nil(t, err)
		assert.Nil(t, sqlmock.ExpectationsWereMet())
		assert.Equal(t, expectedUser, user)
	})
}
