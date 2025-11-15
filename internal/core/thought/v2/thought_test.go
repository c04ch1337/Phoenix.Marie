package thought

import (
	"testing"
	"time"

	"github.com/phoenix-marie/core/internal/core/memory/v2/store"
	"github.com/phoenix-marie/core/internal/core/thought/v2/dream"
	"github.com/phoenix-marie/core/internal/core/thought/v2/integration"
	"github.com/phoenix-marie/core/internal/core/thought/v2/learning"
	"github.com/phoenix-marie/core/internal/core/thought/v2/pattern"
)

func TestThoughtEngineIntegration(t *testing.T) {
	// Setup test components
	storage := setupTestStorage(t)
	patternMgr := setupPatternManager(t)
	learningMgr := setupLearningManager(t)
	dreamProc := setupDreamProcessor(t)
	bridge := setupMemoryBridge(t, storage, patternMgr, learningMgr)

	// Test cases
	t.Run("PatternManagement", func(t *testing.T) {
		// Create test pattern
		testPattern := pattern.Pattern{
			ID:         "test_pattern_1",
			Type:       "test",
			Data:       map[string]interface{}{"value": 42},
			Confidence: 0.8,
			Timestamp:  time.Now(),
			References: []string{},
			Metadata:   map[string]interface{}{},
		}

		// Test pattern detection
		patterns := patternMgr.DetectPatterns(testPattern.Data)
		if len(patterns) == 0 {
			t.Error("No patterns detected")
		}

		// Test pattern registration
		err := patternMgr.RegisterPattern(testPattern)
		if err != nil {
			t.Errorf("Pattern registration failed: %v", err)
		}

		// Test pattern analysis
		analysis := patternMgr.AnalyzePatterns()
		if len(analysis.Patterns) == 0 {
			t.Error("Pattern analysis returned no patterns")
		}

		// Test pattern confidence
		confidence := patternMgr.GetConfidence(testPattern)
		if confidence != testPattern.Confidence {
			t.Errorf("Confidence mismatch: got %f, want %f", confidence, testPattern.Confidence)
		}
	})

	t.Run("LearningSystem", func(t *testing.T) {
		// Test learning from data
		testData := map[string]interface{}{
			"input": "test data",
			"type":  "learning",
		}

		err := learningMgr.Learn(testData)
		if err != nil {
			t.Errorf("Learning failed: %v", err)
		}

		// Test feedback processing
		feedback := learning.Feedback{
			PatternID: "test_pattern_1",
			Score:     0.9,
			Source:    "test",
			Timestamp: time.Now(),
		}

		err = learningMgr.Adapt(feedback)
		if err != nil {
			t.Errorf("Feedback adaptation failed: %v", err)
		}

		// Test learning progress
		progress := learningMgr.GetProgress()
		if progress < 0 || progress > 1 {
			t.Errorf("Invalid learning progress: %f", progress)
		}

		// Test learning stats
		stats := learningMgr.GetStats()
		if stats.TotalPatterns == 0 {
			t.Error("No patterns recorded in learning stats")
		}
	})

	t.Run("DreamProcessing", func(t *testing.T) {
		// Configure dream processor
		config := dream.DreamConfig{
			MaxDuration:    time.Second * 2,
			MinConfidence:  0.5,
			BatchSize:      10,
			EnableLearning: true,
		}

		err := dreamProc.Configure(config)
		if err != nil {
			t.Errorf("Dream processor configuration failed: %v", err)
		}

		// Start dream processing
		err = dreamProc.Start()
		if err != nil {
			t.Errorf("Dream processor start failed: %v", err)
		}

		// Create test context
		context := dream.Context{
			Patterns: []pattern.Pattern{
				{
					ID:         "dream_test_1",
					Type:       "dream",
					Confidence: 0.7,
					Timestamp:  time.Now(),
				},
			},
			State: map[string]interface{}{
				"mode": "test",
			},
			StartTime: time.Now(),
		}

		// Process dream
		result := dreamProc.ProcessDream(context)
		if len(result.Insights) == 0 {
			t.Error("No insights generated from dream processing")
		}

		// Stop dream processing
		err = dreamProc.Stop()
		if err != nil {
			t.Errorf("Dream processor stop failed: %v", err)
		}
	})

	t.Run("MemoryIntegration", func(t *testing.T) {
		// Test pattern storage
		testPattern := pattern.Pattern{
			ID:         "integration_test_1",
			Type:       "integration",
			Data:       map[string]interface{}{"test": true},
			Confidence: 0.85,
			Timestamp:  time.Now(),
		}

		err := bridge.StorePattern(testPattern)
		if err != nil {
			t.Errorf("Pattern storage failed: %v", err)
		}

		// Test pattern retrieval
		retrieved, err := bridge.RetrievePattern(testPattern.ID)
		if err != nil {
			t.Errorf("Pattern retrieval failed: %v", err)
		}
		if retrieved.ID != testPattern.ID {
			t.Error("Retrieved pattern does not match stored pattern")
		}

		// Test learning state storage
		err = bridge.StoreLearningState()
		if err != nil {
			t.Errorf("Learning state storage failed: %v", err)
		}

		// Test pattern synchronization
		err = bridge.SyncPatterns()
		if err != nil {
			t.Errorf("Pattern synchronization failed: %v", err)
		}
	})

	t.Run("ErrorHandling", func(t *testing.T) {
		// Test invalid pattern registration
		err := patternMgr.RegisterPattern(pattern.Pattern{})
		if err == nil {
			t.Error("Expected error for invalid pattern registration")
		}

		// Test invalid learning input
		err = learningMgr.Learn(nil)
		if err == nil {
			t.Error("Expected error for nil learning input")
		}

		// Test invalid dream configuration
		err = dreamProc.Configure(dream.DreamConfig{MaxDuration: -1})
		if err == nil {
			t.Error("Expected error for invalid dream configuration")
		}

		// Test invalid pattern retrieval
		_, err = bridge.RetrievePattern("non_existent")
		if err == nil {
			t.Error("Expected error for non-existent pattern retrieval")
		}
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

func setupPatternManager(t *testing.T) *pattern.Manager {
	return pattern.NewManager()
}

func setupLearningManager(t *testing.T) *learning.Manager {
	config := map[string]interface{}{
		"max_history": 100,
	}
	return learning.NewManager(config)
}

func setupDreamProcessor(t *testing.T) *dream.Processor {
	config := dream.DreamConfig{
		MaxDuration:    time.Second * 5,
		MinConfidence:  0.5,
		BatchSize:      10,
		EnableLearning: true,
	}
	return dream.NewProcessor(config)
}

func setupMemoryBridge(t *testing.T, storage store.StorageEngine, patterns *pattern.Manager, learning *learning.Manager) *integration.MemoryBridge {
	config := integration.BridgeConfig{
		CacheTTL:      time.Minute,
		BatchSize:     10,
		SyncInterval:  time.Second * 5,
		RetryAttempts: 3,
	}
	return integration.NewMemoryBridge(storage, nil, patterns, learning, config)
}
