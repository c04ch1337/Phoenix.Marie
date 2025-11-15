package v2

import (
	"log"
	"os"
	"strconv"
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
