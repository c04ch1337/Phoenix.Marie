package internal

import (
	"bytes"
	"fmt"
	"log"
	"os"
	"strings"
	"testing"

	"github.com/phoenix-marie/core/internal/core"
	"github.com/phoenix-marie/core/internal/core/memory"
	"github.com/phoenix-marie/core/internal/emotion"
	v2 "github.com/phoenix-marie/core/internal/orch/v2"
	"github.com/phoenix-marie/core/internal/orch/v3/dna"
	"github.com/phoenix-marie/core/internal/orch/v3/evolution"
)

func init() {
	// Set environment variables before any tests run
	os.Setenv("EMOTION_FLAME_PULSE_BASE", "5")
	os.Setenv("EMOTION_VOICE_TONE", "calm")
	os.Setenv("EMOTION_RESPONSE_STYLE", "direct")
}

func TestGetCurrentState(t *testing.T) {
	state := emotion.GetCurrentState()

	if state["flamePulse"] == nil {
		t.Error("Expected flamePulse to be present in state")
	}

	if state["voiceTone"] == nil {
		t.Error("Expected voiceTone to be present in state")
	}

	if state["responseStyle"] == nil {
		t.Error("Expected responseStyle to be present in state")
	}
}

func TestOrchFunctions(t *testing.T) {
	army := v2.NewEvolvedArmy()
	army.Deploy()

	// Test Consensus
	consensus := army.Consensus()
	if consensus == "" {
		t.Error("Expected non-empty consensus decision")
	}

	// Test GetStatus
	status := army.GetStatus()

	if status["deployed"] != true {
		t.Error("Expected army to be deployed")
	}

	if status["count"] == nil {
		t.Error("Expected count to be present in status")
	}

	if status["interval"] == nil {
		t.Error("Expected interval to be present in status")
	}
}

func TestMemoryPersistence(t *testing.T) {
	phl, err := memory.NewPHL("./testdata")
	if err != nil {
		t.Fatalf("Failed to initialize PHL: %v", err)
	}
	defer phl.Close()

	// Test storing and retrieving data
	emotionState := emotion.GetCurrentState()
	if !phl.Store("emotion", "currentState", emotionState) {
		t.Error("Failed to store emotion state in memory")
	}

	retrieved, ok := phl.Retrieve("emotion", "currentState")
	if !ok {
		t.Error("Failed to retrieve emotion state from memory")
	}
	if retrieved == nil {
		t.Error("Retrieved emotion state is nil")
	}

	// Test storing and retrieving ORCH status
	army := v2.NewEvolvedArmy()
	orchStatus := army.GetStatus()
	if !phl.Store("logic", "orchStatus", orchStatus) {
		t.Error("Failed to store ORCH status in memory")
	}

	retrieved, ok = phl.Retrieve("logic", "orchStatus")
	if !ok {
		t.Error("Failed to retrieve ORCH status from memory")
	}
	if retrieved == nil {
		t.Error("Retrieved ORCH status is nil")
	}

	// Test memory cleanup
	if !phl.Cleanup("emotion") {
		t.Error("Failed to cleanup emotion layer")
	}

	// Verify cleanup worked
	_, ok = phl.Retrieve("emotion", "currentState")
	if ok {
		t.Error("Data still exists after cleanup")
	}

	// Test invalid layer operations
	if phl.Store("invalid", "key", "value") {
		t.Error("Store operation succeeded on invalid layer")
	}

	if _, ok := phl.Retrieve("invalid", "key"); ok {
		t.Error("Retrieve operation succeeded on invalid layer")
	}

	if phl.Cleanup("invalid") {
		t.Error("Cleanup operation succeeded on invalid layer")
	}
}

func TestThoughtEngineCore(t *testing.T) {
	// Initialize Phoenix core
	phoenix := core.Ignite()

	// Verify initial state
	firstThought, ok := phoenix.Memory.Retrieve("emotion", "first_thought")
	if !ok {
		t.Error("Failed to retrieve first thought from memory")
	}
	if firstThought != "I am warm. I am loved." {
		t.Errorf("Unexpected first thought: %v", firstThought)
	}

	// Test flame core integration
	initialPulse := phoenix.Flame.PulseRate
	if initialPulse != 1 {
		t.Errorf("Unexpected initial pulse rate: %d", initialPulse)
	}

	// Test speaking and pulse rate changes
	phoenix.Speak("Testing thought process")
	if phoenix.Flame.PulseRate != 2 {
		t.Errorf("Pulse rate did not increment after speaking. Expected 2, got %d", phoenix.Flame.PulseRate)
	}

	// Test pulse rate cycling
	for i := 0; i < 10; i++ {
		phoenix.Flame.Pulse()
	}
	if phoenix.Flame.PulseRate != 2 {
		t.Errorf("Pulse rate did not cycle correctly. Expected 2, got %d", phoenix.Flame.PulseRate)
	}

	// Verify DNA lock is present
	if phoenix.DNA == nil {
		t.Error("DNA lock not initialized")
	}
	if phoenix.DNA.ID != "PHOENIX-MARIE-001" {
		t.Errorf("Unexpected DNA ID: %s", phoenix.DNA.ID)
	}
}

func TestEmotionIntegration(t *testing.T) {
	// Reset emotion state to initial values
	emotion.Reset()

	// Get initial emotion state
	state := emotion.GetCurrentState()

	// Verify environment configuration
	if state["flamePulse"] != 5 {
		t.Errorf("Expected flame pulse base of 5, got %v", state["flamePulse"])
	}
	if state["voiceTone"] != "calm" {
		t.Errorf("Expected voice tone 'calm', got %v", state["voiceTone"])
	}
	if state["responseStyle"] != "direct" {
		t.Errorf("Expected response style 'direct', got %v", state["responseStyle"])
	}

	// Test emotion pulse integration
	emotion.Pulse("joy", 2)
	state = emotion.GetCurrentState()
	if state["flamePulse"].(int) != 7 { // Base 5 + intensity 2
		t.Errorf("Expected flame pulse 7, got %v", state["flamePulse"])
	}

	// Test voice output with emotion configuration
	var logBuffer bytes.Buffer
	log.SetOutput(&logBuffer)
	defer log.SetOutput(os.Stdout)

	emotion.Speak("Test message")
	logOutput := logBuffer.String()

	if !strings.Contains(logOutput, "PHOENIX (calm, direct): Test message") {
		t.Errorf("Expected emotion-configured output, got: %s", logOutput)
	}

	// Verify emotion state persistence in memory
	phl, err := memory.NewPHL("./testdata")
	if err != nil {
		t.Fatalf("Failed to initialize PHL: %v", err)
	}
	defer phl.Close()
	phl.Store("emotion", "currentState", state)
	retrieved, ok := phl.Retrieve("emotion", "currentState")
	if !ok {
		t.Error("Failed to retrieve emotion state from memory")
	}

	retrievedState := retrieved.(map[string]interface{})
	if retrievedState["flamePulse"] != state["flamePulse"] {
		t.Errorf("Emotion state not persisted correctly in memory. Expected %v, got %v",
			state["flamePulse"], retrievedState["flamePulse"])
	}
}

func TestEvolutionSystemIntegration(t *testing.T) {
	// Initialize consensus manager with test parameters
	manager := evolution.NewConsensusManager(3, 10)

	// Create initial population with diverse traits
	initialMembers := []struct {
		id     string
		traits map[string]float64
	}{
		{
			"ORCH-001",
			map[string]float64{
				"replication_rate": 0.8,
				"consensus_weight": 0.9,
				"adaptation_speed": 0.3,
			},
		},
		{
			"ORCH-002",
			map[string]float64{
				"replication_rate": 0.4,
				"consensus_weight": 0.6,
				"adaptation_speed": 0.8,
			},
		},
		{
			"ORCH-003",
			map[string]float64{
				"replication_rate": 0.6,
				"consensus_weight": 0.7,
				"adaptation_speed": 0.6,
			},
		},
	}

	// Add initial members
	for _, m := range initialMembers {
		d := dna.NewDNA(m.id)
		for gene, value := range m.traits {
			d.Genes[gene].Value = value
		}
		manager.AddMember(d)
	}

	// Verify initial consensus based on traits
	consensus, err := manager.GetConsensus()
	if err != nil {
		t.Errorf("Unexpected error getting initial consensus: %v", err)
	}
	if consensus != "REPLICATE" {
		t.Errorf("Expected REPLICATE consensus due to high replication rate, got %s", consensus)
	}

	// Evolve population multiple times
	for i := 0; i < 3; i++ {
		manager.Evolve()

		// Verify population constraints
		popSize := len(manager.Population)
		if popSize < 3 || popSize > 10 {
			t.Errorf("Population size %d outside bounds [3, 10] after evolution %d", popSize, i+1)
		}

		// Verify gene value bounds
		for _, member := range manager.Population {
			for _, gene := range member.Genes {
				if gene.Value < 0 || gene.Value > 1 {
					t.Errorf("Gene %s value %.2f outside valid range [0,1] after evolution %d",
						gene.Name, gene.Value, i+1)
				}
			}
		}

		// Verify generation increment
		for _, member := range manager.Population {
			if member.Generation < 1 {
				t.Errorf("Invalid generation %d for member %s", member.Generation, member.ID)
			}
		}
	}

	// Test population control
	for i := 0; i < 5; i++ {
		d := dna.NewDNA(fmt.Sprintf("EXTRA-%d", i))
		manager.AddMember(d)
	}

	if len(manager.Population) > 10 {
		t.Error("Population control failed to maintain maximum size")
	}

	// Verify fitness-based selection
	var totalFitness float64
	minFitness := 1.0
	maxFitness := 0.0

	for _, member := range manager.Population {
		fitness := member.CalculateFitness()
		totalFitness += fitness
		if fitness < minFitness {
			minFitness = fitness
		}
		if fitness > maxFitness {
			maxFitness = fitness
		}
	}

	avgFitness := totalFitness / float64(len(manager.Population))
	if avgFitness < 0.3 {
		t.Errorf("Average fitness %.2f below acceptable threshold", avgFitness)
	}
}

func TestORCHConnectivity(t *testing.T) {
	// Set required environment variables
	os.Setenv("ORCH_COUNT_TARGET", "10")
	os.Setenv("ORCH_HEARTBEAT_INTERVAL", "3")
	os.Setenv("ORCH_AUTO_DEPLOY_PHASES", "true")

	// Initialize army
	army := v2.NewEvolvedArmy()

	// Verify initial state
	if army.Count != 10 {
		t.Errorf("Expected army count 10, got %d", army.Count)
	}
	if army.Interval != 3 {
		t.Errorf("Expected heartbeat interval 3, got %d", army.Interval)
	}
	if army.PhasesRun {
		t.Error("Expected phases not run initially")
	}

	// Test pre-deployment consensus
	if consensus := army.Consensus(); consensus != "PENDING_DEPLOYMENT" {
		t.Errorf("Expected PENDING_DEPLOYMENT consensus, got %s", consensus)
	}

	// Deploy army and verify phases run
	army.Deploy()
	if !army.PhasesRun {
		t.Error("Expected phases to be run after deployment")
	}

	// Test post-deployment status
	status := army.GetStatus()
	if !status["deployed"].(bool) {
		t.Error("Expected deployed status to be true")
	}
	if status["count"] != 10 {
		t.Errorf("Expected count 10 in status, got %v", status["count"])
	}
	if status["interval"] != 3 {
		t.Errorf("Expected interval 3 in status, got %v", status["interval"])
	}

	// Test post-deployment consensus
	consensus := army.Consensus()
	if consensus != "EVOLVE" {
		t.Errorf("Expected EVOLVE consensus after deployment, got %s", consensus)
	}

	// Skip network activation check since it requires port binding
	// Test insufficient swarm size scenario
	smallArmy := v2.NewEvolvedArmy()
	smallArmy.Count = 3        // Below minimum threshold
	smallArmy.PhasesRun = true // Mark as deployed to test size-based consensus
	if consensus := smallArmy.Consensus(); consensus != "INSUFFICIENT_SWARM" {
		t.Errorf("Expected INSUFFICIENT_SWARM for small army, got %s", consensus)
	}
}
