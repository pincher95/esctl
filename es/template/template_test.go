package template

import (
	"context"
	"encoding/json"
	"testing"

	"github.com/pincher95/esctl/internal/testutil"
	"github.com/pincher95/esctl/shared"
)

func TestList(t *testing.T) {
	mockResp := map[string]interface{}{
		"index_templates": []interface{}{
			map[string]interface{}{
				"name": "logs-template",
				"index_template": map[string]interface{}{
					"index_patterns": []interface{}{"logs-*"},
					"priority":       100,
					"version":        1,
				},
			},
			map[string]interface{}{
				"name": "metrics-template",
				"index_template": map[string]interface{}{
					"index_patterns": []interface{}{"metrics-*"},
					"priority":       50,
				},
			},
		},
	}

	respJSON, _ := json.Marshal(mockResp)
	srv, cli := testutil.NewMockServer(string(respJSON), "/_index_template")
	defer srv.Close()
	shared.SetClient(cli)

	templates, err := List(context.Background())
	if err != nil {
		t.Fatalf("List() error = %v", err)
	}

	if len(templates) != 2 {
		t.Errorf("Expected 2 templates, got %d", len(templates))
	}

	if tmpl, ok := templates["logs-template"]; !ok {
		t.Error("Expected logs-template to be present")
	} else {
		if tmpl.Priority != 100 {
			t.Errorf("Expected priority 100, got %d", tmpl.Priority)
		}
		if len(tmpl.IndexPatterns) != 1 || tmpl.IndexPatterns[0] != "logs-*" {
			t.Errorf("Expected index pattern 'logs-*', got %v", tmpl.IndexPatterns)
		}
	}
}

func TestGet(t *testing.T) {
	mockResp := map[string]interface{}{
		"index_templates": []interface{}{
			map[string]interface{}{
				"name": "logs-template",
				"index_template": map[string]interface{}{
					"index_patterns": []interface{}{"logs-*", "app-logs-*"},
					"priority":       100,
					"version":        1,
					"template": map[string]interface{}{
						"settings": map[string]interface{}{
							"number_of_shards": 3,
						},
					},
				},
			},
		},
	}

	respJSON, _ := json.Marshal(mockResp)
	srv, cli := testutil.NewMockServer(string(respJSON), "/_index_template/logs-template")
	defer srv.Close()
	shared.SetClient(cli)

	tmpl, err := Get(context.Background(), "logs-template")
	if err != nil {
		t.Fatalf("Get() error = %v", err)
	}

	if tmpl.Name != "logs-template" {
		t.Errorf("Expected name 'logs-template', got %s", tmpl.Name)
	}

	if len(tmpl.IndexPatterns) != 2 {
		t.Errorf("Expected 2 index patterns, got %d", len(tmpl.IndexPatterns))
	}

	if tmpl.Priority != 100 {
		t.Errorf("Expected priority 100, got %d", tmpl.Priority)
	}
}

func TestGetNotFound(t *testing.T) {
	srv, cli := testutil.NewMockServerWithStatus(404, `{"error":"not found"}`, "/_index_template/nonexistent")
	defer srv.Close()
	shared.SetClient(cli)

	_, err := Get(context.Background(), "nonexistent")
	if err == nil {
		t.Error("Expected error for nonexistent template")
	}
}

func TestPut(t *testing.T) {
	mockResp := map[string]interface{}{
		"acknowledged": true,
	}

	respJSON, _ := json.Marshal(mockResp)
	srv, cli := testutil.NewMockServer(string(respJSON), "/_index_template/new-template")
	defer srv.Close()
	shared.SetClient(cli)

	tmpl := Template{
		IndexPatterns: []string{"test-*"},
		Priority:      50,
		Version:       1,
	}

	err := Put(context.Background(), "new-template", tmpl)
	if err != nil {
		t.Fatalf("Put() error = %v", err)
	}
}

func TestDelete(t *testing.T) {
	mockResp := map[string]interface{}{
		"acknowledged": true,
	}

	respJSON, _ := json.Marshal(mockResp)
	srv, cli := testutil.NewMockServer(string(respJSON), "/_index_template/old-template")
	defer srv.Close()
	shared.SetClient(cli)

	err := Delete(context.Background(), "old-template")
	if err != nil {
		t.Fatalf("Delete() error = %v", err)
	}
}

func TestExists(t *testing.T) {
	srv, cli := testutil.NewMockServerWithStatus(200, "", "/_index_template/existing-template")
	defer srv.Close()
	shared.SetClient(cli)

	exists, err := Exists(context.Background(), "existing-template")
	if err != nil {
		t.Fatalf("Exists() error = %v", err)
	}

	if !exists {
		t.Error("Expected template to exist")
	}
}

func TestNotExists(t *testing.T) {
	srv, cli := testutil.NewMockServerWithStatus(404, "", "/_index_template/nonexistent")
	defer srv.Close()
	shared.SetClient(cli)

	exists, err := Exists(context.Background(), "nonexistent")
	if err != nil {
		t.Fatalf("Exists() error = %v", err)
	}

	if exists {
		t.Error("Expected template to not exist")
	}
}
