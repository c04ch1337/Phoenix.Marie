package memory

import (
	"fmt"
	"log"
	"os"
)

type PHL struct {
	Layers      map[string]*Layer
	log         *log.Logger
	storage     *Storage
	interaction *LayerInteraction
	processors  *ProcessorManager
	validator   *LayerValidator
}

type Layer struct {
	Name string
	Data map[string]any
}

func NewPHL(dataDir string) (*PHL, error) {
	logger := log.New(os.Stdout, "PHL_MEMORY: ", log.Ldate|log.Ltime|log.Lmicroseconds)

	storage, err := NewStorage(dataDir)
	if err != nil {
		return nil, fmt.Errorf("failed to initialize storage: %w", err)
	}

	layers := map[string]*Layer{
		"sensory": {Name: "Sensory", Data: make(map[string]any)},
		"emotion": {Name: "Emotion", Data: make(map[string]any)},
		"logic":   {Name: "Logic", Data: make(map[string]any)},
		"dream":   {Name: "Dream", Data: make(map[string]any)},
		"eternal": {Name: "Eternal", Data: make(map[string]any)},
	}

	phl := &PHL{
		Layers:     layers,
		log:        logger,
		storage:    storage,
		processors: NewProcessorManager(),
		validator:  NewLayerValidator(),
	}

	phl.interaction = NewLayerInteraction(phl)
	return phl, nil
}

func (p *PHL) Store(layer, key string, value any) bool {
	// Validate layer and key
	if err := ValidateLayer(layer); err != nil {
		p.log.Printf("Layer validation failed: %v", err)
		return false
	}
	if err := ValidateKey(key); err != nil {
		p.log.Printf("Key validation failed: %v", err)
		return false
	}

	// Validate data for the specific layer
	if err := p.validator.ValidateLayerData(layer, value); err != nil {
		p.log.Printf("Data validation failed for %s layer: %v", layer, err)
		return false
	}

	if l, ok := p.Layers[layer]; ok {
		// Process the data using the appropriate processor
		processed, err := p.processors.ProcessData(layer, value)
		if err != nil {
			p.log.Printf("Failed to process data for %s layer: %s (%v)", layer, key, err)
			return false
		}

		l.Data[key] = processed
		if err := p.storage.Store(layer, key, processed); err != nil {
			p.log.Printf("Failed to persist in %s layer: %s (%v)", layer, key, err)
			return false
		}
		p.log.Printf("Stored in %s layer: %s", layer, key)
		return true
	}
	p.log.Printf("Failed to store in %s layer: %s (layer not found)", layer, key)
	return false
}

func (p *PHL) Retrieve(layer, key string) (any, bool) {
	// Validate layer and key
	if err := ValidateLayer(layer); err != nil {
		p.log.Printf("Layer validation failed: %v", err)
		return nil, false
	}
	if err := ValidateKey(key); err != nil {
		p.log.Printf("Key validation failed: %v", err)
		return nil, false
	}

	if l, ok := p.Layers[layer]; ok {
		// Try memory first
		if value, exists := l.Data[key]; exists {
			p.log.Printf("Retrieved from %s layer memory: %s", layer, key)
			return value, true
		}

		// Try persistent storage
		if value, err := p.storage.Retrieve(layer, key); err == nil && value != nil {
			l.Data[key] = value // Cache in memory
			p.log.Printf("Retrieved from %s layer storage: %s", layer, key)
			return value, true
		}

		p.log.Printf("Key not found in %s layer: %s", layer, key)
		return nil, false
	}
	p.log.Printf("Failed to retrieve from %s layer: %s (layer not found)", layer, key)
	return nil, false
}

func (p *PHL) Cleanup(layer string) bool {
	if l, ok := p.Layers[layer]; ok {
		l.Data = make(map[string]any)
		if err := p.storage.DeleteLayer(layer); err != nil {
			p.log.Printf("Failed to cleanup %s layer storage: %v", layer, err)
			return false
		}
		p.log.Printf("Cleaned up %s layer", layer)
		return true
	}
	p.log.Printf("Failed to cleanup %s layer (not found)", layer)
	return false
}

func (p *PHL) Close() error {
	return p.storage.Close()
}

// Backup creates a backup of the memory database
func (p *PHL) Backup(path string) error {
	return p.storage.Backup(path)
}

// GetStorage returns the storage instance (for advanced operations)
func (p *PHL) GetStorage() *Storage {
	return p.storage
}

func (p *PHL) PropagateData(sourceLayer, key string) error {
	return p.interaction.PropagateData(sourceLayer, key)
}

func (p *PHL) AddLayerRoute(sourceLayer, targetLayer string) error {
	return p.interaction.AddRoute(sourceLayer, targetLayer)
}

func (p *PHL) GetLayerRoutes(sourceLayer string) ([]string, error) {
	return p.interaction.GetRoutes(sourceLayer)
}
