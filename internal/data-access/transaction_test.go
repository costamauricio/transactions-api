package dataAccess_test

import (
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/costamauricio/transactions-api/internal/data-access"
	"github.com/costamauricio/transactions-api/internal/data-access-test"
	"github.com/costamauricio/transactions-api/internal/models"
)

// Test the creation of a new transaction with success
func TestNewTransactionShouldCreateWithSuccess(t *testing.T) {
	mockedAccount, mock, closeFunc := dataAccessTest.GetMockedDatabase[dataAccess.Transaction](t)
	defer closeFunc()

	mockedId := int64(3)
	mockedAccountId := int(1)
	mockedType := models.TRANSACTION_TYPE_WITHDRAW
	mockedAmount := float64(20.5)

	mock.ExpectPrepare("INSERT INTO transactions").
		WillBeClosed().
		ExpectExec().
		WithArgs(mockedAccountId, mockedType, mockedAmount).
		WillReturnResult(
			sqlmock.NewResult(mockedId, 1),
		)

	insertedId, err := mockedAccount.NewTransaction(mockedAccountId, mockedType, mockedAmount)

	if err != nil {
		t.Errorf("Unexpected error: %s", err)
	}

	if insertedId != mockedId {
		t.Errorf("Expected insertedId to be %d but got %d", mockedId, insertedId)
	}
}
