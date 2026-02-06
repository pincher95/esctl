package logger

import (
	"bytes"
	"strings"
	"testing"
)

func TestInit(t *testing.T) {
	// Test normal init
	Init(false)
	if Logger == nil {
		t.Error("Logger should be initialized")
	}

	// Test debug init
	Init(true)
	if Logger == nil {
		t.Error("Logger should be initialized in debug mode")
	}
}

func TestInitWithWriter(t *testing.T) {
	buf := &bytes.Buffer{}
	InitWithWriter(buf, false)

	if Logger == nil {
		t.Error("Logger should be initialized")
	}

	// Test that logging works
	Info("test message")
	output := buf.String()

	if !strings.Contains(output, "test message") {
		t.Errorf("Expected 'test message' in output, got: %s", output)
	}
}

func TestDebugLogging(t *testing.T) {
	buf := &bytes.Buffer{}

	// Test with debug disabled
	InitWithWriter(buf, false)
	Debug("debug message")

	if strings.Contains(buf.String(), "debug message") {
		t.Error("Debug message should not appear when debug is disabled")
	}

	// Test with debug enabled
	buf.Reset()
	InitWithWriter(buf, true)
	Debug("debug message")

	if !strings.Contains(buf.String(), "debug message") {
		t.Error("Debug message should appear when debug is enabled")
	}
}

func TestInfoLogging(t *testing.T) {
	buf := &bytes.Buffer{}
	InitWithWriter(buf, false)

	Info("info message")
	output := buf.String()

	if !strings.Contains(output, "info message") {
		t.Errorf("Expected 'info message' in output, got: %s", output)
	}
}

func TestWarnLogging(t *testing.T) {
	buf := &bytes.Buffer{}
	InitWithWriter(buf, false)

	Warn("warning message")
	output := buf.String()

	if !strings.Contains(output, "warning message") {
		t.Errorf("Expected 'warning message' in output, got: %s", output)
	}
}

func TestErrorLogging(t *testing.T) {
	buf := &bytes.Buffer{}
	InitWithWriter(buf, false)

	Error("error message")
	output := buf.String()

	if !strings.Contains(output, "error message") {
		t.Errorf("Expected 'error message' in output, got: %s", output)
	}
}

func TestWith(t *testing.T) {
	buf := &bytes.Buffer{}
	InitWithWriter(buf, false)

	contextLogger := With("component", "test", "user", "alice")
	contextLogger.Info("test message")

	output := buf.String()

	if !strings.Contains(output, "component=test") {
		t.Errorf("Expected 'component=test' in output, got: %s", output)
	}

	if !strings.Contains(output, "user=alice") {
		t.Errorf("Expected 'user=alice' in output, got: %s", output)
	}
}

func TestLoggingWithAttributes(t *testing.T) {
	buf := &bytes.Buffer{}
	InitWithWriter(buf, false)

	Info("operation complete", "duration", "2s", "records", 100)
	output := buf.String()

	if !strings.Contains(output, "operation complete") {
		t.Errorf("Expected 'operation complete' in output, got: %s", output)
	}

	if !strings.Contains(output, "duration=2s") {
		t.Errorf("Expected 'duration=2s' in output, got: %s", output)
	}

	if !strings.Contains(output, "records=100") {
		t.Errorf("Expected 'records=100' in output, got: %s", output)
	}
}
