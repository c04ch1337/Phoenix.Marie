package monitoring

import (
	"fmt"
	"runtime"
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
	metrics := []struct {
		name      string
		typ       MetricType
		unit      string
		component string
	}{
		{"storage.operations", Counter, "operations", "storage"},
		{"storage.latency", Gauge, "milliseconds", "storage"},
		{"storage.memory_usage", Gauge, "bytes", "storage"},
		{"storage.gc_pause", Gauge, "milliseconds", "storage"},

		// Pattern recognition metrics
		{"patterns.detected", Counter, "patterns", "patterns"},
		{"patterns.confidence", Gauge, "ratio", "patterns"},
		{"patterns.processing_time", Gauge, "milliseconds", "patterns"},
		{"patterns.batch_size", Gauge, "count", "patterns"},

		// Learning system metrics
		{"learning.progress", Gauge, "ratio", "learning"},
		{"learning.success_rate", Gauge, "ratio", "learning"},
		{"learning.model_updates", Counter, "updates", "learning"},
		{"learning.feedback_latency", Gauge, "milliseconds", "learning"},

		// Dream processing metrics
		{"dream.processing_time", Gauge, "milliseconds", "dream"},
		{"dream.insight_count", Counter, "insights", "dream"},
		{"dream.memory_usage", Gauge, "bytes", "dream"},
		{"dream.batch_efficiency", Gauge, "ratio", "dream"},

		// Concurrency metrics
		{"concurrent.goroutines", Gauge, "count", "system"},
		{"concurrent.mutex_contentions", Counter, "contentions", "system"},
		{"concurrent.memory_sync", Gauge, "operations", "system"},
	}

	for _, m := range metrics {
		if err := mi.collector.RegisterMetric(
			m.name,
			m.typ,
			m.unit,
			map[string]string{"component": m.component},
		); err != nil {
			return fmt.Errorf("failed to register metric %s: %w", m.name, err)
		}
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

	start := time.Now()

	// Update storage metrics
	if stats := mi.getStorageStats(); stats != nil {
		mi.collector.UpdateMetric("storage.operations", float64(stats.Operations))
		mi.collector.UpdateMetric("storage.latency", stats.AverageLatency)
		mi.collector.UpdateMetric("storage.memory_usage", float64(stats.MemoryUsage))
		mi.collector.UpdateMetric("storage.gc_pause", float64(stats.GCPauseNs)/1e6)
	}

	// Update pattern metrics
	if analysis := mi.patterns.AnalyzePatterns(); analysis.Patterns != nil {
		mi.collector.UpdateMetric("patterns.detected", float64(len(analysis.Patterns)))
		mi.collector.UpdateMetric("patterns.confidence", analysis.AverageConf)
		mi.collector.UpdateMetric("patterns.processing_time", float64(analysis.ProcessingTime.Milliseconds()))
		mi.collector.UpdateMetric("patterns.batch_size", float64(analysis.BatchSize))
	}

	// Update learning metrics
	mi.collector.UpdateMetric("learning.progress", mi.learning.GetProgress())
	if stats := mi.learning.GetStats(); stats.TotalPatterns > 0 {
		mi.collector.UpdateMetric("learning.success_rate", stats.SuccessRate)
		mi.collector.UpdateMetric("learning.model_updates", float64(stats.ModelUpdates))
		mi.collector.UpdateMetric("learning.feedback_latency", float64(stats.FeedbackLatency.Milliseconds()))
	}

	// Update dream metrics
	if dreamStats := mi.getDreamStats(); dreamStats != nil {
		mi.collector.UpdateMetric("dream.processing_time", float64(dreamStats.ProcessingTime.Milliseconds()))
		mi.collector.UpdateMetric("dream.insight_count", float64(dreamStats.InsightCount))
		mi.collector.UpdateMetric("dream.memory_usage", float64(dreamStats.MemoryUsage))
		mi.collector.UpdateMetric("dream.batch_efficiency", dreamStats.BatchEfficiency)
	}

	// Update concurrency metrics
	mi.collector.UpdateMetric("concurrent.goroutines", float64(runtime.NumGoroutine()))
	mi.collector.UpdateMetric("concurrent.mutex_contentions", float64(runtime.NumMutexes()))
	mi.collector.UpdateMetric("concurrent.memory_sync", float64(runtime.NumMemSync()))

	// Collect snapshot with execution time
	snapshot := mi.collector.CollectSnapshot()
	snapshot.Metrics["monitoring.execution_time"] = float64(time.Since(start).Milliseconds())
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
