package template

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/pincher95/esctl/shared"
)

// ComponentTemplate represents a component template
type ComponentTemplate struct {
	Name     string                 `json:"-"`
	Version  int                    `json:"version,omitempty"`
	Template ComponentDefinition    `json:"template"`
	Meta     map[string]interface{} `json:"_meta,omitempty"`
}

// ComponentDefinition contains component template settings, mappings, and aliases
type ComponentDefinition struct {
	Settings map[string]interface{} `json:"settings,omitempty"`
	Mappings map[string]interface{} `json:"mappings,omitempty"`
	Aliases  map[string]interface{} `json:"aliases,omitempty"`
}

// ComponentListResponse represents the response from listing component templates
type ComponentListResponse map[string]ComponentTemplate

// ListComponents retrieves all component templates
func ListComponents(ctx context.Context) (ComponentListResponse, error) {
	var result map[string]interface{}

	resp, err := shared.Client.R().
		SetContext(ctx).
		SetHeader("Content-Type", "application/json").
		SetResult(&result).
		Get("_component_template")

	if err != nil {
		return nil, fmt.Errorf("failed to list component templates: %w", err)
	}

	if resp.StatusCode() != 200 {
		return nil, fmt.Errorf("failed to list component templates: %s", resp.Status())
	}

	// Parse the response
	templates := make(ComponentListResponse)
	if componentTemplates, ok := result["component_templates"].([]interface{}); ok {
		for _, item := range componentTemplates {
			if itemMap, ok := item.(map[string]interface{}); ok {
				if name, ok := itemMap["name"].(string); ok {
					if tmplData, ok := itemMap["component_template"].(map[string]interface{}); ok {
						tmplBytes, _ := json.Marshal(tmplData)
						var tmpl ComponentTemplate
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

// GetComponent retrieves a specific component template
func GetComponent(ctx context.Context, name string) (*ComponentTemplate, error) {
	var result map[string]interface{}

	endpoint := fmt.Sprintf("_component_template/%s", name)
	resp, err := shared.Client.R().
		SetContext(ctx).
		SetHeader("Content-Type", "application/json").
		SetResult(&result).
		Get(endpoint)

	if err != nil {
		return nil, fmt.Errorf("failed to get component template: %w", err)
	}

	if resp.StatusCode() == 404 {
		return nil, fmt.Errorf("component template not found: %s", name)
	}

	if resp.StatusCode() != 200 {
		return nil, fmt.Errorf("failed to get component template: %s", resp.Status())
	}

	// Parse the response
	if componentTemplates, ok := result["component_templates"].([]interface{}); ok && len(componentTemplates) > 0 {
		if itemMap, ok := componentTemplates[0].(map[string]interface{}); ok {
			if tmplData, ok := itemMap["component_template"].(map[string]interface{}); ok {
				tmplBytes, _ := json.Marshal(tmplData)
				var tmpl ComponentTemplate
				if err := json.Unmarshal(tmplBytes, &tmpl); err == nil {
					tmpl.Name = name
					return &tmpl, nil
				}
			}
		}
	}

	return nil, fmt.Errorf("failed to parse component template response")
}

// PutComponent creates or updates a component template
func PutComponent(ctx context.Context, name string, template ComponentTemplate) error {
	endpoint := fmt.Sprintf("_component_template/%s", name)

	resp, err := shared.Client.R().
		SetContext(ctx).
		SetHeader("Content-Type", "application/json").
		SetBody(template).
		Put(endpoint)

	if err != nil {
		return fmt.Errorf("failed to put component template: %w", err)
	}

	if resp.StatusCode() != 200 {
		return fmt.Errorf("failed to put component template: %s - %s", resp.Status(), string(resp.Body()))
	}

	return nil
}

// DeleteComponent removes a component template
func DeleteComponent(ctx context.Context, name string) error {
	endpoint := fmt.Sprintf("_component_template/%s", name)

	resp, err := shared.Client.R().
		SetContext(ctx).
		SetHeader("Content-Type", "application/json").
		Delete(endpoint)

	if err != nil {
		return fmt.Errorf("failed to delete component template: %w", err)
	}

	if resp.StatusCode() == 404 {
		return fmt.Errorf("component template not found: %s", name)
	}

	if resp.StatusCode() != 200 {
		return fmt.Errorf("failed to delete component template: %s", resp.Status())
	}

	return nil
}

// ExistsComponent checks if a component template exists
func ExistsComponent(ctx context.Context, name string) (bool, error) {
	endpoint := fmt.Sprintf("_component_template/%s", name)

	resp, err := shared.Client.R().
		SetContext(ctx).
		Head(endpoint)

	if err != nil {
		return false, fmt.Errorf("failed to check component template existence: %w", err)
	}

	return resp.StatusCode() == 200, nil
}
