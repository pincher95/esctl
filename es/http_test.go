package es

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/pincher95/esctl/internal/client"
	eserrors "github.com/pincher95/esctl/internal/errors"
	"github.com/pincher95/esctl/shared"
)

func withTestClient(t *testing.T, handler http.HandlerFunc) {
	t.Helper()
	oldClient := shared.Client
	server := httptest.NewServer(handler)
	shared.SetClient(client.NewClient(&client.Config{BaseURL: server.URL}))
	t.Cleanup(func() {
		server.Close()
		shared.SetClient(oldClient)
	})
}

func TestHttpRequestAcceptsAny2xx(t *testing.T) {
	withTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		w.Write([]byte(`{"ok":true}`))
	})

	var target map[string]any
	if err := httpRequest(context.Background(), "POST", "created", map[string]string{"a": "b"}, &target); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestHttpRequestReturnsESError(t *testing.T) {
	withTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(`{"error":{"type":"illegal_argument_exception","reason":"bad request","root_cause":[{"type":"illegal_argument_exception","reason":"bad request"}]},"status":400}`))
	})

	var target map[string]any
	err := httpRequest(context.Background(), "GET", "bad", nil, &target)
	if err == nil {
		t.Fatalf("expected error")
	}

	esErr, ok := err.(*eserrors.ESError)
	if !ok {
		t.Fatalf("expected ESError, got %T", err)
	}
	if esErr.StatusCode != http.StatusBadRequest {
		t.Fatalf("expected status 400, got %d", esErr.StatusCode)
	}
	if esErr.Type != "illegal_argument_exception" {
		t.Fatalf("expected error type, got %q", esErr.Type)
	}
	if esErr.Reason != "bad request" {
		t.Fatalf("expected reason, got %q", esErr.Reason)
	}
}

func TestHttpRequestNonJSONBody(t *testing.T) {
	withTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("boom"))
	})

	var target map[string]any
	err := httpRequest(context.Background(), "GET", "error", nil, &target)
	if err == nil {
		t.Fatalf("expected error")
	}

	esErr, ok := err.(*eserrors.ESError)
	if !ok {
		t.Fatalf("expected ESError, got %T", err)
	}
	if esErr.Message == "" {
		t.Fatalf("expected raw message to be set")
	}
}
