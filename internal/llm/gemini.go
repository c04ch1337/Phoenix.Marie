package llm

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

// GeminiClient handles communication with Google Gemini API
type GeminiClient struct {
	apiKey     string
	baseURL    string
	httpClient *http.Client
	config     *Config
}

// NewGeminiClient creates a new Gemini client
func NewGeminiClient(config *Config) *GeminiClient {
	baseURL := config.GeminiBaseURL
	if baseURL == "" {
		baseURL = "https://generativelanguage.googleapis.com/v1"
	}

	return &GeminiClient{
		apiKey:  config.GeminiAPIKey,
		baseURL: baseURL,
		httpClient: &http.Client{
			Timeout: time.Duration(config.RequestTimeout) * time.Second,
		},
		config: config,
	}
}

// GetName returns the provider name
func (c *GeminiClient) GetName() string {
	return "gemini"
}

// IsAvailable checks if the provider is available
func (c *GeminiClient) IsAvailable() bool {
	return c.apiKey != ""
}

// GeminiRequest represents the request format for Gemini
type GeminiRequest struct {
	Contents []struct {
		Parts []struct {
			Text string `json:"text"`
		} `json:"parts"`
	} `json:"contents"`
	GenerationConfig struct {
		MaxOutputTokens int     `json:"maxOutputTokens,omitempty"`
		Temperature     float64 `json:"temperature,omitempty"`
		TopP            float64 `json:"topP,omitempty"`
	} `json:"generationConfig,omitempty"`
}

// GeminiResponse represents the response from Gemini
type GeminiResponse struct {
	Candidates []struct {
		Content struct {
			Parts []struct {
				Text string `json:"text"`
			} `json:"parts"`
		} `json:"content"`
		FinishReason string `json:"finishReason"`
	} `json:"candidates"`
	UsageMetadata struct {
		PromptTokenCount     int `json:"promptTokenCount"`
		CandidatesTokenCount int `json:"candidatesTokenCount"`
		TotalTokenCount      int `json:"totalTokenCount"`
	} `json:"usageMetadata"`
}

// Call makes a request to Gemini API
func (c *GeminiClient) Call(modelID string, messages []Message, maxTokens int, temperature float64) (*Response, error) {
	startTime := time.Now()

	if maxTokens == 0 {
		maxTokens = c.config.DefaultMaxTokens
	}
	if temperature == 0.0 {
		temperature = c.config.DefaultTemperature
	}

	// Convert messages to Gemini format
	contents := make([]struct {
		Parts []struct {
			Text string `json:"text"`
		} `json:"parts"`
	}, 0)

	for _, msg := range messages {
		content := struct {
			Parts []struct {
				Text string `json:"text"`
			} `json:"parts"`
		}{
			Parts: []struct {
				Text string `json:"text"`
			}{
				{Text: msg.Content},
			},
		}
		contents = append(contents, content)
	}

	reqBody := GeminiRequest{
		Contents: contents,
		GenerationConfig: struct {
			MaxOutputTokens int     `json:"maxOutputTokens,omitempty"`
			Temperature     float64 `json:"temperature,omitempty"`
			TopP            float64 `json:"topP,omitempty"`
		}{
			MaxOutputTokens: maxTokens,
			Temperature:     temperature,
			TopP:            c.config.DefaultTopP,
		},
	}

	jsonData, err := json.Marshal(reqBody)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request: %w", err)
	}

	url := fmt.Sprintf("%s/models/%s:generateContent?key=%s", c.baseURL, modelID, c.apiKey)
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
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

	var geminiResp GeminiResponse
	if err := json.NewDecoder(resp.Body).Decode(&geminiResp); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	if len(geminiResp.Candidates) == 0 {
		return nil, fmt.Errorf("no candidates in response")
	}

	responseTime := time.Since(startTime)

	// Calculate cost
	model, exists := GetModel(modelID)
	if !exists {
		model = Model{InputPrice: 1.0, OutputPrice: 1.0}
	}

	cost := c.calculateCost(
		geminiResp.UsageMetadata.PromptTokenCount,
		geminiResp.UsageMetadata.CandidatesTokenCount,
		model.InputPrice,
		model.OutputPrice,
	)

	content := geminiResp.Candidates[0].Content.Parts[0].Text
	for i := 1; i < len(geminiResp.Candidates[0].Content.Parts); i++ {
		content += "\n" + geminiResp.Candidates[0].Content.Parts[i].Text
	}

	return &Response{
		Content: content,
		Model:   modelID,
		TokensUsed: TokenUsage{
			PromptTokens:     geminiResp.UsageMetadata.PromptTokenCount,
			CompletionTokens: geminiResp.UsageMetadata.CandidatesTokenCount,
			TotalTokens:      geminiResp.UsageMetadata.TotalTokenCount,
		},
		Cost:         cost,
		ResponseTime: responseTime,
		FinishReason: geminiResp.Candidates[0].FinishReason,
	}, nil
}

// calculateCost calculates the cost based on token usage
func (c *GeminiClient) calculateCost(promptTokens, completionTokens int, inputPrice, outputPrice float64) float64 {
	promptCost := (float64(promptTokens) / 1_000_000.0) * inputPrice
	completionCost := (float64(completionTokens) / 1_000_000.0) * outputPrice
	return promptCost + completionCost
}

// CallWithRetry makes a request with retry logic
func (c *GeminiClient) CallWithRetry(modelID string, messages []Message, maxTokens int, temperature float64) (*Response, error) {
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

