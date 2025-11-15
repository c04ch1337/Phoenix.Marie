package orch

import (
	"log"
	v2 "github.com/phoenix-marie/core/internal/orch/v2"
)

// Legacy Army struct for backward compatibility
// This redirects to the evolved v2 system
type Army struct {
	Count    int
	Interval int
	evolved  *v2.EvolvedArmy
}

// NewArmy creates a legacy army that redirects to v2
func NewArmy() *Army {
	log.Println("ORCH: Legacy army → upgrading to v2")
	evolved := v2.NewEvolvedArmy()
	return &Army{
		Count:    evolved.Count,
		Interval: evolved.Interval,
		evolved:  evolved,
	}
}

// Deploy redirects to v2 deployment
func (a *Army) Deploy() {
	log.Println("ORCH: Legacy deploy → redirecting to v2 evolved army")
	if a.evolved != nil {
		a.evolved.Deploy()
	} else {
		evolved := v2.NewEvolvedArmy()
		evolved.Deploy()
	}
}
