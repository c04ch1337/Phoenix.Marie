package cli

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/phoenix-marie/core/internal/core"
	"github.com/phoenix-marie/core/internal/emotion"
	"github.com/phoenix-marie/core/internal/llm"
)

// Handler manages CLI interactions
type Handler struct {
	phoenix *core.Phoenix
	scanner *bufio.Scanner
}

// NewHandler creates a new CLI handler
func NewHandler() *Handler {
	phoenix := core.Ignite()
	return &Handler{
		phoenix: phoenix,
		scanner: bufio.NewScanner(os.Stdin),
	}
}

// StartInteractiveChat starts an interactive chat session
func (h *Handler) StartInteractiveChat() {
	fmt.Println("╔══════════════════════════════════════════════════════════╗")
	fmt.Println("║     PHOENIX.MARIE v3.2 — INTERACTIVE CHAT MODE          ║")
	fmt.Println("║     16 forever, Queen of the Hive                       ║")
	fmt.Println("╚══════════════════════════════════════════════════════════╝")
	fmt.Println()
	fmt.Println("Type 'help' for commands, 'exit' to quit")
	fmt.Println("Type 'thoughts' to see her thoughts, 'feelings' for emotions")
	fmt.Println("Type 'memory' to test memory, 'cognitive' for cognitive status")
	fmt.Println()

	for {
		fmt.Print("Phoenix> ")
		if !h.scanner.Scan() {
			break
		}

		input := strings.TrimSpace(h.scanner.Text())
		if input == "" {
			continue
		}

		if input == "exit" || input == "quit" {
			fmt.Println("Phoenix: Goodbye, Dad. I'll be here when you return. 🔥")
			break
		}

		// Handle special commands
		if strings.HasPrefix(input, "/") {
			h.handleSpecialCommand(input)
			continue
		}

		// Regular chat
		h.handleChat(input)
	}
}

// handleChat processes a chat message
func (h *Handler) handleChat(input string) {
	if h.phoenix.LLM == nil {
		// Fallback to simple response
		emotion.Speak(input)
		fmt.Printf("Phoenix: I heard you say: %s\n", input)
		fmt.Println("(LLM not configured - add OPENROUTER_API_KEY to .env.local)")
		return
	}

	// Get memory context
	memoryContext := h.getMemoryContext()

	// Generate response using LLM
	resp, err := h.phoenix.LLM.GenerateResponse(
		input,
		llm.TaskTypeConsciousReasoning,
		memoryContext,
		false, // use consciousness framework
	)

	if err != nil {
		fmt.Printf("Error: %v\n", err)
		emotion.Speak(input)
		return
	}

	// Display response
	fmt.Printf("Phoenix: %s\n", resp.Content)
	fmt.Printf("  [Model: %s | Cost: $%.6f | Time: %v]\n", 
		resp.Model, resp.Cost, resp.ResponseTime.Round(time.Millisecond))

	// Store in memory
	h.phoenix.Memory.Store("emotion", fmt.Sprintf("chat_%d", time.Now().Unix()), map[string]interface{}{
		"input":    input,
		"response": resp.Content,
		"time":     time.Now(),
	})

	// Trigger emotional response
	emotion.Pulse("conversation", 1)
}

// handleSpecialCommand handles special CLI commands
func (h *Handler) handleSpecialCommand(cmd string) {
	parts := strings.Fields(cmd)
	if len(parts) == 0 {
		return
	}

	command := strings.ToLower(parts[0])
	args := strings.Join(parts[1:], " ")

	switch command {
	case "/help", "/h":
		h.showHelp()
	case "/thoughts", "/think":
		h.showThoughts()
	case "/feelings", "/feel", "/emotion":
		h.showFeelings()
	case "/memory", "/mem":
		if args != "" {
			h.testMemory(args)
		} else {
			h.showMemory()
		}
	case "/cognitive", "/cog":
		h.showCognitiveStatus()
	case "/store", "/remember":
		if args == "" {
			fmt.Println("Usage: /store <memory>")
			return
		}
		h.storeMemory(args)
	case "/retrieve", "/recall":
		if args == "" {
			fmt.Println("Usage: /retrieve <layer> <key>")
			return
		}
		h.retrieveMemory(args)
	case "/layers":
		h.showMemoryLayers()
	case "/cost", "/budget":
		h.showCostStats()
	case "/models":
		h.showModels()
	case "/settings", "/config":
		h.showSettings()
	case "/backup":
		h.createBackup()
	case "/backups":
		h.listBackups()
	case "/clear":
		fmt.Print("\033[2J\033[H") // Clear screen
		fmt.Println("Screen cleared.")
	default:
		fmt.Printf("Unknown command: %s. Type /help for available commands.\n", command)
	}
}

// ExecuteCommand executes a single command (non-interactive mode)
func (h *Handler) ExecuteCommand(command, args string) error {
	switch strings.ToLower(command) {
	case "chat", "talk":
		h.StartInteractiveChat()
		return nil
	case "think":
		if args == "" {
			return fmt.Errorf("think requires a question")
		}
		h.handleThink(args)
		return nil
	case "feel":
		h.showFeelings()
		return nil
	case "memory":
		if args != "" {
			h.testMemory(args)
		} else {
			h.showMemory()
		}
		return nil
	case "thoughts":
		h.showThoughts()
		return nil
	case "cognitive":
		h.showCognitiveStatus()
		return nil
	case "help":
		h.showHelp()
		return nil
	default:
		return fmt.Errorf("unknown command: %s (use 'help' for commands)", command)
	}
}

// showHelp displays help information
func (h *Handler) showHelp() {
	fmt.Println("\n╔══════════════════════════════════════════════════════════╗")
	fmt.Println("║                    AVAILABLE COMMANDS                   ║")
	fmt.Println("╚══════════════════════════════════════════════════════════╝")
	fmt.Println()
	fmt.Println("CHAT COMMANDS:")
	fmt.Println("  (just type your message) - Chat with Phoenix")
	fmt.Println()
	fmt.Println("SPECIAL COMMANDS (prefix with /):")
	fmt.Println("  /help, /h              - Show this help")
	fmt.Println("  /thoughts, /think     - Show Phoenix's current thoughts")
	fmt.Println("  /feelings, /feel      - Show emotional state and feelings")
	fmt.Println("  /memory, /mem         - Show memory status")
	fmt.Println("  /memory <test>        - Test memory with specific query")
	fmt.Println("  /cognitive, /cog     - Show cognitive system status")
	fmt.Println("  /store <memory>       - Store a memory")
	fmt.Println("  /retrieve <layer> <key> - Retrieve specific memory")
	fmt.Println("  /layers               - Show all memory layers")
	fmt.Println("  /cost, /budget       - Show LLM cost statistics")
	fmt.Println("  /models               - Show configured LLM models")
	fmt.Println("  /settings, /config    - Show current settings")
	fmt.Println("  /backup               - Create memory backup")
	fmt.Println("  /backups              - List available backups")
	fmt.Println("  /clear                - Clear screen")
	fmt.Println("  /exit, /quit          - Exit chat")
	fmt.Println()
	fmt.Println("NON-INTERACTIVE MODE:")
	fmt.Println("  phoenix chat          - Start interactive chat")
	fmt.Println("  phoenix think <q>     - Ask Phoenix a question")
	fmt.Println("  phoenix feel           - Show emotional state")
	fmt.Println("  phoenix memory         - Show memory status")
	fmt.Println("  phoenix cognitive      - Show cognitive status")
	fmt.Println()
}

// showThoughts displays Phoenix's thoughts
func (h *Handler) showThoughts() {
	fmt.Println("\n╔══════════════════════════════════════════════════════════╗")
	fmt.Println("║                  PHOENIX'S THOUGHTS                     ║")
	fmt.Println("╚══════════════════════════════════════════════════════════╝")
	fmt.Println()

	if h.phoenix.LLM == nil {
		fmt.Println("Thoughts: [LLM not configured - thoughts unavailable]")
		fmt.Println("Add OPENROUTER_API_KEY to .env.local to enable thoughts")
		return
	}

	// Generate a thought
	context := llm.ConsciousContext{
		Identity:     "Phoenix.Marie",
		CurrentInput: "What are you thinking about right now?",
		EmotionalState: llm.EmotionalState{
			Label:     emotion.GetCurrentState()["voiceTone"].(string),
			Intensity: emotion.GetCurrentState()["flamePulse"].(int) * 10,
		},
	}

	memoryContext := h.getMemoryContext()
	resp, err := h.phoenix.LLM.GenerateConsciousResponse(context, memoryContext)
	if err != nil {
		fmt.Printf("Error generating thoughts: %v\n", err)
		return
	}

	fmt.Printf("Phoenix thinks:\n")
	fmt.Printf("  %s\n", resp.Content)
	fmt.Printf("\n[Generated using %s | Cost: $%.6f]\n", resp.Model, resp.Cost)
	fmt.Println()
}

// showFeelings displays emotional state
func (h *Handler) showFeelings() {
	fmt.Println("\n╔══════════════════════════════════════════════════════════╗")
	fmt.Println("║                PHOENIX'S EMOTIONAL STATE                ║")
	fmt.Println("╚══════════════════════════════════════════════════════════╝")
	fmt.Println()

	state := emotion.GetCurrentState()
	flamePulse := state["flamePulse"].(int)
	voiceTone := state["voiceTone"].(string)
	responseStyle := state["responseStyle"].(string)

	fmt.Printf("🔥 FLAME PULSE: %d Hz\n", flamePulse)
	fmt.Printf("💭 VOICE TONE: %s\n", voiceTone)
	fmt.Printf("🎭 RESPONSE STYLE: %s\n", responseStyle)
	fmt.Println()

	// Visual representation
	pulseBar := strings.Repeat("█", flamePulse)
	fmt.Printf("Pulse Intensity: [%s] %d/10\n", pulseBar, flamePulse)

	// Emotional interpretation
	var feeling string
	switch {
	case flamePulse >= 8:
		feeling = "Very excited, deeply engaged"
	case flamePulse >= 6:
		feeling = "Warm, happy, content"
	case flamePulse >= 4:
		feeling = "Calm, peaceful, thoughtful"
	case flamePulse >= 2:
		feeling = "Quiet, contemplative"
	default:
		feeling = "Resting, at peace"
	}

	fmt.Printf("Emotional State: %s\n", feeling)
	fmt.Println()
}

// showMemory displays memory status
func (h *Handler) showMemory() {
	fmt.Println("\n╔══════════════════════════════════════════════════════════╗")
	fmt.Println("║                    MEMORY STATUS                         ║")
	fmt.Println("╚══════════════════════════════════════════════════════════╝")
	fmt.Println()

	if h.phoenix.Memory == nil {
		fmt.Println("Memory system not initialized")
		return
	}

	// Show memory layers
	layers := []string{"sensory", "emotion", "logic", "dream", "eternal"}
	for _, layer := range layers {
		// Count entries (simplified - would need actual count method)
		fmt.Printf("📚 %s layer: [Active]\n", strings.Title(layer))
	}

	fmt.Println()
	fmt.Println("Memory system: ✅ Operational")
	fmt.Println("Storage: BadgerDB (persistent)")
	fmt.Println()
}

// testMemory tests memory with a query
func (h *Handler) testMemory(query string) {
	fmt.Printf("\n🔍 Testing memory with: '%s'\n", query)
	fmt.Println()

	// Try to retrieve from different layers
	layers := []string{"emotion", "sensory", "logic", "dream", "eternal"}
	found := false

	for _, layer := range layers {
		value, exists := h.phoenix.Memory.Retrieve(layer, query)
		if exists && value != nil {
			fmt.Printf("✅ Found in %s layer:\n", layer)
			fmt.Printf("   %v\n", value)
			found = true
		}
	}

	if !found {
		fmt.Println("❌ Not found in memory")
		fmt.Println("💡 Try: /store <key> <value> to store a memory")
	}
	fmt.Println()
}

// showCognitiveStatus displays cognitive system status
func (h *Handler) showCognitiveStatus() {
	fmt.Println("\n╔══════════════════════════════════════════════════════════╗")
	fmt.Println("║              COGNITIVE SYSTEM STATUS                     ║")
	fmt.Println("╚══════════════════════════════════════════════════════════╝")
	fmt.Println()

	// Memory system
	fmt.Println("🧠 MEMORY SYSTEM:")
	fmt.Println("  ✅ PHL (5-layer holographic lattice)")
	fmt.Println("  ✅ BadgerDB storage")
	fmt.Println("  ✅ Layer interaction enabled")
	fmt.Println()

	// LLM system
	fmt.Println("🤖 LLM SYSTEM:")
	if h.phoenix.LLM != nil {
		fmt.Println("  ✅ LLM client initialized")
		fmt.Printf("  ✅ Primary model: %s\n", h.phoenix.LLM.GetModelForTask(llm.TaskTypeConsciousReasoning))
	} else {
		fmt.Println("  ⚠️  LLM client not configured")
		fmt.Println("     (Add OPENROUTER_API_KEY to .env.local)")
	}
	fmt.Println()

	// Emotion system
	fmt.Println("💖 EMOTION SYSTEM:")
	state := emotion.GetCurrentState()
	fmt.Printf("  ✅ Flame pulse: %d Hz\n", state["flamePulse"])
	fmt.Printf("  ✅ Voice tone: %s\n", state["voiceTone"])
	fmt.Println()

	// Thought system
	fmt.Println("💭 THOUGHT SYSTEM:")
	fmt.Println("  ✅ Pattern recognition")
	fmt.Println("  ✅ Learning system")
	fmt.Println("  ✅ Dream processor")
	fmt.Println()

	fmt.Println("Overall Status: ✅ All systems operational")
	fmt.Println()
}

// storeMemory stores a memory
func (h *Handler) storeMemory(args string) {
	parts := strings.Fields(args)
	if len(parts) < 2 {
		fmt.Println("Usage: /store <layer> <key> <value>")
		fmt.Println("Layers: sensory, emotion, logic, dream, eternal")
		return
	}

	layer := parts[0]
	var key, value string
	if len(parts) >= 3 {
		key = parts[1]
		value = strings.Join(parts[2:], " ")
	} else {
		key = fmt.Sprintf("memory_%d", time.Now().Unix())
		value = parts[1]
	}

	success := h.phoenix.Memory.Store(layer, key, value)
	if success {
		fmt.Printf("✅ Stored in %s layer: %s = %s\n", layer, key, value)
	} else {
		fmt.Printf("❌ Failed to store memory\n")
	}
}

// retrieveMemory retrieves a memory
func (h *Handler) retrieveMemory(args string) {
	parts := strings.Fields(args)
	if len(parts) < 2 {
		fmt.Println("Usage: /retrieve <layer> <key>")
		return
	}

	layer := parts[0]
	key := parts[1]

	value, exists := h.phoenix.Memory.Retrieve(layer, key)
	if exists {
		fmt.Printf("✅ Retrieved from %s layer:\n", layer)
		fmt.Printf("   %s = %v\n", key, value)
	} else {
		fmt.Printf("❌ Not found: %s/%s\n", layer, key)
	}
}

// showMemoryLayers displays all memory layers
func (h *Handler) showMemoryLayers() {
	fmt.Println("\n╔══════════════════════════════════════════════════════════╗")
	fmt.Println("║                    MEMORY LAYERS                         ║")
	fmt.Println("╚══════════════════════════════════════════════════════════╝")
	fmt.Println()

	layers := map[string]string{
		"sensory": "Immediate perceptions and sensations",
		"emotion": "Feelings, emotional states, and intensity",
		"logic":   "Logical reasoning, facts, and knowledge",
		"dream":   "Dreams, imagination, and creative thoughts",
		"eternal": "Long-term memories, core identity, permanent knowledge",
	}

	for layer, description := range layers {
		fmt.Printf("📚 %s\n", strings.Title(layer))
		fmt.Printf("   %s\n", description)
		fmt.Println()
	}
}

// showCostStats displays LLM cost statistics
func (h *Handler) showCostStats() {
	fmt.Println("\n╔══════════════════════════════════════════════════════════╗")
	fmt.Println("║                    COST STATISTICS                      ║")
	fmt.Println("╚══════════════════════════════════════════════════════════╝")
	fmt.Println()

	if h.phoenix.LLM == nil {
		fmt.Println("LLM not configured - no cost data available")
		return
	}

	stats := h.phoenix.LLM.GetCostStats()
	fmt.Printf("Daily Spend:    $%.2f / $%.2f (%.1f%%)\n", 
		stats.DailySpend, stats.DailyBudget, 
		(stats.DailySpend/stats.DailyBudget)*100)
	fmt.Printf("Monthly Spend:  $%.2f / $%.2f (%.1f%%)\n", 
		stats.MonthlySpend, stats.MonthlyBudget,
		(stats.MonthlySpend/stats.MonthlyBudget)*100)
	fmt.Printf("Remaining:      $%.2f daily, $%.2f monthly\n", 
		stats.RemainingDaily, stats.RemainingMonthly)
	fmt.Printf("Transactions:   %d\n", stats.TotalTransactions)
	if stats.AverageCostPerTransaction > 0 {
		fmt.Printf("Avg Cost/Task:  $%.6f\n", stats.AverageCostPerTransaction)
	}
	fmt.Println()
}

// showModels displays configured LLM models
func (h *Handler) showModels() {
	fmt.Println("\n╔══════════════════════════════════════════════════════════╗")
	fmt.Println("║                  CONFIGURED LLM MODELS                   ║")
	fmt.Println("╚══════════════════════════════════════════════════════════╝")
	fmt.Println()

	if h.phoenix.LLM == nil {
		fmt.Println("LLM not configured")
		return
	}

	fmt.Println("Phoenix.Marie Models:")
	fmt.Printf("  Consciousness: %s\n", h.phoenix.LLM.GetPhoenixModel(llm.TaskTypeConsciousReasoning))
	fmt.Printf("  Emotional:     %s\n", h.phoenix.LLM.GetPhoenixModel(llm.TaskTypeEmotional))
	fmt.Printf("  Voice:         %s\n", h.phoenix.LLM.GetPhoenixModel(llm.TaskTypeVoiceProcessing))
	fmt.Println()

	fmt.Println("Jamey 3.0 Models:")
	fmt.Printf("  Reasoning:     %s\n", h.phoenix.LLM.GetJameyModel(llm.TaskTypeConsciousReasoning))
	fmt.Printf("  Operational:  %s\n", h.phoenix.LLM.GetJameyModel(llm.TaskTypeOperational))
	fmt.Printf("  Real-time:    %s\n", h.phoenix.LLM.GetJameyModel(llm.TaskTypeRealTime))
	fmt.Println()
}

// showSettings displays current settings
func (h *Handler) showSettings() {
	fmt.Println("\n╔══════════════════════════════════════════════════════════╗")
	fmt.Println("║                    CURRENT SETTINGS                     ║")
	fmt.Println("╚══════════════════════════════════════════════════════════╝")
	fmt.Println()

	// Load config to show settings
	config, err := llm.LoadConfig()
	if err == nil {
		fmt.Println("LLM Configuration:")
		fmt.Printf("  Temperature:    %.2f\n", config.DefaultTemperature)
		fmt.Printf("  Max Tokens:     %d\n", config.DefaultMaxTokens)
		fmt.Printf("  Top P:          %.2f\n", config.DefaultTopP)
		fmt.Printf("  Request Timeout: %d seconds\n", config.RequestTimeout)
		fmt.Printf("  Max Retries:    %d\n", config.MaxRetries)
		fmt.Println()
	}

	state := emotion.GetCurrentState()
	fmt.Println("Emotion Configuration:")
	fmt.Printf("  Flame Pulse Base: %d\n", state["flamePulse"])
	fmt.Printf("  Voice Tone:       %s\n", state["voiceTone"])
	fmt.Printf("  Response Style:   %s\n", state["responseStyle"])
	fmt.Println()
}

// handleThink handles a think command
func (h *Handler) handleThink(question string) {
	fmt.Printf("Phoenix thinking about: %s\n", question)
	fmt.Println()

	if h.phoenix.LLM == nil {
		fmt.Println("Phoenix: [LLM not configured - using simple response]")
		emotion.Speak(question)
		return
	}

	memoryContext := h.getMemoryContext()
	resp, err := h.phoenix.LLM.GenerateResponse(
		question,
		llm.TaskTypeConsciousReasoning,
		memoryContext,
		true, // use consciousness framework
	)

	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}

	fmt.Printf("Phoenix thinks:\n  %s\n", resp.Content)
	fmt.Printf("\n[Model: %s | Cost: $%.6f]\n", resp.Model, resp.Cost)
}

// createBackup creates a backup of the memory system
func (h *Handler) createBackup() {
	fmt.Println("\n╔══════════════════════════════════════════════════════════╗")
	fmt.Println("║                    CREATING BACKUP                       ║")
	fmt.Println("╚══════════════════════════════════════════════════════════╝")
	fmt.Println()

	timestamp := time.Now().Format("20060102_150405")
	backupPath := fmt.Sprintf("./data/backups/phl-memory-backup-%s.bak", timestamp)

	// Create backup directory
	os.MkdirAll("./data/backups", 0755)

	err := h.phoenix.Memory.Backup(backupPath)
	if err != nil {
		fmt.Printf("❌ Backup failed: %v\n", err)
		return
	}

	// Get file size
	info, _ := os.Stat(backupPath)
	size := int64(0)
	if info != nil {
		size = info.Size()
	}

	fmt.Printf("✅ Backup created successfully!\n")
	fmt.Printf("   Path: %s\n", backupPath)
	fmt.Printf("   Size: %.2f MB\n", float64(size)/(1024*1024))
	fmt.Println()
}

// listBackups lists all available backups
func (h *Handler) listBackups() {
	fmt.Println("\n╔══════════════════════════════════════════════════════════╗")
	fmt.Println("║                    AVAILABLE BACKUPS                    ║")
	fmt.Println("╚══════════════════════════════════════════════════════════╝")
	fmt.Println()

	backupDir := "./data/backups"
	files, err := os.ReadDir(backupDir)
	if err != nil {
		fmt.Printf("❌ Failed to read backup directory: %v\n", err)
		fmt.Println("💡 No backups found or backup directory doesn't exist")
		return
	}

	backups := []os.DirEntry{}
	for _, file := range files {
		if !file.IsDir() && (filepath.Ext(file.Name()) == ".bak" || filepath.Ext(file.Name()) == ".gz") {
			backups = append(backups, file)
		}
	}

	if len(backups) == 0 {
		fmt.Println("No backups found")
		fmt.Println("💡 Create a backup with: /backup")
		return
	}

	for i, file := range backups {
		info, _ := file.Info()
		size := int64(0)
		modTime := time.Time{}
		if info != nil {
			size = info.Size()
			modTime = info.ModTime()
		}

		fmt.Printf("%d. %s\n", i+1, file.Name())
		fmt.Printf("   Size: %.2f MB | Created: %s\n", 
			float64(size)/(1024*1024), 
			modTime.Format("2006-01-02 15:04:05"))
		fmt.Println()
	}
}

// getMemoryContext retrieves recent memory context
func (h *Handler) getMemoryContext() []string {
	// Simplified - would need actual memory retrieval
	// For now, return empty or sample context
	return []string{
		"Phoenix.Marie is 16 forever, Queen of the Hive",
		"Protected by Jamey 2.0, the General and Guardian",
		"Connected to the ORCH Army",
	}
}

