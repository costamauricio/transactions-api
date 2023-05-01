package handlers

import (
	"github.com/costamauricio/transactions-api/internal/models"
)

// Define a log interface that is used by the handlers
type Logger interface {
	Print(msg ...any)
	Errorf(format string, args ...interface{})
}

// Define the data access object the handlers expect to handle the account operations
type AccountDAOInterface interface {
	NewAccount(documentNumber string) (int64, error)
	GetAccount(id int) (*models.Account, error)
}

// Define the data access object the handlers expect to handle the transactions operations
type TransactionDAOInterface interface {
	NewTransaction(accountId int, operation models.TransactionType, amount float64) (int64, error)
}

// Dependencies of the application handler
// The application handlers will be methods of this struct
type ApplicationHandlers struct {
	Logger        Logger
	AccountDAO    AccountDAOInterface
	TransactioDAO TransactionDAOInterface
}
