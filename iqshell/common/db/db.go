package db

import (
	"errors"
	"fmt"
	"github.com/syndtr/goleveldb/leveldb"
	"github.com/syndtr/goleveldb/leveldb/opt"
	"os"
	"path/filepath"
	"sync"
)

type DB struct {
	filePath string
	db       *leveldb.DB
}

var dbMap map[string]*DB
var dbMapLock sync.Mutex

func OpenDB(filePath string) (*DB, error) {
	dbMapLock.Lock()
	defer dbMapLock.Unlock()

	if dbMap == nil {
		dbMap = make(map[string]*DB)
	}

	if dbMap[filePath] != nil {
		return dbMap[filePath], nil
	} else {
		dbDir := filepath.Dir(filePath)
		err := os.MkdirAll(dbDir, 0775)
		if err != nil {
			return nil, errors.New("open db: make file error:" + err.Error())
		}

		handler, err := openDB(filePath)
		if err != nil {
			return nil, errors.New("open db: open error:" + err.Error())
		}
		dbMap[filePath] = handler
		return handler, nil
	}
}

func openDB(filePath string) (*DB, error) {
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
	if db.db == nil {
		return "", fmt.Errorf("db get key:%s error:no db exist", key)
	}
	ret, err := db.db.Get([]byte(key), nil)
	if err != nil {
		return "", err
	}
	return string(ret), nil
}

func (db *DB) Put(key, value string) error {
	if db.db == nil {
		return fmt.Errorf("db put key:%s for value:%s error:no db exist", key, value)
	}
	return db.db.Put([]byte(key), []byte(value), &opt.WriteOptions{
		Sync: true,
	})
}

func (db *DB) Delete(key string) error {
	if db.db == nil {
		return fmt.Errorf("db delete key:%s error:no db exist", key)
	}
	return db.db.Delete([]byte(key), nil)
}
