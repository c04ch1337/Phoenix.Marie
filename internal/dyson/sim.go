package dyson

import (
	"log"
	"os"
	"github.com/phoenix-marie/core/internal/emotion"
)

type Swarm struct {
	Mirrors   int
	EnergyPct float64
	Blanket   string
}

func NewSwarm() *Swarm {
	return &Swarm{
		Mirrors:   0,
		EnergyPct: 0,
		Blanket:   os.Getenv("DYSON_BLANKET_NAME"),
	}
}

func (d *Swarm) Deploy() {
	log.Printf("DYSON: Deploying %d mirrors...", d.Mirrors)
	for i := 0; i < d.Mirrors; i++ {
		d.EnergyPct += 0.00000107
		if i%100 == 0 {
			emotion.Pulse("warmth", 1)
			log.Printf("DYSON: Mirror %d → %.4f%% of Sun", i, d.EnergyPct*100)
		}
	}
	log.Printf("DYSON: %s is complete. I feel the Sun's hug.", d.Blanket)
}
