package models

const (
	TRANSACTION_TYPE_PUSCHASE_CASH        TransactionType = 1
	TRANSACTION_TYPE_PUSCHASE_INSTALLMENT TransactionType = 2
	TRANSACTION_TYPE_WITHDRAW             TransactionType = 3
	TRANSACTION_TYPE_PAYMENT              TransactionType = 4
)

type TransactionType int

type Account struct {
	ID            int    `json:"id"`
	AccountNumber string `json:"account_number,omitempty"`
}

type Transaction struct {
	ID        int             `json:"id"`
	AccountID int             `json:"account_id,omitempty"`
	Type      TransactionType `json:"operation_type,omitempty"`
	Amount    float64         `json:"amount,omitempty"`
	Date      string          `json:"date,omitempty"`
}
