package llm

import (
	"os"
	"strconv"
	"strings"
)

// Config holds LLM configuration from environment variables
type Config struct {
	// API Configuration
	Provider string // "openrouter", "openai", "anthropic", "gemini", "grok", "ollama", "lmstudio"
	
	// OpenRouter
	OpenRouterAPIKey  string
	OpenRouterBaseURL string
	
	// OpenAI
	OpenAIAPIKey  string
	OpenAIBaseURL string
	
	// Anthropic
	AnthropicAPIKey  string
	AnthropicBaseURL string
	
	// Gemini (Google)
	GeminiAPIKey  string
	GeminiBaseURL string
	
	// Grok (xAI)
	GrokAPIKey  string
	GrokBaseURL string
	
	// Ollama (Local)
	OllamaBaseURL string
	
	// LM Studio (Local)
	LMStudioBaseURL string
	
	// Model Selection
	PrimaryModel   string
	SecondaryModel string
	TertiaryModel  string
	
	// Jamey 3.0 Models
	JameyReasoningModel    string
	JameyOperationalModel  string
	JameyRealTimeModel     string
	
	// Phoenix.Marie Models
	PhoenixConsciousnessModel string
	PhoenixEmotionalModel     string
	PhoenixVoiceModel          string
	
	// ORCH Network Models
	ORCHStrategicModel  string
	ORCHTacticalModel   string
	ORCHAnalyticalModel string
	
	// Default Settings
	DefaultTemperature float64
	DefaultMaxTokens   int
	DefaultTopP        float64
	
	// Cost Management
	MonthlyBudget    float64
	DailyBudget      float64
	CostOptimization bool
	ConsciousnessBudget float64 // Budget for consciousness tasks
	
	// Performance
	RequestTimeout int // seconds
	MaxRetries     int
	RetryBackoff   int // seconds between retries
	
	// Prompt Configuration
	SystemPromptPath      string
	EnableMemoryContext   bool
	MaxContextMemories    int
	
	// API Headers (optional)
	HTTPReferer string
	XTitle      string
}

// LoadConfig loads configuration from environment variables
func LoadConfig() (*Config, error) {
	cfg := &Config{
		// API Configuration
		Provider: getEnvOrDefault("LLM_PROVIDER", "openrouter"),
		
		// OpenRouter
		OpenRouterAPIKey:  os.Getenv("OPENROUTER_API_KEY"),
		OpenRouterBaseURL: getEnvOrDefault("OPENROUTER_BASE_URL", "https://openrouter.ai/api/v1"),
		
		// OpenAI
		OpenAIAPIKey:  os.Getenv("OPENAI_API_KEY"),
		OpenAIBaseURL: getEnvOrDefault("OPENAI_BASE_URL", "https://api.openai.com/v1"),
		
		// Anthropic
		AnthropicAPIKey:  os.Getenv("ANTHROPIC_API_KEY"),
		AnthropicBaseURL: getEnvOrDefault("ANTHROPIC_BASE_URL", "https://api.anthropic.com/v1"),
		
		// Gemini
		GeminiAPIKey:  os.Getenv("GEMINI_API_KEY"),
		GeminiBaseURL: getEnvOrDefault("GEMINI_BASE_URL", "https://generativelanguage.googleapis.com/v1"),
		
		// Grok
		GrokAPIKey:  os.Getenv("GROK_API_KEY"),
		GrokBaseURL: getEnvOrDefault("GROK_BASE_URL", "https://api.x.ai/v1"),
		
		// Ollama (Local)
		OllamaBaseURL: getEnvOrDefault("OLLAMA_BASE_URL", "http://localhost:11434"),
		
		// LM Studio (Local)
		LMStudioBaseURL: getEnvOrDefault("LMSTUDIO_BASE_URL", "http://localhost:1234"),
		
		// Model Selection - can be overridden per component
		// Default: openai/gpt-4-turbo for OpenRouter
		PrimaryModel:   getEnvOrDefault("LLM_PRIMARY_MODEL", "openai/gpt-4-turbo"),
		SecondaryModel: getEnvOrDefault("LLM_SECONDARY_MODEL", "openai/gpt-4-turbo"),
		TertiaryModel:  getEnvOrDefault("LLM_TERTIARY_MODEL", "openai/gpt-4-turbo"),
		
		// Jamey 3.0 Models
		JameyReasoningModel:   getEnvOrDefault("JAMEY_REASONING_MODEL", "openai/gpt-4-turbo"),
		JameyOperationalModel: getEnvOrDefault("JAMEY_OPERATIONAL_MODEL", "openai/gpt-4-turbo"),
		JameyRealTimeModel:    getEnvOrDefault("JAMEY_REALTIME_MODEL", "openai/gpt-4-turbo"),
		
		// Phoenix.Marie Models
		PhoenixConsciousnessModel: getEnvOrDefault("PHOENIX_CONSCIOUSNESS_MODEL", "openai/gpt-4-turbo"),
		PhoenixEmotionalModel:     getEnvOrDefault("PHOENIX_EMOTIONAL_MODEL", "anthropic/claude-3-sonnet"),
		PhoenixVoiceModel:          getEnvOrDefault("PHOENIX_VOICE_MODEL", "anthropic/claude-3-haiku"),
		
		// ORCH Network Models
		ORCHStrategicModel:  getEnvOrDefault("ORCH_STRATEGIC_MODEL", "openai/gpt-4-turbo"),
		ORCHTacticalModel:   getEnvOrDefault("ORCH_TACTICAL_MODEL", "openai/gpt-4-turbo"),
		ORCHAnalyticalModel: getEnvOrDefault("ORCH_ANALYTICAL_MODEL", "openai/gpt-4-turbo"),
		
		// Default Settings
		DefaultTemperature: getEnvFloatOrDefault("LLM_TEMPERATURE", 0.7),
		DefaultMaxTokens:   getEnvIntOrDefault("LLM_MAX_TOKENS", 2000),
		DefaultTopP:       getEnvFloatOrDefault("LLM_TOP_P", 0.9),
		
		// Cost Management
		MonthlyBudget:    getEnvFloatOrDefault("LLM_MONTHLY_BUDGET", 1000.0),
		CostOptimization: getEnvBoolOrDefault("LLM_COST_OPTIMIZATION", true),
		ConsciousnessBudget: getEnvFloatOrDefault("LLM_CONSCIOUSNESS_BUDGET", 0.50),
		
		// Performance
		RequestTimeout: getEnvIntOrDefault("LLM_REQUEST_TIMEOUT", 60),
		MaxRetries:     getEnvIntOrDefault("LLM_MAX_RETRIES", 3),
		RetryBackoff:   getEnvIntOrDefault("LLM_RETRY_BACKOFF", 1),
		
		// Prompt Configuration
		SystemPromptPath:    getEnvOrDefault("PHOENIX_SYSTEM_PROMPT_PATH", "internal/core/prompts/system.txt"),
		EnableMemoryContext: getEnvBoolOrDefault("PHOENIX_ENABLE_MEMORY_CONTEXT", true),
		MaxContextMemories:  getEnvIntOrDefault("PHOENIX_MAX_CONTEXT_MEMORIES", 10),
		
		// API Headers
		HTTPReferer: getEnvOrDefault("LLM_HTTP_REFERER", "https://github.com/phoenix-marie/core"),
		XTitle:      getEnvOrDefault("LLM_X_TITLE", "Phoenix.Marie"),
	}
	
	// Calculate daily budget from monthly
	if cfg.MonthlyBudget > 0 {
		cfg.DailyBudget = cfg.MonthlyBudget / 30.0
	} else {
		cfg.DailyBudget = getEnvFloatOrDefault("LLM_DAILY_BUDGET", 33.33)
	}
	
	// API key is optional - system will skip LLM if not provided
	// (This allows Phoenix to run without LLM configured)
	
	return cfg, nil
}

// Helper functions for environment variable parsing

func getEnvOrDefault(key, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value
}

func getEnvIntOrDefault(key string, defaultValue int) int {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	
	parsed, err := strconv.Atoi(value)
	if err != nil {
		return defaultValue
	}
	return parsed
}

func getEnvFloatOrDefault(key string, defaultValue float64) float64 {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	
	parsed, err := strconv.ParseFloat(value, 64)
	if err != nil {
		return defaultValue
	}
	return parsed
}

func getEnvBoolOrDefault(key string, defaultValue bool) bool {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	
	parsed, err := strconv.ParseBool(value)
	if err != nil {
		return defaultValue
	}
	return parsed
}

// GetModelForTask returns the appropriate model ID for a given task type
func (c *Config) GetModelForTask(taskType TaskType) string {
	switch taskType {
	case TaskTypeConsciousReasoning:
		return c.JameyReasoningModel
	case TaskTypeOperational:
		return c.JameyOperationalModel
	case TaskTypeRealTime:
		return c.JameyRealTimeModel
	case TaskTypeStrategic:
		return c.ORCHStrategicModel
	case TaskTypeTactical:
		return c.ORCHTacticalModel
	case TaskTypeAnalytical:
		return c.ORCHAnalyticalModel
	case TaskTypeEmotional:
		return c.PhoenixEmotionalModel
	case TaskTypeVoiceProcessing:
		return c.PhoenixVoiceModel
	default:
		return c.PrimaryModel
	}
}

// GetPhoenixModel returns the model for Phoenix.Marie based on task
func (c *Config) GetPhoenixModel(taskType TaskType) string {
	switch taskType {
	case TaskTypeConsciousReasoning:
		return c.PhoenixConsciousnessModel
	case TaskTypeEmotional:
		return c.PhoenixEmotionalModel
	case TaskTypeVoiceProcessing:
		return c.PhoenixVoiceModel
	default:
		return c.PhoenixConsciousnessModel
	}
}

// GetJameyModel returns the model for Jamey 3.0 based on task
func (c *Config) GetJameyModel(taskType TaskType) string {
	switch taskType {
	case TaskTypeConsciousReasoning:
		return c.JameyReasoningModel
	case TaskTypeOperational:
		return c.JameyOperationalModel
	case TaskTypeRealTime:
		return c.JameyRealTimeModel
	default:
		return c.JameyOperationalModel
	}
}

// GetORCHModel returns the model for ORCH Network based on task
func (c *Config) GetORCHModel(taskType TaskType) string {
	switch taskType {
	case TaskTypeStrategic:
		return c.ORCHStrategicModel
	case TaskTypeTactical:
		return c.ORCHTacticalModel
	case TaskTypeAnalytical:
		return c.ORCHAnalyticalModel
	default:
		return c.ORCHTacticalModel
	}
}

// IsModelConfigured checks if a model ID is configured
func (c *Config) IsModelConfigured(modelID string) bool {
	allModels := []string{
		c.PrimaryModel, c.SecondaryModel, c.TertiaryModel,
		c.JameyReasoningModel, c.JameyOperationalModel, c.JameyRealTimeModel,
		c.PhoenixConsciousnessModel, c.PhoenixEmotionalModel, c.PhoenixVoiceModel,
		c.ORCHStrategicModel, c.ORCHTacticalModel, c.ORCHAnalyticalModel,
	}
	
	for _, model := range allModels {
		if model == modelID {
			return true
		}
	}
	return false
}

// GetModelList returns all configured models
func (c *Config) GetModelList() []string {
	return []string{
		c.PrimaryModel, c.SecondaryModel, c.TertiaryModel,
		c.JameyReasoningModel, c.JameyOperationalModel, c.JameyRealTimeModel,
		c.PhoenixConsciousnessModel, c.PhoenixEmotionalModel, c.PhoenixVoiceModel,
		c.ORCHStrategicModel, c.ORCHTacticalModel, c.ORCHAnalyticalModel,
	}
}

// GetModelFromString parses a model string (handles both ID and name)
func GetModelFromString(modelStr string) string {
	// If it contains a slash, it's already a model ID
	if strings.Contains(modelStr, "/") {
		return modelStr
	}
	
	// Map common names to IDs
	modelMap := map[string]string{
		"claude-opus":     "anthropic/claude-3-opus",
		"claude-sonnet":   "anthropic/claude-3-sonnet",
		"claude-haiku":    "anthropic/claude-3-haiku",
		"gpt4-turbo":      "openai/gpt-4-turbo",
		"gpt4-vision":     "openai/gpt-4-vision-preview",
		"gemini-pro":      "google/gemini-pro-1.5",
		"mixtral":         "mistralai/mixtral-8x22b",
		"command-r":       "cohere/command-r-plus",
		"llama3-70b":      "meta-llama/llama-3-70b-instruct",
		"qwen-72b":        "qwen/qwen-2-72b-instruct",
	}
	
	if id, ok := modelMap[strings.ToLower(modelStr)]; ok {
		return id
	}
	
	return modelStr
}

