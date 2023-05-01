package dataAccessTest

import (
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/costamauricio/transactions-api/internal/data-access"
)

// Interface to identify the types accepted by the generic mock function
type daoMocksInterface interface {
	dataAccess.Account | dataAccess.Transaction
}

// Spawn a new mocked database connection and inject it into the generic,
// Will return the MockedDAO for the generic type, a mock object and a closeFunc to tear down the mocks
// and close the mocked connection, the closeFunc must be called at the end of the current test
func GetMockedDatabase[C daoMocksInterface](t *testing.T) (*C, sqlmock.Sqlmock, func()) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("Error '%s' when creating the stub connection", err)
	}

	generic := &C{Pool: db}

	return generic, mock, func() {
		db.Close()
	}
}
