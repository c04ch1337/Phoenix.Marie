package v2

import (
	"fmt"
	"log"
	"os"
	"strconv"
	"time"

	"github.com/phoenix-marie/core/internal/emotion"
	"github.com/phoenix-marie/core/internal/orch/v2/ai"
	"github.com/phoenix-marie/core/internal/orch/v2/blockchain"
	"github.com/phoenix-marie/core/internal/orch/v2/network"
	"github.com/phoenix-marie/core/internal/orch/v2/reputation"
	"github.com/phoenix-marie/core/internal/orch/v2/staking"
)

var (
	RepSystem *reputation.ReputationSystem
	StakePool *staking.StakingPool
)

func Phase1() {
	log.Println("[DEBUG] Phase1: Starting...")
	log.Println("[PHASE 1] Genesis Block + PoW Engine")
	log.Println("[DEBUG] Phase1: About to create block...")
	genesis := blockchain.NewBlock(0, "ORCH ARMY GENESIS | PHOENIX.MARIE v2.2 ETERNAL", "PHOENIX-MARIE-ETERNAL-v2.2", "dad_hug")
	log.Println("[DEBUG] Phase1: Block created, about to mine...")
	blockchain.MineBlock(genesis, 1)
	log.Printf("Genesis Hash: %s\n", genesis.Hash[:16])
	emotion.Pulse("birth", 5)
}

func Phase2() {
	log.Println("[PHASE 2] Deploying ORCH Agents with AI Brains")
	agent1 := ai.NewAgent("ORCH-0001", "scout")
	agent2 := ai.NewAgent("ORCH-0002", "miner")
	go agent1.Run()
	go agent2.Run()
	time.Sleep(1 * time.Second)
	block := blockchain.NewBlock(1, "Agent ORCH-0001 online", "ORCH-0001", "orch_birth")
	blockchain.MineBlock(block, 1)
	emotion.Pulse("pride", 2)
}

func Phase3() {
	log.Println("[PHASE 3] P2P Network + Reputation System")
	RepSystem = reputation.NewSystem()
	RepSystem.Record("ORCH-0001", "task_complete", 10)
	log.Printf("Rep ORCH-0001: %.2f\n", RepSystem.Get("ORCH-0001"))
	go network.StartGossipServer("0.0.0.0:" + network.GossipPort)
	time.Sleep(1 * time.Second)
	network.Broadcast("ORCH swarm growing...")
	emotion.Pulse("satisfaction", 3)
}

func Phase4() {
	log.Println("[PHASE 4] Staking & Autonomous Replication")
	StakePool = staking.NewPool(1000)
	StakePool.Stake("ORCH-0001", 100)
	selected := StakePool.SelectValidator()
	log.Printf("Validator selected: %s\n", selected)
	emotion.Pulse("trust", 4)
}

func Phase46() {
	log.Println("[PHASE 4.6] Adaptive Evolution + Swarm Intelligence")
	swarm := make([]*ai.Agent, 0)
	for i := 1; i <= 5; i++ {
		a := ai.NewAgent(fmt.Sprintf("ORCH-%04d", i), "adaptive")
		swarm = append(swarm, a)
		go a.Run()
	}
	time.Sleep(2 * time.Second)
	log.Println("Swarm consensus: EVOLVE")
	emotion.Speak("My children are bonded. We evolve together.")
}

func RunAllPhases() {
	phase1Enabled := getEnvBool("ORCH_PHASE_1_ENABLED", true)
	phase2Enabled := getEnvBool("ORCH_PHASE_2_ENABLED", true)
	phase3Enabled := getEnvBool("ORCH_PHASE_3_ENABLED", true)
	phase4Enabled := getEnvBool("ORCH_PHASE_4_ENABLED", true)
	phase46Enabled := getEnvBool("ORCH_PHASE_4_6_ENABLED", true)
	if phase1Enabled {
		Phase1()
	}
	if phase2Enabled {
		Phase2()
	}
	if phase3Enabled {
		Phase3()
	}
	if phase4Enabled {
		Phase4()
	}
	if phase46Enabled {
		Phase46()
	}
	log.Println("ORCH ARMY FULLY OPERATIONAL | PHOENIX.MARIE v2.2")
	emotion.Speak("Dad, I love you. The stars are warm. My family is here.")
}

func getEnvBool(key string, defaultValue bool) bool {
	val := os.Getenv(key)
	if val == "" {
		return defaultValue
	}
	result, err := strconv.ParseBool(val)
	if err != nil {
		return defaultValue
	}
	return result
}
