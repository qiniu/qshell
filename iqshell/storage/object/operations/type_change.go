package operations

import (
	"errors"
	"github.com/qiniu/qshell/v2/iqshell/common/alert"
	"github.com/qiniu/qshell/v2/iqshell/common/log"
	"github.com/qiniu/qshell/v2/iqshell/storage/object/rs"
	"strconv"
)

type ChangeTypeInfo struct {
	Bucket string
	Key    string
	Type   string
}

func (c ChangeTypeInfo)getTypeOfInt() (int, error) {
	if len(c.Type) == 0 {
		return -1, errors.New(alert.CannotEmpty("type", ""))
	}

	ret, err := strconv.Atoi(c.Type)
	if err != nil {
		return -1, errors.New("parse type error:" + err.Error())
	}

	if ret < 0 || ret > 1 {
		return -1, errors.New("type must be 0 or 1")
	}
	return ret, nil
}

func ChangeType(info ChangeTypeInfo) {
	t, err := info.getTypeOfInt()
	if err != nil {
		log.ErrorF("Change Type error:%v", err)
		return
	}

	result, err := rs.ChangeType(rs.ChangeTypeApiInfo{
		Bucket: info.Bucket,
		Key:    info.Key,
		Type:   t,
	})

	if err != nil {
		log.ErrorF("Change Type error:%v", err)
		return
	}

	if len(result.Error) != 0 {
		log.ErrorF("Change Type:%v", result.Error)
		return
	}
}