package config

import (
	"testing"
)

func TestNewChannelMapping(t *testing.T) {
	mapping := NewChannelMapping()

	// Test that mapping is created
	if len(mapping) == 0 {
		t.Fatal("expected non-empty channel mapping")
	}

	// Test specific channels
	tests := []struct {
		name     string
		chID     int
		jikkyoID string
	}{
		{"NHK総合", 1, "jk1"},
		{"テレビ東京", 7, "jk7"},
		{"TOKYO MX", 19, "jk9"},
		{"AT-X", 20, "jk333"},
	}

	for _, test := range tests {
		ch, ok := mapping[test.name]
		if !ok {
			t.Errorf("expected to find channel %s in mapping", test.name)
			continue
		}
		if ch.ChID != test.chID {
			t.Errorf("expected ChID %d for %s, got %d", test.chID, test.name, ch.ChID)
		}
		if ch.JikkyoID != test.jikkyoID {
			t.Errorf("expected JikkyoID %s for %s, got %s", test.jikkyoID, test.name, ch.JikkyoID)
		}
		if ch.ChName != test.name {
			t.Errorf("expected ChName %s, got %s", test.name, ch.ChName)
		}
	}
}

func TestTimeFormat(t *testing.T) {
	// TimeFormat should be a valid Go time format string
	expected := "2006-01-02 15:04:05"
	if TimeFormat != expected {
		t.Errorf("expected TimeFormat %q, got %q", expected, TimeFormat)
	}
}

func TestAPIEndpoints(t *testing.T) {
	// Test that API endpoints are defined
	if SyoboiTitleSearchURL == "" {
		t.Error("SyoboiTitleSearchURL should not be empty")
	}
	if SyoboiProgLookupURL == "" {
		t.Error("SyoboiProgLookupURL should not be empty")
	}
	if JikkyoBaseURL == "" {
		t.Error("JikkyoBaseURL should not be empty")
	}

	// Test that URLs are properly formatted
	if SyoboiTitleSearchURL != "http://cal.syoboi.jp/json" {
		t.Errorf("unexpected SyoboiTitleSearchURL: %s", SyoboiTitleSearchURL)
	}
	if SyoboiProgLookupURL != "http://cal.syoboi.jp/db" {
		t.Errorf("unexpected SyoboiProgLookupURL: %s", SyoboiProgLookupURL)
	}
	if JikkyoBaseURL != "https://jikkyo.tsukumijima.net/api/kakolog" {
		t.Errorf("unexpected JikkyoBaseURL: %s", JikkyoBaseURL)
	}
}
