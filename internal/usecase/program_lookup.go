package usecase

import (
	"fmt"
	"log/slog"
	"strconv"
	"time"

	"github.com/hogecode/JikkyoUtil/internal/api"
	"github.com/hogecode/JikkyoUtil/internal/config"
	"github.com/hogecode/JikkyoUtil/internal/models"
)

// ProgramLookupUseCase handles program lookup
type ProgramLookupUseCase struct {
	client         *api.Client
	logger         *slog.Logger
	channelMapping models.ChannelMapping
}

// NewProgramLookupUseCase creates a new program lookup use case
func NewProgramLookupUseCase(client *api.Client, logger *slog.Logger, channelMapping models.ChannelMapping) *ProgramLookupUseCase {
	return &ProgramLookupUseCase{
		client:         client,
		logger:         logger,
		channelMapping: channelMapping,
	}
}

// LookupProgram fetches program information
func (uc *ProgramLookupUseCase) LookupProgram(title *models.Title, episode int) (*models.ProgItem, string, error) {

	uc.logger.Debug("looking up program with ChIDs", slog.String("title", title.TID))

	// Call ProgLookup API with all fixed ChIDs (1-10) in a single request
	// The API supports multiple ChIDs as comma-separated values: ChID=1,2,3,4,5,6,7,8,9,10
	chIDsParam := "1,2,3,4,5,6,7,8,9,19" // Using 1-9 and 19 (TOKYO MX) as they are the main channels with good Jikkyo support
	resp, err := uc.client.ProgLookup(title.TID, chIDsParam, episode)
	if err != nil {
		return nil, "", fmt.Errorf("failed to call ProgLookup API: %w", err)
	}

	if len(resp.ProgItems) == 0 {
		return nil, "", fmt.Errorf("no programs found for TID: %s, episode: %d across all fixed ChIDs (1-10)", title.TID, episode)
	}

	// Build JikkyoID mapping from channel mapping
	var jikkyoIDMap = make(map[string]string) // ChID -> JikkyoID mapping
	for channelName, channel := range uc.channelMapping {
		jikkyoIDMap[strconv.Itoa(channel.ChID)] = channel.JikkyoID
		uc.logger.Debug("built channel mapping", slog.String("channelName", channelName), slog.Int("chid", channel.ChID), slog.String("jikkyoid", channel.JikkyoID))
	}

	var allProgItems []models.ProgItem

	// Collect all program items
	for i := range resp.ProgItems {
		allProgItems = append(allProgItems, resp.ProgItems[i])
	}

	if len(allProgItems) == 0 {
		return nil, "", fmt.Errorf("no program items found in response for TID: %s, episode: %d", title.TID, episode)
	}

	// Find best non-deleted ProgItem
	// Prefer channels with popular Jikkyo support (jk1-jk9) for better comment availability
	var progItem *models.ProgItem
	var fallbackProgItem *models.ProgItem

	for i := range allProgItems {
		if allProgItems[i].Deleted == 0 {
			// Get the JikkyoID for this item
			jikkyoID := jikkyoIDMap[allProgItems[i].ChID]
			if jikkyoID != "" && config.IsPopularChannel(jikkyoID) {
				progItem = &allProgItems[i]
				break // Found a popular channel, use it
			}
			// Keep first non-deleted as fallback
			if fallbackProgItem == nil {
				fallbackProgItem = &allProgItems[i]
			}
		}
	}

	// Use fallback if no popular channel found
	if progItem == nil {
		progItem = fallbackProgItem
	}

	if progItem == nil {
		return nil, "", fmt.Errorf("all program items are marked as deleted")
	}

	// Get the appropriate JikkyoID for the selected program
	selectedJikkyoID := jikkyoIDMap[progItem.ChID]
	if selectedJikkyoID == "" {
		// Fallback to the mapped channel's JikkyoID if available
		if channel, ok := uc.channelMapping[title.FirstCh]; ok {
			selectedJikkyoID = channel.JikkyoID
		}
	}

	uc.logger.Info("found program",
		slog.String("sttime", progItem.StTime),
		slog.String("edtime", progItem.EdTime),
		slog.String("subtitle", progItem.STSubTitle),
		slog.String("chid", progItem.ChID))

	return progItem, selectedJikkyoID, nil
}

// ParseProgItemTimes parses StTime and EdTime from ProgItem and converts to Unix timestamps
// The times are assumed to be in JST (Japan Standard Time, UTC+9)
func ParseProgItemTimes(progItem *models.ProgItem) (int64, int64, error) {
	// Parse StTime: "2021-01-28 19:30:00" - assumed to be JST
	stTime, err := time.ParseInLocation(config.TimeFormat, progItem.StTime, config.JST)
	if err != nil {
		return 0, 0, fmt.Errorf("failed to parse StTime: %w", err)
	}

	// Parse EdTime: "2021-01-28 20:00:00" - assumed to be JST
	edTime, err := time.ParseInLocation(config.TimeFormat, progItem.EdTime, config.JST)
	if err != nil {
		return 0, 0, fmt.Errorf("failed to parse EdTime: %w", err)
	}

	return stTime.Unix(), edTime.Unix(), nil
}

// FormatUnixTimestamp converts Unix timestamp back to formatted string in JST
func FormatUnixTimestamp(unixTime int64) string {
	return time.Unix(unixTime, 0).In(config.JST).Format(config.TimeFormat)
}

// ConvertStTimeToUnix converts the StTime string to Unix timestamp (assumes JST)
func ConvertStTimeToUnix(stTime string) (int64, error) {
	t, err := time.ParseInLocation(config.TimeFormat, stTime, config.JST)
	if err != nil {
		return 0, fmt.Errorf("failed to parse time: %w", err)
	}
	return t.Unix(), nil
}

// ConvertStringToInt converts a string to int
func ConvertStringToInt(s string) (int, error) {
	return strconv.Atoi(s)
}
