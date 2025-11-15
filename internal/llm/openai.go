package llm

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

// OpenAIClient handles communication with OpenAI API
type OpenAIClient struct {
	apiKey     string
	baseURL    string
	httpClient *http.Client
	config     *Config
}

// NewOpenAIClient creates a new OpenAI client
func NewOpenAIClient(config *Config) *OpenAIClient {
	baseURL := config.OpenAIBaseURL
	if baseURL == "" {
		baseURL = "https://api.openai.com/v1"
	}

	return &OpenAIClient{
		apiKey:  config.OpenAIAPIKey,
		baseURL: baseURL,
		httpClient: &http.Client{
			Timeout: time.Duration(config.RequestTimeout) * time.Second,
		},
		config: config,
	}
}

// GetName returns the provider name
func (c *OpenAIClient) GetName() string {
	return "openai"
}

// IsAvailable checks if the provider is available
func (c *OpenAIClient) IsAvailable() bool {
	return c.apiKey != ""
}

// OpenAIRequest represents the request format for OpenAI
type OpenAIRequest struct {
	Model       string    `json:"model"`
	Messages    []Message `json:"messages"`
	MaxTokens   int       `json:"max_tokens,omitempty"`
	Temperature float64   `json:"temperature,omitempty"`
	TopP        float64   `json:"top_p,omitempty"`
}

// OpenAIResponse represents the response from OpenAI
type OpenAIResponse struct {
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

// Call makes a request to OpenAI API
func (c *OpenAIClient) Call(modelID string, messages []Message, maxTokens int, temperature float64) (*Response, error) {
	startTime := time.Now()

	if maxTokens == 0 {
		maxTokens = c.config.DefaultMaxTokens
	}
	if temperature == 0.0 {
		temperature = c.config.DefaultTemperature
	}

	reqBody := OpenAIRequest{
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

	var openAIResp OpenAIResponse
	if err := json.NewDecoder(resp.Body).Decode(&openAIResp); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	if len(openAIResp.Choices) == 0 {
		return nil, fmt.Errorf("no choices in response")
	}

	responseTime := time.Since(startTime)

	// Calculate cost (OpenAI pricing)
	model, exists := GetModel(modelID)
	if !exists {
		model = Model{InputPrice: 1.0, OutputPrice: 1.0}
	}

	cost := c.calculateCost(
		openAIResp.Usage.PromptTokens,
		openAIResp.Usage.CompletionTokens,
		model.InputPrice,
		model.OutputPrice,
	)

	return &Response{
		Content: openAIResp.Choices[0].Message.Content,
		Model:   openAIResp.Model,
		TokensUsed: TokenUsage{
			PromptTokens:     openAIResp.Usage.PromptTokens,
			CompletionTokens: openAIResp.Usage.CompletionTokens,
			TotalTokens:      openAIResp.Usage.TotalTokens,
		},
		Cost:         cost,
		ResponseTime: responseTime,
		FinishReason: openAIResp.Choices[0].FinishReason,
	}, nil
}

// calculateCost calculates the cost based on token usage
func (c *OpenAIClient) calculateCost(promptTokens, completionTokens int, inputPrice, outputPrice float64) float64 {
	promptCost := (float64(promptTokens) / 1_000_000.0) * inputPrice
	completionCost := (float64(completionTokens) / 1_000_000.0) * outputPrice
	return promptCost + completionCost
}

// CallWithRetry makes a request with retry logic
func (c *OpenAIClient) CallWithRetry(modelID string, messages []Message, maxTokens int, temperature float64) (*Response, error) {
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

