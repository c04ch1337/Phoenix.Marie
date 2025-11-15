package pattern

import (
	"fmt"
	"sync"
	"time"
)

// Pattern represents a detected thought pattern
type Pattern struct {
	ID         string
	Type       string
	Data       interface{}
	Confidence float64
	Timestamp  time.Time
	References []string
	Metadata   map[string]interface{}
}

// PatternState represents the current state of pattern detection
type PatternState struct {
	Active        bool
	LastUpdate    time.Time
	PatternCount  int64
	ConfidenceSum float64
}

// PatternAnalysis contains analysis results of pattern detection
type PatternAnalysis struct {
	Patterns     []Pattern
	Distribution map[string]int
	AverageConf  float64
	TopPatterns  []Pattern
	Timestamp    time.Time
}

// Manager handles pattern detection and management
type Manager struct {
	patterns      map[string]Pattern
	state         PatternState
	analysisCache *PatternAnalysis
	cacheTTL      time.Duration
	mu            sync.RWMutex
}

// NewManager creates a new pattern manager instance
func NewManager() *Manager {
	return &Manager{
		patterns: make(map[string]Pattern),
		state: PatternState{
			Active:     true,
			LastUpdate: time.Now(),
		},
		cacheTTL: time.Minute * 5,
	}
}

// DetectPatterns analyzes input data for patterns
func (m *Manager) DetectPatterns(input interface{}) []Pattern {
	m.mu.Lock()
	defer m.mu.Unlock()

	var detected []Pattern

	// Extract metadata if available
	metadata := extractMetadata(input)

	// Create new pattern from input
	pattern := Pattern{
		ID:         generatePatternID(input),
		Type:       determinePatternType(input),
		Data:       input,
		Confidence: calculateInitialConfidence(input),
		Timestamp:  time.Now(),
		Metadata:   metadata,
	}

	// Check for similar patterns
	similar := m.findSimilarPatterns(pattern)
	if len(similar) > 0 {
		// Merge with existing patterns
		merged := m.mergePatterns(pattern, similar)
		detected = append(detected, merged)
	} else {
		// Store new pattern
		m.patterns[pattern.ID] = pattern
		detected = append(detected, pattern)
	}

	// Update state
	m.updateState(detected)

	return detected
}

// RegisterPattern registers a new pattern
func (m *Manager) RegisterPattern(pattern Pattern) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	if _, exists := m.patterns[pattern.ID]; exists {
		return fmt.Errorf("pattern already exists: %s", pattern.ID)
	}

	m.patterns[pattern.ID] = pattern
	m.updateState([]Pattern{pattern})

	return nil
}

// UpdatePattern updates an existing pattern
func (m *Manager) UpdatePattern(pattern Pattern) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	if _, exists := m.patterns[pattern.ID]; !exists {
		return fmt.Errorf("pattern not found: %s", pattern.ID)
	}

	m.patterns[pattern.ID] = pattern
	m.updateState([]Pattern{pattern})

	return nil
}

// AnalyzePatterns performs analysis on detected patterns
func (m *Manager) AnalyzePatterns() PatternAnalysis {
	m.mu.RLock()
	defer m.mu.RUnlock()

	// Check cache
	if m.analysisCache != nil && time.Since(m.analysisCache.Timestamp) < m.cacheTTL {
		return *m.analysisCache
	}

	// Perform new analysis
	analysis := PatternAnalysis{
		Patterns:     make([]Pattern, 0, len(m.patterns)),
		Distribution: make(map[string]int),
		Timestamp:    time.Now(),
	}

	var totalConf float64
	for _, p := range m.patterns {
		analysis.Patterns = append(analysis.Patterns, p)
		analysis.Distribution[p.Type]++
		totalConf += p.Confidence
	}

	if len(m.patterns) > 0 {
		analysis.AverageConf = totalConf / float64(len(m.patterns))
	}

	// Find top patterns
	analysis.TopPatterns = m.findTopPatterns(5)

	// Cache analysis
	m.analysisCache = &analysis

	return analysis
}

// GetConfidence returns the confidence score for a pattern
func (m *Manager) GetConfidence(pattern Pattern) float64 {
	m.mu.RLock()
	defer m.mu.RUnlock()

	if existing, exists := m.patterns[pattern.ID]; exists {
		return existing.Confidence
	}

	return 0
}

// GetState returns the current pattern detection state
func (m *Manager) GetState() PatternState {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return m.state
}

// Reset resets the pattern manager state
func (m *Manager) Reset() error {
	m.mu.Lock()
	defer m.mu.Unlock()

	m.patterns = make(map[string]Pattern)
	m.state = PatternState{
		Active:     true,
		LastUpdate: time.Now(),
	}
	m.analysisCache = nil

	return nil
}

// Helper methods

func (m *Manager) updateState(patterns []Pattern) {
	m.state.LastUpdate = time.Now()
	m.state.PatternCount += int64(len(patterns))

	for _, p := range patterns {
		m.state.ConfidenceSum += p.Confidence
	}

	// Invalidate cache
	m.analysisCache = nil
}

func (m *Manager) findSimilarPatterns(pattern Pattern) []Pattern {
	var similar []Pattern

	for _, existing := range m.patterns {
		if isSimilar(pattern, existing) {
			similar = append(similar, existing)
		}
	}

	return similar
}

func (m *Manager) mergePatterns(new Pattern, similar []Pattern) Pattern {
	merged := new

	// Combine confidence scores
	totalConf := new.Confidence
	for _, p := range similar {
		totalConf += p.Confidence
		merged.References = append(merged.References, p.ID)
	}

	merged.Confidence = totalConf / float64(len(similar)+1)

	return merged
}

func (m *Manager) findTopPatterns(n int) []Pattern {
	if len(m.patterns) == 0 {
		return nil
	}

	// Convert map to slice for sorting
	patterns := make([]Pattern, 0, len(m.patterns))
	for _, p := range m.patterns {
		patterns = append(patterns, p)
	}

	// Sort by confidence
	sortPatternsByConfidence(patterns)

	// Return top N patterns
	if len(patterns) < n {
		return patterns
	}
	return patterns[:n]
}

// Utility functions

func generatePatternID(input interface{}) string {
	// Implementation would depend on the input type
	// This is a simple placeholder
	return fmt.Sprintf("pattern_%d", time.Now().UnixNano())
}

func determinePatternType(input interface{}) string {
	// Implementation would depend on the input type
	// This is a simple placeholder
	return fmt.Sprintf("%T", input)
}

func calculateInitialConfidence(input interface{}) float64 {
	// Implementation would depend on the input type
	// This is a simple placeholder
	return 0.5
}

func extractMetadata(input interface{}) map[string]interface{} {
	// Implementation would depend on the input type
	// This is a simple placeholder
	return make(map[string]interface{})
}

func isSimilar(p1, p2 Pattern) bool {
	// Implementation would depend on pattern comparison logic
	// This is a simple placeholder
	return p1.Type == p2.Type
}

func sortPatternsByConfidence(patterns []Pattern) {
	// Implementation would use sort.Slice with confidence comparison
	// This is a placeholder
}
