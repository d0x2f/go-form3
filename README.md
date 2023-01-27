# Form3 Accounts Client Library

Me: Dylan McGannon

A Form3 accounts API client library for go.

## Exports

This library exports a `Client` struct with three methods to manage Form3 accounts, `Create`,
`Delete` & `Get` with the following signatures:

```go
type Client struct {
	HttpClient http.Client // Can be used to configure extra http details
	BaseUrl    string      // The base URL of the Form3 accounts API
	EndPoint   string      // The accounts endpoint
}
```
Sensible defaults will be used when no values are given.

```go
// Creates a new account
client.Create(
	ctx context.Context, // To enable standard golang context management
	account AccountData  // The account to create
) (
	*AccountData, // The created account returned from the API
	error         // An error object if one occurred
)

// Deletes an account with the given id and version
client.Delete(
	ctx context.Context, // To enable standard golang context management
	accountUuid string,  // The UUID of the account you wish to delete
	version int64        // Used to prevent concurrency clashes (presumably)
) error // An error object if one occurred

// Fetches an account with the given id
client.Get(
	ctx context.Context,
	accountUuid string
) (
	*AccountData, // The account returned from the API
	error         // An error object if one occurred
)
```

## Example Usage

```go
package main

import (
	"context"

	accountsv1 "github.com/d0x2f/go-form3/accountsv1"
)

func main() {
	ctx := context.Background()

	c := accountsv1.Client{
		BaseUrl: "http://localhost:8080",
	}

	// Define an account
	country := "AU"
	createAccountRequest := accountsv1.AccountData{
		OrganisationID: "3de43a70-1e5a-4e03-942d-96fb6f345ffa",
		ID:             "654a4fe5-766f-452c-b7e3-8c7c5030668b",
		Type:           "accounts",
		Attributes: &accountsv1.AccountAttributes{
			Name:    []string{"dylan"},
			Country: &country,
		},
	}

	// Create
	createdAccount, err := c.Create(ctx, createAccountRequest)
	if (err != nil) {
		// handle error
		// statusCode := err.HttpResponse.StatusCode
		// errorMessage := err.Message
	}

	// Fetch
	account, err := c.Get(ctx, "654a4fe5-766f-452c-b7e3-8c7c5030668b")
	if (err != nil) {
		// handle error
	}

	// Delete
	err = c.Delete(ctx, "654a4fe5-766f-452c-b7e3-8c7c5030668b", *account.Version)
	if (err != nil) {
		// handle error
	}
}
```

## Features

 - You can provide your own http.Client to configure requests how you wish.
 - Methods accept a context parameter to enable cancellable requests.
 - API generated errors include the raw HTTP response to allow you to perform deeper inspection.


## Testing

All tests can be run against a mocked api or against a running api service.
By default tests are run against both and are configured to run using the provided docker-compose.yml.
Running the integration tests can be disabled using the --offline flag like so:

```bash
$ go test --offline
```

You can also configure the API URL for integration tests to run against using the --base-url flag:

```bash
$ go test --base-url http://localhost:8080
```

Tests run automatically using a GitHub workflow.