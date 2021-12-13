package operations

import (
	"github.com/qiniu/qshell/v2/iqshell/common/log"
	"github.com/qiniu/qshell/v2/iqshell/storage/object/rs"
)

type ForbiddenInfo struct {
	Bucket      string
	Key         string
	UnForbidden bool
}

func (c ForbiddenInfo) getStatus() int {
	// 0:启用  1:禁用
	if c.UnForbidden {
		return 0
	} else {
		return 1
	}
}

func ForbiddenObject(info ForbiddenInfo) {
	result, err := rs.ChangeStatus(rs.ChangeStatusApiInfo{
		Bucket: info.Bucket,
		Key:    info.Key,
		Status: info.getStatus(),
	})

	if err != nil {
		log.ErrorF("change stat error:%v", err)
		return
	}

	if len(result.Error) > 0 {
		log.ErrorF("change stat error:%s", result.Error)
		return
	}
}

func BatchChangeStatus() {

}
