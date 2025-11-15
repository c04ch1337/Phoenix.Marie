package thought

import (
	"fmt"
	"math"
	"sync"
)

// Pattern represents a recognized thought pattern
type Pattern struct {
	ID         string
	Elements   []interface{}
	Confidence float64
	Frequency  int
}

// PatternManager handles pattern recognition and management
type PatternManager struct {
	patterns   map[string]*Pattern
	minConf    float64
	mu         sync.RWMutex
	totalConf  float64
	patternCnt int
}

// NewPatternManager creates a new pattern manager
func NewPatternManager(minConfidence float64) *PatternManager {
	return &PatternManager{
		patterns: make(map[string]*Pattern),
		minConf:  minConfidence,
	}
}

// ProcessInput processes new input for pattern recognition
func (pm *PatternManager) ProcessInput(input interface{}) {
	pm.mu.Lock()
	defer pm.mu.Unlock()

	// Extract elements from input
	elements := pm.extractElements(input)

	// Generate pattern ID
	patternID := pm.generatePatternID(elements)

	// Update existing pattern or create new one
	if pattern, exists := pm.patterns[patternID]; exists {
		pattern.Frequency++
		pattern.Confidence = pm.calculateConfidence(pattern)
		pm.updateTotalConfidence()
	} else {
		newPattern := &Pattern{
			ID:         patternID,
			Elements:   elements,
			Frequency:  1,
			Confidence: pm.initialConfidence(elements),
		}
		if newPattern.Confidence >= pm.minConf {
			pm.patterns[patternID] = newPattern
			pm.patternCnt++
			pm.updateTotalConfidence()
		}
	}
}

// GetPatterns returns all recognized patterns
func (pm *PatternManager) GetPatterns() []*Pattern {
	pm.mu.RLock()
	defer pm.mu.RUnlock()

	patterns := make([]*Pattern, 0, len(pm.patterns))
	for _, pattern := range pm.patterns {
		patterns = append(patterns, pattern)
	}
	return patterns
}

// GetPatternCount returns the number of recognized patterns
func (pm *PatternManager) GetPatternCount() int {
	pm.mu.RLock()
	defer pm.mu.RUnlock()
	return pm.patternCnt
}

// GetAverageConfidence returns the average confidence across all patterns
func (pm *PatternManager) GetAverageConfidence() float64 {
	pm.mu.RLock()
	defer pm.mu.RUnlock()

	if pm.patternCnt == 0 {
		return 0
	}
	return pm.totalConf / float64(pm.patternCnt)
}

// extractElements extracts pattern elements from input
func (pm *PatternManager) extractElements(input interface{}) []interface{} {
	switch v := input.(type) {
	case []interface{}:
		return v
	case map[string]interface{}:
		elements := make([]interface{}, 0, len(v))
		for _, val := range v {
			elements = append(elements, val)
		}
		return elements
	default:
		return []interface{}{input}
	}
}

// generatePatternID generates a unique ID for a pattern based on its elements
func (pm *PatternManager) generatePatternID(elements []interface{}) string {
	// Simple implementation - can be enhanced with more sophisticated hashing
	hash := ""
	for _, elem := range elements {
		hash += fmt.Sprintf("%v", elem)
	}
	return hash
}

// calculateConfidence calculates confidence for a pattern
func (pm *PatternManager) calculateConfidence(p *Pattern) float64 {
	// Base confidence on frequency and complexity
	baseConf := float64(p.Frequency) * 0.1
	complexityFactor := float64(len(p.Elements)) * 0.05
	return math.Min(baseConf+complexityFactor, 1.0)
}

// initialConfidence calculates initial confidence for a new pattern
func (pm *PatternManager) initialConfidence(elements []interface{}) float64 {
	return math.Max(0.1*float64(len(elements)), pm.minConf)
}

// updateTotalConfidence updates the total confidence sum
func (pm *PatternManager) updateTotalConfidence() {
	pm.totalConf = 0
	for _, pattern := range pm.patterns {
		pm.totalConf += pattern.Confidence
	}
}
