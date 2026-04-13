package usecase

import (
	"testing"
)

func TestConvertToInt64(t *testing.T) {
	tests := []struct {
		input    interface{}
		expected int64
	}{
		{int64(12345), int64(12345)},
		{int(100), int64(100)},
		{float64(999.5), int64(999)},
		{"54321", int64(54321)},
		{"invalid", int64(0)},
		{nil, int64(0)},
	}

	for _, test := range tests {
		result := convertToInt64(test.input)
		if result != test.expected {
			t.Errorf("convertToInt64(%v) = %d, expected %d", test.input, result, test.expected)
		}
	}
}

func TestFindMostCommonTime(t *testing.T) {
	// Test empty map
	result := findMostCommonTime(make(map[int64]int))
	if result != 0 {
		t.Errorf("expected 0 for empty map, got %d", result)
	}

	// Test map with single entry
	markers := map[int64]int{
		int64(1000): 1,
	}
	result = findMostCommonTime(markers)
	if result != 1000 {
		t.Errorf("expected 1000, got %d", result)
	}

	// Test map with multiple entries
	markers = map[int64]int{
		int64(1000): 2,
		int64(2000): 5,
		int64(3000): 3,
	}
	result = findMostCommonTime(markers)
	if result != 2000 {
		t.Errorf("expected 2000 (highest count), got %d", result)
	}

	// Test map with tied counts (should return one of them)
	markers = map[int64]int{
		int64(1000): 5,
		int64(2000): 5,
	}
	result = findMostCommonTime(markers)
	if result != 1000 && result != 2000 {
		t.Errorf("expected 1000 or 2000, got %d", result)
	}
}
