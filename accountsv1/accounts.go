package accountsv1

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

// Client config defaults
const (
	DefaultApiUri   = "https://accountapi:8080"
	DefaultEndpoint = "/v1/organisation/accounts"
)

type Client struct {
	HttpClient http.Client
	BaseUrl    string
	EndPoint   string
}

func (c *Client) baseUrl() string {
	if c.BaseUrl == "" {
		c.BaseUrl = DefaultApiUri
	}
	return c.BaseUrl
}

func (c *Client) endpoint() string {
	if c.EndPoint == "" {
		c.EndPoint = DefaultEndpoint
	}
	return c.EndPoint
}

func checkForError(res *http.Response) error {
	// Exclude successful status codes
	if res.StatusCode >= http.StatusOK && res.StatusCode < http.StatusBadRequest {
		return nil
	}

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return err
	}

	apiError := Error{}

	// We can ignore the error here, if we can't parse the body it's available to
	// the user via `HttpResponse`.
	json.Unmarshal(body, &apiError)

	// Include the raw response
	apiError.HttpResponse = res
	return &apiError
}

func (c *Client) sendRequest(req *http.Request, v interface{}) error {
	req.Header.Set("Content-Type", "application/json; charset=utf-8")
	req.Header.Set("Accept", "application/json; charset=utf-8")

	res, err := c.HttpClient.Do(req)
	if err != nil {
		return err
	}

	defer res.Body.Close()

	if err := checkForError(res); err != nil {
		return err
	}

	if v != nil {
		result := successResponse{
			Data: v,
		}

		body, err := io.ReadAll(res.Body)
		if err != nil {
			return err
		}

		if err := json.Unmarshal(body, &result); err != nil {
			return err
		}
	}

	return nil
}

// Create a new account
func (c *Client) Create(ctx context.Context, account AccountData) (*AccountData, error) {
	payload, err := json.Marshal(requestBody{Data: &account})
	if err != nil {
		return nil, err
	}
	req, err := http.NewRequest(
		"POST",
		fmt.Sprintf("%s%s", c.baseUrl(), c.endpoint()),
		bytes.NewReader(payload),
	)
	if err != nil {
		return nil, err
	}

	req = req.WithContext(ctx)

	res := AccountData{}
	if err := c.sendRequest(req, &res); err != nil {
		return nil, err
	}

	return &res, nil
}

// Delete an account
func (c *Client) Delete(ctx context.Context, accountUuid string, version int64) error {
	req, err := http.NewRequest(
		"DELETE",
		fmt.Sprintf(
			"%s%s/%s?version=%d",
			c.baseUrl(),
			c.endpoint(),
			accountUuid,
			version,
		),
		nil,
	)
	if err != nil {
		return err
	}

	req = req.WithContext(ctx)
	if err := c.sendRequest(req, nil); err != nil {
		return err
	}

	return nil
}

// Fetch an account
func (c *Client) Fetch(ctx context.Context, accountUuid string) (*AccountData, error) {
	req, err := http.NewRequest(
		"GET",
		fmt.Sprintf(
			"%s%s/%s",
			c.baseUrl(),
			c.endpoint(),
			accountUuid,
		),
		nil,
	)
	if err != nil {
		return nil, err
	}

	req = req.WithContext(ctx)

	res := AccountData{}
	if err := c.sendRequest(req, &res); err != nil {
		return nil, err
	}

	return &res, nil
}
