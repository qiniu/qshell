package download

import (
	"errors"
	"github.com/qiniu/qshell/v2/iqshell/common/log"
	"github.com/qiniu/qshell/v2/iqshell/common/utils"
	"github.com/qiniu/qshell/v2/iqshell/storage/bucket"
	"os"
)

type Checker struct {
	File                string // 被检测的文件 【必填】
	Bucket              string // 文件所在 bucket 用于检查 hash【选填】
	Key                 string // 文件被保存的 key  用于检查 hash【选填】
	FileHash            string // 文件 hash，有值则会检测 hash【选填】
	FileSize            int64  // 文件大小，有值则会检测文件大小【选填】
	RemoveFileWhenError bool   // 当遇到错误时是否该移除文件【选填】
}

func (c *Checker) Check() (err error) {
	defer func() {
		if err != nil && c.RemoveFileWhenError {
			e := os.Remove(c.File)
			if e != nil && !os.IsNotExist(e) {
				log.WarningF("download check: remove file error:%v", e)
			}
		}
	}()

	err = c.CheckFileSize()
	if err != nil {
		return err
	}

	err = c.CheckFileHash()
	return
}

func (c *Checker) CheckFileSize() error {
	if c.FileSize <= 0 {
		return nil
	}

	tempFileStatus, err := os.Stat(c.File)
	if err != nil {
		return err
	}

	if tempFileStatus == nil {
		return errors.New("download check: can't get file status:" + c.File)
	}

	if c.FileSize != tempFileStatus.Size() {
		return errors.New("download check: download file size is unexpected:" + c.File)
	}

	return nil
}

func (c *Checker) CheckFileHash() error {
	if len(c.FileHash) == 0 {
		return nil
	}

	hashFile, err := os.Open(c.File)
	if err != nil {
		return errors.New("download check: check hash get temp file error:" + err.Error())
	}

	var hash string
	if utils.IsSignByEtagV2(c.FileHash) {
		log.Debug("download check: get etag by v2 for key:" + c.Key)
		if len(c.Bucket) == 0 || len(c.Key) == 0 {
			return errors.New("download check: etag v2 check should provide bucket and key")
		}

		bucketManager, err := bucket.GetBucketManager()
		if err != nil {
			return errors.New("download check: etag v2 get bucket manager error:" + err.Error())
		}

		stat, err := bucketManager.Stat(c.Bucket, c.Key)
		if err != nil {
			return errors.New("download check: etag v2 get file status error:" + err.Error())
		}

		hash, err = utils.EtagV2(hashFile, stat.Parts)
	} else {
		log.Debug("download check: get etag by v1 for key:" + c.Key)
		hash, err = utils.EtagV1(hashFile)
	}

	if err != nil {
		return errors.New("download check: get file etag error:" + err.Error())
	}

	if hash != c.FileHash {
		return errors.New("download: file hash doesn't match for key:" + c.Key + "download file hash:" + hash + " excepted:" + c.FileHash)
	}

	return nil
}
