package cmd

import (
	"fmt"
	"log/slog"
	"os"
	"path/filepath"

	"github.com/hogecode/JikkyoUtil/internal/api"
	"github.com/hogecode/JikkyoUtil/internal/presentation"
	"github.com/hogecode/JikkyoUtil/internal/usecase"
	"github.com/spf13/cobra"
)

var (
	title      string
	episode    int
	verbose    bool
	logFile    string
	outputDir  string
)

// jikkyoCmd represents the jikkyo command
var jikkyoCmd = &cobra.Command{
	Use:   "jikkyo",
	Short: "Get anime broadcast segment markers",
	Long:  `Search for an anime title and get the start times of broadcast segments (ｷﾀ, A, B, C).`,
	Args:  cobra.NoArgs,
	RunE: func(cmd *cobra.Command, args []string) error {
		return runjikkyo()
	},
}

// runjikkyo is the main execution logic
func runjikkyo() error {
	// Validate inputs
	if title == "" {
		return errorf("title is required (use -t or --title)")
	}
	if episode <= 0 {
		return errorf("episode must be a positive number (use -e or --episode)")
	}

	// Setup logger
	loggerCfg := presentation.LoggerConfig{
		Verbose: verbose,
		Output:  os.Stderr,
		LogFile: logFile,
	}
	logger, err := presentation.NewLogger(loggerCfg)
	if err != nil {
		return errorf("failed to create logger: %v", err)
	}

	// Create API client
	client := api.NewClient()

	// Create core use case
	coreUC := usecase.NewCoreUseCase(client, logger, os.Stdin)

	// Execute workflow
	result, err := coreUC.Execute(title, episode)
	if err != nil {
		return err
	}

	// Output result
	outputter := presentation.NewOutputFormatter(verbose)
	outputter.PrintResult(result)

	// Write program info file if output directory is specified
	if outputDir != "" && result.ProgramFileName != "" && result.ProgramContent != "" {
		filePath := filepath.Join(outputDir, result.ProgramFileName)
		err := os.WriteFile(filePath, []byte(result.ProgramContent), 0644)
		if err != nil {
			logger.Error("failed to write program info file",
				slog.String("path", filePath),
				slog.String("error", err.Error()))
		} else {
			logger.Info("program info file written successfully",
				slog.String("path", filePath))
		}
	}

	return nil
}

// errorf returns an error with a formatted message
func errorf(msg string, args ...interface{}) error {
	return returnError{fmt.Sprintf(msg, args...)}
}

type returnError struct {
	msg string
}

func (e returnError) Error() string {
	return e.msg
}

func init() {
	rootCmd.AddCommand(jikkyoCmd)

	// Define flags
	jikkyoCmd.Flags().StringVarP(&title, "title", "t", "", "Anime title to search for (required)")
	jikkyoCmd.Flags().IntVarP(&episode, "episode", "e", 0, "Episode number (required)")
	jikkyoCmd.Flags().BoolVarP(&verbose, "verbose", "v", false, "Enable verbose logging")
	jikkyoCmd.Flags().StringVarP(&logFile, "log-file", "l", "", "Log file path (optional)")
	jikkyoCmd.Flags().StringVarP(&outputDir, "output-dir", "o", "", "Output directory for program info file (optional)")

	// Mark required flags
	jikkyoCmd.MarkFlagRequired("title")
	jikkyoCmd.MarkFlagRequired("episode")
}
