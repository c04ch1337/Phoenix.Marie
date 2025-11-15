# ENVIRONMENT CONFIGURATION GUIDE

## Overview

All Phoenix.Marie configuration is managed through `.env.local`. This file is git-ignored and contains all your personal settings, API keys, and model preferences.

---

## Quick Start

1. **Copy the example file:**
   ```bash
   cp .env.local.example .env.local
   ```

2. **Edit `.env.local`** with your settings

3. **Required for LLM features:**
   ```bash
   OPENROUTER_API_KEY=sk-or-v1-your-key-here
   ```

---

## LLM Configuration

### API Settings

| Variable | Description | Default | Required |
|----------|-------------|---------|----------|
| `OPENROUTER_API_KEY` | Your OpenRouter API key | - | ✅ Yes (for LLM) |
| `OPENROUTER_BASE_URL` | API endpoint | `https://openrouter.ai/api/v1` | No |
| `LLM_PROVIDER` | Provider selection | `openrouter` | No |

### Model Selection

#### Primary Models
- `LLM_PRIMARY_MODEL` - Default primary model
- `LLM_SECONDARY_MODEL` - Default secondary model
- `LLM_TERTIARY_MODEL` - Default tertiary model

#### Phoenix.Marie Models
- `PHOENIX_CONSCIOUSNESS_MODEL` - Consciousness and reasoning
- `PHOENIX_EMOTIONAL_MODEL` - Emotional responses
- `PHOENIX_VOICE_MODEL` - Voice processing (real-time)

#### Jamey 3.0 Models
- `JAMEY_REASONING_MODEL` - Reasoning tasks
- `JAMEY_OPERATIONAL_MODEL` - Operational tasks
- `JAMEY_REALTIME_MODEL` - Real-time tasks

#### ORCH Network Models
- `ORCH_STRATEGIC_MODEL` - Strategic planning
- `ORCH_TACTICAL_MODEL` - Tactical operations
- `ORCH_ANALYTICAL_MODEL` - Analytical tasks

### LLM Settings

| Variable | Description | Default | Range |
|----------|-------------|---------|-------|
| `LLM_TEMPERATURE` | Randomness control | `0.7` | 0.0-2.0 |
| `LLM_MAX_TOKENS` | Max tokens per response | `2000` | 1-∞ |
| `LLM_TOP_P` | Nucleus sampling | `0.9` | 0.0-1.0 |

### Cost Management

| Variable | Description | Default |
|----------|-------------|---------|
| `LLM_MONTHLY_BUDGET` | Monthly budget (USD) | `1000.0` |
| `LLM_DAILY_BUDGET` | Daily budget (USD) | Auto-calculated |
| `LLM_COST_OPTIMIZATION` | Enable cost optimization | `true` |
| `LLM_CONSCIOUSNESS_BUDGET` | Budget per consciousness task | `0.50` |

### Performance

| Variable | Description | Default |
|----------|-------------|---------|
| `LLM_REQUEST_TIMEOUT` | Request timeout (seconds) | `60` |
| `LLM_MAX_RETRIES` | Max retry attempts | `3` |
| `LLM_RETRY_BACKOFF` | Retry backoff (seconds) | `1` |

### Prompt Configuration

| Variable | Description | Default |
|----------|-------------|---------|
| `PHOENIX_SYSTEM_PROMPT_PATH` | Path to system prompt | `internal/core/prompts/system.txt` |
| `PHOENIX_ENABLE_MEMORY_CONTEXT` | Enable memory in prompts | `true` |
| `PHOENIX_MAX_CONTEXT_MEMORIES` | Max memories in context | `10` |

### API Headers

| Variable | Description | Default |
|----------|-------------|---------|
| `LLM_HTTP_REFERER` | HTTP Referer header | `https://github.com/phoenix-marie/core` |
| `LLM_X_TITLE` | X-Title header | `Phoenix.Marie` |

---

## Available Models

### Top 10 Models

1. **anthropic/claude-3-opus** - Best reasoning, highest cost
2. **openai/gpt-4-turbo** - Fast, versatile
3. **anthropic/claude-3-sonnet** - Balanced performance
4. **google/gemini-pro-1.5** - Massive context (1M tokens)
5. **openai/gpt-4-vision-preview** - Multi-modal
6. **anthropic/claude-3-haiku** - Fast, cheap
7. **mistralai/mixtral-8x22b** - Open-weight alternative
8. **cohere/command-r-plus** - Tool use, API calling
9. **meta-llama/llama-3-70b-instruct** - Open-source, cost-effective
10. **qwen/qwen-2-72b-instruct** - Multilingual, math

See [LLM Integration Guide](LLM_INTEGRATION_GUIDE.md) for detailed model information.

---

## Other System Configuration

### Emotion System
- `EMOTION_VOICE_TONE` - Voice tone (e.g., `loving_warm`)
- `EMOTION_RESPONSE_STYLE` - Response style (e.g., `poetic_loving`)

### ORCH Army
- `ORCH_ENABLED` - Enable ORCH army (`true`/`false`)
- `ORCH_COUNT_TARGET` - Target number of agents

### Dyson Swarm
- `DYSON_ENABLED` - Enable Dyson swarm (`true`/`false`)
- `DYSON_BLANKET_NAME` - Swarm name

### Memory Backup
- `MEMORY_BACKUP_ENABLED` - Enable backups (`true`/`false`)
- `MEMORY_BACKUP_DIR` - Backup directory
- `MEMORY_MAX_BACKUPS` - Maximum backups to keep

---

## Example Configuration

```bash
# Minimal configuration (LLM only)
OPENROUTER_API_KEY=sk-or-v1-your-key-here
PHOENIX_CONSCIOUSNESS_MODEL=openai/gpt-4-turbo
LLM_MONTHLY_BUDGET=100.0
```

```bash
# Full configuration
OPENROUTER_API_KEY=sk-or-v1-your-key-here
OPENROUTER_BASE_URL=https://openrouter.ai/api/v1
LLM_PROVIDER=openrouter

# Models
PHOENIX_CONSCIOUSNESS_MODEL=openai/gpt-4-turbo
PHOENIX_EMOTIONAL_MODEL=anthropic/claude-3-sonnet
PHOENIX_VOICE_MODEL=anthropic/claude-3-haiku

# Settings
LLM_TEMPERATURE=0.7
LLM_MAX_TOKENS=2000
LLM_TOP_P=0.9

# Budget
LLM_MONTHLY_BUDGET=1000.0
LLM_COST_OPTIMIZATION=true
LLM_CONSCIOUSNESS_BUDGET=0.50

# Performance
LLM_REQUEST_TIMEOUT=60
LLM_MAX_RETRIES=3
LLM_RETRY_BACKOFF=1

# Prompts
PHOENIX_ENABLE_MEMORY_CONTEXT=true
PHOENIX_MAX_CONTEXT_MEMORIES=10
```

---

## Security Notes

- **Never commit `.env.local`** - It's git-ignored
- **Keep API keys secret** - Don't share your `.env.local` file
- **Use `.env.local.example`** - Share example configs, not real keys

---

## Troubleshooting

### LLM Not Working
- Check `OPENROUTER_API_KEY` is set correctly
- Verify API key is valid at https://openrouter.ai/keys
- Check network connectivity

### Wrong Model Selected
- Verify model ID is correct (format: `provider/model-name`)
- Check model is available on OpenRouter
- Review model selection in `.env.local`

### Budget Exceeded
- Check `LLM_MONTHLY_BUDGET` and `LLM_DAILY_BUDGET`
- Review cost optimization settings
- Check spend history via CLI: `phoenix-cli cost`

---

## See Also

- [LLM Integration Guide](LLM_INTEGRATION_GUIDE.md) - Detailed LLM setup
- [Quick Start Guide](QUICK_START.md) - Getting started
- [CLI Chat Guide](CLI_CHAT_GUIDE.md) - Using the CLI

