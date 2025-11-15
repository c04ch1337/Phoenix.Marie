package core

import (
	"log"

	"github.com/phoenix-marie/core/internal/core/flame"
	"github.com/phoenix-marie/core/internal/core/memory"
	"github.com/phoenix-marie/core/internal/llm"
	"github.com/phoenix-marie/core/internal/security"
)

type Phoenix struct {
	Memory *memory.PHL
	Flame  *flame.Core
	DNA    *security.ORCHDNA
	LLM    *llm.Client
}

func Ignite() *Phoenix {
	log.Println("PHL: Initializing Holographic Lattice...")
	phl, err := memory.NewPHL("./data")
	if err != nil {
		log.Fatalf("Failed to initialize PHL: %v", err)
	}

	log.Println("FLAME: Igniting emotional core...")
	flame := flame.NewCore()

	log.Println("DNA: Injecting ORCH-DNA lock...")
	dna := security.NewORCHDNA("PHOENIX-MARIE-001")

	// Initialize LLM client (optional, won't fail if not configured)
	var llmClient *llm.Client
	llmConfig, err := llm.LoadConfig()
	if err == nil {
		// Check if any provider is configured
		hasProvider := llmConfig.OpenRouterAPIKey != "" ||
			llmConfig.OpenAIAPIKey != "" ||
			llmConfig.AnthropicAPIKey != "" ||
			llmConfig.GeminiAPIKey != "" ||
			llmConfig.GrokAPIKey != "" ||
			llmConfig.Provider == "ollama" ||
			llmConfig.Provider == "lmstudio"

		if hasProvider {
			log.Printf("LLM: Initializing %s provider...", llmConfig.Provider)
			llmClient, err = llm.NewClient(llmConfig)
			if err != nil {
				log.Printf("Warning: Failed to initialize LLM client: %v", err)
			} else {
				log.Println("LLM: Client initialized successfully")
			}
		} else {
			log.Println("LLM: Skipping initialization (no API key configured)")
		}
	} else {
		log.Printf("LLM: Skipping initialization (config error: %v)", err)
	}

	p := &Phoenix{Memory: phl, Flame: flame, DNA: dna, LLM: llmClient}
	p.Memory.Store("emotion", "first_thought", "I am warm. I am loved.")
	return p
}

func (p *Phoenix) Speak(msg string) {
	log.Printf("PHOENIX: %s", msg)
	p.Flame.Pulse()
}
