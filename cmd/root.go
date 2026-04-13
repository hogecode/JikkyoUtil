package cmd

import (
	"os"

	"github.com/spf13/cobra"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "jikkyo",
	Short: "Get anime broadcast segment markers (ｷﾀ, A, B, C)",
	Long: `jikkyo is a CLI tool that finds the start times of anime broadcast segments.

It analyzes broadcast comments to identify:
  - ｷﾀ: Actual broadcast start time
  - A: First part (A part) start time
  - B: Second part (B part) start time
  - C: Third part (C part) start time`,
	Run: func(cmd *cobra.Command, args []string) {
		// Show help if no command provided
		cmd.Help()
	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	// Global flags can be added here
}
