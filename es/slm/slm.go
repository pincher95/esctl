// Package slm wraps the Snapshot Lifecycle Management (_slm) APIs for scheduling
// and retaining snapshots.
package slm

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/pincher95/esctl/shared"
)

// Policy represents a single SLM policy as returned by the API.
type Policy struct {
	Name          string           `json:"-"`
	Version       int              `json:"version,omitempty"`
	ModifiedDate  string           `json:"modified_date,omitempty"`
	NextExecution string           `json:"next_execution,omitempty"`
	Policy        PolicyDefinition `json:"policy"`
}

// PolicyDefinition is the user-supplied SLM policy configuration.
type PolicyDefinition struct {
	Name       string         `json:"name,omitempty"`
	Schedule   string         `json:"schedule,omitempty"`
	Repository string         `json:"repository,omitempty"`
	Config     map[string]any `json:"config,omitempty"`
	Retention  map[string]any `json:"retention,omitempty"`
}

// ListResponse maps policy id to policy.
type ListResponse map[string]Policy

// ExecuteResponse is returned when a policy is triggered on demand.
type ExecuteResponse struct {
	SnapshotName string `json:"snapshot_name"`
}

// List retrieves all SLM policies.
func List(ctx context.Context) (ListResponse, error) {
	var result ListResponse
	resp, err := shared.Client.R().
		SetContext(ctx).
		SetHeader("Content-Type", "application/json").
		SetResult(&result).
		Get("_slm/policy")
	if err != nil {
		return nil, fmt.Errorf("failed to list SLM policies: %w", err)
	}
	if resp.StatusCode() != 200 {
		return nil, fmt.Errorf("failed to list SLM policies: %s", resp.Status())
	}
	for name, policy := range result {
		p := policy
		p.Name = name
		result[name] = p
	}
	return result, nil
}

// Get retrieves a single SLM policy by id.
func Get(ctx context.Context, name string) (*Policy, error) {
	var result map[string]Policy
	endpoint := fmt.Sprintf("_slm/policy/%s", name)
	resp, err := shared.Client.R().
		SetContext(ctx).
		SetHeader("Content-Type", "application/json").
		SetResult(&result).
		Get(endpoint)
	if err != nil {
		return nil, fmt.Errorf("failed to get SLM policy: %w", err)
	}
	if resp.StatusCode() == 404 {
		return nil, fmt.Errorf("SLM policy not found: %s", name)
	}
	if resp.StatusCode() != 200 {
		return nil, fmt.Errorf("failed to get SLM policy: %s", resp.Status())
	}
	if policy, ok := result[name]; ok {
		policy.Name = name
		return &policy, nil
	}
	return nil, fmt.Errorf("SLM policy not found: %s", name)
}

// Put creates or updates an SLM policy from a raw JSON definition.
func Put(ctx context.Context, name string, body []byte) error {
	var def map[string]any
	if err := json.Unmarshal(body, &def); err != nil {
		return fmt.Errorf("failed to parse policy JSON: %w", err)
	}

	endpoint := fmt.Sprintf("_slm/policy/%s", name)
	resp, err := shared.Client.R().
		SetContext(ctx).
		SetHeader("Content-Type", "application/json").
		SetBody(def).
		Put(endpoint)
	if err != nil {
		return fmt.Errorf("failed to put SLM policy: %w", err)
	}
	if resp.StatusCode() != 200 {
		return fmt.Errorf("failed to put SLM policy: %s - %s", resp.Status(), string(resp.Body()))
	}
	return nil
}

// Delete removes an SLM policy.
func Delete(ctx context.Context, name string) error {
	endpoint := fmt.Sprintf("_slm/policy/%s", name)
	resp, err := shared.Client.R().
		SetContext(ctx).
		SetHeader("Content-Type", "application/json").
		Delete(endpoint)
	if err != nil {
		return fmt.Errorf("failed to delete SLM policy: %w", err)
	}
	if resp.StatusCode() == 404 {
		return fmt.Errorf("SLM policy not found: %s", name)
	}
	if resp.StatusCode() != 200 {
		return fmt.Errorf("failed to delete SLM policy: %s", resp.Status())
	}
	return nil
}

// Execute triggers a policy immediately and returns the created snapshot name.
func Execute(ctx context.Context, name string) (*ExecuteResponse, error) {
	endpoint := fmt.Sprintf("_slm/policy/%s/_execute", name)
	var out ExecuteResponse
	resp, err := shared.Client.R().
		SetContext(ctx).
		SetHeader("Content-Type", "application/json").
		SetResult(&out).
		Post(endpoint)
	if err != nil {
		return nil, fmt.Errorf("failed to execute SLM policy: %w", err)
	}
	if resp.StatusCode() != 200 {
		return nil, fmt.Errorf("failed to execute SLM policy: %s - %s", resp.Status(), string(resp.Body()))
	}
	return &out, nil
}
