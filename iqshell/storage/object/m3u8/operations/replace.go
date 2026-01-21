package operations

import (
	"os"

	"github.com/qiniu/qshell/v2/iqshell"
	"github.com/qiniu/qshell/v2/iqshell/common/alert"
	"github.com/qiniu/qshell/v2/iqshell/common/data"
	"github.com/qiniu/qshell/v2/iqshell/common/log"
	"github.com/qiniu/qshell/v2/iqshell/storage/object/m3u8"
)

type ReplaceDomainInfo m3u8.ReplaceDomainApiInfo

func (info *ReplaceDomainInfo) Check() *data.CodeError {
	if len(info.Bucket) == 0 {
		return alert.CannotEmptyError("Bucket", "")
	}
	if len(info.Key) == 0 {
		return alert.CannotEmptyError("Key", "")
	}
	return nil
}

func ReplaceDomain(cfg *iqshell.Config, info ReplaceDomainInfo) {
	if shouldContinue := iqshell.CheckAndLoad(cfg, iqshell.CheckAndLoadInfo{
		Checker: &info,
	}); !shouldContinue {
		return
	}

	err := m3u8.ReplaceDomain(m3u8.ReplaceDomainApiInfo(info))
	if err != nil {
		log.ErrorF("m3u8 replace domain error: %v", err)
		os.Exit(data.StatusError)
	}
}
