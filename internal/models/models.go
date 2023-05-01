package models

const (
	TRANSACTION_TYPE_PUSCHASE_CASH        TransactionType = 1
	TRANSACTION_TYPE_PUSCHASE_INSTALLMENT TransactionType = 2
	TRANSACTION_TYPE_WITHDRAW             TransactionType = 3
	TRANSACTION_TYPE_PAYMENT              TransactionType = 4
)

// Define the transaction type
type TransactionType int

// Validate if the transaction type is between the defined types
func (t TransactionType) IsValid() bool {
	switch t {
	case TRANSACTION_TYPE_WITHDRAW,
		TRANSACTION_TYPE_PAYMENT,
		TRANSACTION_TYPE_PUSCHASE_CASH,
		TRANSACTION_TYPE_PUSCHASE_INSTALLMENT:
		return true
	}
	return false
}

type Account struct {
	ID             int    `json:"id"`
	DocumentNumber string `json:"document_number,omitempty"`
}

type Transaction struct {
	ID        int             `json:"id"`
	AccountID int             `json:"account_id,omitempty"`
	Type      TransactionType `json:"operation_type,omitempty"`
	Amount    float64         `json:"amount,omitempty"`
	Date      string          `json:"date,omitempty"`
}
