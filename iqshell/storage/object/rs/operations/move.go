package operations

import (
	"github.com/qiniu/qshell/v2/iqshell/common/log"
	"github.com/qiniu/qshell/v2/iqshell/storage/object/rs"
)

type MoveInfo rs.MoveApiInfo

func Move(info MoveInfo) {
	result, err := rs.Move(rs.MoveApiInfo(info))
	if err != nil {
		log.ErrorF("Move error:%v", err)
		return
	}

	if len(result.Error) != 0 {
		log.ErrorF("Move error:%v", result.Error)
		return
	}
}
