package testutil

import (
	"log/slog"
	"os"
	"testing"
)

// SetupTestLogging configures a silent logger for tests to reduce verbose output
func SetupTestLogging(t *testing.T) *slog.Logger {
	// Create a logger that writes to a null handler (discards all output)
	logger := slog.New(slog.NewTextHandler(os.Stderr, &slog.HandlerOptions{
		Level: slog.LevelError, // Only show errors during tests
	}))

	// Set as default logger
	slog.SetDefault(logger)

	return logger
}
