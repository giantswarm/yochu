// Package fetchclient provides the FetchClient interface.
package fetchclient

// FetchClient is an interface for getting data.
type FetchClient interface {
	Get(string) ([]byte, error)
}
