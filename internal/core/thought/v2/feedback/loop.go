package feedback

import (
	"fmt"
	"sync"
	"time"

	"github.com/phoenix-marie/core/internal/core/thought/v2/dream"
	"github.com/phoenix-marie/core/internal/core/thought/v2/integration"
	"github.com/phoenix-marie/core/internal/core/thought/v2/learning"
	"github.com/phoenix-marie/core/internal/core/thought/v2/pattern"
)

// FeedbackLoop manages the continuous feedback cycle between components
type FeedbackLoop struct {
	patterns *pattern.Manager
	learning *learning.Manager
	dream    *dream.Processor
	bridge   *integration.MemoryBridge
	config   LoopConfig
	state    LoopState
	metrics  LoopMetrics
	stopChan chan struct{}
	mu       sync.RWMutex
}

// LoopConfig contains configuration for the feedback loop
type LoopConfig struct {
	UpdateInterval time.Duration
	BatchSize      int
	MinConfidence  float64
	MaxQueueSize   int
	EnableDreaming bool
}

// LoopState represents the current state of the feedback loop
type LoopState struct {
	Active         bool
	LastUpdate     time.Time
	ProcessedCount int64
	QueueSize      int
	DreamingActive bool
}

// LoopMetrics contains performance metrics for the feedback loop
type LoopMetrics struct {
	ProcessingTime time.Duration
	SuccessRate    float64
	ErrorRate      float64
	ThroughputRate float64
	QueueLatency   time.Duration
}

// NewFeedbackLoop creates a new feedback loop instance
func NewFeedbackLoop(
	patterns *pattern.Manager,
	learning *learning.Manager,
	dream *dream.Processor,
	bridge *integration.MemoryBridge,
	config LoopConfig,
) *FeedbackLoop {
	return &FeedbackLoop{
		patterns: patterns,
		learning: learning,
		dream:    dream,
		bridge:   bridge,
		config:   config,
		stopChan: make(chan struct{}),
	}
}

// Start activates the feedback loop
func (fl *FeedbackLoop) Start() error {
	fl.mu.Lock()
	defer fl.mu.Unlock()

	if fl.state.Active {
		return fmt.Errorf("feedback loop already active")
	}

	fl.state.Active = true
	fl.state.LastUpdate = time.Now()

	// Start processing goroutine
	go fl.processLoop()

	// Start dream processing if enabled
	if fl.config.EnableDreaming {
		if err := fl.startDreaming(); err != nil {
			fl.Stop()
			return fmt.Errorf("dream processing start failed: %w", err)
		}
	}

	return nil
}

// Stop deactivates the feedback loop
func (fl *FeedbackLoop) Stop() error {
	fl.mu.Lock()
	defer fl.mu.Unlock()

	if !fl.state.Active {
		return fmt.Errorf("feedback loop not active")
	}

	// Signal processing loop to stop
	close(fl.stopChan)

	// Stop dream processing if active
	if fl.state.DreamingActive {
		if err := fl.stopDreaming(); err != nil {
			return fmt.Errorf("dream processing stop failed: %w", err)
		}
	}

	fl.state.Active = false
	return nil
}

// ProcessFeedback processes new feedback through the loop
func (fl *FeedbackLoop) ProcessFeedback(feedback learning.Feedback) error {
	fl.mu.Lock()
	defer fl.mu.Unlock()

	startTime := time.Now()

	// Validate feedback
	if err := fl.validateFeedback(feedback); err != nil {
		fl.updateMetrics(startTime, err)
		return fmt.Errorf("invalid feedback: %w", err)
	}

	// Process through learning system
	if err := fl.learning.Adapt(feedback); err != nil {
		fl.updateMetrics(startTime, err)
		return fmt.Errorf("learning adaptation failed: %w", err)
	}

	// Update pattern confidence
	pattern, err := fl.bridge.RetrievePattern(feedback.PatternID)
	if err != nil {
		fl.updateMetrics(startTime, err)
		return fmt.Errorf("pattern retrieval failed: %w", err)
	}

	pattern.Confidence = fl.learning.GetProgress()
	if err := fl.patterns.UpdatePattern(pattern); err != nil {
		fl.updateMetrics(startTime, err)
		return fmt.Errorf("pattern update failed: %w", err)
	}

	// Store updated pattern
	if err := fl.bridge.StorePattern(pattern); err != nil {
		fl.updateMetrics(startTime, err)
		return fmt.Errorf("pattern storage failed: %w", err)
	}

	fl.updateMetrics(startTime, nil)
	return nil
}

// GetState returns the current state of the feedback loop
func (fl *FeedbackLoop) GetState() LoopState {
	fl.mu.RLock()
	defer fl.mu.RUnlock()
	return fl.state
}

// GetMetrics returns the current metrics of the feedback loop
func (fl *FeedbackLoop) GetMetrics() LoopMetrics {
	fl.mu.RLock()
	defer fl.mu.RUnlock()
	return fl.metrics
}

// Helper methods

func (fl *FeedbackLoop) processLoop() {
	ticker := time.NewTicker(fl.config.UpdateInterval)
	defer ticker.Stop()

	for {
		select {
		case <-fl.stopChan:
			return
		case <-ticker.C:
			if err := fl.updateLoop(); err != nil {
				// Log error but continue processing
				fmt.Printf("Loop update failed: %v\n", err)
			}
		}
	}
}

func (fl *FeedbackLoop) updateLoop() error {
	fl.mu.Lock()
	defer fl.mu.Unlock()

	// Sync patterns with memory
	if err := fl.bridge.SyncPatterns(); err != nil {
		return fmt.Errorf("pattern sync failed: %w", err)
	}

	// Store learning state
	if err := fl.bridge.StoreLearningState(); err != nil {
		return fmt.Errorf("learning state storage failed: %w", err)
	}

	// Update state
	fl.state.LastUpdate = time.Now()
	fl.state.ProcessedCount++

	return nil
}

func (fl *FeedbackLoop) startDreaming() error {
	if err := fl.dream.Start(); err != nil {
		return fmt.Errorf("dream processor start failed: %w", err)
	}

	fl.state.DreamingActive = true
	return nil
}

func (fl *FeedbackLoop) stopDreaming() error {
	if err := fl.dream.Stop(); err != nil {
		return fmt.Errorf("dream processor stop failed: %w", err)
	}

	fl.state.DreamingActive = false
	return nil
}

func (fl *FeedbackLoop) validateFeedback(feedback learning.Feedback) error {
	if feedback.PatternID == "" {
		return fmt.Errorf("pattern ID is required")
	}
	if feedback.Score < 0 || feedback.Score > 1 {
		return fmt.Errorf("score must be between 0 and 1")
	}
	return nil
}

func (fl *FeedbackLoop) updateMetrics(startTime time.Time, err error) {
	duration := time.Since(startTime)
	fl.metrics.ProcessingTime += duration

	if err != nil {
		fl.metrics.ErrorRate = (fl.metrics.ErrorRate*float64(fl.state.ProcessedCount) + 1) /
			float64(fl.state.ProcessedCount+1)
	} else {
		fl.metrics.SuccessRate = (fl.metrics.SuccessRate*float64(fl.state.ProcessedCount) + 1) /
			float64(fl.state.ProcessedCount+1)
	}

	fl.metrics.ThroughputRate = float64(fl.state.ProcessedCount) /
		time.Since(fl.state.LastUpdate).Seconds()
}
