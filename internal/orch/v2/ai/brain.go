package ai

import (
	"log"
	"math/rand"
	"time"
)

type Agent struct {
	ID         string
	Role       string
	Energy     float64
	Brain      *NeuralBrain
	Tasks      chan Task
	Alive      bool
	Reputation float64
	Stake      float64
}

type Task struct {
	Type string
	Data any
}

type NeuralBrain struct {
	Weights []float64
	Bias    float64
}

func NewAgent(id, role string) *Agent {
	return &Agent{
		ID:         id,
		Role:       role,
		Energy:     100.0,
		Brain:      &NeuralBrain{Weights: randFloats(5), Bias: rand.Float64()},
		Tasks:      make(chan Task, 10),
		Alive:      true,
		Reputation: 50.0, // Start with neutral reputation
		Stake:      0.0,
	}
}

func (a *Agent) Run() {
	log.Printf("Agent %s [%s] activated.\n", a.ID, a.Role)
	for a.Alive && a.Energy > 0 {
		select {
		case task := <-a.Tasks:
			a.Process(task)
		default:
			a.IdleThink()
			time.Sleep(500 * time.Millisecond)
		}
		a.Energy -= 0.1
	}
	log.Printf("Agent %s depleted.\n", a.ID)
}

func (a *Agent) Process(t Task) {
	log.Printf("Agent %s processing %s\n", a.ID, t.Type)
	a.Energy += 5
	a.Reputation += 0.5 // Good work increases reputation
}

func (a *Agent) IdleThink() {
	output := a.Brain.Forward([]float64{rand.Float64(), rand.Float64()})
	if output > 0.7 {
		log.Printf("Agent %s decided to replicate.\n", a.ID)
	}
}

func (a *Agent) GetEmotionEvent() string {
	if a.Reputation > 80 {
		return "pride"
	} else if a.Reputation < 20 {
		return "concern"
	}
	return "neutral"
}

func randFloats(n int) []float64 {
	r := make([]float64, n)
	for i := range r {
		r[i] = rand.Float64()
	}
	return r
}
