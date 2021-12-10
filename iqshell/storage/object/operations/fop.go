package operations

import (
	"github.com/qiniu/qshell/v2/iqshell/common/log"
	"github.com/qiniu/qshell/v2/iqshell/storage/object/rs"
)

type PreFopStatusInfo struct {
	Id string
}

func PreFopStatus(info PreFopStatusInfo) {
	ret, err := rs.PreFopStatus(info.Id)
	if err != nil {
		log.ErrorF("pre fog status error:%v", err)
		return
	}

	log.Alert(ret.String())
}

type PreFopInfo rs.PreFopApiInfo

func PreFop(info PreFopInfo) {
	persistentId, err := rs.PreFop(rs.PreFopApiInfo(info))
	if err != nil {
		log.ErrorF("pre fog error:%v", err)
		return
	}
	log.Alert(persistentId)
}
