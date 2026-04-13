package cmd

import (
	"os"

	"github.com/spf13/cobra"
)

var (
	title      string
	episode    int
	verbose    bool
	logFile    string
	outputDir  string
)

// NewRootCommand creates and returns the root command
func NewRootCommand() *cobra.Command {
	rootCmd := &cobra.Command{
		Use:   "jikkyo",
		Short: "Get anime broadcast segment markers (ｷﾀ, A, B, C)",
		Long: `jikkyo is a CLI tool that finds the start times of anime broadcast segments.

It analyzes broadcast comments to identify:
  - ｷﾀ: Actual broadcast start time
  - A: First part (A part) start time
  - B: Second part (B part) start time
  - C: Third part (C part) start time`,
		RunE: func(cmd *cobra.Command, args []string) error {
			// If no flags are provided, show help
			if len(args) == 0 && !cmd.Flags().Changed("title") {
				return cmd.Help()
			}
			return runjikkyo()
		},
	}

	// Define flags on root command
	rootCmd.Flags().StringVarP(&title, "title", "t", "", "Anime title to search for (required)")
	rootCmd.Flags().IntVarP(&episode, "episode", "e", 0, "Episode number (required)")
	rootCmd.Flags().BoolVarP(&verbose, "verbose", "v", false, "Enable verbose logging")
	rootCmd.Flags().StringVarP(&logFile, "log-file", "l", "", "Log file path (optional)")
	rootCmd.Flags().StringVarP(&outputDir, "output-dir", "o", "", "Output directory for program info file (optional)")

	// Mark required flags
	rootCmd.MarkFlagRequired("title")
	rootCmd.MarkFlagRequired("episode")

	return rootCmd
}

// Execute adds all child commands to the root command and sets flags appropriately.
func Execute() {
	rootCmd := NewRootCommand()
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}
