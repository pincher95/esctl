package get

import "time"

var (
	flagActions             []string
	flagColumns             []string
	flagIndex               string
	flagNode                string
	flagSortBy              string
	flagRefreshInterval     time.Duration
	flagShard               int
	flagInitializing        bool
	flagPrimary             bool
	flagRelocating          bool
	flagReplica             bool
	flagStarted             bool
	flagUnassigned          bool
	flagRefresh             bool
	flagIncludeDiskInfo     bool
	flagIncludeYesDecisions bool
)
