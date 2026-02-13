package template

import (
	"context"
	"encoding/json"
	"testing"

	"github.com/pincher95/esctl/internal/testutil"
	"github.com/pincher95/esctl/shared"
)

func TestListComponents(t *testing.T) {
	mockResp := map[string]interface{}{
		"component_templates": []interface{}{
			map[string]interface{}{
				"name": "test-component",
				"component_template": map[string]interface{}{
					"version": 1,
					"template": map[string]interface{}{
						"settings": map[string]interface{}{
							"number_of_shards": 1,
						},
					},
				},
			},
		},
	}
	respJSON, _ := json.Marshal(mockResp)
	srv, cli := testutil.NewMockServer(string(respJSON), "/_component_template")
	defer srv.Close()
	shared.SetClient(cli)

	templates, err := ListComponents(context.Background())
	if err != nil {
		t.Fatalf("ListComponents() error = %v", err)
	}
	if len(templates) != 1 {
		t.Errorf("Expected 1 template, got %d", len(templates))
	}
	if _, ok := templates["test-component"]; !ok {
		t.Error("Expected template 'test-component' not found")
	}
}

func TestGetComponent(t *testing.T) {
	mockResp := map[string]interface{}{
		"component_templates": []interface{}{
			map[string]interface{}{
				"name": "test-component",
				"component_template": map[string]interface{}{
					"version": 1,
					"template": map[string]interface{}{
						"settings": map[string]interface{}{
							"number_of_shards": 1,
						},
						"mappings": map[string]interface{}{
							"properties": map[string]interface{}{
								"field1": map[string]interface{}{
									"type": "text",
								},
							},
						},
					},
				},
			},
		},
	}
	respJSON, _ := json.Marshal(mockResp)
	srv, cli := testutil.NewMockServer(string(respJSON), "/_component_template/test-component")
	defer srv.Close()
	shared.SetClient(cli)

	tmpl, err := GetComponent(context.Background(), "test-component")
	if err != nil {
		t.Fatalf("GetComponent() error = %v", err)
	}
	if tmpl.Name != "test-component" {
		t.Errorf("Expected name 'test-component', got %s", tmpl.Name)
	}
	if tmpl.Version != 1 {
		t.Errorf("Expected version 1, got %d", tmpl.Version)
	}
}

func TestGetComponentNotFound(t *testing.T) {
	srv, cli := testutil.NewMockServerWithStatus(404, "", "/_component_template/nonexistent")
	defer srv.Close()
	shared.SetClient(cli)

	_, err := GetComponent(context.Background(), "nonexistent")
	if err == nil {
		t.Error("Expected error for non-existent template, got nil")
	}
}

func TestPutComponent(t *testing.T) {
	mockResp := map[string]interface{}{
		"acknowledged": true,
	}
	respJSON, _ := json.Marshal(mockResp)
	srv, cli := testutil.NewMockServer(string(respJSON), "/_component_template/test-component")
	defer srv.Close()
	shared.SetClient(cli)

	tmpl := ComponentTemplate{
		Version: 1,
		Template: ComponentDefinition{
			Settings: map[string]interface{}{
				"number_of_shards": 1,
			},
		},
	}

	err := PutComponent(context.Background(), "test-component", tmpl)
	if err != nil {
		t.Fatalf("PutComponent() error = %v", err)
	}
}

func TestDeleteComponent(t *testing.T) {
	mockResp := map[string]interface{}{
		"acknowledged": true,
	}
	respJSON, _ := json.Marshal(mockResp)
	srv, cli := testutil.NewMockServer(string(respJSON), "/_component_template/test-component")
	defer srv.Close()
	shared.SetClient(cli)

	err := DeleteComponent(context.Background(), "test-component")
	if err != nil {
		t.Fatalf("DeleteComponent() error = %v", err)
	}
}

func TestDeleteComponentNotFound(t *testing.T) {
	srv, cli := testutil.NewMockServerWithStatus(404, "", "/_component_template/nonexistent")
	defer srv.Close()
	shared.SetClient(cli)

	err := DeleteComponent(context.Background(), "nonexistent")
	if err == nil {
		t.Error("Expected error for non-existent template, got nil")
	}
}

func TestExistsComponent(t *testing.T) {
	srv, cli := testutil.NewMockServer("", "/_component_template/test-component")
	defer srv.Close()
	shared.SetClient(cli)

	exists, err := ExistsComponent(context.Background(), "test-component")
	if err != nil {
		t.Fatalf("ExistsComponent() error = %v", err)
	}
	if !exists {
		t.Error("Expected template to exist")
	}
}

func TestExistsComponentNotFound(t *testing.T) {
	srv, cli := testutil.NewMockServerWithStatus(404, "", "/_component_template/nonexistent")
	defer srv.Close()
	shared.SetClient(cli)

	exists, err := ExistsComponent(context.Background(), "nonexistent")
	if err != nil {
		t.Fatalf("ExistsComponent() error = %v", err)
	}
	if exists {
		t.Error("Expected template to not exist")
	}
}
