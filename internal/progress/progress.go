package progress

import (
	"fmt"
	"io"
	"os"
	"strings"
	"sync"
	"time"
)

// Bar represents a progress bar for long-running operations
type Bar struct {
	total   int64
	current int64
	start   time.Time
	prefix  string
	writer  io.Writer
	mu      sync.Mutex
	done    bool
}

// New creates a new progress bar
func New(total int64) *Bar {
	return NewWithOptions(total, Options{
		Prefix: "Progress",
		Writer: os.Stdout,
	})
}

// Options for creating a progress bar
type Options struct {
	Prefix string
	Writer io.Writer
}

// NewWithOptions creates a new progress bar with custom options
func NewWithOptions(total int64, opts Options) *Bar {
	prefix := opts.Prefix
	if prefix == "" {
		prefix = "Progress"
	}

	writer := opts.Writer
	if writer == nil {
		writer = os.Stdout
	}

	return &Bar{
		total:  total,
		start:  time.Now(),
		prefix: prefix,
		writer: writer,
	}
}

// Add increments the progress by n
func (b *Bar) Add(n int64) {
	b.mu.Lock()
	defer b.mu.Unlock()

	if b.done {
		return
	}

	b.current += n
	if b.current > b.total {
		b.current = b.total
	}

	b.render()
}

// Set sets the current progress to n
func (b *Bar) Set(n int64) {
	b.mu.Lock()
	defer b.mu.Unlock()

	if b.done {
		return
	}

	b.current = min(n, b.total)

	b.render()
}

// Finish marks the progress as complete
func (b *Bar) Finish() {
	b.mu.Lock()
	defer b.mu.Unlock()

	if b.done {
		return
	}

	b.current = b.total
	b.render()
	fmt.Fprintln(b.writer) // New line after completion
	b.done = true
}

// render draws the progress bar (must be called with lock held)
func (b *Bar) render() {
	if b.total == 0 {
		return
	}

	percent := float64(b.current) / float64(b.total) * 100
	elapsed := time.Since(b.start)

	// Calculate ETA
	var eta string
	if b.current > 0 {
		rate := float64(b.current) / elapsed.Seconds()
		if rate > 0 {
			remaining := float64(b.total-b.current) / rate
			eta = fmt.Sprintf(" - ETA: %s", time.Duration(remaining*float64(time.Second)).Round(time.Second))
		}
	}

	// Create progress bar
	barWidth := 40
	filled := min(int(float64(barWidth)*percent/100), barWidth)

	bar := strings.Repeat("█", filled) + strings.Repeat("░", barWidth-filled)

	fmt.Fprintf(b.writer, "\r%s: [%s] %.1f%% (%d/%d) - Elapsed: %s%s",
		b.prefix, bar, percent, b.current, b.total, elapsed.Round(time.Second), eta)
}

// Spinner represents a simple spinner for operations with unknown duration
type Spinner struct {
	message string
	frames  []string
	current int
	stop    chan bool
	done    bool
	writer  io.Writer
	mu      sync.Mutex
}

// NewSpinner creates a new spinner
func NewSpinner(message string) *Spinner {
	return &Spinner{
		message: message,
		frames:  []string{"⠋", "⠙", "⠹", "⠸", "⠼", "⠴", "⠦", "⠧", "⠇", "⠏"},
		stop:    make(chan bool),
		writer:  os.Stdout,
	}
}

// NewSpinnerWithWriter creates a new spinner with a custom writer
func NewSpinnerWithWriter(message string, writer io.Writer) *Spinner {
	return &Spinner{
		message: message,
		frames:  []string{"⠋", "⠙", "⠹", "⠸", "⠼", "⠴", "⠦", "⠧", "⠇", "⠏"},
		stop:    make(chan bool),
		writer:  writer,
	}
}

// Start starts the spinner animation
func (s *Spinner) Start() {
	go func() {
		ticker := time.NewTicker(100 * time.Millisecond)
		defer ticker.Stop()

		for {
			select {
			case <-s.stop:
				return
			case <-ticker.C:
				s.mu.Lock()
				if !s.done {
					frame := s.frames[s.current%len(s.frames)]
					fmt.Fprintf(s.writer, "\r%s %s", frame, s.message)
					s.current++
				}
				s.mu.Unlock()
			}
		}
	}()
}

// Stop stops the spinner animation
func (s *Spinner) Stop() {
	s.mu.Lock()
	if s.done {
		s.mu.Unlock()
		return
	}

	s.done = true
	s.mu.Unlock()

	// Send stop signal without holding the lock
	select {
	case s.stop <- true:
	default:
		// Channel might be closed or full, that's OK
	}

	fmt.Fprintf(s.writer, "\r✓ %s\n", s.message)
}

// UpdateMessage updates the spinner message
func (s *Spinner) UpdateMessage(message string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.message = message
}

// Counter represents a simple counter for tracking operations
type Counter struct {
	value  int64
	label  string
	writer io.Writer
	mu     sync.Mutex
}

// NewCounter creates a new counter
func NewCounter(label string) *Counter {
	return &Counter{
		label:  label,
		writer: os.Stdout,
	}
}

// Inc increments the counter
func (c *Counter) Inc() {
	c.Add(1)
}

// Add adds n to the counter
func (c *Counter) Add(n int64) {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.value += n
	fmt.Fprintf(c.writer, "\r%s: %d", c.label, c.value)
}

// Value returns the current value
func (c *Counter) Value() int64 {
	c.mu.Lock()
	defer c.mu.Unlock()
	return c.value
}

// Finish completes the counter display
func (c *Counter) Finish() {
	c.mu.Lock()
	defer c.mu.Unlock()
	fmt.Fprintf(c.writer, "\r%s: %d\n", c.label, c.value)
}
