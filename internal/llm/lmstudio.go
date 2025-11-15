package llm

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

// LMStudioClient handles communication with local LM Studio API
type LMStudioClient struct {
	baseURL    string
	httpClient *http.Client
	config     *Config
}

// NewLMStudioClient creates a new LM Studio client
func NewLMStudioClient(config *Config) *LMStudioClient {
	baseURL := config.LMStudioBaseURL
	if baseURL == "" {
		baseURL = "http://localhost:1234"
	}

	return &LMStudioClient{
		baseURL: baseURL,
		httpClient: &http.Client{
			Timeout: time.Duration(config.RequestTimeout) * time.Second,
		},
		config: config,
	}
}

// GetName returns the provider name
func (c *LMStudioClient) GetName() string {
	return "lmstudio"
}

// IsAvailable checks if the provider is available
func (c *LMStudioClient) IsAvailable() bool {
	// Try to ping LM Studio
	resp, err := http.Get(c.baseURL + "/v1/models")
	if err != nil {
		return false
	}
	defer resp.Body.Close()
	return resp.StatusCode == http.StatusOK
}

// LMStudioRequest represents the request format for LM Studio (OpenAI-compatible)
type LMStudioRequest struct {
	Model       string    `json:"model"`
	Messages    []Message `json:"messages"`
	MaxTokens   int       `json:"max_tokens,omitempty"`
	Temperature float64   `json:"temperature,omitempty"`
	TopP        float64   `json:"top_p,omitempty"`
}

// LMStudioResponse represents the response from LM Studio
type LMStudioResponse struct {
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

// Call makes a request to LM Studio API
func (c *LMStudioClient) Call(modelID string, messages []Message, maxTokens int, temperature float64) (*Response, error) {
	startTime := time.Now()

	if maxTokens == 0 {
		maxTokens = c.config.DefaultMaxTokens
	}
	if temperature == 0.0 {
		temperature = c.config.DefaultTemperature
	}

	reqBody := LMStudioRequest{
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

	req, err := http.NewRequest("POST", c.baseURL+"/v1/chat/completions", bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

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

	var lmStudioResp LMStudioResponse
	if err := json.NewDecoder(resp.Body).Decode(&lmStudioResp); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	if len(lmStudioResp.Choices) == 0 {
		return nil, fmt.Errorf("no choices in response")
	}

	responseTime := time.Since(startTime)

	// LM Studio is free (local), so cost is 0
	return &Response{
		Content: lmStudioResp.Choices[0].Message.Content,
		Model:   lmStudioResp.Model,
		TokensUsed: TokenUsage{
			PromptTokens:     lmStudioResp.Usage.PromptTokens,
			CompletionTokens: lmStudioResp.Usage.CompletionTokens,
			TotalTokens:      lmStudioResp.Usage.TotalTokens,
		},
		Cost:         0.0, // Local, no cost
		ResponseTime: responseTime,
		FinishReason: lmStudioResp.Choices[0].FinishReason,
	}, nil
}

// CallWithRetry makes a request with retry logic
func (c *LMStudioClient) CallWithRetry(modelID string, messages []Message, maxTokens int, temperature float64) (*Response, error) {
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

