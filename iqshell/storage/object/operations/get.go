package operations

import (
	"github.com/qiniu/qshell/v2/iqshell/common/data"
	"github.com/qiniu/qshell/v2/iqshell/common/log"
	"github.com/qiniu/qshell/v2/iqshell/storage/object"
	"os"
)

type GetInfo object.GetApiInfo

func GetObject(info GetInfo) {
	err := object.GetObject(object.GetApiInfo(info))
	if err != nil {
		log.ErrorF("Get error: %v\n", err)
		os.Exit(data.StatusError)
	}
}
