package analysis

import (
	"context"
	"fmt"

	"github.com/pincher95/esctl/shared"
)

type AnalyzeRequest struct {
	Analyzer string   `json:"analyzer,omitempty"`
	Field    string   `json:"field,omitempty"`
	Text     []string `json:"text"`
}

type AnalyzeResponse struct {
	Tokens []any `json:"tokens"`
}

// Analyze executes the _analyze API.
func Analyze(ctx context.Context, index string, req AnalyzeRequest) (AnalyzeResponse, error) {
	var respData AnalyzeResponse
	endpoint := "_analyze"
	if index != "" {
		endpoint = fmt.Sprintf("%s/_analyze", index)
	}

	resp, err := shared.Client.R().
		SetContext(ctx).
		SetHeader("Content-Type", "application/json").
		SetBody(req).
		SetResult(&respData).
		Post(endpoint)
	if err != nil {
		return respData, err
	}
	if resp.StatusCode() != 200 {
		return respData, fmt.Errorf("analyze failed: %s", resp.Status())
	}
	return respData, nil
}

type ExplainRequest struct {
	Query map[string]any `json:"query"`
}

// Explain executes the _explain API.
func Explain(ctx context.Context, index, id string, req ExplainRequest) (map[string]any, error) {
	var respData map[string]any
	endpoint := fmt.Sprintf("%s/_explain/%s", index, id)
	resp, err := shared.Client.R().
		SetContext(ctx).
		SetHeader("Content-Type", "application/json").
		SetBody(req).
		SetResult(&respData).
		Post(endpoint)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode() != 200 {
		return nil, fmt.Errorf("explain failed: %s", resp.Status())
	}
	return respData, nil
}

// ProfileSearch executes a search with profiling enabled.
func ProfileSearch(ctx context.Context, index string, query map[string]any) (map[string]any, error) {
	var respData map[string]any
	body := map[string]any{
		"profile": true,
		"query":   query,
	}
	endpoint := fmt.Sprintf("%s/_search", index)
	resp, err := shared.Client.R().
		SetContext(ctx).
		SetHeader("Content-Type", "application/json").
		SetBody(body).
		SetResult(&respData).
		Post(endpoint)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode() != 200 {
		return nil, fmt.Errorf("profile search failed: %s", resp.Status())
	}
	return respData, nil
}
