package common

// FailuresCause captures nested error information returned by Elasticsearch for shard or node failures.
// It is shared across index, cluster and other packages so we only define it once.
type FailuresCause struct {
	Type   string `json:"type"`
	Reason string `json:"reason"`
	NodeID string `json:"node_id,omitempty"`
	Cause  *struct {
		Type   string `json:"type"`
		Reason string `json:"reason"`
		Cause  *struct {
			Type   string  `json:"type"`
			Reason *string `json:"reason"`
		} `json:"caused_by,omitempty"`
	} `json:"caused_by,omitempty"`
}
