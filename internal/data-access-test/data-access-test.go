package dataAccessTest

import (
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/costamauricio/transactions-api/internal/data-access"
)

// Interface to identify the types accepted by the generic mock function
type databaseMocks interface {
	dataAccess.Account | dataAccess.Transaction
}

// Returns a MockedDAO of type, sqlmock and a closeFunc to tear down the mocks
func GetMockedDatabase[C databaseMocks](t *testing.T) (*C, sqlmock.Sqlmock, func()) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("Error '%s' when creating the stub connection", err)
	}

	generic := &C{Pool: db}

	return generic, mock, func() {
		db.Close()
	}
}
