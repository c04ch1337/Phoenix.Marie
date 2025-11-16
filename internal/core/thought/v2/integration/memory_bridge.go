package integration

import (
	"context"
	"encoding/json"
	"fmt"
	"sync"
	"time"

	"github.com/phoenix-marie/core/internal/core/memory/v2/processor"
	"github.com/phoenix-marie/core/internal/core/memory/v2/store"
	"github.com/phoenix-marie/core/internal/core/thought/v2/learning"
	"github.com/phoenix-marie/core/internal/core/thought/v2/pattern"
)

// MemoryBridge provides integration between thought engine and memory system
type MemoryBridge struct {
	store     store.StorageEngine
	processor processor.LayerProcessor
	patterns  *pattern.Manager
	learning  *learning.Manager
	cache     sync.Map
	cacheTTL  time.Duration
	txManager *TransactionManager
}

// BridgeConfig contains configuration for memory integration
type BridgeConfig struct {
	CacheTTL      time.Duration
	BatchSize     int
	SyncInterval  time.Duration
	RetryAttempts int
}

// NewMemoryBridge creates a new memory integration bridge
func NewMemoryBridge(
	store store.StorageEngine,
	processor processor.LayerProcessor,
	patterns *pattern.Manager,
	learning *learning.Manager,
	config BridgeConfig,
) *MemoryBridge {
	return &MemoryBridge{
		store:     store,
		processor: processor,
		patterns:  patterns,
		learning:  learning,
		cacheTTL:  config.CacheTTL,
		txManager: NewTransactionManager(store),
	}
}

// StorePattern stores a pattern in memory using transaction
func (mb *MemoryBridge) StorePattern(ctx context.Context, p pattern.Pattern) error {
	if err := validatePattern(p); err != nil {
		return fmt.Errorf("invalid pattern: %w", err)
	}

	// Process pattern data
	processed, err := mb.processor.Process(p)
	if err != nil {
		return fmt.Errorf("pattern processing failed: %w", err)
	}

	// Create transaction operation
	op := TransactionOp{
		Type:  "store",
		Layer: "patterns",
		Key:   p.ID,
		Value: processed.Data,
	}

	// Execute transaction
	if err := mb.txManager.ExecuteTransaction([]TransactionOp{op}); err != nil {
		return fmt.Errorf("pattern storage failed: %w", err)
	}

	// Update cache
	mb.updateCache(p.ID, processed.Data)

	return nil
}

// RetrievePattern retrieves a pattern from memory
func (mb *MemoryBridge) RetrievePattern(ctx context.Context, id string) (pattern.Pattern, error) {
	if id == "" {
		return pattern.Pattern{}, fmt.Errorf("pattern ID cannot be empty")
	}

	// Check cache first
	if data, exists := mb.checkCache(id); exists {
		return mb.deserializePattern(data)
	}

	// Retrieve from storage
	data, err := mb.store.Retrieve("patterns", id)
	if err != nil {
		return pattern.Pattern{}, fmt.Errorf("pattern retrieval failed: %w", err)
	}

	// Update cache
	mb.updateCache(id, data)

	return mb.deserializePattern(data)
}

// StoreLearningState stores learning system state
func (mb *MemoryBridge) StoreLearningState(ctx context.Context) error {
	// Get current learning stats
	stats := mb.learning.GetStats()

	// Create transaction operation
	op := TransactionOp{
		Type:  "store",
		Layer: "learning",
		Key:   "state",
		Value: stats,
	}

	// Execute transaction
	return mb.txManager.ExecuteTransaction([]TransactionOp{op})
}

// SyncPatterns synchronizes patterns between thought engine and memory
func (mb *MemoryBridge) SyncPatterns(ctx context.Context) error {
	// Get all patterns from memory
	patterns, err := mb.retrieveAllPatterns(ctx)
	if err != nil {
		return fmt.Errorf("pattern retrieval failed: %w", err)
	}

	// Update pattern manager in batches
	var ops []TransactionOp
	for _, p := range patterns {
		if err := mb.patterns.UpdatePattern(p); err != nil {
			return fmt.Errorf("pattern update failed: %w", err)
		}

		ops = append(ops, TransactionOp{
			Type:  "store",
			Layer: "patterns",
			Key:   p.ID,
			Value: p,
		})

		// Execute batch when it reaches the limit
		if len(ops) >= 100 {
			if err := mb.txManager.ExecuteTransaction(ops); err != nil {
				return fmt.Errorf("batch update failed: %w", err)
			}
			ops = ops[:0]
		}
	}

	// Execute remaining operations
	if len(ops) > 0 {
		if err := mb.txManager.ExecuteTransaction(ops); err != nil {
			return fmt.Errorf("final batch update failed: %w", err)
		}
	}

	// Store analysis
	analysis := mb.patterns.AnalyzePatterns()
	op := TransactionOp{
		Type:  "store",
		Layer: "patterns",
		Key:   "analysis",
		Value: analysis,
	}

	return mb.txManager.ExecuteTransaction([]TransactionOp{op})
}

// ProcessMemoryFeedback processes feedback from memory system
func (mb *MemoryBridge) ProcessMemoryFeedback(ctx context.Context, feedback []learning.Feedback) error {
	var ops []TransactionOp

	for _, f := range feedback {
		// Process feedback through learning system
		if err := mb.learning.Adapt(f); err != nil {
			return fmt.Errorf("feedback processing failed: %w", err)
		}

		// Update pattern confidence based on feedback
		pattern, err := mb.RetrievePattern(ctx, f.PatternID)
		if err != nil {
			continue
		}

		pattern.Confidence = mb.learning.GetProgress()

		// Add pattern update to transaction
		ops = append(ops, TransactionOp{
			Type:  "store",
			Layer: "patterns",
			Key:   pattern.ID,
			Value: pattern,
		})

		// Execute batch when it reaches the limit
		if len(ops) >= 100 {
			if err := mb.txManager.ExecuteTransaction(ops); err != nil {
				return fmt.Errorf("batch update failed: %w", err)
			}
			ops = ops[:0]
		}
	}

	// Execute remaining operations
	if len(ops) > 0 {
		if err := mb.txManager.ExecuteTransaction(ops); err != nil {
			return fmt.Errorf("final batch update failed: %w", err)
		}
	}

	return nil
}

// Helper methods

func (mb *MemoryBridge) checkCache(key string) (interface{}, bool) {
	value, exists := mb.cache.Load(key)
	if !exists {
		return nil, false
	}

	metadata, ok := value.(cacheEntry)
	if !ok || time.Since(metadata.timestamp) > mb.cacheTTL {
		mb.cache.Delete(key)
		return nil, false
	}

	return metadata.data, true
}

func (mb *MemoryBridge) updateCache(key string, data interface{}) {
	mb.cache.Store(key, cacheEntry{
		data:      data,
		timestamp: time.Now(),
	})
}

func (mb *MemoryBridge) retrieveAllPatterns(ctx context.Context) ([]pattern.Pattern, error) {
	var allPatterns []pattern.Pattern
	var lastKey string
	const batchSize = 100

	for {
		// Use BatchRetrieveByPrefix to get a batch of patterns
		batch, err := mb.store.BatchRetrieveByPrefix("patterns", lastKey, batchSize)
		if err != nil {
			return nil, fmt.Errorf("batch retrieval failed: %w", err)
		}

		if len(batch) == 0 {
			break
		}

		// Process each pattern in the batch
		for key, data := range batch {
			p, err := mb.deserializePattern(data)
			if err != nil {
				return nil, fmt.Errorf("pattern deserialization failed for %s: %w", key, err)
			}
			allPatterns = append(allPatterns, p)
			lastKey = key
		}

		// Check context cancellation
		select {
		case <-ctx.Done():
			return nil, ctx.Err()
		default:
		}
	}

	return allPatterns, nil
}

func validatePattern(p pattern.Pattern) error {
	if p.ID == "" {
		return fmt.Errorf("pattern ID cannot be empty")
	}
	if p.Type == "" {
		return fmt.Errorf("pattern type cannot be empty")
	}
	if p.Data == nil {
		return fmt.Errorf("pattern data cannot be nil")
	}
	return nil
}

func (mb *MemoryBridge) deserializePattern(data interface{}) (pattern.Pattern, error) {
	switch v := data.(type) {
	case pattern.Pattern:
		return v, nil
	case map[string]interface{}:
		// Convert to JSON and back to ensure proper type conversion
		jsonData, err := json.Marshal(v)
		if err != nil {
			return pattern.Pattern{}, fmt.Errorf("failed to marshal pattern data: %w", err)
		}

		var p pattern.Pattern
		if err := json.Unmarshal(jsonData, &p); err != nil {
			return pattern.Pattern{}, fmt.Errorf("failed to unmarshal pattern data: %w", err)
		}

		return p, nil
	default:
		return pattern.Pattern{}, fmt.Errorf("unsupported data type for pattern deserialization: %T", data)
	}
}

// Cache entry with timestamp
type cacheEntry struct {
	data      interface{}
	timestamp time.Time
}

// TransactionManager handles memory transactions
type TransactionManager struct {
	store store.StorageEngine
	mu    sync.Mutex
}

// NewTransactionManager creates a new transaction manager
func NewTransactionManager(store store.StorageEngine) *TransactionManager {
	return &TransactionManager{
		store: store,
	}
}

// ExecuteTransaction executes a memory transaction
func (tm *TransactionManager) ExecuteTransaction(ops []TransactionOp) error {
	if len(ops) == 0 {
		return nil
	}

	tm.mu.Lock()
	defer tm.mu.Unlock()

	tx, err := tm.store.BeginTx()
	if err != nil {
		return fmt.Errorf("transaction start failed: %w", err)
	}

	for _, op := range ops {
		if err := tm.executeOperation(tx, op); err != nil {
			tx.Rollback()
			return fmt.Errorf("operation failed: %w", err)
		}
	}

	return tx.Commit()
}

// TransactionOp represents a memory operation
type TransactionOp struct {
	Type  string // "store", "retrieve", "delete"
	Layer string
	Key   string
	Value interface{}
}

func (tm *TransactionManager) executeOperation(tx store.Transaction, op TransactionOp) error {
	switch op.Type {
	case "store":
		return tx.Store(op.Layer, op.Key, op.Value)
	case "delete":
		return tx.Delete(op.Layer, op.Key)
	default:
		return fmt.Errorf("unsupported operation: %s", op.Type)
	}
}
