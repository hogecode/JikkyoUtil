package presentation

import (
	"io"
	"log/slog"
	"os"
)

// LoggerConfig holds logger configuration
type LoggerConfig struct {
	Verbose bool
	Output  io.Writer
	LogFile string // Optional log file path
}

// NewLogger creates a new logger with the specified configuration
func NewLogger(config LoggerConfig) (*slog.Logger, error) {
	output := config.Output
	if output == nil {
		output = os.Stderr
	}

	var level slog.Level
	if config.Verbose {
		level = slog.LevelDebug
	} else {
		level = slog.LevelInfo
	}

	// If log file is specified, also write to file
	if config.LogFile != "" {
		file, err := os.OpenFile(config.LogFile, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
		if err != nil {
			return nil, err
		}
		// Use MultiWriter to write to both stderr and file
		output = io.MultiWriter(output, file)
	}

	opts := &slog.HandlerOptions{
		Level: level,
	}

	handler := slog.NewTextHandler(output, opts)
	return slog.New(handler), nil
}
