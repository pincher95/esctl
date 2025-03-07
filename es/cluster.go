package es

// type ClusterHealth struct {
// 	ClusterName                 string  `json:"cluster_name" yaml:"clusterName"`
// 	Status                      string  `json:"status" yaml:"status"`
// 	TimedOut                    bool    `json:"timed_out" yaml:"timedOut"`
// 	NumberOfNodes               int     `json:"number_of_nodes" yaml:"numberOfNodes"`
// 	NumberOfDataNodes           int     `json:"number_of_data_nodes" yaml:"numberOfDataNodes"`
// 	ActivePrimaryShards         int     `json:"active_primary_shards" yaml:"activePrimaryShards"`
// 	ActiveShards                int     `json:"active_shards" yaml:"activeShards"`
// 	RelocatingShards            int     `json:"relocating_shards" yaml:"relocatingShards"`
// 	InitializingShards          int     `json:"initializing_shards" yaml:"initializingShards"`
// 	UnassignedShards            int     `json:"unassigned_shards" yaml:"unassignedShards"`
// 	DelayedUnassignedShards     int     `json:"delayed_unassigned_shards" yaml:"delayedUnassignedShards"`
// 	NumberOfPendingTasks        int     `json:"number_of_pending_tasks" yaml:"numberOfPendingTasks"`
// 	NumberOfInFlightFetch       int     `json:"number_of_in_flight_fetch" yaml:"numberOfInFlightFetch"`
// 	TaskMaxWaitingInQueueMillis int     `json:"task_max_waiting_in_queue_millis" yaml:"taskMaxWaitingInQueueMillis"`
// 	ActiveShardsPercentAsNumber float64 `json:"active_shards_percent_as_number" yaml:"activeShardsPercentAsNumber"`
// }

// type ClusterStats struct {
// 	ClusterUUID string `json:"cluster_uuid" yaml:"clusterUUID"`
// 	Nodes       struct {
// 		Total      int `json:"total" yaml:"total"`
// 		Successful int `json:"successful" yaml:"successful"`
// 		Failed     int `json:"failed" yaml:"failed"`
// 	} `json:"_nodes" yaml:"nodes"`
// 	Indices struct {
// 		Count  int `json:"count" yaml:"count"`
// 		Shards struct {
// 			Total       int     `json:"total" yaml:"total"`
// 			Primaries   int     `json:"primaries" yaml:"primaries"`
// 			Replication float64 `json:"replication" yaml:"replication"`
// 			Index       struct {
// 				Shards struct {
// 					Min float64 `json:"min" yaml:"min"`
// 					Max float64 `json:"max" yaml:"max"`
// 					Avg float64 `json:"avg" yaml:"avg"`
// 				} `json:"shards" yaml:"shards"`
// 				Primaries struct {
// 					Min float64 `json:"min" yaml:"min"`
// 					Max float64 `json:"max" yaml:"max"`
// 					Avg float64 `json:"avg" yaml:"avg"`
// 				} `json:"primaries" yaml:"primaries"`
// 				Replication struct {
// 					Min float64 `json:"min" yaml:"min"`
// 					Max float64 `json:"max" yaml:"max"`
// 					Avg float64 `json:"avg" yaml:"avg"`
// 				} `json:"replication" yaml:"replication"`
// 			} `json:"index" yaml:"index"`
// 		} `json:"shards" yaml:"shards"`
// 		Store struct {
// 			SizeInBytes             int `json:"size_in_bytes" yaml:"sizeInBytes"`
// 			TotalDataSetSizeInBytes int `json:"total_data_set_size_in_bytes" yaml:"totalDataSetSizeInBytes"`
// 			ReservedInBytes         int `json:"reserved_in_bytes" yaml:"reservedInBytes"`
// 		} `json:"store" yaml:"store"`
// 	} `json:"indices" yaml:"indices"`
// }

// type ClusterSettings map[string]any

// type Cluster struct {
// 	Health   ClusterHealth   `json:"health"`
// 	Stats    ClusterStats    `json:"stats"`
// 	Settings ClusterSettings `json:"settings"`
// }

// func GetCluster() (*Cluster, error) {
// 	var cluster Cluster
// 	if err := getJSONResponse("_cluster/stats", &cluster.Stats); err != nil {
// 		return nil, err
// 	}

// 	if err := getJSONResponse("_cluster/health", &cluster.Health); err != nil {
// 		return nil, err
// 	}

// 	if err := getJSONResponse("_cluster/settings", &cluster.Settings); err != nil {
// 		return nil, err
// 	}

// 	return &cluster, nil
// }

// func GetClusterSettings(flat, defaults bool) (ClusterSettings, error) {
// 	var settings ClusterSettings

// 	endpoint := "_cluster/settings"

// 	if defaults && !flat {
// 		endpoint += fmt.Sprintf("?%s&%s", "flat_settings", "include_defaults")
// 	} else if defaults {
// 		endpoint += fmt.Sprintf("?%s", "include_defaults")
// 	} else if !flat {
// 		endpoint += fmt.Sprintf("?%s", "flat_settings")
// 	}

// 	if err := getJSONResponse(endpoint, &settings); err != nil {
// 		return nil, err
// 	}

// 	return settings, nil
// }

// func GetClusterHealth(level, index string) (ClusterSettings, error) {
// 	var settings ClusterSettings

// 	endpoint := "_cluster/health"

// 	if index != "" {
// 		endpoint += fmt.Sprintf("/%s", index)
// 	}

// 	if level != "" {
// 		endpoint += fmt.Sprintf("?%s=%s", "level", level)
// 	}

// 	if err := getJSONResponse(endpoint, &settings); err != nil {
// 		return nil, err
// 	}

// 	return settings, nil
// }

// func GetClusterStats() (ClusterStats, error) {
// 	var stats ClusterStats

// 	if err := getJSONResponse("_cluster/stats", &stats); err != nil {
// 		return stats, err
// 	}

// 	return stats, nil
// }

// func GetAllocationExplain(flagIncludeDiskInfo, flagIncludeYesDecisions bool) (any, error) {
// 	var response JsonResponse

// 	endpoint := "_cluster/allocation/explain"

// 	if flagIncludeDiskInfo && flagIncludeYesDecisions {
// 		endpoint += fmt.Sprintf("?%s&%s", "include_disk_info", "include_yes_decisions")
// 	} else if flagIncludeDiskInfo {
// 		endpoint += fmt.Sprintf("?%s", "include_disk_info")
// 	} else if flagIncludeYesDecisions {
// 		endpoint += fmt.Sprintf("?%s", "include_yes_decisions")
// 	}

// 	if err := getJSONResponse(endpoint, &response); err != nil {
// 		return nil, err
// 	}

// 	return response, nil
// }

// func UpdateReroute(flagDryRun, flagExplain, flagRetryFailed bool, flagMertic string) (any, error) {
// 	var response JsonResponse

// 	endpoint := fmt.Sprintf("_cluster/reroute?dry_run=%t&explain=%t&retry_failed=%t&metric=%s", flagDryRun, flagExplain, flagRetryFailed, flagMertic)

// 	if err := postWithoutBody(endpoint, &response); err != nil {
// 		return nil, err
// 	}

// 	return response, nil
// }
