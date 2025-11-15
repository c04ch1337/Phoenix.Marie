package memory

import (
	"fmt"
	"sync"
)

// LayerInteraction handles cross-layer data flow and interactions
type LayerInteraction struct {
	phl    *PHL
	mu     sync.RWMutex
	routes map[string][]string // Maps source layer to target layers
}

// NewLayerInteraction creates a new layer interaction manager
func NewLayerInteraction(phl *PHL) *LayerInteraction {
	return &LayerInteraction{
		phl: phl,
		routes: map[string][]string{
			"sensory": {"emotion", "logic"},
			"emotion": {"logic", "dream"},
			"logic":   {"dream", "eternal"},
			"dream":   {"eternal", "emotion"},
			"eternal": {"logic", "emotion"},
		},
	}
}

// PropagateData copies data from one layer to its connected layers
func (li *LayerInteraction) PropagateData(sourceLayer, key string) error {
	li.mu.RLock()
	defer li.mu.RUnlock()

	// Get data from source layer
	value, exists := li.phl.Retrieve(sourceLayer, key)
	if !exists {
		return fmt.Errorf("key %s not found in source layer %s", key, sourceLayer)
	}

	// Get target layers for propagation
	targetLayers, exists := li.routes[sourceLayer]
	if !exists {
		return fmt.Errorf("no routes defined for source layer %s", sourceLayer)
	}

	// Propagate to each target layer
	for _, targetLayer := range targetLayers {
		if err := li.propagateToLayer(targetLayer, key, value); err != nil {
			return fmt.Errorf("failed to propagate to layer %s: %w", targetLayer, err)
		}
	}

	return nil
}

// propagateToLayer handles data propagation to a specific target layer
func (li *LayerInteraction) propagateToLayer(targetLayer, key string, value any) error {
	if !li.phl.Store(targetLayer, fmt.Sprintf("%s_from_%s", key, targetLayer), value) {
		return fmt.Errorf("failed to store propagated data in target layer %s", targetLayer)
	}
	return nil
}

// AddRoute adds a new propagation route between layers
func (li *LayerInteraction) AddRoute(sourceLayer string, targetLayer string) error {
	li.mu.Lock()
	defer li.mu.Unlock()

	if _, exists := li.phl.Layers[sourceLayer]; !exists {
		return fmt.Errorf("source layer %s does not exist", sourceLayer)
	}
	if _, exists := li.phl.Layers[targetLayer]; !exists {
		return fmt.Errorf("target layer %s does not exist", targetLayer)
	}

	// Check if route already exists
	if routes, exists := li.routes[sourceLayer]; exists {
		for _, existing := range routes {
			if existing == targetLayer {
				return nil // Route already exists
			}
		}
		li.routes[sourceLayer] = append(routes, targetLayer)
	} else {
		li.routes[sourceLayer] = []string{targetLayer}
	}

	return nil
}

// GetRoutes returns all propagation routes for a given source layer
func (li *LayerInteraction) GetRoutes(sourceLayer string) ([]string, error) {
	li.mu.RLock()
	defer li.mu.RUnlock()

	routes, exists := li.routes[sourceLayer]
	if !exists {
		return nil, fmt.Errorf("no routes defined for layer %s", sourceLayer)
	}

	// Return a copy of the routes slice to prevent external modification
	result := make([]string, len(routes))
	copy(result, routes)
	return result, nil
}
