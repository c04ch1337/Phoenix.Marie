package thought

import (
	"testing"
	"time"

	"github.com/phoenix-marie/core/internal/core/thought/v2/pattern"
)

func TestPatternDetection(t *testing.T) {
	patternMgr := pattern.NewManager()

	testPattern := pattern.Pattern{
		ID:         "test_pattern",
		Type:       "test",
		Data:       map[string]interface{}{"test": "data"},
		Confidence: 0.8,
		Timestamp:  time.Now(),
		Metadata:   map[string]interface{}{"source": "test"},
	}

	// Test pattern registration
	err := patternMgr.RegisterPattern(testPattern)
	if err != nil {
		t.Errorf("Failed to register pattern: %v", err)
	}

	// Test pattern detection
	patterns, err := patternMgr.DetectPatterns(testPattern.Data)
	if err != nil {
		t.Errorf("Failed to detect patterns: %v", err)
	}
	if len(patterns) == 0 {
		t.Error("No patterns detected")
	}

	// Test pattern analysis
	analysis := patternMgr.AnalyzePatterns()
	if len(analysis.Patterns) == 0 {
		t.Error("No patterns in analysis")
	}

	// Test confidence calculation
	conf := patternMgr.GetConfidence(testPattern)
	if conf <= 0 {
		t.Error("Invalid confidence score")
	}

	// Test state management
	state := patternMgr.GetState()
	if !state.Active {
		t.Error("Pattern manager not active")
	}

	// Test pattern update
	testPattern.Confidence = 0.9
	err = patternMgr.UpdatePattern(testPattern)
	if err != nil {
		t.Errorf("Failed to update pattern: %v", err)
	}

	// Test reset
	err = patternMgr.Reset()
	if err != nil {
		t.Errorf("Failed to reset pattern manager: %v", err)
	}

	state = patternMgr.GetState()
	if state.PatternCount != 0 {
		t.Error("Pattern count not reset")
	}
}

func TestPatternValidation(t *testing.T) {
	patternMgr := pattern.NewManager()

	// Test nil input
	_, err := patternMgr.DetectPatterns(nil)
	if err == nil {
		t.Error("Expected error for nil input")
	}

	// Test invalid pattern
	invalidPattern := pattern.Pattern{
		// Missing required fields
	}
	err = patternMgr.RegisterPattern(invalidPattern)
	if err == nil {
		t.Error("Expected error for invalid pattern")
	}

	// Test invalid confidence
	invalidPattern = pattern.Pattern{
		ID:         "test",
		Type:       "test",
		Data:       "test",
		Confidence: 2.0, // Invalid confidence value
		Timestamp:  time.Now(),
	}
	err = patternMgr.RegisterPattern(invalidPattern)
	if err == nil {
		t.Error("Expected error for invalid confidence")
	}
}

func TestPatternSimilarity(t *testing.T) {
	patternMgr := pattern.NewManager()

	pattern1 := pattern.Pattern{
		ID:   "test1",
		Type: "test",
		Data: map[string]interface{}{
			"field1": "value1",
			"field2": "value2",
		},
		Confidence: 0.8,
		Timestamp:  time.Now(),
		Metadata: map[string]interface{}{
			"source": "test",
			"type":   "similarity",
		},
	}

	pattern2 := pattern.Pattern{
		ID:   "test2",
		Type: "test",
		Data: map[string]interface{}{
			"field1": "value1",
			"field2": "different",
		},
		Confidence: 0.7,
		Timestamp:  time.Now(),
		Metadata: map[string]interface{}{
			"source": "test",
			"type":   "similarity",
		},
	}

	// Register first pattern
	err := patternMgr.RegisterPattern(pattern1)
	if err != nil {
		t.Errorf("Failed to register pattern1: %v", err)
	}

	// Detect similar pattern
	patterns, err := patternMgr.DetectPatterns(pattern2.Data)
	if err != nil {
		t.Errorf("Failed to detect patterns: %v", err)
	}
	if len(patterns) == 0 {
		t.Error("No similar patterns detected")
	}

	// Verify merged confidence
	if len(patterns) > 0 && patterns[0].Confidence <= pattern1.Confidence {
		t.Error("Pattern confidence not properly merged")
	}
}
