package operations

import (
	"github.com/qiniu/qshell/v2/iqshell/common/data"
	"github.com/qiniu/qshell/v2/iqshell/common/log"
	"github.com/qiniu/qshell/v2/iqshell/common/utils"
	"github.com/qiniu/qshell/v2/iqshell/storage/object/rs"
	"os"
)

type FetchInfo rs.FetchApiInfo

func Fetch(info FetchInfo) {
	result, err := rs.Fetch(rs.FetchApiInfo(info))
	if err != nil {
		log.ErrorF("Fetch error: %v", err)
		os.Exit(data.STATUS_ERROR)
	} else {
		log.AlertF("Key:%s", result.Key)
		log.AlertF("Hash:%s", result.Hash)
		log.AlertF("Fsize: %d (%s)", result.Fsize, utils.FormatFileSize(result.Fsize))
		log.AlertF("Mime:%s", result.MimeType)
	}
}
