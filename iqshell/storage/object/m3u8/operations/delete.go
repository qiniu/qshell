package operations

import (
	"github.com/qiniu/qshell/v2/iqshell/common/data"
	"github.com/qiniu/qshell/v2/iqshell/common/log"
	"github.com/qiniu/qshell/v2/iqshell/storage/object/m3u8"
	"os"
)

type DeleteInfo m3u8.DeleteApiInfo

func Delete(info DeleteInfo) {
	err := m3u8.Delete(m3u8.DeleteApiInfo(info))
	if err != nil {
		log.ErrorF("m3u8 delete error: %v", err)
		os.Exit(data.STATUS_ERROR)
	}
}
