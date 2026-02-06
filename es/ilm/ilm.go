package ilm

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/pincher95/esctl/internal/logger"
	"github.com/pincher95/esctl/shared"
)

// Policy represents an ILM policy
type Policy struct {
	Name     string                 `json:"-"`
	Version  int                    `json:"version,omitempty"`
	Modified string                 `json:"modified_date,omitempty"`
	Policy   PolicyDefinition       `json:"policy"`
	Meta     map[string]interface{} `json:"_meta,omitempty"`
}

// PolicyDefinition contains the actual ILM policy configuration
type PolicyDefinition struct {
	Phases map[string]Phase `json:"phases"`
}

// Phase represents a phase in the ILM policy
type Phase struct {
	MinAge  string                 `json:"min_age,omitempty"`
	Actions map[string]interface{} `json:"actions,omitempty"`
}

// ListResponse represents the response from listing ILM policies
type ListResponse map[string]Policy

// ExplainResponse represents the explain response for an index
type ExplainResponse struct {
	Indices map[string]IndexExplain `json:"indices"`
}

// IndexExplain contains ILM explain information for a single index
type IndexExplain struct {
	Index           string                 `json:"index"`
	Managed         bool                   `json:"managed"`
	Policy          string                 `json:"policy,omitempty"`
	Phase           string                 `json:"phase,omitempty"`
	Action          string                 `json:"action,omitempty"`
	Step            string                 `json:"step,omitempty"`
	FailedStep      string                 `json:"failed_step,omitempty"`
	StepTime        string                 `json:"step_time_millis,omitempty"`
	PhaseTime       string                 `json:"phase_time_millis,omitempty"`
	ActionTime      string                 `json:"action_time_millis,omitempty"`
	Age             string                 `json:"age,omitempty"`
	FailedStepRetry int                    `json:"failed_step_retry_count,omitempty"`
	PhaseExecution  map[string]interface{} `json:"phase_execution,omitempty"`
}

// List retrieves all ILM policies
func List(ctx context.Context) (ListResponse, error) {
	logger.Debug("listing ILM policies")

	var result ListResponse

	resp, err := shared.Client.R().
		SetContext(ctx).
		SetHeader("Content-Type", "application/json").
		SetResult(&result).
		Get("_ilm/policy")

	if err != nil {
		logger.Error("failed to list ILM policies", "error", err)
		return nil, fmt.Errorf("failed to list ILM policies: %w", err)
	}

	if resp.StatusCode() != 200 {
		logger.Error("unexpected status listing ILM policies", "status", resp.StatusCode())
		return nil, fmt.Errorf("failed to list ILM policies: %s", resp.Status())
	}

	// Set policy names from keys
	for name, policy := range result {
		p := policy
		p.Name = name
		result[name] = p
	}

	logger.Debug("listed ILM policies", "count", len(result))
	return result, nil
}

// Get retrieves a specific ILM policy
func Get(ctx context.Context, name string) (*Policy, error) {
	logger.Debug("getting ILM policy", "name", name)

	var result map[string]Policy

	endpoint := fmt.Sprintf("_ilm/policy/%s", name)
	resp, err := shared.Client.R().
		SetContext(ctx).
		SetHeader("Content-Type", "application/json").
		SetResult(&result).
		Get(endpoint)

	if err != nil {
		logger.Error("failed to get ILM policy", "name", name, "error", err)
		return nil, fmt.Errorf("failed to get ILM policy: %w", err)
	}

	if resp.StatusCode() == 404 {
		return nil, fmt.Errorf("ILM policy not found: %s", name)
	}

	if resp.StatusCode() != 200 {
		logger.Error("unexpected status getting ILM policy", "name", name, "status", resp.StatusCode())
		return nil, fmt.Errorf("failed to get ILM policy: %s", resp.Status())
	}

	if policy, ok := result[name]; ok {
		policy.Name = name
		logger.Debug("retrieved ILM policy", "name", name)
		return &policy, nil
	}

	return nil, fmt.Errorf("failed to parse ILM policy response")
}

// Put creates or updates an ILM policy
func Put(ctx context.Context, name string, policy Policy) error {
	logger.Debug("putting ILM policy", "name", name)

	endpoint := fmt.Sprintf("_ilm/policy/%s", name)

	// Only send the policy definition, not the wrapper
	body := map[string]interface{}{
		"policy": policy.Policy,
	}

	resp, err := shared.Client.R().
		SetContext(ctx).
		SetHeader("Content-Type", "application/json").
		SetBody(body).
		Put(endpoint)

	if err != nil {
		logger.Error("failed to put ILM policy", "name", name, "error", err)
		return fmt.Errorf("failed to put ILM policy: %w", err)
	}

	if resp.StatusCode() != 200 {
		logger.Error("unexpected status putting ILM policy", "name", name, "status", resp.StatusCode())
		return fmt.Errorf("failed to put ILM policy: %s - %s", resp.Status(), string(resp.Body()))
	}

	logger.Info("ILM policy created/updated", "name", name)
	return nil
}

// Delete removes an ILM policy
func Delete(ctx context.Context, name string) error {
	logger.Debug("deleting ILM policy", "name", name)

	endpoint := fmt.Sprintf("_ilm/policy/%s", name)

	resp, err := shared.Client.R().
		SetContext(ctx).
		SetHeader("Content-Type", "application/json").
		Delete(endpoint)

	if err != nil {
		logger.Error("failed to delete ILM policy", "name", name, "error", err)
		return fmt.Errorf("failed to delete ILM policy: %w", err)
	}

	if resp.StatusCode() == 404 {
		return fmt.Errorf("ILM policy not found: %s", name)
	}

	if resp.StatusCode() != 200 {
		logger.Error("unexpected status deleting ILM policy", "name", name, "status", resp.StatusCode())
		return fmt.Errorf("failed to delete ILM policy: %s", resp.Status())
	}

	logger.Info("ILM policy deleted", "name", name)
	return nil
}

// Exists checks if an ILM policy exists
func Exists(ctx context.Context, name string) (bool, error) {
	logger.Debug("checking ILM policy existence", "name", name)

	endpoint := fmt.Sprintf("_ilm/policy/%s", name)

	resp, err := shared.Client.R().
		SetContext(ctx).
		Head(endpoint)

	if err != nil {
		return false, fmt.Errorf("failed to check ILM policy existence: %w", err)
	}

	exists := resp.StatusCode() == 200
	logger.Debug("ILM policy existence check", "name", name, "exists", exists)
	return exists, nil
}

// Explain returns ILM explain information for indices
func Explain(ctx context.Context, index string) (*ExplainResponse, error) {
	logger.Debug("explaining ILM status", "index", index)

	endpoint := fmt.Sprintf("%s/_ilm/explain", index)

	var result ExplainResponse
	resp, err := shared.Client.R().
		SetContext(ctx).
		SetHeader("Content-Type", "application/json").
		SetResult(&result).
		Get(endpoint)

	if err != nil {
		logger.Error("failed to explain ILM status", "index", index, "error", err)
		return nil, fmt.Errorf("failed to explain ILM status: %w", err)
	}

	if resp.StatusCode() != 200 {
		logger.Error("unexpected status explaining ILM", "index", index, "status", resp.StatusCode())
		return nil, fmt.Errorf("failed to explain ILM status: %s", resp.Status())
	}

	logger.Debug("explained ILM status", "index", index, "count", len(result.Indices))
	return &result, nil
}

// Retry retries the failed step for an index
func Retry(ctx context.Context, index string) error {
	logger.Debug("retrying ILM failed step", "index", index)

	endpoint := fmt.Sprintf("%s/_ilm/retry", index)

	resp, err := shared.Client.R().
		SetContext(ctx).
		SetHeader("Content-Type", "application/json").
		Post(endpoint)

	if err != nil {
		logger.Error("failed to retry ILM step", "index", index, "error", err)
		return fmt.Errorf("failed to retry ILM step: %w", err)
	}

	if resp.StatusCode() != 200 {
		logger.Error("unexpected status retrying ILM", "index", index, "status", resp.StatusCode())
		return fmt.Errorf("failed to retry ILM step: %s - %s", resp.Status(), string(resp.Body()))
	}

	logger.Info("ILM retry initiated", "index", index)
	return nil
}

// MoveToStep moves an index to a specific step in its ILM policy
func MoveToStep(ctx context.Context, index string, currentStep, nextStep map[string]string) error {
	logger.Debug("moving index to ILM step", "index", index)

	endpoint := fmt.Sprintf("_ilm/move/%s", index)

	body := map[string]interface{}{
		"current_step": currentStep,
		"next_step":    nextStep,
	}

	resp, err := shared.Client.R().
		SetContext(ctx).
		SetHeader("Content-Type", "application/json").
		SetBody(body).
		Post(endpoint)

	if err != nil {
		logger.Error("failed to move index to step", "index", index, "error", err)
		return fmt.Errorf("failed to move index to step: %w", err)
	}

	if resp.StatusCode() != 200 {
		logger.Error("unexpected status moving index", "index", index, "status", resp.StatusCode())
		return fmt.Errorf("failed to move index to step: %s - %s", resp.Status(), string(resp.Body()))
	}

	logger.Info("index moved to step", "index", index)
	return nil
}

// ParsePolicyFromJSON parses a policy from JSON bytes
func ParsePolicyFromJSON(data []byte) (*Policy, error) {
	var policy Policy
	if err := json.Unmarshal(data, &policy); err != nil {
		return nil, fmt.Errorf("failed to parse policy JSON: %w", err)
	}
	return &policy, nil
}
