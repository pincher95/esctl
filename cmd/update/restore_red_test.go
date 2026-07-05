package update

import "testing"

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
