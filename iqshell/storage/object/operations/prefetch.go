package operations

import (
	"github.com/qiniu/qshell/v2/iqshell/common/data"
	"github.com/qiniu/qshell/v2/iqshell/common/log"
	"github.com/qiniu/qshell/v2/iqshell/storage"
	"os"
)

type MirrorUpdateInfo storage.PrefetchApiInfo

func MirrorUpdate(info MirrorUpdateInfo) {
	err := storage.Prefetch(storage.PrefetchApiInfo(info))
	if err != nil {
		log.ErrorF("mirror update error: %v", err)
		os.Exit(data.STATUS_ERROR)
	}
}
