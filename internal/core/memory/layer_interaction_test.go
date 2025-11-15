package memory

import (
	"testing"
)

func TestLayerInteraction(t *testing.T) {
	phl, err := NewPHL("./testdata")
	if err != nil {
		t.Fatalf("Failed to create PHL: %v", err)
	}
	defer phl.Close()

	// Test data propagation
	t.Run("Test data propagation", func(t *testing.T) {
		// Store test data in sensory layer
		testData := map[string]string{"test": "value"}
		if !phl.Store("sensory", "test_key", testData) {
			t.Fatal("Failed to store test data")
		}

		// Propagate data from sensory to connected layers
		err := phl.PropagateData("sensory", "test_key")
		if err != nil {
			t.Fatalf("Failed to propagate data: %v", err)
		}

		// Verify data in emotion layer
		emotionData, exists := phl.Retrieve("emotion", "test_key_from_emotion")
		if !exists {
			t.Error("Data not propagated to emotion layer")
		}
		if emotionData == nil {
			t.Error("Propagated data is nil in emotion layer")
		}

		// Verify data in logic layer
		logicData, exists := phl.Retrieve("logic", "test_key_from_logic")
		if !exists {
			t.Error("Data not propagated to logic layer")
		}
		if logicData == nil {
			t.Error("Propagated data is nil in logic layer")
		}
	})

	t.Run("Test route management", func(t *testing.T) {
		// Test adding new route
		err := phl.AddLayerRoute("sensory", "dream")
		if err != nil {
			t.Fatalf("Failed to add new route: %v", err)
		}

		// Verify new route
		routes, err := phl.GetLayerRoutes("sensory")
		if err != nil {
			t.Fatalf("Failed to get routes: %v", err)
		}

		foundNewRoute := false
		for _, route := range routes {
			if route == "dream" {
				foundNewRoute = true
				break
			}
		}
		if !foundNewRoute {
			t.Error("New route not found in routes list")
		}

		// Test invalid source layer
		err = phl.AddLayerRoute("invalid", "emotion")
		if err == nil {
			t.Error("Expected error for invalid source layer")
		}

		// Test invalid target layer
		err = phl.AddLayerRoute("sensory", "invalid")
		if err == nil {
			t.Error("Expected error for invalid target layer")
		}
	})

	t.Run("Test multi-layer propagation", func(t *testing.T) {
		// Store test data in emotion layer
		testData := map[string]string{"emotion": "happy"}
		if !phl.Store("emotion", "emotion_state", testData) {
			t.Fatal("Failed to store emotion test data")
		}

		// Propagate data from emotion layer
		err := phl.PropagateData("emotion", "emotion_state")
		if err != nil {
			t.Fatalf("Failed to propagate emotion data: %v", err)
		}

		// Verify data propagated to logic layer
		logicData, exists := phl.Retrieve("logic", "emotion_state_from_logic")
		if !exists {
			t.Error("Emotion data not propagated to logic layer")
		}

		// Verify data propagated to dream layer
		dreamData, exists := phl.Retrieve("dream", "emotion_state_from_dream")
		if !exists {
			t.Error("Emotion data not propagated to dream layer")
		}

		// Verify data consistency
		if logicData == nil || dreamData == nil {
			t.Error("Propagated data is nil in target layers")
		}
	})
}
