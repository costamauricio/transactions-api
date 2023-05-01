package handlers

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/costamauricio/transactions-api/internal/models"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
)

type newAccountPayload struct {
	DocumentNumber string `json:"document_number"`
}

// Implements the render.Binder interface to validate the payload fields
func (account *newAccountPayload) Bind(r *http.Request) error {
	if len(account.DocumentNumber) == 0 {
		return errors.New("field 'document_number' is required and shouldn't be empty")
	}

	return nil
}

// Handler for the account creation
func (handler *ApplicationHandlers) CreateAccountHandler(w http.ResponseWriter, r *http.Request) {
	account := &newAccountPayload{}

	err := parseRequestBody(r, account)
	if err != nil {
		render.Render(w, r, newErrorResponse(
			err,
			http.StatusBadRequest,
		))
		return
	}

	id, err := handler.AccountDAO.NewAccount(account.DocumentNumber)

	if err != nil {
		handler.Logger.Errorf("Error creating account on database: %s", err)
		render.Render(w, r, newErrorResponse(
			errors.New("failed to create account"),
			http.StatusInternalServerError,
		))
		return
	}

	render.Render(w, r, newAccountResponse(
		&models.Account{
			ID: int(id),
		},
		http.StatusCreated,
	))
}

// Handler to retrieve the account
func (handler *ApplicationHandlers) GetAccountHandler(w http.ResponseWriter, r *http.Request) {
	idParam := chi.URLParam(r, "AccountId")

	accountId, err := strconv.ParseInt(idParam, 10, 0)
	if err != nil {
		handler.Logger.Errorf("Error converting to int: %s", idParam)
		render.Render(w, r, newErrorResponse(
			errors.New("failed to identify account"),
			http.StatusInternalServerError,
		))
		return
	}

	account, err := handler.AccountDAO.GetAccount(int(accountId))
	if err != nil {
		handler.Logger.Errorf("Failed to retrieve the account from the database: %s", err)
		render.Render(w, r, newErrorResponse(
			errors.New("failed to retrieve the account"),
			http.StatusInternalServerError,
		))
		return
	}

	if account == nil {
		render.Render(w, r, newErrorResponse(
			errors.New("account not found"),
			http.StatusNotFound,
		))
		return
	}

	render.Render(w, r, newAccountResponse(account, http.StatusOK))
}
