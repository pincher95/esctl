package searchtemplate

import (
	"context"
	"testing"

	"github.com/pincher95/esctl/internal/testutil"
	"github.com/pincher95/esctl/shared"
)

func TestList(t *testing.T) {
	mockResponse := `{
		"template1": {
			"found": true,
			"_id": "template1",
			"template": {
				"query": {
					"match": {
						"title": "{{query_string}}"
					}
				}
			}
		},
		"template2": {
			"found": true,
			"_id": "template2",
			"template": {
				"query": {
					"range": {
						"price": {
							"gte": "{{min_price}}",
							"lte": "{{max_price}}"
						}
					}
				}
			}
		}
	}`

	srv, cli := testutil.NewMockServer(mockResponse, "/_scripts")
	defer srv.Close()

	shared.SetClient(cli)

	templates, err := List(context.Background())
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if len(templates) != 2 {
		t.Fatalf("expected 2 templates, got %d", len(templates))
	}

	if templates["template1"].ID != "template1" {
		t.Errorf("expected template1 ID, got %s", templates["template1"].ID)
	}
}

func TestGet(t *testing.T) {
	mockResponse := `{
		"found": true,
		"_id": "my-template",
		"template": {
			"query": {
				"match": {
					"title": "{{query_string}}"
				}
			}
		}
	}`

	srv, cli := testutil.NewMockServer(mockResponse, "/_scripts/my-template")
	defer srv.Close()

	shared.SetClient(cli)

	template, err := Get(context.Background(), "my-template")
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if template.ID != "my-template" {
		t.Errorf("expected ID to be my-template, got %s", template.ID)
	}

	if template.Template == nil {
		t.Error("expected template to have content")
	}
}

func TestGetNotFound(t *testing.T) {
	srv, cli := testutil.NewMockServerWithStatus(404, `{"error":"not found"}`, "/_scripts/nonexistent")
	defer srv.Close()

	shared.SetClient(cli)

	_, err := Get(context.Background(), "nonexistent")
	if err == nil {
		t.Fatal("expected error for non-existent template")
	}

	if err.Error() != "search template 'nonexistent' not found" {
		t.Errorf("unexpected error message: %v", err)
	}
}

func TestGetEmptyID(t *testing.T) {
	_, err := Get(context.Background(), "")
	if err == nil {
		t.Fatal("expected error for empty ID")
	}

	if err.Error() != "search template ID is required" {
		t.Errorf("unexpected error message: %v", err)
	}
}

func TestPut(t *testing.T) {
	mockResponse := `{"acknowledged": true}`

	srv, cli := testutil.NewMockServer(mockResponse, "/_scripts/my-template")
	defer srv.Close()

	shared.SetClient(cli)

	template := SearchTemplate{
		Template: map[string]any{
			"query": map[string]any{
				"match": map[string]any{
					"title": "{{query_string}}",
				},
			},
		},
	}

	err := Put(context.Background(), "my-template", template)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
}

func TestPutValidation(t *testing.T) {
	tests := []struct {
		name        string
		id          string
		template    SearchTemplate
		expectedErr string
	}{
		{
			name:        "empty ID",
			id:          "",
			template:    SearchTemplate{Template: map[string]any{"query": "test"}},
			expectedErr: "search template ID is required",
		},
		{
			name:        "empty template",
			id:          "test",
			template:    SearchTemplate{Template: nil},
			expectedErr: "template definition is required",
		},
		{
			name:        "empty template map",
			id:          "test",
			template:    SearchTemplate{Template: map[string]any{}},
			expectedErr: "template definition is required",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := Put(context.Background(), tt.id, tt.template)
			if err == nil {
				t.Fatal("expected error")
			}
			if err.Error() != tt.expectedErr {
				t.Errorf("expected error %q, got %q", tt.expectedErr, err.Error())
			}
		})
	}
}

func TestDelete(t *testing.T) {
	mockResponse := `{"acknowledged": true}`

	srv, cli := testutil.NewMockServer(mockResponse, "/_scripts/my-template")
	defer srv.Close()

	shared.SetClient(cli)

	err := Delete(context.Background(), "my-template")
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
}

func TestDeleteNotFound(t *testing.T) {
	srv, cli := testutil.NewMockServerWithStatus(404, `{"error":"not found"}`, "/_scripts/nonexistent")
	defer srv.Close()

	shared.SetClient(cli)

	err := Delete(context.Background(), "nonexistent")
	if err == nil {
		t.Fatal("expected error for non-existent template")
	}

	if err.Error() != "search template 'nonexistent' not found" {
		t.Errorf("unexpected error message: %v", err)
	}
}

func TestDeleteEmptyID(t *testing.T) {
	err := Delete(context.Background(), "")
	if err == nil {
		t.Fatal("expected error for empty ID")
	}

	if err.Error() != "search template ID is required" {
		t.Errorf("unexpected error message: %v", err)
	}
}

func TestRender(t *testing.T) {
	mockResponse := `{
		"template_output": {
			"query": {
				"match": {
					"title": "search text"
				}
			}
		}
	}`

	srv, cli := testutil.NewMockServer(mockResponse, "/_render/template")
	defer srv.Close()

	shared.SetClient(cli)

	params := map[string]any{
		"query_string": "search text",
	}

	result, err := Render(context.Background(), "my-template", params)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if result.TemplateOutput == nil {
		t.Error("expected template output")
	}
}

func TestRenderEmptyID(t *testing.T) {
	_, err := Render(context.Background(), "", nil)
	if err == nil {
		t.Fatal("expected error for empty ID")
	}

	if err.Error() != "search template ID is required" {
		t.Errorf("unexpected error message: %v", err)
	}
}

func TestRenderInline(t *testing.T) {
	mockResponse := `{
		"template_output": {
			"query": {
				"match": {
					"title": "search text"
				}
			}
		}
	}`

	srv, cli := testutil.NewMockServer(mockResponse, "/_render/template")
	defer srv.Close()

	shared.SetClient(cli)

	template := map[string]any{
		"query": map[string]any{
			"match": map[string]any{
				"title": "{{query_string}}",
			},
		},
	}

	params := map[string]any{
		"query_string": "search text",
	}

	result, err := RenderInline(context.Background(), template, params)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if result.TemplateOutput == nil {
		t.Error("expected template output")
	}
}

func TestRenderInlineEmptyTemplate(t *testing.T) {
	_, err := RenderInline(context.Background(), nil, nil)
	if err == nil {
		t.Fatal("expected error for empty template")
	}

	if err.Error() != "template is required" {
		t.Errorf("unexpected error message: %v", err)
	}
}
