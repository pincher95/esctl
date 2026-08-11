package update

import (
	"reflect"
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
