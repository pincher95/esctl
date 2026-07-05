package snapshots

import (
	"context"
	"encoding/json"
	"strings"
	"testing"

	"github.com/pincher95/esctl/internal/testutil"
	"github.com/pincher95/esctl/shared"
)

func TestRestoreSnapshotRequestMarshalsAliasRename(t *testing.T) {
	req := RestoreSnapshotRequest{
		Indices:                "a,b",
		IncludeAliases:         true,
		RenameAliasPattern:     "logz-(.+)-write-alias",
		RenameAliasReplacement: "old-$1-alias",
	}
	b, err := json.Marshal(req)
	if err != nil {
		t.Fatalf("marshal: %v", err)
	}
	s := string(b)
	for _, want := range []string{
		`"rename_alias_pattern":"logz-(.+)-write-alias"`,
		`"rename_alias_replacement":"old-$1-alias"`,
		`"include_aliases":true`,
	} {
		if !strings.Contains(s, want) {
			t.Errorf("marshaled request missing %s in: %s", want, s)
		}
	}
}

func TestRestoreSnapshotWithTimeout(t *testing.T) {
	srv, cli := testutil.NewMockServer(`{"accepted":true}`, "/_snapshot/my-repo/snap-1/_restore")
	defer srv.Close()
	shared.SetClient(cli)

	err := RestoreSnapshotWithTimeout(context.Background(), "my-repo", "snap-1",
		RestoreSnapshotRequest{Indices: "a"}, true, "5m")
	if err != nil {
		t.Fatalf("RestoreSnapshotWithTimeout error = %v", err)
	}
}
