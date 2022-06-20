package upload

import (
	"github.com/qiniu/qshell/v2/iqshell/common/data"
	"github.com/qiniu/qshell/v2/iqshell/common/log"
	"github.com/qiniu/qshell/v2/iqshell/common/utils"
	"github.com/qiniu/qshell/v2/iqshell/storage/object"
	"github.com/qiniu/qshell/v2/iqshell/storage/object/batch"
	"os"
)

type serverChecker struct {
	Bucket     string
	Key        string
	FilePath   string // 本地文件存储位置
	CheckExist bool   // 检查服务端是否已存在
	CheckHash  bool   // 是否检查 hash, 检查是会对比服务端文件 hash
	CheckSize  bool   // 是否检查文件大小，检查是会对比服务端文件大小
	FileSize   int64  // 文件大小，核对文件大小时使用
}

func (c *serverChecker) check() (exist, match bool, err *data.CodeError) {
	fileServerStatus, err := object.Status(object.StatusApiInfo{
		Bucket: c.Bucket,
		Key:    c.Key,
	})
	if err != nil {
		err = data.NewEmptyError().AppendDescF("get file status, %s", err)
		return false, false, err
	}

	checkHash := c.CheckHash
	if checkHash && utils.IsNetworkSource(c.FilePath) {
		checkHash = false
		log.WarningF("network resource doesn't support check hash: %s", c.FilePath)
	}

	if checkHash {
		return c.checkHash(fileServerStatus.OperationResult)
	} else if c.CheckSize {
		return c.checkServerSize(fileServerStatus.OperationResult)
	} else {
		//return true, true, nil
		return c.checkServerSize(fileServerStatus.OperationResult)
	}
}

func (c *serverChecker) checkHash(fileServerStatus batch.OperationResult) (exist bool, match bool, err *data.CodeError) {
	file, oErr := os.Open(c.FilePath)
	if oErr != nil {
		return true, false, data.NewEmptyError().AppendDescF("check hash: open local file:%s error, %s", c.FilePath, oErr)
	}
	defer func() {
		if e := file.Close(); e != nil {
			log.ErrorF("check hash: close file:%s error:%v", c.FilePath, e)
		}
	}()

	localHash := ""
	if utils.IsSignByEtagV2(fileServerStatus.Hash) {
		localHash, err = utils.EtagV2(file, fileServerStatus.Parts)
		if err != nil {
			return true, false, data.NewEmptyError().AppendDescF("check hash: get etag v2:%s error, %v", c.FilePath, err)
		}
	} else {
		localHash, err = utils.EtagV1(file)
		if err != nil {
			log.ErrorF("====== %v hash:%+v", err, localHash)
			return true, false, data.NewEmptyError().AppendDescF("check hash: get etag v1:%s error, %v", c.FilePath, err)
		}
	}

	if localHash == fileServerStatus.Hash {
		return true, true, nil
	} else {
		log.WarningF("File:%s exist at [%s:%s], but hash not match[%s|%s]",
			c.FilePath, c.Bucket, c.Key, localHash, fileServerStatus.Hash)
		return true, false, nil
	}
}

func (c *serverChecker) checkServerSize(fileServerStatus batch.OperationResult) (bool, bool, *data.CodeError) {
	if c.FileSize == fileServerStatus.FSize {
		return true, true, nil
	} else {
		log.WarningF("File:%s exist at [%s:%s], but size not match[%d|%d]",
			c.FilePath, c.Bucket, c.Key, c.FileSize, fileServerStatus.FSize)
		return true, false, nil
	}
}
