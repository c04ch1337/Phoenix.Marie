package llm

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

// OllamaClient handles communication with local Ollama API
type OllamaClient struct {
	baseURL    string
	httpClient *http.Client
	config     *Config
}

// NewOllamaClient creates a new Ollama client
func NewOllamaClient(config *Config) *OllamaClient {
	baseURL := config.OllamaBaseURL
	if baseURL == "" {
		baseURL = "http://localhost:11434"
	}

	return &OllamaClient{
		baseURL: baseURL,
		httpClient: &http.Client{
			Timeout: time.Duration(config.RequestTimeout) * time.Second,
		},
		config: config,
	}
}

// GetName returns the provider name
func (c *OllamaClient) GetName() string {
	return "ollama"
}

// IsAvailable checks if the provider is available
func (c *OllamaClient) IsAvailable() bool {
	// Try to ping Ollama
	resp, err := http.Get(c.baseURL + "/api/tags")
	if err != nil {
		return false
	}
	defer resp.Body.Close()
	return resp.StatusCode == http.StatusOK
}

// OllamaRequest represents the request format for Ollama
type OllamaRequest struct {
	Model       string    `json:"model"`
	Messages    []Message `json:"messages"`
	Stream      bool      `json:"stream"`
	Options     struct {
		Temperature float64 `json:"temperature,omitempty"`
		TopP        float64 `json:"top_p,omitempty"`
		NumPredict  int     `json:"num_predict,omitempty"`
	} `json:"options,omitempty"`
}

// OllamaResponse represents the response from Ollama
type OllamaResponse struct {
	Model     string `json:"model"`
	Message   Message `json:"message"`
	Done      bool   `json:"done"`
	TotalDuration int64 `json:"total_duration"`
	PromptEvalCount int `json:"prompt_eval_count"`
	EvalCount int    `json:"eval_count"`
}

// Call makes a request to Ollama API
func (c *OllamaClient) Call(modelID string, messages []Message, maxTokens int, temperature float64) (*Response, error) {
	startTime := time.Now()

	if maxTokens == 0 {
		maxTokens = c.config.DefaultMaxTokens
	}
	if temperature == 0.0 {
		temperature = c.config.DefaultTemperature
	}

	reqBody := OllamaRequest{
		Model:   modelID,
		Messages: messages,
		Stream:  false,
	}
	reqBody.Options.Temperature = temperature
	reqBody.Options.TopP = c.config.DefaultTopP
	reqBody.Options.NumPredict = maxTokens

	jsonData, err := json.Marshal(reqBody)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request: %w", err)
	}

	req, err := http.NewRequest("POST", c.baseURL+"/api/chat", bytes.NewBuffer(jsonData))
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

	var ollamaResp OllamaResponse
	if err := json.NewDecoder(resp.Body).Decode(&ollamaResp); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	responseTime := time.Since(startTime)

	// Ollama is free (local), so cost is 0
	return &Response{
		Content: ollamaResp.Message.Content,
		Model:   ollamaResp.Model,
		TokensUsed: TokenUsage{
			PromptTokens:     ollamaResp.PromptEvalCount,
			CompletionTokens: ollamaResp.EvalCount,
			TotalTokens:      ollamaResp.PromptEvalCount + ollamaResp.EvalCount,
		},
		Cost:         0.0, // Local, no cost
		ResponseTime: responseTime,
		FinishReason: "stop",
	}, nil
}

// CallWithRetry makes a request with retry logic
func (c *OllamaClient) CallWithRetry(modelID string, messages []Message, maxTokens int, temperature float64) (*Response, error) {
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

