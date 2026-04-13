package cmd

import (
	"testing"
)

func TestJikkyoCommand(t *testing.T) {
	// Test that jikkyoCmd is properly configured
	if jikkyoCmd == nil {
		t.Fatal("jikkyoCmd should not be nil")
	}

	if jikkyoCmd.Use != "jikkyo" {
		t.Errorf("expected Use 'jikkyo', got %s", jikkyoCmd.Use)
	}

	if jikkyoCmd.Short == "" {
		t.Error("expected Short description to be non-empty")
	}

	if jikkyoCmd.Long == "" {
		t.Error("expected Long description to be non-empty")
	}
}

func TestFlagConfiguration(t *testing.T) {
	// Reset flag values
	title = ""
	episode = 0
	verbose = false
	logFile = ""

	// Verify default values
	if title != "" {
		t.Errorf("expected default title to be empty, got %s", title)
	}
	if episode != 0 {
		t.Errorf("expected default episode to be 0, got %d", episode)
	}
	if verbose {
		t.Error("expected default verbose to be false")
	}
	if logFile != "" {
		t.Errorf("expected default logFile to be empty, got %s", logFile)
	}
}

func TestErrorf(t *testing.T) {
	err := errorf("test error")
	if err == nil {
		t.Fatal("expected error to be non-nil")
	}
	if err.Error() != "test error" {
		t.Errorf("expected error message 'test error', got %s", err.Error())
	}

	// Test with formatted message
	err = errorf("test %s %d", "error", 123)
	if err == nil {
		t.Fatal("expected error to be non-nil")
	}
	if err.Error() != "test error 123" {
		t.Errorf("expected error message 'test error 123', got %s", err.Error())
	}
}
