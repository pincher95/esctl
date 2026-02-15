package datastream

import (
	"context"
	"fmt"
	"net/url"

	"github.com/pincher95/esctl/shared"
)

// DataStream represents a data stream
type DataStream struct {
	Name           string            `json:"name"`
	TimestampField TimestampField    `json:"timestamp_field"`
	Indices        []DataStreamIndex `json:"indices"`
	Generation     int               `json:"generation"`
	Status         string            `json:"status"`
	Template       string            `json:"template"`
	ILMPolicy      string            `json:"ilm_policy,omitempty"`
	Hidden         bool              `json:"hidden"`
}

// TimestampField contains the timestamp field configuration
type TimestampField struct {
	Name string `json:"name"`
}

// DataStreamIndex represents an index in a data stream
type DataStreamIndex struct {
	IndexName string `json:"index_name"`
	IndexUUID string `json:"index_uuid"`
}

// ListResponse represents the response from listing data streams
type ListResponse struct {
	DataStreams []DataStream `json:"data_streams"`
}

// List retrieves all data streams or filters by name pattern
func List(ctx context.Context, name string) ([]DataStream, error) {
	u := url.URL{}
	if name != "" {
		u.Path = fmt.Sprintf("_data_stream/%s", name)
	} else {
		u.Path = "_data_stream"
	}

	var result ListResponse
	resp, err := shared.Client.R().
		SetContext(ctx).
		SetHeader("Content-Type", "application/json").
		SetResult(&result).
		Get(u.String())

	if err != nil {
		return nil, fmt.Errorf("failed to list data streams: %w", err)
	}

	if resp.StatusCode() != 200 {
		return nil, fmt.Errorf("failed to list data streams: %s", resp.Status())
	}

	return result.DataStreams, nil
}

// Get retrieves a specific data stream
func Get(ctx context.Context, name string) (*DataStream, error) {
	endpoint := fmt.Sprintf("_data_stream/%s", name)

	var result ListResponse
	resp, err := shared.Client.R().
		SetContext(ctx).
		SetHeader("Content-Type", "application/json").
		SetResult(&result).
		Get(endpoint)

	if err != nil {
		return nil, fmt.Errorf("failed to get data stream: %w", err)
	}

	if resp.StatusCode() == 404 {
		return nil, fmt.Errorf("data stream not found: %s", name)
	}

	if resp.StatusCode() != 200 {
		return nil, fmt.Errorf("failed to get data stream: %s", resp.Status())
	}

	if len(result.DataStreams) == 0 {
		return nil, fmt.Errorf("data stream not found: %s", name)
	}

	return &result.DataStreams[0], nil
}

// Delete removes a data stream
func Delete(ctx context.Context, name string) error {
	endpoint := fmt.Sprintf("_data_stream/%s", name)

	resp, err := shared.Client.R().
		SetContext(ctx).
		SetHeader("Content-Type", "application/json").
		Delete(endpoint)

	if err != nil {
		return fmt.Errorf("failed to delete data stream: %w", err)
	}

	if resp.StatusCode() == 404 {
		return fmt.Errorf("data stream not found: %s", name)
	}

	if resp.StatusCode() != 200 {
		return fmt.Errorf("failed to delete data stream: %s", resp.Status())
	}

	return nil
}

// RolloverResponse represents the response from a rollover operation
type RolloverResponse struct {
	Acknowledged       bool   `json:"acknowledged"`
	ShardsAcknowledged bool   `json:"shards_acknowledged"`
	OldIndex           string `json:"old_index"`
	NewIndex           string `json:"new_index"`
	RolledOver         bool   `json:"rolled_over"`
	DryRun             bool   `json:"dry_run"`
}

// Rollover creates a new backing index for the data stream
func Rollover(ctx context.Context, name string) (*RolloverResponse, error) {
	endpoint := fmt.Sprintf("_data_stream/%s/_rollover", name)

	var result RolloverResponse
	resp, err := shared.Client.R().
		SetContext(ctx).
		SetHeader("Content-Type", "application/json").
		SetResult(&result).
		Post(endpoint)

	if err != nil {
		return nil, fmt.Errorf("failed to rollover data stream: %w", err)
	}

	if resp.StatusCode() != 200 {
		return nil, fmt.Errorf("failed to rollover data stream: %s", resp.Status())
	}

	return &result, nil
}

// StatsResponse represents the response from data stream stats
type StatsResponse struct {
	Shards      map[string]any    `json:"_shards"`
	DataStreams []DataStreamStats `json:"data_streams"`
}

// DataStreamStats contains statistics for a data stream
type DataStreamStats struct {
	DataStream       string `json:"data_stream"`
	BackingIndices   int    `json:"backing_indices"`
	StoreSize        string `json:"store_size,omitempty"`
	StoreSizeBytes   int64  `json:"store_size_bytes"`
	MaximumTimestamp int64  `json:"maximum_timestamp"`
}

// GetStats retrieves statistics for one or more data streams
func GetStats(ctx context.Context, name string) (*StatsResponse, error) {
	u := url.URL{}
	if name != "" {
		u.Path = fmt.Sprintf("_data_stream/%s/_stats", name)
	} else {
		u.Path = "_data_stream/_stats"
	}

	var result StatsResponse
	resp, err := shared.Client.R().
		SetContext(ctx).
		SetHeader("Content-Type", "application/json").
		SetResult(&result).
		Get(u.String())

	if err != nil {
		return nil, fmt.Errorf("failed to get data stream stats: %w", err)
	}

	if resp.StatusCode() != 200 {
		return nil, fmt.Errorf("failed to get data stream stats: %s", resp.Status())
	}

	return &result, nil
}
