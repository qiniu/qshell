package download

import (
	"fmt"
	"github.com/qiniu/qshell/v2/iqshell/common/data"
	"github.com/qiniu/qshell/v2/iqshell/common/db"
	"github.com/qiniu/qshell/v2/iqshell/common/log"
	"strconv"
	"strings"
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

func (d *dbHandler) init() (err *data.CodeError) {
	if len(d.DBFilePath) == 0 {
		return nil
	}

	d.dbHandler, err = db.OpenDB(d.DBFilePath)
	if err != nil {
		return data.NewEmptyError().AppendDesc("download init error:" + err.Error())
	}
	return
}

func (d *dbHandler) checkInfoOfDB() *data.CodeError {
	if d.dbHandler == nil {
		log.Debug("check db: no db handler set")
		return nil
	}

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
		return data.NewEmptyError().AppendDesc("get file modify time error from db, error:" + err.Error())
	}
	// db 数据：文件大小
	fileSize, err := strconv.ParseInt(items[1], 10, 64)
	if err != nil {
		return data.NewEmptyError().AppendDesc("get file size error from db, error:" + err.Error())
	}
	// db 数据：文件 hash
	fileHash := items[2]

	// 验证文件大小
	if d.FileSize > 0 && fileSize != d.FileSize {
		return data.NewEmptyError().AppendDescF("local file size doesn't match server, fileSize: %d|%d", fileSize, d.FileSize)
	}

	// 验证数据库 hash
	if len(d.FileHash) > 0 && len(fileHash) > 0 && fileHash != d.FileHash {
		return data.NewEmptyError().AppendDescF("local file hash doesn't match server, fileHash: %s|%s", fileHash, d.FileHash)
	}

	// 验证修改时间
	if fileServerUpdateTime > 0 && fileServerUpdateTime != d.FileServerUpdateTime {
		return data.NewEmptyError().AppendDescF("local file update time doesn't match server, updateTime: %d|%s", fileServerUpdateTime, d.FileServerUpdateTime)
	}

	return nil
}

func (d *dbHandler) saveInfoToDB() (err *data.CodeError) {
	if d.dbHandler == nil {
		return nil
	}

	value := fmt.Sprintf("%d|%d|%s", d.FileServerUpdateTime, d.FileSize, d.FileHash)
	return d.dbHandler.Put(d.FilePath, value)
}
