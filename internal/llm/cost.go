package llm

import (
	"fmt"
	"sync"
	"time"
)

// CostManager manages LLM API costs and budgets
type CostManager struct {
	config        *Config
	dailySpend    float64
	monthlySpend  float64
	lastReset     time.Time
	spendHistory  []CostRecord
	mu            sync.RWMutex
}

// CostRecord tracks a single cost transaction
type CostRecord struct {
	Timestamp time.Time
	Model     string
	Cost      float64
	TaskType  TaskType
}

// NewCostManager creates a new cost manager
func NewCostManager(config *Config) *CostManager {
	return &CostManager{
		config:       config,
		dailySpend:   0.0,
		monthlySpend: 0.0,
		lastReset:    time.Now(),
		spendHistory: make([]CostRecord, 0),
	}
}

// CanAffordModel checks if we can afford a model for a task
func (cm *CostManager) CanAffordModel(task Task, model Model) (bool, error) {
	cm.mu.RLock()
	defer cm.mu.RUnlock()
	
	// Reset daily spend if it's a new day
	cm.checkAndReset()
	
	// Estimate cost
	estimatedCost := cm.estimateTaskCost(task, model)
	
	// Check daily budget
	projectedDaily := cm.dailySpend + estimatedCost
	dailyBudget := cm.config.DailyBudget
	
	// Allow 10% overage buffer
	if projectedDaily > dailyBudget*1.1 {
		return false, fmt.Errorf("would exceed daily budget: $%.2f / $%.2f", projectedDaily, dailyBudget)
	}
	
	// Check monthly budget
	projectedMonthly := cm.monthlySpend + estimatedCost
	monthlyBudget := cm.config.MonthlyBudget
	
	if projectedMonthly > monthlyBudget*1.1 {
		return false, fmt.Errorf("would exceed monthly budget: $%.2f / $%.2f", projectedMonthly, monthlyBudget)
	}
	
	return true, nil
}

// RecordCost records a cost transaction
func (cm *CostManager) RecordCost(modelID string, cost float64, taskType TaskType) {
	cm.mu.Lock()
	defer cm.mu.Unlock()
	
	cm.checkAndReset()
	
	cm.dailySpend += cost
	cm.monthlySpend += cost
	
	cm.spendHistory = append(cm.spendHistory, CostRecord{
		Timestamp: time.Now(),
		Model:     modelID,
		Cost:      cost,
		TaskType:  taskType,
	})
	
	// Keep only last 1000 records
	if len(cm.spendHistory) > 1000 {
		cm.spendHistory = cm.spendHistory[len(cm.spendHistory)-1000:]
	}
}

// GetDailySpend returns current daily spend
func (cm *CostManager) GetDailySpend() float64 {
	cm.mu.RLock()
	defer cm.mu.RUnlock()
	cm.checkAndReset()
	return cm.dailySpend
}

// GetMonthlySpend returns current monthly spend
func (cm *CostManager) GetMonthlySpend() float64 {
	cm.mu.RLock()
	defer cm.mu.RUnlock()
	return cm.monthlySpend
}

// GetRemainingDailyBudget returns remaining daily budget
func (cm *CostManager) GetRemainingDailyBudget() float64 {
	cm.mu.RLock()
	defer cm.mu.RUnlock()
	cm.checkAndReset()
	return cm.config.DailyBudget - cm.dailySpend
}

// GetRemainingMonthlyBudget returns remaining monthly budget
func (cm *CostManager) GetRemainingMonthlyBudget() float64 {
	cm.mu.RLock()
	defer cm.mu.RUnlock()
	return cm.config.MonthlyBudget - cm.monthlySpend
}

// GetCostEffectiveAlternative returns a cheaper alternative model
func (cm *CostManager) GetCostEffectiveAlternative(task Task, currentModelID string) (string, error) {
	hierarchy := GetModelHierarchy()
	
	currentModel, exists := GetModel(currentModelID)
	if !exists {
		return "", fmt.Errorf("current model not found")
	}
	
	currentCost := cm.estimateTaskCost(task, currentModel)
	
	for _, modelID := range hierarchy {
		model, exists := GetModel(modelID)
		if !exists {
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
			newCost := cm.estimateTaskCost(task, model)
			if newCost < currentCost {
				return modelID, nil
			}
		}
	}
	
	return "", fmt.Errorf("no cheaper alternative found")
}

// estimateTaskCost estimates the cost of a task with a given model
func (cm *CostManager) estimateTaskCost(task Task, model Model) float64 {
	// Estimate tokens (rough approximation: 1 token ≈ 4 characters)
	estimatedPromptTokens := len(task.Prompt) / 4
	estimatedCompletionTokens := task.MaxTokens
	
	if estimatedCompletionTokens == 0 {
		estimatedCompletionTokens = cm.config.DefaultMaxTokens
	}
	
	promptCost := (float64(estimatedPromptTokens) / 1_000_000.0) * model.InputPrice
	completionCost := (float64(estimatedCompletionTokens) / 1_000_000.0) * model.OutputPrice
	
	return promptCost + completionCost
}

// checkAndReset checks if we need to reset daily spend (must be called with lock held)
func (cm *CostManager) checkAndReset() {
	now := time.Now()
	if now.Day() != cm.lastReset.Day() || now.Month() != cm.lastReset.Month() || now.Year() != cm.lastReset.Year() {
		cm.dailySpend = 0.0
		cm.lastReset = now
	}
	
	// Reset monthly spend on first day of month
	if now.Day() == 1 && now.Month() != cm.lastReset.Month() {
		cm.monthlySpend = 0.0
	}
}

// GetSpendHistory returns recent spend history
func (cm *CostManager) GetSpendHistory(limit int) []CostRecord {
	cm.mu.RLock()
	defer cm.mu.RUnlock()
	
	if limit <= 0 || limit > len(cm.spendHistory) {
		limit = len(cm.spendHistory)
	}
	
	// Return most recent records
	start := len(cm.spendHistory) - limit
	if start < 0 {
		start = 0
	}
	
	result := make([]CostRecord, limit)
	copy(result, cm.spendHistory[start:])
	return result
}

// GetStats returns cost statistics
func (cm *CostManager) GetStats() CostStats {
	cm.mu.RLock()
	defer cm.mu.RUnlock()
	
	cm.checkAndReset()
	
	stats := CostStats{
		DailySpend:          cm.dailySpend,
		MonthlySpend:         cm.monthlySpend,
		DailyBudget:         cm.config.DailyBudget,
		MonthlyBudget:       cm.config.MonthlyBudget,
		RemainingDaily:      cm.config.DailyBudget - cm.dailySpend,
		RemainingMonthly:    cm.config.MonthlyBudget - cm.monthlySpend,
		TotalTransactions:  len(cm.spendHistory),
	}
	
	// Calculate average cost per transaction
	if len(cm.spendHistory) > 0 {
		total := 0.0
		for _, record := range cm.spendHistory {
			total += record.Cost
		}
		stats.AverageCostPerTransaction = total / float64(len(cm.spendHistory))
	}
	
	return stats
}

// CostStats contains cost statistics
type CostStats struct {
	DailySpend                  float64
	MonthlySpend                float64
	DailyBudget                 float64
	MonthlyBudget               float64
	RemainingDaily              float64
	RemainingMonthly            float64
	TotalTransactions            int
	AverageCostPerTransaction   float64
}

