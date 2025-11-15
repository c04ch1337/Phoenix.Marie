package thought

import (
	"context"
	"sync"
	"time"

	"github.com/phoenix-marie/core/internal/core/memory"
)

// ThoughtEngine represents the core thought processing system
type ThoughtEngine struct {
	memory   *memory.PHL
	patterns *PatternManager
	learner  *LearningManager
	dreamer  *DreamManager
	monitor  *MonitorManager
	mu       sync.RWMutex
	isActive bool
	ctx      context.Context
	cancel   context.CancelFunc
}

// Config holds the configuration for the thought engine
type Config struct {
	MemoryPath     string
	LearningRate   float64
	DreamInterval  time.Duration
	PatternMinConf float64
}

// NewThoughtEngine creates a new thought engine instance
func NewThoughtEngine(cfg *Config) (*ThoughtEngine, error) {
	ctx, cancel := context.WithCancel(context.Background())

	// Initialize memory system
	mem, err := memory.NewPHL(cfg.MemoryPath)
	if err != nil {
		cancel()
		return nil, err
	}

	engine := &ThoughtEngine{
		memory:   mem,
		ctx:      ctx,
		cancel:   cancel,
		isActive: false,
	}

	// Initialize managers
	engine.patterns = NewPatternManager(cfg.PatternMinConf)
	engine.learner = NewLearningManager(cfg.LearningRate)
	engine.dreamer = NewDreamManager(cfg.DreamInterval)
	engine.monitor = NewMonitorManager()

	return engine, nil
}

// Start activates the thought engine and begins processing
func (te *ThoughtEngine) Start() error {
	te.mu.Lock()
	defer te.mu.Unlock()

	if te.isActive {
		return nil
	}

	te.isActive = true
	go te.processThoughts()
	go te.dreamer.Start(te.ctx, te.memory)
	go te.monitor.Start(te.ctx)

	return nil
}

// Stop gracefully shuts down the thought engine
func (te *ThoughtEngine) Stop() error {
	te.mu.Lock()
	defer te.mu.Unlock()

	if !te.isActive {
		return nil
	}

	te.cancel()
	te.isActive = false

	return te.memory.Close()
}

// processThoughts is the main thought processing loop
func (te *ThoughtEngine) processThoughts() {
	ticker := time.NewTicker(100 * time.Millisecond)
	defer ticker.Stop()

	for {
		select {
		case <-te.ctx.Done():
			return
		case <-ticker.C:
			te.processCycle()
		}
	}
}

// processCycle executes one complete thought processing cycle
func (te *ThoughtEngine) processCycle() {
	// Process sensory input
	if input, exists := te.memory.Retrieve("sensory", "current_input"); exists {
		te.patterns.ProcessInput(input)
	}

	// Update learning models
	te.learner.Update(te.patterns.GetPatterns())

	// Propagate insights to memory
	if insights := te.learner.GetInsights(); len(insights) > 0 {
		te.memory.Store("logic", "current_insights", insights)
	}

	// Monitor and record metrics
	te.monitor.RecordMetrics(map[string]float64{
		"pattern_confidence": te.patterns.GetAverageConfidence(),
		"learning_progress":  te.learner.GetProgress(),
	})
}

// GetStatus returns the current status of the thought engine
func (te *ThoughtEngine) GetStatus() map[string]interface{} {
	te.mu.RLock()
	defer te.mu.RUnlock()

	return map[string]interface{}{
		"active":            te.isActive,
		"patterns_detected": te.patterns.GetPatternCount(),
		"learning_progress": te.learner.GetProgress(),
		"dream_active":      te.dreamer.IsActive(),
		"metrics":           te.monitor.GetMetrics(),
	}
}

// InjectThought injects a thought directly into the processing system
func (te *ThoughtEngine) InjectThought(thought interface{}) error {
	return te.memory.Store("logic", "injected_thought", thought)
}

// GetInsights retrieves current insights from the thought process
func (te *ThoughtEngine) GetInsights() ([]interface{}, error) {
	if insights, exists := te.memory.Retrieve("logic", "current_insights"); exists {
		if insightSlice, ok := insights.([]interface{}); ok {
			return insightSlice, nil
		}
	}
	return nil, nil
}
