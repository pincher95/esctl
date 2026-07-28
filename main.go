package main

import (
	"context"
	"errors"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"time"

	"github.com/pincher95/esctl/cmd"
	"github.com/pincher95/esctl/internal/client"
	eserrors "github.com/pincher95/esctl/internal/errors"
	"github.com/pincher95/esctl/shared"
)

// timeoutHint returns the timeout that was in effect for this run and a longer one
// to suggest, so the hint stays accurate when --timeout was already raised.
func timeoutHint() (effective, suggested time.Duration) {
	effective = shared.TimeoutDuration
	if effective <= 0 {
		effective = client.DefaultTimeout
	}
	return effective, max(2*effective, 5*time.Minute)
}

func main() {
	// Cancel the root context on SIGINT/SIGTERM. Commands (and the watch loop)
	// observe the cancellation via ctx, run their deferred cleanup — restoring the
	// terminal in watch mode — and return. This is the single source of signal
	// handling so interrupt behavior is deterministic.
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	if err := cmd.Execute(ctx); err != nil {
		// A user-initiated interrupt cancels the context; exit quietly with the
		// conventional 128+SIGINT status rather than printing a scary error.
		if errors.Is(err, context.Canceled) {
			os.Exit(130)
		}
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		if eserrors.IsTimeout(err) {
			effective, suggested := timeoutHint()
			fmt.Fprintf(os.Stderr,
				"hint: the request exceeded the %s client timeout. Some endpoints (listing snapshots in a\n"+
					"      large repository, stats on a big cluster) need longer: retry with --timeout %s.\n",
				effective, suggested)
		}
		os.Exit(1)
	}
}
