package searchtemplate

import (
	"context"
	"fmt"
	"net/url"

	"github.com/pincher95/esctl/shared"
)

// SearchTemplate represents a stored search template
type SearchTemplate struct {
	ID       string         `json:"id,omitempty"`
	Template map[string]any `json:"template"`
}

// SearchTemplateResponse is the response when getting a search template
type SearchTemplateResponse struct {
	Found    bool           `json:"found"`
	ID       string         `json:"_id"`
	Template map[string]any `json:"template"`
}

// RenderResponse is the response from rendering a search template
type RenderResponse struct {
	TemplateOutput map[string]any `json:"template_output"`
}

// storedScript represents a script entry in the cluster state metadata
type storedScript struct {
	Lang   string `json:"lang"`
	Source string `json:"source"`
}

// clusterStateResponse represents the cluster state metadata response for stored scripts
type clusterStateResponse struct {
	Metadata struct {
		StoredScripts map[string]storedScript `json:"stored_scripts"`
	} `json:"metadata"`
}

// List retrieves all stored search templates via the cluster state API
func List(ctx context.Context) (map[string]SearchTemplateResponse, error) {
	var result clusterStateResponse
	resp, err := shared.Client.R().
		SetContext(ctx).
		SetHeader("Content-Type", "application/json").
		SetResult(&result).
		Get("/_cluster/state/metadata?filter_path=metadata.stored_scripts")

	if err != nil {
		return nil, fmt.Errorf("failed to list search templates: %w", err)
	}

	if resp.StatusCode() != 200 {
		return nil, fmt.Errorf("error listing search templates: %s", resp.Status())
	}

	// Filter to only mustache templates (search templates)
	templates := make(map[string]SearchTemplateResponse)
	for id, s := range result.Metadata.StoredScripts {
		if s.Lang == "mustache" {
			templates[id] = SearchTemplateResponse{
				Found: true,
				ID:    id,
				Template: map[string]any{
					"lang":   s.Lang,
					"source": s.Source,
				},
			}
		}
	}

	return templates, nil
}

// Get retrieves a specific stored search template by ID
func Get(ctx context.Context, id string) (*SearchTemplateResponse, error) {
	if id == "" {
		return nil, fmt.Errorf("search template ID is required")
	}

	escapedID := url.PathEscape(id)

	resp, err := shared.Client.R().
		SetContext(ctx).
		SetHeader("Content-Type", "application/json").
		SetResult(&SearchTemplateResponse{}).
		Get(fmt.Sprintf("/_scripts/%s", escapedID))

	if err != nil {
		return nil, fmt.Errorf("failed to get search template: %w", err)
	}

	if resp.StatusCode() == 404 {
		return nil, fmt.Errorf("search template '%s' not found", id)
	}

	if resp.StatusCode() != 200 {
		return nil, fmt.Errorf("error getting search template: %s", resp.Status())
	}

	result := resp.Result().(*SearchTemplateResponse)
	result.ID = id
	return result, nil
}

// Put creates or updates a stored search template
func Put(ctx context.Context, id string, template SearchTemplate) error {
	if id == "" {
		return fmt.Errorf("search template ID is required")
	}

	if template.Template == nil || len(template.Template) == 0 {
		return fmt.Errorf("template definition is required")
	}

	escapedID := url.PathEscape(id)

	// Search templates are stored as scripts with lang=mustache
	body := map[string]any{
		"script": map[string]any{
			"lang":   "mustache",
			"source": template.Template,
		},
	}

	resp, err := shared.Client.R().
		SetContext(ctx).
		SetHeader("Content-Type", "application/json").
		SetBody(body).
		Put(fmt.Sprintf("/_scripts/%s", escapedID))

	if err != nil {
		return fmt.Errorf("failed to put search template: %w", err)
	}

	if resp.StatusCode() != 200 {
		return fmt.Errorf("error putting search template: %s", resp.Status())
	}

	return nil
}

// Delete removes a stored search template
func Delete(ctx context.Context, id string) error {
	if id == "" {
		return fmt.Errorf("search template ID is required")
	}

	escapedID := url.PathEscape(id)

	resp, err := shared.Client.R().
		SetContext(ctx).
		SetHeader("Content-Type", "application/json").
		Delete(fmt.Sprintf("/_scripts/%s", escapedID))

	if err != nil {
		return fmt.Errorf("failed to delete search template: %w", err)
	}

	if resp.StatusCode() == 404 {
		return fmt.Errorf("search template '%s' not found", id)
	}

	if resp.StatusCode() != 200 {
		return fmt.Errorf("error deleting search template: %s", resp.Status())
	}

	return nil
}

// Render renders a search template with the given parameters
func Render(ctx context.Context, id string, params map[string]any) (*RenderResponse, error) {
	if id == "" {
		return nil, fmt.Errorf("search template ID is required")
	}

	body := map[string]any{
		"id": id,
	}
	if params != nil && len(params) > 0 {
		body["params"] = params
	}

	resp, err := shared.Client.R().
		SetContext(ctx).
		SetHeader("Content-Type", "application/json").
		SetBody(body).
		SetResult(&RenderResponse{}).
		Post("/_render/template")

	if err != nil {
		return nil, fmt.Errorf("failed to render search template: %w", err)
	}

	if resp.StatusCode() == 404 {
		return nil, fmt.Errorf("search template '%s' not found", id)
	}

	if resp.StatusCode() != 200 {
		return nil, fmt.Errorf("error rendering search template: %s", resp.Status())
	}

	result := resp.Result().(*RenderResponse)
	return result, nil
}

// RenderInline renders an inline search template (not stored)
func RenderInline(ctx context.Context, template map[string]any, params map[string]any) (*RenderResponse, error) {
	if template == nil || len(template) == 0 {
		return nil, fmt.Errorf("template is required")
	}

	body := map[string]any{
		"source": template,
	}
	if params != nil && len(params) > 0 {
		body["params"] = params
	}

	resp, err := shared.Client.R().
		SetContext(ctx).
		SetHeader("Content-Type", "application/json").
		SetBody(body).
		SetResult(&RenderResponse{}).
		Post("/_render/template")

	if err != nil {
		return nil, fmt.Errorf("failed to render inline template: %w", err)
	}

	if resp.StatusCode() != 200 {
		return nil, fmt.Errorf("error rendering inline template: %s", resp.Status())
	}

	result := resp.Result().(*RenderResponse)
	return result, nil
}
