package llm

import (
	"fmt"
	"sync"
	"time"
)

// ProviderHealth tracks the health status of a provider
type ProviderHealth struct {
	ProviderName    string
	IsAvailable     bool
	LastChecked     time.Time
	LastSuccess     time.Time
	LastFailure     time.Time
	ConsecutiveFailures int
	TotalRequests   int64
	SuccessfulRequests int64
	FailedRequests  int64
	AverageResponseTime time.Duration
	mu              sync.RWMutex
}

// HealthMonitor monitors the health of all providers
type HealthMonitor struct {
	providers map[string]*ProviderHealth
	mu        sync.RWMutex
}

// NewHealthMonitor creates a new health monitor
func NewHealthMonitor() *HealthMonitor {
	return &HealthMonitor{
		providers: make(map[string]*ProviderHealth),
	}
}

// RegisterProvider registers a provider for health monitoring
func (hm *HealthMonitor) RegisterProvider(providerName string) {
	hm.mu.Lock()
	defer hm.mu.Unlock()
	
	if _, exists := hm.providers[providerName]; !exists {
		hm.providers[providerName] = &ProviderHealth{
			ProviderName: providerName,
			IsAvailable: false,
		}
	}
}

// UpdateHealth updates the health status of a provider
func (hm *HealthMonitor) UpdateHealth(providerName string, success bool, responseTime time.Duration) {
	hm.mu.Lock()
	defer hm.mu.Unlock()
	
	health, exists := hm.providers[providerName]
	if !exists {
		health = &ProviderHealth{
			ProviderName: providerName,
		}
		hm.providers[providerName] = health
	}
	
	health.mu.Lock()
	defer health.mu.Unlock()
	
	health.LastChecked = time.Now()
	health.TotalRequests++
	
	if success {
		health.IsAvailable = true
		health.LastSuccess = time.Now()
		health.SuccessfulRequests++
		health.ConsecutiveFailures = 0
		
		// Update average response time
		if health.AverageResponseTime == 0 {
			health.AverageResponseTime = responseTime
		} else {
			// Exponential moving average
			health.AverageResponseTime = time.Duration(
				float64(health.AverageResponseTime)*0.7 + float64(responseTime)*0.3,
			)
		}
	} else {
		health.LastFailure = time.Now()
		health.FailedRequests++
		health.ConsecutiveFailures++
		
		// Mark as unavailable after 3 consecutive failures
		if health.ConsecutiveFailures >= 3 {
			health.IsAvailable = false
		}
	}
}

// GetHealth returns the health status of a provider
func (hm *HealthMonitor) GetHealth(providerName string) (*ProviderHealth, bool) {
	hm.mu.RLock()
	defer hm.mu.RUnlock()
	
	health, exists := hm.providers[providerName]
	if !exists {
		return nil, false
	}
	
	health.mu.RLock()
	defer health.mu.RUnlock()
	
	// Return a copy to avoid race conditions
	return &ProviderHealth{
		ProviderName:        health.ProviderName,
		IsAvailable:         health.IsAvailable,
		LastChecked:         health.LastChecked,
		LastSuccess:         health.LastSuccess,
		LastFailure:         health.LastFailure,
		ConsecutiveFailures: health.ConsecutiveFailures,
		TotalRequests:       health.TotalRequests,
		SuccessfulRequests: health.SuccessfulRequests,
		FailedRequests:     health.FailedRequests,
		AverageResponseTime: health.AverageResponseTime,
	}, true
}

// GetAllHealth returns health status of all providers
func (hm *HealthMonitor) GetAllHealth() map[string]*ProviderHealth {
	hm.mu.RLock()
	defer hm.mu.RUnlock()
	
	result := make(map[string]*ProviderHealth)
	for name, health := range hm.providers {
		health.mu.RLock()
		result[name] = &ProviderHealth{
			ProviderName:        health.ProviderName,
			IsAvailable:         health.IsAvailable,
			LastChecked:         health.LastChecked,
			LastSuccess:         health.LastSuccess,
			LastFailure:         health.LastFailure,
			ConsecutiveFailures: health.ConsecutiveFailures,
			TotalRequests:       health.TotalRequests,
			SuccessfulRequests: health.SuccessfulRequests,
			FailedRequests:     health.FailedRequests,
			AverageResponseTime: health.AverageResponseTime,
		}
		health.mu.RUnlock()
	}
	return result
}

// GetAvailableProviders returns a list of available provider names
func (hm *HealthMonitor) GetAvailableProviders() []string {
	hm.mu.RLock()
	defer hm.mu.RUnlock()
	
	var available []string
	for name, health := range hm.providers {
		health.mu.RLock()
		if health.IsAvailable {
			available = append(available, name)
		}
		health.mu.RUnlock()
	}
	return available
}

// CheckProviderHealth performs a health check on a provider
func (hm *HealthMonitor) CheckProviderHealth(provider Provider) bool {
	if provider == nil {
		return false
	}
	
	startTime := time.Now()
	available := provider.IsAvailable()
	responseTime := time.Since(startTime)
	
	hm.UpdateHealth(provider.GetName(), available, responseTime)
	return available
}

// GetProviderStatus returns a human-readable status string
func (h *ProviderHealth) GetProviderStatus() string {
	if h.IsAvailable {
		successRate := float64(0)
		if h.TotalRequests > 0 {
			successRate = float64(h.SuccessfulRequests) / float64(h.TotalRequests) * 100
		}
		return fmt.Sprintf("✅ Available (%.1f%% success, avg: %v)", 
			successRate, h.AverageResponseTime)
	}
	
	if h.ConsecutiveFailures > 0 {
		return fmt.Sprintf("❌ Unavailable (%d consecutive failures)", h.ConsecutiveFailures)
	}
	
	return "⚠️  Unknown"
}

