package cluster_test

import (
	"context"
	"testing"

	"github.com/pincher95/esctl/es/cluster"
	"github.com/pincher95/esctl/internal/testutil"
	"github.com/pincher95/esctl/shared"
)

func TestClusterPendingTasks(t *testing.T) {
	jsonResp := `{"tasks":[
		{"insert_order":101,"priority":"URGENT","source":"create-index [foo]","executing":true,"time_in_queue_millis":86,"time_in_queue":"86ms"},
		{"insert_order":102,"priority":"HIGH","source":"put-mapping","executing":false,"time_in_queue_millis":12,"time_in_queue":"12ms"}
	]}`
	srv, cli := testutil.NewMockServer(jsonResp, "/_cluster/pending_tasks")
	defer srv.Close()
	shared.SetClient(cli)

	resp, err := cluster.ClusterPendingTasks(context.Background())
	if err != nil {
		t.Fatalf("ClusterPendingTasks() error = %v", err)
	}
	if len(resp.Tasks) != 2 {
		t.Fatalf("expected 2 tasks, got %d", len(resp.Tasks))
	}
	if resp.Tasks[0].Priority != "URGENT" || !resp.Tasks[0].Executing {
		t.Errorf("unexpected first task: %+v", resp.Tasks[0])
	}
}
