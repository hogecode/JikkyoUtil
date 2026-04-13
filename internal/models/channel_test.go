package models

import (
	"testing"
)

func TestChannelMapping(t *testing.T) {
	// Test that Channel struct is properly defined
	ch := &Channel{
		ChID:     1,
		ChGID:    11,
		ChName:   "NHK総合",
		JikkyoID: "jk1",
	}

	if ch.ChID != 1 {
		t.Errorf("expected ChID 1, got %d", ch.ChID)
	}
	if ch.ChName != "NHK総合" {
		t.Errorf("expected ChName 'NHK総合', got %s", ch.ChName)
	}
	if ch.JikkyoID != "jk1" {
		t.Errorf("expected JikkyoID 'jk1', got %s", ch.JikkyoID)
	}

	// Test ChannelMapping
	mapping := ChannelMapping{
		"NHK総合": ch,
	}

	found, ok := mapping["NHK総合"]
	if !ok {
		t.Error("expected to find channel in mapping")
	}
	if found.ChID != 1 {
		t.Errorf("expected ChID 1, got %d", found.ChID)
	}
}
