package memory

import (
	"testing"
)

func TestProcessorManager(t *testing.T) {
	pm := NewProcessorManager()

	t.Run("Test processor registration", func(t *testing.T) {
		// Test registering nil processor
		err := pm.RegisterProcessor("test", nil)
		if err == nil {
			t.Error("Expected error when registering nil processor")
		}

		// Test registering duplicate processor
		err = pm.RegisterProcessor("sensory", &SensoryProcessor{})
		if err == nil {
			t.Error("Expected error when registering duplicate processor")
		}
	})

	t.Run("Test sensory processor", func(t *testing.T) {
		testData := "test sensory input"
		processed, err := pm.ProcessData("sensory", testData)
		if err != nil {
			t.Fatalf("Failed to process sensory data: %v", err)
		}

		result, ok := processed.(map[string]interface{})
		if !ok {
			t.Fatal("Processed data is not a map")
		}

		if result["type"] != "sensory" {
			t.Errorf("Expected type 'sensory', got %v", result["type"])
		}
		if result["raw"] != testData {
			t.Errorf("Expected raw data '%v', got %v", testData, result["raw"])
		}
	})

	t.Run("Test emotion processor", func(t *testing.T) {
		testData := map[string]interface{}{
			"emotion":   "joy",
			"intensity": 0.8,
		}
		processed, err := pm.ProcessData("emotion", testData)
		if err != nil {
			t.Fatalf("Failed to process emotion data: %v", err)
		}

		result, ok := processed.(map[string]interface{})
		if !ok {
			t.Fatal("Processed data is not a map")
		}

		if result["type"] != "emotion" {
			t.Errorf("Expected type 'emotion', got %v", result["type"])
		}
		if result["processed_at"].(int64) <= 0 {
			t.Error("Invalid processed_at timestamp")
		}
	})

	t.Run("Test logic processor", func(t *testing.T) {
		testData := map[string]interface{}{
			"condition": true,
			"value":     42,
		}
		processed, err := pm.ProcessData("logic", testData)
		if err != nil {
			t.Fatalf("Failed to process logic data: %v", err)
		}

		result, ok := processed.(map[string]interface{})
		if !ok {
			t.Fatal("Processed data is not a map")
		}

		if result["type"] != "logic" {
			t.Errorf("Expected type 'logic', got %v", result["type"])
		}
		if result["complexity"].(int) != 2 {
			t.Errorf("Expected complexity 2, got %v", result["complexity"])
		}
	})

	t.Run("Test dream processor", func(t *testing.T) {
		testData := map[string]interface{}{
			"pattern": "test pattern",
			"source":  "memory",
		}
		processed, err := pm.ProcessData("dream", testData)
		if err != nil {
			t.Fatalf("Failed to process dream data: %v", err)
		}

		result, ok := processed.(map[string]interface{})
		if !ok {
			t.Fatal("Processed data is not a map")
		}

		if result["type"] != "dream" {
			t.Errorf("Expected type 'dream', got %v", result["type"])
		}

		patterns, ok := result["patterns"].([]string)
		if !ok {
			t.Fatal("Patterns is not a string slice")
		}
		if len(patterns) != 2 {
			t.Errorf("Expected 2 patterns, got %d", len(patterns))
		}
	})

	t.Run("Test eternal processor", func(t *testing.T) {
		testData := map[string]interface{}{
			"data":       "important memory",
			"importance": 5,
		}
		processed, err := pm.ProcessData("eternal", testData)
		if err != nil {
			t.Fatalf("Failed to process eternal data: %v", err)
		}

		result, ok := processed.(map[string]interface{})
		if !ok {
			t.Fatal("Processed data is not a map")
		}

		if result["type"] != "eternal" {
			t.Errorf("Expected type 'eternal', got %v", result["type"])
		}
		if !result["permanent"].(bool) {
			t.Error("Expected permanent to be true")
		}
		if result["importance"] != 5 {
			t.Errorf("Expected importance 5, got %v", result["importance"])
		}
	})

	t.Run("Test invalid layer", func(t *testing.T) {
		_, err := pm.ProcessData("invalid", "test")
		if err == nil {
			t.Error("Expected error for invalid layer")
		}
	})
}
