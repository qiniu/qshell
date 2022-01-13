package download

import (
	"errors"
	"fmt"
	"github.com/qiniu/qshell/v2/iqshell/common/db"
	"github.com/qiniu/qshell/v2/iqshell/common/log"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"sync"
)

//db 金检查是否已下载过，且下载后保存的数据符合预期；

const infoSegment = "|"

type dbHandler struct {
	DBFilePath           string // db 文件路径 【必填】
	FilePath             string // 被检测的文件 【必填】
	FileHash             string // 文件 hash，有值则会检测 hash【必填】
	FileSize             int64  // 文件大小，有值则会检测文件大小【必填】
	FileServerUpdateTime int64  // 服务端文件修改时间 【必填】
	dbHandler            *db.DB
}

func (d *dbHandler) init() (err error) {
	if len(d.DBFilePath) == 0 {
		return nil
	}

	d.dbHandler, err = openDB(d.DBFilePath)
	if err != nil {
		return errors.New("download init error:" + err.Error())
	}
	return
}

func (d *dbHandler) checkInfoOfDB() error {
	// 数据库中存在也验证数据库信息，数据库不存在则仅验证本地文件信息
	value, _ := d.dbHandler.Get(d.FilePath)
	items := strings.Split(value, infoSegment)
	if len(items) == 0 || len(items) < 3 {
		log.Warning("get invalid file info from db:" + value)
		return nil
	}

	// db 数据：服务端文件修改时间
	fileServerUpdateTime, err := strconv.ParseInt(items[0], 10, 64)
	if err != nil {
		return errors.New("get file modify time error from db, error:" + err.Error())
	}
	// db 数据：文件大小
	fileSize, err := strconv.ParseInt(items[1], 10, 64)
	if err != nil {
		return errors.New("get file size error from db, error:" + err.Error())
	}
	// db 数据：文件 hash
	fileHash := items[2]

	// 验证文件大小
	if d.FileSize > 0 && fileSize != d.FileSize {
		return fmt.Errorf("local file size doesn't match server, fileSize: %d|%d", fileSize, d.FileSize)
	}

	// 验证数据库 hash
	if len(d.FileHash) > 0 && len(fileHash) > 0 && fileHash != d.FileHash {
		return fmt.Errorf("local file hash doesn't match server, fileHash: %s|%s", fileHash, d.FileHash)
	}

	// 验证修改时间
	if fileServerUpdateTime > 0 && fileServerUpdateTime != d.FileServerUpdateTime {
		return fmt.Errorf("local file update time doesn't match server, updateTime: %d|%s", fileServerUpdateTime, d.FileServerUpdateTime)
	}

	return nil
}

func (d *dbHandler) saveInfoToDB() (err error) {
	value := fmt.Sprintf("%d|%d|%s", d.FileServerUpdateTime, d.FileSize, d.FileHash)
	return d.dbHandler.Put(d.FilePath, value)
}

var dbMap map[string]*db.DB
var dbMapLock sync.Mutex

func openDB(filePath string) (*db.DB, error) {
	dbMapLock.Lock()
	defer dbMapLock.Unlock()

	if dbMap == nil {
		dbMap = make(map[string]*db.DB)
	}

	if dbMap[filePath] != nil {
		return dbMap[filePath], nil
	} else {
		dbDir := filepath.Dir(filePath)
		err := os.MkdirAll(dbDir, 0775)
		if err != nil {
			return nil, errors.New("download db make file error:" + err.Error())
		}

		handler, err := db.OpenDB(filePath)
		if err != nil {
			return nil, errors.New("download db open error:" + err.Error())
		}
		dbMap[filePath] = handler
		return handler, nil
	}
}
