package operations

import (
	"github.com/qiniu/qshell/v2/iqshell/common/log"
	"github.com/qiniu/qshell/v2/iqshell/storage/object/rs"
)

func Move(info rs.MoveApiInfo) {
	result, err := rs.Move(info)
	if err != nil {
		log.ErrorF("Move error:%v", err)
		return
	}

	if len(result.Error) != 0 {
		log.ErrorF("Move error:%v", result.Error)
		return
	}
}
