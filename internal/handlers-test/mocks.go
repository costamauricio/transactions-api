package handlersTest

import (
	"github.com/costamauricio/transactions-api/internal/models"
	"github.com/stretchr/testify/mock"
)

type AccountDAOMock struct {
	mock.Mock
}

func (acc *AccountDAOMock) NewAccount(documentNumber string) (int64, error) {
	args := acc.Called(documentNumber)
	return int64(args.Int(0)), args.Error(1)
}

func (acc *AccountDAOMock) GetAccount(id int) (*models.Account, error) {
	args := acc.Called(id)

	// Returns nil if we don't have the account
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}

	return args.Get(0).(*models.Account), args.Error(1)
}

type TransactionDAOMock struct {
	mock.Mock
}

func (tr *TransactionDAOMock) NewTransaction(accountId int, operation models.TransactionType, amount float64) (int64, error) {
	args := tr.Called(accountId, operation, amount)

	return int64(args.Int(0)), args.Error(1)
}

type LoggerMock struct{}

func (l *LoggerMock) Errorf(format string, args ...interface{}) {}
func (l *LoggerMock) Print(msg ...any)                          {}
