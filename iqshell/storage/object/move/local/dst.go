package local

import (
	"fmt"
	"github.com/qiniu/qshell/v2/iqshell/common/alert"
	"github.com/qiniu/qshell/v2/iqshell/common/data"
	"github.com/qiniu/qshell/v2/iqshell/storage/object/move"
	"os"
	"path/filepath"
)

type DstInfo struct {
	FilePath string
}

func (i *DstInfo) Check() *data.CodeError {
	if len(i.FilePath) == 0 {
		return alert.CannotEmptyError("FilePath", "")
	}

	dir := filepath.Dir(i.FilePath)
	if len(dir) == 0 {
		return data.NewEmptyError().AppendDescF("dir is invalid:%s", i.FilePath)
	}
	return nil
}

func NewDst(info DstInfo) (move.Dst, *data.CodeError) {
	if err := info.Check(); err != nil {
		return nil, err
	}

	return &localFile{
		filePath: info.FilePath,
	}, nil
}

func (l *localFile) PrepareToWrite() (offset int64, err error) {
	if stat, sErr := os.Stat(l.filePath); sErr == nil {
		// 文件已存在
		if handle, e := os.OpenFile(l.filePath, os.O_APPEND|os.O_WRONLY, 0655); e != nil {
			return 0, fmt.Errorf("open file error:%v", e)
		} else {
			l.fileHandler = handle
			return stat.Size(), nil
		}
	} else if os.IsNotExist(err) {
		if handle, e := os.Create(l.filePath); e != nil {
			return 0, fmt.Errorf("create file error:%v", e)
		} else {
			l.fileHandler = handle
			return 0, nil
		}
	} else {
		return 0, fmt.Errorf("get file stat error:%v", err)
	}
}

func (l *localFile) Write(p []byte) (n int, err error) {
	return l.fileHandler.Write(p)
}

func (l *localFile) CompleteWrite() error {
	return nil
}
