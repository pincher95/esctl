package utils

import (
	"bufio"
	"bytes"
	"os"
	"os/signal"
	"syscall"
	"time"
)

// WatchLoop runs draw() every interval inside an alternate screen buffer.
// It restores the original screen and cursor even on Ctrl-C.
func WatchLoop(interval time.Duration, draw func() error) error {
	// switch to alternate screen & hide cursor
	// os.Stdout.WriteString("\x1b[?1049h\x1b[?25l")
	os.Stdout.WriteString("\x1b[H\x1b[J")
	defer os.Stdout.WriteString("\x1b[?1049l\x1b[?25h")

	sig := make(chan os.Signal, 1)
	signal.Notify(sig, syscall.SIGINT, syscall.SIGTERM)

	// channel to signal termination
	done := make(chan struct{})

	go func() {
		<-sig
		close(done)
	}()

	// Listen for 'q' + Enter to quit (no raw mode)
	go func() {
		reader := bufio.NewReader(os.Stdin)
		for {
			line, err := reader.ReadBytes('\n')
			if err != nil {
				return
			}
			if len(bytes.TrimSpace(line)) == 1 && bytes.TrimSpace(line)[0] == 'q' {
				close(done)
				return
			}
		}
	}()

	for {
		os.Stdout.WriteString("\x1b[H") // move cursor home
		if err := draw(); err != nil {
			return err
		}
		// clear rest of screen to remove leftover lines if output shrank
		os.Stdout.WriteString("\x1b[J")
		select {
		case <-time.After(interval):
		case <-done:
			return nil
		}
	}
}
