package logger

import (
	"io"
	"log/slog"
	"os"
)

var Logger *slog.Logger

// Init initializes the global logger with the specified debug level
func Init(debug bool) {
	level := slog.LevelInfo
	if debug {
		level = slog.LevelDebug
	}

	handler := slog.NewTextHandler(os.Stderr, &slog.HandlerOptions{
		Level: level,
		// Add timestamp and source location for debug mode
		AddSource: debug,
	})

	Logger = slog.New(handler)
}

// InitWithWriter initializes the logger with a custom writer (useful for testing)
func InitWithWriter(writer io.Writer, debug bool) {
	level := slog.LevelInfo
	if debug {
		level = slog.LevelDebug
	}

	handler := slog.NewTextHandler(writer, &slog.HandlerOptions{
		Level:     level,
		AddSource: debug,
	})

	Logger = slog.New(handler)
}

// Debug logs a debug message
func Debug(msg string, args ...any) {
	if Logger != nil {
		Logger.Debug(msg, args...)
	}
}

// Info logs an info message
func Info(msg string, args ...any) {
	if Logger != nil {
		Logger.Info(msg, args...)
	}
}

// Warn logs a warning message
func Warn(msg string, args ...any) {
	if Logger != nil {
		Logger.Warn(msg, args...)
	}
}

// Error logs an error message
func Error(msg string, args ...any) {
	if Logger != nil {
		Logger.Error(msg, args...)
	}
}

// With returns a logger with the given attributes
func With(args ...any) *slog.Logger {
	if Logger != nil {
		return Logger.With(args...)
	}
	return slog.Default()
}
