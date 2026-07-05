package main

import (
	"context"
	"errors"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/pincher95/esctl/cmd"
)

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
		os.Exit(1)
	}
}
