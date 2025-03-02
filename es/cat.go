package es

// type NodeDetails map[string]NodeSummary

// type NodeSummary struct {
// 	IP      string                  `json:"ip"`
// 	Role    string                  `json:"role"`
// 	Master  string                  `json:"master"`
// 	Indices map[string]IndexSummary `json:"indices"`
// }

// type IndexSummary struct {
// 	Health       string                  `json:"health"`
// 	Pri          string                  `json:"primary"`
// 	Rep          string                  `json:"replica"`
// 	PriStoreSize string                  `json:"pri-store-size"`
// 	Shards       map[string]ShardSummary `json:"shards"`
// }

// type ShardSummary struct {
// 	PriRep        string `json:"pri-rep"`
// 	State         string `json:"state"`
// 	SegmentsCount string `json:"segments-count"`
// }

// func GetNodeDetails(nodeName string) (*NodeDetails, error) {
// 	nodes, err := GetNodes("")
// 	if err != nil {
// 		return nil, err
// 	}

// 	nodeDetails := make(NodeDetails)

// 	for _, node := range nodes {
// 		if nodeName != "" && node.Name != nodeName {
// 			continue
// 		}

// 		nodeSummary, err := getNodeSummary(node)
// 		if err != nil {
// 			return nil, err
// 		}

// 		nodeDetails[node.Name] = nodeSummary
// 	}

// 	return &nodeDetails, nil
// }

// func getNodeSummary(node Node) (NodeSummary, error) {
// 	indices, err := GetIndices("")
// 	if err != nil {
// 		return NodeSummary{}, err
// 	}

// 	shards, err := GetShards("")
// 	if err != nil {
// 		return NodeSummary{}, err
// 	}

// 	indexSummaries := make(map[string]IndexSummary)
// 	for _, index := range indices {
// 		shardSummaries := make(map[string]ShardSummary)
// 		for _, shard := range shards {
// 			if shard.Index == index.Index && shard.Node == node.Name {
// 				shardSummary := ShardSummary{
// 					PriRep:        shard.PriRep,
// 					State:         shard.State,
// 					SegmentsCount: shard.SegmentsCount,
// 				}
// 				shardSummaries[shard.Shard] = shardSummary
// 			}
// 		}

// 		indexSummary := IndexSummary{
// 			Health:       index.Health,
// 			Pri:          index.Pri,
// 			Rep:          index.Rep,
// 			PriStoreSize: index.PriStoreSize,
// 			Shards:       shardSummaries,
// 		}

// 		indexSummaries[index.Index] = indexSummary
// 	}

// 	nodeSummary := NodeSummary{
// 		IP:      node.IP,
// 		Role:    node.NodeRole,
// 		Master:  node.Master,
// 		Indices: indexSummaries,
// 	}

// 	return nodeSummary, nil
// }
