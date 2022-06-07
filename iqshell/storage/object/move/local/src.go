package local

import (
	"fmt"
	"github.com/qiniu/qshell/v2/iqshell/common/alert"
	"github.com/qiniu/qshell/v2/iqshell/common/data"
	"github.com/qiniu/qshell/v2/iqshell/storage/object/move"
	"os"
	"path/filepath"
)

type SrcInfo struct {
	FilePath string
}

func (i *SrcInfo) Check() *data.CodeError {
	if len(i.FilePath) == 0 {
		return alert.CannotEmptyError("FilePath", "")
	}

	dir := filepath.Dir(i.FilePath)
	if len(dir) == 0 {
		return data.NewEmptyError().AppendDescF("dir is invalid:%s", i.FilePath)
	}
	return nil
}

func NewSrc(info SrcInfo) (move.Src, *data.CodeError) {
	if err := info.Check(); err != nil {
		return nil, err
	}

	return &localFile{
		filePath: info.FilePath,
	}, nil
}

func (l *localFile) PrepareToRead(offset int64) error {
	handle, err := os.OpenFile(l.filePath, os.O_RDONLY, 0655)
	if err != nil{
		return fmt.Errorf("open file error:%v", err)
	}

	if offset > 0 {
		if _, err = handle.Seek(offset, 0); err != nil {
			return fmt.Errorf("seek error:%v", err)
		}
	}

	l.fileHandler = handle
	return nil
}

func (l *localFile) Read(p []byte) (n int, err error) {
	return l.fileHandler.Read(p)
}

func (l *localFile) CompleteRead() error {
	return nil
}
