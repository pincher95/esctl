package update

import (
	"reflect"
	"testing"
	"time"
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
