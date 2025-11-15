package llm

import "time"

// Model represents an LLM model configuration
type Model struct {
	ID            string
	Name          string
	Provider      string
	ContextLength int
	InputPrice    float64  // Price per million input tokens
	OutputPrice   float64  // Price per million output tokens
	Capabilities  Capabilities
}

// Capabilities describes what a model can do
type Capabilities struct {
	Reasoning   bool
	Creativity   bool
	Speed        bool
	ToolUse      bool
	Multimodal   bool
	Multilingual bool
	Math         bool
}

// Task represents an LLM task request
type Task struct {
	Type            TaskType
	Prompt          string
	ContextLength   int
	RequiresReasoning bool
	RequiresCreativity bool
	RequiresSpeed     bool
	RequiresToolUse   bool
	MaxTokens       int
	Temperature     float64
	Budget          float64 // Maximum cost for this task
}

// TaskType represents the type of task
type TaskType string

const (
	TaskTypeConsciousReasoning TaskType = "conscious_reasoning"
	TaskTypeOperational        TaskType = "operational"
	TaskTypeRealTime          TaskType = "real_time"
	TaskTypeStrategic         TaskType = "strategic"
	TaskTypeTactical          TaskType = "tactical"
	TaskTypeAnalytical        TaskType = "analytical"
	TaskTypeEmotional         TaskType = "emotional"
	TaskTypeVoiceProcessing   TaskType = "voice_processing"
)

// Response, TokenUsage, and Message are defined in provider.go

// ConsciousContext provides context for consciousness-aware prompts
type ConsciousContext struct {
	Identity      string
	CurrentInput  string
	EmotionalState EmotionalState
	MemoryContext  []MemoryEvent
}

// EmotionalState represents current emotional state
type EmotionalState struct {
	Label     string
	Intensity int // 0-100
}

// MemoryEvent represents a memory event for context
type MemoryEvent struct {
	Summary string
	Time    time.Time
}

// ModelPerformance tracks model performance metrics
type ModelPerformance struct {
	Model            string
	ConsciousnessScore float64 // 0-100
	ResponseTime     time.Duration
	CostPerTask      float64
	Reliability      float64 // 0-1
	TasksCompleted   int
	TasksFailed      int
}

