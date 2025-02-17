package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/pincher95/esctl/cmd"
)

func main() {
	// Create a cancellable context that listens for SIGINT or SIGTERM
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		<-sigChan
		fmt.Println("Received an interrupt, cancelling context...")
		cancel()
	}()

	// Execute our root command with the context
	if err := cmd.Execute(ctx); err != nil {
		fmt.Println("Error:", err)
		os.Exit(1)
	}
}
