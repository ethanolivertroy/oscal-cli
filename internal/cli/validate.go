package cli

import (
	"fmt"
	"os"

	"github.com/ethantroy/oscal-cli/pkg/validate"
	"github.com/spf13/cobra"
)

var (
	validateQuiet bool
)

var validateCmd = &cobra.Command{
	Use:   "validate <input-file>",
	Short: "Validate OSCAL documents against JSON Schema",
	Long: `Validate OSCAL documents against the official NIST JSON schemas.

The input format (XML, JSON, YAML) is auto-detected from the file extension or content.
The document type is auto-detected from the document content.

Embedded schemas correspond to OSCAL version ` + validate.SchemaVersion + `.

Examples:
  oscal validate catalog.json
  oscal validate ssp.xml
  oscal validate profile.yaml --quiet`,
	Args: cobra.ExactArgs(1),
	RunE: runValidate,
}

func init() {
	rootCmd.AddCommand(validateCmd)

	validateCmd.Flags().BoolVarP(&validateQuiet, "quiet", "q", false,
		"Only output errors, suppress success messages")
}

func runValidate(cmd *cobra.Command, args []string) error {
	inputPath := args[0]

	// Check file exists
	if _, err := os.Stat(inputPath); os.IsNotExist(err) {
		return fmt.Errorf("file not found: %s", inputPath)
	}

	// Create validator
	validator, err := validate.NewValidator()
	if err != nil {
		return fmt.Errorf("failed to initialize validator: %w", err)
	}

	// Validate
	result, err := validator.ValidateFile(inputPath)
	if err != nil {
		return fmt.Errorf("validation failed: %w", err)
	}

	// Output results
	if result.Valid {
		if !validateQuiet {
			cmd.Printf("Valid %s document: %s\n", result.DocumentType, inputPath)
		}
		return nil
	}

	// Print errors to stderr
	cmd.PrintErrln(result.String())

	// Exit with error code
	os.Exit(1)
	return nil
}
