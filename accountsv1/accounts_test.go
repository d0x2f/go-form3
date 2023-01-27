package accountsv1

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

// Test creating a new account
func testCreate(t *testing.T, c Client, testCase TestCase) {
	ctx := context.Background()
	accountId := testCase.InputAccount.ID

	got, err := c.Create(ctx, *testCase.InputAccount)
	assert.NoError(t, err)
	assert.Equal(t, testCase.ExpectedAccount, got)

	// Clean-up
	err = c.Delete(ctx, accountId, *got.Version)
	assert.NoError(t, err)
}

// Test account deletion
func testDelete(t *testing.T, c Client, testCase TestCase) {
	ctx := context.Background()
	accountId := testCase.InputAccount.ID

	account, err := c.Create(ctx, *testCase.InputAccount)
	assert.NoError(t, err)

	err = c.Delete(ctx, accountId, *account.Version)
	assert.NoError(t, err)

	_, err = c.Get(ctx, accountId)
	assert.Equal(
		t,
		fmt.Sprintf("404 Not Found - record %s does not exist", accountId),
		err.Error(),
	)
}

// Test fetching an account
func testGet(t *testing.T, c Client, testCase TestCase) {
	ctx := context.Background()
	accountId := testCase.InputAccount.ID

	_, err := c.Create(ctx, *testCase.InputAccount)
	assert.NoError(t, err)

	// Fetch the account
	got, err := c.Get(ctx, accountId)
	assert.NoError(t, err)
	assert.Equal(t, testCase.ExpectedAccount, got)

	// Clean-up
	err = c.Delete(ctx, accountId, *got.Version)
	assert.NoError(t, err)
}

type TestCase struct {
	InputAccount    *AccountData `json:"input"`
	ExpectedAccount *AccountData `json:"expected"`
}

func loadJsonFixture(name string) TestCase {
	testCaseJson, err := os.ReadFile(fmt.Sprintf("./fixtures/%s.json", name))
	if err != nil {
		panic(err)
	}

	testCase := TestCase{}
	if err := json.Unmarshal(testCaseJson, &testCase); err != nil {
		panic(err)
	}

	return testCase
}

func testCreateError(t *testing.T, c Client) {
	t.Run("CreateError", func(t *testing.T) {
		ctx := context.Background()
		emptyAccount := AccountData{
			ID:             "3de43a70-1e5a-4e03-942d-96fb6f345ffa",
			OrganisationID: "654a4fe5-766f-452c-b7e3-8c7c5030668b",
			Type:           "accounts",
		}
		_, err := c.Create(ctx, emptyAccount)
		assert.EqualError(
			t,
			err,
			"400 Bad Request - validation failure list:\nvalidation failure list:\nattributes in body is required",
		)
	})
}

func testGetError(t *testing.T, c Client) {
	t.Run("GetError", func(t *testing.T) {
		ctx := context.Background()
		_, err := c.Get(ctx, "00000000-0000-0000-0000-deadbeefcafe")
		assert.EqualError(
			t,
			err,
			"404 Not Found - record 00000000-0000-0000-0000-deadbeefcafe does not exist",
		)
	})
}

func testDeleteError(t *testing.T, c Client) {
	t.Run("DeleteError", func(t *testing.T) {
		ctx := context.Background()
		err := c.Delete(ctx, "00000000-0000-0000-0000-deadbeefcafe", 0)
		assert.EqualError(
			t,
			err,
			"404 Not Found",
		)
	})
}

// Run each of the tests against the given API URL
func testAgainstUrl(t *testing.T, apiUrl string) {
	c := Client{
		BaseUrl: apiUrl,
	}

	fixtures := []string{
		"full_account",
		"minimal_account",
		"mutated_version",
		"lots_of_names",
	}

	var tests = []struct {
		name string
		test func(*testing.T, Client, TestCase)
	}{
		{name: "New", test: testCreate},
		{name: "Delete", test: testDelete},
		{name: "Get", test: testGet},
	}

	// For each fixture, run each test
	for _, fixture := range fixtures {
		t.Run(fixture, func(t *testing.T) {
			testCase := loadJsonFixture(fixture)
			for _, spec := range tests {
				t.Run(spec.name, func(t *testing.T) { spec.test(t, c, testCase) })
			}
		})
	}

	// Extra tests
	testCreateError(t, c)
	testGetError(t, c)
	testDeleteError(t, c)
}

// Run the tests against a mocked HTTP service
func TestMocked(t *testing.T) {
	hs := httptest.NewServer(http.HandlerFunc(mockApiHandler))
	defer hs.Close()
	testAgainstUrl(t, hs.URL)
}
