package dna

import (
	"strings"
	"testing"
)

func TestNewDNA(t *testing.T) {
	dna := NewDNA("TEST-001")

	if dna.ID != "TEST-001" {
		t.Errorf("Expected ID TEST-001, got %s", dna.ID)
	}

	expectedGenes := []string{"replication_rate", "consensus_weight", "adaptation_speed"}
	for _, geneName := range expectedGenes {
		if gene, exists := dna.Genes[geneName]; !exists {
			t.Errorf("Expected gene %s not found", geneName)
		} else if gene.Value < 0 || gene.Value > 1 {
			t.Errorf("Gene %s value %.2f outside valid range [0,1]", geneName, gene.Value)
		}
	}
}

func TestDNAMutation(t *testing.T) {
	dna := NewDNA("TEST-002")
	originalValues := make(map[string]float64)

	// Store original values
	for name, gene := range dna.Genes {
		originalValues[name] = gene.Value
	}

	// Perform multiple mutations to ensure at least some genes change
	for i := 0; i < 100; i++ {
		dna.Mutate()
	}

	changed := false
	for name, gene := range dna.Genes {
		if gene.Value != originalValues[name] {
			changed = true
			if gene.Value < 0 || gene.Value > 1 {
				t.Errorf("Mutation produced invalid value %.2f for gene %s", gene.Value, name)
			}
		}
	}

	if !changed {
		t.Error("No genes were modified after 100 mutations")
	}
}

func TestDNACrossover(t *testing.T) {
	parent1 := NewDNA("PARENT-001")
	parent2 := NewDNA("PARENT-002")

	// Set distinct values for parents
	parent1.Genes["replication_rate"].Value = 0.2
	parent2.Genes["replication_rate"].Value = 0.8

	child := Crossover(parent1, parent2)

	// Verify child ID format
	if !strings.Contains(child.ID, "-") {
		t.Errorf("Invalid child ID format: %s", child.ID)
	}

	// Verify generation increment
	if child.Generation != 2 {
		t.Errorf("Expected child generation 2, got %d", child.Generation)
	}

	// Verify gene inheritance
	replicationRate := child.Genes["replication_rate"].Value
	if replicationRate < 0.2 || replicationRate > 0.8 {
		t.Errorf("Child gene value %.2f outside expected range [0.2, 0.8]", replicationRate)
	}
}

func TestFitnessCalculation(t *testing.T) {
	dna := NewDNA("TEST-003")

	// Set known values for predictable fitness calculation
	totalValue := 0.0
	geneCount := len(dna.Genes)
	for _, gene := range dna.Genes {
		gene.Value = 0.5
		totalValue += 0.5
	}

	expectedFitness := totalValue / float64(geneCount)
	actualFitness := dna.CalculateFitness()

	if actualFitness != expectedFitness {
		t.Errorf("Expected fitness %.2f, got %.2f", expectedFitness, actualFitness)
	}
}
