package monitoring

import (
	"fmt"
	"sync"
	"time"

	"runtime"
)

// MetricsCollector handles system-wide performance monitoring
type MetricsCollector struct {
	metrics      map[string]*Metric
	snapshots    []MetricsSnapshot
	maxSnapshots int
	startTime    time.Time
	mu           sync.RWMutex
}

// Metric represents a single performance metric
type Metric struct {
	Name        string
	Value       float64
	Type        MetricType
	Unit        string
	Labels      map[string]string
	LastUpdated time.Time
}

// MetricType defines the type of metric being collected
type MetricType string

const (
	Counter   MetricType = "counter"
	Gauge     MetricType = "gauge"
	Histogram MetricType = "histogram"
)

// MetricsSnapshot represents a point-in-time collection of metrics
type MetricsSnapshot struct {
	Timestamp time.Time
	Metrics   map[string]float64
	System    SystemMetrics
}

// SystemMetrics contains system-level performance metrics
type SystemMetrics struct {
	CPUUsage    float64
	MemoryUsage uint64
	GCStats     runtime.MemStats
	Goroutines  int
}

// NewMetricsCollector creates a new metrics collector instance
func NewMetricsCollector(maxSnapshots int) *MetricsCollector {
	return &MetricsCollector{
		metrics:      make(map[string]*Metric),
		snapshots:    make([]MetricsSnapshot, 0),
		maxSnapshots: maxSnapshots,
		startTime:    time.Now(),
	}
}

// RegisterMetric registers a new metric for collection
func (mc *MetricsCollector) RegisterMetric(name string, metricType MetricType, unit string, labels map[string]string) error {
	mc.mu.Lock()
	defer mc.mu.Unlock()

	if _, exists := mc.metrics[name]; exists {
		return fmt.Errorf("metric already registered: %s", name)
	}

	mc.metrics[name] = &Metric{
		Name:        name,
		Type:        metricType,
		Unit:        unit,
		Labels:      labels,
		LastUpdated: time.Now(),
	}

	return nil
}

// UpdateMetric updates the value of a metric
func (mc *MetricsCollector) UpdateMetric(name string, value float64) error {
	mc.mu.Lock()
	defer mc.mu.Unlock()

	metric, exists := mc.metrics[name]
	if !exists {
		return fmt.Errorf("metric not found: %s", name)
	}

	metric.Value = value
	metric.LastUpdated = time.Now()
	return nil
}

// IncrementCounter increments a counter metric
func (mc *MetricsCollector) IncrementCounter(name string, value float64) error {
	mc.mu.Lock()
	defer mc.mu.Unlock()

	metric, exists := mc.metrics[name]
	if !exists {
		return fmt.Errorf("metric not found: %s", name)
	}

	if metric.Type != Counter {
		return fmt.Errorf("metric %s is not a counter", name)
	}

	metric.Value += value
	metric.LastUpdated = time.Now()
	return nil
}

// CollectSnapshot takes a snapshot of current metrics
func (mc *MetricsCollector) CollectSnapshot() MetricsSnapshot {
	mc.mu.Lock()
	defer mc.mu.Unlock()

	snapshot := MetricsSnapshot{
		Timestamp: time.Now(),
		Metrics:   make(map[string]float64),
		System:    mc.collectSystemMetrics(),
	}

	for name, metric := range mc.metrics {
		snapshot.Metrics[name] = metric.Value
	}

	mc.snapshots = append(mc.snapshots, snapshot)
	if len(mc.snapshots) > mc.maxSnapshots {
		mc.snapshots = mc.snapshots[1:]
	}

	return snapshot
}

// GetMetric retrieves a specific metric
func (mc *MetricsCollector) GetMetric(name string) (*Metric, error) {
	mc.mu.RLock()
	defer mc.mu.RUnlock()

	metric, exists := mc.metrics[name]
	if !exists {
		return nil, fmt.Errorf("metric not found: %s", name)
	}

	return metric, nil
}

// GetSnapshots returns all collected snapshots
func (mc *MetricsCollector) GetSnapshots() []MetricsSnapshot {
	mc.mu.RLock()
	defer mc.mu.RUnlock()
	return mc.snapshots
}

// AnalyzePerformance analyzes performance trends
func (mc *MetricsCollector) AnalyzePerformance() PerformanceAnalysis {
	mc.mu.RLock()
	defer mc.mu.RUnlock()

	analysis := PerformanceAnalysis{
		StartTime:    mc.startTime,
		EndTime:      time.Now(),
		MetricTrends: make(map[string]MetricTrend),
		Anomalies:    make([]Anomaly, 0),
	}

	// Analyze each metric
	for name := range mc.metrics {
		trend := mc.analyzeMetricTrend(name)
		analysis.MetricTrends[name] = trend

		// Check for anomalies
		if anomaly := mc.detectAnomaly(name, trend); anomaly != nil {
			analysis.Anomalies = append(analysis.Anomalies, *anomaly)
		}
	}

	return analysis
}

// Helper methods

func (mc *MetricsCollector) collectSystemMetrics() SystemMetrics {
	var stats runtime.MemStats
	runtime.ReadMemStats(&stats)

	return SystemMetrics{
		CPUUsage:    mc.calculateCPUUsage(),
		MemoryUsage: stats.Alloc,
		GCStats:     stats,
		Goroutines:  runtime.NumGoroutine(),
	}
}

func (mc *MetricsCollector) calculateCPUUsage() float64 {
	// Implementation would calculate actual CPU usage
	// This is a placeholder
	return 0.0
}

func (mc *MetricsCollector) analyzeMetricTrend(name string) MetricTrend {
	if len(mc.snapshots) < 2 {
		return MetricTrend{
			Direction: "stable",
			Rate:      0,
		}
	}

	latest := mc.snapshots[len(mc.snapshots)-1].Metrics[name]
	previous := mc.snapshots[len(mc.snapshots)-2].Metrics[name]
	change := latest - previous

	return MetricTrend{
		Direction: mc.determineTrendDirection(change),
		Rate:      change,
	}
}

func (mc *MetricsCollector) determineTrendDirection(change float64) string {
	switch {
	case change > 0:
		return "increasing"
	case change < 0:
		return "decreasing"
	default:
		return "stable"
	}
}

func (mc *MetricsCollector) detectAnomaly(name string, trend MetricTrend) *Anomaly {
	// Implementation would use more sophisticated anomaly detection
	// This is a simple threshold-based detection
	if abs(trend.Rate) > 0.5 { // Arbitrary threshold
		return &Anomaly{
			Metric:    name,
			Timestamp: time.Now(),
			Severity:  "warning",
			Message:   fmt.Sprintf("Rapid %s trend detected", trend.Direction),
		}
	}
	return nil
}

// Types for performance analysis

type PerformanceAnalysis struct {
	StartTime    time.Time
	EndTime      time.Time
	MetricTrends map[string]MetricTrend
	Anomalies    []Anomaly
}

type MetricTrend struct {
	Direction string
	Rate      float64
}

type Anomaly struct {
	Metric    string
	Timestamp time.Time
	Severity  string
	Message   string
}

// Utility functions

func abs(x float64) float64 {
	if x < 0 {
		return -x
	}
	return x
}
