package validation

import (
	"fmt"
	"reflect"
	"time"
)

// Schema defines the structure and rules for data validation
type Schema struct {
	Fields    map[string]FieldDefinition
	Required  []string
	Version   string
	UpdatedAt time.Time
}

// FieldDefinition defines the validation rules for a single field
type FieldDefinition struct {
	Type      reflect.Kind
	Required  bool
	MinValue  interface{}
	MaxValue  interface{}
	Pattern   string
	Validator func(interface{}) error
}

// ValidationError represents a validation failure
type ValidationError struct {
	Field   string
	Message string
	Value   interface{}
	Rule    string
}

// ValidationEngine implements the ValidationEngine interface
type ValidationEngine struct {
	schemas map[string]Schema
	errors  []ValidationError
}

// NewValidationEngine creates a new validation engine instance
func NewValidationEngine() *ValidationEngine {
	return &ValidationEngine{
		schemas: make(map[string]Schema),
		errors:  make([]ValidationError, 0),
	}
}

// ValidateData validates data against the registered schema for a layer
func (ve *ValidationEngine) ValidateData(layer string, data interface{}) error {
	schema, exists := ve.schemas[layer]
	if !exists {
		return fmt.Errorf("no schema registered for layer: %s", layer)
	}

	ve.clearErrors()
	value := reflect.ValueOf(data)

	// Handle nil data
	if !value.IsValid() {
		ve.addError("", "data is nil", nil, "required")
		return fmt.Errorf("data validation failed")
	}

	// Handle pointer types
	if value.Kind() == reflect.Ptr {
		value = value.Elem()
	}

	// Validate against schema
	if value.Kind() == reflect.Struct {
		for fieldName, fieldDef := range schema.Fields {
			field := value.FieldByName(fieldName)
			if !field.IsValid() {
				if fieldDef.Required {
					ve.addError(fieldName, "required field missing", nil, "required")
				}
				continue
			}

			if err := ve.validateField(fieldName, field.Interface(), fieldDef); err != nil {
				ve.addError(fieldName, err.Error(), field.Interface(), "validation")
			}
		}
	} else {
		return fmt.Errorf("data must be a struct type, got %v", value.Kind())
	}

	if len(ve.errors) > 0 {
		return fmt.Errorf("validation failed with %d errors", len(ve.errors))
	}

	return nil
}

// ValidateSchema validates the schema definition itself
func (ve *ValidationEngine) ValidateSchema(layer string, schema Schema) error {
	if layer == "" {
		return fmt.Errorf("layer name cannot be empty")
	}

	if schema.Fields == nil {
		return fmt.Errorf("schema must define fields")
	}

	for fieldName, field := range schema.Fields {
		if fieldName == "" {
			return fmt.Errorf("field name cannot be empty")
		}

		if err := ve.validateFieldDefinition(field); err != nil {
			return fmt.Errorf("invalid field definition for %s: %w", fieldName, err)
		}
	}

	return nil
}

// RegisterSchema registers a new schema for a layer
func (ve *ValidationEngine) RegisterSchema(layer string, schema Schema) error {
	if err := ve.ValidateSchema(layer, schema); err != nil {
		return fmt.Errorf("invalid schema: %w", err)
	}

	ve.schemas[layer] = schema
	return nil
}

// UpdateSchema updates an existing schema
func (ve *ValidationEngine) UpdateSchema(layer string, schema Schema) error {
	if _, exists := ve.schemas[layer]; !exists {
		return fmt.Errorf("schema not found for layer: %s", layer)
	}

	if err := ve.ValidateSchema(layer, schema); err != nil {
		return fmt.Errorf("invalid schema: %w", err)
	}

	ve.schemas[layer] = schema
	return nil
}

// GetValidationErrors returns all validation errors
func (ve *ValidationEngine) GetValidationErrors() []ValidationError {
	return ve.errors
}

// ClearErrors clears all validation errors
func (ve *ValidationEngine) ClearErrors() error {
	ve.clearErrors()
	return nil
}

// Helper methods

func (ve *ValidationEngine) validateField(name string, value interface{}, def FieldDefinition) error {
	if value == nil && def.Required {
		return fmt.Errorf("required field is nil")
	}

	if value == nil {
		return nil
	}

	val := reflect.ValueOf(value)
	if val.Kind() != def.Type {
		return fmt.Errorf("invalid type: expected %v, got %v", def.Type, val.Kind())
	}

	if def.Validator != nil {
		if err := def.Validator(value); err != nil {
			return fmt.Errorf("custom validation failed: %w", err)
		}
	}

	return ve.validateRange(value, def)
}

func (ve *ValidationEngine) validateRange(value interface{}, def FieldDefinition) error {
	if def.MinValue != nil || def.MaxValue != nil {
		switch value.(type) {
		case int, int32, int64:
			val := reflect.ValueOf(value).Int()
			if def.MinValue != nil && val < reflect.ValueOf(def.MinValue).Int() {
				return fmt.Errorf("value below minimum")
			}
			if def.MaxValue != nil && val > reflect.ValueOf(def.MaxValue).Int() {
				return fmt.Errorf("value above maximum")
			}
		case float32, float64:
			val := reflect.ValueOf(value).Float()
			if def.MinValue != nil && val < reflect.ValueOf(def.MinValue).Float() {
				return fmt.Errorf("value below minimum")
			}
			if def.MaxValue != nil && val > reflect.ValueOf(def.MaxValue).Float() {
				return fmt.Errorf("value above maximum")
			}
		}
	}
	return nil
}

func (ve *ValidationEngine) validateFieldDefinition(def FieldDefinition) error {
	switch def.Type {
	case reflect.Bool, reflect.Int, reflect.Int32, reflect.Int64,
		reflect.Float32, reflect.Float64, reflect.String:
		return nil
	default:
		return fmt.Errorf("unsupported field type: %v", def.Type)
	}
}

func (ve *ValidationEngine) addError(field, message string, value interface{}, rule string) {
	ve.errors = append(ve.errors, ValidationError{
		Field:   field,
		Message: message,
		Value:   value,
		Rule:    rule,
	})
}

func (ve *ValidationEngine) clearErrors() {
	ve.errors = make([]ValidationError, 0)
}
