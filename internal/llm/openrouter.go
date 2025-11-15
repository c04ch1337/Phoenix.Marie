package llm

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

// OpenRouterClient handles communication with OpenRouter API
type OpenRouterClient struct {
	apiKey     string
	baseURL    string
	httpClient *http.Client
	config     *Config
}

// NewOpenRouterClient creates a new OpenRouter client
func NewOpenRouterClient(config *Config) *OpenRouterClient {
	return &OpenRouterClient{
		apiKey:  config.OpenRouterAPIKey,
		baseURL: "https://openrouter.ai/api/v1",
		httpClient: &http.Client{
			Timeout: time.Duration(config.RequestTimeout) * time.Second,
		},
		config: config,
	}
}

// OpenRouterRequest represents the request format for OpenRouter
type OpenRouterRequest struct {
	Model       string    `json:"model"`
	Messages    []Message `json:"messages"`
	MaxTokens   int       `json:"max_tokens,omitempty"`
	Temperature float64   `json:"temperature,omitempty"`
	TopP        float64   `json:"top_p,omitempty"`
}

// Message represents a chat message
type Message struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

// OpenRouterResponse represents the response from OpenRouter
type OpenRouterResponse struct {
	ID      string `json:"id"`
	Model   string `json:"model"`
	Choices []struct {
		Message struct {
			Role    string `json:"role"`
			Content string `json:"content"`
		} `json:"message"`
		FinishReason string `json:"finish_reason"`
	} `json:"choices"`
	Usage struct {
		PromptTokens     int `json:"prompt_tokens"`
		CompletionTokens int `json:"completion_tokens"`
		TotalTokens      int `json:"total_tokens"`
	} `json:"usage"`
}

// Call makes a request to OpenRouter API
func (c *OpenRouterClient) Call(modelID string, messages []Message, maxTokens int, temperature float64) (*Response, error) {
	startTime := time.Now()
	
	// Use defaults from config if not specified
	if maxTokens == 0 {
		maxTokens = c.config.DefaultMaxTokens
	}
	if temperature == 0.0 {
		temperature = c.config.DefaultTemperature
	}
	
	reqBody := OpenRouterRequest{
		Model:       modelID,
		Messages:    messages,
		MaxTokens:   maxTokens,
		Temperature: temperature,
		TopP:        c.config.DefaultTopP,
	}
	
	jsonData, err := json.Marshal(reqBody)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request: %w", err)
	}
	
	req, err := http.NewRequest("POST", c.baseURL+"/chat/completions", bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}
	
	req.Header.Set("Authorization", "Bearer "+c.apiKey)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("HTTP-Referer", "https://github.com/phoenix-marie/core") // Optional
	req.Header.Set("X-Title", "Phoenix.Marie")                               // Optional
	
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to make request: %w", err)
	}
	defer resp.Body.Close()
	
	if resp.StatusCode != http.StatusOK {
		bodyBytes, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("API error (status %d): %s", resp.StatusCode, string(bodyBytes))
	}
	
	var openRouterResp OpenRouterResponse
	if err := json.NewDecoder(resp.Body).Decode(&openRouterResp); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}
	
	if len(openRouterResp.Choices) == 0 {
		return nil, fmt.Errorf("no choices in response")
	}
	
	responseTime := time.Since(startTime)
	
	// Calculate cost
	model, exists := GetModel(modelID)
	if !exists {
		// Use default pricing if model not found
		model = Model{InputPrice: 1.0, OutputPrice: 1.0}
	}
	
	cost := c.calculateCost(
		openRouterResp.Usage.PromptTokens,
		openRouterResp.Usage.CompletionTokens,
		model.InputPrice,
		model.OutputPrice,
	)
	
	return &Response{
		Content: openRouterResp.Choices[0].Message.Content,
		Model:   openRouterResp.Model,
		TokensUsed: TokenUsage{
			PromptTokens:     openRouterResp.Usage.PromptTokens,
			CompletionTokens: openRouterResp.Usage.CompletionTokens,
			TotalTokens:      openRouterResp.Usage.TotalTokens,
		},
		Cost:         cost,
		ResponseTime: responseTime,
		FinishReason: openRouterResp.Choices[0].FinishReason,
	}, nil
}

// calculateCost calculates the cost based on token usage
func (c *OpenRouterClient) calculateCost(promptTokens, completionTokens int, inputPrice, outputPrice float64) float64 {
	promptCost := (float64(promptTokens) / 1_000_000.0) * inputPrice
	completionCost := (float64(completionTokens) / 1_000_000.0) * outputPrice
	return promptCost + completionCost
}

// CallWithRetry makes a request with retry logic
func (c *OpenRouterClient) CallWithRetry(modelID string, messages []Message, maxTokens int, temperature float64) (*Response, error) {
	var lastErr error
	
	for attempt := 0; attempt < c.config.MaxRetries; attempt++ {
		if attempt > 0 {
			// Exponential backoff
			backoff := time.Duration(attempt) * time.Second
			time.Sleep(backoff)
		}
		
		resp, err := c.Call(modelID, messages, maxTokens, temperature)
		if err == nil {
			return resp, nil
		}
		
		lastErr = err
	}
	
	return nil, fmt.Errorf("failed after %d retries: %w", c.config.MaxRetries, lastErr)
}

