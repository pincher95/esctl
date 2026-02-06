package progress

import (
	"bytes"
	"strings"
	"testing"
	"time"
)

func TestNewBar(t *testing.T) {
	bar := New(100)
	if bar == nil {
		t.Error("New() should return a non-nil bar")
	}

	if bar.total != 100 {
		t.Errorf("Expected total 100, got %d", bar.total)
	}
}

func TestBarAdd(t *testing.T) {
	buf := &bytes.Buffer{}
	bar := NewWithOptions(100, Options{
		Prefix: "Test",
		Writer: buf,
	})

	bar.Add(50)

	if bar.current != 50 {
		t.Errorf("Expected current 50, got %d", bar.current)
	}

	output := buf.String()
	if !strings.Contains(output, "50.0%") {
		t.Errorf("Expected output to contain 50.0%%, got: %s", output)
	}
}

func TestBarSet(t *testing.T) {
	buf := &bytes.Buffer{}
	bar := NewWithOptions(100, Options{
		Prefix: "Test",
		Writer: buf,
	})

	bar.Set(75)

	if bar.current != 75 {
		t.Errorf("Expected current 75, got %d", bar.current)
	}
}

func TestBarFinish(t *testing.T) {
	buf := &bytes.Buffer{}
	bar := NewWithOptions(100, Options{
		Prefix: "Test",
		Writer: buf,
	})

	bar.Add(50)
	bar.Finish()

	if bar.current != 100 {
		t.Errorf("Expected current 100 after finish, got %d", bar.current)
	}

	if !bar.done {
		t.Error("Expected done to be true after finish")
	}

	output := buf.String()
	if !strings.Contains(output, "100.0%") {
		t.Errorf("Expected output to contain 100.0%%, got: %s", output)
	}
}

func TestBarOverflow(t *testing.T) {
	buf := &bytes.Buffer{}
	bar := NewWithOptions(100, Options{
		Prefix: "Test",
		Writer: buf,
	})

	bar.Add(150) // More than total

	if bar.current > 100 {
		t.Errorf("Current should not exceed total, got %d", bar.current)
	}
}

func TestBarMultipleAdd(t *testing.T) {
	buf := &bytes.Buffer{}
	bar := NewWithOptions(100, Options{
		Prefix: "Test",
		Writer: buf,
	})

	bar.Add(25)
	bar.Add(25)
	bar.Add(25)

	if bar.current != 75 {
		t.Errorf("Expected current 75, got %d", bar.current)
	}
}

func TestNewSpinner(t *testing.T) {
	spinner := NewSpinner("Loading")
	if spinner == nil {
		t.Error("NewSpinner() should return a non-nil spinner")
	}

	if spinner.message != "Loading" {
		t.Errorf("Expected message 'Loading', got %s", spinner.message)
	}
}

func TestSpinnerStartStop(t *testing.T) {
	buf := &bytes.Buffer{}
	spinner := NewSpinnerWithWriter("Processing", buf)

	spinner.Start()
	time.Sleep(300 * time.Millisecond) // Let it spin a few times
	spinner.Stop()

	output := buf.String()
	if !strings.Contains(output, "Processing") {
		t.Errorf("Expected output to contain 'Processing', got: %s", output)
	}

	if !strings.Contains(output, "✓") {
		t.Errorf("Expected output to contain checkmark, got: %s", output)
	}
}

func TestSpinnerUpdateMessage(t *testing.T) {
	buf := &bytes.Buffer{}
	spinner := NewSpinnerWithWriter("Initial", buf)

	spinner.Start()
	time.Sleep(100 * time.Millisecond)
	spinner.UpdateMessage("Updated")
	time.Sleep(100 * time.Millisecond)
	spinner.Stop()

	output := buf.String()
	if !strings.Contains(output, "Updated") {
		t.Errorf("Expected output to contain 'Updated', got: %s", output)
	}
}

func TestNewCounter(t *testing.T) {
	counter := NewCounter("Items")
	if counter == nil {
		t.Error("NewCounter() should return a non-nil counter")
	}

	if counter.label != "Items" {
		t.Errorf("Expected label 'Items', got %s", counter.label)
	}
}

func TestCounterInc(t *testing.T) {
	buf := &bytes.Buffer{}
	counter := NewCounter("Items")
	counter.writer = buf

	counter.Inc()
	counter.Inc()

	if counter.Value() != 2 {
		t.Errorf("Expected value 2, got %d", counter.Value())
	}
}

func TestCounterAdd(t *testing.T) {
	buf := &bytes.Buffer{}
	counter := NewCounter("Items")
	counter.writer = buf

	counter.Add(5)
	counter.Add(10)

	if counter.Value() != 15 {
		t.Errorf("Expected value 15, got %d", counter.Value())
	}
}

func TestCounterFinish(t *testing.T) {
	buf := &bytes.Buffer{}
	counter := NewCounter("Items")
	counter.writer = buf

	counter.Add(42)
	counter.Finish()

	output := buf.String()
	if !strings.Contains(output, "42") {
		t.Errorf("Expected output to contain '42', got: %s", output)
	}
}

func TestBarWithZeroTotal(t *testing.T) {
	buf := &bytes.Buffer{}
	bar := NewWithOptions(0, Options{
		Prefix: "Test",
		Writer: buf,
	})

	bar.Add(10) // Should not panic

	// With zero total, render should be a no-op
	output := buf.String()
	if output != "" {
		t.Errorf("Expected empty output with zero total, got: %s", output)
	}
}

func TestBarIdempotentFinish(t *testing.T) {
	buf := &bytes.Buffer{}
	bar := NewWithOptions(100, Options{
		Prefix: "Test",
		Writer: buf,
	})

	bar.Finish()
	bar.Finish() // Should not panic or change state

	if bar.current != 100 {
		t.Errorf("Expected current to remain 100, got %d", bar.current)
	}
}

func TestSpinnerIdempotentStop(t *testing.T) {
	buf := &bytes.Buffer{}
	spinner := NewSpinnerWithWriter("Test", buf)

	spinner.Start()
	time.Sleep(100 * time.Millisecond)
	spinner.Stop()
	spinner.Stop() // Should not panic

	if !spinner.done {
		t.Error("Expected spinner to be done")
	}
}
