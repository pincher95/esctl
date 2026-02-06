package pipeline

import (
	"context"
	"fmt"

	"github.com/pincher95/esctl/shared"
)

type Pipeline struct {
	Description string                 `json:"description,omitempty"`
	Processors  []map[string]any       `json:"processors"`
	OnFailure   []map[string]any       `json:"on_failure,omitempty"`
	Version     int                    `json:"version,omitempty"`
	Meta        map[string]interface{} `json:"_meta,omitempty"`
}

type PipelineResponse map[string]Pipeline

type SimulateRequest struct {
	Pipeline *Pipeline `json:"pipeline,omitempty"`
	Docs     []SimDoc  `json:"docs"`
}

type SimDoc struct {
	ID     string         `json:"_id,omitempty"`
	Index  string         `json:"_index,omitempty"`
	Source map[string]any `json:"_source"`
}

type SimulateResponse struct {
	Docs []SimResult `json:"docs"`
}

type SimResult struct {
	Doc              SimDoc            `json:"doc"`
	ProcessorResults []ProcessorResult `json:"processor_results,omitempty"`
}

type ProcessorResult struct {
	ProcessorType string  `json:"processor_type,omitempty"`
	Status        string  `json:"status,omitempty"`
	Doc           *SimDoc `json:"doc,omitempty"`
}

// ListPipelines lists all ingest pipelines
func ListPipelines(ctx context.Context) (PipelineResponse, error) {
	var result PipelineResponse
	resp, err := shared.Client.R().
		SetContext(ctx).
		SetHeader("Content-Type", "application/json").
		SetResult(&result).
		Get("_ingest/pipeline")
	if err != nil {
		return nil, err
	}
	if resp.StatusCode() != 200 {
		return nil, fmt.Errorf("failed to list pipelines: %s", resp.Status())
	}
	return result, nil
}

// GetPipeline gets a specific pipeline
func GetPipeline(ctx context.Context, id string) (PipelineResponse, error) {
	var result PipelineResponse
	resp, err := shared.Client.R().
		SetContext(ctx).
		SetHeader("Content-Type", "application/json").
		SetResult(&result).
		Get(fmt.Sprintf("_ingest/pipeline/%s", id))
	if err != nil {
		return nil, err
	}
	if resp.StatusCode() != 200 {
		return nil, fmt.Errorf("failed to get pipeline: %s", resp.Status())
	}
	return result, nil
}

// PutPipeline creates or updates a pipeline
func PutPipeline(ctx context.Context, id string, pipeline Pipeline) error {
	resp, err := shared.Client.R().
		SetContext(ctx).
		SetHeader("Content-Type", "application/json").
		SetBody(pipeline).
		Put(fmt.Sprintf("_ingest/pipeline/%s", id))
	if err != nil {
		return err
	}
	if resp.StatusCode() != 200 {
		return fmt.Errorf("failed to put pipeline: %s", resp.Status())
	}
	return nil
}

// DeletePipeline deletes a pipeline
func DeletePipeline(ctx context.Context, id string) error {
	resp, err := shared.Client.R().
		SetContext(ctx).
		SetHeader("Content-Type", "application/json").
		Delete(fmt.Sprintf("_ingest/pipeline/%s", id))
	if err != nil {
		return err
	}
	if resp.StatusCode() != 200 {
		return fmt.Errorf("failed to delete pipeline: %s", resp.Status())
	}
	return nil
}

// SimulatePipeline simulates pipeline execution
func SimulatePipeline(ctx context.Context, pipelineID string, request SimulateRequest) (SimulateResponse, error) {
	var result SimulateResponse
	endpoint := "_ingest/pipeline/_simulate"
	if pipelineID != "" {
		endpoint = fmt.Sprintf("_ingest/pipeline/%s/_simulate", pipelineID)
	}
	resp, err := shared.Client.R().
		SetContext(ctx).
		SetHeader("Content-Type", "application/json").
		SetBody(request).
		SetResult(&result).
		Post(endpoint)
	if err != nil {
		return result, err
	}
	if resp.StatusCode() != 200 {
		return result, fmt.Errorf("failed to simulate pipeline: %s", resp.Status())
	}
	return result, nil
}
