package handlers_test

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/costamauricio/transactions-api/internal/handlers-test"
	"github.com/costamauricio/transactions-api/internal/models"
)

func TestCreateTransactionWithSuccess(t *testing.T) {
	mockedApi := handlersTest.NewMockedServer()

	mockedBody := strings.NewReader(`{"account_id": 1, "operation_type": 4, "amount": 50.2}`)
	request := httptest.NewRequest("POST", "/transactions", mockedBody)
	request.Header.Set("Content-Type", "application/json")

	mockedAccount := &models.Account{
		ID:             1,
		DocumentNumber: "123",
	}
	mockedApi.MockedAccountDAO.On("GetAccount", 1).Return(mockedAccount, nil)
	mockedApi.MockedTransactionDAO.On("NewTransaction", 1, models.TRANSACTION_TYPE_PAYMENT, 50.2).Return(7, nil)

	response := mockedApi.ExecuteRequest(request)

	if response.Code != http.StatusCreated {
		t.Errorf("Expected response status code %d, got %d", http.StatusCreated, response.Code)
	}

	var body map[string]int

	if err := json.Unmarshal(response.Body.Bytes(), &body); err != nil {
		t.Fatalf("Unexpected error parsing response: %s", err)
	}

	if body["id"] != 7 {
		t.Errorf(`Expected response id to be 7, received %v`, body["id"])
	}
}

func TestCreateTransactionHandleDatabaseErrorsAndAccountNotFound(t *testing.T) {
	mappedScenarions := map[int]int{
		0: http.StatusBadRequest, // Tests for payload validation error
		1: http.StatusInternalServerError,
		2: http.StatusInternalServerError,
		3: http.StatusBadRequest, // Account Not found
	}

	mockedApi := handlersTest.NewMockedServer()
	mockedApi.MockedAccountDAO.On("GetAccount", 1).Return(nil, errors.New("database error"))
	mockedAccount := &models.Account{
		ID:             2,
		DocumentNumber: "123",
	}
	mockedApi.MockedAccountDAO.On("GetAccount", 2).Return(mockedAccount, nil)
	mockedApi.MockedAccountDAO.On("GetAccount", 3).Return(nil, nil)
	mockedApi.MockedTransactionDAO.On("NewTransaction", 2, models.TRANSACTION_TYPE_PAYMENT, 50.2).Return(0, errors.New("database error"))

	for accountId, statusCode := range mappedScenarions {
		mockedBody := strings.NewReader(fmt.Sprintf(`{"account_id": %d, "operation_type": 4, "amount": 50.2}`, accountId))

		request := httptest.NewRequest("POST", "/transactions", mockedBody)
		request.Header.Set("Content-Type", "application/json")

		response := mockedApi.ExecuteRequest(request)

		if response.Code != statusCode {
			t.Errorf("Expected response status code %d, got %d", statusCode, response.Code)
		}

	}
}

func TestCreateTransactionHandleTheAmountValueCorrectly(t *testing.T) {
	mappedScenarions := map[models.TransactionType]struct {
		Amount         float64
		ExpectedAmount float64
	}{
		models.TRANSACTION_TYPE_PUSCHASE_CASH:        {20.2, -20.2},
		models.TRANSACTION_TYPE_PUSCHASE_INSTALLMENT: {23.5, -23.5},
		models.TRANSACTION_TYPE_WITHDRAW:             {13, -13},
		models.TRANSACTION_TYPE_PAYMENT:              {234.4, 234.4},
	}

	mockedApi := handlersTest.NewMockedServer()
	mockedAccount := &models.Account{
		ID:             1,
		DocumentNumber: "123",
	}
	mockedApi.MockedAccountDAO.On("GetAccount", 1).Return(mockedAccount, nil)
	mockedApi.MockedTransactionDAO.On("NewTransaction", 2, models.TRANSACTION_TYPE_PAYMENT, 50.2).Return(0, errors.New("database error"))

	for operationType, params := range mappedScenarions {
		mockedBody := strings.NewReader(fmt.Sprintf(`{"account_id": 1, "operation_type": %d, "amount": %f}`, operationType, params.Amount))

		request := httptest.NewRequest("POST", "/transactions", mockedBody)
		request.Header.Set("Content-Type", "application/json")

		mockedApi.MockedTransactionDAO.On("NewTransaction", 1, operationType, params.ExpectedAmount).Return(2, nil)

		response := mockedApi.ExecuteRequest(request)

		if response.Code != http.StatusCreated {
			t.Errorf("Expected response status code %d, got %d", http.StatusCreated, response.Code)
		}
	}
}
