package dataAccess

import (
	"database/sql"

	"github.com/costamauricio/transactions-api/internal/models"
)

type Account struct {
	Pool *sql.DB
}

func (db *Account) NewAccount(accountNumber string) (int64, error) {
	query, err := db.Pool.Prepare("INSERT INTO accounts(account_number) VALUES(?)")
	if err != nil {
		return 0, err
	}
	defer query.Close()

	result, err := query.Exec(accountNumber)
	if err != nil {
		return 0, err
	}

	id, _ := result.LastInsertId()

	return id, nil
}

func (db *Account) GetAccount(id int) (*models.Account, error) {
	query, err := db.Pool.Prepare("SELECT id, account_number FROM accounts WHERE id = ?")
	if err != nil {
		return nil, err
	}
	defer query.Close()

	account := &models.Account{}

	row := query.QueryRow(id)
	err = row.Scan(&account.ID, &account.AccountNumber)

	if err != nil && err == sql.ErrNoRows {
		return nil, nil
	}

	if err != nil {
		return nil, err
	}

	return account, nil
}
