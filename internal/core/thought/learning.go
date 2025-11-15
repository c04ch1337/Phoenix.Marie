package thought

import (
	"fmt"
	"math"
	"sync"
	"time"
)

// LearningManager handles the learning and adaptation capabilities
type LearningManager struct {
	learningRate float64
	insights     []interface{}
	models       map[string]*LearningModel
	progress     float64
	mu           sync.RWMutex
}

// LearningModel represents a specific learning model
type LearningModel struct {
	Name       string
	Weights    map[string]float64
	Bias       float64
	LastUpdate time.Time
}

// NewLearningManager creates a new learning manager
func NewLearningManager(learningRate float64) *LearningManager {
	return &LearningManager{
		learningRate: learningRate,
		insights:     make([]interface{}, 0),
		models:       make(map[string]*LearningModel),
		progress:     0.0,
	}
}

// Update updates learning models based on observed patterns
func (lm *LearningManager) Update(patterns []*Pattern) {
	lm.mu.Lock()
	defer lm.mu.Unlock()

	if len(patterns) == 0 {
		return
	}

	// Process each pattern
	for _, pattern := range patterns {
		model := lm.getOrCreateModel(pattern.ID)
		lm.updateModel(model, pattern)
	}

	// Generate insights from updated models
	lm.generateInsights(patterns)

	// Update learning progress
	lm.updateProgress()
}

// GetProgress returns the current learning progress
func (lm *LearningManager) GetProgress() float64 {
	lm.mu.RLock()
	defer lm.mu.RUnlock()
	return lm.progress
}

// GetInsights returns current learning insights
func (lm *LearningManager) GetInsights() []interface{} {
	lm.mu.RLock()
	defer lm.mu.RUnlock()

	// Return a copy of insights
	insights := make([]interface{}, len(lm.insights))
	copy(insights, lm.insights)
	return insights
}

// getOrCreateModel gets an existing model or creates a new one
func (lm *LearningManager) getOrCreateModel(id string) *LearningModel {
	if model, exists := lm.models[id]; exists {
		return model
	}

	model := &LearningModel{
		Name:       id,
		Weights:    make(map[string]float64),
		Bias:       0.0,
		LastUpdate: time.Now(),
	}
	lm.models[id] = model
	return model
}

// updateModel updates a learning model based on pattern data
func (lm *LearningManager) updateModel(model *LearningModel, pattern *Pattern) {
	// Update weights based on pattern elements
	for i := range pattern.Elements {
		key := fmt.Sprintf("weight_%d", i)
		currentWeight := model.Weights[key]

		// Apply learning rate to weight update
		delta := lm.learningRate * pattern.Confidence
		model.Weights[key] = currentWeight + delta
	}

	// Update bias
	model.Bias += lm.learningRate * (pattern.Confidence - model.Bias)
	model.LastUpdate = time.Now()
}

// generateInsights generates insights from learning models
func (lm *LearningManager) generateInsights(patterns []*Pattern) {
	insights := make([]interface{}, 0)

	// Generate insights based on pattern correlations
	for _, pattern := range patterns {
		model := lm.models[pattern.ID]
		if model == nil {
			continue
		}

		// Create insight if pattern shows significant learning
		if pattern.Confidence > 0.7 && len(model.Weights) > 0 {
			insight := map[string]interface{}{
				"pattern_id":   pattern.ID,
				"confidence":   pattern.Confidence,
				"weight_count": len(model.Weights),
				"bias":         model.Bias,
				"timestamp":    time.Now().UnixNano(),
			}
			insights = append(insights, insight)
		}
	}

	lm.insights = insights
}

// updateProgress updates the overall learning progress
func (lm *LearningManager) updateProgress() {
	if len(lm.models) == 0 {
		lm.progress = 0.0
		return
	}

	var totalProgress float64
	for _, model := range lm.models {
		// Calculate model progress based on weights and bias
		weightSum := 0.0
		for _, weight := range model.Weights {
			weightSum += math.Abs(weight)
		}

		modelProgress := (weightSum + math.Abs(model.Bias)) / float64(len(model.Weights)+1)
		totalProgress += modelProgress
	}

	lm.progress = math.Min(totalProgress/float64(len(lm.models)), 1.0)
}
