package dataAccess

import (
	"database/sql"

	"github.com/costamauricio/transactions-api/internal/models"
)

type Transaction struct {
	Pool *sql.DB
}

// Insert a new transaction into the database with the argument values and return the inserted id
func (db *Transaction) NewTransaction(accountId int, operation models.TransactionType, amount float64) (int64, error) {
	// We use the 'datetime()' function from the database to save the current time as the created_at value
	query, err := db.Pool.Prepare(
		"INSERT INTO transactions(account_id, operation_type, amount, created_at) VALUES(?, ?, ?, datetime())",
	)
	if err != nil {
		return 0, err
	}
	defer query.Close()

	result, err := query.Exec(accountId, operation, amount)
	if err != nil {
		return 0, err
	}

	id, _ := result.LastInsertId()

	return id, nil
}
