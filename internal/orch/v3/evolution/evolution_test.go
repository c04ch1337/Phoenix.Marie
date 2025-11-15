package evolution

import (
	"testing"

	"github.com/phoenix-marie/core/internal/orch/v3/dna"
)

func TestConsensusManager(t *testing.T) {
	// Initialize manager with test parameters
	cm := NewConsensusManager(3, 10)
	if cm == nil {
		t.Fatal("Failed to create ConsensusManager")
	}

	// Test initial state
	if len(cm.Population) != 0 {
		t.Errorf("Expected empty population, got size %d", len(cm.Population))
	}

	consensus, err := cm.GetConsensus()
	if err != nil {
		t.Errorf("Unexpected error getting consensus: %v", err)
	}
	if consensus != "INSUFFICIENT_POPULATION" {
		t.Errorf("Expected INSUFFICIENT_POPULATION, got %s", consensus)
	}
}

func TestPopulationManagement(t *testing.T) {
	cm := NewConsensusManager(3, 5)

	// Add members up to max population
	for i := 1; i <= 7; i++ {
		d := dna.NewDNA(string(rune('A' + i)))
		cm.AddMember(d)
	}

	// Check population control
	if len(cm.Population) > 5 {
		t.Errorf("Population control failed: size %d exceeds max 5", len(cm.Population))
	}

	// Test member removal
	cm.RemoveMember("B")
	if _, exists := cm.Population["B"]; exists {
		t.Error("Failed to remove member B from population")
	}
}

func TestEvolutionProcess(t *testing.T) {
	cm := NewConsensusManager(3, 5)

	// Add initial population with different fitness levels
	members := []struct {
		id    string
		rates map[string]float64
	}{
		{"A", map[string]float64{"replication_rate": 0.8, "consensus_weight": 0.9, "adaptation_speed": 0.7}},
		{"B", map[string]float64{"replication_rate": 0.4, "consensus_weight": 0.5, "adaptation_speed": 0.3}},
		{"C", map[string]float64{"replication_rate": 0.6, "consensus_weight": 0.7, "adaptation_speed": 0.5}},
	}

	for _, m := range members {
		d := dna.NewDNA(m.id)
		for gene, value := range m.rates {
			d.Genes[gene].Value = value
		}
		cm.AddMember(d)
	}

	// Test consensus before evolution
	consensus, err := cm.GetConsensus()
	if err != nil {
		t.Errorf("Unexpected error getting consensus: %v", err)
	}
	if consensus != "REPLICATE" {
		t.Errorf("Expected REPLICATE consensus due to high replication rate, got %s", consensus)
	}

	// Test evolution process
	cm.Evolve()

	// Verify population size is maintained within bounds
	if len(cm.Population) < cm.minPopulation || len(cm.Population) > cm.maxPopulation {
		t.Errorf("Population size %d outside bounds [%d, %d]",
			len(cm.Population), cm.minPopulation, cm.maxPopulation)
	}

	// Verify new generation contains evolved members
	newGenFound := false
	for _, member := range cm.Population {
		if member.Generation > 1 {
			newGenFound = true
			break
		}
	}
	if !newGenFound {
		t.Error("No evolved members found after evolution")
	}
}

func TestConsensusVoting(t *testing.T) {
	cm := NewConsensusManager(3, 5)

	// Add members with specific traits to test voting
	testCases := []struct {
		id       string
		traits   map[string]float64
		expected string
	}{
		{
			"HIGH_REP",
			map[string]float64{
				"replication_rate": 0.9,
				"consensus_weight": 1.0,
				"adaptation_speed": 0.3,
			},
			"REPLICATE",
		},
		{
			"HIGH_ADAPT",
			map[string]float64{
				"replication_rate": 0.3,
				"consensus_weight": 1.0,
				"adaptation_speed": 0.9,
			},
			"EVOLVE",
		},
		{
			"BALANCED",
			map[string]float64{
				"replication_rate": 0.5,
				"consensus_weight": 1.0,
				"adaptation_speed": 0.5,
			},
			"MAINTAIN",
		},
	}

	// Add test members
	for _, tc := range testCases {
		d := dna.NewDNA(tc.id)
		for gene, value := range tc.traits {
			d.Genes[gene].Value = value
		}
		cm.AddMember(d)
	}

	// Test consensus
	consensus, err := cm.GetConsensus()
	if err != nil {
		t.Errorf("Unexpected error getting consensus: %v", err)
	}

	// Verify consensus reflects weighted voting
	if consensus == "INSUFFICIENT_POPULATION" {
		t.Error("Got INSUFFICIENT_POPULATION despite having enough members")
	}
}
