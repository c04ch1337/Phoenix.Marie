package memory

import (
	"fmt"
	"time"
)

// EternalMemoryManager manages long-term eternal memories
type EternalMemoryManager struct {
	phl *PHL
}

// NewEternalMemoryManager creates a new eternal memory manager
func NewEternalMemoryManager(phl *PHL) *EternalMemoryManager {
	return &EternalMemoryManager{phl: phl}
}

// StoreEternal stores a memory in the eternal layer (long-term preservation)
func (emm *EternalMemoryManager) StoreEternal(key string, value any, importance int) bool {
	// Eternal memories should have importance > 0
	if importance <= 0 {
		importance = 1
	}

	eternalData := map[string]interface{}{
		"data":       value,
		"importance": importance,
		"stored_at":  time.Now().Unix(),
		"layer":      "eternal",
	}

	return emm.phl.Store("eternal", key, eternalData)
}

// RetrieveEternal retrieves a memory from the eternal layer
func (emm *EternalMemoryManager) RetrieveEternal(key string) (any, bool) {
	return emm.phl.Retrieve("eternal", key)
}

// ListEternalMemories lists all eternal memories
func (emm *EternalMemoryManager) ListEternalMemories() ([]EternalMemory, error) {
	// Note: This would require iterating over all keys in eternal layer
	// BadgerDB doesn't have a direct "list all keys" method without iteration
	// This is a placeholder for future implementation
	
	return []EternalMemory{}, fmt.Errorf("list operation not yet implemented - use Retrieve with specific keys")
}

// GetEternalStats returns statistics about eternal memories
func (emm *EternalMemoryManager) GetEternalStats() EternalStats {
	// This would require counting keys in eternal layer
	// Placeholder implementation
	return EternalStats{
		TotalMemories: 0,
		Layer:         "eternal",
		Status:        "operational",
	}
}

// EternalMemory represents an eternal memory entry
type EternalMemory struct {
	Key        string
	Data       any
	Importance int
	StoredAt   time.Time
}

// EternalStats contains statistics about eternal memories
type EternalStats struct {
	TotalMemories int
	Layer         string
	Status        string
}

// PromoteToEternal promotes a memory from another layer to eternal
func (emm *EternalMemoryManager) PromoteToEternal(sourceLayer, key string, importance int) bool {
	value, exists := emm.phl.Retrieve(sourceLayer, key)
	if !exists {
		return false
	}

	// Store in eternal layer
	success := emm.StoreEternal(key, value, importance)
	if success {
		// Optionally remove from source layer (or keep as reference)
		// emm.phl.Cleanup(sourceLayer) // Only if we want to move, not copy
	}

	return success
}

