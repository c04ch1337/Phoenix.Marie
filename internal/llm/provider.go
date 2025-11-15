package llm

import (
	"fmt"
	"time"
)

// Provider defines the interface for LLM providers
type Provider interface {
	// Call makes a request to the LLM API
	Call(modelID string, messages []Message, maxTokens int, temperature float64) (*Response, error)
	
	// CallWithRetry makes a request with retry logic
	CallWithRetry(modelID string, messages []Message, maxTokens int, temperature float64) (*Response, error)
	
	// GetName returns the provider name
	GetName() string
	
	// IsAvailable checks if the provider is available
	IsAvailable() bool
}

// ProviderFactory creates providers based on configuration
type ProviderFactory struct {
	config *Config
}

// NewProviderFactory creates a new provider factory
func NewProviderFactory(config *Config) *ProviderFactory {
	return &ProviderFactory{config: config}
}

// CreateProvider creates a provider based on the configured provider type
func (pf *ProviderFactory) CreateProvider() (Provider, error) {
	switch pf.config.Provider {
	case "openrouter":
		return NewOpenRouterClient(pf.config), nil
	case "openai":
		return NewOpenAIClient(pf.config), nil
	case "anthropic":
		return NewAnthropicClient(pf.config), nil
	case "gemini":
		return NewGeminiClient(pf.config), nil
	case "grok":
		return NewGrokClient(pf.config), nil
	case "ollama":
		return NewOllamaClient(pf.config), nil
	case "lmstudio":
		return NewLMStudioClient(pf.config), nil
	default:
		return nil, fmt.Errorf("unknown provider: %s", pf.config.Provider)
	}
}

// Message represents a chat message
type Message struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

// Response represents an LLM response
type Response struct {
	Content      string
	Model        string
	TokensUsed   TokenUsage
	Cost         float64
	ResponseTime time.Duration
	FinishReason string
}

// TokenUsage tracks token consumption
type TokenUsage struct {
	PromptTokens     int
	CompletionTokens int
	TotalTokens      int
}

