package es

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"strings"

	"github.com/go-resty/resty/v2"
	"github.com/pincher95/esctl/shared"
)

type EsError struct {
	Error struct {
		Type   string `json:"type"`
		Reason string `json:"reason"`
	} `json:"error"`
	Status int `json:"status"`
}

func debugLog(format string, args ...any) {
	if shared.Debug {
		fmt.Fprintf(os.Stderr, "DEBUG: "+format+"\n", args...)
	}
}

func httpRequest(ctx context.Context, method, endpoint string, body, target any, expectedStatusCode int) error {
	// Build URL relative to base configured in shared.Client (which already has base URL set).
	// If shared.Client is nil, return error.
	if shared.Client == nil {
		return errors.New("shared.Client not initialised")
	}

	// Ensure endpoint has no leading slash to avoid double slashes with baseURL.
	endpoint = strings.TrimPrefix(endpoint, "/")

	req := shared.Client.R().
		SetContext(ctx).
		SetHeader("Content-Type", "application/json").
		SetResult(target)

	if body != nil {
		req = req.SetBody(body)
	}

	var resp *resty.Response
	var err error

	switch strings.ToUpper(method) {
	case "GET":
		resp, err = req.Get(endpoint)
	case "POST":
		resp, err = req.Post(endpoint)
	case "PUT":
		resp, err = req.Put(endpoint)
	case "DELETE":
		resp, err = req.Delete(endpoint)
	default:
		return fmt.Errorf("unsupported method %s", method)
	}

	if err != nil {
		return err
	}

	if resp.StatusCode() != expectedStatusCode {
		var esError EsError
		if err := json.Unmarshal(resp.Body(), &esError); err != nil {
			return fmt.Errorf("unexpected http status: %d", resp.StatusCode())
		}
		return errors.New(esError.Error.Reason)
	}

	return nil
}

func getJSONResponse(ctx context.Context, endpoint string, target any) error {
	return httpRequest(ctx, "GET", endpoint, nil, target, 200)
}

// func getJSONResponseWithBody(endpoint string, target any, body any) error {
// 	return httpRequest(http.MethodGet, endpoint, body, target, http.StatusOK)
// }

func postJSONResponseWithBody(ctx context.Context, endpoint string, target any, body any) error {
	return httpRequest(ctx, "POST", endpoint, body, target, 200)
}

func postWithoutBody(ctx context.Context, endpoint string, target any) error {
	return httpRequest(ctx, "POST", endpoint, nil, target, 200)
}

func getNestedPath(field string, nestedPaths []string) (string, bool) {
	for _, path := range nestedPaths {
		if strings.HasPrefix(field, path) {
			return path, true
		}
	}
	return "", false
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}
