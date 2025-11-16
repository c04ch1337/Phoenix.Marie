package store

// StorageEngine defines the interface for storage operations
type StorageEngine interface {
	// Core operations
	Store(layer, key string, value any) error
	Retrieve(layer, key string) (any, error)
	Delete(layer, key string) error

	// Batch operations
	BatchStore(operations []StoreOperation) error
	BatchRetrieve(queries []Query) ([]QueryResult, error)
	BatchRetrieveByPrefix(layer, prefix string, limit int) (map[string]any, error)

	// Transaction management
	BeginTx() (Transaction, error)

	// Maintenance
	Compact() error
	Backup(path string) error

	// Metrics
	GetStats() StorageStats
}

// Transaction represents an atomic set of storage operations
type Transaction interface {
	Store(layer, key string, value any) error
	Delete(layer, key string) error
	Commit() error
	Rollback() error
}

// StoreOperation represents a single storage operation
type StoreOperation struct {
	Layer string
	Key   string
	Value any
}

// Query represents a storage query
type Query struct {
	Layer string
	Key   string
}

// QueryResult represents a query result
type QueryResult struct {
	Key   string
	Value any
	Error error
}

// StorageStats contains storage engine metrics
type StorageStats struct {
	TotalEntries   int64
	TotalSize      int64
	CacheHitRate   float64
	OperationCount int64
	ErrorCount     int64
}
