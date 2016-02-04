package httpclient

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
)

// TestGet tests the Get function
func TestGet(t *testing.T) {
	Configure(t.Logf)

	testData := "this is some test data"

	ts := httptest.NewServer(http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			fmt.Fprint(w, testData)
		}))
	defer ts.Close()

	httpClient, err := NewHTTPClient(ts.URL)
	if err != nil {
		t.Fatal("Could not create http client: ", err)
	}

	data, err := httpClient.Get("/")
	if err != nil {
		t.Fatal("Could not get from server: ", err)
	}

	if string(data) != testData {
		t.Fatalf("Response does not match: got %s, should be: %s", data, testData)
	}
}
