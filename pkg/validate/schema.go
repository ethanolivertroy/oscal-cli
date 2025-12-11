// Package validate provides OSCAL document validation using JSON Schema.
package validate

import (
	"embed"
	"fmt"

	"github.com/ethantroy/oscal-cli/pkg/oscal/model"
)

//go:embed schemas/*.json
var schemaFS embed.FS

// SchemaVersion is the OSCAL schema version embedded in this build.
const SchemaVersion = "1.1.3"

// documentTypeToSchema maps OSCAL document types to schema file names.
var documentTypeToSchema = map[model.DocumentType]string{
	model.DocumentTypeCatalog:             "schemas/oscal_catalog_schema.json",
	model.DocumentTypeProfile:             "schemas/oscal_profile_schema.json",
	model.DocumentTypeSSP:                 "schemas/oscal_ssp_schema.json",
	model.DocumentTypeComponentDefinition: "schemas/oscal_component_schema.json",
	model.DocumentTypeAssessmentPlan:      "schemas/oscal_assessment-plan_schema.json",
	model.DocumentTypeAssessmentResults:   "schemas/oscal_assessment-results_schema.json",
	model.DocumentTypePOAM:                "schemas/oscal_poam_schema.json",
}

// LoadSchema loads the embedded JSON schema for a document type.
func LoadSchema(docType model.DocumentType) ([]byte, error) {
	schemaPath, ok := documentTypeToSchema[docType]
	if !ok {
		return nil, fmt.Errorf("no schema found for document type: %s", docType)
	}

	return schemaFS.ReadFile(schemaPath)
}

// GetSchemaURI returns the schema identifier URI for a document type.
func GetSchemaURI(docType model.DocumentType) string {
	return fmt.Sprintf("oscal://schemas/%s.json", docType)
}

// SupportedDocumentTypes returns all document types that have schemas available.
func SupportedDocumentTypes() []model.DocumentType {
	types := make([]model.DocumentType, 0, len(documentTypeToSchema))
	for docType := range documentTypeToSchema {
		types = append(types, docType)
	}
	return types
}
