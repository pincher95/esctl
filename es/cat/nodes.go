package cat

import (
	"context"
	"fmt"
	"net/url"
	"slices"
	"strings"

	"github.com/pincher95/esctl/shared"
)

type CatNodesResponse struct {
	// 32-bit scalars first (minimise padding)
	Port int `json:"port,string"`
	CPU  int `json:"cpu,string,omitempty"`

	// Embedded metric groups (contain a mix but large)
	DiskStats     // disk.*
	HeapStats     // heap.*
	RAMStats      // ram.*
	FileDescStats // file_desc.*
	LoadStats     // load_*.

	// Strings / reference fields
	ID          string `json:"id"`
	PID         string `json:"pid,omitempty"`
	IP          string `json:"ip"`
	HTTPAddress string `json:"http_address"`
	Version     string `json:"version"`
	Type        string `json:"type,omitempty"`
	Build       string `json:"build,omitempty"`
	JDK         string `json:"jdk,omitempty"`
	Uptime      string `json:"uptime,omitempty"`

	NodeIdentity // node.role(s), name, master
}

func (c *cat) CatNodes(ctx context.Context, endpoint, nodeName, bytes, timeUnit string) ([]CatNodesResponse, error) {
	if endpoint == "" {
		u := url.URL{Path: "_cat/nodes"}
		q := u.Query()
		q.Set("format", "json")
		q.Set("h", "name,ip,node.role,node.roles,master,heap.percent,cpu,load_1m,load_5m,load_15m,ram.percent")
		if bytes != "" {
			q.Set("bytes", bytes)
		}
		if timeUnit != "" {
			q.Set("time", timeUnit)
		}
		u.RawQuery = q.Encode()
		endpoint = u.String()
	}

	nodes := make([]CatNodesResponse, 0)

	resp, err := shared.Client.R().
		SetContext(ctx).
		SetHeader("Content-Type", "application/json").
		SetResult(&nodes).
		Get(endpoint)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode() != 200 {
		return nil, fmt.Errorf("failed to get nodes: %s", resp.Status())
	}

	if nodeName != "" {
		filtered := make([]CatNodesResponse, 0, len(nodes))
		for _, n := range nodes {
			if strings.Contains(n.Name, nodeName) {
				filtered = append(filtered, n)
			}
		}
		if len(filtered) == 0 {
			return nil, fmt.Errorf("node not found: %s", nodeName)
		}
		nodes = filtered
	}

	return nodes, nil
}

// RolesList returns the individual single-letter roles exposed by _cat/nodes.
func (n CatNodesResponse) RolesList() []string {
	var set = make(map[string]struct{})

	raw := n.Roles
	if raw == "" {
		raw = n.Role // fallback
	}

	if strings.Contains(raw, ",") {
		tokens := strings.SplitSeq(raw, ",")
		for t := range tokens {
			t = strings.TrimSpace(t)
			if t == "" {
				continue
			}
			set[roleAbbrev(t)] = struct{}{}
		}
	} else {
		// If no comma and length >1 we treat as sequence of single-letter codes
		if len(raw) > 1 {
			for _, r := range raw {
				set[string(r)] = struct{}{}
			}
		} else if raw != "" {
			set[raw] = struct{}{}
		}
	}

	delete(set, "r") // remote_cluster_client not counted in health output

	roles := make([]string, 0, len(set))
	for r := range set {
		roles = append(roles, r)
	}
	return roles
}

// roleAbbrev maps verbose role tokens to single-letter abbreviations.
func roleAbbrev(token string) string {
	switch token {
	case "data", "data_hot", "data_warm", "data_cold", "data_frozen", "data_content":
		return "d"
	case "master", "cluster_manager":
		return "m"
	case "coordinating_only":
		return "c"
	case "ingest":
		return "i"
	case "remote_cluster_client":
		return "r"
	case "search":
		return "s"
	default:
		if len(token) > 0 {
			return strings.ToLower(token[:1])
		}
		return token
	}
}

// HasRole checks whether the node advertises the given single-letter role (e.g. "d" for data).
func (n CatNodesResponse) HasRole(r string) bool {
	return slices.Contains(n.RolesList(), r)
}

// IsData returns true when the node has the data ("d") role.
func (n CatNodesResponse) IsData() bool { return n.HasRole("d") }

// IsMaster returns true when the node is master-eligible ("m") or currently elected (master == "*").
func (n CatNodesResponse) IsMaster() bool {
	return n.HasRole("m") || n.Master == "*" || n.ClusterManager == "*"
}
