package operations

import (
	"github.com/qiniu/qshell/v2/iqshell/common/alert"
	"github.com/qiniu/qshell/v2/iqshell/common/data"
	"github.com/qiniu/qshell/v2/iqshell/common/log"
	"github.com/qiniu/qshell/v2/iqshell/storage/object/m3u8"
	"os"
)

type DeleteInfo m3u8.DeleteApiInfo

func (info *DeleteInfo) Check() error {
	if len(info.Bucket) == 0 {
		return alert.CannotEmptyError("Bucket", "")
	}
	if len(info.Key) == 0 {
		return alert.CannotEmptyError("Key", "")
	}
	return nil
}

func Delete(info DeleteInfo) {
	results, err := m3u8.Delete(m3u8.DeleteApiInfo(info))
	for _, result := range results {
		if result.Code != 200 || len(result.Error) > 0 {
			log.ErrorF("result error:%s", result.Error)
		}
	}

	if err != nil {
		log.ErrorF("m3u8 delete error: %v", err)
		os.Exit(data.StatusError)
	}
}
