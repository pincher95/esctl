package cat

// DiskStats groups disk space related fields reported by _cat APIs.
type DiskStats struct {
	DiskTotal       string `json:"disk.total,omitempty"`
	DiskUsed        string `json:"disk.used,omitempty"`
	DiskAvail       string `json:"disk.avail,omitempty"`
	DiskUsedPercent string `json:"disk.used_percent,omitempty"`
}

// HeapStats groups JVM heap metrics.
type HeapStats struct {
	HeapCurrent string `json:"heap.current,omitempty"`
	HeapPercent int    `json:"heap.percent,string,omitempty"`
	HeapMax     string `json:"heap.max,omitempty"`
}

// RAMStats groups RAM usage metrics.
type RAMStats struct {
	RAMCurrent string `json:"ram.current,omitempty"`
	RAMPercent int    `json:"ram.percent,string,omitempty"`
	RAMMax     string `json:"ram.max,omitempty"`
}

// FileDescStats groups file descriptor metrics.
type FileDescStats struct {
	FileDescCurrent int `json:"file_desc.current,string,omitempty"`
	FileDescPercent int `json:"file_desc.percent,string,omitempty"`
	FileDescMax     int `json:"file_desc.max,string,omitempty"`
}

// LoadStats groups system load averages.
type LoadStats struct {
	Load1M  string `json:"load_1m,omitempty"`
	Load5M  string `json:"load_5m,omitempty"`
	Load15M string `json:"load_15m,omitempty"`
}

// NodeIdentity groups node role/name related fields.
type NodeIdentity struct {
	Role           string `json:"node.role"`
	Roles          string `json:"node.roles"`
	Master         string `json:"master"`
	ClusterManager string `json:"cluster_manager"`
	Name           string `json:"name"`
}

// ShardIdentity groups identifiers and core state of a shard.
type ShardIdentity struct {
	Index  string `json:"index"`
	Shard  int    `json:"shard,string"`
	Prirep string `json:"prirep"`
	State  string `json:"state"`
}

// ShardStorageStats groups document and size metrics.
type ShardStorageStats struct {
	Docs  string `json:"docs,omitempty"`
	Store string `json:"store,omitempty"`
}

// ShardNodeInfo groups node location information for a shard.
type ShardNodeInfo struct {
	IP   string `json:"ip,omitempty"`
	ID   string `json:"id,omitempty"`
	Node string `json:"node,omitempty"`
}

// ShardUnassignedInfo groups unassigned shard metadata.
type ShardUnassignedInfo struct {
	UnassignedReason  string `json:"unassigned.reason,omitempty"`
	UnassignedAt      string `json:"unassigned.at,omitempty"`
	UnassignedFor     string `json:"unassigned.for,omitempty"`
	UnassignedDetails string `json:"unassigned.details,omitempty"`
}
