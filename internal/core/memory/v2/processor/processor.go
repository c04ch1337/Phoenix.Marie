package processor

import (
	"fmt"
	"time"
)

// ProcessorState represents the current state of a layer processor
type ProcessorState struct {
	Active         bool
	LastRun        time.Time
	ProcessedCount int64
	ErrorCount     int64
}

// ProcessorConfig contains configuration for a layer processor
type ProcessorConfig struct {
	BatchSize     int
	Timeout       time.Duration
	RetryAttempts int
	ValidateInput bool
	CacheEnabled  bool
}

// ProcessorMetrics contains performance metrics for a processor
type ProcessorMetrics struct {
	ProcessingTime time.Duration
	SuccessRate    float64
	ErrorRate      float64
	AverageLatency time.Duration
	ResourceUsage  float64
}

// ProcessedData represents the output of a processor
type ProcessedData struct {
	Data       any
	Metadata   map[string]interface{}
	Timestamp  time.Time
	Confidence float64
}

// BaseProcessor provides common functionality for layer processors
type BaseProcessor struct {
	state   ProcessorState
	config  ProcessorConfig
	metrics ProcessorMetrics
}

// Process implements the core processing logic
func (bp *BaseProcessor) Process(data any) (ProcessedData, error) {
	startTime := time.Now()

	// Validate input if enabled
	if bp.config.ValidateInput {
		if err := bp.Validate(data); err != nil {
			bp.state.ErrorCount++
			return ProcessedData{}, fmt.Errorf("validation failed: %w", err)
		}
	}

	// Process the data (to be implemented by specific processors)
	processed := ProcessedData{
		Data:      data,
		Metadata:  make(map[string]interface{}),
		Timestamp: time.Now(),
	}

	// Update metrics
	bp.metrics.ProcessingTime += time.Since(startTime)
	bp.state.ProcessedCount++
	bp.updateMetrics()

	return processed, nil
}

// Validate performs basic data validation
func (bp *BaseProcessor) Validate(data any) error {
	if data == nil {
		return fmt.Errorf("data cannot be nil")
	}
	return nil
}

// GetState returns the current processor state
func (bp *BaseProcessor) GetState() ProcessorState {
	return bp.state
}

// Reset resets the processor state
func (bp *BaseProcessor) Reset() error {
	bp.state = ProcessorState{
		Active:         true,
		LastRun:        time.Now(),
		ProcessedCount: 0,
		ErrorCount:     0,
	}
	return nil
}

// Configure updates the processor configuration
func (bp *BaseProcessor) Configure(config ProcessorConfig) error {
	if config.BatchSize <= 0 {
		return fmt.Errorf("batch size must be positive")
	}
	if config.Timeout <= 0 {
		return fmt.Errorf("timeout must be positive")
	}
	bp.config = config
	return nil
}

// GetMetrics returns the current processor metrics
func (bp *BaseProcessor) GetMetrics() ProcessorMetrics {
	return bp.metrics
}

// updateMetrics updates the processor metrics
func (bp *BaseProcessor) updateMetrics() {
	totalCount := bp.state.ProcessedCount + bp.state.ErrorCount
	if totalCount > 0 {
		bp.metrics.SuccessRate = float64(bp.state.ProcessedCount) / float64(totalCount)
		bp.metrics.ErrorRate = float64(bp.state.ErrorCount) / float64(totalCount)
	}
}

// LayerProcessor defines the interface for processing data in memory layers
type LayerProcessor interface {
	// Core processing
	Process(data any) (ProcessedData, error)
	Validate(data any) error

	// State management
	GetState() ProcessorState
	Reset() error

	// Configuration
	Configure(config ProcessorConfig) error

	// Metrics
	GetMetrics() ProcessorMetrics
}
