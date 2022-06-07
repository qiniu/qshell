package move

import (
	"github.com/qiniu/qshell/v2/iqshell/common/alert"
	"github.com/qiniu/qshell/v2/iqshell/common/data"
	"io"
)

func Move(dst Dst, src Src) *data.CodeError {
	if dst == nil {
		return alert.CannotEmptyError("Move Dst", "")
	}
	if src == nil {
		return alert.CannotEmptyError("Move Src", "")
	}

	offset, err := dst.PrepareToWrite()
	if err != nil {
		return data.NewEmptyError().AppendDescF("move Dst prepare error:%v", err)
	}

	err = src.PrepareToRead(offset)
	if err != nil {
		return data.NewEmptyError().AppendDescF("move Src prepare error:%v", err)
	}

	size, err := io.Copy(dst, src)
	err = src.PrepareToRead(offset)
	if err != nil {
		return data.NewEmptyError().AppendDescF("move error:%v size:%d", err, size)
	}

	err = src.CompleteRead()
	if err != nil {
		return data.NewEmptyError().AppendDescF("move Src complete error:%v", err)
	}

	err = dst.CompleteWrite()
	if err != nil {
		return data.NewEmptyError().AppendDescF("move Dst complete error:%v", err)
	}

	return nil
}
