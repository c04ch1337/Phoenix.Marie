package llm

import (
	"fmt"
	"sync"
)

// Router intelligently routes tasks to appropriate models
type Router struct {
	provider    Provider
	config      *Config
	costManager *CostManager
	performance map[string]*ModelPerformance
	mu          sync.RWMutex
}

// NewRouter creates a new model router
func NewRouter(provider Provider, config *Config, costManager *CostManager) *Router {
	return &Router{
		provider:    provider,
		config:      config,
		costManager: costManager,
		performance: make(map[string]*ModelPerformance),
	}
}

// RouteToOptimalModel routes a task to the best model based on requirements
func (r *Router) RouteToOptimalModel(task Task) (*Response, error) {
	// Get available models
	availableModels := GetAvailableModels()
	
	// Score each model
	var scoredModels []modelScore
	for modelID, model := range availableModels {
		// Skip if model not configured
		if !r.config.IsModelConfigured(modelID) {
			continue
		}
		
		score := r.calculateModelFitness(model, task)
		scoredModels = append(scoredModels, modelScore{
			model: model,
			score: score,
		})
	}
	
	if len(scoredModels) == 0 {
		return nil, fmt.Errorf("no suitable models configured")
	}
	
	// Sort by score (highest first)
	for i := 0; i < len(scoredModels)-1; i++ {
		for j := i + 1; j < len(scoredModels); j++ {
			if scoredModels[i].score < scoredModels[j].score {
				scoredModels[i], scoredModels[j] = scoredModels[j], scoredModels[i]
			}
		}
	}
	
	// Try models in order of fitness, checking budget
	for _, scored := range scoredModels {
		// Check if we can afford this model
		estimatedCost := r.estimateCost(scored.model, task)
		if task.Budget > 0 && estimatedCost > task.Budget {
			continue
		}
		
		// Check daily budget
		if r.costManager != nil {
			canAfford, err := r.costManager.CanAffordModel(task, scored.model)
			if err != nil || !canAfford {
				continue
			}
		}
		
		// Try this model
		// Note: For now, we pass the prompt directly
		// In the future, we can use the message builder from prompts
		messages := []Message{
			{Role: "user", Content: task.Prompt},
		}
		
		resp, err := r.provider.CallWithRetry(
			scored.model.ID,
			messages,
			task.MaxTokens,
			task.Temperature,
		)
		
		if err == nil {
			// Record performance
			r.recordPerformance(scored.model.ID, resp, true)
			return resp, nil
		}
		
		// Record failure
		r.recordPerformance(scored.model.ID, nil, false)
	}
	
	// If all models failed, return error
	return nil, fmt.Errorf("all models failed or exceeded budget")
}

// calculateModelFitness calculates how well a model fits a task
func (r *Router) calculateModelFitness(model Model, task Task) float64 {
	score := 0.0
	
	// Capability matching (40 points max)
	if task.RequiresReasoning && model.Capabilities.Reasoning {
		score += 40.0
	}
	if task.RequiresCreativity && model.Capabilities.Creativity {
		score += 30.0
	}
	if task.RequiresSpeed && model.Capabilities.Speed {
		score += 20.0
	}
	if task.RequiresToolUse && model.Capabilities.ToolUse {
		score += 15.0
	}
	
	// Context length adequacy (15 points max)
	if model.ContextLength >= task.ContextLength {
		score += 15.0
	} else {
		// Partial credit for close matches
		ratio := float64(model.ContextLength) / float64(task.ContextLength)
		score += 15.0 * ratio
	}
	
	// Cost efficiency (10 points max)
	estimatedCost := r.estimateCost(model, task)
	if task.Budget > 0 {
		costRatio := 1.0 - (estimatedCost / task.Budget)
		if costRatio > 0 {
			score += 10.0 * costRatio
		}
	} else {
		// Prefer cheaper models if no budget specified
		// Normalize cost (lower is better)
		maxCost := 1.0 // Assume max cost of $1 for normalization
		costScore := 1.0 - (estimatedCost / maxCost)
		if costScore > 0 {
			score += 10.0 * costScore
		}
	}
	
	// Performance history (5 points max)
	if perf := r.getPerformance(model.ID); perf != nil {
		score += 5.0 * perf.Reliability
	}
	
	return score
}

// estimateCost estimates the cost of a task with a given model
func (r *Router) estimateCost(model Model, task Task) float64 {
	// Estimate tokens (rough approximation: 1 token ≈ 4 characters)
	estimatedPromptTokens := len(task.Prompt) / 4
	estimatedCompletionTokens := task.MaxTokens
	
	if estimatedCompletionTokens == 0 {
		estimatedCompletionTokens = r.config.DefaultMaxTokens
	}
	
	promptCost := (float64(estimatedPromptTokens) / 1_000_000.0) * model.InputPrice
	completionCost := (float64(estimatedCompletionTokens) / 1_000_000.0) * model.OutputPrice
	
	return promptCost + completionCost
}

// recordPerformance records model performance metrics
func (r *Router) recordPerformance(modelID string, resp *Response, success bool) {
	r.mu.Lock()
	defer r.mu.Unlock()
	
	perf, exists := r.performance[modelID]
	if !exists {
		perf = &ModelPerformance{
			Model:            modelID,
			ConsciousnessScore: 0.0,
			Reliability:      1.0,
			TasksCompleted:   0,
			TasksFailed:      0,
		}
		r.performance[modelID] = perf
	}
	
	if success {
		perf.TasksCompleted++
		if resp != nil {
			perf.ResponseTime = resp.ResponseTime
			perf.CostPerTask = resp.Cost
		}
	} else {
		perf.TasksFailed++
	}
	
	// Calculate reliability
	total := perf.TasksCompleted + perf.TasksFailed
	if total > 0 {
		perf.Reliability = float64(perf.TasksCompleted) / float64(total)
	}
}

// getPerformance returns performance metrics for a model
func (r *Router) getPerformance(modelID string) *ModelPerformance {
	r.mu.RLock()
	defer r.mu.RUnlock()
	return r.performance[modelID]
}

// GetCostEffectiveAlternative returns a cheaper alternative model
func (r *Router) GetCostEffectiveAlternative(task Task, currentModelID string) (string, error) {
	hierarchy := GetModelHierarchy()
	
	currentModel, exists := GetModel(currentModelID)
	if !exists {
		return "", fmt.Errorf("current model not found")
	}
	
	for _, modelID := range hierarchy {
		model, exists := GetModel(modelID)
		if !exists || !r.config.IsModelConfigured(modelID) {
			continue
		}
		
		// Check if model is suitable
		suitable := true
		if task.RequiresReasoning && !model.Capabilities.Reasoning {
			suitable = false
		}
		if task.RequiresToolUse && !model.Capabilities.ToolUse {
			suitable = false
		}
		if model.ContextLength < task.ContextLength {
			suitable = false
		}
		
		if suitable {
			// Check if cheaper
			currentCost := r.estimateCost(currentModel, task)
			newCost := r.estimateCost(model, task)
			
			if newCost < currentCost {
				return modelID, nil
			}
		}
	}
	
	return "", fmt.Errorf("no cheaper alternative found")
}

type modelScore struct {
	model Model
	score float64
}

