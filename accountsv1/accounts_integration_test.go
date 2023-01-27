package accountsv1

import (
	"flag"
	"testing"
)

var offline = flag.Bool(
	"offline",
	false,
	"don't depend on external services while testing",
)

var apiUrl = flag.String(
	"base-url",
	"http://accountapi:8080",
	"apiUrl to run integration tests against",
)

// Run the tests against a locally running service if the offline flag wasn't
// given.
func TestRealApi(t *testing.T) {
	if !*offline {
		testAgainstUrl(t, *apiUrl)
	}
}
