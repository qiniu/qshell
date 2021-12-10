package operations

import (
	"github.com/qiniu/qshell/v2/iqshell/common/log"
	"github.com/qiniu/qshell/v2/iqshell/storage/object/rs"
)

func ChangeMime(info rs.ChangeMimeApiInfo) {
	result, err := rs.ChangeMime(info)
	if err != nil {
		log.ErrorF("Change Mime error:%v", err)
		return
	}

	if len(result.Error) != 0 {
		log.ErrorF("Change Mime:%v", result.Error)
		return
	}
}
