package presentation

import (
	"testing"
)

func TestOutputFormatter(t *testing.T) {
	formatter := NewOutputFormatter(false)

	// Verify the formatter is created successfully
	if formatter == nil {
		t.Error("failed to create output formatter")
	}
	if formatter.verbose {
		t.Error("expected verbose=false for non-verbose formatter")
	}
}

func TestOutputFormatterWithEmpty(t *testing.T) {
	formatter := NewOutputFormatter(true)
	if formatter == nil {
		t.Error("failed to create verbose formatter")
	}
	if !formatter.verbose {
		t.Error("expected formatter to be verbose")
	}
}

func TestNewOutputFormatter(t *testing.T) {
	verbose := NewOutputFormatter(true)
	if !verbose.verbose {
		t.Error("expected verbose formatter to have verbose=true")
	}

	normal := NewOutputFormatter(false)
	if normal.verbose {
		t.Error("expected normal formatter to have verbose=false")
	}
}
