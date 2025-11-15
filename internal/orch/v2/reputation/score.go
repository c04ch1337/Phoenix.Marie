package reputation

import (
	"sync"
)

type ReputationSystem struct {
	Scores map[string]float64
	mu     sync.RWMutex
}

func NewSystem() *ReputationSystem {
	return &ReputationSystem{
		Scores: make(map[string]float64),
	}
}

func (r *ReputationSystem) Record(agent, event string, value float64) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.Scores[agent] += value
	if r.Scores[agent] < 0 {
		r.Scores[agent] = 0
	}
	if r.Scores[agent] > 100 {
		r.Scores[agent] = 100
	}
}

func (r *ReputationSystem) Get(agent string) float64 {
	r.mu.RLock()
	defer r.mu.RUnlock()
	if score, ok := r.Scores[agent]; ok {
		return score
	}
	return 50.0 // Default neutral reputation
}

func (r *ReputationSystem) GetEmotionForReputation(agent string) string {
	rep := r.Get(agent)
	if rep > 80 {
		return "pride"
	} else if rep > 60 {
		return "satisfaction"
	} else if rep < 30 {
		return "concern"
	} else if rep < 20 {
		return "sadness"
	}
	return "neutral"
}
