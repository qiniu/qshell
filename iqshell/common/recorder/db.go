package recorder

import (
	"github.com/qiniu/qshell/v2/iqshell/common/data"
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

func CreateDBRecorder(filePath string) (Recorder, *data.CodeError) {
	dbMapLock.Lock()
	defer dbMapLock.Unlock()

	if dbMap == nil {
		dbMap = make(map[string]*dbRecorder)
	}

	if dbMap[filePath] != nil {
		return dbMap[filePath], nil
	} else {
		dbDir := filepath.Dir(filePath)
		if err := os.MkdirAll(dbDir, os.ModePerm); err != nil {
			return nil, data.NewEmptyError().AppendDesc("open db: make file").AppendError(err)
		}

		handler, err := openDB(filePath)
		if err != nil {
			return nil, data.NewEmptyError().AppendDesc("open db: open").AppendError(err)
		}
		dbMap[filePath] = handler
		return handler, nil
	}
}

func openDB(filePath string) (*dbRecorder, *data.CodeError) {
	db, err := leveldb.OpenFile(filePath, nil)
	if err != nil {
		return nil, data.NewEmptyError().AppendError(err)
	}

	return &dbRecorder{
		filePath: filePath,
		db:       db,
	}, nil
}

func (db *dbRecorder) Get(key string) (string, *data.CodeError) {
	if db.db == nil {
		return "", data.NewEmptyError().AppendDescF("db get key:%s error:no db exist", key)
	}
	ret, err := db.db.Get([]byte(key), nil)
	if err != nil {
		return "", data.NewEmptyError().AppendError(err)
	}
	return string(ret), nil
}

func (db *dbRecorder) Put(key, value string) *data.CodeError {
	if db.db == nil {
		return data.NewEmptyError().AppendDescF("db put key:%s for value:%s error:no db exist", key, value)
	}
	err := db.db.Put([]byte(key), []byte(value), &opt.WriteOptions{
		Sync: false,
	})
	if err != nil {
		return data.NewEmptyError().AppendError(err)
	}
	return nil
}

func (db *dbRecorder) Delete(key string) *data.CodeError {
	if db.db == nil {
		return data.NewEmptyError().AppendDescF("db delete key:%s error:no db exist", key)
	}
	err := db.db.Delete([]byte(key), nil)
	if err != nil {
		return data.NewEmptyError().AppendError(err)
	}
	return nil
}
