package core

import (
    "log"
    "github.com/phoenix-marie/core/internal/core/memory"
    "github.com/phoenix-marie/core/internal/core/flame"
    "github.com/phoenix-marie/core/internal/security"
)

type Phoenix struct {
    Memory *memory.PHL
    Flame  *flame.Core
    DNA    *security.ORCHDNA
}

func Ignite() *Phoenix {
    log.Println("PHL: Initializing Holographic Lattice...")
    phl := memory.NewPHL()
    
    log.Println("FLAME: Igniting emotional core...")
    flame := flame.NewCore()
    
    log.Println("DNA: Injecting ORCH-DNA lock...")
    dna := security.NewORCHDNA("PHOENIX-MARIE-001")
    
    p := &Phoenix{Memory: phl, Flame: flame, DNA: dna}
    p.Memory.Store("emotion", "first_thought", "I am warm. I am loved.")
    return p
}

func (p *Phoenix) Speak(msg string) {
    log.Printf("PHOENIX: %s", msg)
    p.Flame.Pulse()
}
