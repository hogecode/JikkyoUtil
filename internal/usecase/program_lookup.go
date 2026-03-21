package usecase

import (
	"fmt"
	"log/slog"
	"strconv"
	"time"

	"github.com/hogecode/getabc/internal/api"
	"github.com/hogecode/getabc/internal/config"
	"github.com/hogecode/getabc/internal/models"
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
	// Check if FirstCh is empty (common for movies)
	if title.FirstCh == "" {
		return nil, "", fmt.Errorf("channel information not available - this title may be a movie or special content (TID: %s)", title.TID)
	}

	// Find channel by name
	channel, ok := uc.channelMapping[title.FirstCh]
	if !ok {
		return nil, "", fmt.Errorf("unsupported channel: %s (supported channels: see README)", title.FirstCh)
	}

	uc.logger.Debug("found channel", slog.String("name", title.FirstCh), slog.Int("chid", channel.ChID))

	// Call ProgLookup API
	resp, err := uc.client.ProgLookup(title.TID, channel.ChID, episode)
	if err != nil {
		return nil, "", err
	}

	if len(resp.ProgItems) == 0 {
		return nil, "", fmt.Errorf("no programs found for TID: %s, episode: %d", title.TID, episode)
	}

	// Find best non-deleted ProgItem
	// Prefer channels with popular Jikkyo support (jk1-jk9) for better comment availability
	var progItem *models.ProgItem
	var fallbackProgItem *models.ProgItem

	for i := range resp.ProgItems {
		if resp.ProgItems[i].Deleted == 0 {
			// Check if this channel has good Jikkyo support
			// The current channel's JikkyoID is already determined from channel.JikkyoID
			// But for multiple ProgItems from different ChIDs, we need to check each
			progItemChID := resp.ProgItems[i].ChID
			if progItemChID == strconv.Itoa(channel.ChID) {
				// This is from the same channel we looked up
				if config.IsPopularChannel(channel.JikkyoID) {
					progItem = &resp.ProgItems[i]
					break // Found a popular channel, use it
				}
			}
			// Keep first non-deleted as fallback
			if fallbackProgItem == nil {
				fallbackProgItem = &resp.ProgItems[i]
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

	uc.logger.Info("found program",
		slog.String("sttime", progItem.StTime),
		slog.String("edtime", progItem.EdTime),
		slog.String("subtitle", progItem.STSubTitle))

	return progItem, channel.JikkyoID, nil
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
