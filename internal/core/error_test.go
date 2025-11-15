package core

import (
	"fmt"
	"reflect"
	"testing"
	"time"

	"github.com/phoenix-marie/core/internal/core/memory/v2/store"
	"github.com/phoenix-marie/core/internal/core/memory/v2/validation"
	"github.com/phoenix-marie/core/internal/core/thought/v2/dream"
	"github.com/phoenix-marie/core/internal/core/thought/v2/integration"
	"github.com/phoenix-marie/core/internal/core/thought/v2/learning"
	"github.com/phoenix-marie/core/internal/core/thought/v2/pattern"
)

func TestErrorHandling(t *testing.T) {
	// Setup test components
	storage := setupTestStorage(t)
	validator := setupTestValidator(t)
	patterns := setupTestPatterns(t)
	learning := setupTestLearning(t)
	dreaming := setupTestDreaming(t)
	bridge := setupTestBridge(t, storage, patterns, learning)

	t.Run("StorageErrors", func(t *testing.T) {
		// Test invalid layer
		_, err := storage.Retrieve("invalid_layer", "key")
		if err == nil {
			t.Error("Expected error for invalid layer")
		}

		// Test invalid key
		_, err = storage.Retrieve("test_layer", "")
		if err == nil {
			t.Error("Expected error for empty key")
		}

		// Test transaction errors
		tx, _ := storage.BeginTx()
		tx.Commit()       // Commit first time
		err = tx.Commit() // Try to commit again
		if err == nil {
			t.Error("Expected error for double commit")
		}

		// Test batch operation errors
		err = storage.BatchStore([]store.StoreOperation{
			{Layer: "", Key: "key", Value: "value"}, // Invalid layer
		})
		if err == nil {
			t.Error("Expected error for invalid batch operation")
		}
	})

	t.Run("ValidationErrors", func(t *testing.T) {
		// Test invalid schema
		err := validator.RegisterSchema("test", validation.Schema{})
		if err == nil {
			t.Error("Expected error for empty schema")
		}

		// Test schema validation errors
		testSchema := validation.Schema{
			Fields: map[string]validation.FieldDefinition{
				"age": {
					Type:     reflect.Int,
					Required: true,
					MinValue: 0,
					MaxValue: 150,
				},
			},
		}
		err = validator.RegisterSchema("person", testSchema)
		if err != nil {
			t.Fatalf("Failed to register schema: %v", err)
		}

		err = validator.ValidateData("person", map[string]interface{}{
			"age": -1, // Invalid age
		})
		if err == nil {
			t.Error("Expected validation error for invalid age")
		}

		// Test required field errors
		err = validator.ValidateData("person", map[string]interface{}{})
		if err == nil {
			t.Error("Expected error for missing required field")
		}
	})

	t.Run("PatternErrors", func(t *testing.T) {
		// Test invalid pattern registration
		err := patterns.RegisterPattern(pattern.Pattern{})
		if err == nil {
			t.Error("Expected error for empty pattern")
		}

		// Test duplicate pattern
		p := pattern.Pattern{
			ID:   "test_pattern",
			Type: "test",
		}
		patterns.RegisterPattern(p)
		err = patterns.RegisterPattern(p) // Try to register same pattern again
		if err == nil {
			t.Error("Expected error for duplicate pattern")
		}

		// Test pattern update errors
		err = patterns.UpdatePattern(pattern.Pattern{ID: "non_existent"})
		if err == nil {
			t.Error("Expected error for non-existent pattern update")
		}
	})

	t.Run("LearningErrors", func(t *testing.T) {
		// Test nil input
		err := learning.Learn(nil)
		if err == nil {
			t.Error("Expected error for nil learning input")
		}

		// Test invalid feedback
		err = learning.Learn(map[string]interface{}{
			"pattern_id": "test",
			"score":      2.0, // Invalid score > 1.0
			"source":     "test",
		})
		if err == nil {
			t.Error("Expected error for invalid feedback score")
		}

		// Test model errors
		err = learning.LoadModel("non_existent_model")
		if err == nil {
			t.Error("Expected error for non-existent model")
		}
	})

	t.Run("DreamProcessingErrors", func(t *testing.T) {
		// Test invalid configuration
		err := dreaming.Configure(dream.DreamConfig{
			MaxDuration: -1, // Invalid duration
		})
		if err == nil {
			t.Error("Expected error for invalid dream configuration")
		}

		// Test processing errors
		err = dreaming.Start()
		if err != nil {
			t.Fatal(err)
		}
		result := dreaming.ProcessDream(dream.Context{}) // Empty context
		if len(result.Insights) > 0 {
			t.Error("Expected no insights from empty context")
		}

		// Test double start
		err = dreaming.Start()
		if err == nil {
			t.Error("Expected error for double start")
		}
	})

	t.Run("IntegrationErrors", func(t *testing.T) {
		// Test pattern storage errors
		err := bridge.StorePattern(pattern.Pattern{})
		if err == nil {
			t.Error("Expected error for empty pattern storage")
		}

		// Test pattern retrieval errors
		_, err = bridge.RetrievePattern("non_existent")
		if err == nil {
			t.Error("Expected error for non-existent pattern retrieval")
		}

		// Test sync errors
		err = bridge.SyncPatterns()
		if err == nil {
			t.Error("Expected error for syncing empty patterns")
		}
	})

	t.Run("ConcurrencyErrors", func(t *testing.T) {
		// Test concurrent pattern operations
		done := make(chan bool)
		for i := 0; i < 10; i++ {
			go func(id int) {
				p := pattern.Pattern{
					ID:   fmt.Sprintf("concurrent_%d", id),
					Type: "test",
				}
				patterns.RegisterPattern(p)
				done <- true
			}(i)
		}

		// Wait for all goroutines
		for i := 0; i < 10; i++ {
			<-done
		}

		// Verify no data corruption
		analysis := patterns.AnalyzePatterns()
		if len(analysis.Patterns) != 10 {
			t.Errorf("Expected 10 patterns, got %d", len(analysis.Patterns))
		}
	})

	t.Run("RecoveryBehavior", func(t *testing.T) {
		// Test storage recovery
		tx, _ := storage.BeginTx()
		tx.Store("test", "key", "value")
		tx.Rollback()

		// Verify rollback worked
		_, err := storage.Retrieve("test", "key")
		if err == nil {
			t.Error("Expected error after rollback")
		}

		// Test dream processing recovery
		dreaming.Stop()
		err = dreaming.Start()
		if err != nil {
			t.Error("Expected successful restart after stop")
		}

		// Test pattern manager recovery
		patterns.GetState() // Should not panic after concurrent operations
	})
}

// Helper functions

func setupTestStorage(t *testing.T) store.StorageEngine {
	storage, err := store.NewBadgerStore(t.TempDir())
	if err != nil {
		t.Fatalf("Failed to create test storage: %v", err)
	}
	return storage
}

func setupTestValidator(t *testing.T) *validation.ValidationEngine {
	return validation.NewValidationEngine()
}

func setupTestPatterns(t *testing.T) *pattern.Manager {
	return pattern.NewManager()
}

func setupTestLearning(t *testing.T) *learning.Manager {
	config := map[string]interface{}{
		"max_history": 100,
	}
	return learning.NewManager(config)
}

func setupTestDreaming(t *testing.T) *dream.Processor {
	config := dream.DreamConfig{
		MaxDuration:    time.Second * 5,
		MinConfidence:  0.5,
		BatchSize:      10,
		EnableLearning: true,
	}
	return dream.NewProcessor(config)
}

func setupTestBridge(t *testing.T, storage store.StorageEngine, patterns *pattern.Manager, learning *learning.Manager) *integration.MemoryBridge {
	config := integration.BridgeConfig{
		CacheTTL:      time.Minute,
		BatchSize:     10,
		SyncInterval:  time.Second * 5,
		RetryAttempts: 3,
	}
	return integration.NewMemoryBridge(storage, nil, patterns, learning, config)
}
