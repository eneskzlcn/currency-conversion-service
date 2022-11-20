package auth_test

import (
	"context"
	"database/sql/driver"
	"errors"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/eneskzlcn/currency-conversion-service/app/auth"
	"github.com/eneskzlcn/currency-conversion-service/app/model"
	"github.com/eneskzlcn/currency-conversion-service/postgres"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
	"regexp"
	"testing"
	"time"
)

func TestPostgresRepository_GetUserByUsernameAndPassword(t *testing.T) {
	db, sqlMock := postgres.NewMockPostgres()
	repository := auth.NewPostgresRepository(db, zap.S())
	query := regexp.QuoteMeta(`
		SELECT id, username, password, created_at, updated_at
		FROM users WHERE username = $1 AND password = $2`)

	t.Run("given not existing username and password then it should return error", func(t *testing.T) {
		givenUsername := "notexistinguser"
		givenPassword := "notexistingpass"
		//
		//sqlMock.ExpectQuery(query).
		//	WithArgs(givenUsername, givenPassword).
		//	WillReturnError(errors.New("user not exist"))
		ExpectSqlQueryThatReturnError(sqlMock, query, givenUsername, givenPassword)
		user, err := repository.GetUserByUsernameAndPassword(context.Background(), givenUsername, givenPassword)
		assert.NotNil(t, err)
		assert.Empty(t, user)
		assert.Nil(t, sqlMock.ExpectationsWereMet())
	})
	t.Run("given existing username and password then it should return user without error", func(t *testing.T) {
		givenUsername := "existinguser"
		givenPassword := "existingpass"
		expectedUser := model.User{
			ID:        1,
			Username:  givenUsername,
			Password:  givenPassword,
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		}
		expectedRows := sqlMock.NewRows([]string{"id", "username", "password", "created_at", "updated_at"}).
			AddRow(expectedUser.ID, expectedUser.Username,
				expectedUser.Password, expectedUser.CreatedAt, expectedUser.UpdatedAt)

		sqlMock.ExpectQuery(query).
			WithArgs(givenUsername, givenPassword).
			WillReturnRows(expectedRows)

		user, err := repository.GetUserByUsernameAndPassword(context.Background(), givenUsername, givenPassword)
		assert.Nil(t, err)
		assert.Nil(t, sqlMock.ExpectationsWereMet())
		assert.Equal(t, expectedUser, user)
	})
}
func ExpectSqlQueryThatReturnError(sqlMock sqlmock.Sqlmock, query string, args ...driver.Value) {
	sqlMock.ExpectQuery(query).WithArgs(args...).WillReturnError(errors.New("error occurred"))
}
