package processor

import (
	"fmt"
	"math"
	"time"
)

// EmotionalData represents emotional input data
type EmotionalData struct {
	Type      string                 // Type of emotion (joy, fear, etc.)
	Intensity float64                // Intensity level (0.0 to 1.0)
	Context   map[string]interface{} // Contextual information
	Source    string                 // Source of emotional data
	Timestamp time.Time              // When the emotion was detected
}

// EmotionProcessor specializes in processing emotional data
type EmotionProcessor struct {
	BaseProcessor
	intensityThreshold float64
	decayRate          float64
	contextWeight      float64
	emotionHistory     []EmotionalData
	historySize        int
}

// NewEmotionProcessor creates a new emotion processor instance
func NewEmotionProcessor(config EmotionProcessorConfig) *EmotionProcessor {
	return &EmotionProcessor{
		intensityThreshold: config.IntensityThreshold,
		decayRate:          config.DecayRate,
		contextWeight:      config.ContextWeight,
		emotionHistory:     make([]EmotionalData, 0),
		historySize:        config.HistorySize,
	}
}

// EmotionProcessorConfig contains configuration for emotion processing
type EmotionProcessorConfig struct {
	IntensityThreshold float64 // Minimum intensity to process
	DecayRate          float64 // Rate at which emotions decay
	ContextWeight      float64 // Weight given to contextual information
	HistorySize        int     // Size of emotion history to maintain
}

// Process implements specialized processing for emotional data
func (ep *EmotionProcessor) Process(data interface{}) (ProcessedData, error) {
	startTime := time.Now()

	// Type assertion
	emotionalData, ok := data.(EmotionalData)
	if !ok {
		return ProcessedData{}, fmt.Errorf("invalid data type: expected EmotionalData")
	}

	// Validate intensity
	if emotionalData.Intensity < 0 || emotionalData.Intensity > 1.0 {
		return ProcessedData{}, fmt.Errorf("invalid intensity: must be between 0 and 1")
	}

	// Apply intensity threshold
	if emotionalData.Intensity < ep.intensityThreshold {
		return ProcessedData{
			Data: emotionalData,
			Metadata: map[string]interface{}{
				"filtered_out": true,
				"reason":       "below_threshold",
			},
			Timestamp:  time.Now(),
			Confidence: 1.0,
		}, nil
	}

	// Process emotion with context and history
	processedEmotion := ep.processEmotionWithContext(emotionalData)

	// Update history
	ep.updateEmotionHistory(emotionalData)

	// Create processed data
	processed := ProcessedData{
		Data: processedEmotion,
		Metadata: map[string]interface{}{
			"original_type":      emotionalData.Type,
			"original_intensity": emotionalData.Intensity,
			"processed_time":     time.Now(),
			"history_size":       len(ep.emotionHistory),
		},
		Timestamp:  time.Now(),
		Confidence: ep.calculateConfidence(processedEmotion),
	}

	// Update metrics
	ep.metrics.ProcessingTime += time.Since(startTime)
	ep.state.ProcessedCount++

	return processed, nil
}

// Helper methods

func (ep *EmotionProcessor) processEmotionWithContext(data EmotionalData) EmotionalData {
	// Apply temporal decay
	adjustedIntensity := ep.applyTemporalDecay(data)

	// Consider context
	if contextualIntensity, ok := data.Context["related_intensity"].(float64); ok {
		adjustedIntensity = (adjustedIntensity * (1 - ep.contextWeight)) +
			(contextualIntensity * ep.contextWeight)
	}

	// Create processed emotion
	processed := EmotionalData{
		Type:      data.Type,
		Intensity: adjustedIntensity,
		Context:   data.Context,
		Source:    data.Source,
		Timestamp: time.Now(),
	}

	return processed
}

func (ep *EmotionProcessor) applyTemporalDecay(data EmotionalData) float64 {
	timeDiff := time.Since(data.Timestamp).Seconds()
	decayFactor := math.Exp(-ep.decayRate * timeDiff)
	return data.Intensity * decayFactor
}

func (ep *EmotionProcessor) updateEmotionHistory(data EmotionalData) {
	ep.emotionHistory = append(ep.emotionHistory, data)

	// Maintain history size
	for len(ep.emotionHistory) > ep.historySize {
		ep.emotionHistory = ep.emotionHistory[1:]
	}
}

func (ep *EmotionProcessor) calculateConfidence(data EmotionalData) float64 {
	// Base confidence starts high
	confidence := 0.9

	// Adjust based on history
	if len(ep.emotionHistory) > 0 {
		// Check for emotional consistency
		consistencyScore := ep.calculateConsistencyScore(data)
		confidence *= consistencyScore
	}

	// Adjust based on context completeness
	if len(data.Context) > 0 {
		contextScore := 0.1 * float64(len(data.Context))
		confidence = math.Min(confidence+contextScore, 1.0)
	}

	return confidence
}

func (ep *EmotionProcessor) calculateConsistencyScore(data EmotionalData) float64 {
	if len(ep.emotionHistory) == 0 {
		return 1.0
	}

	var consistentCount float64
	for _, historical := range ep.emotionHistory {
		// Check type consistency
		if historical.Type == data.Type {
			// Check intensity variation
			intensityDiff := math.Abs(historical.Intensity - data.Intensity)
			if intensityDiff < 0.3 { // Allow some variation
				consistentCount++
			}
		}
	}

	return consistentCount / float64(len(ep.emotionHistory))
}

// BatchProcess handles multiple emotional inputs at once
func (ep *EmotionProcessor) BatchProcess(data []interface{}) ([]ProcessedData, error) {
	results := make([]ProcessedData, 0, len(data))

	for _, item := range data {
		processed, err := ep.Process(item)
		if err != nil {
			return nil, fmt.Errorf("batch processing failed: %w", err)
		}
		results = append(results, processed)
	}

	return results, nil
}

// EmotionAnalyzer provides analysis of emotional patterns
type EmotionAnalyzer struct {
	patterns map[string]float64
}

// AnalyzeEmotions analyzes a set of emotional data for patterns
func (ea *EmotionAnalyzer) AnalyzeEmotions(data []EmotionalData) map[string]interface{} {
	analysis := make(map[string]interface{})

	// Calculate emotion distribution
	distribution := make(map[string]int)
	intensities := make(map[string]float64)

	for _, emotion := range data {
		distribution[emotion.Type]++
		intensities[emotion.Type] += emotion.Intensity
	}

	// Calculate averages and dominant emotion
	dominantEmotion := ""
	maxCount := 0

	for emotionType, count := range distribution {
		if count > maxCount {
			maxCount = count
			dominantEmotion = emotionType
		}
		intensities[emotionType] /= float64(count)
	}

	analysis["distribution"] = distribution
	analysis["average_intensities"] = intensities
	analysis["dominant_emotion"] = dominantEmotion
	analysis["sample_size"] = len(data)

	return analysis
}
