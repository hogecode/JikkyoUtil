package usecase

import (
	"fmt"
	"io"
	"log/slog"
	"strconv"
	"strings"

	"github.com/hogecode/getabc/internal/api"
	"github.com/hogecode/getabc/internal/models"
)

// TitleSearchUseCase handles title search and hogecode selection
type TitleSearchUseCase struct {
	client *api.Client
	logger *slog.Logger
	input  io.Reader
}

// NewTitleSearchUseCase creates a new title search use case
func NewTitleSearchUseCase(client *api.Client, logger *slog.Logger, input io.Reader) *TitleSearchUseCase {
	return &TitleSearchUseCase{
		client: client,
		logger: logger,
		input:  input,
	}
}

// SearchAndSelect searches for titles and returns the selected one
func (uc *TitleSearchUseCase) SearchAndSelect(query string) (*models.Title, error) {
	// Call API
	resp, err := uc.client.TitleSearch(query)
	if err != nil {
		return nil, err
	}

	if len(resp.Titles) == 0 {
		return nil, fmt.Errorf("no titles found matching '%s'", query)
	}

	// Convert map to slice for easier selection
	titles := make([]*models.Title, 0, len(resp.Titles))
	for _, title := range resp.Titles {
		t := title
		titles = append(titles, &t)
	}

	uc.logger.Debug("found titles", slog.Int("count", len(titles)))

	// If only one result, return it
	if len(titles) == 1 {
		uc.logger.Info("only one title found", slog.String("title", titles[0].Title))
		return titles[0], nil
	}

	// Multiple results - show selection menu
	return uc.promptUserSelection(titles)
}

// promptUserSelection displays a menu and gets hogecode selection
func (uc *TitleSearchUseCase) promptUserSelection(titles []*models.Title) (*models.Title, error) {
	fmt.Println("\nMultiple titles found:")
	for i, title := range titles {
		fmt.Printf("%d. %s (TID: %s)\n", i+1, title.Title, title.TID)
	}

	fmt.Print("\nSelect (1-" + strconv.Itoa(len(titles)) + "): ")

	// Read hogecode input
	var input string
	_, err := fmt.Fscanln(uc.input, &input)
	if err != nil {
		return nil, fmt.Errorf("failed to read input: %w", err)
	}

	// Parse selection
	selection, err := strconv.Atoi(strings.TrimSpace(input))
	if err != nil {
		return nil, fmt.Errorf("invalid selection: %w", err)
	}

	if selection < 1 || selection > len(titles) {
		return nil, fmt.Errorf("selection out of range (1-%d)", len(titles))
	}

	selected := titles[selection-1]
	uc.logger.Info("title selected", slog.String("title", selected.Title), slog.String("tid", selected.TID))

	return selected, nil
}
