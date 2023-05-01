package handlersTest

import (
	"net/http"
	"net/http/httptest"

	"github.com/costamauricio/transactions-api/internal/handlers"
	"github.com/costamauricio/transactions-api/internal/server"
)

type MockedServer struct {
	MockedAccountDAO     *AccountDAOMock
	MockedTransactionDAO *TransactionDAOMock
	server               *server.Server
}

func (mock *MockedServer) ExecuteRequest(req *http.Request) *httptest.ResponseRecorder {
	response := httptest.NewRecorder()
	mock.server.Router.ServeHTTP(response, req)

	return response
}

// Create a new ApplicationHandlers with the mocks instead of the implementations
// Create a new Server and attach the Mocked dependencies to it
// and return a new MockedServer that can be used to trigger requests
func NewMockedServer() *MockedServer {
	api := server.New("80")

	accountMock := &AccountDAOMock{}
	transactionMock := &TransactionDAOMock{}

	mockedHandlers := &handlers.ApplicationHandlers{
		Logger:        &LoggerMock{},
		AccountDAO:    accountMock,
		TransactioDAO: transactionMock,
	}

	api.AttachHandlers(mockedHandlers)

	return &MockedServer{
		MockedAccountDAO:     accountMock,
		MockedTransactionDAO: transactionMock,
		server:               api,
	}
}
