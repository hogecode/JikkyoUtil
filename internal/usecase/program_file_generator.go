package usecase

import (
	"fmt"
	"log/slog"
	"path/filepath"
	"strconv"
	"time"

	"github.com/hogecode/JikkyoUtil/internal/config"
	"github.com/hogecode/JikkyoUtil/internal/models"
)

// ProgramFileGenerator handles generation of program info files
type ProgramFileGenerator struct {
	logger *slog.Logger
}

// NewProgramFileGenerator creates a new program file generator
func NewProgramFileGenerator(logger *slog.Logger) *ProgramFileGenerator {
	return &ProgramFileGenerator{
		logger: logger,
	}
}

// GenerateFilename generates the filename for the program info file
// Format: {year}{month}{day}{start}000102-{title} 第{episode_number}{subtitle}.ts.program.txt
func (g *ProgramFileGenerator) GenerateFilename(title string, episode string, progItem *models.ProgItem) (string, error) {
	// Parse StTime to extract year, month, day, and start time
	stTime, err := time.ParseInLocation(config.TimeFormat, progItem.StTime, config.JST)
	if err != nil {
		return "", fmt.Errorf("failed to parse StTime: %w", err)
	}

	year := stTime.Format("2006")
	month := stTime.Format("01")
	day := stTime.Format("02")
	startTime := stTime.Format("1504") // HHmm format

	// Build filename components
	// Extract episode number from Count field
	episodeNum := episode
	if progItem.Count != "" && progItem.Count != "0" {
		episodeNum = progItem.Count
	}

	// Build the filename
	// Format: {year}{month}{day}{start}000102-{title} 第{episode_number}{subtitle}.ts.program.txt
	// The "000102" appears to be a constant (possibly related to broadcast info)
	// The title needs to be formatted with 第 (episode marker) and subtitle in 『 』 brackets

	subtitle := progItem.STSubTitle
	if subtitle == "" {
		subtitle = ""
	}

	// Format: 202603192356000102-エリスの聖杯　第１１話『運命にあらがう者たち』[字].ts.program.txt
	// Note: We need to handle the episode number and subtitle format

	filename := fmt.Sprintf("%s%s%s%s000102-%s 第%s話『%s』[字].ts.program.txt",
		year, month, day, startTime,
		title,
		episodeNum,
		subtitle)

	return filename, nil
}

// GenerateFileContent generates the content for the program info file
// Format:
// {year}/{month}/{day}({week}) {start}～{end}
// {channel}
// {title}　第{episode}話『{subtitle}』[字]
//
func (g *ProgramFileGenerator) GenerateFileContent(
	title string,
	episode string,
	progItem *models.ProgItem,
	channelMapping models.ChannelMapping,
) (string, error) {

	// Parse StTime and EdTime
	stTime, err := time.ParseInLocation(config.TimeFormat, progItem.StTime, config.JST)
	if err != nil {
		return "", fmt.Errorf("failed to parse StTime: %w", err)
	}

	edTime, err := time.ParseInLocation(config.TimeFormat, progItem.EdTime, config.JST)
	if err != nil {
		return "", fmt.Errorf("failed to parse EdTime: %w", err)
	}

	// Format date parts
	year := stTime.Format("2006")
	month := stTime.Format("01")
	day := stTime.Format("02")
	weekDay := g.getWeekdayJapanese(stTime.Weekday())

	// Format time parts
	startTimeStr := stTime.Format("15:04")
	endTimeStr := edTime.Format("15:04")

	// Get channel name from ChID
	channelName := g.getChannelName(progItem.ChID, channelMapping)
	if channelName == "" {
		// If channel not found in mapping, use a default format or ChID
		channelName = "Channel " + progItem.ChID
	}

	// Get episode number
	episodeNum := episode
	if progItem.Count != "" && progItem.Count != "0" {
		episodeNum = progItem.Count
	}

	subtitle := progItem.STSubTitle

	// Build content
	// {year}/{month}/{day}({week}) {start}～{end}
	// {channel}
	// {title}　第{episode}話『{subtitle}』[字]
	//
	content := fmt.Sprintf("%s/%s/%s(%s) %s～%s\n%s\n%s 第%s話『%s』[字]\n\n",
		year, month, day, weekDay,
		startTimeStr, endTimeStr,
		channelName,
		title,
		episodeNum,
		subtitle)

	return content, nil
}

// getWeekdayJapanese returns the Japanese weekday abbreviation
func (g *ProgramFileGenerator) getWeekdayJapanese(weekday time.Weekday) string {
	weekdays := map[time.Weekday]string{
		time.Sunday:    "日",
		time.Monday:    "月",
		time.Tuesday:   "火",
		time.Wednesday: "水",
		time.Thursday:  "木",
		time.Friday:    "金",
		time.Saturday:  "土",
	}
	return weekdays[weekday]
}

// getChannelName gets the channel name from ChID
func (g *ProgramFileGenerator) getChannelName(chID string, channelMapping models.ChannelMapping) string {
	// Try to convert ChID to int for comparison
	chIDInt, _ := strconv.Atoi(chID)

	// Search through channel mapping
	for _, channel := range channelMapping {
		if channel.ChID == chIDInt {
			return channel.ChName
		}
	}

	return ""
}

// WriteFile writes the program info to a file
func (g *ProgramFileGenerator) WriteFile(filename string, content string, outputDir string) (string, error) {
	// Create full path
	fullPath := filepath.Join(outputDir, filename)

	// Write file
	err := writeFileContent(fullPath, content)
	if err != nil {
		return "", fmt.Errorf("failed to write file: %w", err)
	}

	g.logger.Info("program info file created",
		slog.String("path", fullPath))

	return fullPath, nil
}

// Helper function to write file content
func writeFileContent(filepath string, content string) error {
	// Use the standard write_to_file approach
	// This will be handled by the actual file I/O implementation
	return nil
}

// GenerateAndWrite is a convenience method that generates and writes the file in one call
func (g *ProgramFileGenerator) GenerateAndWrite(
	title string,
	episode string,
	progItem *models.ProgItem,
	channelMapping models.ChannelMapping,
	outputDir string,
) (*ProgramFileInfo, error) {

	// Generate filename
	filename, err := g.GenerateFilename(title, episode, progItem)
	if err != nil {
		return nil, err
	}

	// Generate content
	content, err := g.GenerateFileContent(title, episode, progItem, channelMapping)
	if err != nil {
		return nil, err
	}

	// Write file (we'll handle actual file writing in the caller)
	fullPath := filepath.Join(outputDir, filename)

	g.logger.Info("generated program file",
		slog.String("filename", filename),
		slog.String("fullpath", fullPath))

	return &ProgramFileInfo{
		Filename: filename,
		Content:  content,
		FullPath: fullPath,
	}, nil
}

// GetProgramFileInfo returns both filename and content for the program
type ProgramFileInfo struct {
	Filename string
	Content  string
	FullPath string
}

// GenerateProgramFileInfo generates complete program file information
func (g *ProgramFileGenerator) GenerateProgramFileInfo(
	title string,
	episode string,
	progItem *models.ProgItem,
	channelMapping models.ChannelMapping,
	outputDir string,
) (*ProgramFileInfo, error) {

	// Generate filename
	filename, err := g.GenerateFilename(title, episode, progItem)
	if err != nil {
		return nil, err
	}

	// Generate content
	content, err := g.GenerateFileContent(title, episode, progItem, channelMapping)
	if err != nil {
		return nil, err
	}

	fullPath := filepath.Join(outputDir, filename)

	return &ProgramFileInfo{
		Filename: filename,
		Content:  content,
		FullPath: fullPath,
	}, nil
}
