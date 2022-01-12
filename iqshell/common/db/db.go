package db

import (
	"github.com/syndtr/goleveldb/leveldb"
	"github.com/syndtr/goleveldb/leveldb/opt"
)

type DB struct {
	filePath string
	db       *leveldb.DB
}

func OpenDB(filePath string) (*DB, error) {
	db, err := leveldb.OpenFile(filePath, nil)
	if err != nil {
		return nil, err
	}

	return &DB{
		filePath: filePath,
		db:       db,
	}, nil
}

func (db *DB) Get(key string) (string, error) {
	ret, err := db.db.Get([]byte(key), nil)
	if err != nil {
		return "", err
	}
	return string(ret), nil
}

func (db *DB) Put(key, value string) error {
	return db.db.Put([]byte(key), []byte(value), &opt.WriteOptions{
		Sync: true,
	})
}

func (db *DB) Delete(key string) error {
	return db.db.Delete([]byte(key), nil)
}
