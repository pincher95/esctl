package template

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/pincher95/esctl/shared"
)

// Template represents an index template
type Template struct {
	Name          string             `json:"-"`
	IndexPatterns []string           `json:"index_patterns"`
	Template      TemplateDefinition `json:"template"`
	Priority      int                `json:"priority,omitempty"`
	Version       int                `json:"version,omitempty"`
	Meta          map[string]any     `json:"_meta,omitempty"`
	ComposedOf    []string           `json:"composed_of,omitempty"`
}

// TemplateDefinition contains template settings and mappings
type TemplateDefinition struct {
	Settings map[string]any `json:"settings,omitempty"`
	Mappings map[string]any `json:"mappings,omitempty"`
	Aliases  map[string]any `json:"aliases,omitempty"`
}

// ListResponse represents the response from listing templates
type ListResponse map[string]Template

// List retrieves all index templates
func List(ctx context.Context) (ListResponse, error) {
	var result map[string]any

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
	if indexTemplates, ok := result["index_templates"].([]any); ok {
		for _, item := range indexTemplates {
			if itemMap, ok := item.(map[string]any); ok {
				if name, ok := itemMap["name"].(string); ok {
					if tmplData, ok := itemMap["index_template"].(map[string]any); ok {
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

// ListLegacy retrieves all legacy index templates from the _template endpoint
func ListLegacy(ctx context.Context) (ListResponse, error) {
	var result map[string]any

	resp, err := shared.Client.R().
		SetContext(ctx).
		SetHeader("Content-Type", "application/json").
		SetResult(&result).
		Get("_template")

	if err != nil {
		return nil, fmt.Errorf("failed to list legacy templates: %w", err)
	}

	if resp.StatusCode() != 200 {
		return nil, fmt.Errorf("failed to list legacy templates: %s", resp.Status())
	}

	// Legacy template response is a flat map: { "name": { "order":..., "index_patterns":..., "settings":..., ... } }
	templates := make(ListResponse)
	for name, item := range result {
		itemMap, ok := item.(map[string]any)
		if !ok {
			continue
		}

		tmplBytes, _ := json.Marshal(itemMap)

		// Legacy format has settings/mappings/aliases at top level and "order" instead of "priority"
		var raw struct {
			Order         int            `json:"order"`
			IndexPatterns []string       `json:"index_patterns"`
			Settings      map[string]any `json:"settings"`
			Mappings      map[string]any `json:"mappings"`
			Aliases       map[string]any `json:"aliases"`
			Version       int            `json:"version"`
		}
		if err := json.Unmarshal(tmplBytes, &raw); err != nil {
			continue
		}

		tmpl := Template{
			Name:          name,
			IndexPatterns: raw.IndexPatterns,
			Priority:      raw.Order,
			Version:       raw.Version,
			Template: TemplateDefinition{
				Settings: raw.Settings,
				Mappings: raw.Mappings,
				Aliases:  raw.Aliases,
			},
		}
		templates[name] = tmpl
	}

	return templates, nil
}

// Get retrieves a specific index template (tries composable first, then legacy)
func Get(ctx context.Context, name string) (*Template, error) {
	var result map[string]any

	endpoint := fmt.Sprintf("_index_template/%s", name)
	resp, err := shared.Client.R().
		SetContext(ctx).
		SetHeader("Content-Type", "application/json").
		SetResult(&result).
		Get(endpoint)

	if err != nil {
		return nil, fmt.Errorf("failed to get template: %w", err)
	}

	if resp.StatusCode() == 200 {
		// Parse composable template response
		if indexTemplates, ok := result["index_templates"].([]any); ok && len(indexTemplates) > 0 {
			if itemMap, ok := indexTemplates[0].(map[string]any); ok {
				if tmplData, ok := itemMap["index_template"].(map[string]any); ok {
					tmplBytes, _ := json.Marshal(tmplData)
					var tmpl Template
					if err := json.Unmarshal(tmplBytes, &tmpl); err == nil {
						tmpl.Name = name
						return &tmpl, nil
					}
				}
			}
		}
	}

	// Fall back to legacy template endpoint
	return getLegacy(ctx, name)
}

// getLegacy retrieves a specific legacy index template
func getLegacy(ctx context.Context, name string) (*Template, error) {
	var result map[string]any

	endpoint := fmt.Sprintf("_template/%s", name)
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

	// Legacy response: { "name": { "order":..., "index_patterns":..., ... } }
	itemMap, ok := result[name].(map[string]any)
	if !ok {
		return nil, fmt.Errorf("failed to parse legacy template response")
	}

	tmplBytes, _ := json.Marshal(itemMap)
	var raw struct {
		Order         int            `json:"order"`
		IndexPatterns []string       `json:"index_patterns"`
		Settings      map[string]any `json:"settings"`
		Mappings      map[string]any `json:"mappings"`
		Aliases       map[string]any `json:"aliases"`
		Version       int            `json:"version"`
	}
	if err := json.Unmarshal(tmplBytes, &raw); err != nil {
		return nil, fmt.Errorf("failed to parse legacy template response: %w", err)
	}

	tmpl := &Template{
		Name:          name,
		IndexPatterns: raw.IndexPatterns,
		Priority:      raw.Order,
		Version:       raw.Version,
		Template: TemplateDefinition{
			Settings: raw.Settings,
			Mappings: raw.Mappings,
			Aliases:  raw.Aliases,
		},
	}
	return tmpl, nil
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
