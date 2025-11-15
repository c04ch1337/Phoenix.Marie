package learning

import (
	"fmt"
	"sync"
	"time"

	"github.com/phoenix-marie/core/internal/core/thought/v2/pattern"
)

// Feedback represents learning feedback data
type Feedback struct {
	PatternID string
	Score     float64
	Source    string
	Context   map[string]interface{}
	Timestamp time.Time
}

// LearningStats contains metrics about the learning process
type LearningStats struct {
	TotalPatterns   int64
	LearnedPatterns int64
	AverageScore    float64
	SuccessRate     float64
	LastUpdate      time.Time
}

// Model represents a learning model
type Model struct {
	Version    string
	Patterns   map[string]pattern.Pattern
	Weights    map[string]float64
	Parameters map[string]interface{}
	UpdatedAt  time.Time
}

// Manager handles the learning system
type Manager struct {
	model       Model
	stats       LearningStats
	feedbackLog []Feedback
	maxHistory  int
	mu          sync.RWMutex
}

// NewManager creates a new learning manager instance
func NewManager(config map[string]interface{}) *Manager {
	maxHistory := 1000
	if val, ok := config["max_history"].(int); ok {
		maxHistory = val
	}

	return &Manager{
		model: Model{
			Version:    "1.0",
			Patterns:   make(map[string]pattern.Pattern),
			Weights:    make(map[string]float64),
			Parameters: config,
			UpdatedAt:  time.Now(),
		},
		feedbackLog: make([]Feedback, 0),
		maxHistory:  maxHistory,
	}
}

// Learn processes new data for learning
func (m *Manager) Learn(data interface{}) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	startTime := time.Now()

	// Extract patterns from data
	patterns := extractPatterns(data)
	if len(patterns) == 0 {
		return fmt.Errorf("no patterns found in data")
	}

	// Process each pattern
	for _, p := range patterns {
		if err := m.learnPattern(p); err != nil {
			return fmt.Errorf("failed to learn pattern %s: %w", p.ID, err)
		}
	}

	// Update stats
	m.updateStats(patterns, startTime)

	return nil
}

// Adapt processes feedback and adapts the learning model
func (m *Manager) Adapt(feedback Feedback) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	// Validate feedback
	if err := m.validateFeedback(feedback); err != nil {
		return fmt.Errorf("invalid feedback: %w", err)
	}

	// Store feedback
	m.storeFeedback(feedback)

	// Update pattern weights based on feedback
	if err := m.updateWeights(feedback); err != nil {
		return fmt.Errorf("failed to update weights: %w", err)
	}

	// Optimize model if needed
	if m.shouldOptimize() {
		if err := m.optimize(); err != nil {
			return fmt.Errorf("optimization failed: %w", err)
		}
	}

	return nil
}

// Optimize improves the learning model
func (m *Manager) Optimize() error {
	m.mu.Lock()
	defer m.mu.Unlock()

	return m.optimize()
}

// SaveModel saves the current learning model
func (m *Manager) SaveModel(path string) error {
	m.mu.RLock()
	defer m.mu.RUnlock()

	// Implementation would serialize the model to the specified path
	return nil
}

// LoadModel loads a learning model from storage
func (m *Manager) LoadModel(path string) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	// Implementation would deserialize the model from the specified path
	return nil
}

// GetProgress returns the learning progress
func (m *Manager) GetProgress() float64 {
	m.mu.RLock()
	defer m.mu.RUnlock()

	if m.stats.TotalPatterns == 0 {
		return 0
	}

	return float64(m.stats.LearnedPatterns) / float64(m.stats.TotalPatterns)
}

// GetStats returns the current learning statistics
func (m *Manager) GetStats() LearningStats {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return m.stats
}

// Helper methods

func (m *Manager) learnPattern(p pattern.Pattern) error {
	// Store pattern in model
	m.model.Patterns[p.ID] = p

	// Initialize or update pattern weight
	if _, exists := m.model.Weights[p.ID]; !exists {
		m.model.Weights[p.ID] = calculateInitialWeight(p)
	}

	return nil
}

func (m *Manager) validateFeedback(f Feedback) error {
	if f.PatternID == "" {
		return fmt.Errorf("pattern ID is required")
	}
	if f.Score < 0 || f.Score > 1 {
		return fmt.Errorf("score must be between 0 and 1")
	}
	if _, exists := m.model.Patterns[f.PatternID]; !exists {
		return fmt.Errorf("pattern not found: %s", f.PatternID)
	}
	return nil
}

func (m *Manager) storeFeedback(f Feedback) {
	m.feedbackLog = append(m.feedbackLog, f)

	// Maintain history size
	if len(m.feedbackLog) > m.maxHistory {
		m.feedbackLog = m.feedbackLog[1:]
	}
}

func (m *Manager) updateWeights(f Feedback) error {
	currentWeight := m.model.Weights[f.PatternID]

	// Apply feedback score with decay
	learningRate := 0.1 // Could be configurable
	newWeight := currentWeight + (learningRate * (f.Score - currentWeight))

	m.model.Weights[f.PatternID] = newWeight
	return nil
}

func (m *Manager) shouldOptimize() bool {
	// Implement optimization triggering logic
	// This is a simple example based on feedback count
	return len(m.feedbackLog) >= 100
}

func (m *Manager) optimize() error {
	if len(m.feedbackLog) == 0 {
		return nil
	}

	// Analyze feedback patterns
	patterns := analyzePatterns(m.feedbackLog)

	// Update model parameters
	for patternID, score := range patterns {
		if weight, exists := m.model.Weights[patternID]; exists {
			m.model.Weights[patternID] = (weight + score) / 2
		}
	}

	// Update model timestamp
	m.model.UpdatedAt = time.Now()

	return nil
}

func (m *Manager) updateStats(patterns []pattern.Pattern, startTime time.Time) {
	m.stats.TotalPatterns += int64(len(patterns))
	m.stats.LearnedPatterns += int64(len(patterns))
	m.stats.LastUpdate = time.Now()

	// Calculate success rate
	if m.stats.TotalPatterns > 0 {
		m.stats.SuccessRate = float64(m.stats.LearnedPatterns) / float64(m.stats.TotalPatterns)
	}
}

// Utility functions

func extractPatterns(data interface{}) []pattern.Pattern {
	// Implementation would depend on the data type
	// This is a placeholder
	return nil
}

func calculateInitialWeight(p pattern.Pattern) float64 {
	// Implementation would depend on pattern characteristics
	// This is a simple placeholder
	return 0.5
}

func analyzePatterns(feedback []Feedback) map[string]float64 {
	patterns := make(map[string]float64)
	counts := make(map[string]int)

	for _, f := range feedback {
		patterns[f.PatternID] += f.Score
		counts[f.PatternID]++
	}

	// Calculate averages
	for id, total := range patterns {
		if count := counts[id]; count > 0 {
			patterns[id] = total / float64(count)
		}
	}

	return patterns
}
