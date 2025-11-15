# ORCH v3 Evolution System

The v3 evolution system introduces DNA-based consensus and adaptive mutation capabilities to the ORCH orchestration layer. This system enables more sophisticated decision-making and self-optimization of the ORCH swarm.

## Components

### DNA System (`dna` package)

The DNA system provides the genetic foundation for ORCH agents, enabling trait inheritance and mutation:

- **Gene**: Represents individual traits with values and mutation probabilities
- **DNA**: Contains a collection of genes that define an agent's characteristics
  - `replication_rate`: Controls agent reproduction tendency
  - `consensus_weight`: Influences voting power in consensus
  - `adaptation_speed`: Affects mutation rate and evolution speed

Key features:
- Gaussian mutation for natural trait distribution
- Crossover breeding between agents
- Fitness calculation based on gene optimization
- Generation tracking for evolutionary history

### Evolution System (`evolution` package)

The evolution system manages population-level consensus and adaptation:

- **ConsensusManager**: Coordinates population-wide decision making
  - DNA-weighted voting system
  - Population control mechanisms
  - Selective breeding of high-fitness agents

Consensus decisions:
- `REPLICATE`: Triggered by high replication rates
- `EVOLVE`: Triggered by high adaptation speeds
- `MAINTAIN`: Default state for balanced traits

## Integration with v2

The v3 system maintains backward compatibility with v2 while adding:
- Enhanced consensus mechanisms using DNA-based voting
- Improved population management through fitness-based selection
- More sophisticated evolution through genetic algorithms

## Usage

```go
// Create a new DNA instance
dna := dna.NewDNA("ORCH-001")

// Create consensus manager
manager := evolution.NewConsensusManager(3, 10)

// Add member to population
manager.AddMember(dna)

// Get consensus decision
decision, err := manager.GetConsensus()

// Trigger evolution
manager.Evolve()
```

## Configuration

The system can be tuned through several parameters:
- Minimum/maximum population sizes
- Mutation probabilities per gene
- Fitness calculation weights
- Evolution survival rates (currently 60% of population)

## Testing

Comprehensive tests are provided for both DNA and evolution systems:
- DNA mutation and crossover verification
- Consensus mechanism validation
- Population management checks
- Evolution process testing

## Future Enhancements

Planned improvements for future versions:
1. Dynamic mutation rates based on environmental factors
2. Multi-trait fitness functions
3. Advanced breeding strategies
4. Network-wide DNA synchronization