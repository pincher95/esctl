package errors

import (
	"context"
	"encoding/json"
	stderrors "errors"
	"fmt"
	"net"
	"net/http"
	"os"
)

// ESError represents a structured Elasticsearch/OpenSearch error
type ESError struct {
	StatusCode int
	Message    string
	Type       string
	Reason     string
	RootCause  []RootCause
}

type RootCause struct {
	Type   string `json:"type"`
	Reason string `json:"reason"`
}

type errorResponse struct {
	Error struct {
		Type      string      `json:"type"`
		Reason    string      `json:"reason"`
		RootCause []RootCause `json:"root_cause"`
	} `json:"error"`
	Status int `json:"status"`
}

func (e *ESError) Error() string {
	if e.Reason != "" {
		return fmt.Sprintf("%s (status: %d): %s", e.Type, e.StatusCode, e.Reason)
	}
	return fmt.Sprintf("elasticsearch error (status: %d): %s", e.StatusCode, e.Message)
}

// NewESError creates a structured error from HTTP response
func NewESError(statusCode int, body []byte) *ESError {
	err := &ESError{
		StatusCode: statusCode,
		Message:    string(body),
	}

	var parsed errorResponse
	if unmarshalErr := json.Unmarshal(body, &parsed); unmarshalErr == nil {
		if parsed.Status != 0 {
			err.StatusCode = parsed.Status
		}
		if parsed.Error.Type != "" {
			err.Type = parsed.Error.Type
		}
		if parsed.Error.Reason != "" {
			err.Reason = parsed.Error.Reason
		}
		if len(parsed.Error.RootCause) > 0 {
			err.RootCause = parsed.Error.RootCause
		}
	}

	return err
}

// IsTimeout reports whether err was caused by the client giving up waiting (the
// request exceeded the client/--timeout deadline) rather than by a response from
// the server. Such a request may still be running on the cluster.
func IsTimeout(err error) bool {
	if err == nil {
		return false
	}
	if stderrors.Is(err, context.DeadlineExceeded) || stderrors.Is(err, os.ErrDeadlineExceeded) {
		return true
	}
	var netErr net.Error
	if stderrors.As(err, &netErr) {
		return netErr.Timeout()
	}
	return false
}

// IsAuthError checks if the error is authentication-related
func IsAuthError(err error) bool {
	if esErr, ok := err.(*ESError); ok {
		return esErr.StatusCode == http.StatusUnauthorized ||
			esErr.StatusCode == http.StatusForbidden
	}
	return false
}

// IsNotFoundError checks if the error is a 404
func IsNotFoundError(err error) bool {
	if esErr, ok := err.(*ESError); ok {
		return esErr.StatusCode == http.StatusNotFound
	}
	return false
}

// IsConflictError checks if the error is a conflict (version mismatch, etc.)
func IsConflictError(err error) bool {
	if esErr, ok := err.(*ESError); ok {
		return esErr.StatusCode == http.StatusConflict
	}
	return false
}

// IsBadRequestError checks if the error is a bad request
func IsBadRequestError(err error) bool {
	if esErr, ok := err.(*ESError); ok {
		return esErr.StatusCode == http.StatusBadRequest
	}
	return false
}
