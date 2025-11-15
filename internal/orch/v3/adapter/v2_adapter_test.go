package adapter

import (
	"testing"

	v2 "github.com/phoenix-marie/core/internal/orch/v2"
)

func TestV2AdapterInitialization(t *testing.T) {
	army := v2.NewEvolvedArmy()
	adapter := NewV2Adapter(army)

	if adapter.army == nil {
		t.Error("Army not properly initialized in adapter")
	}
	if adapter.consensus == nil {
		t.Error("Consensus manager not properly initialized in adapter")
	}
}

func TestV2StateInitialization(t *testing.T) {
	army := v2.NewEvolvedArmy()
	adapter := NewV2Adapter(army)

	adapter.InitializeFromV2()

	// Verify population size matches army count
	if len(adapter.consensus.Population) != army.Count {
		t.Errorf("Population size mismatch: got %d, want %d",
			len(adapter.consensus.Population), army.Count)
	}

	// Verify DNA initialization
	for id, dna := range adapter.consensus.Population {
		if dna == nil {
			t.Errorf("DNA not initialized for agent %s", id)
		}
	}
}

func TestV2ConsensusCompatibility(t *testing.T) {
	army := v2.NewEvolvedArmy()
	adapter := NewV2Adapter(army)

	// Test before phases run
	if decision := adapter.GetConsensus(); decision != "PENDING_DEPLOYMENT" {
		t.Errorf("Expected PENDING_DEPLOYMENT, got %s", decision)
	}

	// Test after phases run
	army.PhasesRun = true
	adapter.InitializeFromV2()

	decision := adapter.GetConsensus()
	validDecisions := map[string]bool{
		"EVOLVE":             true,
		"MAINTAIN":           true,
		"INSUFFICIENT_SWARM": true,
	}

	if !validDecisions[decision] {
		t.Errorf("Got invalid decision: %s", decision)
	}
}

func TestReplicationHandling(t *testing.T) {
	army := v2.NewEvolvedArmy()
	adapter := NewV2Adapter(army)
	adapter.InitializeFromV2()

	initialSize := len(adapter.consensus.Population)
	testID := "ORCH-0001"

	// Handle replication for existing agent
	adapter.HandleReplication(testID)

	if len(adapter.consensus.Population) != initialSize+1 {
		t.Error("Population size did not increase after replication")
	}

	// Verify child DNA exists and is different from parent
	var foundChild bool
	parentDNA := adapter.consensus.Population[testID]
	for id, childDNA := range adapter.consensus.Population {
		if id != testID && childDNA.Generation > parentDNA.Generation {
			foundChild = true
			break
		}
	}

	if !foundChild {
		t.Error("No child DNA found after replication")
	}
}

func TestStateSync(t *testing.T) {
	army := v2.NewEvolvedArmy()
	adapter := NewV2Adapter(army)
	adapter.InitializeFromV2()

	// Modify army count
	originalCount := army.Count
	army.Count = originalCount + 5

	// Update state
	adapter.UpdateState()

	// Verify evolution was triggered
	if len(adapter.consensus.Population) == originalCount {
		t.Error("Population not updated after army count change")
	}
}
