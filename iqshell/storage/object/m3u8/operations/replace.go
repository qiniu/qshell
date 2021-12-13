package operations

import (
	"github.com/qiniu/qshell/v2/iqshell/common/data"
	"github.com/qiniu/qshell/v2/iqshell/common/log"
	"github.com/qiniu/qshell/v2/iqshell/storage/object/m3u8"
	"os"
)

type ReplaceDomainInfo m3u8.ReplaceDomainApiInfo

func ReplaceDomain(info ReplaceDomainInfo) {
	err := m3u8.ReplaceDomain(m3u8.ReplaceDomainApiInfo(info))
	if err != nil {
		log.ErrorF("m3u8 replace domain error: %v", err)
		os.Exit(data.STATUS_ERROR)
	}
}
