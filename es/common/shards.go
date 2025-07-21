package common

// ResponseShards represents the common _shards object returned by many Elasticsearch APIs.
type ResponseShards struct {
	Total      int                     `json:"total"`
	Successful int                     `json:"successful"`
	Failed     int                     `json:"failed"`
	Failures   []ResponseShardsFailure `json:"failures"`
	Skipped    int                     `json:"skipped"`
}

// ResponseShardsFailure captures a single failure entry inside the _shards.failures array.
type ResponseShardsFailure struct {
	Shard  int           `json:"shard"`
	Index  any           `json:"index"`
	Node   string        `json:"node"`
	Reason FailuresCause `json:"reason"`
}

// FailuresShard is used by some responses that flatten shard failure information.
type FailuresShard struct {
	Shard  int           `json:"shard"`
	Index  string        `json:"index"`
	Status string        `json:"status"`
	Reason FailuresCause `json:"reason"`
}
