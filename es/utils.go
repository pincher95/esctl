package es

import (
	"context"
	"errors"
	"fmt"
	"os"
	"slices"
	"strings"

	"github.com/go-resty/resty/v2"
	eserrors "github.com/pincher95/esctl/internal/errors"
	"github.com/pincher95/esctl/shared"
)

func debugLog(format string, args ...any) {
	if shared.Debug {
		fmt.Fprintf(os.Stderr, "DEBUG: "+format+"\n", args...)
	}
}

func httpRequest(ctx context.Context, method, endpoint string, body, target any, expectedStatusCodes ...int) error {
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

	// Emit debug output before the request is sent.
	debugLog("HTTP %s %s", strings.ToUpper(method), endpoint)

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

	// Emit debug output after receiving the response.
	debugLog("HTTP %s %s -> %d", strings.ToUpper(method), endpoint, resp.StatusCode())

	if !isExpectedStatus(resp.StatusCode(), expectedStatusCodes) {
		return eserrors.NewESError(resp.StatusCode(), resp.Body())
	}

	return nil
}

func GetJSONResponse(ctx context.Context, endpoint string, target any) error {
	return httpRequest(ctx, "GET", endpoint, nil, target)
}

// func getJSONResponseWithBody(endpoint string, target any, body any) error {
// 	return httpRequest(http.MethodGet, endpoint, body, target, http.StatusOK)
// }

func postJSONResponseWithBody(ctx context.Context, endpoint string, target any, body any) error {
	return httpRequest(ctx, "POST", endpoint, body, target)
}

func postWithoutBody(ctx context.Context, endpoint string, target any) error {
	return httpRequest(ctx, "POST", endpoint, nil, target)
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

func isExpectedStatus(code int, expected []int) bool {
	if len(expected) == 0 {
		return code >= 200 && code < 300
	}
	if slices.Contains(expected, code) {
		return true
	}
	return false
}
