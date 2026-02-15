package snapshots

import (
	"context"
	"fmt"
	"net/url"

	"github.com/pincher95/esctl/shared"
)

type Repository struct {
	Type     string         `json:"type"`
	Settings map[string]any `json:"settings"`
}

type RepositoryResponse map[string]Repository

type SnapshotInfo struct {
	Snapshot           string         `json:"snapshot"`
	UUID               string         `json:"uuid"`
	Repository         string         `json:"repository"`
	VersionID          int            `json:"version_id"`
	Version            string         `json:"version"`
	Indices            []string       `json:"indices"`
	DataStreams        []string       `json:"data_streams"`
	IncludeGlobalState bool           `json:"include_global_state"`
	State              string         `json:"state"`
	StartTime          string         `json:"start_time"`
	StartTimeInMillis  int64          `json:"start_time_in_millis"`
	EndTime            string         `json:"end_time"`
	EndTimeInMillis    int64          `json:"end_time_in_millis"`
	DurationInMillis   int64          `json:"duration_in_millis"`
	Failures           []any          `json:"failures"`
	Shards             SnapshotShards `json:"shards"`
	Metadata           map[string]any `json:"metadata"`
}

type SnapshotShards struct {
	Total      int `json:"total"`
	Failed     int `json:"failed"`
	Successful int `json:"successful"`
}

type SnapshotResponse struct {
	Snapshots []SnapshotInfo `json:"snapshots"`
}

type SnapshotStatusResponse struct {
	Snapshots []SnapshotStatusInfo `json:"snapshots"`
}

type SnapshotStatusInfo struct {
	Snapshot    string                   `json:"snapshot"`
	Repository  string                   `json:"repository"`
	UUID        string                   `json:"uuid"`
	State       string                   `json:"state"`
	ShardsStats SnapshotShards           `json:"shards_stats"`
	Stats       SnapshotStats            `json:"stats"`
	Indices     map[string]SnapshotIndex `json:"indices"`
}

type SnapshotStats struct {
	Incremental       SnapshotFileStats `json:"incremental"`
	Total             SnapshotFileStats `json:"total"`
	StartTimeInMillis int64             `json:"start_time_in_millis"`
	TimeInMillis      int64             `json:"time_in_millis"`
}

type SnapshotFileStats struct {
	FileCount   int64 `json:"file_count"`
	SizeInBytes int64 `json:"size_in_bytes"`
}

type SnapshotIndex struct {
	ShardsStats SnapshotShards           `json:"shards_stats"`
	Stats       SnapshotStats            `json:"stats"`
	Shards      map[string]SnapshotShard `json:"shards"`
}

type SnapshotShard struct {
	Stage string        `json:"stage"`
	Stats SnapshotStats `json:"stats"`
}

type CreateSnapshotRequest struct {
	Indices            string         `json:"indices,omitempty"`
	IgnoreUnavailable  bool           `json:"ignore_unavailable,omitempty"`
	IncludeGlobalState *bool          `json:"include_global_state,omitempty"`
	Partial            bool           `json:"partial,omitempty"`
	Metadata           map[string]any `json:"metadata,omitempty"`
}

type RestoreSnapshotRequest struct {
	Indices             string         `json:"indices,omitempty"`
	IgnoreUnavailable   bool           `json:"ignore_unavailable,omitempty"`
	IncludeGlobalState  bool           `json:"include_global_state,omitempty"`
	RenamePattern       string         `json:"rename_pattern,omitempty"`
	RenameReplacement   string         `json:"rename_replacement,omitempty"`
	IncludeAliases      bool           `json:"include_aliases,omitempty"`
	IndexSettings       map[string]any `json:"index_settings,omitempty"`
	IgnoreIndexSettings []string       `json:"ignore_index_settings,omitempty"`
}

// ListRepositories gets all snapshot repositories
func ListRepositories(ctx context.Context) (RepositoryResponse, error) {
	var result RepositoryResponse

	resp, err := shared.Client.R().
		SetContext(ctx).
		SetHeader("Content-Type", "application/json").
		SetResult(&result).
		Get("_snapshot")
	if err != nil {
		return nil, err
	}

	if resp.StatusCode() != 200 {
		return nil, fmt.Errorf("failed to list repositories: %s", resp.Status())
	}

	return result, nil
}

// GetRepository gets a specific snapshot repository
func GetRepository(ctx context.Context, repository string) (RepositoryResponse, error) {
	var result RepositoryResponse

	resp, err := shared.Client.R().
		SetContext(ctx).
		SetHeader("Content-Type", "application/json").
		SetResult(&result).
		Get(fmt.Sprintf("_snapshot/%s", repository))
	if err != nil {
		return nil, err
	}

	if resp.StatusCode() != 200 {
		return nil, fmt.Errorf("failed to get repository %s: %s", repository, resp.Status())
	}

	return result, nil
}

// CreateRepository creates a new snapshot repository
func CreateRepository(ctx context.Context, repository string, repoType string, settings map[string]any) error {
	body := Repository{
		Type:     repoType,
		Settings: settings,
	}

	resp, err := shared.Client.R().
		SetContext(ctx).
		SetHeader("Content-Type", "application/json").
		SetBody(body).
		Put(fmt.Sprintf("_snapshot/%s", repository))
	if err != nil {
		return err
	}

	if resp.StatusCode() != 200 {
		return fmt.Errorf("failed to create repository %s: %s", repository, resp.Status())
	}

	return nil
}

// DeleteRepository deletes a snapshot repository
func DeleteRepository(ctx context.Context, repository string) error {
	resp, err := shared.Client.R().
		SetContext(ctx).
		SetHeader("Content-Type", "application/json").
		Delete(fmt.Sprintf("_snapshot/%s", repository))
	if err != nil {
		return err
	}

	if resp.StatusCode() != 200 {
		return fmt.Errorf("failed to delete repository %s: %s", repository, resp.Status())
	}

	return nil
}

// ListSnapshots lists snapshots in a repository
func ListSnapshots(ctx context.Context, repository string) (SnapshotResponse, error) {
	var result SnapshotResponse

	resp, err := shared.Client.R().
		SetContext(ctx).
		SetHeader("Content-Type", "application/json").
		SetResult(&result).
		Get(fmt.Sprintf("_snapshot/%s/_all", repository))
	if err != nil {
		return result, err
	}

	if resp.StatusCode() != 200 {
		return result, fmt.Errorf("failed to list snapshots: %s", resp.Status())
	}

	return result, nil
}

// GetSnapshot gets a specific snapshot
func GetSnapshot(ctx context.Context, repository, snapshot string) (SnapshotResponse, error) {
	var result SnapshotResponse

	resp, err := shared.Client.R().
		SetContext(ctx).
		SetHeader("Content-Type", "application/json").
		SetResult(&result).
		Get(fmt.Sprintf("_snapshot/%s/%s", repository, snapshot))
	if err != nil {
		return result, err
	}

	if resp.StatusCode() != 200 {
		return result, fmt.Errorf("failed to get snapshot %s: %s", snapshot, resp.Status())
	}

	return result, nil
}

// CreateSnapshot creates a new snapshot
func CreateSnapshot(ctx context.Context, repository, snapshot string, request CreateSnapshotRequest, waitForCompletion bool) error {
	u := url.URL{Path: fmt.Sprintf("_snapshot/%s/%s", repository, snapshot)}
	q := u.Query()
	if waitForCompletion {
		q.Set("wait_for_completion", "true")
	}
	u.RawQuery = q.Encode()

	resp, err := shared.Client.R().
		SetContext(ctx).
		SetHeader("Content-Type", "application/json").
		SetBody(request).
		Put(u.String())
	if err != nil {
		return err
	}

	if resp.StatusCode() != 200 && resp.StatusCode() != 202 {
		return fmt.Errorf("failed to create snapshot %s: %s", snapshot, resp.Status())
	}

	return nil
}

// DeleteSnapshot deletes a snapshot
func DeleteSnapshot(ctx context.Context, repository, snapshot string) error {
	resp, err := shared.Client.R().
		SetContext(ctx).
		SetHeader("Content-Type", "application/json").
		Delete(fmt.Sprintf("_snapshot/%s/%s", repository, snapshot))
	if err != nil {
		return err
	}

	if resp.StatusCode() != 200 {
		return fmt.Errorf("failed to delete snapshot %s: %s", snapshot, resp.Status())
	}

	return nil
}

// SnapshotStatus gets the status of snapshots
func SnapshotStatus(ctx context.Context, repository, snapshot string) (SnapshotStatusResponse, error) {
	var result SnapshotStatusResponse
	var endpoint string

	if repository != "" && snapshot != "" {
		endpoint = fmt.Sprintf("_snapshot/%s/%s/_status", repository, snapshot)
	} else if repository != "" {
		endpoint = fmt.Sprintf("_snapshot/%s/_status", repository)
	} else {
		endpoint = "_snapshot/_status"
	}

	resp, err := shared.Client.R().
		SetContext(ctx).
		SetHeader("Content-Type", "application/json").
		SetResult(&result).
		Get(endpoint)
	if err != nil {
		return result, err
	}

	if resp.StatusCode() != 200 {
		return result, fmt.Errorf("failed to get snapshot status: %s", resp.Status())
	}

	return result, nil
}

// RestoreSnapshot restores a snapshot
func RestoreSnapshot(ctx context.Context, repository, snapshot string, request RestoreSnapshotRequest, waitForCompletion bool) error {
	u := url.URL{Path: fmt.Sprintf("_snapshot/%s/%s/_restore", repository, snapshot)}
	q := u.Query()
	if waitForCompletion {
		q.Set("wait_for_completion", "true")
	}
	u.RawQuery = q.Encode()

	resp, err := shared.Client.R().
		SetContext(ctx).
		SetHeader("Content-Type", "application/json").
		SetBody(request).
		Post(u.String())
	if err != nil {
		return err
	}

	if resp.StatusCode() != 200 && resp.StatusCode() != 202 {
		return fmt.Errorf("failed to restore snapshot %s: %s", snapshot, resp.Status())
	}

	return nil
}
