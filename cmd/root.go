package cmd

import "github.com/spf13/cobra"

var rootCmd = &cobra.Command{
	Use:   "sorteiagram",
	Short: "An auto-sweepstakes for Instagram",
	Long: `SorteiaGram is a CLI application which interfaces with the non-official Instagram's third-party API, GoInsta.
This tool will participate automatically in Instagram's sweepstakes following a set of pre-defined rules.`,
}

// Execute - Start the CLI
func Execute() error {
	return rootCmd.Execute()
}
