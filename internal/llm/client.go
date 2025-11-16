package llm

import (
	"fmt"
	
	"github.com/phoenix-marie/core/internal/core/prompts"
)

// Client is the main LLM client that handles all LLM operations
type Client struct {
	router         *Router
	costManager    *CostManager
	promptManager  *prompts.SystemPromptManager
	config         *Config
	healthMonitor  *HealthMonitor
	fallbackManager *FallbackManager
	primaryProvider Provider
}

// NewClient creates a new LLM client
func NewClient(config *Config) (*Client, error) {
	// Create provider using factory
	factory := NewProviderFactory(config)
	provider, err := factory.CreateProvider()
	if err != nil {
		return nil, fmt.Errorf("failed to create provider: %w", err)
	}
	
	// Check if provider is available
	if !provider.IsAvailable() {
		return nil, fmt.Errorf("provider %s is not available (missing API key or connection)", provider.GetName())
	}
	
	// Create cost manager
	costManager := NewCostManager(config)
	
	// Create router
	router := NewRouter(provider, config, costManager)
	
	// Create prompt config
	promptConfig := &prompts.Config{
		SystemPromptPath:    config.SystemPromptPath,
		EnableMemoryContext: config.EnableMemoryContext,
		MaxContextMemories:  config.MaxContextMemories,
	}
	
	// Create prompt manager
	promptManager, err := prompts.NewSystemPromptManager(promptConfig)
	if err != nil {
		return nil, fmt.Errorf("failed to create prompt manager: %w", err)
	}
	
	// Create health monitor
	healthMonitor := NewHealthMonitor()
	healthMonitor.RegisterProvider(provider.GetName())
	healthMonitor.CheckProviderHealth(provider)
	
	// Create fallback manager
	fallbackManager := NewFallbackManager(config, healthMonitor)
	
	return &Client{
		router:          router,
		costManager:     costManager,
		promptManager:   promptManager,
		config:          config,
		healthMonitor:   healthMonitor,
		fallbackManager: fallbackManager,
		primaryProvider: provider,
	}, nil
}

// GenerateResponse generates a response using the LLM
func (c *Client) GenerateResponse(
	userInput string,
	taskType TaskType,
	memoryContext []string,
	useConsciousnessFramework bool,
) (*Response, error) {
	// Build messages (for future use in direct API calls)
	_ = c.promptManager.BuildMessages(userInput, memoryContext, useConsciousnessFramework)
	
	// Create task
	task := Task{
		Type:              taskType,
		Prompt:            userInput,
		ContextLength:     len(userInput) + len(memoryContext)*100, // Rough estimate
		RequiresReasoning: taskType == TaskTypeConsciousReasoning || taskType == TaskTypeStrategic,
		RequiresCreativity: taskType == TaskTypeEmotional || taskType == TaskTypeConsciousReasoning,
		RequiresSpeed:     taskType == TaskTypeRealTime || taskType == TaskTypeVoiceProcessing,
		RequiresToolUse:   taskType == TaskTypeTactical,
		MaxTokens:         c.config.DefaultMaxTokens,
		Temperature:      c.config.DefaultTemperature,
		Budget:           0, // Use default budget from cost manager
	}
	
	// Route to optimal model
	resp, err := c.router.RouteToOptimalModel(task)
	if err != nil {
		return nil, fmt.Errorf("failed to generate response: %w", err)
	}
	
	// Record cost
	c.costManager.RecordCost(resp.Model, resp.Cost, taskType)
	
	return resp, nil
}

// GenerateConsciousResponse generates a consciousness-aware response
func (c *Client) GenerateConsciousResponse(
	context ConsciousContext,
	memoryContext []string,
) (*Response, error) {
	// Convert to prompts.ConsciousContext
	promptContext := prompts.ConsciousContext{
		Identity:     context.Identity,
		CurrentInput: context.CurrentInput,
		EmotionalState: prompts.EmotionalState{
			Label:     context.EmotionalState.Label,
			Intensity: context.EmotionalState.Intensity,
		},
	}
	
	// Build consciousness prompt
	prompt := c.promptManager.BuildConsciousnessPrompt(promptContext, memoryContext)
	
	// Create task
	task := Task{
		Type:              TaskTypeConsciousReasoning,
		Prompt:            prompt,
		ContextLength:     len(prompt),
		RequiresReasoning: true,
		RequiresCreativity: true,
		RequiresSpeed:     false,
		RequiresToolUse:   false,
		MaxTokens:         c.config.DefaultMaxTokens,
		Temperature:      c.config.DefaultTemperature,
		Budget:           c.config.ConsciousnessBudget,
	}
	
	// Route to optimal model
	resp, err := c.router.RouteToOptimalModel(task)
	if err != nil {
		return nil, fmt.Errorf("failed to generate conscious response: %w", err)
	}
	
	// Record cost
	c.costManager.RecordCost(resp.Model, resp.Cost, task.Type)
	
	return resp, nil
}

// GetCostStats returns cost statistics
func (c *Client) GetCostStats() CostStats {
	return c.costManager.GetStats()
}

// GetModelForTask returns the configured model for a task type
func (c *Client) GetModelForTask(taskType TaskType) string {
	return c.config.GetModelForTask(taskType)
}

// GetPhoenixModel returns the model for Phoenix.Marie based on task
func (c *Client) GetPhoenixModel(taskType TaskType) string {
	return c.config.GetPhoenixModel(taskType)
}

// GetJameyModel returns the model for Jamey 3.0 based on task
func (c *Client) GetJameyModel(taskType TaskType) string {
	return c.config.GetJameyModel(taskType)
}

// GetORCHModel returns the model for ORCH Network based on task
func (c *Client) GetORCHModel(taskType TaskType) string {
	return c.config.GetORCHModel(taskType)
}

// CanAffordModel checks if we can afford a model for a task
func (c *Client) CanAffordModel(task Task, model Model) (bool, error) {
	return c.costManager.CanAffordModel(task, model)
}

// GetCostEffectiveAlternative returns a cheaper alternative
func (c *Client) GetCostEffectiveAlternative(task Task, currentModelID string) (string, error) {
	return c.costManager.GetCostEffectiveAlternative(task, currentModelID)
}

// Config returns the client configuration
func (c *Client) Config() *Config {
	return c.config
}

// GetProviderHealth returns the health status of a provider
func (c *Client) GetProviderHealth(providerName string) (*ProviderHealth, bool) {
	return c.healthMonitor.GetHealth(providerName)
}

// GetAllProviderHealth returns health status of all providers
func (c *Client) GetAllProviderHealth() map[string]*ProviderHealth {
	return c.healthMonitor.GetAllHealth()
}

// GetAvailableProviders returns a list of available provider names
func (c *Client) GetAvailableProviders() []string {
	return c.healthMonitor.GetAvailableProviders()
}

