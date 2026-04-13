package usecase

import (
	"fmt"
	"io"
	"log/slog"
	"strconv"

	"github.com/hogecode/JikkyoUtil/internal/api"
	"github.com/hogecode/JikkyoUtil/internal/config"
	"github.com/hogecode/JikkyoUtil/internal/models"
)

// CoreUseCase orchestrates the main workflow
type CoreUseCase struct {
	apiClient      *api.Client
	logger         *slog.Logger
	channelMapping models.ChannelMapping
	titleSearch    *TitleSearchUseCase
	programLookup  *ProgramLookupUseCase
	jikkyoAnalysis *JikkyoAnalysisUseCase
	programFileGen *ProgramFileGenerator
	input          io.Reader
}

// NewCoreUseCase creates a new core use case
func NewCoreUseCase(client *api.Client, logger *slog.Logger, input io.Reader) *CoreUseCase {
	channelMapping := config.NewChannelMapping()

	// Set logger for HTTP request/response logging
	client.SetLogger(logger)

	return &CoreUseCase{
		apiClient:      client,
		logger:         logger,
		channelMapping: channelMapping,
		titleSearch:    NewTitleSearchUseCase(client, logger, input),
		programLookup:  NewProgramLookupUseCase(client, logger, channelMapping),
		jikkyoAnalysis: NewJikkyoAnalysisUseCase(client, logger),
		programFileGen: NewProgramFileGenerator(logger),
		input:          input,
	}
}

// Execute runs the main workflow
func (uc *CoreUseCase) Execute(titleQuery string, episode int) (*models.JikkyoResult, error) {
	uc.logger.Debug("starting jikkyo workflow",
		slog.String("title_query", titleQuery),
		slog.Int("episode", episode))

	// Step 1: Search and select title
	title, err := uc.titleSearch.SearchAndSelect(titleQuery)
	if err != nil {
		return nil, fmt.Errorf("title search failed: %w", err)
	}

	// Step 2: Lookup program information
	progItem, jikkyoID, err := uc.programLookup.LookupProgram(title, episode)
	if err != nil {
		return nil, fmt.Errorf("program lookup failed: %w", err)
	}

	// Step 3: Parse times
	stUnix, edUnix, err := ParseProgItemTimes(progItem)
	if err != nil {
		return nil, fmt.Errorf("time parsing failed: %w", err)
	}

	// Step 4: Analyze Jikkyo comments
	analysis, err := uc.jikkyoAnalysis.AnalyzeComments(jikkyoID, stUnix, edUnix)
	if err != nil {
		return nil, fmt.Errorf("comment analysis failed: %w", err)
	}

	// Step 5: Build result
	result := uc.buildResult(title, episode, progItem, analysis)

	// Step 6: Generate program info file
	programFileInfo, err := uc.programFileGen.GenerateProgramFileInfo(
		title.Title,
		strconv.Itoa(episode),
		progItem,
		uc.channelMapping,
		"", // outputDir will be handled by CLI
	)
	if err != nil {
		uc.logger.Warn("failed to generate program file info",
			slog.String("error", err.Error()))
		// Continue without program file generation
	} else {
		result.ProgramFileName = programFileInfo.Filename
		result.ProgramContent = programFileInfo.Content
	}

	uc.logger.Info("workflow completed successfully")

	return result, nil
}

// buildResult constructs the final JikkyoResult
func (uc *CoreUseCase) buildResult(
	title *models.Title,
	episode int,
	progItem *models.ProgItem,
	analysis *models.CommentAnalysis,
) *models.JikkyoResult {

	// Convert episode to int if needed
	episodeInt, _ := strconv.Atoi(strconv.Itoa(episode))

	result := &models.JikkyoResult{
		Title:    title.Title,
		Episode:  episodeInt,
		SubTitle: progItem.STSubTitle,
		Start:    progItem.StTime,
	}

	// Calculate actual times based on comment markers
	if analysis.KitaTime > 0 {
		result.RealStartTime = FormatUnixTimestamp(analysis.KitaTime)
	} else {
		// Fallback to StTime if no ｷﾀ comment found
		result.RealStartTime = progItem.StTime
	}

	if analysis.ATime > 0 {
		result.A = FormatUnixTimestamp(analysis.ATime)
	}

	if analysis.BTime > 0 {
		result.B = FormatUnixTimestamp(analysis.BTime)
	}

	if analysis.CTime > 0 {
		result.C = FormatUnixTimestamp(analysis.CTime)
	}

	return result
}
