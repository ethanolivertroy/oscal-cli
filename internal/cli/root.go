package cli

import (
	"github.com/spf13/cobra"
)

var (
	// Version information set at build time
	Version   = "dev"
	GitCommit = "unknown"
	BuildDate = "unknown"
)

var rootCmd = &cobra.Command{
	Use:   "oscal",
	Short: "OSCAL CLI - Work with OSCAL documents",
	Long: `OSCAL CLI is a command-line tool for working with Open Security Controls
Assessment Language (OSCAL) documents.

It supports reading and writing OSCAL documents in XML, JSON, and YAML formats,
validating documents against OSCAL schemas, and converting between formats.`,
}

func Execute() error {
	return rootCmd.Execute()
}

func init() {
	rootCmd.AddCommand(versionCmd)
}

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print version information",
	Run: func(cmd *cobra.Command, args []string) {
		cmd.Printf("oscal-cli %s\n", Version)
		cmd.Printf("  commit: %s\n", GitCommit)
		cmd.Printf("  built:  %s\n", BuildDate)
	},
}
