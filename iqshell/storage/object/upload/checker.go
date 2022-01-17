package upload

import (
	"fmt"
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
}

func (c *serverChecker) isNeedCheck() bool {
	return c.CheckHash || c.CheckSize
}

func (c *serverChecker) check() (exist, match bool, err error) {
	fileServerStatus, err := object.Status(object.StatusApiInfo{
		Bucket: c.Bucket,
		Key:    c.Key,
	})
	if err != nil {
		err = fmt.Errorf("upoad check hash: get file [%s:%s] stat error, %s", c.Bucket, c.Key, err)
		return
	}

	if c.CheckHash {
		return c.checkHash(fileServerStatus)
	} else if c.CheckSize {
		return c.checkServerSize(fileServerStatus)
	}

	return false, false, nil
}

func (c *serverChecker) checkHash(fileServerStatus batch.OperationResult) (bool, bool, error) {
	file, err := os.Open(c.FilePath)
	if err != nil {
		return false, false, fmt.Errorf("upoad check hash: open local file:%s error, %s", c.FilePath, err)
	}
	defer func() {
		if e := file.Close(); e != nil {
			log.ErrorF("upoad check hash: close file:%s error:%v", c.FilePath, e)
		}
	}()

	localHash := ""
	if utils.IsSignByEtagV2(fileServerStatus.Hash) {
		localHash, err = utils.EtagV2(file, fileServerStatus.Parts)
		if err != nil {
			return false, false, fmt.Errorf("upoad check hash: get etag v2:%s error, %s", c.FilePath, err)
		}
	} else {
		localHash, err = utils.EtagV1(file)
		if err != nil {
			return false, false, fmt.Errorf("upoad check hash: get etag v1:%s error, %s", c.FilePath, err)
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

func (c *serverChecker) checkServerSize(fileServerStatus batch.OperationResult) (bool, bool, error) {
	localFileStatus, err := os.Stat(c.FilePath)
	if err != nil {
		return false, false, fmt.Errorf("get file:%s status error:%v", c.FilePath, err)
	}

	if localFileStatus.Size() == fileServerStatus.FSize {
		return true, true, nil
	} else {
		log.WarningF("File:%s exist at [%s:%s], but size not match[%d|%d]",
			c.FilePath, c.Bucket, c.Key, localFileStatus.Size(), fileServerStatus.FSize)
		return true, false, nil
	}
}
