package handlers

import (
	"errors"
	"net/http"

	"github.com/costamauricio/transactions-api/internal/models"
	"github.com/go-chi/render"
)

type newTransactionPayload struct {
	AccountId     int                    `json:"account_id"`
	OperationType models.TransactionType `json:"operation_type"`
	Amount        float64                `json:"amount"`
}

// Implements the render.Binder interface to validate the payload fields
func (transaction *newTransactionPayload) Bind(r *http.Request) error {
	if transaction.AccountId == 0 {
		return errors.New("field 'account_id' shouldn't be empty")
	}

	if !transaction.OperationType.IsValid() {
		return errors.New("invalid value for 'operation_type' field")
	}

	if transaction.Amount <= 0 {
		return errors.New("field 'amount' is required and must be a positive number")
	}

	return nil
}

// Handler for the transaction
func (handler *ApplicationHandlers) CreateTransactionHandler(w http.ResponseWriter, r *http.Request) {
	transaction := &newTransactionPayload{}

	err := parseRequestBody(r, transaction)
	if err != nil {
		render.Render(w, r, newErrorResponse(
			err,
			http.StatusBadRequest,
		))
		return
	}

	// Here we get the account from the database to validate it is a real account
	// in this api context it should be ok, but that adds more latency to the endpoint
	// since it will be more sequential database operations occuring without any optimization
	account, err := handler.AccountDAO.GetAccount(transaction.AccountId)
	if err != nil {
		handler.Logger.Errorf("Failed to retrieve the account from the database: %s", err)
		render.Render(w, r, newErrorResponse(
			errors.New("failed to retrieve the account"),
			http.StatusInternalServerError,
		))
		return
	}

	// If the account is nil, then the account was not found for the account id
	if account == nil {
		render.Render(w, r, newErrorResponse(
			errors.New("account not found for the provided 'account_id'"),
			http.StatusBadRequest,
		))
		return
	}

	// If the transaction is different from payment eg.(purchase, withdraw), the amount must be stored with a negative value
	if transaction.OperationType != models.TRANSACTION_TYPE_PAYMENT {
		transaction.Amount *= -1
	}

	id, err := handler.TransactioDAO.NewTransaction(transaction.AccountId, transaction.OperationType, transaction.Amount)
	if err != nil {
		handler.Logger.Errorf("Error creating account on database: %s", err)
		render.Render(w, r, newErrorResponse(
			errors.New("failed to create transaction"),
			http.StatusInternalServerError,
		))
		return
	}

	render.Render(w, r, newTransactionResponse(
		&models.Transaction{
			ID: int(id),
		},
		http.StatusCreated,
	))
}
