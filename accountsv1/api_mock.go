package accountsv1

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
)

// A map of accounts stored in the mocked API
var mockAccounts = make(map[string]AccountData)

type mockErrorResponse struct {
	Message interface{} `json:"error_message"`
}

// Mock of POST /
func mockPostHandler(rw http.ResponseWriter, req *http.Request) bool {
	if req.Method != "POST" {
		return false
	}

	account := requestBody{}
	if err := json.NewDecoder(req.Body).Decode(&account); err != nil {
		panic(err)
	}

	if account.Data.Attributes == nil {
		rw.WriteHeader(http.StatusBadRequest)
		response := mockErrorResponse{
			Message: "validation failure list:\nvalidation failure list:\nattributes in body is required",
		}

		responseJson, err := json.Marshal(response)
		if err != nil {
			panic(err)
		}

		io.WriteString(rw, string(responseJson))
		return true
	}

	accountId := account.Data.ID

	mockAccounts[accountId] = account.Data.hydrate()

	response := successResponse{
		Data: *account.Data,
	}

	responseJson, err := json.Marshal(response)
	if err != nil {
		panic(err)
	}

	io.WriteString(rw, string(responseJson))

	return true
}

// Mock of DELETE /<uuid>
func mockDeleteHandler(rw http.ResponseWriter, req *http.Request) bool {
	if req.Method != "DELETE" {
		return false
	}

	accountId := strings.Split(req.URL.Path, "/")[4]

	if _, found := mockAccounts[accountId]; found {
		rw.WriteHeader(http.StatusNoContent)
		delete(mockAccounts, accountId)
	} else {
		rw.WriteHeader(http.StatusNotFound)
	}

	return true
}

// Mock of GET /<uuid>
func mockGetHandler(rw http.ResponseWriter, req *http.Request) bool {
	if req.Method != "GET" {
		return false
	}

	accountId := strings.Split(req.URL.Path, "/")[4]

	var response interface{}
	if _, found := mockAccounts[accountId]; found {
		response = successResponse{
			Data: mockAccounts[accountId],
		}
	} else {
		rw.WriteHeader(http.StatusNotFound)
		response = Error{
			Message: fmt.Sprintf("record %s does not exist", accountId),
		}
	}

	responseJson, err := json.Marshal(response)
	if err != nil {
		panic(err)
	}

	io.WriteString(rw, string(responseJson))

	return true
}

// Take a request and route it to the appropriate mock handler using the
// chain of responsibility pattern
func mockApiHandler(rw http.ResponseWriter, req *http.Request) {
	handlers := []func(http.ResponseWriter, *http.Request) bool{
		mockPostHandler,
		mockDeleteHandler,
		mockGetHandler,
	}

	handled := false

	for _, handler := range handlers {
		if handler(rw, req) {
			handled = true
			break
		}
	}

	if !handled {
		panic(fmt.Sprintf("Mock request not handled: %s %s", req.Method, req.URL.Path))
	}
}

// Fill an account with default data matching that which the real api service
// would do
func (account *AccountData) hydrate() AccountData {
	// Version is always returned as 0 for new accounts
	var version int64 = 0
	account.Version = &version
	return *account
}
