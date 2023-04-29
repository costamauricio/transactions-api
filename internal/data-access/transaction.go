package dataAccess

import (
	"database/sql"

	"github.com/costamauricio/transactions-api/internal/models"
)

type Transaction struct {
	Pool *sql.DB
}

func (db *Transaction) NewTransaction(accountId int, operation models.TransactionType, amount float64) (int64, error) {
	query, err := db.Pool.Prepare(
		"INSERT INTO transactions(account_id, operation_type, amount, created_at) VALUES(?, ?, ?, now())",
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
