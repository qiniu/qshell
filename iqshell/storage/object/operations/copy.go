package operations

import (
	"github.com/qiniu/qshell/v2/iqshell/common/log"
	"github.com/qiniu/qshell/v2/iqshell/storage/object/rs"
)

type CopyInfo rs.CopyApiInfo

func Copy(info CopyInfo) {
	result, err := rs.Copy(rs.CopyApiInfo(info))
	if err != nil {
		log.ErrorF("Copy error:%v", err)
		return
	}

	if len(result.Error) != 0 {
		log.ErrorF("Copy error:%v", result.Error)
		return
	}
}
