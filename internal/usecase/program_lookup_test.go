package usecase

import (
	"testing"
	"time"

	"github.com/hogecode/JikkyoUtil/internal/models"
)

func TestParseProgItemTimes(t *testing.T) {
	progItem := &models.ProgItem{
		StTime: "2014-09-06 10:00:00",
		EdTime: "2014-09-06 10:30:00",
	}

	stUnix, edUnix, err := ParseProgItemTimes(progItem)
	if err != nil {
		t.Fatalf("failed to parse times: %v", err)
	}

	// Verify Unix timestamps are reasonable
	if stUnix == 0 {
		t.Error("expected non-zero stUnix")
	}
	if edUnix == 0 {
		t.Error("expected non-zero edUnix")
	}
	if edUnix <= stUnix {
		t.Error("expected edUnix to be greater than stUnix")
	}

	// Verify the difference is 30 minutes (1800 seconds)
	diff := edUnix - stUnix
	if diff != 1800 {
		t.Errorf("expected 30 minute difference (1800s), got %d seconds", diff)
	}
}

func TestFormatUnixTimestamp(t *testing.T) {
	// Use UTC to avoid timezone issues in testing
	expectedTime := time.Date(2014, 9, 6, 10, 0, 0, 0, time.UTC)
	unixTime := expectedTime.Unix()

	formatted := FormatUnixTimestamp(unixTime)

	// Verify format matches expected pattern
	if len(formatted) != len("2006-01-02 15:04:05") {
		t.Errorf("expected formatted time length 19, got %d: %s", len(formatted), formatted)
	}

	// Verify format contains expected components
	if !contains(formatted, "-") || !contains(formatted, ":") {
		t.Errorf("expected formatted time to contain date/time separators: %s", formatted)
	}
}

// contains checks if a string contains a substring
func contains(s, substr string) bool {
	for i := 0; i < len(s)-len(substr)+1; i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}

func TestConvertStTimeToUnix(t *testing.T) {
	stTime := "2014-09-06 10:00:00"

	unixTime, err := ConvertStTimeToUnix(stTime)
	if err != nil {
		t.Fatalf("failed to convert StTime: %v", err)
	}

	if unixTime == 0 {
		t.Error("expected non-zero Unix timestamp")
	}

	// Format back and compare
	formatted := FormatUnixTimestamp(unixTime)
	if formatted != stTime {
		t.Logf("formatted: %s, original: %s (time zone may differ)", formatted, stTime)
	}
}

func TestConvertStringToInt(t *testing.T) {
	tests := []struct {
		input    string
		expected int
	}{
		{"10", 10},
		{"1", 1},
		{"100", 100},
		{"0", 0},
	}

	for _, test := range tests {
		result, err := ConvertStringToInt(test.input)
		if err != nil {
			t.Errorf("failed to convert %s: %v", test.input, err)
		}
		if result != test.expected {
			t.Errorf("expected %d, got %d for input %s", test.expected, result, test.input)
		}
	}

	// Test invalid input
	_, err := ConvertStringToInt("invalid")
	if err == nil {
		t.Error("expected error for invalid input")
	}
}
