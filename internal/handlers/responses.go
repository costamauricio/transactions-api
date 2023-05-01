package handlers

import (
	"net/http"

	"github.com/costamauricio/transactions-api/internal/models"
	"github.com/go-chi/render"
)

// Base response that will be returned without any struct embbeded to it
type BaseResponse struct {
	HttpStatusCode int `json:"-"`
}

// Implements the render.Renderer interface on the BaseResponse
func (resp *BaseResponse) Render(w http.ResponseWriter, r *http.Request) error {
	// Set the returned StatusCode
	render.Status(r, resp.HttpStatusCode)
	return nil
}

// Response with a single models.Account
type accountResponse struct {
	*models.Account
	*BaseResponse
}

// Build the response with a single models.Account
func newAccountResponse(account *models.Account, status int) render.Renderer {
	return &accountResponse{
		account,
		&BaseResponse{
			HttpStatusCode: status,
		},
	}
}

// Response with a single models.Transaction
type transactionResponse struct {
	*models.Transaction
	*BaseResponse
}

// Build the response with a single model.Transaction
func newTransactionResponse(transaction *models.Transaction, status int) render.Renderer {
	return &transactionResponse{
		transaction,
		&BaseResponse{
			HttpStatusCode: status,
		},
	}
}
