package utils

import (
	"bufio"
	"context"
	"fmt"
	"os"
	"os/exec"
	"strings"
	"sync"
	"time"
)

// WatchLoop runs draw() every interval inside an alternate screen buffer.
// It restores the original screen and cursor even on Ctrl-C.
func WatchLoop(interval time.Duration, draw func() error) error {
	return WatchLoopContext(context.Background(), interval, draw)
}

// WatchLoopContext runs draw() every interval until the user quits (q/Esc/Ctrl-C),
// a termination signal arrives, or ctx is done. It returns nil on a user-initiated
// quit and ctx.Err() when the parent context is cancelled (e.g. --timeout).
//
// When stdout is a terminal it renders in the alternate screen buffer and, when
// stdin is also a terminal, switches the terminal to cbreak mode so a single
// keypress quits without requiring Enter. When either is not a terminal (piped or
// redirected), it degrades to a plain redraw loop driven by the interval, signals,
// and ctx only.
func WatchLoopContext(ctx context.Context, interval time.Duration, draw func() error) error {
	ttyOut := isTerminal(os.Stdout)
	ttyIn := isTerminal(os.Stdin)

	if ttyOut {
		// Enter alternate screen buffer and hide the cursor; restore on exit.
		os.Stdout.WriteString("\x1b[?1049h\x1b[?25l")
		defer os.Stdout.WriteString("\x1b[?1049l\x1b[?25h")
	}

	// cbreak mode lets us read a single keypress (q/Esc) without Enter, while
	// keeping output post-processing so draw()'s newlines render correctly.
	if ttyIn {
		if restore := enterCbreak(); restore != nil {
			defer restore()
		}
	}

	// quit is closed once (guarded by sync.Once) when the user presses a quit key.
	// Termination signals (Ctrl-C / SIGTERM) are delivered via ctx cancellation by
	// main(), so we don't register a second signal handler here; that keeps exit
	// behavior deterministic and guarantees the deferred terminal restore runs.
	quit := make(chan struct{})
	var once sync.Once
	stop := func() { once.Do(func() { close(quit) }) }

	if ttyIn {
		go func() {
			reader := bufio.NewReader(os.Stdin)
			for {
				b, err := reader.ReadByte()
				if err != nil {
					return
				}
				switch b {
				case 'q', 'Q', 0x1b /* Esc */, 0x03 /* Ctrl-C */ :
					stop()
					return
				}
				select {
				case <-quit:
					return
				case <-ctx.Done():
					return
				default:
				}
			}
		}()
	}

	for {
		if ttyOut {
			os.Stdout.WriteString("\x1b[H") // move cursor home
		}
		if err := draw(); err != nil {
			return err
		}
		if ttyOut {
			os.Stdout.WriteString("\x1b[J") // clear from cursor to end of screen
			fmt.Fprintf(os.Stdout, "\n[refreshing every %s — press q or Ctrl-C to quit]", interval)
		}

		select {
		case <-time.After(interval):
		case <-quit:
			return nil
		case <-ctx.Done():
			return ctx.Err()
		}
	}
}

// isTerminal reports whether f is attached to a character device (a terminal).
func isTerminal(f *os.File) bool {
	fi, err := f.Stat()
	if err != nil {
		return false
	}
	return fi.Mode()&os.ModeCharDevice != 0
}

// enterCbreak switches the controlling terminal to cbreak mode (no canonical
// line buffering, no echo) and returns a function that restores the previous
// settings. It returns nil if the terminal cannot be configured (e.g. stty is
// unavailable), in which case the caller falls back to signal/ctx-only control.
func enterCbreak() func() {
	saved, err := sttyOutput("-g")
	if err != nil {
		return nil
	}
	if err := stty("-echo", "-icanon", "min", "1", "time", "0"); err != nil {
		return nil
	}
	return func() {
		saved = strings.TrimSpace(saved)
		if saved == "" || stty(saved) != nil {
			_ = stty("sane")
		}
	}
}

func stty(args ...string) error {
	cmd := exec.Command("stty", args...)
	cmd.Stdin = os.Stdin
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

func sttyOutput(args ...string) (string, error) {
	cmd := exec.Command("stty", args...)
	cmd.Stdin = os.Stdin
	out, err := cmd.Output()
	return string(out), err
}
