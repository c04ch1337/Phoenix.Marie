package memory

import (
	"encoding/json"
	"fmt"
	"reflect"
)

// ValidationError represents a validation error
type ValidationError struct {
	Field   string
	Message string
}

func (e *ValidationError) Error() string {
	return fmt.Sprintf("validation error on field '%s': %s", e.Field, e.Message)
}

// Validator defines the interface for data validators
type Validator interface {
	Validate(data any) error
}

// LayerValidator provides validation rules for each memory layer
type LayerValidator struct {
	validators map[string]Validator
}

// NewLayerValidator creates a new layer validator with specific validation rules
func NewLayerValidator() *LayerValidator {
	return &LayerValidator{
		validators: map[string]Validator{
			"sensory": &SensoryValidator{},
			"emotion": &EmotionValidator{},
			"logic":   &LogicValidator{},
			"dream":   &DreamValidator{},
			"eternal": &EternalValidator{},
		},
	}
}

// ValidateLayerData validates data for a specific layer
func (lv *LayerValidator) ValidateLayerData(layer string, data any) error {
	validator, exists := lv.validators[layer]
	if !exists {
		return fmt.Errorf("no validator found for layer: %s", layer)
	}
	return validator.Validate(data)
}

// SensoryValidator implements validation rules for sensory data
type SensoryValidator struct{}

func (v *SensoryValidator) Validate(data any) error {
	if data == nil {
		return &ValidationError{Field: "data", Message: "sensory data cannot be nil"}
	}
	return nil
}

// EmotionValidator implements validation rules for emotional data
type EmotionValidator struct{}

func (v *EmotionValidator) Validate(data any) error {
	switch val := data.(type) {
	case map[string]interface{}:
		if intensity, ok := val["intensity"].(float64); ok {
			if intensity < 0 || intensity > 1 {
				return &ValidationError{Field: "intensity", Message: "intensity must be between 0 and 1"}
			}
		}
	case string:
		if len(val) == 0 {
			return &ValidationError{Field: "emotion", Message: "emotion string cannot be empty"}
		}
	default:
		return &ValidationError{Field: "type", Message: "invalid emotion data type"}
	}
	return nil
}

// LogicValidator implements validation rules for logical data
type LogicValidator struct{}

func (v *LogicValidator) Validate(data any) error {
	if data == nil {
		return &ValidationError{Field: "data", Message: "logic data cannot be nil"}
	}

	// Validate data structure
	switch val := data.(type) {
	case map[string]interface{}:
		if len(val) == 0 {
			return &ValidationError{Field: "data", Message: "logic data map cannot be empty"}
		}
	case string:
		if len(val) == 0 {
			return &ValidationError{Field: "data", Message: "logic data string cannot be empty"}
		}
	}
	return nil
}

// DreamValidator implements validation rules for dream data
type DreamValidator struct{}

func (v *DreamValidator) Validate(data any) error {
	if data == nil {
		return &ValidationError{Field: "data", Message: "dream data cannot be nil"}
	}

	// Ensure data can be serialized
	if _, err := json.Marshal(data); err != nil {
		return &ValidationError{Field: "data", Message: "dream data must be serializable"}
	}
	return nil
}

// EternalValidator implements validation rules for eternal data
type EternalValidator struct{}

func (v *EternalValidator) Validate(data any) error {
	if data == nil {
		return &ValidationError{Field: "data", Message: "eternal data cannot be nil"}
	}

	switch val := data.(type) {
	case map[string]interface{}:
		if importance, ok := val["importance"].(float64); ok {
			if importance <= 0 {
				return &ValidationError{Field: "importance", Message: "importance must be greater than 0"}
			}
		}
	}

	// Check if data is of a supported type
	kind := reflect.TypeOf(data).Kind()
	switch kind {
	case reflect.Map, reflect.Struct, reflect.String, reflect.Slice:
		return nil
	default:
		return &ValidationError{
			Field:   "type",
			Message: fmt.Sprintf("unsupported data type for eternal storage: %v", kind),
		}
	}
}

// ValidateKey validates a key string
func ValidateKey(key string) error {
	if key == "" {
		return &ValidationError{Field: "key", Message: "key cannot be empty"}
	}
	if len(key) > 255 {
		return &ValidationError{Field: "key", Message: "key length cannot exceed 255 characters"}
	}
	return nil
}

// ValidateLayer validates a layer name
func ValidateLayer(layer string) error {
	validLayers := map[string]bool{
		"sensory": true,
		"emotion": true,
		"logic":   true,
		"dream":   true,
		"eternal": true,
	}

	if !validLayers[layer] {
		return &ValidationError{
			Field:   "layer",
			Message: fmt.Sprintf("invalid layer name: %s", layer),
		}
	}
	return nil
}
