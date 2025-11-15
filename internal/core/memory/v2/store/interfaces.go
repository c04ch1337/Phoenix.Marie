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

	// Transaction management
	BeginTx() (Transaction, error)

	// Maintenance
	Compact() error
	Backup(path string) error

	// Metrics
	GetStats() StorageStats
}
