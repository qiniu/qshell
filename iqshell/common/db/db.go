package db

import (
	"github.com/qiniu/qshell/v2/iqshell/common/data"
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

func OpenDB(filePath string) (*DB, *data.CodeError) {
	dbMapLock.Lock()
	defer dbMapLock.Unlock()

	if dbMap == nil {
		dbMap = make(map[string]*DB)
	}

	if dbMap[filePath] != nil {
		return dbMap[filePath], nil
	} else {
		dbDir := filepath.Dir(filePath)
		if e := os.MkdirAll(dbDir, os.ModePerm); e != nil {
			return nil, data.NewEmptyError().AppendDesc("open db, make file").AppendError(e)
		}

		handler, err := openDB(filePath)
		if err != nil {
			return nil, data.NewEmptyError().AppendDesc("open db: open").AppendError(err)
		}
		dbMap[filePath] = handler
		return handler, nil
	}
}

func openDB(filePath string) (*DB, *data.CodeError) {
	db, err := leveldb.OpenFile(filePath, nil)
	if err != nil {
		return nil, data.NewEmptyError().AppendError(err)
	}

	return &DB{
		filePath: filePath,
		db:       db,
	}, nil
}

func (db *DB) Get(key string) (string, *data.CodeError) {
	if db.db == nil {
		return "", data.NewEmptyError().AppendDescF("db get key:%s error:no db exist", key)
	}
	ret, err := db.db.Get([]byte(key), nil)
	if err != nil {
		return "", data.NewEmptyError().AppendError(err)
	}
	return string(ret), nil
}

func (db *DB) Put(key, value string) *data.CodeError {
	if db.db == nil {
		return data.NewEmptyError().AppendDescF("db put key:%s for value:%s error:no db exist", key, value)
	}
	if e := db.db.Put([]byte(key), []byte(value), &opt.WriteOptions{
		Sync: false,
	}); e != nil {
		return data.NewEmptyError().AppendError(e)
	}
	return nil
}

func (db *DB) Delete(key string) *data.CodeError {
	if db.db == nil {
		return data.NewEmptyError().AppendDescF("db delete key:%s error:no db exist", key)
	}
	if e := db.db.Delete([]byte(key), nil); e != nil {
		return data.NewEmptyError().AppendError(e)
	}
	return nil
}
