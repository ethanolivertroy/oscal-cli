package cli

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/ethantroy/oscal-cli/pkg/oscal/io"
	"github.com/spf13/cobra"
)

var (
	convertTo        string
	convertOutput    string
	convertOverwrite bool
)

var convertCmd = &cobra.Command{
	Use:   "convert <input-file>",
	Short: "Convert OSCAL documents between formats",
	Long: `Convert OSCAL documents between XML, JSON, and YAML formats.

The input format is auto-detected from the file extension or content.
The output format must be specified with --to.

Examples:
  oscal convert catalog.xml --to json
  oscal convert catalog.json --to yaml --output catalog.yaml
  oscal convert catalog.xml --to json > catalog.json`,
	Args: cobra.ExactArgs(1),
	RunE: runConvert,
}

func init() {
	rootCmd.AddCommand(convertCmd)

	convertCmd.Flags().StringVarP(&convertTo, "to", "t", "", "Output format (json, xml, yaml) [required]")
	convertCmd.Flags().StringVarP(&convertOutput, "output", "o", "", "Output file (default: stdout)")
	convertCmd.Flags().BoolVar(&convertOverwrite, "overwrite", false, "Overwrite output file if it exists")

	convertCmd.MarkFlagRequired("to")
}

func runConvert(cmd *cobra.Command, args []string) error {
	inputPath := args[0]

	// Parse output format
	outputFormat, err := io.ParseFormat(convertTo)
	if err != nil {
		return fmt.Errorf("invalid output format: %w", err)
	}

	// Load input document
	loader := io.NewLoader()
	doc, err := loader.Load(inputPath)
	if err != nil {
		return fmt.Errorf("failed to load document: %w", err)
	}

	// Create writer
	writer := io.NewWriter()

	// Determine output destination
	if convertOutput != "" {
		// Check if output file exists
		if !convertOverwrite {
			if _, err := os.Stat(convertOutput); err == nil {
				return fmt.Errorf("output file already exists: %s (use --overwrite to replace)", convertOutput)
			}
		}

		// Ensure output directory exists
		outputDir := filepath.Dir(convertOutput)
		if outputDir != "." && outputDir != "" {
			if err := os.MkdirAll(outputDir, 0755); err != nil {
				return fmt.Errorf("failed to create output directory: %w", err)
			}
		}

		// Write to file
		if err := writer.Write(doc, convertOutput, outputFormat); err != nil {
			return fmt.Errorf("failed to write document: %w", err)
		}

		cmd.Printf("Converted %s to %s\n", inputPath, convertOutput)
	} else {
		// Write to stdout
		if err := writer.WriteTo(doc, os.Stdout, outputFormat); err != nil {
			return fmt.Errorf("failed to write document: %w", err)
		}
	}

	return nil
}

// inferOutputPath generates an output path from input path and format.
func inferOutputPath(inputPath string, format io.Format) string {
	ext := filepath.Ext(inputPath)
	base := strings.TrimSuffix(inputPath, ext)
	return base + format.Extension()
}
