package script

import (
	"context"
	"fmt"
	"net/url"

	"github.com/pincher95/esctl/shared"
)

// Script represents a stored script
type Script struct {
	ID     string         `json:"id,omitempty"`
	Lang   string         `json:"lang"`
	Source string         `json:"source"`
	Params map[string]any `json:"params,omitempty"`
}

// ScriptResponse is the response when getting a script
type ScriptResponse struct {
	Found  bool   `json:"found"`
	ID     string `json:"_id"`
	Script Script `json:"script"`
}

// ListResponse contains all stored scripts
type ListResponse struct {
	Scripts map[string]ScriptResponse `json:"scripts"`
}

// List retrieves all stored scripts
func List(ctx context.Context) (map[string]ScriptResponse, error) {
	resp, err := shared.Client.R().
		SetContext(ctx).
		SetHeader("Content-Type", "application/json").
		SetResult(&map[string]ScriptResponse{}).
		Get("/_scripts")

	if err != nil {
		return nil, fmt.Errorf("failed to list scripts: %w", err)
	}

	if resp.StatusCode() != 200 {
		return nil, fmt.Errorf("error listing scripts: %s", resp.Status())
	}

	result := resp.Result().(*map[string]ScriptResponse)
	return *result, nil
}

// Get retrieves a specific stored script by ID
func Get(ctx context.Context, id string) (*ScriptResponse, error) {
	if id == "" {
		return nil, fmt.Errorf("script ID is required")
	}

	escapedID := url.PathEscape(id)

	resp, err := shared.Client.R().
		SetContext(ctx).
		SetHeader("Content-Type", "application/json").
		SetResult(&ScriptResponse{}).
		Get(fmt.Sprintf("/_scripts/%s", escapedID))

	if err != nil {
		return nil, fmt.Errorf("failed to get script: %w", err)
	}

	if resp.StatusCode() == 404 {
		return nil, fmt.Errorf("script '%s' not found", id)
	}

	if resp.StatusCode() != 200 {
		return nil, fmt.Errorf("error getting script: %s", resp.Status())
	}

	result := resp.Result().(*ScriptResponse)
	result.ID = id
	return result, nil
}

// Put creates or updates a stored script
func Put(ctx context.Context, id string, script Script) error {
	if id == "" {
		return fmt.Errorf("script ID is required")
	}

	if script.Lang == "" {
		return fmt.Errorf("script language is required")
	}

	if script.Source == "" {
		return fmt.Errorf("script source is required")
	}

	escapedID := url.PathEscape(id)

	body := map[string]any{
		"script": script,
	}

	resp, err := shared.Client.R().
		SetContext(ctx).
		SetHeader("Content-Type", "application/json").
		SetBody(body).
		Put(fmt.Sprintf("/_scripts/%s", escapedID))

	if err != nil {
		return fmt.Errorf("failed to put script: %w", err)
	}

	if resp.StatusCode() != 200 {
		return fmt.Errorf("error putting script: %s", resp.Status())
	}

	return nil
}

// Delete removes a stored script
func Delete(ctx context.Context, id string) error {
	if id == "" {
		return fmt.Errorf("script ID is required")
	}

	escapedID := url.PathEscape(id)

	resp, err := shared.Client.R().
		SetContext(ctx).
		SetHeader("Content-Type", "application/json").
		Delete(fmt.Sprintf("/_scripts/%s", escapedID))

	if err != nil {
		return fmt.Errorf("failed to delete script: %w", err)
	}

	if resp.StatusCode() == 404 {
		return fmt.Errorf("script '%s' not found", id)
	}

	if resp.StatusCode() != 200 {
		return fmt.Errorf("error deleting script: %s", resp.Status())
	}

	return nil
}
