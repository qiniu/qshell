package operations

import (
	"github.com/qiniu/qshell/v2/iqshell/common/log"
	"github.com/qiniu/qshell/v2/iqshell/storage/object/rs"
)

func Copy(info rs.CopyApiInfo) {
	result, err := rs.Copy(info)
	if err != nil {
		log.ErrorF("Copy error:%v", err)
		return
	}

	if len(result.Error) != 0 {
		log.ErrorF("Copy error:%v", result.Error)
		return
	}
}