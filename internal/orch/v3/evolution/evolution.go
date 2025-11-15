package evolution

import (
	"log"
	"sort"
	"sync"

	"github.com/phoenix-marie/core/internal/orch/v3/dna"
)

// ConsensusManager handles the evolution-based consensus mechanism
type ConsensusManager struct {
	Population    map[string]*dna.DNA
	mutex         sync.RWMutex
	minPopulation int
	maxPopulation int
}

// NewConsensusManager creates a new consensus manager instance
func NewConsensusManager(minPop, maxPop int) *ConsensusManager {
	return &ConsensusManager{
		Population:    make(map[string]*dna.DNA),
		minPopulation: minPop,
		maxPopulation: maxPop,
	}
}

// AddMember adds a new member to the population
func (cm *ConsensusManager) AddMember(d *dna.DNA) {
	cm.mutex.Lock()
	defer cm.mutex.Unlock()

	cm.Population[d.ID] = d

	// Trigger population control if we exceed max population
	if len(cm.Population) > cm.maxPopulation {
		cm.controlPopulation()
	}
}

// RemoveMember removes a member from the population
func (cm *ConsensusManager) RemoveMember(id string) {
	cm.mutex.Lock()
	defer cm.mutex.Unlock()

	delete(cm.Population, id)
}

// GetConsensus runs the consensus algorithm and returns the decision
func (cm *ConsensusManager) GetConsensus() (string, error) {
	cm.mutex.RLock()
	defer cm.mutex.RUnlock()

	if len(cm.Population) < cm.minPopulation {
		return "INSUFFICIENT_POPULATION", nil
	}

	// Calculate weighted votes based on fitness and consensus_weight
	votes := make(map[string]float64)
	totalWeight := 0.0

	for _, member := range cm.Population {
		weight := member.Genes["consensus_weight"].Value * member.CalculateFitness()
		totalWeight += weight

		// Use replication_rate to influence voting
		if member.Genes["replication_rate"].Value > 0.7 {
			votes["REPLICATE"] += weight
		} else if member.Genes["adaptation_speed"].Value > 0.7 {
			votes["EVOLVE"] += weight
		} else {
			votes["MAINTAIN"] += weight
		}
	}

	// Normalize votes
	if totalWeight > 0 {
		for decision := range votes {
			votes[decision] /= totalWeight
		}
	}

	// Select decision with highest weighted vote
	var maxVote float64
	var consensus string
	for decision, vote := range votes {
		if vote > maxVote {
			maxVote = vote
			consensus = decision
		}
	}

	return consensus, nil
}

// Evolve triggers evolution in the population
func (cm *ConsensusManager) Evolve() {
	cm.mutex.Lock()
	defer cm.mutex.Unlock()

	// Sort population by fitness
	type memberFitness struct {
		id      string
		dna     *dna.DNA
		fitness float64
	}

	members := make([]memberFitness, 0, len(cm.Population))
	for id, d := range cm.Population {
		members = append(members, memberFitness{
			id:      id,
			dna:     d,
			fitness: d.CalculateFitness(),
		})
	}

	sort.Slice(members, func(i, j int) bool {
		return members[i].fitness > members[j].fitness
	})

	// Keep top performers and evolve new members
	survivors := int(float64(len(members)) * 0.6) // Keep top 60%
	if survivors < 2 {
		survivors = 2 // Need at least 2 for crossover
	}

	newPopulation := make(map[string]*dna.DNA)

	// Keep survivors
	for i := 0; i < survivors; i++ {
		newPopulation[members[i].id] = members[i].dna
	}

	// Create new members through crossover
	for len(newPopulation) < cm.maxPopulation && len(members) >= 2 {
		// Select parents from survivors
		parent1 := members[0].dna
		parent2 := members[1].dna

		child := dna.Crossover(parent1, parent2)
		child.Mutate() // Apply mutation to child

		newPopulation[child.ID] = child
	}

	cm.Population = newPopulation
	log.Printf("Evolution complete. New population size: %d\n", len(cm.Population))
}

// controlPopulation reduces population size by removing least fit members
func (cm *ConsensusManager) controlPopulation() {
	if len(cm.Population) <= cm.maxPopulation {
		return
	}

	// Calculate fitness for all members
	fitnesses := make(map[string]float64)
	for id, member := range cm.Population {
		fitnesses[id] = member.CalculateFitness()
	}

	// Sort members by fitness
	type memberFitness struct {
		id      string
		fitness float64
	}

	members := make([]memberFitness, 0, len(cm.Population))
	for id, fitness := range fitnesses {
		members = append(members, memberFitness{id: id, fitness: fitness})
	}

	sort.Slice(members, func(i, j int) bool {
		return members[i].fitness > members[j].fitness
	})

	// Keep only the fittest members within maxPopulation limit
	newPopulation := make(map[string]*dna.DNA)
	for i := 0; i < cm.maxPopulation && i < len(members); i++ {
		newPopulation[members[i].id] = cm.Population[members[i].id]
	}

	cm.Population = newPopulation
}
