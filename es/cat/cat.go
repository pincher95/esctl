package cat

type Cat interface {
	// // CatAliases is a wrapper for the `/_cat/aliases` endpoint
	// CatAliases(endpoint *string) (*CatAliasesResponse, error)
	// CatAllocation is a wrapper for the `/_cat/allocation` endpoint
	CatAllocation(endpoint, nodeID, bytes *string) (*[]CatAllocationResponse, error)
	// CatCount is a wrapper for the `/_cat/count` endpoint
	// CatCount(endpoint *string) (*CatCountResponse, error)
	// CatFielddata is a wrapper for the `/_cat/fielddata` endpoint
	CatFielddata(endpoint, fields, bytes *string) (*[]CatFielddataResponse, error)
	// // CatHealth is a wrapper for the `/_cat/health` endpoint
	// CatHealth(endpoint *string) (*CatHealthResponse, error)
	// CatIndices is a wrapper for the `/_cat/indices` endpoint
	CatIndices(endpoint, index, bytes *string) (*[]CatIndiceResponse, error)
	// // CatMaster is a wrapper for the `/_cat/master` endpoint
	// CatMaster(endpoint *string) (*CatMasterResponse, error)
	// CatNodes is a wrapper for the `/_cat/nodes` endpoint
	CatNodes(endpoint, nodeName, bytes, time *string) (*[]CatNodesResponse, error)
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
