package interaction

import (
	"fmt"
	"sync"
	"time"

	"github.com/phoenix-marie/core/internal/core/memory/v2/processor"
	"github.com/phoenix-marie/core/internal/core/memory/v2/store"
	"github.com/phoenix-marie/core/internal/core/memory/v2/validation"
)

// LayerConfig defines configuration for a memory layer
type LayerConfig struct {
	Name             string
	ProcessorType    string
	ValidationSchema validation.Schema
	CacheEnabled     bool
	BatchSize        int
	RetryAttempts    int
}

// LayerState represents the current state of a memory layer
type LayerState struct {
	Active        bool
	LastSync      time.Time
	ProcessedData int64
	ErrorCount    int64
}

// LayerMetrics contains performance metrics for a layer
type LayerMetrics struct {
	ProcessingTime time.Duration
	StorageLatency time.Duration
	CacheHitRate   float64
	ErrorRate      float64
}

// LayerManager handles interactions between memory layers
type LayerManager struct {
	store      store.StorageEngine
	processors map[string]processor.LayerProcessor
	validator  *validation.ValidationEngine
	configs    map[string]LayerConfig
	states     map[string]LayerState
	metrics    map[string]LayerMetrics
	mu         sync.RWMutex
}

// NewLayerManager creates a new layer manager instance
func NewLayerManager(store store.StorageEngine) *LayerManager {
	return &LayerManager{
		store:      store,
		processors: make(map[string]processor.LayerProcessor),
		validator:  validation.NewValidationEngine(),
		configs:    make(map[string]LayerConfig),
		states:     make(map[string]LayerState),
		metrics:    make(map[string]LayerMetrics),
	}
}

// RegisterLayer registers a new memory layer
func (lm *LayerManager) RegisterLayer(config LayerConfig) error {
	lm.mu.Lock()
	defer lm.mu.Unlock()

	if _, exists := lm.configs[config.Name]; exists {
		return fmt.Errorf("layer already registered: %s", config.Name)
	}

	// Register validation schema
	if err := lm.validator.RegisterSchema(config.Name, config.ValidationSchema); err != nil {
		return fmt.Errorf("failed to register validation schema: %w", err)
	}

	// Initialize processor
	proc, err := lm.createProcessor(config)
	if err != nil {
		return fmt.Errorf("failed to create processor: %w", err)
	}

	// Initialize layer state
	lm.configs[config.Name] = config
	lm.processors[config.Name] = proc
	lm.states[config.Name] = LayerState{
		Active:   true,
		LastSync: time.Now(),
	}
	lm.metrics[config.Name] = LayerMetrics{}

	return nil
}

// ProcessData processes data through a specific layer
func (lm *LayerManager) ProcessData(layer string, data interface{}) error {
	lm.mu.Lock()
	defer lm.mu.Unlock()

	// Validate layer exists
	proc, exists := lm.processors[layer]
	if !exists {
		return fmt.Errorf("layer not found: %s", layer)
	}

	startTime := time.Now()

	// Validate data
	if err := lm.validator.ValidateData(layer, data); err != nil {
		lm.updateMetrics(layer, time.Since(startTime), err)
		return fmt.Errorf("validation failed: %w", err)
	}

	// Process data
	processed, err := proc.Process(data)
	if err != nil {
		lm.updateMetrics(layer, time.Since(startTime), err)
		return fmt.Errorf("processing failed: %w", err)
	}

	// Store processed data
	if err := lm.store.Store(layer, processed.Metadata["key"].(string), processed.Data); err != nil {
		lm.updateMetrics(layer, time.Since(startTime), err)
		return fmt.Errorf("storage failed: %w", err)
	}

	lm.updateMetrics(layer, time.Since(startTime), nil)
	return nil
}

// GetLayerState returns the current state of a layer
func (lm *LayerManager) GetLayerState(layer string) (LayerState, error) {
	lm.mu.RLock()
	defer lm.mu.RUnlock()

	state, exists := lm.states[layer]
	if !exists {
		return LayerState{}, fmt.Errorf("layer not found: %s", layer)
	}

	return state, nil
}

// GetLayerMetrics returns metrics for a specific layer
func (lm *LayerManager) GetLayerMetrics(layer string) (LayerMetrics, error) {
	lm.mu.RLock()
	defer lm.mu.RUnlock()

	metrics, exists := lm.metrics[layer]
	if !exists {
		return LayerMetrics{}, fmt.Errorf("layer not found: %s", layer)
	}

	return metrics, nil
}

// Helper methods

func (lm *LayerManager) createProcessor(config LayerConfig) (processor.LayerProcessor, error) {
	proc := &processor.BaseProcessor{}
	err := proc.Configure(processor.ProcessorConfig{
		BatchSize:     config.BatchSize,
		Timeout:       time.Second * 30,
		RetryAttempts: config.RetryAttempts,
		ValidateInput: true,
		CacheEnabled:  config.CacheEnabled,
	})
	if err != nil {
		return nil, err
	}
	return proc, nil
}

func (lm *LayerManager) updateMetrics(layer string, duration time.Duration, err error) {
	metrics := lm.metrics[layer]
	metrics.ProcessingTime += duration

	state := lm.states[layer]
	state.ProcessedData++

	if err != nil {
		state.ErrorCount++
		metrics.ErrorRate = float64(state.ErrorCount) / float64(state.ProcessedData)
	}

	lm.metrics[layer] = metrics
	lm.states[layer] = state
}
