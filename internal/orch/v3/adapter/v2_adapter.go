package adapter

import (
	"fmt"

	v2 "github.com/phoenix-marie/core/internal/orch/v2"
	"github.com/phoenix-marie/core/internal/orch/v3/dna"
	"github.com/phoenix-marie/core/internal/orch/v3/evolution"
)

// V2Adapter provides compatibility between v2 and v3 systems
type V2Adapter struct {
	army      *v2.EvolvedArmy
	consensus *evolution.ConsensusManager
}

// NewV2Adapter creates a new adapter instance
func NewV2Adapter(army *v2.EvolvedArmy) *V2Adapter {
	// Initialize consensus manager with army parameters
	consensus := evolution.NewConsensusManager(5, army.Count)

	return &V2Adapter{
		army:      army,
		consensus: consensus,
	}
}

// InitializeFromV2 creates initial DNA population from v2 army state
func (a *V2Adapter) InitializeFromV2() {
	// Create DNA for each existing agent
	for i := 1; i <= a.army.Count; i++ {
		id := fmt.Sprintf("ORCH-%04d", i)
		d := dna.NewDNA(id)

		// Set initial DNA values based on v2 state
		if a.army.PhasesRun {
			// If phases are run, set higher adaptation rate
			d.Genes["adaptation_speed"].Value = 0.7
		}

		a.consensus.AddMember(d)
	}
}

// GetConsensus returns a consensus decision compatible with v2
func (a *V2Adapter) GetConsensus() string {
	if !a.army.PhasesRun {
		return "PENDING_DEPLOYMENT"
	}

	decision, err := a.consensus.GetConsensus()
	if err != nil {
		return "INSUFFICIENT_SWARM"
	}

	// Map v3 decisions to v2 format
	switch decision {
	case "REPLICATE", "EVOLVE":
		return "EVOLVE"
	case "MAINTAIN":
		return "MAINTAIN"
	default:
		return "INSUFFICIENT_SWARM"
	}
}

// UpdateState synchronizes v2 and v3 states
func (a *V2Adapter) UpdateState() {
	// Trigger evolution if army size changed
	if len(a.consensus.Population) != a.army.Count {
		a.consensus.Evolve()
	}
}

// HandleReplication manages replication events from v2
func (a *V2Adapter) HandleReplication(agentID string) {
	if d := a.consensus.Population[agentID]; d != nil {
		// Create child DNA through mutation
		child := dna.NewDNA(fmt.Sprintf("%s-CHILD", agentID))
		for name, gene := range d.Genes {
			child.Genes[name].Value = gene.Value
		}
		child.Mutate()

		a.consensus.AddMember(child)
	}
}
