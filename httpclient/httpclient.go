// Package httpclient provides an HTTP client.
package httpclient

import (
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/giantswarm/yochu/fetchclient"
)

var vLogger = func(f string, v ...interface{}) {}

// HTTPClient is an HTTP client configured with an endpoint for all requests to use.
type HTTPClient struct {
	// endpoint is the base for all paths supplied by Get to use
	endpoint string
}

// Configure sets the logger for this package.
func Configure(vl func(f string, v ...interface{})) {
	vLogger = vl
}

// NewHTTPClient returns a HTTPClient with endpoint specified.
func NewHTTPClient(endpoint string) (fetchclient.FetchClient, error) {
	return &HTTPClient{
		endpoint: endpoint,
	}, nil
}

// Get appends the supplied path to the HTTPClient's endpoint,
// performs a GET request against the resulting URL,
// and returns the response body.
func (h *HTTPClient) Get(path string) ([]byte, error) {
	vLogger("  call HTTPClient.Get(endpoint, path): %v - %v", h.endpoint, path)

	url := fmt.Sprintf("%s/%s", h.endpoint, path)
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	return ioutil.ReadAll(resp.Body)
}
