package llm

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

// GrokClient handles communication with xAI Grok API
type GrokClient struct {
	apiKey     string
	baseURL    string
	httpClient *http.Client
	config     *Config
}

// NewGrokClient creates a new Grok client
func NewGrokClient(config *Config) *GrokClient {
	baseURL := config.GrokBaseURL
	if baseURL == "" {
		baseURL = "https://api.x.ai/v1"
	}

	return &GrokClient{
		apiKey:  config.GrokAPIKey,
		baseURL: baseURL,
		httpClient: &http.Client{
			Timeout: time.Duration(config.RequestTimeout) * time.Second,
		},
		config: config,
	}
}

// GetName returns the provider name
func (c *GrokClient) GetName() string {
	return "grok"
}

// IsAvailable checks if the provider is available
func (c *GrokClient) IsAvailable() bool {
	return c.apiKey != ""
}

// GrokRequest represents the request format for Grok
type GrokRequest struct {
	Model       string    `json:"model"`
	Messages    []Message `json:"messages"`
	MaxTokens   int       `json:"max_tokens,omitempty"`
	Temperature float64   `json:"temperature,omitempty"`
	TopP        float64   `json:"top_p,omitempty"`
}

// GrokResponse represents the response from Grok
type GrokResponse struct {
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

// Call makes a request to Grok API
func (c *GrokClient) Call(modelID string, messages []Message, maxTokens int, temperature float64) (*Response, error) {
	startTime := time.Now()

	if maxTokens == 0 {
		maxTokens = c.config.DefaultMaxTokens
	}
	if temperature == 0.0 {
		temperature = c.config.DefaultTemperature
	}

	reqBody := GrokRequest{
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

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to make request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		bodyBytes, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("API error (status %d): %s", resp.StatusCode, string(bodyBytes))
	}

	var grokResp GrokResponse
	if err := json.NewDecoder(resp.Body).Decode(&grokResp); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	if len(grokResp.Choices) == 0 {
		return nil, fmt.Errorf("no choices in response")
	}

	responseTime := time.Since(startTime)

	// Calculate cost
	model, exists := GetModel(modelID)
	if !exists {
		model = Model{InputPrice: 1.0, OutputPrice: 1.0}
	}

	cost := c.calculateCost(
		grokResp.Usage.PromptTokens,
		grokResp.Usage.CompletionTokens,
		model.InputPrice,
		model.OutputPrice,
	)

	return &Response{
		Content: grokResp.Choices[0].Message.Content,
		Model:   grokResp.Model,
		TokensUsed: TokenUsage{
			PromptTokens:     grokResp.Usage.PromptTokens,
			CompletionTokens: grokResp.Usage.CompletionTokens,
			TotalTokens:      grokResp.Usage.TotalTokens,
		},
		Cost:         cost,
		ResponseTime: responseTime,
		FinishReason: grokResp.Choices[0].FinishReason,
	}, nil
}

// calculateCost calculates the cost based on token usage
func (c *GrokClient) calculateCost(promptTokens, completionTokens int, inputPrice, outputPrice float64) float64 {
	promptCost := (float64(promptTokens) / 1_000_000.0) * inputPrice
	completionCost := (float64(completionTokens) / 1_000_000.0) * outputPrice
	return promptCost + completionCost
}

// CallWithRetry makes a request with retry logic
func (c *GrokClient) CallWithRetry(modelID string, messages []Message, maxTokens int, temperature float64) (*Response, error) {
	var lastErr error

	for attempt := 0; attempt < c.config.MaxRetries; attempt++ {
		if attempt > 0 {
			backoff := time.Duration(attempt) * time.Duration(c.config.RetryBackoff) * time.Second
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

