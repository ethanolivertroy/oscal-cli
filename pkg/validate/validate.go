package validate

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"github.com/ethantroy/oscal-cli/pkg/oscal/io"
	"github.com/ethantroy/oscal-cli/pkg/oscal/model"
	"github.com/santhosh-tekuri/jsonschema/v6"
	"gopkg.in/yaml.v3"
)

// Validator provides OSCAL document validation against JSON schemas.
type Validator struct {
	compiler *jsonschema.Compiler
	schemas  map[model.DocumentType]*jsonschema.Schema
}

// NewValidator creates a new Validator with pre-compiled embedded schemas.
func NewValidator() (*Validator, error) {
	v := &Validator{
		compiler: jsonschema.NewCompiler(),
		schemas:  make(map[model.DocumentType]*jsonschema.Schema),
	}

	// Pre-compile all schemas
	for docType := range documentTypeToSchema {
		if err := v.loadSchema(docType); err != nil {
			return nil, fmt.Errorf("failed to load schema for %s: %w", docType, err)
		}
	}

	return v, nil
}

// loadSchema loads and compiles a schema for a document type.
func (v *Validator) loadSchema(docType model.DocumentType) error {
	schemaData, err := LoadSchema(docType)
	if err != nil {
		return err
	}

	// Parse schema JSON
	var schemaDoc any
	if err := json.Unmarshal(schemaData, &schemaDoc); err != nil {
		return fmt.Errorf("failed to parse schema JSON: %w", err)
	}

	// Get the schema's $id for proper resolution
	schemaMap, ok := schemaDoc.(map[string]any)
	if !ok {
		return fmt.Errorf("schema is not a JSON object")
	}
	schemaID, _ := schemaMap["$id"].(string)
	if schemaID == "" {
		schemaID = GetSchemaURI(docType)
	}

	// Add schema as a resource to the compiler
	if err := v.compiler.AddResource(schemaID, schemaDoc); err != nil {
		return fmt.Errorf("failed to add schema resource: %w", err)
	}

	// Compile the schema
	schema, err := v.compiler.Compile(schemaID)
	if err != nil {
		return fmt.Errorf("failed to compile schema: %w", err)
	}

	v.schemas[docType] = schema
	return nil
}

// ValidateFile validates an OSCAL document file.
func (v *Validator) ValidateFile(path string) (*ValidationResult, error) {
	// Detect format
	format, err := io.DetectFormat(path)
	if err != nil {
		return nil, fmt.Errorf("failed to detect format: %w", err)
	}

	// Read file
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("failed to read file: %w", err)
	}

	return v.ValidateBytes(data, format, path)
}

// ValidateBytes validates OSCAL document bytes.
func (v *Validator) ValidateBytes(data []byte, format io.Format, sourcePath string) (*ValidationResult, error) {
	// Convert to JSON if needed (schema validation requires JSON)
	jsonData, err := v.toJSON(data, format)
	if err != nil {
		return nil, fmt.Errorf("failed to convert to JSON: %w", err)
	}

	// Detect document type from content
	docType, err := v.detectDocumentType(jsonData)
	if err != nil {
		return nil, fmt.Errorf("failed to detect document type: %w", err)
	}

	// Get compiled schema
	schema, ok := v.schemas[docType]
	if !ok {
		return nil, fmt.Errorf("no schema available for document type: %s", docType)
	}

	// Parse JSON for validation
	var doc any
	if err := json.Unmarshal(jsonData, &doc); err != nil {
		return nil, fmt.Errorf("failed to parse JSON: %w", err)
	}

	// Validate
	result := &ValidationResult{
		FilePath:      sourcePath,
		DocumentType:  string(docType),
		SchemaVersion: SchemaVersion,
		Valid:         true,
	}

	err = schema.Validate(doc)
	if err != nil {
		result.Valid = false
		result.Errors = v.extractErrors(err)
	}

	return result, nil
}

// toJSON converts data from any supported format to JSON.
func (v *Validator) toJSON(data []byte, format io.Format) ([]byte, error) {
	switch format {
	case io.FormatJSON:
		return data, nil

	case io.FormatYAML:
		var obj any
		if err := yaml.Unmarshal(data, &obj); err != nil {
			return nil, fmt.Errorf("failed to parse YAML: %w", err)
		}
		// Convert YAML's map[string]any to JSON-compatible format
		obj = convertYAMLToJSON(obj)
		return json.Marshal(obj)

	case io.FormatXML:
		// Load as Go struct then marshal to JSON
		loader := io.NewLoader()
		doc, err := loader.LoadFromBytes(data, format)
		if err != nil {
			return nil, fmt.Errorf("failed to load XML: %w", err)
		}
		writer := io.NewWriter()
		return writer.Marshal(doc, io.FormatJSON)

	default:
		return nil, fmt.Errorf("unsupported format: %s", format)
	}
}

// convertYAMLToJSON recursively converts YAML types to JSON-compatible types.
// YAML uses map[string]any but can also have map[any]any which JSON doesn't support.
func convertYAMLToJSON(v any) any {
	switch val := v.(type) {
	case map[string]any:
		result := make(map[string]any)
		for k, v := range val {
			result[k] = convertYAMLToJSON(v)
		}
		return result
	case map[any]any:
		result := make(map[string]any)
		for k, v := range val {
			key := fmt.Sprintf("%v", k)
			result[key] = convertYAMLToJSON(v)
		}
		return result
	case []any:
		result := make([]any, len(val))
		for i, v := range val {
			result[i] = convertYAMLToJSON(v)
		}
		return result
	default:
		return v
	}
}

// detectDocumentType determines the OSCAL document type from JSON content.
func (v *Validator) detectDocumentType(jsonData []byte) (model.DocumentType, error) {
	var wrapper map[string]any
	if err := json.Unmarshal(jsonData, &wrapper); err != nil {
		return "", fmt.Errorf("failed to parse JSON: %w", err)
	}

	// Check for known root elements
	typeMap := map[string]model.DocumentType{
		"catalog":                       model.DocumentTypeCatalog,
		"profile":                       model.DocumentTypeProfile,
		"system-security-plan":          model.DocumentTypeSSP,
		"component-definition":          model.DocumentTypeComponentDefinition,
		"assessment-plan":               model.DocumentTypeAssessmentPlan,
		"assessment-results":            model.DocumentTypeAssessmentResults,
		"plan-of-action-and-milestones": model.DocumentTypePOAM,
	}

	for key, docType := range typeMap {
		if _, ok := wrapper[key]; ok {
			return docType, nil
		}
	}

	return "", fmt.Errorf("could not determine document type from content")
}

// extractErrors converts jsonschema errors to ValidationErrors.
func (v *Validator) extractErrors(err error) []ValidationError {
	var errors []ValidationError

	if validationErr, ok := err.(*jsonschema.ValidationError); ok {
		errors = v.flattenErrors(validationErr, errors)
	} else {
		errors = append(errors, ValidationError{
			Message: err.Error(),
		})
	}

	return errors
}

// flattenErrors recursively extracts errors from the validation error tree.
func (v *Validator) flattenErrors(err *jsonschema.ValidationError, errors []ValidationError) []ValidationError {
	// If this error has causes, recurse into them
	if len(err.Causes) > 0 {
		for _, cause := range err.Causes {
			errors = v.flattenErrors(cause, errors)
		}
		return errors
	}

	// Convert instance location ([]string) to JSON pointer path
	path := "/" + strings.Join(err.InstanceLocation, "/")
	if path == "/" {
		path = ""
	}

	// Leaf error - add it
	errors = append(errors, ValidationError{
		Path:       path,
		Message:    err.Error(),
		SchemaPath: err.SchemaURL,
	})

	return errors
}
