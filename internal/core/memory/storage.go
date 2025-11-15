package memory

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	badger "github.com/dgraph-io/badger/v3"
)

type Storage struct {
	db *badger.DB
}

// Backup creates a backup of the database
func (s *Storage) Backup(path string) error {
	file, err := os.Create(path)
	if err != nil {
		return fmt.Errorf("failed to create backup file: %w", err)
	}
	defer file.Close()

	_, err = s.db.Backup(file, 0)
	if err != nil {
		return fmt.Errorf("failed to create backup: %w", err)
	}

	return nil
}

// GetDB returns the underlying BadgerDB instance (for advanced operations)
func (s *Storage) GetDB() *badger.DB {
	return s.db
}

func NewStorage(dataDir string) (*Storage, error) {
	opts := badger.DefaultOptions(filepath.Join(dataDir, "phl-memory"))
	opts.Logger = nil // Disable Badger's internal logger

	db, err := badger.Open(opts)
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %w", err)
	}

	return &Storage{db: db}, nil
}

func (s *Storage) Close() error {
	return s.db.Close()
}

func (s *Storage) Store(layer, key string, value any) error {
	data, err := json.Marshal(value)
	if err != nil {
		return fmt.Errorf("failed to marshal value: %w", err)
	}

	dbKey := []byte(fmt.Sprintf("%s:%s", layer, key))
	return s.db.Update(func(txn *badger.Txn) error {
		return txn.Set(dbKey, data)
	})
}

func (s *Storage) Retrieve(layer, key string) (any, error) {
	var value any
	dbKey := []byte(fmt.Sprintf("%s:%s", layer, key))

	err := s.db.View(func(txn *badger.Txn) error {
		item, err := txn.Get(dbKey)
		if err != nil {
			return err
		}

		return item.Value(func(val []byte) error {
			return json.Unmarshal(val, &value)
		})
	})

	if err == badger.ErrKeyNotFound {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve value: %w", err)
	}

	return value, nil
}

func (s *Storage) DeleteLayer(layer string) error {
	prefix := []byte(fmt.Sprintf("%s:", layer))

	return s.db.Update(func(txn *badger.Txn) error {
		opts := badger.DefaultIteratorOptions
		opts.PrefetchValues = false

		it := txn.NewIterator(opts)
		defer it.Close()

		for it.Seek(prefix); it.ValidForPrefix(prefix); it.Next() {
			key := it.Item().Key()
			if err := txn.Delete(key); err != nil {
				return fmt.Errorf("failed to delete key %s: %w", key, err)
			}
		}
		return nil
	})
}
