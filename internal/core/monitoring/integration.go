package monitoring

import (
	"fmt"
	"sync"
	"time"

	"github.com/phoenix-marie/core/internal/core/memory/v2/store"
	"github.com/phoenix-marie/core/internal/core/thought/v2/learning"
	"github.com/phoenix-marie/core/internal/core/thought/v2/pattern"
)

// MonitoringIntegration provides monitoring for core system components
type MonitoringIntegration struct {
	collector   *MetricsCollector
	storage     store.StorageEngine
	patterns    *pattern.Manager
	learning    *learning.Manager
	initialized bool
	stopChan    chan struct{}
	interval    time.Duration
	mu          sync.RWMutex
}

// NewMonitoringIntegration creates a new monitoring integration
func NewMonitoringIntegration(
	storage store.StorageEngine,
	patterns *pattern.Manager,
	learning *learning.Manager,
	interval time.Duration,
) *MonitoringIntegration {
	return &MonitoringIntegration{
		collector: NewMetricsCollector(1000), // Keep 1000 snapshots
		storage:   storage,
		patterns:  patterns,
		learning:  learning,
		interval:  interval,
		stopChan:  make(chan struct{}),
	}
}

// Start initializes and starts monitoring
func (mi *MonitoringIntegration) Start() error {
	mi.mu.Lock()
	defer mi.mu.Unlock()

	if mi.initialized {
		return fmt.Errorf("monitoring already initialized")
	}

	// Register core metrics
	if err := mi.registerCoreMetrics(); err != nil {
		return fmt.Errorf("failed to register metrics: %w", err)
	}

	mi.initialized = true

	// Start monitoring goroutine
	go mi.monitor()

	return nil
}

// Stop stops the monitoring system
func (mi *MonitoringIntegration) Stop() error {
	mi.mu.Lock()
	defer mi.mu.Unlock()

	if !mi.initialized {
		return fmt.Errorf("monitoring not initialized")
	}

	close(mi.stopChan)
	mi.initialized = false

	return nil
}

// GetAnalysis returns current performance analysis
func (mi *MonitoringIntegration) GetAnalysis() (PerformanceAnalysis, error) {
	mi.mu.RLock()
	defer mi.mu.RUnlock()

	if !mi.initialized {
		return PerformanceAnalysis{}, fmt.Errorf("monitoring not initialized")
	}

	return mi.collector.AnalyzePerformance(), nil
}

// GetLatestSnapshot returns the most recent metrics snapshot
func (mi *MonitoringIntegration) GetLatestSnapshot() (MetricsSnapshot, error) {
	mi.mu.RLock()
	defer mi.mu.RUnlock()

	if !mi.initialized {
		return MetricsSnapshot{}, fmt.Errorf("monitoring not initialized")
	}

	snapshots := mi.collector.GetSnapshots()
	if len(snapshots) == 0 {
		return MetricsSnapshot{}, fmt.Errorf("no snapshots available")
	}

	return snapshots[len(snapshots)-1], nil
}

// Helper methods

func (mi *MonitoringIntegration) registerCoreMetrics() error {
	// Storage metrics
	if err := mi.collector.RegisterMetric(
		"storage.operations",
		Counter,
		"operations",
		map[string]string{"component": "storage"},
	); err != nil {
		return err
	}

	if err := mi.collector.RegisterMetric(
		"storage.latency",
		Gauge,
		"milliseconds",
		map[string]string{"component": "storage"},
	); err != nil {
		return err
	}

	// Pattern metrics
	if err := mi.collector.RegisterMetric(
		"patterns.detected",
		Counter,
		"patterns",
		map[string]string{"component": "patterns"},
	); err != nil {
		return err
	}

	if err := mi.collector.RegisterMetric(
		"patterns.confidence",
		Gauge,
		"ratio",
		map[string]string{"component": "patterns"},
	); err != nil {
		return err
	}

	// Learning metrics
	if err := mi.collector.RegisterMetric(
		"learning.progress",
		Gauge,
		"ratio",
		map[string]string{"component": "learning"},
	); err != nil {
		return err
	}

	if err := mi.collector.RegisterMetric(
		"learning.success_rate",
		Gauge,
		"ratio",
		map[string]string{"component": "learning"},
	); err != nil {
		return err
	}

	return nil
}

func (mi *MonitoringIntegration) monitor() {
	ticker := time.NewTicker(mi.interval)
	defer ticker.Stop()

	for {
		select {
		case <-mi.stopChan:
			return
		case <-ticker.C:
			mi.updateMetrics()
		}
	}
}

func (mi *MonitoringIntegration) updateMetrics() {
	mi.mu.Lock()
	defer mi.mu.Unlock()

	// Update storage metrics
	if stats := mi.getStorageStats(); stats != nil {
		mi.collector.UpdateMetric("storage.operations", float64(stats.Operations))
		mi.collector.UpdateMetric("storage.latency", stats.AverageLatency)
	}

	// Update pattern metrics
	if analysis := mi.patterns.AnalyzePatterns(); analysis.Patterns != nil {
		mi.collector.UpdateMetric("patterns.detected", float64(len(analysis.Patterns)))
		mi.collector.UpdateMetric("patterns.confidence", analysis.AverageConf)
	}

	// Update learning metrics
	mi.collector.UpdateMetric("learning.progress", mi.learning.GetProgress())
	if stats := mi.learning.GetStats(); stats.TotalPatterns > 0 {
		successRate := stats.SuccessRate
		mi.collector.UpdateMetric("learning.success_rate", successRate)
	}

	// Collect snapshot
	mi.collector.CollectSnapshot()
}

func (mi *MonitoringIntegration) getStorageStats() *StorageStats {
	// Implementation would get actual storage stats
	// This is a placeholder
	return nil
}

// StorageStats represents storage performance statistics
type StorageStats struct {
	Operations     int64
	AverageLatency float64
}
