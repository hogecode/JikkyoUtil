package usecase

import (
	"log/slog"
	"testing"
	"time"

	"github.com/hogecode/JikkyoUtil/internal/config"
	"github.com/hogecode/JikkyoUtil/internal/models"
)

func TestGenerateFilename(t *testing.T) {
	logger := slog.New(slog.NewTextHandler(nil, nil))
	gen := NewProgramFileGenerator(logger)

	progItem := &models.ProgItem{
		StTime: "2026-03-19 23:56:00",
		Count:  "11",
		STSubTitle: "運命にあらがう者たち",
	}

	filename, err := gen.GenerateFilename("エリスの聖杯", "11", progItem)
	if err != nil {
		t.Fatalf("failed to generate filename: %v", err)
	}

	expected := "202603192356000102-エリスの聖杯 第11話『運命にあらがう者たち』[字].ts.program.txt"
	if filename != expected {
		t.Errorf("filename mismatch\nexpected: %s\ngot: %s", expected, filename)
	}
}

func TestGenerateFileContent(t *testing.T) {
	logger := slog.New(slog.NewTextHandler(nil, nil))
	gen := NewProgramFileGenerator(logger)

	progItem := &models.ProgItem{
		StTime:    "2026-03-19 23:56:00",
		EdTime:    "2026-03-20 00:26:00",
		Count:     "11",
		STSubTitle: "運命にあらがう者たち",
		ChID:      "5", // TBS
	}

	channelMapping := config.NewChannelMapping()

	content, err := gen.GenerateFileContent("エリスの聖杯", "11", progItem, channelMapping)
	if err != nil {
		t.Fatalf("failed to generate content: %v", err)
	}

	// Verify content contains expected elements
	if !contains(content, "2026/03/19(木) 23:56～00:26") {
		t.Errorf("content does not contain date/time: %s", content)
	}

	if !contains(content, "TBS") {
		t.Errorf("content does not contain channel name TBS: %s", content)
	}

	if !contains(content, "エリスの聖杯 第11話『運命にあらがう者たち』[字]") {
		t.Errorf("content does not contain title/episode/subtitle: %s", content)
	}

	// Verify content ends with two newlines
	if !contains(content, "\n\n") {
		t.Errorf("content does not end with double newline: %s", content)
	}
}

func TestGenerateProgramFileInfo(t *testing.T) {
	logger := slog.New(slog.NewTextHandler(nil, nil))
	gen := NewProgramFileGenerator(logger)

	progItem := &models.ProgItem{
		StTime:    "2026-03-19 23:56:00",
		EdTime:    "2026-03-20 00:26:00",
		Count:     "11",
		STSubTitle: "運命にあらがう者たち",
		ChID:      "5", // TBS
	}

	channelMapping := config.NewChannelMapping()

	info, err := gen.GenerateProgramFileInfo("エリスの聖杯", "11", progItem, channelMapping, "/tmp")
	if err != nil {
		t.Fatalf("failed to generate program file info: %v", err)
	}

	if info.Filename == "" {
		t.Error("filename is empty")
	}

	if info.Content == "" {
		t.Error("content is empty")
	}

	if info.FullPath == "" {
		t.Error("fullPath is empty")
	}

	// Check that fullPath contains the filename
	if !contains(info.FullPath, info.Filename) {
		t.Errorf("fullPath does not contain filename: %s", info.FullPath)
	}
}

func TestGetWeekdayJapanese(t *testing.T) {
	logger := slog.New(slog.NewTextHandler(nil, nil))
	gen := NewProgramFileGenerator(logger)

	tests := []struct {
		weekday time.Weekday
		expected string
	}{
		{time.Sunday, "日"},
		{time.Monday, "月"},
		{time.Tuesday, "火"},
		{time.Wednesday, "水"},
		{time.Thursday, "木"},
		{time.Friday, "金"},
		{time.Saturday, "土"},
	}

	for _, tt := range tests {
		if result := gen.getWeekdayJapanese(tt.weekday); result != tt.expected {
			t.Errorf("weekday %v: expected %s, got %s", tt.weekday, tt.expected, result)
		}
	}
}

func TestGetChannelName(t *testing.T) {
	logger := slog.New(slog.NewTextHandler(nil, nil))
	gen := NewProgramFileGenerator(logger)

	channelMapping := config.NewChannelMapping()

	tests := []struct {
		chID     string
		expected string
	}{
		{"1", "NHK総合"},
		{"5", "TBS"},
		{"19", "TOKYO MX"},
		{"999", ""}, // Non-existent channel
	}

	for _, tt := range tests {
		if result := gen.getChannelName(tt.chID, channelMapping); result != tt.expected {
			t.Errorf("chID %s: expected %s, got %s", tt.chID, tt.expected, result)
		}
	}
}
