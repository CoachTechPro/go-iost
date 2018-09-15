package kv

import (
	"github.com/iost-official/Go-IOS-Protocol/db/kv/leveldb"
	"github.com/iost-official/Go-IOS-Protocol/db/kv/rocksdb"
)

// StorageType is the type of storage, include leveldb and rocksdb
type StorageType uint8

// Storage type constant
const (
	_ StorageType = iota
	LevelDBStorage
	RocksDBStorage
)

// StorageBackend is the storage backend interface
type StorageBackend interface {
	Get(key []byte) ([]byte, error)
	Put(key []byte, value []byte) error
	Has(key []byte) (bool, error)
	Delete(key []byte) error
	Keys(prefix []byte) ([][]byte, error)
	BeginBatch() error
	CommitBatch() error
	Close() error
}

// Storage is a kv database
type Storage struct {
	StorageBackend
}

// NewStorage return the storage of the specify type
func NewStorage(path string, t StorageType) (*Storage, error) {
	switch t {
	case LevelDBStorage:
		sb, err := leveldb.NewDB(path)
		if err != nil {
			return nil, err
		}
		return &Storage{StorageBackend: sb}, nil
	case RocksDBStorage:
		sb, err := rocksdb.NewDB(path)
		if err != nil {
			return nil, err
		}
		return &Storage{StorageBackend: sb}, nil
	default:
		sb, err := leveldb.NewDB(path)
		if err != nil {
			return nil, err
		}
		return &Storage{StorageBackend: sb}, nil
	}
}
