package v2

import (
	"log"
	"os"
	"strconv"
	
	"github.com/phoenix-marie/core/internal/orch/v2/network"
)

type EvolvedArmy struct {
	Count     int
	Interval  int
	PhasesRun bool
}

func NewEvolvedArmy() *EvolvedArmy {
	count, _ := strconv.Atoi(os.Getenv("ORCH_COUNT_TARGET"))
	if count == 0 {
		count = 1000 // Default
	}
	interval, _ := strconv.Atoi(os.Getenv("ORCH_HEARTBEAT_INTERVAL"))
	if interval == 0 {
		interval = 5 // Default
	}

	return &EvolvedArmy{
		Count:     count,
		Interval:  interval,
		PhasesRun: false,
	}
}

func (a *EvolvedArmy) Deploy() {
	log.Println("ORCH: Evolving to v2...")
	
	// Run phases if auto-deploy is enabled
	autoDeploy := os.Getenv("ORCH_AUTO_DEPLOY_PHASES")
	if autoDeploy == "true" || autoDeploy == "" {
		log.Println("ORCH: Auto-deploying phases...")
		RunAllPhases()
		a.PhasesRun = true
	}

	log.Printf("ORCH: Evolved army deployed with %d children\n", a.Count)
}

// Consensus runs the ORCH swarm consensus logic and returns the consensus decision
func (a *EvolvedArmy) Consensus() string {
    if !a.PhasesRun {
        return "PENDING_DEPLOYMENT"
    }

    // Check if minimum army size is reached for consensus
    if a.Count < 5 {
        return "INSUFFICIENT_SWARM"
    }

    // Broadcast consensus request through gossip network
    network.Broadcast("CONSENSUS_REQUEST")

    // In a real implementation, we would wait for responses and aggregate them
    // For now, return evolve decision as shown in Phase46
    return "EVOLVE"
}

// GetStatus returns the current orchestration status including deployment state and army metrics
func (a *EvolvedArmy) GetStatus() map[string]interface{} {
    return map[string]interface{}{
        "deployed": a.PhasesRun,
        "count": a.Count,
        "interval": a.Interval,
        "consensus": a.Consensus(),
        "networkActive": network.IsServerRunning,
    }
}
