package utils

import (
	"context"
	"errors"
	"io"
	"os"
	"sync/atomic"
	"testing"
	"time"
)

// withPipes points os.Stdin/os.Stdout at pipes (not character devices) so
// WatchLoopContext takes its non-interactive path and never touches the real
// terminal during tests.
func withPipes(t *testing.T, fn func()) {
	t.Helper()
	inR, inW, err := os.Pipe()
	if err != nil {
		t.Fatalf("pipe: %v", err)
	}
	outR, outW, err := os.Pipe()
	if err != nil {
		t.Fatalf("pipe: %v", err)
	}
	go io.Copy(io.Discard, outR)

	oldIn, oldOut := os.Stdin, os.Stdout
	os.Stdin, os.Stdout = inR, outW
	defer func() {
		os.Stdin, os.Stdout = oldIn, oldOut
		inR.Close()
		inW.Close()
		outR.Close()
		outW.Close()
	}()
	fn()
}

func TestWatchLoopContextStopsOnCancel(t *testing.T) {
	withPipes(t, func() {
		ctx, cancel := context.WithCancel(context.Background())
		var count int32
		go func() {
			for atomic.LoadInt32(&count) < 2 {
				time.Sleep(time.Millisecond)
			}
			cancel()
		}()

		err := WatchLoopContext(ctx, 2*time.Millisecond, func() error {
			atomic.AddInt32(&count, 1)
			return nil
		})
		if !errors.Is(err, context.Canceled) {
			t.Fatalf("expected context.Canceled, got %v", err)
		}
		if atomic.LoadInt32(&count) < 2 {
			t.Fatalf("expected draw to run at least twice, ran %d", count)
		}
	})
}

func TestWatchLoopContextReturnsDrawError(t *testing.T) {
	withPipes(t, func() {
		boom := errors.New("boom")
		err := WatchLoopContext(context.Background(), time.Hour, func() error {
			return boom
		})
		if !errors.Is(err, boom) {
			t.Fatalf("expected draw error to propagate, got %v", err)
		}
	})
}
