package memory

import (
	"testing"
)

func TestValidation(t *testing.T) {
	validator := NewLayerValidator()

	t.Run("Test key validation", func(t *testing.T) {
		// Test empty key
		if err := ValidateKey(""); err == nil {
			t.Error("Expected error for empty key")
		}

		// Test long key
		longKey := make([]byte, 256)
		for i := range longKey {
			longKey[i] = 'a'
		}
		if err := ValidateKey(string(longKey)); err == nil {
			t.Error("Expected error for key longer than 255 characters")
		}

		// Test valid key
		if err := ValidateKey("valid_key"); err != nil {
			t.Errorf("Unexpected error for valid key: %v", err)
		}
	})

	t.Run("Test layer validation", func(t *testing.T) {
		// Test invalid layer
		if err := ValidateLayer("invalid"); err == nil {
			t.Error("Expected error for invalid layer")
		}

		// Test valid layers
		validLayers := []string{"sensory", "emotion", "logic", "dream", "eternal"}
		for _, layer := range validLayers {
			if err := ValidateLayer(layer); err != nil {
				t.Errorf("Unexpected error for valid layer %s: %v", layer, err)
			}
		}
	})

	t.Run("Test sensory validation", func(t *testing.T) {
		// Test nil data
		if err := validator.ValidateLayerData("sensory", nil); err == nil {
			t.Error("Expected error for nil sensory data")
		}

		// Test valid data
		validData := map[string]interface{}{
			"input": "test data",
		}
		if err := validator.ValidateLayerData("sensory", validData); err != nil {
			t.Errorf("Unexpected error for valid sensory data: %v", err)
		}
	})

	t.Run("Test emotion validation", func(t *testing.T) {
		// Test invalid intensity
		invalidData := map[string]interface{}{
			"intensity": 1.5,
		}
		if err := validator.ValidateLayerData("emotion", invalidData); err == nil {
			t.Error("Expected error for invalid emotion intensity")
		}

		// Test valid intensity
		validData := map[string]interface{}{
			"intensity": 0.5,
		}
		if err := validator.ValidateLayerData("emotion", validData); err != nil {
			t.Errorf("Unexpected error for valid emotion data: %v", err)
		}

		// Test empty string
		if err := validator.ValidateLayerData("emotion", ""); err == nil {
			t.Error("Expected error for empty emotion string")
		}
	})

	t.Run("Test logic validation", func(t *testing.T) {
		// Test nil data
		if err := validator.ValidateLayerData("logic", nil); err == nil {
			t.Error("Expected error for nil logic data")
		}

		// Test empty map
		if err := validator.ValidateLayerData("logic", map[string]interface{}{}); err == nil {
			t.Error("Expected error for empty logic map")
		}

		// Test valid data
		validData := map[string]interface{}{
			"condition": true,
		}
		if err := validator.ValidateLayerData("logic", validData); err != nil {
			t.Errorf("Unexpected error for valid logic data: %v", err)
		}
	})

	t.Run("Test dream validation", func(t *testing.T) {
		// Test nil data
		if err := validator.ValidateLayerData("dream", nil); err == nil {
			t.Error("Expected error for nil dream data")
		}

		// Test valid data
		validData := map[string]interface{}{
			"pattern": "test pattern",
		}
		if err := validator.ValidateLayerData("dream", validData); err != nil {
			t.Errorf("Unexpected error for valid dream data: %v", err)
		}
	})

	t.Run("Test eternal validation", func(t *testing.T) {
		// Test nil data
		if err := validator.ValidateLayerData("eternal", nil); err == nil {
			t.Error("Expected error for nil eternal data")
		}

		// Test invalid importance
		invalidData := map[string]interface{}{
			"importance": 0.0,
		}
		if err := validator.ValidateLayerData("eternal", invalidData); err == nil {
			t.Error("Expected error for invalid importance")
		}

		// Test valid data
		validData := map[string]interface{}{
			"importance": 1.0,
			"data":       "test data",
		}
		if err := validator.ValidateLayerData("eternal", validData); err != nil {
			t.Errorf("Unexpected error for valid eternal data: %v", err)
		}

		// Test unsupported type
		if err := validator.ValidateLayerData("eternal", 123); err == nil {
			t.Error("Expected error for unsupported eternal data type")
		}
	})

	t.Run("Test invalid layer validation", func(t *testing.T) {
		if err := validator.ValidateLayerData("invalid", "data"); err == nil {
			t.Error("Expected error for invalid layer")
		}
	})
}
