package operations

import (
	"github.com/qiniu/qshell/v2/iqshell/common/data"
	"github.com/qiniu/qshell/v2/iqshell/common/log"
	"github.com/qiniu/qshell/v2/iqshell/storage/object/rs"
	"os"
)

type GetInfo rs.GetApiInfo

func GetObject(info GetInfo) {
	err := rs.GetObject(rs.GetApiInfo(info))
	if err != nil {
		log.ErrorF("Get error: %v\n", err)
		os.Exit(data.STATUS_ERROR)
	}
}
