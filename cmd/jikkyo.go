package cmd

import (
	"fmt"
	"log/slog"
	"os"
	"path/filepath"

	"github.com/hogecode/JikkyoUtil/internal/api"
	"github.com/hogecode/JikkyoUtil/internal/presentation"
	"github.com/hogecode/JikkyoUtil/internal/usecase"
)

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
