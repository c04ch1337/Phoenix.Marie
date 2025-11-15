package store

import (
	"encoding/json"
	"fmt"
	"os"

	badger "github.com/dgraph-io/badger/v3"
)

// Transaction defines the interface for storage transactions
type Transaction interface {
	Store(layer, key string, value any) error
	Retrieve(layer, key string) (any, error)
	Delete(layer, key string) error
	Commit() error
	Rollback() error
}

// BadgerStore implements the StorageEngine interface using BadgerDB
type BadgerStore struct {
	db      *badger.DB
	options *badger.Options
}

// StoreOperation represents a single store operation for batch processing
type StoreOperation struct {
	Layer string
	Key   string
	Value any
}

// Query represents a single retrieval query
type Query struct {
	Layer string
	Key   string
}

// QueryResult represents the result of a query operation
type QueryResult struct {
	Value any
	Error error
}

// StorageStats contains metrics about the storage engine
type StorageStats struct {
	ItemCount     uint64
	LSMSize       int64
	VLogSize      int64
	PendingWrites int64
}

// NewBadgerStore creates a new BadgerDB storage instance
func NewBadgerStore(path string) (*BadgerStore, error) {
	opts := badger.DefaultOptions(path)
	opts.NumCompactors = 2
	opts.NumLevelZeroTables = 3
	opts.NumMemtables = 2
	opts.ValueLogFileSize = 1 << 28 // 256MB

	db, err := badger.Open(opts)
	if err != nil {
		return nil, fmt.Errorf("failed to open BadgerDB: %w", err)
	}

	return &BadgerStore{
		db:      db,
		options: &opts,
	}, nil
}

// Store implements the Store method of StorageEngine
func (bs *BadgerStore) Store(layer, key string, value any) error {
	data, err := json.Marshal(value)
	if err != nil {
		return fmt.Errorf("failed to marshal value: %w", err)
	}

	compositeKey := fmt.Sprintf("%s:%s", layer, key)
	return bs.db.Update(func(txn *badger.Txn) error {
		return txn.Set([]byte(compositeKey), data)
	})
}

// Retrieve implements the Retrieve method of StorageEngine
func (bs *BadgerStore) Retrieve(layer, key string) (any, error) {
	var value any
	compositeKey := fmt.Sprintf("%s:%s", layer, key)

	err := bs.db.View(func(txn *badger.Txn) error {
		item, err := txn.Get([]byte(compositeKey))
		if err != nil {
			return err
		}

		return item.Value(func(val []byte) error {
			return json.Unmarshal(val, &value)
		})
	})

	if err != nil {
		return nil, fmt.Errorf("failed to retrieve value: %w", err)
	}

	return value, nil
}

// Delete implements the Delete method of StorageEngine
func (bs *BadgerStore) Delete(layer, key string) error {
	compositeKey := fmt.Sprintf("%s:%s", layer, key)
	return bs.db.Update(func(txn *badger.Txn) error {
		return txn.Delete([]byte(compositeKey))
	})
}

// BatchStore implements the BatchStore method of StorageEngine
func (bs *BadgerStore) BatchStore(operations []StoreOperation) error {
	wb := bs.db.NewWriteBatch()
	defer wb.Cancel()

	for _, op := range operations {
		data, err := json.Marshal(op.Value)
		if err != nil {
			return fmt.Errorf("failed to marshal value for key %s: %w", op.Key, err)
		}

		compositeKey := fmt.Sprintf("%s:%s", op.Layer, op.Key)
		if err := wb.Set([]byte(compositeKey), data); err != nil {
			return fmt.Errorf("failed to batch set key %s: %w", compositeKey, err)
		}
	}

	return wb.Flush()
}

// BatchRetrieve implements the BatchRetrieve method of StorageEngine
func (bs *BadgerStore) BatchRetrieve(queries []Query) ([]QueryResult, error) {
	results := make([]QueryResult, len(queries))

	err := bs.db.View(func(txn *badger.Txn) error {
		for i, query := range queries {
			compositeKey := fmt.Sprintf("%s:%s", query.Layer, query.Key)
			item, err := txn.Get([]byte(compositeKey))
			if err != nil {
				results[i] = QueryResult{nil, err}
				continue
			}

			var value any
			err = item.Value(func(val []byte) error {
				return json.Unmarshal(val, &value)
			})

			results[i] = QueryResult{value, err}
		}
		return nil
	})

	if err != nil {
		return nil, fmt.Errorf("batch retrieve failed: %w", err)
	}

	return results, nil
}

// BeginTx implements the BeginTx method of StorageEngine
func (bs *BadgerStore) BeginTx() (Transaction, error) {
	txn := bs.db.NewTransaction(true)
	return &BadgerTransaction{txn: txn}, nil
}

// Compact implements the Compact method of StorageEngine
func (bs *BadgerStore) Compact() error {
	return bs.db.RunValueLogGC(0.5)
}

// Backup implements the Backup method of StorageEngine
func (bs *BadgerStore) Backup(path string) error {
	file, err := os.Create(path)
	if err != nil {
		return fmt.Errorf("failed to create backup file: %w", err)
	}
	defer file.Close()

	_, err = bs.db.Backup(file, 0)
	if err != nil {
		return fmt.Errorf("failed to create backup: %w", err)
	}

	return nil
}

// GetStats implements the GetStats method of StorageEngine
func (bs *BadgerStore) GetStats() StorageStats {
	lsmSize, vlogSize := bs.db.Size()
	return StorageStats{
		LSMSize:       lsmSize,
		VLogSize:      vlogSize,
		PendingWrites: bs.db.MaxBatchCount(),
	}
}

// BadgerTransaction implements the Transaction interface
type BadgerTransaction struct {
	txn *badger.Txn
}

func (bt *BadgerTransaction) Store(layer, key string, value any) error {
	data, err := json.Marshal(value)
	if err != nil {
		return fmt.Errorf("failed to marshal value: %w", err)
	}

	compositeKey := fmt.Sprintf("%s:%s", layer, key)
	return bt.txn.Set([]byte(compositeKey), data)
}

func (bt *BadgerTransaction) Retrieve(layer, key string) (any, error) {
	var value any
	compositeKey := fmt.Sprintf("%s:%s", layer, key)

	item, err := bt.txn.Get([]byte(compositeKey))
	if err != nil {
		return nil, err
	}

	err = item.Value(func(val []byte) error {
		return json.Unmarshal(val, &value)
	})

	if err != nil {
		return nil, fmt.Errorf("failed to retrieve value: %w", err)
	}

	return value, nil
}

func (bt *BadgerTransaction) Delete(layer, key string) error {
	compositeKey := fmt.Sprintf("%s:%s", layer, key)
	return bt.txn.Delete([]byte(compositeKey))
}

func (bt *BadgerTransaction) Commit() error {
	return bt.txn.Commit()
}

func (bt *BadgerTransaction) Rollback() error {
	bt.txn.Discard()
	return nil
}
