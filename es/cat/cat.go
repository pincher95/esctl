package cat

import (
	"fmt"

	"github.com/pincher95/esctl/shared"
)

type Cat interface {
	// // CatAliases is a wrapper for the `/_cat/aliases` endpoint
	// CatAliases(endpoint *string) (*CatAliasesResponse, error)
	// // CatAllocation is a wrapper for the `/_cat/allocation` endpoint
	// CatAllocation(endpoint *string) (*CatAllocationResponse, error)
	// // CatCount is a wrapper for the `/_cat/count` endpoint
	// CatCount(endpoint *string) (*CatCountResponse, error)
	// // CatFielddata is a wrapper for the `/_cat/fielddata` endpoint
	// CatFielddata(endpoint *string) (*CatFielddataResponse, error)
	// // CatHealth is a wrapper for the `/_cat/health` endpoint
	// CatHealth(endpoint *string) (*CatHealthResponse, error)
	// // CatIndices is a wrapper for the `/_cat/indices` endpoint
	// CatIndices(endpoint *string) (*CatIndicesResponse, error)
	// // CatMaster is a wrapper for the `/_cat/master` endpoint
	// CatMaster(endpoint *string) (*CatMasterResponse, error)
	// // CatNodes is a wrapper for the `/_cat/nodes` endpoint
	// CatNodes(endpoint *string) (*CatNodesResponse, error)
	// // CatPendingTasks is a wrapper for the `/_cat/pending_tasks` endpoint
	// CatPendingTasks(endpoint *string) (*CatPendingTasksResponse, error)
	// // CatRecovery is a wrapper for the `/_cat/recovery` endpoint
	// CatRecovery(endpoint *string) (*CatRecoveryResponse, error)
	// // CatRepositories is a wrapper for the `/_cat/repositories` endpoint
	// CatRepositories(endpoint *string) (*CatRepositoriesResponse, error)
	// // CatSegments is a wrapper for the `/_cat/segments` endpoint
	// CatSegments(endpoint *string) (*CatSegmentsResponse, error)
	// // CatShards is a wrapper for the `/_cat/shards` endpoint
	// CatShards(endpoint *string) (*CatShardsResponse, error)
	// CatSnapshots is a wrapper for the `/_cat/snapshots` endpoint
	CatSnapshots(endpoint, repository *string) (*[]CatSnapshotResponse, error)
	// CatTasks is a wrapper for the `/_cat/tasks` endpoint
	// 	CatTasks(endpoint *string) (*CatTasksResponse, error)
	// 	// CatTemplates is a wrapper for the `/_cat/templates` endpoint
	// 	CatTemplates(endpoint *string) (*CatTemplatesResponse, error)
	// 	// CatThreadPool is a wrapper for the `/_cat/thread_pool` endpoint
	// 	CatThreadPool(endpoint *string) (*CatThreadPoolResponse, error)
}

type cat struct {
	Cat
}

func NewCat() Cat {
	return &cat{}
}

type CatSnapshotResponse struct {
	ID               string `json:"id"`
	Status           string `json:"status"`
	StartEpoch       int    `json:"start_epoch,string"`
	StartTime        string `json:"start_time"`
	EndEpoch         int    `json:"end_epoch,string"`
	EndTime          string `json:"end_time"`
	Duration         string `json:"duration"`
	Indices          int    `json:"indices,string"`
	SuccessfulShards int    `json:"successful_shards,string"`
	FailedShards     int    `json:"failed_shards,string"`
	TotalShards      int    `json:"total_shards,string"`
	Reason           string `json:"reason"`
}

func (c *cat) CatSnapshots(endpoint, repository *string) (*[]CatSnapshotResponse, error) {
	if endpoint == nil {
		endpoint = new(string)
		*endpoint = "_cat/snapshots?format=json"
	}

	if *repository != "" {
		*endpoint = fmt.Sprintf("_cat/snapshots/%s?format=json", *repository)
	}

	snapshots := make([]CatSnapshotResponse, 0)

	resp, err := shared.Client.R().SetHeader("Content-Type", "application/json").SetResult(&snapshots).Get(*endpoint)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode() != 200 {
		return nil, fmt.Errorf("failed to get snapshots: %s", resp.Status())
	}

	return &snapshots, nil
}
