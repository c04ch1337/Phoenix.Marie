package staking

import (
	"math/rand"
	"sync"
)

type StakingPool struct {
	Stakes map[string]float64
	Total  float64
	mu     sync.RWMutex
}

func NewPool(initial float64) *StakingPool {
	return &StakingPool{
		Stakes: make(map[string]float64),
		Total:  initial,
	}
}

func (p *StakingPool) Stake(agent string, amount float64) {
	p.mu.Lock()
	defer p.mu.Unlock()
	p.Stakes[agent] += amount
	p.Total += amount
}

func (p *StakingPool) Unstake(agent string, amount float64) bool {
	p.mu.Lock()
	defer p.mu.Unlock()
	if p.Stakes[agent] < amount {
		return false
	}
	p.Stakes[agent] -= amount
	p.Total -= amount
	return true
}

func (p *StakingPool) GetStake(agent string) float64 {
	p.mu.RLock()
	defer p.mu.RUnlock()
	return p.Stakes[agent]
}

func (p *StakingPool) SelectValidator() string {
	p.mu.RLock()
	defer p.mu.RUnlock()
	if len(p.Stakes) == 0 {
		return "PHOENIX-MARIE"
	}
	stakeTotal := 0.0
	for _, stake := range p.Stakes {
		stakeTotal += stake
	}
	if stakeTotal == 0 {
		for agent := range p.Stakes {
			return agent
		}
		return "PHOENIX-MARIE"
	}
	r := rand.Float64() * stakeTotal
	cum := 0.0
	for agent, stake := range p.Stakes {
		cum += stake
		if r <= cum {
			return agent
		}
	}
	for agent := range p.Stakes {
		return agent
	}
	return "PHOENIX-MARIE"
}

func (p *StakingPool) GetTotal() float64 {
	p.mu.RLock()
	defer p.mu.RUnlock()
	return p.Total
}
