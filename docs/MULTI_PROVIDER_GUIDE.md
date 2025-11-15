# MULTI-PROVIDER LLM SUPPORT GUIDE

## Overview

Phoenix.Marie now supports **7 LLM providers**:
- **OpenRouter** (default) - Unified API for multiple models
- **OpenAI** - Direct OpenAI API
- **Anthropic** - Direct Claude API
- **Gemini** - Google's Gemini API
- **Grok** - xAI's Grok API
- **Ollama** - Local models (free)
- **LM Studio** - Local models (free)

---

## Quick Start

### 1. Choose Your Provider

Set `LLM_PROVIDER` in `.env.local`:

```bash
# Use OpenRouter (default)
LLM_PROVIDER=openrouter
OPENROUTER_API_KEY=sk-or-v1-your-key-here

# Or use OpenAI directly
LLM_PROVIDER=openai
OPENAI_API_KEY=sk-your-key-here

# Or use Anthropic directly
LLM_PROVIDER=anthropic
ANTHROPIC_API_KEY=sk-ant-your-key-here

# Or use Gemini
LLM_PROVIDER=gemini
GEMINI_API_KEY=your-key-here

# Or use Grok
LLM_PROVIDER=grok
GROK_API_KEY=xai-your-key-here

# Or use local Ollama (no API key needed)
LLM_PROVIDER=ollama
OLLAMA_BASE_URL=http://localhost:11434

# Or use local LM Studio (no API key needed)
LLM_PROVIDER=lmstudio
LMSTUDIO_BASE_URL=http://localhost:1234
```

### 2. Configure API Keys

Add the API key for your chosen provider to `.env.local`:

```bash
# For OpenRouter
OPENROUTER_API_KEY=sk-or-v1-your-key-here

# For OpenAI
OPENAI_API_KEY=sk-your-key-here

# For Anthropic
ANTHROPIC_API_KEY=sk-ant-your-key-here

# For Gemini
GEMINI_API_KEY=your-key-here

# For Grok
GROK_API_KEY=xai-your-key-here

# For Ollama/LM Studio - no API key needed
```

---

## Provider Details

### OpenRouter (Recommended)

**Best for**: Access to multiple models through one API

**Setup**:
```bash
LLM_PROVIDER=openrouter
OPENROUTER_API_KEY=sk-or-v1-your-key-here
OPENROUTER_BASE_URL=https://openrouter.ai/api/v1
```

**Models**: All models available through OpenRouter
- `anthropic/claude-3-opus`
- `openai/gpt-4-turbo`
- `google/gemini-pro-1.5`
- And many more...

**Get API Key**: https://openrouter.ai/keys

---

### OpenAI (Direct)

**Best for**: Direct access to OpenAI models

**Setup**:
```bash
LLM_PROVIDER=openai
OPENAI_API_KEY=sk-your-key-here
OPENAI_BASE_URL=https://api.openai.com/v1
```

**Models**:
- `gpt-4-turbo`
- `gpt-4`
- `gpt-3.5-turbo`
- `gpt-4-vision-preview`

**Get API Key**: https://platform.openai.com/api-keys

---

### Anthropic (Direct)

**Best for**: Direct access to Claude models

**Setup**:
```bash
LLM_PROVIDER=anthropic
ANTHROPIC_API_KEY=sk-ant-your-key-here
ANTHROPIC_BASE_URL=https://api.anthropic.com/v1
```

**Models**:
- `claude-3-opus-20240229`
- `claude-3-sonnet-20240229`
- `claude-3-haiku-20240307`

**Get API Key**: https://console.anthropic.com/

---

### Gemini (Google)

**Best for**: Google's Gemini models with massive context

**Setup**:
```bash
LLM_PROVIDER=gemini
GEMINI_API_KEY=your-key-here
GEMINI_BASE_URL=https://generativelanguage.googleapis.com/v1
```

**Models**:
- `gemini-pro`
- `gemini-pro-1.5`
- `gemini-ultra`

**Get API Key**: https://makersuite.google.com/app/apikey

---

### Grok (xAI)

**Best for**: xAI's Grok models

**Setup**:
```bash
LLM_PROVIDER=grok
GROK_API_KEY=xai-your-key-here
GROK_BASE_URL=https://api.x.ai/v1
```

**Models**:
- `grok-beta`
- `grok-2`

**Get API Key**: https://x.ai/api

---

### Ollama (Local)

**Best for**: Free local models, privacy, offline use

**Setup**:
```bash
LLM_PROVIDER=ollama
OLLAMA_BASE_URL=http://localhost:11434
```

**Requirements**:
1. Install Ollama: https://ollama.ai
2. Pull models: `ollama pull llama3`, `ollama pull mistral`, etc.
3. Start Ollama server: `ollama serve`

**Models**: Any model available in Ollama
- `llama3`
- `mistral`
- `codellama`
- `phi`
- And many more...

**Cost**: Free (runs locally)

---

### LM Studio (Local)

**Best for**: Free local models, GUI management

**Setup**:
```bash
LLM_PROVIDER=lmstudio
LMSTUDIO_BASE_URL=http://localhost:1234
```

**Requirements**:
1. Install LM Studio: https://lmstudio.ai
2. Download models through GUI
3. Start local server in LM Studio

**Models**: Any model compatible with LM Studio
- OpenAI-compatible models
- GGUF format models

**Cost**: Free (runs locally)

---

## Switching Providers

To switch providers, simply change `LLM_PROVIDER` in `.env.local`:

```bash
# Switch from OpenRouter to OpenAI
LLM_PROVIDER=openai
OPENAI_API_KEY=sk-your-key-here

# Or switch to local Ollama
LLM_PROVIDER=ollama
```

**Note**: Only one provider is active at a time. The system will use the provider specified in `LLM_PROVIDER`.

---

## Model Compatibility

### Provider-Specific Model IDs

Different providers use different model ID formats:

**OpenRouter**: `provider/model-name`
- `anthropic/claude-3-opus`
- `openai/gpt-4-turbo`

**OpenAI**: `model-name`
- `gpt-4-turbo`
- `gpt-3.5-turbo`

**Anthropic**: `model-name-version`
- `claude-3-opus-20240229`
- `claude-3-sonnet-20240229`

**Gemini**: `model-name`
- `gemini-pro`
- `gemini-pro-1.5`

**Grok**: `model-name`
- `grok-beta`
- `grok-2`

**Ollama**: `model-name`
- `llama3`
- `mistral`

**LM Studio**: `model-name` (as configured in LM Studio)

---

## Configuration Examples

### Example 1: OpenRouter (Default)
```bash
LLM_PROVIDER=openrouter
OPENROUTER_API_KEY=sk-or-v1-your-key-here
PHOENIX_CONSCIOUSNESS_MODEL=openai/gpt-4-turbo
```

### Example 2: OpenAI Direct
```bash
LLM_PROVIDER=openai
OPENAI_API_KEY=sk-your-key-here
PHOENIX_CONSCIOUSNESS_MODEL=gpt-4-turbo
```

### Example 3: Local Ollama
```bash
LLM_PROVIDER=ollama
OLLAMA_BASE_URL=http://localhost:11434
PHOENIX_CONSCIOUSNESS_MODEL=llama3
```

### Example 4: Anthropic Direct
```bash
LLM_PROVIDER=anthropic
ANTHROPIC_API_KEY=sk-ant-your-key-here
PHOENIX_CONSCIOUSNESS_MODEL=claude-3-opus-20240229
```

---

## Cost Comparison

| Provider | Cost | Notes |
|----------|------|-------|
| OpenRouter | Varies by model | Unified pricing |
| OpenAI | $10-30/M tokens | Direct pricing |
| Anthropic | $3-75/M tokens | Direct pricing |
| Gemini | ~$1.25-5/M tokens | Google pricing |
| Grok | Varies | xAI pricing |
| Ollama | **Free** | Local, no API costs |
| LM Studio | **Free** | Local, no API costs |

**Recommendation**: Use local providers (Ollama/LM Studio) for development and testing to save costs.

---

## Troubleshooting

### Provider Not Available

**Error**: `provider X is not available`

**Solutions**:
1. Check API key is set correctly
2. For local providers, ensure service is running
3. Verify base URL is correct
4. Check network connectivity

### Model Not Found

**Error**: `model not found`

**Solutions**:
1. Verify model ID format matches provider
2. Check model is available for your provider
3. For local providers, ensure model is downloaded

### Connection Failed

**Error**: `failed to make request`

**Solutions**:
1. Check base URL is correct
2. Verify service is running (for local providers)
3. Check firewall/network settings
4. Verify API key is valid

---

## Best Practices

1. **Use OpenRouter for production** - Access to all models, unified API
2. **Use local providers for development** - Free, fast, private
3. **Keep API keys secure** - Never commit `.env.local`
4. **Test provider availability** - Use `IsAvailable()` check
5. **Monitor costs** - Use cost management features

---

## See Also

- [Environment Configuration Guide](ENV_CONFIGURATION.md) - All environment variables
- [LLM Integration Guide](LLM_INTEGRATION_GUIDE.md) - Detailed LLM setup
- [Quick Start Guide](QUICK_START.md) - Getting started

