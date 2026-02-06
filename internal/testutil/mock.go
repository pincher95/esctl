package testutil

import (
	"net/http"
	"net/http/httptest"

	"github.com/pincher95/esctl/internal/client"
)

// NewMockServer spins up an httptest.Server that returns the provided JSON with status 200
// and asserts the requested path (no query string).
// It returns the server pointer and an ESClient configured to talk to it.
func NewMockServer(responseJSON string, expectedPath string) (*httptest.Server, client.ESClient) {
	return NewMockServerWithStatus(http.StatusOK, responseJSON, expectedPath)
}

// NewMockServerWithStatus spins up an httptest.Server that returns the provided JSON with the specified status code
// and asserts the requested path (no query string).
// It returns the server pointer and an ESClient configured to talk to it.
func NewMockServerWithStatus(statusCode int, responseJSON string, expectedPath string) (*httptest.Server, client.ESClient) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != expectedPath {
			http.Error(w, "unexpected path", http.StatusBadRequest)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(statusCode)
		w.Write([]byte(responseJSON))
	}))

	cfg := &client.Config{BaseURL: srv.URL}
	cli := client.NewClient(cfg)
	return srv, cli
}
