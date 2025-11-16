package thought

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/phoenix-marie/core/internal/core/memory"
)

// DreamManager is a stub for dream processing (to be implemented)
type DreamManager struct {
	interval time.Duration
	isActive bool
}

// MonitorManager is a stub for monitoring (to be implemented)
type MonitorManager struct {
	metrics map[string]float64
}

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
	engine.dreamer = &DreamManager{interval: cfg.DreamInterval, isActive: false}
	engine.monitor = &MonitorManager{metrics: make(map[string]float64)}

	return engine, nil
}

// NewThoughtEngineWithMemory creates a thought engine using an existing memory instance
// This allows sharing memory between Phoenix and ThoughtEngine
func NewThoughtEngineWithMemory(mem *memory.PHL, learningRate float64, patternMinConf float64) (*ThoughtEngine, error) {
	ctx, cancel := context.WithCancel(context.Background())

	engine := &ThoughtEngine{
		memory:   mem, // Use shared memory instance
		ctx:      ctx,
		cancel:   cancel,
		isActive: false,
	}

	// Initialize managers
	engine.patterns = NewPatternManager(patternMinConf)
	engine.learner = NewLearningManager(learningRate)
	engine.dreamer = &DreamManager{interval: 5 * time.Minute, isActive: false}
	engine.monitor = &MonitorManager{metrics: make(map[string]float64)}

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
	// TODO: Implement dream manager start
	// go te.dreamer.Start(te.ctx, te.memory)
	// TODO: Implement monitor manager start
	// go te.monitor.Start(te.ctx)

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

	// Don't close memory if it's shared (we'll let Phoenix handle it)
	// Only close if we created it ourselves
	// For now, we'll skip closing to avoid issues with shared memory
	return nil
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
	te.monitor.metrics["pattern_confidence"] = te.patterns.GetAverageConfidence()
	te.monitor.metrics["learning_progress"] = te.learner.GetProgress()
}

// GetStatus returns the current status of the thought engine
func (te *ThoughtEngine) GetStatus() map[string]interface{} {
	te.mu.RLock()
	defer te.mu.RUnlock()

	return map[string]interface{}{
		"active":            te.isActive,
		"patterns_detected": te.patterns.GetPatternCount(),
		"learning_progress": te.learner.GetProgress(),
		"dream_active":      te.dreamer.isActive,
		"metrics":           te.monitor.metrics,
	}
}

// InjectThought injects a thought directly into the processing system
func (te *ThoughtEngine) InjectThought(thought interface{}) error {
	success := te.memory.Store("logic", "injected_thought", thought)
	if !success {
		return fmt.Errorf("failed to store injected thought")
	}
	return nil
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
