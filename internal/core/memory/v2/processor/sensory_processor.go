package processor

import (
	"fmt"
	"time"
)

// SensoryData represents input from sensory systems
type SensoryData struct {
	Type      string                 // Type of sensory input (visual, auditory, etc.)
	Value     interface{}            // The actual sensory data
	Metadata  map[string]interface{} // Additional context about the data
	Timestamp time.Time              // When the data was captured
}

// SensoryProcessor specializes in processing sensory input data
type SensoryProcessor struct {
	BaseProcessor
	filters     []SensoryFilter
	aggregators []DataAggregator
}

// SensoryFilter defines a filter function for sensory data
type SensoryFilter func(SensoryData) (SensoryData, error)

// DataAggregator defines an aggregation function for processed data
type DataAggregator func([]ProcessedData) (ProcessedData, error)

// NewSensoryProcessor creates a new sensory processor instance
func NewSensoryProcessor() *SensoryProcessor {
	return &SensoryProcessor{
		filters:     make([]SensoryFilter, 0),
		aggregators: make([]DataAggregator, 0),
	}
}

// AddFilter adds a new filter to the processing pipeline
func (sp *SensoryProcessor) AddFilter(filter SensoryFilter) {
	sp.filters = append(sp.filters, filter)
}

// AddAggregator adds a new aggregator to the processing pipeline
func (sp *SensoryProcessor) AddAggregator(aggregator DataAggregator) {
	sp.aggregators = append(sp.aggregators, aggregator)
}

// Process implements specialized processing for sensory data
func (sp *SensoryProcessor) Process(data interface{}) (ProcessedData, error) {
	startTime := time.Now()

	// Type assertion
	sensoryData, ok := data.(SensoryData)
	if !ok {
		return ProcessedData{}, fmt.Errorf("invalid data type: expected SensoryData")
	}

	// Apply filters
	filtered, err := sp.applyFilters(sensoryData)
	if err != nil {
		return ProcessedData{}, fmt.Errorf("filter processing failed: %w", err)
	}

	// Create processed data
	processed := ProcessedData{
		Data: filtered.Value,
		Metadata: map[string]interface{}{
			"type":      filtered.Type,
			"timestamp": filtered.Timestamp,
			"filtered":  true,
		},
		Timestamp:  time.Now(),
		Confidence: sp.calculateConfidence(filtered),
	}

	// Merge original metadata
	for k, v := range filtered.Metadata {
		processed.Metadata[k] = v
	}

	// Update metrics
	sp.metrics.ProcessingTime += time.Since(startTime)
	sp.state.ProcessedCount++

	return processed, nil
}

// BatchProcess handles multiple sensory inputs at once
func (sp *SensoryProcessor) BatchProcess(data []interface{}) ([]ProcessedData, error) {
	results := make([]ProcessedData, 0, len(data))

	for _, item := range data {
		processed, err := sp.Process(item)
		if err != nil {
			return nil, fmt.Errorf("batch processing failed: %w", err)
		}
		results = append(results, processed)
	}

	// Apply aggregators if configured
	if len(sp.aggregators) > 0 {
		aggregated, err := sp.applyAggregators(results)
		if err != nil {
			return nil, fmt.Errorf("aggregation failed: %w", err)
		}
		return []ProcessedData{aggregated}, nil
	}

	return results, nil
}

// Helper methods

func (sp *SensoryProcessor) applyFilters(data SensoryData) (SensoryData, error) {
	current := data

	for _, filter := range sp.filters {
		var err error
		current, err = filter(current)
		if err != nil {
			return SensoryData{}, fmt.Errorf("filter application failed: %w", err)
		}
	}

	return current, nil
}

func (sp *SensoryProcessor) applyAggregators(data []ProcessedData) (ProcessedData, error) {
	current := data

	for _, aggregator := range sp.aggregators {
		aggregated, err := aggregator(current)
		if err != nil {
			return ProcessedData{}, fmt.Errorf("aggregation failed: %w", err)
		}
		current = []ProcessedData{aggregated}
	}

	return current[0], nil
}

func (sp *SensoryProcessor) calculateConfidence(data SensoryData) float64 {
	// Basic confidence calculation
	// This could be enhanced with more sophisticated algorithms
	baseConfidence := 0.8 // Base confidence level

	// Adjust based on metadata if available
	if signalStrength, ok := data.Metadata["signal_strength"].(float64); ok {
		baseConfidence *= signalStrength
	}

	// Adjust based on processing time
	if time.Since(data.Timestamp) > time.Second {
		baseConfidence *= 0.9 // Slight reduction for older data
	}

	return baseConfidence
}

// Common filter implementations

// NoiseFilter reduces noise in sensory data
func NoiseFilter(threshold float64) SensoryFilter {
	return func(data SensoryData) (SensoryData, error) {
		// Implementation would depend on the type of sensory data
		// This is a placeholder for demonstration
		if val, ok := data.Value.(float64); ok {
			if val < threshold {
				data.Value = 0.0
			}
		}
		return data, nil
	}
}

// TimeFilter filters out old data
func TimeFilter(maxAge time.Duration) SensoryFilter {
	return func(data SensoryData) (SensoryData, error) {
		if time.Since(data.Timestamp) > maxAge {
			return SensoryData{}, fmt.Errorf("data too old")
		}
		return data, nil
	}
}

// Common aggregator implementations

// AverageAggregator combines multiple data points by averaging
func AverageAggregator(data []ProcessedData) (ProcessedData, error) {
	if len(data) == 0 {
		return ProcessedData{}, fmt.Errorf("no data to aggregate")
	}

	var sum float64
	for _, d := range data {
		if val, ok := d.Data.(float64); ok {
			sum += val
		}
	}

	return ProcessedData{
		Data: sum / float64(len(data)),
		Metadata: map[string]interface{}{
			"aggregation_type": "average",
			"sample_count":     len(data),
		},
		Timestamp:  time.Now(),
		Confidence: calculateAggregateConfidence(data),
	}, nil
}

func calculateAggregateConfidence(data []ProcessedData) float64 {
	if len(data) == 0 {
		return 0.0
	}

	var sumConfidence float64
	for _, d := range data {
		sumConfidence += d.Confidence
	}

	return sumConfidence / float64(len(data))
}
