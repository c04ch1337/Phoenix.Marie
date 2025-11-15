package dna

import (
	"crypto/rand"
	"math"
	"math/big"
)

// Gene represents a single trait in the DNA
type Gene struct {
	Name       string
	Value      float64
	MutateProb float64
}

// DNA represents the genetic makeup of an ORCH agent
type DNA struct {
	ID         string
	Genes      map[string]*Gene
	Fitness    float64
	Generation int
}

// NewDNA creates a new DNA instance with default genes
func NewDNA(id string) *DNA {
	dna := &DNA{
		ID:         id,
		Genes:      make(map[string]*Gene),
		Generation: 1,
	}

	// Initialize default genes
	dna.Genes["replication_rate"] = &Gene{
		Name:       "replication_rate",
		Value:      0.5,
		MutateProb: 0.1,
	}
	dna.Genes["consensus_weight"] = &Gene{
		Name:       "consensus_weight",
		Value:      1.0,
		MutateProb: 0.05,
	}
	dna.Genes["adaptation_speed"] = &Gene{
		Name:       "adaptation_speed",
		Value:      0.3,
		MutateProb: 0.15,
	}

	return dna
}

// Mutate applies random mutations to genes based on their mutation probabilities
func (d *DNA) Mutate() {
	for _, gene := range d.Genes {
		if shouldMutate(gene.MutateProb) {
			// Apply gaussian mutation
			mutation := gaussianMutation(0.1) // 0.1 is the standard deviation
			gene.Value = math.Max(0, math.Min(1, gene.Value+mutation))
		}
	}
}

// Crossover creates a new DNA by combining genes from two parents
func Crossover(dna1, dna2 *DNA) *DNA {
	childID := generateChildID(dna1.ID, dna2.ID)
	child := NewDNA(childID)
	child.Generation = max(dna1.Generation, dna2.Generation) + 1

	for name, gene1 := range dna1.Genes {
		gene2 := dna2.Genes[name]
		if gene2 != nil {
			// Interpolate between parent values with random weight
			weight := randomFloat()
			child.Genes[name].Value = gene1.Value*weight + gene2.Value*(1-weight)
			// Average mutation probabilities
			child.Genes[name].MutateProb = (gene1.MutateProb + gene2.MutateProb) / 2
		}
	}

	return child
}

// CalculateFitness computes the fitness score based on gene values
func (d *DNA) CalculateFitness() float64 {
	// Basic fitness calculation - can be expanded based on specific requirements
	fitness := 0.0
	for _, gene := range d.Genes {
		fitness += gene.Value
	}
	d.Fitness = fitness / float64(len(d.Genes))
	return d.Fitness
}

// Helper functions

func shouldMutate(probability float64) bool {
	n, _ := rand.Int(rand.Reader, big.NewInt(100))
	return float64(n.Int64())/100 < probability
}

func gaussianMutation(stdDev float64) float64 {
	// Box-Muller transform for gaussian distribution
	u1 := randomFloat()
	u2 := randomFloat()
	z0 := math.Sqrt(-2*math.Log(u1)) * math.Cos(2*math.Pi*u2)
	return z0 * stdDev
}

func randomFloat() float64 {
	n, _ := rand.Int(rand.Reader, big.NewInt(1000))
	return float64(n.Int64()) / 1000
}

func generateChildID(parent1ID, parent2ID string) string {
	// Simple concatenation of first parts of parent IDs
	return parent1ID[:4] + "-" + parent2ID[4:]
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}
