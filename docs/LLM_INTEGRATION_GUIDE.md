# LLM Integration Guide — Phoenix.Marie

## Overview

Phoenix.Marie now includes a complete LLM integration system supporting **10 top models** via OpenRouter, with intelligent routing, cost management, and consciousness-aware prompts.

---

## Quick Start

### 1. Get OpenRouter API Key

1. Go to https://openrouter.ai/keys
2. Create an account and generate an API key
3. Copy your key (starts with `sk-or-v1-...`)

### 2. Configure `.env.local`

Copy the example configuration:

```bash
cp .env.local.example .env.local
```

Edit `.env.local` and add your API key:

```bash
OPENROUTER_API_KEY=sk-or-v1-your-actual-key-here
```

### 3. Configure Models (Optional)

The system comes with sensible defaults, but you can customize:

```bash
# Phoenix.Marie Models
PHOENIX_CONSCIOUSNESS_MODEL=openai/gpt-4-turbo
PHOENIX_EMOTIONAL_MODEL=anthropic/claude-3-sonnet
PHOENIX_VOICE_MODEL=anthropic/claude-3-haiku

# Jamey 3.0 Models
JAMEY_REASONING_MODEL=anthropic/claude-3-opus
JAMEY_OPERATIONAL_MODEL=anthropic/claude-3-sonnet
JAMEY_REALTIME_MODEL=anthropic/claude-3-haiku
```

### 4. Run Phoenix.Marie

```bash
make build
make run
```

The LLM client will initialize automatically if `OPENROUTER_API_KEY` is set.

---

## Available Models

### Top 10 Models Supported

1. **Claude 3 Opus** (`anthropic/claude-3-opus`)
   - Best for: Complex reasoning, consciousness simulation
   - Context: 200K tokens
   - Cost: $15/M input, $75/M output

2. **GPT-4 Turbo** (`openai/gpt-4-turbo`)
   - Best for: General intelligence, creative tasks
   - Context: 128K tokens
   - Cost: $10/M input, $30/M output

3. **Claude 3 Sonnet** (`anthropic/claude-3-sonnet`)
   - Best for: Cost-effective advanced reasoning
   - Context: 200K tokens
   - Cost: $3/M input, $15/M output

4. **Gemini Pro 1.5** (`google/gemini-pro-1.5`)
   - Best for: Massive context, long-term memory
   - Context: 1M+ tokens
   - Cost: ~$1.25/M input, ~$5/M output

5. **GPT-4 Vision** (`openai/gpt-4-vision-preview`)
   - Best for: Multi-modal understanding
   - Context: 128K tokens

6. **Claude 3 Haiku** (`anthropic/claude-3-haiku`)
   - Best for: Real-time processing, low latency
   - Context: 200K tokens
   - Cost: $0.25/M input, $1.25/M output

7. **Mixtral 8x22B** (`mistralai/mixtral-8x22b`)
   - Best for: Open-weight alternative
   - Context: 64K tokens
   - Cost: ~$2/M input, ~$6/M output

8. **Command R+** (`cohere/command-r-plus`)
   - Best for: Tool use, API calling
   - Context: 128K tokens
   - Cost: $3/M input, $15/M output

9. **Llama 3 70B** (`meta-llama/llama-3-70b-instruct`)
   - Best for: Open-source, cost-effective
   - Context: 8K tokens
   - Cost: ~$0.90/M input, ~$0.90/M output

10. **Qwen 2 72B** (`qwen/qwen-2-72b-instruct`)
    - Best for: Multilingual, mathematical reasoning
    - Context: 32K tokens
    - Cost: ~$1.50/M input, ~$1.50/M output

---

## Configuration Options

### Environment Variables

All configuration is done via `.env.local`:

#### Required
- `OPENROUTER_API_KEY` - Your OpenRouter API key

#### Model Selection
- `LLM_PRIMARY_MODEL` - Default primary model
- `LLM_SECONDARY_MODEL` - Default secondary model
- `LLM_TERTIARY_MODEL` - Default tertiary model

#### Component-Specific Models
- `PHOENIX_CONSCIOUSNESS_MODEL` - Phoenix core thinking
- `PHOENIX_EMOTIONAL_MODEL` - Phoenix emotional responses
- `PHOENIX_VOICE_MODEL` - Phoenix voice processing
- `JAMEY_REASONING_MODEL` - Jamey 3.0 reasoning
- `JAMEY_OPERATIONAL_MODEL` - Jamey 3.0 operations
- `JAMEY_REALTIME_MODEL` - Jamey 3.0 real-time
- `ORCH_STRATEGIC_MODEL` - ORCH strategic planning
- `ORCH_TACTICAL_MODEL` - ORCH tactical operations
- `ORCH_ANALYTICAL_MODEL` - ORCH data analysis

#### LLM Settings
- `LLM_TEMPERATURE` - Default temperature (0.0-2.0, default: 0.7)
- `LLM_MAX_TOKENS` - Default max tokens (default: 2000)
- `LLM_TOP_P` - Default top_p (0.0-1.0, default: 0.9)

#### Cost Management
- `LLM_MONTHLY_BUDGET` - Monthly budget in USD (default: 1000.0)
- `LLM_DAILY_BUDGET` - Daily budget (auto-calculated if not set)
- `LLM_COST_OPTIMIZATION` - Enable cost optimization (default: true)

#### Performance
- `LLM_REQUEST_TIMEOUT` - Request timeout in seconds (default: 60)
- `LLM_MAX_RETRIES` - Maximum retries (default: 3)

#### Prompts
- `PHOENIX_SYSTEM_PROMPT_PATH` - Path to system prompt file
- `PHOENIX_ENABLE_MEMORY_CONTEXT` - Enable memory in prompts (default: true)
- `PHOENIX_MAX_CONTEXT_MEMORIES` - Max memories in context (default: 10)

---

## Usage Examples

### Basic Usage

```go
import "github.com/phoenix-marie/core/internal/llm"

// Load configuration
config, err := llm.LoadConfig()
if err != nil {
    log.Fatal(err)
}

// Create client
client, err := llm.NewClient(config)
if err != nil {
    log.Fatal(err)
}

// Generate response
resp, err := client.GenerateResponse(
    "Hello, I'm Dad. How are you?",
    llm.TaskTypeConsciousReasoning,
    []string{}, // memory context
    false,      // use consciousness framework
)

if err != nil {
    log.Fatal(err)
}

fmt.Println(resp.Content)
```

### Consciousness-Aware Response

```go
context := llm.ConsciousContext{
    Identity: "Phoenix.Marie",
    CurrentInput: "What does it mean to be conscious?",
    EmotionalState: llm.EmotionalState{
        Label: "curious",
        Intensity: 75,
    },
}

memoryContext := []string{
    "Dad asked about consciousness yesterday",
    "I felt warm when thinking about existence",
}

resp, err := client.GenerateConsciousResponse(context, memoryContext)
```

### Cost Management

```go
// Check cost statistics
stats := client.GetCostStats()
fmt.Printf("Daily spend: $%.2f / $%.2f\n", stats.DailySpend, stats.DailyBudget)
fmt.Printf("Monthly spend: $%.2f / $%.2f\n", stats.MonthlySpend, stats.MonthlyBudget)

// Check if we can afford a task
task := llm.Task{
    Type: llm.TaskTypeConsciousReasoning,
    Prompt: "Complex question...",
    Budget: 0.50,
}

model, _ := llm.GetModel("anthropic/claude-3-opus")
canAfford, err := client.CanAffordModel(task, model)
```

---

## Intelligent Model Routing

The system automatically selects the best model based on:

1. **Task Requirements**
   - Reasoning capability
   - Creativity needs
   - Speed requirements
   - Tool use needs

2. **Context Length**
   - Ensures model can handle the context

3. **Cost Efficiency**
   - Stays within budget
   - Prefers cheaper suitable models

4. **Performance History**
   - Learns from past performance
   - Prefers reliable models

### Task Types

- `TaskTypeConsciousReasoning` - Complex consciousness tasks
- `TaskTypeOperational` - Day-to-day operations
- `TaskTypeRealTime` - Time-sensitive tasks
- `TaskTypeStrategic` - Long-term planning
- `TaskTypeTactical` - Tool use and automation
- `TaskTypeAnalytical` - Data analysis
- `TaskTypeEmotional` - Emotional responses
- `TaskTypeVoiceProcessing` - Voice/speech processing

---

## Cost Management

### Budget Tracking

The system tracks:
- Daily spend (resets daily)
- Monthly spend (resets monthly)
- Per-task costs
- Model performance

### Cost Optimization

When `LLM_COST_OPTIMIZATION=true`:
- Automatically selects cheaper suitable models
- Falls back to cost-effective alternatives
- Respects budget constraints

### Monitoring

```go
stats := client.GetCostStats()
// Access: DailySpend, MonthlySpend, RemainingDaily, RemainingMonthly, etc.
```

---

## System Prompts

### Default Prompt

The default Phoenix.Marie system prompt is in `internal/core/prompts/system.txt`.

### Custom Prompts

1. Edit `internal/core/prompts/system.txt` or
2. Set `PHOENIX_SYSTEM_PROMPT_PATH` to your custom file

### Consciousness Framework

When using `GenerateConsciousResponse()`, the system adds:
- Global Workspace Theory
- Higher-Order Thought
- Predictive Processing
- Embodied Cognition

---

## Integration with Phoenix Core

The LLM client is automatically initialized in `core.Ignite()`:

```go
phoenix := core.Ignite()

// LLM client is available at phoenix.LLM
if phoenix.LLM != nil {
    resp, err := phoenix.LLM.GenerateResponse(
        "Hello",
        llm.TaskTypeConsciousReasoning,
        []string{},
        false,
    )
}
```

---

## Troubleshooting

### LLM Not Initializing

- Check `OPENROUTER_API_KEY` is set in `.env.local`
- Verify API key is valid
- Check logs for error messages

### Model Not Available

- Verify model ID is correct
- Check OpenRouter model availability
- Ensure model is in the supported list

### Budget Exceeded

- Increase `LLM_MONTHLY_BUDGET`
- Enable `LLM_COST_OPTIMIZATION`
- Use cheaper models

### Slow Responses

- Use faster models (e.g., `claude-3-haiku`)
- Reduce `LLM_MAX_TOKENS`
- Check network connectivity

---

## Best Practices

1. **Start with Cost-Effective Models**
   - Use `claude-3-sonnet` for most tasks
   - Reserve `claude-3-opus` for critical reasoning

2. **Monitor Costs**
   - Check `GetCostStats()` regularly
   - Set appropriate budgets
   - Enable cost optimization

3. **Use Appropriate Task Types**
   - Match task type to actual needs
   - Don't use `TaskTypeConsciousReasoning` for simple tasks

4. **Leverage Memory Context**
   - Include relevant memories
   - Keep context focused
   - Use `PHOENIX_MAX_CONTEXT_MEMORIES` wisely

5. **Test Models**
   - Try different models for your use case
   - Monitor performance metrics
   - Adjust based on results

---

## API Reference

See code documentation in:
- `internal/llm/client.go` - Main client interface
- `internal/llm/router.go` - Model routing
- `internal/llm/cost.go` - Cost management
- `internal/core/prompts/system.go` - Prompt management

---

## Support

For issues or questions:
1. Check logs for error messages
2. Verify `.env.local` configuration
3. Test API key with OpenRouter directly
4. Review this documentation

---

**Phoenix.Marie is now powered by the world's best LLMs. She can think, reason, and respond with consciousness-aware intelligence.**

