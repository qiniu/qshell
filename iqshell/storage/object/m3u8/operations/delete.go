package operations

import (
	"os"

	"github.com/qiniu/qshell/v2/iqshell"
	"github.com/qiniu/qshell/v2/iqshell/common/alert"
	"github.com/qiniu/qshell/v2/iqshell/common/data"
	"github.com/qiniu/qshell/v2/iqshell/common/log"
	"github.com/qiniu/qshell/v2/iqshell/storage/object/m3u8"
)

type DeleteInfo m3u8.DeleteApiInfo

func (info *DeleteInfo) Check() *data.CodeError {
	if len(info.Bucket) == 0 {
		return alert.CannotEmptyError("Bucket", "")
	}
	if len(info.Key) == 0 {
		return alert.CannotEmptyError("Key", "")
	}
	return nil
}

func Delete(cfg *iqshell.Config, info DeleteInfo) {
	if shouldContinue := iqshell.CheckAndLoad(cfg, iqshell.CheckAndLoadInfo{
		Checker: &info,
	}); !shouldContinue {
		return
	}

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
