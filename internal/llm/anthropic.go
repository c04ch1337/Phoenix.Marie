package llm

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

// AnthropicClient handles communication with Anthropic API
type AnthropicClient struct {
	apiKey     string
	baseURL    string
	httpClient *http.Client
	config     *Config
}

// NewAnthropicClient creates a new Anthropic client
func NewAnthropicClient(config *Config) *AnthropicClient {
	baseURL := config.AnthropicBaseURL
	if baseURL == "" {
		baseURL = "https://api.anthropic.com/v1"
	}

	return &AnthropicClient{
		apiKey:  config.AnthropicAPIKey,
		baseURL: baseURL,
		httpClient: &http.Client{
			Timeout: time.Duration(config.RequestTimeout) * time.Second,
		},
		config: config,
	}
}

// GetName returns the provider name
func (c *AnthropicClient) GetName() string {
	return "anthropic"
}

// IsAvailable checks if the provider is available
func (c *AnthropicClient) IsAvailable() bool {
	return c.apiKey != ""
}

// AnthropicRequest represents the request format for Anthropic
type AnthropicRequest struct {
	Model       string    `json:"model"`
	Messages    []Message `json:"messages"`
	MaxTokens   int       `json:"max_tokens"`
	Temperature float64   `json:"temperature,omitempty"`
	TopP        float64   `json:"top_p,omitempty"`
}

// AnthropicResponse represents the response from Anthropic
type AnthropicResponse struct {
	ID      string `json:"id"`
	Model   string `json:"model"`
	Content []struct {
		Type string `json:"type"`
		Text string `json:"text"`
	} `json:"content"`
	StopReason string `json:"stop_reason"`
	Usage      struct {
		InputTokens  int `json:"input_tokens"`
		OutputTokens int `json:"output_tokens"`
	} `json:"usage"`
}

// Call makes a request to Anthropic API
func (c *AnthropicClient) Call(modelID string, messages []Message, maxTokens int, temperature float64) (*Response, error) {
	startTime := time.Now()

	if maxTokens == 0 {
		maxTokens = c.config.DefaultMaxTokens
	}
	if temperature == 0.0 {
		temperature = c.config.DefaultTemperature
	}

	reqBody := AnthropicRequest{
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

	req, err := http.NewRequest("POST", c.baseURL+"/messages", bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("x-api-key", c.apiKey)
	req.Header.Set("anthropic-version", "2023-06-01")
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

	var anthropicResp AnthropicResponse
	if err := json.NewDecoder(resp.Body).Decode(&anthropicResp); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	if len(anthropicResp.Content) == 0 {
		return nil, fmt.Errorf("no content in response")
	}

	responseTime := time.Since(startTime)

	// Calculate cost
	model, exists := GetModel(modelID)
	if !exists {
		model = Model{InputPrice: 1.0, OutputPrice: 1.0}
	}

	cost := c.calculateCost(
		anthropicResp.Usage.InputTokens,
		anthropicResp.Usage.OutputTokens,
		model.InputPrice,
		model.OutputPrice,
	)

	content := anthropicResp.Content[0].Text
	for i := 1; i < len(anthropicResp.Content); i++ {
		content += "\n" + anthropicResp.Content[i].Text
	}

	return &Response{
		Content: content,
		Model:   anthropicResp.Model,
		TokensUsed: TokenUsage{
			PromptTokens:     anthropicResp.Usage.InputTokens,
			CompletionTokens: anthropicResp.Usage.OutputTokens,
			TotalTokens:      anthropicResp.Usage.InputTokens + anthropicResp.Usage.OutputTokens,
		},
		Cost:         cost,
		ResponseTime: responseTime,
		FinishReason: anthropicResp.StopReason,
	}, nil
}

// calculateCost calculates the cost based on token usage
func (c *AnthropicClient) calculateCost(promptTokens, completionTokens int, inputPrice, outputPrice float64) float64 {
	promptCost := (float64(promptTokens) / 1_000_000.0) * inputPrice
	completionCost := (float64(completionTokens) / 1_000_000.0) * outputPrice
	return promptCost + completionCost
}

// CallWithRetry makes a request with retry logic
func (c *AnthropicClient) CallWithRetry(modelID string, messages []Message, maxTokens int, temperature float64) (*Response, error) {
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

