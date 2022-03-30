package recorder

import (
	"errors"
	"fmt"
	"github.com/syndtr/goleveldb/leveldb"
	"github.com/syndtr/goleveldb/leveldb/opt"
	"os"
	"path/filepath"
	"sync"
)

type dbRecorder struct {
	filePath string
	db       *leveldb.DB
}

var dbMap map[string]*dbRecorder
var dbMapLock sync.Mutex

func CreateDBRecorder(filePath string) (Recorder, error) {
	dbMapLock.Lock()
	defer dbMapLock.Unlock()

	if dbMap == nil {
		dbMap = make(map[string]*dbRecorder)
	}

	if dbMap[filePath] != nil {
		return dbMap[filePath], nil
	} else {
		dbDir := filepath.Dir(filePath)
		err := os.MkdirAll(dbDir, os.ModePerm)
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

func openDB(filePath string) (*dbRecorder, error) {
	db, err := leveldb.OpenFile(filePath, nil)
	if err != nil {
		return nil, err
	}

	return &dbRecorder{
		filePath: filePath,
		db:       db,
	}, nil
}

func (db *dbRecorder) Get(key string) (string, error) {
	if db.db == nil {
		return "", fmt.Errorf("db get key:%s error:no db exist", key)
	}
	ret, err := db.db.Get([]byte(key), nil)
	if err != nil {
		return "", err
	}
	return string(ret), nil
}

func (db *dbRecorder) Put(key, value string) error {
	if db.db == nil {
		return fmt.Errorf("db put key:%s for value:%s error:no db exist", key, value)
	}
	return db.db.Put([]byte(key), []byte(value), &opt.WriteOptions{
		Sync: true,
	})
}

func (db *dbRecorder) Delete(key string) error {
	if db.db == nil {
		return fmt.Errorf("db delete key:%s error:no db exist", key)
	}
	return db.db.Delete([]byte(key), nil)
}
