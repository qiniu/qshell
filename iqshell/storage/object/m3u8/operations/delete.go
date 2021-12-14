package operations

import (
	"github.com/qiniu/qshell/v2/iqshell/common/data"
	"github.com/qiniu/qshell/v2/iqshell/common/log"
	"github.com/qiniu/qshell/v2/iqshell/storage/object/m3u8"
	"os"
)

type DeleteInfo m3u8.DeleteApiInfo

func Delete(info DeleteInfo) {
	results, err := m3u8.Delete(m3u8.DeleteApiInfo(info))
	for _, result := range results {
		if result.Code != 200 || len(result.Error) > 0 {
			log.ErrorF("result error:%s", result.Error)
		}
	}

	if err != nil {
		log.ErrorF("m3u8 delete error: %v", err)
		os.Exit(data.STATUS_ERROR)
	}
}
