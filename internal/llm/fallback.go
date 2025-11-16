package llm

import (
	"fmt"
)

// FallbackManager manages provider fallback logic
type FallbackManager struct {
	config        *Config
	healthMonitor *HealthMonitor
	fallbackOrder []string
}

// NewFallbackManager creates a new fallback manager
func NewFallbackManager(config *Config, healthMonitor *HealthMonitor) *FallbackManager {
	// Define fallback order: try primary provider first, then alternatives
	fallbackOrder := []string{
		config.Provider, // Primary provider
	}
	
	// Add alternative providers in order of preference
	// OpenRouter is preferred for model variety, then direct providers, then local
	alternatives := []string{"openrouter", "openai", "anthropic", "gemini", "grok", "ollama", "lmstudio"}
	for _, alt := range alternatives {
		if alt != config.Provider {
			fallbackOrder = append(fallbackOrder, alt)
		}
	}
	
	return &FallbackManager{
		config:        config,
		healthMonitor: healthMonitor,
		fallbackOrder: fallbackOrder,
	}
}

// GetNextProvider returns the next available provider in the fallback chain
func (fm *FallbackManager) GetNextProvider(currentProvider string) (Provider, error) {
	// Find current provider index
	currentIndex := -1
	for i, p := range fm.fallbackOrder {
		if p == currentProvider {
			currentIndex = i
			break
		}
	}
	
	// Try next providers in fallback order
	for i := currentIndex + 1; i < len(fm.fallbackOrder); i++ {
		providerName := fm.fallbackOrder[i]
		
		// Check if provider is healthy
		if health, exists := fm.healthMonitor.GetHealth(providerName); exists {
			if !health.IsAvailable {
				continue // Skip unavailable providers
			}
		}
		
		// Create provider
		factory := NewProviderFactory(fm.config)
		// Temporarily switch provider to test availability
		originalProvider := fm.config.Provider
		fm.config.Provider = providerName
		
		provider, err := factory.CreateProvider()
		fm.config.Provider = originalProvider // Restore
		
		if err != nil {
			continue
		}
		
		if provider.IsAvailable() {
			return provider, nil
		}
	}
	
	return nil, fmt.Errorf("no available fallback providers")
}

// TryWithFallback attempts a request with the primary provider, falling back if needed
func (fm *FallbackManager) TryWithFallback(
	primaryProvider Provider,
	modelID string,
	messages []Message,
	maxTokens int,
	temperature float64,
) (*Response, error) {
	// Try primary provider first
	resp, err := primaryProvider.CallWithRetry(modelID, messages, maxTokens, temperature)
	if err == nil {
		// Record success
		fm.healthMonitor.UpdateHealth(
			primaryProvider.GetName(),
			true,
			resp.ResponseTime,
		)
		return resp, nil
	}
	
	// Record failure
	fm.healthMonitor.UpdateHealth(
		primaryProvider.GetName(),
		false,
		0,
	)
	
	// Try fallback providers
	fallbackProvider, fallbackErr := fm.GetNextProvider(primaryProvider.GetName())
	if fallbackErr != nil {
		return nil, fmt.Errorf("primary provider failed and no fallback available: %w", err)
	}
	
	// Try fallback
	resp, fallbackErr = fallbackProvider.CallWithRetry(modelID, messages, maxTokens, temperature)
	if fallbackErr == nil {
		// Record fallback success
		fm.healthMonitor.UpdateHealth(
			fallbackProvider.GetName(),
			true,
			resp.ResponseTime,
		)
		return resp, nil
	}
	
	// Record fallback failure
	fm.healthMonitor.UpdateHealth(
		fallbackProvider.GetName(),
		false,
		0,
	)
	
	return nil, fmt.Errorf("all providers failed: primary=%v, fallback=%v", err, fallbackErr)
}

// GetFallbackChain returns the current fallback chain
func (fm *FallbackManager) GetFallbackChain() []string {
	return fm.fallbackOrder
}

// UpdateFallbackOrder updates the fallback order based on provider health
func (fm *FallbackManager) UpdateFallbackOrder() {
	// Reorder based on health: healthy providers first
	allHealth := fm.healthMonitor.GetAllHealth()
	
	// Sort by availability and success rate
	type providerScore struct {
		name  string
		score float64
	}
	
	var scores []providerScore
	for _, name := range fm.fallbackOrder {
		health, exists := allHealth[name]
		if !exists {
			scores = append(scores, providerScore{name: name, score: 0.5})
			continue
		}
		
		score := 0.0
		if health.IsAvailable {
			score += 1.0
		}
		if health.TotalRequests > 0 {
			successRate := float64(health.SuccessfulRequests) / float64(health.TotalRequests)
			score += successRate
		}
		
		scores = append(scores, providerScore{name: name, score: score})
	}
	
	// Sort by score (simple bubble sort)
	for i := 0; i < len(scores)-1; i++ {
		for j := i + 1; j < len(scores); j++ {
			if scores[i].score < scores[j].score {
				scores[i], scores[j] = scores[j], scores[i]
			}
		}
	}
	
	// Update fallback order
	newOrder := make([]string, len(scores))
	for i, s := range scores {
		newOrder[i] = s.name
	}
	
	fm.fallbackOrder = newOrder
}

