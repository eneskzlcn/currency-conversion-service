package seed_test

import (
	"context"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/eneskzlcn/currency-conversion-service/seed"
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"regexp"
	"testing"
)

func TestMigrateTables(t *testing.T) {
	db, mock, _ := sqlmock.New()
	fileBytes, err := ioutil.ReadFile("./create_seed.sql")
	assert.Nil(t, err)
	mock.ExpectExec(regexp.QuoteMeta(string(fileBytes))).WillReturnResult(sqlmock.NewResult(1, 1))
	err = seed.MigrateTables(context.Background(), db)
	assert.Nil(t, err)
	err = mock.ExpectationsWereMet()
	assert.Nil(t, err)
}

func TestDropTables(t *testing.T) {
	db, mock, _ := sqlmock.New()
	fileBytes, err := ioutil.ReadFile("./drop_seed.sql")
	assert.Nil(t, err)
	mock.ExpectExec(regexp.QuoteMeta(string(fileBytes))).
		WillReturnResult(sqlmock.NewResult(1, 1))
	err = seed.DropTables(context.Background(), db)
	assert.Nil(t, err)
	err = mock.ExpectationsWereMet()
	assert.Nil(t, err)
}
