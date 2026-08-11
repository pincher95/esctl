package update

import (
	"encoding/json"
	"reflect"
	"strings"
	"testing"
	"time"

	"github.com/pincher95/esctl/es/cat"
)

func TestBatchStrings(t *testing.T) {
	got := batchStrings([]string{"a", "b", "c", "d", "e"}, 2)
	if len(got) != 3 {
		t.Fatalf("expected 3 batches, got %d (%v)", len(got), got)
	}
	if len(got[0]) != 2 || len(got[1]) != 2 || len(got[2]) != 1 {
		t.Errorf("unexpected batch sizes: %v", got)
	}

	// size <= 0 falls back to the default (one batch here)
	if b := batchStrings([]string{"a", "b"}, 0); len(b) != 1 {
		t.Errorf("expected 1 batch for size 0, got %d", len(b))
	}

	// empty input yields no batches
	if b := batchStrings(nil, 10); len(b) != 0 {
		t.Errorf("expected 0 batches for empty input, got %d", len(b))
	}
}

func TestDateExclusions(t *testing.T) {
	now := time.Date(2026, 8, 11, 15, 0, 0, 0, time.UTC)
	if got, want := dateExclusions(now), []string{"260811", "260812"}; !reflect.DeepEqual(got, want) {
		t.Errorf("dateExclusions(%v) = %v, want %v", now, got, want)
	}

	// Tomorrow is the next calendar day, across month and year boundaries.
	eoy := time.Date(2026, 12, 31, 23, 30, 0, 0, time.UTC)
	if got, want := dateExclusions(eoy), []string{"261231", "270101"}; !reflect.DeepEqual(got, want) {
		t.Errorf("dateExclusions(%v) = %v, want %v", eoy, got, want)
	}
}

func TestContainsAny(t *testing.T) {
	if !containsAny("logz-abc-260811-000001", []string{"260811", "260812"}) {
		t.Error("expected match on 260811")
	}
	if containsAny("logz-abc-260810-000001", []string{"260811", "260812"}) {
		t.Error("did not expect a match for 260810")
	}
	if containsAny("logz-abc-260811", nil) {
		t.Error("nil exclusions must never match")
	}
}

func catIdx(name, health, status string) cat.CatIndiceResponse {
	return cat.CatIndiceResponse{Index: name, Health: health, Status: status}
}

func TestSelectIndicesPatternMode(t *testing.T) {
	indices := []cat.CatIndiceResponse{
		catIdx("logz-red", "red", "open"),
		catIdx("logz-green", "green", "open"),
		catIdx("logz-closed", "red", "close"),
		catIdx("logz-red-nosnap", "red", "open"),
	}
	inSnap := map[string]bool{"logz-red": true, "logz-closed": true}

	sel := selectIndices(indices, inSnap, false, false, nil)
	if want := []string{"logz-red"}; !reflect.DeepEqual(sel.toRestore, want) {
		t.Errorf("toRestore = %v, want %v", sel.toRestore, want)
	}
	if want := []string{"logz-closed"}; !reflect.DeepEqual(sel.closedSkipped, want) {
		t.Errorf("closedSkipped = %v, want %v", sel.closedSkipped, want)
	}
	if want := []string{"logz-red-nosnap"}; !reflect.DeepEqual(sel.notInSnapshot, want) {
		t.Errorf("notInSnapshot = %v, want %v", sel.notInSnapshot, want)
	}

	// includeClosed promotes the closed index to a restore target.
	sel = selectIndices(indices, inSnap, false, true, nil)
	if want := []string{"logz-red", "logz-closed"}; !reflect.DeepEqual(sel.toRestore, want) {
		t.Errorf("toRestore with includeClosed = %v, want %v", sel.toRestore, want)
	}
	if len(sel.closedSkipped) != 0 {
		t.Errorf("closedSkipped with includeClosed = %v, want empty", sel.closedSkipped)
	}
}

func TestSelectIndicesAliasMode(t *testing.T) {
	indices := []cat.CatIndiceResponse{
		catIdx("logz-green", "green", "open"),
		catIdx("logz-red", "red", "open"),
		catIdx("logz-closed", "green", "close"),
		catIdx("logz-nosnap", "green", "open"),
	}
	inSnap := map[string]bool{"logz-green": true, "logz-red": true, "logz-closed": true}

	// includeClosed=false on purpose: alias mode must take closed indices anyway.
	sel := selectIndices(indices, inSnap, true, false, nil)
	if want := []string{"logz-green", "logz-red", "logz-closed"}; !reflect.DeepEqual(sel.toRestore, want) {
		t.Errorf("toRestore = %v, want %v", sel.toRestore, want)
	}
	if want := []string{"logz-nosnap"}; !reflect.DeepEqual(sel.notInSnapshot, want) {
		t.Errorf("notInSnapshot = %v, want %v", sel.notInSnapshot, want)
	}
	if len(sel.closedSkipped) != 0 || len(sel.dateExcluded) != 0 {
		t.Errorf("unexpected closedSkipped=%v dateExcluded=%v", sel.closedSkipped, sel.dateExcluded)
	}
	// openToRestore counts only the open indices among toRestore: logz-green and
	// logz-red are open, logz-closed is closed, logz-nosnap never lands in
	// toRestore because it is absent from the snapshot.
	if sel.openToRestore != 2 {
		t.Errorf("openToRestore = %d, want 2", sel.openToRestore)
	}
}

func TestSelectIndicesOpenToRestoreZeroInPatternMode(t *testing.T) {
	indices := []cat.CatIndiceResponse{
		catIdx("logz-red", "red", "open"),
		catIdx("logz-closed", "red", "close"),
	}
	inSnap := map[string]bool{"logz-red": true, "logz-closed": true}

	sel := selectIndices(indices, inSnap, false, true, nil)
	if want := []string{"logz-red", "logz-closed"}; !reflect.DeepEqual(sel.toRestore, want) {
		t.Fatalf("toRestore = %v, want %v", sel.toRestore, want)
	}
	if sel.openToRestore != 0 {
		t.Errorf("openToRestore = %d, want 0 in pattern mode", sel.openToRestore)
	}
}

// TestSelectIndicesPatternModeClosedNotInSnapshotDropped pins existing behavior:
// a closed index that is neither included via --include-closed nor present in
// the snapshot is silently dropped from every bucket, since closedSkipped is
// only populated for closed indices that ARE in the snapshot.
func TestSelectIndicesPatternModeClosedNotInSnapshotDropped(t *testing.T) {
	indices := []cat.CatIndiceResponse{
		catIdx("logz-closed-nosnap", "red", "close"),
	}
	sel := selectIndices(indices, map[string]bool{}, false, false, nil)
	if len(sel.toRestore) != 0 {
		t.Errorf("toRestore = %v, want empty", sel.toRestore)
	}
	if len(sel.closedSkipped) != 0 {
		t.Errorf("closedSkipped = %v, want empty", sel.closedSkipped)
	}
	if len(sel.notInSnapshot) != 0 {
		t.Errorf("notInSnapshot = %v, want empty", sel.notInSnapshot)
	}
}

// TestSelectIndicesPatternModeExclusionOverridesSnapshotMatch confirms date
// exclusion is checked before mode-specific candidacy: a red, in-snapshot
// index whose name carries an excluded date stamp is set aside as
// dateExcluded rather than restored, regardless of selection mode.
func TestSelectIndicesPatternModeExclusionOverridesSnapshotMatch(t *testing.T) {
	indices := []cat.CatIndiceResponse{
		catIdx("logz-a-260811-000001", "red", "open"),
	}
	inSnap := map[string]bool{"logz-a-260811-000001": true}

	sel := selectIndices(indices, inSnap, false, false, []string{"260811"})
	if len(sel.toRestore) != 0 {
		t.Errorf("toRestore = %v, want empty", sel.toRestore)
	}
	if want := []string{"logz-a-260811-000001"}; !reflect.DeepEqual(sel.dateExcluded, want) {
		t.Errorf("dateExcluded = %v, want %v", sel.dateExcluded, want)
	}
}

func TestSelectIndicesDateExclusion(t *testing.T) {
	indices := []cat.CatIndiceResponse{
		catIdx("logz-a-260811-000001", "red", "open"),
		catIdx("logz-a-260812-000001", "red", "open"),
		catIdx("logz-a-260810-000001", "red", "open"),
	}
	inSnap := map[string]bool{
		"logz-a-260811-000001": true,
		"logz-a-260812-000001": true,
		"logz-a-260810-000001": true,
	}
	sel := selectIndices(indices, inSnap, true, false, []string{"260811", "260812"})
	if want := []string{"logz-a-260810-000001"}; !reflect.DeepEqual(sel.toRestore, want) {
		t.Errorf("toRestore = %v, want %v", sel.toRestore, want)
	}
	if want := []string{"logz-a-260811-000001", "logz-a-260812-000001"}; !reflect.DeepEqual(sel.dateExcluded, want) {
		t.Errorf("dateExcluded = %v, want %v", sel.dateExcluded, want)
	}
}

func TestBuildRestoreRequestSettings(t *testing.T) {
	req := buildRestoreRequest("a,b", true, "logz-(.+)-write-alias", "old-$1-alias",
		0, "default,ingestion",
		[]string{"index.routing.allocation.total_shards_per_node", "index.routing.allocation.require._ip"})
	b, err := json.Marshal(req)
	if err != nil {
		t.Fatalf("marshal: %v", err)
	}
	s := string(b)
	for _, want := range []string{
		`"indices":"a,b"`,
		`"include_aliases":true`,
		`"rename_alias_pattern":"logz-(.+)-write-alias"`,
		`"rename_alias_replacement":"old-$1-alias"`,
		`"index.number_of_replicas":0`,
		`"index.routing.allocation.include.box_type":"default,ingestion"`,
		`"ignore_index_settings":["index.routing.allocation.total_shards_per_node","index.routing.allocation.require._ip"]`,
	} {
		if !strings.Contains(s, want) {
			t.Errorf("marshaled request missing %s in: %s", want, s)
		}
	}
}

func TestBuildRestoreRequestNoOverrides(t *testing.T) {
	req := buildRestoreRequest("a", true, "", "", -1, "", nil)
	if req.IndexSettings != nil {
		t.Errorf("IndexSettings = %v, want nil", req.IndexSettings)
	}
	b, err := json.Marshal(req)
	if err != nil {
		t.Fatalf("marshal: %v", err)
	}
	for _, banned := range []string{"index_settings", "ignore_index_settings"} {
		if strings.Contains(string(b), banned) {
			t.Errorf("marshaled request must not contain %q: %s", banned, b)
		}
	}
}

func TestBuildRestoreRequestReplicasOnly(t *testing.T) {
	req := buildRestoreRequest("a", true, "", "", 2, "", nil)
	want := map[string]any{"index.number_of_replicas": 2}
	if !reflect.DeepEqual(req.IndexSettings, want) {
		t.Errorf("IndexSettings = %v, want %v", req.IndexSettings, want)
	}
}

func TestBuildRestoreRequestBoxTypeOnly(t *testing.T) {
	req := buildRestoreRequest("a", true, "", "", -1, "warm", nil)
	want := map[string]any{"index.routing.allocation.include.box_type": "warm"}
	if !reflect.DeepEqual(req.IndexSettings, want) {
		t.Errorf("IndexSettings = %v, want %v", req.IndexSettings, want)
	}
}

func TestPollPatternFor(t *testing.T) {
	// Alias mode polls "*": a restored index whose write-alias was renamed no
	// longer matches the alias pattern and would look complete while red.
	if got := pollPatternFor(true, "logz-*-write-alias"); got != "*" {
		t.Errorf("alias mode poll pattern = %q, want *", got)
	}
	if got := pollPatternFor(false, "logz-*"); got != "logz-*" {
		t.Errorf("pattern mode poll pattern = %q, want logz-*", got)
	}
}

func TestPreviewList(t *testing.T) {
	if got, want := previewList([]string{"a", "b"}, 5), "a, b"; got != want {
		t.Errorf("fewer than max: previewList = %q, want %q", got, want)
	}
	if got, want := previewList([]string{"a", "b", "c", "d", "e"}, 5), "a, b, c, d, e"; got != want {
		t.Errorf("exactly max: previewList = %q, want %q", got, want)
	}
	if got, want := previewList([]string{"a", "b", "c", "d", "e", "f", "g"}, 5), "a, b, c, d, e, and 2 more"; got != want {
		t.Errorf("more than max: previewList = %q, want %q", got, want)
	}
}

func TestValidateSelectionFlags(t *testing.T) {
	if err := validateSelectionFlags("", ""); err == nil {
		t.Error("expected error when neither flag is set")
	}
	if err := validateSelectionFlags("logz-*", "logz-*-write-alias"); err == nil {
		t.Error("expected error when both flags are set")
	}
	if err := validateSelectionFlags("logz-*", ""); err != nil {
		t.Errorf("pattern only: unexpected error %v", err)
	}
	if err := validateSelectionFlags("", "logz-*-write-alias"); err != nil {
		t.Errorf("alias pattern only: unexpected error %v", err)
	}
}
