package dataAccess_test

import (
	"errors"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/costamauricio/transactions-api/internal/data-access"
	"github.com/costamauricio/transactions-api/internal/data-access-test"
)

// Test the creation of a new account with success
func TestNewAccountShouldCreateWithSuccess(t *testing.T) {
	mockedAccount, mock, closeFunc := dataAccessTest.GetMockedDatabase[dataAccess.Account](t)
	defer closeFunc()

	mockedId := int64(4)

	mock.ExpectPrepare("INSERT INTO accounts").
		WillBeClosed().
		ExpectExec().
		WithArgs("test").
		WillReturnResult(
			sqlmock.NewResult(mockedId, 1),
		)

	insertedId, err := mockedAccount.NewAccount("test")

	if err != nil {
		t.Errorf("Unexpected error: %s", err)
	}

	if insertedId != mockedId {
		t.Errorf("Expected insertedId to be %d but got %d", mockedId, insertedId)
	}
}

// Test if the NewAccount will handle the prepare statement error
func TestNewAccountShouldHandlePrepareError(t *testing.T) {
	mockedAccount, mock, closeFunc := dataAccessTest.GetMockedDatabase[dataAccess.Account](t)
	defer closeFunc()

	mockedError := errors.New("failed to prepare statmenet")

	mock.ExpectPrepare("INSERT INTO accounts").
		WillReturnError(mockedError)

	insertedId, err := mockedAccount.NewAccount("test")

	if err != mockedError {
		t.Errorf("Expecting error: %s. No error received", mockedError)
	}

	if insertedId != 0 {
		t.Errorf("Expected insertedId to be 0 but got %d", insertedId)
	}
}

// Test if the NewAccount will handle a query error
func TestNewAccountShouldHandleQueryError(t *testing.T) {
	mockedAccount, mock, closeFunc := dataAccessTest.GetMockedDatabase[dataAccess.Account](t)
	defer closeFunc()

	mockedError := errors.New("failed to run the query")

	mock.ExpectPrepare("INSERT INTO accounts").
		ExpectExec().
		WillReturnError(mockedError)

	insertedId, err := mockedAccount.NewAccount("test")

	if err != mockedError {
		t.Errorf("Expecting error: %s. No error received", mockedError)
	}

	if insertedId != 0 {
		t.Errorf("Expected insertedId to be 0 but got %d", insertedId)
	}
}

// Test if the GetAccount will return the account correctly
func TestGetAccountShouldReturnWithSuccess(t *testing.T) {
	mockedAccount, mock, closeFunc := dataAccessTest.GetMockedDatabase[dataAccess.Account](t)
	defer closeFunc()

	mockedId := int(2)

	mockedAccountRows := sqlmock.NewRows([]string{"id", "account_number"}).
		AddRow("2", "test")

	mock.ExpectPrepare("SELECT id, document_number FROM accounts").
		WillBeClosed().
		ExpectQuery().
		WithArgs(mockedId).
		WillReturnRows(
			mockedAccountRows,
		)

	account, err := mockedAccount.GetAccount(mockedId)

	if err != nil {
		t.Errorf("Unexpected error: %s", err)
	}

	if account.ID != mockedId {
		t.Errorf("Expected ID to be %d but got %d", mockedId, account.ID)
	}

	if account.DocumentNumber != "test" {
		t.Errorf("Expected DocumentNumber %s", account.DocumentNumber)
	}
}

// Test if the GetAccount will handle the prepare statement error
func TestGetAccountShouldHandlePrepareError(t *testing.T) {
	mockedAccount, mock, closeFunc := dataAccessTest.GetMockedDatabase[dataAccess.Account](t)
	defer closeFunc()

	mockedId := int(2)
	mockedError := errors.New("failed to prepare statement")

	mock.ExpectPrepare("SELECT id, document_number FROM accounts").
		WillReturnError(mockedError)

	account, err := mockedAccount.GetAccount(mockedId)

	if err != mockedError {
		t.Errorf("Expecting error: %s. No error received", mockedError)
	}

	if account != nil {
		t.Error("Expected account to be nil")
	}
}
