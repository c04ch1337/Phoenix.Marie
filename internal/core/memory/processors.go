package memory

import (
	"encoding/json"
	"fmt"
	"sync"
	"time"
)

// ProcessorManager manages the specialized processors for each layer
type ProcessorManager struct {
	processors map[string]Processor
	mu         sync.RWMutex
}

// Processor interface defines the methods each specialized processor must implement
type Processor interface {
	Process(data any) (any, error)
	GetType() string
}

// NewProcessorManager creates a new processor manager with specialized processors
func NewProcessorManager() *ProcessorManager {
	pm := &ProcessorManager{
		processors: make(map[string]Processor),
	}

	// Initialize specialized processors for each layer
	pm.processors["sensory"] = &SensoryProcessor{}
	pm.processors["emotion"] = &EmotionProcessor{}
	pm.processors["logic"] = &LogicProcessor{}
	pm.processors["dream"] = &DreamProcessor{}
	pm.processors["eternal"] = &EternalProcessor{}

	return pm
}

// SensoryProcessor handles incoming sensory data
type SensoryProcessor struct{}

func (p *SensoryProcessor) Process(data any) (any, error) {
	processed := map[string]interface{}{
		"timestamp": time.Now().UnixNano(),
		"type":      "sensory",
		"raw":       data,
	}
	return processed, nil
}

func (p *SensoryProcessor) GetType() string {
	return "sensory"
}

// EmotionProcessor handles emotional data and state
type EmotionProcessor struct{}

func (p *EmotionProcessor) Process(data any) (any, error) {
	var emotionData map[string]interface{}

	switch v := data.(type) {
	case map[string]interface{}:
		emotionData = v
	case string:
		if err := json.Unmarshal([]byte(v), &emotionData); err != nil {
			emotionData = map[string]interface{}{"value": v}
		}
	default:
		emotionData = map[string]interface{}{"value": v}
	}

	emotionData["processed_at"] = time.Now().UnixNano()
	emotionData["type"] = "emotion"

	return emotionData, nil
}

func (p *EmotionProcessor) GetType() string {
	return "emotion"
}

// LogicProcessor handles logical operations and decision making
type LogicProcessor struct{}

func (p *LogicProcessor) Process(data any) (any, error) {
	processed := map[string]interface{}{
		"timestamp": time.Now().UnixNano(),
		"type":      "logic",
	}

	switch v := data.(type) {
	case map[string]interface{}:
		processed["data"] = v
		processed["complexity"] = len(v)
	case string:
		processed["data"] = v
		processed["complexity"] = len(v)
	default:
		processed["data"] = fmt.Sprintf("%v", v)
		processed["complexity"] = 1
	}

	return processed, nil
}

func (p *LogicProcessor) GetType() string {
	return "logic"
}

// DreamProcessor handles dream state processing and pattern recognition
type DreamProcessor struct{}

func (p *DreamProcessor) Process(data any) (any, error) {
	dreamState := map[string]interface{}{
		"timestamp": time.Now().UnixNano(),
		"type":      "dream",
		"source":    data,
		"patterns":  make([]string, 0),
	}

	// Extract patterns from data
	switch v := data.(type) {
	case map[string]interface{}:
		for key, val := range v {
			pattern := fmt.Sprintf("%s:%v", key, val)
			dreamState["patterns"] = append(dreamState["patterns"].([]string), pattern)
		}
	case string:
		dreamState["patterns"] = append(dreamState["patterns"].([]string), v)
	}

	return dreamState, nil
}

func (p *DreamProcessor) GetType() string {
	return "dream"
}

// EternalProcessor handles long-term memory storage and retrieval
type EternalProcessor struct{}

func (p *EternalProcessor) Process(data any) (any, error) {
	eternal := map[string]interface{}{
		"created_at": time.Now().UnixNano(),
		"type":       "eternal",
		"permanent":  true,
	}

	switch v := data.(type) {
	case map[string]interface{}:
		eternal["data"] = v
		if importance, ok := v["importance"]; ok {
			eternal["importance"] = importance
		} else {
			eternal["importance"] = 1
		}
	default:
		eternal["data"] = v
		eternal["importance"] = 1
	}

	return eternal, nil
}

func (p *EternalProcessor) GetType() string {
	return "eternal"
}

// ProcessData processes data using the appropriate processor for the given layer
func (pm *ProcessorManager) ProcessData(layer string, data any) (any, error) {
	pm.mu.RLock()
	processor, exists := pm.processors[layer]
	pm.mu.RUnlock()

	if !exists {
		return nil, fmt.Errorf("no processor found for layer: %s", layer)
	}

	return processor.Process(data)
}

// RegisterProcessor registers a new processor for a layer
func (pm *ProcessorManager) RegisterProcessor(layer string, processor Processor) error {
	if processor == nil {
		return fmt.Errorf("cannot register nil processor")
	}

	pm.mu.Lock()
	defer pm.mu.Unlock()

	if _, exists := pm.processors[layer]; exists {
		return fmt.Errorf("processor already registered for layer: %s", layer)
	}

	pm.processors[layer] = processor
	return nil
}
