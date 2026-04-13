package usecase

import (
	"log/slog"
	"strconv"
	"strings"

	"github.com/hogecode/JikkyoUtil/internal/api"
	"github.com/hogecode/JikkyoUtil/internal/models"
)

// JikkyoAnalysisUseCase analyzes Jikkyo comments
type JikkyoAnalysisUseCase struct {
	client *api.Client
	logger *slog.Logger
}

// NewJikkyoAnalysisUseCase creates a new Jikkyo analysis use case
func NewJikkyoAnalysisUseCase(client *api.Client, logger *slog.Logger) *JikkyoAnalysisUseCase {
	return &JikkyoAnalysisUseCase{
		client: client,
		logger: logger,
	}
}

// AnalyzeComments fetches and analyzes Jikkyo comments for markers (ｷﾀ, A, B, C)
func (uc *JikkyoAnalysisUseCase) AnalyzeComments(jikkyoID string, stUnix, edUnix int64) (*models.CommentAnalysis, error) {
	// Fetch comments
	resp, err := uc.client.GetJikkyoComments(jikkyoID, stUnix, edUnix)
	if err != nil {
		return nil, err
	}

	if len(resp.Packets) == 0 {
		uc.logger.Warn("no comments found in time range",
			slog.String("jikkyo_id", jikkyoID),
			slog.Int64("start_unix", stUnix),
			slog.Int64("end_unix", edUnix))
		// Return empty analysis with zero times
		// This can happen for older broadcasts or less popular channels
		return &models.CommentAnalysis{}, nil
	}

	uc.logger.Debug("fetched comments", slog.Int("count", len(resp.Packets)))

	// Count occurrences of markers (ｷﾀ, A, B, C) by timestamp
	kitaMarkers := make(map[int64]int)
	aMarkers := make(map[int64]int)
	bMarkers := make(map[int64]int)
	cMarkers := make(map[int64]int)

	for _, packet := range resp.Packets {
		content := packet.Chat.Content
		date := convertToInt64(packet.Chat.Date)
		if date == 0 {
			continue // Skip if unable to parse date
		}

		if strings.Contains(content, "ｷﾀ") {
			kitaMarkers[date]++
		}
		trimmed := strings.TrimSpace(content)
		if trimmed == "A" {
			aMarkers[date]++
		}
		if trimmed == "B" {
			bMarkers[date]++
		}
		if trimmed == "C" {
			cMarkers[date]++
		}
	}

	// Find the most common timestamp for each marker
	analysis := &models.CommentAnalysis{
		KitaTime: findMostCommonTime(kitaMarkers),
		ATime:    findMostCommonTime(aMarkers),
		BTime:    findMostCommonTime(bMarkers),
		CTime:    findMostCommonTime(cMarkers),
	}

	uc.logger.Info("comment analysis complete",
		slog.Int64("kita_time", analysis.KitaTime),
		slog.Int64("a_time", analysis.ATime),
		slog.Int64("b_time", analysis.BTime),
		slog.Int64("c_time", analysis.CTime))

	return analysis, nil
}

// findMostCommonTime finds the timestamp with the most occurrences
func findMostCommonTime(markers map[int64]int) int64 {
	if len(markers) == 0 {
		return 0
	}

	var maxTime int64
	maxCount := 0

	for timestamp, count := range markers {
		if count > maxCount {
			maxCount = count
			maxTime = timestamp
		}
	}

	return maxTime
}

// convertToInt64 converts date field (which can be string or int) to int64
func convertToInt64(v interface{}) int64 {
	switch val := v.(type) {
	case int64:
		return val
	case int:
		return int64(val)
	case float64:
		return int64(val)
	case string:
		parsed, err := strconv.ParseInt(val, 10, 64)
		if err != nil {
			return 0
		}
		return parsed
	default:
		return 0
	}
}
