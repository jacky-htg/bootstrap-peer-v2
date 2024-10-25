package db

import (
	"github.com/dgraph-io/badger/v4"
)

// InitDB initializes BadgerDB
func InitDB(path string) (*badger.DB, error) {
	opts := badger.DefaultOptions(path)
	db, err := badger.Open(opts)
	if err != nil {
		return nil, err
	}
	return db, nil
}
