package client

import (
	"net/http"
	"net/http/httptest"
	"sync/atomic"
	"testing"
	"time"
)

// newCountingServer returns a server that always fails with 500 and counts hits.
func newCountingServer(t *testing.T, hits *int32) *httptest.Server {
	t.Helper()
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		atomic.AddInt32(hits, 1)
		http.Error(w, "boom", http.StatusInternalServerError)
	}))
}

func newTestClient(t *testing.T, baseURL string) ESClient {
	t.Helper()
	cli, err := NewClient(&Config{
		BaseURL:       baseURL,
		RetryCount:    3,
		RetryWaitTime: time.Millisecond,
	})
	if err != nil {
		t.Fatalf("NewClient() error = %v", err)
	}
	return cli
}

// A mutating request must never be retried: re-sending a POST that the server may
// already have accepted can duplicate the operation (e.g. a snapshot restore).
func TestNoRetryOnMutatingRequest(t *testing.T) {
	var hits int32
	srv := newCountingServer(t, &hits)
	defer srv.Close()

	for _, method := range []string{"POST", "PUT", "DELETE"} {
		atomic.StoreInt32(&hits, 0)
		req := newTestClient(t, srv.URL).R()

		var err error
		switch method {
		case "POST":
			_, err = req.Post("/_snapshot/repo/snap/_restore")
		case "PUT":
			_, err = req.Put("/_snapshot/repo/snap")
		case "DELETE":
			_, err = req.Delete("/some-index")
		}
		if err != nil {
			t.Fatalf("%s: unexpected transport error: %v", method, err)
		}
		if got := atomic.LoadInt32(&hits); got != 1 {
			t.Errorf("%s was sent %d times; want exactly 1 (no retry on mutations)", method, got)
		}
	}
}

// A client-side timeout means the server is slow, not that the connection blipped.
// Re-sending re-runs the same expensive query and adds load, so it is not retried
// even for a read-only request.
func TestNoRetryOnTimeout(t *testing.T) {
	var hits int32
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		atomic.AddInt32(&hits, 1)
		time.Sleep(300 * time.Millisecond)
	}))
	defer srv.Close()

	cli, err := NewClient(&Config{
		BaseURL:       srv.URL,
		RetryCount:    3,
		RetryWaitTime: time.Millisecond,
		Timeout:       50 * time.Millisecond,
	})
	if err != nil {
		t.Fatalf("NewClient() error = %v", err)
	}

	if _, err := cli.R().Get("/_cat/snapshots/repo"); err == nil {
		t.Fatal("expected a timeout error, got nil")
	}
	if got := atomic.LoadInt32(&hits); got != 1 {
		t.Errorf("timed-out GET was sent %d times; want exactly 1 (no retry on timeout)", got)
	}
}

// Read-only requests are safe to repeat, so 5xx responses are still retried.
func TestRetryOnReadOnlyRequest(t *testing.T) {
	var hits int32
	srv := newCountingServer(t, &hits)
	defer srv.Close()

	if _, err := newTestClient(t, srv.URL).R().Get("/_cat/indices"); err != nil {
		t.Fatalf("unexpected transport error: %v", err)
	}
	// 1 initial attempt + 3 retries
	if got := atomic.LoadInt32(&hits); got != 4 {
		t.Errorf("GET was sent %d times; want 4 (1 attempt + 3 retries)", got)
	}
}
