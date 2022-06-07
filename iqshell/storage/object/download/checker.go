package download

import (
	"github.com/qiniu/qshell/v2/iqshell/common/data"
	"github.com/qiniu/qshell/v2/iqshell/common/log"
	"github.com/qiniu/qshell/v2/iqshell/storage/object"
	"os"
)

type FileChecker struct {
	File                string // 被检测的文件 【必填】
	Bucket              string // 文件所在 bucket 用于检查 hash【选填】
	Key                 string // 文件被保存的 key  用于检查 hash【选填】
	FileHash            string // 文件 hash，有值则会检测 hash【选填】
	FileSize            int64  // 文件大小，有值则会检测文件大小【选填】
	RemoveFileWhenError bool   // 当遇到错误时是否该移除文件【选填】
}

func (l *FileChecker) IsFileMatch() (err *data.CodeError) {
	defer func() {
		if err != nil && l.RemoveFileWhenError {
			e := os.Remove(l.File)
			if e != nil && !os.IsNotExist(e) {
				log.WarningF("download file check: remove file error:%v", e)
			}
		}
	}()

	err = l.IsFileSizeMatch()
	if err != nil {
		return err
	}

	err = l.IsFileHashMatch()
	return
}

func (l *FileChecker) IsFileSizeMatch() *data.CodeError {
	if l.FileSize <= 0 {
		log.Debug("download file check size: needn't to check")
		return nil
	}

	tempFileStatus, err := os.Stat(l.File)
	if err != nil {
		return data.ConvertError(err)
	}

	if tempFileStatus == nil {
		return data.NewEmptyError().AppendDesc("download file check: can't get file status:" + l.File)
	}

	if l.FileSize != tempFileStatus.Size() {
		return data.NewEmptyError().AppendDesc("download file check: download file size is unexpected:" + l.File)
	}

	return nil
}

func (l *FileChecker) IsFileHashMatch() *data.CodeError {
	if len(l.FileHash) == 0 {
		log.Debug("download file check hash: needn't to check")
		return nil
	}

	_, err := object.Match(object.MatchApiInfo{
		Bucket:    l.Bucket,
		Key:       l.Key,
		FileHash:  l.FileHash,
		LocalFile: l.File,
	})
	return err
}
