package upload

import (
	"errors"
	"fmt"
	"github.com/qiniu/qshell/v2/iqshell/common/db"
	"github.com/qiniu/qshell/v2/iqshell/common/log"
	"strconv"
	"strings"
)

const infoSegment = "|"

type dbHandler struct {
	DBFilePath     string // db 文件路径 【必填】
	FilePath       string // 被检测的文件 【必填】
	FileUpdateTime int64  // 本地文件修改的时间 【必填】
	dbHandler      *db.DB
}

func (d *dbHandler) init() (err error) {
	if len(d.DBFilePath) == 0 {
		return nil
	}

	d.dbHandler, err = db.OpenDB(d.DBFilePath)
	if err != nil {
		return errors.New("download init error:" + err.Error())
	}
	return
}

// 当数据库中不再相应文件信息 或 文件信息不匹配 则返回 error, (exist, match, error)
func (d *dbHandler) checkInfoOfDB() (bool, bool, error) {
	if d.dbHandler == nil {
		return false, false, errors.New("upload db: no set upload db path")
	}

	// 数据库中存在也验证数据库信息，数据库不存在则仅验证本地文件信息
	value, _ := d.dbHandler.Get(d.FilePath)
	items := strings.Split(value, infoSegment)
	if len(items) == 0 || len(items[0]) == 0 {
		return false, false, errors.New("upload db: get invalid file info from db:" + value)
	}

	// db 数据：服务端文件修改时间
	fileUpdateTime, err := strconv.ParseInt(items[0], 10, 64)
	if err != nil {
		return true, false, errors.New("upload db: get file modify time error from db, error:" + err.Error())
	}

	// 验证修改时间
	if fileUpdateTime != d.FileUpdateTime {
		log.WarningF("upload db: local file has update, updateTime: %d|%d", d.FileUpdateTime, fileUpdateTime)
		return true, false, nil
	} else {
		return true, true, nil
	}
}

func (d *dbHandler) saveInfoToDB() (err error) {
	if d.dbHandler == nil {
		log.Debug("upload save status to db error:no db handler")
		return
	}

	value := fmt.Sprintf("%d", d.FileUpdateTime)
	return d.dbHandler.Put(d.FilePath, value)
}
