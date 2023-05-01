package dataAccess

import (
	"database/sql"

	"github.com/costamauricio/transactions-api/internal/models"
)

type Account struct {
	Pool *sql.DB
}

// Insert a new account on the database and returns the insertedId
func (db *Account) NewAccount(documentNumber string) (int64, error) {
	query, err := db.Pool.Prepare("INSERT INTO accounts(document_number) VALUES(?)")
	if err != nil {
		return 0, err
	}
	defer query.Close()

	result, err := query.Exec(documentNumber)
	if err != nil {
		return 0, err
	}

	id, _ := result.LastInsertId()

	return id, nil
}

// Retrieve an account from the database for the id argument
func (db *Account) GetAccount(id int) (*models.Account, error) {
	query, err := db.Pool.Prepare("SELECT id, document_number FROM accounts WHERE id = ?")
	if err != nil {
		return nil, err
	}
	defer query.Close()

	account := &models.Account{}

	row := query.QueryRow(id)
	err = row.Scan(&account.ID, &account.DocumentNumber)

	// If the error is sql.ErrNoRows means that there is no account for the ID
	// So we won't return an error, just a nil account
	if err != nil && err == sql.ErrNoRows {
		return nil, nil
	}

	if err != nil {
		return nil, err
	}

	return account, nil
}
