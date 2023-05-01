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

func TestCreateAccountWithSuccess(t *testing.T) {
	mockedApi := handlersTest.NewMockedServer()

	mockedBody := strings.NewReader(`{"document_number": "123"}`)
	request := httptest.NewRequest("POST", "/accounts", mockedBody)
	request.Header.Set("Content-Type", "application/json")

	mockedApi.MockedAccountDAO.On("NewAccount", "123").Return(2, nil)

	response := mockedApi.ExecuteRequest(request)

	if response.Code != http.StatusCreated {
		t.Errorf("Expected response status code %d, got %d", http.StatusCreated, response.Code)
	}

	var body map[string]int

	if err := json.Unmarshal(response.Body.Bytes(), &body); err != nil {
		t.Fatalf("Unexpected error parsing response: %s", err)
	}

	if body["id"] != 2 {
		t.Errorf(`Expected response id to be 2, received %v`, body["id"])
	}
}

func TestCreateAccountHandleErrorAndMalformedPayload(t *testing.T) {
	mockedScenarions := map[string]int{
		`"document_number": "123"`: http.StatusInternalServerError,
		`"other": "123"`:           http.StatusBadRequest,
	}

	mockedApi := handlersTest.NewMockedServer()

	mockedApi.MockedAccountDAO.On("NewAccount", "123").Return(0, errors.New("error on database"))

	for payload, expectedStatus := range mockedScenarions {
		mockedBody := strings.NewReader(fmt.Sprintf("{%s}", payload))

		request := httptest.NewRequest("POST", "/accounts", mockedBody)
		request.Header.Set("Content-Type", "application/json")

		response := mockedApi.ExecuteRequest(request)

		if response.Code != expectedStatus {
			t.Errorf("Expected response status code %d, got %d", expectedStatus, response.Code)
		}
	}
}

func TestGetAccountWithSuccess(t *testing.T) {
	mockedApi := handlersTest.NewMockedServer()

	request := httptest.NewRequest("GET", "/accounts/1", nil)

	mockedAccount := &models.Account{
		ID:             1,
		DocumentNumber: "123",
	}
	mockedApi.MockedAccountDAO.On("GetAccount", 1).Return(mockedAccount, nil)

	response := mockedApi.ExecuteRequest(request)

	if response.Code != http.StatusOK {
		t.Errorf("Expected response status code %d, got %d", http.StatusOK, response.Code)
	}
	var responseAccount models.Account

	if err := json.Unmarshal(response.Body.Bytes(), &responseAccount); err != nil {
		t.Fatalf("Unexpected error parsing response: %s", err)
	}

	if responseAccount.ID != mockedAccount.ID {
		t.Errorf(`Expected account to be %v, received %s`, mockedAccount, response.Body.String())
	}
}

func TestGetAccountHandleDatabaseErrorsAndAccountNotFound(t *testing.T) {
	mappedScenarions := map[string]int{
		"1": http.StatusInternalServerError,
		"2": http.StatusNotFound,

		"112312312312312312312312312312312312": http.StatusInternalServerError, // Tests for integer conversion error
	}
	mockedApi := handlersTest.NewMockedServer()

	mockedApi.MockedAccountDAO.On("GetAccount", 1).Return(nil, errors.New("database error"))
	mockedApi.MockedAccountDAO.On("GetAccount", 2).Return(nil, nil)

	for accountId, expectedStatus := range mappedScenarions {
		request := httptest.NewRequest("GET", fmt.Sprintf("/accounts/%s", accountId), nil)

		response := mockedApi.ExecuteRequest(request)

		if response.Code != expectedStatus {
			t.Errorf("Expected response status code %d, got %d", expectedStatus, response.Code)
		}
	}
}
