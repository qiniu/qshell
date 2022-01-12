package download

import (
	"errors"
	"fmt"
	"github.com/qiniu/qshell/v2/iqshell/common/db"
	"os"
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
	FileServerModifyTime int64  // 服务端文件修改时间 【必填】
	dbHandler            *db.DB
}

func (d *dbHandler) init() (err error) {
	if len(d.DBFilePath) == 0 {
		return nil
	}

	d.dbHandler, err = db.OpenDB(d.DBFilePath)
	return
}

func (d *dbHandler) checkInfoOfDB() (exist bool, err error) {
	fileStatus, err := os.Stat(d.FilePath)
	if err != nil {
		if os.IsNotExist(err) {
			// 文件不存在
			exist = false
			_ = d.dbHandler.Delete(d.FilePath)
			return exist, nil
		} else {
			// 文件存在，但访问错误
			exist = true
			return exist, errors.New("check db: get local file status error, error:" + err.Error())
		}
	}

	// 数据库中存在也验证数据库信息，不存在仅验证本地文件信息
	value, _ := d.dbHandler.Get(d.FilePath)
	items := strings.Split(value, infoSegment)
	// db 数据：服务端文件修改时间
	fileServerModifyTime := int64(0)
	if len(items) > 0 && len(items[0]) > 0 {
		fileServerModifyTime, err = strconv.ParseInt(items[0], 10, 64)
		if err != nil {
			return exist, errors.New("get file modify time error from db, error:" + err.Error())
		}
	}
	// db 数据：文件大小
	fileSize := int64(0)
	if len(items) > 1 {
		fileSize, err = strconv.ParseInt(items[1], 10, 64)
		if err != nil {
			return exist, errors.New("get file size error from db, error:" + err.Error())
		}
	}
	// db 数据：文件 hash
	fileHash := ""
	if len(items) > 2 {
		fileHash = items[2]
	}
	// db 数据：文件修改时间
	fileModifyTime := int64(0)
	if len(items) > 3 {
		fileModifyTime, err = strconv.ParseInt(items[3], 10, 64)
		if err != nil {
			return exist, errors.New("get file modify time error from db, error:" + err.Error())
		}
	}

	// 验证本地文件信息是否和 check 数据一致
	if fileStatus.Size() != d.FileSize {
		return exist, fmt.Errorf("local file info doesn't match server, fileSize: %d|%d", fileStatus.Size(), d.FileSize)
	}

	// 验证本地文件是否和数据库存储保存数据一致，数据库中不存在则跳过验证
	if (fileModifyTime > 0 && fileStatus.ModTime().Unix() != fileModifyTime) ||
		(fileSize > 0 && fileStatus.Size() != fileSize) {
		return exist, fmt.Errorf("local file info doesn't match db, modTime: %d|%d  fileSize: %d|%d",
			fileStatus.ModTime().Unix(), fileModifyTime, fileStatus.Size(), fileSize)
	}

	// 验证数据库保存信息是否和 check 数据一致，除了 hash 数据库中不存在则跳过验证
	if (len(d.FileHash) > 0 && fileHash != d.FileHash) || /* 文件 hash，严格验证，只要存在就会验证，本地数据库中没有则直接报错 */
		(fileSize > 0 && fileSize != d.FileSize) || /* 文件大小 */
		(fileServerModifyTime > 0 && fileServerModifyTime != d.FileServerModifyTime) /* 服务端文件修改时间 */ {
		return exist, fmt.Errorf("local file info doesn't match db, modTime: %d|%d  fileSize: %d|%d",
			fileStatus.ModTime().Unix(), fileModifyTime, fileStatus.Size(), d.FileSize)
	}

	return exist, nil
}

func (d *dbHandler) saveInfoToDB() (err error) {
	fileStatus, err := os.Stat(d.FilePath)
	if err != nil {
		return errors.New("save db: get local file status error, error:" + err.Error())
	}

	value := fmt.Sprintf("%d|%d|%s|%d", d.FileServerModifyTime, d.FileSize, d.FileHash, fileStatus.ModTime().Unix())
	return d.dbHandler.Put(d.FilePath, value)
}
