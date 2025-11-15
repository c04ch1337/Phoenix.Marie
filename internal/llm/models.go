package llm

// GetAvailableModels returns all available models with their configurations
func GetAvailableModels() map[string]Model {
	return map[string]Model{
		"anthropic/claude-3-opus": {
			ID:            "anthropic/claude-3-opus",
			Name:          "Claude 3 Opus",
			Provider:      "anthropic",
			ContextLength: 200000,
			InputPrice:    15.0,
			OutputPrice:   75.0,
			Capabilities: Capabilities{
				Reasoning:   true,
				Creativity:   true,
				Speed:        false,
				ToolUse:      true,
				Multimodal:   false,
				Multilingual: true,
				Math:         true,
			},
		},
		"openai/gpt-4-turbo": {
			ID:            "openai/gpt-4-turbo",
			Name:          "GPT-4 Turbo",
			Provider:      "openai",
			ContextLength: 128000,
			InputPrice:    10.0,
			OutputPrice:   30.0,
			Capabilities: Capabilities{
				Reasoning:   true,
				Creativity:   true,
				Speed:        true,
				ToolUse:      true,
				Multimodal:   false,
				Multilingual: true,
				Math:         true,
			},
		},
		"anthropic/claude-3-sonnet": {
			ID:            "anthropic/claude-3-sonnet",
			Name:          "Claude 3 Sonnet",
			Provider:      "anthropic",
			ContextLength: 200000,
			InputPrice:    3.0,
			OutputPrice:   15.0,
			Capabilities: Capabilities{
				Reasoning:   true,
				Creativity:   true,
				Speed:        true,
				ToolUse:      true,
				Multimodal:   false,
				Multilingual: true,
				Math:         true,
			},
		},
		"google/gemini-pro-1.5": {
			ID:            "google/gemini-pro-1.5",
			Name:          "Gemini Pro 1.5",
			Provider:      "google",
			ContextLength: 1000000, // Theoretical, actual may vary
			InputPrice:    1.25,    // Approximate
			OutputPrice:   5.0,     // Approximate
			Capabilities: Capabilities{
				Reasoning:   true,
				Creativity:   true,
				Speed:        true,
				ToolUse:      true,
				Multimodal:   true,
				Multilingual: true,
				Math:         true,
			},
		},
		"openai/gpt-4-vision-preview": {
			ID:            "openai/gpt-4-vision-preview",
			Name:          "GPT-4 Vision",
			Provider:      "openai",
			ContextLength: 128000,
			InputPrice:    10.0, // Varies with images
			OutputPrice:   30.0,
			Capabilities: Capabilities{
				Reasoning:   true,
				Creativity:   true,
				Speed:        true,
				ToolUse:      true,
				Multimodal:   true,
				Multilingual: true,
				Math:         true,
			},
		},
		"anthropic/claude-3-haiku": {
			ID:            "anthropic/claude-3-haiku",
			Name:          "Claude 3 Haiku",
			Provider:      "anthropic",
			ContextLength: 200000,
			InputPrice:    0.25,
			OutputPrice:   1.25,
			Capabilities: Capabilities{
				Reasoning:   false,
				Creativity:   false,
				Speed:        true,
				ToolUse:      true,
				Multimodal:   false,
				Multilingual: true,
				Math:         false,
			},
		},
		"mistralai/mixtral-8x22b": {
			ID:            "mistralai/mixtral-8x22b",
			Name:          "Mixtral 8x22B",
			Provider:      "mistralai",
			ContextLength: 64000,
			InputPrice:    2.0,
			OutputPrice:   6.0,
			Capabilities: Capabilities{
				Reasoning:   true,
				Creativity:   true,
				Speed:        true,
				ToolUse:      false,
				Multimodal:   false,
				Multilingual: true,
				Math:         true,
			},
		},
		"cohere/command-r-plus": {
			ID:            "cohere/command-r-plus",
			Name:          "Command R+",
			Provider:      "cohere",
			ContextLength: 128000,
			InputPrice:    3.0,
			OutputPrice:   15.0,
			Capabilities: Capabilities{
				Reasoning:   false,
				Creativity:   false,
				Speed:        true,
				ToolUse:      true,
				Multimodal:   false,
				Multilingual: true,
				Math:         false,
			},
		},
		"meta-llama/llama-3-70b-instruct": {
			ID:            "meta-llama/llama-3-70b-instruct",
			Name:          "Llama 3 70B",
			Provider:      "meta",
			ContextLength: 8000,
			InputPrice:    0.90,
			OutputPrice:   0.90,
			Capabilities: Capabilities{
				Reasoning:   true,
				Creativity:   true,
				Speed:        true,
				ToolUse:      false,
				Multimodal:   false,
				Multilingual: true,
				Math:         true,
			},
		},
		"qwen/qwen-2-72b-instruct": {
			ID:            "qwen/qwen-2-72b-instruct",
			Name:          "Qwen 2 72B",
			Provider:      "qwen",
			ContextLength: 32000,
			InputPrice:    1.50,
			OutputPrice:   1.50,
			Capabilities: Capabilities{
				Reasoning:   true,
				Creativity:   false,
				Speed:        true,
				ToolUse:      false,
				Multimodal:   false,
				Multilingual: true,
				Math:         true,
			},
		},
	}
}

// GetModel returns a model by ID
func GetModel(modelID string) (Model, bool) {
	models := GetAvailableModels()
	model, exists := models[modelID]
	return model, exists
}

// GetModelHierarchy returns models in order of capability (best to cheapest)
func GetModelHierarchy() []string {
	return []string{
		"anthropic/claude-3-opus",
		"openai/gpt-4-turbo",
		"anthropic/claude-3-sonnet",
		"google/gemini-pro-1.5",
		"mistralai/mixtral-8x22b",
		"cohere/command-r-plus",
		"meta-llama/llama-3-70b-instruct",
		"anthropic/claude-3-haiku",
		"qwen/qwen-2-72b-instruct",
	}
}

