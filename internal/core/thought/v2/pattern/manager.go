package pattern

import (
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"sort"
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
func (m *Manager) DetectPatterns(input interface{}) ([]Pattern, error) {
	if err := validateInput(input); err != nil {
		return nil, fmt.Errorf("invalid input: %w", err)
	}

	m.mu.Lock()
	defer m.mu.Unlock()

	var detected []Pattern

	// Extract metadata if available
	metadata, err := extractMetadata(input)
	if err != nil {
		return nil, fmt.Errorf("metadata extraction failed: %w", err)
	}

	// Create new pattern from input
	patternID, err := generatePatternID()
	if err != nil {
		return nil, fmt.Errorf("pattern ID generation failed: %w", err)
	}

	pattern := Pattern{
		ID:         patternID,
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

	return detected, nil
}

// RegisterPattern registers a new pattern
func (m *Manager) RegisterPattern(pattern Pattern) error {
	if err := validatePattern(pattern); err != nil {
		return fmt.Errorf("invalid pattern: %w", err)
	}

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
	if err := validatePattern(pattern); err != nil {
		return fmt.Errorf("invalid pattern: %w", err)
	}

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

	patterns := make([]Pattern, 0, len(m.patterns))
	for _, p := range m.patterns {
		patterns = append(patterns, p)
	}

	sortPatternsByConfidence(patterns)

	if len(patterns) < n {
		return patterns
	}
	return patterns[:n]
}

// Utility functions

func generatePatternID() (string, error) {
	// Generate 16 random bytes for the ID
	bytes := make([]byte, 16)
	if _, err := rand.Read(bytes); err != nil {
		return "", fmt.Errorf("failed to generate random bytes: %w", err)
	}

	// Encode as URL-safe base64
	return base64.URLEncoding.EncodeToString(bytes), nil
}

func validateInput(input interface{}) error {
	if input == nil {
		return fmt.Errorf("input cannot be nil")
	}

	// Validate input can be marshaled to JSON
	if _, err := json.Marshal(input); err != nil {
		return fmt.Errorf("input must be JSON serializable: %w", err)
	}

	return nil
}

func validatePattern(p Pattern) error {
	if p.ID == "" {
		return fmt.Errorf("pattern ID cannot be empty")
	}

	if p.Type == "" {
		return fmt.Errorf("pattern type cannot be empty")
	}

	if p.Data == nil {
		return fmt.Errorf("pattern data cannot be nil")
	}

	if p.Confidence < 0 || p.Confidence > 1 {
		return fmt.Errorf("confidence must be between 0 and 1")
	}

	if p.Timestamp.IsZero() {
		return fmt.Errorf("timestamp cannot be zero")
	}

	return nil
}

func determinePatternType(input interface{}) string {
	return fmt.Sprintf("%T", input)
}

func calculateInitialConfidence(input interface{}) float64 {
	// Implementation would depend on the input type
	// This is a simple placeholder
	return 0.5
}

func extractMetadata(input interface{}) (map[string]interface{}, error) {
	metadata := make(map[string]interface{})

	// Convert input to JSON for metadata extraction
	data, err := json.Marshal(input)
	if err != nil {
		return nil, err
	}

	// Extract basic metadata
	var rawMap map[string]interface{}
	if err := json.Unmarshal(data, &rawMap); err == nil {
		if len(rawMap) > 0 {
			metadata["fields"] = len(rawMap)
			metadata["size"] = len(data)
		}
	}

	return metadata, nil
}

func isSimilar(p1, p2 Pattern) bool {
	if p1.Type != p2.Type {
		return false
	}

	// Compare metadata for similarity
	if len(p1.Metadata) > 0 && len(p2.Metadata) > 0 {
		matches := 0
		for k, v1 := range p1.Metadata {
			if v2, exists := p2.Metadata[k]; exists && v1 == v2 {
				matches++
			}
		}
		return float64(matches)/float64(len(p1.Metadata)) > 0.7
	}

	return false
}

func sortPatternsByConfidence(patterns []Pattern) {
	sort.Slice(patterns, func(i, j int) bool {
		return patterns[i].Confidence > patterns[j].Confidence
	})
}
