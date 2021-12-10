package operations

import (
	"github.com/qiniu/qshell/v2/iqshell/common/data"
	"github.com/qiniu/qshell/v2/iqshell/common/log"
	"github.com/qiniu/qshell/v2/iqshell/storage/object/rs"
	"os"
)

type MirrorUpdateInfo rs.PrefetchApiInfo

func MirrorUpdate(info MirrorUpdateInfo) {
	err := rs.Prefetch(rs.PrefetchApiInfo(info))
	if err != nil {
		log.ErrorF("Prefetch error: %v", err)
		os.Exit(data.STATUS_ERROR)
	}
}
