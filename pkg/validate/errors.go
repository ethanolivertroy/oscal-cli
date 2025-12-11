package validate

import (
	"fmt"
	"strings"
)

// ValidationError represents a single validation error.
type ValidationError struct {
	Path       string // JSON pointer to the error location (e.g., "/catalog/metadata/title")
	Message    string // Human-readable error message
	SchemaPath string // Path in the schema that caused the error
}

// Error implements the error interface.
func (e *ValidationError) Error() string {
	if e.Path == "" || e.Path == "/" {
		return e.Message
	}
	return fmt.Sprintf("%s: %s", e.Path, e.Message)
}

// ValidationResult holds the complete validation outcome.
type ValidationResult struct {
	Valid        bool
	Errors       []ValidationError
	DocumentType string
	FilePath     string
	SchemaVersion string
}

// ErrorCount returns the number of validation errors.
func (r *ValidationResult) ErrorCount() int {
	return len(r.Errors)
}

// String formats the validation result for display.
func (r *ValidationResult) String() string {
	if r.Valid {
		return fmt.Sprintf("Valid %s document: %s", r.DocumentType, r.FilePath)
	}

	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("Document %s is invalid (%d error", r.FilePath, len(r.Errors)))
	if len(r.Errors) != 1 {
		sb.WriteString("s")
	}
	sb.WriteString("):\n")

	for i, err := range r.Errors {
		sb.WriteString(fmt.Sprintf("  %d. %s\n", i+1, err.Error()))
	}

	return sb.String()
}
