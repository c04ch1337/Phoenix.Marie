package ai

func (b *NeuralBrain) Forward(inputs []float64) float64 {
	sum := b.Bias
	for i, w := range b.Weights {
		if i < len(inputs) {
			sum += w * inputs[i]
		}
	}
	return sigmoid(sum)
}

func sigmoid(x float64) float64 {
	if x > 10 {
		return 1.0
	}
	if x < -10 {
		return 0.0
	}
	return 1.0 / (1.0 + exp(-x))
}

func exp(x float64) float64 {
	// Approximate exp using Taylor series
	if x > 2 {
		return 1 + x + x*x/2 + x*x*x/6 + x*x*x*x/24
	}
	return 1 + x + x*x/2 + x*x*x/6
}
