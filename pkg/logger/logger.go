package logger

// use this packge for logging debug information using the slog.Info or other functions
// Note: this is not intended for user experience

import (
	"log/slog"
	"os"
)

func init() {
	// Enable AddSource: true to automatically include file and line
	opts := &slog.HandlerOptions{
		Level:     slog.LevelInfo,
		AddSource: true, // <--- This enables file/line reporting
	}
	// Create a handler that writes to Stdout
	// can switch to slog.NewJSONHandler for machine-readable logs
	handler := slog.NewTextHandler(os.Stdout, opts)

	// Create the logger
	logger := slog.New(handler)

	// Set this as the global default so I can use slog.Info(...) anywhere
	slog.SetDefault(logger)
}
