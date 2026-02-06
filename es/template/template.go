package template

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/pincher95/esctl/shared"
)

// Template represents an index template
type Template struct {
	Name          string                 `json:"-"`
	IndexPatterns []string               `json:"index_patterns"`
	Template      TemplateDefinition     `json:"template,omitempty"`
	Priority      int                    `json:"priority,omitempty"`
	Version       int                    `json:"version,omitempty"`
	Meta          map[string]interface{} `json:"_meta,omitempty"`
	ComposedOf    []string               `json:"composed_of,omitempty"`
}

// TemplateDefinition contains template settings and mappings
type TemplateDefinition struct {
	Settings map[string]interface{} `json:"settings,omitempty"`
	Mappings map[string]interface{} `json:"mappings,omitempty"`
	Aliases  map[string]interface{} `json:"aliases,omitempty"`
}

// ListResponse represents the response from listing templates
type ListResponse map[string]Template

// List retrieves all index templates
func List(ctx context.Context) (ListResponse, error) {
	var result map[string]interface{}

	resp, err := shared.Client.R().
		SetContext(ctx).
		SetHeader("Content-Type", "application/json").
		SetResult(&result).
		Get("_index_template")

	if err != nil {
		return nil, fmt.Errorf("failed to list templates: %w", err)
	}

	if resp.StatusCode() != 200 {
		return nil, fmt.Errorf("failed to list templates: %s", resp.Status())
	}

	// Parse the response
	templates := make(ListResponse)
	if indexTemplates, ok := result["index_templates"].([]interface{}); ok {
		for _, item := range indexTemplates {
			if itemMap, ok := item.(map[string]interface{}); ok {
				if name, ok := itemMap["name"].(string); ok {
					if tmplData, ok := itemMap["index_template"].(map[string]interface{}); ok {
						tmplBytes, _ := json.Marshal(tmplData)
						var tmpl Template
						if err := json.Unmarshal(tmplBytes, &tmpl); err == nil {
							tmpl.Name = name
							templates[name] = tmpl
						}
					}
				}
			}
		}
	}

	return templates, nil
}

// Get retrieves a specific index template
func Get(ctx context.Context, name string) (*Template, error) {
	var result map[string]interface{}

	endpoint := fmt.Sprintf("_index_template/%s", name)
	resp, err := shared.Client.R().
		SetContext(ctx).
		SetHeader("Content-Type", "application/json").
		SetResult(&result).
		Get(endpoint)

	if err != nil {
		return nil, fmt.Errorf("failed to get template: %w", err)
	}

	if resp.StatusCode() == 404 {
		return nil, fmt.Errorf("template not found: %s", name)
	}

	if resp.StatusCode() != 200 {
		return nil, fmt.Errorf("failed to get template: %s", resp.Status())
	}

	// Parse the response
	if indexTemplates, ok := result["index_templates"].([]interface{}); ok && len(indexTemplates) > 0 {
		if itemMap, ok := indexTemplates[0].(map[string]interface{}); ok {
			if tmplData, ok := itemMap["index_template"].(map[string]interface{}); ok {
				tmplBytes, _ := json.Marshal(tmplData)
				var tmpl Template
				if err := json.Unmarshal(tmplBytes, &tmpl); err == nil {
					tmpl.Name = name
					return &tmpl, nil
				}
			}
		}
	}

	return nil, fmt.Errorf("failed to parse template response")
}

// Put creates or updates an index template
func Put(ctx context.Context, name string, template Template) error {
	endpoint := fmt.Sprintf("_index_template/%s", name)

	resp, err := shared.Client.R().
		SetContext(ctx).
		SetHeader("Content-Type", "application/json").
		SetBody(template).
		Put(endpoint)

	if err != nil {
		return fmt.Errorf("failed to put template: %w", err)
	}

	if resp.StatusCode() != 200 {
		return fmt.Errorf("failed to put template: %s - %s", resp.Status(), string(resp.Body()))
	}

	return nil
}

// Delete removes an index template
func Delete(ctx context.Context, name string) error {
	endpoint := fmt.Sprintf("_index_template/%s", name)

	resp, err := shared.Client.R().
		SetContext(ctx).
		SetHeader("Content-Type", "application/json").
		Delete(endpoint)

	if err != nil {
		return fmt.Errorf("failed to delete template: %w", err)
	}

	if resp.StatusCode() == 404 {
		return fmt.Errorf("template not found: %s", name)
	}

	if resp.StatusCode() != 200 {
		return fmt.Errorf("failed to delete template: %s", resp.Status())
	}

	return nil
}

// Exists checks if a template exists
func Exists(ctx context.Context, name string) (bool, error) {
	endpoint := fmt.Sprintf("_index_template/%s", name)

	resp, err := shared.Client.R().
		SetContext(ctx).
		Head(endpoint)

	if err != nil {
		return false, fmt.Errorf("failed to check template existence: %w", err)
	}

	return resp.StatusCode() == 200, nil
}
