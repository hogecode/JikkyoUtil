package presentation

import (
	"fmt"

	"github.com/hogecode/getabc/internal/models"
)

// OutputFormatter handles result formatting and display
type OutputFormatter struct {
	verbose bool
}

// NewOutputFormatter creates a new output formatter
func NewOutputFormatter(verbose bool) *OutputFormatter {
	return &OutputFormatter{verbose: verbose}
}

// PrintResult prints the result in the specified format
func (of *OutputFormatter) PrintResult(result *models.GetABCResult) {
	fmt.Printf("(title) %s\n", result.Title)
	fmt.Printf("(episode) %d\n", result.Episode)
	fmt.Printf("(subtitle) %s\n", result.SubTitle)
	fmt.Printf("(start) %s\n", result.Start)
	fmt.Printf("(real_start_time) %s\n", result.RealStartTime)

	if result.A != "" {
		fmt.Printf("(A) %s\n", result.A)
	}

	if result.B != "" {
		fmt.Printf("(B) %s\n", result.B)
	}

	if result.C != "" {
		fmt.Printf("(C) %s\n", result.C)
	}
}
