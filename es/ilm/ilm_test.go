package ilm

import (
	"context"
	"encoding/json"
	"testing"

	"github.com/pincher95/esctl/internal/testutil"
	"github.com/pincher95/esctl/shared"
)

func TestList(t *testing.T) {
	mockResp := map[string]any{
		"hot_delete_policy": map[string]any{
			"version":       1,
			"modified_date": "2024-01-01T00:00:00.000Z",
			"policy": map[string]any{
				"phases": map[string]any{
					"hot": map[string]any{
						"actions": map[string]any{
							"rollover": map[string]any{
								"max_age":  "7d",
								"max_size": "50gb",
							},
						},
					},
					"delete": map[string]any{
						"min_age": "30d",
						"actions": map[string]any{
							"delete": map[string]any{},
						},
					},
				},
			},
		},
		"cold_policy": map[string]any{
			"version": 2,
			"policy": map[string]any{
				"phases": map[string]any{
					"cold": map[string]any{
						"min_age": "14d",
						"actions": map[string]any{
							"freeze": map[string]any{},
						},
					},
				},
			},
		},
	}

	respJSON, _ := json.Marshal(mockResp)
	srv, cli := testutil.NewMockServer(string(respJSON), "/_ilm/policy")
	defer srv.Close()
	shared.SetClient(cli)

	policies, err := List(context.Background())
	if err != nil {
		t.Fatalf("List() error = %v", err)
	}

	if len(policies) != 2 {
		t.Errorf("Expected 2 policies, got %d", len(policies))
	}

	if _, ok := policies["hot_delete_policy"]; !ok {
		t.Error("Expected hot_delete_policy to be present")
	}

	if _, ok := policies["cold_policy"]; !ok {
		t.Error("Expected cold_policy to be present")
	}

	// Check that names were set
	if policies["hot_delete_policy"].Name != "hot_delete_policy" {
		t.Errorf("Expected name to be 'hot_delete_policy', got '%s'", policies["hot_delete_policy"].Name)
	}
}

func TestGet(t *testing.T) {
	mockResp := map[string]any{
		"hot_delete_policy": map[string]any{
			"version":       1,
			"modified_date": "2024-01-01T00:00:00.000Z",
			"policy": map[string]any{
				"phases": map[string]any{
					"hot": map[string]any{
						"actions": map[string]any{
							"rollover": map[string]any{
								"max_age": "7d",
							},
						},
					},
				},
			},
		},
	}

	respJSON, _ := json.Marshal(mockResp)
	srv, cli := testutil.NewMockServer(string(respJSON), "/_ilm/policy/hot_delete_policy")
	defer srv.Close()
	shared.SetClient(cli)

	policy, err := Get(context.Background(), "hot_delete_policy")
	if err != nil {
		t.Fatalf("Get() error = %v", err)
	}

	if policy.Name != "hot_delete_policy" {
		t.Errorf("Expected name 'hot_delete_policy', got '%s'", policy.Name)
	}

	if policy.Version != 1 {
		t.Errorf("Expected version 1, got %d", policy.Version)
	}

	if len(policy.Policy.Phases) != 1 {
		t.Errorf("Expected 1 phase, got %d", len(policy.Policy.Phases))
	}

	if _, ok := policy.Policy.Phases["hot"]; !ok {
		t.Error("Expected 'hot' phase to be present")
	}
}

func TestGetNotFound(t *testing.T) {
	srv, cli := testutil.NewMockServerWithStatus(404, `{"error":"not found"}`, "/_ilm/policy/nonexistent")
	defer srv.Close()
	shared.SetClient(cli)

	_, err := Get(context.Background(), "nonexistent")
	if err == nil {
		t.Error("Expected error for nonexistent policy")
	}

	if err.Error() != "ILM policy not found: nonexistent" {
		t.Errorf("Expected 'ILM policy not found' error, got: %s", err.Error())
	}
}

func TestPut(t *testing.T) {
	mockResp := map[string]any{
		"acknowledged": true,
	}

	respJSON, _ := json.Marshal(mockResp)
	srv, cli := testutil.NewMockServer(string(respJSON), "/_ilm/policy/new-policy")
	defer srv.Close()
	shared.SetClient(cli)

	policy := Policy{
		Policy: PolicyDefinition{
			Phases: map[string]Phase{
				"hot": {
					Actions: map[string]any{
						"rollover": map[string]any{
							"max_age": "7d",
						},
					},
				},
			},
		},
	}

	err := Put(context.Background(), "new-policy", policy)
	if err != nil {
		t.Fatalf("Put() error = %v", err)
	}
}

func TestDelete(t *testing.T) {
	mockResp := map[string]any{
		"acknowledged": true,
	}

	respJSON, _ := json.Marshal(mockResp)
	srv, cli := testutil.NewMockServer(string(respJSON), "/_ilm/policy/old-policy")
	defer srv.Close()
	shared.SetClient(cli)

	err := Delete(context.Background(), "old-policy")
	if err != nil {
		t.Fatalf("Delete() error = %v", err)
	}
}

func TestDeleteNotFound(t *testing.T) {
	srv, cli := testutil.NewMockServerWithStatus(404, `{"error":"not found"}`, "/_ilm/policy/nonexistent")
	defer srv.Close()
	shared.SetClient(cli)

	err := Delete(context.Background(), "nonexistent")
	if err == nil {
		t.Error("Expected error for nonexistent policy")
	}
}

func TestExists(t *testing.T) {
	srv, cli := testutil.NewMockServerWithStatus(200, "", "/_ilm/policy/existing-policy")
	defer srv.Close()
	shared.SetClient(cli)

	exists, err := Exists(context.Background(), "existing-policy")
	if err != nil {
		t.Fatalf("Exists() error = %v", err)
	}

	if !exists {
		t.Error("Expected policy to exist")
	}
}

func TestNotExists(t *testing.T) {
	srv, cli := testutil.NewMockServerWithStatus(404, "", "/_ilm/policy/nonexistent")
	defer srv.Close()
	shared.SetClient(cli)

	exists, err := Exists(context.Background(), "nonexistent")
	if err != nil {
		t.Fatalf("Exists() error = %v", err)
	}

	if exists {
		t.Error("Expected policy to not exist")
	}
}

func TestExplain(t *testing.T) {
	mockResp := map[string]any{
		"indices": map[string]any{
			"test-index-000001": map[string]any{
				"index":   "test-index-000001",
				"managed": true,
				"policy":  "hot_delete_policy",
				"phase":   "hot",
				"action":  "rollover",
				"step":    "check-rollover-ready",
				"age":     "5d",
			},
		},
	}

	respJSON, _ := json.Marshal(mockResp)
	srv, cli := testutil.NewMockServer(string(respJSON), "/test-index-*/_ilm/explain")
	defer srv.Close()
	shared.SetClient(cli)

	result, err := Explain(context.Background(), "test-index-*")
	if err != nil {
		t.Fatalf("Explain() error = %v", err)
	}

	if len(result.Indices) != 1 {
		t.Errorf("Expected 1 index, got %d", len(result.Indices))
	}

	indexInfo, ok := result.Indices["test-index-000001"]
	if !ok {
		t.Fatal("Expected test-index-000001 to be present")
	}

	if !indexInfo.Managed {
		t.Error("Expected index to be managed")
	}

	if indexInfo.Policy != "hot_delete_policy" {
		t.Errorf("Expected policy 'hot_delete_policy', got '%s'", indexInfo.Policy)
	}

	if indexInfo.Phase != "hot" {
		t.Errorf("Expected phase 'hot', got '%s'", indexInfo.Phase)
	}
}

func TestRetry(t *testing.T) {
	mockResp := map[string]any{
		"acknowledged": true,
	}

	respJSON, _ := json.Marshal(mockResp)
	srv, cli := testutil.NewMockServer(string(respJSON), "/test-index/_ilm/retry")
	defer srv.Close()
	shared.SetClient(cli)

	err := Retry(context.Background(), "test-index")
	if err != nil {
		t.Fatalf("Retry() error = %v", err)
	}
}

func TestMoveToStep(t *testing.T) {
	mockResp := map[string]any{
		"acknowledged": true,
	}

	respJSON, _ := json.Marshal(mockResp)
	srv, cli := testutil.NewMockServer(string(respJSON), "/_ilm/move/test-index")
	defer srv.Close()
	shared.SetClient(cli)

	currentStep := map[string]string{
		"phase":  "hot",
		"action": "rollover",
		"name":   "check-rollover-ready",
	}

	nextStep := map[string]string{
		"phase":  "hot",
		"action": "rollover",
		"name":   "attempt-rollover",
	}

	err := MoveToStep(context.Background(), "test-index", currentStep, nextStep)
	if err != nil {
		t.Fatalf("MoveToStep() error = %v", err)
	}
}

func TestParsePolicyFromJSON(t *testing.T) {
	jsonData := `{
		"policy": {
			"phases": {
				"hot": {
					"actions": {
						"rollover": {
							"max_age": "7d"
						}
					}
				}
			}
		}
	}`

	policy, err := ParsePolicyFromJSON([]byte(jsonData))
	if err != nil {
		t.Fatalf("ParsePolicyFromJSON() error = %v", err)
	}

	if len(policy.Policy.Phases) != 1 {
		t.Errorf("Expected 1 phase, got %d", len(policy.Policy.Phases))
	}

	if _, ok := policy.Policy.Phases["hot"]; !ok {
		t.Error("Expected 'hot' phase to be present")
	}
}

func TestParsePolicyFromJSONInvalid(t *testing.T) {
	jsonData := `{invalid json}`

	_, err := ParsePolicyFromJSON([]byte(jsonData))
	if err == nil {
		t.Error("Expected error for invalid JSON")
	}
}
